# Generic EVM API (Go)

A Go-based API service for interacting with EVM-compatible blockchains. This API provides various endpoints for querying smart contract data and performing contract calls.

## Features

- Contract code size retrieval
- Contract bytecode fetching
- Contract storage data reading
- View function calls
- Contract balance checking
- Version information

## Prerequisites

- Go 1.22 or higher
- Node.js 18 or higher
- npm

## Installation & Setup

1. Clone the repository:
```bash
git clone https://github.com/yourusername/generic-evm-api-go.git
cd generic-evm-api-go
```

2. Install Go dependencies:
```bash
go mod download
```

3. Install Node.js dependencies (for testing):
```bash
cd test
npm install
```

## Running the API

1. Start the API server:
```bash
go run main.go
```

2. Run the test suite:
```bash
cd test
./run.sh
```

## API Reference

Base URL: `generic-evm-api-go.vercel.app/api`

All endpoints use the query format: `?query=<endpoint-name>&<parameters>`

### Available Endpoints

#### 1. Get Contract External Code Size
- Endpoint: `?query=evm-contract-ext-code-size`
- Parameters:
  - `chain-id`: Chain ID (required)
  - `json-rpc`: JSON-RPC endpoint (optional)
  - `contract-address`: Contract address (required)

#### 2. Get Contract Code
- Endpoint: `?query=evm-contract-code`
- Parameters:
  - `chain-id`: Chain ID (required)
  - `json-rpc`: JSON-RPC endpoint (optional)
  - `contract-address`: Contract address (required)

#### 3. Get Contract Storage Data
- Endpoint: `?query=evm-contract-data-at-memory`
- Parameters:
  - `chain-id`: Chain ID (required)
  - `json-rpc`: JSON-RPC endpoint (optional)
  - `contract-address`: Contract address (required)
  - `storage-at`: Storage slot (required)

#### 4. Call Contract View Function
- Endpoint: `?query=evm-contract-call-view`
- Parameters:
  - `chain-id`: Chain ID (required)
  - `json-rpc`: JSON-RPC endpoint (optional)
  - `contract-address`: Contract address (required)
  - `method-name`: Function name (optional)
  - `method-inputs`: Array of parameter objects (optional)
    - Each parameter object contains:
      - `type`: Parameter type
      - `value`: Parameter value

#### 5. Get Contract Balance
- Endpoint: `?query=get-contract-balance`
- Parameters:
  - `chain-id`: Chain ID (required)
  - `json-rpc`: JSON-RPC endpoint (optional)
  - `address`: Contract address (required)

#### 6. Get Version
- Endpoint: `?query=version`
- No additional parameters required

## Example Usage

Examples of API usage can be found in `test/test.ts`. Here's a basic example:

```typescript
// Example of calling a view function
const response = await fetch(
  'generic-evm-api-go.vercel.app/api?query=evm-contract-call-view&chain-id=1&contract-address=0x123...&method-name=balanceOf&method-inputs=[{"type":"address","value":"0x456..."}]'
);
const data = await response.json();
```

## Hosting

You can either:
1. Use the hosted version at `generic-evm-api-go.vercel.app`
2. Self-host by running the Go server locally

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)

This project is licensed under the GNU General Public License v3.0 (GPL-3.0) - see [LICENSE](LICENSE) for details.

Copyright (c) 2025