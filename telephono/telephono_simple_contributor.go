package telephono

import "errors"

type SimpleContributor struct {
	prefix  string
	backing map[string]string
}

func (s SimpleContributor) Contribute() (string, interface{}, error) {
	if len(s.prefix) == 0 {
		return "", 0, errors.New("Must have a prefix")
	}
	return s.prefix, s.backing, nil
}

func (s *SimpleContributor) Set(key, value string) error {
	s.backing[key] = value
	return nil
}

func New(prefix string) SimpleContributor {
	if len(prefix) == 0 {
		// TODO AH: Should we panic?
		panic("Must have a prefix")
	}

	return SimpleContributor{
		prefix:  prefix,
		backing: map[string]string{},
	}
}
