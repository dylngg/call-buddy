package telephono

import (
	"os"
	"strings"
)

type ContextContributor interface {
	/*
		Contribute returns the name of the contribution and the object (Map or struct)

		For example, the environment variable contributor should return "Env", a map of key values, nil
	*/
	Contribute() (string, interface{}, error)
}

type SimpleContributor struct {
	backing map[string]string
}

func (s *SimpleContributor) Contribute() (string, interface{}, error) {
	panic("implement me")
}

type EnvironmentContributor struct {
	cache  map[string]string
	cached bool
}

//refresh will go pull all of the environment variables
func (e *EnvironmentContributor) refresh() error {
	e.cache = make(map[string]string)

	for _, kv := range os.Environ() {
		parts := strings.SplitN(kv, "=", 2)
		e.cache[parts[0]] = parts[1]
	}

	return nil
}

func (e EnvironmentContributor) Contribute() (string, interface{}, error) {
	if !e.cached {
		_ = e.refresh()
		e.cached = true
	}

	return "Env", e.cache, nil
}
