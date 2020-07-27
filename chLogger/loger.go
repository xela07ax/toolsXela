package chLogger

import (
	"fmt"
	"github.com/xela07ax/toolsXela/tp"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Config struct {
	IntervalMs int // интервал обновления
	ConsolFilterFn map[string]int // map[funcName]mode
	Mode int
	Dir string // Папка для сохранений
}

type mutexErrors struct {
	mu sync.Mutex
	errEtls  map[string]error
}
func (c *mutexErrors) AddError(etlName string,err error) {
	c.mu.Lock()
	c.errEtls[etlName] = err
	c.mu.Unlock()
	return
}
func (c *mutexErrors) CheckErr(etlName string) bool {
	c.mu.Lock()
	 _, ok := c.errEtls[etlName]
	c.mu.Unlock()
	return ok
}
func (c *mutexErrors) GetErrors() (errors map[string]error) {
	c.mu.Lock()
	errors = c.errEtls
	c.mu.Unlock()
	return
}

type mutexStopper struct {
	mu sync.Mutex
	x  bool
}
func (c *mutexStopper) stopSignal()  {
	c.mu.Lock()
	c.x = true
	c.mu.Unlock()
}

func (p *ChLoger) Stop()  {
	p.stopPrepare.stopSignal()
	for i:=0; i <  p.GetCores(); i++ {
		<-p.stopX
	}
	fmt.Println(tp.Getime(), "| ** END ** | Logger |")
	return
}

func (p *ChLoger) SignalStoper(prepareExit <-chan bool)  {
	<-prepareExit
	p.Stop()
	p.GoodBy <- true
	fmt.Println("Отправили *GoodBy* Logger")
	return
}

type ChLoger struct {
	//logPath string
	Options *Config
	batchCnt *mutexCounter
	process *mutexRunner
	Errors *mutexErrors
	ChInLog chan [4]string //Приемник строк
	stopPrepare *mutexStopper //Сигнал, что пока закругляться
	stopX chan bool //Миньон завершила свою работу
	GoodBy chan bool
}

func NewChLoger(cfg *Config) *ChLoger {
	time := tp.Getime()
	time = strings.Replace(time, " ", "_", -1)
	time = strings.Replace(time, ":", "", -1)
	tp.FckText(fmt.Sprintf("Создание папки логов: %s",cfg.Dir),tp.CheckMkdir(cfg.Dir))
	merr := new(mutexErrors)
	merr.errEtls = make(map[string]error)
	return &ChLoger{Options:cfg,ChInLog:make(chan [4]string,100),stopPrepare:new(mutexStopper),process:new(mutexRunner),stopX:make(chan bool,1000),batchCnt:new(mutexCounter),GoodBy:make(chan bool),Errors:merr}
}

type mutexRunner struct {
	mu sync.Mutex
	x  int
}
func (c *mutexRunner) addRun()(i int) {
	c.mu.Lock()
	c.x++
	i = c.x
	c.mu.Unlock()
	return
}
func (c *mutexRunner) getCores()(cores int)   {
	c.mu.Lock()
	cores = c.x
	c.mu.Unlock()
	return 
}

func (c *ChLoger) GetCores()(cores int)  {
	cores = c.process.getCores()
	return
}

func (p *ChLoger) RunMinion()  {
	go p.runMinion(p.process.addRun())
	return
}


type mutexCounter struct {
	mu sync.Mutex
	x  int
}
func (c *mutexCounter) cntPlus()  {
	c.mu.Lock()
	c.x++
	c.mu.Unlock()
}



func (c *ChLoger) GetCountPackage()(cnt int)   {
	c.batchCnt.mu.Lock()
	cnt = c.batchCnt.x
	c.batchCnt.mu.Unlock()
	return
}
//ch <- chan [3][]string //Приемник строк
//[0] - Имя функции
//[1] - Имя Etl
//[2] - Текст
//[3] - Тип "error","ok"

func (p *ChLoger) runMinion(gopher int)  {
	// Загрузка будет проходить в папку с датой //
	var gg time.Duration = time.Duration(p.Options.IntervalMs)*time.Millisecond
	for {
		select {
		case elem := <-p.ChInLog:
			var etlFlush string
			var funcFlush string
		ln := len(elem[2])-1
		if elem[2][ln:][0] == 10 {
			elem[2] =elem[2][:ln]
		}

			//fmt.Println(elem[len(elem[2])-2:])
			// Посмотрим есть ли эта функция в правилах
			if val, ok := p.Options.ConsolFilterFn[elem[0]]; ok {
				// Если режим совпадает, то печатаем  или скрываем
				if p.Options.Mode != val {
					continue
				}
			}
			etlFlush = fmt.Sprintf("%s | FUNC: %s | TEXT: %s\n",tp.Getime(),elem[0],elem[2])
			funcFlush =fmt.Sprintf("%s | FUNC:%s | UNIT: %s | TEXT: %s\n",tp.Getime(),elem[0],elem[1],elem[2])
			if elem[1] != "nil" { // Если юнит не задан
				fileEtl := fmt.Sprintf("%s.txt",elem[1])
				fEtl, err := tp.OpenWriteFile(filepath.Join(p.Options.Dir,fileEtl))
				tp.FckText(fmt.Sprintf("Логгер | Открытие файла: %s",fileEtl),err)
				fEtl.Write([]byte(etlFlush))
				fEtl.Close()
			}
			if elem[3] == "1" { // Если ошибка
				errFlush := fmt.Sprintf("%s | UNIT: %s | FUNC: %s | TEXT: %s\n",tp.Getime(),elem[1],elem[0],elem[2])
				fileEtl := fmt.Sprintf("%s.txt","Errors")
				fEtl, err := tp.OpenWriteFile(filepath.Join(p.Options.Dir,fileEtl))
				tp.FckText(fmt.Sprintf("Логгер | Открытие файла: %s",fileEtl),err)
				fmt.Fprintf(os.Stderr, errFlush)
				fEtl.Write([]byte(errFlush))
				fEtl.Close()
			} else {
				fmt.Print(funcFlush)
			}
			fileFunc := fmt.Sprintf("%s.txt",elem[0])
			fFunc, err := tp.OpenWriteFile(filepath.Join(p.Options.Dir,fileFunc))
			tp.FckText(fmt.Sprintf("Логгер | Открытие файла: %s",fileFunc),err)
			fFunc.Write([]byte(funcFlush))
			fFunc.Close()

		default:
			if p.stopPrepare.x {
				if len(p.ChInLog) == 0 {
					p.stopX <- true
					return
				}
			}
			time.Sleep(gg)
		}
	}
}

