# Archive not official library

В связи с постоянной путаницы путей при создании zip архивов, принято решение рассмотреть альтернативные библиотеки.

1. Самая ранняя мной изученная и показала, что работает
https://github.com/mholt/archiver
- Относительные пути
```go
package main

import (
	"compress/flate"
	"fmt"
	"github.com/mholt/archiver/v3"
)

func main() {
	//Создаем каталоги для измененных проектов
	z := archiver.Zip{
		CompressionLevel:       flate.DefaultCompression,
		MkdirAll:               true,
		SelectiveCompression:   true,
		ContinueOnError:        false,
		OverwriteExisting:      false,
		ImplicitTopLevelFolder: false,
	}

	err := z.Archive([]string{"."}, "test.zip")
	fmt.Println(err)
}
```
``
archDll.go  go.mod  go.sum  README.md  tmp/
``
- Абсолютные пути
```go
err := z.Archive([]string{"D:\\Projects\\toolsXela\\toolsXela\\arch2"}, "test.zip")
```
``arch2``



