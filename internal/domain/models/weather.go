package models

type Current struct {
	Celsius    float64 `json:"temp_c"`
	Fahrenheit float64 `json:"temp_f"`
}

type Weather struct {
	Current Current `json:"current"`
}
