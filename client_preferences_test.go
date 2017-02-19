package irc

import (
	"testing"
)

func TestConfiguration_ReadFile(t *testing.T) {
	var err error
	var p = NewPreferences()
	if err = p.ReadFile("./testdata/valid-client-preferences.json"); err != nil {
		t.Error(err)
	}
	if p.JsonSchema != JsonSchemaUrl {
		t.Errorf(`unexpected JSON schema URI encountered: "%s"`, p.JsonSchema)
	}
}
