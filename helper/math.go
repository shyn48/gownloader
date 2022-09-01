package helper

import (
	"fmt"
	"math"
)

func IntToFloatString(value int) string {
	floatValue := float64(value)
	floatValue = floatValue / math.Pow(10, 6)

	return fmt.Sprintf("%f", floatValue)
}
