package routes

import (
	"net/http"

	"github.com/matdorneles/leitor_csv/controllers"
)

func SetupRoutes() {
	//ler o css do index.html
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	//página inicial
	http.HandleFunc("/", controllers.Index)
	//ao clicar enviar, direcionará para este endpoint
	http.HandleFunc("/upload", controllers.LerArquivo)
	//comando para iniciar servidor
	http.ListenAndServe(":8000", nil)
}
