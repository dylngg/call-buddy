package telephono

import "testing"

func TestExpander_Expand(t *testing.T) {
	under := Expander{}
	t.Run("Adding Contributors", func(s *testing.T) {
		under.AddContributor(EnvironmentContributor{})
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
}
