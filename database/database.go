package database

import (
	"log"

	"github.com/matdorneles/leitor_csv/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
)

func ConectarBancoDeDados() {
	dsn := "host=localhost user=root password=root dbname=root port=5432 sslmode=disable TimeZone=America/Sao_Paulo"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panic("Erro ao conectar banco de dados")
	}
	DB.AutoMigrate(&models.Transacao{})
}
