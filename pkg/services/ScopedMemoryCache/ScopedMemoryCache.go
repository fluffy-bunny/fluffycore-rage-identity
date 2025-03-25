package ScopedMemoryCache

import (
	"sync"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_cache "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cache"
)

type (
	service struct {
		store sync.Map
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_cache.IScopedMemoryCache = stemService
}
func (s *service) Ctor() (contracts_cache.IScopedMemoryCache, error) {
	return &service{}, nil
}

func AddScopedIScopedMemoryCache(cb di.ContainerBuilder) {
	di.AddScoped[contracts_cache.IScopedMemoryCache](cb, stemService.Ctor)
}

func (s *service) Set(key string, value any) {
	s.store.Store(key, value)

}
func (s *service) Get(key string) (any, bool) {
	return s.store.Load(key)
}
func (s *service) Delete(key string) {
	s.store.Delete(key)
}
func (s *service) Clear() {
	s.store.Clear()
}
