package handlers

import (
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
	"sync"
)

func insertBlock(db *leveldb.DB, block types.BlockStuct) error {
	data, err := json.Marshal(block)
	if err != nil {
		return err
	}

	err = db.Put([]byte(fmt.Sprint(block.Blockheight)), data, nil)
	if err != nil {
		return err
	}

	return nil
}

func BlockSave(wg *sync.WaitGroup, blockIndex int, ldb *leveldb.DB, ldt *leveldb.DB) {
	payloadJSON, payloadJSONErr := json.Marshal(
		map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      1,
			"method":  "getBlock",
			"params": []interface{}{
				blockIndex,
				map[string]interface{}{
					"encoding":                       "json",
					"maxSupportedTransactionVersion": 0,
					"transactionDetails":             "full",
					"rewards":                        false,
				},
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
		logs.LogMessage("ERROR:", fmt.Sprintf("Resquesting in solana RPC : %s"+res.Status))
		os.Exit(0)
	}

	var response types.JsonRpcResponse
	err := json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		logs.LogMessage("ERROR:", fmt.Sprintf(err.Error()))
		os.Exit(0)
	}

	var signatures []string

	for _, transaction := range response.Result.Transactions {
		signatures = append(signatures, transaction.Transaction.Signatures[0])
	}

	blockData := types.BlockStuct{
		Blockheight:       response.Result.BlockHeight,
		Blockhash:         response.Result.Blockhash,
		Parentslot:        response.Result.ParentSlot,
		Previousblockhash: response.Result.PreviousBlockhash,
		Timestamp:         response.Result.BlockTime,
		TransactionLength: len(response.Result.Transactions),
		Transactions:      signatures,
	}

	insetBlocker := insertBlock(ldb, blockData)
	if insetBlocker != nil {
		logs.LogMessage("ERROR:", fmt.Sprintf("Failed to insert block: %s"+insetBlocker.Error()))
		os.Exit(0)
	}

	for index := 0; index < len(response.Result.Transactions); index++ {
		TxnSave(blockData.Blockhash, blockData.Transactions[index], ldt)
	}

	err = os.WriteFile("data/blockCount.txt", []byte(strconv.Itoa(blockIndex+1)), 0666)
	if err != nil {
		logs.LogMessage("ERROR:", fmt.Sprintf("Failed to update block file: %s"+err.Error()))
		os.Exit(0)
	}

	logs.LogMessage("SUCCESS:", fmt.Sprintf("Successfully saved Block %s in the latest phase", strconv.Itoa(blockIndex+1)))

	BlockCheck(wg, ldb, ldt)
}
