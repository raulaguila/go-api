package helper

// PanicIfErr triggers a panic if the provided error is not nil.
func PanicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
