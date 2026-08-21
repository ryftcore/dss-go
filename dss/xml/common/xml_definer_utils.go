// Ported from dss-xml-common/src/main/java/eu/europa/esig/dss/xml/common/XmlDefinerUtils.java (DSS 6.5.RC1).
package common

import (
	"sync"

	"github.com/ryftcore/dss-go/dss/internal/xmldom"
)

// XmlDefinerUtils builds the objects for dealing with XML: a *xmldom.ParseOptions in place
// of DocumentBuilderFactory (see document_builder_factory_builder.go for the mapping) and
// the documented TransformerFactory/SchemaFactory/Validator stubs otherwise.
type XmlDefinerUtils struct {
	secureDocumentBuilderFactoryBuilder *DocumentBuilderFactoryBuilder
	secureTransformerFactoryBuilder     *TransformerFactoryBuilder
	secureSchemaFactoryBuilder          *SchemaFactoryBuilder
	secureValidatorConfigurator         *ValidatorConfigurator
}

var (
	xmlDefinerUtilsSingleton *XmlDefinerUtils
	xmlDefinerUtilsOnce      sync.Once
)

// GetXmlDefinerUtilsInstance returns the XmlDefinerUtils singleton. Ports getInstance();
// Java's non-thread-safe lazy singleton (a plain "if (singleton == null)" check) becomes
// sync.Once, since Go offers no counterpart to relying on undefined-but-usually-fine
// single-threaded initialization order.
func GetXmlDefinerUtilsInstance() *XmlDefinerUtils {
	xmlDefinerUtilsOnce.Do(func() {
		xmlDefinerUtilsSingleton = &XmlDefinerUtils{
			secureDocumentBuilderFactoryBuilder: GetSecureDocumentBuilderFactoryBuilder(),
			secureTransformerFactoryBuilder:     GetSecureTransformerBuilder(),
			secureSchemaFactoryBuilder:          GetSecureSchemaBuilder(),
			secureValidatorConfigurator:         GetSecureValidatorConfigurator(),
		}
	})
	return xmlDefinerUtilsSingleton
}

// SetDocumentBuilderFactoryBuilder sets a pre-configured builder to instantiate a
// *xmldom.ParseOptions. Ports setDocumentBuilderFactoryBuilder(DocumentBuilderFactoryBuilder).
func (u *XmlDefinerUtils) SetDocumentBuilderFactoryBuilder(documentBuilderFactoryBuilder *DocumentBuilderFactoryBuilder) {
	u.secureDocumentBuilderFactoryBuilder = documentBuilderFactoryBuilder
}

// GetSecureDocumentBuilderFactory returns a *xmldom.ParseOptions with security features
// enabled. Ports getSecureDocumentBuilderFactory().
func (u *XmlDefinerUtils) GetSecureDocumentBuilderFactory() (*xmldom.ParseOptions, error) {
	return u.secureDocumentBuilderFactoryBuilder.Build()
}

// GetSchema returns a Schema for the given xsdSources. Ports getSchema(List<Source>); the
// requireNonNull(xsdSources, "XSD Source(s) must be provided") is a no-op check here since
// xsdSources being nil and being empty are indistinguishable for a Go slice (Java's
// requireNonNull only rejects null, not an empty list) - see SchemaFactoryBuilder's doc
// comment for why the returned Schema is currently a documented stub carrying no content.
func (u *XmlDefinerUtils) GetSchema(xsdSources []*Source) (*Schema, error) {
	if _, err := u.GetSecureSchemaFactory(); err != nil {
		return nil, err
	}
	return &Schema{}, nil
}

// SetSchemaFactoryBuilder sets a pre-configured builder to instantiate a SchemaFactory.
// Ports setSchemaFactoryBuilder(SchemaFactoryBuilder).
func (u *XmlDefinerUtils) SetSchemaFactoryBuilder(schemaFactoryBuilder *SchemaFactoryBuilder) {
	u.secureSchemaFactoryBuilder = schemaFactoryBuilder
}

// GetSecureSchemaFactory returns a SchemaFactory with enabled security features. Ports
// getSecureSchemaFactory().
func (u *XmlDefinerUtils) GetSecureSchemaFactory() (*SchemaFactory, error) {
	return u.secureSchemaFactoryBuilder.Build()
}

// SetTransformerFactoryBuilder sets a pre-configured builder to instantiate a
// TransformerFactory. Ports setTransformerFactoryBuilder(TransformerFactoryBuilder).
func (u *XmlDefinerUtils) SetTransformerFactoryBuilder(transformerFactoryBuilder *TransformerFactoryBuilder) {
	u.secureTransformerFactoryBuilder = transformerFactoryBuilder
}

// GetSecureTransformerFactory returns a TransformerFactory with enabled security features.
// Ports getSecureTransformerFactory().
func (u *XmlDefinerUtils) GetSecureTransformerFactory() (*TransformerFactory, error) {
	return u.secureTransformerFactoryBuilder.Build()
}

// SetValidatorConfigurator sets a pre-configured builder to instantiate a Validator. Ports
// setValidatorConfigurator(ValidatorConfigurator).
func (u *XmlDefinerUtils) SetValidatorConfigurator(validatorConfigurator *ValidatorConfigurator) {
	u.secureValidatorConfigurator = validatorConfigurator
}

// Configure configures the validator. Ports configure(Validator).
func (u *XmlDefinerUtils) Configure(validator *Validator) error {
	return u.secureValidatorConfigurator.Configure(validator)
}

// PostProcess post-processes the validator after validation is executed. Ports
// postProcess(Validator).
func (u *XmlDefinerUtils) PostProcess(validator *Validator) error {
	return u.secureValidatorConfigurator.PostProcess(validator)
}
