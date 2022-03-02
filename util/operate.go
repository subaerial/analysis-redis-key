package operate

import (
	"bufio"
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

var keySet = make(map[string]int32)

func readCSVByLine(opencsv *os.File) {
	reader := bufio.NewReader(opencsv)
	count := 0
	limit := -1

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
		prefix := parsingField(split)
		num := keySet[prefix]
		if num == 0 {
			keySet[prefix] = 1
		} else {
			keySet[prefix]++
		}
	}

	log.Printf("CSV文件读取完成, 共计: %d 行", count)
}
