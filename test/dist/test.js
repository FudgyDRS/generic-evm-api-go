"use strict";
// test.ts
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    var desc = Object.getOwnPropertyDescriptor(m, k);
    if (!desc || ("get" in desc ? !m.__esModule : desc.writable || desc.configurable)) {
      desc = { enumerable: true, get: function() { return m[k]; } };
    }
    Object.defineProperty(o, k2, desc);
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || (function () {
    var ownKeys = function(o) {
        ownKeys = Object.getOwnPropertyNames || function (o) {
            var ar = [];
            for (var k in o) if (Object.prototype.hasOwnProperty.call(o, k)) ar[ar.length] = k;
            return ar;
        };
        return ownKeys(o);
    };
    return function (mod) {
        if (mod && mod.__esModule) return mod;
        var result = {};
        if (mod != null) for (var k = ownKeys(mod), i = 0; i < k.length; i++) if (k[i] !== "default") __createBinding(result, mod, k[i]);
        __setModuleDefault(result, mod);
        return result;
    };
})();
Object.defineProperty(exports, "__esModule", { value: true });
const axios_1 = __importStar(require("axios"));
// Configuration
const BASE_URL = 'http://localhost:8080/api/api';
const CHAIN_ID = '56'; // BSC Mainnet
const RPC_URL = 'https://binance.llamarpc.com';
const CONTRACTS = {
    USDC: '0x8965349fb649A33a30cbFDa057D8eC2C48AbE2A2',
    PANCAKE_FACTORY: '0xcA143Ce32Fe78f1f7019d7d551a6402fC5350c73',
};
function formatError(error) {
    const message = error.response?.data?.message || error.response?.data?.error || error.message || 'Unknown error';
    return `${error.response?.status || 'Unknown'} Error: ${message}`;
}
// Test version endpoint
async function testVersion() {
    try {
        console.log('Testing API Version...');
        const response = await axios_1.default.get(BASE_URL, {
            params: {
                query: 'version',
            },
        });
        console.log('✅ API Version:', response.data.version);
        console.log('✅ Status Code:', response.status);
    }
    catch (error) {
        if (error instanceof axios_1.AxiosError) {
            console.error('❌ API Error:', error.response?.data || error.message);
            console.error('Status Code:', error.response?.status);
        }
        else {
            console.error('❌ Unexpected error:', error);
        }
        process.exit(1);
    }
}
// Helper function for API calls
async function makeRequest(query, params) {
    try {
        const response = await axios_1.default.get(BASE_URL, {
            params: {
                query,
                ...params,
            },
        });
        return response.data;
    }
    catch (error) {
        if (error instanceof axios_1.AxiosError) {
            console.error(`Error in ${query}:`, error.response?.data || error.message);
        }
        else {
            console.error(`Unexpected error in ${query}:`, error);
        }
        throw error;
    }
}
// Test functions
async function testExtCodeSize(address) {
    console.log(`\nTesting getExtCodeSize for ${address}`);
    const response = await makeRequest('evm-contract-ext-code-size', {
        'chain-id': CHAIN_ID,
        'contract-address': address,
    });
    console.log('Contract size:', response['contract-size']);
}
async function testContractCode(address) {
    console.log(`\nTesting getContractCode for ${address}`);
    const response = await makeRequest('evm-contract-code', {
        'chain-id': CHAIN_ID,
        'contract-address': address,
    });
    console.log('Contract code length:', response['contract-code'].length);
}
async function testStorageSlot(address, slot) {
    console.log(`\nTesting storage slot ${slot} for ${address}`);
    const response = await makeRequest('evm-contract-data-at-memory', {
        'chain-id': CHAIN_ID,
        'contract-address': address,
        'storage-at': slot,
    });
    console.log('Storage data:', response.bytes);
}
async function testContractCall(address, methodName, params = []) {
    console.log(`\nTesting ${methodName} call for ${address}`);
    const response = await makeRequest('evm-contract-call-view', {
        'chain-id': CHAIN_ID,
        'contract-address': address,
        'method-name': methodName,
        'method-inputs': params,
    });
    console.log(`${methodName} response:`, response.response);
}
async function testBalance(address) {
    console.log(`\nTesting balance for ${address}`);
    const response = await makeRequest('get-contract-balance', {
        'chain-id': CHAIN_ID,
        address: address,
    });
    console.log('Balance:', response.balance);
}
const delay = (ms) => new Promise(resolve => setTimeout(resolve, ms));
// Main test suite
async function runTests() {
    console.log('Starting API Tests...\n');
    try {
        console.log('Testing API Version...');
        const response = await axios_1.default.get(BASE_URL, {
            params: {
                query: 'version',
            },
        });
        console.log('✅ API Version:', response.data.version);
        console.log('✅ Status Code:', response.status);
    }
    catch (error) {
        if (error instanceof axios_1.AxiosError) {
            console.error('❌ Error:', formatError(error));
            console.error('URL:', error.config?.url);
        }
        else {
            console.error('❌ Unexpected error:', error);
        }
    }
    try {
        // USDC Tests
        console.log('🪙 Testing USDC Contract');
        await testExtCodeSize(CONTRACTS.USDC);
        await testContractCode(CONTRACTS.USDC);
        await testBalance(CONTRACTS.USDC);
        await testContractCall(CONTRACTS.USDC, 'name');
        await testContractCall(CONTRACTS.USDC, 'symbol');
        await testContractCall(CONTRACTS.USDC, 'decimals');
        await testContractCall(CONTRACTS.USDC, 'totalSupply');
        // await testContractCall(CONTRACTS.USDC, 'balanceOf', [{ type: 'address', value: CONTRACTS.USDC }]);
    }
    catch (error) {
        if (error instanceof axios_1.AxiosError) {
            console.error('❌ Error:', formatError(error));
            console.error('URL:', error.config?.url);
        }
        else {
            console.error('❌ Unexpected error:', error);
        }
    }
    await delay(2000);
    try {
        // PancakeSwap Factory Tests
        console.log('\n🥞 Testing PancakeSwap Factory Contract');
        await testExtCodeSize(CONTRACTS.PANCAKE_FACTORY);
        await testContractCode(CONTRACTS.PANCAKE_FACTORY);
        await testBalance(CONTRACTS.PANCAKE_FACTORY);
        await testContractCall(CONTRACTS.PANCAKE_FACTORY, 'feeTo');
        await testContractCall(CONTRACTS.PANCAKE_FACTORY, 'feeToSetter');
        await testContractCall(CONTRACTS.PANCAKE_FACTORY, 'allPairsLength');
        await testContractCall(CONTRACTS.PANCAKE_FACTORY, 'getPair', [
            { type: 'address', value: '0xbA2aE424d960c26247Dd6c32edC70B295c744C43' },
            { type: 'address', value: '0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c' },
        ]);
        // Test some storage slots
        await testStorageSlot(CONTRACTS.PANCAKE_FACTORY, '0'); // feeTo
        await testStorageSlot(CONTRACTS.PANCAKE_FACTORY, '1'); // feeToSetter
    }
    catch (error) {
        if (error instanceof axios_1.AxiosError) {
            console.error('❌ Error:', formatError(error));
            console.error('URL:', error.config?.url);
        }
        else {
            console.error('❌ Unexpected error:', error);
        }
    }
    console.log('\n✅ All tests completed');
}
// Run the test suite
runTests();
