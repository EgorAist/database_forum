package service

import (
	"log"

	"github.com/EgorAist/database_forum/internal/models"
	"github.com/EgorAist/database_forum/internal/modules/database"
	"github.com/valyala/fasthttp"
)

func (self ForumPgsql) ForumCreate(ctx *fasthttp.RequestCtx) {
	forum := &models.Forum{}
	err := forum.UnmarshalJSON(ctx.PostBody())
	if err != nil {
		log.Println(err)
		return
	}
	err = database.CreateForum(self.db, forum)
	if err != nil {
		switch err {
		case database.ErrForumNotFound:
			resp(ctx, Error, fasthttp.StatusNotFound)
			return
		case database.ErrForumConflict:
			err := database.SelectForum(self.db, forum)
			if err != nil {
				resp(ctx, Error, fasthttp.StatusInternalServerError)
				return
			}
			resp(ctx, forum, fasthttp.StatusConflict)
			return
		default:
			resp(ctx, Error, fasthttp.StatusInternalServerError)
			return
		}
	}

	resp(ctx, forum, fasthttp.StatusCreated)
}

func (self *ForumPgsql) ForumGetOne(ctx *fasthttp.RequestCtx) {
	forum := &models.Forum{}
	forum.Slug = ctx.UserValue("slug").(string)
	err := database.SelectForum(self.db, forum)
	if err != nil {
		if err == database.ErrForumNotFound {
			resp(ctx, Error, fasthttp.StatusNotFound)
			return
		}
		resp(ctx, Error, fasthttp.StatusInternalServerError)
		return
	}

	resp(ctx, forum, fasthttp.StatusOK)
}

func (self *ForumPgsql) ForumGetThreads(ctx *fasthttp.RequestCtx) {
	threads := &models.Threads{}

	err := database.SelectAllThreadsByForum(self.db, ctx.UserValue("slug").(string),
		ctx.QueryArgs().GetUintOrZero("limit"),
		getBool("desc", ctx.QueryArgs()),
		string(ctx.QueryArgs().Peek("since")), threads)

	if err != nil {
		if err == database.ErrForumNotFound {
			resp(ctx, Error, fasthttp.StatusNotFound)
			return
		}
		resp(ctx, Error, fasthttp.StatusInternalServerError)
		return
	}
	resp(ctx, *threads, fasthttp.StatusOK)
	return
}

func (self *ForumPgsql) ForumGetUsers(ctx *fasthttp.RequestCtx) {
	users := &models.Users{}
	err := database.SelectAllUsersByForum(self.db, ctx.UserValue("slug").(string),
		ctx.QueryArgs().GetUintOrZero("limit"),
		getBool("desc", ctx.QueryArgs()),
		string(ctx.QueryArgs().Peek("since")), users)

	if err != nil {
		if err == database.ErrForumNotFound {
			resp(ctx, Error, fasthttp.StatusNotFound)
			return
		}
		resp(ctx, Error, fasthttp.StatusInternalServerError)
		return
	}
	resp(ctx, *users, fasthttp.StatusOK)
	return
}
