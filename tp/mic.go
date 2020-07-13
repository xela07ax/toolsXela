package tp

import (
	"hash/fnv"
	"os"
	"fmt"
	"path/filepath"
	"strings"
	"time"
	"math"
)


func Fck(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s", err)
		os.Exit(1)
	}
}
func FckText(text string,err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR HEAD: %s\n ERROT TEXT: %s", text,err)
		os.Exit(1)
	}
}
func DeleteFile(path string)error {
	// delete file
	var err = os.Remove(path)
	if err != nil {
		return  err
	}
	return nil
}

func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}

	return (err != nil)
}

func Getime()string  {
	return time.Now().Format("2006-01-02 15:04:05")
}

func Round(x float64, prec int) float64 {
	var rounder float64
	pow := math.Pow(10, float64(prec))
	intermed := x * pow
	_, frac := math.Modf(intermed)
	if frac >= 0.5 {
		rounder = math.Ceil(intermed)
	} else {
		rounder = math.Floor(intermed)
	}

	return rounder / pow
}

// fmt.Sprintf("Ошибка, директория расположения программы не считывается: %v", err)
func BinDir() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "" , err
	}
	return filepath.Dir(ex), nil
}

func WorkDir() (string, error) {
	ex, err := os.Getwd()
	if err != nil {
		return "" , err
	}
	return filepath.Dir(ex), nil
}

// fmt.Sprintf("Ошибка, директория не может быть создана: %v", err)
func CheckMkdir(workFolder string)error  {
	if _, err := os.Stat(workFolder); os.IsNotExist(err) {
		err = os.Mkdir(workFolder, 0777)
		if err != nil {
			return err
		}
	}
	return nil
}

func Hash32(s []byte) uint32 {
	h := fnv.New32a()
	h.Write(s)
	return h.Sum32()
}
func FindReplace(refStr string,old string,new string) string  {
	str := refStr
	for true {
		i := strings.Index(str, old) // Индекс начала
		if i > -1 { //Если мы нашли вхождение
			chars := str[:i] //Выбираем все от начала найденного индекса
			arefun := str[len(old)+len(chars):] //
			str = chars+new+arefun
		} else {
			break
		}
	}
	return str
}

