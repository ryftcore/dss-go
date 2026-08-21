// Ported from
// dss-jades/src/main/java/eu/europa/esig/dss/jades/validation/JAdESRevocationRefExtractionUtils.java
// (DSS 6.5.RC1).
//
// Deviation: upstream reads 'ocspId'/'crlId'/'responderId' with a raw `(Map<?, ?>) ... .get(...)`
// cast that throws ClassCastException (caught by the surrounding try/catch(Exception), aborting
// the whole createOCSPRef/createCRLRef call and returning null) when the value is present but not
// a JSON object. This file instead reads them through DSSJsonUtilsGetAsMap, which degrades to an
// empty map on a type mismatch rather than aborting - the same tolerant-cast convention every
// other DSSJsonUtils-based reader in this port already uses (see JAdESCertificateSource,
// JAdESCRLSource, JAdESOCSPSource). A malformed non-map value under 'ocspId'/'crlId'/'responderId'
// is accordingly treated as absent rather than failing the whole reference, which can only ever
// produce a reference DSS itself would otherwise have discarded entirely.
package jades

import (
	"math/big"
	"strconv"
	"time"

	"github.com/ryftcore/dss-go/dss/internal/jose"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/utils"
)

// JAdESRevocationRefExtractionUtilsCreateOCSPRef extracts an OCSPRef from the 'ocspRefs' header.
// Port of the public static createOCSPRef(Map).
func JAdESRevocationRefExtractionUtilsCreateOCSPRef(ocpRef *jose.Object) *spi.OCSPRef {
	var responderId *spi.ResponderId
	var producedAt time.Time

	ocspId := DSSJsonUtilsGetAsMap(ocpRef, JAdESHeaderParameterNamesOcspId)
	if ocspId.Size() != 0 {
		producedAt = spi.DSSUtilsParseRFCDate(DSSJsonUtilsGetAsString(ocspId, JAdESHeaderParameterNamesProducedAt))

		var err error
		responderId, err = jadesRevocationRefExtractionUtilsResponderId(ocspId)
		if err != nil {
			// Upstream logs "Unable to extract OCSPRef. Reason : {}" and discards the whole
			// result (the outer try/catch(Exception) wraps the entire method body).
			return nil
		}
	}

	digest, ok := DSSJsonUtilsDigest(ocpRef)
	if !ok {
		// Upstream logs "Missing digest information in OCSPRef".
		return nil
	}
	return spi.NewOCSPRef(digest, producedAt, responderId)
}

// jadesRevocationRefExtractionUtilsResponderId ports the private getResponderId(Map).
func jadesRevocationRefExtractionUtilsResponderId(ocspId *jose.Object) (*spi.ResponderId, error) {
	responderIdMap := DSSJsonUtilsGetAsMap(ocspId, JAdESHeaderParameterNamesResponderId)
	if responderIdMap.Size() == 0 {
		return nil, nil
	}

	var subjectX500Principal *model.X500Principal
	var ski []byte

	byNameB64 := DSSJsonUtilsGetAsString(responderIdMap, JAdESHeaderParameterNamesByName)
	if utils.IsStringNotEmpty(byNameB64) && utils.IsBase64Encoded(byNameB64) {
		principal, err := spi.DSSASN1UtilsToX500Principal(utils.FromBase64(byNameB64))
		if err != nil {
			return nil, err
		}
		subjectX500Principal = principal
	}

	byKeyB64 := DSSJsonUtilsGetAsString(responderIdMap, JAdESHeaderParameterNamesByKey)
	if utils.IsStringNotEmpty(byKeyB64) && utils.IsBase64Encoded(byKeyB64) {
		ski = utils.FromBase64(byKeyB64)
	}

	if subjectX500Principal != nil || utils.IsArrayNotEmpty(ski) {
		return spi.NewResponderId(subjectX500Principal, ski), nil
	}
	return nil, nil
}

// JAdESRevocationRefExtractionUtilsCreateCRLRef extracts a CRLRef from the 'crlRefs' header.
// Port of the public static createCRLRef(Map).
func JAdESRevocationRefExtractionUtilsCreateCRLRef(crlRefMap *jose.Object) *spi.CRLRef {
	var crlIssuer *model.X500Principal
	var crlIssuedTime time.Time
	var crlNumber *big.Int

	crlId := DSSJsonUtilsGetAsMap(crlRefMap, JAdESHeaderParameterNamesCrlId)
	if crlId.Size() != 0 {
		issuerB64 := DSSJsonUtilsGetAsString(crlId, JAdESHeaderParameterNamesIssuer)
		if utils.IsStringNotEmpty(issuerB64) && utils.IsBase64Encoded(issuerB64) {
			principal, err := spi.DSSASN1UtilsToX500Principal(utils.FromBase64(issuerB64))
			if err != nil {
				// Upstream logs "Unable to extract a CRLRef. Reason : {}" and discards the
				// whole result (the outer try/catch(Exception) wraps the entire method body).
				return nil
			}
			crlIssuer = principal
		}

		issueTimeStr := DSSJsonUtilsGetAsString(crlId, JAdESHeaderParameterNamesIssueTime)
		if utils.IsStringNotEmpty(issueTimeStr) {
			crlIssuedTime = spi.DSSUtilsParseRFCDate(issueTimeStr)
		}

		crlNumberString := DSSJsonUtilsGetAsString(crlId, JAdESHeaderParameterNamesNumber)
		if utils.IsStringNotEmpty(crlNumberString) {
			crlNumberValue, err := strconv.ParseInt(crlNumberString, 10, 64)
			if err != nil {
				// Java's Long.parseLong(String) throws NumberFormatException, caught by the
				// outer try/catch(Exception): the whole result is discarded.
				return nil
			}
			crlNumber = big.NewInt(crlNumberValue)
		}
	}

	digest, ok := DSSJsonUtilsDigest(crlRefMap)
	if !ok {
		// Upstream logs "Missing digest information in CRLRef".
		return nil
	}
	return spi.NewCRLRefWithNumber(digest, crlIssuer, crlIssuedTime, crlNumber)
}
