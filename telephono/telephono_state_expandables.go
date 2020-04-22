package telephono

////
// Basic Expandable
////
type BasicExpandable struct {
	backing string
}

func (basic *BasicExpandable) Expand(expandable Expander) (string, error) {
	return expandable.Expand(basic.GetUnexpanded())
}

func (basic *BasicExpandable) GetUnexpanded() string {
	return basic.backing
}

func (basic *BasicExpandable) SetUnexpanded(toSet string) {
	basic.backing = toSet
}

func NewExpandable(unexpanded string) BasicExpandable {
	return BasicExpandable{unexpanded}
}
