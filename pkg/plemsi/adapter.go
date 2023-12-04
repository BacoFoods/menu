package plemsi

import "net/http"

type Adapter interface {
	TestConnection() error
	EmitFinalConsumerInvoice(finalConsumerInvoice *Invoice) error
	EmitConsumerInvoice(consumerInvoice *Invoice) error
}

type adapter struct {
	httpclient *http.Client
}

func NewPlemsi(httpclient *http.Client) *adapter {
	return &adapter{
		httpclient: httpclient,
	}
}

func (a *adapter) TestConnection() error {
	// TODO: implement
	return nil
}
