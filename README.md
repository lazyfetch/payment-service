# Payment service beta v0.0.1

**Payment service** — A lightweight Golang wrapper for accepting payments on your website.  
Generates payment links, stores data in Postgres, and pushes events to message brokers like Kafka or Redis.

- Link generation is exposed via [gRPC](https://grpc.io/)
- Webhook processing is handled through [CHI](https://github.com/go-chi/chi)

---

## 🔧 Start Local Development

### 1. Install go-task (requires Go 1.17+)

- **macOS:**
  ```bash
  brew install go-task/tap/go-task
```

- **Linux (Ubuntu/Debian/Arch/etc):**
    ```bash
    sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d
    sudo mv ./task /usr/local/bin/
    ```
### 2. Go to the root project directory

```bash
cd payment-service
```

### 3. Launch local development environment

- Ensure `docker` and `docker-compose` are installed.
    

```bash
go-task dev-env
```

### 4. Start the application

```bash
go-task dev-app
```

---

## 🚀 What This Service Can Do

- Full abstraction for spinning up workers that publish messages to a broker.
    
- Easy payment gateway integration: just copy `internal-govnokassa` and implement functions like:
    
    - `GeneratePaymentURL(data *models.DBPayment) (string, error)` — Generate payment link.
        
    - `ValidateData(rawData []byte) (*GovnoPayment, error)` — Validate webhooks from 3rd-party payment systems.
        

---

## 📋 TODO

-  MIDDLEWARE FOR GRPC ✅
    
-  IP LIMITER / BAN + Webhook middleware
    
-  Better logs + error handling
    
-  Cleanup `// TEMP` code
    
-  Database migrations
    
-  Improve this README (it's perfect now 😎)
    
-  docker-compose + deploy setup
    
-  CI/CD pipeline
    
-  Workerpool for outbox-pattern
    
-  gRPC interceptor IP limiter
