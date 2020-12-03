package ttp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func HttpReadBody(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return []byte{}, err
	}
	return b, nil
}

func UploadFileBytes(w http.ResponseWriter, r *http.Request) (dat []byte, err error) {
	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	// _ -> handler
	file, _, err := r.FormFile("myFile")
	if err != nil {
		return dat, fmt.Errorf("retrieving r body the File|ERTX:%s\n",err)
	}
	defer file.Close()
	//fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	//fmt.Printf("File Size: %+v\n", handler.Size)
	//fmt.Printf("MIME Header: %+v\n", handler.Header)
	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return fileBytes, fmt.Errorf("считывание байт с file |ERTX:%s\n",err)
	}
	return fileBytes, nil
}

