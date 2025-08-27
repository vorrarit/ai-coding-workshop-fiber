# Entity-Relationship Diagram

This document contains the Entity-Relationship diagram for the Fiber API database schema.

## Database Schema Overview

The system consists of two main entities:
- **Users**: Manages user accounts with authentication and point balances
- **Transfers**: Tracks point transfer transactions between users

## ER Diagram

```plantuml
@startuml
!define ENTITY entity
!define PK <color:red><b>PK</b></color>
!define FK <color:blue><b>FK</b></color>
!define UK <color:green><b>UK</b></color>
!define NN <color:orange><b>NN</b></color>

ENTITY User {
  PK id : UINT <<generated>>
  UK email : VARCHAR(255) NN
  password : VARCHAR(255) NN
  first_name : VARCHAR(255) NN
  last_name : VARCHAR(255) NN
  phone_number : VARCHAR(255)
  dob : DATETIME
  UK lbk_code : VARCHAR(255) NN
  point_balance : UINT <<default: 0>>
  created_at : DATETIME NN
  updated_at : DATETIME NN
}

ENTITY Transfer {
  PK id : UINT <<generated>>
  FK from_user_id : UINT NN
  FK to_user_id : UINT NN
  amount : UINT NN
  message : TEXT
  status : VARCHAR(50) <<default: 'completed'>>
  created_at : DATETIME NN
  updated_at : DATETIME NN
}

User ||--o{ Transfer : "from_user_id"
User ||--o{ Transfer : "to_user_id"

note top of User : Primary entity for user management\nStores authentication data and point balances
note top of Transfer : Transaction log for point transfers\nAudit trail with complete transfer history

note right of User::lbk_code : Unique identification code\nUsed for searching and transferring points
note right of User::point_balance : Current available points\nUpdated through transfer transactions
note right of Transfer::status : Transfer status\nValues: completed, failed, pending

@enduml
```

## Entity Descriptions

### User Entity
The User entity represents registered users in the system.

**Primary Key**: `id` (Auto-incrementing unsigned integer)

**Unique Constraints**:
- `email` - Ensures each user has a unique email address
- `lbk_code` - Unique LBK identification code for point transfers

**Required Fields** (NOT NULL):
- `email` - User's email address for authentication
- `password` - Hashed password for security
- `first_name` - User's first name
- `last_name` - User's last name
- `lbk_code` - Unique LBK identification code
- `created_at` - Record creation timestamp
- `updated_at` - Last modification timestamp

**Optional Fields**:
- `phone_number` - User's contact number
- `dob` - User's date of birth

**Default Values**:
- `point_balance` - Defaults to 0, represents available points

### Transfer Entity
The Transfer entity tracks all point transfer transactions between users.

**Primary Key**: `id` (Auto-incrementing unsigned integer)

**Foreign Keys**:
- `from_user_id` - References User.id (sender)
- `to_user_id` - References User.id (recipient)

**Required Fields** (NOT NULL):
- `from_user_id` - ID of the user sending points
- `to_user_id` - ID of the user receiving points  
- `amount` - Number of points transferred
- `created_at` - Transaction timestamp
- `updated_at` - Last modification timestamp

**Optional Fields**:
- `message` - Optional message attached to transfer

**Default Values**:
- `status` - Defaults to 'completed'

## Relationships

### User → Transfer (One-to-Many)
- **From Relationship**: One User can send multiple Transfers (`User.id` → `Transfer.from_user_id`)
- **To Relationship**: One User can receive multiple Transfers (`User.id` → `Transfer.to_user_id`)
- **Cardinality**: 1:N for both relationships
- **Referential Integrity**: Foreign key constraints ensure transfer references exist

## Business Rules

1. **User Registration**:
   - Email must be unique across the system
   - LBK code must be unique for point transfer identification
   - New users start with 1000 points

2. **Point Transfers**:
   - Users cannot transfer points to themselves
   - Sender must have sufficient point balance
   - All transfers are logged for audit purposes
   - Transfers are atomic (both balances updated or transaction fails)

3. **Data Integrity**:
   - User deletion should be handled carefully due to transfer references
   - Transfer records should be preserved for audit trail
   - Point balances must always be non-negative

## Database Indexes

**Recommended Indexes** (for performance):
- `users.email` - Unique index for login queries
- `users.lbk_code` - Unique index for user search
- `transfers.from_user_id` - Index for user's sent transfers
- `transfers.to_user_id` - Index for user's received transfers
- `transfers.created_at` - Index for chronological queries

## Schema Evolution

The current schema supports:
- ✅ User authentication and profile management
- ✅ Point balance tracking
- ✅ Point transfer transactions
- ✅ Transfer audit trail
- ✅ User search by LBK code

**Future Considerations**:
- Additional user profile fields
- Transfer categories or types
- Point earning transactions
- User roles and permissions
- Transfer limits and rules
