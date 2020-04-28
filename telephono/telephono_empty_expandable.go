package telephono

type emptyExpandable struct{}

func (e emptyExpandable) GetUnexpanded() string {
	return ""
}

func (e emptyExpandable) SetUnexpanded(s string) {}

func (e emptyExpandable) Expand(expandable Expander) (string, error) {
	return "", nil
}

// EmptyExpandable always returns an empty string
var EmptyExpandable = emptyExpandable{}
