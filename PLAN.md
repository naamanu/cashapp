# Microservice Migration Plan for CashApp

## 1. Executive Summary

The current application is a monolithic Go application handling both User Management and Payment Processing. To scale this effectively and allow independent deployment cycles, we propose splitting this into two primary microservices:

1. **User Service**: Handles identity, user profiles, and wallet provision.
2. **Ledger Service (Payment)**: Handles transaction processing, balances, and history.

## 2. Infrastructure & Modernization Improvements

Before splitting, we recommend modernizing the foundation:

- **Upgrade Go Version**: Current version 1.15 is end-of-life. Upgrade to **Go 1.22+**.
- **Structured Logging**: Replace standard `log` with `uber-go/zap` or `rs/zerolog` for better observability in distributed systems.
- **Configuration**: Use `viper` or `godotenv` with strict schema validation for environment variables.
- **API Documentation**: Continue using Swagger, but generate separate specs for each service.

## 3. Architecture Design

### A. Service Boundaries

#### User Service

- **Responsibilities**:
  - User Registration/Login.
  - KYC (if applicable).
  - Wallet creation (assigning a wallet ID to a user).
  - Storing `User` and `Wallet` entities.
- **Database**: `cashapp_users` (Postgres)
- **API Exposed**:
  - `POST /users`: Create user.
  - `GET /users/:id/wallet`: Get primary wallet ID for a user.

#### Ledger Service (Payment)

- **Responsibilities**:
  - Processing Transfers, Deposits, Withdrawals.
  - Maintaining the "Double Entry" Ledger (`transaction_events`).
  - managing `Transaction` state.
- **Database**: `cashapp_ledger` (Postgres)
- **API Exposed**:
  - `POST /payments/send`: Initiate transfer.
  - `GET /payments/history`: Transaction history.
  - `GET /wallets/:id/balance`: Get current balance.

### B. Communication Strategy

- **Synchronous (RPC/HTTP)**:
  - The Ledger Service needs to validate if a sender/receiver exists. It should call the User Service via gRPC (preferred for internal) or internal HTTP REST.
- **Asynchronous (Events)**:
  - Use a message broker (RabbitMQ/Kafka/Redis PubSub) for non-critical side effects (e.g., ensuring a "Welcome Email" is sent after wallet creation, or analytical data syncing).

### C. Data Migration Strategy

The hardest part is untangling the database.

1. **Phase 1 (Code Split)**: Refactor code into modular structure (`/services/user`, `/services/ledger`) within the repo, enforcing boundaries (no direct DB joins across domains).
2. **Phase 2 (Logical Split)**: Deploy two instances of the app, one as "User" and one as "Ledger", both pointing to the same DB but only using their respective tables.
3. **Phase 3 (Physical Split)**: Migrate tables to separate databases.

## 4. Handling Distributed Transactions

Currently, `MoveMoneyBetweenWallets` runs in a single ACID transaction. When splitting:

- **Single Service Domain**: Luckily, `Transaction` and `TransactionEvent` (debit/credit) both live in the Ledger Service. The money movement remains ACID compliant within the Ledger Service.
- **Cross-Service Validation**: The Ledger service only needs the _ID_ of the wallet. It trusts the ID exists or verifies it via API. It does _not_ need to lock the User row. This simplifies things significantly.

## 5. Deployment (Docker & Kubernetes)

- Create a `Dockerfile.user` and `Dockerfile.ledger`.
- Use **Docker Compose** for local dev with 2 DB containers and 2 Service containers.
- Use **Kubernetes (Helm)** for production, adding an Ingress Controller (Nginx) to route `/users/*` to User Svc and `/payments/*` to Ledger Svc.

## 6. Next Steps

1. Upgrade `go.mod`.
2. Refactor directory structure to separate `apps/` or `services/`.
3. Abstract the `Repo` layer so `PaymentService` calls a `UserClient` interface instead of `UserRepo`.
