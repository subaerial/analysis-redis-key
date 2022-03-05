package sqlite

func dropRedisKeyTable() {
	sql := `drop table if exists redis_keys;
			drop table if exists redis_key_prefixes;`
	connect.Exec(sql)
}

// createRedisKeyTable 创建redis key存储表
func createRedisKeyTable() {
	sql := `create table if not exists redis_keys
			(
				id                  INTEGER primary key AUTOINCREMENT,
				db                  INTEGER,
				key_type            text,
				key                 text,
				size_in_byte        INTEGER,
				num_elements        INTEGER,
				len_largest_element INTEGER,
				expire              TEXT
			);`
	connect.Exec(sql)
}

// createRedisKeyPrefixTable 创建redis key prefix分析存储表
func createRedisKeyPrefixTable() {
	sql := `create table if not exists redis_key_prefixes
			(
				id                  INTEGER primary key AUTOINCREMENT,
				db                  INTEGER,
				key_type            text,
				key                 text,
				size_in_byte        INTEGER,
				num_elements        INTEGER,
				len_largest_element INTEGER,
				expire              TEXT,
                count               INTEGER
			);`
	connect.Exec(sql)
}
