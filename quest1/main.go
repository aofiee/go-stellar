package main

import (
	"log"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
)

func main() {
	pair, _ := keypair.Parse("")
	addr := pair.Address()
	client := horizonclient.DefaultTestNetClient
	accountRequest := horizonclient.AccountRequest{AccountID: addr}
	account1, err := client.AccountDetail(accountRequest)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Account:", account1)
	des, _ := keypair.Random()

	createAccountOp := txnbuild.CreateAccount{
		Destination: des.Address(),
		Amount:      "1000",
	}
	txParams := txnbuild.TransactionParams{
		SourceAccount:        &account1,
		IncrementSequenceNum: true,
		Operations:           []txnbuild.Operation{&createAccountOp},
		BaseFee:              txnbuild.MinBaseFee,
		Timebounds:           txnbuild.NewTimeout(300),
	}
	tx, _ := txnbuild.NewTransaction(txParams)
	signedTx, err := tx.Sign(network.TestNetworkPassphrase, pair.(*keypair.Full))
	if err != nil {
		log.Fatalln(err)
	}
	txeBase64, err := signedTx.Base64()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Transaction base64: ", txeBase64)

	resp, err := client.SubmitTransactionXDR(txeBase64)
	if err != nil {
		hError := err.(*horizonclient.Error)
		log.Fatal("Error submitting transaction:", hError.Problem)
	}

	log.Println("\nTransaction response: ", resp)
}
