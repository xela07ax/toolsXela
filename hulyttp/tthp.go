package hulyttp

import (
	"github.com/xela07ax/toolsXela/tp"
	"net/http"
)

type Notify struct {
	FuncName string
	Text string
	Status int
	Show bool
	UpdNum int
}

func Resp (w http.ResponseWriter, r *http.Request, funcName string,text string, status int, show bool) {
	Notify := Notify {
		FuncName: funcName,
		Text:     text,
		Status:   status,
		Show:     show,
	}
	tp.Httpjson(w, r, Notify)
}