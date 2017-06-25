package kanaelib

import (
	"math"

	"github.com/kardianos/osext"
)

func Check(e error) {
	if e != nil {
		//raven.CaptureErrorAndWait(e, nil)
		panic(e)
	}
}

func GetExecutablePath() string {
	path, err := osext.ExecutableFolder()
	Check(err)
	return path
}

// https://stackoverflow.com/questions/18390266/how-can-we-truncate-float64-type-to-a-particular-precision-in-golang
func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func ToPoloniexFixed(num float64) float64 {
	prec := 8
	return ToFixed(num, prec)
}
