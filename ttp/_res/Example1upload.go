package main

import (
	"flag"
	"fmt"
	"github.com/xela07ax/toolsXela/tp"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8187", "http service address")
func main()  {
	http.HandleFunc("/upload", ploadFile)
	http.HandleFunc("/home", handler)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
func handler(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, r.Host)
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <meta http-equiv="X-UA-Compatible" content="ie=edge" />
    <title>Document</title>
  </head>
  <body>
    <form
      enctype="multipart/form-data"
      action="http://{{.}}/upload"
      method="post"
    >
      <input type="file" name="myFile" />
      <input type="submit" value="upload" />
    </form>
  </body>
</html>
`))

func UploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File Upload Endpoint Hit")

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)
	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("считывание байт с file |ERTX:%s\n",err)
		return
	}
	f, err := tp.CreateOpenFile("tmp.psd")
	if err != nil {
		fmt.Printf("создание файла |ERTX:%s\n",err)
		return
	}
	defer f.Close()
	// write this byte array to our temporary file
	f.Write(fileBytes)
	// return that we have successfully uploaded our file!
	fmt.Fprintf(w, "Successfully Uploaded File\n")
}
