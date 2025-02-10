package config

import (
	"database/sql"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	_ "github.com/go-sql-driver/mysql"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsnRoot := "root:@tcp(127.0.0.1:3306)/"
	sqlDB, err := sql.Open("mysql", dsnRoot)
	if err != nil {
		log.Fatal("Koneksi ke MySQL gagal:", err)
	}
	defer sqlDB.Close()

	_, err = sqlDB.Exec("CREATE DATABASE IF NOT EXISTS article")
	if err != nil {
		log.Fatal("Gagal membuat database:", err)
	}

	fmt.Println("✅ Database 'article' berhasil dibuat atau sudah ada!")

	dsn := "root:@tcp(127.0.0.1:3306)/article?charset=utf8mb4&parseTime=True&loc=Local"
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Koneksi ke database gagal:", err)
	}

	DB = database
	log.Println("✅ Database berhasil terkoneksi!")
}
