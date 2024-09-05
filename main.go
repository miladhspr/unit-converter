package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

type ConversionData struct {
	Result string
}

func main() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/convert", convert)
	fmt.Println("server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.Execute(w, nil)
}

func convert(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("index.html"))
	if r.Method == "POST" {
		valueStr := r.FormValue("value")
		fromUnit := r.FormValue("from")
		toUnit := r.FormValue("to")

		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		result, err := convertUnits(value, fromUnit, toUnit)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		data := ConversionData{Result: fmt.Sprintf("%.2f %s", result, toUnit)}
		tmpl.Execute(w, data)
		return
	}
	tmpl.Execute(w, nil)
}
func convertUnits(value float64, fromUnit, toUnit string) (float64, error) {
	lengthConversions := map[string]float64{
		"millimeter": 0.001,
		"centimeter": 0.01,
		"meter":      1,
		"kilometer":  1000,
		"inch":       0.0254,
		"foot":       0.3048,
		"yard":       0.9144,
		"mile":       1609.34,
	}
	weightConversions := map[string]float64{
		"milligram": 0.000001,
		"gram":      0.001,
		"kilogram":  1,
		"ounce":     0.0283495,
		"pound":     0.453592,
	}

	f, done := tempratureCalc(value, fromUnit, toUnit)
	if done {
		return f, nil
	}

	// Length conversions
	if factorFrom, okFrom := lengthConversions[fromUnit]; okFrom {
		if factorTo, okTo := lengthConversions[toUnit]; okTo {
			return value * factorFrom / factorTo, nil
		}
	}

	// Weight conversions
	if factorFrom, okFrom := weightConversions[fromUnit]; okFrom {
		if factorTo, okTo := weightConversions[toUnit]; okTo {
			return value * factorFrom / factorTo, nil
		}
	}

	// If no conversion is possible, return an error
	return 0, fmt.Errorf("unsupported units: %s -> %s", fromUnit, toUnit)
}

func tempratureCalc(value float64, fromUnit string, toUnit string) (float64, bool) {
	if fromUnit == "celsius" && toUnit == "fahrenheit" {
		return value*9/5 + 32, true
	}
	if fromUnit == "fahrenheit" && toUnit == "celsius" {
		return (value - 32) * 5 / 9, true
	}
	if fromUnit == "celsius" && toUnit == "kelvin" {
		return value + 273.15, true
	}
	if fromUnit == "kelvin" && toUnit == "celsius" {
		return value - 273.15, true
	}
	if fromUnit == "fahrenheit" && toUnit == "kelvin" {
		return (value-32)*5/9 + 273.15, true
	}
	if fromUnit == "kelvin" && toUnit == "fahrenheit" {
		return (value-273.15)*9/5 + 32, true
	}
	return 0, false
}
