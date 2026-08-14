// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/tsl/TrustProperties.java (DSS 6.5.RC1).
package tsl

import "github.com/utain/esig/dss/model/timedependent"

// TrustProperties contains the Trust properties for certificates.
//
// java.io.Serializable has no Go counterpart and is dropped.
type TrustProperties struct {
	// lotlInfo is the LOTL id.
	lotlInfo *LOTLInfo
	// tlInfo is the TL id.
	tlInfo *TLInfo
	// trustServiceProvider is the trustServiceProvider.
	trustServiceProvider *TrustServiceProvider
	// trustService is the trustService.
	trustService *timedependent.TimeDependentValues[*TrustServiceStatusAndInformationExtensions]
}

// NewTrustProperties creates a TrustProperties object for extracted information from an
// "independent" trusted list (Java's (TLInfo, TrustServiceProvider, TimeDependentValues)
// constructor, which delegates to the LOTL-aware one with a nil lotlInfo).
//
// Panics with the Java message ("tlInfo cannot be null!", "trustServiceProvider cannot be
// null!", "trustService cannot be null!") when the respective argument is nil, mirroring
// Objects.requireNonNull.
func NewTrustProperties(tlInfo *TLInfo, trustServiceProvider *TrustServiceProvider,
	trustService *timedependent.TimeDependentValues[*TrustServiceStatusAndInformationExtensions]) *TrustProperties {
	return NewTrustPropertiesWithLOTL(nil, tlInfo, trustServiceProvider, trustService)
}

// NewTrustPropertiesWithLOTL creates a TrustProperties object linked to a LOTL with MRA.
//
// Panics with the Java message ("tlInfo cannot be null!", "trustServiceProvider cannot be
// null!", "trustService cannot be null!") when the respective argument is nil, mirroring
// Objects.requireNonNull.
func NewTrustPropertiesWithLOTL(lotlInfo *LOTLInfo, tlInfo *TLInfo, trustServiceProvider *TrustServiceProvider,
	trustService *timedependent.TimeDependentValues[*TrustServiceStatusAndInformationExtensions]) *TrustProperties {
	if tlInfo == nil {
		panic("tlInfo cannot be null!")
	}
	if trustServiceProvider == nil {
		panic("trustServiceProvider cannot be null!")
	}
	if trustService == nil {
		panic("trustService cannot be null!")
	}
	return &TrustProperties{
		lotlInfo:             lotlInfo,
		tlInfo:               tlInfo,
		trustServiceProvider: trustServiceProvider,
		trustService:         trustService,
	}
}

// LOTLInfo gets the LOTL.
func (t *TrustProperties) LOTLInfo() *LOTLInfo {
	return t.lotlInfo
}

// TLInfo gets the TL.
func (t *TrustProperties) TLInfo() *TLInfo {
	return t.tlInfo
}

// TrustServiceProvider gets the trust service provider.
func (t *TrustProperties) TrustServiceProvider() *TrustServiceProvider {
	return t.trustServiceProvider
}

// TrustService gets the trust service.
func (t *TrustProperties) TrustService() *timedependent.TimeDependentValues[*TrustServiceStatusAndInformationExtensions] {
	return t.trustService
}
