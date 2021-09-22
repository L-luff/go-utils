package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"testing"
	"time"
)

func TestCreateFile(t *testing.T) {
	path := "D:\\my\\test_agent\\filesync2\\tests"
	CreateFile(1, 1, 0, 100, path)
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

func TestListFile(t *testing.T) {
	path := "D:\\my\\test_agent\\filesync2"
	ans, err := ListFile(path, true)
	if err != nil {
		t.Fatal(err)
	}
	if ans == nil || len(ans) == 0 {
		t.Fatal("no files")
	}
	//fmt.Println(ans)
}

func TestCountOfDir(t *testing.T) {
	dir := "D:\\my\\test_agent\\"
	count, err := CountOfDir(dir, true)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("number of subdirectories:", count)
}

func TestWriteFile(t *testing.T) {
	filePath := "D:\\my\\test_agent\\a.txt"
	file, err := os.OpenFile(filePath, os.O_RDWR, 0666)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	// 设置了seek,之后写入文件时，这个这个seek的偏移量来写
	//
	seek, err := file.Seek(2, io.SeekEnd)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("seek position : %v \n", seek)
	writer := bufio.NewWriter(file)
	writer.WriteString("test--")
	writer.Flush()
}

func TestFileSize(t *testing.T) {
	filePath := "D:\\my\\test_agent\\a.txt"
	stat, err := os.Stat(filePath)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("file size : ", stat.Size())
}

func TestCreateDir(t *testing.T) {
	depthsCount := []int{20, 20}
	CreateDir("D:\\my\\test_agent\\filesync1", depthsCount, false)
}

func TestRandomUpdateFile(t *testing.T) {
	filePath := "D:\\my\\test_agent\\a.txt"
	RandomUpdateFile(filePath)
}

func TestMatht(t *testing.T) {
	n := 3
	k := 7
	val := (int((math.Pow(float64(n), float64(k+1)))-1) / (n - 1)) - 1
	fmt.Println(val)
	ans := 0
	for i := 1; i <= k; i++ {
		ans += int(math.Pow(float64(n), float64(i)))
	}
	fmt.Println(ans)
}
