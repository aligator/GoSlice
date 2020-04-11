package util

type coord int

func bool2Int(flag bool) int {
	if flag {
		return 1
	}
	return 0
}

func (c coord) intToMM() float64 {
	return float64(c) / 1000.0
}

func (c coord) intToMM2() float64 {
	return float64(c) / 1000000.0
}

func (c coord) intToMicron() float64 {
	return float64(c) / 1
}

func MicronToInt(n int) float64 {
	return float64(n) * 1
}

func mmToInt(n float64) coord {
	return coord(n*100 + 0.5*(float64(bool2Int(n > 0))-float64(bool2Int(n < 0))))
}

func mm2ToInt(n float64) coord {
	return coord(n*1000000 + 0.5*(float64(bool2Int(n > 0))-float64(bool2Int(n < 0))))
}

func mm3ToInt(n float64) coord {
	return coord(n*1000000000 + 0.5*(float64(bool2Int(n > 0))-float64(bool2Int(n < 0))))
}
