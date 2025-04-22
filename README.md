# Payment service raw v0.0.2

**Payment service** - A wrapper on Golang for accepting payments on your website that can generate payments and send them to Postgres, and there's a stint on a message broker like Kafka or Redis.

- Interface of link generation as [GRPC](https://grpc.io/), for webhook processing using is [CHI](https://github.com/go-chi/chi)

## Start local development
- Install go-task (required Go 1.17+)
- MacOS `brew install go-task` 
- Ubuntu / Debian`sudo apt install go-task`
- Arch Linux `sudo pacman -S go-task`

### Go to main directory
```sh
$ cd payment-service
```
### Run the go-task utility
- Make sure you download `docker-compose`
```sh
$ sudo go-task dev-env
```
### For start application
```sh
$ go-task dev-app
```

## What's service can do:
- Full abstaction for making worker which will produce message for broker
- Easy integration with payment service. Just move `internal-govnokassa` to your implementation and given the name of the functions like:
- - `func (g *Govnokassa) GeneratePaymentURL(data *models.DBPayment) (string, error)` - for generate payment link
  - `func (g *Govnokassa) ValidateData(rawData []byte) (*GovnoPayment, error)` - for validate incoming webhook from thirdy-payment service

## TODO
- TODO: ~~MIDDLEWARE FOR GRPC~~ IP LIMITER/BAN + WEBHOOK middleware
- TODO: ~~ЛУЧШЕ ЛОГИ + ОБРАБОТКА ОШИБОК~~
- TODO: ~~ПРОБЕЖАТЬСЯ ПО // TEMP~~
- TODO: ~~ТЕСТЫ~~ 
- TODO: ~~МИГРАЦИИ~~
- TODO: ~~README.MD получше~~ он идеален
- TODO: ~~DOCKER-COMPOSE OR AUTO-DEPLOY etc.~~ 
- TODO: CI/CD
- TODO: Make workerpool, гибкая обработка outbox-pattern'a
- TODO: ip-limiter for GRPC interceptors