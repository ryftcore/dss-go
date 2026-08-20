// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/function/converter/CriteriaListConverter.java (DSS 6.5.RC1).
//
// otherCriteriaList (eu.europa.esig.xades.jaxb.xades132.AnyType, outside this manifest) is bound
// by dss/trustedlist/jaxb as an opaque captured token stream (unexported xadesAnyType) rather
// than dispatched through the "lax" xs:any mechanism AnyType/ExtensionType use elsewhere - see
// jaxb_common.go's CriteriaListType header. This port therefore re-decodes the three concrete
// element shapes CriteriaListConverter's Java code type-switches on
// (CertSubjectDNAttributeType/ExtendedKeyUsageType/QcStatementListType) directly off the
// captured token stream by their known global element names (CertSubjectDNAttribute,
// ExtendedKeyUsage, QcStatementSet - see dss/trustedlist/jaxb's wildcardElements table), instead
// of a JAXBElement type switch.
package tsl

import (
	"encoding/xml"
	"io"

	"github.com/utain/esig/dss/enumerations"
	tslmodel "github.com/utain/esig/dss/model/tsl"
	"github.com/utain/esig/dss/spi"
	"github.com/utain/esig/dss/trustedlist/jaxb"
	"github.com/utain/esig/dss/utils"
)

// CriteriaListConverter converts a CriteriaListType to a Condition.
type CriteriaListConverter struct{}

// NewCriteriaListConverter is the default constructor. Port of CriteriaListConverter().
func NewCriteriaListConverter() *CriteriaListConverter {
	return &CriteriaListConverter{}
}

// Apply ports apply(CriteriaListType).
func (c *CriteriaListConverter) Apply(criteriaList *jaxb.CriteriaListType) tslmodel.Condition {
	var matchingCriteriaIndicator enumerations.Assert
	if criteriaList.Assert != nil {
		matchingCriteriaIndicator = enumerations.Assert(*criteriaList.Assert)
	}
	condition := NewCompositeConditionWithMatchingCriteriaIndicator(matchingCriteriaIndicator)

	c.addKeyUsageConditionsIfPresent(criteriaList.KeyUsage, condition)
	c.addPolicyIdConditionsIfPresent(criteriaList.PolicySet, condition)
	if criteriaList.OtherCriteriaList != nil {
		c.addOtherCriteriaListConditionsIfPresent(criteriaList.OtherCriteriaList.Tokens, condition)
	}
	c.addCriteriaListConditionsIfPresent(criteriaList.CriteriaList, condition)

	return condition
}

func (c *CriteriaListConverter) addKeyUsageConditionsIfPresent(keyUsages []*jaxb.KeyUsageType, criteriaCondition *CompositeCondition) {
	if utils.IsCollectionNotEmpty(keyUsages) {
		for _, keyUsageType := range keyUsages {
			condition := NewCompositeCondition()
			for _, keyUsageBit := range keyUsageType.KeyUsageBit {
				var bit enumerations.KeyUsageBit
				if keyUsageBit.Name != nil {
					bit = enumerations.KeyUsageBit(*keyUsageBit.Name)
				}
				condition.AddChild(NewKeyUsageCondition(bit, keyUsageBit.Value))
			}
			criteriaCondition.AddChild(condition)
		}
	}
}

func (c *CriteriaListConverter) addPolicyIdConditionsIfPresent(policySet []*jaxb.PoliciesListType, criteriaCondition *CompositeCondition) {
	if utils.IsCollectionNotEmpty(policySet) {
		for _, policiesListType := range policySet {
			condition := NewCompositeCondition()
			for _, oidType := range policiesListType.PolicyIdentifier {
				if id, ok := objectIdentifierTypeValue(oidType); ok && utils.IsStringNotEmpty(id) {
					condition.AddChild(NewPolicyIdCondition(id))
				}
			}
			criteriaCondition.AddChild(condition)
		}
	}
}

// addOtherCriteriaListConditionsIfPresent ports the private
// addOtherCriteriaListConditionsIfPresent(AnyType, CompositeCondition), ETSI TS 119 612 V1.1.1 /
// 5.5.9.2.2.3 - see this file's header for the re-decoding strategy.
func (c *CriteriaListConverter) addOtherCriteriaListConditionsIfPresent(tokens []xml.Token, condition *CompositeCondition) {
	for _, start := range topLevelStartElements(tokens) {
		switch start.start.Name.Local {
		case "CertSubjectDNAttribute":
			var v jaxb.CertSubjectDNAttributeType
			if err := decodeTokenSubtree(tokens[start.index:], &start.start, &v); err == nil {
				condition.AddChild(NewCertSubjectDNAttributeCondition(c.extractOids(v.AttributeOID)))
			}
		case "ExtendedKeyUsage":
			var v jaxb.ExtendedKeyUsageType
			if err := decodeTokenSubtree(tokens[start.index:], &start.start, &v); err == nil {
				condition.AddChild(NewExtendedKeyUsageCondition(c.extractOids(v.KeyPurposeId)))
			}
		case "QcStatementSet":
			var v jaxb.QcStatementListType
			if err := decodeTokenSubtree(tokens[start.index:], &start.start, &v); err == nil {
				composite := NewCompositeConditionWithMatchingCriteriaIndicator(enumerations.Assert_ALL)
				for _, qcStatementType := range v.QcStatement {
					var oid, legislation, typ string
					if id, ok := objectIdentifierTypeValue(qcStatementType.QcStatementId); ok {
						oid = id
					}
					if qcStatementType.QcStatementInfo != nil {
						if qcStatementType.QcStatementInfo.QcCClegislation != nil {
							legislation = *qcStatementType.QcStatementInfo.QcCClegislation
						}
						if qcStatementType.QcStatementInfo.QcType != nil {
							if id, ok := objectIdentifierTypeValue(qcStatementType.QcStatementInfo.QcType); ok && utils.IsStringNotEmpty(id) {
								typ = id
							}
						}
					}
					composite.AddChild(NewQCStatementCondition(oid, typ, legislation))
				}
				condition.AddChild(composite)
			}
		}
	}
}

// extractOids ports the private extractOids(List<ObjectIdentifierType>).
func (c *CriteriaListConverter) extractOids(oits []*jaxb.ObjectIdentifierType) []string {
	var oids []string
	if utils.IsCollectionNotEmpty(oits) {
		for _, objectIdentifierType := range oits {
			if id, ok := objectIdentifierTypeValue(objectIdentifierType); ok && utils.IsStringNotEmpty(id) {
				if spi.DSSUtilsIsOidCode(id) {
					oids = append(oids, id)
				}
			}
		}
	}
	return oids
}

func (c *CriteriaListConverter) addCriteriaListConditionsIfPresent(criteriaList []*jaxb.CriteriaListType, condition *CompositeCondition) {
	if utils.IsCollectionNotEmpty(criteriaList) {
		for _, criteriaListType := range criteriaList {
			condition.AddChild(c.Apply(criteriaListType))
		}
	}
}

// objectIdentifierTypeValue extracts the OID value from a raw-captured ObjectIdentifierType,
// mirroring DSSUtils.getObjectIdentifierValue(identifier.getValue(), identifier.getQualifier())
// applied to its nested xades132 Identifier element. The qualifier argument is a documentation-
// only parameter in this port's DSSUtilsObjectIdentifierValueWithQualifier (see dss_utils.go),
// so it is not separately extracted here.
func objectIdentifierTypeValue(o *jaxb.ObjectIdentifierType) (string, bool) {
	if o == nil {
		return "", false
	}
	value, ok := findChildElementText(o.Tokens, "Identifier")
	if !ok {
		return "", false
	}
	return spi.DSSUtilsObjectIdentifierValueWithQualifier(value, ""), true
}

// findChildElementText scans a captured token stream for a top-level child element with the
// given local name and returns its character-data content.
func findChildElementText(tokens []xml.Token, localName string) (string, bool) {
	for _, start := range topLevelStartElements(tokens) {
		if start.start.Name.Local != localName {
			continue
		}
		var text string
		depth := 0
		for i := start.index + 1; i < len(tokens); i++ {
			switch t := tokens[i].(type) {
			case xml.StartElement:
				depth++
			case xml.EndElement:
				if depth == 0 {
					return text, true
				}
				depth--
			case xml.CharData:
				if depth == 0 {
					text += string(t)
				}
			}
		}
		return text, true
	}
	return "", false
}

// topLevelStartElement pairs a top-level (depth 0) StartElement of a captured token stream with
// its index in that stream.
type topLevelStartElement struct {
	start xml.StartElement
	index int
}

// topLevelStartElements walks a captured token stream (foreignContent.Tokens) and returns every
// direct-child StartElement, skipping over nested descendants.
func topLevelStartElements(tokens []xml.Token) []topLevelStartElement {
	var result []topLevelStartElement
	depth := 0
	for i, tok := range tokens {
		switch t := tok.(type) {
		case xml.StartElement:
			if depth == 0 {
				result = append(result, topLevelStartElement{start: t, index: i})
			}
			depth++
		case xml.EndElement:
			depth--
		}
	}
	return result
}

// tokenSliceReader replays a captured []xml.Token as an xml.TokenReader, so a subtree already
// captured by foreignContent.unmarshal can be re-decoded into a concrete struct via
// xml.NewTokenDecoder + Decoder.DecodeElement.
type tokenSliceReader struct {
	tokens []xml.Token
	pos    int
}

// Token returns the next captured token.
func (r *tokenSliceReader) Token() (xml.Token, error) {
	if r.pos >= len(r.tokens) {
		return nil, io.EOF
	}
	t := r.tokens[r.pos]
	r.pos++
	return t, nil
}

// decodeTokenSubtree decodes the subtree starting at start (already the first element of
// tokensFromStart) into target.
func decodeTokenSubtree(tokensFromStart []xml.Token, start *xml.StartElement, target any) error {
	dec := xml.NewTokenDecoder(&tokenSliceReader{tokens: tokensFromStart})
	return dec.DecodeElement(target, start)
}
