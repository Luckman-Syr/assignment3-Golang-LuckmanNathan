package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

func serveHTTP(w http.ResponseWriter, r *http.Request) {
	jsonData, err := ioutil.ReadFile("../status.json")
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
	http.HandleFunc("/", serveHTTP)
	fmt.Println("Server started at http://localhost:8000")
	http.ListenAndServe(":8000", nil)
}
