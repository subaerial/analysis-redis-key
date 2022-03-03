package operate

import (
	"log"
	"regexp"
	"strconv"
	"strings"
)

type RedisKeyInfo struct {
	// db 数据库默认0
	// keyType 类型: 基本的redis key类型: string、hash、list、set、zset
	// key 键
	// sizeInByte 使用的内存: 包括键、值和其他开销
	// numElements key中的value的个数
	// lenLargestElement key中的value的长度
	// expire 是否设置过期时间c
	db                int32
	keyType           string
	key               string
	sizeInByte        int32
	numElements       int32
	lenLargestElement int32
	expire            bool
}

// 对应每行字段
// 1. big key 不对key进行单独分析, 若长度过大10240则直接返回
// 2. 分析后的key, 保存内存占用最大的top10, 以及总量最大的top10
func parsingField(line []string) *RedisKeyInfo {
	log.Printf("%v", line)
	// csv中数据分割后每个值对应的含义
	db, _ := strconv.Atoi(line[0])
	sizeInByte, _ := strconv.Atoi(line[3])
	numElements, _ := strconv.Atoi(line[5])
	lenLargestElement, _ := strconv.Atoi(line[6])
	expire := strings.Trim(line[7], " \t\n\r\v\f")
	isExpire := false
	if expire != "" || len(expire) > 0 {
		isExpire = true
	}

	info := &RedisKeyInfo{
		db:                int32(db),
		keyType:           line[1],
		key:               strings.TrimSpace(line[2]),
		sizeInByte:        int32(sizeInByte),
		numElements:       int32(numElements),
		lenLargestElement: int32(lenLargestElement),
		expire:            isExpire,
	}
	return info
}

// 第一次分析, 分析出较宽泛的key
func analysisRedisKey(key string) string {
	separator, arrays := getBestMatchSeparator(key)
	prefix := arrays[0] + separator
	for i, arr := range arrays[1:] {
		matched, _ := regexp.MatchString(".*[0-9]+.*", arr)
		if !matched && i < len(arrays) {
			prefix += arr + separator
		}
	}
	return prefix
}

// redis key 的拼接符
var separators = []string{":", "-", "_"}

// 获取最佳的拼接符
func getBestMatchSeparator(key string) (string, []string) {
	// 最佳匹配分隔符
	bestSeparator := ""
	matchNums := -1
	for _, separator := range separators {
		size := len(strings.Split(key, separator))
		if size > 1 && size > matchNums && bestSeparator != ":" {
			bestSeparator = separator
			matchNums = size
		}
	}
	return bestSeparator, strings.Split(key, bestSeparator)
}
