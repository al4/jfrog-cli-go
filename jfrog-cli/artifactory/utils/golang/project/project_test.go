package project

import "testing"

func TestParseModuleName(t *testing.T) {
	content := `

         	module github.com/jfrog/vgo-example

        require rsc.io/quote v1.5.2
	`

	expected := "github.com/jfrog/vgo-example"
	actual, err := parseModuleName(content)

	if err != nil {
		t.Error(err)
	}

	if actual != expected {
		t.Errorf("Expected: %s, got: %s.", expected, actual)
	}
}
