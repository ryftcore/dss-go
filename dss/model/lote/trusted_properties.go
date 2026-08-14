// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/lote/TrustedProperties.java (DSS 6.5.RC1).
package lote

import "github.com/utain/esig/dss/model/timedependent"

// TrustedProperties contains a list of trusted certificates and their properties.
//
// java.io.Serializable has no Go counterpart and is dropped.
type TrustedProperties struct {
	// loloteInfo is the LoLoTE id.
	loloteInfo *LoLoTEInfo
	// listInfo is the LoTE id.
	listInfo *LoTEInfo
	// trustedEntity is the trustedEntity.
	trustedEntity *TrustedEntity
	// trustedServices is the current trust service.
	trustedServices *timedependent.TimeDependentValues[ServiceStatusAndInformationExtensions]
}

// NewTrustedProperties creates a TrustedProperties object for extracted information from an
// "independent" list.
func NewTrustedProperties(listInfo *LoTEInfo, trustedEntity *TrustedEntity,
	trustedServices *timedependent.TimeDependentValues[ServiceStatusAndInformationExtensions]) *TrustedProperties {
	return NewTrustedPropertiesWithLoLoTE(nil, listInfo, trustedEntity, trustedServices)
}

// NewTrustedPropertiesWithLoLoTE creates a TrustedProperties object for extracted information
// with a related List of Lists.
//
// Panics with "tlInfo cannot be null!", "trustedEntity cannot be null!", or "trustedServices
// cannot be null!" when the respective argument is nil, mirroring Objects.requireNonNull (the
// Java message names "tlInfo" even though the parameter is listInfo, kept verbatim).
func NewTrustedPropertiesWithLoLoTE(loloteInfo *LoLoTEInfo, listInfo *LoTEInfo, trustedEntity *TrustedEntity,
	trustedServices *timedependent.TimeDependentValues[ServiceStatusAndInformationExtensions]) *TrustedProperties {
	if listInfo == nil {
		panic("tlInfo cannot be null!")
	}
	if trustedEntity == nil {
		panic("trustedEntity cannot be null!")
	}
	if trustedServices == nil {
		panic("trustedServices cannot be null!")
	}
	return &TrustedProperties{
		loloteInfo:      loloteInfo,
		listInfo:        listInfo,
		trustedEntity:   trustedEntity,
		trustedServices: trustedServices,
	}
}

// LoLoTEInfo gets LoLoTE.
func (t *TrustedProperties) LoLoTEInfo() *LoLoTEInfo {
	return t.loloteInfo
}

// LoTEInfo gets List.
func (t *TrustedProperties) LoTEInfo() *LoTEInfo {
	return t.listInfo
}

// TrustedEntity gets trusted entity.
func (t *TrustedProperties) TrustedEntity() *TrustedEntity {
	return t.trustedEntity
}

// TrustedServices gets trust service.
func (t *TrustedProperties) TrustedServices() *timedependent.TimeDependentValues[ServiceStatusAndInformationExtensions] {
	return t.trustedServices
}
