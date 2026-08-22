// Ported from the generated JAXB adapters Adapter1 and Adapter2 of
// eu.europa.esig.dss.policy.jaxb, which delegate to
// eu.europa.esig.dss.jaxb.parsers.LevelParser and
// eu.europa.esig.dss.jaxb.parsers.ValidationModelParser respectively. As in
// dss/diagnostic/jaxb, each adapter becomes a named Go type over the
// corresponding enumerations constant so that encoding/xml marshals it through
// encoding.TextMarshaler with the exact lexical form the parser prints (here,
// simply the Java enum's name()), and rejects any other lexical form on the way
// in, mirroring Level.valueOf/ValidationModel.valueOf throwing
// IllegalArgumentException for an unknown name.
package jaxb

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
)

// LevelValue is the Level adapter: Adapter1 print()/parse() map to the Java
// enum name() (eu.europa.esig.dss.enumerations.Level has no field/override
// beyond the default toString()).
type LevelValue enumerations.Level

// Level returns the underlying enumeration constant.
func (v LevelValue) Level() enumerations.Level {
	return enumerations.Level(v)
}

// MarshalText writes the lexical form produced by the Java parser.
func (v LevelValue) MarshalText() ([]byte, error) {
	return []byte(string(v)), nil
}

// UnmarshalText resolves the lexical form back to the enumeration constant.
func (v *LevelValue) UnmarshalText(text []byte) error {
	s := string(text)
	level, err := enumerations.LevelValueOf(s)
	if err != nil {
		return err
	}
	*v = LevelValue(level)
	return nil
}

// ValidationModelValue is the ValidationModel adapter: Adapter2 print()/parse()
// map to the Java enum name().
type ValidationModelValue enumerations.ValidationModel

// ValidationModel returns the underlying enumeration constant.
func (v ValidationModelValue) ValidationModel() enumerations.ValidationModel {
	return enumerations.ValidationModel(v)
}

// MarshalText writes the lexical form produced by the Java parser.
func (v ValidationModelValue) MarshalText() ([]byte, error) {
	return []byte(string(v)), nil
}

// UnmarshalText resolves the lexical form back to the enumeration constant.
//
// dss/enumerations (frozen - see PORTING.md) does not export a
// ValidationModelValueOf lookup alongside ValidationModel, unlike most other
// enumerations types, so the lookup is inlined here over
// enumerations.ValidationModelValues() instead of adding one to that frozen
// package.
func (v *ValidationModelValue) UnmarshalText(text []byte) error {
	s := string(text)
	for _, candidate := range enumerations.ValidationModelValues() {
		if string(candidate) == s {
			*v = ValidationModelValue(candidate)
			return nil
		}
	}
	return fmt.Errorf("no enum constant ValidationModel.%s", s)
}

// TimeUnit is the Go form of the generated JAXB enum TimeUnit
// (eu.europa.esig.dss.policy.jaxb.TimeUnit, a plain @XmlEnum with no
// @XmlEnumValue overrides, so JAXB binds it by name() directly - unlike Level
// and ValidationModel above, no XmlAdapter is generated for it).
//
// This is a schema-local enum, distinct from java.util.concurrent.TimeUnit
// (which RuleUtils.convertDuration converts to/from - see ../rule_utils.go)
// and from any eu.europa.esig.dss.enumerations type.
type TimeUnit string

const (
	TimeUnitDays         TimeUnit = "DAYS"
	TimeUnitHours        TimeUnit = "HOURS"
	TimeUnitMinutes      TimeUnit = "MINUTES"
	TimeUnitSeconds      TimeUnit = "SECONDS"
	TimeUnitMilliseconds TimeUnit = "MILLISECONDS"
)

// TimeUnitValues returns all TimeUnit constants in declaration order.
func TimeUnitValues() []TimeUnit {
	return []TimeUnit{
		TimeUnitDays,
		TimeUnitHours,
		TimeUnitMinutes,
		TimeUnitSeconds,
		TimeUnitMilliseconds,
	}
}

// TimeUnitValueOf returns the TimeUnit matching the given Java enum name.
// Ports the generated TimeUnit#fromValue, itself Enum#valueOf.
func TimeUnitValueOf(name string) (TimeUnit, error) {
	for _, v := range TimeUnitValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", fmt.Errorf("no enum constant TimeUnit.%s", name)
}

// Value returns the id code of the referred value. Ports TimeUnit#value.
func (t TimeUnit) Value() string {
	return string(t)
}
