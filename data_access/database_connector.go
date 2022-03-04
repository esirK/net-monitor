package data_access

import (
	"context"
	"fmt"

	"github.com/esirk/net-monitor/types"
	"github.com/jackc/pgx/v4"
)

func SaveState(ch chan types.PingResult) {
	for result := range ch {
		saveData(result.State, result.Ping_time)
	}
}

func saveData(state int, ping_time float64) {
	conn, err := pgx.Connect(context.Background(), "postgresql://bluetail:bluetail@localhost:5432/net_mon")
	if err != nil {
		fmt.Println("Error on Connect ", err)
	}
	defer conn.Close(context.Background())
	var ping int
	sqlStatement := `INSERT INTO pings (status, ping_time) VALUES ($1, $2) RETURNING id`

	err = conn.QueryRow(context.Background(), sqlStatement, state, ping_time).Scan(&ping)
	if err != nil {
		fmt.Println("Error on QueryRow ", err)
	}
	fmt.Println("Ping: ", ping)
}
