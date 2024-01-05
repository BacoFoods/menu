package plemsi_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPlemsi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Plemsi Suite")
}
