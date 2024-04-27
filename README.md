
## 🔗 ERC Identifier shows if a contract complies with ERC20, ERC721 or ERC1155 standards 🤖

### Overview

ERC Identifier is a Go-based service designed to analyze Ethereum smart contracts and determine if they comply with specific ERC (Ethereum Request for Comments) standards, namely ERC20, ERC721, and ERC1155. This tool is particularly useful for developers, auditors, and users who need to ensure the functionality and security of smart contracts before interacting with them on the Ethereum blockchain.

### Features

- Analyzes smart contracts to check for compliance with ERC20, ERC721, and ERC1155 standards.
- Fetches contract ABI from Etherscan using the contract address.
- Parses the ABI to check if it matches the expected methods and events for each ERC standard.

### Usage

1. Ensure you have Go installed on your system.
2. Clone this repository or download the source code.
3. Set up your environment variables in the `example.env` file or directly in your environment. You will need an Etherscan API key (Not necessary but will give you higher request rate).
4. Rename `example.env` to `.env`, remember not to save your API keys or sensitive data publicly.
5. Enter the contract address in the `main` function, you wish to analyze.
6. Run the application using `go run main.go`.


### Example
```bash
# Run the program
$ go run main.go

# Output
Checking ERC-20...
Contract ABI complies with the ERC-20 standard

```

