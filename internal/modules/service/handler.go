package service

import (
	"github.com/valyala/fasthttp"
)

type ForumHandler interface {
	ForumCreate(*fasthttp.RequestCtx)
	UserCreate(*fasthttp.RequestCtx)
	ForumGetOne(*fasthttp.RequestCtx)
	UserGetOne(*fasthttp.RequestCtx)
	UserUpdate(*fasthttp.RequestCtx)
	ThreadCreate(*fasthttp.RequestCtx)
	ThreadVote(*fasthttp.RequestCtx)
	ThreadGetOne(*fasthttp.RequestCtx)
	ThreadUpdate(*fasthttp.RequestCtx)
	ForumGetThreads(*fasthttp.RequestCtx)
	PostsCreate(*fasthttp.RequestCtx)
	Clear(*fasthttp.RequestCtx)
	Status(*fasthttp.RequestCtx)
	PostGetOne(*fasthttp.RequestCtx)
	PostUpdate(*fasthttp.RequestCtx)
	ThreadGetPosts(*fasthttp.RequestCtx)
	ForumGetUsers(*fasthttp.RequestCtx)
}
