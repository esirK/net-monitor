package main

import (
	"net/http"
	"time"

	"github.com/esirk/net-monitor/internal/data"
)

func (app *application) handlePing(w http.ResponseWriter, r *http.Request){
	// 1 way to create json response
	// response := `{"status": %d, "ping_latency": %f, "created_at": %q}`
	// response = fmt.Sprintf(response, 1, 95.03, time.Now())

	// 2nd way to create json response
	// data := map[string]interface{}{
	// 	"status": 0,
	// 	"ping_latency": 95.03,
	// 	"created_at": time.Now(),
	// }
	
	data := data.Ping{
		Status: 1,
		PingLatency: 95.03,
		CreatedAt: time.Now(),
	}
	app.writeJson(w, http.StatusOK, data, nil)
}
