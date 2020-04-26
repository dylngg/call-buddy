package telephono

import (
	"context"
	"net/http"
)

//here is a test comment to make sure i am editing and commiting correctly - coop diddy

type HttpMethod string

const (
	Post   HttpMethod = "POST"
	Get               = "GET"
	Put               = "PUT"
	Delete            = "DELETE"
	Head              = "HEAD"
)

func (m HttpMethod) asMethodString() string {
	return string(m)
}

type expandable interface {
	//GetUnexpanded gives the string as it is now
	GetUnexpanded() string
	//SetUnexpanded will set the unexpanded string
	SetUnexpanded(string)

	//Expand takes the expander and will return the expanded string
	Expand(expandable Expander) (string, error)
}

/*CallBuddyState is the full shippable state of call buddy
environments, call templates, possibly history, variables, etc. are all in here
It can be shipped to remote servers to be run
*/
type CallBuddyState struct {
	Collections  []CallBuddyCollection
	Environments []CallBuddyEnvironment
}

type CallBuddyCollection struct {
	// TODO AH: Call templates
	RequestTemplates []RequestTemplate
}

type CallBuddyEnvironment struct {
	StoredVariables SimpleContributor
}

type CallBuddyInternalState struct {
	client      *http.Client
	callContext context.Context

	freeFunc context.CancelFunc
}

var globalState CallBuddyInternalState

func init() {

	// TODO AH: Make this timeout longer and configurable. Maybe have a check for number of received bytes on each call
	//timeoutContext, cancelFunc := context.@WithTimeout(context.Background(), time.Minute*3)
	// goddamn I love garbage collection
	globalState.callContext = context.Background()
	globalState.freeFunc = func() {}

	// create the client
	globalState.client = &http.Client{
		Transport:     http.DefaultTransport,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       0,
	}
}
