package cmdConveer

import "time"

type Command struct {
	Name string
	Executer string
	DirHome string
	Args []string
	Timeout time.Duration
	Cfg map[string]string
}

// inCh0[0] = ExecCommand
// inCh0[1] = Домашняя директория
// inCh0[2] = Команды

type Reject struct {
	Name string
	Status int
	MyStdout string
	MyStderr string
	Comment string
	Command Command
}

// outCh0[0] = Статус
// "0" = error
// "1" = ok
// outCh0[1] = Комментарий
// Текст ошибки
// outCh0[2] = Текст