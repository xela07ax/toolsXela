ToolPack v2 ToolsXela
===

Пакет ToolPack предназначен для упрощения рутинных задач в проектах написанных на языке Go!
## Возможности
	Работа с файлами
- WriteFile - Запись байтов в файл
- CheckExistsFile - Проверка файла на существование
- BinDir - Возвращает путь до исполняемого бинарного файла
- FindFile - Цель этого модуля искать файл в папке где лежит исполняемый файл и в рабочей директории
- CheckMkdir - Создать папку если ее нет
- ReadfileWorkDirIntel - И ищет и возвращает текст
- ReadFile - Просто читает файл

	Работа с конфигурациями
- OpenINIfile - открытие INI файлов, но нужно предварительно задать структуру
    OpenConf("file.ini", &cfg_struct)

# Логирование
## Logger - Предоставляет возможность вести логирование сразу в файл и в консоль
Создаем объект в который будем передавать что логировать и его же передавать пакетам которые будут туда писать

```go
// Получим текущую директорию
	dir, err := tp.BinDir()
	tp.Fck(err)
// Запускаем логер
	logEr := chLogger.NewChLoger(&chLogger.Config{
		IntervalMs:     300,
		ConsolFilterFn: map[string]int{"Front Http Server":  0},
		ConsolFilterUn: map[string]int{"Pooling": 1},
		Mode:           0,
		Dir:            dir,
	})
	logEr.RunMinion()
	logEr.ChInLog <- [4]string{"Welcome","nil",fmt.Sprintf("Вас приветствует \"Silika-FileManager Контроллер\" v1.1 (11112020) \n")}

```
- Round - Округление float64
    f = tp.Round(5.867868, 2) //5.86
- SaveStruct и LoadStruct - Сохранение и загрузка структуры из Json

```sh
$ go get github.com/xela07ax/toolPack
```


Please add `-u` flag to update in the future.

## Getting Help
# Примеры

```go
import (
	"github.com/tebeka/selenium"
	"github.com/xela07ax/toolsXela/tp"
	)
		var err error
		we,err = wb.FindElement(selenium.ByID, "lookup_page1_tbod-tbd")
		if err != nil {return true,nil}
		bt, _ := we.Screenshot(true)	
		f,_ := tp.OpenWriteFile("file.png")
		f.Write(bt)
```
```go
package main

import (
	"../../lib/toolsXela/tp"
	"../../model"
	"encoding/json"
	"fmt"
	"path/filepath"
)
type DictRoles struct {
	Chat       string
	Dispatcher string
	Operator   string
	Dashboard  string
}
type Conf struct {
	Name int // Порт программы
	Roles DictRoles
	datail map[int]string
}

func main()  {
	dir, err := tp.BinDir()
	tp.Fck(err)
	// Открываем конфиг
	var config model.Config
	configDir := filepath.Join(dir,"config.json")
	fi,err := tp.OpenReadFile(configDir)
	if err != nil {
		fmt.Printf("Ошибка при открытии конфигурации %s: %s\n",configDir,err)
		tp.ExitWithSecTimeout(1)
	}
	err = json.Unmarshal(fi,&config)
	if err != nil {
		fmt.Printf("Ошибка чтения JSON %s: %s\n",configDir,err)
		tp.ExitWithSecTimeout(1)
	}
	fmt.Println(config)
	dat := Conf{
		Name:  453,
		Roles: DictRoles{
			Chat:       "Привет",
			Dispatcher: "мой",
			Operator:   "прекрвсный",
			Dashboard:  "мир",
		},
	}
	raw, _ := json.Marshal(dat)
	fFunc, err := tp.OpenWriteFile(filepath.Join(dir,"tstConfig.json"))
	tp.FckText(fmt.Sprintf("Логгер | Открытие файла: %s",fFunc),err)
	fFunc.Write(raw)
	fFunc.Close()
}
```
## License

This project is under Apache v2 License. See the [LICENSE](LICENSE) file for the full license text.
