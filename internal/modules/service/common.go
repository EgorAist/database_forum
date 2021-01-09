package service

import (
	"strconv"

	"github.com/valyala/fasthttp"
)

func getBool(k string, args *fasthttp.Args) bool {
	v := args.Peek(k)
	if v != nil && v[0] == 't' {
		return true
	}
	return false
}

func postIDToInt(ctx *fasthttp.RequestCtx) int32 {
	id, _ := strconv.Atoi(ctx.UserValue("id").(string))
	return int32(id)
}
