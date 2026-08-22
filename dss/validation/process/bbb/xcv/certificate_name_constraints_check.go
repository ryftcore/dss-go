// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/checks/CertificateNameConstraintsCheck.java (DSS 6.5.RC1).
//
// Java's slf4j logging (LOG.debug/LOG.warn) has no Go equivalent and is not ported.
package xcv

import (
	"strings"

	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	diagjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// CertificateNameConstraintsCheck verifies the validity of the certificate in
// regard to "Name constraint" certificate extension's value in its
// certificate chain.
//
// NOTE: only directoryName general name type is supported by this class.
type CertificateNameConstraintsCheck struct {
	*process.ChainItemBase[*jaxb.XmlSubXCV]

	// certificate is the certificate to check.
	certificate *diagnostic.CertificateWrapper
}

// NewCertificateNameConstraintsCheck is the default constructor. Port of
// CertificateNameConstraintsCheck(I18nProvider, XmlSubXCV, CertificateWrapper, LevelRule).
func NewCertificateNameConstraintsCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSubXCV],
	certificate *diagnostic.CertificateWrapper, constraint policy.LevelRule) *CertificateNameConstraintsCheck {
	c := &CertificateNameConstraintsCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		certificate:   certificate,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *CertificateNameConstraintsCheck) Process() bool {
	/*
	 * 6.1.2. Initialization
	 */
	certificateChain := append([]*diagnostic.CertificateWrapper{c.certificate}, c.certificate.CertificateChain()...)
	/*
	 * (b) permitted_subtrees: a set of root names for each name type
	 * (e.g., X.500 distinguished names, email addresses, or IP
	 * addresses) defining a set of subtrees within which all
	 * subject names in subsequent certificates in the certification
	 * path MUST fall. This variable includes a set for each name
	 * type, and the initial value is initial-permitted-subtrees.
	 */
	var permittedSubtrees map[enumerations.GeneralNameType][]*diagjaxb.XmlGeneralName
	/*
	 * (c) excluded_subtrees: a set of root names for each name type
	 * (e.g., X.500 distinguished names, email addresses, or IP
	 * addresses) defining a set of subtrees within which no subject
	 * name in subsequent certificates in the certification path may
	 * fall. This variable includes a set for each name type, and
	 * the initial value is initial-excluded-subtrees.
	 */
	var excludedSubtrees map[enumerations.GeneralNameType][]*diagjaxb.XmlGeneralName

	/*
	 * 6.1.3. Basic Certificate Processing
	 * The basic path processing actions to be performed for certificate i
	 * (for all i in [1..n]) are listed below.
	 */
	for i := len(certificateChain) - 1; i > -1; i-- {
		cert := certificateChain[i]
		/*
		 * (b) If certificate i is self-issued and it is not the final
		 * certificate in the path, skip this step for certificate i.
		 * Otherwise, verify that the subject name is within one of the
		 * permitted_subtrees for X.500 distinguished names, and verify
		 * that each of the alternative names in the subjectAltName
		 * extension (critical or non-critical) is within one of the
		 * permitted_subtrees for that name type.
		 */
		// perform validation only for the current certificate to support flexible validation policy
		if i == 0 {
			certDN := cert.CertificateDN()
			subAltNames := append([]*diagjaxb.XmlGeneralName{}, cert.SubjectAlternativeNames()...)
			/*
			 * Legacy implementations exist where an electronic mail address is
			 * embedded in the subject distinguished name in an attribute of type
			 * emailAddress (Section 4.1.2.6). When constraints are imposed on the
			 * alternative name, the rfc822Name constraint MUST be applied to the
			 * attribute of type emailAddress in the subject distinguished name.
			 */
			if !containsRFC822SubjectAlternativeName(subAltNames) {
				subAltNames = append(subAltNames, emailAddressDNIfPresent(certDN)...)
			}

			if permittedSubtrees != nil {
				dnGeneralNames := permittedSubtrees[enumerations.GeneralNameTypeDirectoryName]
				if dnGeneralNames != nil && !isWithinDNSubtrees(certDN, dnGeneralNames) {
					return false
				}
				for _, subAltName := range subAltNames {
					subtreesOfType := permittedSubtrees[generalNameTypeOf(subAltName)]
					if subtreesOfType != nil && !isWithinSubtrees(subAltName, subtreesOfType) {
						return false
					}
				}
			}
			/*
			 * (c) If certificate i is self-issued and it is not the final
			 * certificate in the path, skip this step for certificate i.
			 * Otherwise, verify that the subject name is not within any of
			 * the excluded_subtrees for X.500 distinguished names, and
			 * verify that each of the alternative names in the
			 * subjectAltName extension (critical or non-critical) is not
			 * within any of the excluded_subtrees for that name type.
			 */
			if excludedSubtrees != nil {
				dnGeneralNames := excludedSubtrees[enumerations.GeneralNameTypeDirectoryName]
				if dnGeneralNames != nil && isWithinDNSubtrees(certDN, dnGeneralNames) {
					return false
				}
				for _, subAltName := range subAltNames {
					subtreesOfType := excludedSubtrees[generalNameTypeOf(subAltName)]
					if subtreesOfType != nil && isWithinSubtrees(subAltName, subtreesOfType) {
						return false
					}
				}
			}
		}
		/*
		 * 6.1.4. Preparation for Certificate i+1
		 *
		 * (g) If a name constraints extension is included in the
		 * certificate, modify the permitted_subtrees and
		 * excluded_subtrees state variables as follows:
		 */
		certPermittedSubtrees := toGeneralNameMap(cert.PermittedSubtrees())
		certExcludedSubtrees := toGeneralNameMap(cert.ExcludedSubtrees())
		/*
		 * (1) If permittedSubtrees is present in the certificate, set
		 * the permitted_subtrees state variable to the intersection
		 * of its previous value and the value indicated in the
		 * extension field. If permittedSubtrees does not include a
		 * particular name type, the permitted_subtrees state
		 * variable is unchanged for that name type. For example,
		 * the intersection of example.com and foo.example.com is
		 * foo.example.com. And the intersection of example.com and
		 * example.net is the empty set.
		 */
		if certPermittedSubtrees != nil {
			if permittedSubtrees != nil {
				permittedSubtrees = intersectNameConstraints(permittedSubtrees, certPermittedSubtrees)
			} else {
				permittedSubtrees = certPermittedSubtrees
			}
		}

		/*
		 * (2) If excludedSubtrees is present in the certificate, set the
		 * excluded_subtrees state variable to the union of its
		 * previous value and the value indicated in the extension
		 * field. If excludedSubtrees does not include a particular
		 * name type, the excluded_subtrees state variable is
		 * unchanged for that name type. For example, the union of
		 * the name spaces example.com and foo.example.com is
		 * example.com. And the union of example.com and example.net
		 * is both name spaces.
		 */
		if certExcludedSubtrees != nil {
			if excludedSubtrees != nil {
				excludedSubtrees = unionNameConstraints(excludedSubtrees, certExcludedSubtrees)
			} else {
				excludedSubtrees = certExcludedSubtrees
			}
		}
	}

	return true
}

// containsRFC822SubjectAlternativeName ports the private
// containsRFC822SubjectAlternativeName(List).
func containsRFC822SubjectAlternativeName(subAltNames []*diagjaxb.XmlGeneralName) bool {
	for _, n := range subAltNames {
		if generalNameTypeOf(n) == enumerations.GeneralNameTypeRFC822Name {
			return true
		}
	}
	return false
}

// emailAddressDNIfPresent ports the private getEmailAddressDNIfPresent(String).
func emailAddressDNIfPresent(certDN string) []*diagjaxb.XmlGeneralName {
	dnMap := toDNMap(certDN)
	emailAddressValues := dnMap["1.2.840.113549.1.9.1"] // emailAddress

	var result []*diagjaxb.XmlGeneralName
	for emailAddress := range emailAddressValues {
		generalNameType := diagjaxb.GeneralNameTypeValue(enumerations.GeneralNameTypeRFC822Name)
		result = append(result, &diagjaxb.XmlGeneralName{
			XmlGeneralNameContent: diagjaxb.XmlGeneralNameContent{Value: emailAddress},
			XmlGeneralNameAttrs:   diagjaxb.XmlGeneralNameAttrs{Type: &generalNameType},
		})
	}
	return result
}

// toDNMap builds a DN map based on an RFC 2253 encoded string. Port of the
// private toDNMap(String).
//
// NOTE: see sun.security.x509.X500Name.parseRFC2253DN(String dnString).
func toDNMap(rfc2253EncodedString string) map[string]map[string]struct{} {
	if utils.IsStringEmpty(rfc2253EncodedString) {
		return map[string]map[string]struct{}{}
	}
	result := map[string]map[string]struct{}{}
	nextEnd := strings.IndexByte(rfc2253EncodedString, ',')
	searchOffset := 0
	dnOffset := 0
	for nextEnd >= 0 {
		if nextEnd > 0 && rfc2253EncodedString[nextEnd-1] != '\\' {
			nextStr := rfc2253EncodedString[dnOffset:nextEnd]
			if key, value, ok := getRDN(nextStr); ok {
				enrichDNMap(result, key, value)
			}
			dnOffset = nextEnd + 1
		}
		searchOffset = nextEnd + 1
		nextEnd = indexByteFrom(rfc2253EncodedString, ',', searchOffset)
	}
	// get last value entry
	substring := rfc2253EncodedString[dnOffset:]
	if key, value, ok := getRDN(substring); ok {
		enrichDNMap(result, key, value)
	}
	return result
}

// indexByteFrom mirrors Java's String#indexOf(int, int): the index of the
// first occurrence of c at or after fromIndex, or -1.
func indexByteFrom(s string, c byte, fromIndex int) int {
	if fromIndex >= len(s) {
		return -1
	}
	idx := strings.IndexByte(s[fromIndex:], c)
	if idx == -1 {
		return -1
	}
	return idx + fromIndex
}

// enrichDNMap ports the private enrichMap(Map, String, String).
func enrichDNMap(dnMap map[string]map[string]struct{}, rdnKey, rdnValue string) {
	values, ok := dnMap[rdnKey]
	if !ok {
		values = map[string]struct{}{}
		dnMap[rdnKey] = values
	}
	values[rdnValue] = struct{}{}
}

// getRDN ports the private getRDN(String).
func getRDN(str string) (key string, value string, ok bool) {
	nextEquals := strings.IndexByte(str, '=')
	if nextEquals >= 0 && len(str) >= nextEquals+1 {
		return str[:nextEquals], str[nextEquals+1:], true
	}
	return "", "", false
}

// isWithinSubtrees ports the private isWithinSubtrees(XmlGeneralName, Set).
func isWithinSubtrees(generalName *diagjaxb.XmlGeneralName, permittedSubtrees []*diagjaxb.XmlGeneralName) bool {
	for _, permittedSubtree := range permittedSubtrees {
		if isWithinSubtree(generalName, permittedSubtree) {
			return true
		}
	}
	return false
}

// isWithinDNSubtrees ports the private isWithinDNSubtrees(String, Set).
func isWithinDNSubtrees(certDN string, permittedSubtrees []*diagjaxb.XmlGeneralName) bool {
	if utils.IsStringEmpty(certDN) && utils.IsCollectionEmpty(permittedSubtrees) {
		return true
	}
	for _, permittedSubtree := range permittedSubtrees {
		if isWithinDNSubtree(certDN, permittedSubtree.Value) {
			return true
		}
	}
	return false
}

// toGeneralNameMap ports the private toXmlGeneralNameMap(Collection).
func toGeneralNameMap(generalSubtrees []*diagjaxb.XmlGeneralSubtree) map[enumerations.GeneralNameType][]*diagjaxb.XmlGeneralName {
	if utils.IsCollectionEmpty(generalSubtrees) {
		return nil
	}
	result := map[enumerations.GeneralNameType][]*diagjaxb.XmlGeneralName{}
	for _, generalSubtree := range generalSubtrees {
		t := generalSubtreeTypeOf(generalSubtree)
		result[t] = append(result[t], &diagjaxb.XmlGeneralName{
			XmlGeneralNameContent: generalSubtree.XmlGeneralNameContent,
			XmlGeneralNameAttrs:   generalSubtree.XmlGeneralNameAttrs,
		})
	}
	return result
}

// intersectNameConstraints ports the private intersectNew(Map, Map).
func intersectNameConstraints(originalConstraints, currentConstraints map[enumerations.GeneralNameType][]*diagjaxb.XmlGeneralName) map[enumerations.GeneralNameType][]*diagjaxb.XmlGeneralName {
	result := map[enumerations.GeneralNameType][]*diagjaxb.XmlGeneralName{}
	for _, t := range generalNameTypesOf(currentConstraints) {
		currentGeneralNames := currentConstraints[t]
		var intersection []*diagjaxb.XmlGeneralName

		originalSubtrees := originalConstraints[t]
		if utils.IsCollectionNotEmpty(originalSubtrees) {
			for _, currentGeneralName := range currentGeneralNames {
				for _, originalGeneralName := range originalSubtrees {
					if isWithinSubtree(originalGeneralName, currentGeneralName) {
						intersection = append(intersection, originalGeneralName)
					} else if isWithinSubtree(currentGeneralName, originalGeneralName) {
						intersection = append(intersection, currentGeneralName)
					}
				}
			}
		} else {
			intersection = append(intersection, currentGeneralNames...)
		}
		result[t] = intersection
	}

	for _, t := range generalNameTypesOf(originalConstraints) {
		if _, ok := result[t]; !ok {
			result[t] = append(result[t], originalConstraints[t]...)
		}
	}

	return result
}

// unionNameConstraints ports the private unionNew(Map, Map).
func unionNameConstraints(originalConstraints, currentConstraints map[enumerations.GeneralNameType][]*diagjaxb.XmlGeneralName) map[enumerations.GeneralNameType][]*diagjaxb.XmlGeneralName {
	result := map[enumerations.GeneralNameType][]*diagjaxb.XmlGeneralName{}
	for _, t := range generalNameTypesOf(currentConstraints) {
		currentGeneralNames := currentConstraints[t]
		var union []*diagjaxb.XmlGeneralName

		originalSubtrees := originalConstraints[t]
		if utils.IsCollectionNotEmpty(originalSubtrees) {
			for _, currentGeneralName := range currentGeneralNames {
				for _, originalGeneralName := range originalSubtrees {
					if isWithinSubtree(originalGeneralName, currentGeneralName) {
						union = append(union, currentGeneralName)
					} else if isWithinSubtree(currentGeneralName, originalGeneralName) {
						union = append(union, originalGeneralName)
					} else {
						union = append(union, currentGeneralName, originalGeneralName)
					}
				}
			}
		} else {
			union = append(union, currentGeneralNames...)
		}
		result[t] = union
	}

	for _, t := range generalNameTypesOf(originalConstraints) {
		if _, ok := result[t]; !ok {
			result[t] = append(result[t], originalConstraints[t]...)
		}
	}

	return result
}

// generalNameTypesOf returns the keys of a general-name-type map in a
// stable, deterministic order (avoiding Go's randomized map-iteration order
// in the output).
func generalNameTypesOf(m map[enumerations.GeneralNameType][]*diagjaxb.XmlGeneralName) []enumerations.GeneralNameType {
	order := []enumerations.GeneralNameType{
		enumerations.GeneralNameTypeOtherName, enumerations.GeneralNameTypeRFC822Name,
		enumerations.GeneralNameTypeDNSName, enumerations.GeneralNameTypeX400Address,
		enumerations.GeneralNameTypeDirectoryName, enumerations.GeneralNameTypeEDIPartyName,
		enumerations.GeneralNameTypeUniformResourceIdentifier, enumerations.GeneralNameTypeIPAddress,
		enumerations.GeneralNameTypeRegisteredID,
	}
	var types []enumerations.GeneralNameType
	for _, t := range order {
		if _, ok := m[t]; ok {
			types = append(types, t)
		}
	}
	return types
}

// isWithinSubtree ports the private isWithinSubtree(XmlGeneralName, XmlGeneralName).
func isWithinSubtree(generalName, subtreeGeneralName *diagjaxb.XmlGeneralName) bool {
	if utils.IsStringEmpty(generalName.Value) || utils.IsStringEmpty(subtreeGeneralName.Value) {
		return false
	}
	generalNameType := generalNameTypeOf(generalName)
	if generalNameType != enumerations.GeneralNameTypeIPAddress &&
		len(subtreeGeneralName.Value) > len(generalName.Value) {
		return false
	}

	switch generalNameType {
	case enumerations.GeneralNameTypeUniformResourceIdentifier:
		return isWithinURISubtree(generalName.Value, subtreeGeneralName.Value)
	case enumerations.GeneralNameTypeRFC822Name:
		return isWithinEmailSubtree(generalName.Value, subtreeGeneralName.Value)
	case enumerations.GeneralNameTypeDNSName:
		return isWithinDNSSubtree(generalName.Value, subtreeGeneralName.Value)
	case enumerations.GeneralNameTypeDirectoryName:
		return isWithinDNSubtree(generalName.Value, subtreeGeneralName.Value)
	case enumerations.GeneralNameTypeIPAddress:
		return isWithinIPAddressSubtree(generalName.Value, subtreeGeneralName.Value)
	case enumerations.GeneralNameTypeOtherName, enumerations.GeneralNameTypeX400Address,
		enumerations.GeneralNameTypeEDIPartyName, enumerations.GeneralNameTypeRegisteredID:
		// The NameConstraint of this type is not supported. Full comparison is executed.
		return isWithinOtherNameSubtree(generalName.Value, subtreeGeneralName.Value)
	}
	return false
}

// isWithinURISubtree ports the private isWithinURISubtree(String, String).
func isWithinURISubtree(value, subtree string) bool {
	domainNameSubtree := process.GetDomainName(subtree)
	domainNameValue := process.GetDomainName(value)
	return isWithinDomain(domainNameValue, domainNameSubtree)
}

// isWithinDomain ports the private isWithinDomain(String, String).
func isWithinDomain(value, domain string) bool {
	if strings.HasPrefix(domain, ".") {
		return strings.HasSuffix(strings.ToLower(value), strings.ToLower(domain))
	}
	return strings.EqualFold(domain, value)
}

// isWithinEmailSubtree ports the private isWithinEmailSubtree(String, String).
func isWithinEmailSubtree(value, subtree string) bool {
	if isEmail(subtree) {
		return strings.EqualFold(value, subtree)
	}
	if isEmail(value) {
		value = domainNameFromEmail(value)
	}
	return isWithinDomain(value, subtree)
}

// isEmail ports the private isEmail(String).
func isEmail(str string) bool {
	return strings.IndexByte(str, '@') != -1
}

// domainNameFromEmail ports the private getDomainNameFromEmail(String).
func domainNameFromEmail(email string) string {
	return email[strings.IndexByte(email, '@')+1:]
}

// isWithinDNSSubtree ports the private isWithinDNSSubtree(String, String).
func isWithinDNSSubtree(value, subtree string) bool {
	valueArray := strings.Split(value, ".")
	subTreeArray := strings.Split(subtree, ".")
	diff := len(valueArray) - len(subTreeArray)
	if diff == 0 {
		return equalStringSlices(subTreeArray, valueArray)
	} else if diff > 0 {
		for i := len(subTreeArray) - 1; i > -1; i-- {
			if subTreeArray[i] != valueArray[i+diff] {
				return false
			}
		}
		return true
	}
	return false
}

// equalStringSlices ports java.util.Arrays#equals(Object[], Object[]).
func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// isWithinDNSubtree ports the private isWithinDNSubtree(String, String).
func isWithinDNSubtree(value, subtree string) bool {
	dnMap := toDNMap(value)
	subtreeMap := toDNMap(subtree)
	for subtreeKey, subtreeValues := range subtreeMap {
		dnValues, ok := dnMap[subtreeKey]
		if !ok || !containsAllStrings(dnValues, subtreeValues) {
			return false
		}
	}
	return true
}

// containsAllStrings reports whether set contains every member of subset.
func containsAllStrings(set, subset map[string]struct{}) bool {
	for v := range subset {
		if _, ok := set[v]; !ok {
			return false
		}
	}
	return true
}

// isWithinIPAddressSubtree ports the private isWithinIPAddressSubtree(String, String).
func isWithinIPAddressSubtree(value, subtree string) bool {
	ipAddress := toByteArrayIPAddress(value)
	constraint := toByteArrayIPAddress(subtree)

	length := len(ipAddress)
	if length != len(constraint)/2 {
		return false
	}

	subnetMask := utils.Subarray(constraint, length, len(constraint))
	constraintSubnetAddress := make([]byte, length)
	ipSubnetAddress := make([]byte, length)

	// the resulting IP address by applying the subnet mask
	for i := 0; i < length; i++ {
		constraintSubnetAddress[i] = constraint[i] & subnetMask[i]
		ipSubnetAddress[i] = ipAddress[i] & subnetMask[i]
	}

	return equalByteSlices(constraintSubnetAddress, ipSubnetAddress)
}

// equalByteSlices ports java.util.Arrays#equals(byte[], byte[]).
func equalByteSlices(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// toByteArrayIPAddress ports the private toByteArrayIPAddress(String).
func toByteArrayIPAddress(ipAddress string) []byte {
	// consider internal hex-encoded values
	if strings.HasPrefix(ipAddress, "#") {
		stripped := strings.ReplaceAll(ipAddress, "#", "")
		if utils.IsHexEncoded(stripped) {
			decoded, err := utils.FromHex(stripped)
			if err == nil {
				return decoded
			}
		}
	}
	return []byte(ipAddress)
}

// isWithinOtherNameSubtree ports the private isWithinOtherNameSubtree(String, String).
func isWithinOtherNameSubtree(value, subtree string) bool {
	return subtree == value
}

// generalNameTypeOf reads the type attribute of an XmlGeneralName, mapping a
// null attribute (schema-absent) to the empty GeneralNameType.
func generalNameTypeOf(generalName *diagjaxb.XmlGeneralName) enumerations.GeneralNameType {
	if generalName.Type == nil {
		return ""
	}
	return generalName.Type.GeneralNameType()
}

// generalSubtreeTypeOf reads the type attribute of an XmlGeneralSubtree.
func generalSubtreeTypeOf(generalSubtree *diagjaxb.XmlGeneralSubtree) enumerations.GeneralNameType {
	if generalSubtree.Type == nil {
		return ""
	}
	return generalSubtree.Type.GeneralNameType()
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *CertificateNameConstraintsCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBXCVDCSBSINC
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *CertificateNameConstraintsCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBXCVDCSBSINCANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *CertificateNameConstraintsCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *CertificateNameConstraintsCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationCertificateChainGeneralFailure
}
