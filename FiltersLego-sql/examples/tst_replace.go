package main

import (
	"../tools"
	"../types"
	"fmt"
)
func main()  {
	text := "to+#Code_out# and #GG#"
	elementArr2 := []types.FindVarTyp{
		types.FindVarTyp{
			L:3,
			R:13,
			FindValue:"Code_out",
			ReplaceFinal:"1246-erregrtghf-hgfhfghfg-6666",
		},
		types.FindVarTyp{
			L:18,
			R:22,
			FindValue:"GG",
			ReplaceFinal:"0000000-e-00000",
		},
	}

	all := tools.Replace(text,elementArr2)
	fmt.Println("Out text: ",all) // all:  to+1246-erregrtghf-hgfhfghfg-6666 and 0000000-e-00000
}
