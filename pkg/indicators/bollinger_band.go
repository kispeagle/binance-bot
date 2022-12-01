package indicators

import "math"

func bollingerBand(closes []float64) (float64, float64, float64) {
	var avg, upper, lower float64
	avg = SMA(closes)
	sd := standardDeviation(avg, closes)
	upper = avg + 2*sd
	lower = avg - 2*sd

	return avg, upper, lower
}

func standardDeviation(avg float64, closes []float64) float64 {
	sum := 0.0
	for _, c := range closes {
		sum += (c - avg) * (c - avg)
	}
	sum /= float64(len(closes))
	return math.Sqrt(sum)
}
