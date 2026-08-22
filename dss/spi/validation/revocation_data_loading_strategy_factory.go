// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/RevocationDataLoadingStrategyFactory.java (DSS 6.5.RC1).
package validation

// RevocationDataLoadingStrategyFactory is used to initialize a new RevocationDataLoadingStrategy.
type RevocationDataLoadingStrategyFactory interface {
	// Create initializes a new RevocationDataLoadingStrategy. Port of create().
	Create() *RevocationDataLoadingStrategy
}
