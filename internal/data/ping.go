package data

import "time"

type Ping struct {
	Status 				int 		`json:"status"`
	PingLatency 		float64 	`json:"ping_latency"`
	CreatedAt 			time.Time 	`json:"created_at"`
}
