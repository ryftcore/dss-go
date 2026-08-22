// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/OCSPFirstRevocationDataLoadingStrategyFactory.java (DSS 6.5.RC1).
//
// OCSPFirstRevocationDataLoadingStrategy embeds RevocationDataLoadingStrategy exactly like
// CRLFirstRevocationDataLoadingStrategy does, so Create() below returns the embedded
// *RevocationDataLoadingStrategy, letting callers dispatch to the right concrete
// RevocationToken() implementation uniformly.
package validation

// OCSPFirstRevocationDataLoadingStrategyFactory initializes an
// OCSPFirstRevocationDataLoadingStrategy.
type OCSPFirstRevocationDataLoadingStrategyFactory struct{}

// NewOCSPFirstRevocationDataLoadingStrategyFactory is the default constructor.
func NewOCSPFirstRevocationDataLoadingStrategyFactory() *OCSPFirstRevocationDataLoadingStrategyFactory {
	return &OCSPFirstRevocationDataLoadingStrategyFactory{}
}

// Create returns a new OCSPFirstRevocationDataLoadingStrategy instance. Port of create().
func (f *OCSPFirstRevocationDataLoadingStrategyFactory) Create() *RevocationDataLoadingStrategy {
	s := NewOCSPFirstRevocationDataLoadingStrategy()
	return &s.RevocationDataLoadingStrategy
}
