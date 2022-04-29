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
	"path"
	"regexp"

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

//lê o arquivo CSV e o retorna em linhas JSON
func UploadArquivo(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(10 << 20) // Tamanho máximo do arquivo = 10mb

	// Handler para nome do arquivo, tamanho, header
	file, handler, err := r.FormFile("arquivo")
	if err != nil {
		fmt.Println("Não foi possível completar o upload")
		fmt.Println(err)
		return
	}
	defer file.Close()

	filename := path.Base(handler.Filename)

	fmt.Printf("Arquivo enviado: %+v\n", handler.Filename)
	fmt.Printf("Tamanho do arquivo: %+v\n", handler.Size)
	fmt.Printf("Header: %+v\n", handler.Header)

	if handler.Size <= 2 {
		http.Error(w, "Arquivo está vazio", http.StatusBadRequest)
	}

	// Criando arquivo
	csv, err := os.Create(filename)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer csv.Close()

	// Copiando arquivo do upload para o criado no sistema
	if _, err := io.Copy(csv, file); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Arquivo enviado com sucesso!")
	LerArquivo(filename)

}

// Lendo arquivo, separando linhas e virgulas e adicionando a cada atributo da classe Transacao
func LerArquivo(arquivo string) {
	arquivoCsv, err := os.Open(arquivo)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer arquivoCsv.Close()

	leitor := csv.NewReader(arquivoCsv)
	var transacoes []models.Transacao

	dadosCSV, err := leitor.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	//verificando data da primeira linha
	re := regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
	fmt.Printf("Procurando pelo padrão: %v\n", re.String())

	dtPrimeiraTransacao := re.FindString(dadosCSV[0][7])
	fmt.Printf("A data encontrada foi: %v\n", dtPrimeiraTransacao)

	// dtPrimeiraTransacao, err := time.Parse("2006-01-02T15:04:05", dadosCSV[0][7])
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	//lendo linha por linha e guardando dados para o DB
	for _, linha := range dadosCSV {

		if contains(linha, "") {
			continue
		}

		dtTransacao := re.FindString(linha[7])

		if dtTransacao != dtPrimeiraTransacao {
			continue
		}

		// dtTransacao, err := time.Parse("2006-01-02T15:04:05", linha[7])
		// if err != nil {
		// 	continue
		// }

		// if dtTransacao.Day() != dtPrimeiraTransacao.Day() && dtTransacao.Month() != dtPrimeiraTransacao.Month() && dtTransacao.Year() != dtPrimeiraTransacao.Year() {
		// 	continue
		// }

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
