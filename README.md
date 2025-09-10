# Payment Service

<p align="left">
  <a href="#"><img src="https://img.shields.io/badge/go-1.24.1+-blue.svg" alt="Go Version"></a>
  <a href="https://github.com/YOUR_USERNAME/payment-service/blob/main/LICENSE"><img src="https://img.shields.io/github/license/YOUR_USERNAME/payment-service" alt="License"></a>
  <a href="https://github.com/YOUR_USERNAME/payment-service/actions"><img src="https://img.shields.io/github/actions/workflow/status/YOUR_USERNAME/payment-service/go.yml?branch=main" alt="Build Status"></a>
</p>

A production-ready boilerplate for a payment processing microservice written in Golang. This service is designed to be a lightweight, reliable backend component for accepting payments, generating unique payment URLs, and processing webhooks from third-party gateways.

It uses a clean, extensible architecture perfect for a microservice environment, featuring gRPC for internal APIs and a RESTful endpoint for asynchronous callbacks.

## ✨ Core Features

- **🚀 High-Performance API:** Exposes a **gRPC** endpoint for fast, internal service-to-service communication to create payment requests.
- **🔌 Reliable Webhook Processing:** Uses a robust **chi** router to handle incoming webhooks from payment providers.
- **📦 Guaranteed Event Delivery:** Implements the **Outbox Pattern** to ensure that events (like `payment_successful`) are reliably published to a message broker (**Kafka** or **Redis**) even under high load or in case of transient failures.
- **💾 Persistent & Scalable:** Leverages **PostgreSQL** for data storage, with a clean database schema and integrated migrations.
- **🐳 Fully Containerized:** Comes with a complete `docker-compose` setup for one-command local development.
- **✅ Extensible Provider System:** Easily integrate any payment gateway by implementing a simple provider interface.

## 🏗️ Architectural Highlights

The service follows a standard microservice pattern:

1.  An internal service calls the `payment-service` via **gRPC** to create a payment invoice.
2.  The service generates a unique payment URL from the integrated payment provider and stores the transaction details in **Postgres**.
3.  When the user pays, the third-party provider sends a webhook to the service's **HTTP endpoint**.
4.  The webhook is validated, the database state is updated, and a success/failure event is written to an `outbox` table in the same transaction.
5.  A separate worker process (**Outbox Worker**) reads from this table and reliably pushes the event to **Kafka/Redis**, ensuring at-least-once delivery.

## 🔧 Getting Started (Local Development)

### Prerequisites
- Go 1.24.1+
- Docker & Docker Compose
- [go-task](https://taskfile.dev/installation/)

### Installation & Launch

1.  **Install `go-task`:**

    *   **macOS:** `brew install go-task`
    *   **Linux/Other:** `sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d && sudo mv ./bin/task /usr/local/bin/`

2.  **Clone and navigate to the project root:**

    ```bash
    git clone https://github.com/YOUR_USERNAME/payment-service.git
    cd payment-service
    ```

3.  **Launch the development environment (Postgres, Redis, etc.):**

    > This command spins up all necessary services defined in `docker-compose.yml`.
    ```bash
    task dev-env
    ```

4.  **Start the application server:**
    > The app will connect to the services launched in the previous step.
    ```bash
    task dev-app
    ```

## 🔌 How to Add a New Payment Provider

The service is designed to be easily extendable. To add a new provider (e.g., `NewPay`):

1.  Create a new directory: `internal/app/providers/newpay`.
2.  Inside, implement the `Provider` interface, which requires two main functions:
    ```go
    // GeneratePaymentURL creates a unique link for the customer.
    GeneratePaymentURL(data *models.DBPayment) (string, error)

    // ValidateData parses and validates the incoming webhook from the provider.
    ValidateData(rawData []byte) (*ValidatedPaymentData, error)
    ```
3.  Register your new provider in the main application service.

## 🗺️ Roadmap

Here is the current status and future plans for the project.

#### ✅ Completed
- [x] Middleware for gRPC
- [x] Middleware for Webhook security (basic)
- [x] Improved logging and structured error handling
- [x] Database migrations setup
- [x] Full Docker Compose development environment

#### 🛠️ Up Next
- [ ] IP Limiter middleware for gRPC and HTTP endpoints
- [ ] Implement a worker pool for the Outbox publisher for better throughput
- [ ] Add CI/CD pipeline (e.g., GitHub Actions for build/test)
- [ ] Add OpenTelemetry for distributed tracing