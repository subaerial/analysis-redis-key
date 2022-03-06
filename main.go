package main

import (
	"analysis.redis/config"
	operate "analysis.redis/csv"
	"analysis.redis/mail"
	"analysis.redis/sqlite"
	"flag"
	"fmt"
	"os"
)

func init() {
	cmd()
	config.InitProperties()
	sqlite.InitDB()
}

func main() {
	runtime := operate.StartAnalysis(csv)
	mail.SendEndEmail(csv, runtime)
}

var csv string

func cmd() {
	csvUsage := "csv filename, eg: -file=2022-02-22_10.100.2.33_7001.csv or -file=/opt/2022-02-22_10.100.2.33_7001.csv"
	flag.StringVar(&csv, "file", "", csvUsage)
	flag.Parse()
	if csv == "" {
		fmt.Println("input", csvUsage)
		os.Exit(0)
	}
}
