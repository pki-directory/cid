package base16

const (
	tableLower = "0123456789abcdef"
	tableUpper = "0123456789ABCDEF"
)

func Encode(dst, src []byte) {
	j := 0
	for i := 0; i < len(src); i++ {
		b := src[i]
		dst[j] = tableLower[b>>4]
		dst[j+1] = tableLower[b&0x0f]
		j += 2
	}
}

func EncodeUpper(dst, src []byte) {
	j := 0
	for i := 0; i < len(src); i++ {
		b := src[i]
		dst[j] = tableUpper[b>>4]
		dst[j+1] = tableUpper[b&0x0f]
		j += 2
	}
}

func EncodedLen(n int) int {
	return n * 2
}
