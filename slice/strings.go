package slice

func ContainsString(ss []string, s string) bool {
	for i := range ss  {
		if s == ss[i] {
			return true
		}
	}
	return false
}