package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/Zauro25/Capstone-PerpusKominfosan/models"
)

func main() {
	_ = godotenv.Load()

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_USER", "admin"),
		getEnv("DB_PASSWORD", "perpustakaandpk123"),
		getEnv("DB_NAME", "perpustakaan_db"),
		getEnv("DB_PORT", "5432"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("gagal koneksi database: %v", err)
	}

	if err := db.AutoMigrate(&models.AdminDPK{}); err != nil {
		log.Fatalf("gagal migrasi tabel admin_dpk: %v", err)
	}

	username := "budi"
	plainPassword := "budi123"
	now := time.Now().Unix()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("gagal hash password: %v", err)
	}

	defaultData := models.AdminDPK{
		Username:    username,
		Password:    string(hashedPassword),
		NamaLengkap: "Budi Santoso",
		Email:       fmt.Sprintf("budi.admin.dpk.%d@sidapus.local", now),
		NoTelepon:   "081234567890",
		HakAkses:    "admin",
		IsActive:    true,
	}

	var existing models.AdminDPK
	err = db.Where("username = ?", username).First(&existing).Error
	if err == nil {
		existing.Password = string(hashedPassword)
		existing.IsActive = true
		if existing.NamaLengkap == "" {
			existing.NamaLengkap = defaultData.NamaLengkap
		}
		if existing.Email == "" {
			existing.Email = defaultData.Email
		}
		if existing.NoTelepon == "" {
			existing.NoTelepon = defaultData.NoTelepon
		}
		if existing.HakAkses == "" {
			existing.HakAkses = "admin"
		}

		if saveErr := db.Save(&existing).Error; saveErr != nil {
			log.Fatalf("gagal update admin DPK existing: %v", saveErr)
		}

		log.Printf("admin DPK sudah ada, password direset. username=%s", username)
		return
	}

	if err != gorm.ErrRecordNotFound {
		log.Fatalf("gagal cek admin DPK existing: %v", err)
	}

	if createErr := db.Create(&defaultData).Error; createErr != nil {
		log.Fatalf("gagal membuat admin DPK: %v", createErr)
	}

	log.Printf("admin DPK berhasil dibuat. username=%s", username)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
