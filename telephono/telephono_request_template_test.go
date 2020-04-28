package telephono_test

import (
	. "github.com/call-buddy/call-buddy/telephono"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestBasicCallTemplate(t *testing.T) {
	setUpServer()
	template := NewHeadersTemplate()
	template.Set("User", "{{Var.A}}")

	//Expander and contributor code
	expander := Expander{}
	expander.AddContributor(EnvironmentContributor{})
	contributor := NewSimpleContributor("Var")
	//contributor.Set("host", "https://httpbin.org")
	contributor.Set("host", GlobalTestState.getPrefix())
	contributor.Set("A", "AAAA")
	contributor.Set("B", "BBBB")
	contributor.Set("Status", "329")
	expander.AddContributor(contributor)

	// Execute call
	templateUnderTest := RequestTemplate{
		Method:  Post,
		Url:     NewExpandable("{{Var.host}}/postbalogna"),
		Headers: template,
		ExpandableBody: NewExpandable(
			`{
	"path": "{{Var.B}}"
}`),
	}

	response, requestErr := templateUnderTest.ExecuteWithClientAndExpander(http.DefaultClient, expander)

	if requestErr != nil {
		t.Fatal("Got an error!\n", requestErr.Error())
	}

	t.Run("Check Status Code", func(sub *testing.T) {
		if response.StatusCode != 200 {
			t.Fatal("Got an error!", response.StatusCode)
		}
	})

	if allBytes, readErr := ioutil.ReadAll(response.Body); readErr == nil {
		asString := string(allBytes)
		t.Log("Received body below")
		t.Log(prefixAllLinesOfString(asString, '|'))
		if !(strings.Contains(asString, "BBBB") && strings.Contains(asString, "AAAA") && strings.Contains(asString, "postbalogna")) {
			t.Fatal("Didn't find AAAA and BBBB")
		}
	} else {
		t.Fatal(readErr.Error())
	}
}
