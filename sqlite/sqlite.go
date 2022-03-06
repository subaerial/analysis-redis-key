package sqlite

import (
	"analysis.redis/config"
	"analysis.redis/mail"
	"analysis.redis/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"path/filepath"
)

var connect *gorm.DB

// db url
var dbUrl string

// InitDB 初始化数据库
func InitDB() {
	dbUrl = config.Properties.Database.Url
	createSqliteFile()
	connectDb()
	// 先清空表, 保证每次值存储当天分析数据
	dropRedisKeyTable()
	// 创建数据表
	createRedisKeyTable()
	createRedisKeyPrefixTable()
}

// 创建sqlite文件
func createSqliteFile() {
	if _, err2 := os.Stat(dbUrl); err2 == nil {
		log.Println("sqlite数据库文件已存在")
		return
	}

	dir, _ := filepath.Split(dbUrl)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		mail.SendErrorEmail(err)
		log.Fatalf("dir: %s 创建失败", dir)
	}

	if _, err := os.Create(dbUrl); err != nil {
		mail.SendErrorEmail(err)
		log.Fatalf("sqlite文件创建失败")
	}
}

// connect 创建数据库链接
func connectDb() {
	var err error
	if connect == nil {
		connect, err = gorm.Open(sqlite.Open(dbUrl), &gorm.Config{})
	}
	if err != nil {
		mail.SendErrorEmail(err)
		panic("failed to connect database")
	}
}

// InsertRedisKey 新增数据
func InsertRedisKey(keyInfo *model.RedisKeyInfo) {
	info := model.RedisKey{KeyInfo: keyInfo}
	if key := SelectOneByKey(info.KeyInfo.Key); key.KeyInfo == nil {
		connect.Create(&info)
		log.Printf("%s key 已入库", keyInfo.Key)
	}
}

// SelectOneByKey 根据key查找
func SelectOneByKey(key string) *model.RedisKey {
	var info *model.RedisKey
	connect.Where("key = ?", key).First(&info)
	return info
}

// InsertRedisKeyPrefix 新增数据
func InsertRedisKeyPrefix(keyInfo *model.RedisKeyPrefix) {
	if key := SelectOneByKeyPrefix(keyInfo.KeyInfo.Key); key.KeyInfo == nil {
		connect.Create(&keyInfo)
		log.Printf("%s key-prefix 已入库", keyInfo.KeyInfo.Key)
	}
}

// SelectOneByKeyPrefix 根据key查找
func SelectOneByKeyPrefix(key string) *model.RedisKeyPrefix {
	var info *model.RedisKeyPrefix
	connect.Where("key = ?", key).First(&info)
	return info
}

// SelectTop100BigKeyByMemory 查找内存占用TOP100的bigKey
func SelectTop100BigKeyByMemory(isExpire bool) *[]model.RedisKey {
	var infos *[]model.RedisKey
	connect.Where("expire = ?", isExpire).Order("size_in_byte desc").Limit(config.Properties.Redis.BigKey.Memory).Find(&infos)
	return infos
}

// SelectTop100BigKeyByLen 查找长度占用TOP100的bigKey
func SelectTop100BigKeyByLen(isExpire bool) *[]model.RedisKey {
	var infos *[]model.RedisKey
	connect.Where("expire = ?", isExpire).Order("num_elements desc").Limit(config.Properties.Redis.BigKey.Len).Find(&infos)
	return infos
}
