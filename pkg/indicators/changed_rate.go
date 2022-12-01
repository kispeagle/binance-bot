package indicators

func ChangedPercentage(c, p float64) float64 {
	return c/p - 1
}

func ChangedRate(c, p float64) float64 {
	return c - p
}
