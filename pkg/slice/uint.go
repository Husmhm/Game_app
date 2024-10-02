package slice

func DoesExist(List []uint, value uint) bool {
	for _, v := range List {
		if v == value {
			return true
		}
	}
	return false
}
