package database

import (
	"github.com/blyndusk/cofy/api/models"
	"github.com/sirupsen/logrus"
)

func Migrate() {
	_ = Db.AutoMigrate(&models.User{})
	logrus.Info("Migrations done !")
}
