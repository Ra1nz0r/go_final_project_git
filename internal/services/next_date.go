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

	// Создаем мапу цифр из переданных значений в запросе.
	numRepeatTask, errMap := repNumsParse(clearRep)
	if errMap != nil {
		return "", errMap
	}

	resDate := startDate // дата с которой производятся расчеты и сравнения в функции

	if clearRep[0] == "d" && len(clearRep) == 2 {
		result, errDay := dayRepeatCount(numRepeatTask, currentDate, startDate, resDate)
		if errDay != nil {
			return result, errDay
		}
		return result, nil
	}

	if clearRep[0] == "y" && len(clearRep) == 1 {
		result := yearRepeatCount(currentDate, startDate, resDate)
		return result, nil
	}

	resMap := make(map[int]time.Time) // результирующая для первых чисел в repeat

	if clearRep[0] == "w" && len(clearRep) == 2 {
		result, errWeek := weekRepeatCount(numRepeatTask, resMap, currentDate, startDate, resDate)
		if errWeek != nil {
			return result, errWeek
		}
		return result, nil
	}

	if clearRep[0] == "m" && (len(clearRep) == 3 || len(clearRep) == 2) {
		result, errMonth := monthRepeatCount(numRepeatTask, resMap, currentDate, startDate, resDate)
		if errMonth != nil {
			return result, errMonth
		}
		return result, nil
	}
	return "", fmt.Errorf("failed: incorrect REPEAT format '%s'", ruleRepeat)
}

// Вычисляет следующую дату для запроса REPEAT с буквой 'd' - задача переносится на указанное число дней.
// Максимально допустимое число равно 400. Примеры:
// d 1 — каждый день;
// d 7 — для вычисления следующей даты добавляем семь дней;
// d 60 — переносим на 60 дней.
func dayRepeatCount(numRepeatTask map[int][]int, currentDate, startDate, resDate time.Time) (string, error) {
	// dNum - число переданное в REPEAT, для примера [d 56] оно будет равно 56,
	// должно быть в диапазоне 1-400 и может быть только одно ([d 56 1], ошибка).
	dNum := numRepeatTask[0][0]
	if dNum < 1 || dNum > 400 {
		return "", fmt.Errorf("failed: value DAY must be between 1 and 400 (your val '%d')", dNum)
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

// Вычисляет следующую дату для запроса REPEAT с буквой 'y' - задача переносится год.
func yearRepeatCount(currentDate, startDate, resDate time.Time) string {
	// Пока текущая дата идет после расчётной, прибавляем год к расчётной.
	for currentDate.After(resDate) {
		resDate = resDate.AddDate(1, 0, 0)
	}
	// Если расчётная дата не изменилась и равна стартовой, то прибавляем год к расчётной.
	if resDate == startDate {
		resDate = resDate.AddDate(1, 0, 0)
	}
	return resDate.Format("20060102")
}

// Вычисляет следующую дату для запроса REPEAT с буквой 'w' - задача назначается в указанные дни недели.
// w <через запятую от 1 до 7> , где 1 — понедельник, 7 — воскресенье. Например:
// w 7 — задача перенесётся на ближайшее воскресенье;
// w 1,4,5 — задача перенесётся на ближайший понедельник, четверг или пятницу;
// w 2,3 — задача перенесётся на ближайший вторник или среду.
func weekRepeatCount(numRepeatTask map[int][]int, resMap map[int]time.Time, currentDate, startDate, resDate time.Time) (string, error) {
	// Достаем числа, переданные в REPEAT, из мапы и работаем с ними.
	for _, value := range numRepeatTask {
		for _, wNum := range value {
			var cntDayPass int // подсчёт пройденных дней
			if wNum < 1 || wNum > 7 {
				return "", fmt.Errorf("failed: value DAY_WEEK must be between 1 and 7 (your val '%d')", wNum)
			}
			// Если текущая дата идет после стратовой, меняем значение расчётной на текущую.
			if currentDate.After(startDate) {
				resDate = currentDate
			}
			// Пока день недели расчётной даты не равен указанной в REPEAT, прибавляем дни и считаем количество.
			for ok := true; ok; ok = (resDate.Weekday() != time.Weekday(wNum%7)) { // Находим остаток, потому что Воскресенье = 0, а не 7.
				resDate = resDate.AddDate(0, 0, 1)
				cntDayPass++
			}
			// Записываем количество пройденных дней и итоговую дату в результирующую мапу.
			resMap[cntDayPass] = resDate
			// Сбрасываем итоговую дату на стартовую.
			resDate = startDate
		}
	}
	return resMap[MinValue(resMap)].Format("20060102"), nil
}

// Вычисляет следующую дату для запроса REPEAT с буквой 'm' - задача назначается в указанные дни недели,
// m <через запятую от 1 до 31,-1,-2> [через запятую от 1 до 12].
// При этом вторая последовательность чисел опциональна и указывает на определённые месяцы. Например:
// m 4 — задача назначается на четвёртое число каждого месяца;
// m 1,15,25 — задача назначается на 1-е, 15-е и 25-е число каждого месяца;
// m -1 — задача назначается на последний день месяца;
// m -2 — задача назначается на предпоследний день месяца;
// m 3 1,3,6 — задача назначается на 3-е число января, марта и июня;
// m 1,-1 2,8 — задача назначается на 1-е и последнее число число февраля и авгуcта.
func monthRepeatCount(numRepeatTask map[int][]int, resMap map[int]time.Time, currentDate, startDate, resDate time.Time) (string, error) {
	resMapMonth := make(map[int]time.Month) // результирующая для вторых чисел в repeat, используется при передаче m значения

	// перебор с конца, чтобы сначала работать с месяцами, найти ближайший и передать значение для работы с первыми числами (дни месяца)
	for key := len(numRepeatTask) - 1; key >= 0; key-- {
		for _, mNum := range numRepeatTask[key] {
			switch key {
			case 0:
				if mNum < -2 || mNum > 31 {
					return "", fmt.Errorf("failed: value DAY_MONTH must be between -2 and 31 (your val '%d')", mNum)
				}
				// Если текущая дата идёт после стартовой, то перезаписываем итоговую на текущую.
				if currentDate.After(startDate) {
					resDate = currentDate
				}
				var cntDayPass int
				switch {
				case mNum < 0:
					// Увеличиваем месяц на один от итогового, чтобы сравнивать и искать последний день месяца.
					compareMonth := resDate.Month() + 1
					// Пока итоговый месяц не равен сравниваемому, то прибавляем дни и считаем их кол-во.
					for ok := true; ok; ok = (resDate.Month() != compareMonth) {
						resDate = resDate.AddDate(0, 0, 1)
						cntDayPass++
					}
					// Записываем количество пройденных дней и итоговую дату в результирующую мапу,
					// вычитая один(два) день(я) для нахождения последнего(предпоследнего) дня месяца.
					resMap[cntDayPass] = resDate.AddDate(0, 0, mNum)
				default:
					// Если мапа с месяцами не пустая, то обновляем значение итоговой даты,
					if len(resMapMonth) > 0 {
						resDate = resDate.AddDate(0, MinMonthCnt(resMapMonth)-1, 0)
					}
					// Пока день итоговый даты не равен переданному, то прибавляем дни и считаем их кол-во.
					for ok := true; ok; ok = (resDate.Day() != mNum) {
						resDate = resDate.AddDate(0, 0, 1)
						cntDayPass++
					}
					// Записываем количество пройденных дней и итоговую дату в результирующую мапу.
					resMap[cntDayPass] = resDate
				}
				// Сбрасываем итоговую дату на стартовую.
				resDate = startDate
			case 1:
				// см. описание функции
				if errNear := nearestMonth(resMapMonth, mNum, currentDate, startDate); errNear != nil {
					return "", errNear
				}
			}
		}
	}
	return resMap[MinValue(resMap)].Format("20060102"), nil
}

// Вычисляет количество пройденных месяцев от стартовой даты, записывает результат
// в формате - map['число пройденных месяцев']'месяц'. Если текущая дата идёт после стартовой,
// то стартовая дата перезаписывается на текущую и отсчёт идет от неё.
// Возвращает ошибку, если mNum не равен ни одному числу месяца.
func nearestMonth(resMapMonth map[int]time.Month, mNum int, currentDate, startDate time.Time) error {
	if mNum < 1 || mNum > 12 {
		return fmt.Errorf("failed: value MONTH must be between 1 and 12 (your val '%d')", mNum)
	}
	var cntMonthPass int // количество пройденных месяцев

	if currentDate.After(startDate) {
		startDate = currentDate
	}
	// Пока итоговый месяц не равен переданному, то прибавляем месяца и считаем их кол-во.
	for ok := true; ok; ok = (startDate.Month() != time.Month(mNum)) {
		startDate = startDate.AddDate(0, 1, 0)
		cntMonthPass++
	}
	// Записываем количество пройденных месяцев и итоговый месяц в результирующую мапу.
	resMapMonth[cntMonthPass] = time.Month(mNum)

	return nil
}

// Проверка и создание числовой карты для работы, из чисел в REPEAT-запросе,
// исключая первый символ, со сменой типа string на int.
// А также частичная проверка на корректность запроса.
// В случае некорректности запроса, возвращает ошибку.
func repNumsParse(clearRepeat []string) (map[int][]int, error) {
	numResMap := make(map[int][]int)
	for key, value := range clearRepeat[1:] {
		var totalResult []int
		for _, j := range strings.Split(value, ",") {
			num, errAtoi := strconv.Atoi(j)
			if errAtoi != nil {
				return nil, fmt.Errorf("failed: incorrect symbols in second/third REPEAT part (%s)", clearRepeat)
			}
			totalResult = append(totalResult, num)
		}
		numResMap[key] = totalResult
	}
	return numResMap, nil
}
