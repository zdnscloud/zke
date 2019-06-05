package helper

import (
	//"fmt"
	ut "github.com/zdnscloud/cement/unittest"
	"testing"
)

func TestYamlDocParse(t *testing.T) {
	cases := []struct {
		yaml         string
		expectedDocs []string
	}{
		{
			`
---
---
good
---
apiVersion: v1
data:
  notary-signer-ca.crt: |
    -----BEGIN CERTIFICATE-----
    MIIDAzCCAeugAwIBAgIRAO5ZGsMfcfIrCCx93tA0QdwwDQYJKoZIhvcNAQELBQAw
---`,
			[]string{"good",
				`apiVersion: v1
data:
  notary-signer-ca.crt: |
    -----BEGIN CERTIFICATE-----
    MIIDAzCCAeugAwIBAgIRAO5ZGsMfcfIrCCx93tA0QdwwDQYJKoZIhvcNAQELBQAw`},
		},
	}

	for _, tc := range cases {
		var docs []string
		mapOnYamlDocument(tc.yaml, func(doc []byte) error {
			docs = append(docs, string(doc))
			return nil
		})

		ut.Equal(t, len(docs), len(tc.expectedDocs))
		for i, doc := range docs {
			ut.Equal(t, doc, tc.expectedDocs[i])
		}
	}
}
