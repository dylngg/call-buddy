package telephono

import "errors"

/*
 Simple contributor key/value store

May be replaced with a better struct more suited to serialization
*/
type SimpleContributor struct {
	Prefix  string
	Backing map[string]string
}

func (s SimpleContributor) Contribute() (string, interface{}, error) {
	if len(s.Prefix) == 0 {
		return "", 0, errors.New("Must have a prefix")
	}
	return s.Prefix, s.Backing, nil
}

func (s *SimpleContributor) Set(key, value string) {
	s.Backing[key] = value
}

//NewSimpleContributor creates a new simple contribut at the `prefix` address, i.e. accessible at {{prefix.Value}}
func NewSimpleContributor(prefix string) SimpleContributor {
	if len(prefix) == 0 {
		// TODO AH: Should we panic?
		panic("Must have a prefix")
	}

	return SimpleContributor{
		Prefix:  prefix,
		Backing: map[string]string{},
	}
}
