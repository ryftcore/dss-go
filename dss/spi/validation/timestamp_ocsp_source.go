// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/tsp/TimestampOCSPSource.java (DSS 6.5.RC1).
package validation

import (
	"github.com/ryftcore/dss-go/dss/internal/cmscore"
	"github.com/ryftcore/dss-go/dss/spi"
)

// TimestampOCSPSource is an OCSPSource that retrieves information embedded in a TimeStampToken.
type TimestampOCSPSource struct {
	*spi.CMSOCSPSource
}

// newTimestampOCSPSource extracts the OCSP responses and references of the given time-stamp
// token. Port of the package-private TimestampOCSPSource(TimeStampToken) constructor, whose
// only caller is TimestampToken.
//
// Upstream selects the two otherRevocationInfo stores by OID -
// CMSObjectIdentifiers.id_ri_ocsp_response and OCSPObjectIdentifiers.id_pkix_ocsp_basic - which
// cmscore exposes as OCSPResponses and OCSPBasicResponses (internal/cmscore/oids.go:
// OIDRIOCSPResponse, OIDPKIXOCSPBasic).
func newTimestampOCSPSource(timeStampToken *cmscore.TimeStampToken) (*TimestampOCSPSource, error) {
	cms := timeStampToken.CMS()
	signerInfo := spi.DSSASN1UtilsFirstSignerInformation(cms.SignerInfos())
	base, err := spi.NewCMSOCSPSource(cms.OCSPResponses(), cms.OCSPBasicResponses(),
		signerInfo.UnsignedAttributes)
	if err != nil {
		return nil, err
	}
	source := &TimestampOCSPSource{CMSOCSPSource: base}
	// The outermost concrete source registers itself, so that the RevocationToken dispatch of
	// OfflineRevocationSourceBase reaches OfflineOCSPSourceBase.RevocationTokens; this class
	// does not override it.
	source.InitOfflineRevocationSource(source)
	return source, nil
}
