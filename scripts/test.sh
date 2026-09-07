#!/bin/bash
set -e

API_KEY="${API_KEY:-test}"

HOST="${HOST:-localhost}"
PORT="${PORT:-2525}"

BASE_URL="https://$HOST:$PORT"

# ---

echo "test 01"
echo ""

curl --cacert ./certs/cert.pem "$BASE_URL/livez"
curl --cacert ./certs/cert.pem "$BASE_URL/readyz"
curl --cacert ./certs/cert.pem "$BASE_URL/healthz"

# ---

echo ""
echo "test 02"
echo ""

curl --cacert ./certs/cert.pem "$BASE_URL/api/v1/rate"
curl --cacert ./certs/cert.pem -H "X-API-Key: wrong" "$BASE_URL/api/v1/rate"

# ---

echo ""
echo "test 03"
echo ""

curl --cacert ./certs/cert.pem -H "X-API-Key: $API_KEY" "$BASE_URL/api/v1/rate"

echo "---"

curl --cacert ./certs/cert.pem -H "X-API-Key: $API_KEY" "$BASE_URL/api/v1/rate?fromCurrency=US&toCurrency=EUR"
curl --cacert ./certs/cert.pem -H "X-API-Key: $API_KEY" "$BASE_URL/api/v1/rate?fromCurrency=USD&toCurrency=EU"

echo "---"

curl --cacert ./certs/cert.pem -H "X-API-Key: $API_KEY" "$BASE_URL/api/v1/rate?fromCurrency=USD&toCurrency=EU1"
curl --cacert ./certs/cert.pem -H "X-API-Key: $API_KEY" "$BASE_URL/api/v1/rate?fromCurrency=US1&toCurrency=EUR"

# ---

echo ""
echo "test 04"
echo ""

curl --cacert ./certs/cert.pem -H "X-API-Key: $API_KEY" "$BASE_URL/api/v1/rates"

echo "---"

curl --cacert ./certs/cert.pem -H "X-API-Key: $API_KEY" "$BASE_URL/api/v1/rates?baseCurrency=US"
curl --cacert ./certs/cert.pem -H "X-API-Key: $API_KEY" "$BASE_URL/api/v1/rates?baseCurrency=US1"

# ---

echo ""
echo "test 05"
echo ""

curl --cacert ./certs/cert.pem -X POST -H "X-API-Key: $API_KEY" "$BASE_URL/api/v1/convert"

echo "---"

curl --cacert ./certs/cert.pem -X POST -H "X-API-Key: $API_KEY" -H "Content-Type: application/json" -d '{"toCurrency":"EUR","amount":100}' "$BASE_URL/api/v1/convert"
curl --cacert ./certs/cert.pem -X POST -H "X-API-Key: $API_KEY" -H "Content-Type: application/json" -d '{"fromCurrency":"USD","amount":100}' "$BASE_URL/api/v1/convert"
curl --cacert ./certs/cert.pem -X POST -H "X-API-Key: $API_KEY" -H "Content-Type: application/json" -d '{"fromCurrency":"USD","toCurrency":"EUR"}' "$BASE_URL/api/v1/convert"

echo "---"

curl --cacert ./certs/cert.pem -X POST -H "X-API-Key: $API_KEY" -H "Content-Type: application/json" -d '{"fromCurrency":"US","toCurrency":"EUR","amount":100}' "$BASE_URL/api/v1/convert"
curl --cacert ./certs/cert.pem -X POST -H "X-API-Key: $API_KEY" -H "Content-Type: application/json" -d '{"fromCurrency":"USD","toCurrency":"EU","amount":100}' "$BASE_URL/api/v1/convert"

echo "---"

curl --cacert ./certs/cert.pem -X POST -H "X-API-Key: $API_KEY" -H "Content-Type: application/json" -d '{"fromCurrency":"US1","toCurrency":"EUR","amount":100}' "$BASE_URL/api/v1/convert"
curl --cacert ./certs/cert.pem -X POST -H "X-API-Key: $API_KEY" -H "Content-Type: application/json" -d '{"fromCurrency":"USD","toCurrency":"EU1","amount":100}' "$BASE_URL/api/v1/convert"

echo "---"

curl --cacert ./certs/cert.pem -X POST -H "X-API-Key: $API_KEY" -H "Content-Type: application/json" -d '{"fromCurrency":"USD","toCurrency":"EUR","amount":-1}' "$BASE_URL/api/v1/convert"
curl --cacert ./certs/cert.pem -X POST -H "X-API-Key: $API_KEY" -H "Content-Type: application/json" -d '{"fromCurrency":"USD","toCurrency":"EUR","amount":1000000000000000000000000000000000000000000000000000000000000000000}' "$BASE_URL/api/v1/convert"

# ---

echo ""
echo "test 06"
echo ""

curl --cacert ./certs/cert.pem -H "X-API-Key: $API_KEY" "$BASE_URL/api/v1/rate?fromCurrency=USD&toCurrency=EUR"
curl --cacert ./certs/cert.pem -H "X-API-Key: $API_KEY" "$BASE_URL/api/v1/rates?baseCurrency=USD"
curl --cacert ./certs/cert.pem -X POST -H "X-API-Key: $API_KEY" -H "Content-Type: application/json" -d '{"fromCurrency":"USD","toCurrency":"EUR","amount":100}' "$BASE_URL/api/v1/convert"
