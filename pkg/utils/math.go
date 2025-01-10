package utils

func Max[T Numbers](numbers ...T) (maximum T) {
	for i, number := range numbers {
		if i == 0 {
			maximum = number
			continue
		}
		if number > maximum {
			maximum = number
		}
	}
	return
}

func Min[T Numbers](numbers ...T) (minimum T) {
	for i, number := range numbers {
		if i == 0 {
			minimum = number
			continue
		}
		if number < minimum {
			minimum = number
		}
	}
	return
}

func Sum[T Numbers](numbers ...T) (sum T) {
	for _, number := range numbers {
		sum += number
	}
	return
}

func Avg[T Numbers](numbers ...T) (average float64) {
	if len(numbers) > 0 {
		average = float64(Sum(numbers...)) / float64(len(numbers))
	}
	return
}
