package telephono

type HttpMethod string
const (
	Post HttpMethod = "POST"
	Get = "GET"
	Put = "PUT"
	Delete = "DELETE"
	Head = "HEAD"
)

type Expandable interface {
	//GetUnexpanded gives the string as it is now
	GetUnexpanded() string
	//SetUnexpanded will set the unexpanded string
	SetUnexpanded(string)

	//Expand takes the expander and will return the expanded string
	Expand(expandable *Expander) (string, error)
}

type CallTemplate struct {

}


