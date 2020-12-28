package qarr

func StringsIn(in []string, str string) bool {
	for _, s := range in {
		if s == str {
			return true
		}
	}
	return false
}
