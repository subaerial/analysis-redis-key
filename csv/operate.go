package operate

import (
	"bufio"
	"container/list"
	"io"
	"log"
	"os"
	"strings"
)

const SapComma = ","

// ReadFile 读取csv文件内容, 并手动分割
// path 文件路径
func ReadFile(path string) {
	opencsv, err := os.Open(path)
	if err != nil {
		log.Printf("csv文件读取失败: %s", path)
		return
	}

	defer func(reader *os.File) {
		err := reader.Close()
		if err != nil {
			log.Printf("csv文件流关闭失败: %s", path)
		}
	}(opencsv)

	readCSVByLine(opencsv)

}

func readCSVByLine(opencsv *os.File) {
	reader := bufio.NewReader(opencsv)
	count := 0
	limit := -1

	bigKeys := *list.New()
	for {
		line, err := reader.ReadString('\n')
		if count == 0 && limit == -1 {
			limit = len(strings.Split(line, SapComma))
			continue
		}
		count++
		if err == io.EOF {
			break
		} else if err != nil {
			log.Println(err)
			continue
		}
		split := strings.SplitN(line, SapComma, limit)
		if len(split) < limit-1 {
			continue
		}
		keyInfo := parsingField(split)
		if isBigKey(keyInfo) {
			bigKeys.PushBack(keyInfo)
		}
	}

	log.Printf("CSV文件读取完成, 共计: %d 行", count)
}

// key的阈值为10240，
//也就是对于string类型的value大于10240的认为是大key，
//对于list的话如果list长度大于10240认为是大key，
//对于hash的话如果field的数目大于10240认为是大key
func isBigKey(keyInfo *RedisKeyInfo) bool {
	length := int32(0)
	if keyInfo.keyType == "string" {
		length = keyInfo.sizeInByte
	} else if keyInfo.keyType == "hash" {
		length = keyInfo.numElements
	} else if keyInfo.keyType == "list" {
		length = keyInfo.numElements
	} else if keyInfo.keyType == "set" || keyInfo.keyType == "sortedset" {
		length = keyInfo.numElements
	} else if keyInfo.keyType == "zset" {
		length = keyInfo.numElements
	}
	return length > 10240
}
