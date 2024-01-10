package main

import (
	"time" // Предоставляет функциональность для работы с временем и задержками в программе
	"fmt" // Предоставляет функциональность для форматирования ввода-вывода, включая печать текста и чтение ввода с консоли
)

type Token struct {
	Message   string // Данные сообщения
	Recipient int    // Номер узла-получателя сообщения
	TTL       int    // Время жизни токена
}

func processToken(id int, in <-chan Token, out chan<- Token) {
	for {
		token := <-in // Получаем токен из входного канала

		// Выводим информацию о полученном токене
		fmt.Printf("Узел %d получил токен. TTL: %d\n", id, token.TTL)

		if token.Recipient == id { // Если узел-получатель -- текущий узел, то сообщение достигло цели
			fmt.Printf("Узел %d получил сообщение: %s\n", id, token.Message)
			continue
		}

		if token.TTL > 0 { // Если время жизни токена еще не истекло, передаем токен следующему узлу
			token.TTL-- // Уменьшаем время жизни токена
			out <- token
		} else { // Если время жизни токена истекло, выводим соответствующее сообщение
			fmt.Printf("Узел %d: Истек срок действия для токена: %s\n", id, token.Message)
		}
	}
}

func initializeChannels(N int) []chan Token {
	channels := make([]chan Token, N) // Создание массива каналов заданного размера N
	for i := range channels {
		channels[i] = make(chan Token) // Инициализация каждого элемента канала в массиве
	}
	return channels // Возвращает инициализированный массив каналов
}

func startNodeGoroutines(channels []chan Token, N int) {
	for i := 0; i < N-1; i++ { // Запускаем горутины для каждого узла (кроме последнего)
		go processToken(i, channels[i], channels[i+1])
	}
	go processToken(N-1, channels[N-1], channels[0]) // Запускаем горутину для последнего узла
}

func getUserInput(N int) (string, int, int) {
	var message string
	var recipient, ttl int

	fmt.Println("Введите данные сообщения:")
	fmt.Scanln(&message)

	fmt.Println("Введите номер узла-получателя (от 0 до", N-1, "):")
	fmt.Scanln(&recipient)
	if recipient < 0 || recipient > N-1 {
		fmt.Println("Неверный номер узла-получателя")
		return "", 0, 0
	}

	fmt.Println("Введите время жизни (TTL):")
	fmt.Scanln(&ttl)

	return message, recipient, ttl
}

func runProgram() {
	var numNodes int
	fmt.Println("Введите количество узлов:")
	fmt.Scanln(&numNodes)

	channels := initializeChannels(numNodes) // Инициализируем каналы для связи между узлами
	startNodeGoroutines(channels, numNodes)  // Запускаем горутины для каждого узла

	message, recipient, ttl := getUserInput(numNodes) // Получаем данные сообщения от пользователя
	if message == "" {
		return
	}

	initialToken := Token{
		Message:   message,
		Recipient: recipient,
		TTL:       ttl,
	}

	fmt.Printf("Отправка токена с узла 0 на узел %d\n", recipient)
	channels[0] <- initialToken // Передаем исходный токен в первый узел (исходный узел)

	time.Sleep(time.Second * 2) // Пауза для завершения работы программы
}

func main() {
	runProgram()
}
