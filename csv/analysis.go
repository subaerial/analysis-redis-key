package operate

import (
	"analysis.redis/model"
	"regexp"
	"strconv"
	"strings"
)

// 对应每行字段
// 1. big key 不对key进行单独分析, 若长度过大10240则直接返回
// 2. 分析后的key, 保存内存占用最大的top10, 以及总量最大的top10
func parsingField(line []string) *model.RedisKeyInfo {
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

	info := &model.RedisKeyInfo{
		Db:                int32(db),
		KeyType:           line[1],
		Key:               strings.TrimSpace(line[2]),
		SizeInByte:        int32(sizeInByte),
		NumElements:       int32(numElements),
		LenLargestElement: int32(lenLargestElement),
		Expire:            isExpire,
	}
	return info
}

// 分析redis key前缀
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
