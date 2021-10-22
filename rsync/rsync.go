package rsync

import "crypto/md5"

func WeakCheckSum(data []byte) uint32 {
	var high, low, mask uint32 = 0, 1, 65521
	for _, b := range data {
		low = (uint32(b) + low) % mask
		high = (high + low) % mask
	}
	return (high << 16) | low
}

func CheckAngGeneratorTransBlock(src string) {

}

func StrongCheckSum(data []byte) []byte {
	hash := md5.New()
	hash.Write(data)
	return hash.Sum(nil)
}

type Block struct {
	BlockIndex     uint
	BlockSize      uint
	WeakCheckSum   uint32
	StrongCheckSum []byte
}

type Trans struct {
	BlockHint  bool
	BlockIndex uint
	data       []byte
}
