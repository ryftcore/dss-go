// Ported from dss-xml-common/src/main/java/eu/europa/esig/dss/xml/common/SchemaFactoryBuilder.java (DSS 6.5.RC1).
package common

// SchemaFactory is a documented stub standing in for javax.xml.validation.SchemaFactory.
// XSD schema validation is out of scope for this phase (internal/xmldom's doc.go lists it
// among its explicit non-goals), so this type carries no state; see
// TransformerFactory's doc comment in transformer_factory_builder.go for the general
// rationale, which applies identically here.
type SchemaFactory struct {
	// SchemaLanguage records the schema language the factory was built for (Java:
	// SchemaFactory.newInstance(schemaLanguage)); it has no effect since no real schema
	// factory exists yet.
	SchemaLanguage string
}

// Schema is a documented stub standing in for javax.xml.validation.Schema, the compiled
// result of SchemaFactory#newSchema(Source[]).
type Schema struct{}

// Source is a documented stub standing in for javax.xml.transform.Source. dss-xml-common
// never constructs its content itself (SchemaFactory#newSchema(Source[]) only ever consumes
// it whole), so it carries no fields; a later phase that ports XSD schema validation
// replaces this with a real type without changing GetSchema's signature.
type Source struct{}

// xmlConstantsW3CXMLSchemaNSURI is javax.xml.XMLConstants.W3C_XML_SCHEMA_NS_URI, the default
// schema language.
const xmlConstantsW3CXMLSchemaNSURI = "http://www.w3.org/2001/XMLSchema"

// SchemaFactoryBuilder builds a SchemaFactory.
//
// # Documented stub
//
// SchemaFactory has no Go counterpart (see its doc comment), so this builder's
// setSecurityFeature/setSecurityAttribute are no-ops that always succeed: there is no real
// object to configure, and no real object to fail configuring either. Build() therefore
// always succeeds. The feature/attribute configuration machinery inherited from
// AbstractConfigurator is preserved and fully functional, ready for a later phase to give
// setSecurityFeature/setSecurityAttribute real behaviour without changing this type's
// shape.
type SchemaFactoryBuilder struct {
	*AbstractFactoryBuilder[*SchemaFactory]
	schemaLanguage string
}

// newSchemaFactoryBuilder builds a secure pre-configured instance. Ports the protected
// default constructor.
func newSchemaFactoryBuilder() *SchemaFactoryBuilder {
	b := &SchemaFactoryBuilder{schemaLanguage: xmlConstantsW3CXMLSchemaNSURI}
	b.AbstractFactoryBuilder = NewAbstractFactoryBuilder[*SchemaFactory](
		func() *SchemaFactory { return &SchemaFactory{SchemaLanguage: b.schemaLanguage} },
		b.setSecurityFeature,
		b.setSecurityAttribute,
	)
	b.EnableFeature(xmlConstantsFeatureSecureProcessing)
	b.SetAttribute(xmlConstantsAccessExternalDTD, "")
	b.SetAttribute(xmlConstantsAccessExternalSchema, "")
	return b
}

// GetSecureSchemaBuilder instantiates a pre-configured, secure SchemaFactoryBuilder. Ports
// the static factory getSecureSchemaBuilder().
func GetSecureSchemaBuilder() *SchemaFactoryBuilder {
	return newSchemaFactoryBuilder()
}

// Build builds the configured SchemaFactory. Ports build().
func (b *SchemaFactoryBuilder) Build() (*SchemaFactory, error) {
	factory := b.InstantiateFactory()
	if err := b.SetSecurityFeatures(factory); err != nil {
		return nil, err
	}
	if err := b.SetSecurityAttributes(factory); err != nil {
		return nil, err
	}
	return factory, nil
}

// SetSchemaLanguage sets the schema language to instantiate the SchemaFactory with. Ports
// setSchemaLanguage(String).
func (b *SchemaFactoryBuilder) SetSchemaLanguage(schemaLanguage string) {
	b.schemaLanguage = schemaLanguage
}

// EnableFeature overrides AbstractConfigurator.EnableFeature purely for the covariant
// return, exactly as Java's override does.
func (b *SchemaFactoryBuilder) EnableFeature(feature string) *SchemaFactoryBuilder {
	b.AbstractFactoryBuilder.EnableFeature(feature)
	return b
}

// DisableFeature overrides AbstractConfigurator.DisableFeature for the covariant return.
func (b *SchemaFactoryBuilder) DisableFeature(feature string) *SchemaFactoryBuilder {
	b.AbstractFactoryBuilder.DisableFeature(feature)
	return b
}

// SetAttribute overrides AbstractConfigurator.SetAttribute for the covariant return.
func (b *SchemaFactoryBuilder) SetAttribute(attribute string, value any) *SchemaFactoryBuilder {
	b.AbstractFactoryBuilder.SetAttribute(attribute, value)
	return b
}

// RemoveAttribute overrides AbstractConfigurator.RemoveAttribute for the covariant return.
func (b *SchemaFactoryBuilder) RemoveAttribute(attribute string) *SchemaFactoryBuilder {
	b.AbstractFactoryBuilder.RemoveAttribute(attribute)
	return b
}

// setSecurityFeature is a documented no-op (see the type doc comment).
func (b *SchemaFactoryBuilder) setSecurityFeature(factory *SchemaFactory, feature string, value bool) error {
	return nil
}

// setSecurityAttribute is a documented no-op (see the type doc comment).
func (b *SchemaFactoryBuilder) setSecurityAttribute(factory *SchemaFactory, attribute string, value any) error {
	return nil
}
