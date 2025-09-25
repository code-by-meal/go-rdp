package core

func Reverse(arr []byte) []byte {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}

	return arr
}

// Joining to byte slices
func Join(a, b []byte) []byte {
	a = append(a, b...)

	return a
}
