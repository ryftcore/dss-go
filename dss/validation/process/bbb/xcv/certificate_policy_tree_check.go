// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/checks/CertificatePolicyTreeCheck.java (DSS 6.5.RC1).
package xcv

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// CertificatePolicyTreeCheck verifies if the certificate has a valid policy
// tree according to its certification path in regard to RFC 5280.
type CertificatePolicyTreeCheck struct {
	*process.ChainItemBase[*jaxb.XmlSubXCV]

	// certificate is the certificate to check.
	certificate *diagnostic.CertificateWrapper
}

// NewCertificatePolicyTreeCheck is the default constructor. Port of
// CertificatePolicyTreeCheck(I18nProvider, XmlSubXCV, CertificateWrapper, LevelRule).
func NewCertificatePolicyTreeCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSubXCV],
	certificate *diagnostic.CertificateWrapper, constraint policy.LevelRule) *CertificatePolicyTreeCheck {
	c := &CertificatePolicyTreeCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		certificate:   certificate,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *CertificatePolicyTreeCheck) Process() bool {
	/*
	 * 6.1.2. Initialization
	 */
	certificateChain := append([]*diagnostic.CertificateWrapper{c.certificate}, c.certificate.CertificateChain()...)
	/*
	 * (a) valid_policy_tree: A tree of certificate policies with their
	 * optional qualifiers; each of the leaves of the tree
	 * represents a valid policy at this stage in the certification
	 * path validation.
	 * ...
	 * The initial value of the valid_policy_tree is a single node with
	 * valid_policy anyPolicy, an empty qualifier_set, and an
	 * expected_policy_set with the single value anyPolicy.
	 * This node is considered to be at depth zero.
	 */
	validPolicyTree := InitPolicyTree()
	/*
	 * (d) explicit_policy: an integer that indicates if a non-NULL
	 * valid_policy_tree is required.
	 * ...
	 * If initial-explicit-policy is set, then the
	 * initial value is 0, otherwise the initial value is n+1.
	 */
	explicitPolicy := len(certificateChain) + 1
	/*
	 * (e) inhibit_anyPolicy: an integer that indicates whether the
	 * anyPolicy policy identifier is considered a match.
	 * ...
	 * If initial-any-policy-inhibit is set, then the initial value is 0,
	 * otherwise the initial value is n+1.
	 */
	inhibitAnyPolicy := len(certificateChain) + 1

	// internal variable to facilitate processing
	previousLevelNodes := []*PolicyTreeNode{validPolicyTree}

	/*
	 * 6.1.3. Basic Certificate Processing
	 * The basic path processing actions to be performed for certificate i
	 * (for all i in [1..n]) are listed below.
	 */
	for i := len(certificateChain) - 1; i > -1; i-- {
		cert := certificateChain[i]
		certRequireExplicitPolicy := cert.RequireExplicitPolicy()
		certInhibitAnyPolicy := cert.InhibitAnyPolicy()
		/*
		 * (d) If the certificate policies extension is present in the
		 * certificate and the valid_policy_tree is not NULL, process
		 * the policy information by performing the following steps in
		 * order:
		 */
		var currentLevelNodes []*PolicyTreeNode
		certificatePolicies := cert.CertificatePolicies()
		if utils.IsCollectionNotEmpty(certificatePolicies) && utils.IsCollectionNotEmpty(previousLevelNodes) {
			for _, certificatePolicy := range certificatePolicies {
				var cpsURL string
				if certificatePolicy.CpsUrl != nil {
					cpsURL = *certificatePolicy.CpsUrl
				}
				policyNode := NewPolicyTreeNode(certificatePolicy.Value, cpsURL)
				/*
				 * (1) For each policy P not equal to anyPolicy in the
				 * certificate policies extension, let P-OID denote the OID
				 * for policy P and P-Q denote the qualifier set for policy P.
				 * Perform the following steps in order:
				 */
				if !policyNode.IsAnyPolicy() {
					for _, node := range previousLevelNodes {
						if node.AddChildNodeIfMatch(policyNode) {
							currentLevelNodes = appendPolicyTreeNodeOnce(currentLevelNodes, policyNode)
						}
					}
					/*
					 * (2) If the certificate policies extension includes the policy
					 * anyPolicy with the qualifier set AP-Q and either (a)
					 * inhibit_anyPolicy is greater than 0 or (b) i<n and the
					 * certificate is self-issued, then:
					 */
				} else if inhibitAnyPolicy > 0 || (i != 0 && cert.IsSelfSigned()) {
					for _, node := range previousLevelNodes {
						children := node.CreateAnyPolicyChildren()
						for _, child := range children {
							currentLevelNodes = appendPolicyTreeNodeOnce(currentLevelNodes, child)
						}
					}
				}
			}
			/*
			 * (3) If there is a node in the valid_policy_tree of depth i-1
			 * or less without any child nodes, delete that node. Repeat
			 * this step until there are no nodes of depth i-1 or less
			 * without children.
			 */
			if validPolicyTree != nil {
				validPolicyTree = validPolicyTree.DeleteNodesAtLevelWithoutChildren(len(certificateChain) - 1 - i)
			}
			/*
			 * (e) If the certificate policies extension is not present, set the
			 * valid_policy_tree to NULL.
			 */
		} else if utils.IsCollectionEmpty(certificatePolicies) {
			validPolicyTree = nil
		}
		/*
		 * (f) Verify that either explicit_policy is greater than 0 or the
		 * valid_policy_tree is not equal to NULL;
		 */
		// skip this step in order to support flexible validation policy
		/*
		 * If i is not equal to n, continue by performing the preparatory steps
		 * listed in Section 6.1.4. If i is equal to n, perform the wrap-up
		 * steps listed in Section 6.1.5.
		 */
		// descending order
		if i != 0 {
			/*
			 * 6.1.4. Preparation for Certificate i+1
			 */
			previousLevelNodes = currentLevelNodes
			/*
			 * (h) If certificate i is not self-issued:
			 */
			if !cert.IsSelfSigned() {
				/*
				 * (1) If explicit_policy is not 0, decrement explicit_policy by 1.
				 * ...
				 */
				if explicitPolicy != 0 {
					explicitPolicy--
				}
				/*
				 * (3) If inhibit_anyPolicy is not 0, decrement inhibit_anyPolicy by 1.
				 */
				if inhibitAnyPolicy != 0 {
					inhibitAnyPolicy--
				}
			}
			/*
			 * (i) (1) If requireExplicitPolicy is present and is less than
			 * explicit_policy, set explicit_policy to the value of
			 * requireExplicitPolicy.
			 */
			if certRequireExplicitPolicy != -1 && certRequireExplicitPolicy < explicitPolicy {
				explicitPolicy = certRequireExplicitPolicy
			}
			/*
			 * (j) If the inhibitAnyPolicy extension is included in the
			 * certificate and is less than inhibit_anyPolicy, set
			 * inhibit_anyPolicy to the value of inhibitAnyPolicy.
			 */
			if certInhibitAnyPolicy != -1 && certInhibitAnyPolicy < inhibitAnyPolicy {
				inhibitAnyPolicy = certInhibitAnyPolicy
			}
			/*
			 * 6.1.5. Wrap-Up Procedure
			 * To complete the processing of the target certificate, perform the
			 * following steps for certificate n:
			 */
		} else {
			/*
			 * (a) If explicit_policy is not 0, decrement explicit_policy by 1.
			 */
			if explicitPolicy != 0 {
				explicitPolicy--
			}
			/*
			 * (b) If a policy constraints extension is included in the
			 * certificate and requireExplicitPolicy is present and has a
			 * value of 0, set the explicit_policy state variable to 0.
			 */
			if certRequireExplicitPolicy == 0 {
				explicitPolicy = 0
			}
			/*
			 * If either (1) the value of explicit_policy variable is greater than
			 * zero or (2) the valid_policy_tree is not NULL, then path processing
			 * has succeeded.
			 */
			if explicitPolicy == 0 && validPolicyTree == nil {
				return false
			}
		}
	}

	return true
}

// appendPolicyTreeNodeOnce mirrors adding to a Java HashSet<PolicyTreeNode>
// with default identity equality: pointer identity de-duplication.
func appendPolicyTreeNodeOnce(nodes []*PolicyTreeNode, node *PolicyTreeNode) []*PolicyTreeNode {
	for _, n := range nodes {
		if n == node {
			return nodes
		}
	}
	return append(nodes, node)
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *CertificatePolicyTreeCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_ICPTV
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *CertificatePolicyTreeCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_ICPTV_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *CertificatePolicyTreeCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *CertificatePolicyTreeCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_CERTIFICATE_CHAIN_GENERAL_FAILURE
}
