# Architecture & Migration Report

## Overview

This document details the architectural transformation of the CashApp codebase from a monolithic service into a distributed microservice architecture. The goal of this refactor was to improve scalability, enforce domain boundaries, and modernize the infrastructure stack.

## 1. Architectural Changes

### Microservice Split

The application has been split into two distinct services:

1.  **User Service** (`cmd/user`)

    - **Responsibility**: Manages Identity (Users) and Wallets.
    - **Port**: `5454`
    - **Database**: Owns `users` and `wallets` tables.
    - **Endpoints**: `/users`, `/users/:id/wallet`

2.  **Ledger Service** (`cmd/ledger`)
    - **Responsibility**: Manages core banking logic, Transactions, and the Ledger.
    - **Port**: `5455`
    - **Database**: Owns `transactions` and `transaction_events` tables.
    - **Endpoints**: `/payments`, `/wallets/:id/balance`

### Decoupling Logic

Previously, the `PaymentService` directly joined `User` and `Wallet` tables in SQL queries. This has been decoupled:

- The Ledger Service no longer imports User models.
- **WalletLookup Interface**: A new `WalletLookupRepo` interface was created in the Ledger service. Currently, it still reads from the shared database (Phase 1 of migration), but it is structurally ready to be replaced by a gRPC/HTTP client to the User Service in Phase 2.

## 2. Infrastructure Modernization

- **Go 1.23**: Upgraded from Go 1.15 to the latest stable Go 1.23 for better performance and language features.
- **Structured Logging (Zap)**: Replaced the standard `log` package with `uber-go/zap` for high-performance, structured JSON logging.
- **Configuration (Viper)**: Implemented `spf13/viper` for robust configuration management, supporting environment variables, `.env` files, and defaults.
- **Swagger/OpenAPI**: Individual Swagger documentation is now generated for each microservice.

## 3. Directory Structure

The project has been reorganized to follow the Standard Go Project Layout:

```
├── cmd/
│   ├── user/        # Entrypoint for User Service
│   └── ledger/      # Entrypoint for Ledger Service
├── internal/
│   ├── user/        # Private application logic for User Domain
│   └── ledger/      # Private application logic for Ledger Domain
├── core/            # Shared kernels (Config, Logger, Common Types)
├── Dockerfile.user  # Docker build for User Svc
├── Dockerfile.ledger # Docker build for Ledger Svc
└── docker-compose.yml # Orchestration
```

## 4. Key Improvements & Bug Fixes

- **Transaction Processor Fix**: Identified and fixed a critical bug in `processTransaction` where `Withdrawal` logic was unreachable due to a duplicate `switch` case for `Transfer`.
- **Type Safety**: Introduced specific types (`Status`, `Direction`, `Purpose`) in `core/types.go` to replace loose string typing.
- **Dependnecy Injection**: Both services now use a clear Repository/Service pattern with interface-based dependency injection, making testing easier.

## 5. How to Run

### Docker (Recommended)

Spin up the entire stack (Postgres, Redis, User Svc, Ledger Svc):

```bash
make docker-build
make docker-up
```

### Local Development

You can run services individually (requires running Postgres/Redis separately):

```bash
# Terminal 1
make run-user

# Terminal 2
make run-ledger
```

### Documentation

- User API: http://localhost:5454/swagger/index.html
- Ledger API: http://localhost:5455/swagger/index.html
