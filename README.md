# Payment service raw v0.0.1

**Payment service** - A wrapper on Golang for accepting payments on your website that can generate payments and send them to Postgres, and there's a stint on a message broker like Kafka or Redis.

- Interface of link generation as [GRPC](https://grpc.io/), for webhook processing using is [CHI](https://github.com/go-chi/chi)

## ENV VARIABLE'S LIST:
1. `POSTGRES-PASSWORD` - NEED FOR POSTGRES-PASSWORD *xd*
2. `CONFIG_PATH` or `--CONFIG-PATH` - Config path *xd*

## What's service can do:
- Full abstaction for making worker which will produce message for broker
- Easy integration with payment service. Just move `internal-govnokassa` to your implementation and given the name of the functions like:
- - `func (g *Govnokassa) GeneratePaymentURL(data *models.DBPayment) (string, error)` - for generate payment link
  - `func (g *Govnokassa) ValidateData(rawData []byte) (*GovnoPayment, error)` - for validate incoming webhook from thirdy-payment service

*really dont know for what i write this, anyway maybe its be useful*