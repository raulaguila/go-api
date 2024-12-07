package utils

// PanicIfErr triggers a panic if the given error is not nil.
func PanicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
