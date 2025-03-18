package confirmsrv

type PaymentProvider interface {
	UpdateStatus() error
}

// Тут ебануть нормальные интерфейсы, подумать над названием пакета, репы, потом...
