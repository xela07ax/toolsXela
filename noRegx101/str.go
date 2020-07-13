package noRegx101

import (
	"strings"
)

//Нахождение совпадения без использования регулярок
func FindStrinTemplateNoRx(fullText string, findText string)bool  {
	// True - Найдено
	// False - Не найдено
	resu:= strings.Index(fullText, findText)
	// Prints true
	// fmt.Printf("%v", resu != -1)
	if resu != -1 {
		return true
	}
	return false
}
