package sqlite

import (
	"analysis.redis/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

var connect *gorm.DB

// InitDB 初始化数据库
func InitDB() {
	connectDb()
	// 先清空表, 保证每次值存储当天分析数据
	dropRedisKeyTable()
	// 创建数据表
	createRedisKeyTable()
	createRedisKeyPrefixTable()
}

// connect 创建数据库链接
func connectDb() {
	var err error
	if connect == nil {
		connect, err = gorm.Open(sqlite.Open("data/bigkey.db"), &gorm.Config{})
	}
	if err != nil {
		panic("failed to connect database")
	}
}

// InsertRedisKey 新增数据
func InsertRedisKey(keyInfo *model.RedisKeyInfo) {
	info := model.RedisKey{
		KeyInfo: keyInfo,
	}
	key := SelectOneByKey(info.KeyInfo.Key)
	if key.KeyInfo == nil {
		connect.Create(&info)
		log.Printf("%s key 已入库", keyInfo.Key)
	} else {
		log.Printf("%s key 已存在", keyInfo.Key)
	}

}

// SelectOneByKey 根据key查找
func SelectOneByKey(key string) *model.RedisKey {
	var info *model.RedisKey
	connect.Where("key = ?", key).First(&info)
	return info
}

func SelectTop100BigKeyByMemory(isExpire bool) *[]model.RedisKey {
	var infos *[]model.RedisKey
	connect.Where("expire = ?", isExpire).Order("size_in_byte desc").Limit(100).Find(&infos)
	return infos
}

func SelectTop100BigKeyByLen(isExpire bool) *[]model.RedisKey {
	var infos *[]model.RedisKey
	connect.Where("expire = ?", isExpire).Order("num_elements desc").Limit(100).Find(&infos)
	return infos
}

// InsertRedisKeyPrefix 新增数据
func InsertRedisKeyPrefix(keyInfo *model.RedisKeyPrefix) {
	key := SelectOneByKeyPrefix(keyInfo.KeyInfo.Key)
	if key.KeyInfo == nil {
		connect.Create(&keyInfo)
		log.Printf("%s key-prefix 已入库", keyInfo.KeyInfo.Key)
	} else {
		log.Printf("%s key-prefix 已存在", keyInfo.KeyInfo.Key)
	}
}

// SelectOneByKeyPrefix 根据key查找
func SelectOneByKeyPrefix(key string) *model.RedisKeyPrefix {
	var info *model.RedisKeyPrefix
	connect.Where("key = ?", key).First(&info)
	return info
}
