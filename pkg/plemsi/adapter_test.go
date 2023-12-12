package plemsi_test

import (
	"fmt"
	"github.com/BacoFoods/menu/internal"
	"github.com/BacoFoods/menu/pkg/invoice"
	"github.com/BacoFoods/menu/pkg/plemsi"
	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/go-faker/faker/v4"
	"github.com/go-resty/resty/v2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Adapter", func() {
	var httpclient shared.RestClient
	var adapter plemsi.Adapter
	var invoice_ invoice.Invoice

	BeforeEach(func() {
		internal.Config.AppConfig.PlemsiHost = "https://pruebas.plemsi.com/api"
		internal.Config.AppConfig.PlemsiToken = "c7f49e2170762271392564b4"
		httpclient = shared.NewRestClient(resty.New())
		adapter = plemsi.NewPlemsi(httpclient)
		_ = faker.FakeData(&invoice_)
	})

	Describe("Endpoint Test Connection", func() {
		Context("With token", func() {
			It("should have a success response", func() {
				err := adapter.TestConnection()
				Expect(err).To(BeNil())
			})
		})

		Context("With out token", func() {
			internal.Config.AppConfig.PlemsiToken = ""
			It("shouldn't connect", func() {
				err := adapter.TestConnection()
				Expect(err.Error()).To(Equal(plemsi.ErrorPlemsiTestConnection))
			})
		})
	})

	Describe("Endpoint End Consumer", func() {
		Context("With nil invoice", func() {
			It("should return an error", func() {
				err := adapter.EmitFinalConsumerInvoice(nil)
				Expect(err.Error()).To(Equal(plemsi.ErrorPlemsiEmptyInvoice))
			})
		})

		Context("With zero invoice", func() {
			It("should return an error", func() {
				fakeInvoice := &invoice.Invoice{}
				if err := faker.FakeData(fakeInvoice); err != nil {
					Expect(err).To(BeNil())
				}

				plemsiInvoice, err := invoice.ToPlemsiInvoice(fakeInvoice, fmt.Sprint(faker.SetRandomNumberBoundaries(1000000000, 9999999999)))
				Expect(err).To(BeNil())

				if err := adapter.EmitFinalConsumerInvoice(plemsiInvoice); err != nil {
					Expect(err.Error()).To(Equal(plemsi.ErrorPlemsiTestConnection))
				}
			})
		})
	})
})
