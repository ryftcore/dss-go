// Ported from dss-enumerations/.../TSLTypeEnum.java (DSS 6.5.RC1).
package enumerations

// TSLTypeEnum defines common TSLType values supported by the
// implementation.
type TSLTypeEnum string

const (
	// TSLTypeEnum_EUlistofthelists is the EU List of the Trusted Lists.
	TSLTypeEnum_EUlistofthelists TSLTypeEnum = "EUlistofthelists"
	// TSLTypeEnum_EUgeneric is the EU Trusted Lists.
	TSLTypeEnum_EUgeneric TSLTypeEnum = "EUgeneric"
	// TSLTypeEnum_AdESlistofthelists is the AdES List of the Trusted Lists.
	TSLTypeEnum_AdESlistofthelists TSLTypeEnum = "AdESlistofthelists"
)

type tslTypeEnumFields struct {
	uri   string
	label string
}

// tslTypeEnumData holds the (uri, label) tuple for each constant.
var tslTypeEnumData = map[TSLTypeEnum]tslTypeEnumFields{
	TSLTypeEnum_EUlistofthelists:   {"http://uri.etsi.org/TrstSvc/TrustedList/TSLType/EUlistofthelists", "EU List of the Trusted Lists"},
	TSLTypeEnum_EUgeneric:          {"http://uri.etsi.org/TrstSvc/TrustedList/TSLType/EUgeneric", "EU Trusted List"},
	TSLTypeEnum_AdESlistofthelists: {"http://ec.europa.eu/tools/lotl/mra/ades-lotl-tsl-type", "AdES List of the Trusted Lists"},
}

// TSLTypeEnumValues returns all constants in declaration order.
func TSLTypeEnumValues() []TSLTypeEnum {
	return []TSLTypeEnum{
		TSLTypeEnum_EUlistofthelists,
		TSLTypeEnum_EUgeneric,
		TSLTypeEnum_AdESlistofthelists,
	}
}

// URI returns the URI associated with the TSPType.
func (t TSLTypeEnum) URI() string {
	return tslTypeEnumData[t].uri
}

// Label returns the name of the TSLType.
func (t TSLTypeEnum) Label() string {
	return tslTypeEnumData[t].label
}

// TSLTypeEnumValueOf returns the constant matching the given Java enum
// name.
func TSLTypeEnumValueOf(name string) (TSLTypeEnum, error) {
	for _, v := range TSLTypeEnumValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", &tslTypeEnumInvalidValueError{name}
}

type tslTypeEnumInvalidValueError struct {
	name string
}

func (e *tslTypeEnumInvalidValueError) Error() string {
	return "no enum constant TSLTypeEnum." + e.name
}

// Compile-time assertion that TSLTypeEnum implements TSLType.
var _ TSLType = TSLTypeEnum_EUgeneric
