package main

import (
	"github.com/matdorneles/leitor_csv/database"
	"github.com/matdorneles/leitor_csv/routes"
)

func main() {
	database.ConectarBancoDeDados()
	routes.SetupRoutes()
}
