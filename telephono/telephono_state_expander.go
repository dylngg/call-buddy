package telephono

import (
	"bytes"
	"io"
	"log"
	"sort"
)

type couldntResolve struct {
	name string
}

func (c couldntResolve) Error() string {
	return "couldn't resolve variable: " + c.name
}

func newCouldntResolve(name string) couldntResolve {
return couldntResolve{name}
}

type Expander struct {
	resolvers OrderedVariableResolver
}


func (e Expander) resolve(name string) (string, error) {
	for _, resolver := range e.resolvers {
		if resolver.CanHandle(name) {
			if resolved, resolverErr := resolver.Resolve(name); resolverErr != nil {
				return resolved, nil
			}
		}
	}

	return "", newCouldntResolve(name)
}

//Expand the content and return the expanded string or an error if it failed
func (e Expander) Expand(content string) (string, error) {
	type expandState struct {
		readingVarName bool
		varName string
		chunks []Expandable
	}
	expanded := bytes.Buffer{}
	bufferedCurrent := bytes.NewBufferString(content)

	currentState := expandState{false, "", []Expandable{}}

	for true {
		this, _, err := bufferedCurrent.ReadRune()
		if err == io.EOF {
			break
		} else if err != nil {
			// If we get some other kind of error log it and move on
			log.Fatal(err.Error())
		}

		if currentState.readingVarName {
			if this == '{' {
				continue
			} else if this == '}'{
				currentState.readingVarName = false
				if resolved, resolverError := e.resolve(currentState.varName); resolverError == nil {
					expanded.WriteString(resolved)
				} else {
					// TODO AH: Logging
					continue
				}
			} else {
				currentState.varName = currentState.varName + string(this)
			}

		} else {
			if this == '$' {
				currentState.readingVarName = true
				continue
			} else {

			}
		}
	}

	return expanded.String(), nil
}

func (e *Expander) AddResolver(resolver VariableResolver) {
	e.resolvers = append(e.resolvers, resolver)
	sort.Sort(e.resolvers)
}
