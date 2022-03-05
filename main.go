package main

import (
	operate "analysis.redis/csv"
	"analysis.redis/mail"
	"analysis.redis/sqlite"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	cmd()
	start := time.Now().UnixMilli()
	sqlite.InitDB()
	operate.ReadFile(csv)
	end := time.Now().UnixMilli()
	runtime := (end - start) / 1000
	log.Printf("key分析完成, 共耗时: %d s", runtime)

	ip, _ := getLocalIP()
	body := fmt.Sprintf(`redis key全量分析脚本一执行结束, 服务器: %s , 共耗时: %ds`, ip, runtime)
	mail.SendEmail("moxfan@126.com", "wangpf09@126.com", "redis全量分析脚本执行通知", body)
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

// 获取本机网卡IP
func getLocalIP() (ipv4 string, err error) {
	var (
		addrs   []net.Addr
		addr    net.Addr
		ipNet   *net.IPNet // IP地址
		isIpNet bool
	)
	// 获取所有网卡
	if addrs, err = net.InterfaceAddrs(); err != nil {
		return
	}
	// 取第一个非lo的网卡IP
	for _, addr = range addrs {
		// 这个网络地址是IP地址: ipv4, ipv6
		if ipNet, isIpNet = addr.(*net.IPNet); isIpNet && !ipNet.IP.IsLoopback() {
			// 跳过IPV6
			if ipNet.IP.To4() != nil {
				ipv4 = ipNet.IP.String() // 192.168.1.1
				return
			}
		}
	}
	hostname, _ := os.Hostname()
	return hostname, err
}
