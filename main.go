package main


import (
	"bytes"
	"compress/flate"
	"fmt"
	"github.com/xela07ax/toolsXela/archiver"
	"github.com/xela07ax/toolsXela/tp"
)

func main()  {
	ts, _ := tp.BinDir()
	fmt.Printf("Bit: %s\n", ts)
	archiver.Print()
	z := archiver.Zip{
		CompressionLevel:       flate.DefaultCompression,
		MkdirAll:               true,
		SelectiveCompression:   true,
		ContinueOnError:        false,
		OverwriteExisting:      false,
		ImplicitTopLevelFolder: false,  //Неявная папка верхнего уровня
	}
	//err := z.Archive([]string{"D:\\Projects\\toolsXela\\toolsXela\\arch2"}, "test.zip")
	var b bytes.Buffer // A Buffer needs no initialization.
	err := z.ArchiveWriter([]string{"D:\\Projects\\toolsXela\\toolsXela\\arch2"}, &b)
	f, err := tp.CreateOpenFile("tx.zip")
		//b.WriteTo(os.Stdout)
		b.WriteTo(f)
	fmt.Println(err)
}