package main

import (
	"../models"
	"../toolsXela/chLogger"
	"../toolsXela/regx101"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)


func ReadLetters(dir string,pythonPath string, inCh0 chan <- [3]string, outCh0 <- chan [3]string, loger *chLogger.ChLoger)  {
	// inCh0[0] = ExecCommand
	// inCh0[1] = Домашняя директория
	// inCh0[2] = Команда
		// Разделитель аргументов |

	// outCh0[0] = Статус
		// "0" = error
		// "1" = ok
	// outCh0[1] = Комментарий
		// Текст ошибки
	// outCh0[2] = Текст


	fmt.Println(dir)
	pathMail101 := filepath.Join(dir,`../Mail101/py_Project`)
	//cmd := exec.Command("/usr/local/bin/python3","-V")
	inCh0 <- [3]string{pythonPath,pathMail101,"main.py"}

	// Открываем
	cmdOutputStr := <- outCh0
	loger.ChInLog <- [4]string{"Balancer","nil",fmt.Sprintf("Ответ пришел: %s",cmdOutputStr[0])}

	if cmdOutputStr[0] == "0" {
		loger.ChInLog <- [4]string{"Balancer","nil",fmt.Sprintf("Ошибка: %s",cmdOutputStr[1])}
		return
	}


	stopStr := "Goob By! | "
	resu:= strings.Index(cmdOutputStr[2], stopStr)

	// Prints true
	if resu == -1 {
		loger.ChInLog <- [4]string{"Conveer","nil","Писем нет :("}
	}
	pathLetter := cmdOutputStr[2][resu+len(stopStr):]
	loger.ChInLog <- [4]string{"Conveer","nil",fmt.Sprintf("Письмо получено: %s",pathLetter)}
	time.Sleep(2*time.Second)


	pathStp := filepath.Join(dir,`../Contain_stp/py_Project`)
	//cmd := exec.Command("/usr/local/bin/python3","-V")
	inCh0 <- [3]string{pythonPath,pathStp,"create_ticket.py"}
}

func WriteLetter(dir string,config map[string]string, inCh0 chan <- models.Command, outCh0 <- chan models.Reject, stopper chan <- bool, goodBy <- chan bool, loger *chLogger.ChLoger)  {
	// Инициальзируем Config.json
	// Находим
	var pythonPath string
	if val, ok := config["PythonPath"]; ok {
		pythonPath = val
	} else {
		loger.ChInLog <- [4]string{"Configurator", "nil", fmt.Sprintf("Ошибка при чтении параметра \"PythonPath\" в конфигурации"), "1"}
		time.Sleep(1 * time.Second)
		os.Exit(1)
	}

	var timeotCheckNewLetter int
	if val, ok := config["TimeotSecondCheckNewLetter"]; ok {
		num, err := strconv.Atoi(val)
		if err != nil {
			loger.ChInLog<- [4]string{"Configurator","nil",fmt.Sprintf("TimeotCheckNewLetter не является числом в конфигурации"),"1"}
			time.Sleep(1 * time.Second)
			os.Exit(1)
		}
		timeotCheckNewLetter = num
	} else {
		loger.ChInLog <- [4]string{"Configurator", "nil", fmt.Sprintf("Ошибка при чтении параметра \"TimeotCheckNewLetter\" в конфигурации"), "1"}
		time.Sleep(1 * time.Second)
		os.Exit(1)
	}

	pathMail101 := filepath.Join(dir,`../Mail101/py_Project`)
	pathStp := filepath.Join(dir,`../Contain_stp/py_Project`)
	//com_sys_Pythones := models.Command{Name:"SYS-Pythone3",Executer:"type",DirHome:".",Args:[]string{"-a", "python3"},Timeout:5}
	com_101_Login :=  models.Command{Name:"101-login",Executer:pythonPath,DirHome:pathMail101,Args:[]string{"login.py"},Timeout:5}
	com_101_GetLetter :=  models.Command{Name:"101-extractLetter",Executer:pythonPath,DirHome:pathMail101,Args:[]string{"extractLetter.py"},Timeout:60}
	com_101_MoveLetter :=  models.Command{Name:"101-MoveLetter",Executer:pythonPath,DirHome:pathMail101,Args:[]string{"moveLetter.py","Необходимо вставить значение"},Timeout:60}
	com_stp_createTicket := models.Command{Name:"STP-createTicket",Executer:pythonPath,DirHome:pathStp,Args:[]string{"create_ticket.py"},Timeout:60}


	// Сначала нам нужно залогиниться в почте
	loger.ChInLog <- [4]string{"Balancer","nil",fmt.Sprintf("Проверяем Python")}

	inCh0 <- com_101_GetLetter
	 // Нельзя логиниться первой операцией! Будет рекурсия.
	loger.ChInLog <- [4]string{"Balancer","nil",fmt.Sprintf("Команда отправлена")}
	for {
		select {
		case msg1 := <-outCh0:
			if msg1.Status == 0 {
				// Если операция успешна
				loger.ChInLog <- [4]string{"Balancer", "nil", fmt.Sprintf("Ответ пришел | Имя: %s | Статус %v | Текст: %s",msg1.Name, msg1.Status, msg1.MyStdout)}
				switch name := msg1.Name; name {
				case "101-login":
					loger.ChInLog<- [4]string{"Balancer","tryLogin",fmt.Sprintf("Мы в системе. Выполняем команду: %s",msg1.Command.Name)}
					inCh0 <- msg1.Command
				case "101-extractLetter":
					/*
					Существующие папки:
					 Входящие -> ($Inbox)
					 Черновики -> ($Drafts)
					 InsecureRequestWarning)
					 Отправленные -> ($Sent)
					 На контроль -> ($Follow-Up)
					 Все документы -> ($All)
					 Нежелательная -> ($JunkMail)
					 Корзина -> ($SoftDeletions)
					 Цепочки сообщений -> Threads
					 Дана консультация -> 147718583c02b951432583bd004dfe34
					 Необработанные -> 82f18a94e8694f8a432583e700547569
					 Обработанные -> 31ab09d826714b15432583bd004e03e5
					 Отработанные -> d00cc3eb446f964e432583e7005485bf
					 Архивация -> Archive
					 Правила -> (Rules)
					 Бланк -> Stationery
					 */
					loger.ChInLog<- [4]string{"Balancer",name,fmt.Sprintf("Считываем информацию о ходе работ")}
					// Сначала проверим количество писем в ящике
					var compRegEx = regexp.MustCompile(`Find\sletters\sin\sfolder:(?P<CntLetters>\d+)`)
					paramsMap := regx101.ParseRxNammed(compRegEx,msg1.MyStdout)
					if stat, ok := paramsMap["CntLetters"]; ok {
						num, err := strconv.Atoi(stat)
						if err != nil {
							loger.ChInLog<- [4]string{"Balancer",name,fmt.Sprintf("Количество писем не является числом"),"1"}
							stopper <- true
							continue
						}
						if num == 0 {
							loger.ChInLog<- [4]string{"Balancer",name,fmt.Sprintf("Писем нет, ожидаем")}
							time.Sleep(time.Duration(timeotCheckNewLetter) * time.Second)
							inCh0 <- com_101_GetLetter
							continue
						}
					} else {
						loger.ChInLog<- [4]string{"Balancer",name,fmt.Sprintf("Не удалось выяснить количество писем в ящике"),"1"}
						stopper <- true
						continue
					}
					compRegEx = regexp.MustCompile(`101Parse\sstatus:\s(?P<ParseStatus>\w+)\s\|\sfolderId:(?P<FolderID>\w+)\s\|\sUserID:(?P<UserID>\w+)\s\|\sPath:(?P<PathCfg>\w+)`)
					paramsMap = regx101.ParseRxNammed(compRegEx,msg1.MyStdout)
					if stat, ok := paramsMap["ParseStatus"]; ok {
						if stat == "none"{
							loger.ChInLog<- [4]string{"Balancer","101-extractLetter",fmt.Sprintf("Не удалось собрать все поля для заявки, переновим в Необработанные")}
							//1- oldfolderID
							//2- UserID
							//3- destinationFolderID
							com_101_MoveLetter.Args = []string{"moveLetter.py", paramsMap["FolderID"],paramsMap["UserID"],"82f18a94e8694f8a432583e700547569"}
							inCh0 <- com_101_MoveLetter
							continue
						} else if stat == "yes" {
							loger.ChInLog<- [4]string{"Balancer","101-extractLetter",fmt.Sprintf("Удалось собрать все поля, регистрируем заявку и переновим в Отработанные")}
							com_stp_createTicket = models.Command{Name:"STP-createTicket",Executer:pythonPath,DirHome:pathStp,Args:[]string{"create_ticket.py",pathStp,paramsMap["PathCfg"]},Timeout:60,Cfg:paramsMap}
							inCh0 <- com_stp_createTicket
							continue
						}

					}
					loger.ChInLog<- [4]string{"Balancer",name,fmt.Sprintf("Неизвестная ошибка"),"1"}
					stopper <- true
				case "101-MoveLetter":
					loger.ChInLog<- [4]string{"Balancer",name,fmt.Sprintf("Считываем информацию о ходе работ")}
					var compRegEx = regexp.MustCompile(`Move\sletter:(?P<UserID>\w+)\s\|\sstatus:(?P<Status>\w+)`)
					paramsMap := regx101.ParseRxNammed(compRegEx,msg1.MyStdout)
					if stat, ok := paramsMap["Status"]; ok {
						if stat == "success"{
							loger.ChInLog<- [4]string{"Balancer","101-MoveLetter",fmt.Sprintf("Письмо успешно перемещено, читаем следующее письмо")}
							inCh0 <- com_101_GetLetter
							continue
						}
					}
					loger.ChInLog<- [4]string{"Balancer",name,fmt.Sprintf("Неизвестная ошибка"),"1"}
					stopper <- true
				case "STP-createTicket":
					loger.ChInLog<- [4]string{"Balancer",name,fmt.Sprintf("Считываем информацию о ходе работ")}
					resu := strings.Index(msg1.MyStdout, "Ticket created success")
					if resu != -1 {
						loger.ChInLog<- [4]string{"Balancer","nil",fmt.Sprintf("Заявка создана успешно, перемещаем")}
						com_101_MoveLetter.Args = []string{"moveLetter.py", msg1.Command.Cfg["FolderID"],msg1.Command.Cfg["UserID"],"d00cc3eb446f964e432583e7005485bf"}
						inCh0 <- com_101_MoveLetter
						continue
					}
					loger.ChInLog<- [4]string{"Balancer",name,fmt.Sprintf("Неизвестная ошибка"),"1"}
					stopper <- true
				default:
					loger.ChInLog<- [4]string{"Balancer","nil",fmt.Sprintf("Не удалось идентифицировать ответ от Phyton")}
					stopper <- true
					// freebsd, openbsd,
					// plan9, windows...
				}
			} else if msg1.Status == 1 {
				// Если операция не успешна
				loger.ChInLog<- [4]string{"Balancer","nil",fmt.Sprintf("Ответ пришел с ОШИБКОЙ | Имя: %s | Текст: %s", msg1.Name, msg1.MyStderr)}
				if msg1.Comment == "Timeout" {
					inCh0 <- msg1.Command
				}
				switch name := msg1.Name; name {
				case "101-login":
					loger.ChInLog<- [4]string{"Balancer","tryLogin",fmt.Sprintf("Пробуем залогиниться еще раз")}
					inCh0 <- com_101_Login
				default:
					resu := strings.Index(msg1.MyStderr, "ERROR AUTORIZATION")
					if resu != -1 {
						loger.ChInLog<- [4]string{"Balancer","nil",fmt.Sprintf("Логинимся")}
						inCh0 <- com_101_Login
						loger.ChInLog <- [4]string{"Loginner","tryLogin",fmt.Sprintf("Команда залогиниться отправлена")}
						continue
					}
					resu = strings.Index(msg1.MyStderr, "Failed to establish a new connection")
					if resu != -1 {
						loger.ChInLog<- [4]string{"Balancer",name,fmt.Sprintf("Проблемы с сервером, пробуем повторно | Out:%s",msg1.MyStdout),"1"}
						inCh0 <- msg1.Command
						continue
					}
					loger.ChInLog<- [4]string{"Balancer",name,fmt.Sprintf("Неизвестная ошибка при выполнении скрипта Phyton | Out:%s",msg1.MyStdout),"1"}
					stopper <- true
				}
			}
			time.Sleep(1 * time.Second)
			//stopper <- true // На продакшене закоментировать, что б не отключалась программа
		case <-goodBy:
			return
		}
	}

	fmt.Println(pathMail101, pathStp,com_101_Login,com_101_GetLetter,com_stp_createTicket)
	//attach := "tmp/attach/523976EDCACCD41E3CB5779C00BCBC32/config.json"
	//pathStp := filepath.Join(dir,`../Contain_stp/py_Project`)
	// Открываем

	//if cmdOutputStr.Status == 0 {
		//loger.ChInLog <- [4]string{"Balancer","nil",fmt.Sprintf("Ошибка выполнения: %s \n Текст выполнения: %s",cmdOutputStr.MyStderr,cmdOutputStr.MyStdout)}
		//return
	//}

}