package main

import (
	"./.."
	"encoding/json"
	"fmt"
	"log"
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
	Версия 1.3
		Другой формат входных данных
	 */

	//----------------------Часть 1 Прием значений
	/* Ожидаем структуру в формате Json
	SQLTEMPLATE - Шаблон SQL запроса
	FILTERS -
	VALUES -
	MULTYVALUES -

	*/
	/*
		Руководство использования шаблонов в теле фильтра:
		Слова заключенные в символ решетки (#слово#) это аргументы, значение которых будет искаться в структуре HTTP ответа, имя должно совпатать с Tag-ом в "json:" префиксе поля в структуре
		Слова заключенные в символ решетки имеет возможность задать аргументы через запятую(#слово(аргумент1,аргумент2)#)
		Слова заключенные в символ доллара ($слово$) это аргументы, затается в скобках и в кавычках рядом с именем переменной имени фильтра.
		Слова заключенные в символ доллара (@слово@) это итерация по массиву, внутри имя массива и в звездочках *,* разделитель (в примере запятая)
	*/



	httpResponse := `{
  "STRINGTEMPLATE": "with @ acc as #DatesTemplate(closedate,opendate)#\n    @ and #Code_out# = 6 and #Code_in# = 8 @\n  @ and #Code_null# = not visible @\n    #Level# -- <-- Этого фильтра не должно быть!!!\n          #Code_in#\n  select servicetofkcode as tofk_code\n             , count(distinct accnum) as acc_count\n\t#DatesTemplate(apc_end_date,apc_bgn_date)#\n          #Level# -- <-- Этого фильтра не должно быть!!!\n          #Code#\n\t\t      #Code_in#\n        group by tofk_code;",
  "FILTERS": {
    "DatesTemplate": "AND $1$ >= date'#StartDate#' --<< дата с\n        AND $2$ <= date'#EndDate#' --<< дата по",
    "Level": "and budget_level = nvl('#bud get_level#',budget_level)  --<< уровень бюджета",
    "Code": " and budget_code  = nvl('#budget_code#',budget_code)   --<< код бюджета",
    "ArrParamsTemplate": "AND $1$ >= in'#Code_in#' --<< дата с\n        AND $2$ <= out'#Code_out#' --<< дата по",
    "Code_in": " and budget_code  in (@#budget_code_array#*,*@)   --<< код бюджета",
    "Code_out": " and budget_code  in (@#budget_code_array#*,*@)   --<< код бюджета"
  },
  "VALUES": {
    "budget_code":
      "99010001"
    ,
    "budget_level": 
      "1"
    ,
    "StartDate": 
      "16.07.2018"
    ,
    "EndDate": 
      "20.08.2018"
  },
  "MULTYVALUES": {
    "Code_in": [
      "1000001",
      "20000002"
    ],
   "Code_out": [
      "99999-9991",
      "77777-7773"
    ],
   "Array_Null": [
    ]
  }
}`

fmt.Println(httpResponse)
	// Пройдем по списку и поищем полученные данные, если есть мы подставим, если нет удалим из списка актуальных фильтров

	var httpData Lego.Constructor
	err := json.Unmarshal([]byte(httpResponse), &httpData)
	if err != nil {
		log.Fatalf("Error Unmarshal body: %v", err)
	}
	//httpResponse := InterfaceToMap(httpData)
	//fmt.Println("httpResponse: ",httpResponse)

	//----------------------Часть 2 Парсинг SQL запроса

	variables := Lego.ScanText(httpData.StringTemplate,httpData)

	//variables := FindVariable(httpData.SqlTemplate,"@",0)
	fmt.Printf("SQL Variables: %v\n", variables) //out: {with acc as #DatesTemplate(closedate,opendate)# ...[{12 47 DatesTemplate  map[1:closedate 2:opendate]} {58 65 Level  map[]} {139 148 Code_in  map[]} {239 281 DatesTemplate  map[1:apc_end_date 2:apc_bgn_date]} {292 299 Level  map[]} {373 379 Code  map[]} {388 397 Code_in  map[]}]
	//SQL Variables: [{5 24  to+#Code_out#*$*  map[]} {72 112  and #Code_out# = 6 and #Code_in# = 8   map[]}]



	//----------------------Часть 3 Парсинг фильтров

	// Проходим по пропарсиным в теле запроса переменным с параметрами и собираем значение для замены
	// Переменные сначала пытаются найти замену в том что у нее в параметрах, остальное из Http
}


