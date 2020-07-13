package main

import (
	"./balancer"
	"./models"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)
import "./ConveerMinions/Conveer_Mail101"
import "./toolsXela/tp"
import "./toolsXela/chLogger"

const (
	pathTmp = "tmp/"
	tableName = "go_program"
)

func main() {
	// Создаем логер
	dir, err := tp.BinDir()
	tp.Fck(err)
	FullLogPath := filepath.Join(dir,"log")
	err = tp.CheckMkdir(FullLogPath)
	tp.Fck(err)
	logEr := chLogger.NewChLoger(FullLogPath)
	logEr.RunMinion()
	logEr.ChInLog<- [4]string{"Welcome","nil",fmt.Sprintf("Вас приветствует \"СТП Контроллер\" v1.0 (230919) \n")}

	// Открываем конфиг
	configDir := filepath.Join(dir,"../config.json")
	fi,err := tp.OpenReadFile(configDir)
	if err != nil {
		logEr.ChInLog <- [4]string{"Configurator","nil",fmt.Sprintf("Ошибка при открытии конфигурации %s: %s\n",configDir,err),"1"}
		time.Sleep(1*time.Second)
		os.Exit(1)
	}

	config := make(map[string]string)
	json.Unmarshal(fi,&config)

	// Создадим конвеер для запуска фоновой программы
	inCh0 := make(chan models.Command, 10000) // Создаем тестовый канал с буфером в 10000, в который наполним немного синтетики
	conveerOne := Conveer_Mail101.NewBalancer("Mail101",logEr, inCh0) // Инициализируем объект конвеера, передадим канал для работы и еще что нибудь
	logEr.ChInLog <- [4]string{"Main","nil",fmt.Sprintf("Запускаем первого миньона")}
	conveerOne.RunMinion() // Запускаем первого миньона

	//balancer.ReadLetters(dir,pythonPath, inCh0,conveerOne.OutCh ,logEr)
	stopConveerOne := make(chan bool)
	go conveerOne.SignalStoper(stopConveerOne)
	logEr.ChInLog <- [4]string{"Main","nil",fmt.Sprintf("Запускаем балансир")}
	balancer.WriteLetter(dir, config, inCh0,conveerOne.OutCh, stopConveerOne, conveerOne.GoodBy, logEr)

	time.Sleep(1*time.Second) //Для того, что бы в консоле не перемешались тексты выводв, добавим небольшую паузу
	fmt.Println("Конвеер conveerOne полностью запущен и готов работать")

	// Метод завершения конвеера 2
	//go conveerOne.Stop()  //Запускаем функцию остановки в фоне



	fmt.Println("Всем спасибо за внимание")
}
