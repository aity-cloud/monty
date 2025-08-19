package bootstrap_test

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net"
	"runtime"
	"time"

	bootstrapv2 "github.com/aity-cloud/monty/pkg/apis/bootstrap/v2"
	"github.com/aity-cloud/monty/pkg/bootstrap"
	"github.com/aity-cloud/monty/pkg/config/reactive"
	"github.com/aity-cloud/monty/pkg/config/reactive/reactivetest"
	configv1 "github.com/aity-cloud/monty/pkg/config/v1"
	"github.com/aity-cloud/monty/pkg/ident"
	"github.com/aity-cloud/monty/pkg/pkp"
	"github.com/aity-cloud/monty/pkg/storage"
	"github.com/aity-cloud/monty/pkg/storage/inmemory"
	mock_ident "github.com/aity-cloud/monty/pkg/test/mock/ident"
	mock_storage "github.com/aity-cloud/monty/pkg/test/mock/storage"
	"github.com/aity-cloud/monty/pkg/test/testdata"
	"github.com/aity-cloud/monty/pkg/test/testlog"
	"github.com/aity-cloud/monty/pkg/test/testutil"
	"github.com/aity-cloud/monty/pkg/tokens"
	"github.com/aity-cloud/monty/pkg/trust"
	"github.com/aity-cloud/monty/pkg/util"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/samber/lo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/reflect/protopath"
)

func pkpTrustStrategy(cert *x509.Certificate) trust.Strategy {
	conf := trust.StrategyConfig{
		PKP: &trust.PKPConfig{
			Pins: trust.NewPinSource([]*pkp.PublicKeyPin{pkp.NewSha256(cert)}),
		},
	}
	return util.Must(conf.Build())
}

var _ = Describe("Client V2", Ordered, Label("unit"), func() {
	var fooIdent ident.Provider
	var cert *tls.Certificate
	var store struct {
		storage.Backend
		storage.LockManagerBroker
	}

	var endpoint string

	BeforeAll(func() {
		if runtime.GOOS != "linux" {
			Skip("skipping test on non-linux OS")
		}
		fooIdent = mock_ident.NewTestIdentProvider(ctrl, "foo")
		var err error
		crt, err := tls.X509KeyPair(testdata.TestData("localhost.crt"), testdata.TestData("localhost.key"))
		Expect(err).NotTo(HaveOccurred())
		crt.Leaf, err = x509.ParseCertificate(crt.Certificate[0])
		Expect(err).NotTo(HaveOccurred())
		cert = &crt
		store.Backend = mock_storage.NewTestStorageBackend(context.Background(), ctrl)
		store.LockManagerBroker = inmemory.NewLockManagerBroker()

		srv := grpc.NewServer(grpc.Creds(credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{*cert},
		})))

		ctrl, ctx, ca := reactivetest.InMemoryController[*configv1.GatewayConfigSpec](reactivetest.WithExistingActiveConfig(
			&configv1.GatewayConfigSpec{
				Certs: &configv1.CertsSpec{
					CaCertData:      lo.ToPtr(string(testdata.TestData("root_ca.crt"))),
					ServingCertData: lo.ToPtr(string(testdata.TestData("localhost.crt"))),
					ServingKeyData:  lo.ToPtr(string(testdata.TestData("localhost.key"))),
				},
			},
		))

		server := bootstrap.NewServerV2(ctx, store, reactive.Message[*configv1.CertsSpec](ctrl.Reactive(protopath.Path(configv1.ProtoPath().Certs()))), testlog.Log)
		bootstrapv2.RegisterBootstrapServer(srv, server)

		listener, err := net.Listen("tcp4", "127.0.0.1:0")
		Expect(err).NotTo(HaveOccurred())
		endpoint = listener.Addr().String()

		go srv.Serve(listener)

		DeferCleanup(func() {
			srv.Stop()
			ca()
		})
	})

	It("should bootstrap with the server", func() {
		token, _ := store.CreateToken(context.Background(), 1*time.Minute)
		cc := bootstrap.ClientConfigV2{
			Token:         testutil.Must(tokens.FromBootstrapToken(token)),
			Endpoint:      endpoint,
			TrustStrategy: pkpTrustStrategy(cert.Leaf),
		}

		_, err := cc.Bootstrap(context.Background(), fooIdent)
		Expect(err).NotTo(HaveOccurred())
	})
	Context("error handling", func() {
		When("no token is given", func() {
			It("should error", func() {
				cc := bootstrap.ClientConfigV2{}
				kr, err := cc.Bootstrap(context.Background(), fooIdent)
				Expect(kr).To(BeNil())
				Expect(err).To(MatchError(bootstrap.ErrNoToken))
			})
		})
		When("an invalid endpoint is given", func() {
			It("should error", func() {
				token, _ := store.CreateToken(context.Background(), 1*time.Minute)
				cc := bootstrap.ClientConfigV2{
					Token:         testutil.Must(tokens.FromBootstrapToken(token)),
					Endpoint:      "\x7f",
					TrustStrategy: pkpTrustStrategy(cert.Leaf),
				}
				kr, err := cc.Bootstrap(context.Background(), fooIdent)
				Expect(kr).To(BeNil())
				Expect(err.Error()).To(ContainSubstring("net/url"))
			})
		})
		When("the client fails to send a request to the server", func() {
			It("should error", func() {
				token, _ := store.CreateToken(context.Background(), 1*time.Minute)
				cc := bootstrap.ClientConfigV2{
					Token:         testutil.Must(tokens.FromBootstrapToken(token)),
					Endpoint:      "localhost:65545",
					TrustStrategy: pkpTrustStrategy(cert.Leaf),
				}
				kr, err := cc.Bootstrap(context.Background(), fooIdent)
				Expect(kr).To(BeNil())
				Expect(err.Error()).To(ContainSubstring("invalid port"))
			})
		})
	})
})
