package tx

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"log"
	"math/big"
	"testing"
	"time"
)

func Test_CreateKey(t *testing.T) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	priv := hexutil.Encode(privateKeyBytes)[2:]
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Printf("priv: %v address:%v\n", priv, address)
}

func Test_balance(t *testing.T) {

	client, err := ethclient.Dial("https://opbnb-testnet-rpc.bnbchain.org")
	if err != nil {
		log.Fatal(err)
	}
	value, err := client.BalanceAt(context.Background(), common.HexToAddress(to), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v,balance:%v\n", to, value)
}

var to = ""
var pk = ""
var from = ""

func Test_sendNpEIP155Tx(t *testing.T) {

	client, err := ethclient.Dial("https://opbnb-testnet-rpc.bnbchain.org")
	if err != nil {
		log.Fatal(err)
	}
	value, err := client.BalanceAt(context.Background(), common.HexToAddress(to), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v,balance:%v\n", to, value)
	hash, err := sendNoEip155Transaction(client, pk, from, to, 1)
	fmt.Printf("Test_sendTx result:%v\n", err)

	for {
		time.Sleep(time.Second * 5)
		recept, err := client.TransactionReceipt(context.Background(), hash)
		if err == nil {
			fmt.Printf("%v\n", recept.Status)
		} else {
			fmt.Printf("%v\n", err)
		}
	}
	value, err = client.BalanceAt(context.Background(), common.HexToAddress(to), nil)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v,balance:%v\n", to, value)
}

func sendNoEip155Transaction(cl *ethclient.Client, privateKey string, from string, to string, token int64) (common.Hash, error) {
	value := new(big.Int).Mul(big.NewInt(token), big.NewInt(params.Ether))
	sk := crypto.ToECDSAUnsafe(common.FromHex(privateKey))
	toAddress := common.HexToAddress(to)
	sender := common.HexToAddress(from)
	nonce, _ := cl.PendingNonceAt(context.Background(), sender)
	var gasLimit uint64 = 30000
	gasPrice, _ := cl.SuggestGasPrice(context.Background())
	tx := types.NewTx(
		&types.LegacyTx{
			Nonce:    nonce,
			GasPrice: gasPrice,
			Gas:      gasLimit,
			To:       &toAddress,
			Value:    value,
			Data:     nil,
		})
	signedTx, _ := types.SignTx(tx, types.FrontierSigner{}, sk)

	return tx.Hash(), cl.SendTransaction(context.Background(), signedTx)
}
