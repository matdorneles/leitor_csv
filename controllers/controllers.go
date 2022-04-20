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
	temp.ExecuteTemplate(w, "Index", nil)
}

//lê o arquivo CSV e o retorna em linhas JSON
func LerArquivo(w http.ResponseWriter, r *http.Request) {
	arquivoCsv, err := os.Open("transacoes-2022-01-01.csv")
	if err != nil {
		fmt.Println(err)
	}

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
	defer arquivoCsv.Close()
}
