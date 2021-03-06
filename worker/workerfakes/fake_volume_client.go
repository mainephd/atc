// This file was generated by counterfeiter
package workerfakes

import (
	"sync"

	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc/dbng"
	"github.com/concourse/atc/worker"
)

type FakeVolumeClient struct {
	CreateVolumeForResourceCacheStub        func(lager.Logger, worker.VolumeSpec, *dbng.UsedResourceCache) (worker.Volume, error)
	createVolumeForResourceCacheMutex       sync.RWMutex
	createVolumeForResourceCacheArgsForCall []struct {
		arg1 lager.Logger
		arg2 worker.VolumeSpec
		arg3 *dbng.UsedResourceCache
	}
	createVolumeForResourceCacheReturns struct {
		result1 worker.Volume
		result2 error
	}
	createVolumeForResourceCacheReturnsOnCall map[int]struct {
		result1 worker.Volume
		result2 error
	}
	FindOrCreateVolumeForContainerStub        func(lager.Logger, worker.VolumeSpec, dbng.CreatingContainer, int, string) (worker.Volume, error)
	findOrCreateVolumeForContainerMutex       sync.RWMutex
	findOrCreateVolumeForContainerArgsForCall []struct {
		arg1 lager.Logger
		arg2 worker.VolumeSpec
		arg3 dbng.CreatingContainer
		arg4 int
		arg5 string
	}
	findOrCreateVolumeForContainerReturns struct {
		result1 worker.Volume
		result2 error
	}
	findOrCreateVolumeForContainerReturnsOnCall map[int]struct {
		result1 worker.Volume
		result2 error
	}
	FindOrCreateCOWVolumeForContainerStub        func(lager.Logger, worker.VolumeSpec, dbng.CreatingContainer, worker.Volume, int, string) (worker.Volume, error)
	findOrCreateCOWVolumeForContainerMutex       sync.RWMutex
	findOrCreateCOWVolumeForContainerArgsForCall []struct {
		arg1 lager.Logger
		arg2 worker.VolumeSpec
		arg3 dbng.CreatingContainer
		arg4 worker.Volume
		arg5 int
		arg6 string
	}
	findOrCreateCOWVolumeForContainerReturns struct {
		result1 worker.Volume
		result2 error
	}
	findOrCreateCOWVolumeForContainerReturnsOnCall map[int]struct {
		result1 worker.Volume
		result2 error
	}
	FindOrCreateVolumeForBaseResourceTypeStub        func(lager.Logger, worker.VolumeSpec, int, string) (worker.Volume, error)
	findOrCreateVolumeForBaseResourceTypeMutex       sync.RWMutex
	findOrCreateVolumeForBaseResourceTypeArgsForCall []struct {
		arg1 lager.Logger
		arg2 worker.VolumeSpec
		arg3 int
		arg4 string
	}
	findOrCreateVolumeForBaseResourceTypeReturns struct {
		result1 worker.Volume
		result2 error
	}
	findOrCreateVolumeForBaseResourceTypeReturnsOnCall map[int]struct {
		result1 worker.Volume
		result2 error
	}
	FindInitializedVolumeForResourceCacheStub        func(lager.Logger, *dbng.UsedResourceCache) (worker.Volume, bool, error)
	findInitializedVolumeForResourceCacheMutex       sync.RWMutex
	findInitializedVolumeForResourceCacheArgsForCall []struct {
		arg1 lager.Logger
		arg2 *dbng.UsedResourceCache
	}
	findInitializedVolumeForResourceCacheReturns struct {
		result1 worker.Volume
		result2 bool
		result3 error
	}
	findInitializedVolumeForResourceCacheReturnsOnCall map[int]struct {
		result1 worker.Volume
		result2 bool
		result3 error
	}
	LookupVolumeStub        func(lager.Logger, string) (worker.Volume, bool, error)
	lookupVolumeMutex       sync.RWMutex
	lookupVolumeArgsForCall []struct {
		arg1 lager.Logger
		arg2 string
	}
	lookupVolumeReturns struct {
		result1 worker.Volume
		result2 bool
		result3 error
	}
	lookupVolumeReturnsOnCall map[int]struct {
		result1 worker.Volume
		result2 bool
		result3 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeVolumeClient) CreateVolumeForResourceCache(arg1 lager.Logger, arg2 worker.VolumeSpec, arg3 *dbng.UsedResourceCache) (worker.Volume, error) {
	fake.createVolumeForResourceCacheMutex.Lock()
	ret, specificReturn := fake.createVolumeForResourceCacheReturnsOnCall[len(fake.createVolumeForResourceCacheArgsForCall)]
	fake.createVolumeForResourceCacheArgsForCall = append(fake.createVolumeForResourceCacheArgsForCall, struct {
		arg1 lager.Logger
		arg2 worker.VolumeSpec
		arg3 *dbng.UsedResourceCache
	}{arg1, arg2, arg3})
	fake.recordInvocation("CreateVolumeForResourceCache", []interface{}{arg1, arg2, arg3})
	fake.createVolumeForResourceCacheMutex.Unlock()
	if fake.CreateVolumeForResourceCacheStub != nil {
		return fake.CreateVolumeForResourceCacheStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.createVolumeForResourceCacheReturns.result1, fake.createVolumeForResourceCacheReturns.result2
}

func (fake *FakeVolumeClient) CreateVolumeForResourceCacheCallCount() int {
	fake.createVolumeForResourceCacheMutex.RLock()
	defer fake.createVolumeForResourceCacheMutex.RUnlock()
	return len(fake.createVolumeForResourceCacheArgsForCall)
}

func (fake *FakeVolumeClient) CreateVolumeForResourceCacheArgsForCall(i int) (lager.Logger, worker.VolumeSpec, *dbng.UsedResourceCache) {
	fake.createVolumeForResourceCacheMutex.RLock()
	defer fake.createVolumeForResourceCacheMutex.RUnlock()
	return fake.createVolumeForResourceCacheArgsForCall[i].arg1, fake.createVolumeForResourceCacheArgsForCall[i].arg2, fake.createVolumeForResourceCacheArgsForCall[i].arg3
}

func (fake *FakeVolumeClient) CreateVolumeForResourceCacheReturns(result1 worker.Volume, result2 error) {
	fake.CreateVolumeForResourceCacheStub = nil
	fake.createVolumeForResourceCacheReturns = struct {
		result1 worker.Volume
		result2 error
	}{result1, result2}
}

func (fake *FakeVolumeClient) CreateVolumeForResourceCacheReturnsOnCall(i int, result1 worker.Volume, result2 error) {
	fake.CreateVolumeForResourceCacheStub = nil
	if fake.createVolumeForResourceCacheReturnsOnCall == nil {
		fake.createVolumeForResourceCacheReturnsOnCall = make(map[int]struct {
			result1 worker.Volume
			result2 error
		})
	}
	fake.createVolumeForResourceCacheReturnsOnCall[i] = struct {
		result1 worker.Volume
		result2 error
	}{result1, result2}
}

func (fake *FakeVolumeClient) FindOrCreateVolumeForContainer(arg1 lager.Logger, arg2 worker.VolumeSpec, arg3 dbng.CreatingContainer, arg4 int, arg5 string) (worker.Volume, error) {
	fake.findOrCreateVolumeForContainerMutex.Lock()
	ret, specificReturn := fake.findOrCreateVolumeForContainerReturnsOnCall[len(fake.findOrCreateVolumeForContainerArgsForCall)]
	fake.findOrCreateVolumeForContainerArgsForCall = append(fake.findOrCreateVolumeForContainerArgsForCall, struct {
		arg1 lager.Logger
		arg2 worker.VolumeSpec
		arg3 dbng.CreatingContainer
		arg4 int
		arg5 string
	}{arg1, arg2, arg3, arg4, arg5})
	fake.recordInvocation("FindOrCreateVolumeForContainer", []interface{}{arg1, arg2, arg3, arg4, arg5})
	fake.findOrCreateVolumeForContainerMutex.Unlock()
	if fake.FindOrCreateVolumeForContainerStub != nil {
		return fake.FindOrCreateVolumeForContainerStub(arg1, arg2, arg3, arg4, arg5)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.findOrCreateVolumeForContainerReturns.result1, fake.findOrCreateVolumeForContainerReturns.result2
}

func (fake *FakeVolumeClient) FindOrCreateVolumeForContainerCallCount() int {
	fake.findOrCreateVolumeForContainerMutex.RLock()
	defer fake.findOrCreateVolumeForContainerMutex.RUnlock()
	return len(fake.findOrCreateVolumeForContainerArgsForCall)
}

func (fake *FakeVolumeClient) FindOrCreateVolumeForContainerArgsForCall(i int) (lager.Logger, worker.VolumeSpec, dbng.CreatingContainer, int, string) {
	fake.findOrCreateVolumeForContainerMutex.RLock()
	defer fake.findOrCreateVolumeForContainerMutex.RUnlock()
	return fake.findOrCreateVolumeForContainerArgsForCall[i].arg1, fake.findOrCreateVolumeForContainerArgsForCall[i].arg2, fake.findOrCreateVolumeForContainerArgsForCall[i].arg3, fake.findOrCreateVolumeForContainerArgsForCall[i].arg4, fake.findOrCreateVolumeForContainerArgsForCall[i].arg5
}

func (fake *FakeVolumeClient) FindOrCreateVolumeForContainerReturns(result1 worker.Volume, result2 error) {
	fake.FindOrCreateVolumeForContainerStub = nil
	fake.findOrCreateVolumeForContainerReturns = struct {
		result1 worker.Volume
		result2 error
	}{result1, result2}
}

func (fake *FakeVolumeClient) FindOrCreateVolumeForContainerReturnsOnCall(i int, result1 worker.Volume, result2 error) {
	fake.FindOrCreateVolumeForContainerStub = nil
	if fake.findOrCreateVolumeForContainerReturnsOnCall == nil {
		fake.findOrCreateVolumeForContainerReturnsOnCall = make(map[int]struct {
			result1 worker.Volume
			result2 error
		})
	}
	fake.findOrCreateVolumeForContainerReturnsOnCall[i] = struct {
		result1 worker.Volume
		result2 error
	}{result1, result2}
}

func (fake *FakeVolumeClient) FindOrCreateCOWVolumeForContainer(arg1 lager.Logger, arg2 worker.VolumeSpec, arg3 dbng.CreatingContainer, arg4 worker.Volume, arg5 int, arg6 string) (worker.Volume, error) {
	fake.findOrCreateCOWVolumeForContainerMutex.Lock()
	ret, specificReturn := fake.findOrCreateCOWVolumeForContainerReturnsOnCall[len(fake.findOrCreateCOWVolumeForContainerArgsForCall)]
	fake.findOrCreateCOWVolumeForContainerArgsForCall = append(fake.findOrCreateCOWVolumeForContainerArgsForCall, struct {
		arg1 lager.Logger
		arg2 worker.VolumeSpec
		arg3 dbng.CreatingContainer
		arg4 worker.Volume
		arg5 int
		arg6 string
	}{arg1, arg2, arg3, arg4, arg5, arg6})
	fake.recordInvocation("FindOrCreateCOWVolumeForContainer", []interface{}{arg1, arg2, arg3, arg4, arg5, arg6})
	fake.findOrCreateCOWVolumeForContainerMutex.Unlock()
	if fake.FindOrCreateCOWVolumeForContainerStub != nil {
		return fake.FindOrCreateCOWVolumeForContainerStub(arg1, arg2, arg3, arg4, arg5, arg6)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.findOrCreateCOWVolumeForContainerReturns.result1, fake.findOrCreateCOWVolumeForContainerReturns.result2
}

func (fake *FakeVolumeClient) FindOrCreateCOWVolumeForContainerCallCount() int {
	fake.findOrCreateCOWVolumeForContainerMutex.RLock()
	defer fake.findOrCreateCOWVolumeForContainerMutex.RUnlock()
	return len(fake.findOrCreateCOWVolumeForContainerArgsForCall)
}

func (fake *FakeVolumeClient) FindOrCreateCOWVolumeForContainerArgsForCall(i int) (lager.Logger, worker.VolumeSpec, dbng.CreatingContainer, worker.Volume, int, string) {
	fake.findOrCreateCOWVolumeForContainerMutex.RLock()
	defer fake.findOrCreateCOWVolumeForContainerMutex.RUnlock()
	return fake.findOrCreateCOWVolumeForContainerArgsForCall[i].arg1, fake.findOrCreateCOWVolumeForContainerArgsForCall[i].arg2, fake.findOrCreateCOWVolumeForContainerArgsForCall[i].arg3, fake.findOrCreateCOWVolumeForContainerArgsForCall[i].arg4, fake.findOrCreateCOWVolumeForContainerArgsForCall[i].arg5, fake.findOrCreateCOWVolumeForContainerArgsForCall[i].arg6
}

func (fake *FakeVolumeClient) FindOrCreateCOWVolumeForContainerReturns(result1 worker.Volume, result2 error) {
	fake.FindOrCreateCOWVolumeForContainerStub = nil
	fake.findOrCreateCOWVolumeForContainerReturns = struct {
		result1 worker.Volume
		result2 error
	}{result1, result2}
}

func (fake *FakeVolumeClient) FindOrCreateCOWVolumeForContainerReturnsOnCall(i int, result1 worker.Volume, result2 error) {
	fake.FindOrCreateCOWVolumeForContainerStub = nil
	if fake.findOrCreateCOWVolumeForContainerReturnsOnCall == nil {
		fake.findOrCreateCOWVolumeForContainerReturnsOnCall = make(map[int]struct {
			result1 worker.Volume
			result2 error
		})
	}
	fake.findOrCreateCOWVolumeForContainerReturnsOnCall[i] = struct {
		result1 worker.Volume
		result2 error
	}{result1, result2}
}

func (fake *FakeVolumeClient) FindOrCreateVolumeForBaseResourceType(arg1 lager.Logger, arg2 worker.VolumeSpec, arg3 int, arg4 string) (worker.Volume, error) {
	fake.findOrCreateVolumeForBaseResourceTypeMutex.Lock()
	ret, specificReturn := fake.findOrCreateVolumeForBaseResourceTypeReturnsOnCall[len(fake.findOrCreateVolumeForBaseResourceTypeArgsForCall)]
	fake.findOrCreateVolumeForBaseResourceTypeArgsForCall = append(fake.findOrCreateVolumeForBaseResourceTypeArgsForCall, struct {
		arg1 lager.Logger
		arg2 worker.VolumeSpec
		arg3 int
		arg4 string
	}{arg1, arg2, arg3, arg4})
	fake.recordInvocation("FindOrCreateVolumeForBaseResourceType", []interface{}{arg1, arg2, arg3, arg4})
	fake.findOrCreateVolumeForBaseResourceTypeMutex.Unlock()
	if fake.FindOrCreateVolumeForBaseResourceTypeStub != nil {
		return fake.FindOrCreateVolumeForBaseResourceTypeStub(arg1, arg2, arg3, arg4)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.findOrCreateVolumeForBaseResourceTypeReturns.result1, fake.findOrCreateVolumeForBaseResourceTypeReturns.result2
}

func (fake *FakeVolumeClient) FindOrCreateVolumeForBaseResourceTypeCallCount() int {
	fake.findOrCreateVolumeForBaseResourceTypeMutex.RLock()
	defer fake.findOrCreateVolumeForBaseResourceTypeMutex.RUnlock()
	return len(fake.findOrCreateVolumeForBaseResourceTypeArgsForCall)
}

func (fake *FakeVolumeClient) FindOrCreateVolumeForBaseResourceTypeArgsForCall(i int) (lager.Logger, worker.VolumeSpec, int, string) {
	fake.findOrCreateVolumeForBaseResourceTypeMutex.RLock()
	defer fake.findOrCreateVolumeForBaseResourceTypeMutex.RUnlock()
	return fake.findOrCreateVolumeForBaseResourceTypeArgsForCall[i].arg1, fake.findOrCreateVolumeForBaseResourceTypeArgsForCall[i].arg2, fake.findOrCreateVolumeForBaseResourceTypeArgsForCall[i].arg3, fake.findOrCreateVolumeForBaseResourceTypeArgsForCall[i].arg4
}

func (fake *FakeVolumeClient) FindOrCreateVolumeForBaseResourceTypeReturns(result1 worker.Volume, result2 error) {
	fake.FindOrCreateVolumeForBaseResourceTypeStub = nil
	fake.findOrCreateVolumeForBaseResourceTypeReturns = struct {
		result1 worker.Volume
		result2 error
	}{result1, result2}
}

func (fake *FakeVolumeClient) FindOrCreateVolumeForBaseResourceTypeReturnsOnCall(i int, result1 worker.Volume, result2 error) {
	fake.FindOrCreateVolumeForBaseResourceTypeStub = nil
	if fake.findOrCreateVolumeForBaseResourceTypeReturnsOnCall == nil {
		fake.findOrCreateVolumeForBaseResourceTypeReturnsOnCall = make(map[int]struct {
			result1 worker.Volume
			result2 error
		})
	}
	fake.findOrCreateVolumeForBaseResourceTypeReturnsOnCall[i] = struct {
		result1 worker.Volume
		result2 error
	}{result1, result2}
}

func (fake *FakeVolumeClient) FindInitializedVolumeForResourceCache(arg1 lager.Logger, arg2 *dbng.UsedResourceCache) (worker.Volume, bool, error) {
	fake.findInitializedVolumeForResourceCacheMutex.Lock()
	ret, specificReturn := fake.findInitializedVolumeForResourceCacheReturnsOnCall[len(fake.findInitializedVolumeForResourceCacheArgsForCall)]
	fake.findInitializedVolumeForResourceCacheArgsForCall = append(fake.findInitializedVolumeForResourceCacheArgsForCall, struct {
		arg1 lager.Logger
		arg2 *dbng.UsedResourceCache
	}{arg1, arg2})
	fake.recordInvocation("FindInitializedVolumeForResourceCache", []interface{}{arg1, arg2})
	fake.findInitializedVolumeForResourceCacheMutex.Unlock()
	if fake.FindInitializedVolumeForResourceCacheStub != nil {
		return fake.FindInitializedVolumeForResourceCacheStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	return fake.findInitializedVolumeForResourceCacheReturns.result1, fake.findInitializedVolumeForResourceCacheReturns.result2, fake.findInitializedVolumeForResourceCacheReturns.result3
}

func (fake *FakeVolumeClient) FindInitializedVolumeForResourceCacheCallCount() int {
	fake.findInitializedVolumeForResourceCacheMutex.RLock()
	defer fake.findInitializedVolumeForResourceCacheMutex.RUnlock()
	return len(fake.findInitializedVolumeForResourceCacheArgsForCall)
}

func (fake *FakeVolumeClient) FindInitializedVolumeForResourceCacheArgsForCall(i int) (lager.Logger, *dbng.UsedResourceCache) {
	fake.findInitializedVolumeForResourceCacheMutex.RLock()
	defer fake.findInitializedVolumeForResourceCacheMutex.RUnlock()
	return fake.findInitializedVolumeForResourceCacheArgsForCall[i].arg1, fake.findInitializedVolumeForResourceCacheArgsForCall[i].arg2
}

func (fake *FakeVolumeClient) FindInitializedVolumeForResourceCacheReturns(result1 worker.Volume, result2 bool, result3 error) {
	fake.FindInitializedVolumeForResourceCacheStub = nil
	fake.findInitializedVolumeForResourceCacheReturns = struct {
		result1 worker.Volume
		result2 bool
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeVolumeClient) FindInitializedVolumeForResourceCacheReturnsOnCall(i int, result1 worker.Volume, result2 bool, result3 error) {
	fake.FindInitializedVolumeForResourceCacheStub = nil
	if fake.findInitializedVolumeForResourceCacheReturnsOnCall == nil {
		fake.findInitializedVolumeForResourceCacheReturnsOnCall = make(map[int]struct {
			result1 worker.Volume
			result2 bool
			result3 error
		})
	}
	fake.findInitializedVolumeForResourceCacheReturnsOnCall[i] = struct {
		result1 worker.Volume
		result2 bool
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeVolumeClient) LookupVolume(arg1 lager.Logger, arg2 string) (worker.Volume, bool, error) {
	fake.lookupVolumeMutex.Lock()
	ret, specificReturn := fake.lookupVolumeReturnsOnCall[len(fake.lookupVolumeArgsForCall)]
	fake.lookupVolumeArgsForCall = append(fake.lookupVolumeArgsForCall, struct {
		arg1 lager.Logger
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("LookupVolume", []interface{}{arg1, arg2})
	fake.lookupVolumeMutex.Unlock()
	if fake.LookupVolumeStub != nil {
		return fake.LookupVolumeStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	return fake.lookupVolumeReturns.result1, fake.lookupVolumeReturns.result2, fake.lookupVolumeReturns.result3
}

func (fake *FakeVolumeClient) LookupVolumeCallCount() int {
	fake.lookupVolumeMutex.RLock()
	defer fake.lookupVolumeMutex.RUnlock()
	return len(fake.lookupVolumeArgsForCall)
}

func (fake *FakeVolumeClient) LookupVolumeArgsForCall(i int) (lager.Logger, string) {
	fake.lookupVolumeMutex.RLock()
	defer fake.lookupVolumeMutex.RUnlock()
	return fake.lookupVolumeArgsForCall[i].arg1, fake.lookupVolumeArgsForCall[i].arg2
}

func (fake *FakeVolumeClient) LookupVolumeReturns(result1 worker.Volume, result2 bool, result3 error) {
	fake.LookupVolumeStub = nil
	fake.lookupVolumeReturns = struct {
		result1 worker.Volume
		result2 bool
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeVolumeClient) LookupVolumeReturnsOnCall(i int, result1 worker.Volume, result2 bool, result3 error) {
	fake.LookupVolumeStub = nil
	if fake.lookupVolumeReturnsOnCall == nil {
		fake.lookupVolumeReturnsOnCall = make(map[int]struct {
			result1 worker.Volume
			result2 bool
			result3 error
		})
	}
	fake.lookupVolumeReturnsOnCall[i] = struct {
		result1 worker.Volume
		result2 bool
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeVolumeClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.createVolumeForResourceCacheMutex.RLock()
	defer fake.createVolumeForResourceCacheMutex.RUnlock()
	fake.findOrCreateVolumeForContainerMutex.RLock()
	defer fake.findOrCreateVolumeForContainerMutex.RUnlock()
	fake.findOrCreateCOWVolumeForContainerMutex.RLock()
	defer fake.findOrCreateCOWVolumeForContainerMutex.RUnlock()
	fake.findOrCreateVolumeForBaseResourceTypeMutex.RLock()
	defer fake.findOrCreateVolumeForBaseResourceTypeMutex.RUnlock()
	fake.findInitializedVolumeForResourceCacheMutex.RLock()
	defer fake.findInitializedVolumeForResourceCacheMutex.RUnlock()
	fake.lookupVolumeMutex.RLock()
	defer fake.lookupVolumeMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeVolumeClient) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ worker.VolumeClient = new(FakeVolumeClient)
