package internal

// ImpossibleError panics if err is not nil. It is used to document errors that
// are thought to be impossible without ignoring them if they actually happen.
func ImpossibleError(err error) {
	if err != nil {
		panic(err)
	}
}
