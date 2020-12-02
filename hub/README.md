# Websocket hub
Хаб Вебсокет клиентов которые пришли к нам.


Реализация простого логирования в вебсокеты, консоль и файл

```go

func main() {
    
	hub := newHub()
    go hub.run()
    
	
	logEr := bcl.NewChLoger(&bcl.Config{
		IntervalMs:     300,
		ConsolFilterFn: map[string]int{"Front Http Server":  0},
		ConsolFilterUn: map[string]int{"Pooling": 1},
		Mode:           0,
		Dir:            "x-loger",
		Broadcast: hub.Input,
	})
	logEr.RunMinion()
	Logx("-main->end[newHub]")
	Logx("-main->start[hub.run]-p2")
	time.Sleep(1*time.Second)
	logEr.ChInLog <- [4]string{"Welcome","nil",fmt.Sprintf("Вас приветствует \"Silika-FileManager Контроллер\" v1.1 (11112020) \n")}

	Logx("-main->end[hub.run]")
	// По вебсокетам у нас будет логер
```
Полный пример
```go
package main


import (
	"bytes"
	"compress/flate"
	"fmt"
	"github.com/xela07ax/toolsXela/archiver"
	"github.com/xela07ax/toolsXela/tp"
)

func main()  {
	ts, _ := tp.BinDir()
	fmt.Printf("Bit: %s\n", ts)
	archiver.Print()
	z := archiver.Zip{
		CompressionLevel:       flate.DefaultCompression,
		MkdirAll:               true,
		SelectiveCompression:   true,
		ContinueOnError:        false,
		OverwriteExisting:      true,
		ImplicitTopLevelFolder: false,  //Неявная папка верхнего уровня
	}
	//err := z.Archive([]string{"D:\\Projects\\toolsXela\\toolsXela\\arch2"}, "test.zip")
	var b bytes.Buffer // A Buffer needs no initialization.
	err := z.ArchiveWriter([]string{"D:\\Projects\\toolsXela\\toolsXela\\arch2\\"}, &b)
	//f, err := tp.CreateOpenFile("tx.zip")
		//b.WriteTo(os.Stdout)
	fmt.Println(b.Len())
	//	b.WriteTo(f)
	//fmt.Println(err)
	//err = z.Unarchive("tx.zip", "txzip")
	//fmt.Println(err)
	r := bytes.NewReader(b.Bytes())
	fmt.Println(b.Len())
	err = z.UnarchReader(r,b.Len(),"txzip2")
	fmt.Println(err)
}
```
1. Запустим сервер
```sh
cd ../chat-pokemon
λ  go run .
```
2. Откроем страницу в браузере  
http://localhost:8180/
3. Отправим текст для логера, что бы он транслировал в вэб окно и тд.
```sh
$ curl -H "Content-Type: application/json" -X POST http://localhost:8180/wsx/sendMsg -d "Hello World"
```

<img src="ws-chat-loger.jpg" width="550" />
