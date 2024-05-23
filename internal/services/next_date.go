package services

import (
	"strconv"
	"strings"
	"time"

	"fmt"
)

// Функция возвращает следующую дату для задачи в планировщике.
// Возращается дата в соответствии с параметрами, указанными в REPEAT.
// Возвращаемая дата больше даты, указанной в переменной currentDate.
func NextDate(currentDate time.Time, beginDate string, ruleRepeat string) (string, error) {
	startDate, errPars := time.Parse("20060102", beginDate)
	if errPars != nil {
		return "", fmt.Errorf("failed: incorrect DATE (%v)", errPars)
	}

	// Подготовка REPEAT для работы, очистка от пробелов вокруг символов и
	// разбивание на элементы по символу " ".
	clearRep := strings.Split(strings.TrimSpace(ruleRepeat), " ")

	// Вычисления для d-случаев.
	if clearRep[0] == "d" && len(clearRep) == 2 {
		// Получаем числа дней из REPEAT.
		days, errD := repNumsParse(clearRep[1])
		if errD != nil {
			return "", errD
		}

		// Вычисляем и модифицируем даты в соответствии с переданными в days.
		resDate, errD := dayRepeatCount(days, currentDate, startDate)
		if errD != nil {
			return "", errD
		}
		return resDate, nil
	}

	// Вычисления для y-случаев.
	if clearRep[0] == "y" && len(clearRep) == 1 {
		result := yearRepeatCount(currentDate, startDate)
		return result, nil
	}

	// Вычисления для w-случаев.
	if clearRep[0] == "w" && len(clearRep) == 2 {
		// Получаем числа дней из REPEAT.
		weekDay, errD := repNumsParse(clearRep[1])
		if errD != nil {
			return "", errD
		}

		// Вычисляем и модифицируем даты в соответствии с переданными в weekDay.
		daysRes, errD := weekRepeatCount(weekDay, currentDate, startDate)
		if errD != nil {
			return "", errD
		}

		// Из полученных дат, находим следующую ближайщую после стартовой.
		resDate := findNearestDate(daysRes, startDate)
		return resDate, nil
	}

	// Если текущая дата идет после стартовой, меняем значение расчётной на текущую.
	if currentDate.After(startDate) {
		startDate = currentDate
	}

	// Вычисления для m-случаев, только с переданными днями месяцев без указания конкретных месяцев.
	if clearRep[0] == "m" && len(clearRep) == 2 {
		// Получаем числа дней из REPEAT.
		monthDays, errD := repNumsParse(clearRep[1])
		if errD != nil {
			return "", errD
		}

		// Вычисляем и модифицируем даты в соответствии с переданными в monthDays.
		modiDateRes, errD := modifyDate(monthDays, currentDate, startDate)
		if errD != nil {
			return "", errD
		}

		// Из полученных дат, находим следующую ближайщую после стартовой.
		resDate := findNearestDate(modiDateRes, startDate)
		return resDate, nil
	}

	// Вычисления для m-случаев, с переданными днями месяцев и с указанием конкретных месяцев.
	if clearRep[0] == "m" && len(clearRep) == 3 {
		// Получаем месяца из REPEAT.
		months, errM := repNumsParse(clearRep[2])
		if errM != nil {
			return "", errM
		}
		// Вычисляем даты месяцев из переданных в REPAT и меняем день на первый.
		monthsDateRes, errM := findingMonth(months, currentDate, startDate)
		if errM != nil {
			return "", errM
		}
		// Получаем числа дней из REPEAT.
		monthDays, errD := repNumsParse(clearRep[1])
		if errD != nil {
			return "", errD
		}
		// Модифицируем даты monthRes, изменяя дни на переданные в monthDays.
		modiDateRes, errD := modifyDayMonth(monthsDateRes, monthDays)
		if errD != nil {
			return "", errD
		}
		// Из полученных дат, находим следующую ближайщую после стартовой.
		resDate := findNearestDate(modiDateRes, startDate)
		return resDate, nil
	}
	return "", fmt.Errorf("failed: incorrect REPEAT format '%s'", ruleRepeat)
}

// Находит следующую ближайшую дату к введённой в resDate из среза дат daysRes.
// Если дата из среза меньше введёной, то она пропускается.
// Возвращает следующую ближайшую дату в виде строки в формате "20060102".
// При передаче resDate в функцию, он равен startDate, используется для сравнений и основных подсчётов.
func findNearestDate(daysRes []time.Time, startDate time.Time) string {
	var ttlDat time.Time
	h, _ := time.ParseDuration("999999h")
	for _, ttl := range daysRes {
		dif := ttl.Sub(startDate)
		if dif.Hours() < h.Hours() && ttl.After(startDate) {
			h = dif
			ttlDat = ttl
		}
	}
	return ttlDat.Format("20060102")
}

// Разбивает по "," строку с числами ("1,2,3", "6,2,1", ...), конвертируя их в int
// и возвращает срез с этими числами.
func repNumsParse(clearRepeat string) ([]int, error) {
	var totalResult []int
	for _, value := range strings.Split(clearRepeat, ",") {
		num, errAtoi := strconv.Atoi(value)
		if errAtoi != nil {
			return nil, fmt.Errorf("failed: incorrect symbols in second/third REPEAT part (%s)", clearRepeat)
		}
		totalResult = append(totalResult, num)
	}
	return totalResult, nil
}
