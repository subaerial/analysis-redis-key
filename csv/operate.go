package operate

import (
	"analysis.redis/model"
	"analysis.redis/sqlite"
	"bufio"
	"io"
	"log"
	"os"
	"strings"
	"sync"
)

const SapComma = ","

// ReadFile 读取csv文件内容, 并手动分割
// path 文件路径
func ReadFile(path string) {
	opencsv, err := os.Open(path)
	if err != nil {
		log.Fatalf("csv文件读取失败: %s", path)
	}
	defer func(reader *os.File) {
		if err := reader.Close(); err != nil {
			log.Fatalf("csv文件关闭失败: %s", path)
		}
	}(opencsv)

	readCSVByLine(opencsv)
}

// readCSVByLine 按行读取csv文件
func readCSVByLine(opencsv *os.File) {
	reader := bufio.NewReader(opencsv)
	count := 0
	limit := -1

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if (count == 0 && limit == -1) || err != nil {
			limit = len(strings.Split(line, SapComma))
			continue
		}
		count++

		split := strings.SplitN(line, SapComma, limit)
		if len(split) < limit-1 {
			continue
		}
		// 将分割后的字符进行字段对应, 并进行存储
		keyInfo := parsingField(split)
		saveBigKeyInfo2Db(keyInfo)
		go analysisKeyPrefix(keyInfo)
	}

	keySet.Range(func(key, value interface{}) bool {
		prefix := value.(*model.RedisKeyPrefix)
		if isBigKey(prefix.KeyInfo) || prefix.Count > 40 {
			sqlite.InsertRedisKeyPrefix(prefix)
		}
		return true
	})

	log.Printf("CSV文件读取完成, 共计: %d 行", count)
}

// saveBigKeyInfo2Db 保存bigKey信息到sqlite
func saveBigKeyInfo2Db(keyInfo *model.RedisKeyInfo) {
	if isBigKey(keyInfo) {
		sqlite.InsertRedisKey(keyInfo)
	}
}

var keySet = sync.Map{}

// 保存分析后的key前缀到map中
func analysisKeyPrefix(keyInfo *model.RedisKeyInfo) {
	keyInfo.Key = analysisRedisKey(keyInfo.Key)
	info, _ := keySet.Load(keyInfo.Key)
	if info == nil {
		info = &model.RedisKeyPrefix{
			KeyInfo: &model.RedisKeyInfo{
				Db:                keyInfo.Db,
				KeyType:           keyInfo.KeyType,
				Key:               keyInfo.Key,
				SizeInByte:        keyInfo.SizeInByte,
				NumElements:       keyInfo.NumElements,
				LenLargestElement: keyInfo.LenLargestElement,
				Expire:            keyInfo.Expire,
			},
			Count: 1,
		}
	} else {
		temp := info.(*model.RedisKeyPrefix)
		temp.Count++
		temp.KeyInfo = &model.RedisKeyInfo{
			Db:                keyInfo.Db,
			KeyType:           keyInfo.KeyType,
			Key:               keyInfo.Key,
			SizeInByte:        temp.KeyInfo.SizeInByte + keyInfo.SizeInByte,
			NumElements:       temp.KeyInfo.NumElements + keyInfo.NumElements,
			LenLargestElement: temp.KeyInfo.LenLargestElement + keyInfo.LenLargestElement,
			Expire:            keyInfo.Expire,
		}
		info = temp
	}
	keySet.Store(keyInfo.Key, info)
}

// key的阈值为10240，
//也就是对于string类型的value大于10240的认为是大key，
//对于list的话如果list长度大于10240认为是大key，
//对于hash的话如果field的数目大于10240认为是大key
func isBigKey(keyInfo *model.RedisKeyInfo) bool {
	length := int32(0)
	if keyInfo.KeyType == "string" {
		length = keyInfo.SizeInByte
	} else if keyInfo.KeyType == "hash" {
		length = keyInfo.NumElements
	} else if keyInfo.KeyType == "list" {
		length = keyInfo.NumElements
	} else if keyInfo.KeyType == "set" || keyInfo.KeyType == "sortedset" {
		length = keyInfo.NumElements
	} else if keyInfo.KeyType == "zset" {
		length = keyInfo.NumElements
	}
	return length > 10240
}
