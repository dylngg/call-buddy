package telephono

import (
	"strings"
	"testing"
)

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
