// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/RevocationDataLoadingStrategyFactory.java (DSS 6.5.RC1).
//
// CROSS-CHUNK DEPENDENCY: RevocationDataLoadingStrategy (the abstract class returned here) is
// owned by sibling chunk VAL-D, per the header comment of crl_first_revocation_data_loading_strategy.go
// and ocsp_first_revocation_data_loading_strategy_factory.go, both already landed in this package.
package validation

// RevocationDataLoadingStrategyFactory is used to initialize a new RevocationDataLoadingStrategy.
type RevocationDataLoadingStrategyFactory interface {
	// Create initializes a new RevocationDataLoadingStrategy. Port of create().
	Create() *RevocationDataLoadingStrategy
}
