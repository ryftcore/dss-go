// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/signature/AdvancedSignature.java (DSS 6.5.RC1).
//
// FORWARD DEPENDENCY: EAA (Java spi.eaa.EAA), EvidenceRecord (Java spi.x509.evidencerecord.
// EvidenceRecord) and CertificateVerifier (Java spi.validation.CertificateVerifier) are all
// flattened into this same Go package by sibling chunks of phase 2b and are therefore
// referenced unqualified, following the precedent set by timestamp_source.go's EvidenceRecord
// forward dependency. Their assumed shapes are not needed here: every use below is either a
// parameter or return value passed through opaquely.
package validation

import (
	"time"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/model/scope"
	"github.com/utain/esig/dss/model/signature"
	"github.com/utain/esig/dss/model/x509/revocation"
	"github.com/utain/esig/dss/spi"
)

// AdvancedSignature provides an abstraction for an Advanced Electronic Signature. This eases the
// validation process. Every signature format: XAdES, CAdES and PAdES are treated the same.
//
// This is the central interface of the whole library: every format-specific signature
// implementation (CAdES/XAdES/JAdES/PAdES, ported in later phases) satisfies it. java.io.
// Serializable is dropped (no Go counterpart).
type AdvancedSignature interface {
	// model.IdentifierBasedObject supplies DSSID() model.Identifier, the Go counterpart of the
	// inherited IdentifierBasedObject#getDSSId(). Java additionally narrows the interface's own
	// getDSSId() to return the covariant SignatureIdentifier; Go has no covariant returns, so
	// callers needing the narrower type assert on the result, e.g.
	// sigID, ok := signature.DSSID().(*SignatureIdentifier), matching the convention documented
	// on model.IdentifierBasedObject and used throughout model.Token/TokenBase.
	model.IdentifierBasedObject

	// Filename returns the signature filename (useful for ASiC and multiple signature files).
	// Port of getFilename().
	Filename() string
	// SetFilename allows setting the signature filename (useful in case of ASiC).
	// Port of setFilename(String).
	SetFilename(filename string)

	// DetachedContents returns, in the case of the detached signature, the list of signed
	// contents. Port of getDetachedContents().
	DetachedContents() []model.DSSDocument
	// SetDetachedContents allows setting the signed contents in the case of the detached
	// signature. Port of setDetachedContents(List).
	SetDetachedContents(detachedContents []model.DSSDocument)

	// ContainerContents returns, in case of ASiC-S signature, a list of an archive container
	// documents. Port of getContainerContents().
	ContainerContents() []model.DSSDocument
	// SetContainerContents allows setting the archive container contents in the case of ASiC-S
	// signature. Port of setContainerContents(List).
	SetContainerContents(containerContents []model.DSSDocument)

	// ManifestFile returns a related ManifestFile in the case of ASiC-E signature.
	// Port of getManifestFile().
	ManifestFile() *model.ManifestFile
	// SetManifestFile allows setting a manifest file in the case of ASiC-E signature.
	// Port of setManifestFile(ManifestFile).
	SetManifestFile(manifestFile *model.ManifestFile)

	// SigningCertificateSource gets a signing certificate source, when provided.
	// Port of getSigningCertificateSource().
	SigningCertificateSource() spi.CertificateSource
	// SetSigningCertificateSource sets a certificate source which allows finding the signing
	// certificate by kid or certificate's digest.
	// Port of setSigningCertificateSource(CertificateSource).
	SetSigningCertificateSource(signingCertificateSource spi.CertificateSource)

	// SignatureForm specifies the format of the signature. Port of getSignatureForm().
	SignatureForm() enumerations.SignatureForm

	// SignatureAlgorithm retrieves the signature algorithm (or cipher) used for generating the
	// signature. Port of getSignatureAlgorithm().
	SignatureAlgorithm() enumerations.SignatureAlgorithm

	// EncryptionAlgorithm retrieves the encryption algorithm used for generating the signature.
	// Port of getEncryptionAlgorithm().
	EncryptionAlgorithm() enumerations.EncryptionAlgorithm

	// DigestAlgorithm retrieves the digest algorithm used for generating the signature.
	// Port of getDigestAlgorithm().
	DigestAlgorithm() enumerations.DigestAlgorithm

	// SigningTime returns the signing time included within the signature, or nil.
	// Port of getSigningTime(); the javadoc documents an explicit "or null" contract, hence the
	// pointer, matching the *time.Time convention used for other nullable Java Dates in this
	// port (e.g. model.BLevelParameters.SigningDate).
	SigningTime() *time.Time

	// CertificateSource gets a certificate source which contains ALL certificates embedded in
	// the signature. Port of getCertificateSource().
	CertificateSource() *spi.SignatureCertificateSource

	// CompleteCertificateSource gets a ListCertificateSource representing a merged source from
	// the signature certificate source and all included signature timestamp objects.
	// Port of getCompleteCertificateSource().
	CompleteCertificateSource() *spi.ListCertificateSource

	// CRLSource gets a CRL source which contains ALL CRLs embedded in the signature.
	// Port of getCRLSource().
	CRLSource() spi.OfflineRevocationSource[revocation.CRL]

	// OCSPSource gets an OCSP source which contains ALL OCSP responses embedded in the
	// signature. Port of getOCSPSource().
	OCSPSource() spi.OfflineRevocationSource[revocation.OCSP]

	// CompleteCRLSource gets a ListRevocationSource representing a merged source from the
	// signature CRL source and all included signature timestamp objects.
	// Port of getCompleteCRLSource().
	CompleteCRLSource() *spi.ListRevocationSource[revocation.CRL]

	// CompleteOCSPSource gets a ListRevocationSource representing a merged source from the
	// signature OCSP source and all included signature timestamp objects.
	// Port of getCompleteOCSPSource().
	CompleteOCSPSource() *spi.ListRevocationSource[revocation.OCSP]

	// TimestampSource gets a Signature Timestamp source which contains ALL timestamps embedded
	// in the signature. Port of getTimestampSource().
	TimestampSource() TimestampSource

	// CandidatesForSigningCertificate gets an object containing the signing certificate or
	// information indicating why it is impossible to extract it from the signature. If the
	// signing certificate is identified then it is cached and the subsequent calls to this
	// method return this cached value. This method never returns nil.
	// Port of getCandidatesForSigningCertificate().
	CandidatesForSigningCertificate() *spi.CandidatesForSigningCertificate

	// InitBaselineRequirementsChecker creates an offline copy of certificateVerifier and
	// instantiates a BaselineRequirementsChecker.
	// Port of initBaselineRequirementsChecker(CertificateVerifier).
	InitBaselineRequirementsChecker(certificateVerifier CertificateVerifier)

	// MasterSignature gets the master signature. Port of getMasterSignature().
	MasterSignature() AdvancedSignature
	// SetMasterSignature indicates the master signature: this means that the current signature
	// is a counter signature. Port of setMasterSignature(AdvancedSignature).
	SetMasterSignature(masterSignature AdvancedSignature)

	// EAA gets the EAA of an EAA issuing or key binding signature. Port of getEAA().
	EAA() EAA
	// SetEAA sets the EAA presentation of the EAA issuing or key binding signature.
	// Port of setEAA(EAA).
	SetEAA(eaa EAA)

	// IsCounterSignature checks if the current signature is a counter signature (i.e. has a
	// master signature). Port of isCounterSignature().
	IsCounterSignature() bool

	// IsKeyBindingSignature checks if the current signature is a key binding signature.
	// NOTE: used for EAA tokens. Port of isKeyBindingSignature().
	IsKeyBindingSignature() bool
	// SetKeyBindingSignature sets whether the current signature is a key binding signature.
	// NOTE: used for EAA tokens. Port of setKeyBindingSignature(boolean).
	SetKeyBindingSignature(keyBindingSignature bool)

	// SigningCertificateToken returns the signing certificate token or nil if there is no valid
	// signing certificate. Note that to determine the signing certificate the signature must be
	// validated: the method CheckSignatureIntegrity must be called.
	// Port of getSigningCertificateToken().
	SigningCertificateToken() *model.CertificateToken

	// CheckSignatureIntegrity verifies the signature integrity; checks if the signed content has
	// not been tampered with. In the case of a non-AdES signature not including the signing
	// certificate then the latter must be provided by calling
	// SetProvidedSigningCertificateToken. In the case of a detached signature the signed content
	// must be provided by calling SetProvidedSigningCertificateToken.
	// Port of checkSignatureIntegrity().
	CheckSignatureIntegrity()

	// SignatureCryptographicVerification gets the signature's cryptographic validation result.
	// Port of getSignatureCryptographicVerification().
	SignatureCryptographicVerification() *signature.SignatureCryptographicVerification

	// SignaturePolicy returns the Signature Policy OID from the signature.
	// Port of getSignaturePolicy().
	SignaturePolicy() *signature.SignaturePolicy

	// SignaturePolicyStore returns the Signature Policy Store from the signature.
	// Port of getSignaturePolicyStore().
	SignaturePolicyStore() *model.SignaturePolicyStore

	// SignatureProductionPlace returns information about the place where the signature was
	// generated. Port of getSignatureProductionPlace().
	SignatureProductionPlace() *signature.SignatureProductionPlace

	// CommitmentTypeIndications obtains the information concerning commitment type indication
	// linked to the signature. Port of getCommitmentTypeIndications().
	CommitmentTypeIndications() []*signature.CommitmentTypeIndication

	// ContentType returns the value of the signed attribute content-type.
	// Port of getContentType().
	ContentType() string

	// MimeType returns the value of the signed attribute mime-type. Port of getMimeType().
	MimeType() string

	// SignatureType returns the value of the signature type protected header (JAdES, CB-AdES).
	// Port of getSignatureType().
	SignatureType() string

	// SignerRoles returns the list of roles of the signer. Port of getSignerRoles().
	SignerRoles() []*signature.SignerRole

	// SignedAssertions returns the list of embedded signed assertions.
	// Port of getSignedAssertions().
	SignedAssertions() []*signature.SignerRole

	// ClaimedSignerRoles returns the claimed roles of the signer. Port of getClaimedSignerRoles().
	ClaimedSignerRoles() []*signature.SignerRole

	// CertifiedSignerRoles returns the certified roles of the signer.
	// Port of getCertifiedSignerRoles().
	CertifiedSignerRoles() []*signature.SignerRole

	// Certificates gets certificates embedded in the signature.
	// Port of getCertificates().
	Certificates() []*model.CertificateToken

	// ContentTimestamps returns the content timestamps. Port of getContentTimestamps().
	ContentTimestamps() []*TimestampToken

	// SignatureTimestamps returns the signature timestamps. Port of getSignatureTimestamps().
	SignatureTimestamps() []*TimestampToken

	// TimestampsX1 returns the time-stamp which is placed on the digital signature (XAdES
	// example: ds:SignatureValue element), the signature time-stamp(s) present in the AdES-T
	// form, the certification path references and the revocation status references.
	// Port of getTimestampsX1().
	TimestampsX1() []*TimestampToken

	// TimestampsX2 returns the time-stamp which is computed over the concatenation of
	// CompleteCertificateRefs and CompleteRevocationRefs elements (XAdES example).
	// Port of getTimestampsX2().
	TimestampsX2() []*TimestampToken

	// ArchiveTimestamps returns the archive timestamps. Port of getArchiveTimestamps().
	ArchiveTimestamps() []*TimestampToken

	// DocumentTimestamps returns a list of timestamps defined with the 'DocTimeStamp' type.
	// NOTE: applicable only for PAdES. Port of getDocumentTimestamps().
	DocumentTimestamps() []*TimestampToken

	// DetachedTimestamps returns a list of detached timestamps.
	// NOTE: used for ASiC with CAdES only. Port of getDetachedTimestamps().
	DetachedTimestamps() []*TimestampToken

	// AllTimestamps returns a list of all timestamps found in the signature.
	// Port of getAllTimestamps().
	AllTimestamps() []*TimestampToken

	// AddExternalTimestamp allows adding an external timestamp. The given timestamp must be
	// processed before. NOTE: supported only for CAdES signatures.
	// Port of addExternalTimestamp(TimestampToken).
	AddExternalTimestamp(timestamp *TimestampToken)

	// CounterSignatures returns a list of counter signatures applied to this signature.
	// Port of getCounterSignatures().
	CounterSignatures() []AdvancedSignature

	// EmbeddedEvidenceRecords returns a list of embedded evidence records.
	// Port of getEmbeddedEvidenceRecords().
	EmbeddedEvidenceRecords() []EvidenceRecord

	// AddExternalEvidenceRecord adds an evidence record covering the signature file.
	// Port of addExternalEvidenceRecord(EvidenceRecord).
	AddExternalEvidenceRecord(evidenceRecord EvidenceRecord)

	// DetachedEvidenceRecords returns a list of detached evidence records.
	// Port of getDetachedEvidenceRecords().
	DetachedEvidenceRecords() []EvidenceRecord

	// AllEvidenceRecords returns a list of all evidence records.
	// Port of getAllEvidenceRecords().
	AllEvidenceRecords() []EvidenceRecord

	// ID returns the DSS unique signature id. It allows unambiguously identifying each
	// signature. Port of getId().
	ID() string

	// DAIdentifier returns an identifier provided by the Driving Application (DA).
	// NOTE: used only for XAdES. Port of getDAIdentifier().
	DAIdentifier() string

	// DataFoundUpToLevel returns the signature level. Port of getDataFoundUpToLevel().
	DataFoundUpToLevel() enumerations.SignatureLevel

	// HasAdESProfile checks if the signature is conformant to the corresponding AdES profile.
	// Port of hasAdESProfile().
	HasAdESProfile() bool

	// HasBProfile checks if the signature is conformant to AdES-BASELINE-B level.
	// Port of hasBProfile().
	HasBProfile() bool

	// HasTProfile checks if the T-level is present in the signature. Port of hasTProfile().
	HasTProfile() bool

	// HasLTProfile checks if the LT-level is present in the signature. Port of hasLTProfile().
	HasLTProfile() bool

	// HasLTAProfile checks if the LTA-level is present in the signature.
	// Port of hasLTAProfile().
	HasLTAProfile() bool

	// HasBESProfile checks the presence of signing certificate covered by the signature, what is
	// the proof of the -BES profile existence. Port of hasBESProfile().
	HasBESProfile() bool

	// HasEPESProfile checks the presence of SignaturePolicyIdentifier element in the signature,
	// what is the proof of the -EPES profile existence. Port of hasEPESProfile().
	HasEPESProfile() bool

	// HasExtendedTProfile checks the presence of SignatureTimeStamp element in the signature,
	// what is the proof of the -T profile existence. Port of hasExtendedTProfile().
	HasExtendedTProfile() bool

	// HasCProfile checks the presence of CompleteCertificateRefs and CompleteRevocationRefs
	// segments in the signature, what is the proof of the -C profile existence.
	// Port of hasCProfile().
	HasCProfile() bool

	// HasXProfile checks the presence of SigAndRefsTimeStamp segment in the signature, what is
	// the proof of the -X profile existence. Port of hasXProfile().
	HasXProfile() bool

	// HasXLProfile checks the presence of CertificateValues/RevocationValues segment in the
	// signature, what is the proof of the -XL profile existence. Port of hasXLProfile().
	HasXLProfile() bool

	// HasAProfile checks the presence of ArchiveTimeStamp element in the signature, what is the
	// proof of the -A profile existence. Port of hasAProfile().
	HasAProfile() bool

	// HasERSProfile checks the presence of SealingEvidenceRecord element in the signature, what
	// is the proof of the -ERS profile existence. Port of hasERSProfile().
	HasERSProfile() bool

	// AreAllSelfSignedCertificates checks if all certificate chains present in the signature are
	// self-signed. Port of areAllSelfSignedCertificates().
	AreAllSelfSignedCertificates() bool

	// StructureValidationResult returns a message if the structure validation fails: a list of
	// error messages if validation fails, an empty list if structural validation succeeds.
	// Port of getStructureValidationResult().
	StructureValidationResult() []string

	// SignatureScopes returns a list of found SignatureScopes. Port of getSignatureScopes().
	SignatureScopes() []scope.SignatureScope

	// IsDocHashOnlyValidation returns true if the validation of the signature has been performed
	// only on Signer's Document Representation (SDR). (An SDR typically is built on a
	// cryptographic hash of the Signer's Document). Port of isDocHashOnlyValidation().
	IsDocHashOnlyValidation() bool

	// IsHashOnlyValidation returns true if the validation of the signature has been performed
	// only on Data To Be Signed Representation (DTBSR).
	//
	// EN 319 102-1 v1.1.1 (4.2.8 Data to be signed representation (DTBSR)): The DTBS
	// preparation component shall take the DTBSF and hash it according to the hash algorithm
	// specified in the cryptographic suite. The result of this process is the DTBSR, which is
	// then used to create the signature. NOTE: In order for the produced hash to be
	// representative of the DTBSF, the hashing function has the property that it is
	// computationally infeasible to find collisions for the expected signature lifetime. Should
	// the hash function become weak in the future, additional security measures, such as
	// applying time-stamp tokens, can be taken.
	// Port of isHashOnlyValidation().
	IsHashOnlyValidation() bool

	// SignatureValue returns the digital signature value. Port of getSignatureValue().
	SignatureValue() []byte

	// ReferenceValidations returns individual validation for each reference (XAdES, JAdES) or
	// for the message-imprint (CAdES). Port of getReferenceValidations().
	ReferenceValidations() []*model.ReferenceValidation

	// SignatureDigestReference returns a signature reference element as defined in TS 119 442 -
	// V1.1.1 - Electronic Signatures and Infrastructures (ESI), ch. 5.1.4.2.1.3 XML component.
	// Port of getSignatureDigestReference(DigestAlgorithm).
	SignatureDigestReference(digestAlgorithm enumerations.DigestAlgorithm) *signature.SignatureDigestReference

	// DataToBeSignedRepresentation returns the DTBSR, which is then used to create the
	// signature.
	//
	// TS 119 102-1 (4.2.8 Data to be signed representation (DTBSR)): The DTBS preparation
	// component shall take the DTBSF and hash it according to the hash algorithm specified in
	// the cryptographic suite.
	// Port of getDataToBeSignedRepresentation().
	DataToBeSignedRepresentation() model.Digest
}
