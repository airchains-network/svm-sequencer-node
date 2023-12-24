package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/airchains-network/svm-sequencer-node/common/logs"
	"github.com/btcsuite/btcutil/base58"
	"github.com/joho/godotenv"
	"math/big"
	"net/http"
	"os"
	"strconv"
)

func ENVCheck() (string, string) {

	err := godotenv.Load()
	if err != nil {
		return "", fmt.Sprintf("Error in loading .env file")
	}

	ExecutionClientRPC = os.Getenv("ExecutionClientRPC")
	if ExecutionClientRPC == "" {
		return "", fmt.Sprintf("ExecutionClientRPC is not set")
	}

	DaClientRPC = os.Getenv("DaClientRPC")
	if ExecutionClientRPC == "" {
		return "", fmt.Sprintf("DaClientRPC is not set")
	}

	SettlementClientRPC = os.Getenv("SettlementClientRPC")
	if ExecutionClientRPC == "" {
		return "", fmt.Sprintf("SettlementClientRPC is not set")
	}

	return "ENV check", ""
}

func Base52Decoder(value string) string {
	decodedBytes := base58.Decode(value)
	decodedBigInt := new(big.Int).SetBytes(decodedBytes)
	return decodedBigInt.String()
}

func AccountNouceCheck(accountKey string) string {
	payloadJSON, payloadJSONErr := json.Marshal(
		map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      1,
			"method":  "getConfirmedSignaturesForAddress2",
			"params": []interface{}{
				accountKey,
			},
		},
	)
	if payloadJSONErr != nil {
		logs.LogMessage("ERROR:", fmt.Sprintf("Failed to read file: %s"+payloadJSONErr.Error()))
		os.Exit(0)
	}

	client := &http.Client{}
	req, reqErr := http.NewRequest("POST", ExecutionClientRPC, bytes.NewBuffer(payloadJSON))
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

	var accountInfo struct {
		JSONRPC string `json:"jsonrpc"`
		Result  []struct {
			BlockTime          int         `json:"blockTime"`
			ConfirmationStatus string      `json:"confirmationStatus"`
			Err                interface{} `json:"err"`
			Memo               interface{} `json:"memo"`
			Signature          string      `json:"signature"`
			Slot               int         `json:"slot"`
		} `json:"result"`
		ID int `json:"id"`
	}

	decodeError := json.NewDecoder(res.Body).Decode(&accountInfo)
	if decodeError != nil {
		logs.LogMessage("ERROR:", fmt.Sprintf(decodeError.Error()))
		os.Exit(0)
	}

	return strconv.Itoa(len(accountInfo.Result))
}
