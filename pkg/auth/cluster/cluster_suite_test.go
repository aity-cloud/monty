package cluster_test

import (
	"testing"

	"github.com/aity-cloud/monty/pkg/ecdh"
	"github.com/aity-cloud/monty/pkg/keyring"
	_ "github.com/aity-cloud/monty/pkg/test/setup"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var (
	testSharedKeys   *keyring.SharedKeys
	testServerKey    []byte
	testClientKey    []byte
	invalidKey       []byte
	testSharedSecret []byte

	ctrl *gomock.Controller
)

func TestClusterAuth(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cluster Suite")
}

var _ = BeforeSuite(func() {
	ctrl = gomock.NewController(GinkgoT())
	kp1 := ecdh.NewEphemeralKeyPair()
	kp2 := ecdh.NewEphemeralKeyPair()
	sec, err := ecdh.DeriveSharedSecret(kp1, ecdh.PeerPublicKey{
		PublicKey: kp2.PublicKey,
		PeerType:  ecdh.PeerTypeClient,
	})
	if err != nil {
		panic(err)
	}
	testSharedKeys = keyring.NewSharedKeys(sec)
	testServerKey = testSharedKeys.ServerKey
	testClientKey = testSharedKeys.ClientKey
	testSharedSecret = sec
})
