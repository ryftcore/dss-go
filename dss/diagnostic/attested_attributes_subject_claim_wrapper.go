// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/claim/AttestedAttributesSubjectClaimWrapper.java (DSS 6.5.RC1).
package diagnostic

import "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"

// AttestedAttributesSubjectClaimWrapper wraps a jaxb.XmlAttestedAttributesSubjectClaim. Unlike
// most claim subtype wrappers, Java declares only the single-argument constructor here (no
// parent-taking overload).
type AttestedAttributesSubjectClaimWrapper struct {
	ClaimWrapper

	// wrapped is the concrete XmlAttestedAttributesSubjectClaim, shadowing the promoted
	// (synthetic) field of the embedded ClaimWrapper; see the covariant-getWrapped note in
	// claim_wrapper.go.
	wrapped *jaxb.XmlAttestedAttributesSubjectClaim
}

// NewAttestedAttributesSubjectClaimWrapper is the default constructor. Port of
// AttestedAttributesSubjectClaimWrapper(XmlAttestedAttributesSubjectClaim).
func NewAttestedAttributesSubjectClaimWrapper(wrapped *jaxb.XmlAttestedAttributesSubjectClaim) *AttestedAttributesSubjectClaimWrapper {
	w := &AttestedAttributesSubjectClaimWrapper{
		ClaimWrapper: *NewClaimWrapper(claimBase(wrapped.XmlClaimContent, wrapped.XmlClaimAttrs)),
		wrapped:      wrapped,
	}
	w.InitClaimOverrides(w)
	return w
}

// SubjectId gets the status's unique index identifier (sic, the Java Javadoc is copy-pasted
// across all three accessors of this class; reproduced as-is). Port of getSubjectId().
func (w *AttestedAttributesSubjectClaimWrapper) SubjectId() *AttestedAttributesSubjectClaimIdWrapper {
	if w.wrapped.SubjectId != nil {
		return NewAttestedAttributesSubjectClaimIdWrapperWithParent(w.wrapped.SubjectId, &w.ClaimWrapper)
	}
	return nil
}

// SubjectPseudonym gets the status's unique index identifier (sic, see the SubjectId note).
// Port of getSubjectPseudonym().
func (w *AttestedAttributesSubjectClaimWrapper) SubjectPseudonym() *ClaimWrapper {
	if w.wrapped.SubjectPseudonym != nil {
		return NewClaimWrapperWithParent(w.wrapped.SubjectPseudonym, &w.ClaimWrapper)
	}
	return nil
}

// Attributes gets the status's unique index identifier (sic, see the SubjectId note). Port of
// getAttributes().
func (w *AttestedAttributesSubjectClaimWrapper) Attributes() *ClaimWrapper {
	if w.wrapped.Attributes != nil {
		return NewClaimWrapperWithParent(w.wrapped.Attributes, &w.ClaimWrapper)
	}
	return nil
}

// Map is the override, assembling the map from the dedicated subject-id/pseudonym/attributes
// child claims rather than the generic Entry list. Port of the overridden getMap(); NOTE: unlike
// most claim subtypes that override getMap(), Java does not additionally override isMap() here -
// IsMap therefore keeps the base's generic Entry-derived behaviour (which is false for a value
// built from these dedicated fields), reproduced as ported, not "fixed".
func (w *AttestedAttributesSubjectClaimWrapper) Map() map[string]*ClaimWrapper {
	result := map[string]*ClaimWrapper{}
	for k, v := range w.ClaimWrapper.Map() {
		result[k] = v
	}
	if subjectID := w.SubjectId(); subjectID != nil {
		result[subjectID.Name()] = subjectID.AsClaim()
	}
	if subjectPseudonym := w.SubjectPseudonym(); subjectPseudonym != nil {
		result[subjectPseudonym.Name()] = subjectPseudonym
	}
	if attributes := w.Attributes(); attributes != nil {
		result[attributes.Name()] = attributes
	}
	return result
}

// Wrapped is the covariant override of ClaimWrapper.Wrapped(). Port of the covariant getWrapped().
func (w *AttestedAttributesSubjectClaimWrapper) Wrapped() *jaxb.XmlAttestedAttributesSubjectClaim {
	return w.wrapped
}

// AsClaim views the wrapper as its ClaimWrapper base type. Map() is baked in even though IsMap()
// is not overridden, matching Java where getMap() is reachable (and returns the assembled map)
// through a ClaimWrapper-typed reference regardless of what isMap() reports; see the package
// note in claim_wrapper.go.
func (w *AttestedAttributesSubjectClaimWrapper) AsClaim() *ClaimWrapper {
	cw := w.ClaimWrapper
	cw.mapOverride = w.Map()
	return &cw
}
