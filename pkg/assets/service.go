package assets

// "github.com/360EntSecGroup-Skylar/excelize"

type AssetRepositoryInterface interface {
	FindByPlaca(placa string) (*Asset, error)
	FindAll(limit int, offset int) ([]Asset, error)
	Create(asset *Asset) error
	Update(asset Asset) error
}

type AssetService struct {
	repo AssetRepositoryInterface
}

// NewAssetService initializes a new asset service
func NewAssetService(repo AssetRepositoryInterface) *AssetService {
	return &AssetService{repo: repo}
}

// FindByPlaca fetches an asset by its Placa value
func (svc *AssetService) FindByPlaca(placa string) (*Asset, error) {
	return svc.repo.FindByPlaca(placa)
}

// FindAll fetches all assets with pagination
func (svc *AssetService) FindAll(limit int, offset int) ([]Asset, error) {
	return svc.repo.FindAll(limit, offset)
}

// CreateAsset creates a new asset
func (svc *AssetService) CreateAsset(asset *Asset) error {
	return svc.repo.Create(asset)
}

// ImportFromExcel reads the Excel file, parses its content, and inserts the results into the database
// func (svc *AssetService) ImportFromExcel(filePath string) error {
// 	f, err := excelize.OpenFile(filePath)
// 	if err != nil {
// 		return err
// 	}

// 	// Read the data from the "ACTIVOS BACU" sheet
// 	rows, err := f.GetRows("ACTIVOS BACU")
// 	if err != nil {
// 		return err
// 	}

// 	// Parse rows to assets and insert them into the database
// 	for _, row := range rows[2:] { // Skipping header rows
// 		asset := Asset{
// 			// Map the columns to the asset fields
// 			// Example: Reference: row[0]
// 		}
// 		err := svc.repo.Create(&asset)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// Update updates an asset
func (svc *AssetService) Update(asset Asset) error {
	return svc.repo.Update(asset)
}
