package assets

import (
	"net/http"

	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/gin-gonic/gin"
)

const (
	ErrorAssetNotFound = "asset not found"
)

// Handler is the struct that handles the requests
// An example curl to create an asset locally:
// curl -X POST http://localhost:8080/api/menu/v1/assets -H 'Authorization: <token>' -H 'Content-Type: application/json' -d '{"reference": "030076", "placa": "003", "asset_name": "MICROSCOPIO", "asset_name_siesa": "MICROSCOPIO", "operation_code_siesa": "MUEBLES Y ENSERES", "operation_name_siesa": "BIENES DIVERSOS", "purchase_date": "2023-01-01T00:00:00Z", "purchase_invoice": "FACTURA", "provider_nit": "1122334455", "provider_name": "PROVEEDOR", "official_placa": "PLACA", "current_location": "LABORATORIO", "current_cost_center": 100, "custodian": "RESPONSABLE", "invoice": "FACTURA", "invoice_link": "http://example.com/invoice", "actual_purchase_date": "2023-01-01T00:00:00Z", "contract": "CONTRATO", "category": "CATEGORIA", "renting": false, "price": 1000.5, "barcode": "BARCODE", "type": "SIESA", "invoice_initial_location": "LABORATORIO"}'
//
// An example curl to get an asset locally:
// curl http://localhost:8080/api/menu/v1/assets/003 -H 'Content-Type: application/json' -H 'Authorization: <token>'
type Handler struct {
	svc *AssetService
}

func NewHandler(svc *AssetService) *Handler {
	return &Handler{svc}
}

// GetByPlaca to find assets by code
// @Summary Find assets by code
// @Description Find assets by code
// @Tags Assets
// @Param code path string true "placa of the asset"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string}
// @Failure 401 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /assets/{code} [get]
func (h *Handler) GetByPlaca(c *gin.Context) {
	placa := c.Param("code")
	asset, err := h.svc.FindByPlaca(placa)
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared.ErrorResponse(err.Error()))
		return
	}

	if asset == nil {
		c.JSON(http.StatusNotFound, shared.ErrorResponse(ErrorAssetNotFound))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(asset))
}

// CreateAsset to find assets by code
// @Summary Create a new asset
// @Description Create a new asset
// @Tags Assets
// @Accept  json
// @Produce  json
// @Param asset body Asset true "Asset to be created"
// @Success 201 {object} Asset
// @Router /assets [post]
func (h *Handler) CreateAsset(c *gin.Context) {
	var asset Asset
	err := c.BindJSON(&asset)
	if err != nil {
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(err.Error()))
		return
	}

	err = h.svc.CreateAsset(&asset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, shared.SuccessResponse(asset))
}
