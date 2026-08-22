// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/validation/scope/XAdESTimestampScopeFinder.java (DSS 6.5.RC1).
//
// # GAP flagged for integrator: not reachable from SignatureTimestampSource's internal flow
//
// Java's XAdESTimestampSource#getTimestampScopes(TimestampToken) override (a concrete-but-not-
// abstract protected method on the frozen spi/validation/timestamp.SignatureTimestampSource,
// outside this manifest) constructs and uses ITS OWN TimestampScopeFinder instead of the
// base's plain EncapsulatedTimestampScopeFinder. Go's SignatureTimestampSource[AS, SA] (already
// landed, spi/validation/timestamp/signature_timestamp_source.go) has no override hook for this
// method at all - unlike every fully-abstract method (routed through
// SignatureTimestampSourceOverrides) or even the three CAdES-flagged concrete-but-overridable
// hooks (IncorporateArchiveTimestampReferences/GetSignatureSignedDataReferences/
// GetCounterSignatureReferences, see cades_timestamp_source.go's file header for the identical
// pattern), its private getTimestampScopes always constructs a bare
// validationscope.NewEncapsulatedTimestampScopeFinder() and calls s.signature (not through
// s.overrides). This file's TimestampScopeFinder is therefore correct and independently
// usable (xades_timestamp_source.go's GetTimestampScopes method, exported for exactly this
// reason, constructs and uses it directly), but is NOT reached by
// SignatureTimestampSource.validateTimestamps()'s own calls to getTimestampScopes(timestampToken)
// for content/archive timestamps - those still get the base's XAdES-agnostic
// EncapsulatedTimestampScopeFinder (which returns ALL of the signature's scopes rather than
// XAdES's include-filtered subset for an IndividualDataObjectsTimestamp) until
// SignatureTimestampSourceOverrides gains a GetTimestampScopes(*TimestampToken) entry and
// validateTimestamps's call site is changed to go through it. Flagged prominently in the porter
// report; the integrator arbitrates.
//
// Type-asserts the embedded AdvancedSignature to *Signature to call its
// XAdESReferenceValidations() (see xades_signature_scope_finder.go).
package xades

import (
	"github.com/ryftcore/dss-go/dss/model/scope"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	spiscope "github.com/ryftcore/dss-go/dss/spi/validation/scope"
	"github.com/ryftcore/dss-go/dss/utils"
)

// XAdESTimestampScopeFinder finds a timestamp scope for a XAdES encapsulated timestamp. Port of
// the class TimestampScopeFinder, extending spiscope.EncapsulatedTimestampScopeFinder.
type TimestampScopeFinder struct {
	*spiscope.EncapsulatedTimestampScopeFinder
}

// NewXAdESTimestampScopeFinder is the port of the default constructor.
func NewTimestampScopeFinder() *TimestampScopeFinder {
	return &TimestampScopeFinder{EncapsulatedTimestampScopeFinder: spiscope.NewEncapsulatedTimestampScopeFinder()}
}

// FindTimestampScope returns a timestamp scope for the given TimestampToken, dispatching to this
// type's own FilterCoveredSignatureScopes override. Port of findTimestampScope(TimestampToken) as
// Java's virtual dispatch actually resolves it for an XAdESTimestampScopeFinder instance: the
// base's FindTimestampScope (spiscope.EncapsulatedTimestampScopeFinder, embedded by pointer here)
// calls f.FilterCoveredSignatureScopes statically bound to ITS OWN body, so this override
// re-implements the base's short body instead of delegating to it, giving this type's
// FilterCoveredSignatureScopes below the call it needs - the same embedding-has-no-virtual-
// dispatch limitation flagged throughout this port (e.g. xades_reference_validation.go's
// TransformationNames).
func (f *TimestampScopeFinder) FindTimestampScope(timestampToken *validation.TimestampToken) []scope.SignatureScope {
	if timestampToken.IsMessageImprintDataIntact() {
		return f.FilterCoveredSignatureScopes(timestampToken)
	}
	return []scope.SignatureScope{}
}

// FilterCoveredSignatureScopes filters and returns covered SignatureScopes by the current
// timestamp. Port of the protected filterCoveredSignatureScopes(TimestampToken) override.
func (f *TimestampScopeFinder) FilterCoveredSignatureScopes(timestampToken *validation.TimestampToken) []scope.SignatureScope {
	timestampIncludes := timestampToken.TimestampIncludes()
	if utils.IsCollectionNotEmpty(timestampIncludes) {
		individualSignatureScopes := make([]scope.SignatureScope, 0)
		signatureScopes := f.Signature.SignatureScopes()
		if utils.IsCollectionNotEmpty(signatureScopes) {
			if xadesSignature, ok := f.Signature.(*Signature); ok {
				for _, xadesReferenceValidation := range xadesSignature.XAdESReferenceValidations() {
					if xadesTimestampScopeFinderIsContentTimestampedReference(xadesReferenceValidation, timestampIncludes) {
						for _, signatureScope := range signatureScopes {
							if utils.EndsWithIgnoreCase(xadesReferenceValidation.Uri(), signatureScope.DocumentName()) {
								individualSignatureScopes = append(individualSignatureScopes, signatureScope)
							}
						}
					}
				}
			}
		}
		return individualSignatureScopes
	}
	return f.EncapsulatedTimestampScopeFinder.FilterCoveredSignatureScopes(timestampToken)
}

// xadesTimestampScopeFinderIsContentTimestampedReference ports the private
// isContentTimestampedReference(ReferenceValidation, List<TimestampInclude>).
func xadesTimestampScopeFinderIsContentTimestampedReference(xadesReferenceValidation *ReferenceValidation,
	includes []*validation.TimestampInclude) bool {
	if xadesReferenceValidation.Id() != "" {
		for _, timestampInclude := range includes {
			if xadesReferenceValidation.Id() == timestampInclude.URI() {
				return true
			}
		}
	}
	return false
}

// compile-time interface assertion.
var _ spiscope.TimestampScopeFinder = (*TimestampScopeFinder)(nil)
