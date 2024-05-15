package services

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Проверка и создание числовой карты для работы, из чисел в REPEAT-запросе,
// исключая первый символ, со сменой типа string на int.
// А также частичная проверка на корректность запроса.
// В случае некорректности запроса, возвращает ошибку.
func repNumsParse(clearRepeat []string) (map[int][]int, error) {
	var totalResult []int
	numResMap := make(map[int][]int)
	for key, value := range clearRepeat[1:] {
		for _, j := range strings.Split(value, ",") {
			num, errAtoi := strconv.Atoi(j)
			if errAtoi != nil {
				return nil, fmt.Errorf("failed: incorrect symbols in second/third REPEAT part (%s)", clearRepeat)
			}
			totalResult = append(totalResult, num)
		}
		numResMap[key] = totalResult
		totalResult = []int{}
	}
	return numResMap, nil
}

// Функция возвращает следующую дату для задачи в планировщике.
// Возращается дата в соответствии с параметрами, указанными в REPEAT.
// Возвращаемая дата больше даты, указанной в переменной currentDate.
func NextDate(currentDate time.Time, beginDate string, ruleRepeat string) (string, error) {
	startDate, errPars := time.Parse("20060102", beginDate)
	if errPars != nil {
		return "", fmt.Errorf("failed: incorrect DATE (%v)", errPars)
	}

	// Подготовка REPEAT-запроса для работы, очистка от пробелов вокруг символов.
	clearRep := strings.Split(strings.TrimSpace(ruleRepeat), " ")
	if clearRep == nil {
		return "", fmt.Errorf("failed: incorrect REPEAT format")
	}

	// Создаем мапу цифр из переданных значений в запросе.
	numRepeatTask, errMap := repNumsParse(clearRep)
	if errMap != nil {
		return "", errMap
	}
	//-------------------------------------------------

	resDate := startDate // дата с которой производятся расчеты и сравнения в функции

	resMap := make(map[uint16]time.Time) // результирующая для первых чисел в repeat
	resMin := ^uint16(0)                 // хранит наименьшее значение пройденных дней из мапы resMap

	//resMapMonth := make(map[uint8]time.Month) // результирующая для вторых чисел в repeat, используется при передаче m значения
	//resMinMonth := ^uint8(0)                  // хранит наименьшее значение пройденных месяцев из мапы resMapMonth

	if clearRep[0] == "d" && len(clearRep) == 2 {
		// dNum - число переданное в REPEAT, для примера [d 56] оно будет равно 56,
		// должно быть в диапазоне 1-400 и может быть только одно ([d 56 1], ошибка).
		dNum := numRepeatTask[0][0]

		if dNum < 1 || dNum > 400 {
			return "", fmt.Errorf("failed: value (%d) DAY must be between 1 and 400", dNum)
		}
		// Пока текущая дата идет после расчётной, прибавляем указанные дни к расчётной.
		for currentDate.After(resDate) {
			resDate = resDate.AddDate(0, 0, dNum)
		}
		// Если итоговая дата не изменилась и равна расчётной, то прибавляем дни к расчётной.
		if resDate == startDate {
			resDate = resDate.AddDate(0, 0, dNum)
		}

		return resDate.Format("20060102"), nil
	}

	if clearRep[0] == "y" && len(clearRep) == 1 {
		// Пока текущая дата идет после расчётной, прибавляем год к расчётной.
		for currentDate.After(resDate) {
			resDate = resDate.AddDate(1, 0, 0)
		}
		// Если расчётная дата не изменилась и равна стартовой, то прибавляем год к расчётной.
		if resDate == startDate {
			resDate = resDate.AddDate(1, 0, 0)
		}
		return resDate.Format("20060102"), nil
	}

	if clearRep[0] == "w" && len(clearRep) == 2 {
		// Достаем числа, переданные в REPEAT, из мапы и работаем с ними.
		for _, value := range numRepeatTask {
			for _, wNum := range value {
				var cntDayPass uint16 // подсчёт пройденных дней
				if wNum < 1 || wNum > 7 {
					return "", fmt.Errorf("failed: value (%d) DAY_WEEK must be between 1 and 7", wNum)
				}
				// Если текущая дата идет после стратовой, меняем значение расчётной на текущую.
				if currentDate.After(startDate) {
					resDate = currentDate
				}
				// Для правильности подсчёта дня недели Воскресенья, меняем цифру на 0.
				if wNum == 7 {
					wNum = 0
				}
				// Пока день недели расчётной даты не равен указанной в REPEAT, прибавляем дни и считаем количество.
				for ok := true; ok; ok = (resDate.Weekday() != time.Weekday(wNum)) {
					resDate = resDate.AddDate(0, 0, 1)
					cntDayPass++
				}
				// Записываем количество пройденных дней и итоговую дату в результирующую мапу.
				resMap[cntDayPass] = resDate
				// Если текущее кол-во пройденных дней меньше сохранненого, то перезаписываем значение.
				if cntDayPass < resMin {
					resMin = cntDayPass
				}
				// Сбрасываем итоговую дату на стартовую.
				resDate = startDate
			}
		}
		return resMap[resMin].Format("20060102"), nil
	}

	return "", fmt.Errorf("failed: incorrect REPEAT format '%s'", ruleRepeat)
}
