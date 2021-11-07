package main

import (
	"log"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
)

func main() {
	distributor := ""

	issuer := ""

	client := horizonclient.DefaultTestNetClient
	issuerKeyPair := keypair.MustParseFull(issuer)
	distributorPair := keypair.MustParseFull(distributor)

	request := horizonclient.AccountRequest{AccountID: issuerKeyPair.Address()}
	issuerAccount, err := client.AccountDetail(request)
	if err != nil {
		panic(err)
	}
	request = horizonclient.AccountRequest{AccountID: distributorPair.Address()}
	distributorAccount, err := client.AccountDetail(request)
	if err != nil {
		panic(err)
	}

	finalGravity := txnbuild.CreditAsset{Code: "FinalGravity", Issuer: issuerKeyPair.Address()}
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &distributorAccount,
			IncrementSequenceNum: true,
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewInfiniteTimeout(),
			Operations: []txnbuild.Operation{
				&txnbuild.ChangeTrust{
					Line:  finalGravity.MustToChangeTrustAsset(),
					Limit: "1000000",
				},
			},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	signedTx, err := tx.Sign(network.TestNetworkPassphrase, distributorPair)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.SubmitTransaction(signedTx)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Trust: %s\n", resp.Hash)
	}

	destAccountRequest := horizonclient.AccountRequest{AccountID: distributorPair.Address()}
	destinationAccount, err := client.AccountDetail(destAccountRequest)
	if err != nil {
		panic(err)
	}

	tx, err = txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &issuerAccount,
			IncrementSequenceNum: true,
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewInfiniteTimeout(),
			Operations: []txnbuild.Operation{
				&txnbuild.Payment{
					Destination: destinationAccount.AccountID,
					Asset:       finalGravity,
					Amount:      "1000000",
				},
			},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	signedTx, err = tx.Sign(network.TestNetworkPassphrase, issuerKeyPair)
	if err != nil {
		log.Fatal(err)
	}
	resp, err = client.SubmitTransaction(signedTx)

	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Pay: %s\n", resp.Hash)
	}
}
