package database

import (
	"database/sql"

	"github.com/EgorAist/database_forum/internal/models"
	pgx "gopkg.in/jackc/pgx.v2"
)

func SelectPostFull(db *pgx.ConnPool, related []string, pf *models.PostFull) error {
	isIncludeUser, isIncludeForum, isIncludeThread := false, false, false
	for _, rel := range related {
		switch rel {
		case "user":
			pf.Author = &models.User{}
			isIncludeUser = true
		case "forum":
			pf.Forum = &models.Forum{}
			isIncludeForum = true
		case "thread":
			pf.Thread = &models.Thread{}
			isIncludeThread = true
		}
	}

	var err error
	if isIncludeForum && isIncludeUser && isIncludeThread {
		err = selectPostWithForumUserThread(db, pf)
	} else if !isIncludeForum && isIncludeUser && isIncludeThread {
		err = selectPostWithUserThread(db, pf)
	} else if isIncludeForum && !isIncludeUser && isIncludeThread {
		err = selectPostWithForumThread(db, pf)
	} else if isIncludeForum && isIncludeUser && !isIncludeThread {
		err = selectPostWithForumUser(db, pf)
	} else if !isIncludeForum && !isIncludeUser && isIncludeThread {
		err = selectPostWithThread(db, pf)
	} else if !isIncludeForum && isIncludeUser && !isIncludeThread {
		err = selectPostWithUser(db, pf)
	} else if isIncludeForum && !isIncludeUser && !isIncludeThread {
		err = selectPostWithForum(db, pf)
	} else if !isIncludeForum && !isIncludeUser && !isIncludeThread {
		err = selectPost(db, pf.Post)
	}

	if err != nil {
		if err == pgx.ErrNoRows {
			return ErrPostNotFound
		}
		return err
	}
	return nil
}

const (
	selectPostWithForumUserThreadQuery = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, p.forum_slug,
	f.user_nick, f.slug, f.title, f.thread_count, f.post_count,
	t.id, t.slug, t.user_nick, t.created, t.forum_slug, t.title, t.message, t.votes,
	u.nickname, u.fullname, u.about, u.email
	FROM post p 
	JOIN thread t ON p.thread_id = t.id
	JOIN "user" u ON p.author = u.nickname
	JOIN forum f ON p.forum_slug = f.slug
	WHERE p.id = $1`
)

func selectPostWithForumUserThread(db *pgx.ConnPool, pf *models.PostFull) error {
	parent := sql.NullInt64{}
	slugThread := sql.NullString{}
	err := db.QueryRow(selectPostWithForumUserThreadQuery, pf.Post.ID).Scan(
		&pf.Post.ID,
		&pf.Post.Author,
		&pf.Post.Created,
		&pf.Post.IsEdited,
		&pf.Post.Message,
		&parent,
		&pf.Post.Thread,
		&pf.Post.Forum,
		&pf.Forum.User,
		&pf.Forum.Slug,
		&pf.Forum.Title,
		&pf.Forum.Threads,
		&pf.Forum.Posts,
		&pf.Thread.ID,
		&slugThread,
		&pf.Thread.Author,
		&pf.Thread.Created,
		&pf.Thread.Forum,
		&pf.Thread.Title,
		&pf.Thread.Message,
		&pf.Thread.Votes,
		&pf.Author.Nickname,
		&pf.Author.Fullname,
		&pf.Author.About,
		&pf.Author.Email,
	)
	if err != nil {
		return err
	}

	if parent.Valid {
		pf.Post.Parent = parent.Int64
	} else {
		pf.Post.Parent = 0
	}
	if slugThread.Valid {
		pf.Thread.Slug = slugThread.String
	} else {
		pf.Thread.Slug = ""
	}
	return nil
}

const (
	selectPostWithUserThreadQuery = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, p.forum_slug,
	t.id, t.slug, t.user_nick, t.created, t.forum_slug, t.title, t.message, t.votes,
	u.nickname, u.fullname, u.about, u.email
	FROM post p 
	JOIN thread t ON p.thread_id = t.id
	JOIN "user" u ON p.author = u.nickname
	WHERE p.id = $1`
)

func selectPostWithUserThread(db *pgx.ConnPool, pf *models.PostFull) error {
	parent := sql.NullInt64{}
	slugThread := sql.NullString{}
	err := db.QueryRow(selectPostWithUserThreadQuery, pf.Post.ID).Scan(
		&pf.Post.ID,
		&pf.Post.Author,
		&pf.Post.Created,
		&pf.Post.IsEdited,
		&pf.Post.Message,
		&parent,
		&pf.Post.Thread,
		&pf.Post.Forum,
		&pf.Thread.ID,
		&slugThread,
		&pf.Thread.Author,
		&pf.Thread.Created,
		&pf.Thread.Forum,
		&pf.Thread.Title,
		&pf.Thread.Message,
		&pf.Thread.Votes,
		&pf.Author.Nickname,
		&pf.Author.Fullname,
		&pf.Author.About,
		&pf.Author.Email,
	)

	if err != nil {
		return err
	}

	if parent.Valid {
		pf.Post.Parent = parent.Int64
	} else {
		pf.Post.Parent = 0
	}

	if slugThread.Valid {
		pf.Thread.Slug = slugThread.String
	} else {
		pf.Thread.Slug = ""
	}
	return nil
}

const (
	selectPostWithForumThreadQuery = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, p.forum_slug,
	f.user_nick, f.slug, f.title, f.thread_count, f.post_count,
	t.id, t.slug, t.user_nick, t.created, t.forum_slug, t.title, t.message, t.votes
	FROM post p 
	JOIN thread t ON p.thread_id = t.id
	JOIN forum f ON p.forum_slug = f.slug
	WHERE p.id = $1`
)

func selectPostWithForumThread(db *pgx.ConnPool, pf *models.PostFull) error {
	parent := sql.NullInt64{}
	slugThread := sql.NullString{}
	err := db.QueryRow(selectPostWithForumThreadQuery, pf.Post.ID).Scan(
		&pf.Post.ID,
		&pf.Post.Author,
		&pf.Post.Created,
		&pf.Post.IsEdited,
		&pf.Post.Message,
		&parent,
		&pf.Post.Thread,
		&pf.Post.Forum,
		&pf.Forum.User,
		&pf.Forum.Slug,
		&pf.Forum.Title,
		&pf.Forum.Threads,
		&pf.Forum.Posts,
		&pf.Thread.ID,
		&slugThread,
		&pf.Thread.Author,
		&pf.Thread.Created,
		&pf.Thread.Forum,
		&pf.Thread.Title,
		&pf.Thread.Message,
		&pf.Thread.Votes,
	)

	if err != nil {
		return err
	}

	if parent.Valid {
		pf.Post.Parent = parent.Int64
	} else {
		pf.Post.Parent = 0
	}

	if slugThread.Valid {
		pf.Thread.Slug = slugThread.String
	} else {
		pf.Thread.Slug = ""
	}
	return nil
}

const (
	selectPostWithForumUserQuery = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, p.forum_slug,
	f.user_nick, f.slug, f.title, f.thread_count, f.post_count,
	u.nickname, u.fullname, u.about, u.email
	FROM post p 
	JOIN "user" u ON p.author = u.nickname
	JOIN forum f ON p.forum_slug = f.slug
	WHERE p.id = $1`
)

func selectPostWithForumUser(db *pgx.ConnPool, pf *models.PostFull) error {
	parent := sql.NullInt64{}
	err := db.QueryRow(selectPostWithForumUserQuery, pf.Post.ID).Scan(
		&pf.Post.ID,
		&pf.Post.Author,
		&pf.Post.Created,
		&pf.Post.IsEdited,
		&pf.Post.Message,
		&parent,
		&pf.Post.Thread,
		&pf.Post.Forum,
		&pf.Forum.User,
		&pf.Forum.Slug,
		&pf.Forum.Title,
		&pf.Forum.Threads,
		&pf.Forum.Posts,
		&pf.Author.Nickname,
		&pf.Author.Fullname,
		&pf.Author.About,
		&pf.Author.Email,
	)

	if err != nil {
		return err
	}

	if parent.Valid {
		pf.Post.Parent = parent.Int64
	} else {
		pf.Post.Parent = 0
	}
	return nil
}

const (
	selectPostWithThreadQuery = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, p.forum_slug,
	t.id, t.slug, t.user_nick, t.created, t.forum_slug, t.title, t.message, t.votes
	FROM post p 
	JOIN thread t ON p.thread_id = t.id
	WHERE p.id = $1`
)

func selectPostWithThread(db *pgx.ConnPool, pf *models.PostFull) error {
	parent := sql.NullInt64{}
	slugThread := sql.NullString{}
	err := db.QueryRow(selectPostWithThreadQuery, pf.Post.ID).Scan(
		&pf.Post.ID,
		&pf.Post.Author,
		&pf.Post.Created,
		&pf.Post.IsEdited,
		&pf.Post.Message,
		&parent,
		&pf.Post.Thread,
		&pf.Post.Forum,
		&pf.Thread.ID,
		&slugThread,
		&pf.Thread.Author,
		&pf.Thread.Created,
		&pf.Thread.Forum,
		&pf.Thread.Title,
		&pf.Thread.Message,
		&pf.Thread.Votes,
	)

	if err != nil {
		return err
	}

	if parent.Valid {
		pf.Post.Parent = parent.Int64
	} else {
		pf.Post.Parent = 0
	}

	if slugThread.Valid {
		pf.Thread.Slug = slugThread.String
	} else {
		pf.Thread.Slug = ""
	}
	return nil
}

const (
	selectPostWithForumQuery = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, p.forum_slug,
	f.user_nick, f.slug, f.title, f.thread_count, f.post_count
	FROM post p 
	JOIN thread t ON p.thread_id = t.id
	JOIN forum f ON p.forum_slug = f.slug
	WHERE p.id = $1`
)

func selectPostWithForum(db *pgx.ConnPool, pf *models.PostFull) error {
	parent := sql.NullInt64{}
	err := db.QueryRow(selectPostWithForumQuery, pf.Post.ID).Scan(
		&pf.Post.ID,
		&pf.Post.Author,
		&pf.Post.Created,
		&pf.Post.IsEdited,
		&pf.Post.Message,
		&parent,
		&pf.Post.Thread,
		&pf.Post.Forum,
		&pf.Forum.User,
		&pf.Forum.Slug,
		&pf.Forum.Title,
		&pf.Forum.Threads,
		&pf.Forum.Posts,
	)

	if err != nil {
		return err
	}

	if parent.Valid {
		pf.Post.Parent = parent.Int64
	} else {
		pf.Post.Parent = 0
	}

	return nil
}

const (
	selectPostWithUserQuery = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, p.forum_slug,
	u.nickname, u.fullname, u.about, u.email
	FROM post p 
	JOIN "user" u ON p.author = u.nickname
	WHERE p.id = $1`
)

func selectPostWithUser(db *pgx.ConnPool, pf *models.PostFull) error {
	parent := sql.NullInt64{}
	err := db.QueryRow(selectPostWithUserQuery, pf.Post.ID).Scan(
		&pf.Post.ID,
		&pf.Post.Author,
		&pf.Post.Created,
		&pf.Post.IsEdited,
		&pf.Post.Message,
		&parent,
		&pf.Post.Thread,
		&pf.Post.Forum,
		&pf.Author.Nickname,
		&pf.Author.Fullname,
		&pf.Author.About,
		&pf.Author.Email,
	)

	if err != nil {
		return err
	}

	if parent.Valid {
		pf.Post.Parent = parent.Int64
	} else {
		pf.Post.Parent = 0
	}
	return nil
}

const (
	selectPostQuery = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, p.forum_slug
	FROM post p 
	WHERE p.id = $1`
)

func selectPost(db *pgx.ConnPool, pf *models.Post) error {
	return scanPost(db.QueryRow(selectPostQuery, pf.ID), pf)
}

func SelectAllPostsByThread(db *pgx.ConnPool, slugOrIDThread string, limit int, desc bool,
	since int, sort string, posts *models.Posts) error {

	isExist := false
	var err error
	threadID := 0
	if id, isID := isID(slugOrIDThread); isID {
		threadID = id
		isExist, err = isThreadExist(db, threadID)
	} else {
		threadID, isExist, err = ifThreadExistGetID(db, slugOrIDThread)
	}

	if !isExist {
		return ErrThreadNotFound
	}

	if err != nil {
		return err
	}

	return selectAllPostsByThreadID(db, threadID, limit, desc, since, sort, posts)
}

const selectPostsFlatLimitByID = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, p.forum_slug
	FROM post p
	WHERE p.thread_id = $1
	ORDER BY p.created, p.id
	LIMIT $2
`

const selectPostsFlatLimitDescByID = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, p.forum_slug
	FROM post p
	WHERE p.thread_id = $1
	ORDER BY p.created DESC, p.id DESC
	LIMIT $2
`

const selectPostsFlatLimitSinceByID = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, p.forum_slug
	FROM post p
	WHERE p.thread_id = $1 and p.id > $2
	ORDER BY p.created, p.id
	LIMIT $3
`
const selectPostsFlatLimitSinceDescByID = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, p.forum_slug
	FROM post p
	WHERE p.thread_id = $1 and p.id < $2
	ORDER BY p.created DESC, p.id DESC
	LIMIT $3
`

const selectPostsTreeLimitByID = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, p.forum_slug
	FROM post p
	WHERE p.thread_id = $1
	ORDER BY p.path
	LIMIT $2
`

const selectPostsTreeLimitDescByID = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, p.forum_slug
	FROM post p
	WHERE p.thread_id = $1
	ORDER BY path DESC
	LIMIT $2
`

const selectPostsTreeLimitSinceByID = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, p.forum_slug
	FROM post p
	WHERE p.thread_id = $1 and (p.path > (SELECT p2.path from post p2 where p2.id = $2))
	ORDER BY p.path
	LIMIT $3
`

const selectPostsTreeLimitSinceDescByID = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, p.forum_slug
	FROM post p
	WHERE p.thread_id = $1 and (p.path < (SELECT p2.path from post p2 where p2.id = $2))
	ORDER BY p.path DESC
	LIMIT $3
`

const selectPostsParentTreeLimitByID = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, p.forum_slug
	FROM post p
	WHERE p.thread_id = $1 and p.path[1] IN (
		SELECT p2.path[1]
		FROM post p2
		WHERE p2.thread_id = $2 AND p2.parent_id IS NULL
		ORDER BY p2.path
		LIMIT $3
	)
	ORDER BY path
`

const selectPostsParentTreeLimitDescByID = `
SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, p.forum_slug
FROM post p
WHERE p.thread_id = $1 and p.path[1] IN (
    SELECT p2.path[1]
    FROM post p2
	WHERE p2.parent_id IS NULL and p2.thread_id = $2
	ORDER BY p2.path DESC
    LIMIT $3
)
ORDER BY p.path[1] DESC, p.path[2:]
`

const selectPostsParentTreeLimitSinceByID = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, p.forum_slug
	FROM post p
	WHERE p.thread_id = $1 and p.path[1] IN (
		SELECT p2.path[1]
		FROM post p2
		WHERE p2.thread_id = $2 AND p2.parent_id IS NULL and p2.path[1] > (SELECT p3.path[1] from post p3 where p3.id = $3)
		ORDER BY p2.path
		LIMIT $4
	)
	ORDER BY p.path
`

const selectPostsParentTreeLimitSinceDescByID = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, p.forum_slug
	FROM post p
	WHERE p.thread_id = $1 and p.path[1] IN (
		SELECT p2.path[1]
		FROM post p2
		WHERE p2.thread_id = $2 AND p2.parent_id IS NULL and p2.path[1] < (SELECT p3.path[1] from post p3 where p3.id = $3)
		ORDER BY p2.path DESC
		LIMIT $4
	)
	ORDER BY p.path[1] DESC, p.path[2:]
`

func selectAllPostsByThreadID(db *pgx.ConnPool, id int, limit int, desc bool,
	since int, sort string, posts *models.Posts) error {

	rows, err := doQuery(db, id, limit, desc, since, sort)
	if err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() {
		post := &models.Post{}
		err := scanPostRows(rows, post)
		if err != nil {
			return err
		}

		*posts = append(*posts, post)
	}

	return nil
}

func doQuery(db *pgx.ConnPool, id int, limit int, desc bool,
	since int, sort string) (*pgx.Rows, error) {
	var rows *pgx.Rows
	var err error
	switch sort {
	case "":
		fallthrough
	case "flat":
		if since > 0 {
			if desc {
				rows, err = db.Query(selectPostsFlatLimitSinceDescByID, id,
					since, limit)
			} else {
				rows, err = db.Query(selectPostsFlatLimitSinceByID, id,
					since, limit)
			}
		} else {
			if desc == true {
				rows, err = db.Query(selectPostsFlatLimitDescByID, id, limit)
			} else {
				rows, err = db.Query(selectPostsFlatLimitByID, id, limit)
			}
		}
	case "tree":
		if since > 0 {
			if desc {
				rows, err = db.Query(selectPostsTreeLimitSinceDescByID, id,
					since, limit)
			} else {
				rows, err = db.Query(selectPostsTreeLimitSinceByID, id,
					since, limit)
			}
		} else {
			if desc {
				rows, err = db.Query(selectPostsTreeLimitDescByID, id, limit)
			} else {
				rows, err = db.Query(selectPostsTreeLimitByID, id, limit)
			}
		}
	case "parent_tree":
		if since > 0 {
			if desc {
				rows, err = db.Query(selectPostsParentTreeLimitSinceDescByID, id, id,
					since, limit)
			} else {
				rows, err = db.Query(selectPostsParentTreeLimitSinceByID, id, id,
					since, limit)
			}
		} else {
			if desc {
				rows, err = db.Query(selectPostsParentTreeLimitDescByID, id, id,
					limit)
			} else {
				rows, err = db.Query(selectPostsParentTreeLimitByID, id, id,
					limit)
			}
		}
	}

	return rows, err
}
