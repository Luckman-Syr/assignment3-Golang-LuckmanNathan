package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

type Status struct {
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

type StatusInfoo struct {
	Status      Status `json:"status"`
	WaterStatus string `json:"water_status"`
	WindStatus  string `json:"wind_status"`
}

func generateRandomValue() (int, int) {
	rand.Seed(time.Now().UnixNano())
	water := rand.Intn(100) + 1
	wind := rand.Intn(100) + 1
	return water, wind
}

func getStatus(water int, wind int) (string, string) {
	waterStatus := ""
	windStatus := ""

	// Determine water status
	if water < 5 {
		waterStatus = "Aman"
	} else if water >= 6 && water <= 8 {
		waterStatus = "Siaga"
	} else {
		waterStatus = "Bahaya"
	}

	// Determine wind status
	if wind < 6 {
		windStatus = "Aman"
	} else if wind >= 7 && wind <= 15 {
		windStatus = "Siaga"
	} else {
		windStatus = "Bahaya"
	}

	return waterStatus, windStatus
}

func updateJSONFile() {
	for {
		water, wind := generateRandomValue()
		waterStatus, windStatus := getStatus(water, wind)

		status := StatusInfoo{
			Status: Status{
				Water: water,
				Wind:  wind,
			},
			WaterStatus: waterStatus,
			WindStatus:  windStatus,
		}

		jsonData, _ := json.MarshalIndent(status, "", "    ")
		_ = ioutil.WriteFile("../status.json", jsonData, 0644)

		time.Sleep(15 * time.Second)
		fmt.Println("Status updated")
	}
}

func main() {
	updateJSONFile()

	fmt.Println("Client started at http://localhost:8001")
	http.ListenAndServe(":8001", nil)
}
