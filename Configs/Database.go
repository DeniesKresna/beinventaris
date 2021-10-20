package Configs

import (
	"log"
	"os"
	"strconv"

	"github.com/DeniesKresna/beinventaris/Models"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func DatabaseInit() (err error) {
	er := godotenv.Load(".env")
	if er != nil {
		log.Fatalf("Error loading .env file")
	}
	port, er := strconv.Atoi(os.Getenv("DB_PORT"))
	if er != nil {
		log.Fatal("Error Convert PORT")
		return
	}

	dsn := os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@tcp(" + os.Getenv("DB_HOST") + ":" + strconv.Itoa(port) + ")/" + os.Getenv("DB_NAME") + "?charset=utf8mb4&parseTime=True&loc=Local"
	//log.Fatal(dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	DB = db
	return err
}

func DatabaseMigrate() {
	DB.AutoMigrate(&Models.User{}, &Models.Role{}, &Models.Academy{}, &Models.Unit{}, &Models.Room{}, &Models.GoodsType{}, &Models.Inventory{},
		&Models.Condition{}, &Models.History{}, &Models.Period{})
}

func init() {
	os.Setenv("TZ", "Asia/Jakarta")
}
