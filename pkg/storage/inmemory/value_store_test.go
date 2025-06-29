package inmemory_test

import (
	"context"
	"fmt"
	"strings"

	"github.com/aity-cloud/monty/pkg/storage"
	"github.com/aity-cloud/monty/pkg/storage/inmemory"
	"github.com/aity-cloud/monty/pkg/test/testutil"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
)

var _ = Describe("Value Store", Ordered, Label("unit"), func() {
	var (
		ctx             context.Context
		ca              context.CancelFunc
		valueStore      storage.ValueStoreT[string]
		updateC         <-chan storage.WatchEvent[storage.KeyRevision[string]]
		expectedUpdates []lo.Tuple2[string, string]
	)

	BeforeEach(func() {
		ctx, ca = context.WithCancel(context.Background())

		expectedUpdates = nil
		valueStore = inmemory.NewValueStore(strings.Clone)

		var err error
		updateC, err = valueStore.Watch(ctx)
		Expect(err).NotTo(HaveOccurred())
	})
	AfterEach(func() {
		ca()
		receivedUpdates := lo.ChannelToSlice(updateC)
		if len(receivedUpdates) != len(expectedUpdates) {
			Fail(fmt.Sprintf("received %d updates, expected %d: \nwant: %v\ngot: %v", len(receivedUpdates), len(expectedUpdates), expectedUpdates, receivedUpdates))
		}
		for i, pair := range expectedUpdates {
			prev, cur := pair.Unpack()
			update := receivedUpdates[i]
			if update.Previous == nil {
				Expect(prev).To(BeEmpty())
			} else {
				Expect(update.Previous.Value()).To(Equal(prev))
			}
			if update.Current == nil {
				Expect(cur).To(BeEmpty())
			} else {
				Expect(update.Current.Value()).To(Equal(cur))
			}
		}
	})

	expectUpdates := func(pairs ...lo.Tuple2[string, string]) {
		expectedUpdates = append(expectedUpdates, pairs...)
	}

	Describe("Put", func() {
		When("a value is put without a specific revision", func() {
			It("should put the value successfully", func() {
				Expect(valueStore.Put(ctx, "value")).To(Succeed())
				value, err := valueStore.Get(ctx)
				Expect(err).NotTo(HaveOccurred())
				Expect(value).To(Equal("value"))

				expectUpdates(lo.T2("", "value"))
			})
		})

		When("a value is put with a specific revision", func() {
			It("should fail if the revision does not match", func() {
				rev := int64(10)
				Expect(valueStore.Put(ctx, "value", storage.WithRevision(rev))).To(MatchError(storage.ErrConflict))
			})

			It("should succeed if the revision matches", func() {
				revOut := int64(0)
				Expect(valueStore.Put(ctx, "value")).To(Succeed())
				Expect(valueStore.Put(ctx, "new-value", storage.WithRevisionOut(&revOut))).To(Succeed())
				Expect(valueStore.Put(ctx, "new-value2", storage.WithRevision(revOut))).To(Succeed())
				value, err := valueStore.Get(ctx)
				Expect(err).NotTo(HaveOccurred())
				Expect(value).To(Equal("new-value2"))

				expectUpdates(
					lo.T2("", "value"),
					lo.T2("value", "new-value"),
					lo.T2("new-value", "new-value2"),
				)
			})
		})
	})

	Describe("Get", func() {
		When("retrieving a value that does not exist", func() {
			It("should return an error", func() {
				_, err := valueStore.Get(ctx)
				Expect(err).To(Equal(storage.ErrNotFound))
			})
		})

		When("retrieving a value with a specific revision", func() {
			When("the value store is empty", func() {
				It("should return an OutOfRange error for any nonzero revision", func() {
					_, err := valueStore.Get(ctx, storage.WithRevision(1))
					Expect(err).To(testutil.MatchStatusCode(codes.OutOfRange))
				})
			})
			It("should retrieve the value for the correct revision", func() {
				revOut := int64(0)
				Expect(valueStore.Put(ctx, "value1")).To(Succeed())
				Expect(valueStore.Put(ctx, "value2", storage.WithRevisionOut(&revOut))).To(Succeed())
				value, err := valueStore.Get(ctx, storage.WithRevision(revOut))
				Expect(err).NotTo(HaveOccurred())
				Expect(value).To(Equal("value2"))

				expectUpdates(
					lo.T2("", "value1"),
					lo.T2("value1", "value2"),
				)
			})
			When("the value at the specified revision has been deleted", func() {
				It("should return a NotFound error", func() {
					Expect(valueStore.Put(ctx, "value1")).To(Succeed())
					Expect(valueStore.Put(ctx, "value2")).To(Succeed())
					Expect(valueStore.Delete(ctx)).To(Succeed())
					Expect(valueStore.Put(ctx, "value3")).To(Succeed())
					_, err := valueStore.Get(ctx, storage.WithRevision(3))
					Expect(err).To(Equal(storage.ErrNotFound))

					expectUpdates(
						lo.T2("", "value1"),
						lo.T2("value1", "value2"),
						lo.T2("value2", ""),
						lo.T2("", "value3"),
					)
				})
			})
		})
	})

	Describe("Delete", func() {
		When("deleting a value that does not exist", func() {
			It("should return an error", func() {
				Expect(valueStore.Delete(ctx)).To(Equal(storage.ErrNotFound))
			})
		})

		When("deleting a value with a specific revision", func() {
			It("should delete the value if the revision matches", func() {
				Expect(valueStore.Put(ctx, "value")).To(Succeed())
				Expect(valueStore.Delete(ctx, storage.WithRevision(1))).To(Succeed())
				_, err := valueStore.Get(ctx)
				Expect(err).To(Equal(storage.ErrNotFound))

				expectUpdates(
					lo.T2("", "value"),
					lo.T2("value", ""),
				)
			})

			It("should not delete the value if the revision does not match", func() {
				rev := int64(10)
				Expect(valueStore.Put(ctx, "value")).To(Succeed())
				Expect(valueStore.Delete(ctx, storage.WithRevision(rev))).To(HaveOccurred())
				value, err := valueStore.Get(ctx)
				Expect(err).NotTo(HaveOccurred())
				Expect(value).To(Equal("value"))

				expectUpdates(lo.T2("", "value"))
			})
		})

		When("deleting a value that has been marked as deleted", func() {
			It("should return an error", func() {
				Expect(valueStore.Put(ctx, "value")).To(Succeed())
				Expect(valueStore.Delete(ctx)).To(Succeed())
				Expect(valueStore.Delete(ctx)).To(Equal(storage.ErrNotFound))

				expectUpdates(
					lo.T2("", "value"),
					lo.T2("value", ""),
				)
			})
		})

		When("deleting a value that doesn't exist", func() {
			It("should return an error", func() {
				Expect(valueStore.Delete(ctx)).To(Equal(storage.ErrNotFound))
			})
		})
	})

	Describe("History", func() {
		When("retrieving history for a non-existent key", func() {
			It("should return an error", func() {
				_, err := valueStore.History(ctx)
				Expect(err).To(Equal(storage.ErrNotFound))
			})
		})

		When("retrieving history with a specific revision", func() {
			It("should retrieve the history up to the specified revision", func() {
				revOut := int64(0)
				Expect(valueStore.Put(ctx, "value1")).To(Succeed())
				Expect(valueStore.Put(ctx, "value2", storage.WithRevisionOut(&revOut))).To(Succeed())
				Expect(valueStore.Put(ctx, "value3")).To(Succeed())
				history, err := valueStore.History(ctx, storage.WithRevision(revOut))
				Expect(err).NotTo(HaveOccurred())
				Expect(len(history)).To(Equal(2))

				expectUpdates(
					lo.T2("", "value1"),
					lo.T2("value1", "value2"),
					lo.T2("value2", "value3"),
				)
			})
		})

		When("retrieving history with a revision that does not exist", func() {
			It("should return a NotFound error", func() {
				_, err := valueStore.History(ctx, storage.WithRevision(-1))
				Expect(err).To(Equal(storage.ErrNotFound))
			})
		})

		When("retrieving history where the current element has been deleted", func() {
			It("should return a NotFound error", func() {
				Expect(valueStore.Put(ctx, "value1")).To(Succeed())
				Expect(valueStore.Put(ctx, "value2")).To(Succeed())
				Expect(valueStore.Delete(ctx)).To(Succeed())
				_, err := valueStore.History(ctx, storage.IncludeValues(true), storage.WithRevision(3))
				Expect(err).To(Equal(storage.ErrNotFound))

				expectUpdates(
					lo.T2("", "value1"),
					lo.T2("value1", "value2"),
					lo.T2("value2", ""),
				)
			})

			When("a new element is added after the deleted element", func() {
				It("should not include the deleted element in the history", func() {
					Expect(valueStore.Put(ctx, "value1")).To(Succeed())
					Expect(valueStore.Put(ctx, "value2")).To(Succeed())
					Expect(valueStore.Delete(ctx)).To(Succeed())

					Expect(valueStore.Put(ctx, "value3")).To(Succeed())
					history, err := valueStore.History(ctx, storage.IncludeValues(true))
					Expect(err).NotTo(HaveOccurred())
					Expect(len(history)).To(Equal(1))
					Expect(history[0].Value()).To(Equal("value3"))

					expectUpdates(
						lo.T2("", "value1"),
						lo.T2("value1", "value2"),
						lo.T2("value2", ""),
						lo.T2("", "value3"),
					)
				})
			})
		})
	})
})
