package main

import (
	"crypto/rand"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"math/big"
	"sync"
	"time"
)

type Promotion struct {
	ID     string
	ProID  int
	Code   string
	Status bool
}

func main() {
	db, err := gorm.Open(postgres.Open("postgres://postgres:postgres@localhost:5435/test?sslmode=disable"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&Promotion{})
	log.Println("Basla")
	var wg sync.WaitGroup
	for i := 0; i < 500; i++ {
		wg.Add(1)
		go gen(db, &wg)
	}
	wg.Wait()
}

const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func code() string {
	ret := make([]byte, 15)
	for i := 0; i < 15; i++ {
		if i%5 == 0 && i != 0 {
			ret[i] = byte('-')
		} else {
			num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
			ret[i] = letters[num.Int64()]
		}
	}
	return string(ret)
}

func gen(db *gorm.DB, wg *sync.WaitGroup) {
	defer wg.Done()
	var data []Promotion
	basla := time.Now()
	for i := 0; i < 10000; i++ {
		data = append(data, Promotion{
			uuid.NewString(),
			2,
			code(),
			false,
		})
	}

	err := db.Create(&data).Error
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Bitti", time.Since(basla))
}
