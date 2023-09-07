package secrets

import (
	"encoding/json"
	"io/ioutil"
)

// Secrets holds the secret fetcher command configuration
type Secrets struct {
	Secrets []Secret `json:"secrets"`
}

// Secret is a tuple containing a name and a value
type Secret struct {
	Key   string `json:"name"`
	Value string `json:"value"`
}

func NewFromPath(path string) (*Secrets, error) {

	c, err := readSecretsFile(path)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// Merge takes the values in the input Secrets and adds them to s. Values in s will be overridden if present.
func (s *Secrets) Merge(override *Secrets) {

	if override == nil {
		return
	}

	// Override secrets in s.Secrets with secrets in override.Secrets, if present
	for _, a := range s.Secrets {
		for _, b := range override.Secrets {
			if a.Key == b.Key && b.Value != "" {
				a.Value = b.Value
			}
		}
	}

	// Add additional secrets from override
	for _, o := range override.Secrets {
		if !isMember(s.Secrets, o) {
			s.Secrets = append(s.Secrets, o)
		}
	}

}

func isMember(ss []Secret, m Secret) bool {
	for _, s := range ss {
		if s.Key == m.Key {
			return true
		}
	}
	return false
}

func readSecretsFile(f string) (*Secrets, error) {

	bs, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}

	ss := &Secrets{}
	err = json.Unmarshal(bs, ss)

	return ss, err
}
