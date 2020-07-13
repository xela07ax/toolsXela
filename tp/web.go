package tp

import (
	"encoding/json"
	"net/http"
)

func HttpBytes(w http.ResponseWriter, r *http.Request, res []byte) (err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(res)
	if err != nil{
		return err
	}
	return nil
}
func Httpjson(w http.ResponseWriter, r *http.Request, res interface{}) (err error) {
		detailjson, _ := json.Marshal(res)
		return HttpBytes(w,r,detailjson)
}