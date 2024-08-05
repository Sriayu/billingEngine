package postgres

import (
	"database/sql"
	"fmt"
	"log"

	gormLogger "gorm.io/gorm/logger"

	config "BillingEngine/src/infra"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDb struct {
	DB    *gorm.DB
	SqlDB *sql.DB
}

var gormOpen = gorm.Open

func New(conf config.SqlDbConf) *PostgresDb {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		conf.Host,
		conf.Username,
		conf.Password,
		conf.Name,
		conf.Port,
	)

	if conf.Password == "" {
		dsn = fmt.Sprintf(
			"host=%s user=%s dbname=%s port=%s sslmode=disable",
			conf.Host,
			conf.Username,
			conf.Name,
			conf.Port,
		)
	}

	queryLogger := gormLogger.Default.LogMode(gormLogger.Silent)
	db, err := gormOpen(postgres.New(postgres.Config{
		DSN: dsn,
	}), &gorm.Config{
		Logger: queryLogger,
	})
	if err != nil {
		panic("Failed to connect to database!")
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("database err: %s", err)
	}

	err = sqlDB.Ping()
	if err != nil {
		panic(err)
	}

	log.Printf("sql database connection %s success", db.Name())

	return &PostgresDb{
		DB:    db,
		SqlDB: sqlDB,
	}
}
