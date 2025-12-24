package database

import (
	"fmt"
	"log"
	"os"
	"procurement-system/models"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	// 1. Load file .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// 2. Ambil data dari .env
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// 3. Susun DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)

	// 4. Buka Koneksi
	connection, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panic("Gagal koneksi database:", err)
	}

	DB = connection

	// 5. Auto Migrate (Membuat Tabel Otomatis)
	// Pastikan file models/setup.go sudah Anda buat seperti pesan saya sebelumnya
	log.Println("Migrating Database...")
	connection.AutoMigrate(
		&models.User{},
		&models.Supplier{},
		&models.Item{},
		&models.Purchasing{},
		&models.PurchasingDetail{},
	)

	// --- SEEDING DATA ---
	// Cek apakah sudah ada user di database?
	var userCount int64
	connection.Model(&models.User{}).Count(&userCount)

	// Jika kosong, buat user default
	if userCount == 0 {
		password, _ := bcrypt.GenerateFromPassword([]byte("password123"), 14)

		dummyUser := models.User{
			Username: "admin",
			Password: string(password),
			Role:     "admin",
		}
		connection.Create(&dummyUser)
		log.Println("Data Seeding: User 'admin' dengan password 'password123' berhasil dibuat!")
	}

	log.Println("Database Connected & Migrated!")
}
