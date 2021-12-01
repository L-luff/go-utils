package rsync

import (
	"fmt"
	"testing"
)

func TestWeakCheckSum(t *testing.T) {
	bs := []byte{128, 128, 128, 128, 128, 128, 128}
	bs2 := []byte{128, 128, 128, 128, 128, 128, 128}
	sum := WeakCheckSum(bs)
	sum2 := WeakCheckSum(bs2)
	if sum != sum2 {
		t.Fatal("check sum error")
	}
}

func TestStrongCheckSum(t *testing.T) {
	bs := []byte{128, 128, 128, 128, 128, 128, 128}
	sum := StrongCheckSum(bs)
	fmt.Println(len(sum))
}

func TestSliceAppend(t *testing.T) {

	var a []int = nil
	a = append(a, 1, 2, 3)
	fmt.Println(a)
}
