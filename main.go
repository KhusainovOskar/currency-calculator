package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv" // Импорт пакета для работы с .env файлами
)

// ExchangeRates структура для парсинга JSON-ответа от API
type ExchangeRates struct {
	BaseCode string             `json:"base_code"`        // Базовая валюта
	Rates    map[string]float64 `json:"conversion_rates"` // Курсы валют
}

func main() {
	// Загрузка переменных окружения из файла .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	// Парсинг аргументов командной строки
	amount := flag.Float64("amount", 0, "Сумма для конвертации")
	from := flag.String("from", "", "Исходная валюта (например: USD)")
	to := flag.String("to", "", "Целевая валюта (например: EUR)")
	flag.Parse()

	// Проверка валидности входных данных
	if *amount <= 0 || *from == "" || *to == "" {
		fmt.Println("Ошибка: неверные параметры")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Получение курсов валют
	rates, err := getExchangeRates(*from)
	if err != nil {
		log.Fatalf("Ошибка получения курсов: %v", err)
	}

	// Проверка наличия целевой валюты в списке курсов
	rate, exists := rates[*to]
	if !exists {
		log.Fatalf("Валюта %s не найдена", *to)
	}

	// Вычисление и вывод результата
	converted := *amount * rate
	fmt.Printf("%.2f %s = %.2f %s\n", *amount, *from, converted, *to)
}

// getExchangeRates получает курсы валют от API
func getExchangeRates(base string) (map[string]float64, error) {
	apiKey := os.Getenv("EXCHANGE_RATE_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("API ключ не найден в .env файле")
	}

	// Формирование URL для запроса
	url := fmt.Sprintf("https://v6.exchangerate-api.com/v6/%s/latest/%s", apiKey, base)

	// Выполнение HTTP GET запроса
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса: %v", err)
	}
	defer resp.Body.Close()

	// Проверка статуса ответа
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API вернул код ошибки: %d", resp.StatusCode)
	}

	// Чтение тела ответа
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения ответа: %v", err)
	}

	// Парсинг JSON
	var result ExchangeRates
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("ошибка парсинга JSON: %v", err)
	}

	return result.Rates, nil
}
