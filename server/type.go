package main

type Transaction struct {
	FromUID string  `json:"from_uid"`
	ToUID   string  `json:"to_uid"`
	Amount  float64 `json:"amount"`
}

type FixedTransaction struct {
	FromUID   string  `json:"from_uid"`
	ToUID     string  `json:"to_uid"`
	Amount    float64 `json:"amount"`
	Timestamp string  `json:"timestamp"`
}

type AddTransactionResponse struct {
	Timestamp string `json:"timestamp"`
}

type GetTaskRequest struct{}

type Task struct {
	Exists      bool              `json:"exists"`
	Transaction *FixedTransaction `json:"transaction"`
	PrevHash    string            `json:"prev_hash"`
}

type Nonce struct {
	PrevHash string `json:"prev_hash"`
	Nonce    string `json:"nonce"`
}

type ResolveNonceResponse struct {
	Succeeded bool    `json:"succeeded"`
	Amount    float64 `json:"amount"`
}
