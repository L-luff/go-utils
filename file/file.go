package main

import (
	"bufio"
	"flag"
	"fmt"
	"go-utils/random"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	small = iota
	medium
	big
)

const (
	KB = iota
	MB
	GB
)

const (
	Separator = os.PathSeparator // 路径分隔符（分隔路径元素）
)

const (
	KBP = 10
)

const GOT = 2000

var typeMap = map[int]int{KB: small, MB: medium, GB: big}

func init() {
	rand.Seed(time.Now().Unix())
}

// -ms 100 -mx 200 -u kb -c 100 -p filePath

func main() {
	var (
		ms int
		mx int
		u  int
		c  int
		p  string
		t  int
	)
	flag.IntVar(&ms, "ms", 1, "file min size")
	flag.IntVar(&mx, "mx", 1024, "file max size ")
	flag.IntVar(&u, "u", KB, "file size units,0:kb 1:mb 2:gb please type 0-2 , default 0")
	flag.IntVar(&c, "c", 1, "file count")
	flag.IntVar(&t, "t", 0, "type")
	flag.StringVar(&p, "p", "", "file path")
	flag.Parse()
	fmt.Printf("file path:%s,minSize %v,maxSize %v\n", p, ms, mx)
	if len([]rune(p)) == 0 {
		panic("please type file path")
	}
	if ms > mx {
		panic("minSize greater than maxSize,please type correct size")
	}
	if t == 0 {
		CreateFile(ms, mx, u, c, p)
	} else {
		CreateFile2(ms, mx, u, c, p)
	}
}

func CreateFile(ms int, mx int, u int, c int, p string) {
	var wg sync.WaitGroup
	t := c / GOT
	tr := c % GOT
	startTime := time.Now()
	for idx := 0; idx < GOT; idx++ {
		if t == 0 && tr == 0 {
			break
		}
		count := t
		wg.Add(1)
		if tr > 0 {
			tr--
			count++
		}
		go GeneratorFileBatch(count, &wg, p, ms, mx, u)
	}
	wg.Wait()
	times := time.Since(startTime).Seconds()
	fmt.Println(fmt.Sprintf("spend time %v s\n", times))
}

func CreateFile2(ms int, mx int, u int, c int, p string) {
	var wg sync.WaitGroup
	startTime := time.Now()
	for ; c > 0; c-- {
		wg.Add(1)
		go GeneratorFileBatch(1, &wg, p, ms, mx, u)
	}
	wg.Wait()
	times := time.Since(startTime).Seconds()
	fmt.Println(fmt.Sprintf("spend time %vs\n", times))
}

func GeneratorFileBatch(count int, wg *sync.WaitGroup, path string, minSize int, maxSize int, uints int) error {
	defer wg.Done()
	var err error
	for ; count > 0; count-- {
		// ignore error
		errs := GeneratorFileRangeSize(path, minSize, maxSize, uints)
		if errs != nil {
			fmt.Println(err)
			err = errs
		}
	}
	return err
}

func GeneratorFileRangeSize(path string, minSize int, maxSize int, units int) error {
	if !strings.HasSuffix(path, strconv.Itoa(Separator)) {
		path = path + string(Separator)
	}
	if exists, err := DirExits(path); err != nil || !exists {
		if err = os.MkdirAll(path, os.ModePerm); err != nil {
			return fmt.Errorf("create dir %s error %v", path, err)
		}
	}
	filePath := path + generatorFileName()
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	switch units {
	case KB:
		writeFileDataKB(file, minSize, maxSize)
	case MB:
		writeFileDataMB(file, minSize, maxSize)
	case GB:
		writeFileDataGB(file, minSize, maxSize)
	default:
		return fmt.Errorf("not support")
	}
	return nil
}

func DirExits(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, err
	}
	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func writeFileDataKB(file *os.File, minSize int, maxSize int) {
	size := random.RandomIntMN(minSize, maxSize)
	//fmt.Printf("write file size %v kb\n",size)
	writeKB(file, size)
}

func writeKB(file *os.File, size int) {
	size = size * 1024
	perSize := size / KBP
	remainSize := size % KBP
	writer := bufio.NewWriter(file)
	for idx := 0; idx < KBP; idx++ {
		writer.WriteString(random.RandomStr(perSize))
	}
	if remainSize > 0 {
		writer.WriteString(random.RandomStr(remainSize))
	}
	writer.Flush()
}
func writeFileDataMB(file *os.File, minSize int, maxSize int) {
	size := random.RandomIntMN(minSize, maxSize)
	size = size * 1024
	writeKB(file, size)
}

func writeFileDataGB(file *os.File, minSize int, maxSize int) {
	size := random.RandomIntMN(minSize, maxSize)
	size = size * 1024 * 1024
	writeKB(file, size)
}

//
//func GeneratorFile(path string,types int) error  {
//	var minSize int= 1
//	var maxSize int= 1024
//	return GeneratorFileRangeSize(path,minSize,maxSize,types)
//}

func generatorFileName() string {
	return random.RandomStr(10) + "_" + time.Now().Format("20060102150405")
}
