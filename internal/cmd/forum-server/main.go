package main

import (
	"flag"
	"log"

	"gopkg.in/jackc/pgx.v2"

	"github.com/EgorAist/database_forum/internal/modules/service"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

const port = "5000"

var databaseDsn = flag.String("db", "", "dsn database")
var portFlag = flag.String("port", "5000", "port")

func main() {
	flag.Parse()
	err := ListenAndServe(*portFlag, *databaseDsn)
	if err != nil {
		log.Fatal(err)
	}
}

func redirect(router *fasthttprouter.Router, handler service.ForumHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		path := string(ctx.Path())
		if path == "/api/forum/create" {
			handler.ForumCreate(ctx)
			return
		}
		router.Handler(ctx)
	}
}

func ListenAndServe(port, dsn string) error {
	config, err := pgx.ParseConnectionString(dsn)
	if err != nil {
		return err
	}
	handler := service.NewForumPgsql(&config)
	router := fasthttprouter.New()
	router.POST("/api/forum/:slug/create", handler.ThreadCreate)
	router.POST("/api/user/:nickname/create", handler.UserCreate)
	router.GET("/api/forum/:slug/details", handler.ForumGetOne)
	router.GET("/api/user/:nickname/profile", handler.UserGetOne)
	router.POST("/api/user/:nickname/profile", handler.UserUpdate)
	router.POST("/api/thread/:slug_or_id/vote", handler.ThreadVote)
	router.GET("/api/thread/:slug_or_id/details", handler.ThreadGetOne)
	router.POST("/api/thread/:slug_or_id/details", handler.ThreadUpdate)
	router.GET("/api/forum/:slug/threads", handler.ForumGetThreads)
	router.POST("/api/thread/:slug_or_id/create", handler.PostsCreate)
	router.POST("/api/service/clear", handler.Clear)
	router.GET("/api/service/status", handler.Status)
	router.POST("/api/post/:id/details", handler.PostUpdate)
	router.GET("/api/post/:id/details", handler.PostGetOne)
	router.GET("/api/thread/:slug_or_id/posts", handler.ThreadGetPosts)
	router.GET("/api/forum/:slug/users", handler.ForumGetUsers)

	return fasthttp.ListenAndServe(":"+port, redirect(router, handler))
}
