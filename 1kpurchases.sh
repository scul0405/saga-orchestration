#!/bin/bash

API_ENDPOINT="http://localhost:80/api/v1/purchases"
DATA='{
    "order_items": [
        {
            "product_id": 505340720819732500,
            "quantity": 1
        },
        {
            "product_id": 505340722648449044,
            "quantity": 1
        }
    ],
    "payment": {
        "currency_code": "VND"
    }
}'

for ((i=1; i<=1000; i++)); do
    curl --location --request POST "$API_ENDPOINT" \
    --header 'Content-Type: application/json' \
    --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJDdXN0b21lcklEIjo1MDUzNDA2NjUwMzU0ODkyOTEsIlJlZnJlc2giOmZhbHNlLCJleHAiOjE3MTA3MzY1OTd9.2GUIqEDZUMMC1_xle1JHBBFt3DoMNO5ummB-14X-HGM' \
    --data-raw "$DATA" &
done

wait