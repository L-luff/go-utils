package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"github.com/L-luff/go-utils/file"
	"github.com/L-luff/go-utils/log"
	"strconv"
	"strings"
	"time"
)

// -o 1 -ms 100 -mx 200 -u kb -c 100 -p dirPath
// -o 0 -p dirPath -dc 20,20
// -o 2 -p dirPath
// -o 3 -p dirPath -r true -c 100
// -o 4  -ms 100 -mx 200 -u kb -c 100 -p dirPath
// -o 5 -p dirPath
// -o 6 -p dirPath
// -0 7 -p dirPath(or filePath) -c 100 -fc
func main() {
	var (
		ms              int
		mx              int
		u               int
		c               int
		p               string
		t               int
		o               int // 操作类型
		recursive       bool
		depthsCountVar  string
		logLevel        int  //  是否debug
		forceConsitency bool // 完全一致性，会检查目录
	)

	flag.IntVar(&ms, "ms", 1, "file min size")
	flag.IntVar(&mx, "mx", 1024, "file max size ")
	flag.IntVar(&u, "u", file.KB, "file size units,0:kb 1:mb 2:gb please type 0-2 , default 0")
	flag.IntVar(&c, "c", 1, "file count")
	flag.IntVar(&t, "t", 0, "type")
	flag.StringVar(&p, "p", "", "file path")
	flag.IntVar(&o, "o", 1, "operation type, 0: CREATE_DIR,1:CREATE_FILE 2:COUNT_DIR 3: LIST_FILE,4:RANDOM_CREATE_FILE 5:HASH 6:COUNT_FILE 7: DELETE_FILE ")
	flag.BoolVar(&recursive, "r", false, "recursive")
	flag.StringVar(&depthsCountVar, "dc", "1", "dir depths count")
	flag.IntVar(&logLevel, "l", 1, "log level")
	flag.BoolVar(&forceConsitency, "fc", false, "true: check dir name and file content consistency or delete dir  flase: just check file content")

	flag.Parse()
	if len([]rune(p)) == 0 {
		panic("please type  path")
	}
	log.SetLevel(logLevel)
	switch o {
	case file.CREATE_DIR:
		splitS := strings.Split(depthsCountVar, ",")
		depthsCount := make([]int, 0)
		for i := 0; i < len(splitS); i++ {
			val, err := strconv.Atoi(splitS[i])
			if err != nil {
				panic(err)
			}
			depthsCount = append(depthsCount, val)
		}
		log.Debug("depths count val:", depthsCount)
		err := file.CreateDir(p, depthsCount, true)
		if err != nil {
			panic(err)
		}
	case file.CREATE_FILE:
		if ms > mx {
			panic("minSize greater than maxSize,please type correct size")
		}
		log.Debugf("file path:%s,minSize %v,maxSize %v\n", p, ms, mx)
		startTime := time.Now()
		if t == 0 {
			file.CreateFile(ms, mx, u, c, p)
		} else {
			file.CreateFile2(ms, mx, u, c, p)
		}
		times := time.Since(startTime).Seconds()
		fmt.Println(times)
	case file.COUNT_DIR:
		count, err := file.CountOfDir(p, recursive)
		if err != nil {
			panic(err)
		}
		fmt.Printf("path : %s,dir count:%d\n", p, count)
	case file.RANDOM_FILE_WRITE:
		if c <= 0 {
			panic("please type correct update file count")
		}
		file.RandomUpdateFilesOnDir(p, c, recursive)
	case file.RANDOM_CREATE_FILE:
		err := file.RandomCreateFile(ms, mx, u, c, p)
		if err != nil {
			panic(err)
		}
	case file.HASH:
		hash, err := file.Hash(p, forceConsitency)
		if err != nil {
			panic(err)
		}
		val, err := binary.ReadUvarint(bytes.NewBuffer(hash))
		if err != nil {
			panic(err)
		}
		fmt.Println(val)
	case file.COUNT_FILE:
		count, err := file.CountOfFile(p, recursive)
		if err != nil {
			panic(err)
		}
		fmt.Println(count)
	case file.DELETE_FILE:
		startTime := time.Now()
		err := file.Delete(p, c, forceConsitency)
		if err != nil {
			panic(err)
		}
		fmt.Println(time.Since(startTime).Seconds())
	default:
		panic("not support")
	}
}
