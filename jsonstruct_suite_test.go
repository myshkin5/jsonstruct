package jsonstruct_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestJSON(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "JSONStruct Suite")
}
