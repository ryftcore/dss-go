// Ported from dss-xml-common/src/main/java/eu/europa/esig/dss/xml/common/DocumentBuilderFactoryBuilder.java (DSS 6.5.RC1).
package common

import "github.com/utain/esig/dss/internal/xmldom"

// The JAXP feature URIs DocumentBuilderFactoryBuilder's constructor toggles.
const (
	documentBuilderFeatureNamespaces              = "http://xml.org/sax/features/namespaces"
	documentBuilderFeatureCreateEntityRefNodes    = "http://apache.org/xml/features/dom/create-entity-ref-nodes"
	documentBuilderFeatureDisallowDoctypeDecl     = "http://apache.org/xml/features/disallow-doctype-decl"
	documentBuilderFeatureExternalGeneralEntities = "http://xml.org/sax/features/external-general-entities"
	documentBuilderFeatureExternalParamEntities   = "http://xml.org/sax/features/external-parameter-entities"
	documentBuilderFeatureLoadExternalDTD         = "http://apache.org/xml/features/nonvalidating/load-external-dtd"
	xmlConstantsAccessExternalDTD                 = "http://javax.xml.XMLConstants/property/accessExternalDTD"
	xmlConstantsAccessExternalSchema              = "http://javax.xml.XMLConstants/property/accessExternalSchema"
)

// DocumentBuilderFactoryBuilder builds a *xmldom.ParseOptions configured with DSS's secure
// parser defaults.
//
// # Mapping onto internal/xmldom
//
// Upstream configures javax.xml.parsers.DocumentBuilderFactory/Xerces through JAXP feature
// and attribute URIs. xmldom has no JAXP layer; xmldom.ParseOptions already encodes the
// equivalent secure posture directly as plain Go fields, and several of the JAXP toggles
// have no xmldom counterpart at all, because the behaviour they turn off does not exist in
// xmldom to begin with. This table is the complete, exhaustive mapping - every
// feature/attribute the constructor sets is accounted for:
//
//	JAXP feature/attribute                                              xmldom mapping
//	http://xml.org/sax/features/namespaces (enabled)                    no-op: xmldom is namespace-aware unconditionally (xmldom/doc.go)
//	http://apache.org/xml/features/dom/create-entity-ref-nodes (enabled) no-op: xmldom has no DTD/entity layer to create such nodes from (Parse §1.5 non-goals)
//	http://apache.org/xml/features/disallow-doctype-decl (enabled)      ParseOptions.AllowDoctype = false (already the zero-value default)
//	http://xml.org/sax/features/external-general-entities (disabled)    no-op: xmldom never resolves external entities (no DTD support at all)
//	http://xml.org/sax/features/external-parameter-entities (disabled)  no-op: same
//	http://apache.org/xml/features/nonvalidating/load-external-dtd (disabled) no-op: same
//	XMLConstants.ACCESS_EXTERNAL_DTD = ""                                no-op: same
//	XMLConstants.ACCESS_EXTERNAL_SCHEMA = ""                             no-op: same
//
// Every toggle above is therefore a documented no-op except disallow-doctype-decl, which
// maps onto the one ParseOptions field that exists for it. Unlike upstream, an unrecognized
// feature/attribute name given to EnableFeature/DisableFeature/SetAttribute is also a
// silent no-op rather than a logged warning (slf4j is dropped per PORTING.md), mirroring
// DocumentBuilderFactoryBuilder#setSecurityFeature/#setSecurityAttribute, which upstream
// catches and logs locally rather than escalating through SecurityConfigurationException -
// so this builder's Build never fails.
type DocumentBuilderFactoryBuilder struct {
	*AbstractFactoryBuilder[*xmldom.ParseOptions]
}

// newDocumentBuilderFactoryBuilder builds a secure pre-configured instance. Ports the
// protected default constructor.
func newDocumentBuilderFactoryBuilder() *DocumentBuilderFactoryBuilder {
	b := &DocumentBuilderFactoryBuilder{}
	b.AbstractFactoryBuilder = NewAbstractFactoryBuilder[*xmldom.ParseOptions](
		func() *xmldom.ParseOptions { return &xmldom.ParseOptions{} },
		b.setSecurityFeature,
		b.setSecurityAttribute,
	)
	b.EnableFeature(documentBuilderFeatureNamespaces)           // .setNamespaceAware(true)
	b.EnableFeature(documentBuilderFeatureCreateEntityRefNodes) // .setExpandEntityReferences(false)
	b.EnableFeature(documentBuilderFeatureDisallowDoctypeDecl)
	b.DisableFeature(documentBuilderFeatureExternalGeneralEntities)
	b.DisableFeature(documentBuilderFeatureExternalParamEntities)
	b.DisableFeature(documentBuilderFeatureLoadExternalDTD)
	b.SetAttribute(xmlConstantsAccessExternalDTD, "")
	b.SetAttribute(xmlConstantsAccessExternalSchema, "")
	return b
}

// GetSecureDocumentBuilderFactoryBuilder instantiates a pre-configured, secure
// DocumentBuilderFactoryBuilder. Ports the static factory
// getSecureDocumentBuilderFactoryBuilder().
func GetSecureDocumentBuilderFactoryBuilder() *DocumentBuilderFactoryBuilder {
	return newDocumentBuilderFactoryBuilder()
}

// Build builds the configured *xmldom.ParseOptions. Ports build().
func (b *DocumentBuilderFactoryBuilder) Build() (*xmldom.ParseOptions, error) {
	parseOptions := b.InstantiateFactory()
	if err := b.SetSecurityFeatures(parseOptions); err != nil {
		return nil, err
	}
	if err := b.SetSecurityAttributes(parseOptions); err != nil {
		return nil, err
	}
	return parseOptions, nil
}

// EnableFeature overrides AbstractConfigurator.EnableFeature purely for the covariant
// return, exactly as Java's override does.
func (b *DocumentBuilderFactoryBuilder) EnableFeature(feature string) *DocumentBuilderFactoryBuilder {
	b.AbstractFactoryBuilder.EnableFeature(feature)
	return b
}

// DisableFeature overrides AbstractConfigurator.DisableFeature for the covariant return.
func (b *DocumentBuilderFactoryBuilder) DisableFeature(feature string) *DocumentBuilderFactoryBuilder {
	b.AbstractFactoryBuilder.DisableFeature(feature)
	return b
}

// SetAttribute overrides AbstractConfigurator.SetAttribute for the covariant return.
func (b *DocumentBuilderFactoryBuilder) SetAttribute(attribute string, value any) *DocumentBuilderFactoryBuilder {
	b.AbstractFactoryBuilder.SetAttribute(attribute, value)
	return b
}

// RemoveAttribute overrides AbstractConfigurator.RemoveAttribute for the covariant return.
func (b *DocumentBuilderFactoryBuilder) RemoveAttribute(attribute string) *DocumentBuilderFactoryBuilder {
	b.AbstractFactoryBuilder.RemoveAttribute(attribute)
	return b
}

// setSecurityFeature maps a JAXP feature URI onto *xmldom.ParseOptions per the table in the
// type doc comment. It never fails (see that comment for why), matching upstream, whose
// setSecurityFeature catches ParserConfigurationException locally and only logs it.
func (b *DocumentBuilderFactoryBuilder) setSecurityFeature(factory *xmldom.ParseOptions, feature string, value bool) error {
	if feature == documentBuilderFeatureDisallowDoctypeDecl {
		factory.AllowDoctype = !value
	}
	// Every other recognized feature is a documented no-op (see the type doc comment);
	// an unrecognized feature is silently ignored, matching upstream's local catch+log.
	return nil
}

// setSecurityAttribute is a documented no-op for both attributes this builder sets (see the
// type doc comment); it never fails, matching upstream's local catch+log.
func (b *DocumentBuilderFactoryBuilder) setSecurityAttribute(factory *xmldom.ParseOptions, attribute string, value any) error {
	return nil
}
