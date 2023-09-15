package models

import "math/big"

type ECDHReqDTO struct {
	Token       string   `json:"token"`
	PubKeyWasmX *big.Int `json:"pub_key_wasm_x"`
	PubKeyWasmY *big.Int `json:"pub_key_wasm_y"`
}

type ECDHKeyExchangeOutput struct {
	PubKeyServerX *big.Int `json:"pub_key_server_x"`
	PubKeyServerY *big.Int `json:"pub_key_server_y"`
}
