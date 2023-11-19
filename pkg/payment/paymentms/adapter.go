package paymentms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/BacoFoods/menu/pkg/payment"
)

type PaymentsAPI struct {
	Httpclient HTTPClient

	paymentsHost string
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewPaymentsAPI(Httpclient HTTPClient, paymentsHost string) *PaymentsAPI {
	return &PaymentsAPI{Httpclient, paymentsHost}
}

func (p *PaymentsAPI) PaylotStatus(paylotID string) (*payment.PaylotStatus, error) {
	url := fmt.Sprintf("%s/api/payments/v1/paylot/%s/status", p.paymentsHost, paylotID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	var response payment.PaylotStatusResponse
	resp, err := p.Httpclient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, fmt.Errorf("error getting order info: nil")
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&response)

	if resp != nil && resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("error getting order info: %s", resp.Status)
	}

	if err != nil {
		return nil, err
	}

	return &response.Data, nil
}

func (p *PaymentsAPI) CreatePaylot(r payment.PaylotReq) (*payment.Paylot, error) {
	url := fmt.Sprintf("%s/api/payments/v1/paylot", p.paymentsHost)
	body, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		url,
		bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header["Content-Type"] = []string{"application/json"}

	var response payment.PaylotResponse
	resp, err := p.Httpclient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, fmt.Errorf("error getting order info: nil")
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&response)

	if resp != nil && (resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK) {
		return nil, fmt.Errorf("error getting order info: %s", resp.Status)
	}

	if err != nil {
		return nil, err
	}

	return &response.Data, nil
}
