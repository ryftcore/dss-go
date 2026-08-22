// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/signature/SignaturePolicyStoreBuilder.java (DSS 6.5.RC1).
//
// Java extends the abstract ExtensionBuilder through its empty constructor; the Go port embeds
// ExtensionBuilder and registers itself with InitExtensionBuilder, the
// TokenBase.InitToken(self) convention of PORTING.md.
//
// Java's two addSignaturePolicyStore overloads cannot share one Go name:
//
//	addSignaturePolicyStore(DSSDocument, SignaturePolicyStore)          -> AddSignaturePolicyStore
//	addSignaturePolicyStore(DSSDocument, SignaturePolicyStore, String)  -> AddSignaturePolicyStoreForSignature
//
// checkDigest wraps every failure of the policy-digest recomputation in a DSSException, which
// in Go means recovering from the panic AbstractSignaturePolicyValidator.GetComputedDigest
// raises (see spi/policy) as well as taking the error XMLSignaturePolicyValidator returns.
// slf4j logging is dropped (PORTING.md).
package xades

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/signature"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/exception"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/xades/definition"
	xmlutils "github.com/ryftcore/dss-go/dss/xml/utils"
)

// SignaturePolicyStoreBuilder builds a XAdES SignaturePolicyStore.
type SignaturePolicyStoreBuilder struct {
	ExtensionBuilder
}

// NewSignaturePolicyStoreBuilder is the default constructor.
// Port of the SignaturePolicyStoreBuilder() constructor.
func NewSignaturePolicyStoreBuilder() *SignaturePolicyStoreBuilder {
	builder := &SignaturePolicyStoreBuilder{}
	builder.InitExtensionBuilder(builder)
	return builder
}

// AddSignaturePolicyStore adds a signaturePolicyStore to all signatures inside the document,
// matching the incorporated signature policy.
// Port of the #addSignaturePolicyStore(DSSDocument, SignaturePolicyStore) overload.
func (b *SignaturePolicyStoreBuilder) AddSignaturePolicyStore(signatureDocument model.DSSDocument,
	signaturePolicyStore *model.SignaturePolicyStore) (model.DSSDocument, error) {
	if signatureDocument == nil {
		panic("Signature document must be provided!")
	}
	if err := signaturePolicyStoreBuilderAssertConfigurationValid(signaturePolicyStore); err != nil {
		return nil, err
	}

	documentAnalyzer, err := b.InitDocumentAnalyzer(signatureDocument)
	if err != nil {
		return nil, err
	}

	signatures := documentAnalyzer.Signatures()
	signaturePolicyStoreAdded := false
	for _, signature := range signatures {
		xadesSignature, ok := signature.(*Signature)
		if !ok {
			// Java's (XAdESSignature) cast; a non-XAdES signature would raise a ClassCastException.
			return nil, fmt.Errorf("unexpected signature type %T", signature)
		}
		added, err := b.AddSignaturePolicyStoreIfDigestMatch(xadesSignature, b.DocumentDom, signaturePolicyStore)
		if err != nil {
			return nil, err
		}
		signaturePolicyStoreAdded = signaturePolicyStoreAdded || added
	}
	if !signaturePolicyStoreAdded {
		return nil, exception.NewIllegalInputException(
			"The process did not find a signature to add SignaturePolicyStore!")
	}

	return b.CreateXmlDocument()
}

// AddSignaturePolicyStoreForSignature adds a signaturePolicyStore to the signature with the
// given signatureId, if the signature policy identifier matches the policy provided within the
// SignaturePolicyStore.
// Port of the #addSignaturePolicyStore(DSSDocument, SignaturePolicyStore, String) overload.
func (b *SignaturePolicyStoreBuilder) AddSignaturePolicyStoreForSignature(
	signatureDocument model.DSSDocument, signaturePolicyStore *model.SignaturePolicyStore,
	signatureId string) (model.DSSDocument, error) {
	if signatureDocument == nil {
		panic("Signature document must be provided!")
	}
	if err := signaturePolicyStoreBuilderAssertConfigurationValid(signaturePolicyStore); err != nil {
		return nil, err
	}

	documentAnalyzer, err := b.InitDocumentAnalyzer(signatureDocument)
	if err != nil {
		return nil, err
	}
	signature := documentAnalyzer.SignatureByID(signatureId)
	if signature == nil {
		return nil, exception.NewIllegalInputException(
			fmt.Sprintf("Unable to find a signature with Id : %s!", signatureId))
	}
	xadesSignature, ok := signature.(*Signature)
	if !ok {
		return nil, exception.NewIllegalInputException(
			fmt.Sprintf("Unable to find a signature with Id : %s!", signatureId))
	}
	added, err := b.AddSignaturePolicyStoreIfDigestMatch(xadesSignature, b.DocumentDom, signaturePolicyStore)
	if err != nil {
		return nil, err
	}
	if !added {
		return nil, exception.NewIllegalInputException(fmt.Sprintf(
			"The process was not able to add SignaturePolicyStore to a signature with Id : %s!", signatureId))
	}

	return b.CreateXmlDocument()
}

// AddSignaturePolicyStoreIfDigestMatch adds a SignaturePolicyStore to documentDom if required,
// and reports whether it has been added for the particular signature.
// Port of the protected #addSignaturePolicyStoreIfDigestMatch.
func (b *SignaturePolicyStoreBuilder) AddSignaturePolicyStoreIfDigestMatch(
	xadesSignature *Signature, documentDom *xmldom.Node,
	signaturePolicyStore *model.SignaturePolicyStore) (bool, error) {
	if err := b.AssertUnsignedPropertiesExtensionPossible(xadesSignature); err != nil {
		return false, err
	}

	initialized, err := b.InitializeSignatureBuilder(xadesSignature)
	if err != nil {
		return false, err
	}
	xadesSignature = initialized

	if err := b.EnsureUnsignedProperties(); err != nil {
		return false, err
	}
	if err := b.EnsureUnsignedSignatureProperties(); err != nil {
		return false, err
	}

	digestMatch, err := b.CheckDigest(xadesSignature, signaturePolicyStore)
	if err != nil {
		return false, err
	}
	if digestMatch {
		signaturePolicyStoreElement := xmlutils.DomUtilsAddElement(documentDom,
			b.UnsignedSignaturePropertiesDom, b.Xades141Namespace(),
			definition.XAdES141ElementSignaturePolicyStore)

		if signaturePolicyStore.Id() != "" {
			signaturePolicyStoreElement.SetAttr(
				xmldom.Name{Local: definition.XAdES141AttributeID.AttributeName()},
				signaturePolicyStore.Id())
		}

		spDocSpecification := signaturePolicyStore.SpDocSpecification()
		if err := b.IncorporateSPDocSpecification(signaturePolicyStoreElement, spDocSpecification); err != nil {
			return false, err
		}

		signaturePolicyContent := signaturePolicyStore.SignaturePolicyContent()
		if signaturePolicyContent != nil {
			policyDocElement := xmlutils.DomUtilsAddElement(documentDom, signaturePolicyStoreElement,
				b.Xades141Namespace(), definition.XAdES141ElementSignaturePolicyDocument)

			policyBytes, err := spi.DSSUtilsToByteArrayOfDocument(signaturePolicyContent)
			if err != nil {
				return false, err
			}
			xmlutils.DomUtilsSetTextNode(documentDom, policyDocElement, utils.ToBase64(policyBytes))
		}

		sigPolDocLocalURI := signaturePolicyStore.SigPolDocLocalURI()
		if utils.IsStringNotEmpty(sigPolDocLocalURI) {
			xmlutils.DomUtilsAddTextElement(documentDom, signaturePolicyStoreElement,
				b.Xades141Namespace(), definition.XAdES141ElementSigPolDocLocalURI, sigPolDocLocalURI)
		}

		return true, nil
	}

	return false, nil
}

// CheckDigest verifies whether the digests computed in the provided SignaturePolicyStore match
// the digest defined in the incorporated signature policy identifier, i.e. whether the
// SignaturePolicyStore can be embedded. Port of the protected #checkDigest.
func (b *SignaturePolicyStoreBuilder) CheckDigest(xadesSignature *Signature,
	signaturePolicyStore *model.SignaturePolicyStore) (bool, error) {
	// Upstream reports the signature Id in each of the warnings below.
	signaturePolicy := xadesSignature.SignaturePolicy()
	if signaturePolicy == nil {
		// "No defined SignaturePolicyIdentifier for signature with Id : {}"
		return false, nil
	}
	digest := signaturePolicy.Digest()
	if digest.IsEmpty() {
		// "No defined digest for signature with Id : {}"; model.Digest is a value type, so
		// Java's null check becomes the zero-Digest check IsEmpty already expresses.
		return false, nil
	}

	signaturePolicyContent := signaturePolicyStore.SignaturePolicyContent()
	if signaturePolicyContent == nil {
		// "No policy document has been provided. Digests are not checked!"
		return true, nil
	}
	signaturePolicy.SetPolicyContent(signaturePolicyContent)

	// Java downcasts xadesSignature.getSignaturePolicy() to XAdESSignaturePolicy here (its
	// concrete runtime type). Go's Policy() promotion returns the plain base pointer
	// (frozen spi/validation.AdvancedSignature interface), so recover the concrete value via
	// the registry xades_signature_policy.go maintains for exactly this purpose.
	xadesSignaturePolicy, _ := SignaturePolicyFor(signaturePolicy)

	computedDigest, err := b.computeSignaturePolicyDigest(signaturePolicy, xadesSignaturePolicy, digest, signaturePolicyContent)
	if err != nil {
		return false, model.NewDSSError(fmt.Sprintf(
			"Unable to compute digest for a SignaturePolicyStore. Reason : %s", err.Error()))
	}

	digestMatch := digest.Equals(computedDigest)
	// Upstream warns when the digests differ.
	return digestMatch, nil
}

// computeSignaturePolicyDigest carries the body of Java's try block in checkDigest, including
// its catch-all: AbstractSignaturePolicyValidator.GetComputedDigest signals failure by panicking
// (see spi/policy), which Java's `catch (Exception e)` would have caught, so it is recovered
// here and turned into the error CheckDigest wraps.
func (b *SignaturePolicyStoreBuilder) computeSignaturePolicyDigest(signaturePolicy *signature.Policy,
	xadesSignaturePolicy *SignaturePolicy,
	digest model.Digest, signaturePolicyContent model.DSSDocument) (computedDigest model.Digest, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			if recoveredErr, ok := recovered.(error); ok {
				err = recoveredErr
			} else {
				err = fmt.Errorf("%v", recovered)
			}
		}
	}()

	validator := b.DocumentAnalyzer.SignaturePolicyValidatorLoader().LoadValidator(signaturePolicy)
	if xmlSignaturePolicyValidator, ok := validator.(*XMLSignaturePolicyValidator); ok {
		var transforms *xmldom.Node
		if xadesSignaturePolicy != nil {
			transforms = xadesSignaturePolicy.Transforms()
		}
		return xmlSignaturePolicyValidator.GetDigestAfterTransforms(signaturePolicyContent,
			digest.Algorithm(), transforms)
	}
	// Java re-reads signaturePolicyStore.getSignaturePolicyContent() here; it is the same
	// document checkDigest already holds.
	return validator.GetComputedDigest(signaturePolicyContent, digest.Algorithm()), nil
}

// signaturePolicyStoreBuilderAssertConfigurationValid ports the private
// assertConfigurationValid; Objects.requireNonNull becomes a panic carrying the Java message,
// the IllegalArgumentException a returned error.
func signaturePolicyStoreBuilderAssertConfigurationValid(
	signaturePolicyStore *model.SignaturePolicyStore) error {
	if signaturePolicyStore == nil {
		panic("SignaturePolicyStore must be provided")
	}
	if signaturePolicyStore.SpDocSpecification() == nil {
		panic("SpDocSpecification must be provided")
	}
	if signaturePolicyStore.SpDocSpecification().Id() == "" {
		panic("ID (OID or URI) for SpDocSpecification must be provided")
	}

	signaturePolicyContentPresent := signaturePolicyStore.SignaturePolicyContent() != nil
	sigPolDocLocalURIPresent := signaturePolicyStore.SigPolDocLocalURI() != ""
	if signaturePolicyContentPresent == sigPolDocLocalURIPresent {
		return fmt.Errorf("SignaturePolicyStore shall contain either " +
			"SignaturePolicyContent document or sigPolDocLocalURI!")
	}
	return nil
}
