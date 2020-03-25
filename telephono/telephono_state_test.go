package telephono

import (
	"testing"
)

func TestSimpleInMemResolver(t *testing.T) {
	under := newSimpleInMemResolver()
	t.Run("Order", func(t *testing.T) {
		order := under.Order()
		if order >= 0 {
			t.Error("In memory resolver should be below 0")
		}
	})

	t.Run("Resolver", func(t *testing.T) {
		tester := map[string]string{}
		resolutions := []struct {
			key string
			val string
		}{
			{"1", "1-1"},
			{"1", "1-2"},
			{"2", "2-2"},
		}

		for _, resolution := range resolutions {
			tester[resolution.key] = resolution.val
			under.addResolution(resolution.key, resolution.val)
		}

		for k, v := range tester {
			if canDo := under.CanHandle(k); canDo {
				if resolved, resolverErr := under.Resolve(k); resolverErr == nil {
					if resolved != v {
						t.Error("Resolved ", k, " incorrectly; expected: ", v, " got ", resolved)
					}
				} else {
					t.Error("Got error: ", resolverErr.Error())
				}
			} else {
				t.Error("Couldn't resolve: ", k)
			}
		}
	})
}
