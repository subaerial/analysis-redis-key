package main

import operate "analysis.redis/util"

func main() {
	filename := "2022-02-22_10.100.2.33_7001.csv"
	path := "/opt/tpapp/" + filename
	operate.ReadFile(path)
}
