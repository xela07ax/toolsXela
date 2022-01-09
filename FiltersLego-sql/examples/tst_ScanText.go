package main

import (
	"./.."
	"fmt"
)
func main()  {
	args := Lego.Constructor{
		Filters: map[string]string{
			"DatesTemplate":     "AND $1$ >= date'#StartDate#' --<< дата с\n        AND $2$ <= date'#EndDate#' --<< дата по",
			"Level":             "and budget_level = nvl('#budget_level#',budget_level)  --<< уровень бюджета",
			"Code":              " and budget_code  = nvl('#budget_code#',budget_code)   --<< код бюджета",
			"ArrParamsTemplate": "AND $1$ >= in'#Code_in#' --<< дата с\n        AND $2$ <= out'#Code_out#' --<< дата по",
			"Code_out":          " and budget_code  in (@#budget_code_array#*,*@)   --<< код бюджета"},
		Values: map[string]string{
			"budget_code":  "99010001",
			"budget_level": "",
			"StartDate":    "16.07.2018",
			"EndDate":      "20.08.2018",
		},
		MultyValues: map[string][]string{
			"budget_code_array": []string{"1000001", "20000002",},
			"Code_in":           []string{"99999-9991", "77777-7773"},
			"Array_Null":        []string{},
		},
	}
	// text := `with acc as #DatesTemplate(closedate,opendate)#  group by tofk_code #budget_code#;`
	 //text := "select * FROM t WHERE 1=1 #Code_out# limit 100"
	// text := "select * FROM t WHERE 1=1 @#DatesTemplate(closedate,opendate)#@ @#DatesTemplate(col2,col3)#@ limit 1000;"
	 text := ` @#DatesTemplate(closedate,opendate)#@, @ wau#Array_Null#@with @ to+#Code_out#*$*@ acc as #DatesTemplate(closedate,opendate)#
    @ and #Code_out# = 6 and #Code_in# = 8 @
  @ and #Code_null# = not visible @
    #Level# -- <-- Этого фильтра не должно быть!!!
          #Code_in#
  select servicetofkcode as tofk_code
             , count(distinct accnum) as acc_count
	#DatesTemplate(apc_end_date,apc_bgn_date)#
          #Level# -- <-- Этого фильтра не должно быть!!!
          #Code#
		      #Code_in#
        group by tofk_code;`
	//fmt.Println("Insert text: ", text)
	all := Lego.ScanText(text,args)
	fmt.Println("all: ",all) // all:  to+1246-erregrtghf-hgfhfghfg-6666 and 0000000-e-00000
}
