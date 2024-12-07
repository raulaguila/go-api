package utils

// Max computes and returns the maximum value from a variadic list of numbers of a specified numeric type.
// If the list is empty, the zero value of the type is returned.
// The function accepts arguments that implement the Numbers interface, which includes int, float types, etc.
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

// Min computes and returns the minimum value from a variadic list of numbers of a specified numeric type.
// If the list is empty, the zero value of the type is returned.
// The function accepts arguments that implement the Numbers interface, which includes int, float types, etc.
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

// Sum computes and returns the sum value from a variadic list of numbers of a specified numeric type.
// If the list is empty, the zero value of the type is returned.
// The function accepts arguments that implement the Numbers interface, which includes int, float types, etc.
func Sum[T Numbers](numbers ...T) (sum T) {
	for _, number := range numbers {
		sum += number
	}
	return
}

// Avg computes and returns the average value from a variadic list of numbers of a specified numeric type.
// If the list is empty, the zero value of the type is returned.
// The function accepts arguments that implement the Numbers interface, which includes int, float types, etc.
func Avg[T Numbers](numbers ...T) (average float64) {
	if len(numbers) > 0 {
		average = float64(Sum(numbers...)) / float64(len(numbers))
	}
	return
}
