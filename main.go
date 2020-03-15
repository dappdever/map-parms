package main

import (
	"encoding/hex"
	"fmt"
	"github.com/eoscanada/eos-go"
	"log"
)

// PubMapActionPayload ...
type PubMapActionPayload struct {
	PayloadMap map[string]string `json:"payload"`
}

func newPubMap(payload map[string]string) *eos.Action {
	return &eos.Action{
		Account: eos.AN("messengerbus"),
		Name:    eos.ActN("pubmap"),
		Authorization: []eos.PermissionLevel{
			{Actor: eos.AN("messengerbus"), Permission: eos.PN("active")},
		},
		ActionData: eos.NewActionData(PubMapActionPayload{
			PayloadMap: payload,
		}),
	}
}

func PublishMapToBlockchain(payload map[string]string) (string, error) {
	api := eos.New("https://kylin.eosusa.news")

	keyBag := &eos.KeyBag{}
	err := keyBag.ImportPrivateKey("5KAP1zytghuvowgprSPLNasajibZcxf4KMgdgNbrNj98xhcGAUa")
	if err != nil {
		log.Printf("import private key: %s", err)
	}
	api.SetSigner(keyBag)

	txOpts := &eos.TxOptions{}
	if err := txOpts.FillFromChain(api); err != nil {
		log.Printf("Error filling tx opts: %s", err)
		return "error", err
	}

	tx := eos.NewTransaction([]*eos.Action{newPubMap(payload)}, txOpts)
	_, packedTx, err := api.SignTransaction(tx, txOpts.ChainID, eos.CompressionNone)
	if err != nil {
		log.Printf("Error signing transaction: %s", err)
		return "error", err
	}

	response, err := api.PushTransaction(packedTx)
	if err != nil {
		log.Printf("Error pushing transaction: %s", err)
		return "error", err
	}
	return hex.EncodeToString(response.Processed.ID), nil
}

func main() {

	payload := make(map[string]string)
	payload["foo"] = "bar"
	payload["abc"] = "xyz"

	trxID, err := PublishMapToBlockchain(payload)
	if err != nil {
		panic(err)
	}
	fmt.Println("Transaction ID: ", trxID)
}
