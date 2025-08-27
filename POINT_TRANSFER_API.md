# Point Transfer API Documentation

This document describes the point transfer system APIs that have been added to the Fiber REST API server.

## New Endpoints

### 1. Get Point Balance
**GET** `/points/balance`

Get the current user's point balance and LBK information.

**Headers:**
- `Authorization: Bearer <jwt_token>`

**Response:**
```json
{
  "lbk_code": "LBK001234",
  "point_balance": 1000,
  "first_name": "John",
  "last_name": "Doe"
}
```

### 2. Search User by LBK Code
**GET** `/users/search?lbk_code=LBK001234`

Search for a user by their LBK code to get their basic information for transfers.

**Headers:**
- `Authorization: Bearer <jwt_token>`

**Query Parameters:**
- `lbk_code` (required): The LBK code to search for

**Response:**
```json
{
  "lbk_code": "LBK001234",
  "first_name": "John",
  "last_name": "Doe"
}
```

### 3. Transfer Points
**POST** `/points/transfer`

Transfer points from the current user to another user identified by their LBK code.

**Headers:**
- `Authorization: Bearer <jwt_token>`
- `Content-Type: application/json`

**Request Body:**
```json
{
  "to_lbk_code": "LBK001234",
  "amount": 100,
  "message": "Optional transfer message"
}
```

**Response:**
```json
{
  "transfer_id": 1,
  "message": "Transfer completed successfully",
  "from_user": {
    "lbk_code": "LBK001235",
    "first_name": "Jane",
    "last_name": "Smith"
  },
  "to_user": {
    "lbk_code": "LBK001234",
    "first_name": "John",
    "last_name": "Doe"
  },
  "amount": 100,
  "status": "completed"
}
```

### 4. Get Transfer History
**GET** `/points/history`

Get the transfer history for the current user (both sent and received transfers).

**Headers:**
- `Authorization: Bearer <jwt_token>`

**Response:**
```json
{
  "transfers": [
    {
      "id": 1,
      "from_user_id": 2,
      "to_user_id": 1,
      "from_user": {
        "id": 2,
        "lbk_code": "LBK001235",
        "first_name": "Jane",
        "last_name": "Smith"
      },
      "to_user": {
        "id": 1,
        "lbk_code": "LBK001234",
        "first_name": "John",
        "last_name": "Doe"
      },
      "amount": 100,
      "message": "Transfer message",
      "status": "completed",
      "created_at": "2025-08-27T14:30:00Z"
    }
  ],
  "count": 1
}
```

## Updated User Model

The User model has been updated to include:
- `lbk_code`: Unique LBK identification code (automatically generated)
- `point_balance`: Current point balance (new users start with 1000 points)

## New Database Tables

### Transfer Table
Stores all point transfer transactions with the following fields:
- `id`: Primary key
- `from_user_id`: ID of the sender
- `to_user_id`: ID of the recipient
- `amount`: Number of points transferred
- `message`: Optional transfer message
- `status`: Transfer status (completed, failed, pending)
- `created_at`, `updated_at`: Timestamps

## Example Usage

### 1. Check Your Point Balance
```bash
curl -H "Authorization: Bearer <your_jwt_token>" \
     http://localhost:3000/points/balance
```

### 2. Search for a User
```bash
curl -H "Authorization: Bearer <your_jwt_token>" \
     "http://localhost:3000/users/search?lbk_code=LBK001234"
```

### 3. Transfer Points
```bash
curl -X POST \
     -H "Authorization: Bearer <your_jwt_token>" \
     -H "Content-Type: application/json" \
     -d '{"to_lbk_code":"LBK001234","amount":100,"message":"Thanks!"}' \
     http://localhost:3000/points/transfer
```

### 4. View Transfer History
```bash
curl -H "Authorization: Bearer <your_jwt_token>" \
     http://localhost:3000/points/history
```

## Error Responses

All endpoints return appropriate HTTP status codes and error messages:

- `400 Bad Request`: Invalid request data
- `401 Unauthorized`: Missing or invalid JWT token
- `404 Not Found`: User or resource not found
- `500 Internal Server Error`: Server-side error

Example error response:
```json
{
  "error": "Insufficient points"
}
```

## Security Features

- All point transfer endpoints require JWT authentication
- Transfers are protected by database transactions to ensure consistency
- Users cannot transfer points to themselves
- Point balances cannot go negative
- Transfer history is limited to 50 recent transactions per request
