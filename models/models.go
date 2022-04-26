package models

type Transacao struct {
	BancoOrigem       string `json:"banco-origem" validate:"nonzero"`
	AgenciaOrigem     string `json:"agencia-origem" validate:"nonzero"`
	ContaOrigem       string `json:"conta-origem" validate:"nonzero"`
	BancoDestino      string `json:"banco-destino" validate:"nonzero"`
	AgenciaDestino    string `json:"agencia-destino" validate:"nonzero"`
	ContaDestino      string `json:"conta-destino" validate:"nonzero"`
	ValorTransacao    string `json:"valor-transacao" validate:"nonzero"`
	DataHoraTransacao string `json:"data-hora-transacao" validate:"nonzero"`
}
