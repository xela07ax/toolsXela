## Канальный логер

### Примеры
- [Example 1](examples/main_webPage_поВебсокетам.go)
`Трансляция логов в web page по вебсокетам`

Запуск
```sh
go run main.go
```
Демонстрация
[![Demo docker api Xela golang](examples/logerWsWepPage-1.gif)](examples/logerWsWepPage-1.mp4)

  Возможна отправка сообщений с виртуальной консоли, обработка более продвинутым способом.  
Получение сообщения в `client.go`
```go
...
func (c *Client) writePump() {
...
		case message, ok := <-c.send:
            log.Printf("msg:%s",message)
            // Сообщение в текстовом формате
            	c.hub.WebSocketOutput <- message
            // Отправляем сообщение другому обработчику
			...
        case <-ticker.C:
...
```

- [Example 2](examples/main_webPage_поВебсокетам_игнор.go)
Сообщения которые отправлены по Вэбсокету отправляются назад, будем игноррировать их, если не пришло ничего интересного.


Демонстрация
[![Demo docker api Xela golang](examples/logerWsWepPage-2.gif)](examples/logerWsWepPage-2.mp4)

В примере видно как отправленная структура распозналась, а все остальное игноррировалось