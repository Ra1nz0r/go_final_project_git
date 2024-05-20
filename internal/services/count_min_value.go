package services

import "time"

func MinValue(resMap map[int]time.Time) int {
	var minNumber int
	for minNumber = range resMap {
		break
	}
	for n := range resMap {
		if n < minNumber {
			minNumber = n
		}
	}
	return minNumber
}

func MinMonthCnt(resMap map[int]time.Month) int {
	var minNumber int
	for minNumber = range resMap {
		break
	}
	for n := range resMap {
		if n < minNumber {
			minNumber = n
		}
	}
	return minNumber
}
