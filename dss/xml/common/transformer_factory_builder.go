// Ported from dss-xml-common/src/main/java/eu/europa/esig/dss/xml/common/TransformerFactoryBuilder.java (DSS 6.5.RC1).
package common

// TransformerFactory is a documented stub standing in for
// javax.xml.transform.TransformerFactory. Go has no XSLT engine and none is planned (not a
// project goal; internal/xmldom's doc.go lists XSLT among its explicit non-goals), so this
// type carries no state. It exists so TransformerFactoryBuilder's public contract - the
// type name a caller would hold a value of - is ready for a later phase, without an API
// break, if XSLT-adjacent behaviour is ever needed (DSS itself uses TransformerFactory only
// for javax.xml.transform.Transformer-based serialization, which internal/xmldom.Serialize
// already replaces without XSLT).
type TransformerFactory struct{}

// The JAXP feature/attribute URIs TransformerFactoryBuilder's constructor sets.
const (
	xmlConstantsFeatureSecureProcessing  = "http://javax.xml.XMLConstants/feature/secure-processing"
	xmlConstantsAccessExternalStylesheet = "http://javax.xml.XMLConstants/property/accessExternalStylesheet"
)

// TransformerFactoryBuilder builds a TransformerFactory.
//
// # Documented stub
//
// TransformerFactory has no Go counterpart (see its doc comment), so this builder's
// setSecurityFeature/setSecurityAttribute are no-ops that always succeed: there is no real
// object to configure, and no real object to fail configuring either. Build() therefore
// always succeeds and returns an empty *TransformerFactory. The feature/attribute
// configuration machinery inherited from AbstractConfigurator (EnableFeature, SetAttribute,
// etc.) is preserved and fully functional, ready for a later phase to give
// setSecurityFeature/setSecurityAttribute real behaviour without changing this type's shape.
type TransformerFactoryBuilder struct {
	*AbstractFactoryBuilder[*TransformerFactory]
}

// newTransformerFactoryBuilder builds a secure pre-configured instance. Ports the protected
// default constructor.
func newTransformerFactoryBuilder() *TransformerFactoryBuilder {
	b := &TransformerFactoryBuilder{}
	b.AbstractFactoryBuilder = NewAbstractFactoryBuilder[*TransformerFactory](
		func() *TransformerFactory { return &TransformerFactory{} },
		b.setSecurityFeature,
		b.setSecurityAttribute,
	)
	b.EnableFeature(xmlConstantsFeatureSecureProcessing)
	b.SetAttribute(xmlConstantsAccessExternalDTD, "")
	b.SetAttribute(xmlConstantsAccessExternalStylesheet, "")
	return b
}

// GetSecureTransformerBuilder instantiates a pre-configured, secure
// TransformerFactoryBuilder. Ports the static factory getSecureTransformerBuilder().
func GetSecureTransformerBuilder() *TransformerFactoryBuilder {
	return newTransformerFactoryBuilder()
}

// Build builds the configured TransformerFactory. Ports build().
func (b *TransformerFactoryBuilder) Build() (*TransformerFactory, error) {
	factory := b.InstantiateFactory()
	if err := b.SetSecurityFeatures(factory); err != nil {
		return nil, err
	}
	if err := b.SetSecurityAttributes(factory); err != nil {
		return nil, err
	}
	return factory, nil
}

// EnableFeature overrides AbstractConfigurator.EnableFeature purely for the covariant
// return, exactly as Java's override does.
func (b *TransformerFactoryBuilder) EnableFeature(feature string) *TransformerFactoryBuilder {
	b.AbstractFactoryBuilder.EnableFeature(feature)
	return b
}

// DisableFeature overrides AbstractConfigurator.DisableFeature for the covariant return.
func (b *TransformerFactoryBuilder) DisableFeature(feature string) *TransformerFactoryBuilder {
	b.AbstractFactoryBuilder.DisableFeature(feature)
	return b
}

// SetAttribute overrides AbstractConfigurator.SetAttribute for the covariant return.
func (b *TransformerFactoryBuilder) SetAttribute(attribute string, value any) *TransformerFactoryBuilder {
	b.AbstractFactoryBuilder.SetAttribute(attribute, value)
	return b
}

// RemoveAttribute overrides AbstractConfigurator.RemoveAttribute for the covariant return.
func (b *TransformerFactoryBuilder) RemoveAttribute(attribute string) *TransformerFactoryBuilder {
	b.AbstractFactoryBuilder.RemoveAttribute(attribute)
	return b
}

// setSecurityFeature is a documented no-op (see the type doc comment).
func (b *TransformerFactoryBuilder) setSecurityFeature(factory *TransformerFactory, feature string, value bool) error {
	return nil
}

// setSecurityAttribute is a documented no-op (see the type doc comment).
func (b *TransformerFactoryBuilder) setSecurityAttribute(factory *TransformerFactory, attribute string, value any) error {
	return nil
}
