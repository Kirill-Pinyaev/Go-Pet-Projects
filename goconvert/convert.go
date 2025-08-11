package main

import (
	"fmt"
	"strings"
)

const kmToMi = 0.621371
const miToKm = 1.609344
const kgToLb = 2.2046226218
const lbToKg = 0.45359237
const cToF = 9/5 + 32
const fToC = 5 / 9

func convert(num float64, unit string) (float64, string, error) {
	unit = strings.ToLower(unit)
	switch unit {
	case "km":
		return num * kmToMi, "mi", nil
	case "mi":
		return num * miToKm, "km", nil
	case "kg":
		return num * kgToLb, "lb", nil
	case "lb":
		return num * lbToKg, "kg", nil
	case "c":
		return num * cToF, "F", nil
	case "f":
		return (num - 32) * fToC, "C", nil
	default:
		return 0, "", fmt.Errorf("unknow unit: %s", unit)
	}

}
