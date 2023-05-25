package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*2000)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		fmt.Println("Error on request", err)
		return
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error on response", err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error on reading response body", err)
		return
	}

	var data map[string]string
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error on unmarshaling JSON", err)
		return
	}

	cotacao := data["Dólar"]

	file, err := os.Create("cotacao.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error to create the file text: %v", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString("Dólar: " + cotacao)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error to write file: %v", err)
		return
	}

	fmt.Println("Cotação do Dólar salva com sucesso!")
}
