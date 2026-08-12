// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/policy/SignaturePolicyProvider.java (DSS 6.5.RC1).
package policy

import (
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi/client/http"
	"github.com/utain/esig/dss/utils"
)

// SignaturePolicyProvider retrieves a policy by its
// SignaturePolicyIdentifier.
type SignaturePolicyProvider struct {
	// dataLoader is used to obtain a policy from a source (e.g. online).
	dataLoader http.DataLoader

	// signaturePoliciesByID is the map of signature policy documents by
	// IDs.
	signaturePoliciesByID map[string]model.DSSDocument

	// signaturePoliciesByURL is the map of signature policy documents by
	// URLs.
	signaturePoliciesByURL map[string]model.DSSDocument
}

// NewSignaturePolicyProvider is the default constructor, instantiating the
// object with a nil data loader and empty maps.
func NewSignaturePolicyProvider() *SignaturePolicyProvider {
	return &SignaturePolicyProvider{
		signaturePoliciesByID:  make(map[string]model.DSSDocument),
		signaturePoliciesByURL: make(map[string]model.DSSDocument),
	}
}

// SetDataLoader sets the DataLoader used to retrieve signature policy
// documents (e.g. from online).
func (p *SignaturePolicyProvider) SetDataLoader(dataLoader http.DataLoader) {
	p.dataLoader = dataLoader
}

// SetSignaturePoliciesByID sets the map of signature policy documents to
// retrieve by IDs.
func (p *SignaturePolicyProvider) SetSignaturePoliciesByID(signaturePoliciesByID map[string]model.DSSDocument) {
	p.signaturePoliciesByID = signaturePoliciesByID
}

// SetSignaturePoliciesByURL sets the map of signature policy documents to
// retrieve by URLs.
func (p *SignaturePolicyProvider) SetSignaturePoliciesByURL(signaturePoliciesByURL map[string]model.DSSDocument) {
	p.signaturePoliciesByURL = signaturePoliciesByURL
}

// GetSignaturePolicyByID gets a signature policy document with the
// corresponding policyID from signaturePoliciesByID. Returns nil if not
// found.
func (p *SignaturePolicyProvider) GetSignaturePolicyByID(policyID string) model.DSSDocument {
	return p.signaturePoliciesByID[policyID]
}

// GetSignaturePolicyByURL gets a signature policy document with the
// corresponding url from signaturePoliciesByURL; if not found, retrieves the
// data from url with the DataLoader. Returns nil if not found.
func (p *SignaturePolicyProvider) GetSignaturePolicyByURL(url string) model.DSSDocument {
	dssDocument := p.signaturePoliciesByURL[url]
	if dssDocument == nil && utils.IsStringNotBlank(url) && p.dataLoader != nil {
		bytes := signaturePolicyProviderTryGet(p.dataLoader, url)
		if utils.IsArrayEmpty(bytes) {
			// LOG.warn("Empty content for url '{}'", url) dropped: not
			// load-bearing per PORTING.md.
			return nil
		}
		dssDocument = model.NewInMemoryDocument(bytes)
	}
	return dssDocument
}

// signaturePolicyProviderTryGet calls dataLoader.Get(url), recovering from
// any panic it raises (e.g. *exception.DSSExternalResourceException, as
// NativeHTTPDataLoader.Get panics on failure) and returning nil in that
// case. Ports the try/catch(Exception) guarding dataLoader.get(url).
func signaturePolicyProviderTryGet(dataLoader http.DataLoader, url string) (data []byte) {
	defer func() {
		if recover() != nil {
			data = nil
		}
	}()
	return dataLoader.Get(url)
}

// GetSignaturePolicy gets a signature policy by all available ways (id and
// uri). Returns nil if not found.
func (p *SignaturePolicyProvider) GetSignaturePolicy(policyID, url string) model.DSSDocument {
	dssDocument := p.GetSignaturePolicyByID(policyID)
	if dssDocument == nil {
		dssDocument = p.GetSignaturePolicyByURL(url)
	}
	return dssDocument
}
