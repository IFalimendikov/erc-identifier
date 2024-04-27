package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/joho/godotenv"
)

type EtherscanResp struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  string `json:"result"`
}

// getABI fetches the ABI for a given contract address from Etherscan.
func getABI(addr string) (abi.ABI, error) {
	api := os.Getenv("ETHERSCAN_API_KEY")
	url := fmt.Sprintf("https://api.etherscan.io/api?module=contract&action=getabi&address=%s&apikey=%s", addr, api)

	// Perform the HTTP GET request to the Etherscan API.
	resp, err := http.Get(url)
	if err != nil {
		return abi.ABI{}, fmt.Errorf("failed to make HTTP request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return abi.ABI{}, fmt.Errorf("failed to read response body: %v", err)
	}

	var ethResp EtherscanResp
	err = json.Unmarshal(body, &ethResp)
	if err != nil {
		return abi.ABI{}, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	// Check for unverified or zero contract response.
	if ethResp.Status == "0" {
		if ethResp.Result == "Contract source code not verified" {
			return abi.ABI{}, nil
		}
		return abi.ABI{}, fmt.Errorf("not real EVM address: %s", ethResp.Message)
	}

	pAbi, err := abi.JSON(strings.NewReader(ethResp.Result))
	if err != nil {
		return abi.ABI{}, fmt.Errorf("failed to parse ABI JSON: %v", err)
	}

	return pAbi, nil
}

// parseABI parses the ABI of a contract address and checks if it complies with any of the ERC standards.
func parseABI(addr string) (string, error) {
	abi, err := getABI(addr)
	if err != nil {
		return "", fmt.Errorf("failed to get contract ABI: %v", err)
	}

	if len(abi.Methods) == 0 && len(abi.Events) == 0 {
		return "Address is either a wallet or the contract source code is not verified", nil
	}
	// Define the sequence of token standards explicitly for ordered checking.
	order := []string{"ERC-20", "ERC-721", "ERC-1155"}

	// Map token standards to their corresponding ABI files.
	abiFiles := map[string]string{
		"ERC-20":   "abi/erc20.abi.json",
		"ERC-721":  "abi/erc721.abi.json",
		"ERC-1155": "abi/erc1155.abi.json",
	}

	// Check ABIs in the specified order.
	for _, erc := range order {
		fmt.Printf("Checking %s...\n", erc)

		// Read the ABI from the corresponding file.
		path := abiFiles[erc]
		abiBytes, err := os.ReadFile(path)
		if err != nil {
			return "", fmt.Errorf("failed to read %s: %v", path, err)
		}

		var methodsEvents []map[string]interface{}
		if err := json.Unmarshal(abiBytes, &methodsEvents); err != nil {
			return "", fmt.Errorf("failed to parse %s ABI JSON: %v", erc, err)
		}

		// Create a map to check the existence of each method and event in the fetched ABI.
		allExists := true
		for _, item := range methodsEvents {
			if name, ok := item["name"].(string); ok {
				if _, exists := abi.Methods[name]; !exists {
					if _, exists := abi.Events[name]; !exists {
						allExists = false
						break
					}
				}
			}
		}

		// If all items exist, return the standard.
		if allExists {
			return fmt.Sprintf("Contract ABI complies with the %s standard", erc), nil
		}
	}

	// If no standards matched, return that none suffice.
	return "Contract ABI does not comply with any of the ERC-20, ERC-721, ERC-1155 standards", nil
}

func main() {
	godotenv.Load()
	contract := "0xdAC17F958D2ee523a2206206994597C13D831ec7"

	erc, err := parseABI(contract)
	if err != nil {
		fmt.Println("Failed to parse contract ABI:", err)
	}

	fmt.Println(erc)
}
