package server_test

import (
	"context"

	controlv1 "github.com/aity-cloud/monty/pkg/apis/control/v1"
	"github.com/aity-cloud/monty/pkg/config/reactive"
	"github.com/aity-cloud/monty/pkg/config/reactive/reactivetest"
	"github.com/aity-cloud/monty/pkg/config/reactive/subtle"
	configv1 "github.com/aity-cloud/monty/pkg/config/v1"
	_ "github.com/aity-cloud/monty/pkg/oci/noop"
	"github.com/aity-cloud/monty/pkg/test/testlog"
	"github.com/aity-cloud/monty/pkg/update"
	"github.com/aity-cloud/monty/pkg/update/kubernetes"
	"github.com/aity-cloud/monty/pkg/update/kubernetes/server"
	"github.com/aity-cloud/monty/pkg/urn"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protopath"
)

const (
	imageDigest = "sha256:15e2b0d3c33891ebb0f1ef609ec419420c20e320ce94c65fbc8c3312448eb225"
)

var _ = Describe("Kubernetes sync server", Label("unit"), func() {
	var k8sServer update.UpdateTypeHandler
	incomingContext := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
		controlv1.UpdateStrategyKeyForType(urn.Agent), "noop",
	))

	BeforeEach(func() {
		ctrl, ctx, ca := reactivetest.InMemoryController[*configv1.GatewayConfigSpec](
			reactivetest.WithExistingActiveConfig(&configv1.GatewayConfigSpec{
				Upgrades: &configv1.UpgradesSpec{
					Agents: &configv1.AgentUpgradesSpec{
						Kubernetes: &configv1.KubernetesAgentUpgradeSpec{
							ImageResolver: configv1.KubernetesAgentUpgradeSpec_Noop.Enum(),
						},
					},
				},
			}),
		)
		DeferCleanup(ca)

		var err error
		k8sServer, err = server.NewKubernetesSyncServer(ctx,
			reactive.Message[*configv1.KubernetesAgentUpgradeSpec](
				ctrl.Reactive(protopath.Path(configv1.ProtoPath().Upgrades().Agents().Kubernetes()))),
			testlog.Log,
		)
		Expect(err).NotTo(HaveOccurred())

		subtle.WaitOne(ctx, ctrl.Reactive(protopath.Path(configv1.ProtoPath().Upgrades().Agents().Kubernetes())))
	})

	When("unknown package type is provided", func() {
		packageURN := urn.NewMontyURN(urn.Agent, kubernetes.UpdateStrategy, "unknown")
		manifest := &controlv1.UpdateManifest{
			Items: []*controlv1.UpdateManifestEntry{
				{
					Package: packageURN.String(),
					Path:    "registry.aity.tech/monty/monty",
					Digest:  "latest",
				},
			},
		}
		It("should return an error", func() {
			_, err := k8sServer.CalculateUpdate(incomingContext, manifest)
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
		})
	})
	When("an invalid URN is provided", func() {
		manifest := &controlv1.UpdateManifest{
			Items: []*controlv1.UpdateManifestEntry{
				{
					Package: "urn:malformed",
					Path:    "registry.aity.tech/monty/monty",
					Digest:  "latest",
				},
			},
		}
		It("should return an error", func() {
			_, err := k8sServer.CalculateUpdate(incomingContext, manifest)
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
		})
	})
	When("URNs are valid", func() {
		var packageURN1, packageURN2 urn.MontyURN
		var manifest *controlv1.UpdateManifest
		BeforeEach(func() {
			packageURN1 = urn.NewMontyURN(urn.Agent, kubernetes.UpdateStrategy, "agent")
			packageURN2 = urn.NewMontyURN(urn.Agent, kubernetes.UpdateStrategy, "controller")
		})
		When("manifest matches the current version", func() {
			JustBeforeEach(func() {
				manifest = &controlv1.UpdateManifest{
					Items: []*controlv1.UpdateManifestEntry{
						{
							Package: packageURN1.String(),
							Path:    "example.io/monty-noop",
							Digest:  imageDigest,
						},
						{
							Package: packageURN2.String(),
							Path:    "example.io/monty-noop",
							Digest:  imageDigest,
						},
					},
				}
			})
			It("should return noop updates", func() {
				patchList, err := k8sServer.CalculateUpdate(incomingContext, manifest)
				Expect(err).NotTo(HaveOccurred())
				Expect(patchList.GetItems()).To(HaveLen(2))
				Expect(patchList.Items[0].GetOp()).To(Equal(controlv1.PatchOp_None))
				Expect(patchList.Items[1].GetOp()).To(Equal(controlv1.PatchOp_None))
				Expect(patchList.Items[0].GetPackage()).To(Equal(packageURN1.String()))
				Expect(patchList.Items[1].GetPackage()).To(Equal(packageURN2.String()))
			})
		})
		When("digest matches the current version but registry is different", func() {
			JustBeforeEach(func() {
				manifest = &controlv1.UpdateManifest{
					Items: []*controlv1.UpdateManifestEntry{
						{
							Package: packageURN1.String(),
							Path:    "quay.io/monty-noop",
							Digest:  imageDigest,
						},
						{
							Package: packageURN2.String(),
							Path:    "monty-noop",
							Digest:  imageDigest,
						},
					},
				}
			})
			It("should return noop updates", func() {
				patchList, err := k8sServer.CalculateUpdate(incomingContext, manifest)
				Expect(err).NotTo(HaveOccurred())
				Expect(patchList.GetItems()).To(HaveLen(2))
				Expect(patchList.Items[0].GetOp()).To(Equal(controlv1.PatchOp_None))
				Expect(patchList.Items[1].GetOp()).To(Equal(controlv1.PatchOp_None))
				Expect(patchList.Items[0].GetPackage()).To(Equal(packageURN1.String()))
				Expect(patchList.Items[1].GetPackage()).To(Equal(packageURN2.String()))
			})
		})
		When("one digest does not match", func() {
			JustBeforeEach(func() {
				manifest = &controlv1.UpdateManifest{
					Items: []*controlv1.UpdateManifestEntry{
						{
							Package: packageURN1.String(),
							Path:    "example.io/monty-noop",
							Digest:  imageDigest,
						},
						{
							Package: packageURN2.String(),
							Path:    "example.io/monty-noop",
							Digest:  "latest",
						},
					},
				}
			})
			It("should return one change update", func() {
				patchList, err := k8sServer.CalculateUpdate(incomingContext, manifest)
				Expect(err).NotTo(HaveOccurred())
				for _, patch := range patchList.GetItems() {
					if patch.GetPackage() == packageURN2.String() {
						Expect(patch.GetOp()).To(Equal(controlv1.PatchOp_Update))
						Expect(patch.GetPath()).To(Equal("example.io/monty-noop"))
						Expect(patch.GetNewDigest()).To(Equal(imageDigest))
						Expect(patch.GetOldDigest()).To(Equal("latest"))
					}
				}
			})
			It("should return one noop update", func() {
				patchList, err := k8sServer.CalculateUpdate(incomingContext, manifest)
				Expect(err).NotTo(HaveOccurred())
				Expect(func() bool {
					for _, patch := range patchList.GetItems() {
						if patch.GetOp() == controlv1.PatchOp_None {
							return true
						}
					}
					return false
				}()).To(BeTrue())
			})
		})
		When("one digest does not match and the registry is different", func() {
			JustBeforeEach(func() {
				manifest = &controlv1.UpdateManifest{
					Items: []*controlv1.UpdateManifestEntry{
						{
							Package: packageURN1.String(),
							Path:    "example.io/monty-noop",
							Digest:  imageDigest,
						},
						{
							Package: packageURN2.String(),
							Path:    "monty.io/monty-noop",
							Digest:  "latest",
						},
					},
				}
			})
			It("should return one change update with the registry changed", func() {
				patchList, err := k8sServer.CalculateUpdate(incomingContext, manifest)
				Expect(err).NotTo(HaveOccurred())
				for _, patch := range patchList.GetItems() {
					if patch.GetPackage() == packageURN2.String() {
						Expect(patch.GetOp()).To(Equal(controlv1.PatchOp_Update))
						Expect(patch.GetPath()).To(Equal("example.io/monty-noop"))
						Expect(patch.GetNewDigest()).To(Equal(imageDigest))
						Expect(patch.GetOldDigest()).To(Equal("latest"))
					}
				}
			})
		})
		When("image repository is different", func() {
			JustBeforeEach(func() {
				manifest = &controlv1.UpdateManifest{
					Items: []*controlv1.UpdateManifestEntry{
						{
							Package: packageURN1.String(),
							Path:    "example.io/monty-noop",
							Digest:  imageDigest,
						},
						{
							Package: packageURN2.String(),
							Path:    "example.io/rancher/test",
							Digest:  imageDigest,
						},
					},
				}
			})
			It("should return one change update with the correct repo", func() {
				patchList, err := k8sServer.CalculateUpdate(incomingContext, manifest)
				Expect(err).NotTo(HaveOccurred())
				for _, patch := range patchList.GetItems() {
					if patch.GetPackage() == packageURN2.String() {
						Expect(patch.GetOp()).To(Equal(controlv1.PatchOp_Update))
						Expect(patch.GetPath()).To(Equal("example.io/monty-noop"))
						Expect(patch.GetNewDigest()).To(Equal(imageDigest))
						Expect(patch.GetOldDigest()).To(Equal(imageDigest))
					}
				}
			})
		})
		When("both images should be updated", func() {
			JustBeforeEach(func() {
				manifest = &controlv1.UpdateManifest{
					Items: []*controlv1.UpdateManifestEntry{
						{
							Package: packageURN1.String(),
							Path:    "example.io/monty-noop",
							Digest:  "latest",
						},
						{
							Package: packageURN2.String(),
							Path:    "example.io/rancher/test",
							Digest:  imageDigest,
						},
					},
				}
			})
			It("should return patch updates", func() {
				patchList, err := k8sServer.CalculateUpdate(incomingContext, manifest)
				Expect(err).NotTo(HaveOccurred())
				Expect(func() bool {
					for _, patch := range patchList.GetItems() {
						if patch.GetOp() != controlv1.PatchOp_Update {
							return false
						}
					}
					return true
				}()).To(BeTrue())
			})
		})
	})
})
