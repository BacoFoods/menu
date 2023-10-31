package assets_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/BacoFoods/menu/pkg/assets"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AssetHandler", func() {
	var (
		router       *gin.Engine
		mockRepo     *MockAssetRepository
		assetHandler *assets.Handler
	)

	BeforeEach(func() {
		router = gin.Default()
		mockRepo = &MockAssetRepository{[]assets.Asset{}}
		assetService := assets.NewAssetService(mockRepo)
		assetHandler = assets.NewHandler(assetService)

		// Setting up the routes
		assetsRoutes := assets.NewRoutes(assetHandler)
		assetsRoutes.RegisterRoutes(router.Group(""))
	})

	Describe("CreateAsset", func() {
		It("should create an asset successfully", func() {
			asset := assets.Asset{
				Reference: "TestReference",
				// ... other fields
			}
			body, _ := json.Marshal(asset)
			req, _ := http.NewRequest("POST", "/assets", bytes.NewBuffer(body))
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)
			Expect(resp.Code).To(Equal(http.StatusCreated))
			Expect(mockRepo.Assets).To(HaveLen(1))
		})
	})

	Describe("GetByPlaca", func() {
		It("should fetch an asset by its Placa value", func() {
			placa := "TESTPLACA"
			asset := assets.Asset{
				ID:        1,
				Reference: "TestReference",
				Placa:     placa,
				// ... other fields
			}

			mockRepo.Create(&asset)
			req, _ := http.NewRequest("GET", "/assets/TESTPLACA", nil)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)
			Expect(resp.Code).To(Equal(http.StatusOK), resp.Body.String())

			var ret map[string]any
			_ = json.Unmarshal(resp.Body.Bytes(), &ret)
			returnedAsset := ret["data"].(map[string]any)
			Expect(returnedAsset["placa"]).To(Equal("TESTPLACA"))
			Expect(returnedAsset["reference"]).To(Equal("TestReference"))
		})
	})
})
