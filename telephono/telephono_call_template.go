package telephono

type Body BasicExpandable


type RequestTemplate struct {
	method HttpMethod
	headers HeadersTemplate
	body BasicExpandable
}

func (r *RequestTemplate) Execute() error {
	return nil
}
