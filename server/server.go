package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Money struct {
	ID         int    `gorm:"primaryKey"`
	Code       string `json:"code"`
	CodeIn     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

func main() {
	http.HandleFunc("/cotacao", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	println("Requisição inicializada")
	defer println("Requisição finalizada")

	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	ctxDB, cancel2 := context.WithTimeout(ctx, time.Second*1)
	defer cancel()
	defer cancel2()

	db, err := gorm.Open(sqlite.Open("../db.sqlite"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to get database instance")
	}
	defer sqlDB.Close()

	err = db.AutoMigrate(&Money{})
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
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

	var data map[string]Money
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error on unmarshaling JSON", err)
		return
	}

	result := map[string]string{"Dólar": data["USDBRL"].Bid}
	jsonValue, err := json.Marshal(result)
	if err != nil {
		fmt.Println("Error on marshaling JSON", err)
		return
	}

	money := data["USDBRL"]
	err = db.WithContext(ctxDB).Create(&money).Error
	if err != nil {
		fmt.Println("Error on saving data to database", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonValue)
}
