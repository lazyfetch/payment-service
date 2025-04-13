# Payment service raw v0.0.2

**Payment service** - A wrapper on Golang for accepting payments on your website that can generate payments and send them to Postgres, and there's a stint on a message broker like Kafka or Redis.

- Interface of link generation as [GRPC](https://grpc.io/), for webhook processing using is [CHI](https://github.com/go-chi/chi)

## ENV VARIABLE'S LIST:
1. `POSTGRES-PASSWORD` - NEED FOR POSTGRES-PASSWORD
2. `CONFIG_PATH` or `--CONFIG-PATH` - Config path 

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
- TODO: DOCKER-COMPOSE OR AUTO-DEPLOY etc.

*really dont know for what i write this, anyway maybe its be useful*

