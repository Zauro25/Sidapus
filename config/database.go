package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Zauro25/Capstone-PerpusKominfosan/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error

	// Database connection string
	dsn := buildDSN(
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_USER", "admin"),
		getEnv("DB_PASSWORD", "perpustakaandpk123"),
		getEnv("DB_NAME", "perpustakaan_db"),
		getEnv("DB_PORT", "5432"),
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connection established")

	if getEnv("AUTO_MIGRATE", "true") == "true" {
		if err := runAutoMigrate(); err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "permission denied for schema") {
				log.Fatalf("Failed to migrate database: %v\nHint: set MIGRATION_DB_USER/MIGRATION_DB_PASSWORD with schema create privileges, or grant CREATE on schema public to DB_USER.", err)
			}
			log.Fatal("Failed to migrate database:", err)
		}
		log.Println("Database migrated successfully")
	}

	log.Println("Database initialization completed")
}

func AutoMigrateDB() error {
	if DB == nil {
		return fmt.Errorf("database is not initialized")
	}

	// Auto migrate tables
	return DB.AutoMigrate(
		&models.AdminPerpustakaan{},
		&models.AdminDPK{},
		&models.Executive{},
		&models.Perpustakaan{},
		&models.Koleksi{},
		&models.SDM{},
		&models.Pengunjung{},
		&models.Anggota{},
		&models.Verifikasi{},
		&models.Laporan{},
		&models.LogRevisi{},
		&models.Notifikasi{},
		&models.AuditLog{},
	)
}

func runAutoMigrate() error {
	migrationUser := os.Getenv("MIGRATION_DB_USER")
	migrationPassword := os.Getenv("MIGRATION_DB_PASSWORD")

	if migrationUser == "" {
		return AutoMigrateDB()
	}

	migrationDSN := buildDSN(
		getEnv("DB_HOST", "localhost"),
		migrationUser,
		migrationPassword,
		getEnv("DB_NAME", "perpustakaan_db"),
		getEnv("DB_PORT", "5432"),
	)

	migrationDB, err := gorm.Open(postgres.Open(migrationDSN), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect migration user: %w", err)
	}

	return migrationDB.AutoMigrate(
		&models.AdminPerpustakaan{},
		&models.AdminDPK{},
		&models.Executive{},
		&models.Perpustakaan{},
		&models.Koleksi{},
		&models.SDM{},
		&models.Pengunjung{},
		&models.Anggota{},
		&models.Verifikasi{},
		&models.Laporan{},
		&models.LogRevisi{},
		&models.Notifikasi{},
		&models.AuditLog{},
	)
}

func buildDSN(host, user, password, dbName, port string) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", host, user, password, dbName, port)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
