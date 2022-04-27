package controllers

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/matdorneles/leitor_csv/database"
	"github.com/matdorneles/leitor_csv/models"
)

//função para verificar se há algum campo vazio
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

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

func UploadArquivo(w http.ResponseWriter, r *http.Request) {

	//lê o arquivo CSV e o retorna em linhas JSON
	r.ParseMultipartForm(10 << 20) // recebe arquivo enviado do HTML, tamanho máximo do arquivo = 10mb

	// Handler para nome do arquivo, tamanho, header
	file, handler, err := r.FormFile("arquivo")
	if err != nil {
		fmt.Println("Não foi possível completar o upload")
		fmt.Println(err)
		return
	}
	defer file.Close()

	// Verificando se o arquivo não está vazio, um CSV vazio possui ~2 bytes
	var verificarSize bytes.Buffer
	tamanhoArquivo, err := verificarSize.ReadFrom(file)
	tamanhoArquivoConv := float64(tamanhoArquivo)
	if err != nil || tamanhoArquivoConv <= 2 {
		http.Error(w, "O arquivo está vazio ou é menor/igual a 2 bytes", http.StatusBadRequest)
		return
	}

	fmt.Printf("Arquivo enviado: %+v\n", handler.Filename)
	fmt.Printf("Tamanho do arquivo: %+v\n", handler.Size)
	fmt.Printf("Header: %+v\n", handler.Header)

	// Criando arquivo
	fileCsv, err := os.Create(handler.Filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Copiando arquivo do upload para o criado no sistema
	if _, err := io.Copy(fileCsv, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Arquivo enviado com sucesso!")

	arquivoCsv, err := os.Open(handler.Filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer arquivoCsv.Close()

	leitor := csv.NewReader(arquivoCsv)
	var transacoes []models.Transacao

	for {
		dadosCSV, err := leitor.ReadAll()
		if err != nil {
			log.Fatal(err)
		}

		//verificando data da primeira linha
		dtPrimeiraTransacao, err := time.Parse("2006-01-02T15:04:05", dadosCSV[0][7])
		if err != nil {
			fmt.Println(err)
			return
		}

		//lendo linha por linha e guardando dados para o DB
		for _, linha := range dadosCSV {

			if contains(linha, "") {
				continue
			}

			dtTransacao, err := time.Parse("2006-01-02T15:04:05", linha[7])
			if err != nil {
				continue
			}

			if dtTransacao.Day() != dtPrimeiraTransacao.Day() && dtTransacao.Month() != dtPrimeiraTransacao.Month() && dtTransacao.Year() != dtPrimeiraTransacao.Year() {
				continue
			}

			transacoes = append(transacoes, models.Transacao{
				BancoOrigem:       linha[0],
				AgenciaOrigem:     linha[1],
				ContaOrigem:       linha[2],
				BancoDestino:      linha[3],
				AgenciaDestino:    linha[4],
				ContaDestino:      linha[5],
				ValorTransacao:    linha[6],
				DataHoraTransacao: dtTransacao,
			})
		}
		transacaoJson, _ := json.Marshal(transacoes)
		fmt.Println(string(transacaoJson))

		database.DB.Create(&transacoes)
	}
}
