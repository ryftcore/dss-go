package dss

import (
	"errors"
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	spivalidation "github.com/ryftcore/dss-go/dss/spi/validation"

	// Blank-imported for their init()-registered document validator and
	// document analyzer factories, which Validate's format auto-detection
	// dispatches through. Java discovers the same set with
	// java.util.ServiceLoader; see the package doc.
	_ "github.com/ryftcore/dss-go/dss/asic/cades"
	_ "github.com/ryftcore/dss-go/dss/asic/xades"
	_ "github.com/ryftcore/dss-go/dss/cades"
	_ "github.com/ryftcore/dss-go/dss/jades"
	_ "github.com/ryftcore/dss-go/dss/pades"
	_ "github.com/ryftcore/dss-go/dss/xades"
)

// Aliases of the ported types the facade takes and returns. They are aliases,
// not new types, so a value obtained from the underlying packages can be
// passed straight into the facade and vice versa.
type (
	// Document is a document handed to or produced by the library. See
	// [model.DSSDocument] for the implementations (in-memory, file-backed,
	// digest-only).
	Document = model.DSSDocument

	// DigestAlgorithm is a message digest algorithm, for example
	// [DigestSHA256].
	DigestAlgorithm = enumerations.DigestAlgorithm

	// SignaturePackaging says how a signature relates to the data it covers:
	// [PackagingEnveloped], [PackagingEnveloping] or [PackagingDetached].
	SignaturePackaging = enumerations.SignaturePackaging

	// SignatureLevel is the ETSI signature level of a signature. Its value is
	// the upstream enum name, "XAdES_BASELINE_LTA"; its String method renders
	// the dash spelling, "XAdES-BASELINE-LTA", exactly as Java's toString()
	// does, which is also what the reports show. Use [Format.BaselineLevel]
	// to obtain one from a [Format] and a [Level].
	SignatureLevel = enumerations.SignatureLevel

	// ContainerType is the ASiC container flavour: [ContainerASiCS] or
	// [ContainerASiCE].
	ContainerType = enumerations.ASiCContainerType

	// JWSSerializationType is the JAdES serialization: compact, JSON or
	// flattened JSON.
	JWSSerializationType = enumerations.JWSSerializationType

	// ValidationLevel is how far the validation process is taken; the default
	// is [ValidationArchivalData].
	ValidationLevel = enumerations.ValidationLevel

	// TokenExtractionStrategy selects which tokens are embedded, base64
	// encoded, into the diagnostic data.
	TokenExtractionStrategy = enumerations.TokenExtractionStrategy

	// Indication is the top-level validation verdict of EN 319 102-1:
	// [IndicationTotalPassed], [IndicationIndeterminate] or
	// [IndicationTotalFailed].
	Indication = enumerations.Indication

	// SubIndication refines an [Indication] with the reason behind it.
	SubIndication = enumerations.SubIndication

	// SignatureQualification is the eIDAS qualification determined for a
	// signature, for example "QESig". It is "NA" when nothing in the trusted
	// lists applies.
	SignatureQualification = enumerations.SignatureQualification

	// TimestampQualification is the eIDAS qualification determined for a
	// time-stamp token, for example "QTSA".
	TimestampQualification = enumerations.TimestampQualification

	// CertificateToken is an X.509 certificate as the library models it.
	CertificateToken = model.CertificateToken

	// CertificateSource is a source of certificates - a trust store, a
	// trusted-list certificate source produced by the TSL job, and so on.
	CertificateSource = spi.CertificateSource

	// CertificateVerifier carries every external source the validation and
	// augmentation processes consult: trust anchors, CRL/OCSP sources, AIA,
	// and the alert policy. Build one with
	// [spivalidation.NewCommonCertificateVerifier] when the facade options are
	// not enough.
	CertificateVerifier = spivalidation.CertificateVerifier

	// TSPSource issues RFC 3161 time-stamp tokens. Required for levels T, LT
	// and LTA; see the package doc on what the port does and does not ship.
	TSPSource = spivalidation.TSPSource
)

// Frequently used enumeration values, re-exported so that the common paths need
// no import of the enumerations package.
const (
	DigestSHA256 = enumerations.DigestAlgorithmSHA256
	DigestSHA384 = enumerations.DigestAlgorithmSHA384
	DigestSHA512 = enumerations.DigestAlgorithmSHA512

	PackagingEnveloped  = enumerations.SignaturePackagingEnveloped
	PackagingEnveloping = enumerations.SignaturePackagingEnveloping
	PackagingDetached   = enumerations.SignaturePackagingDetached

	ContainerASiCS = enumerations.ASiCContainerTypeASiCS
	ContainerASiCE = enumerations.ASiCContainerTypeASiCE

	JWSCompact       = enumerations.JWSSerializationTypeCompactSerialization
	JWSJSON          = enumerations.JWSSerializationTypeJSONSerialization
	JWSFlattenedJSON = enumerations.JWSSerializationTypeFlattenedJSONSerialization

	ValidationBasicSignatures = enumerations.ValidationLevelBasicSignatures
	ValidationTimestamps      = enumerations.ValidationLevelTimestamps
	ValidationLongTermData    = enumerations.ValidationLevelLongTermData
	ValidationArchivalData    = enumerations.ValidationLevelArchivalData

	IndicationTotalPassed      = enumerations.IndicationTotalPassed
	IndicationIndeterminate    = enumerations.IndicationIndeterminate
	IndicationTotalFailed      = enumerations.IndicationTotalFailed
	IndicationNoSignatureFound = enumerations.IndicationNoSignatureFound
)

// Errors the facade returns in addition to whatever the underlying services
// report. Match them with errors.Is.
var (
	// ErrUnsupportedFormat is returned for a [Format] value the facade does
	// not know.
	ErrUnsupportedFormat = errors.New("dss: unsupported signature format")

	// ErrUnsupportedLevel is returned for a [Level] value the facade does not
	// know.
	ErrUnsupportedLevel = errors.New("dss: unsupported signature level")

	// ErrTSPSourceRequired is returned when levels T, LT or LTA are requested
	// without a TSPSource to obtain the time-stamp tokens from.
	ErrTSPSourceRequired = errors.New("dss: a TSPSource is required for levels T, LT and LTA")

	// ErrMultipleDocuments is returned by [SignMultiple] when more than one
	// document is passed to a format that signs a single document.
	ErrMultipleDocuments = errors.New("dss: only the ASiC formats can sign several documents at once")

	// ErrNoDocument is returned when no document to sign was provided.
	ErrNoDocument = errors.New("dss: at least one document to sign is required")

	// ErrPasswordProtectionNotApplicable is returned when a PasswordProtection
	// is given for anything but an encrypted PDF: a [Format] other than
	// [FormatPAdES] in [SignOptions] or [ExtendOptions], or a document
	// [Validate] detected as something other than a PDF.
	ErrPasswordProtectionNotApplicable = errors.New("dss: a PasswordProtection applies to PDF documents (PAdES) only")
)

// Format identifies a signature format family. It selects which of the ported
// signature services the facade delegates to.
type Format string

// The signature formats the facade can create, extend and validate.
const (
	// FormatCAdES is CMS Advanced Electronic Signatures (ETSI EN 319 122),
	// the format for binary content; delegates to
	// [github.com/ryftcore/dss-go/dss/cades].
	FormatCAdES Format = "CAdES"

	// FormatXAdES is XML Advanced Electronic Signatures (ETSI EN 319 132);
	// delegates to [github.com/ryftcore/dss-go/dss/xades].
	FormatXAdES Format = "XAdES"

	// FormatPAdES is PDF Advanced Electronic Signatures (ETSI EN 319 142);
	// delegates to [github.com/ryftcore/dss-go/dss/pades].
	FormatPAdES Format = "PAdES"

	// FormatJAdES is JSON Advanced Electronic Signatures (ETSI TS 119 182);
	// delegates to [github.com/ryftcore/dss-go/dss/jades].
	FormatJAdES Format = "JAdES"

	// FormatASiCWithCAdES is an ASiC container (ETSI EN 319 162) holding
	// CAdES signatures; delegates to
	// [github.com/ryftcore/dss-go/dss/asic/cades].
	FormatASiCWithCAdES Format = "ASiC-CAdES"

	// FormatASiCWithXAdES is an ASiC container (ETSI EN 319 162) holding
	// XAdES signatures; delegates to
	// [github.com/ryftcore/dss-go/dss/asic/xades].
	FormatASiCWithXAdES Format = "ASiC-XAdES"
)

// String returns the format name.
func (f Format) String() string { return string(f) }

// IsContainer reports whether the format produces an ASiC container, which is
// the case exactly for [FormatASiCWithCAdES] and [FormatASiCWithXAdES]. Only
// container formats accept more than one document in [SignMultiple].
func (f Format) IsContainer() bool {
	return f == FormatASiCWithCAdES || f == FormatASiCWithXAdES
}

// Level is a baseline signature level as defined by the ETSI baseline
// profiles. Higher levels build on the lower ones.
type Level string

// The four baseline levels.
const (
	// LevelB is the baseline B-B level: the signature itself, with the signed
	// attributes the profile mandates. No time-stamp, no revocation data.
	LevelB Level = "B"

	// LevelT is the baseline B-T level: B plus a trusted time-stamp over the
	// signature value, proving the signature existed at that time. Requires a
	// [TSPSource].
	LevelT Level = "T"

	// LevelLT is the baseline B-LT level: T plus the certificates and
	// revocation data a verifier needs long after the fact. Requires a
	// [TSPSource] and revocation sources on the [CertificateVerifier].
	LevelLT Level = "LT"

	// LevelLTA is the baseline B-LTA level: LT plus an archival time-stamp,
	// which can be renewed to keep the signature verifiable past the
	// cryptographic lifetime of the algorithms used. Requires a [TSPSource].
	LevelLTA Level = "LTA"
)

// String returns the level name.
func (l Level) String() string { return string(l) }

// NeedsTimestamp reports whether the level requires a [TSPSource], which is
// the case for every level above [LevelB].
func (l Level) NeedsTimestamp() bool { return l != LevelB }

// baselineLevels maps a facade (format, level) pair onto the upstream
// enumerations.SignatureLevel the ported services take. The two ASiC formats
// carry the levels of the signature format inside the container, exactly as
// the ASiC services expect.
var baselineLevels = map[Format]map[Level]SignatureLevel{
	FormatCAdES: {
		LevelB:   enumerations.SignatureLevelCAdESBaselineB,
		LevelT:   enumerations.SignatureLevelCAdESBaselineT,
		LevelLT:  enumerations.SignatureLevelCAdESBaselineLT,
		LevelLTA: enumerations.SignatureLevelCAdESBaselineLTA,
	},
	FormatXAdES: {
		LevelB:   enumerations.SignatureLevelXAdESBaselineB,
		LevelT:   enumerations.SignatureLevelXAdESBaselineT,
		LevelLT:  enumerations.SignatureLevelXAdESBaselineLT,
		LevelLTA: enumerations.SignatureLevelXAdESBaselineLTA,
	},
	FormatPAdES: {
		LevelB:   enumerations.SignatureLevelPAdESBaselineB,
		LevelT:   enumerations.SignatureLevelPAdESBaselineT,
		LevelLT:  enumerations.SignatureLevelPAdESBaselineLT,
		LevelLTA: enumerations.SignatureLevelPAdESBaselineLTA,
	},
	FormatJAdES: {
		LevelB:   enumerations.SignatureLevelJAdESBaselineB,
		LevelT:   enumerations.SignatureLevelJAdESBaselineT,
		LevelLT:  enumerations.SignatureLevelJAdESBaselineLT,
		LevelLTA: enumerations.SignatureLevelJAdESBaselineLTA,
	},
	FormatASiCWithCAdES: {
		LevelB:   enumerations.SignatureLevelCAdESBaselineB,
		LevelT:   enumerations.SignatureLevelCAdESBaselineT,
		LevelLT:  enumerations.SignatureLevelCAdESBaselineLT,
		LevelLTA: enumerations.SignatureLevelCAdESBaselineLTA,
	},
	FormatASiCWithXAdES: {
		LevelB:   enumerations.SignatureLevelXAdESBaselineB,
		LevelT:   enumerations.SignatureLevelXAdESBaselineT,
		LevelLT:  enumerations.SignatureLevelXAdESBaselineLT,
		LevelLTA: enumerations.SignatureLevelXAdESBaselineLTA,
	},
}

// BaselineLevel returns the upstream [SignatureLevel] the format and level
// pair denotes, for example FormatXAdES.BaselineLevel(LevelLTA) is
// "XAdES_BASELINE_LTA". It reports [ErrUnsupportedFormat] or
// [ErrUnsupportedLevel] for values outside the two enumerations.
func (f Format) BaselineLevel(l Level) (SignatureLevel, error) {
	levels, ok := baselineLevels[f]
	if !ok {
		return "", fmt.Errorf("%w: %q", ErrUnsupportedFormat, string(f))
	}
	level, ok := levels[l]
	if !ok {
		return "", fmt.Errorf("%w: %q", ErrUnsupportedLevel, string(l))
	}
	return level, nil
}

// defaultPackaging is the signature packaging the facade applies when
// SignOptions.Packaging is left empty. The ported services have no default of
// their own (Java DSS leaves the field null and asserts on it), so the facade
// picks the usual one per format: an enveloping CMS or JWS, and an enveloped
// XML signature. PAdES and ASiC do not use the field.
func defaultPackaging(f Format) SignaturePackaging {
	switch f {
	case FormatXAdES:
		return PackagingEnveloped
	case FormatCAdES, FormatJAdES:
		return PackagingEnveloping
	default:
		return ""
	}
}

// OpenDocument reads the file at path as a [Document]. The content is read
// lazily, on demand, not slurped into memory. Delegates to
// [model.NewFileDocument].
func OpenDocument(path string) (Document, error) {
	return model.NewFileDocument(path)
}

// NewDocument wraps content as an in-memory [Document] carrying the given
// name. The name matters: it ends up in the ASiC container entries and in the
// report's document filename. Delegates to
// [model.NewInMemoryDocumentWithName].
func NewDocument(name string, content []byte) Document {
	return model.NewInMemoryDocumentWithName(content, name)
}

// recovered runs fn, turning a panic raised by the ported services - which
// follow Java DSS and throw unchecked exceptions - into an error. A panic
// value that is already an error is wrapped, so errors.As and errors.Is
// against the port's error types keep working.
func recovered(op string, fn func() error) (err error) {
	defer func() {
		if p := recover(); p != nil {
			if cause, ok := p.(error); ok {
				err = fmt.Errorf("dss: %s: %w", op, cause)
				return
			}
			err = fmt.Errorf("dss: %s: %v", op, p)
		}
	}()
	return fn()
}
