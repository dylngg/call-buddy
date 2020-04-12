package telephono

import (
	"fmt"
	"net/http"
)

type HeadersTemplate struct {
	backing map[string]BasicExpandable
}

type HeaderResolutionErr struct {
	headerKey string
	// TODO AH: The below may not be necessary
	//isKey indicates the resolution error was with the key or the value
	isKey bool

	underlying error
}

func (h HeaderResolutionErr) Error() string {
	var keyOrVal string
	if h.isKey {
		keyOrVal = "key"
	} else {
		keyOrVal = "value"
	}
	return fmt.Sprintf("%-30s unable to resolve %s", h.headerKey, keyOrVal)
}

func (h HeaderResolutionErr) Unwrap() error {
	return h.underlying
}

func (headersTemplate HeadersTemplate) ResolveAllAsHeader(expander Expander) (http.Header, []error) {
	returnErrors := make([]error, 0, len(headersTemplate.backing))
	toReturn := http.Header{}

	// TODO AH: Error handling here
	for k, v := range headersTemplate.backing {
		if valueResolved, valueErr := v.Expand(expander); valueErr == nil {
			// FIXME AH: Multiple headers?
			toReturn.Set(k, valueResolved)
		} else {
			returnErrors = append(returnErrors, HeaderResolutionErr{
				headerKey:  k,
				isKey:      true,
				underlying: valueErr,
			})
		}
	}

	return toReturn, returnErrors
}
