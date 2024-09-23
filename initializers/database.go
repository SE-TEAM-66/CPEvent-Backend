package initializers

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		log.Fatal("DB_URL environment variable not set")
	}

	// Set the maximum time to retry (e.g., 24 hours)
	maxRetryDuration := 24 * time.Hour
	if retryDurationEnv := os.Getenv("DB_MAX_RETRY_DURATION"); retryDurationEnv != "" {
		if duration, err := time.ParseDuration(retryDurationEnv); err == nil {
			maxRetryDuration = duration
		}
	}

	// Retry delay (e.g., 2 seconds between retries)
	retryDelay := 2 * time.Second
	if delayEnv := os.Getenv("DB_RETRY_DELAY"); delayEnv != "" {
		if d, err := time.ParseDuration(delayEnv); err == nil {
			retryDelay = d
		}
	}

	startTime := time.Now()

	// Retry loop with time-based cutoff
	for time.Since(startTime) < maxRetryDuration {
		var err error
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Printf("Failed to connect to database: %v. Retrying in %v...\n", err, retryDelay)
			time.Sleep(retryDelay)
		} else {
			log.Println("Connected to database successfully")
			return
		}
	}

	log.Fatalf("Failed to connect to database after %v", maxRetryDuration)
}
