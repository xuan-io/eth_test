package tx

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"log"
	"math/big"
	"testing"
	"time"
)

func Test_sendEIP1559Tx(t *testing.T) {
	var to = ""
	var pk = ""
	var from = ""
	client, err := ethclient.Dial("https://opbnb-testnet-rpc.bnbchain.org")
	if err != nil {
		log.Fatal(err)
	}
	value, err := client.BalanceAt(context.Background(), common.HexToAddress(to), nil)
	if err != nil {
		log.Fatal(err)
	}
	nonce, err := client.NonceAt(context.Background(), common.HexToAddress(from), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("send before %v,balance:%v  nonce:%v\n", to, value, nonce)
	hash, err := sendEip1559Transaction(client, pk, from, to, 900)
	fmt.Printf("Test_sendTx result:%v\n", err)

	for {
		fmt.Printf("tx hash: %v\n", hash)
		time.Sleep(time.Second * 10)
		recept, err := client.TransactionReceipt(context.Background(), hash)
		if err == nil {
			fmt.Printf("%v\n", recept.Status)
		} else {
			fmt.Printf("%v\n", err)

			nonce, err := client.NonceAt(context.Background(), common.HexToAddress(from), nil)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("send after %v,balance:%v  nonce:%v\n", to, value, nonce)
			tx, isPending, err := client.TransactionByHash(context.Background(), hash)
			if err != nil {
				fmt.Printf("TransactionByHash Error:%v\n", err)
			} else {
				fmt.Printf("send after result hash :%v,isPending:%v, error:%v\n", tx.Hash(), isPending, err)

			}

		}
	}
	value, err = client.BalanceAt(context.Background(), common.HexToAddress(to), nil)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v,balance:%v\n", to, value)
}

func sendEip1559Transaction(cl *ethclient.Client, privateKey string, from string, to string, token int64) (common.Hash, error) {
	value := new(big.Int).Mul(big.NewInt(token), big.NewInt(params.GWei))
	sk := crypto.ToECDSAUnsafe(common.FromHex(privateKey))
	toAddress := common.HexToAddress(to)
	sender := common.HexToAddress(from)
	nonce, _ := cl.PendingNonceAt(context.Background(), sender)
	var gasLimit uint64 = 0
	gasPrice, _ := cl.SuggestGasPrice(context.Background())
	tx := types.NewTx(
		&types.DynamicFeeTx{
			Nonce:     nonce,
			GasTipCap: big.NewInt(0),
			GasFeeCap: gasPrice,
			Gas:       gasLimit,
			To:        &toAddress,
			Value:     value,
			Data:      nil,
		})
	chainid, err := cl.ChainID(context.Background())
	if err != nil {
		panic(err)
	}
	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainid), sk)
	if err != nil {
		fmt.Errorf("Error: %v\n", err)
	}

	return tx.Hash(), cl.SendTransaction(context.Background(), signedTx)
}
