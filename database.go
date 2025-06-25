package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"os"
	"time"
)

func UpdateTimeScaleDb(vitals *Vitals, price *HourlyPrice) {
	ctx := context.Background()
	connStr := os.Getenv("TIMESCALE_DB_URI")
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		fmt.Printf("Unable to connect to database: %v\n", err)
		return
	}
	defer func(conn *pgx.Conn, ctx context.Context) {
		_ = conn.Close(ctx)
	}(conn, ctx)

	// Get the previous wh reading
	var prevWh float64
	err = conn.QueryRow(ctx, "SELECT total_wh FROM charging_history ORDER BY time DESC LIMIT 1").Scan(&prevWh)
	if err != nil {
		fmt.Printf("SELECT Error: %s\n", err)
		return
	}
	var whInLastHour float64
	if vitals.SessionEnergyWh-prevWh < 0 {
		whInLastHour = vitals.SessionEnergyWh
	} else {
		whInLastHour = vitals.SessionEnergyWh - prevWh
	}
	// Insert the new wh reading
	insertQuery := `INSERT INTO comed_price (time, price) VALUES ($1, $2)`
	_, err = conn.Exec(ctx, insertQuery, time.Now(), price.Price)
	if err != nil {
		fmt.Printf("INSERT comed_price Error: %s\n", err)
	}
	// Insert the new wh reading
	insertQuery = `INSERT INTO charging_history (time, total_wh, wh_difference) VALUES ($1, $2, $3)`
	_, err = conn.Exec(ctx, insertQuery, time.Now(), vitals.SessionEnergyWh, whInLastHour)
	if err != nil {
		fmt.Printf("INSERT charging_history Error: %s\n", err)
	}
}
