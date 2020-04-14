package telephono

import mustache "github.com/cbroglie/mustache"

type couldntResolve struct {
	name string
}

func (c couldntResolve) Error() string {
	return "couldn't resolve variable: " + c.name
}

func newCouldntResolve(name string) couldntResolve {
	return couldntResolve{name}
}

/*
 * Expander will take an expandable and push it
 */
type Expander struct {
	contributors    []ContextContributor
	leaveUnresolved bool
}

//Expand the content and return the expanded string or an error if it failed
func (e Expander) Expand(content string) (string, error) {
	var compiled *mustache.Template
	var templateErr error
	if compiled, templateErr = mustache.ParseString(content); templateErr != nil {
		return "", templateErr
	}

	contexts := make(map[string]interface{})
	noNameContexts := make([]interface{}, 0)

	for _, contributor := range e.contributors {
		if name, ctx, err := contributor.Contribute(); err == nil {
			if name != "" {
				contexts[name] = ctx
			} else {
				noNameContexts = append(noNameContexts, ctx)
			}
		}
	}

	var rendered string
	var renderErr error
	if rendered, renderErr = compiled.Render(contexts); renderErr != nil {
		return "", renderErr
	}

	return rendered, nil
}

func (e *Expander) AddContributor(contributor ContextContributor) {
	if e.contributors == nil {
		e.contributors = make([]ContextContributor, 0, 3)
	}

	e.contributors = append(e.contributors, contributor)
}
