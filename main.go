package main


import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/xela07ax/toolsXela/chLogger"
	"github.com/xela07ax/toolsXela/hub"
	"github.com/xela07ax/toolsXela/hub/blog"
	"net/http"
	"time"
)

var addr = flag.String("addr", ":8187", "http service address")
func main()  {
	// 1. Настройка хаба для клиентов
	flag.Parse()
	fmt.Println(*addr)
	fmt.Println("-main->start[newHub]")
	hubib := hub.NewHub(false)
	go hubib.Run()


	// Для коннекта нужно прокинуть наружу
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		blog.Home(w, r)
		//serveWs(hubib, w, r)
	})
	http.HandleFunc("/wsLog", func(w http.ResponseWriter, r *http.Request) {
		hubib.ServeWs(w, r)
		// blog.Home(w, r)
		// serveWs(hubib, w, r)
	})
	go http.ListenAndServe(*addr, nil)
	//if err != nil {
	//	log.Fatal("ListenAndServe: ", err)
	//}

	// 2.  Настройка логера
	fmt.Println("-main->end[Настройка логера]")
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
	// 3. Настройка обработчика входящих сообщений
	// канал для приема <- hubib.WebSocketOutput
	type Notyfy struct {
		Name string
		Text string
		Data []byte
	}
// {"Name":"Run","Text":"Проект запуск"}
	go func() {
		for {
			msg := <- hubib.WebSocketOutput
			//fmt.Printf("%s",msg)
			// Проблема в том, что отправленное сообщение немедленно возвращается
			// А потому надо принимать только команды, будем парсить формат
			// Сюда должна прийти структура новой команды, или игноррируем
			var command Notyfy
			err := json.Unmarshal(msg, &command)
			if err != nil {
				continue
			}
			logEr.ChInLog <- [4]string{"Anonimouse","nil",fmt.Sprintf("%s",msg)}
		}

	}()

	fmt.Println("-main->end[RunMinion]")
	//fmt.Println("-main->start[hub.run]-p2")
	time.Sleep(5*time.Second)
	fmt.Println("-main->end[newHub]")
	logEr.ChInLog <- [4]string{"Welcome","nil",fmt.Sprintf("Вас приветствует \"Silika-FileManager Контроллер\" v1.1 (11112020) \n")}
	fmt.Println("-main->wait")
	time.Sleep(5*time.Second)
	fmt.Println("-main->end[newHub]")
	logEr.ChInLog <- [4]string{"Welcome","nil",fmt.Sprintf("Вас приветствует \"Silika-FileManager Контроллер\" v1.1 (11112020) \n")}
	fmt.Println("-main->wait")
	time.Sleep(5*time.Second)
	fmt.Println("-main->end[newHub]")
	logEr.ChInLog <- [4]string{"Welcome","nil",fmt.Sprintf("Вас приветствует \"Silika-FileManager Контроллер\" v1.1 (11112020) \n")}
	fmt.Println("-main->wait")
	time.Sleep(1*time.Second)
}