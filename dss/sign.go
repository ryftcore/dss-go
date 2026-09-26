package dss

import (
	"fmt"
	"time"

	asiccades "github.com/ryftcore/dss-go/dss/asic/cades"
	asicxades "github.com/ryftcore/dss-go/dss/asic/xades"
	"github.com/ryftcore/dss-go/dss/cades"
	"github.com/ryftcore/dss-go/dss/document"
	"github.com/ryftcore/dss-go/dss/jades"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/pades"
	spivalidation "github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/xades"
)

// SignOptions configures [Sign], [SignMultiple] and - for the fields it shares
// - [Extend]. Only Format and Level are required; everything else has a
// documented default.
type SignOptions struct {
	// Format selects the signature format family and therefore which ported
	// service does the work. Required.
	Format Format

	// Level is the baseline level to produce. Required. Levels above
	// [LevelB] need TSPSource; LT and LTA additionally need revocation data,
	// which means a CertificateVerifier carrying CRL and/or OCSP sources.
	Level Level

	// DigestAlgorithm is the digest used for the signature and for the
	// references it covers. Defaults to [DigestSHA256].
	DigestAlgorithm DigestAlgorithm

	// Packaging says how the signature relates to the signed data. It applies
	// to CAdES, XAdES and JAdES only; the facade never passes it on for
	// PAdES (whose signature is always embedded in the PDF) or for the ASiC
	// containers (whose services decide it themselves). Defaults to
	// [PackagingEnveloped] for XAdES and [PackagingEnveloping] for CAdES and
	// JAdES.
	Packaging SignaturePackaging

	// DetachedContents supplies the documents a detached signature covers.
	// Required when signing a document that already carries a detached
	// signature being counter-parallel-signed; for a first detached signature
	// the document being signed is the content and this can stay empty.
	DetachedContents []Document

	// ContainerType selects the ASiC flavour, [ContainerASiCS] (one signed
	// document) or [ContainerASiCE] (several). Applies to the two ASiC
	// formats only; defaults to [ContainerASiCE] for several documents and
	// [ContainerASiCS] for one.
	ContainerType ContainerType

	// JWSSerialization selects the JAdES serialization. Applies to JAdES
	// only; defaults to [JWSCompact]. A detached or a multi-signature JAdES
	// needs [JWSJSON] or [JWSFlattenedJSON].
	JWSSerialization JWSSerializationType

	// SigningTime pins the claimed signing time placed in the signed
	// attributes. Defaults to the moment the signature is built.
	SigningTime *time.Time

	// TSPSource issues the RFC 3161 time-stamp tokens levels T, LT and LTA
	// need. The port ships no HTTP TSA client; see the package doc.
	TSPSource TSPSource

	// CertificateVerifier carries the trust anchors, revocation sources and
	// alert policy the signing services consult - which matters from [LevelLT]
	// upwards, where the signature has to embed validation data. Defaults to
	// a [spivalidation.NewCommonCertificateVerifierSimple] with no network
	// access, which is enough for [LevelB] and [LevelT].
	CertificateVerifier CertificateVerifier

	// PasswordProtection is the password that opens an encrypted
	// (password-protected) PDF, so that it can be signed: the signature is
	// appended as an incremental update encrypted under the document's own
	// security handler, and the result stays protected by the same password.
	// The owner password bypasses the document's permission (/P) check; the
	// user password opens the document, but signing then also needs /P to
	// allow creating or filling a signature field - otherwise it is refused,
	// as upstream's PdfPermissionsChecker refuses it. The password is UTF-8
	// text, the characters of upstream's char[], and is hashed the way
	// pdfbox hashes it for the document's encryption revision (ISO-8859-1
	// for RC4 and AES-128 documents, UTF-8 for AES-256, after SASLprep for
	// the current /R 6 revision), so a document that opens with a given
	// password in Java DSS opens with the same one here. Applies to
	// [FormatPAdES] only - every other format
	// returns [ErrPasswordProtectionNotApplicable] when it is set. Leave it
	// empty for a PDF that is not encrypted; one given for such a PDF is
	// ignored, as pdfbox ignores it. This is the facade's name for
	// PAdESSignatureParameters.setPasswordProtection; encrypting a PDF that
	// is not yet encrypted is not something the library does, upstream or
	// here.
	PasswordProtection []byte
}

// digestAlgorithm returns the configured digest or the facade default.
func (o SignOptions) digestAlgorithm() DigestAlgorithm {
	if o.DigestAlgorithm == "" {
		return DigestSHA256
	}
	return o.DigestAlgorithm
}

// packaging returns the configured packaging or the per-format default, and
// an empty value for the formats that do not take one.
func (o SignOptions) packaging() SignaturePackaging {
	if defaultPackaging(o.Format) == "" {
		return ""
	}
	if o.Packaging == "" {
		return defaultPackaging(o.Format)
	}
	return o.Packaging
}

// certificateVerifier returns the configured verifier or an offline default.
func (o SignOptions) certificateVerifier() CertificateVerifier {
	if o.CertificateVerifier != nil {
		return o.CertificateVerifier
	}
	return spivalidation.NewCommonCertificateVerifierSimple(true)
}

// containerType returns the configured ASiC flavour or the default for the
// number of documents being signed.
func (o SignOptions) containerType(documentCount int) ContainerType {
	if o.ContainerType != "" {
		return o.ContainerType
	}
	if documentCount > 1 {
		return ContainerASiCE
	}
	return ContainerASiCS
}

// jwsSerialization returns the configured JAdES serialization or the default.
func (o SignOptions) jwsSerialization() JWSSerializationType {
	if o.JWSSerialization == "" {
		return JWSCompact
	}
	return o.JWSSerialization
}

// validate checks the options that cannot produce a meaningful signature.
func (o SignOptions) validate() error {
	if _, err := o.Format.BaselineLevel(o.Level); err != nil {
		return err
	}
	if o.Level.NeedsTimestamp() && o.TSPSource == nil {
		return fmt.Errorf("%w (level %s requested)", ErrTSPSourceRequired, o.Level)
	}
	if len(o.PasswordProtection) > 0 && o.Format != FormatPAdES {
		return fmt.Errorf("%w (format %s)", ErrPasswordProtectionNotApplicable, o.Format)
	}
	return nil
}

// Sign signs doc and returns the signed document. It runs the two-step DSS
// signing flow: the selected service computes the data to be signed, the
// [Signer]'s token produces the signature value over it, and the service
// embeds that value in the final signature.
//
// The returned document is not written anywhere; call its Save method or read
// it as a stream.
func Sign(doc Document, signer *Signer, opts SignOptions) (Document, error) {
	if doc == nil {
		return nil, ErrNoDocument
	}
	return SignMultiple([]Document{doc}, signer, opts)
}

// SignMultiple signs several documents into one ASiC container. Only
// [FormatASiCWithCAdES] and [FormatASiCWithXAdES] can cover more than one
// document; every other format returns [ErrMultipleDocuments] when given more
// than one, and behaves exactly like [Sign] when given one.
func SignMultiple(docs []Document, signer *Signer, opts SignOptions) (Document, error) {
	if len(docs) == 0 {
		return nil, ErrNoDocument
	}
	if signer == nil || signer.conn == nil || signer.key == nil {
		return nil, fmt.Errorf("dss: a Signer is required (use OpenPKCS12, OpenPKCS12Bytes or NewSigner)")
	}
	if err := opts.validate(); err != nil {
		return nil, err
	}
	if len(docs) > 1 && !opts.Format.IsContainer() {
		return nil, fmt.Errorf("%w (format %s, %d documents)", ErrMultipleDocuments, opts.Format, len(docs))
	}
	level, err := opts.Format.BaselineLevel(opts.Level)
	if err != nil {
		return nil, err
	}

	var signed Document
	err = recovered("sign", func() error {
		var runErr error
		switch opts.Format {
		case FormatCAdES:
			parameters := cades.NewSignatureParameters()
			applyCommonParameters(parameters, opts, level, signer)
			service := cades.NewService(opts.certificateVerifier())
			applyTSPSource(service, opts)
			signed, runErr = signOne(service, parameters, docs[0], signer, opts.digestAlgorithm())

		case FormatXAdES:
			parameters := xades.NewSignatureParameters()
			applyCommonParameters(parameters, opts, level, signer)
			service := xades.NewService(opts.certificateVerifier())
			applyTSPSource(service, opts)
			signed, runErr = signOne(service, parameters, docs[0], signer, opts.digestAlgorithm())

		case FormatPAdES:
			parameters := pades.NewSignatureParameters()
			applyCommonParameters(parameters, opts, level, signer)
			parameters.SetPasswordProtection(opts.PasswordProtection)
			service := pades.NewService(opts.certificateVerifier())
			applyTSPSource(service, opts)
			signed, runErr = signOne(service, parameters, docs[0], signer, opts.digestAlgorithm())

		case FormatJAdES:
			parameters := jades.NewSignatureParameters()
			applyCommonParameters(parameters, opts, level, signer)
			parameters.SetJwsSerializationType(opts.jwsSerialization())
			service := jades.NewService(opts.certificateVerifier())
			applyTSPSource(service, opts)
			signed, runErr = signOne(service, parameters, docs[0], signer, opts.digestAlgorithm())

		case FormatASiCWithCAdES:
			parameters := asiccades.NewASiCWithCAdESSignatureParameters()
			applyCommonParameters(parameters, opts, level, signer)
			parameters.ASiC().SetContainerType(opts.containerType(len(docs)))
			service := asiccades.NewASiCWithCAdESService(opts.certificateVerifier())
			applyTSPSource(service, opts)
			signed, runErr = signMany(service.MultipleDocumentsService(), parameters, docs, signer, opts.digestAlgorithm())

		case FormatASiCWithXAdES:
			parameters := asicxades.NewASiCWithXAdESSignatureParameters()
			applyCommonParameters(parameters, opts, level, signer)
			parameters.ASiC().SetContainerType(opts.containerType(len(docs)))
			service := asicxades.NewASiCWithXAdESService(opts.certificateVerifier())
			applyTSPSource(service, opts)
			signed, runErr = signMany(service.MultipleDocumentsService(), parameters, docs, signer, opts.digestAlgorithm())

		default:
			return fmt.Errorf("%w: %q", ErrUnsupportedFormat, string(opts.Format))
		}
		return runErr
	})
	if err != nil {
		return nil, err
	}
	return signed, nil
}

// ExtendOptions configures [Extend]. Level is the level to raise the existing
// signatures to; the other fields carry the same meaning as in [SignOptions].
type ExtendOptions struct {
	// Format is the format family of the signatures in the document.
	// Required: extension has no auto-detection, because the caller has to
	// name the target level in that format's terms anyway.
	Format Format

	// Level is the level to reach. Required, and it must be above the level
	// the signatures currently have.
	Level Level

	// TSPSource issues the time-stamp tokens. Required, since every level
	// reachable by extension is above [LevelB].
	TSPSource TSPSource

	// CertificateVerifier carries the revocation sources and trust anchors
	// [LevelLT] and [LevelLTA] need to collect validation data, and the alert
	// policy that decides what an incomplete collection does. Defaults to an
	// offline [spivalidation.NewCommonCertificateVerifierSimple], which is
	// enough only for [LevelT].
	CertificateVerifier CertificateVerifier

	// DetachedContents supplies the documents a detached signature covers,
	// without which its validation data cannot be collected.
	DetachedContents []Document

	// PasswordProtection is the password that opens an encrypted
	// (password-protected) PDF whose signatures are being extended. It
	// carries the same meaning as [SignOptions.PasswordProtection]: PAdES
	// only, [ErrPasswordProtectionNotApplicable] for every other format.
	PasswordProtection []byte
}

// certificateVerifier returns the configured verifier or an offline default.
func (o ExtendOptions) certificateVerifier() CertificateVerifier {
	if o.CertificateVerifier != nil {
		return o.CertificateVerifier
	}
	return spivalidation.NewCommonCertificateVerifierSimple(true)
}

// Extend raises every signature in doc to the requested level - B to T, T to
// LT, LT to LTA - and returns the augmented document. The signing key is not
// involved: extension only adds time-stamps and validation data around the
// existing signature value.
func Extend(doc Document, opts ExtendOptions) (Document, error) {
	if doc == nil {
		return nil, ErrNoDocument
	}
	level, err := opts.Format.BaselineLevel(opts.Level)
	if err != nil {
		return nil, err
	}
	if opts.TSPSource == nil {
		return nil, fmt.Errorf("%w (level %s requested)", ErrTSPSourceRequired, opts.Level)
	}
	if len(opts.PasswordProtection) > 0 && opts.Format != FormatPAdES {
		return nil, fmt.Errorf("%w (format %s)", ErrPasswordProtectionNotApplicable, opts.Format)
	}

	var extended Document
	err = recovered("extend", func() error {
		switch opts.Format {
		case FormatCAdES:
			parameters := cades.NewSignatureParameters()
			applyExtendParameters(parameters, opts, level)
			service := cades.NewService(opts.certificateVerifier())
			service.SetTspSource(opts.TSPSource)
			extended = service.ExtendDocument(doc, parameters)

		case FormatXAdES:
			parameters := xades.NewSignatureParameters()
			applyExtendParameters(parameters, opts, level)
			service := xades.NewService(opts.certificateVerifier())
			service.SetTspSource(opts.TSPSource)
			extended = service.ExtendDocument(doc, parameters)

		case FormatPAdES:
			parameters := pades.NewSignatureParameters()
			applyExtendParameters(parameters, opts, level)
			parameters.SetPasswordProtection(opts.PasswordProtection)
			service := pades.NewService(opts.certificateVerifier())
			service.SetTspSource(opts.TSPSource)
			extended = service.ExtendDocument(doc, parameters)

		case FormatJAdES:
			parameters := jades.NewSignatureParameters()
			applyExtendParameters(parameters, opts, level)
			service := jades.NewService(opts.certificateVerifier())
			service.SetTspSource(opts.TSPSource)
			extended = service.ExtendDocument(doc, parameters)

		case FormatASiCWithCAdES:
			parameters := asiccades.NewASiCWithCAdESSignatureParameters()
			applyExtendParameters(parameters, opts, level)
			service := asiccades.NewASiCWithCAdESService(opts.certificateVerifier())
			service.SetTspSource(opts.TSPSource)
			extended = service.ExtendDocument(doc, parameters)

		case FormatASiCWithXAdES:
			parameters := asicxades.NewASiCWithXAdESSignatureParameters()
			applyExtendParameters(parameters, opts, level)
			service := asicxades.NewASiCWithXAdESService(opts.certificateVerifier())
			service.SetTspSource(opts.TSPSource)
			extended = service.ExtendDocument(doc, parameters)

		default:
			return fmt.Errorf("%w: %q", ErrUnsupportedFormat, string(opts.Format))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return extended, nil
}

// signatureParameters is the slice of the ported signature-parameter types the
// facade sets. Every concrete parameters type - SignatureParameters,
// SignatureParameters, SignatureParameters,
// SignatureParameters and the two ASiC ones - satisfies it, several of
// them with their own overriding SetSignatureLevel, which is why the facade
// calls through this interface rather than through a shared embedded struct.
type signatureParameters interface {
	SetSignatureLevel(signatureLevel SignatureLevel)
	SetDigestAlgorithm(digestAlgorithm DigestAlgorithm)
	SetSignaturePackaging(packaging SignaturePackaging)
	SetSigningCertificate(signingCertificate *model.CertificateToken)
	SetCertificateChainFromTokens(certificateChain ...*model.CertificateToken)
	SetDetachedContents(detachedContents []model.DSSDocument)
	BLevel() *model.BLevelParameters
}

// tspSourceSetter is the SetTspSource of document.AbstractSignatureService,
// which every ported signature service embeds.
type tspSourceSetter interface {
	SetTspSource(tspSource spivalidation.TSPSource)
}

// applyCommonParameters fills in the parameters every format shares.
func applyCommonParameters(parameters signatureParameters, opts SignOptions, level SignatureLevel, signer *Signer) {
	parameters.SetSignatureLevel(level)
	parameters.SetDigestAlgorithm(opts.digestAlgorithm())
	if packaging := opts.packaging(); packaging != "" {
		parameters.SetSignaturePackaging(packaging)
	}
	parameters.SetSigningCertificate(signer.Certificate())
	parameters.SetCertificateChainFromTokens(signer.CertificateChain()...)
	if len(opts.DetachedContents) > 0 {
		parameters.SetDetachedContents(opts.DetachedContents)
	}
	if opts.SigningTime != nil {
		parameters.BLevel().SetSigningDate(opts.SigningTime)
	}
}

// applyExtendParameters fills in the parameters extension needs. Extension
// touches an existing signature, so neither the signing certificate nor the
// packaging is set.
func applyExtendParameters(parameters signatureParameters, opts ExtendOptions, level SignatureLevel) {
	parameters.SetSignatureLevel(level)
	if len(opts.DetachedContents) > 0 {
		parameters.SetDetachedContents(opts.DetachedContents)
	}
}

// applyTSPSource hands the TSPSource to the service when one was configured.
func applyTSPSource(service tspSourceSetter, opts SignOptions) {
	if opts.TSPSource != nil {
		service.SetTspSource(opts.TSPSource)
	}
}

// signOne runs the get-data-to-sign / sign / embed flow of
// document.SignatureService.
func signOne[SP model.SerializableSignatureParameters, TP model.SerializableTimestampParameters](
	service document.SignatureService[SP, TP], parameters SP, doc Document,
	signer *Signer, digestAlgorithm DigestAlgorithm) (Document, error) {
	dataToSign := service.GetDataToSign(doc, parameters)
	signatureValue, err := signer.conn.Sign(dataToSign, digestAlgorithm, signer.key)
	if err != nil {
		return nil, fmt.Errorf("dss: signing the data to be signed: %w", err)
	}
	return service.SignDocument(doc, parameters, signatureValue), nil
}

// signMany is signOne for document.MultipleDocumentsSignatureService, which is
// how the ASiC services cover several documents with one signature.
func signMany[SP model.SerializableSignatureParameters, TP model.SerializableTimestampParameters](
	service document.MultipleDocumentsSignatureService[SP, TP], parameters SP, docs []Document,
	signer *Signer, digestAlgorithm DigestAlgorithm) (Document, error) {
	dataToSign := service.GetDataToSign(docs, parameters)
	signatureValue, err := signer.conn.Sign(dataToSign, digestAlgorithm, signer.key)
	if err != nil {
		return nil, fmt.Errorf("dss: signing the data to be signed: %w", err)
	}
	return service.SignDocument(docs, parameters, signatureValue), nil
}
