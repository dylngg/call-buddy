package telephono

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"testing"
)

func TestBasicCallTemplate(t *testing.T) {
	template := NewHeadersTemplate()
	template.Set("User", "{{Var.A}}")

	//Expander and contributor code
	expander := Expander{}
	expander.AddContributor(EnvironmentContributor{})
	contributor := NewSimpleContributor("Var")
	contributor.Set("host", "https://httpbin.org")
	contributor.Set("A", "AAAA")
	contributor.Set("B", "BBBB")
	contributor.Set("Status", "329")
	expander.AddContributor(contributor)

	// Execute call
	templateUnderTest := RequestTemplate{
		method:  Post,
		url:     NewExpandable("{{Var.host}}/post"),
		headers: template,
		expandableBody: NewExpandable(
			`{
	"path": "{{Var.B}}"
}`),
	}

	response, requestErr := templateUnderTest.executeWithClientAndExpander(http.DefaultClient, expander)

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
		log.Println(asString)
		if !(strings.Contains(asString, "BBBB") && strings.Contains(asString, "AAAA")) {
			t.Fatal("Didn't find AAAA and BBBB")
		}
	} else {
		t.Fatal(readErr.Error())
	}
}
