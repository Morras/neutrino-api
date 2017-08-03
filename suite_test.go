package neutrinoapi

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestNeutrinoAPI(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Neutrino api token validator Suite")
}
