// Ported from dss-xml-common/src/main/java/eu/europa/esig/dss/xml/common/AbstractFactoryBuilder.java (DSS 6.5.RC1).
package common

// AbstractFactoryBuilder builds a secure F instance.
type AbstractFactoryBuilder[F any] struct {
	*AbstractConfigurator[F]

	// instantiateFactory stands in for the abstract method instantiateFactory(); see
	// abstract_configurator.go's SecurityFeatureSetter doc for why Go injects it instead of
	// overriding it.
	instantiateFactory func() F
}

// NewAbstractFactoryBuilder creates an AbstractFactoryBuilder given the constructor
// function that stands in for the overridden instantiateFactory() and the setters that
// stand in for setSecurityFeature/setSecurityAttribute.
func NewAbstractFactoryBuilder[F any](instantiateFactory func() F, setFeature SecurityFeatureSetter[F], setAttribute SecurityAttributeSetter[F]) *AbstractFactoryBuilder[F] {
	return &AbstractFactoryBuilder[F]{
		AbstractConfigurator: NewAbstractConfigurator[F](setFeature, setAttribute),
		instantiateFactory:   instantiateFactory,
	}
}

// InstantiateFactory instantiates the corresponding factory. Ports the protected abstract
// method instantiateFactory().
//
// NOTE: in order to change the default behavior, the injected instantiateFactory function
// given to NewAbstractFactoryBuilder should be replaced.
func (b *AbstractFactoryBuilder[F]) InstantiateFactory() F {
	return b.instantiateFactory()
}
