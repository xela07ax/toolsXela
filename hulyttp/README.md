# HulyTTP - Huly ty tut polzaesh
## Http инструменты

- Resp  
Возвращаем rest ответ в стандартизированном варианте Notify
#### Пример
```go
import "github.com/xela07ax/toolsXela/hulyttp"
func sendMsg(w http.ResponseWriter, r *http.Request) {
	hulyttp.resp(w,r,"SendMsgNm", "Hello World", 0, true)
}
```
``
{"FuncName":"SendMsgNm","Text":"Hello World",Status":0,"Show":true,"UpdNum":0}
``