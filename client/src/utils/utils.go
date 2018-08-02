package utils

// Check for an error, panic if one has occurred
func CheckForError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
