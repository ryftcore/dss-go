// Ported from dss-model/.../model/policy/crypto/CryptographicSuiteMetadata.java (DSS 6.5.RC1).
package policy

import (
	"fmt"
	"time"
)

// CryptographicSuiteMetadata contains metadata about the ETSI TS 119 322
// cryptographic suite.
//
// java.io.Serializable is dropped silently (no Go counterpart).
type CryptographicSuiteMetadata struct {
	// policyName is the value of the /dssc:PolicyName/dssc:Name element.
	policyName string

	// policyOID is the value of the
	// /dssc:PolicyName/dssc:ObjectIdentifier element.
	policyOID string

	// policyURI is the value of the /dssc:PolicyName/dssc:URI element.
	policyURI string

	// publisherName is the value of the /dssc:Publisher/dssc:Name
	// element.
	publisherName string

	// publisherAddress is the value of the
	// /dssc:Publisher/dssc:Address element.
	publisherAddress string

	// publisherURI is the value of the /dssc:Publisher/dssc:URI element.
	publisherURI string

	// policyIssueDate is the value of the /dssc:PolicyIssueDate element.
	policyIssueDate *time.Time

	// nextUpdate is the value of the /dssc:NextUpdate element.
	nextUpdate *time.Time

	// usage is the value of the /dssc:Usage element.
	usage string

	// version is the value of the /dssc:version element.
	version string

	// lang is the value of the /dssc:lang element.
	lang string

	// id is the value of the /dssc:id element.
	id string
}

// NewCryptographicSuiteMetadata is the default constructor.
func NewCryptographicSuiteMetadata() *CryptographicSuiteMetadata {
	return &CryptographicSuiteMetadata{}
}

// PolicyName gets the policy name.
func (c *CryptographicSuiteMetadata) PolicyName() string { return c.policyName }

// SetPolicyName sets the value of the /dssc:PolicyName/dssc:Name element.
func (c *CryptographicSuiteMetadata) SetPolicyName(policyName string) { c.policyName = policyName }

// PolicyOID gets the policy OID.
func (c *CryptographicSuiteMetadata) PolicyOID() string { return c.policyOID }

// SetPolicyOID sets the value of the
// /dssc:PolicyName/dssc:ObjectIdentifier element.
func (c *CryptographicSuiteMetadata) SetPolicyOID(policyOID string) { c.policyOID = policyOID }

// PolicyURI gets the policy URI.
func (c *CryptographicSuiteMetadata) PolicyURI() string { return c.policyURI }

// SetPolicyURI sets the value of the /dssc:PolicyName/dssc:URI element.
func (c *CryptographicSuiteMetadata) SetPolicyURI(policyURI string) { c.policyURI = policyURI }

// PublisherName gets the policy publisher's name.
func (c *CryptographicSuiteMetadata) PublisherName() string { return c.publisherName }

// SetPublisherName sets the value of the /dssc:Publisher/dssc:Name
// element.
func (c *CryptographicSuiteMetadata) SetPublisherName(publisherName string) {
	c.publisherName = publisherName
}

// PublisherAddress gets the policy publisher's address.
func (c *CryptographicSuiteMetadata) PublisherAddress() string { return c.publisherAddress }

// SetPublisherAddress sets the value of the
// /dssc:Publisher/dssc:Address element.
func (c *CryptographicSuiteMetadata) SetPublisherAddress(publisherAddress string) {
	c.publisherAddress = publisherAddress
}

// PublisherURI gets the policy publisher's website URI.
func (c *CryptographicSuiteMetadata) PublisherURI() string { return c.publisherURI }

// SetPublisherURI sets the value of the /dssc:Publisher/dssc:URI element.
func (c *CryptographicSuiteMetadata) SetPublisherURI(publisherURI string) {
	c.publisherURI = publisherURI
}

// PolicyIssueDate gets the policy's issue date.
func (c *CryptographicSuiteMetadata) PolicyIssueDate() *time.Time { return c.policyIssueDate }

// SetPolicyIssueDate sets the value of the /dssc:PolicyIssueDate element.
func (c *CryptographicSuiteMetadata) SetPolicyIssueDate(policyIssueDate *time.Time) {
	c.policyIssueDate = policyIssueDate
}

// NextUpdate gets the policy's nest update date.
func (c *CryptographicSuiteMetadata) NextUpdate() *time.Time { return c.nextUpdate }

// SetNextUpdate sets the value of the /dssc:NextUpdate element.
func (c *CryptographicSuiteMetadata) SetNextUpdate(nextUpdate *time.Time) { c.nextUpdate = nextUpdate }

// Usage gets the policy usage.
func (c *CryptographicSuiteMetadata) Usage() string { return c.usage }

// SetUsage sets the value of the /dssc:Usage element.
func (c *CryptographicSuiteMetadata) SetUsage(usage string) { c.usage = usage }

// Version gets the policy version.
func (c *CryptographicSuiteMetadata) Version() string { return c.version }

// SetVersion sets the value of the /dssc:version element.
func (c *CryptographicSuiteMetadata) SetVersion(version string) { c.version = version }

// Lang gets the policy language two-character ISO identifier.
func (c *CryptographicSuiteMetadata) Lang() string { return c.lang }

// SetLang sets the value of the /dssc:lang element.
func (c *CryptographicSuiteMetadata) SetLang(lang string) { c.lang = lang }

// Id gets the policy identifier.
func (c *CryptographicSuiteMetadata) Id() string { return c.id }

// SetId sets the value of the /dssc:id element.
func (c *CryptographicSuiteMetadata) SetId(id string) { c.id = id }

// Equals ports CryptographicSuiteMetadata#equals.
func (c *CryptographicSuiteMetadata) Equals(other *CryptographicSuiteMetadata) bool {
	if c == other {
		return true
	}
	if other == nil {
		return false
	}
	return c.policyName == other.policyName &&
		c.policyOID == other.policyOID &&
		c.policyURI == other.policyURI &&
		c.publisherName == other.publisherName &&
		c.publisherAddress == other.publisherAddress &&
		c.publisherURI == other.publisherURI &&
		cryptographicSuiteMetadataTimeEqual(c.policyIssueDate, other.policyIssueDate) &&
		cryptographicSuiteMetadataTimeEqual(c.nextUpdate, other.nextUpdate) &&
		c.usage == other.usage &&
		c.version == other.version &&
		c.lang == other.lang &&
		c.id == other.id
}

// cryptographicSuiteMetadataTimeEqual compares two possibly-nil
// *time.Time values.
func cryptographicSuiteMetadataTimeEqual(a, b *time.Time) bool {
	if a == nil || b == nil {
		return a == b
	}
	return a.Equal(*b)
}

// String ports CryptographicSuiteMetadata#toString.
func (c *CryptographicSuiteMetadata) String() string {
	return fmt.Sprintf("CryptographicSuiteMetadata [policyName='%s', policyOID='%s', policyURI='%s', "+
		"publisherName='%s', publisherAddress='%s', publisherURI='%s', policyIssueDate=%v, nextUpdate=%v, "+
		"usage='%s', version='%s', lang='%s', id='%s']",
		c.policyName, c.policyOID, c.policyURI, c.publisherName, c.publisherAddress, c.publisherURI,
		c.policyIssueDate, c.nextUpdate, c.usage, c.version, c.lang, c.id)
}
