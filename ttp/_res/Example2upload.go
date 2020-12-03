package main

import (
	"flag"
	"fmt"
	"github.com/xela07ax/toolsXela/tp"
	"github.com/xela07ax/toolsXela/ttp"
	"html/template"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8187", "http service address")
func main()  {
	http.HandleFunc("/upload", func (w http.ResponseWriter, r *http.Request) {
		fmt.Println("File Upload Endpoint Hit")
		dat, err := ttp.UploadFileBytes(w,r)
		if err != nil {
			ertx:= fmt.Sprintf("файл не скачан |ertx:%s",err)
			http.Error(w, ertx, http.StatusBadGateway)
			fmt.Print(ertx)
			return
		}
		fmt.Fprintf(w, "Successfully Uploaded File\n")
		f, err := tp.CreateOpenFile("tmp.psd")
		if err != nil {
			ertx:= fmt.Sprintf("создание файла |ERTX:%s\n",err)
			http.Error(w, ertx, http.StatusBadGateway)
			fmt.Print(ertx)
			return
		}
		defer f.Close()
		// write this byte array to our temporary file
		f.Write(dat)
		// return that we have successfully uploaded our file!
		fmt.Fprintf(w, "End Uploaded File\n")
	})
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
