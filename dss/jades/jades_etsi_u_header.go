// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/validation/JAdESEtsiUHeader.java (DSS 6.5.RC1).
//
// See etsi_u_component.go's file header for the JAdESAttribute/JAdESAttributeIdentifier forward
// dependency this file also relies on (JAdESAttribute.Equals, used by ReplaceComponent).
//
// Java's `unprotected.computeIfAbsent(ETSI_U, k -> new JSONArray())` relies on List reference
// semantics: the returned JSONArray is the SAME object stored in the map, so
// add/remove/ListIterator#set mutate the map's entry in place. A Go []any returned from
// jose.Object.Value does not alias the field it came from (see
// spi/validation/timestamp/signature_timestamp_source.go's identical, already-documented gap for
// TimestampToken.TimestampedReferences), so every mutator below writes the modified slice back
// via unprotected.Put before returning, achieving the same externally observable effect.
//
// DEVIATION: AddComponent/RemoveComponent/ReplaceComponent return an error (always nil today -
// none of these mutations can currently fail) and RemoveLastComponent returns (bool, error),
// rather than Java's void/boolean-only signatures. This matches the call-site contract the
// already-landed SIGN chunk (jades_level_baseline_lt.go, jades_level_baseline_lta.go,
// jades_level_baseline_t.go, jades_counter_signature_builder.go,
// jades_signature_policy_store_builder.go) was written against, presumably to leave room for a
// future fallible mutation without a breaking signature change.
package jades

import (
	"github.com/ryftcore/dss-go/dss/internal/jose"
	"github.com/ryftcore/dss-go/dss/utils"
)

// JAdESEtsiUHeader represents the list of components present inside the unprotected 'etsiU'
// header. Port of the class JAdESEtsiUHeader, implementing validation.SignatureProperties[
// *EtsiUComponent].
//
// java.io.Serializable is dropped (no Go counterpart).
type JAdESEtsiUHeader struct {
	// jws is the JWS signature.
	jws *JWS

	// components is the list of 'etsiU' components.
	components []*EtsiUComponent
}

// NewJAdESEtsiUHeader is the default constructor. Port of the (JWS) constructor.
func NewJAdESEtsiUHeader(jws *JWS) *JAdESEtsiUHeader {
	return &JAdESEtsiUHeader{jws: jws}
}

// IsExist implements validation.SignatureProperties. Port of isExist().
func (h *JAdESEtsiUHeader) IsExist() bool {
	return utils.IsCollectionNotEmpty(h.Attributes())
}

// Attributes implements validation.SignatureProperties. Port of getAttributes().
func (h *JAdESEtsiUHeader) Attributes() []*EtsiUComponent {
	if h.components == nil {
		components := []*EtsiUComponent{}
		etsiUContent := DSSJsonUtilsEtsiU(h.jws)
		if utils.IsCollectionNotEmpty(etsiUContent) {
			for ii, item := range etsiUContent {
				if etsiUComponent := EtsiUComponentBuild(item, ii); etsiUComponent != nil {
					components = append(components, etsiUComponent)
				}
				// else : unable to create, skip
			}
		}
		h.components = components
	}
	return h.components
}

// AddComponent adds a new entry to the 'etsiU' array. Port of addComponent(String, Object,
// boolean). See the file header's DEVIATION note on the error return.
func (h *JAdESEtsiUHeader) AddComponent(headerName string, value any, base64UrlEncoded bool) error {
	etsiU := h.getEtsiUToEdit()
	etsiEntry := jadesEtsiUHeaderComponent(headerName, value, base64UrlEncoded)
	etsiU = append(etsiU, etsiEntry)
	h.setEtsiU(etsiU)
	return nil
}

// jadesEtsiUHeaderComponent ports the private getComponent(String, Object, boolean).
func jadesEtsiUHeaderComponent(name string, value any, base64UrlEncoded bool) any {
	jsonObject := NewJsonObject()
	jsonObject.Put(name, value)
	if base64UrlEncoded {
		return DSSJsonUtilsToBase64UrlObject(jsonObject)
	}
	return jsonObject
}

// RemoveComponent removes the 'etsiU' components with the given headerName. Port of
// removeComponent(String). See the file header's DEVIATION note on the error return.
func (h *JAdESEtsiUHeader) RemoveComponent(headerName string) error {
	etsiU := h.getEtsiUToEdit()
	if utils.IsCollectionNotEmpty(etsiU) {
		kept := make([]any, 0, len(etsiU))
		for _, item := range etsiU {
			if jadesEtsiUHeaderMatches(item, headerName) {
				continue
			}
			kept = append(kept, item)
		}
		h.setEtsiU(kept)
	}
	return nil
}

// RemoveLastComponent removes the last 'etsiU' item if the name matches one of the given
// headerName values. Port of removeLastComponent(String...), returning whether the component has
// been removed. See the file header's DEVIATION note on the error return.
func (h *JAdESEtsiUHeader) RemoveLastComponent(headerName ...string) (bool, error) {
	etsiU := h.getEtsiUToEdit()
	if utils.IsCollectionNotEmpty(etsiU) {
		last := etsiU[len(etsiU)-1]
		for _, name := range headerName {
			if jadesEtsiUHeaderMatches(last, name) {
				h.setEtsiU(etsiU[:len(etsiU)-1])
				return true, nil
			}
		}
	}
	return false, nil
}

// jadesEtsiUHeaderMatches ports the private removeLastIfMatches(ListIterator, String...)'s
// matching predicate (`Arrays.stream(headerName).anyMatch(etsiUComponent::containsKey)`), applied
// to a single etsiU entry rather than a backward ListIterator: RemoveComponent/RemoveLastComponent
// above reproduce the iteration pattern themselves.
func jadesEtsiUHeaderMatches(item any, headerName ...string) bool {
	m, ok := DSSJsonUtilsParseEtsiUComponent(item)
	if !ok || m == nil {
		return false
	}
	for _, name := range headerName {
		if m.ContainsKey(name) {
			return true
		}
	}
	return false
}

// ReplaceComponent replaces the given attribute within the 'etsiU' header array. Port of
// replaceComponent(EtsiUComponent). See the file header's DEVIATION note on the error return.
func (h *JAdESEtsiUHeader) ReplaceComponent(attribute *EtsiUComponent) error {
	etsiU := h.getEtsiUToEdit()
	for position, item := range etsiU {
		currentComponent := EtsiUComponentBuild(item, position)
		if currentComponent != nil && attribute.Equals(&currentComponent.JAdESAttribute) {
			etsiU[position] = attribute.Component()
			break
		}
	}
	h.setEtsiU(etsiU)
	return nil
}

// getEtsiUToEdit ports the private getEtsiUToEdit(), including
// `unprotected.computeIfAbsent(ETSI_U, k -> new JSONArray())`'s side effect of creating the
// 'unprotected' header map when absent (the Go counterpart of jws.setUnprotected(new HashMap<>())
// above it). See the file header on why callers must write the result back via setEtsiU.
func (h *JAdESEtsiUHeader) getEtsiUToEdit() []any {
	unprotected := h.jws.Unprotected()
	if unprotected == nil {
		unprotected = jose.NewObject()
		h.jws.SetUnprotected(unprotected)
	}
	h.clearCachedAttributes()

	etsiU, ok := unprotected.Value(JAdESHeaderParameterNamesEtsiU).([]any)
	if !ok {
		etsiU = []any{}
	}
	return etsiU
}

// setEtsiU writes etsiU back to the 'unprotected' header map under the 'etsiU' key. See the file
// header.
func (h *JAdESEtsiUHeader) setEtsiU(etsiU []any) {
	h.jws.Unprotected().Put(JAdESHeaderParameterNamesEtsiU, etsiU)
}

// clearCachedAttributes ports the private clearCachedAttributes().
func (h *JAdESEtsiUHeader) clearCachedAttributes() {
	h.components = nil
}
