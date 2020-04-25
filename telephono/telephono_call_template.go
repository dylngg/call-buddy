package telephono

import (
	"net/http"
	"strings"
)

type Body BasicExpandable

type CallResponse *http.Response

type RequestTemplate struct {
	Method         HttpMethod
	Url            Expandable
	Headers        HeadersTemplate
	ExpandableBody Expandable
	// TODO AH: specify a body type that's just given a reader.
}

//executeWithClientAndExpander will execute this call template with the specified client and expander, returning a response or an error
func (r *RequestTemplate) ExecuteWithClientAndExpander(client *http.Client, expander Expander) (CallResponse, error) {
	//expand the url
	expandedUrl, urlErr := r.Url.Expand(expander)
	if urlErr != nil {
		return nil, urlErr
	}

	//expand the body
	//TODO AH: file bodies for things like binary data or purposefully unrendered stuff
	//OPTIMIZE AH: Instead of just expanding this, stream it so that we're not loading so many things into memory
	expandedBody, bodyErr := r.ExpandableBody.Expand(expander)
	if bodyErr != nil {
		return nil, bodyErr
	}

	bodyReader := strings.NewReader(expandedBody)

	toExecute, newCallErr := http.NewRequestWithContext(globalState.callContext, r.Method.asMethodString(), expandedUrl, bodyReader)
	if newCallErr != nil {
		return nil, newCallErr
	}

	// add the headers
	if header, errors := r.Headers.ExpandAllAsHeader(expander); len(errors) == 0 {
		toExecute.Header = header
	}

	response, err := client.Do(toExecute)

	return response, err
}
