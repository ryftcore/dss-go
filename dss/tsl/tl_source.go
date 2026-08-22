// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/source/TLSource.java (DSS 6.5.RC1).
//
// CROSS-CHUNK DEPENDENCIES (see this batch's porter notes):
//
//  1. This file embeds eu.europa.esig.dss.validation.job.source.DocumentSource, ported by the
//     dss-validation-job chunk into Go package dss/validation/job, which had not landed when this
//     file was written. The embedded name and its constructor follow PORTING.md's "exported Go
//     identifiers keep the Java name" rule literally (DocumentSource / NewDocumentSource).
//
//  2. Java types the three predicate fields as the raw java.util.function.Predicate<T>. Go has no
//     generic functional interface with a canonical name, and every predicate the ported tree
//     actually assigns here is one of the NAMED interfaces the dss-tsl-validation "function"
//     package declares for exactly these three T's - TrustServiceProviderPredicate
//     (Predicate<TSPType>), TrustServicePredicate (Predicate<TSPServiceType>) and
//     TrustAnchorPeriodPredicate (Predicate<TrustServiceStatusAndInformationExtensions>), all
//     produced by TLPredicateFactory. The fields therefore carry those named interface types
//     (same Go package tsl, ported by the TSLJOB chunk); an arbitrary caller lambda is wrapped in
//     an implementation of the corresponding interface rather than assigned directly, which is
//     the only expressible difference from Java.
package tsl

import "github.com/ryftcore/dss-go/dss/validation/job"

// tlSourceDefaultSupportedTLVersions is the list of default supported TL versions. Port of the
// private static DEFAULT_SUPPORTED_TL_VERSIONS constant; it is copied into every new TLSource so
// that a caller mutating the returned slice cannot corrupt the shared default (Java hands out
// the same immutable Arrays.asList to every instance).
var tlSourceDefaultSupportedTLVersions = []int{5, 6}

// TLSource represents a Trusted List source.
type TLSource struct {
	*job.DocumentSource

	// trustServiceProviderPredicate allows filtering the collected trust service provider(s)
	// with a predicate.
	//
	// Default: all trust service providers are selected.
	trustServiceProviderPredicate TrustServiceProviderPredicate

	// trustServicePredicate allows filtering the collected trust service(s) with a predicate.
	//
	// Default: all trust services are selected.
	trustServicePredicate TrustServicePredicate

	// trustAnchorValidityPredicate defines whether an SDI can be considered as a trust anchor
	// during the given period of time.
	trustAnchorValidityPredicate TrustAnchorPeriodPredicate

	// tlVersions is the list of TL versions accepted for the current TLSource. When defined,
	// an error is returned on structure validation.
	tlVersions []int
}

// NewTLSource is the default constructor, instantiating the object with null values. Port of
// TLSource().
func NewTLSource() *TLSource {
	source := &TLSource{DocumentSource: job.NewDocumentSource()}
	source.tlVersions = append([]int(nil), tlSourceDefaultSupportedTLVersions...)
	return source
}

// TrustServiceProviderPredicate gets the predicate used to filter TrustServiceProviders. Port of
// getTrustServiceProviderPredicate().
func (s *TLSource) TrustServiceProviderPredicate() TrustServiceProviderPredicate {
	return s.trustServiceProviderPredicate
}

// SetTrustServiceProviderPredicate sets the predicate used to filter TrustServiceProviders. Port
// of setTrustServiceProviderPredicate(Predicate).
func (s *TLSource) SetTrustServiceProviderPredicate(trustServiceProviderPredicate TrustServiceProviderPredicate) {
	s.trustServiceProviderPredicate = trustServiceProviderPredicate
}

// TrustServicePredicate gets the predicate used to filter TrustServices. Port of
// getTrustServicePredicate().
func (s *TLSource) TrustServicePredicate() TrustServicePredicate {
	return s.trustServicePredicate
}

// SetTrustServicePredicate sets the predicate used to filter TrustServices. Port of
// setTrustServicePredicate(Predicate).
func (s *TLSource) SetTrustServicePredicate(trustServicePredicate TrustServicePredicate) {
	s.trustServicePredicate = trustServicePredicate
}

// TrustAnchorValidityPredicate gets the predicate filtering
// TrustServiceStatusAndInformationExtensions in order to define an acceptability period of a
// corresponding SDI as a trust anchor. Port of getTrustAnchorValidityPredicate().
func (s *TLSource) TrustAnchorValidityPredicate() TrustAnchorPeriodPredicate {
	return s.trustAnchorValidityPredicate
}

// SetTrustAnchorValidityPredicate sets a predicate allowing to filter
// TrustServiceStatusAndInformationExtensions in order to define an acceptability period of a
// corresponding SDI as a trust anchor. If the predicate is defined and the condition fails, the
// SDI will not be treated as a trust anchor during the validation process. Port of
// setTrustAnchorValidityPredicate(Predicate).
func (s *TLSource) SetTrustAnchorValidityPredicate(trustAnchorValidityPredicate TrustAnchorPeriodPredicate) {
	s.trustAnchorValidityPredicate = trustAnchorValidityPredicate
}

// TLVersions gets the list of TL versions to be accepted for the current TL/LOTL source. Port of
// getTLVersions().
func (s *TLSource) TLVersions() []int {
	return s.tlVersions
}

// SetTLVersions sets the list of acceptable XML Trusted List versions. When defined, an error
// message is returned on structural validation; if not defined, no structural validation is
// performed.
//
// Default: 5, 6 (structure validation is performed for versions 5 and 6, while other versions are
// invalidated). Port of setTLVersions(List).
func (s *TLSource) SetTLVersions(tlVersions []int) {
	s.tlVersions = tlVersions
}
