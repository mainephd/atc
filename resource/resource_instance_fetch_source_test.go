package resource_test

import (
	"errors"
	"os"

	"code.cloudfoundry.org/garden"
	gfakes "code.cloudfoundry.org/garden/gardenfakes"
	"code.cloudfoundry.org/lager/lagertest"
	"github.com/concourse/atc"
	. "github.com/concourse/atc/resource"
	"github.com/concourse/atc/resource/resourcefakes"
	"github.com/concourse/atc/worker/workerfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("VolumeFetchSource", func() {
	var (
		fetchSource FetchSource

		fakeContainer        *workerfakes.FakeContainer
		resourceOptions      *resourcefakes.FakeResourceOptions
		fakeVolume           *workerfakes.FakeVolume
		fakeResourceInstance *resourcefakes.FakeResourceInstance
		fakeWorker           *workerfakes.FakeWorker

		signals <-chan os.Signal
		ready   chan<- struct{}
	)

	BeforeEach(func() {
		logger := lagertest.NewTestLogger("test")
		fakeContainer = new(workerfakes.FakeContainer)
		resourceOptions = new(resourcefakes.FakeResourceOptions)
		signals = make(<-chan os.Signal)
		ready = make(chan<- struct{})

		fakeContainer.PropertyReturns("", errors.New("nope"))
		inProcess := new(gfakes.FakeProcess)
		inProcess.IDReturns("process-id")
		inProcess.WaitStub = func() (int, error) {
			return 0, nil
		}

		fakeContainer.RunStub = func(spec garden.ProcessSpec, io garden.ProcessIO) (garden.Process, error) {
			_, err := io.Stdout.Write([]byte("{}"))
			Expect(err).NotTo(HaveOccurred())

			return inProcess, nil
		}

		fakeWorker = new(workerfakes.FakeWorker)
		fakeWorker.CreateResourceGetContainerReturns(fakeContainer, nil)

		fakeVolume = new(workerfakes.FakeVolume)
		fakeResourceInstance = new(resourcefakes.FakeResourceInstance)
		fakeResourceInstance.CreateOnReturns(fakeVolume, nil)
		fetchSource = NewResourceInstanceFetchSource(
			logger,
			fakeResourceInstance,
			fakeWorker,
			resourceOptions,
			nil,
			atc.Tags{},
			42,
			Session{},
			EmptyMetadata{},
			new(workerfakes.FakeImageFetchingDelegate),
		)
	})

	Describe("IsInitialized", func() {
		Context("when there is initialized volume", func() {
			BeforeEach(func() {
				fakeResourceInstance.FindInitializedOnReturns(fakeVolume, true, nil)
			})

			It("finds initialized volume and sets versioned source", func() {
				initialized, err := fetchSource.IsInitialized()
				Expect(err).NotTo(HaveOccurred())
				Expect(initialized).To(BeTrue())
				Expect(fetchSource.VersionedSource()).NotTo(BeNil())
			})
		})

		Context("when there is no initialized volume", func() {
			BeforeEach(func() {
				fakeResourceInstance.FindInitializedOnReturns(nil, false, nil)
			})

			It("does not find initialized volume", func() {
				initialized, err := fetchSource.IsInitialized()
				Expect(err).NotTo(HaveOccurred())
				Expect(initialized).To(BeFalse())
			})
		})
	})

	Describe("Initialize", func() {
		var initErr error

		BeforeEach(func() {
			resourceOptions.ResourceTypeReturns(ResourceType("fake-resource-type"))
		})

		JustBeforeEach(func() {
			initErr = fetchSource.Initialize(signals, ready)
		})

		Context("when there is initialized volume", func() {
			BeforeEach(func() {
				fakeResourceInstance.FindInitializedOnReturns(fakeVolume, true, nil)
			})

			It("does not fetch resource", func() {
				Expect(initErr).NotTo(HaveOccurred())
				Expect(fakeResourceInstance.CreateOnCallCount()).To(Equal(0))
				Expect(fakeContainer.RunCallCount()).To(Equal(0))
			})

			It("finds initialized volume and sets versioned source", func() {
				Expect(initErr).NotTo(HaveOccurred())
				Expect(fetchSource.VersionedSource()).NotTo(BeNil())
			})
		})

		Context("when there is no initialized volume", func() {
			BeforeEach(func() {
				fakeResourceInstance.FindInitializedOnReturns(nil, false, nil)
			})

			It("creates volume for resource instance on provided worker", func() {
				Expect(initErr).NotTo(HaveOccurred())
				Expect(fakeResourceInstance.CreateOnCallCount()).To(Equal(1))
				_, worker := fakeResourceInstance.CreateOnArgsForCall(0)
				Expect(worker).To(Equal(fakeWorker))
			})

			It("creates container with volume and worker", func() {
				Expect(initErr).NotTo(HaveOccurred())
				Expect(fakeWorker.CreateResourceGetContainerCallCount()).To(Equal(1))
			})

			It("fetches versioned source", func() {
				Expect(initErr).NotTo(HaveOccurred())
				Expect(fakeContainer.RunCallCount()).To(Equal(1))
			})

			It("initializes cache", func() {
				Expect(initErr).NotTo(HaveOccurred())
				Expect(fakeVolume.InitializeCallCount()).To(Equal(1))
			})

			Context("when getting resource fails with ErrAborted", func() {
				BeforeEach(func() {
					fakeContainer.RunReturns(nil, ErrAborted)
				})

				It("returns ErrInterrupted", func() {
					Expect(initErr).To(HaveOccurred())
					Expect(initErr).To(Equal(ErrInterrupted))
				})
			})

			Context("when getting resource fails with other error", func() {
				var disaster error

				BeforeEach(func() {
					disaster = errors.New("failed")
					fakeContainer.RunReturns(nil, disaster)
				})

				It("returns the error", func() {
					Expect(initErr).To(HaveOccurred())
					Expect(initErr).To(Equal(disaster))
				})
			})
		})
	})
})
