// Ported from dss-asic-cades/src/main/java/eu/europa/esig/dss/asic/cades/timestamp/ASiCWithCAdESTimestampService.java (DSS 6.5.RC1).
package cades

import (
	"github.com/ryftcore/dss-go/dss/asic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/validation"
)

// ASiCWithCAdESTimestampService creates a timestamp covering signer files.
type ASiCWithCAdESTimestampService struct {
	// tspSource is used to retrieve a timestamp response.
	tspSource validation.TSPSource

	// asicFilenameFactory defines rules for filename creation for a timestamp file.
	asicFilenameFactory ASiCWithCAdESFilenameFactory
}

// NewASiCWithCAdESTimestampService is the default constructor. Ports
// ASiCWithCAdESTimestampService(TSPSource).
func NewASiCWithCAdESTimestampService(tspSource validation.TSPSource) *ASiCWithCAdESTimestampService {
	return NewASiCWithCAdESTimestampServiceWithFilenameFactory(tspSource, NewDefaultASiCWithCAdESFilenameFactory())
}

// NewASiCWithCAdESTimestampServiceWithFilenameFactory is the constructor with filename factory.
// Ports ASiCWithCAdESTimestampService(TSPSource, ASiCWithCAdESFilenameFactory).
func NewASiCWithCAdESTimestampServiceWithFilenameFactory(tspSource validation.TSPSource,
	asicFilenameFactory ASiCWithCAdESFilenameFactory) *ASiCWithCAdESTimestampService {
	return &ASiCWithCAdESTimestampService{
		tspSource:           tspSource,
		asicFilenameFactory: asicFilenameFactory,
	}
}

// TimestampDocuments timestamps a list of documents and returns the timestamped archive. Ports
// timestamp(List, ASiCWithCAdESTimestampParameters); Go has no overloading, so the Java overload
// set timestamp(List, ...) / timestamp(Content, ...) becomes TimestampDocuments / Timestamp.
func (s *ASiCWithCAdESTimestampService) TimestampDocuments(documents []model.DSSDocument,
	parameters *ASiCWithCAdESTimestampParameters) model.DSSDocument {
	asicContent := NewASiCWithCAdESASiCContentBuilder().
		Build(documents, parameters.ASiC().ContainerType())
	asicContent = s.Timestamp(asicContent, parameters)
	zipArchive, err := asic.ZipUtilsInstance().CreateZipArchiveAt(asicContent, parameters.ZipCreationDate())
	if err != nil {
		panic(err)
	}
	return zipArchive
}

// Timestamp adds a timestamp to the given Content, returning the content with the timestamp
// and the related XML Manifest for an ASiC-E container. Ports
// timestamp(Content, ASiCWithCAdESTimestampParameters).
func (s *ASiCWithCAdESTimestampService) Timestamp(asicContent *asic.Content,
	parameters *ASiCWithCAdESTimestampParameters) *asic.Content {
	dataToSignHelper := NewASiCWithCAdESTimestampDataToSignHelperBuilder(s.asicFilenameFactory).
		Build(asicContent, parameters)

	toBeTimestamped := dataToSignHelper.ToBeSigned()
	if enumerations.ASiCContainerTypeASiCE == parameters.ASiC().ContainerType() {
		// XML Document in case of ASiC-E container
		asicContent.SetManifestDocuments(append(asicContent.ManifestDocuments(), toBeTimestamped))
	}

	digestAlgorithm := parameters.DigestAlgorithm()
	digestValue, err := toBeTimestamped.DigestValue(digestAlgorithm)
	if err != nil {
		panic(err)
	}
	timestampBinary, err := s.tspSource.TimeStampResponse(digestAlgorithm, digestValue)
	if err != nil {
		panic(err)
	}

	derEncoded, err := spi.DSSASN1UtilsDEREncodedTimestampBinary(timestampBinary)
	if err != nil {
		panic(err)
	}
	timestampToken := model.NewInMemoryDocumentWithMimeType(derEncoded,
		s.asicFilenameFactory.TimestampFilename(asicContent), enumerations.MimeTypeEnumTST)
	asicContent.SetTimestampDocuments(asic.UtilsAddOrReplaceDocument(asicContent.TimestampDocuments(), timestampToken))

	return asicContent
}
