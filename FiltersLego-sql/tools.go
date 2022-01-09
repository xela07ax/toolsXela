package Lego

import (
	"fmt"
	"strconv"
	"strings"
)

func findVariable (str string,symbols string ) (out []FindVarTyp)  {

	// Пример того что приходит: Задача найти элементы в  #...("$1","$2")# #Budget_level# и тут текст #переменнаа 2#.
	var lastIndex int
	for true {
		wStr := str[lastIndex:]
		var tmpStruct FindVarTyp
		//tmpStruct.Params = make(map[string]string)

		i := strings.Index(wStr, symbols) // Индекс начала
		if i == -1 { // Если ничего нет, то просто выходим
			break
		}
		tmpStruct.L = i+lastIndex
		// fmt.Println(i) // 12
		i2 := strings.Index(wStr[i+1:], symbols) // Индекс c найденного места
		if i2 == -1 { // Если ничего нет, то просто выходим
			break
		}
		//fmt.Println(str[i+1:]) // 12
		//fmt.Println(i2) // 37
		//fmt.Println(str[i+1:i2+i+1]) // Содержимое между двумя этими символами
		tmpStruct.R = i2+1+i+1+lastIndex
		tmpStruct.FindValue = wStr[i+1:i2+i+1]

		// Вытаскиваем аргументы

		//fmt.Println("FindValue: ",tmpStruct.FindValue)
		out = append(out,tmpStruct)
		// os.Exit(-2)
		lastIndex += i2+1+i+1 //Удаляем все до конца найденного последнего индекса, что бы не зациклиться
	}


	return
	// OUT: []string{"...","Budget_level","переменнаа 2"}
}

func findParams(text string) (newText string,params []string)  {
	//fmt.Println("findParams: ",text ) // Code_out
	// Вытаскиваем аргументы
	// Находим первую скобку
	iScoope := strings.Index(text, "(") // Получаем либо индекс, где заканчивается имя, либо -1
	//var funcVariables []string // Будем хранить переменные из шаблона, если есть
	if iScoope != -1 { // Если ничего нет, то просто выходим
		//Находим закрывающую скобку
		var text_startScoope string = text[iScoope+1:]
		iScoope_end := strings.Index(text_startScoope, ")")
		text_startScoope = text_startScoope[:iScoope_end]
		//fmt.Println("text_startScoope: ",text_startScoope) //out:text_startScoope:  apc_end_date,apc_bgn_date
		params = strings.Split(text_startScoope,",")
		prefix := text[:iScoope]
		sufix := text[strings.Index(text, ")")+1:]
		text = prefix + sufix
		//fmt.Println(prefix," | ",sufix)
		//out: [apc_end_date apc_bgn_date]
	}
	newText = text
	return
}

func findCount(text string, args Constructor) (cnt int)  {
	//fmt.Println("Входящая строка: ",text) // and budget_code  in (@#budget_code_array#*,*@)   --<< код бюджета
	// К нам пришла строка, в этой строке нам надо снять переменные и узнать сколько раз ее нужно повторять
	// Перед тем как искать переменные надо убрать переменные массивов, так как это отдеьно повторяюющийся элемент
	// Массивы не ищем!! Если нашли массив, то у него свой цикл
	// Удаляем массивы перед поиском и возвращаем после
	// Снимаем Дамп
	foundDump := findVariable(text,"@")
	//fmt.Println("foundDump: ",foundDump) //and budget_code  in (@#budget_code_array#*,*@)   --<< код бюджета 50
	for j, _ := range foundDump{
		foundDump[j].ReplaceFinal = fmt.Sprintf("@%s@",make([]byte, len(foundDump[j].FindValue)))

	}
	text = replace(text,foundDump)
	// Шаг3: Вытаскиваем # переменные
	variables := findVariable(text,"#")

	for i2, _ := range variables {
		// Шаг4: Нужно вытащить параметры ( в скобках
		variables[i2].FindValue,variables[i2].Params = findParams(variables[i2].FindValue)
	}

	// В variables у нас готовые для поиска переменные

	for _, v := range variables {
		if find, ok := args.Values[v.FindValue]; ok == true {
			cntFilters := findCount(find,args)
			//fmt.Println("cntFilters: ",cntFilters, " | Insert text: ", find) //0  | Insert text:   and budget_code  in (@#budget_code_array#*,*@)   --<< код бюджета
			if cntFilters > cnt {
				cnt = cntFilters
			}
		} else if find, ok := args.MultyValues[v.FindValue]; ok == true {
			//fmt.Println("MultyValues find")
			if len(find) > cnt {
				cnt = len(find)
			}
		} else if _, ok := args.Values[v.FindValue]; ok == true {
			//fmt.Println("one cnt: ", cnt, " | ",text)
			if cnt < 1 { // Если найденый аргумент не пуст, то присваеваем одно повторение
				cnt = 1
			}
		}
	}

	return
}

func clear(text string, l int,r int) string  {
	// Делаем зачистку
	return text[:l] + text[r:]
}
func replace(text string, variables []FindVarTyp) string  {
	// Поскольку индекс после изменения текста сдвигается в большую или меньшую сторону, сюда буду складывать разницу для начальной точки отсчета
	var plus int
	for _, va := range variables{
		//Итерируемся по переменнымв шаблоне: {15 26 StartDate  []}
		va.L += plus
		va.R += plus
		// Делаем замены
		text = text[:va.L] + va.ReplaceFinal + text[va.R:]

		// Считаем размер байтов переменной в тексте шаблона
		sumIndexVar := va.R - va.L
		// Считаем размер байтов текста на который меняем в тексте шаблона
		sumIndexText := len(va.ReplaceFinal)
		// Подсчитываем разницу и добавляем в сумку
		plus += sumIndexText - sumIndexVar

	}
	return text
}

func findSplitSymbol(text string) (newText string, splitSymbol string)  {
	// in: "#budget_code_array#*,*@"  // out: "#budget_code_array#" ","
	// Заходим разделитель для массивов, разделитель заключен в снежинки *
	// Для того что бы найти разделитель массивов, мы сначала используем универсальную функцию
	// которая найдет символ, а потом функцию которая зачистит в тексте это место
	elementArr := findVariable(text, "*")
	if len(elementArr) > 0 {
		splitSymbol = elementArr[0].FindValue
		newText = clear(text,elementArr[0].L,elementArr[0].R)
	} else {
		newText = text
	}

	return
}


func ScanText(text string, args Constructor) (newText []string) {
	//fmt.Println("ScanText in: ",text)
	arrays := findVariable(text, "@")
	//fmt.Println("foundDump: ",foundDump) //and budget_code  in (@#budget_code_array#*,*@)   --<< код бюджета 50
	if len(arrays) > 0 {
		// Вычисляем вложенные массивы (рекурсивная обработка)
		for j, _ := range arrays {
			arrays[j].FindValue , arrays[j].SplitSymbol = findSplitSymbol(arrays[j].FindValue)
			arrays[j].FindValue, arrays[j].Params = findParams(arrays[j].FindValue)
			//fmt.Println("arrays[j]: ",arrays[j]) //{22 46 #budget_code_array#  [] ,}
			scanned := ScanText(arrays[j].FindValue,args)
			textSplit := strings.Join(scanned,arrays[j].SplitSymbol) // Просто сливаем

			//fmt.Println("~~~~~~~>scanned: ",textSplit)
			// Шаг3: Вытаскиваем $ параметры

			//fmt.Println("textSplitreplace: ", textSplitreplace, " | variablesSplit: ", variablesSplit)


			arrays[j].ReplaceFinal = textSplit
		}
		// Применяем изменения
		text = replace(text,arrays)
	}


	cnt := findCount(text,args)
	//fmt.Println("cnt out : ", cnt)
	if cnt == 0 {
		cnt = 1
	}
	for i:=0;i<cnt;i++ {
		exptract := extractValue(text,args,i)
		newText = append(newText, exptract)
			//fmt.Println("Exptracted",exptract, )
	}

	return
}

func extractValue(text string, args Constructor, index int) string  {
	//fmt.Println("Func extractValue IN: ",text, " | INdex: ", index)
	// Шаг3: Вытаскиваем # переменные
	variables := findVariable(text, "#")
	for i2, _ := range variables {
		// Шаг4: Нужно вытащить параметры ( в скобках
		variables[i2].FindValue, variables[i2].Params = findParams(variables[i2].FindValue)
		//fmt.Println("variables[i2].Params = ",variables[i2].Params, " | FindValue: ",variables[i2].FindValue)
	}

	for i3 , v := range variables {
		if find, ok := args.Values[v.FindValue]; ok == true {
			//cnt := findCount(find, args)
			//fmt.Println("cnt EXRRACT out : ", cnt, " | SEND TO SCANTEXT: ", find)
			finds := ScanText(find, args)
			variables[i3].ReplaceFinal = strings.Join(finds,"") // Просто сливаем

			variablesSplit := findVariable(variables[i3].ReplaceFinal, "$")
			//fmt.Println(variablesSplit)
			for i4, _ := range variablesSplit {
				num, err := strconv.Atoi(variablesSplit[i4].FindValue)
				if err != nil {
					continue // Если не получилось преобразовать индекс параметра, то просто пропускаем
				}
				if len(v.Params) >= num {
					variablesSplit[i4].ReplaceFinal = v.Params[num-1]
				}
			}
			variables[i3].ReplaceFinal = replace(variables[i3].ReplaceFinal, variablesSplit)
			//fmt.Println("new Find variables[i3].ReplaceFinal: ", variables[i3].ReplaceFinal, " | v: ",v)
			//variables[i3].ReplaceFinal
		} else if find, ok := args.MultyValues[v.FindValue]; ok == true {
			//fmt.Println("ExtractValue Multivalue detected: ", find)
			if len(find) > index {
				variables[i3].ReplaceFinal = find[index]
			}
		} else if find, ok := args.Values[v.FindValue]; ok == true {
			variables[i3].ReplaceFinal = find
		}
	}
	//fmt.Println("---->text: ",text, " | Variables: ", variables)
	text = replace(text,variables)
	//fmt.Println("--2-->text: ",text)
	///if len(variables) > 0 {
	///
	///}

	//	fmt.Println("Func extractValue OUT: ",text, " | INdex: ", index)
	return text
}
