package telephono


////
// Basic Expandable
////
type BasicExpandable struct {
	backing string
}

func (basic *BasicExpandable) GetUnexpanded() string {
	return basic.backing
}

func (basic *BasicExpandable) SetUnexpanded(toSet string) {
	basic.backing = toSet
}

func (basic *BasicExpandable) Expand(expander Expander) (string, error) {
	return expander.Expand(basic.GetUnexpanded())
}
