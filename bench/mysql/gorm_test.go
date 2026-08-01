package mysql_test

import (
	"fmt"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

type Account struct {
	ID   int64  `gorm:"primaryKey"`
	Name string `gorm:"size:255"`
}

func initDB() *gorm.DB {
	// 主库（写）：走 6446，自动指向 Primary
	writeDSN := "root:Root@123456@tcp(127.0.0.1:6446)/fishes?charset=utf8mb4"

	// 从库（读）：走 6447，自动负载到 mysql-2 / mysql-3
	readDSN := "root:Root@123456@tcp(127.0.0.1:6447)/fishes?charset=utf8mb4"

	db, err := gorm.Open(mysql.Open(writeDSN), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// 注册读写分离插件，框架自动根据 SQL 类型选择端口
	db.Use(dbresolver.Register(dbresolver.Config{
		Sources:  []gorm.Dialector{mysql.Open(writeDSN)},
		Replicas: []gorm.Dialector{mysql.Open(readDSN)},
	}))

	return db
}

func BenchmarkGormReadWriteSeparation(b *testing.B) {
	// 1. Setup 只执行一次，不计入 benchmark 时间
	db := initDB()
	db.AutoMigrate(&Account{})

	// 预热：确保连接池已建立
	db.Create(&Account{Name: "warmup"})

	// 2. 真正的 benchmark 循环
	for i := 0; b.Loop(); i++ {
		// 写操作 -> 应该走 6446（Primary）
		db.Create(&Account{Name: fmt.Sprintf("user_%d", i)})

		// 读操作 -> 应该走 6447（Secondaries）
		var acc Account
		db.First(&acc, "name = ?", fmt.Sprintf("user_%d", i))
	}

	b.StopTimer() // 停止计时，开始清理
}
