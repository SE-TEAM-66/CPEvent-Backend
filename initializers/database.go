package initializers

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB (){	
	var err error  
	dsn := "doadmin:AVNS_KEpHOcUKV5Q1MSrWEaM@tcp(127.0.0.1:25060)/defaultdb"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database")
	}
}