package trust_test

import (
	"crypto/x509"
	"testing"

	_ "github.com/aity-cloud/monty/pkg/test/setup"
	"github.com/aity-cloud/monty/pkg/test/testdata"
	"github.com/aity-cloud/monty/pkg/util"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTrust(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Trust Suite")
}

func newTestCert() *x509.Certificate {
	certData := testdata.TestData("root_ca.crt")
	cert, err := util.ParsePEMEncodedCert(certData)
	Expect(err).NotTo(HaveOccurred())
	return cert
}
