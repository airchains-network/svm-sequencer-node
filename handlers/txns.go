package handlers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/airchains-network/solana-seq-indexer/common"
	"github.com/airchains-network/solana-seq-indexer/common/logs"
	"github.com/airchains-network/solana-seq-indexer/types"
	"github.com/syndtr/goleveldb/leveldb"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func insertTxn(db *leveldb.DB, txns types.TransactionStruck, transactionNumber int) error {
	data, err := json.Marshal(txns)
	if err != nil {
		return err
	}

	txnsKey := fmt.Sprintf("txns-%d", transactionNumber+1)
	err = db.Put([]byte(txnsKey), data, nil)
	if err != nil {
		return err
	}
	err = os.WriteFile("data/transactionCount.txt", []byte(strconv.Itoa(transactionNumber+1)), 0666)
	if err != nil {
		return err
	}

	return nil
}

func TxnSave(blockHash string, signature string, lb *leveldb.DB) {
	payloadJSON, payloadJSONErr := json.Marshal(
		map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      1,
			"method":  "getTransaction",
			"params": []interface{}{
				signature,
			},
		},
	)
	if payloadJSONErr != nil {
		logs.LogMessage("ERROR:", fmt.Sprintf("Failed to read file: %s"+payloadJSONErr.Error()))
		os.Exit(0)
	}

	client := &http.Client{}
	req, reqErr := http.NewRequest("POST", common.ExecutionClientRPC, bytes.NewBuffer(payloadJSON))
	if reqErr != nil {
		logs.LogMessage("ERROR:", fmt.Sprintf("Resquesting in solana RPC : %s"+reqErr.Error()))
		os.Exit(0)
	}

	req.Header.Set("Content-Type", "application/json")
	res, resErr := client.Do(req)
	if resErr != nil {
		logs.LogMessage("ERROR:", fmt.Sprintf("Resquesting in solana RPC : %s"+resErr.Error()))
		os.Exit(0)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		logs.LogMessage("ERROR:", fmt.Sprintf(res.Status))
		os.Exit(0)
	}

	var response types.TxnResponce
	decodeErr := json.NewDecoder(res.Body).Decode(&response)
	if decodeErr != nil {
		logs.LogMessage("ERROR:", fmt.Sprintf(decodeErr.Error()))
		os.Exit(0)
	}

	txnData := types.TransactionStruck{
		Signature:   signature,
		BlockHeight: response.Result.Slot,
		BlockHash:   blockHash,
		Timestamp:   response.Result.BlockTime,
		Mata:        response.Result.Meta,
		Transaction: response.Result.Transaction,
	}

	fileOpen, err := os.Open("data/transactionCount.txt")
	if err != nil {
		logs.LogMessage("ERROR:", fmt.Sprintf("Failed to read file: %s"+err.Error()))
		os.Exit(0)
	}
	defer fileOpen.Close()

	scanner := bufio.NewScanner(fileOpen)

	transactionNumberBytes := ""

	for scanner.Scan() {
		transactionNumberBytes = scanner.Text()
	}

	transactionNumber, err := strconv.Atoi(strings.TrimSpace(string(transactionNumberBytes)))
	if err != nil {
		logs.LogMessage("ERROR:", fmt.Sprintf("Invalid transaction number : %s"+err.Error()))
		os.Exit(0)
	}

	insetTxnErr := insertTxn(lb, txnData, transactionNumber)
	if insetTxnErr != nil {
		logs.LogMessage("ERROR:", fmt.Sprintf("Failed to insert transaction: %s"+insetTxnErr.Error()))
		os.Exit(0)
	}

	logs.LogMessage("SUCCESS:", fmt.Sprintf("Successfully saved Transation %s in the latest phase", signature))

}
