package main

import (
	"./.."
	"fmt"
)
func main() {
	var args Lego.Constructor
	var text string
	// Test 1
	args = Lego.Constructor{
		Values: map[string]string{
			"TestArg": `This arg1: "$1$" and This arg2: "$2$"`},
	}
	text = `Начало | #TestArg(argument 1,arg2)# | Конец`
	fmt.Printf("Test 1: %s\n", Lego.ScanText(text, args))

	// Test 2
	args = Lego.Constructor{
		MultyValues: map[string][]string{
			"test_code_array": []string{"1000001", "20000002",},
		},
	}
	text = `Начало |  num = (@#test_code_array#*,*@) | Конец`
	fmt.Printf("Test 2: %s\n", Lego.ScanText(text, args))

	// Test4
	args = Lego.Constructor{
		Values: map[string]string{
			"TestArg": `This code: "#budget_code#"`,
			"budget_code": "99010001",},
	}
	text = `Начало |  #TestArg# | Конец`
	fmt.Printf("Test 4: %s\n", Lego.ScanText(text, args))

}