// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/CRLFirstRevocationDataLoadingStrategyFactory.java (DSS 6.5.RC1).
package validation

// CRLFirstRevocationDataLoadingStrategyFactory initializes a CRLFirstRevocationDataLoadingStrategy.
type CRLFirstRevocationDataLoadingStrategyFactory struct{}

// NewCRLFirstRevocationDataLoadingStrategyFactory is the default constructor.
func NewCRLFirstRevocationDataLoadingStrategyFactory() *CRLFirstRevocationDataLoadingStrategyFactory {
	return &CRLFirstRevocationDataLoadingStrategyFactory{}
}

// Create returns a new CRLFirstRevocationDataLoadingStrategy instance. Port of create().
func (f *CRLFirstRevocationDataLoadingStrategyFactory) Create() *RevocationDataLoadingStrategy {
	s := NewCRLFirstRevocationDataLoadingStrategy()
	return &s.RevocationDataLoadingStrategy
}

// compile-time assertion: a CRLFirstRevocationDataLoadingStrategyFactory is a
// RevocationDataLoadingStrategyFactory.
var _ RevocationDataLoadingStrategyFactory = (*CRLFirstRevocationDataLoadingStrategyFactory)(nil)
