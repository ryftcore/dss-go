// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/aia/CompositeAIASource.java (DSS 6.5.RC1).
package aia

import (
	"fmt"
	"sort"

	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/exception"
)

// CompositeSource allows retrieving an AIA with different sources. The composite tries all
// sources until it gets a non-empty response.
type CompositeSource struct {
	// aiaSources is a map of source keys and corresponding AIA sources.
	aiaSources map[string]Source
}

// NewCompositeSource is the default constructor, instantiating the object with a nil map.
func NewCompositeSource() *CompositeSource {
	return &CompositeSource{}
}

// SetAIASources allows providing multiple AIA sources. Be careful: all given sources MUST
// accept the same digest algorithm.
func (c *CompositeSource) SetAIASources(aiaSources map[string]Source) {
	c.aiaSources = aiaSources
}

// CertificatesByAIA tries every configured Source in turn, returning the first non-nil
// response. Panics with a *exception.DSSExternalResourceException when none of the sources
// yields one, mirroring Java's unchecked DSSExternalResourceException propagating out of a
// method with no throws clause.
func (c *CompositeSource) CertificatesByAIA(certificateToken *model.CertificateToken) []*model.CertificateToken {
	for _, sourceKey := range compositeSourceOrderedKeys(c.aiaSources) {
		source := c.aiaSources[sourceKey]
		certificateTokens := compositeSourceTryGet(source, certificateToken)
		if certificateTokens != nil {
			return certificateTokens
		}
	}
	panic(exception.NewDSSExternalResourceException(
		fmt.Sprintf("Unable to retrieve the certificateTokens (%d tries)", len(c.aiaSources))))
}

// compositeAIASourceOrderedKeys returns the keys of aiaSources, sorted lexically. Java's
// HashMap iteration order is arbitrary but stable within a JVM run; Go map iteration is
// instead randomized on every run. Sorting (PORTING.md's "order-sensitive upstream iteration
// -> ... explicit sort" rule) keeps the try order - and therefore CertificatesByAIA's result,
// when more than one source would answer - stable from one run to the next, rather than merely
// "immaterial because every source is tried": two sources can both hold a (possibly
// different) valid AIA response for the same certificate.
func compositeSourceOrderedKeys(aiaSources map[string]Source) []string {
	keys := make([]string, 0, len(aiaSources))
	for key := range aiaSources {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// compositeSourceTryGet calls source.CertificatesByAIA, recovering from any panic it raises
// and returning nil in that case. Ports the try/catch(Exception) guarding
// source.getCertificatesByAIA(certificateToken).
func compositeSourceTryGet(source Source, certificateToken *model.CertificateToken) (certificateTokens []*model.CertificateToken) {
	defer func() {
		if recover() != nil {
			certificateTokens = nil
		}
	}()
	return source.CertificatesByAIA(certificateToken)
}

var _ Source = (*CompositeSource)(nil)
