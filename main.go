package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type Vitals struct {
	// Keep fields commented out if not used
	//ContactorClosed     bool          `json:"contactor_closed"`
	//VehicleConnected    bool          `json:"vehicle_connected"`
	//SessionS            int           `json:"session_s"`
	//GridV               float64       `json:"grid_v"`
	//GridHz              float64       `json:"grid_hz"`
	//VehicleCurrentA     float64       `json:"vehicle_current_a"`
	//CurrentAA           float64       `json:"currentA_a"`
	//CurrentBA           float64       `json:"currentB_a"`
	//CurrentCA           float64       `json:"currentC_a"`
	//CurrentNA           float64       `json:"currentN_a"`
	//VoltageAV           float64       `json:"voltageA_v"`
	//VoltageBV           float64       `json:"voltageB_v"`
	//VoltageCV           float64       `json:"voltageC_v"`
	//RelayCoilV          float64       `json:"relay_coil_v"`
	//PcbaTempC           float64       `json:"pcba_temp_c"`
	//HandleTempC         float64       `json:"handle_temp_c"`
	//McuTempC            float64       `json:"mcu_temp_c"`
	//UptimeS             int           `json:"uptime_s"`
	//InputThermopileUv   int           `json:"input_thermopile_uv"`
	//ProxV               float64       `json:"prox_v"`
	//PilotHighV          float64       `json:"pilot_high_v"`
	//PilotLowV           float64       `json:"pilot_low_v"`
	SessionEnergyWh float64 `json:"session_energy_wh"`
	//ConfigStatus        int           `json:"config_status"`
	//EvseState           int           `json:"evse_state"`
	//CurrentAlerts       []interface{} `json:"current_alerts"`
	//EvseNotReadyReasons []int         `json:"evse_not_ready_reasons"`
}

type HourlyPrice struct {
	MillisUTC string `json:"millisUTC"`
	PriceStr  string `json:"price"`
	Millis    int
	Price     float64
}

func getHourlyPrice() (HourlyPrice, error) {
	resp, err := http.Get("https://hourlypricing.comed.com/api?type=currenthouraverage")
	if err != nil {
		return HourlyPrice{}, fmt.Errorf("http get error: %w", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	var hourlyPrice []HourlyPrice
	err = json.NewDecoder(resp.Body).Decode(&hourlyPrice)
	if err != nil {
		return HourlyPrice{}, fmt.Errorf("json decode error: %w", err)
	} else if len(hourlyPrice) == 0 {
		return HourlyPrice{}, fmt.Errorf("no data found in response")
	}
	if hourlyPrice[0].Millis, err = strconv.Atoi(hourlyPrice[0].MillisUTC); err != nil {
		return HourlyPrice{}, fmt.Errorf("error converting millis to int: %w", err)
	}
	if hourlyPrice[0].Price, err = strconv.ParseFloat(hourlyPrice[0].PriceStr, 64); err != nil {
		return HourlyPrice{}, fmt.Errorf("error converting price to float: %w", err)
	}
	return hourlyPrice[0], nil
}

func getVitals() (Vitals, error) {
	resp, err := http.Get("http://192.168.50.229/api/1/vitals")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return Vitals{}, fmt.Errorf("http get error: %w", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	var vitals Vitals
	err = json.NewDecoder(resp.Body).Decode(&vitals)
	if err != nil {
		fmt.Printf("JSON decode error: %s\n", err)
		return Vitals{}, fmt.Errorf("json decode error: %w", err)
	}
	return vitals, nil
}

// Assume that this will be called at the end of each hour
func main() {
	vitals, err := getVitals()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	price, err := getHourlyPrice()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	UpdateTimeScaleDb(&vitals, &price)
}
