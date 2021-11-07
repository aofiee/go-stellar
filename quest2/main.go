package main

import (
	"log"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
)

func main() {
	source := ""
	destination := ""
	client := horizonclient.DefaultTestNetClient
	destAccountRequest := horizonclient.AccountRequest{AccountID: destination}
	destinationAccount, err := client.AccountDetail(destAccountRequest)
	if err != nil {
		panic(err)
	}

	sourceKeyPair := keypair.MustParseFull(source)
	sourceAccountRequest := horizonclient.AccountRequest{AccountID: sourceKeyPair.Address()}
	sourceAccount, err := client.AccountDetail(sourceAccountRequest)
	if err != nil {
		panic(err)
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: true,
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewInfiniteTimeout(),
			Operations: []txnbuild.Operation{
				&txnbuild.Payment{
					Destination: destinationAccount.AccountID,
					Amount:      "10",
					Asset:       txnbuild.NativeAsset{},
				},
			},
		},
	)

	if err != nil {
		panic(err)
	}

	tx, err = tx.Sign(network.TestNetworkPassphrase, sourceKeyPair)
	if err != nil {
		panic(err)
	}

	resp, err := horizonclient.DefaultTestNetClient.SubmitTransaction(tx)
	if err != nil {
		panic(err)
	}

	log.Println("Successful Transaction:")
	log.Println("Ledger:", resp.Ledger)
	log.Println("Hash:", resp.Hash)
}
