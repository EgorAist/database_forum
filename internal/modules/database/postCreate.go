package database

import (
	"fmt"
	"strings"
	"time"

	"github.com/EgorAist/database_foruminternal/models"
	"gopkg.in/jackc/pgx.v2"
)

func PostsCreate(db *pgx.ConnPool, slugOrIDThread string, posts models.Posts) (models.Posts, error) {
	isExist := false
	var err error
	threadID := 0
	forumSlug := ""
	if id, isID := isID(slugOrIDThread); isID {
		threadID = id
		forumSlug, isExist, err = ifThreadExistAndGetFodumSlugByID(db, threadID)
	} else {
		forumSlug, threadID, isExist, err = ifThreadExistAndGetIDForumSlugBySlug(db, slugOrIDThread)
	}
	if !isExist {
		return nil, ErrThreadNotFound
	}

	if err != nil {
		return nil, err
	}

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec("SET LOCAL synchronous_commit TO OFF")
	if err != nil {
		if txErr := tx.Rollback(); txErr != nil {
			return nil, txErr
		}
		return nil, err
	}

	resultPosts, err := insertPostsTx(tx, threadID, posts, forumSlug)
	if err != nil {
		if txErr := tx.Rollback(); txErr != nil {
			return nil, txErr
		}

		if pqError, ok := err.(pgx.PgError); ok {
			switch pqError.Code {
			case pgErrForeignKeyViolation:
				if pqError.ConstraintName == "post_parent_id_fkey" {
					return nil, ErrPostConflict
				}
				if pqError.ConstraintName == "post_author_fkey" {
					return nil, ErrUserNotFound
				}
			}
		}
		return nil, err
	}

	if commitErr := tx.Commit(); commitErr != nil {
		return nil, commitErr
	}

	return resultPosts, nil
}

const (
	insertPostsStart = `
	INSERT INTO post (forum_slug, author, created, message, edited, parent_id, thread_id)
	VALUES 
	`

	insertPostsEnd = `
	RETURNING id, author, created, edited, message, parent_id, thread_id, forum_slug
	`

	insertForumUsersStart = `
	INSERT INTO forum_user (nickname, forum_slug)
	VALUES
	`

	insertForumUsersEnd = `
	ON CONFLICT ON CONSTRAINT unique_forum_user DO NOTHING
	`

	lockQuery = `SELECT * FROM forum_user FOR UPDATE`
)

func insertPostsTx(tx *pgx.Tx, threadID int, posts models.Posts, forumSlug string) (models.Posts, error) {
	resultPosts := models.Posts{}
	if len(posts) == 0 {
		return resultPosts, nil
	}

	postsArgs := make([]interface{}, 0)
	forumUserArgs := make([]interface{}, 0)
	insertPostsQuery, insertForumUserQuery := formInsertQuery(threadID, posts, &postsArgs, &forumUserArgs, forumSlug)
	rows, queryError := tx.Query(*insertPostsQuery, postsArgs...)
	if queryError != nil {
		return nil, queryError
	}
	for rows.Next() {
		post := &models.Post{}
		err := scanPostRows(rows, post)
		if err != nil {
			rows.Close()
			return nil, err
		}

		resultPosts = append(resultPosts, post)
	}

	if err := rows.Err(); err != nil {
		rows.Close()
		return nil, err
	}

	rows.Close()

	err := forumUpdatePostCountByThreadID(tx, threadID, len(resultPosts))
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(*insertForumUserQuery, forumUserArgs...)
	return resultPosts, err
}

func formInsertQuery(id int, posts models.Posts,
	postsArgs *[]interface{}, forumUserArgs *[]interface{}, forumSlug string) (*string, *string) {

	insertValues := ""
	insertUserValues := ""
	finalInsertValues := strings.Builder{}
	finalUserInsertValues := strings.Builder{}

	for idx, post := range posts {
		insertValues = formInsertValuesID(post.Author, post.Created, post.Message, id,
			post.IsEdited, post.Parent, post.Thread, idx*5+1, postsArgs, forumSlug)

		insertUserValues = formInsertUserValues(post.Author, idx*2+1, forumUserArgs, forumSlug)

		if idx != 0 {
			finalInsertValues.WriteString(",")
			finalUserInsertValues.WriteString(",")
		}
		finalInsertValues.WriteString(insertValues)
		finalUserInsertValues.WriteString(insertUserValues)
	}

	insertPostsQuery := strings.Builder{}
	insertPostsQuery.WriteString(insertPostsStart)
	insertPostsQuery.WriteString(finalInsertValues.String())
	insertPostsQuery.WriteString(insertPostsEnd)

	insertUsersQuery := strings.Builder{}
	insertUsersQuery.WriteString(insertForumUsersStart)
	insertUsersQuery.WriteString(finalUserInsertValues.String())
	insertUsersQuery.WriteString(insertForumUsersEnd)

	resPosts, resUser := insertPostsQuery.String(), insertUsersQuery.String()
	return &resPosts, &resUser
}

func formInsertUserValues(author string, placeholder int, args *[]interface{}, forumSlug string) string {
	values := fmt.Sprintf("($%v, $%v)", placeholder, placeholder+1)
	*args = append(*args, author)
	*args = append(*args, forumSlug)
	return values
}

const insertWithCheckParentID = `(
	SELECT (
		CASE WHEN 
		EXISTS(SELECT 1 from post p where p.id=%v and p.thread_id=%v)
		THEN %v ELSE -1 END)
	)`

func formInsertValuesID(author string, created time.Time, message string, ID int, isEdited bool, parent int64,
	thread int32, placeholderStart int, valuesArgs *[]interface{}, forumSlug string) string {
	values := "("
	valuesArr := []string{}
	placeholder := placeholderStart

	valuesArr = append(valuesArr, fmt.Sprintf(`'%v'`, forumSlug))

	valuesArr = append(valuesArr, fmt.Sprintf("$%v", placeholder))
	placeholder++

	if author == "" {
		*valuesArgs = append(*valuesArgs, "NULL")
	} else {
		*valuesArgs = append(*valuesArgs, author)
	}

	valuesArr = append(valuesArr, fmt.Sprintf("$%v", placeholder))
	placeholder++

	if created.IsZero() {
		*valuesArgs = append(*valuesArgs, "now()")
	} else {
		*valuesArgs = append(*valuesArgs, created)
	}

	valuesArr = append(valuesArr, fmt.Sprintf("$%v", placeholder))
	placeholder++

	if message == "" {
		*valuesArgs = append(*valuesArgs, "NULL")
	} else {
		*valuesArgs = append(*valuesArgs, message)
	}

	valuesArr = append(valuesArr, fmt.Sprintf("$%v", placeholder))
	placeholder++

	*valuesArgs = append(*valuesArgs, isEdited)

	if parent == 0 {
		valuesArr = append(valuesArr, fmt.Sprint("(NULL)"))
	} else {
		valuesArr = append(valuesArr, fmt.Sprintf(insertWithCheckParentID, parent, ID, parent))
	}

	valuesArr = append(valuesArr, fmt.Sprintf("$%v", placeholder))
	if thread == 0 {
		*valuesArgs = append(*valuesArgs, ID)
	} else {
		*valuesArgs = append(*valuesArgs, thread)
	}

	values += strings.Join(valuesArr, ", ")
	values += ")"

	return values
}
