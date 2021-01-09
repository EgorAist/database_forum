package service

import (
	"github.com/EgorAist/database_forum/internal/models"
	"github.com/EgorAist/database_forum/internal/modules/database"
	"github.com/valyala/fasthttp"
)

func (self *ForumPgsql) UserCreate(ctx *fasthttp.RequestCtx) {
	u := &models.User{}
	u.Nickname = ctx.UserValue("nickname").(string)
	u.UnmarshalJSON(ctx.PostBody())
	err := database.CreateUser(self.db, u)

	if err != nil {
		switch err {
		case database.ErrUserConflict:
			users, err := database.SelectUsersWithNickOrEmail(self.db, u.Nickname, u.Email)
			if err != nil {
				resp(ctx, Error, fasthttp.StatusInternalServerError)
				return
			}
			resp(ctx, users, fasthttp.StatusConflict)
			return
		default:
			resp(ctx, Error, fasthttp.StatusInternalServerError)
			return
		}
	}

	resp(ctx, u, fasthttp.StatusCreated)
	return
}

func (self *ForumPgsql) UserGetOne(ctx *fasthttp.RequestCtx) {
	user := &models.User{}
	user.Nickname = ctx.UserValue("nickname").(string)
	err := database.SelectUser(self.db, user)

	if err != nil {
		if err == database.ErrUserNotFound {
			resp(ctx, Error, fasthttp.StatusNotFound)
			return
		}
		resp(ctx, Error, fasthttp.StatusInternalServerError)
		return
	}

	resp(ctx, user, fasthttp.StatusOK)
}

func (self *ForumPgsql) UserUpdate(ctx *fasthttp.RequestCtx) {
	user := &models.User{}
	user.Nickname = ctx.UserValue("nickname").(string)
	userUpdate := &models.UserUpdate{}
	userUpdate.UnmarshalJSON(ctx.PostBody())

	err := database.UpdateUser(self.db, user, userUpdate)
	if err != nil {
		switch err {
		case database.ErrUserNotFound:
			resp(ctx, Error, fasthttp.StatusNotFound)
			return
		case database.ErrUserConflict:
			resp(ctx, Error, fasthttp.StatusConflict)
			return
		}
		resp(ctx, Error, fasthttp.StatusInternalServerError)
		return
	}

	resp(ctx, user, fasthttp.StatusOK)
	return
}
