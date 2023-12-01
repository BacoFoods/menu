package plemsi

type Plemsi interface {
	TestConnection() error
	EmitFinalConsumerInvoice(T) (T, error)
	EmitConsumerInvoice(T) (T, error)
}

type PlemsiInvoice struct {
}
