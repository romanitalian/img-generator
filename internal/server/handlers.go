package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/romanitalian/img-generate/pkg/img"
	"github.com/valyala/fasthttp"
	"log"
	"strings"
)

func rendImg(ctx *fasthttp.RequestCtx, buffer *bytes.Buffer) {
	ctx.SetContentType("image/jpeg")
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(buffer.Bytes())
}

func faviconHandler(ctx *fasthttp.RequestCtx) {
	favicon, err := img.GenerateFavicon()
	if err != nil {
		log.Println(err)
	}
	rendImg(ctx, favicon)
}

func imgHandler(ctx *fasthttp.RequestCtx) {
	params := strings.Split(fmt.Sprintf("%v", ctx.Value("params")), "/")
	buffer, err := img.Generate(params)
	if err != nil {
		log.Println(err)
	}
	rendImg(ctx, buffer)
}

func pingHandler(ctx *fasthttp.RequestCtx) {
	_, err := ctx.WriteString("PONG")
	if err != nil {
		log.Println(err)
	}
}
func robotsHandler(ctx *fasthttp.RequestCtx) {
	_, err := ctx.WriteString("robots")
	if err != nil {
		log.Println(err)
	}
}

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func userHandler(ctx *fasthttp.RequestCtx) {
	user := &User{Name: "roman", Email: "some@example.com"}
	jsonBody, err := json.Marshal(user)
	if err != nil {
		ctx.Error("json marshal fail", 500)
		return
	}
	ctx.SetContentType("application/json; charset=utf-8")
	ctx.SetStatusCode(200)
	ctx.Response.SetBody(jsonBody)

	return
}
