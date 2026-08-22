// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/policy/ValidationPolicyLoader.java (DSS 6.5.RC1).
//
// # ServiceLoader stand-in
//
// Java discovers ValidationPolicyFactory/CryptographicSuiteFactory implementations via
// ServiceLoader.load(...) against 'dss-policy-jaxb'/'dss-policy-crypto-xml'/
// 'dss-policy-crypto-json' on the classpath. Go has no runtime service-provider discovery, so -
// following the DocumentAnalyzerFactory / ASiCContainerExtractorFactory / CryptographicSuite
// precedent already established across this codebase (spi/validation/analyzer,
// asic/default_container_merger.go) - this file exposes package-level registries
// (RegisterValidationPolicyFactory / RegisterCryptographicSuiteFactory), consulted in
// registration order exactly like Java's ServiceLoader iterator. Nothing in THIS chunk populates
// them: the concrete implementations (dss/policy.EtsiValidationPolicyFactory,
// dss/policy/crypto/xml.CryptographicSuiteXmlFactory, dss/policy/crypto/json.
// CryptographicSuiteJsonFactory) already exist as frozen packages this chunk's manifest does not
// include and PORTING.md forbids editing; wiring an init()-time self-registration into them is
// left for a follow-up integration step (or the composing application), matching how
// dss/cades/cms_document_analyzer_factory.go self-registers into
// spi/validation/analyzer.RegisterDocumentAnalyzerFactory from a package that came AFTER the
// registry - here the registry unavoidably comes after the frozen factory packages instead, so
// the direction of the wiring inverts, but the mechanism (a package-level registry + Register
// func, consulted first-match-wins) matches the codebase-wide convention. This is a cross-chunk/
// integration flag: see this batch's porter notes.
//
// # Overload flattening
//
// Java overloads with(Default)CryptographicSuite(For(Context[, SubContext])) across six argument
// types (DSSDocument, InputStream, File, String classpath-path, CryptographicSuite,
// CryptographicSuiteCatalogue) and three specificity levels (none / Context / Context+SubContext)
// - eighteen methods on this type alone, mirrored again on ValidationPolicyLoaderWithCryptoSuite.
// Go cannot overload by argument type, so per PORTING.md every argument-type overload gets its
// own name-suffixed method. File and the classpath-relative String overload converge on the same
// Go behavior (both end up reading a path from the local filesystem, since Go has no classpath
// resource loader - see dss/policy/validation_policy_facade.go's GetValidationPolicyFromPath,
// which already made this exact call for ValidationPolicyFacade); both entry points are kept
// (matching Java's distinguished precondition-check messages, asserted by
// ValidationPolicyLoaderTest#cryptoSuiteNullTest) but share one filesystem-reading
// implementation.
//
// # Map order
//
// Java's private Map<CryptographicSuite, List<ContextAndSubContext>> cryptographicSuitesMap
// (LinkedHashMap) groups repeated with*() calls against the SAME CryptographicSuite instance
// (object identity - e.g. a CryptographicSuiteCatalogue caches and reuses one CryptographicSuite
// per distinct algorithm-list, so "global" and "CERTIFICATE" scopes routinely share one
// instance) under one ordered scope list, then create() replays that insertion order. Per
// PORTING.md's "Map<K,V> (order-sensitive upstream iteration -> slice of pairs or explicit sort)"
// rule, this is a Go map (comparable interface values holding pointer-receiver implementations
// compare by pointer identity, exactly matching Java's default Object equality here) paired with
// an insertion-order key slice.
package policy

import (
	"errors"
	"fmt"
	"io"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	modelpolicy "github.com/ryftcore/dss-go/dss/model/policy"
)

// validationPolicyFactoryRegistry and cryptographicSuiteFactoryRegistry hold the factories
// registered via Register*, consulted in registration order - see the file header.
var (
	validationPolicyFactoryRegistry   []modelpolicy.ValidationPolicyFactory
	cryptographicSuiteFactoryRegistry []modelpolicy.CryptographicSuiteFactory
)

// RegisterValidationPolicyFactory registers a ValidationPolicyFactory to be consulted by
// FromDefaultValidationPolicy / FromValidationPolicyDocument. See the file header.
func RegisterValidationPolicyFactory(f modelpolicy.ValidationPolicyFactory) {
	validationPolicyFactoryRegistry = append(validationPolicyFactoryRegistry, f)
}

// RegisterCryptographicSuiteFactory registers a CryptographicSuiteFactory to be consulted by
// the WithCryptographicSuite*/WithDefaultCryptographicSuite* family. See the file header.
func RegisterCryptographicSuiteFactory(f modelpolicy.CryptographicSuiteFactory) {
	cryptographicSuiteFactoryRegistry = append(cryptographicSuiteFactoryRegistry, f)
}

// ValidationPolicyLoader creates an instance of a ValidationPolicy.
type ValidationPolicyLoader struct {
	// validationPolicy is the validation policy document, when provided.
	validationPolicy modelpolicy.ValidationPolicy

	// cryptoOrder preserves the insertion order of cryptoScopes' keys; see the file header.
	cryptoOrder []modelpolicy.CryptographicSuite

	// cryptoScopes is a map of cryptographic suites and their applicability scopes.
	cryptoScopes map[modelpolicy.CryptographicSuite][]contextAndSubContext
}

// newValidationPolicyLoader ports the protected ValidationPolicyLoader(ValidationPolicy)
// constructor (the protected empty constructor delegates to this one with a nil policy).
func newValidationPolicyLoader(validationPolicy modelpolicy.ValidationPolicy) *ValidationPolicyLoader {
	return &ValidationPolicyLoader{
		validationPolicy: validationPolicy,
		cryptoScopes:     make(map[modelpolicy.CryptographicSuite][]contextAndSubContext),
	}
}

// FromDefaultValidationPolicy creates a ValidationPolicyLoader from a default validation policy.
func FromDefaultValidationPolicy() *ValidationPolicyLoader {
	return FromValidationPolicy(loadDefaultPolicy())
}

// FromValidationPolicyDocument creates a ValidationPolicyLoader from a custom validation policy
// document. Port of the fromValidationPolicy(DSSDocument) overload.
//
// Panics when validationPolicy is nil (Java's Objects.requireNonNull("Validation policy document
// cannot be null!")).
func FromValidationPolicyDocument(validationPolicy model.DSSDocument) *ValidationPolicyLoader {
	if validationPolicy == nil {
		panic("Validation policy document cannot be null!")
	}
	return FromValidationPolicy(loadPolicy(validationPolicy))
}

// FromValidationPolicyReader creates a ValidationPolicyLoader from a custom validation policy
// io.Reader. Port of the fromValidationPolicy(InputStream) overload.
//
// Panics when validationPolicyStream is nil (Java's Objects.requireNonNull("Validation policy
// stream cannot be null!")).
func FromValidationPolicyReader(validationPolicyStream io.Reader) *ValidationPolicyLoader {
	if validationPolicyStream == nil {
		panic("Validation policy stream cannot be null!")
	}
	data, err := io.ReadAll(validationPolicyStream)
	if err != nil {
		panic(fmt.Sprintf("Unable to load the validation policy document. Reason : %s", err.Error()))
	}
	return FromValidationPolicyDocument(model.NewInMemoryDocument(data))
}

// FromValidationPolicyFile creates a ValidationPolicyLoader from a custom validation policy file
// on disk. Port of the fromValidationPolicy(File) overload.
//
// Panics when validationPolicyFile is empty (Java's Objects.requireNonNull("Validation policy
// file cannot be null!")).
func FromValidationPolicyFile(validationPolicyFile string) (*ValidationPolicyLoader, error) {
	if validationPolicyFile == "" {
		panic("Validation policy file cannot be null!")
	}
	doc, err := model.NewFileDocument(validationPolicyFile)
	if err != nil {
		return nil, err
	}
	return FromValidationPolicyDocument(doc), nil
}

// FromValidationPolicyPath creates a ValidationPolicyLoader from a custom validation policy file.
// Port of the fromValidationPolicy(String) overload, whose `validationPolicyFilePath` was a
// classpath resource path in Java; Go has no classpath, so - matching the precedent already set
// by dss/policy/validation_policy_facade.go's GetValidationPolicyFromPath - this reads
// validationPolicyFilePath directly from the filesystem instead, like FromValidationPolicyFile.
//
// Panics when validationPolicyFilePath is empty (Java's Objects.requireNonNull("Validation
// policy file path cannot be null!")).
func FromValidationPolicyPath(validationPolicyFilePath string) (*ValidationPolicyLoader, error) {
	if validationPolicyFilePath == "" {
		panic("Validation policy file path cannot be null!")
	}
	return FromValidationPolicyFile(validationPolicyFilePath)
}

// FromValidationPolicy creates a ValidationPolicyLoader from a custom validation policy. Port of
// the fromValidationPolicy(ValidationPolicy) overload.
//
// Panics when validationPolicy is nil (Java's Objects.requireNonNull("Validation policy cannot
// be null!")).
func FromValidationPolicy(validationPolicy modelpolicy.ValidationPolicy) *ValidationPolicyLoader {
	if validationPolicy == nil {
		panic("Validation policy cannot be null!")
	}
	return newValidationPolicyLoader(validationPolicy)
}

// WithDefaultCryptographicSuite sets a default cryptographic suite, loading the first available
// cryptographic suite implementation. DSS provides two modules with implementations, namely
// 'dss-policy-crypto-xml' and 'dss-policy-crypto-json'.
func (l *ValidationPolicyLoader) WithDefaultCryptographicSuite() *ValidationPolicyLoaderWithCryptoSuite {
	return l.WithDefaultCryptographicSuiteForContextAndSubContextUnchecked("", "")
}

// WithDefaultCryptographicSuiteForContext sets a default cryptographic suite for the given
// Context. The supported contexts are: SIGNATURE, COUNTER_SIGNATURE, TIMESTAMP,
// EVIDENCE_RECORD, REVOCATION. The cryptographic suite will be used only for the specific scope.
func (l *ValidationPolicyLoader) WithDefaultCryptographicSuiteForContext(context enumerations.Context) *ValidationPolicyLoaderWithCryptoSuite {
	return l.WithDefaultCryptographicSuiteForContextAndSubContextUnchecked(context, "")
}

// WithDefaultCryptographicSuiteForContextAndSubContext sets a default cryptographic suite for
// the given Context and SubContext. The supported subContext are: SIGNING_CERT and
// CA_CERTIFICATE. The cryptographic suite will be used only for the specific scope.
//
// Returns an error when context is Context_EVIDENCE_RECORD and subContext is not empty (Java's
// thrown IllegalArgumentException).
func (l *ValidationPolicyLoader) WithDefaultCryptographicSuiteForContextAndSubContext(context enumerations.Context, subContext enumerations.SubContext) (*ValidationPolicyLoaderWithCryptoSuite, error) {
	if enumerations.Context_EVIDENCE_RECORD == context && subContext != "" {
		return nil, errors.New("please use a NULL SubContext for the Context.EVIDENCE_RECORD or " +
			"use #withDefaultCryptographicSuiteForContext(cryptographicSuite, context) method")
	}
	return l.WithDefaultCryptographicSuiteForContextAndSubContextUnchecked(context, subContext), nil
}

// WithDefaultCryptographicSuiteForContextAndSubContextUnchecked is the shared implementation
// behind WithDefaultCryptographicSuite/WithDefaultCryptographicSuiteForContext (whose Java
// counterparts carry no EVIDENCE_RECORD+SubContext guard) and the validated
// WithDefaultCryptographicSuiteForContextAndSubContext.
func (l *ValidationPolicyLoader) WithDefaultCryptographicSuiteForContextAndSubContextUnchecked(context enumerations.Context, subContext enumerations.SubContext) *ValidationPolicyLoaderWithCryptoSuite {
	return l.WithCryptographicSuiteCatalogueForContextAndSubContext(loadDefaultCryptographicSuite(), context, subContext)
}

// WithCryptographicSuiteDocument sets a global cryptographic suite document. The suite will
// overwrite all cryptographic constraints defined in the original ValidationPolicy file. It is
// also used when a cryptographic suite is not provided for a specific scope.
//
// Panics when cryptographicSuite is nil (Java's Objects.requireNonNull("Cryptographic suite
// document cannot be null!")).
func (l *ValidationPolicyLoader) WithCryptographicSuiteDocument(cryptographicSuite model.DSSDocument) *ValidationPolicyLoaderWithCryptoSuite {
	return l.WithCryptographicSuiteDocumentForContextAndSubContext(cryptographicSuite, "", "")
}

// WithCryptographicSuiteDocumentForContext sets a cryptographic suite document for the given
// Context.
func (l *ValidationPolicyLoader) WithCryptographicSuiteDocumentForContext(cryptographicSuite model.DSSDocument, context enumerations.Context) *ValidationPolicyLoaderWithCryptoSuite {
	return l.WithCryptographicSuiteDocumentForContextAndSubContext(cryptographicSuite, context, "")
}

// WithCryptographicSuiteDocumentForContextAndSubContext sets a cryptographic suite document for
// the given Context and SubContext.
//
// Panics when cryptographicSuite is nil (Java's Objects.requireNonNull("Cryptographic suite
// document cannot be null!")).
func (l *ValidationPolicyLoader) WithCryptographicSuiteDocumentForContextAndSubContext(cryptographicSuite model.DSSDocument, context enumerations.Context, subContext enumerations.SubContext) *ValidationPolicyLoaderWithCryptoSuite {
	if cryptographicSuite == nil {
		panic("Cryptographic suite document cannot be null!")
	}
	return l.WithCryptographicSuiteCatalogueForContextAndSubContext(loadCryptographicSuiteCatalogue(cryptographicSuite), context, subContext)
}

// WithCryptographicSuiteReader sets a global cryptographic suite from an io.Reader. See
// WithCryptographicSuiteDocument.
//
// Panics when cryptographicSuiteIS is nil (Java's Objects.requireNonNull("Cryptographic suite
// stream cannot be null!")).
func (l *ValidationPolicyLoader) WithCryptographicSuiteReader(cryptographicSuiteIS io.Reader) *ValidationPolicyLoaderWithCryptoSuite {
	return l.WithCryptographicSuiteReaderForContextAndSubContext(cryptographicSuiteIS, "", "")
}

// WithCryptographicSuiteReaderForContext sets a cryptographic suite from an io.Reader for the
// given Context.
func (l *ValidationPolicyLoader) WithCryptographicSuiteReaderForContext(cryptographicSuiteIS io.Reader, context enumerations.Context) *ValidationPolicyLoaderWithCryptoSuite {
	return l.WithCryptographicSuiteReaderForContextAndSubContext(cryptographicSuiteIS, context, "")
}

// WithCryptographicSuiteReaderForContextAndSubContext sets a cryptographic suite from an
// io.Reader for the given Context and SubContext.
//
// Panics when cryptographicSuiteIS is nil (Java's Objects.requireNonNull("Cryptographic suite
// stream cannot be null!")).
func (l *ValidationPolicyLoader) WithCryptographicSuiteReaderForContextAndSubContext(cryptographicSuiteIS io.Reader, context enumerations.Context, subContext enumerations.SubContext) *ValidationPolicyLoaderWithCryptoSuite {
	if cryptographicSuiteIS == nil {
		panic("Cryptographic suite stream cannot be null!")
	}
	data, err := io.ReadAll(cryptographicSuiteIS)
	if err != nil {
		panic(fmt.Sprintf("Unable to load a cryptographic suite. Reason : %s", err.Error()))
	}
	return l.WithCryptographicSuiteDocumentForContextAndSubContext(model.NewInMemoryDocument(data), context, subContext)
}

// WithCryptographicSuiteFile sets a global cryptographic suite file. See
// WithCryptographicSuiteDocument.
//
// Panics when cryptographicSuiteFile is empty (Java's Objects.requireNonNull("Cryptographic
// suite file cannot be null!")).
func (l *ValidationPolicyLoader) WithCryptographicSuiteFile(cryptographicSuiteFile string) (*ValidationPolicyLoaderWithCryptoSuite, error) {
	return l.WithCryptographicSuiteFileForContextAndSubContext(cryptographicSuiteFile, "", "")
}

// WithCryptographicSuiteFileForContext sets a cryptographic suite file for the given Context.
func (l *ValidationPolicyLoader) WithCryptographicSuiteFileForContext(cryptographicSuiteFile string, context enumerations.Context) (*ValidationPolicyLoaderWithCryptoSuite, error) {
	return l.WithCryptographicSuiteFileForContextAndSubContext(cryptographicSuiteFile, context, "")
}

// WithCryptographicSuiteFileForContextAndSubContext sets a cryptographic suite file for the
// given Context and SubContext.
//
// Panics when cryptographicSuiteFile is empty (Java's Objects.requireNonNull("Cryptographic
// suite file cannot be null!")).
func (l *ValidationPolicyLoader) WithCryptographicSuiteFileForContextAndSubContext(cryptographicSuiteFile string, context enumerations.Context, subContext enumerations.SubContext) (*ValidationPolicyLoaderWithCryptoSuite, error) {
	if cryptographicSuiteFile == "" {
		panic("Cryptographic suite file cannot be null!")
	}
	doc, err := model.NewFileDocument(cryptographicSuiteFile)
	if err != nil {
		return nil, err
	}
	return l.WithCryptographicSuiteDocumentForContextAndSubContext(doc, context, subContext), nil
}

// WithCryptographicSuitePath sets a global cryptographic suite file. Port of the
// withCryptographicSuite(String) overload; see FromValidationPolicyPath's doc comment on the
// classpath-to-filesystem collapse this port applies.
//
// Panics when cryptographicSuiteFilePath is empty (Java's Objects.requireNonNull("Cryptographic
// suite file path cannot be null!")).
func (l *ValidationPolicyLoader) WithCryptographicSuitePath(cryptographicSuiteFilePath string) (*ValidationPolicyLoaderWithCryptoSuite, error) {
	return l.WithCryptographicSuitePathForContextAndSubContext(cryptographicSuiteFilePath, "", "")
}

// WithCryptographicSuitePathForContext sets a cryptographic suite file for the given Context.
func (l *ValidationPolicyLoader) WithCryptographicSuitePathForContext(cryptographicSuiteFilePath string, context enumerations.Context) (*ValidationPolicyLoaderWithCryptoSuite, error) {
	return l.WithCryptographicSuitePathForContextAndSubContext(cryptographicSuiteFilePath, context, "")
}

// WithCryptographicSuitePathForContextAndSubContext sets a cryptographic suite file for the
// given Context and SubContext.
//
// Panics when cryptographicSuiteFilePath is empty (Java's Objects.requireNonNull("Cryptographic
// suite file path cannot be null!")).
func (l *ValidationPolicyLoader) WithCryptographicSuitePathForContextAndSubContext(cryptographicSuiteFilePath string, context enumerations.Context, subContext enumerations.SubContext) (*ValidationPolicyLoaderWithCryptoSuite, error) {
	if cryptographicSuiteFilePath == "" {
		panic("Cryptographic suite file path cannot be null!")
	}
	return l.WithCryptographicSuiteFileForContextAndSubContext(cryptographicSuiteFilePath, context, subContext)
}

// WithCryptographicSuite sets a global cryptographic suite.
//
// Panics when cryptographicSuite is nil (Java's Objects.requireNonNull("Cryptographic suite
// cannot be null!")).
func (l *ValidationPolicyLoader) WithCryptographicSuite(cryptographicSuite modelpolicy.CryptographicSuite) *ValidationPolicyLoaderWithCryptoSuite {
	return l.WithCryptographicSuiteForContext(cryptographicSuite, "")
}

// WithCryptographicSuiteForContext sets a cryptographic suite for the given Context.
func (l *ValidationPolicyLoader) WithCryptographicSuiteForContext(cryptographicSuite modelpolicy.CryptographicSuite, context enumerations.Context) *ValidationPolicyLoaderWithCryptoSuite {
	return l.WithCryptographicSuiteForContextAndSubContext(cryptographicSuite, context, "")
}

// WithCryptographicSuiteForContextAndSubContext sets a cryptographic suite for the given
// Context and SubContext.
//
// Panics when cryptographicSuite is nil (Java's Objects.requireNonNull("Cryptographic suite
// cannot be null!")).
func (l *ValidationPolicyLoader) WithCryptographicSuiteForContextAndSubContext(cryptographicSuite modelpolicy.CryptographicSuite, context enumerations.Context, subContext enumerations.SubContext) *ValidationPolicyLoaderWithCryptoSuite {
	if cryptographicSuite == nil {
		panic("Cryptographic suite cannot be null!")
	}
	l.addCryptographicSuiteScope(cryptographicSuite, contextAndSubContext{context: context, subContext: subContext})
	return newValidationPolicyLoaderWithCryptoSuite(l, []modelpolicy.CryptographicSuite{cryptographicSuite})
}

// WithCryptographicSuiteCatalogue sets a global cryptographic suite catalogue.
//
// Panics when cryptographicSuiteCatalogue is nil (Java's Objects.requireNonNull("Cryptographic
// suite catalogue cannot be null!")).
func (l *ValidationPolicyLoader) WithCryptographicSuiteCatalogue(cryptographicSuiteCatalogue *modelpolicy.CryptographicSuiteCatalogue) *ValidationPolicyLoaderWithCryptoSuite {
	return l.WithCryptographicSuiteCatalogueForContext(cryptographicSuiteCatalogue, "")
}

// WithCryptographicSuiteCatalogueForContext sets a cryptographic suite catalogue for the given
// Context.
func (l *ValidationPolicyLoader) WithCryptographicSuiteCatalogueForContext(cryptographicSuiteCatalogue *modelpolicy.CryptographicSuiteCatalogue, context enumerations.Context) *ValidationPolicyLoaderWithCryptoSuite {
	return l.WithCryptographicSuiteCatalogueForContextAndSubContext(cryptographicSuiteCatalogue, context, "")
}

// WithCryptographicSuiteCatalogueForContextAndSubContext sets a cryptographic suite catalogue
// for the given Context and SubContext. The supported contexts are: SIGNATURE,
// COUNTER_SIGNATURE, KEY_BINDING_SIGNATURE, TIMESTAMP, EVIDENCE_RECORD, REVOCATION, EAA,
// EAA_REVOCATION. The supported subContext are: SIGNING_CERT and CA_CERTIFICATE.
//
// Panics when cryptographicSuiteCatalogue is nil (Java's Objects.requireNonNull("Cryptographic
// suite catalogue cannot be null!")) or context is set but unsupported (Java's
// UnsupportedOperationException("The Context '%s' is not supported!")).
func (l *ValidationPolicyLoader) WithCryptographicSuiteCatalogueForContextAndSubContext(cryptographicSuiteCatalogue *modelpolicy.CryptographicSuiteCatalogue, context enumerations.Context, subContext enumerations.SubContext) *ValidationPolicyLoaderWithCryptoSuite {
	if cryptographicSuiteCatalogue == nil {
		panic("Cryptographic suite catalogue cannot be null!")
	}

	var cryptographicSuites []modelpolicy.CryptographicSuite

	add := func(cs modelpolicy.CryptographicSuite, ctx enumerations.Context, subCtx enumerations.SubContext) {
		l.addCryptographicSuiteScope(cs, contextAndSubContext{context: ctx, subContext: subCtx})
		cryptographicSuites = append(cryptographicSuites, cs)
	}

	if context != "" {
		switch context {
		case enumerations.Context_SIGNATURE, enumerations.Context_CERTIFICATE:
			if subContext != "" {
				add(cryptographicSuiteCatalogue.SignatureCertificatesCryptographicSuite(), context, subContext)
			} else {
				add(cryptographicSuiteCatalogue.SignatureCryptographicSuite(), context, "")
				add(cryptographicSuiteCatalogue.SignatureCertificatesCryptographicSuite(), context, enumerations.SubContext_SIGNING_CERT)
				add(cryptographicSuiteCatalogue.SignatureCertificatesCryptographicSuite(), context, enumerations.SubContext_CA_CERTIFICATE)
			}

		case enumerations.Context_COUNTER_SIGNATURE:
			if subContext != "" {
				add(cryptographicSuiteCatalogue.CounterSignatureCertificatesCryptographicSuite(), context, subContext)
			} else {
				add(cryptographicSuiteCatalogue.CounterSignatureCryptographicSuite(), context, "")
				add(cryptographicSuiteCatalogue.CounterSignatureCertificatesCryptographicSuite(), context, enumerations.SubContext_SIGNING_CERT)
				add(cryptographicSuiteCatalogue.CounterSignatureCertificatesCryptographicSuite(), context, enumerations.SubContext_CA_CERTIFICATE)
			}

		case enumerations.Context_KEY_BINDING_SIGNATURE:
			if subContext != "" {
				add(cryptographicSuiteCatalogue.KeyBindingSignatureCertificatesCryptographicSuite(), context, subContext)
			} else {
				add(cryptographicSuiteCatalogue.KeyBindingSignatureCryptographicSuite(), context, "")
				add(cryptographicSuiteCatalogue.KeyBindingSignatureCertificatesCryptographicSuite(), context, enumerations.SubContext_SIGNING_CERT)
				add(cryptographicSuiteCatalogue.KeyBindingSignatureCertificatesCryptographicSuite(), context, enumerations.SubContext_CA_CERTIFICATE)
			}

		case enumerations.Context_REVOCATION:
			if subContext != "" {
				add(cryptographicSuiteCatalogue.RevocationCertificatesCryptographicSuite(), context, subContext)
			} else {
				add(cryptographicSuiteCatalogue.RevocationCryptographicSuite(), context, "")
				add(cryptographicSuiteCatalogue.RevocationCertificatesCryptographicSuite(), context, enumerations.SubContext_SIGNING_CERT)
				add(cryptographicSuiteCatalogue.RevocationCertificatesCryptographicSuite(), context, enumerations.SubContext_CA_CERTIFICATE)
			}

		case enumerations.Context_TIMESTAMP:
			if subContext != "" {
				add(cryptographicSuiteCatalogue.TimestampCertificatesCryptographicSuite(), context, subContext)
			} else {
				add(cryptographicSuiteCatalogue.TimestampCryptographicSuite(), context, "")
				add(cryptographicSuiteCatalogue.TimestampCertificatesCryptographicSuite(), context, enumerations.SubContext_SIGNING_CERT)
				add(cryptographicSuiteCatalogue.TimestampCertificatesCryptographicSuite(), context, enumerations.SubContext_CA_CERTIFICATE)
			}

		case enumerations.Context_EVIDENCE_RECORD:
			add(cryptographicSuiteCatalogue.EvidenceRecordCryptographicSuite(), context, subContext)

		case enumerations.Context_EAA:
			add(cryptographicSuiteCatalogue.EAACryptographicSuite(), context, subContext)

		case enumerations.Context_EAA_REVOCATION:
			add(cryptographicSuiteCatalogue.EAARevocationCryptographicSuite(), context, subContext)

		default:
			panic(fmt.Sprintf("The Context '%s' is not supported!", context))
		}

	} else {
		// apply all constraints
		add(cryptographicSuiteCatalogue.CryptographicSuite(), "", subContext)
		add(cryptographicSuiteCatalogue.SignatureCryptographicSuite(), enumerations.Context_CERTIFICATE, subContext)
		add(cryptographicSuiteCatalogue.SignatureCertificatesCryptographicSuite(), enumerations.Context_CERTIFICATE, enumerations.SubContext_SIGNING_CERT)
		add(cryptographicSuiteCatalogue.SignatureCertificatesCryptographicSuite(), enumerations.Context_CERTIFICATE, enumerations.SubContext_CA_CERTIFICATE)
		add(cryptographicSuiteCatalogue.SignatureCryptographicSuite(), enumerations.Context_SIGNATURE, subContext)
		add(cryptographicSuiteCatalogue.SignatureCertificatesCryptographicSuite(), enumerations.Context_SIGNATURE, enumerations.SubContext_SIGNING_CERT)
		add(cryptographicSuiteCatalogue.SignatureCertificatesCryptographicSuite(), enumerations.Context_SIGNATURE, enumerations.SubContext_CA_CERTIFICATE)
		add(cryptographicSuiteCatalogue.CounterSignatureCryptographicSuite(), enumerations.Context_COUNTER_SIGNATURE, subContext)
		add(cryptographicSuiteCatalogue.CounterSignatureCertificatesCryptographicSuite(), enumerations.Context_COUNTER_SIGNATURE, enumerations.SubContext_SIGNING_CERT)
		add(cryptographicSuiteCatalogue.CounterSignatureCertificatesCryptographicSuite(), enumerations.Context_COUNTER_SIGNATURE, enumerations.SubContext_CA_CERTIFICATE)
		add(cryptographicSuiteCatalogue.RevocationCryptographicSuite(), enumerations.Context_REVOCATION, subContext)
		add(cryptographicSuiteCatalogue.RevocationCertificatesCryptographicSuite(), enumerations.Context_REVOCATION, enumerations.SubContext_SIGNING_CERT)
		add(cryptographicSuiteCatalogue.RevocationCertificatesCryptographicSuite(), enumerations.Context_REVOCATION, enumerations.SubContext_CA_CERTIFICATE)
		add(cryptographicSuiteCatalogue.TimestampCryptographicSuite(), enumerations.Context_TIMESTAMP, subContext)
		add(cryptographicSuiteCatalogue.TimestampCertificatesCryptographicSuite(), enumerations.Context_TIMESTAMP, enumerations.SubContext_SIGNING_CERT)
		add(cryptographicSuiteCatalogue.TimestampCertificatesCryptographicSuite(), enumerations.Context_TIMESTAMP, enumerations.SubContext_CA_CERTIFICATE)
	}

	return newValidationPolicyLoaderWithCryptoSuite(l, cryptographicSuites)
}

// addCryptographicSuiteScope registers scope against cryptographicSuite, grouping repeated
// registrations of the SAME CryptographicSuite instance (see the file header on map order).
func (l *ValidationPolicyLoader) addCryptographicSuiteScope(cryptographicSuite modelpolicy.CryptographicSuite, scope contextAndSubContext) {
	if _, ok := l.cryptoScopes[cryptographicSuite]; !ok {
		l.cryptoOrder = append(l.cryptoOrder, cryptographicSuite)
	}
	l.cryptoScopes[cryptographicSuite] = append(l.cryptoScopes[cryptographicSuite], scope)
}

// Create builds a ValidationPolicy. Port of create().
//
// Panics if replaying a registered scope violates ValidationPolicyWithCryptographicSuite's
// EVIDENCE_RECORD+SubContext guard (Java's unchecked IllegalArgumentException, reachable here
// only via the unvalidated catalogue-based With*ForContextAndSubContext entry points - see their
// doc comments).
func (l *ValidationPolicyLoader) Create() modelpolicy.ValidationPolicy {
	if len(l.cryptoOrder) == 0 {
		return l.validationPolicy
	}

	validationPolicyWithCryptographicSuite := NewValidationPolicyWithCryptographicSuite(l.validationPolicy)
	for _, cryptographicSuite := range l.cryptoOrder {
		for _, scope := range l.cryptoScopes[cryptographicSuite] {
			if err := validationPolicyWithCryptographicSuite.SetCryptographicSuiteForContextAndSubContext(cryptographicSuite, scope.context, scope.subContext); err != nil {
				panic(err.Error())
			}
		}
	}
	return validationPolicyWithCryptographicSuite
}

// loadDefaultPolicy loads a default validation policy. Port of the private static
// loadDefaultPolicy().
//
// Panics when no ValidationPolicyFactory is registered (Java's
// UnsupportedOperationException("No ValidationPolicyFactory has been found! Please add
// 'dss-policy-jaxb' module to the classpath or create your own implementation.")).
func loadDefaultPolicy() modelpolicy.ValidationPolicy {
	if len(validationPolicyFactoryRegistry) == 0 {
		panic("No ValidationPolicyFactory has been found! Please add 'dss-policy-jaxb' module to the classpath or create your own implementation.")
	}
	return validationPolicyFactoryRegistry[0].LoadDefaultValidationPolicy()
}

// loadPolicy loads a validation policy from the given DSSDocument. Port of the private static
// loadPolicy(DSSDocument).
//
// Panics when no registered ValidationPolicyFactory supports validationPolicyDocument (Java's
// UnsupportedOperationException("The validation policy is not valid or no suitable
// ValidationPolicyFactory has been found! ...")).
func loadPolicy(validationPolicyDocument model.DSSDocument) modelpolicy.ValidationPolicy {
	for _, factory := range validationPolicyFactoryRegistry {
		if factory.IsSupported(validationPolicyDocument) {
			return factory.LoadValidationPolicy(validationPolicyDocument)
		}
	}
	panic("The validation policy is not valid or no suitable ValidationPolicyFactory has been found! " +
		"Please ensure the provided policy file is valid and 'dss-policy-jaxb' module is added to the classpath or " +
		"create your own implementation for a custom policy.")
}

// loadDefaultCryptographicSuite loads a default cryptographic suite catalogue. Port of the
// private static loadDefaultCryptographicSuite().
//
// Panics when no CryptographicSuiteFactory is registered (Java's
// UnsupportedOperationException("No ValidationPolicyFactory has been found! ...")). [sic - Java
// itself reuses the ValidationPolicyFactory wording in this message].
func loadDefaultCryptographicSuite() *modelpolicy.CryptographicSuiteCatalogue {
	if len(cryptographicSuiteFactoryRegistry) == 0 {
		panic("No ValidationPolicyFactory has been found! Please add 'dss-policy-jaxb' module to the classpath or create your own implementation.")
	}
	return cryptographicSuiteFactoryRegistry[0].LoadDefaultCryptographicSuite()
}

// loadCryptographicSuiteCatalogue loads a cryptographic suite from the given DSSDocument. Port
// of the private static loadCryptographicSuiteCatalogue(DSSDocument).
//
// Panics when no registered CryptographicSuiteFactory supports cryptographicSuiteDocument
// (Java's UnsupportedOperationException("The cryptographic suite file is not valid or no
// suitable CryptographicSuiteFactory has been found! ...")).
func loadCryptographicSuiteCatalogue(cryptographicSuiteDocument model.DSSDocument) *modelpolicy.CryptographicSuiteCatalogue {
	for _, factory := range cryptographicSuiteFactoryRegistry {
		if factory.IsSupported(cryptographicSuiteDocument) {
			return factory.LoadCryptographicSuite(cryptographicSuiteDocument)
		}
	}
	panic("The cryptographic suite file is not valid or no suitable CryptographicSuiteFactory has been found! " +
		"Please ensure the provided policy file is valid and 'dss-policy-crypto-xml' or 'dss-policy-crypto-json' module is added to the classpath or " +
		"create your own implementation for a custom cryptographic suite policy.")
}
