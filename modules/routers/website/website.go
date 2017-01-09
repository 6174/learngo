package website

import (
	"github.com/kataras/iris"
)

func Home(ctx *iris.Context) {
	type DT struct{
		Name string
	}
	ctx.Render("index.html", DT{Name: "iris"})
}
