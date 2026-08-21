// Ported from dss-token/src/main/java/eu/europa/esig/dss/token/predicate/ExtendedKeyUsageKeyEntryPredicate.java (DSS 6.5.RC1).
package token

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
)

// NewExtendedKeyUsageKeyEntryPredicate creates a predicate filtering private keys based on the
// certificate ExtendedKeyUsage attribute value, accepting the given ExtendedKeyUsages. Port of
// the ExtendedKeyUsage... constructor.
//
// Panics with the Java message if extendedKeyUsages is nil (Objects.requireNonNull); a nil
// element within it is silently dropped, matching the Java Stream#filter(Objects::nonNull).
func NewExtendedKeyUsageKeyEntryPredicate(extendedKeyUsages ...enumerations.ExtendedKeyUsage) DSSKeyEntryPredicate {
	if extendedKeyUsages == nil {
		panic("ExtendedKeyUsage cannot be null!")
	}
	oids := make(map[string]struct{}, len(extendedKeyUsages))
	for _, eku := range extendedKeyUsages {
		if eku != "" {
			oids[eku.OID()] = struct{}{}
		}
	}
	return newExtendedKeyUsageKeyEntryPredicate(oids)
}

// NewExtendedKeyUsageKeyEntryPredicateForOIDs creates a predicate filtering private keys based on
// the certificate ExtendedKeyUsage attribute value, accepting the given ExtendedKeyUsage OIDs.
// Port of the String... constructor.
//
// Panics with the Java message if extendedKeyUsageOIDs is nil (Objects.requireNonNull).
func NewExtendedKeyUsageKeyEntryPredicateForOIDs(extendedKeyUsageOIDs ...string) DSSKeyEntryPredicate {
	if extendedKeyUsageOIDs == nil {
		panic("ExtendedKeyUsage OIDs cannot be null!")
	}
	oids := make(map[string]struct{}, len(extendedKeyUsageOIDs))
	for _, oid := range extendedKeyUsageOIDs {
		oids[oid] = struct{}{}
	}
	return newExtendedKeyUsageKeyEntryPredicate(oids)
}

func newExtendedKeyUsageKeyEntryPredicate(extendedKeyUsageOIDs map[string]struct{}) DSSKeyEntryPredicate {
	return func(dssPrivateKeyEntry DSSPrivateKeyEntry) bool {
		certificate := dssPrivateKeyEntry.Certificate()
		if certificate == nil {
			return false
		}
		for _, extendedKeyUsage := range extendedKeyUsageKeyEntryPredicateGetExtendedKeyUsages(certificate) {
			if _, ok := extendedKeyUsageOIDs[extendedKeyUsage]; ok {
				return true
			}
		}
		return false
	}
}

// extendedKeyUsageKeyEntryPredicateGetExtendedKeyUsages ports the private getExtendedKeyUsages(CertificateToken).
//
// DEVIATION: Java re-raises java.security.cert.CertificateParsingException as an unchecked
// DSSException, which propagates through Predicate#test unannounced; spi.CertificateExtensionsUtilsExtendedKeyUsage
// already parses the ExtendedKeyUsage extension (see its own header on the BouncyCastle-vs-JDK
// discrepancy this port must reproduce) and reports the equivalent parse failure by returning nil,
// so this panics in that case to preserve the "exception surfaces from a predicate call" behavior.
func extendedKeyUsageKeyEntryPredicateGetExtendedKeyUsages(certificateToken *model.CertificateToken) []string {
	extendedKeyUsage := spi.CertificateExtensionsUtilsExtendedKeyUsage(certificateToken)
	if extendedKeyUsage == nil {
		panic(fmt.Sprintf("DSSException : Unable to extract ExtendedKeyUsage from a certificate token. "+
			"Reason : %s", "unable to parse the ExtendedKeyUsage extension"))
	}
	return extendedKeyUsage.Oids()
}
