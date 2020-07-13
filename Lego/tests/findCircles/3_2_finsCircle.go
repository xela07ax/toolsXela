package main

import (
	"../../tools"
	"../../types"
	"fmt"
)

func main()  {
	args := types.Constructor{
		SqlTemplate:"@#DatesTemplate(closedate,opendate)#@, @ wau#Array_Null#@with @ to+#Code_out#*$*@ acc as #DatesTemplate(closedate,opendate)#\n    @ and #Code_out# = 6 and #Code_in# = 8 @\n  @ and #Code_null# = not visible @\n    #Level# -- <-- Этого фильтра не должно быть!!!\n          #Code_in#\n  select servicetofkcode as tofk_code\n             , count(distinct accnum) as acc_count\n\t#DatesTemplate(apc_end_date,apc_bgn_date)#\n          #Level# -- <-- Этого фильтра не должно быть!!!\n          #Code#\n\t\t      #Code_in#\n        group by tofk_code;",
		Filters: map[string]string{
			"DatesTemplate": "AND $1$ >= date'#StartDate#' --<< дата с\n        AND $2$ <= date'#EndDate#' --<< дата по",
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

	scannd3 := tools.ScanText(args.SqlTemplate,args)
	fmt.Println("scanned3: ", scannd3)

}