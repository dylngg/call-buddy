package telephono

type unwrappable interface {
	error
	Unwrap() error
}
