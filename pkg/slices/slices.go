package slices

func Filter[T any](ss []T, fill func(T) bool) (ret []T) {
	for _, s := range ss {
		if fill(s) {
			ret = append(ret, s)
		}
	}
	return ret
}

func Remove[T any](ss []T, remove func(T) bool) (ret []T) {
	for i := len(ss) - 1; i >= 0; i-- {
		if remove(ss[i]) {
			ss = append(ss[:i], ss[i+1:]...)
		}
	}
	return ss
}

func Contains[T any](ss []T, exists func(T) bool) (int, bool) {
	for i, s := range ss {
		if exists(s) {
			return i, true
		}
	}
	return -1, false
}
