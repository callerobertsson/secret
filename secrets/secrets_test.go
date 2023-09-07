package secrets

import (
	"fmt"
	"testing"
)

func TestMerge(t *testing.T) {

	cases := []struct {
		s Secrets
		m Secrets
		e Secrets
	}{
		{
			s: Secrets{[]Secret{Secret{Key: "k", Value: "v"}}},
			m: Secrets{[]Secret{}},
			e: Secrets{[]Secret{Secret{Key: "k", Value: "v"}}},
		},
		{
			s: Secrets{[]Secret{Secret{Key: "k", Value: "v"}}},
			m: Secrets{[]Secret{Secret{Key: "k", Value: "o"}}},
			e: Secrets{[]Secret{Secret{Key: "k", Value: "x"}}},
		},
	}

	for _, c := range cases {

		c.s.Merge(&c.m)

		if !equalSecrets(c.s, c.e) {
			t.Errorf("Expected %v but got %v\n", c.e, c.s)
		}
	}
}

// func mapToSecrets(m map[string]string) Secrets {
// 	s := Secrets{}
// 	for k, v := range m {
// 		s.Secrets = append(s.Secrets, Secret{Key: k, Value: v})
// 	}
// 	return s
// }

// func newSecrets(ss []Secret) Secrets {
// 	ns := Secrets{Secrets: ss}
// }
func equalSecrets(a, b Secrets) bool {
	// TODO: improve

	as := fmt.Sprintf("%v", a)
	bs := fmt.Sprintf("%v", b)

	return as == bs
}
