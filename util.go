package tchart

func Min(i, j int) int {
	if i == j {
		return i
	}
	if i < j {
		return i
	}
	return j
}

func Max(i, j int) int {
	if i == j {
		return i
	}
	if i > j {
		return i
	}
	return j
}

func MinFloat64(i, j float64) float64 {
	if i == j {
		return i
	}
	if i < j {
		return i
	}
	return j
}

func MaxFloat64(i, j float64) float64 {
	if i == j {
		return i
	}
	if i > j {
		return i
	}
	return j
}
