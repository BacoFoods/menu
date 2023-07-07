package menu_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMenu(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Menu Suite")
}
