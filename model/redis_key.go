package model

type RedisKeyInfo struct {
	// Db 数据库默认0
	// KeyType 类型: 基本的redis key类型: string、hash、list、set、zset
	// Key 键
	// SizeInByte 使用的内存: 包括键、值和其他开销
	// NumElements key中的value的个数
	// LenLargestElement key中的value的长度
	// Expire 是否设置过期时间c
	Db                int32
	KeyType           string
	Key               string
	SizeInByte        int32
	NumElements       int32
	LenLargestElement int32
	Expire            bool
}

type RedisKey struct {
	ID      int32         `gorm:"primaryKey"`
	KeyInfo *RedisKeyInfo `gorm:"embedded"`
}

type RedisKeyPrefix struct {
	ID      int32         `gorm:"primaryKey"`
	KeyInfo *RedisKeyInfo `gorm:"embedded"`
	Count   int32
}
