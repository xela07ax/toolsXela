package Lego

type Constructor struct {
	StringTemplate string `json:"STRINGTEMPLATE"`
	Values map[string]string `json:"VALUES"`
	MultyValues map[string][]string `json:"MULTYVALUES"`
}

type FindVarTyp struct {
	L int
	R int
	FindValue string // Имя переменной в шаблоне
	ReplaceFinal string // Значение, на которое меняем переменную
	Params []string // Параметры к имени переменной в шаблоне
	SplitSymbol string // Для массивов! Разделитель
}