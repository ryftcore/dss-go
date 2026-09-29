package xades

import (
	"sync"
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
)

// TestConcurrentValidationOfIndependentDocuments validates the same fixture from several
// goroutines, each with its own validator, the way a server validates independent uploads. The
// process-global state the XAdES layer touches on every signature - the XPath namespace registry
// (DSSXMLUtilsRegisterXAdESNamespaces), the executor loader, the structure-validator singleton
// and the recovery registries - must be safe for that. Meaningful under go test -race.
func TestConcurrentValidationOfIndependentDocuments(t *testing.T) {
	path := xadesFixturePath(t, "upstream/Signature-X-AT-1.xml")

	const workers = 6
	var wg sync.WaitGroup
	errs := make(chan error, workers)
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			doc, err := model.NewFileDocument(path)
			if err != nil {
				errs <- err
				return
			}
			validator, err := NewXMLDocumentValidator(doc)
			if err != nil {
				errs <- err
				return
			}
			validator.SetCertificateVerifier(permissiveCertificateVerifier())
			validator.SetValidationLevel(enumerations.ValidationLevelBasicSignatures)
			reports, err := validator.ValidateDocument()
			if err == nil && reports.GetSimpleReport().GetSignaturesCount() != 1 {
				t.Errorf("GetSignaturesCount() = %d, want 1", reports.GetSimpleReport().GetSignaturesCount())
			}
			errs <- err
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Errorf("concurrent validation: %v", err)
		}
	}
}
