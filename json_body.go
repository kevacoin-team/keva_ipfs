package main

// keva_ipfs config
type ServerConfig struct {
	Electrum_host	string	`json:"Electrum_host"`
	Electrum_port	int		`json:"Electrum_port"`
	Min_payment		float64	`json:"Min_payment"`
	Payment_address	string	`json:"Payment_address"`
	Tls_enabled		bool	`json:"Tls_enabled"`
	Tls_key			string	`json:"Tls_key"`
	Tls_cert		string	`json:"Tls_cert"`
}

// MediaResponse the reponse from uploadMedia
type MediaResponse struct {
	CID string `json:"CID"`
}

// PinMedia the data body for pinning media to IPFS
type PinMedia struct {
	Tx string `json:"tx"`
}

// PaymentInfo payment information for the server.
type PaymentInfo struct {
	PaymentAddress string `json:"payment_address"`
	MinPayment     string `json:"min_payment"`
}
