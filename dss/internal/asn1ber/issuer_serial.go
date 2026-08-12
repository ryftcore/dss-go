package asn1ber

import "math/big"

// GeneralName is one alternative of the X.509 GeneralName CHOICE, reduced to what
// IssuerSerial needs. TagNo is the context-specific tag number (4 for directoryName) and
// Name is the DER encoding of the chosen alternative's value.
type GeneralName struct {
	// TagNo is the CHOICE alternative, e.g. 4 for directoryName.
	TagNo int
	// Name is the DER encoding of the value, e.g. an X.501 Name SEQUENCE for tag 4.
	Name []byte
}

// DER returns the DER encoding of the GeneralName.
//
// directoryName (tag 4) wraps a CHOICE and is therefore explicitly tagged, as
// org.bouncycastle.asn1.x509.GeneralName encodes it; every other alternative is implicitly
// tagged, i.e. the value's identifier octet is replaced by the context-specific tag.
func (g *GeneralName) DER() []byte {
	if g.TagNo == 4 {
		return WriteTLV(ClassContextSpecific|Constructed|byte(g.TagNo), g.Name)
	}
	element, _, err := Parse(g.Name)
	if err != nil {
		return nil
	}
	identifier := byte(ClassContextSpecific) | byte(g.TagNo)
	if element.IsConstructed() {
		identifier |= Constructed
	}
	return WriteTLV(identifier, element.DERContent())
}

// IssuerSerial is the X.509 attribute-certificate structure
//
//	IssuerSerial ::= SEQUENCE {
//	    issuer   GeneralNames,
//	    serial   CertificateSerialNumber,
//	    issuerUID UniqueIdentifier OPTIONAL }
//
// replacing org.bouncycastle.asn1.x509.IssuerSerial. Only the two fields DSS reads are
// modelled; a third field found while parsing is preserved in IssuerUID.
type IssuerSerial struct {
	// Issuer is the GeneralNames sequence identifying the issuer.
	Issuer []GeneralName
	// Serial is the certificate serial number.
	Serial *big.Int
	// IssuerUID is the DER encoding of the optional issuerUID BIT STRING, nil when absent.
	IssuerUID []byte
}

// DER returns the DER encoding of the IssuerSerial.
func (i *IssuerSerial) DER() []byte {
	var names []byte
	for index := range i.Issuer {
		names = append(names, i.Issuer[index].DER()...)
	}
	body := WriteSequence(names)
	body = append(body, EncodeInteger(i.Serial)...)
	if i.IssuerUID != nil {
		body = append(body, i.IssuerUID...)
	}
	return WriteSequence(body)
}
