package main

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func clientDb(host, user, pass string, port int, verbose bool) *gorm.DB {
	dsnSource := fmt.Sprintf("host=%s user=%s password=%s dbname=registry port=%d sslmode=disable", host, user, pass, port)

	var dbSource *gorm.DB
	var err error
	if verbose {
		dbSource, err = gorm.Open(postgres.Open(dsnSource), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	} else {
		dbSource, err = gorm.Open(postgres.Open(dsnSource))
	}

	if err != nil {
		panic("failed to connect database")
	}
	log.Printf("connected db host:  %s, user: %s, port: %d", host, user, port)

	return dbSource
}
