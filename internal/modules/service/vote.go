package service

import (
	"github.com/EgorAist/database_forum/internal/models"
	"github.com/EgorAist/database_forum/internal/modules/database"
	"github.com/valyala/fasthttp"
)

func (self *ForumPgsql) ThreadVote(ctx *fasthttp.RequestCtx) {
	thread := &models.Thread{}
	vote := &models.Vote{}

	vote.UnmarshalJSON(ctx.PostBody())
	err := database.VoteCreate(self.db, ctx.UserValue("slug_or_id").(string),
		thread, vote)

	if err != nil {
		switch err {
		case database.ErrThreadNotFound:
			resp(ctx, Error, fasthttp.StatusNotFound)
			return
		}
		resp(ctx, Error, fasthttp.StatusInternalServerError)
		return
	}
	resp(ctx, thread, fasthttp.StatusOK)
	return
}
