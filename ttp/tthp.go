package ttp

import (
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
	Httpjson(w, r, Notify)
}