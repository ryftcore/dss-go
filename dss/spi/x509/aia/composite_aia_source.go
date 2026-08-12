// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/aia/CompositeAIASource.java (DSS 6.5.RC1).
package aia

import (
	"fmt"

	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi/exception"
)

// CompositeAIASource allows retrieving an AIA with different sources. The composite tries all
// sources until it gets a non-empty response.
type CompositeAIASource struct {
	// aiaSources is a map of source keys and corresponding AIA sources.
	aiaSources map[string]AIASource
}

// NewCompositeAIASource is the default constructor, instantiating the object with a nil map.
func NewCompositeAIASource() *CompositeAIASource {
	return &CompositeAIASource{}
}

// SetAIASources allows providing multiple AIA sources. Be careful: all given sources MUST
// accept the same digest algorithm.
func (c *CompositeAIASource) SetAIASources(aiaSources map[string]AIASource) {
	c.aiaSources = aiaSources
}

// CertificatesByAIA tries every configured AIASource in turn, returning the first non-nil
// response. Panics with a *exception.DSSExternalResourceException when none of the sources
// yields one, mirroring Java's unchecked DSSExternalResourceException propagating out of a
// method with no throws clause.
func (c *CompositeAIASource) CertificatesByAIA(certificateToken *model.CertificateToken) []*model.CertificateToken {
	for _, sourceKey := range compositeAIASourceOrderedKeys(c.aiaSources) {
		source := c.aiaSources[sourceKey]
		certificateTokens := compositeAIASourceTryGet(source, certificateToken)
		if certificateTokens != nil {
			return certificateTokens
		}
	}
	panic(exception.NewDSSExternalResourceException(
		fmt.Sprintf("Unable to retrieve the certificateTokens (%d tries)", len(c.aiaSources))))
}

// compositeAIASourceOrderedKeys returns the keys of aiaSources. Go map iteration order is
// randomized, unlike Java's LinkedHashMap-free HashMap (whose entrySet() order is likewise
// unspecified); the difference is behaviourally immaterial since every source is tried until
// one succeeds.
func compositeAIASourceOrderedKeys(aiaSources map[string]AIASource) []string {
	keys := make([]string, 0, len(aiaSources))
	for key := range aiaSources {
		keys = append(keys, key)
	}
	return keys
}

// compositeAIASourceTryGet calls source.CertificatesByAIA, recovering from any panic it raises
// and returning nil in that case. Ports the try/catch(Exception) guarding
// source.getCertificatesByAIA(certificateToken).
func compositeAIASourceTryGet(source AIASource, certificateToken *model.CertificateToken) (certificateTokens []*model.CertificateToken) {
	defer func() {
		if recover() != nil {
			certificateTokens = nil
		}
	}()
	return source.CertificatesByAIA(certificateToken)
}

var _ AIASource = (*CompositeAIASource)(nil)
