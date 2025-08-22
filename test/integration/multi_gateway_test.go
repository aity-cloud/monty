package integration_test

import (
	"errors"
	"time"

	capabilityv1 "github.com/aity-cloud/monty/pkg/apis/capability/v1"
	corev1 "github.com/aity-cloud/monty/pkg/apis/core/v1"
	managementv1 "github.com/aity-cloud/monty/pkg/apis/management/v1"
	"github.com/aity-cloud/monty/pkg/config/v1beta1"
	"github.com/aity-cloud/monty/pkg/test"
	"github.com/aity-cloud/monty/pkg/test/testlog"
	"github.com/aity-cloud/monty/pkg/util"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var _ = Describe("Multi Gateway etcd", Ordered, Label("integration"), func() {
	var mgmtClient managementv1.ManagementClient
	BeforeAll(func() {
		By("Starting 3 test environments")
		env1 := &test.Environment{
			Logger: testlog.Log.WithGroup("env1"),
		}
		Expect(env1.Start(
			test.WithStorageBackend(v1beta1.StorageTypeEtcd),
			test.WithInMemoryActiveStore(true),
		)).To(Succeed())
		env2 := &test.Environment{
			Logger: testlog.Log.WithGroup("env2"),
		}
		Expect(env2.Start(
			test.WithStorageBackend(v1beta1.StorageTypeEtcd),
			test.WithRemoteEtcdPort(env1.GetPorts().Etcd),
			test.WithInMemoryActiveStore(true),
		)).To(Succeed())
		env3 := &test.Environment{
			Logger: testlog.Log.WithGroup("env3"),
		}
		Expect(env3.Start(
			test.WithStorageBackend(v1beta1.StorageTypeEtcd),
			test.WithRemoteEtcdPort(env1.GetPorts().Etcd),
			test.WithInMemoryActiveStore(true),
		)).To(Succeed())

		resolver := test.NewEnvironmentResolver(env3, env2, env1)
		cc, err := grpc.Dial("testenv:///management",
			grpc.WithResolvers(resolver),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithDefaultCallOptions(grpc.WaitForReady(true)),
		)
		Expect(err).NotTo(HaveOccurred())
		mgmtClient = managementv1.NewManagementClient(cc)

		By("adding one agent to each environment")
		err = env1.BootstrapNewAgent("agent1", test.WithLocalAgent())
		Expect(err).NotTo(HaveOccurred())

		err = env2.BootstrapNewAgent("agent2")
		Expect(err).NotTo(HaveOccurred())

		err = env3.BootstrapNewAgent("agent3")
		Expect(err).NotTo(HaveOccurred())
		stopAll := func(cause string) error {
			var eg util.MultiErrGroup
			eg.Go(func() error {
				return env1.Stop(cause)
			})
			eg.Go(func() error {
				return env2.Stop(cause)
			})
			eg.Go(func() error {
				return env3.Stop(cause)
			})
			eg.Wait()
			return eg.Error()
		}
		DeferCleanup(stopAll, "Test Suite Finished")
	})

	It("should install capabilities onto all agents", func(ctx SpecContext) {
		By("installing the example capability onto all agents")
		time.Sleep(5 * time.Second)
		resp, err := mgmtClient.InstallCapability(ctx, &capabilityv1.InstallRequest{
			Capability: &corev1.Reference{Id: "example"},
			Agent:      &corev1.Reference{Id: "agent1"},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.Status).To(Equal(capabilityv1.InstallResponseStatus_Success))
		_, err = mgmtClient.InstallCapability(ctx, &capabilityv1.InstallRequest{
			Capability: &corev1.Reference{Id: "example"},
			Agent:      &corev1.Reference{Id: "agent2"},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.Status).To(Equal(capabilityv1.InstallResponseStatus_Success))
		_, err = mgmtClient.InstallCapability(ctx, &capabilityv1.InstallRequest{
			Capability: &corev1.Reference{Id: "example"},
			Agent:      &corev1.Reference{Id: "agent3"},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.Status).To(Equal(capabilityv1.InstallResponseStatus_Success))

		Eventually(func() error {
			for _, agent := range []string{"agent1", "agent2", "agent3"} {
				stat, err := mgmtClient.CapabilityStatus(ctx, &capabilityv1.StatusRequest{
					Capability: &corev1.Reference{Id: "example"},
					Agent:      &corev1.Reference{Id: agent},
				})
				if err != nil {
					return err
				}
				if !stat.Enabled {
					return errors.New("capability not enabled")
				}
			}
			return nil
		}).Should(Succeed())
	})
})

var _ = Describe("Multi Gateway js", Ordered, Label("integration"), func() {
	var mgmtClient managementv1.ManagementClient
	BeforeAll(func() {
		By("Starting 3 test environments")
		env1 := &test.Environment{
			Logger: testlog.Log.WithGroup("env1"),
		}
		Expect(env1.Start(
			test.WithStorageBackend(v1beta1.StorageTypeJetStream),
			test.WithInMemoryActiveStore(true),
		)).To(Succeed())
		seedPath := env1.JetStreamConfig().GetNkeySeedPath()

		time.Sleep(1 * time.Second)

		env2 := &test.Environment{
			Logger: testlog.Log.WithGroup("env2"),
		}
		Expect(env2.Start(
			test.WithStorageBackend(v1beta1.StorageTypeJetStream),
			test.WithRemoteJetStreamPort(env1.GetPorts().Jetstream),
			test.WithRemoteJetStreamSeedPath(seedPath),
			test.WithInMemoryActiveStore(true),
		)).To(Succeed())
		env3 := &test.Environment{
			Logger: testlog.Log.WithGroup("env3"),
		}
		Expect(env3.Start(
			test.WithStorageBackend(v1beta1.StorageTypeJetStream),
			test.WithRemoteJetStreamPort(env1.GetPorts().Jetstream),
			test.WithRemoteJetStreamSeedPath(seedPath),
			test.WithInMemoryActiveStore(true),
		)).To(Succeed())

		resolver := test.NewEnvironmentResolver(env3, env2, env1)
		cc, err := grpc.Dial("testenv:///management", grpc.WithResolvers(resolver), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.WaitForReady(true)))
		Expect(err).NotTo(HaveOccurred())
		mgmtClient = managementv1.NewManagementClient(cc)

		By("adding one agent to each environment")
		err = env1.BootstrapNewAgent("agent1", test.WithLocalAgent())
		Expect(err).NotTo(HaveOccurred())

		err = env2.BootstrapNewAgent("agent2")
		Expect(err).NotTo(HaveOccurred())

		err = env3.BootstrapNewAgent("agent3")
		Expect(err).NotTo(HaveOccurred())

		stopAll := func(cause string) error {
			var eg util.MultiErrGroup
			eg.Go(func() error {
				return env1.Stop(cause)
			})
			eg.Go(func() error {
				return env2.Stop(cause)
			})
			eg.Go(func() error {
				return env3.Stop(cause)
			})
			eg.Wait()
			return eg.Error()
		}
		DeferCleanup(stopAll, "Test Suite Finished")
	})

	It("should install capabilities onto all agents", func(ctx SpecContext) {
		By("installing the example capability onto all agents")
		time.Sleep(5 * time.Second)
		resp, err := mgmtClient.InstallCapability(ctx, &capabilityv1.InstallRequest{
			Capability: &corev1.Reference{Id: "example"},
			Agent:      &corev1.Reference{Id: "agent1"},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.Status).To(Equal(capabilityv1.InstallResponseStatus_Success))
		_, err = mgmtClient.InstallCapability(ctx, &capabilityv1.InstallRequest{
			Capability: &corev1.Reference{Id: "example"},
			Agent:      &corev1.Reference{Id: "agent2"},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.Status).To(Equal(capabilityv1.InstallResponseStatus_Success))
		_, err = mgmtClient.InstallCapability(ctx, &capabilityv1.InstallRequest{
			Capability: &corev1.Reference{Id: "example"},
			Agent:      &corev1.Reference{Id: "agent3"},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.Status).To(Equal(capabilityv1.InstallResponseStatus_Success))

		Eventually(func() error {
			for _, agent := range []string{"agent1", "agent2", "agent3"} {
				stat, err := mgmtClient.CapabilityStatus(ctx, &capabilityv1.StatusRequest{
					Capability: &corev1.Reference{Id: "example"},
					Agent:      &corev1.Reference{Id: agent},
				})
				if err != nil {
					return err
				}
				if !stat.Enabled {
					return errors.New("capability not enabled")
				}
			}
			return nil
		}).Should(Succeed())
	})
})
