package main

import (
	"analysis.redis/csv"
	"analysis.redis/sqlite"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	cmd()
	start := time.Now().UnixMilli()
	sqlite.InitDB()
	operate.ReadFile(csv)
	end := time.Now().UnixMilli()
	log.Printf("key分析完成, 共耗时: %d s", (end-start)/1000)
}

var csv string

func cmd() {
	csvUsage := "csv filename, eg: 2022-02-22_10.100.2.33_7001.csv or /opt/2022-02-22_10.100.2.33_7001.csv"
	flag.StringVar(&csv, "file", "csv", csvUsage)
	flag.Parse()
	if csv == "" {
		fmt.Println("input", csvUsage)
		os.Exit(0)
	}
}
