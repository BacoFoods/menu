package plemsi

import (
	"fmt"
	"github.com/BacoFoods/menu/internal"
	"github.com/BacoFoods/menu/pkg/shared"
	"net/http"
)

const LogAdapter = "pkg/plemsi/adapter"

type Adapter interface {
	TestConnection() error
	EmitFinalConsumerInvoice(finalConsumerInvoice *Invoice) error
	EmitConsumerInvoice(consumerInvoice *Invoice) error
}

type adapter struct {
	httpclient shared.RestClient
}

func NewPlemsi(httpclient shared.RestClient) *adapter {
	return &adapter{httpclient}
}

func (a *adapter) TestConnection() error {
	res := make(map[string]any)
	req := shared.Request{
		Endpoint: fmt.Sprintf("%s/ping", internal.Config.PlemsiHost),
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Accept":        "application/json",
			"Authorization": fmt.Sprintf("Bearer %s", internal.Config.PlemsiToken),
		},
		Response: &res,
	}
	resp, err := a.httpclient.Get(req)
	if err != nil {
		shared.LogError("plemsi error making test connection", LogAdapter, "TestConection", err, req)
		return fmt.Errorf(ErrorPlemsiTestConnection)
	}

	if resp.StatusCode() != http.StatusAccepted {
		err := fmt.Errorf(ErrorPlemsiTestConnection)
		shared.LogError("plemsi error making test connection error status", LogAdapter, "TestConnection", err, req, resp)
		return err
	}

	return nil
}

func (a *adapter) EmitFinalConsumerInvoice(finalConsumerInvoice *Invoice) error {
	if finalConsumerInvoice == nil {
		shared.LogWarn("warning invoice nil", LogAdapter, "EmitFinalConsumerInvoice", nil, nil)
		return fmt.Errorf(ErrorPlemsiEmptyInvoice)
	}

	res := make(map[string]any)
	req := shared.Request{
		Endpoint: fmt.Sprintf("%s/billing/invoice", internal.Config.PlemsiHost),
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Accept":        "application/json",
			"Authorization": fmt.Sprintf("Bearer %s", internal.Config.PlemsiToken),
		},
		Response: &res,
		Body:     finalConsumerInvoice,
	}

	resp, err := a.httpclient.Post(req)
	if err != nil {
		shared.LogError("plemsi error, sending final consumer invoice", LogAdapter, "EmitFinalConsumerInvoice", err, req)
		return fmt.Errorf(ErrorPlemsiEndConsumerInvoice)
	}

	if resp.StatusCode() != http.StatusOK {
		shared.LogError("plemsi error, bad status code consumer invoice", LogAdapter, "EmitFinalConsumerInvoice", err, req, resp)
		return fmt.Errorf(ErrorPlemsiEndConsumerInvoice)
	}

	return nil
}

func (a *adapter) EmitConsumerInvoice(consumerInvoice *Invoice) error {
	// TODO: Implements
	return fmt.Errorf("not implemented")
}
