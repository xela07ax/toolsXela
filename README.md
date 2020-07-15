ToolPack v2
===
OLD Docs!!!
Пакет ToolPack предназначен для упрощения рутинных задач в проектах написанных на языке Go!

Он расширяется предоставляя все больше таких возможностей как Lodash, но для Golang.
	}
## Возможности

	Работа с архивами
- ZipFiles и Unzip - для архивирования и разархивирование файлов в формате ZIP
    files := []string{"example.csv", "data.csv"}
    err := ZipFiles("done.zip", files)

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
	logger = new(tp.Logger) // Создаем объект в который будем передавать что логировать и его же передавать пакетам которые будут туда писать
	logger.SetFilePath("log.txt") // Устанавливаем дирректорию до файла куда писать
	logger.LogPrintln("@STARTING DOWNLOAD@", "Моя программа запишет тот текст в лог") //Пример строки для логирования, строка запишется в файл и выведется на экран
	logerDaley_Ms
```go
// Получим текущую директорию
	dir, err := tp.BinDir()
	tp.Fck(err)
// Создаем логер
	FullLogPath := filepath.Join(dir,"log/")
	err = tp.CheckMkdir(FullLogPath)
	if err != nil {
		fmt.Printf("Ошибка временной папки: %s | %s\n",FullLogPath,err)
		tp.ExitWithSecTimeout(1)
		}
// Запускаем логер
	logEr := chLogger.NewChLoger(FullLogPath,300)
	logEr.RunMinion()
	logEr.ChInLog <- [4]string{"Welcome","nil",fmt.Sprintf("Вас приветствует \"Бот Контроллер\" v1.1 (091219) \n")}
	if err != nil {
		logEr.ChInLog <- [4]string{"Configurator","nil",fmt.Sprintf("Ошибка при открытии 32-битного ключа \".key\" %s: %s\n",configDir,err),"1"}
	}
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
