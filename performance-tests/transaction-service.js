import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend } from 'k6/metrics';


// Custom metrics
const errorRate = new Rate('errors');
const transactionDuration = new Trend('transaction_duration');

export const options = {
    vus: 50,
    duration: '10s',
    thresholds: {
        http_req_duration: ['p(95)<500'], // 95% of requests must complete below 500ms
        errors: ['rate<0.1'],             // Error rate must be less than 10%
    },
};

const BASE_URL = 'http://localhost:8081';


// Initialize base account IDs
// Use microsecond precision timestamp plus random offset to avoid collisions
var baseAccountId = Date.now();

// Helper to create test accounts with unique incrementing IDs
function getUniqueAccountId() {

    // Use VU ID and iteration count to create unique offsets
    const vuOffset = __VU * 1000; // Multiply by 1000 to avoid collisions between VUs
    const iterOffset = __ITER * 2; // Multiply by 2 since we need 2 IDs per iteration
    return {
        sourceId: baseAccountId + vuOffset + iterOffset,
        destId: baseAccountId + vuOffset + iterOffset + 1
    };
}

function setupTestAccounts() {
    const accountService = 'http://localhost:8080';

    const { sourceId, destId } = getUniqueAccountId();

    // Create source account with unique integer ID
    const sourceAccount = {
        account_id: sourceId,
        initial_balance: 10000, // $100.00 initial balance
    };

    // Create destination account with unique integer ID
    const destAccount = {
        account_id: destId,
        initial_balance: 10000, // $100.00 initial balance
    };

    // Create source account
    let sourceAccountCreated = false;
    let destAccountCreated = false;

    // Create source account first
    while (!sourceAccountCreated) {
        const sourceResponse = http.post(
            `${accountService}/accounts`, 
            JSON.stringify(sourceAccount),
            { headers: { 'Content-Type': 'application/json' } }
        );

        if (sourceResponse.status === 201) {
            sourceAccountCreated = true;
        }
        // Source account exists, try new ID
        else if (sourceResponse.status === 400 && sourceResponse.body.includes("already exists")) {
            sourceAccountCreated = true; // If account exists, we can use it as source
        }
        // Log other source account errors
        else {
            console.log(`Error creating source account: ${sourceResponse.status} - ${sourceResponse.body}`);
        }

    }

    // Create destination account after source is created
    while (!destAccountCreated) {
        const destResponse = http.post(
            `${accountService}/accounts`,
            JSON.stringify(destAccount),
            { headers: { 'Content-Type': 'application/json' } }
        );

        if (destResponse.status === 201) {
            destAccountCreated = true;
        }
        // Destination account exists, try new ID
        else if (destResponse.status === 400 && destResponse.body.includes("already exists")) {
            destAccountCreated = true; // If account exists, we can use it as destination
        }
        // Log other destination account errors
        else {
            console.log(`Error creating destination account: ${destResponse.status} - ${destResponse.body}`);
        }

    }

    return [sourceAccount, destAccount];
}

export default function () {
    const [sourceAccount, destAccount] = setupTestAccounts();
    sleep(0.5);

    // Create transaction
    const startTime = new Date();
    
    const transactionData = {
        source_account_id: sourceAccount.account_id,
        destination_account_id: destAccount.account_id,
        amount: 1, // $0.01
    };

    let createResponse;
    try {
        createResponse = http.post(
            `${BASE_URL}/transactions`,
            JSON.stringify(transactionData),
            { headers: { 'Content-Type': 'application/json' } }
        );
    } catch (error) {
        console.log(`Error creating transaction: ${error}`);
    }

    const duration = new Date() - startTime;
    transactionDuration.add(duration);

    const checkResult = check(createResponse, {
        'transaction created successfully': (r) => r.status === 201,
        'response time OK': (r) => r.timings.duration < 500,
        'valid response body': (r) => {
            try {
                const body = JSON.parse(r.body);
                return body.id && body.status;
            } catch (e) {
                return false;
            }
        },
    });

    if (!checkResult && createResponse.status !== 201) {
        console.log(`Error response: ${createResponse.status} - ${createResponse.body}`);
        errorRate.add(1);
    }

    sleep(1);
}
