package file

import (
	"bufio"
	"container/list"
	"crypto/sha256"
	"fmt"
	log "github.com/L-luff/go-utils/log"
	"github.com/L-luff/go-utils/random"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"runtime"
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
const (
	CREATE_DIR = iota
	CREATE_FILE
	COUNT_DIR
	RANDOM_FILE_WRITE
	RANDOM_CREATE_FILE
	HASH
	COUNT_FILE
	DELETE_FILE
)

var GOT = runtime.NumCPU() * 4

var typeMap = map[int]int{KB: small, MB: medium, GB: big}

func init() {
	rand.Seed(time.Now().Unix())
}

func Delete(path string, num int, forces bool) error {
	stat, err := os.Stat(path)
	if err != nil {
		return err
	}
	dirType := stat.IsDir()
	if !dirType {
		//delete file
		return DeleteFile(path)
	}
	if forces {
		return os.Remove(path)
	}
	// random delete files
	files, err := ListFile(path, true)
	if err != nil {
		return err
	}
	size := len(files)
	num = int(math.Min(float64(num), float64(size)))
	idx := random.RandomIntM(size)
	for ; num > 0; num-- {
		err := DeleteFile(files[idx%size])
		if err != nil {
			log.Debugf("delete file %s error %v", files[idx%size], err)
		}
		idx = (idx + 1) % size
	}
	return nil
}

func DeleteFile(file string) error {
	if !IsFile(file) {
		return fmt.Errorf("not a file %s", file)
	}
	return os.Remove(file)
}

// 如果path，计算文件hash,如果是目录，计算目录下的所有文件hash+hash(目录数量）

func Hash(path string, forceConsitency bool) ([]byte, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	dirType := stat.IsDir()
	if !dirType {
		return FileHash(path)
	}
	files, dirs, err := filesOfDir(path, true)
	if err != nil {
		return nil, err
	}
	h := sha256.New()
	for _, file := range files {
		fh, err := FileHash(file)
		if err != nil {
			return nil, err
		}
		_, err = h.Write(fh)
		if err != nil && io.EOF != err {
			return nil, err
		}
	}
	// dir name
	if forceConsitency {
		for _, dir := range dirs {
			_, err = h.Write([]byte(dir))
			if err != nil && err != io.EOF {
				return nil, err
			}
		}
	}

	return h.Sum(nil), nil
}

func FileHash(path string) ([]byte, error) {
	if !IsFile(path) {
		return nil, fmt.Errorf("not a file")
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

// 在目录p下任意选择目录进行文件写入

func RandomCreateFile(ms int, mx int, u int, c int, p string) error {

	_, dirs, err := filesOfDir(p, true)
	if err != nil {
		return err
	}
	dirs = append(dirs, p)
	dirSize := len(dirs)
	for i := 0; i < c; i++ {

		CreateFile(ms, mx, u, 1, dirs[random.RandomIntM(dirSize)])
	}
	return nil
}

//计算目录数目

func CountOfDir(dir string, recursive bool) (int, error) {
	_, dirs, err := filesOfDir(dir, recursive)
	if err != nil {
		return 0, err
	}
	return len(dirs), nil
}

func CountOfFile(dir string, recursive bool) (int, error) {
	files, _, err := filesOfDir(dir, recursive)
	if err != nil {
		return 0, err
	}
	return len(files), nil
}

// 返回该目录下的所有文件
// dir:目录
// recursive: 是否递归

func ListFile(dir string, recursive bool) ([]string, error) {
	if !IsDir(dir) {
		return nil, fmt.Errorf("path:%s is not a dir", dir)
	}

	files, _, err := filesOfDir(dir, recursive)
	if err != nil {
		return nil, err
	}
	return files, err
}

// 目录下的所有文件任意修改, 可写的正常文件

func RandomUpdateFilesOnDir(dir string, updateCount int, recursive bool) error {

	if !IsDir(dir) {
		return fmt.Errorf("please type correct dir path, %s not exits", dir)
	}

	files, err := ListFile(dir, recursive)
	if err != nil {
		return err
	}
	lens := len(files)
	if updateCount > lens {
		updateCount = lens
	}
	log.Debug("update file count is ", updateCount)
	idx := random.RandomIntM(lens)
	for ; updateCount > 0; updateCount-- {
		err = RandomUpdateFile(files[idx%lens])
		//just print
		if err != nil {
			fmt.Println(err)
		}
		idx++
	}
	return nil
}

func RandomUpdateFile(file string) error {
	filePoint, err := os.OpenFile(file, os.O_RDWR, 0766)
	if err != nil {
		return err
	}
	defer filePoint.Close()
	stat, err := os.Stat(file)
	if err != nil {
		return err
	}
	if !stat.Mode().IsRegular() {
		log.Debugf("file {} is not a regular file. do not write data", file)
		return nil
	}
	fileSize := stat.Size()
	var seekPosition int64 = 0
	if fileSize > 0 {
		seekPosition = rand.Int63n(fileSize)
	}
	_, err = filePoint.Seek(seekPosition, io.SeekStart)
	if err != nil {
		return err
	}
	// 随机写入1kb大小数据
	writeDataSize := 1024 * 1
	log.Debugf("start write 1kb data to %s on seek %d \n", file, seekPosition)
	dataStr := random.RandomStr(writeDataSize)
	writer := bufio.NewWriter(filePoint)
	writer.WriteString(dataStr)
	writer.Flush()
	return nil
}

// 创建目录
// dir : 当前所在目录创建子目录
// depthsCount: 每层数量
// globalSeq : 是否全局有序
// todo ignore error

func CreateDir(dir string, depthsCount []int, globalSeq bool) error {
	startTime := time.Now()
	if !strings.HasSuffix(dir, string(Separator)) {
		dir = dir + string(Separator)
	}
	if !IsDir(dir) {
		fmt.Errorf("dir path:%s is not a dir", dir)
	}

	dirNameFunc := func(seq int) string {
		suffixStr := "sub"
		return suffixStr + "_" + strconv.Itoa(seq) + string(Separator)
	}
	// 暂时默认全局有序
	stack := list.New()
	stack.PushBack(dir)
	suffixNumber := 1
	for idx := 0; idx < len(depthsCount); idx++ {
		// 第idx层
		stackLen := stack.Len()
		for i := 0; i < stackLen; i++ {
			removeDir := stack.Remove(stack.Front()).(string)
			for j := 0; j < depthsCount[idx]; j++ {
				childDir := removeDir + dirNameFunc(suffixNumber)
				os.Mkdir(childDir, os.ModePerm)
				suffixNumber++
				stack.PushBack(childDir)
			}
		}
	}

	log.Debugf("create dir success,last suffixNumber is %d \n", suffixNumber-1)
	fmt.Println(time.Since(startTime).Seconds())
	return nil
}

// files and dirs

func filesOfDir(dir string, recursive bool) ([]string, []string, error) {
	if !strings.HasSuffix(dir, string(Separator)) {
		dir = dir + string(Separator)
	}
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, nil, err
	}
	ans := make([]string, 0)
	ansDir := make([]string, 0)
	for idx := 0; idx < len(fileInfos); idx++ {
		if fileInfos[idx].IsDir() {
			ansDir = append(ansDir, dir+fileInfos[idx].Name())
			if !recursive {
				continue
			}
			tmpResFile, tmpResDir, err := filesOfDir(dir+fileInfos[idx].Name(), recursive)
			if err != nil {
				return nil, nil, err
			}
			ans = append(ans, tmpResFile...)
			ansDir = append(ansDir, tmpResDir...)
		} else {
			ans = append(ans, dir+fileInfos[idx].Name())
		}
	}
	return ans, ansDir, nil
}
func CreateFile(ms int, mx int, u int, c int, p string) {
	var wg sync.WaitGroup
	t := c / GOT
	tr := c % GOT
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

func IsDir(dir string) bool {
	f, err := os.Stat(dir)
	if err != nil {
		return false
	}
	return f.IsDir()
}

func IsFile(path string) bool {
	f, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !f.IsDir()
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
