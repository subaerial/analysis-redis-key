package operate

import (
	"regexp"
	"strings"
)

// 对应每行字段
func parsingField(line []string) string {
	if line == nil {
		return ""
	}
	// csv中数据分割后每个值对应的含义
	//keyType := line[1]
	key := strings.TrimSpace(line[2])
	//sizeInByte, _ := strconv.Atoi(line[3])
	//numElements, _ := strconv.Atoi(line[5])
	//lenLargestElement, _ := strconv.Atoi(line[6])
	//expire := strings.Trim(line[7], " \t\n\r\v\f")
	//
	//isExpire := false
	//if expire != "" || len(expire) > 0 {
	//	isExpire = true
	//}

	return analysisRedisKey(key)
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
