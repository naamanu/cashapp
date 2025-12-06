# Infrastructure Changelog

## Overview

This branch (`improvement/infrastructure`) focuses on modernizing the underlying infrastructure and tooling of the CashApp service.

## 1. Go Version Upgrade

- **Old Version**: Go 1.15
- **New Version**: Go 1.23+
- **Reason**: Go 1.15 is end-of-life. Upgrading allows us to use modern language features (generics, structured logging optimization), improved garbage collection, and better module support.

## 2. Structured Logging

- **Tool**: `go.uber.org/zap`
- **Change**: replaced the standard library `log` package with Zap.
- **Implementation**:
  - Created `core/logger.go` to initialize a global `Log` variable.
  - **Development Mode**: Uses human-readable, colorized output (`zap.NewDevelopmentConfig`).
  - **Production Mode**: Uses JSON formatted output (`zap.NewProductionConfig`) for easy ingestion by log aggregators (ELK, Datadog, etc.).

## 3. Configuration Management

- **Tool**: `spf13/viper`
- **Change**: Replaced manual `os.Getenv` and `godotenv` with Viper.
- **Features**:
  - **Defaults**: Secure defaults are set for all variables (e.g., `PG_SSLMODE=disable` for dev).
  - **Env Vars**: Automatically reads environment variables.
  - **Dotenv**: Still supports `.env` files for local development.
  - **Type Safety**: Unmarshals configuration directly into a highly-typed `Config` struct using `mapstructure`.

## 4. Docker Improvements

- Updated `Dockerfile` to use multi-stage builds with `golang:1.23-alpine` as the builder and a minimal `alpine:latest` as the runner.
- Added `swag init` to the build process to ensure API documentation is always up-to-date with the code.
