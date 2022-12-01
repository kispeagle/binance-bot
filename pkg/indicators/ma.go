package indicators

func SMA(closes []float64) float64 {
	sum := 0.0
	for _, c := range closes {
		sum += c
	}
	return sum / float64(len(closes))
}

func EMA(closes []float64) float64 {
	n := len(closes)
	if n == 1.0 {
		return closes[0]
	}
	k := 2 / float64(n+1)

	return k*closes[n-1] + EMA(closes[:n-1])*(1-k)

}
