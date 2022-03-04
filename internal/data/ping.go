package data

type Ping struct {
	Status 				int 		`json:"status"`
	PingLatency 		float64 	`json:"ping_latency"`
	CreatedAt 			CreatedAt 	`json:"created_at"`
	Latency 			Latency 	`json:"latency"`
}
