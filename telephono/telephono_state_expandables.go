package telephono

////
// Basic Expandable
////
type BasicExpandable struct {
	Backing string
}

//func (basic *BasicExpandable) UnmarshalJSON(bytes []byte) error {
//	byteString := string(bytes)
//	basic.backing = strings.Trim(byteString, `"`)
//	return nil
//}
//
//func (basic *BasicExpandable) MarshalJSON() ([]byte, error) {
//	return []byte(`"` + basic.backing + `"`), nil
//}

func (basic *BasicExpandable) Expand(expandable Expander) (string, error) {
	return expandable.Expand(basic.GetUnexpanded())
}

func (basic *BasicExpandable) GetUnexpanded() string {
	return basic.Backing
}

func (basic *BasicExpandable) SetUnexpanded(toSet string) {
	basic.Backing = toSet
}

func NewExpandable(unexpanded string) *BasicExpandable {
	return &BasicExpandable{unexpanded}
}
