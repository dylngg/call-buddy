package telephono_test

import (
	"encoding/json"
	. "github.com/call-buddy/call-buddy/telephono"
	. "github.com/onsi/gomega"
	"strings"
	"testing"
)

func TestSerializeState(t *testing.T) {
	// Environment 1 setup
	env1 := CallBuddyEnvironment{NewSimpleContributor("Var")}
	env1.StoredVariables.Set("Environment", "Env1")

	// Environment 2 setup
	env2 := CallBuddyEnvironment{NewSimpleContributor("Var")}
	env2.StoredVariables.Set("Environment", "Env2")

	// Set up a collection
	headers := NewHeadersTemplate()
	headers.Set("BigBad", "Wolf")
	coll1 := CallBuddyCollection{
		RequestTemplates: []RequestTemplate{
			{
				Method:         Post,
				Url:            NewExpandable("https://google.com/"),
				Headers:        headers,
				ExpandableBody: NewExpandable("{{Var.Environment}}"),
			},
		}}

	state := CallBuddyState{
		Collections:  []CallBuddyCollection{coll1},
		Environments: []CallBuddyEnvironment{env1, env2},
	}

	var marshaledString string
	t.Run("Marshal State", func(sub *testing.T) {
		if marshaledBytes, marshalErr := json.MarshalIndent(state, "|", "  "); marshalErr == nil {
			marshaledString = string(marshaledBytes)
			sub.Log("Successfully marshaled the big fat state")
		}
	})

	t.Run("Check Marshaled string", func(sub *testing.T) {
		RegisterFailHandler(func(message string, callerSkip ...int) {
			sub.Error(message)
		})
		if !(Expect(marshaledString).Should(ContainSubstring("google.com")) &&
			Expect(marshaledString).Should(ContainSubstring("Env1")) &&
			Expect(marshaledString).Should(ContainSubstring("Env2")) &&
			Expect(marshaledString).Should(ContainSubstring("BigBad\"")) &&
			Expect(marshaledString).Should(ContainSubstring("Wolf\"")) &&
			Expect(marshaledString).Should(ContainSubstring("{{Var.Environment}}"))) {
			sub.Fatal("Didn't pass all assertions")
		}
	})
}

func TestExpander_Expand(t *testing.T) {
	under := Expander{}
	var alex = NewSimpleContributor("Alex")
	var cooper = NewSimpleContributor("Cooper")

	t.Run("Adding Contributors", func(s *testing.T) {
		under.AddContributor(EnvironmentContributor{})
		under.AddContributor(alex)
		under.AddContributor(cooper)
	})

	t.Run("Resolve environment variable $PATH", func(sub *testing.T) {
		//language=Mustache
		rendered, renderErr := under.Expand("{{Env.PATH}}")

		// Check if we errored...
		if renderErr != nil {
			sub.Error("Couldn't render environment path: " + renderErr.Error())
		}

		// Check if we rendered the PATH correctly
		if len(rendered) < 2 {
			sub.Error("Rendered string is too small (Should be at least two characters)")
		}

		sub.Log("Rendered $PATH too: " + rendered)
	})

	t.Run("Using Simple Contributor", func(sub *testing.T) {
		alex.Set("iscool", "yes")
		cooper.Set("iscool", "no")

		rendered, renderErr := under.Expand(`
{{#Alex}}
{{iscool}}
{{/Alex}}
{{#Cooper}}
{{iscool}}
{{/Cooper}}`[1:])

		// Check if we errored...
		if renderErr != nil {
			sub.Error("Couldn't render environment path: " + renderErr.Error())
		}

		if !(strings.Contains(rendered, "yes") && strings.Contains(rendered, "no")) {
			sub.Error("Should have contained both yes and no:\n", rendered)
		}

		sub.Log("Rendered:\n", rendered)
	})
}
