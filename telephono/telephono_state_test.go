package telephono

import "testing"

func TestExpander_Expand(t *testing.T) {
	under := Expander{}
	t.Run("Adding Contributors", func(s *testing.T) {
		under.AddContributor(&EnvironmentContributor{})
	})
}