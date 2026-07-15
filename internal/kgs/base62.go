package kgs

const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const base = int64(len(charset))

// EncodeBase62 encodes an integer into a Base62 string.
func EncodeBase62(num int64) string {
	if num == 0 {
		return string(charset[0])
	}

	var result []byte

	// Extract digits in reverse order
	for num > 0 {
		rem := num % base
		result = append(result, charset[rem])
		num = num / base
	}

	// Reverse the result slice
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return string(result)
}
