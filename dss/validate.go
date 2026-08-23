package dss

import (
	"fmt"
	"time"

	"github.com/ryftcore/dss-go/dss/policy"
	cryptojson "github.com/ryftcore/dss-go/dss/policy/crypto/json"
	cryptoxml "github.com/ryftcore/dss-go/dss/policy/crypto/xml"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/x509/aia"
	dssvalidation "github.com/ryftcore/dss-go/dss/validation"
	validationpolicy "github.com/ryftcore/dss-go/dss/validation/policy"
	"github.com/ryftcore/dss-go/dss/validation/reports"

	spivalidation "github.com/ryftcore/dss-go/dss/spi/validation"
)

// init performs the registration Java DSS gets from java.util.ServiceLoader:
// the ETSI validation policy implementation (dss-policy-jaxb) and the XML and
// JSON cryptographic suite implementations (dss-policy-crypto-xml,
// dss-policy-crypto-json). Without it every validation with the default policy
// fails, because the loader has no factory to consult.
//
// Order matters: the loader takes the FIRST registered factory as the default
// policy and the default cryptographic suite, and consults the rest in order
// for a supplied document. The XML cryptographic suite is registered first
// because it is the one the ETSI catalogue ships as XML.
func init() {
	validationpolicy.RegisterValidationPolicyFactory(policy.NewEtsiValidationPolicyFactory())
	validationpolicy.RegisterCryptographicSuiteFactory(cryptoxml.NewCryptographicSuiteXmlFactory())
	validationpolicy.RegisterCryptographicSuiteFactory(cryptojson.NewCryptographicSuiteJsonFactory())
}

// ValidateOptions configures [Validate]. The zero value is meaningful: it
// validates against the default ETSI policy, at the current time, with no
// trust anchors and no network access - which yields INDETERMINATE
// (NO_CERTIFICATE_CHAIN_FOUND, or a similar sub-indication) rather than a
// pass, because nothing anchors the certificate chain. Supply trust anchors to
// get a verdict that means something.
type ValidateOptions struct {
	// DetachedContents supplies the original documents a detached signature
	// covers. Without them a detached signature can only report that its
	// signed data is missing.
	DetachedContents []Document

	// TrustedCertificates are trust anchors given as individual certificates,
	// typically read with [LoadCertificate].
	TrustedCertificates []*CertificateToken

	// TrustedCertificateSources are whole trust stores. The trusted-list
	// certificate source the TSL validation job produces
	// ([github.com/ryftcore/dss-go/dss/spi/tsl.TrustedListsCertificateSource]) goes
	// here; that is what makes eIDAS qualification determination possible,
	// since the qualifiers come from the trusted lists.
	TrustedCertificateSources []CertificateSource

	// CertificateVerifier takes over completely: when set, the facade passes
	// it to the validator untouched and ignores TrustedCertificates,
	// TrustedCertificateSources and EnableAIA. Use it to add CRL/OCSP sources,
	// a revocation-data verifier, or a different alert policy.
	CertificateVerifier CertificateVerifier

	// EnableAIA lets the validator download missing issuer certificates over
	// the Authority Information Access extension. Off by default: a library
	// call should not reach the network unasked. Ignored when
	// CertificateVerifier is set.
	EnableAIA bool

	// Policy is a custom validation policy document (the DSS constraint XML).
	// Defaults to the ETSI policy shipped with the library.
	Policy Document

	// CryptographicSuite is a custom cryptographic suite catalogue (ETSI TS
	// 119 312 XML, or its JSON flavour) constraining algorithms and key sizes
	// over time. Requires Policy to be set as well, mirroring the underlying
	// two-document entry point.
	CryptographicSuite Document

	// ValidationTime pins the moment the validation is performed at, which
	// decides whether certificates were valid and revocation data fresh.
	// Defaults to now.
	ValidationTime *time.Time

	// Level is how far the validation process is taken. Defaults to
	// [ValidationArchivalData], the fullest process.
	Level ValidationLevel

	// TokenExtractionStrategy selects which tokens are embedded, base64
	// encoded, into the diagnostic data. Defaults to embedding none.
	TokenExtractionStrategy TokenExtractionStrategy

	// IncludeSemantics adds the human-readable meaning of each Indication and
	// SubIndication to the reports.
	IncludeSemantics bool

	// Locale is the language of the report messages, as a language tag such
	// as "en" or "fr". Defaults to the library default.
	Locale string

	// PasswordProtection is the password that opens an encrypted
	// (password-protected) PDF, without which its signatures cannot be read.
	// Either the user or the owner password is accepted; it is UTF-8 text,
	// hashed as pdfbox hashes it (see [SignOptions.PasswordProtection]).
	// Applies to PDF documents only: [Validate] returns
	// [ErrPasswordProtectionNotApplicable] when it is set for a document
	// detected as anything else. This is the facade's name for
	// PDFDocumentValidator.setPasswordProtection.
	PasswordProtection []byte
}

// certificateVerifier builds the verifier the validator will use.
func (o ValidateOptions) certificateVerifier() CertificateVerifier {
	if o.CertificateVerifier != nil {
		return o.CertificateVerifier
	}
	verifier := spivalidation.NewCommonCertificateVerifierSimple(true)
	sources := make([]spi.CertificateSource, 0, len(o.TrustedCertificateSources)+1)
	sources = append(sources, o.TrustedCertificateSources...)
	if len(o.TrustedCertificates) > 0 {
		sources = append(sources, TrustStore(o.TrustedCertificates...))
	}
	if len(sources) > 0 {
		verifier.SetTrustedCertSources(sources...)
	}
	if o.EnableAIA {
		verifier.SetAIASource(aia.NewDefaultSource())
	}
	return verifier
}

// level returns the configured validation level or the default.
func (o ValidateOptions) level() ValidationLevel {
	if o.Level == "" {
		return ValidationArchivalData
	}
	return o.Level
}

// Validate validates every signature, time-stamp and evidence record in doc
// and returns the DSS reports. The document format is detected automatically -
// CMS, XML, PDF, JWS or ASiC - so the caller does not name it.
//
// The returned error reports a failure to RUN the validation (an unreadable
// document, an unsupported format, a broken policy). A signature that does not
// verify is not an error: it is a verdict, carried by the reports. Check
// [Reports.Valid] or read [Reports.Verdicts].
func Validate(doc Document, opts ValidateOptions) (*Reports, error) {
	if doc == nil {
		return nil, ErrNoDocument
	}

	var result *reports.Reports
	err := recovered("validate", func() error {
		validator, err := dssvalidation.SignedDocumentValidatorFromDocument(doc)
		if err != nil {
			return err
		}
		if len(opts.PasswordProtection) > 0 {
			protected, ok := validator.(passwordProtectionSetter)
			if !ok {
				return fmt.Errorf("%w (the document was not detected as a PDF)", ErrPasswordProtectionNotApplicable)
			}
			protected.SetPasswordProtection(opts.PasswordProtection)
		}
		validator.SetCertificateVerifier(opts.certificateVerifier())
		validator.SetValidationLevel(opts.level())
		if len(opts.DetachedContents) > 0 {
			validator.SetDetachedContents(opts.DetachedContents)
		}
		if opts.ValidationTime != nil {
			validator.SetValidationTime(*opts.ValidationTime)
		}
		if opts.TokenExtractionStrategy != "" {
			validator.SetTokenExtractionStrategy(opts.TokenExtractionStrategy)
		}
		if opts.IncludeSemantics {
			validator.SetIncludeSemantics(true)
		}
		if opts.Locale != "" {
			validator.SetLocale(opts.Locale)
		}

		switch {
		case opts.Policy != nil && opts.CryptographicSuite != nil:
			result, err = validator.ValidateDocumentWithPolicyAndCryptographicSuiteDocument(
				opts.Policy, opts.CryptographicSuite)
		case opts.Policy != nil:
			result, err = validator.ValidateDocumentWithPolicyDocument(opts.Policy)
		case opts.CryptographicSuite != nil:
			return fmt.Errorf("dss: ValidateOptions.CryptographicSuite requires ValidateOptions.Policy to be set as well")
		default:
			result, err = validator.ValidateDocument()
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	return &Reports{Reports: result}, nil
}

// passwordProtectionSetter is the SetPasswordProtection of
// pades.PDFDocumentValidator, the one validator that reads encrypted
// documents. [Validate] detects the document format itself, so it looks for
// the method on whatever validator the detection produced rather than
// naming the PDF validator's type.
type passwordProtectionSetter interface {
	SetPasswordProtection(passwordProtection []byte)
}

// LoadCertificate reads an X.509 certificate, DER or PEM encoded, from the
// file at path. Delegates to
// [github.com/ryftcore/dss-go/dss/spi.DSSUtilsLoadCertificate].
func LoadCertificate(path string) (*CertificateToken, error) {
	return spi.DSSUtilsLoadCertificate(path)
}

// LoadCertificateBytes reads an X.509 certificate, DER or PEM encoded, from
// memory. Delegates to
// [github.com/ryftcore/dss-go/dss/spi.DSSUtilsLoadCertificateFromBinary].
func LoadCertificateBytes(der []byte) (*CertificateToken, error) {
	return spi.DSSUtilsLoadCertificateFromBinary(der)
}

// TrustStore builds a trusted [CertificateSource] holding the given
// certificates, ready for ValidateOptions.TrustedCertificateSources or for a
// [CertificateVerifier] of your own. Certificates in a trusted source are
// treated as trust anchors: a chain that reaches one is anchored.
//
// A trust store built this way carries no trusted-list information, so it
// anchors chains but cannot make a signature qualified. For that, use the
// trusted-list certificate source produced by the TSL validation job of
// [github.com/ryftcore/dss-go/dss/validation/job].
func TrustStore(certificates ...*CertificateToken) CertificateSource {
	source := spi.NewCommonTrustedCertificateSource()
	for _, certificate := range certificates {
		if certificate != nil {
			source.AddCertificate(certificate)
		}
	}
	return source
}
