package telephono

import (
	"errors"
)

/*VariableResolver is an interface that the various different variable resolvers can use to let an Expander
know that they can handle it

e.g. one could source from an internal store, an environment or the systems EnvironmentVariables
*/
type VariableResolver interface {
	Order() int
	CanHandle(variableName string) bool
	Resolve(variableName string) (string, error)
}

type simpleInMemResolver struct {
	backing map[string]string
}

func newSimpleInMemResolver() simpleInMemResolver {
	return simpleInMemResolver{
		backing: map[string]string{},
	}
}

func (s *simpleInMemResolver) Order() int {
	return -10
}

func (s *simpleInMemResolver) CanHandle(variableName string) bool {
	_, found := s.backing[variableName]
	return found
}

func (s *simpleInMemResolver) Resolve(variableName string) (string, error) {
	if found, nada := s.backing[variableName]; nada == true {
		return found, nil
	} else {
		return "", errors.New("didn't resolve " + variableName)
	}
}

func (s *simpleInMemResolver) addResolution(key, value string) {
	s.backing[key] = value
}
