// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/checks/PolicyTreeNode.java (DSS 6.5.RC1).
//
// Java represents a node's children as a HashSet<PolicyTreeNode> without a
// custom equals/hashCode override, i.e. identity-based membership; a Go slice
// of pointers preserves that same identity-based semantics without needing a
// hashable key.
package xcv

import (
	"slices"

	"github.com/ryftcore/dss-go/dss/utils"
)

// policyTreeAnyPolicyOID represents an anyPolicy OID.
const policyTreeAnyPolicyOID = "2.5.29.32.0"

// PolicyTreeNode represents a valid_policy_tree node (leaf) as per RFC 5280.
type PolicyTreeNode struct {
	// validPolicy represents a valid_policy element of a policy tree node.
	validPolicy string

	// qualifierSet represents a qualifier_set element of a policy tree node.
	qualifierSet []string

	// expectedPolicySet represents a expected_policy_set element of a policy
	// tree node.
	expectedPolicySet []string

	// children represents policy node's children.
	children []*PolicyTreeNode
}

// NewPolicyTreeNode is the default constructor. Port of
// PolicyTreeNode(String, String).
func NewPolicyTreeNode(policyOID string, policyQualifier string) *PolicyTreeNode {
	var qualifierSet []string
	if policyQualifier != "" {
		qualifierSet = []string{policyQualifier}
	}
	return &PolicyTreeNode{
		validPolicy:       policyOID,
		qualifierSet:      qualifierSet,
		expectedPolicySet: []string{policyOID},
	}
}

// InitPolicyTree initializes the first node of the valid policy tree
// (containing anyPolicy as the first element). Port of the static
// initTree().
func InitPolicyTree() *PolicyTreeNode {
	return NewPolicyTreeNode(policyTreeAnyPolicyOID, "")
}

// IsAnyPolicy returns if the current policy node represents anyPolicy. Port
// of isAnyPolicy().
func (n *PolicyTreeNode) IsAnyPolicy() bool {
	return policyTreeAnyPolicyOID == n.validPolicy
}

// AddChildNodeIfMatch adds policyNode to the node's children, when
// applicable. Port of addChildNodeIfMatch(PolicyTreeNode).
func (n *PolicyTreeNode) AddChildNodeIfMatch(policyNode *PolicyTreeNode) bool {
	/*
	 * (i) For each node of depth i-1 in the valid_policy_tree
	 * where P-OID is in the expected_policy_set, create a
	 * child node as follows: set the valid_policy to P-OID,
	 * set the qualifier_set to P-Q, and set the
	 * expected_policy_set to {P-OID}.
	 */
	if slices.Contains(n.expectedPolicySet, policyNode.validPolicy) {
		n.children = append(n.children, policyNode)
		return true
	}
	/*
	 * (ii) If there was no match in step (i) and the
	 * valid_policy_tree includes a node of depth i-1 with
	 * the valid_policy anyPolicy, generate a child node with
	 * the following values: set the valid_policy to P-OID,
	 * set the qualifier_set to P-Q, and set the
	 * expected_policy_set to {P-OID}.
	 */
	if n.IsAnyPolicy() {
		// AnyPolicy child is created within the caller class
		n.children = append(n.children, policyNode)
		return true
	}
	return false
}

// CreateAnyPolicyChildren creates any policy children corresponding to the
// current policy node. Port of createAnyPolicyChildren().
func (n *PolicyTreeNode) CreateAnyPolicyChildren() []*PolicyTreeNode {
	var anyPolicyChildren []*PolicyTreeNode
	for _, expectedPolicy := range n.expectedPolicySet {
		child := NewPolicyTreeNode(expectedPolicy, policyTreeAnyPolicyOID)
		anyPolicyChildren = append(anyPolicyChildren, child)
	}
	n.children = append(n.children, anyPolicyChildren...)
	return anyPolicyChildren
}

// DeleteNodesAtLevelWithoutChildren removes nodes at the given depthLevel if
// not having children nodes. Port of deleteNodesAtLevelWithoutChildren(int).
func (n *PolicyTreeNode) DeleteNodesAtLevelWithoutChildren(depthLevel int) *PolicyTreeNode {
	if depthLevel > 0 {
		var kept []*PolicyTreeNode
		for _, child := range n.children {
			if child.DeleteNodesAtLevelWithoutChildren(depthLevel-1) != nil {
				kept = append(kept, child)
			}
		}
		n.children = kept
	}
	if utils.IsCollectionEmpty(n.children) {
		return nil
	}
	return n
}
