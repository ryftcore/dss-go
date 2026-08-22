// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/signature/AllDataObjectsTimeStampBuilder.java (DSS 6.5.RC1).
//
// Java catches TSPException | IOException | CMSException from the BouncyCastle-backed
// TimestampToken constructor and wraps them in a DSSException; the Go TimestampToken
// constructor already returns a single error, which this port wraps in the same message.
// slf4j logging is dropped (PORTING.md).
//
// The reference-processing loop keeps upstream's structure exactly: every ds:Reference is
// dereferenced through ReferenceProcessor, and its octets are fed to one running digest -
// canonicalized first when (and only when) the reference output is a node set AND the produced
// document actually parses as XML, which is the pair of conditions upstream tests.
package xades

import (
	"errors"
	"io"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/utils"
	xmlutils "github.com/ryftcore/dss-go/dss/xml/utils"
)

// AllDataObjectsTimeStampBuilder allows creating a XAdES content-timestamp which covers all
// documents (AllDataObjectsTimeStamp).
type AllDataObjectsTimeStampBuilder struct {
	// tspSource is the TSPSource to use.
	tspSource validation.TSPSource

	// signatureParameters holds the signature parameters.
	signatureParameters *SignatureParameters
}

// NewAllDataObjectsTimeStampBuilder is the default constructor.
// Port of AllDataObjectsTimeStampBuilder(TSPSource, XAdESSignatureParameters).
func NewAllDataObjectsTimeStampBuilder(tspSource validation.TSPSource,
	signatureParameters *SignatureParameters) *AllDataObjectsTimeStampBuilder {
	return &AllDataObjectsTimeStampBuilder{
		tspSource:           tspSource,
		signatureParameters: signatureParameters,
	}
}

// Build builds a message-imprint from the given document and generates a timestamp.
// Port of the #build(DSSDocument) overload.
func (b *AllDataObjectsTimeStampBuilder) Build(document model.DSSDocument) (*validation.TimestampToken, error) {
	return b.BuildForDocuments([]model.DSSDocument{document})
}

// BuildForDocuments timestamps the list of documents.
// Port of the #build(List<DSSDocument>) overload.
func (b *AllDataObjectsTimeStampBuilder) BuildForDocuments(
	documents []model.DSSDocument) (*validation.TimestampToken, error) {
	if err := allDataObjectsTimeStampBuilderAssertTimestampCreationPossible(documents); err != nil {
		return nil, err
	}

	// Prepare references
	references := b.signatureParameters.References()
	if utils.IsCollectionEmpty(references) {
		referenceIdProvider := NewReferenceIdProvider()
		referenceIdProvider.SetSignatureParameters(b.signatureParameters)
		referenceBuilder := NewReferenceBuilder(documents, b.signatureParameters, referenceIdProvider)
		builtReferences, err := referenceBuilder.Build()
		if err != nil {
			return nil, err
		}
		references = builtReferences
		b.signatureParameters.GetContext().SetReferences(references)
	} else {
		referenceVerifier := NewReferenceVerifier(b.signatureParameters)
		if err := referenceVerifier.CheckReferencesValidity(); err != nil {
			return nil, err
		}
	}

	contentTimestampParameters := b.signatureParameters.GetContentTimestampParameters()
	canonicalizationMethod := contentTimestampParameters.CanonicalizationMethod()

	digestAlgorithm := contentTimestampParameters.DigestAlgorithm()
	digestCalculator, err := spi.NewDSSMessageDigestCalculator(digestAlgorithm)
	if err != nil {
		return nil, err
	}
	for _, reference := range references {
		/*
		 * 1) process the retrieved ds:Reference element according to the reference-processing
		 * model of XMLDSIG [1] clause 4.4.3.2;
		 */
		referenceProcessor := NewReferenceProcessor(b.signatureParameters)
		referenceContent, err := referenceProcessor.ReferenceOutput(reference)
		if err != nil {
			return nil, err
		}

		/*
		 * 2) if the result is a XML node set, canonicalize it as specified in clause 4.5; and
		 * 3) concatenate the resulting octets to those resulting from previously processed
		 * ds:Reference elements in ds:SignedInfo.
		 */
		isNodeSet := ReferenceOutputTypeNodeSet == DSSXMLUtilsGetReferenceOutputType(reference) &&
			xmlutils.DomUtilsIsDOM(referenceContent)

		referenceIs, err := referenceContent.OpenStream()
		if err != nil {
			return nil, model.NewDSSErrorMessageCause(
				"Cannot build an AllDataObjectsTimestamp : An error occurred on reference extraction", err)
		}
		if isNodeSet {
			err = allDataObjectsTimeStampBuilderWriteDigestValueOnCanonicalizedInputStream(
				digestCalculator, referenceIs, canonicalizationMethod)
		} else {
			err = digestCalculator.UpdateReader(referenceIs)
		}
		utils.CloseQuietly(referenceIs)
		if err != nil {
			return nil, model.NewDSSErrorMessageCause(
				"Cannot build an AllDataObjectsTimestamp : An error occurred on reference extraction", err)
		}
	}
	messageDigest := digestCalculator.MessageDigest(digestAlgorithm)
	// Upstream traces "Computed AllDataObjectsTimestampData data digest: {}".

	timeStampResponse, err := b.tspSource.TimeStampResponse(digestAlgorithm, messageDigest.Value())
	if err != nil {
		return nil, err
	}
	token, err := validation.NewTimestampToken(timeStampResponse.Bytes(),
		enumerations.TimestampTypeAllDataObjectsTimestamp)
	if err != nil {
		return nil, model.NewDSSErrorMessageCause("Cannot build an AllDataObjectsTimestamp", err)
	}
	token.SetCanonicalizationMethod(canonicalizationMethod)
	return token, nil
}

// allDataObjectsTimeStampBuilderWriteDigestValueOnCanonicalizedInputStream ports the private
// writeDigestValueOnCanonicalizedInputStream.
func allDataObjectsTimeStampBuilderWriteDigestValueOnCanonicalizedInputStream(
	messageDigestCalculator *spi.DSSMessageDigestCalculator, is io.Reader,
	canonicalizationMethod string) error {
	os := messageDigestCalculator.Writer()
	canonicalizer, err := xmlutils.XMLCanonicalizerCreateInstanceWithMethod(canonicalizationMethod)
	if err != nil {
		utils.CloseQuietly(os)
		return err
	}
	if err := canonicalizer.CanonicalizeStreamTo(is, os); err != nil {
		utils.CloseQuietly(os)
		return err
	}
	return os.Close()
}

// allDataObjectsTimeStampBuilderAssertTimestampCreationPossible ports the private
// assertTimestampCreationPossible; Java's IllegalArgumentException becomes a returned error.
func allDataObjectsTimeStampBuilderAssertTimestampCreationPossible(documents []model.DSSDocument) error {
	for _, document := range documents {
		if _, isDigestDocument := document.(*model.DigestDocument); isDigestDocument {
			return errors.New("Content timestamp creation is not possible with DigestDocument!")
		}
	}
	return nil
}
