// Ported from
// dss-pades/src/main/java/eu/europa/esig/dss/pades/validation/timestamp/PAdESTimestampSource.java
// (DSS 6.5.RC1).
//
// # Package layout
//
// Java's PAdESTimestampSource extends dss-cades's CAdESTimestampSource<CAdESSignature,CAdESAttribute>
// unparametrized (it never re-binds the type parameters to Signature/PAdESAttribute - there
// is no PAdESAttribute), relying on virtual dispatch so that TimestampSource's own methods,
// when invoked on a TimestampSource instance, still see the right overrides. This file
// embeds cades.TimestampSource (concrete, itself embedding
// timestamp.SignatureTimestampSource[*cades.Signature, *cades.Attribute]) the same
// way, and re-targets InitSignatureTimestampSource(s) at construction - the InitTimestampIdentifierBuilder-
// style override-registration pattern used for timestamp sources - so that every
// SignatureTimestampSourceOverrides call the base machinery dispatches
// through s.overrides (all the Is*/Make*/Get* checks in makeTimestampTokensFromUnsignedAttributes)
// correctly reaches this type's own shadowed methods (IsCompleteCertificateRef and friends,
// below) instead of TimestampSource's.
//
// # GAP flagged for integrator: two CAdES/base "concrete but overridable" hooks are unreachable
// # from the base's own internal driver
//
// Unlike every check in makeTimestampTokensFromUnsignedAttributes (all reached through
// s.overrides.Is*/Get*Refs/... and therefore correctly dispatched to this type), two methods
// Java's PAdESTimestampSource overrides are NOT part of SignatureTimestampSourceOverrides and
// are called by the base machinery through its own *unexported, non-virtual* methods:
//
//   - getSignatureTimestampReferences() (base SignatureTimestampSource, private in this port:
//     signature_timestamp_source.go's getSignatureTimestampReferences): the base's own
//     makeTimestampTokensFromUnsignedAttributes calls its *own* private copy, unqualified, when
//     building references for a CMS-embedded "signature-timestamp" unsigned attribute (the
//     mechanism PAdES-BASELINE-T itself uses). Java's override
//     (TimestampSource.getSignatureTimestampReferences(), adding
//     getAdbeRevocationInfoArchivalReferences() on top of super's result) is therefore never
//     reached via that call site: a PAdES-BASELINE-T signature-timestamp's own reference set,
//     when produced by the base's *own* CMS-embedded-attribute pass, will be missing the
//     adbe-revocationInfoArchival references PAdES additionally covers. This file's own
//     GetSignatureTimestampReferences (exported) is a faithful, best-effort reconstruction of
//     Java's override (built from GetSignatureSignedDataReferences/SignerDataReferences, both
//     exported and promoted, plus GetAdbeRevocationInfoArchivalReferences below) and IS reached
//     correctly by this file's own document/VRI-revision loop below (which calls it directly,
//     not through the base's internal dispatch) - see the loop's own comment.
//   - getTimestampScopes(TimestampToken) (base, private: signature_timestamp_source.go's
//     getTimestampScopes) is used the same unreachable way by the base's own validateTimestamps()
//     for content-timestamp and archive-timestamp scope assignment (both PAdES-CMS corner cases:
//     an ordinary PAdES-BASELINE-{B,T,LT,LTA} signature carries neither). Java's override (using
//     TimestampScopeFinder instead of the default EncapsulatedTimestampScopeFinder) is
//     reached correctly by this file's own document-timestamp loop below (which calls
//     TimestampScopeFinder directly, exactly mirroring the override's body), just not from
//     that one base call site.
//
// Both gaps are narrow (CMS-embedded content/archive timestamps are not part of the PAdES
// baseline profiles) and match the severity/shape of the three GAPs cades_timestamp_source.go's
// own header already documents for the identical "concrete-but-not-abstract Java method, no
// escape hatch in the ported interface" situation; flagged prominently here for the same reason.
//
// # Document/VRI timestamp population is this file's own, self-contained driver
//
// Java's makeTimestampTokensFromUnsignedAttributes() override calls
// `super.makeTimestampTokensFromUnsignedAttributes()` (running the base's ordinary CMS-embedded-
// attribute pass, populating signature/content/archive timestamps) and then walks documentRevisions
// to additionally populate documentTimestamps/vriTimestamps - all within ONE Java method, so the
// PDF-specific loop can freely call the (by then already-reachable) getSignatureTimestamps(),
// unsignedPropertiesReferences, etc. The base's own populateTimestampTokens driver
// (signature_timestamp_source.go) has no override hook for "run extra logic after the unsigned-
// attributes pass" (see the GAP above), so this file instead triggers the base's ordinary pass
// via the promoted SignatureTimestamps() getter (which lazily runs createAndValidate() the first
// time any of ContentTimestamps/SignatureTimestamps/.../ArchiveTimestamps is read - fully
// reachable and correctly dispatched, per this header's opening paragraph) and then runs its own
// self-contained loop below, calling GetSignatureTimestampReferences (this file's own,
// reconstructed) and TimestampScopeFinder directly rather than through any base dispatch -
// exactly the two GAP-flagged call sites above, made reachable by not routing through the base
// at all for this part.
//
// populateSources(TimestampToken) (base, private, mutating the base's own certificateSource/
// crlSource/ocspSource fields) is reproduced here without needing access to those private
// fields: CertificateSource()/CRLSource()/OCSPSource() (base, exported, promoted) return the
// live pointers to those same merged sources, so calling .Add(...) on the returned pointers
// mutates the exact state populateSources would have - see padesTSPopulateSources below.
//
// s.signature.PdfRevision() returns *PdfSignatureRevision, compared here by pointer identity
// against each documentRevisions entry, exactly as Java's `padesSignature.getPdfRevision() ==
// pdfRevision` reference comparison does.
package pades

import (
	"github.com/ryftcore/dss-go/dss/cades"
	"github.com/ryftcore/dss-go/dss/crlparser"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/scope"
	"github.com/ryftcore/dss-go/dss/model/x509/revocation"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/spi/validation/timestamp"
	"github.com/ryftcore/dss-go/dss/utils"
)

// PAdESTimestampSource extracts timestamps from a PAdES document. Port of the class
// TimestampSource, extending cades.TimestampSource.
type TimestampSource struct {
	cades.TimestampSource

	// signature is the Signature this source extracts timestamps for.
	signature *Signature

	// documentRevisions is a list of document PdfRevisions, in reverse (most recent first)
	// order - Java's `this.documentRevisions = Utils.reverseList(documentRevisions);`.
	documentRevisions []PdfRevision

	// documentTimestamps contains the list of embedded document timestamps; nil stands for
	// "not computed yet", matching the base's own contentTimestamps/... convention.
	documentTimestamps []*validation.TimestampToken

	// vriTimestamps contains the list of embedded /VRI timestamps corresponding to the
	// signature.
	vriTimestamps []*validation.TimestampToken
}

// NewPAdESTimestampSource is the default constructor to extract timestamps for a signature.
// Port of the PAdESTimestampSource(PAdESSignature, List<PdfRevision>) constructor.
//
// Panics with the Java message when documentRevisions is nil (Objects.requireNonNull).
func NewPAdESTimestampSource(signature *Signature, documentRevisions []PdfRevision) *TimestampSource {
	if documentRevisions == nil {
		panic("List of Document revisions must be provided!")
	}
	cadesBase := cades.NewCAdESTimestampSource(signature.Signature)
	s := &TimestampSource{
		TimestampSource:   *cadesBase,
		signature:         signature,
		documentRevisions: utils.ReverseList(documentRevisions),
	}
	// Re-targets the base's overrides registration at s (see this file's header): cadesBase's
	// own InitSignatureTimestampSource(cadesBase) call, made inside cades.NewCAdESTimestampSource
	// above, is superseded by this one.
	s.InitSignatureTimestampSource(s)
	return s
}

// DocumentTimestamps returns the list of embedded document timestamps, computing them (and the
// /VRI timestamps) on first use. Port of the getDocumentTimestamps() override.
func (s *TimestampSource) DocumentTimestamps() []*validation.TimestampToken {
	if s.documentTimestamps == nil {
		s.populateAndValidateDocumentTimestamps()
	}
	return s.documentTimestamps
}

// VriTimestamps returns a list of incorporated /VRI timestamps for the corresponding signature.
// Port of getVriTimestamps().
func (s *TimestampSource) VriTimestamps() []*validation.TimestampToken {
	if s.vriTimestamps == nil {
		s.populateAndValidateDocumentTimestamps()
	}
	return s.vriTimestamps
}

// AllTimestamps returns a list of all incorporated timestamps. Port of the getAllTimestamps()
// override.
func (s *TimestampSource) AllTimestamps() []*validation.TimestampToken {
	timestampTokens := s.TimestampSource.AllTimestamps()
	timestampTokens = append(timestampTokens, s.DocumentTimestamps()...)
	timestampTokens = append(timestampTokens, s.VriTimestamps()...)
	return timestampTokens
}

// IsCompleteCertificateRef implements timestamp.SignatureTimestampSourceOverrides.
// NOTE: not applicable for PAdES. Port of the isCompleteCertificateRef(CAdESAttribute) override.
func (s *TimestampSource) IsCompleteCertificateRef(unsignedAttribute *cades.Attribute) bool {
	return false
}

// IsAttributeCertificateRef implements timestamp.SignatureTimestampSourceOverrides.
// NOTE: not applicable for PAdES. Port of the isAttributeCertificateRef(CAdESAttribute) override.
func (s *TimestampSource) IsAttributeCertificateRef(unsignedAttribute *cades.Attribute) bool {
	return false
}

// IsCompleteRevocationRef implements timestamp.SignatureTimestampSourceOverrides.
// NOTE: not applicable for PAdES. Port of the isCompleteRevocationRef(CAdESAttribute) override.
func (s *TimestampSource) IsCompleteRevocationRef(unsignedAttribute *cades.Attribute) bool {
	return false
}

// IsAttributeRevocationRef implements timestamp.SignatureTimestampSourceOverrides.
// NOTE: not applicable for PAdES. Port of the isAttributeRevocationRef(CAdESAttribute) override.
func (s *TimestampSource) IsAttributeRevocationRef(unsignedAttribute *cades.Attribute) bool {
	return false
}

// IsRefsOnlyTimestamp implements timestamp.SignatureTimestampSourceOverrides.
// NOTE: not applicable for PAdES. Port of the isRefsOnlyTimestamp(CAdESAttribute) override.
func (s *TimestampSource) IsRefsOnlyTimestamp(unsignedAttribute *cades.Attribute) bool {
	return false
}

// IsSigAndRefsTimestamp implements timestamp.SignatureTimestampSourceOverrides.
// NOTE: not applicable for PAdES. Port of the isSigAndRefsTimestamp(CAdESAttribute) override.
func (s *TimestampSource) IsSigAndRefsTimestamp(unsignedAttribute *cades.Attribute) bool {
	return false
}

// IsCertificateValues implements timestamp.SignatureTimestampSourceOverrides.
// NOTE: not applicable for PAdES. Port of the isCertificateValues(CAdESAttribute) override.
func (s *TimestampSource) IsCertificateValues(unsignedAttribute *cades.Attribute) bool {
	return false
}

// IsRevocationValues implements timestamp.SignatureTimestampSourceOverrides.
// NOTE: not applicable for PAdES. Port of the isRevocationValues(CAdESAttribute) override.
func (s *TimestampSource) IsRevocationValues(unsignedAttribute *cades.Attribute) bool {
	return false
}

// IsArchiveTimestamp implements timestamp.SignatureTimestampSourceOverrides.
// NOTE: not applicable for PAdES. Port of the isArchiveTimestamp(CAdESAttribute) override.
func (s *TimestampSource) IsArchiveTimestamp(unsignedAttribute *cades.Attribute) bool {
	return false
}

// GetCounterSignatures implements timestamp.SignatureTimestampSourceOverrides.
// NOTE: not supported in PAdES. Port of the getCounterSignatures(CAdESAttribute) override.
func (s *TimestampSource) GetCounterSignatures(unsignedAttribute *cades.Attribute) []validation.AdvancedSignature {
	return []validation.AdvancedSignature{}
}

// GetSignatureTimestampReferences is a faithful, self-contained reconstruction of the
// getSignatureTimestampReferences() override (base SignatureTimestampSource.
// getSignatureTimestampReferences(), augmented with GetAdbeRevocationInfoArchivalReferences());
// see this file's header GAP note for why it cannot be wired into the base's own internal
// dispatch, and why it is instead called directly by populateAndValidateDocumentTimestamps below.
func (s *TimestampSource) GetSignatureTimestampReferences() []*validation.TimestampedReference {
	references := []*validation.TimestampedReference{}
	padesTSAddReferences(&references, padesTSEncapsulatedReferencesFromTimestamps(s.ContentTimestamps(), s.CertificateSource(), s.CRLSource(), s.OCSPSource()))
	padesTSAddReferences(&references, s.SignerDataReferences())
	padesTSAddReference(&references, validation.NewTimestampedReference(s.signature.ID(), enumerations.TimestampedObjectTypeSignature))
	signatureCertificateSource := s.signature.CertificateSource()
	padesTSAddReferences(&references, timestamp.CreateReferencesForCertificateRefs(signatureCertificateSource.SigningCertificateRefs(),
		signatureCertificateSource, s.CertificateSource()))
	padesTSAddReferences(&references, s.GetAdbeRevocationInfoArchivalReferences())
	return references
}

// GetAdbeRevocationInfoArchivalReferences returns a list of revocation data TimestampedReferences
// from the adbe-revocationInfoArchival signed attribute. Port of
// getAdbeRevocationInfoArchivalReferences().
func (s *TimestampSource) GetAdbeRevocationInfoArchivalReferences() []*validation.TimestampedReference {
	signedSignatureProperties := s.BuildSignedSignatureProperties()
	// Upstream tests isExist() alone; BuildSignedSignatureProperties, like
	// CAdESTimestampSource#buildSignedSignatureProperties, always returns a value.
	if !signedSignatureProperties.IsExist() {
		return []*validation.TimestampedReference{}
	}
	references := []*validation.TimestampedReference{}
	for _, attribute := range signedSignatureProperties.Attributes() {
		if s.isAdbeRevocationInfoArchival(attribute) {
			revValues := UtilsRevocationInfoArchival(attribute.ASN1Object())
			if revValues != nil {
				crlBinaries := padesTSBuildCRLIdentifiers(revValues.CrlVals())
				padesTSAddReferences(&references, timestamp.CreateReferencesForCRLBinaries(crlBinaries))
				ocspBinaries := padesTSBuildOCSPIdentifiers(revValues.OcspVals())
				ocspReferences, err := timestamp.CreateReferencesForOCSPBinaries(ocspBinaries, s.CertificateSource())
				if err != nil {
					panic(model.NewDSSErrorWithCause(err))
				}
				padesTSAddReferences(&references, ocspReferences)
			}
		}
	}
	return references
}

// isAdbeRevocationInfoArchival checks if signedAttribute is an instance of type
// adbe-revocationInfoArchival. Port of isAdbeRevocationInfoArchival(CAdESAttribute).
func (s *TimestampSource) isAdbeRevocationInfoArchival(signedAttribute *cades.Attribute) bool {
	return spi.OIDAdbeRevocationInfoArchival.Equal(signedAttribute.ASN1Oid())
}

// getTimestampScopesForDocumentTimestamp generates timestamp token scopes for a document/VRI
// time-stamp, using a PAdESTimestampScopeFinder. Port of the getTimestampScopes(TimestampToken)
// override; named distinctly since it is this file's own driver, not part of
// SignatureTimestampSourceOverrides (see this file's header GAP note).
func (s *TimestampSource) getTimestampScopesForDocumentTimestamp(timestampToken *validation.TimestampToken) []scope.SignatureScope {
	timestampScopeFinder := NewPAdESTimestampScopeFinder()
	timestampScopeFinder.SetSignature(s.signature)
	return timestampScopeFinder.FindTimestampScope(timestampToken)
}

// populateAndValidateDocumentTimestamps is this file's own self-contained driver, reconstructing
// Java's makeTimestampTokensFromUnsignedAttributes() override (the document/VRI-revision loop)
// followed by the VRI-timestamp half of the validateTimestamps() override; see this file's
// header for why the base's own internal population flow cannot reach these directly.
func (s *TimestampSource) populateAndValidateDocumentTimestamps() {
	// Triggers the base's own (correctly-dispatched, per this file's header) CMS-embedded pass:
	// content/signature/X1/X2/archive timestamps, and the merged certificate/CRL/OCSP sources
	// this loop below reads through CertificateSource()/CRLSource()/OCSPSource().
	cadesSignatureTimestamps := s.SignatureTimestamps()

	s.documentTimestamps = []*validation.TimestampToken{}
	s.vriTimestamps = []*validation.TimestampToken{}

	unsignedPropertiesReferences := []*validation.TimestampedReference{}
	processedPdfRevisionTimestamps := []*validation.TimestampToken{}

	signatureRevisionReached := false
	dssRevisionReached := false

	for _, pdfRevision := range s.documentRevisions {
		switch revision := pdfRevision.(type) {

		case *PdfDocTimestampRevision:
			individualTimestampReferences := []*validation.TimestampedReference{}

			pdfTimestampToken := revision.TimestampToken()
			timestampToken := pdfTimestampToken.TimestampToken

			timestampScopes := s.getTimestampScopesForDocumentTimestamp(timestampToken)
			timestampToken.SetTimestampScopes(timestampScopes)
			padesTSAddReferences(&individualTimestampReferences, timestamp.SignerDataTimestampedReferences(timestampScopes))

			if dssRevisionReached {
				timestampToken.SetArchiveTimestampType(enumerations.ArchiveTimestampTypePAdES)
			}
			if signatureRevisionReached {
				padesTSAddReferences(&individualTimestampReferences, s.GetSignatureTimestampReferences())
				padesTSAddReferences(&individualTimestampReferences, s.GetSignatureSignedDataReferences())
				padesTSAddReferences(&individualTimestampReferences, padesTSEncapsulatedReferencesFromTimestamps(cadesSignatureTimestamps, s.CertificateSource(), s.CRLSource(), s.OCSPSource()))
			}
			if utils.IsCollectionNotEmpty(unsignedPropertiesReferences) {
				// covers DSS dictionary
				padesTSAddReferences(&individualTimestampReferences, unsignedPropertiesReferences)
			}
			padesTSAddReferences(&individualTimestampReferences, padesTSEncapsulatedReferencesFromTimestamps(processedPdfRevisionTimestamps, s.CertificateSource(), s.CRLSource(), s.OCSPSource()))

			// references embedded to timestamp's content are covered by outer timestamps
			existingReferences := timestampToken.TimestampedReferences()
			padesTSAddReferences(&existingReferences, individualTimestampReferences)
			timestampToken.SetTimestampedReferences(existingReferences)

			if signatureRevisionReached {
				s.documentTimestamps = append(s.documentTimestamps, timestampToken)
			}

			padesTSPopulateSources(s, timestampToken)
			processedPdfRevisionTimestamps = append(processedPdfRevisionTimestamps, timestampToken)

		case *PdfDocDssRevision:
			pdfRevisionTimestampSource := NewPdfRevisionTimestampSource(revision, s.CertificateSource(), s.CRLSource(), s.OCSPSource())
			padesTSAddReferences(&unsignedPropertiesReferences, pdfRevisionTimestampSource.IncorporatedReferences())

			s.CertificateSource().Add(revision.CertificateSource())
			s.CRLSource().Add(revision.CRLSource())
			s.OCSPSource().Add(revision.OCSPSource())

			vriTimestampToken := pdfRevisionTimestampSource.VRITimestampToken(s.signature.VRIKey())
			if vriTimestampToken != nil && !padesTSContainsTimestampToken(s.vriTimestamps, vriTimestampToken) {
				existingReferences := vriTimestampToken.TimestampedReferences()
				padesTSAddReferences(&existingReferences, s.GetSignatureTimestampReferences())
				vriTimestampToken.SetTimestampedReferences(existingReferences)

				if signatureRevisionReached {
					s.vriTimestamps = append(s.vriTimestamps, vriTimestampToken)
				}
				padesTSPopulateSources(s, vriTimestampToken)
				processedPdfRevisionTimestamps = append(processedPdfRevisionTimestamps, vriTimestampToken)
			}
			dssRevisionReached = true

		case *PdfSignatureRevision:
			if s.signature.PdfRevision() == revision {
				signatureRevisionReached = true
			}
		}
	}

	// Port of the VRI-timestamp half of the validateTimestamps() override (the
	// super.validateTimestamps() half already ran inside the SignatureTimestamps() call above).
	for _, timestampToken := range s.vriTimestamps {
		messageDigestBuilder := s.GetTimestampMessageImprintDigestBuilderForToken(timestampToken)
		messageDigest := messageDigestBuilder.SignatureTimestampMessageDigest()
		timestampToken.MatchDataMessageDigest(messageDigest)
	}
}

// padesTSPopulateSources allows populating the merged certificate/CRL/OCSP sources with data
// extracted from timestampToken. Port of the inherited protected populateSources(TimestampToken);
// see this file's header for why the base's own private populateSources is reproduced this way.
func padesTSPopulateSources(s *TimestampSource, timestampToken *validation.TimestampToken) {
	if timestampToken != nil {
		s.CertificateSource().Add(timestampToken.CertificateSource())
		s.CRLSource().Add(timestampToken.CRLSource())
		s.OCSPSource().Add(timestampToken.OCSPSource())
	}
}

// padesTSContainsTimestampToken reports whether timestampTokens already contains candidate, by
// pointer identity - the Go counterpart of Java's `!vriTimestamps.contains(vriTimestampToken)`
// (TimestampToken has no value-based equals() override in this port either way, so reference
// identity is what Java's default Object.equals would compare too).
func padesTSContainsTimestampToken(timestampTokens []*validation.TimestampToken, candidate *validation.TimestampToken) bool {
	for _, timestampToken := range timestampTokens {
		if timestampToken == candidate {
			return true
		}
	}
	return false
}

// -----------------------------------------------------------------------------
// Small helpers mirroring the frozen spi/validation/timestamp package's own unexported
// addReference/addReferences/getEncapsulatedReferencesFromTimestamps/buildCRLIdentifiers/
// buildOCSPIdentifiers (signature_timestamp_source.go / cades_timestamp_source.go), which this
// file cannot call directly (different package / unexported). Prefixed distinctly (padesTS) to
// avoid colliding with any sibling file in this package needing the same helper - see
// cades_timestamp_source.go's cadesTS-prefixed precedent for the same situation.
// -----------------------------------------------------------------------------

// padesTSAddReference adds referenceToAdd to *referenceList without duplicates.
func padesTSAddReference(referenceList *[]*validation.TimestampedReference, referenceToAdd *validation.TimestampedReference) {
	padesTSAddReferences(referenceList, []*validation.TimestampedReference{referenceToAdd})
}

// padesTSAddReferences adds referencesToAdd to *referenceList without duplicates (by
// TimestampedReference.Equals).
func padesTSAddReferences(referenceList *[]*validation.TimestampedReference, referencesToAdd []*validation.TimestampedReference) {
	for _, candidate := range referencesToAdd {
		found := false
		for _, existing := range *referenceList {
			if existing.Equals(candidate) {
				found = true
				break
			}
		}
		if !found {
			*referenceList = append(*referenceList, candidate)
		}
	}
}

// padesTSEncapsulatedReferencesFromTimestamps returns a list of TimestampedReferences for
// tokens encapsulated within timestampTokens. Port of the protected
// getEncapsulatedReferencesFromTimestamps(List).
func padesTSEncapsulatedReferencesFromTimestamps(timestampTokens []*validation.TimestampToken,
	certificateSource *spi.ListCertificateSource, crlSource *spi.ListRevocationSource[revocation.CRL],
	ocspSource *spi.ListRevocationSource[revocation.OCSP]) []*validation.TimestampedReference {
	references := []*validation.TimestampedReference{}
	for _, timestampToken := range timestampTokens {
		fromTimestamp, err := timestamp.ReferencesFromTimestamp(timestampToken, certificateSource, crlSource, ocspSource)
		if err != nil {
			panic(model.NewDSSErrorWithCause(err))
		}
		padesTSAddReferences(&references, fromTimestamp)
	}
	return references
}

// padesTSBuildCRLIdentifiers builds CRLBinary identifiers from the DER-encoded CertificateList
// values of a RevocationInfoArchival. Port of CAdESTimestampSource's inherited protected
// buildCRLIdentifiers(CertificateList...) - reproduced here (unreachable, being private to a
// different package) using RevocationInfoArchival.CrlVals()'s already-DER-encoded bytes, the
// same conversion pdf_cms_crl_source.go's extractRevocationInfoArchival already applies.
func padesTSBuildCRLIdentifiers(crlVals [][]byte) []*crlparser.CRLBinary {
	var crlBinaryIdentifiers []*crlparser.CRLBinary
	for _, crlVal := range crlVals {
		crlBinary, err := crlparser.CRLUtilsBuildCRLBinary(crlVal)
		if err != nil {
			// Upstream logs "Unable to parse CRL binaries : {}".
			continue
		}
		crlBinaryIdentifiers = append(crlBinaryIdentifiers, crlBinary)
	}
	return crlBinaryIdentifiers
}

// padesTSBuildOCSPIdentifiers builds OCSPResponseBinary identifiers from the DER-encoded
// OCSPResponse values of a RevocationInfoArchival. Port of CAdESTimestampSource's inherited
// protected buildOCSPIdentifiers(BasicOCSPResp...), composed with
// DSSASN1Utils.toBasicOCSPResps(OCSPResponse...) (Java's
// `buildOCSPIdentifiers(DSSASN1Utils.toBasicOCSPResps(revValues.getOcspVals()))`): unlike CAdES's
// own revocation-values attribute (already-bare BasicOCSPResponse bytes, see
// cades_timestamp_source.go's buildOCSPIdentifiers), RevocationInfoArchival.OcspVals() carries
// full OCSPResponse DER encodings, so this needs the same two-step OCSPResp -> BasicOCSPResp
// unwrap pdf_cms_ocsp_source.go's extractRevocationInfoArchival already applies.
func padesTSBuildOCSPIdentifiers(ocspVals [][]byte) []*spi.OCSPResponseBinary {
	var ocspIdentifiers []*spi.OCSPResponseBinary
	for _, ocspVal := range ocspVals {
		ocspResp := spi.DSSRevocationUtilsOcspResp(ocspVal)
		if ocspResp == nil {
			continue
		}
		basicOCSPResp := spi.DSSRevocationUtilsFromRespToBasic(ocspResp)
		if basicOCSPResp == nil {
			continue
		}
		binary, err := spi.OCSPResponseBinaryBuild(basicOCSPResp)
		if err != nil {
			// Upstream logs "Unable to parse OCSP response binaries : {}".
			continue
		}
		ocspIdentifiers = append(ocspIdentifiers, binary)
	}
	return ocspIdentifiers
}
