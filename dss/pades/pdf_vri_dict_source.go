// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/validation/dss/PdfVriDictSource.java (DSS 6.5.RC1).
//
// slf4j is dropped, per PORTING.md; the two upstream log statements are kept as comments where
// they fired.
package pades

import (
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/utils"
)

// PdfVriDictSource extracts special information from a VRI dictionary.
type PdfVriDictSource struct {
	// pdfVriDict is the VRI dictionary, nil when the name matched no single dictionary.
	pdfVriDict *PdfVriDict

	// vriDictionaryName is the identifier of the VRI dictionary.
	vriDictionaryName string
}

// NewPdfVriDictSource is the port of the PdfVriDictSource(PdfDssDict, String) constructor.
// vriDictionaryName is the SHA-1 of the signature name.
func NewPdfVriDictSource(dssDictionary PdfDssDict, vriDictionaryName string) *PdfVriDictSource {
	source := &PdfVriDictSource{vriDictionaryName: vriDictionaryName}
	vris := UtilsVRIsWithName(dssDictionary, vriDictionaryName)
	if utils.CollectionSize(vris) == 1 {
		source.pdfVriDict = vris[0]
	}
	return source
}

// VRICreationTime returns the VRI creation time extracted from the 'TU' field, nil when absent.
// Port of getVRICreationTime().
func (s *PdfVriDictSource) VRICreationTime() *time.Time {
	if s.pdfVriDict != nil {
		return s.pdfVriDict.TUTime()
	}
	return nil
}

// TimestampToken returns the time-stamp token extracted from the 'TS' field of the VRI
// dictionary, nil when there is none or when it cannot be parsed. Port of getTimestampToken().
//
// Java wraps the whole body in a try/catch that logs and returns null; in Go only the
// TimestampToken construction can fail, so that is the single branch the catch becomes.
func (s *PdfVriDictSource) TimestampToken() *validation.TimestampToken {
	if s.pdfVriDict != nil {
		tsStream := s.pdfVriDict.TSStream()
		if utils.IsArrayNotEmpty(tsStream) {
			identifierBuilder := NewVriDictionaryTimestampIdentifierBuilder(tsStream, s.vriDictionaryName)
			timestampToken, err := validation.NewTimestampTokenWithIdentifierBuilder(s.pdfVriDict.TSStream(),
				enumerations.TimestampTypeVRITimestamp, []*validation.TimestampedReference{}, identifierBuilder)
			if err != nil {
				// Upstream logs "An error occurred while extracting 'TS' timestamp from the
				// corresponding /VRI dictionary : {}" at warn level.
				return nil
			}
			return timestampToken
		}
	}
	return nil
}
