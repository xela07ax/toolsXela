package main


import (
	"flag"
	"fmt"
	"github.com/xela07ax/toolsXela/hub"
	"github.com/xela07ax/toolsXela/chLogger"
	"log"
	"net/http"
	"time"
)

var addr = flag.String("addr", ":8187", "http service address")
func main()  {
	// 1. Настройка хаба для клиентов
	fmt.Println(addr)
	flag.Parse()
	fmt.Println("-main->start[newHub]")
	hubib := hub.NewHub(true)
	go hubib.Run()
	// Для коннекта нужно прокинуть наружу
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hubib, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	// 2.  Настройка логера
	time.Sleep(500*time.Millisecond)
	logEr := chLogger.NewChLoger(&chLogger.Config{
		IntervalMs:     300,
		ConsolFilterFn: map[string]int{"Front Http Server":  0},
		ConsolFilterUn: map[string]int{"Pooling": 1},
		Mode:           0,
		Dir:            "x-loger",
		Broadcast: hubib.Input,
	})
	logEr.RunMinion()

	//fmt.Println("-main->end[newHub]")
	//fmt.Println("-main->start[hub.run]-p2")
	time.Sleep(5*time.Second)
}