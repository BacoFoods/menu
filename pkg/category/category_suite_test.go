package category_test

import (
	"testing"

	"github.com/BacoFoods/menu/pkg/menu"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
)

func TestCategory(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Category Suite")
}

func TruncateMenu(db *gorm.DB) {
	if err := db.Model(&menu.Menu{}).Exec("TRUNCATE TABLE menus CASCADE").Error; err != nil {
		Expect(err).To(BeNil())
	}
}
