package main

func main()  {
	// Создадим конвеер для запуска фоновой программы
	inCh0 := make(chan cmdConveer.Command, 10000) // Создаем тестовый канал с буфером в 10000, в который наполним немного синтетики
	conveerOne := cmdConveer.NewBalancer("Mail101",app.loger, inCh0) // Инициализируем объект конвеера, передадим канал для работы и еще что нибудь
	conveerOne.RunMinion() // Запускаем первого миньона, иначе работников не будет
}
