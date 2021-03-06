package scheduler_test

import (
	"errors"
	"time"

	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/lager/lagertest"
	"github.com/concourse/atc"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/db/algorithm"
	"github.com/concourse/atc/db/dbfakes"
	"github.com/concourse/atc/engine"
	"github.com/concourse/atc/engine/enginefakes"
	"github.com/concourse/atc/scheduler"
	"github.com/concourse/atc/scheduler/inputmapper/inputmapperfakes"
	"github.com/concourse/atc/scheduler/maxinflight/maxinflightfakes"
	"github.com/concourse/atc/scheduler/schedulerfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("I'm a BuildStarter", func() {
	var (
		fakeDB           *schedulerfakes.FakeBuildStarterDB
		fakeUpdater      *maxinflightfakes.FakeUpdater
		fakeFactory      *schedulerfakes.FakeBuildFactory
		fakeEngine       *enginefakes.FakeEngine
		pendingBuilds    []db.Build
		fakeScanner      *schedulerfakes.FakeScanner
		fakeInputMapper  *inputmapperfakes.FakeInputMapper
		fakeBuildStarter *schedulerfakes.FakeBuildStarter

		buildStarter scheduler.BuildStarter

		disaster error
	)

	BeforeEach(func() {
		fakeDB = new(schedulerfakes.FakeBuildStarterDB)
		fakeUpdater = new(maxinflightfakes.FakeUpdater)
		fakeFactory = new(schedulerfakes.FakeBuildFactory)
		fakeEngine = new(enginefakes.FakeEngine)
		fakeScanner = new(schedulerfakes.FakeScanner)
		fakeInputMapper = new(inputmapperfakes.FakeInputMapper)
		fakeBuildStarter = new(schedulerfakes.FakeBuildStarter)

		buildStarter = scheduler.NewBuildStarter(fakeDB, fakeUpdater, fakeFactory, fakeScanner, fakeInputMapper, fakeEngine)

		disaster = errors.New("bad thing")
	})

	Describe("TryStartPendingBuildsForJob", func() {
		var tryStartErr error
		var createdBuild *dbfakes.FakeBuild
		var jobConfig atc.JobConfig
		var versionedResourceTypes atc.VersionedResourceTypes

		BeforeEach(func() {
			versionedResourceTypes = atc.VersionedResourceTypes{
				{
					ResourceType: atc.ResourceType{Name: "some-resource-type"},
					Version:      atc.Version{"some": "version"},
				},
			}
		})

		Context("when manually triggered", func() {
			BeforeEach(func() {
				jobConfig = atc.JobConfig{Name: "some-job", Plan: atc.PlanSequence{{Get: "input-1"}, {Get: "input-2"}}}

				createdBuild = new(dbfakes.FakeBuild)
				createdBuild.IDReturns(66)
				createdBuild.IsManuallyTriggeredReturns(true)

				pendingBuilds = []db.Build{createdBuild}
			})

			JustBeforeEach(func() {
				tryStartErr = buildStarter.TryStartPendingBuildsForJob(
					lagertest.NewTestLogger("test"),
					jobConfig,
					atc.ResourceConfigs{{Name: "some-resource"}},
					versionedResourceTypes,
					pendingBuilds,
				)
			})

			It("updates max in flight for the job", func() {
				Expect(fakeUpdater.UpdateMaxInFlightReachedCallCount()).To(Equal(1))
				_, actualJobConfig, actualBuildID := fakeUpdater.UpdateMaxInFlightReachedArgsForCall(0)
				Expect(actualJobConfig).To(Equal(jobConfig))
				Expect(actualBuildID).To(Equal(66))
			})

			Context("when max in flight is reached", func() {
				BeforeEach(func() {
					fakeUpdater.UpdateMaxInFlightReachedReturns(true, nil)
				})

				It("does not run resource check", func() {
					Expect(fakeScanner.ScanCallCount()).To(Equal(0))
				})
			})

			Context("when max in flight is not reached", func() {
				BeforeEach(func() {
					fakeUpdater.UpdateMaxInFlightReachedReturns(false, nil)
				})

				It("runs resource check for every job resource", func() {
					Expect(fakeScanner.ScanCallCount()).To(Equal(2))
				})

				Context("when resource checking fails", func() {
					BeforeEach(func() {
						fakeScanner.ScanReturns(disaster)
					})

					It("returns an error", func() {
						Expect(tryStartErr).To(Equal(disaster))
					})
				})

				Context("when resource checking succeeds", func() {
					BeforeEach(func() {
						fakeScanner.ScanStub = func(lager.Logger, string) error {
							defer GinkgoRecover()
							Expect(fakeDB.LoadVersionsDBCallCount()).To(BeZero())
							return nil
						}
					})

					Context("when loading the versions DB fails", func() {
						BeforeEach(func() {
							fakeDB.LoadVersionsDBReturns(nil, disaster)
						})

						It("returns an error", func() {
							Expect(tryStartErr).To(Equal(disaster))
						})

						It("checked for the right resources", func() {
							Expect(fakeScanner.ScanCallCount()).To(Equal(2))
							_, resource1 := fakeScanner.ScanArgsForCall(0)
							_, resource2 := fakeScanner.ScanArgsForCall(1)
							Expect([]string{resource1, resource2}).To(ConsistOf("input-1", "input-2"))
						})

						It("loaded the versions DB after checking all the resources", func() {
							Expect(fakeDB.LoadVersionsDBCallCount()).To(Equal(1))
						})
					})

					Context("when loading the versions DB succeeds", func() {
						var versionsDB *algorithm.VersionsDB

						BeforeEach(func() {
							fakeDB.LoadVersionsDBReturns(&algorithm.VersionsDB{
								ResourceVersions: []algorithm.ResourceVersion{
									{
										VersionID:  73,
										ResourceID: 127,
										CheckOrder: 123,
									},
								},
								BuildOutputs: []algorithm.BuildOutput{
									{
										ResourceVersion: algorithm.ResourceVersion{
											VersionID:  73,
											ResourceID: 127,
											CheckOrder: 123,
										},
										BuildID: 66,
										JobID:   13,
									},
								},
								BuildInputs: []algorithm.BuildInput{
									{
										ResourceVersion: algorithm.ResourceVersion{
											VersionID:  66,
											ResourceID: 77,
											CheckOrder: 88,
										},
										BuildID:   66,
										JobID:     13,
										InputName: "some-input-name",
									},
								},
								JobIDs: map[string]int{
									"bad-luck-job": 13,
								},
								ResourceIDs: map[string]int{
									"resource-127": 127,
								},
								CachedAt: time.Unix(42, 0).UTC(),
							}, nil)

							versionsDB = &algorithm.VersionsDB{JobIDs: map[string]int{"j1": 1}}
							fakeDB.LoadVersionsDBReturns(versionsDB, nil)
						})

						Context("when saving the next input mapping fails", func() {
							BeforeEach(func() {
								fakeInputMapper.SaveNextInputMappingReturns(nil, disaster)
							})

							It("saved the next input mapping for the right job and versions", func() {
								Expect(fakeInputMapper.SaveNextInputMappingCallCount()).To(Equal(1))
								_, actualVersionsDB, actualJobConfig := fakeInputMapper.SaveNextInputMappingArgsForCall(0)
								Expect(actualVersionsDB).To(Equal(versionsDB))
								Expect(actualJobConfig).To(Equal(jobConfig))
							})
						})

						Context("when saving the next input mapping succeeds", func() {
							BeforeEach(func() {
								fakeInputMapper.SaveNextInputMappingStub = func(lager.Logger, *algorithm.VersionsDB, atc.JobConfig) (algorithm.InputMapping, error) {
									defer GinkgoRecover()
									return nil, nil
								}
							})

							It("saved the next input mapping and returns the build", func() {
								Expect(tryStartErr).NotTo(HaveOccurred())
							})
						})
					})
				})
			})
		})

		Context("when not manually triggered", func() {
			JustBeforeEach(func() {
				tryStartErr = buildStarter.TryStartPendingBuildsForJob(
					lagertest.NewTestLogger("test"),
					atc.JobConfig{Name: "some-job"},
					atc.ResourceConfigs{{Name: "some-resource"}},
					atc.VersionedResourceTypes{
						{
							ResourceType: atc.ResourceType{Name: "some-resource-type"},
							Version:      atc.Version{"some": "version"},
						},
					},
					pendingBuilds,
				)
			})

			itReturnsTheError := func() {
				It("returns the error", func() {
					Expect(tryStartErr).To(Equal(disaster))
				})
			}

			itDoesntReturnAnErrorOrMarkTheBuildAsScheduled := func() {
				It("doesn't return an error", func() {
					Expect(tryStartErr).NotTo(HaveOccurred())
				})

				It("doesn't try to mark the build as scheduled", func() {
					Expect(fakeDB.UpdateBuildToScheduledCallCount()).To(BeZero())
				})
			}

			itUpdatedMaxInFlightForAllBuilds := func() {
				It("updated max in flight for the right jobs", func() {
					Expect(fakeUpdater.UpdateMaxInFlightReachedCallCount()).To(Equal(3))
					_, actualJobConfig, actualBuildID := fakeUpdater.UpdateMaxInFlightReachedArgsForCall(0)
					Expect(actualJobConfig).To(Equal(atc.JobConfig{Name: "some-job"}))
					Expect(actualBuildID).To(Equal(99))

					_, actualJobConfig, actualBuildID = fakeUpdater.UpdateMaxInFlightReachedArgsForCall(1)
					Expect(actualJobConfig).To(Equal(atc.JobConfig{Name: "some-job"}))
					Expect(actualBuildID).To(Equal(999))
				})
			}

			itUpdatedMaxInFlightForTheFirstBuild := func() {
				It("updated max in flight for the first jobs", func() {
					Expect(fakeUpdater.UpdateMaxInFlightReachedCallCount()).To(Equal(1))
					_, actualJobConfig, actualBuildID := fakeUpdater.UpdateMaxInFlightReachedArgsForCall(0)
					Expect(actualJobConfig).To(Equal(atc.JobConfig{Name: "some-job"}))
					Expect(actualBuildID).To(Equal(99))
				})
			}

			Context("when the stars align", func() {
				BeforeEach(func() {
					fakeUpdater.UpdateMaxInFlightReachedReturns(false, nil)
					fakeDB.GetNextBuildInputsReturns([]db.BuildInput{{Name: "some-input"}}, true, nil)
					fakeDB.IsPausedReturns(false, nil)
					fakeDB.GetJobReturns(db.SavedJob{Paused: false}, true, nil)
				})

				Context("when there are several pending builds", func() {
					var pendingBuild1 *dbfakes.FakeBuild
					var pendingBuild2 *dbfakes.FakeBuild
					var pendingBuild3 *dbfakes.FakeBuild

					BeforeEach(func() {
						pendingBuild1 = new(dbfakes.FakeBuild)
						pendingBuild1.IDReturns(99)
						pendingBuild2 = new(dbfakes.FakeBuild)
						pendingBuild2.IDReturns(999)
						pendingBuild3 = new(dbfakes.FakeBuild)
						pendingBuild3.IDReturns(555)
						pendingBuilds = []db.Build{pendingBuild1, pendingBuild2, pendingBuild3}
					})

					Context("when marking the build as scheduled fails", func() {
						BeforeEach(func() {
							fakeDB.UpdateBuildToScheduledReturns(false, disaster)
						})

						It("returns the error", func() {
							Expect(tryStartErr).To(Equal(disaster))
						})

						It("marked the right build as scheduled", func() {
							Expect(fakeDB.UpdateBuildToScheduledCallCount()).To(Equal(1))
							Expect(fakeDB.UpdateBuildToScheduledArgsForCall(0)).To(Equal(99))
						})
					})

					Context("when someone else already scheduled the build", func() {
						BeforeEach(func() {
							fakeDB.UpdateBuildToScheduledReturns(false, nil)
						})

						It("doesn't return an error", func() {
							Expect(tryStartErr).NotTo(HaveOccurred())
						})

						It("doesn't try to use inputs for build", func() {
							Expect(fakeDB.UseInputsForBuildCallCount()).To(BeZero())
						})
					})

					Context("when marking the build as scheduled succeeds", func() {
						BeforeEach(func() {
							fakeDB.UpdateBuildToScheduledReturns(true, nil)
						})

						Context("when using inputs for build fails", func() {
							BeforeEach(func() {
								fakeDB.UseInputsForBuildReturns(disaster)
							})

							It("returns the error", func() {
								Expect(tryStartErr).To(Equal(disaster))
							})

							It("used the right inputs for the right build", func() {
								Expect(fakeDB.UseInputsForBuildCallCount()).To(Equal(1))
								actualBuildID, actualInputs := fakeDB.UseInputsForBuildArgsForCall(0)
								Expect(actualBuildID).To(Equal(99))
								Expect(actualInputs).To(Equal([]db.BuildInput{{Name: "some-input"}}))
							})
						})

						Context("when using inputs for build succeeds", func() {
							BeforeEach(func() {
								fakeDB.UseInputsForBuildReturns(nil)
							})

							Context("when creating the build plan fails", func() {
								BeforeEach(func() {
									fakeFactory.CreateReturns(atc.Plan{}, disaster)
								})

								It("stops creating builds for job", func() {
									Expect(fakeFactory.CreateCallCount()).To(Equal(1))
									actualJobConfig, actualResourceConfigs, actualResourceTypes, actualBuildInputs := fakeFactory.CreateArgsForCall(0)
									Expect(actualJobConfig).To(Equal(atc.JobConfig{Name: "some-job"}))
									Expect(actualResourceConfigs).To(Equal(atc.ResourceConfigs{{Name: "some-resource"}}))
									Expect(actualResourceTypes).To(Equal(versionedResourceTypes))
									Expect(actualBuildInputs).To(Equal([]db.BuildInput{{Name: "some-input"}}))
								})

								Context("when marking the build as errored fails", func() {
									BeforeEach(func() {
										pendingBuild1.FinishReturns(disaster)
									})

									It("doesn't return an error", func() {
										Expect(tryStartErr).NotTo(HaveOccurred())
									})

									It("marked the right build as errored", func() {
										Expect(pendingBuild1.FinishCallCount()).To(Equal(1))
										actualStatus := pendingBuild1.FinishArgsForCall(0)
										Expect(actualStatus).To(Equal(db.StatusErrored))
									})
								})

								Context("when marking the build as errored succeeds", func() {
									BeforeEach(func() {
										pendingBuild1.FinishReturns(nil)
									})

									It("doesn't return an error", func() {
										Expect(tryStartErr).NotTo(HaveOccurred())
									})
								})
							})

							Context("when creating the build plan succeeds", func() {
								BeforeEach(func() {
									fakeFactory.CreateReturns(atc.Plan{Task: &atc.TaskPlan{ConfigPath: "some-task-1.yml"}}, nil)
									fakeEngine.CreateBuildReturns(new(enginefakes.FakeBuild), nil)
								})

								It("creates build plans for all builds", func() {
									Expect(fakeFactory.CreateCallCount()).To(Equal(3))
									actualJobConfig, actualResourceConfigs, actualResourceTypes, actualBuildInputs := fakeFactory.CreateArgsForCall(0)
									Expect(actualJobConfig).To(Equal(atc.JobConfig{Name: "some-job"}))
									Expect(actualResourceConfigs).To(Equal(atc.ResourceConfigs{{Name: "some-resource"}}))
									Expect(actualResourceTypes).To(Equal(versionedResourceTypes))
									Expect(actualBuildInputs).To(Equal([]db.BuildInput{{Name: "some-input"}}))

									actualJobConfig, actualResourceConfigs, actualResourceTypes, actualBuildInputs = fakeFactory.CreateArgsForCall(1)
									Expect(actualJobConfig).To(Equal(atc.JobConfig{Name: "some-job"}))
									Expect(actualResourceConfigs).To(Equal(atc.ResourceConfigs{{Name: "some-resource"}}))
									Expect(actualResourceTypes).To(Equal(versionedResourceTypes))
									Expect(actualBuildInputs).To(Equal([]db.BuildInput{{Name: "some-input"}}))

									actualJobConfig, actualResourceConfigs, actualResourceTypes, actualBuildInputs = fakeFactory.CreateArgsForCall(2)
									Expect(actualJobConfig).To(Equal(atc.JobConfig{Name: "some-job"}))
									Expect(actualResourceConfigs).To(Equal(atc.ResourceConfigs{{Name: "some-resource"}}))
									Expect(actualResourceTypes).To(Equal(versionedResourceTypes))
									Expect(actualBuildInputs).To(Equal([]db.BuildInput{{Name: "some-input"}}))
								})

								Context("when creating the engine build fails", func() {
									BeforeEach(func() {
										fakeEngine.CreateBuildReturns(nil, disaster)
									})

									It("doesn't return an error", func() {
										Expect(tryStartErr).NotTo(HaveOccurred())
									})
								})

								Context("when creating the engine build succeeds", func() {
									var engineBuild1 *enginefakes.FakeBuild
									var engineBuild2 *enginefakes.FakeBuild
									var engineBuild3 *enginefakes.FakeBuild

									BeforeEach(func() {
										engineBuild1 = new(enginefakes.FakeBuild)
										engineBuild2 = new(enginefakes.FakeBuild)
										engineBuild3 = new(enginefakes.FakeBuild)
										createBuildCallCount := 0
										fakeEngine.CreateBuildStub = func(lager.Logger, db.Build, atc.Plan) (engine.Build, error) {
											createBuildCallCount++
											switch createBuildCallCount {
											case 1:
												return engineBuild1, nil
											case 2:
												return engineBuild2, nil
											case 3:
												return engineBuild3, nil
											default:
												panic("unexpected-call-count-for-create-build")
											}
										}
									})

									It("doesn't return an error", func() {
										Expect(tryStartErr).NotTo(HaveOccurred())
									})

									itUpdatedMaxInFlightForAllBuilds()

									It("created the engine build with the right build and plan", func() {
										Expect(fakeEngine.CreateBuildCallCount()).To(Equal(3))
										_, actualBuild, actualPlan := fakeEngine.CreateBuildArgsForCall(0)
										Expect(actualBuild).To(Equal(pendingBuild1))
										Expect(actualPlan).To(Equal(atc.Plan{Task: &atc.TaskPlan{ConfigPath: "some-task-1.yml"}}))

										_, actualBuild, actualPlan = fakeEngine.CreateBuildArgsForCall(1)
										Expect(actualBuild).To(Equal(pendingBuild2))
										Expect(actualPlan).To(Equal(atc.Plan{Task: &atc.TaskPlan{ConfigPath: "some-task-1.yml"}}))

										_, actualBuild, actualPlan = fakeEngine.CreateBuildArgsForCall(2)
										Expect(actualBuild).To(Equal(pendingBuild3))
										Expect(actualPlan).To(Equal(atc.Plan{Task: &atc.TaskPlan{ConfigPath: "some-task-1.yml"}}))
									})

									It("starts the engine build (asynchronously)", func() {
										Eventually(engineBuild1.ResumeCallCount).Should(Equal(1))
										Eventually(engineBuild2.ResumeCallCount).Should(Equal(1))
										Eventually(engineBuild3.ResumeCallCount).Should(Equal(1))
									})
								})
							})
						})
					})

					Context("when updating max in flight reached fails", func() {
						BeforeEach(func() {
							fakeUpdater.UpdateMaxInFlightReachedReturns(false, disaster)
						})

						itReturnsTheError()
						itUpdatedMaxInFlightForTheFirstBuild()
					})

					Context("when max in flight is reached", func() {
						BeforeEach(func() {
							fakeUpdater.UpdateMaxInFlightReachedReturns(true, nil)
						})

						itDoesntReturnAnErrorOrMarkTheBuildAsScheduled()
					})

					Context("when getting the next build inputs fails", func() {
						BeforeEach(func() {
							fakeDB.GetNextBuildInputsReturns(nil, false, disaster)
						})

						itReturnsTheError()
						itUpdatedMaxInFlightForTheFirstBuild()
					})

					Context("when there are no next build inputs", func() {
						BeforeEach(func() {
							fakeDB.GetNextBuildInputsReturns(nil, false, nil)
						})

						itDoesntReturnAnErrorOrMarkTheBuildAsScheduled()
						itUpdatedMaxInFlightForTheFirstBuild()
					})

					Context("when checking if the pipeline is paused fails", func() {
						BeforeEach(func() {
							fakeDB.IsPausedReturns(false, disaster)
						})

						itReturnsTheError()
						itUpdatedMaxInFlightForTheFirstBuild()
					})

					Context("when the pipeline is paused", func() {
						BeforeEach(func() {
							fakeDB.IsPausedReturns(true, nil)
						})

						itDoesntReturnAnErrorOrMarkTheBuildAsScheduled()
						itUpdatedMaxInFlightForTheFirstBuild()
					})

					Context("when getting the job fails", func() {
						BeforeEach(func() {
							fakeDB.GetJobReturns(db.SavedJob{}, false, disaster)
						})

						itReturnsTheError()
						itUpdatedMaxInFlightForTheFirstBuild()
					})

					Context("when the job is paused", func() {
						BeforeEach(func() {
							fakeDB.GetJobReturns(db.SavedJob{Paused: true}, true, nil)
						})

						itDoesntReturnAnErrorOrMarkTheBuildAsScheduled()
						itUpdatedMaxInFlightForTheFirstBuild()
					})
				})
			})
		})
	})

})
