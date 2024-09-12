package csa

import (

	"fmt"
	"net/http"
	"crypto/md5"
	"encoding/hex"

	"github.com/julienschmidt/httprouter"

	"github.com/govwa/util"
	"github.com/govwa/user/session"
	"github.com/govwa/util/middleware"
)

type XSS struct{
	Name string
}
func New()XSS{
	return XSS{}
}
func (self XSS)SetRouter(r *httprouter.Router){
	mw := middleware.New()
	r.GET("/csa", mw.LoggingMiddleware(mw.CapturePanic(mw.AuthCheck(csaHandler))))
	r.POST("/verify", mw.LoggingMiddleware(mw.CapturePanic(mw.AuthCheck(verifyHandler))))
}

type JsonRes struct{
	Code int `json:"code"`
}

func csaHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	s := session.New()
	uid := s.GetSession(r, "id")

	data := make(map[string]interface{})
	data["title"] = "Client Side Authentication"

	id := fmt.Sprintf("<script> var uid=%s </script>", uid)

	data["js"] = util.ToHTML(id)
	
	util.SafeRender(w,r, "template.csa", data)
}

func verifyHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	if r.Method == "POST"{
		sotp := "a587cd6bf1e49d2c3928d1f8b86f248b" 
		otp := r.FormValue("otp")
		res := JsonRes{}
		if sotp != Md5Sum(otp){
			res.Code = 0
		}else{
			res.Code = 1
		}
		util.RenderAsJson(w, res)
	}
}

func Md5Sum(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
