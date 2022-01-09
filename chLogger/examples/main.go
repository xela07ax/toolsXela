package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)
import "../../../toolsXela/tp"
import "../../../toolsXela/chLogger"

const (
	pathTmp = "tmp/"
	tableName = "go_program"
)

func main() {
	// Создаем логер
	dir, err := tp.BinDir()
	tp.Fck(err)
	// Открываем конфиг
	configDir := filepath.Join(dir,"config.json")
	fi,err := tp.OpenReadFile(configDir)
	tp.FckText(fmt.Sprintf("Ошибка при открытии конфигурации %s",configDir),err)
	var config chLogger.Config // или make(map[string]string)
	tp.Fck(json.Unmarshal(fi,&config))
	config.Dir = filepath.Join(dir,config.Dir)

	//d2 := 300* time.Millisecond // интервал
	logEr := chLogger.NewChLoger(&config)
	logEr.RunMinion()
	logEr.ChInLog<- [4]string{"Welcome","nil",fmt.Sprintf("Вас приветствует \"СТП Контроллер\" v1.0 (230919) \n")}

	// Открываем конфиг
	if err != nil {
		logEr.ChInLog <- [4]string{"Configurator","nil",fmt.Sprintf("Ошибка при открытии конфигурации %s: %s\n",configDir,err),"1"}
		time.Sleep(1*time.Second)
		os.Exit(1)
	}

	logEr.ChInLog <- [4]string{"Main","nil",fmt.Sprintf("Тестовое сообщение")}
	logEr.ChInLog <- [4]string{"nagi","nil",fmt.Sprintf("Тестовое сообщение1")}
	logEr.ChInLog <- [4]string{"nagi","unit1",fmt.Sprintf("Тестовое сообщение2")}
	logEr.ChInLog <- [4]string{"republic","nil",fmt.Sprintf("Тестовое сообщение3")}
	logEr.ChInLog <- [4]string{"Da","unit1",fmt.Sprintf("Тестовое сообщение3")}
	logEr.ChInLog <- [4]string{"Da","nil",fmt.Sprintf("Тестовое сообщение3")}
	time.Sleep(1*time.Second) //Для того, что бы в консоле не перемешались тексты выводв, добавим небольшую паузу
	fmt.Println("Всем спасибо за внимание")
}
