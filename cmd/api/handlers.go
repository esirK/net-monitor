package main

import (
	"net/http"
	"time"

	"github.com/esirk/net-monitor/internal/data"
)

func (app *application) handlePing(w http.ResponseWriter, r *http.Request) {
	// 1 way to create json response
	// response := `{"status": %d, "ping_latency": %f, "created_at": %q}`
	// response = fmt.Sprintf(response, 1, 95.03, time.Now())

	// 2nd way to create json response
	// data := map[string]interface{}{
	// 	"status": 0,
	// 	"ping_latency": 95.03,
	// 	"created_at": time.Now(),
	// }

	pings := []data.Ping{{
		Status:      1,
		PingLatency: 95.03,
		CreatedAt:   data.CreatedAt{OurTime: time.Now()},
		Latency:     95.03,
	}, {
		Status:      0,
		PingLatency: 80.32,
		CreatedAt:   data.CreatedAt{OurTime: time.Now()},
		Latency:     80.32,
	},
	}
	err := app.writeJson(w, http.StatusOK, envelope{"pings": pings}, nil)
	if err != nil {
		app.logger.Panicf("Error writing JSON: %v", err)
		http.Error(w, "Error writing JSON", http.StatusInternalServerError)
	}
}
