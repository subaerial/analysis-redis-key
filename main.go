package main

import (
	"analysis.redis/csv"
	"log"
	"time"
)

func main() {
	start := time.Now().UnixMilli()

	filename := "2022-02-22_10.100.2.33_7001.csv"
	path := "/opt/tpapp/" + filename
	operate.ReadFile(path)
	end := time.Now().UnixMilli()
	log.Printf("key分析完成, 共耗时: %d s", (end-start)/1000)
}
