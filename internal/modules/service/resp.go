package service

import (
	"encoding/json"
	"log"

	"github.com/valyala/fasthttp"
)

func marshalResp(resp json.Marshaler) []byte {
	json, err := resp.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}

	return json
}

func resp(c *fasthttp.RequestCtx, resp json.Marshaler, status int) {
	c.SetContentType("application/json")
	c.SetStatusCode(status)
	c.Write(marshalResp(resp))
}
