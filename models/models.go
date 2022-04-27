package models

import "time"

type Transacao struct {
	BancoOrigem       string    `json:"banco-origem"`
	AgenciaOrigem     string    `json:"agencia-origem"`
	ContaOrigem       string    `json:"conta-origem"`
	BancoDestino      string    `json:"banco-destino"`
	AgenciaDestino    string    `json:"agencia-destino"`
	ContaDestino      string    `json:"conta-destino"`
	ValorTransacao    string    `json:"valor-transacao"`
	DataHoraTransacao time.Time `json:"data-hora-transacao"`
}
