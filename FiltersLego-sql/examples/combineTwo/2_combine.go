package main

import (
	"../../tools"
	"../../types"
	"fmt"
)

func _main()  {
	elementArr := types.Lego{
		Text: "#DatesTemplate(closedate,opendate)#",
		Found: []types.FindVarTyp{
			types.FindVarTyp{
				L:0,
				R:35,
				FindValue:"DatesTemplate(closedate,opendate)",
			},
		},
	}
	text, params := tools.FindParams(elementArr.Found[0].FindValue)
	fmt.Println("text, params: ", text, params)


	text = "DatesTemplate"
	params = []string{
		"closedate", "opendate",
	}

	//Меняем FindValue:
	//Добавляем Params:
	elementArr.Found[0].FindValue = "DatesTemplate"
	elementArr.Found[0].Params = []string{
		"closedate", "opendate",
	}


	elementArr = types.Lego{
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

}

