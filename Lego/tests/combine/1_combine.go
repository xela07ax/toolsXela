package main

import (
	"../../tools"
	"../../types"
	"fmt"
)


func main() {

	onePart := types.Lego{
		Symbol: "@",
		Text: `@#DatesTemplate(closedate,opendate)#@, @ wau#Array_Null#@with @ to+#Code_out#*$*@ acc as #DatesTemplate(closedate,opendate)#
	@ and #Code_out# = 6 and #Code_in# = 8 @
	#Level# -- <-- Этого фильтра не должно быть!!!
	#Code_in#
	select servicetofkcode as tofk_code
		, count(distinct accnum) as acc_count
		#DatesTemplate(apc_end_date,apc_bgn_date)#
		#Level# -- <-- Этого фильтра не должно быть!!!
		#Code#
		#Code_in#
		group by tofk_code;`,
		Found: []types.FindVarTyp{
			{
				L: 0, R: 37, FindValue: "#DatesTemplate(closedate,opendate)#",
			},
			types.FindVarTyp{
				L: 39, R: 57, FindValue: "wau#Array_Null#",
			},
			types.FindVarTyp{
				L:         62,
				R:         81,
				FindValue: "to+#Code_out#*$*",
			},
			types.FindVarTyp{
				L:         129,
				R:         169,
				FindValue: " and #Code_out# = 6 and #Code_in# = 8 ",
			},
			types.FindVarTyp{
				L:         172,
				R:         205,
				FindValue: "and #Code_null# = not visible ",
			},
		},
	}
	//Для того что бы найти разделитель массивов, мы сначала используем универсальную функцию
	// которая найдет символ, а потом функцию которая зачистит в тексте это место
	onePart.Found[2].FindValue , onePart.Found[2].SplitSymbol = tools.FindSplitSymbol(onePart.Found[2].FindValue)
	onePart.Found[2] = types.FindVarTyp {
					L: 62,
					R: 81,
					FindValue: "to+#Code_out#",
					SplitSymbol:"$",
				}
	fmt.Println(onePart.Found[2])


	args := types.Constructor{
		Filters: map[string]string{"DatesTemplate": "AND $1$ >= date'#StartDate#' --<< дата с\n        AND $2$ <= date'#EndDate#' --<< дата по",
			"Level":             "and budget_level = nvl('#bud get_level#',budget_level)  --<< уровень бюджета",
			"Code":              " and budget_code  = nvl('#budget_code#',budget_code)   --<< код бюджета",
			"ArrParamsTemplate": "AND $1$ >= in'#Code_in#' --<< дата с\n        AND $2$ <= out'#Code_out#' --<< дата по",
			"Code_out":          " and budget_code  in (@#budget_code_array#*,*@)   --<< код бюджета"},
		Values: map[string]string{
			"budget_code":  "99010001",
			"budget_level": "1",
			"StartDate":    "16.07.2018",
			"EndDate":      "20.08.2018",
		},
		MultyValues: map[string][]string{
			"budget_code_array":    []string{"1000001", "20000002",},
			"Code_in":   []string{"99999-9991", "77777-7773"},
			"Array_Null": []string{},
		},
	}

	fmt.Println(args)

}