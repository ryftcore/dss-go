// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/validation/PAdESSignature.java
// (DSS 6.5.RC1).
//
// Signature extends cades.Signature across a package boundary (dss-cades ->
// dss-pades in Java; cades -> pades in this port). Go has no protected
// field visibility across packages, so every Java protected field this class reads/writes
// (offlineCertificateSource, signatureCRLSource, signatureOCSPSource, signatureTimestampSource,
// detachedContents) is reached instead through the accessor pairs
// spi/validation/default_advanced_signature.go already exports for exactly this purpose
// (OfflineCertificateSource/SetOfflineCertificateSource, SignatureCRLSource/
// SetSignatureCRLSource, SignatureOCSPSource/SetSignatureOCSPSource, SignatureTimestampSource/
// SetSignatureTimestampSource, SetDetachedContents), promoted through the embedded
// *cades.Signature. NewPAdESSignature re-registers virtual dispatch onto the Signature
// value itself (InitDefaultAdvancedSignature(s), promoted the same way NewCAdESSignature calls
// it onto itself) so that every DefaultAdvancedSignatureOverrides method Signature does not
// itself define keeps resolving, unchanged, to *cades.Signature's implementation (Go method
// promotion standing in for Java's single-dispatch inheritance), while every method PAdESSignature
// does define below shadows it.
package pades

import (
	"strings"
	"time"

	"github.com/ryftcore/dss-go/dss/cades"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/scope"
	"github.com/ryftcore/dss-go/dss/model/signature"
	"github.com/ryftcore/dss-go/dss/model/x509/revocation"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/utils"
)

// PAdESConstantsSignaturePKCS7SubFilter and PAdESConstantsSignaturePKCS7SHA1SubFilter are the
// two PKCS#7 SubFilter values (eu.europa.esig.dss.pdf.PAdESConstants.SIGNATURE_PKCS7_SUBFILTER /
// SIGNATURE_PKCS7_SHA1_SUBFILTER), matching the PAdESConstants<FieldName> naming already
// established elsewhere in this package (see pdf_signature_dictionary.go's header).
const (
	PAdESConstantsSignaturePKCS7SubFilter     = "adbe.pkcs7.detached"
	PAdESConstantsSignaturePKCS7SHA1SubFilter = "adbe.pkcs7.sha1"
)

// PAdESSignature is the implementation of AdvancedSignature for PAdES. Port of the class
// Signature, extending cades.Signature.
type Signature struct {
	*cades.Signature

	// pdfSignatureRevision is the corresponding PDF revision.
	pdfSignatureRevision *PdfSignatureRevision

	// documentRevisions contains a complete list of validating document revisions.
	documentRevisions []PdfRevision

	// dssCertificateSource is the certificate source obtained from DSS/VRI revisions.
	dssCertificateSource *spi.ListCertificateSource

	// dssCRLSource is the CRL source obtained from DSS/VRI revisions.
	dssCRLSource *spi.ListRevocationSource[revocation.CRL]

	// dssOCSPSource is the OCSP source obtained from DSS/VRI revisions.
	dssOCSPSource *spi.ListRevocationSource[revocation.OCSP]

	// vriKey is the SHA-1 key computed on /Contents of the signature.
	vriKey string
}

// NewPAdESSignature is the default constructor for Signature.
// Port of the protected PAdESSignature(PdfSignatureRevision, List<PdfRevision>) constructor.
func NewPAdESSignature(pdfSignatureRevision *PdfSignatureRevision, documentRevisions []PdfRevision) *Signature {
	cadesSignature := cades.NewCAdESSignature(pdfSignatureRevision.CMS(),
		spi.DSSASN1UtilsFirstSignerInformation(pdfSignatureRevision.CMS().SignerInfos()))
	s := &Signature{
		Signature:            cadesSignature,
		pdfSignatureRevision: pdfSignatureRevision,
		documentRevisions:    documentRevisions,
	}
	s.InitDefaultAdvancedSignature(s)
	s.SetDetachedContents([]model.DSSDocument{pdfSignatureRevision.SignedData()})

	// Eagerly warms the offline certificate/CRL/OCSP source caches - shared embedded
	// DefaultAdvancedSignature fields Signature.CertificateSource/CRLSource/OCSPSource also
	// lazily populate on first call - with THIS type's own PAdES-enriched (DSS-dictionary and
	// VRI aware) sources, before anything else gets a chance to populate them first.
	//
	// Without this, the generic timestamp-source machinery (spi/validation/timestamp's
	// SignatureTimestampSource[*cades.Signature, *cades.Attribute], concretely typed
	// per pades_timestamp_source.go's header on why it mirrors Java's unparametrized
	// TimestampSource<Signature,Attribute> inheritance) calls
	// s.signature.CertificateSource()/CRLSource()/OCSPSource() through the embedded
	// *cades.Signature field the moment DocumentTimestamps() is first queried (directly, or
	// via PDF modification-detection post-processing in PDFDocumentAnalyzer.postProcessing).
	// Since Go has no runtime vtable for a concrete (non-interface) generic type parameter - only
	// Java's virtual dispatch reaches PAdESSignature.getCertificateSource() there - that call
	// resolves to Signature's OWN method body, which lazily constructs a plain CAdES
	// certificate source (CMS-embedded certs only, no /DSS or /VRI dictionary awareness) and
	// PERMANENTLY caches it into the very same shared offlineCertificateSource/
	// signatureCRLSource/signatureOCSPSource fields this type's own CertificateSource()/
	// CRLSource()/OCSPSource() overrides read from - so whichever caller reaches the lazy
	// "if cached == nil" check first wins the cache for the signature's whole lifetime.
	//
	// Confirmed by pades/testdata/upstream/validation/PAdES-LT.pdf in the PAdES cross-validation
	// harness, whose signing certificates are embedded only in the document's /DSS and /VRI
	// dictionaries, not in each signature's own CMS: without this warm-up,
	// SigningCertificateToken()/CryptographicVerification()/DataFoundUpToLevel() for
	// such a signature silently degrade to "certificate not found"/broken signature/PAdES_BES the
	// moment DocumentTimestamps() is queried.
	s.CertificateSource()
	s.CRLSource()
	s.OCSPSource()

	return s
}

// SetDssCertificateSource sets a joint DSS/VRI Certificate Source.
// Port of setDssCertificateSource(ListCertificateSource).
func (s *Signature) SetDssCertificateSource(dssCertificateSource *spi.ListCertificateSource) {
	s.dssCertificateSource = dssCertificateSource
}

// SetDssCRLSource sets a joint DSS/VRI CRL Source. Port of setDssCRLSource(ListRevocationSource).
func (s *Signature) SetDssCRLSource(dssCRLSource *spi.ListRevocationSource[revocation.CRL]) {
	s.dssCRLSource = dssCRLSource
}

// SetDssOCSPSource sets a joint DSS/VRI OCSP Source. Port of setDssOCSPSource(ListRevocationSource).
func (s *Signature) SetDssOCSPSource(dssOCSPSource *spi.ListRevocationSource[revocation.OCSP]) {
	s.dssOCSPSource = dssOCSPSource
}

// SignatureForm specifies the format of the signature. Port of getSignatureForm().
func (s *Signature) SignatureForm() enumerations.SignatureForm {
	if s.hasPKCS7SubFilter() {
		return enumerations.SignatureFormPKCS7
	}
	return enumerations.SignatureFormPAdES
}

// CertificateSource gets a certificate source which contains ALL certificates embedded in the
// signature. Port of getCertificateSource().
func (s *Signature) CertificateSource() *spi.SignatureCertificateSource {
	if s.OfflineCertificateSource() == nil {
		padesCertificateSource, err := NewPAdESCertificateSource(s.pdfSignatureRevision, s.VRIKey(), s.SignerInformation())
		if err != nil {
			panic(err)
		}
		s.SetOfflineCertificateSource(&padesCertificateSource.SignatureCertificateSource)
	}
	return s.OfflineCertificateSource()
}

// CRLSource gets a CRL source which contains ALL CRLs embedded in the signature.
// Port of getCRLSource().
func (s *Signature) CRLSource() spi.OfflineRevocationSource[revocation.CRL] {
	if s.SignatureCRLSource() == nil {
		s.SetSignatureCRLSource(NewPAdESCRLSource(s.pdfSignatureRevision, s.VRIKey(), s.SignerInformation().SignedAttributes))
	}
	return s.SignatureCRLSource()
}

// OCSPSource gets an OCSP source which contains ALL OCSP responses embedded in the signature.
// Port of getOCSPSource().
func (s *Signature) OCSPSource() spi.OfflineRevocationSource[revocation.OCSP] {
	if s.SignatureOCSPSource() == nil {
		s.SetSignatureOCSPSource(NewPAdESOCSPSource(s.pdfSignatureRevision, s.VRIKey(), s.SignerInformation().SignedAttributes))
	}
	return s.SignatureOCSPSource()
}

// CompleteCertificateSource returns a complete certificate source including the joint DSS/VRI
// source, when set. Port of getCompleteCertificateSource().
func (s *Signature) CompleteCertificateSource() *spi.ListCertificateSource {
	completeCertificateSource := s.Signature.CompleteCertificateSource()
	if s.dssCertificateSource != nil {
		completeCertificateSource.AddAll(s.dssCertificateSource)
	}
	return completeCertificateSource
}

// CompleteCRLSource returns a complete CRL source including the joint DSS/VRI source, when set.
// Port of getCompleteCRLSource().
func (s *Signature) CompleteCRLSource() *spi.ListRevocationSource[revocation.CRL] {
	completeCRLSource := s.Signature.CompleteCRLSource()
	if s.dssCRLSource != nil {
		completeCRLSource.AddAllFrom(s.dssCRLSource)
	}
	return completeCRLSource
}

// CompleteOCSPSource returns a complete OCSP source including the joint DSS/VRI source, when
// set. Port of getCompleteOCSPSource().
func (s *Signature) CompleteOCSPSource() *spi.ListRevocationSource[revocation.OCSP] {
	completeOCSPSource := s.Signature.CompleteOCSPSource()
	if s.dssOCSPSource != nil {
		completeOCSPSource.AddAllFrom(s.dssOCSPSource)
	}
	return completeOCSPSource
}

// TimestampSource gets a Signature Timestamp source which contains ALL timestamps embedded in
// the signature. Port of getTimestampSource(), covariant in Java (returns
// TimestampSource); Go has no covariant return, so PAdES-specific callers assert on the
// result, as in cades_signature.go's TimestampSource().
func (s *Signature) TimestampSource() validation.TimestampSource {
	if s.SignatureTimestampSource() == nil {
		s.SetSignatureTimestampSource(NewPAdESTimestampSource(s, s.documentRevisions))
	}
	return s.SignatureTimestampSource()
}

// DocumentTimestamps returns the list of document timestamps. Port of getDocumentTimestamps().
func (s *Signature) DocumentTimestamps() []*validation.TimestampToken {
	return s.TimestampSource().(*TimestampSource).DocumentTimestamps()
}

// VRITimestamps returns a list of timestamps enveloped within /VRI dictionary for the current
// signature. Port of getVRITimestamps().
func (s *Signature) VRITimestamps() []*validation.TimestampToken {
	return s.TimestampSource().(*TimestampSource).VriTimestamps()
}

// FindSignatureScopes finds signature scopes. Port of the protected findSignatureScopes()
// override.
func (s *Signature) FindSignatureScopes() []scope.SignatureScope {
	return NewPAdESSignatureScopeFinder().FindSignatureScope(s)
}

// SigningTime returns the claimed signing time. Port of getSigningTime().
func (s *Signature) SigningTime() *time.Time {
	signingTime := s.pdfSignatureRevision.SigningDate()
	return &signingTime
}

// ContentIdentifier returns nil, not applicable for PAdES. Port of getContentIdentifier().
func (s *Signature) ContentIdentifier() string {
	return ""
}

// ContentHints returns nil, not applicable for PAdES. Port of getContentHints().
func (s *Signature) ContentHints() string {
	return ""
}

// CounterSignatures returns an empty list, not applicable for PAdES.
// Port of getCounterSignatures().
func (s *Signature) CounterSignatures() []validation.AdvancedSignature {
	return []validation.AdvancedSignature{}
}

// OriginalDocument returns the original signed document. Port of getOriginalDocument().
func (s *Signature) OriginalDocument() (model.DSSDocument, error) {
	return s.pdfSignatureRevision.SignedData(), nil
}

// SignerDocumentContent returns the document content used for the digest computation.
// Port of the protected getSignerDocumentContent() override.
//
// ISO 32000-1: adbe.pkcs7.sha1: The SHA-1 digest of the document's byte range shall be
// encapsulated in the CMSSignedData field with ContentInfo of type Data.
func (s *Signature) SignerDocumentContent() (model.DSSDocument, error) {
	signerDocument, err := s.OriginalDocument()
	if err != nil {
		return nil, err
	}
	if signerDocument != nil && s.PdfSignatureDictionary() != nil &&
		PAdESConstantsSignaturePKCS7SHA1SubFilter == s.PdfSignatureDictionary().SubFilter() {
		digestValue, err := signerDocument.DigestValue(enumerations.DigestAlgorithmSHA1)
		if err != nil {
			return nil, err
		}
		signerDocument = model.NewInMemoryDocument(digestValue)
	}
	return signerDocument, nil
}

// SignatureIdentifierBuilder returns a builder to define and build a signature Id.
// Port of the protected getSignatureIdentifierBuilder() override.
func (s *Signature) SignatureIdentifierBuilder() validation.SignatureIdentifierBuilder {
	return NewPAdESSignatureIdentifierBuilder(s)
}

// BuildSignatureDigestReference builds a new DigestReference according to the PAdES
// format rules: the input of the digest value computation is the result of decoding the
// hexadecimal string present within the /Contents field of the Signature PDF dictionary
// enclosing one PAdES digital signature (TS 119 442 - V1.1.1, ch. 5.1.4.2.1.3 XML component).
// Port of buildSignatureDigestReference(DigestAlgorithm).
func (s *Signature) BuildSignatureDigestReference(digestAlgorithm enumerations.DigestAlgorithm) *signature.DigestReference {
	contents := s.PdfSignatureDictionary().Contents()
	digestValue, err := spi.DSSUtilsDigest(digestAlgorithm, contents)
	if err != nil {
		panic(err)
	}
	return signature.NewSignatureDigestReference(model.NewDigest(digestAlgorithm, digestValue))
}

// HasLTVProfile checks if the LTV-level is present in the signature. Port of hasLTVProfile().
func (s *Signature) HasLTVProfile() bool {
	return s.BaselineRequirementsChecker().HasExtendedLTVProfile()
}

// DataFoundUpToLevel returns the signature level up to which the data is found.
// Port of getDataFoundUpToLevel().
func (s *Signature) DataFoundUpToLevel() enumerations.SignatureLevel {
	signatureForm := s.SignatureForm()
	if enumerations.SignatureFormPAdES == signatureForm && s.HasBESProfile() {
		if !s.HasBProfile() {
			if s.HasLTVProfile() {
				return enumerations.SignatureLevelPAdESLTV
			}
			if s.HasEPESProfile() {
				return enumerations.SignatureLevelPAdESEPES
			}
			return enumerations.SignatureLevelPAdESBES
		}
		if !s.HasTProfile() {
			return enumerations.SignatureLevelPAdESBaselineB
		}
		if !s.HasLTProfile() {
			return enumerations.SignatureLevelPAdESBaselineT
		}
		if s.HasLTAProfile() {
			return enumerations.SignatureLevelPAdESBaselineLTA
		}
		return enumerations.SignatureLevelPAdESBaselineLT

	} else if enumerations.SignatureFormPKCS7 == signatureForm && s.HasPKCS7Profile() {
		if !s.HasPKCS7TProfile() {
			return enumerations.SignatureLevelPKCS7B
		}
		if !s.HasPKCS7LTProfile() {
			return enumerations.SignatureLevelPKCS7T
		}
		if s.HasPKCS7LTAProfile() {
			return enumerations.SignatureLevelPKCS7LTA
		}
		return enumerations.SignatureLevelPKCS7LT

	}
	return enumerations.SignatureLevelPDFNotETSI
}

// BaselineRequirementsChecker narrows the return type of the promoted
// CAdESSignature.BaselineRequirementsChecker(). Port of the protected
// getBaselineRequirementsChecker() override.
func (s *Signature) BaselineRequirementsChecker() *BaselineRequirementsChecker {
	return s.DefaultAdvancedSignature.BaselineRequirementsChecker().(*BaselineRequirementsChecker)
}

// CreateBaselineRequirementsChecker instantiates a BaselineRequirementsChecker according to the
// signature format. Port of the protected createBaselineRequirementsChecker(CertificateVerifier)
// override.
func (s *Signature) CreateBaselineRequirementsChecker(certificateVerifier validation.CertificateVerifier) validation.BaselineRequirementsCheckerContract {
	return NewPAdESBaselineRequirementsChecker(s, certificateVerifier)
}

// HasPKCS7Profile checks the presence of PKCS#7 corresponding SubFilter.
// Port of hasPKCS7Profile().
func (s *Signature) HasPKCS7Profile() bool {
	return s.BaselineRequirementsChecker().HasPKCS7Profile()
}

// HasPKCS7TProfile checks the presence of a signature-time-stamp.
// Port of hasPKCS7TProfile().
func (s *Signature) HasPKCS7TProfile() bool {
	return s.BaselineRequirementsChecker().HasPKCS7TProfile()
}

// HasPKCS7LTProfile checks the presence of a validation data. Port of hasPKCS7LTProfile().
func (s *Signature) HasPKCS7LTProfile() bool {
	return s.BaselineRequirementsChecker().HasPKCS7LTProfile()
}

// HasPKCS7LTAProfile checks the presence of an archive-time-stamp.
// Port of hasPKCS7LTAProfile().
func (s *Signature) HasPKCS7LTAProfile() bool {
	return s.BaselineRequirementsChecker().HasPKCS7LTAProfile()
}

// HasAProfile checks the presence of ArchiveTimeStamp element in the signature, what is the
// proof -A profile existence. Port of the hasAProfile() override.
func (s *Signature) HasAProfile() bool {
	return s.BaselineRequirementsChecker().HasExtendedAProfile()
}

// DssDictionary gets the last DSS dictionary for the signature. Port of getDssDictionary().
func (s *Signature) DssDictionary() PdfDssDict {
	return s.pdfSignatureRevision.DssDictionary()
}

// hasPKCS7SubFilter ports the private hasPKCS7SubFilter().
func (s *Signature) hasPKCS7SubFilter() bool {
	if s.pdfSignatureRevision != nil {
		subFilter := s.pdfSignatureRevision.PdfSigDictInfo().SubFilter()
		return PAdESConstantsSignaturePKCS7SubFilter == subFilter || PAdESConstantsSignaturePKCS7SHA1SubFilter == subFilter
	}
	return false
}

// PdfRevision retrieves a PdfRevision (PAdES) related to the current signature.
// Port of getPdfRevision().
func (s *Signature) PdfRevision() *PdfSignatureRevision {
	return s.pdfSignatureRevision
}

// PdfSignatureDictionary gets the PdfSignatureDictionary. Port of getPdfSignatureDictionary().
func (s *Signature) PdfSignatureDictionary() *PdfSignatureDictionary {
	return s.pdfSignatureRevision.PdfSigDictInfo()
}

// VRIKey returns the name of the related to the signature VRI dictionary.
// Port of getVRIKey().
func (s *Signature) VRIKey() string {
	if s.vriKey == "" {
		// By ETSI EN 319 142-1 V1.1.1, VRI dictionary's name is the base-16-encoded (uppercase)
		// SHA1 digest of the signature to which it applies
		digest, err := spi.DSSUtilsDigest(enumerations.DigestAlgorithmSHA1, s.PdfSignatureDictionary().Contents())
		if err != nil {
			panic(err)
		}
		s.vriKey = strings.ToUpper(utils.ToHex(digest))
	}
	return s.vriKey
}

// VRICreationTime returns a VRI creation time defined within 'TU' field of a corresponding
// /VRI dictionary. Port of getVRICreationTime().
func (s *Signature) VRICreationTime() *time.Time {
	dssDictionary := s.DssDictionary()
	if dssDictionary != nil {
		pdfVriDictTimestampSource := NewPdfVriDictSource(dssDictionary, s.VRIKey())
		return pdfVriDictTimestampSource.VRICreationTime()
	}
	return nil
}

// AddExternalTimestamp is not supported for PAdES. Port of addExternalTimestamp(TimestampToken).
func (s *Signature) AddExternalTimestamp(timestamp *validation.TimestampToken) {
	panic("The action is not supported for PAdES!")
}
