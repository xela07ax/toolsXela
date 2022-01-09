package main

import (
	"../../tools"
	"../../types"
	"fmt"
)
func main() {
	_ = types.Lego{
		Text: "#DatesTemplate(closedate,opendate)#",
		Found: []types.FindVarTyp{
			types.FindVarTyp{
				L:0,
				R:35,
				FindValue:"DatesTemplate",
				Params: []string{
					"closedate", "opendate",
				},
			},
		},
	}
	elementArr2 := types.Lego{
		Text: "to+#Code_out#*$*",
		Found: []types.FindVarTyp{
			types.FindVarTyp{
				L:3,
				R:13,
				FindValue:"Code_out",
				SplitSymbol:"$",
			},
		},
	}


	args := types.Constructor{
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

	//text := "to+#Code_out#*$*"


	fmt.Println("output program: ", texts)

}


