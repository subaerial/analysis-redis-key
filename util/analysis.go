package operate

import (
	"log"
	"strings"
)

var separators = []string{":", "-", "_"}

// 对应每行字段
func parsingField(line []string) {
	if line == nil {
		return
	}
	// csv中数据分割后每个值对应的含义
	//keyType := line[1]
	key := line[2]
	//sizeInByte, _ := strconv.Atoi(line[3])
	//numElements, _ := strconv.Atoi(line[5])
	//lenLargestElement, _ := strconv.Atoi(line[6])
	//expire := strings.Trim(line[7], " \t\n\r\v\f")
	//
	//isExpire := false
	//if expire != "" || len(expire) > 0 {
	//	isExpire = true
	//}

	//println(keyType, key, sizeInByte, numElements, lenLargestElement, expire, isExpire)
	analysisRedisKey(key)
}

func analysisRedisKey(key string) string {
	for _, separator := range separators {
		split := strings.Split(key, separator)
		log.Printf("%v", split)
	}

	return ""
}
