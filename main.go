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
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
<title>Upload Files</title>
</head>
<body>
    <h2>File Upload</h2>
    Select file
    <input type="file" id="filename" />
    <br>
    <input type="button" value="Connect" onclick="connectChatServer()" />
    <br>
    <input type="button" value="Upload" onclick="sendFile()" />
    <script>
        var ws;
        function connectChatServer() {
            ws = new WebSocket("ws://{{.}}/common");
            ws.binaryType = "arraybuffer";
            ws.onopen = function() {
                alert("Connected.")
            };
            ws.onmessage = function(evt) {
                alert(evt.msg);
            };
            ws.onclose = function() {
                alert("Connection is closed...");
            };
            ws.onerror = function(e) {
                alert(e.msg);
            }
        }
        function sendFile() {
            var file = document.getElementById('filename').files[0];
            var reader = new FileReader();
            var rawData = new ArrayBuffer();
            reader.loadend = function() {
            }
            reader.onload = function(e) {
                rawData = e.target.result;
                ws.send(rawData);
                alert("the File has been transferred.")
            }
            reader.readAsArrayBuffer(file);
        }
    </script>
</body>
</html>
`))
