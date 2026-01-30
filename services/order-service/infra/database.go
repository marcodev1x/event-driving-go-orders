package infra

import (
	"order-service/infra/config"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database struct {
	Client *gorm.DB
}

var (
	DomainDatabase *gorm.DB
)

func (d *Database) Connect(data *config.DatabaseConfig) (*gorm.DB, error) {
	account := mysql.Config{
		DSN: data.User + ":" + data.Password + "@tcp(" + data.Host + ":" + data.Port + ")/" + data.Name + "?charset=utf8mb4&parseTime=True&loc=Local",
	}

	db, err := gorm.Open(mysql.New(account), &gorm.Config{})

	if err != nil {
		config.Logger().Error("Failed to connect to database", zap.Error(err))
		return nil, err
	}

	d.Client = db
	DomainDatabase = db

	return d.Client, nil
}
