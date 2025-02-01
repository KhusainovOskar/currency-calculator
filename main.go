package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv" 
)

type ExchangeRates struct {
	BaseCode string             `json:"base_code"`        
	Rates    map[string]float64 `json:"conversion_rates"` 
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	amount := flag.Float64("amount", 0, "Сумма для конвертации")
	from := flag.String("from", "", "Исходная валюта")
	to := flag.String("to", "", "Целевая валюта")
	flag.Parse()

	if *amount <= 0 || *from == "" || *to == "" {
		fmt.Println("Ошибка: неверные параметры")
		flag.PrintDefaults()
		os.Exit(1)
	}

	rates, err := getExchangeRates(*from)
	if err != nil {
		log.Fatalf("Ошибка получения курсов: %v", err)
	}

	rate, exists := rates[*to]
	if !exists {
		log.Fatalf("Валюта %s не найдена", *to)
	}

	converted := *amount * rate
	fmt.Printf("%.2f %s = %.2f %s\n", *amount, *from, converted, *to)
}

func getExchangeRates(base string) (map[string]float64, error) {
	apiKey := os.Getenv("EXCHANGE_RATE_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("API ключ не найден в .env файле")
	}

	url := fmt.Sprintf("https://v6.exchangerate-api.com/v6/%s/latest/%s", apiKey, base)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API вернул код ошибки: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения ответа: %v", err)
	}

	var result ExchangeRates
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("ошибка парсинга JSON: %v", err)
	}

	return result.Rates, nil
}
