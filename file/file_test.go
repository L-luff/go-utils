package file

import (
	"bufio"
	"bytes"
	"encoding/binary"
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

func TestRegularFile(t *testing.T) {
	filePath := "D:\\my\\project\\self\\go-utils\\file\\a"
	stat, err := os.Stat(filePath)
	if err != nil {
		t.Fatal(err)
	}
	if !stat.Mode().IsRegular() {
		t.Fatal("not a regular file")
	}
}

func TestRandomCreateFile(t *testing.T) {
	path := "D:\\my\\test_agent"
	err := RandomCreateFile(100, 100, 0, 1000, path)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFileHash(t *testing.T) {
	f1 := "D:\\my\\project\\self\\go-utils\\file\\file.go"
	hash1, err := FileHash(f1)
	if err != nil {
		t.Fatal(err)
	}
	f2 := "D:\\my\\project\\self\\go-utils\\file\\file2"
	hash2, err := FileHash(f2)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(hash1, hash2) {
		t.Fatalf("file1:%s content not equals file2:%s", f1, f2)
	}
}

func TestHash(t *testing.T) {
	f1 := "D:\\my\\project\\self\\go-utils\\file"
	hash1, err := Hash(f1, true)
	if err != nil {
		t.Fatal(err)
	}

	s1, err := binary.ReadUvarint(bytes.NewBuffer(hash1))
	if err != nil {
		t.Fatal(err)
	}
	f2 := "D:\\my\\project\\self\\go-utils\\file"
	hash2, err := Hash(f2, true)
	if err != nil {
		t.Fatal(err)
	}
	s2, err := binary.ReadUvarint(bytes.NewBuffer(hash2))
	fmt.Println(s1 == s2)
	if !bytes.Equal(hash1, hash2) {
		t.Fatalf("file1:%s content not equals file2:%s", f1, f2)
	}
}

func TestDelete(t *testing.T) {
	path := "D:\\my\\test_agent\\filesync1"
	num := 1
	fc := true
	err := Delete(path, num, fc)
	if err != nil {
		t.Fatal(err)
	}

}
