// Ported from ts_119612v020401_xsd.xsd (DSS 6.5.RC1) via the JAXB classes
// generated into eu.europa.esig.trustedlist.jaxb.tsl. The generated classes
// are grouped into schema-area files rather than one file per class; every
// Java class keeps its name and its exact field order so that encoding/xml
// reproduces the JAXB element sequence byte for byte. This file holds the
// TSP service and service- history types (TSPServicesListType and
// everything it nests).
package jaxb

// TSPServiceInformationType is the Go form of the generated JAXB class
// TSPServiceInformationType (complexType TSPServiceInformationType).
// StatusStartingTime is the raw xs:dateTime lexical form, round-tripped
// verbatim - see jaxb_tsl_root.go's header.
type TSPServiceInformationType struct {
	ServiceTypeIdentifier        string                        `xml:"ServiceTypeIdentifier"`
	ServiceName                  *InternationalNamesType       `xml:"ServiceName"`
	ServiceDigitalIdentity       *DigitalIdentityListType      `xml:"ServiceDigitalIdentity"`
	ServiceStatus                string                        `xml:"ServiceStatus"`
	StatusStartingTime           *string                       `xml:"StatusStartingTime,omitempty"`
	SchemeServiceDefinitionURI   *NonEmptyMultiLangURIListType `xml:"SchemeServiceDefinitionURI,omitempty"`
	ServiceSupplyPoints          *ServiceSupplyPointsType      `xml:"ServiceSupplyPoints,omitempty"`
	TSPServiceDefinitionURI      *NonEmptyMultiLangURIListType `xml:"TSPServiceDefinitionURI,omitempty"`
	ServiceInformationExtensions *ExtensionsListType           `xml:"ServiceInformationExtensions,omitempty"`
}

// TSPServiceType is the Go form of the generated JAXB class TSPServiceType
// (complexType TSPServiceType).
type TSPServiceType struct {
	ServiceInformation *TSPServiceInformationType `xml:"ServiceInformation"`
	ServiceHistory     *ServiceHistoryType        `xml:"ServiceHistory,omitempty"`
}

// TSPServicesListType is the Go form of the generated JAXB class
// TSPServicesListType (complexType TSPServicesListType).
type TSPServicesListType struct {
	TSPService []*TSPServiceType `xml:"TSPService"`
}

// ServiceHistoryInstanceType is the Go form of the generated JAXB class
// ServiceHistoryInstanceType (complexType ServiceHistoryInstanceType).
type ServiceHistoryInstanceType struct {
	ServiceTypeIdentifier        string                   `xml:"ServiceTypeIdentifier"`
	ServiceName                  *InternationalNamesType  `xml:"ServiceName"`
	ServiceDigitalIdentity       *DigitalIdentityListType `xml:"ServiceDigitalIdentity"`
	ServiceStatus                string                   `xml:"ServiceStatus"`
	StatusStartingTime           *string                  `xml:"StatusStartingTime,omitempty"`
	ServiceInformationExtensions *ExtensionsListType      `xml:"ServiceInformationExtensions,omitempty"`
}

// ServiceHistoryType is the Go form of the generated JAXB class
// ServiceHistoryType (complexType ServiceHistoryType).
type ServiceHistoryType struct {
	ServiceHistoryInstance []*ServiceHistoryInstanceType `xml:"ServiceHistoryInstance"`
}
