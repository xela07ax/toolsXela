package encod

import (
	"encoding/binary"
	"fmt"
)


func IntToBytes(id int) []byte {
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, uint64(id))
	return key
}

func BytesToInt(key []byte) int {
	return int(binary.BigEndian.Uint64(key))
}

func main()  {
	id := 1563
	fmt.Println(id)
	dat := IntToBytes(id)
	fmt.Printf("Byt:%v|\n",dat)
	id2 := BytesToInt(dat)
	fmt.Printf("%d\n",id2)
	/*
	1563
	Byt:[0 0 0 0 0 0 6 27]|
	1563
	*/
}