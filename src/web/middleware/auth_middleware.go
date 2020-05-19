package middleware

import (
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/sessions"
)

func NewAuth() context.Handler{
	u := &auth{}
	return u.auth
}

type auth struct {

}
func (this * auth) auth(ctx context.Context)   {
	sess:=sessions.Get(ctx)
	if auth,err:=sess.GetBoolean("authenticated");!auth||err!=nil{
		if ctx.IsAjax(){
			ctx.StatusCode(401)
			ctx.StopExecution()
			return
		}
		ctx.Redirect("/login")
		return
	}
	ctx.Next()
}