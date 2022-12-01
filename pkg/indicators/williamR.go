package indicators

func WilliamR(closes []float64) float64 {
	max, min := closes[0], closes[0]
	n := len(closes)

	for _, x := range closes[:n-1] {
		if max < x {
			max = x
		} else if min > x {
			min = x
		}
	}
	return (max - closes[n-1]) / (max - min)

}
