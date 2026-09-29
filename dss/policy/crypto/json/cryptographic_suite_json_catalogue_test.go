package cryptojson

import (
	"strings"
	"testing"
	"time"
)

func catalogueFor(t *testing.T, document string) jsonObject {
	t.Helper()
	obj, err := parseJSONObject(strings.NewReader(document))
	if err != nil {
		t.Fatal(err)
	}
	policy := obj.asObject(jsonConstraintSecuritySuitabilityPolicy)
	if policy == nil {
		t.Fatal("no SecuritySuitabilityPolicy")
	}
	return policy
}

// A malformed validity date makes upstream's buildAlgorithm skip the whole
// algorithm entry (RFC3339DateUtils#getDate throws, buildAlgorithm catches).
// It must not be read as an absent date, which would keep the algorithm with an
// open-ended validity window.
func TestCatalogueSkipsAlgorithmWithMalformedValidityDate(t *testing.T) {
	const document = `{"SecuritySuitabilityPolicy": {"Algorithm": [
		{"AlgorithmIdentifier": {"Name": "GOOD", "ObjectIdentifier": "1.1"},
		 "Evaluation": [{"Validity": {"Start": "2020-01-01", "End": "2030-12-31"}}]},
		{"AlgorithmIdentifier": {"Name": "BAD-END", "ObjectIdentifier": "1.2"},
		 "Evaluation": [{"Validity": {"Start": "2020-01-01", "End": "not-a-date"}}]},
		{"AlgorithmIdentifier": {"Name": "BAD-START", "ObjectIdentifier": "1.3"},
		 "Evaluation": [{"Validity": {"Start": "01/01/2020"}}]},
		{"AlgorithmIdentifier": {"Name": "NO-VALIDITY", "ObjectIdentifier": "1.4"},
		 "Evaluation": [{}]}
	]}}`
	algorithms := newCryptographicSuiteJsonCatalogue(catalogueFor(t, document)).AlgorithmList()

	var names []string
	for _, algorithm := range algorithms {
		names = append(names, algorithm.AlgorithmIdentifierName())
	}
	if got, want := strings.Join(names, ","), "GOOD,NO-VALIDITY"; got != want {
		t.Fatalf("algorithms = %s, want %s", got, want)
	}
	end := algorithms[0].EvaluationList()[0].ValidityEnd()
	if end == nil || !end.Equal(time.Date(2030, 12, 31, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("GOOD validity end = %v, want 2030-12-31", end)
	}
}

// buildMetadata has no try/catch upstream: a malformed policy date leaves
// getMetadata() as an IllegalArgumentException, ported as a panic.
func TestCatalogueMetadataPanicsOnMalformedDate(t *testing.T) {
	for _, field := range []string{jsonConstraintPolicyIssueDate, jsonConstraintNextUpdate} {
		document := `{"SecuritySuitabilityPolicy": {"` + field + `": "yesterday"}}`
		catalogue := newCryptographicSuiteJsonCatalogue(catalogueFor(t, document))
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("Metadata() with a malformed %s did not panic", field)
				}
			}()
			catalogue.Metadata()
		}()
	}

	const wellFormed = `{"SecuritySuitabilityPolicy": {"PolicyIssueDate": "2024-10-13T00:00:00.000+00:00"}}`
	metadata := newCryptographicSuiteJsonCatalogue(catalogueFor(t, wellFormed)).Metadata()
	if metadata.PolicyIssueDate() == nil || metadata.NextUpdate() != nil {
		t.Errorf("well-formed metadata dates = %v / %v", metadata.PolicyIssueDate(), metadata.NextUpdate())
	}
}
