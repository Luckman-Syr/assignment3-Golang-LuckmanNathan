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

type StatusInfo struct {
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

		status := StatusInfo{
			Status: Status{
				Water: water,
				Wind:  wind,
			},
			WaterStatus: waterStatus,
			WindStatus:  windStatus,
		}

		jsonData, _ := json.MarshalIndent(status, "", "    ")
		_ = ioutil.WriteFile("status.json", jsonData, 0644)

		time.Sleep(15 * time.Second)
		fmt.Println("Status updated!")
	}
}

func serveHTTP(w http.ResponseWriter, r *http.Request) {
	jsonData, err := ioutil.ReadFile("status.json")
	if err != nil {
		http.Error(w, "Failed to read status file", http.StatusInternalServerError)
		return
	}

	var status StatusInfo
	err = json.Unmarshal(jsonData, &status)
	if err != nil {
		http.Error(w, "Failed to parse status file", http.StatusInternalServerError)
		return
	}

	html := fmt.Sprintf(`
        <html>
        <head>
            <meta http-equiv="refresh" content="15">
        </head>
        <body>
            <h1>Status:</h1>
            <p>Water Level: %d meter (%s)</p>
            <p>Wind Speed: %d meter/s (%s)</p>
        </body>
        </html>
    `, status.Status.Water, status.WaterStatus, status.Status.Wind, status.WindStatus)

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, html)
}

func main() {
	go updateJSONFile()

	http.HandleFunc("/", serveHTTP)
	fmt.Println("Server started at http://localhost:8000")
	http.ListenAndServe(":8000", nil)
}
