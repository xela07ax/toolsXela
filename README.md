ToolPack v2
===
OLD Docs!!!
Пакет ToolPack предназначен для упрощения рутинных задач в проектах написанных на языке Go!

Он расширяется предоставляя все больше таких возможностей как Lodash, но для Golang.

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

	Логирование
- Logger - Предоставляет возможность вести логирование сразу в файл и в консоль
	logger = new(tp.Logger) // Создаем объект в который будем передавать что логировать и его же передавать пакетам которые будут туда писать
	logger.SetFilePath("log.txt") // Устанавливаем дирректорию до файла куда писать
	logger.LogPrintln("@STARTING DOWNLOAD@", "Моя программа запишет тот текст в лог") //Пример строки для логирования, строка запишется в файл и выведется на экран
- Round - Округление float64
    f = tp.Round(5.867868, 2) //5.86
- SaveStruct и LoadStruct - Сохранение и загрузка структуры из Json

## Installation

The minimum requirement of Go is **1.6**.

To use with latest changes:

```sh
$ go get github.com/xela07ax/toolPack
```
```

Please add `-u` flag to update in the future.

## Getting Help

## License

This project is under Apache v2 License. See the [LICENSE](LICENSE) file for the full license text.
