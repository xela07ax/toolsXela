package main

import (
	"bytes"
	"compress/flate"
	"fmt"
	"github.com/xela07ax/toolsXela/arch2/archiver"
	"github.com/xela07ax/toolsXela/tp"
)

func main()  {
	//Создаем каталоги для измененных проектов
	z := archiver.Zip{
		CompressionLevel:       flate.DefaultCompression,
		MkdirAll:               true,
		SelectiveCompression:   true,
		ContinueOnError:        false,
		OverwriteExisting:      false,
		ImplicitTopLevelFolder: false,  //Неявная папка верхнего уровня
	}
	err := z.Archive([]string{"D:\\Projects\\toolsXela\\toolsXela\\arch2"}, "test.zip")
	// Сделаем буффер чтения
	var buf bytes.Buffer
	err = z.Create(&buf)
	if err != nil {
		fmt.Println(err)
		return
	}
	f, err := tp.OpenWriteFile("my.zip")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	z.Write()
	z.Archive()
	fmt.Fprintf(&buf, "Size: %d MB.", 85)

	s := buf.String()) // s == "Size: 85 MB."
	z.Create()
	fmt.Println(err)


}