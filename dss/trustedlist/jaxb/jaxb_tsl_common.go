// Ported from ts_119612v020401_xsd.xsd (DSS 6.5.RC1) via the JAXB classes
// generated into eu.europa.esig.trustedlist.jaxb.tsl. Per the phase-8a
// generated-JAXB rule the generated classes are grouped into schema-area
// files rather than one file per class; every Java class keeps its name and
// its exact field order so that encoding/xml reproduces the JAXB element
// sequence byte for byte. This file holds the scheme-information leaf types
// and the extensibility list types every TSL area (SchemeExtensions,
// TSPInformationExtensions, ServiceInformationExtensions) shares.
package jaxb

import "encoding/xml"

// MultiLangStringType is the Go form of the generated JAXB class
// MultiLangStringType (complexType MultiLangStringType, @XmlValue + xml:lang).
type MultiLangStringType struct {
	Value string `xml:",chardata"`
	Lang  string `xml:"http://www.w3.org/XML/1998/namespace lang,attr"`
}

// MultiLangNormStringType is the Go form of the generated JAXB class
// MultiLangNormStringType (complexType MultiLangNormStringType): like
// MultiLangStringType but bound through NormalizedStringAdapter, which only
// replaces literal tab/CR/LF with a space and does not otherwise change the
// lexical form the RI writes back - a plain xs:string round-trips it
// unchanged for every value this port marshals (values it only reads back
// from an already-normalized document).
type MultiLangNormStringType struct {
	Value string `xml:",chardata"`
	Lang  string `xml:"http://www.w3.org/XML/1998/namespace lang,attr"`
}

// InternationalNamesType is the Go form of the generated JAXB class
// InternationalNamesType (complexType InternationalNamesType).
type InternationalNamesType struct {
	Name []*MultiLangNormStringType `xml:"Name"`
}

// AttributedNonEmptyURIType is the Go form of the generated JAXB class
// AttributedNonEmptyURIType (complexType AttributedNonEmptyURIType).
type AttributedNonEmptyURIType struct {
	Value string  `xml:",chardata"`
	Type  *string `xml:"type,attr,omitempty"`
}

// NonEmptyMultiLangURIType is the Go form of the generated JAXB class
// NonEmptyMultiLangURIType (complexType NonEmptyMultiLangURIType).
type NonEmptyMultiLangURIType struct {
	Value string `xml:",chardata"`
	Lang  string `xml:"http://www.w3.org/XML/1998/namespace lang,attr"`
}

// NonEmptyMultiLangURIListType is the Go form of the generated JAXB class
// NonEmptyMultiLangURIListType (complexType NonEmptyMultiLangURIListType).
type NonEmptyMultiLangURIListType struct {
	URI []*NonEmptyMultiLangURIType `xml:"URI"`
}

// NonEmptyURIListType is the Go form of the generated JAXB class
// NonEmptyURIListType (complexType NonEmptyURIListType).
type NonEmptyURIListType struct {
	URI []string `xml:"URI"`
}

// PostalAddressType is the Go form of the generated JAXB class
// PostalAddressType (complexType PostalAddressType).
type PostalAddressType struct {
	StreetAddress   string  `xml:"StreetAddress"`
	Locality        string  `xml:"Locality"`
	StateOrProvince *string `xml:"StateOrProvince,omitempty"`
	PostalCode      *string `xml:"PostalCode,omitempty"`
	CountryName     string  `xml:"CountryName"`
	Lang            string  `xml:"http://www.w3.org/XML/1998/namespace lang,attr"`
}

// PostalAddressListType is the Go form of the generated JAXB class
// PostalAddressListType (complexType PostalAddressListType).
type PostalAddressListType struct {
	PostalAddress []*PostalAddressType `xml:"PostalAddress"`
}

// ElectronicAddressType is the Go form of the generated JAXB class
// ElectronicAddressType (complexType ElectronicAddressType).
type ElectronicAddressType struct {
	URI []*NonEmptyMultiLangURIType `xml:"URI"`
}

// AddressType is the Go form of the generated JAXB class AddressType
// (complexType AddressType).
type AddressType struct {
	PostalAddresses   *PostalAddressListType `xml:"PostalAddresses"`
	ElectronicAddress *ElectronicAddressType `xml:"ElectronicAddress"`
}

// PolicyOrLegalnoticeType is the Go form of the generated JAXB class
// PolicyOrLegalnoticeType (complexType PolicyOrLegalnoticeType).
type PolicyOrLegalnoticeType struct {
	TSLPolicy      []*NonEmptyMultiLangURIType `xml:"TSLPolicy"`
	TSLLegalNotice []*MultiLangStringType      `xml:"TSLLegalNotice"`
}

// NextUpdateType is the Go form of the generated JAXB class NextUpdateType
// (complexType NextUpdateType). DateTime is the raw xs:dateTime lexical
// form, round-tripped verbatim (see jaxb_tsl_root.go's header for why this
// port does not model java.time semantics at this layer).
type NextUpdateType struct {
	DateTime *string `xml:"dateTime,omitempty"`
}

// OtherTSLPointerType is the Go form of the generated JAXB class
// OtherTSLPointerType (complexType OtherTSLPointerType).
type OtherTSLPointerType struct {
	ServiceDigitalIdentities *ServiceDigitalIdentityListType `xml:"ServiceDigitalIdentities,omitempty"`
	TSLLocation              string                          `xml:"TSLLocation"`
	AdditionalInformation    *AdditionalInformationType      `xml:"AdditionalInformation,omitempty"`
}

// OtherTSLPointersType is the Go form of the generated JAXB class
// OtherTSLPointersType (complexType OtherTSLPointersType).
type OtherTSLPointersType struct {
	OtherTSLPointer []*OtherTSLPointerType `xml:"OtherTSLPointer"`
}

// AdditionalInformationItem is one item of AdditionalInformationType's
// choice content (@XmlElements: TextualInformation/MultiLangStringType or
// OtherInformation/AnyType) - the same choice-of-named-elements pattern as
// dss/diagnostic/jaxb's CertificateExtensionsWrapper.
type AdditionalInformationItem struct {
	// TextualInformation is set for a "TextualInformation" item.
	TextualInformation *MultiLangStringType
	// OtherInformation is set for an "OtherInformation" item.
	OtherInformation *AnyType
}

// AdditionalInformationType is the Go form of the generated JAXB class
// AdditionalInformationType (complexType AdditionalInformationType). Items
// drives its own MarshalXML/UnmarshalXML; see AdditionalInformationItem.
type AdditionalInformationType struct {
	Items []AdditionalInformationItem
}

// UnmarshalXML builds each item from the element name it arrived under.
func (a *AdditionalInformationType) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for {
		tok, err := d.Token()
		if err != nil {
			return err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			t = t.Copy()
			switch t.Name.Local {
			case "TextualInformation":
				var v MultiLangStringType
				if err := d.DecodeElement(&v, &t); err != nil {
					return err
				}
				a.Items = append(a.Items, AdditionalInformationItem{TextualInformation: &v})
			case "OtherInformation":
				var v AnyType
				if err := d.DecodeElement(&v, &t); err != nil {
					return err
				}
				a.Items = append(a.Items, AdditionalInformationItem{OtherInformation: &v})
			default:
				if err := d.Skip(); err != nil {
					return err
				}
			}
		case xml.EndElement:
			return nil
		}
	}
}

// MarshalXML writes each item under its own element name, in order.
func (a *AdditionalInformationType) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	for _, it := range a.Items {
		switch {
		case it.TextualInformation != nil:
			if err := e.EncodeElement(it.TextualInformation, xml.StartElement{Name: xml.Name{Local: "TextualInformation"}}); err != nil {
				return err
			}
		case it.OtherInformation != nil:
			if err := e.EncodeElement(it.OtherInformation, xml.StartElement{Name: xml.Name{Local: "OtherInformation"}}); err != nil {
				return err
			}
		}
	}
	return e.EncodeToken(start.End())
}

// AdditionalServiceInformationType is the Go form of the generated JAXB
// class AdditionalServiceInformationType (complexType
// AdditionalServiceInformationType).
type AdditionalServiceInformationType struct {
	URI              *NonEmptyMultiLangURIType `xml:"URI"`
	InformationValue *string                   `xml:"InformationValue,omitempty"`
	OtherInformation *AnyType                  `xml:"OtherInformation,omitempty"`
}

// ExtensionsListType is the Go form of the generated JAXB class
// ExtensionsListType (complexType ExtensionsListType): the container of
// SchemeExtensions/TSPInformationExtensions/ServiceInformationExtensions.
type ExtensionsListType struct {
	Extension []*ExtensionType `xml:"Extension"`
}

// ServiceSupplyPointsType is the Go form of the generated JAXB class
// ServiceSupplyPointsType (complexType ServiceSupplyPointsType).
type ServiceSupplyPointsType struct {
	ServiceSupplyPoint []*AttributedNonEmptyURIType `xml:"ServiceSupplyPoint"`
}

// DigitalIdentityType is the Go form of the generated JAXB class
// DigitalIdentityType (complexType DigitalIdentityType): one member of a
// DigitalIdentityListType choice (X509Certificate/X509SubjectName/
// ds:KeyValue/X509SKI/Other - the schema's own comment calls this "choice"
// but xjc still binds it as ordinary optional elements since at most one is
// ever populated on any instance the RI produces).
type DigitalIdentityType struct {
	X509Certificate *Base64Binary `xml:"X509Certificate,omitempty"`
	X509SubjectName *string       `xml:"X509SubjectName,omitempty"`
	KeyValue        *dsigKeyValue `xml:"http://www.w3.org/2000/09/xmldsig# KeyValue,omitempty"`
	X509SKI         *Base64Binary `xml:"X509SKI,omitempty"`
	Other           *AnyType      `xml:"Other,omitempty"`
}

// DigitalIdentityListType is the Go form of the generated JAXB class
// DigitalIdentityListType (complexType DigitalIdentityListType).
type DigitalIdentityListType struct {
	DigitalId []*DigitalIdentityType `xml:"DigitalId"`
}

// ServiceDigitalIdentityListType is the Go form of the generated JAXB class
// ServiceDigitalIdentityListType (complexType ServiceDigitalIdentityListType,
// OtherTSLPointerType.ServiceDigitalIdentities).
type ServiceDigitalIdentityListType struct {
	ServiceDigitalIdentity []*DigitalIdentityListType `xml:"ServiceDigitalIdentity"`
}
