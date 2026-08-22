// Java-oracle tests for the synchronization strategies (accept_all_strategy.go and
// expiration_and_signature_check_strategy.go).
//
// These two decide whether a trusted list's certificates reach the TrustedListsCertificateSource
// at all, so a flipped comparison here silently trusts an expired or invalidly-signed list, or
// silently drops a good one - the exact "gate" this batch's Criticals live in.
// ExpirationAndSignatureCheckStrategy has no dedicated JUnit suite upstream; its expected
// behaviour below is read off the Java source (isSyncSupported(TLInfo, boolean, boolean)) as a
// full truth table, and the two rows upstream's own
// TLValidationJobTest#brokenSigWithSyncStrategyTest pins end to end - a TOTAL_FAILED signature
// with acceptInvalidTrustedList=false is NOT synchronized, and the same list is once its
// signature validates - are the "invalid/not accepted" and "valid" rows of that table.
//
// Test vectors are ported, not the JUnit code, per PORTING.md.
package tsl

import (
	"testing"
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
	tslmodel "github.com/ryftcore/dss-go/dss/model/tsl"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/validation/job"
)

// syncStrategyParsingInfo builds a real TLParsingCacheDTO whose result either is absent
// (resultExist=false) or carries the given next-update date.
func syncStrategyParsingInfo(resultExist bool, nextUpdateDate time.Time) tslmodel.TLParsingInfoRecord {
	if !resultExist {
		return NewTLParsingCacheDTOBuilder(job.NewCachedEntry[job.ParsingResult]()).Build()
	}
	result := NewTLParsingResult()
	result.SetNextUpdateDate(nextUpdateDate)
	return NewTLParsingCacheDTOBuilder(job.NewCachedEntryWithResult[job.ParsingResult](result)).Build()
}

// syncStrategyValidationInfo builds a real ValidationCacheDTO whose result either is absent
// (resultExist=false) or carries the given Indication.
func syncStrategyValidationInfo(resultExist bool, indication enumerations.Indication) *job.ValidationCacheDTO {
	if !resultExist {
		return job.NewValidationCacheDTOBuilder(job.NewCachedEntry[job.ValidationResult]()).Build()
	}
	subIndication := enumerations.SubIndication("")
	if indication == enumerations.IndicationTotalFailed {
		subIndication = enumerations.SubIndicationHashFailure
	}
	// A non-nil certificate source, because ValidationCacheDTOBuilder#build() reads
	// getPotentialSigners(), which dereferences it (as Java's does).
	certificateSource := spi.NewCommonCertificateSource()
	result := NewTLValidationResult(indication, subIndication, time.Unix(0, 0).UTC(), nil, &certificateSource)
	return job.NewValidationCacheDTOBuilder(job.NewCachedEntryWithResult[job.ValidationResult](result)).Build()
}

// TestExpirationAndSignatureCheckStrategy_Oracle walks the full truth table of
// isSyncSupported(TLInfo, boolean, boolean).
func TestExpirationAndSignatureCheckStrategy_Oracle(t *testing.T) {
	future := time.Now().Add(24 * time.Hour)
	past := time.Now().Add(-24 * time.Hour)

	cases := []struct {
		name string
		// parsing side
		parsingResultExist bool
		nextUpdateDate     time.Time
		// validation side
		validationResultExist bool
		indication            enumerations.Indication
		// flags
		acceptExpired, acceptInvalid bool
		want                         bool
	}{
		// Neither flag set: a still-valid, TOTAL_PASSED list synchronizes.
		{"validAndFresh", true, future, true, enumerations.IndicationTotalPassed, false, false, true},
		// Expiry gate. A null next-update date counts as expired upstream
		// (`nextUpdateDate == null || currentDate.after(nextUpdateDate)`); the Go stand-in for
		// that null is the zero time.Time.
		{"expired", true, past, true, enumerations.IndicationTotalPassed, false, false, false},
		{"expiredAccepted", true, past, true, enumerations.IndicationTotalPassed, true, false, true},
		{"noNextUpdate", true, time.Time{}, true, enumerations.IndicationTotalPassed, false, false, false},
		{"noNextUpdateAccepted", true, time.Time{}, true, enumerations.IndicationTotalPassed, true, false, true},
		// No parsing result at all: the expiry gate does not apply (Java guards on
		// isResultExist()), so only the signature gate can refuse.
		{"noParsingResult", false, time.Time{}, true, enumerations.IndicationTotalPassed, false, false, true},
		{"noParsingResultInvalid", false, time.Time{}, true, enumerations.IndicationTotalFailed, false, false, false},
		// Signature gate. isValid() is "indication == TOTAL_PASSED", so both TOTAL_FAILED and
		// INDETERMINATE refuse.
		{"brokenSignature", true, future, true, enumerations.IndicationTotalFailed, false, false, false},
		{"brokenSignatureAccepted", true, future, true, enumerations.IndicationTotalFailed, false, true, true},
		{"indeterminateSignature", true, future, true, enumerations.IndicationIndeterminate, false, false, false},
		{"indeterminateSignatureAccepted", true, future, true, enumerations.IndicationIndeterminate, false, true, true},
		// No validation result: the signature gate does not apply either.
		{"noValidationResult", true, future, false, "", false, false, true},
		// Both gates refuse; accepting only one is not enough.
		{"expiredAndBroken", true, past, true, enumerations.IndicationTotalFailed, false, false, false},
		{"expiredAcceptedStillBroken", true, past, true, enumerations.IndicationTotalFailed, true, false, false},
		{"brokenAcceptedStillExpired", true, past, true, enumerations.IndicationTotalFailed, false, true, false},
		{"bothAccepted", true, past, true, enumerations.IndicationTotalFailed, true, true, true},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			tlInfo := tslmodel.NewTLInfo(nil,
				syncStrategyParsingInfo(c.parsingResultExist, c.nextUpdateDate),
				syncStrategyValidationInfo(c.validationResultExist, c.indication),
				"http://test/tl.xml")

			// The TL flags drive canBeSynchronized(TLInfo)...
			tlStrategy := NewExpirationAndSignatureCheckStrategy()
			tlStrategy.SetAcceptExpiredTrustedList(c.acceptExpired)
			tlStrategy.SetAcceptInvalidTrustedList(c.acceptInvalid)
			if got := tlStrategy.CanBeSynchronizedDocument(tlInfo); got != c.want {
				t.Errorf("CanBeSynchronizedDocument = %v, want %v", got, c.want)
			}
			// ... and must NOT drive canBeSynchronized(LOTLInfo), which reads the separate
			// listOfTrustedLists pair (this is the mix-up the two flag pairs invite).
			if got := tlStrategy.CanBeSynchronizedDocumentList(
				lotlInfoFor(tlInfo)); got != expectedWithLOTLDefaults(c.parsingResultExist, c.nextUpdateDate, c.validationResultExist, c.indication) {
				t.Errorf("CanBeSynchronizedDocumentList with only the TL flags set = %v, want %v",
					got, expectedWithLOTLDefaults(c.parsingResultExist, c.nextUpdateDate, c.validationResultExist, c.indication))
			}

			// The LOTL flags drive canBeSynchronized(LOTLInfo) the same way.
			lotlStrategy := NewExpirationAndSignatureCheckStrategy()
			lotlStrategy.SetAcceptExpiredListOfTrustedLists(c.acceptExpired)
			lotlStrategy.SetAcceptInvalidListOfTrustedLists(c.acceptInvalid)
			if got := lotlStrategy.CanBeSynchronizedDocumentList(lotlInfoFor(tlInfo)); got != c.want {
				t.Errorf("CanBeSynchronizedDocumentList = %v, want %v", got, c.want)
			}
		})
	}
}

// lotlInfoFor wraps a TLInfo as the LOTLInfo canBeSynchronized(LOTLInfo) takes.
func lotlInfoFor(tlInfo *tslmodel.TLInfo) *tslmodel.LOTLInfo {
	lotlInfo := tslmodel.NewLOTLInfo(tlInfo.DownloadCacheInfo(), mustTLParsingInfo(tlInfo),
		tlInfo.ValidationCacheInfo(), tlInfo.Url())
	return &lotlInfo
}

func mustTLParsingInfo(tlInfo *tslmodel.TLInfo) tslmodel.TLParsingInfoRecord {
	parsingCacheInfo, _ := tlInfo.TLParsingCacheInfo()
	return parsingCacheInfo
}

// expectedWithLOTLDefaults is the verdict canBeSynchronized(LOTLInfo) must reach when only the
// TRUSTED LIST flags were set, i.e. with both LOTL flags at their false default.
func expectedWithLOTLDefaults(parsingResultExist bool, nextUpdateDate time.Time,
	validationResultExist bool, indication enumerations.Indication) bool {
	if parsingResultExist && (nextUpdateDate.IsZero() || time.Now().After(nextUpdateDate)) {
		return false
	}
	if validationResultExist {
		return indication == enumerations.IndicationTotalPassed
	}
	return true
}

// TestAcceptAllStrategy_Oracle pins the deprecated tsl.AcceptAllStrategy and its
// validation/job replacement: both accept everything, including a list that expired and whose
// signature failed - which is exactly the difference from
// ExpirationAndSignatureCheckStrategy's defaults.
func TestAcceptAllStrategy_Oracle(t *testing.T) {
	tlInfo := tslmodel.NewTLInfo(nil,
		syncStrategyParsingInfo(true, time.Now().Add(-24*time.Hour)),
		syncStrategyValidationInfo(true, enumerations.IndicationTotalFailed),
		"http://test/tl.xml")
	lotlInfo := lotlInfoFor(tlInfo)

	deprecated := NewAcceptAllStrategy()
	if !deprecated.CanBeSynchronizedDocument(tlInfo) {
		t.Error("tsl.AcceptAllStrategy must synchronize an expired, invalidly-signed TL")
	}
	if !deprecated.CanBeSynchronizedDocumentList(lotlInfo) {
		t.Error("tsl.AcceptAllStrategy must synchronize an expired, invalidly-signed LOTL")
	}

	replacement := job.NewAcceptAllStrategy[*tslmodel.TLInfo, *tslmodel.LOTLInfo]()
	if !replacement.CanBeSynchronizedDocument(tlInfo) {
		t.Error("job.AcceptAllStrategy must synchronize an expired, invalidly-signed TL")
	}
	if !replacement.CanBeSynchronizedDocumentList(lotlInfo) {
		t.Error("job.AcceptAllStrategy must synchronize an expired, invalidly-signed LOTL")
	}

	// The default strategy refuses the same list, so the two are not interchangeable.
	strict := NewExpirationAndSignatureCheckStrategy()
	if strict.CanBeSynchronizedDocument(tlInfo) {
		t.Error("ExpirationAndSignatureCheckStrategy's defaults must refuse an expired, invalidly-signed TL")
	}
}
