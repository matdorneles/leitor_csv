package controllers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/matdorneles/leitor_csv/models"
)

//variável apontando pasta de templates html
var temp = template.Must(template.ParseGlob("templates/*.html"))

//ao acessar a página incial '/' direcionará para index.html
func Index(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		temp.ExecuteTemplate(w, "Index", nil)
	case "POST":
		UploadArquivo(w, r)
	}
}

//lê o arquivo CSV e o retorna em linhas JSON
func UploadArquivo(w http.ResponseWriter, r *http.Request) {
	// Tamanho máximo do arquivo = 10mb
	r.ParseMultipartForm(10 << 20)

	// Handler para nome do arquivo, tamanho, header
	file, handler, err := r.FormFile("arquivo")
	if err != nil {
		fmt.Println("Não foi possível completar o upload")
		fmt.Println(err)
		return
	}

	defer file.Close()
	fmt.Printf("Arquivo enviado: %+v\n", handler.Filename)
	fmt.Printf("Tamanho do arquivo: %+v\n", handler.Size)
	fmt.Printf("Header: %+v\n", handler.Header)

	// Criando arquivo
	csv, err := os.Create(handler.Filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Copiando arquivo do upload para o criado no sistema
	if _, err := io.Copy(csv, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Arquivo enviado com sucesso!")
	LerArquivo(string(handler.Filename))

}

func LerArquivo(arquivo string) {
	arquivoCsv, err := os.Open(arquivo)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer arquivoCsv.Close()

	leitor := csv.NewReader(arquivoCsv)
	var transacoes []models.Transacao

	for {
		linha, err := leitor.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		transacoes = append(transacoes, models.Transacao{
			BancoOrigem:       linha[0],
			AgenciaOrigem:     linha[1],
			ContaOrigem:       linha[2],
			BancoDestino:      linha[3],
			AgenciaDestino:    linha[4],
			ContaDestino:      linha[5],
			ValorTransacao:    linha[6],
			DataHoraTransacao: linha[7],
		})
	}

	transacaoJson, _ := json.Marshal(transacoes)
	fmt.Println(string(transacaoJson))
}
