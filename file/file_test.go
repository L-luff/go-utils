package main

import (
	"bufio"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestCreateFile(t *testing.T) {
	path := "D:\\my\\project\\go\\go-exampleabc\\file"
	CreateFile(100, 101, 0, 200, path)
}

func TestGeneratorFileRangeSize(t *testing.T) {
	path := "D:\\my\\project\\go\\go-exampleabc\\file"
	for i := 0; i < 10000; i++ {
		err := GeneratorFileRangeSize(path, 100, 101, KB)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestGeneratorFileRangeSize2(t *testing.T) {
	s := string(os.PathSeparator)
	fmt.Println(s)
}

func TestRandomIntMN(t *testing.T) {
	for idx := 0; idx < 50; idx++ {
		val := RandomIntMN(100, 101)
		if !(val == 101 || val == 100) {
			t.Fatal("not 100 or 101 val = ", val)
		}
	}
}

func TestDate(t *testing.T) {
	d := time.Now()
	str := d.Format("20060102150405")
	fmt.Println(str)
}
func TestGeneratorFileRangeSize3(t *testing.T) {

	create, err := os.Create("D:\\my\\project\\go\\go-example\\file\\t1.txt")
	if err != nil {
		return
	}
	defer create.Close()
	writer := bufio.NewWriter(create)
	writer.WriteString("abc")
	writer.WriteString("efs")
	writer.Flush()
}
