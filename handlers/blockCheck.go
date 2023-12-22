package handlers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/airchains-network/solana-seq-indexer/common"
	"github.com/airchains-network/solana-seq-indexer/common/logs"
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func BlockCheck(wg *sync.WaitGroup, ldb *leveldb.DB, ldt *leveldb.DB) {
	defer wg.Done()
	fileOpen, err := os.Open("data/blockCount.txt")
	if err != nil {
		logs.LogMessage("ERROR:", fmt.Sprintf("Failed to read file: %s"+err.Error()))
		os.Exit(0)
	}
	defer fileOpen.Close()

	scanner := bufio.NewScanner(fileOpen)

	blockNumberBytes := ""

	for scanner.Scan() {
		blockNumberBytes = scanner.Text()
	}

	blockNumber, blockNumberErr := strconv.Atoi(strings.TrimSpace(string(blockNumberBytes)))
	if blockNumberErr != nil {
		logs.LogMessage("ERROR:", fmt.Sprintf("Invalid block number : %s"+blockNumberErr.Error()))
		os.Exit(0)
	}

	payloadJSON, payloadJSONErr := json.Marshal(
		map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      1,
			"method":  "getBlockHeight",
		},
	)

	if payloadJSONErr != nil {
		log.Fatal(payloadJSONErr)
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
		logs.LogMessage("ERROR:", fmt.Sprintf(resErr.Error()))
		os.Exit(0)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		logs.LogMessage("ERROR:", fmt.Sprintf("Resquesting in solana RPC : %s"+res.Status))
		os.Exit(0)
	}

	var response struct {
		Jsonrpc string `json:"jsonrpc"`
		Result  int    `json:"result"`
		Id      int    `json:"id"`
	}

	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		logs.LogMessage("ERROR:", fmt.Sprintf(err.Error()))
		os.Exit(0)
	}

	if blockNumber == response.Result {
		logs.LogMessage("INFO:", fmt.Sprintf("Block numbers match. Waiting for %d seconds before checking again...", common.BlockDelay))
		time.Sleep(time.Duration(common.BlockDelay) * time.Second)
		BlockCheck(wg, ldb, ldt)
	} else {
		BlockSave(wg, blockNumber, ldb, ldt)
	}
}
