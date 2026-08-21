// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/tsp/TimestampCRLSource.java (DSS 6.5.RC1).
package validation

import (
	"github.com/ryftcore/dss-go/dss/internal/cmscore"
	"github.com/ryftcore/dss-go/dss/spi"
)

// TimestampCRLSource is a CRLSource that retrieves information embedded in a TimeStampToken.
type TimestampCRLSource struct {
	*spi.CMSCRLSource
}

// newTimestampCRLSource extracts the CRLs and CRL references of the given time-stamp token.
// Port of the package-private TimestampCRLSource(TimeStampToken) constructor, whose only caller
// is TimestampToken.
//
// Java passes timeStampToken.toCMSSignedData().getCRLs() and getUnsignedAttributes(); the latter
// is the unsignedAttrs of the TSA signer, which is absent for almost every time-stamp - a nil
// cmscore.Attributes, which is what spi.NewCMSCRLSource reads as Java's null AttributeTable.
func newTimestampCRLSource(timeStampToken *cmscore.TimeStampToken) (*TimestampCRLSource, error) {
	signerInfo := spi.DSSASN1UtilsFirstSignerInformation(timeStampToken.CMS().SignerInfos())
	base, err := spi.NewCMSCRLSource(timeStampToken.CMS().CRLs(), signerInfo.UnsignedAttributes)
	if err != nil {
		return nil, err
	}
	source := &TimestampCRLSource{CMSCRLSource: base}
	// The outermost concrete source registers itself, so that the RevocationToken dispatch of
	// OfflineRevocationSourceBase reaches OfflineCRLSourceBase.RevocationTokens; this class
	// does not override it.
	source.InitOfflineRevocationSource(source)
	return source, nil
}
