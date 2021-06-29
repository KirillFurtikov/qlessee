package qless_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestQlessGo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "QlessGo Suite")
}
