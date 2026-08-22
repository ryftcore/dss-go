// Ported from dss-model/.../model/policy/crypto/CryptographicSuiteCatalogue.java (DSS 6.5.RC1).
package policy

import (
	"strings"

	"github.com/ryftcore/dss-go/dss/enumerations"
)

// CryptographicSuiteCatalogue provides an abstract implementation of an
// ETSI TS 119 322 cryptographic suite catalogue, providing extraction of
// cryptographic suites according to the defined usage. This class uses a
// "smart mapping", decreasing the amount of created objects for various
// validation scopes, when applicable.
//
// Java's protected abstract methods buildMetadata()/buildAlgorithmList()
// (implemented by concrete catalogue subclasses outside this manifest)
// become constructor-supplied function fields, since Go embedding does
// not support virtual dispatch back into an outer type. A subclass in a
// future chunk constructs its catalogue via
// NewCryptographicSuiteCatalogue(buildMetadata, buildAlgorithmList) and
// embeds *CryptographicSuiteCatalogue for the exported accessors below.
type CryptographicSuiteCatalogue struct {
	// cryptographicSuiteMap is a cached map of created cryptographic
	// suites, keyed by a canonical identifier of the algorithm list
	// (Java keys by the List<CryptographicSuiteAlgorithm> itself, which
	// uses order-sensitive content equality/hashCode).
	cryptographicSuiteMap map[string]CryptographicSuite

	// metadata is the metadata extracted from the cryptographic suite
	// document.
	metadata *CryptographicSuiteMetadata

	// algorithmList is the list of cryptographic algorithms and their
	// corresponding validation rules extracted from the cryptographic
	// suite document.
	algorithmList []*CryptographicSuiteAlgorithm

	// buildMetadataFunc ports the abstract buildMetadata() method.
	buildMetadataFunc func() *CryptographicSuiteMetadata

	// buildAlgorithmListFunc ports the abstract buildAlgorithmList()
	// method.
	buildAlgorithmListFunc func() []*CryptographicSuiteAlgorithm
}

// NewCryptographicSuiteCatalogue is the default constructor. buildMetadata
// and buildAlgorithmList implement the abstract Java methods of the same
// name.
func NewCryptographicSuiteCatalogue(buildMetadata func() *CryptographicSuiteMetadata, buildAlgorithmList func() []*CryptographicSuiteAlgorithm) *CryptographicSuiteCatalogue {
	return &CryptographicSuiteCatalogue{
		cryptographicSuiteMap:  make(map[string]CryptographicSuite),
		buildMetadataFunc:      buildMetadata,
		buildAlgorithmListFunc: buildAlgorithmList,
	}
}

// Metadata gets the metadata.
func (c *CryptographicSuiteCatalogue) Metadata() *CryptographicSuiteMetadata {
	if c.metadata == nil {
		c.metadata = c.buildMetadataFunc()
	}
	return c.metadata
}

// AlgorithmList gets the algorithm rules list.
func (c *CryptographicSuiteCatalogue) AlgorithmList() []*CryptographicSuiteAlgorithm {
	if c.algorithmList == nil {
		c.algorithmList = c.buildAlgorithmListFunc()
	}
	return c.algorithmList
}

// CryptographicSuite gets the global CryptographicSuite.
func (c *CryptographicSuiteCatalogue) CryptographicSuite() CryptographicSuite {
	algorithms := c.filterByAlgorithmUsage(c.AlgorithmList(), []enumerations.CryptographicSuiteAlgorithmUsage{
		enumerations.CryptographicSuiteAlgorithmUsage_SIGN_DATA, enumerations.CryptographicSuiteAlgorithmUsage_VALIDATE_DATA,
	})
	return c.cryptographicSuiteFor(c.Metadata(), algorithms)
}

// SignatureCryptographicSuite gets the CryptographicSuite for validation
// of a signature.
func (c *CryptographicSuiteCatalogue) SignatureCryptographicSuite() CryptographicSuite {
	// same as global constraints
	return c.CryptographicSuite()
}

// SignatureCertificatesCryptographicSuite gets the CryptographicSuite for
// validation of signature certificates.
func (c *CryptographicSuiteCatalogue) SignatureCertificatesCryptographicSuite() CryptographicSuite {
	algorithms := c.filterByAlgorithmUsage(c.AlgorithmList(), []enumerations.CryptographicSuiteAlgorithmUsage{
		enumerations.CryptographicSuiteAlgorithmUsage_SIGN_DATA, enumerations.CryptographicSuiteAlgorithmUsage_VALIDATE_DATA,
		enumerations.CryptographicSuiteAlgorithmUsage_SIGN_CERTIFICATES, enumerations.CryptographicSuiteAlgorithmUsage_VALIDATE_CERTIFICATES,
	})
	return c.cryptographicSuiteFor(c.Metadata(), algorithms)
}

// CounterSignatureCryptographicSuite gets the CryptographicSuite for
// validation of a counter signature.
func (c *CryptographicSuiteCatalogue) CounterSignatureCryptographicSuite() CryptographicSuite {
	// same as global constraints
	return c.CryptographicSuite()
}

// CounterSignatureCertificatesCryptographicSuite gets the
// CryptographicSuite for validation of counter signature certificates.
func (c *CryptographicSuiteCatalogue) CounterSignatureCertificatesCryptographicSuite() CryptographicSuite {
	// same as signature constraints
	return c.SignatureCertificatesCryptographicSuite()
}

// KeyBindingSignatureCryptographicSuite gets the CryptographicSuite for
// validation of a key binding signature.
func (c *CryptographicSuiteCatalogue) KeyBindingSignatureCryptographicSuite() CryptographicSuite {
	// same as global constraints
	return c.CryptographicSuite()
}

// KeyBindingSignatureCertificatesCryptographicSuite gets the
// CryptographicSuite for validation of key binding signature
// certificates.
func (c *CryptographicSuiteCatalogue) KeyBindingSignatureCertificatesCryptographicSuite() CryptographicSuite {
	// same as signature constraints
	return c.SignatureCertificatesCryptographicSuite()
}

// RevocationCryptographicSuite gets the CryptographicSuite for validation
// of a revocation data.
func (c *CryptographicSuiteCatalogue) RevocationCryptographicSuite() CryptographicSuite {
	algorithms := c.filterByAlgorithmUsage(c.AlgorithmList(), []enumerations.CryptographicSuiteAlgorithmUsage{
		enumerations.CryptographicSuiteAlgorithmUsage_SIGN_DATA, enumerations.CryptographicSuiteAlgorithmUsage_VALIDATE_DATA,
		enumerations.CryptographicSuiteAlgorithmUsage_SIGN_OCSP, enumerations.CryptographicSuiteAlgorithmUsage_VALIDATE_OCSP,
	})
	return c.cryptographicSuiteFor(c.Metadata(), algorithms)
}

// RevocationCertificatesCryptographicSuite gets the CryptographicSuite
// for validation of revocation data certificates.
func (c *CryptographicSuiteCatalogue) RevocationCertificatesCryptographicSuite() CryptographicSuite {
	algorithms := c.filterByAlgorithmUsage(c.AlgorithmList(), []enumerations.CryptographicSuiteAlgorithmUsage{
		enumerations.CryptographicSuiteAlgorithmUsage_SIGN_DATA, enumerations.CryptographicSuiteAlgorithmUsage_VALIDATE_DATA,
		enumerations.CryptographicSuiteAlgorithmUsage_SIGN_OCSP, enumerations.CryptographicSuiteAlgorithmUsage_VALIDATE_OCSP,
		enumerations.CryptographicSuiteAlgorithmUsage_SIGN_CERTIFICATES, enumerations.CryptographicSuiteAlgorithmUsage_VALIDATE_CERTIFICATES,
	})
	return c.cryptographicSuiteFor(c.Metadata(), algorithms)
}

// TimestampCryptographicSuite gets the CryptographicSuite for validation
// of a timestamp.
func (c *CryptographicSuiteCatalogue) TimestampCryptographicSuite() CryptographicSuite {
	algorithms := c.filterByAlgorithmUsage(c.AlgorithmList(), []enumerations.CryptographicSuiteAlgorithmUsage{
		enumerations.CryptographicSuiteAlgorithmUsage_SIGN_DATA, enumerations.CryptographicSuiteAlgorithmUsage_VALIDATE_DATA,
		enumerations.CryptographicSuiteAlgorithmUsage_SIGN_TIMESTAMPS, enumerations.CryptographicSuiteAlgorithmUsage_VALIDATE_TIMESTAMPS,
	})
	return c.cryptographicSuiteFor(c.Metadata(), algorithms)
}

// TimestampCertificatesCryptographicSuite gets the CryptographicSuite for
// validation of timestamp data certificates.
func (c *CryptographicSuiteCatalogue) TimestampCertificatesCryptographicSuite() CryptographicSuite {
	algorithms := c.filterByAlgorithmUsage(c.AlgorithmList(), []enumerations.CryptographicSuiteAlgorithmUsage{
		enumerations.CryptographicSuiteAlgorithmUsage_SIGN_DATA, enumerations.CryptographicSuiteAlgorithmUsage_VALIDATE_DATA,
		enumerations.CryptographicSuiteAlgorithmUsage_SIGN_TIMESTAMPS, enumerations.CryptographicSuiteAlgorithmUsage_VALIDATE_TIMESTAMPS,
		enumerations.CryptographicSuiteAlgorithmUsage_SIGN_CERTIFICATES, enumerations.CryptographicSuiteAlgorithmUsage_VALIDATE_CERTIFICATES,
	})
	return c.cryptographicSuiteFor(c.Metadata(), algorithms)
}

// EvidenceRecordCryptographicSuite gets the CryptographicSuite for
// validation of an evidence record.
func (c *CryptographicSuiteCatalogue) EvidenceRecordCryptographicSuite() CryptographicSuite {
	// no separate handling
	return c.CryptographicSuite()
}

// EAACryptographicSuite gets the CryptographicSuite for validation of an
// EAA.
func (c *CryptographicSuiteCatalogue) EAACryptographicSuite() CryptographicSuite {
	// no separate handling
	return c.CryptographicSuite()
}

// EAARevocationCryptographicSuite gets the CryptographicSuite for
// validation of an EAA revocation.
func (c *CryptographicSuiteCatalogue) EAARevocationCryptographicSuite() CryptographicSuite {
	// no separate handling
	return c.CryptographicSuite()
}

// filterByAlgorithmUsage ports the private
// filterByAlgorithmUsage(List, List) helper.
func (c *CryptographicSuiteCatalogue) filterByAlgorithmUsage(algorithmList []*CryptographicSuiteAlgorithm, algorithmUsages []enumerations.CryptographicSuiteAlgorithmUsage) []*CryptographicSuiteAlgorithm {
	var result []*CryptographicSuiteAlgorithm
	for _, orig := range algorithmList {
		algorithm := CryptographicSuiteAlgorithmCopy(orig)
		if len(algorithm.EvaluationList()) == 0 {
			result = append(result, algorithm)
			continue
		}

		var evaluationList []*CryptographicSuiteEvaluation
		for _, evaluation := range algorithm.EvaluationList() {
			usage := evaluation.AlgorithmUsage()
			if len(usage) == 0 || cryptographicSuiteCatalogueUsageOverlaps(usage, algorithmUsages) {
				evaluationList = append(evaluationList, evaluation)
			}
		}
		algorithm.SetEvaluationList(evaluationList)
		if len(evaluationList) != 0 {
			result = append(result, algorithm)
		}
	}
	return result
}

// cryptographicSuiteCatalogueUsageOverlaps reports whether any element of
// usage is contained in usages, mirroring
// algorithmUsage.stream().anyMatch(algorithmUsages::contains).
func cryptographicSuiteCatalogueUsageOverlaps(usage, usages []enumerations.CryptographicSuiteAlgorithmUsage) bool {
	for _, u := range usage {
		for _, candidate := range usages {
			if u == candidate {
				return true
			}
		}
	}
	return false
}

// cryptographicSuiteFor builds a cryptographic suite for the given
// content. If the content is already present within the
// cryptographicSuiteMap, the method returns the existing entry. Ports
// the protected getCryptographicSuite(CryptographicSuiteMetadata, List)
// method.
func (c *CryptographicSuiteCatalogue) cryptographicSuiteFor(metadata *CryptographicSuiteMetadata, algorithmList []*CryptographicSuiteAlgorithm) CryptographicSuite {
	key := cryptographicSuiteCatalogueAlgorithmListKey(algorithmList)
	cryptographicSuite, ok := c.cryptographicSuiteMap[key]
	if !ok {
		cryptographicSuite = NewCryptographicSuite19322(metadata, algorithmList)
		c.cryptographicSuiteMap[key] = cryptographicSuite
	}
	return cryptographicSuite
}

// cryptographicSuiteCatalogueAlgorithmListKey builds a canonical string
// identifier for an algorithm list, standing in for Java's order-sensitive
// content-based List#equals/hashCode when used as a map key.
func cryptographicSuiteCatalogueAlgorithmListKey(algorithmList []*CryptographicSuiteAlgorithm) string {
	parts := make([]string, len(algorithmList))
	for i, a := range algorithmList {
		parts[i] = a.String()
	}
	return strings.Join(parts, "\x1f")
}
