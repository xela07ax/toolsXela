package tpargs

func SortArgsMap(arg []string)map[string]string {
	// Сортируем полученные параметры
	// -jobstatus  -param P_START_DATE="03.01.2015 00:00:00"
	// Возвращает [jobstatus]"", [param]"P_START_DATE="03.01.2015 00:00:00""
	// Те аргументы которые повторяются в названии могут заменять старые, те значения которые повторяются будут утеряны!!!
	packet := make(map[string]string)
	var nameDetecter bool
	var previousName string
 	for _, val := range arg {
		parParam := val[0]
		if parParam == 45 {
			nameDetecter = true
			previousName = val[1:]
			packet[previousName] = ""
		} else {
			if nameDetecter {
				packet[previousName] = val
			}
		}
	}
	return packet
}

func main() {
	// Импорт аргументов с командной строки
	// arguments := os.Args[1:]
	args := []string{"-jobstatus", "-param", "P_START_DATE=03.01.2015 00:00:00", "-param", "P_END_DATE=13.01.2015 22:00:00"}
	packet := SortArgsMap(args)
	print(packet)
}