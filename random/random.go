package random

import (
	"log"
	"math"
	"math/rand"
	"strings"
)

const (
	LetterBits     = 6
	LetterBitsMask = 1<<6 - 1
	LetterIndexMax = 63 / 6
)

const LetterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomStr(l int) string {
	sb := strings.Builder{}
	bytesLen := len(LetterBytes)
	for i, val, loopIndex := 0, rand.Int63(), LetterIndexMax; i < l; {
		if loopIndex <= 0 {
			val, loopIndex = rand.Int63(), LetterIndexMax
		}
		if v := int(val & LetterBitsMask); v < bytesLen {
			i++
			sb.WriteByte(LetterBytes[v])
		}
		loopIndex--
		val >>= LetterBits
	}
	return sb.String()
}
func RandomIntMN(min, max int) int {
	if min > max {
		log.Panicln("min is greater than max")
	}
	incrementVal := int(math.Abs(float64(min)))
	minN := min + incrementVal
	maxN := max + incrementVal
	ranVal := RandomIntM(maxN-minN+1) + minN
	return ranVal - incrementVal
}
func RandomIntM(max int) int {

	if max < 0 {
		return max - 1 - rand.Int()
	}
	return rand.Intn(max)
}
