package irc

import (
	"encoding/json"
	"io/ioutil"
)

const (
	// The current revision of the JSON schema that shall be used when performing validation
	// of the preferences file.
	JsonSchemaUrl = "https://headcr4sh.github.io/irc/schema/client-preferences.schema.json#"
)

// ClientPreferences is a structure that can be used to store user-defined data used to customize
// the behavior of the IRC client.
type ClientPreferences struct {
	JsonSchema string             `json:"$schema"`
	Networks   map[string]Network `json:"networks"`
}

type Network struct {
	Servers []Server `json:"servers"`
}

type Server struct {
	Hostname string `json:"hostname"`
	Port     uint   `json:"port"`
}

// NewPreferences creates a new (and empty) structure for holding preferences.
func NewPreferences() *ClientPreferences {
	return &ClientPreferences{
		JsonSchema: JsonSchemaUrl,
	}
}

// ReadFile updates the preferences structure with data stored in the file with the given name.
func (p *ClientPreferences) ReadFile(filename string) (err error) {
	var raw []byte
	if raw, err = ioutil.ReadFile(filename); err == nil {
		err = json.Unmarshal(raw, p)
	}
	return
}

// WriteFile writes the preferences structure as JSON to the file with the given name.
func (p *ClientPreferences) WriteFile(filename string) (err error) {
	var raw []byte
	if raw, err = json.Marshal(p); err == nil {
		ioutil.WriteFile(filename, raw, 0222)
	}
	return
}
