package constructor

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func main()  {
	/*
		Версия 1.2
		План задач:
		 Новый функционал:
		  Возможность добавлять элементы в фильтр "in" или дублировать шаблон фильтра с новыми значениями целиком
		 Исправление:
		  Если приходит пустое значение, то фильтр все равно добавляется, но значение аргумента в нем пусто.
		  Фильтр должен исчезнуть, если значение для него пусто.
	*/

	//----------------------Часть 1 Прием значений
	type GlobalFilters struct {
		Code_in []string `json:"budget_code_array"`
		Code string `json:"budget_code"`
		Degree string `json:"budget_level"`
		StartDate string `json:"StartDate"`
		EndDate string `json:"EndDate"`
	}
	httpData := GlobalFilters{ // не заполнено поле "Degree", значит этот фильтр должен исчезнуть
		Code: "99010001",
		Code_in: []string{"1000001","20000002"},
		StartDate: "16.01.2018",
		EndDate: "18.07.2019",
	}
	// Пройдем по списку и поищем полученные данные, если есть мы подставим, если нет удалим из списка актуальных фильтров

	//httpResponse := InterfaceToMap(httpData)
	//fmt.Println("httpResponse: ",httpResponse)

	//----------------------Часть 2 Парсинг SQL запроса
	SQLtemplate := `with acc as #DatesTemplate(closedate,opendate)#
          #Level# -- <-- Этого фильтра не должно быть!!!
          #Code_in#
  select servicetofkcode        as tofk_code
             , count(distinct accnum) as acc_count

	#DatesTemplate(apc_end_date,apc_bgn_date)#
          #Level# -- <-- Этого фильтра не должно быть!!!
          #Code#
		  #Code_in#
        group by tofk_code;`

	variables := FindVariable(SQLtemplate,"#")

	// fmt.Println("SQL Variables: ", variables)

	//----------------------Часть 3 Парсинг фильтров

	// Проходим по пропарсиным в теле запроса переменным с параметрами и собираем значение для замены
	// Переменные сначала пытаются найти замену в том что у нее в параметрах, остальное из Http

	/*
		Руководство использования шаблонов в теле фильтра:
		Слова заключенные в символ решетки (#слово#) это аргументы, значение которых будет искаться в структуре HTTP ответа, имя должно совпатать с Tag-ом в "json:" префиксе поля в структуре
		Слова заключенные в символ решетки имеет возможность задать аргументы через запятую(#слово(аргумент1,аргумент2)#)
		Слова заключенные в символ доллара ($слово$) это аргументы, затается в скобках и в кавычках рядом с именем переменной имени фильтра.
		Слова заключенные в символ доллара (@слово@) это итерация по массиву, внутри имя массива и в звездочках *,* разделитель (в примере запятая)
	*/
	filters := map[string]string{ "DatesTemplate" : `AND $1$ >= date'#StartDate#' --<< дата с
        AND $2$ <= date'#EndDate#' --<< дата по`,
		"Level" : `and budget_level = nvl('#budget_level#',budget_level)  --<< уровень бюджета`,
		"Code" : ` and budget_code  = nvl('#budget_code#',budget_code)   --<< код бюджета`,
		"DatesTemplate_" : `AND $1$ >= date'#StartDate#' --<< дата с
        AND $2$ <= date'#EndDate#' --<< дата по`,
		"Code_in" : ` and budget_code  in (@#budget_code_array#*,*@)   --<< код бюджета`,
	}

	_ = Replace(SQLtemplate,variables)

	text := RunParseTemplate(SQLtemplate,filters,httpData)
	fmt.Println("Replay: ",text)
	/* out:
	Replay:  with acc as AND closedate >= date'16.01.2018' --<< дата с
	        AND opendate <= date'18.07.2019' --<< дата по
	          and budget_level = nvl('',budget_level)  --<< уровень бюджета
	           and budget_code  = nvl('99010001',budget_code)   --<< код бюджета
	  select servicetofkcode        as tofk_code
	             , count(distinct accnum) as acc_count

		AND apc_end_date >= date'16.01.2018' --<< дата с
	        AND apc_bgn_date <= date'18.07.2019' --<< дата по
	          and budget_level = nvl('',budget_level)  --<< уровень бюджета
	           and budget_code  = nvl('99010001',budget_code)   --<< код бюджета
	        group by tofk_code;
	*/
}

func RunParseTemplate(template string, filters map[string]string, httpData interface{}) string  {
	//-Часть 1 - входные данные
	// Пройдем по списку и поищем полученные данные, если есть мы подставим, если нет удалим из списка актуальных фильтров
	// В Sql запросе (template) находится имя фильтра которое должно присутствовать с справочнике шаблонов фильтров (filters)
	// К нам приходит справочник значений в справочник Http сообщения (httpData)
	// В нем содержится имя, и массив значений.??  Нужно ли оно?, да. Каждое значение должно дублировать шаблон.
	// Складывать готовые SQL фильтры будем в variables
	httpResponse := InterfaceToMap(httpData)
	// fmt.Println("ai")
	// fmt.Println("httpResponse: ",httpResponse) //out: httpResponse:  map[EndDate:18.07.2019 StartDate:16.01.2018 budget_code:99010001 budget_level:]
	// fmt.Println("ai")
	//-Часть2 - текст шаблона

	variables := FindVariable(template,"#")
	//fmt.Println("variables: ", variables)
	/* out: variables:  [{12 51 DatesTemplate  map[1:closedate 2:opendate]} {62 69 Level  map[]} {143 149 Code  map[]}
	{248 294 DatesTemplate  map[1:apc_end_date 2:apc_bgn_date]} {305 312 Level  map[]} {386 392 Code  map[]}]
	*/
	//Разобраный из шаблона: {12 51 DatesTemplate false  [closedate opendate]} следующая итерация {62 76 Budget_level false  []}
	for ione, textFiltersProcess := range variables { // Переменные в шаблоне
		if filterTemplate, ok := filters[textFiltersProcess.FindValue]; ok == true { // Вытаскиваем к ней фильтр
			//Извлекаем шаблон фильтра: "Budget_level" : and budget_level = nvl('#Budget_level#',budget_level)  --<< уровень бюджета
			// Ищем все переменные в теле фильтра
			// Просто пройдем по каждой переменные и попробуем найти чем заменить
			// Ищем значение для переменных, для последующей замены
			//fmt.Println("filterTemplate: ",filterTemplate)
			/* out: filterTemplate:  AND $1$ >= date'#StartDate#' --<< дата с
			   AND $2$ <= date'#EndDate#' --<< дата по
			*/
			/*
				filters := map[string]string{ "DatesTemplate" : `AND $1$ >= date'#StartDate#' --<< дата с
				        AND $2$ <= date'#EndDate#' --<< дата по`,
						"Level" : `and budget_level = nvl('#budget_level#',budget_level)  --<< уровень бюджета`,
						"Code" : ` and budget_code  = nvl('#budget_code#',budget_code)   --<< код бюджета`,
						"DatesTemplate_" : `AND $1$ >= date'#StartDate#' --<< дата с
				        AND $2$ <= date'#EndDate#' --<< дата по`,
						"Code_in" : ` and budget_code  in (@#budget_code_array#*","*@)   --<< код бюджета`, <--- Вытащили фильтр
					}

			*/

			var text string
			var allVariables []FindVarTyp

			// Сначала вытащим массив
			//fmt.Println("filterTemplate_oldArr: ",filterTemplate) //filterTemplate_oldArr:   and budget_code  in (@#budget_code_array#*,*@)   --<< код бюджета
			allArrays := FindVariable(filterTemplate,"@")

			// fmt.Println("@",allVariables) // out: @ [{22 48 #budget_code_array#*","*  map[]}]
			// Если нашли маасив, то вытащим разделитель
			for iarr , varr := range allArrays {
				allSplits := FindVariable(varr.FindValue,"*")
				//fmt.Println("allSplits",allSplits) //[{19 24 ","  map[]}]
				// заполняем разделитель
				var splitVal string // Будем хранить разделитель
				for _ , splitVar := range allSplits {
					splitVal = splitVar.FindValue
				}
				//fmt.Println("varr.FindValue_old: ",varr.FindValue) //out: varr.FindValue_old:  #budget_code_array#*,*
				//очистим лишнее
				varr.FindValue = Replace(varr.FindValue,allSplits)
				// fmt.Println("varr.FindValue_old_clean: ",varr.FindValue) //out: varr.FindValue_old_clean:  #budget_code_array#
				// Пройдем по значениям и сошьем массивы в строку
				new_httpResponse := Stapler(httpResponse,splitVal)

				// разделитель заполнен или остался пуст
				// вытаскиваем значения для пременных
				// сначала вытащим сами переменные
				allVariablesFromArr := FindVariable(varr.FindValue,"#")
				// fmt.Println("iarr: ",iarr, " | SplitVal: ",SplitVal," | allVariablesFromArr: ",allVariablesFromArr)
				// out: iarr:  0  | SplitVal:  ","  | allVariablesFromArr:  [{0 19 budget_code_array  map[]}]
				//fmt.Println("allVariables_array_old: ", allVariablesFromArr[0].FindValue) //out: allVariables_array_old:  budget_code_array
				allVariables = CircleReplace(allVariablesFromArr, new_httpResponse)
				//fmt.Println("allVariables_array: ", allVariables[0].ReplaceFinal) // out: allVariables_array:  1000001,20000002
				varr.FindValue = Replace(varr.FindValue,allVariables)
				// fmt.Println("varr.FindValue: ",varr.FindValue) // out:varr.FindValue:  1000001,20000002
				allArrays[iarr].ReplaceFinal = varr.FindValue
				filterTemplate = Replace(filterTemplate,allArrays)
			}

			//fmt.Println("allArrays_final[0].FindValue: ",allArrays[0].FindValue) //out: allArrays_final[0].FindValue:  #budget_code_array#*,*
			//fmt.Println("allArrays_final[0].ReplaceFinal: ",allArrays[0].ReplaceFinal) //out: allArrays_final[0].ReplaceFinal:  1000001,20000002
			//fmt.Println("filterTemplate_NewArr: ",filterTemplate) //out:filterTemplate_NewArr:   and budget_code  in (1000001,20000002)   --<< код бюджета
			//
			//fmt.Println("filterTemplate_noARrr: ",allArrays)

			// CircleReplace - Принимает map[string]string, конвертируем
			primitiveHttpReply := Stapler(httpResponse,"")
			allVariables = CircleReplace(allVariables, primitiveHttpReply)

			allVariables = FindVariable(filterTemplate,"#") //Разбираем переменные в этом фильтре: [{15 26 StartDate  []} {68 77 EndDate  []}] или [{24 38 Budget_level  []}]
			allVariables = CircleReplace(allVariables, primitiveHttpReply)
			//fmt.Println("allVariables",allVariables)
			//out: allVariables [{16 27 StartDate 16.01.2018 map[]} {70 79 EndDate 18.07.2019 map[]}]

			// Убираем пустой фильтр
			// Можно пройти по массиву и если в качестве замены переменной пусто, то пропускаем всесь фильтр
			for _, v := range allVariables {
				if v.ReplaceFinal == "" {
					//fmt.Println("GOTO: ", textFiltersProcess)
					goto NULLFILTER
				}
			}


			text = Replace(filterTemplate,allVariables)

			//fmt.Println("tempStr_not_replace_params: ", text)
			/* out: tempStr_not_replace_params: AND $1$ >= date'16.01.2018' --<< дата с
			AND $2$ <= date'18.07.2019' --<< дата по
			*/
			//fmt.Println("Params: ",textFiltersProcess.Params) //Params:  [closedate opendate]

			// Ищем все параметры в теле фильтра
			allParams := FindVariable(text,"$") //[{4 7 1 false  []} {58 61 2 false  []}]
			allParams = CircleReplace(allParams, textFiltersProcess.Params)
			text = Replace(text,allParams)
			//fmt.Println("allParams- ",allParams) //out: allParams-  [{4 7 1 closedate map[]} {57 60 2 opendate map[]}]
			//fmt.Println("text- ",text)
			/* out: tempStr_not_replace_params: text-  AND closedate >= date'16.01.2018' --<< дата с
			   AND opendate <= date'18.07.2019' --<< дата по
			*/
			variables[ione].ReplaceFinal = text
		}
	NULLFILTER:
	}
	return Replace(template,variables)
}

// от сюда >#Budget_level("переменная","$2")#< до сюда
type FindVarTyp struct {
	L int
	R int
	FindValue string // Имя переменной в шаблоне
	ReplaceFinal string // Значение, на которое меняем переменную
	Params map[string]string // Параметры к имени переменной в шаблоне
}

func CircleReplace(variables []FindVarTyp, values map[string]string)  []FindVarTyp {
	/* fmt.Println("Circle: ",text)
	Circle:  AND $1$ >= date'#StartDate#' --<< дата с
	        AND $2$ <= date'#EndDate#' --<< дата по
	*/

	/*
		Описание функции
		На вход приходит текст шаблона (Sql запроса), информация обо всех переменных в тексте, значения для переменных.

		Итерируемся по каждой переменной
		Ищем значение, записываем в .ReplaceFinal
		Значит именно в .ReplaceFinal должно быть либо пусто (если аргумент пришел другой)
	*/

	//var findReplacement int // Сумка. Проверяем количество замен, должно совпадать с количеством переменных в шаблоне фильтра
	for iFi, _ := range variables { //Итерируемся по переменнымв шаблоне: {15 26 StartDate  []}
		// Ищем в Http
		if valFromHttp, ok := values[variables[iFi].FindValue]; ok == true {// out: "StartDate" -> valHttp ==> "[16.01.2018]"
			// Если элементов больше чем один, пока ничего делать не буду :), потом разберемся
			variables[iFi].ReplaceFinal = valFromHttp
			//findReplacement++
		}
	}

	return variables
}

func Replace(text string, variables []FindVarTyp) string  {
	// Поскольку индекс после изменения текста сдвигается в большую или меньшую сторону, сюда буду складывать разницу для начальной точки отсчета
	var plus int
	for _, va := range variables { //Итерируемся по переменнымв шаблоне: {15 26 StartDate  []}
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

func FindVariable(str string,symbols string) (out []FindVarTyp)  {
	// Пример того что приходит: Задача найти элементы в  #...("$1","$2")# #Budget_level# и тут текст #переменнаа 2#.
	var lastIndex int
	for true {
		wStr := str[lastIndex:]
		var tmpStruct FindVarTyp
		tmpStruct.Params = make(map[string]string)

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
		// Находим первую скобку
		iScoope := strings.Index(tmpStruct.FindValue, "(") // Получаем либо индекс, где заканчивается имя, либо -1
		//var funcVariables []string // Будем хранить переменные из шаблона, если есть
		if iScoope != -1 { // Если ничего нет, то просто выходим
			//Находим закрывающую скобку
			var text_startScoope string = tmpStruct.FindValue[iScoope+1:]
			iScoope_end := strings.Index(text_startScoope, ")")
			text_startScoope = text_startScoope[:iScoope_end]
			//fmt.Println("text_startScoope: ",text_startScoope) //out:text_startScoope:  apc_end_date,apc_bgn_date
			for imap, vFnd := range strings.Split(text_startScoope,",") {
				tmpStruct.Params[strconv.Itoa(imap+1)] = vFnd
			}
			tmpStruct.FindValue = tmpStruct.FindValue[:iScoope]
			//out: [apc_end_date apc_bgn_date]
		}

		out = append(out,tmpStruct)
		// os.Exit(-2)
		lastIndex += i2+1+i+1 //Удаляем все до конца найденного последнего индекса, что бы не зациклиться
	}
	return
	// OUT: []string{"...","Budget_level","переменнаа 2"}
}

type MultiValue struct {
	Type int //0-пусто, 1-string,2-slice
	ValString string
	ValSlice []string
}

func InterfaceToMap(data interface{})map[string]MultiValue  {
	// Преобразуем любую входящую структуру в мапу
	fType := reflect.TypeOf(data) // fType.NumField():  4
	fValue := reflect.ValueOf(data) //fValue.Field(0):  99010001

	out := make(map[string]MultiValue)


	//fmt.Println(fValue.Field(0).Kind())
	for i := 0; i < fType.NumField(); i++ {
		var bag MultiValue
		field := fType.Field(i)
		if field.Tag != "" {
			fieldName := field.Tag.Get("json") // out: budget_code

			// fmt.Println(fValue.Field(0).Kind()) //out: string
			if fValue.Field(i).Kind() == reflect.String { // out: string
				bag.Type = 1
				bag.ValString = fValue.Field(i).String()
			} else if fValue.Field(i).Kind() == reflect.Slice{ // Что бы код был более универсален
				bag.Type = 2
				for iv := 0 ; iv < fValue.Field(i).Len() ; iv++ {
					bag.ValSlice = append(bag.ValSlice,fValue.Field(i).Index(iv).String())// out : [18.07.2019] ... [18.07.2019]
				}
			}
			out[fieldName] = bag
		}

	}
	return out
}

func Stapler(params map[string]MultiValue, splitSymbol string) map[string]string {
	// Конвертор сложной структуры в примитивную
	// Сшиватель массивов в строку
	newParams := make(map[string]string)
	for key, value := range params {
		if value.Type == 1 {
			newParams[key] = value.ValString
		} else if value.Type == 2 {
			// Надо проверить индекс который передали, вдруг он больше того, сколько элементов
			newParams[key] = strings.Join(value.ValSlice[:], splitSymbol)
		}
		//fmt.Println("Key:", key, "Value:", value)
	}

	return newParams
}