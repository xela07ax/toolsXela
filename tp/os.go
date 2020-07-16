package tp

import (
	"fmt"
	"os"
	"time"
)

func ExitWithSecTimeout(status int)  {
	// Го любит завершать свою работу раньше чем сделать все завершающие операции, но все же он остается очень быстрым одной секунды ему достатотчно
	// 0 - norm
	// 1 - error
	fmt.Println("Завершение работы программы через 2 сек.")
	time.Sleep(500*time.Millisecond)
	os.Exit(status)
}
