// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/JAdESHeaderParameterNames.java (DSS 6.5.RC1).
//
// Java's private constructor (utility class) has no Go counterpart; the constants below are
// package-level, named <TypeName><FieldName> per PORTING.md's enum/constant-namespace
// convention, values copied verbatim from ETSI TS 119 182-1 (never invented, abbreviated or
// "fixed").
//
// Two field name spellings collide with an already-landed sibling chunk that referenced this
// class before it was ported: TST_VD ("tstVD") is consumed as both
// JAdESHeaderParameterNamesTstVD (baseline_requirements_checker.go, certificate/crl/ocsp
// sources, level_baseline_lt.go) and JAdESHeaderParameterNamesTstVd (timestamp_source.go); SIG_PST
// ("sigPSt") is consumed as both JAdESHeaderParameterNamesSigPSt and
// JAdESHeaderParameterNamesSigPst (jades_signature.go / jades_level_baseline_lt.go /
// jades_extension_builder.go). Rather than edit those already-frozen files, both spellings are
// defined here, aliasing the same string value - see the file's notes to the porting lead.
package jades

// JAdESHeaderParameterNamesSigT is the claimed signing time. Port of SIG_T.
const JAdESHeaderParameterNamesSigT = "sigT"

// JAdESHeaderParameterNamesX5tO is the X509 certificate digest. Port of X5T_O.
const JAdESHeaderParameterNamesX5tO = "x5t#o"

// JAdESHeaderParameterNamesSigX5tS is the X509 certificate digests. Port of SIG_X5T_S.
const JAdESHeaderParameterNamesSigX5tS = "sigX5ts"

// JAdESHeaderParameterNamesDigAlg is the digest algorithm. Port of DIG_ALG.
const JAdESHeaderParameterNamesDigAlg = "digAlg"

// JAdESHeaderParameterNamesDigVal is the digest value. Port of DIG_VAL.
const JAdESHeaderParameterNamesDigVal = "digVal"

// JAdESHeaderParameterNamesSrCms is the signer commitments. Port of SR_CMS.
const JAdESHeaderParameterNamesSrCms = "srCms"

// JAdESHeaderParameterNamesCommId is the commitment Id. Port of COMM_ID.
const JAdESHeaderParameterNamesCommId = "commId"

// JAdESHeaderParameterNamesCommQuals is the commitment qualifiers. Port of COMM_QUALS.
const JAdESHeaderParameterNamesCommQuals = "commQuals"

// JAdESHeaderParameterNamesSigPl is the signature production place. Port of SIG_PL.
const JAdESHeaderParameterNamesSigPl = "sigPl"

// JAdESHeaderParameterNamesAddressCountry is the country address. Port of ADDRESS_COUNTRY.
const JAdESHeaderParameterNamesAddressCountry = "addressCountry"

// JAdESHeaderParameterNamesAddressLocality is the locality (city) address. Port of
// ADDRESS_LOCALITY.
const JAdESHeaderParameterNamesAddressLocality = "addressLocality"

// JAdESHeaderParameterNamesAddressRegion is the region (state and province) address. Port of
// ADDRESS_REGION.
const JAdESHeaderParameterNamesAddressRegion = "addressRegion"

// JAdESHeaderParameterNamesPostOfficeBoxNumber is the post office box number. Port of
// POST_OFFICE_BOX_NUMBER.
const JAdESHeaderParameterNamesPostOfficeBoxNumber = "postOfficeBoxNumber"

// JAdESHeaderParameterNamesPostalCode is the postal code. Port of POSTAL_CODE.
const JAdESHeaderParameterNamesPostalCode = "postalCode"

// JAdESHeaderParameterNamesStreetAddress is the street address. Port of STREET_ADDRESS.
const JAdESHeaderParameterNamesStreetAddress = "streetAddress"

// JAdESHeaderParameterNamesQArrays is used for signed assertions and claimed. Port of Q_ARRAYS.
const JAdESHeaderParameterNamesQArrays = "qArrays"

// JAdESHeaderParameterNamesMediaType is the media type. Port of MEDIA_TYPE.
const JAdESHeaderParameterNamesMediaType = "mediaType"

// JAdESHeaderParameterNamesQVals are the values used for Q Arrays. Port of Q_VALS.
const JAdESHeaderParameterNamesQVals = "qVals"

// JAdESHeaderParameterNamesSrAts is the signer attributes. Port of SR_ATS.
const JAdESHeaderParameterNamesSrAts = "srAts"

// JAdESHeaderParameterNamesClaimed is claimed. Port of CLAIMED.
const JAdESHeaderParameterNamesClaimed = "claimed"

// JAdESHeaderParameterNamesCertified is certified. Port of CERTIFIED.
const JAdESHeaderParameterNamesCertified = "certified"

// JAdESHeaderParameterNamesCertifiedAttrs are certified attributes. Port of CERTIFIED_ATTRS.
const JAdESHeaderParameterNamesCertifiedAttrs = "certifiedAttrs"

// JAdESHeaderParameterNamesX509AttrCert is the X509 attribute certificate. Port of
// X509_ATTR_CERT.
const JAdESHeaderParameterNamesX509AttrCert = "x509AttrCert"

// JAdESHeaderParameterNamesOtherAttrCert is the other attribute certificate. Port of
// OTHER_ATTR_CERT.
const JAdESHeaderParameterNamesOtherAttrCert = "otherAttrCert"

// JAdESHeaderParameterNamesSignedAssertions are the signed assertions. Port of
// SIGNED_ASSERTIONS.
const JAdESHeaderParameterNamesSignedAssertions = "signedAssertions"

// JAdESHeaderParameterNamesAdoTst is the signed data time-stamp. Port of ADO_TST.
const JAdESHeaderParameterNamesAdoTst = "adoTst"

// JAdESHeaderParameterNamesSigPid is the signature policy identifier. Port of SIG_PID.
const JAdESHeaderParameterNamesSigPid = "sigPId"

// JAdESHeaderParameterNamesId is the Id. Port of ID.
const JAdESHeaderParameterNamesId = "id"

// JAdESHeaderParameterNamesHashAV is the hash algo and value. Port of HASH_AV.
const JAdESHeaderParameterNamesHashAV = "hashAV"

// JAdESHeaderParameterNamesDigPSp indicates the hash policy is aligned to a specification. Port
// of DIG_PSP.
const JAdESHeaderParameterNamesDigPSp = "digPSp"

// JAdESHeaderParameterNamesSigPQuals are the signature policy qualifiers. Port of SIG_P_QUALS.
const JAdESHeaderParameterNamesSigPQuals = "sigPQuals"

// JAdESHeaderParameterNamesSigPQual is the signature policy qualifier. Port of SIG_P_QUAL.
const JAdESHeaderParameterNamesSigPQual = "sigPQual"

// JAdESHeaderParameterNamesSpURI is the signature policy URL qualifier. Port of SP_URI.
const JAdESHeaderParameterNamesSpURI = "spURI"

// JAdESHeaderParameterNamesSpUserNotice is the signature policy User Notice qualifier. Port of
// SP_USER_NOTICE.
const JAdESHeaderParameterNamesSpUserNotice = "spUserNotice"

// JAdESHeaderParameterNamesNoticeRef is the notice references. Port of NOTICE_REF.
const JAdESHeaderParameterNamesNoticeRef = "noticeRef"

// JAdESHeaderParameterNamesOrgantization is the organization. Port of ORGANTIZATION - the
// upstream field name misspells "Organization" verbatim; kept as-is per PORTING.md ("never
// invent, abbreviate, or fix an ... enum name").
const JAdESHeaderParameterNamesOrgantization = "organization"

// JAdESHeaderParameterNamesNoticeNumbers are the notice numbers. Port of NOTICE_NUMBERS.
const JAdESHeaderParameterNamesNoticeNumbers = "noticeNumbers"

// JAdESHeaderParameterNamesExplText is the explicit text. Port of EXPL_TEXT.
const JAdESHeaderParameterNamesExplText = "explText"

// JAdESHeaderParameterNamesSpDspec is the signature policy Document Specification qualifier.
// Port of SP_DSPEC.
const JAdESHeaderParameterNamesSpDspec = "spDSpec"

// JAdESHeaderParameterNamesSigD is the signed data. Port of SIG_D.
const JAdESHeaderParameterNamesSigD = "sigD"

// JAdESHeaderParameterNamesMId is the signed data referencing mechanism URI. Port of M_ID.
const JAdESHeaderParameterNamesMId = "mId"

// JAdESHeaderParameterNamesPars are the signed data references. Port of PARS.
const JAdESHeaderParameterNamesPars = "pars"

// JAdESHeaderParameterNamesHashM is the signed data digest algorithm identifier. Port of HASH_M.
const JAdESHeaderParameterNamesHashM = "hashM"

// JAdESHeaderParameterNamesHashV is the array of signed data digest algorithm values (hashes).
// Port of HASH_V.
const JAdESHeaderParameterNamesHashV = "hashV"

// JAdESHeaderParameterNamesCtys is the array of signed data types (see 'cty'). Port of CTYS.
const JAdESHeaderParameterNamesCtys = "ctys"

// JAdESHeaderParameterNamesDesc is the description. Port of DESC.
const JAdESHeaderParameterNamesDesc = "desc"

// JAdESHeaderParameterNamesDocRefs are the document references. Port of DOC_REFS.
const JAdESHeaderParameterNamesDocRefs = "docRefs"

// JAdESHeaderParameterNamesCanonAlg is the canonicalization algorithm. Port of CANON_ALG.
const JAdESHeaderParameterNamesCanonAlg = "canonAlg"

// JAdESHeaderParameterNamesTstTokens is the timestamp tokens array. Port of TST_TOKENS.
const JAdESHeaderParameterNamesTstTokens = "tstTokens"

// JAdESHeaderParameterNamesEncoding is the encoding (e.g. DER, ...). Port of ENCODING.
const JAdESHeaderParameterNamesEncoding = "encoding"

// JAdESHeaderParameterNamesVal is the value (i.e. Timestamp base64 value). Port of VAL.
const JAdESHeaderParameterNamesVal = "val"

// JAdESHeaderParameterNamesEtsiU is the ETSI Unsigned properties. Port of ETSI_U.
const JAdESHeaderParameterNamesEtsiU = "etsiU"

// JAdESHeaderParameterNamesSigTst is the signature timestamp. Port of SIG_TST.
const JAdESHeaderParameterNamesSigTst = "sigTst"

// JAdESHeaderParameterNamesTstContainer is the container for timestamps. Port of TST_CONTAINER.
const JAdESHeaderParameterNamesTstContainer = "tstContainer"

// JAdESHeaderParameterNamesXVals are the certificate values. Port of X_VALS.
const JAdESHeaderParameterNamesXVals = "xVals"

// JAdESHeaderParameterNamesAxVals are the certificate values of Attribute Authorities. Port of
// AX_VALS.
const JAdESHeaderParameterNamesAxVals = "axVals"

// JAdESHeaderParameterNamesRVals are the revocation values. Port of R_VALS.
const JAdESHeaderParameterNamesRVals = "rVals"

// JAdESHeaderParameterNamesArVals are the revocation values of Attribute Authorities. Port of
// AR_VALS.
const JAdESHeaderParameterNamesArVals = "arVals"

// JAdESHeaderParameterNamesCrlVals are the CRL values. Port of CRL_VALS.
const JAdESHeaderParameterNamesCrlVals = "crlVals"

// JAdESHeaderParameterNamesOcspVals are the OCSP values. Port of OCSP_VALS.
const JAdESHeaderParameterNamesOcspVals = "ocspVals"

// JAdESHeaderParameterNamesOtherVals are the other values. Port of OTHER_VALS.
const JAdESHeaderParameterNamesOtherVals = "otherVals"

// JAdESHeaderParameterNamesX509Cert is the X.509 certificate. Port of X509_CERT.
const JAdESHeaderParameterNamesX509Cert = "x509Cert"

// JAdESHeaderParameterNamesOtherCert is the other certificate. Port of OTHER_CERT.
const JAdESHeaderParameterNamesOtherCert = "otherCert"

// JAdESHeaderParameterNamesXRefs are the certificate references. Port of X_REFS.
const JAdESHeaderParameterNamesXRefs = "xRefs"

// JAdESHeaderParameterNamesAxRefs are the references to certificates of Attribute Authorities.
// Port of AX_REFS.
const JAdESHeaderParameterNamesAxRefs = "axRefs"

// JAdESHeaderParameterNamesCertId is the certificate identifier (reference). Port of CERT_ID.
const JAdESHeaderParameterNamesCertId = "certId"

// JAdESHeaderParameterNamesRRefs are the revocation references. Port of R_REFS.
const JAdESHeaderParameterNamesRRefs = "rRefs"

// JAdESHeaderParameterNamesArRefs are the references to revocations of Attribute Authorities.
// Port of AR_REFS.
const JAdESHeaderParameterNamesArRefs = "arRefs"

// JAdESHeaderParameterNamesCrlRefs are the CRL references. Port of CRL_REFS.
const JAdESHeaderParameterNamesCrlRefs = "crlRefs"

// JAdESHeaderParameterNamesOcspRefs are the OCSP references. Port of OCSP_REFS.
const JAdESHeaderParameterNamesOcspRefs = "ocspRefs"

// JAdESHeaderParameterNamesOcspId is the OCSP Id. Port of OCSP_ID.
const JAdESHeaderParameterNamesOcspId = "ocspId"

// JAdESHeaderParameterNamesProducedAt is the OCSP Response production time. Port of
// PRODUCED_AT.
const JAdESHeaderParameterNamesProducedAt = "producedAt"

// JAdESHeaderParameterNamesResponderId is the OCSP Responder id. Port of RESPONDER_ID.
const JAdESHeaderParameterNamesResponderId = "responderId"

// JAdESHeaderParameterNamesByName is the OCSP Responder id - by name. Port of BY_NAME.
const JAdESHeaderParameterNamesByName = "byName"

// JAdESHeaderParameterNamesByKey is the OCSP Responder id - by key. Port of BY_KEY.
const JAdESHeaderParameterNamesByKey = "byKey"

// JAdESHeaderParameterNamesCrlId is the CRL Id. Port of CRL_ID.
const JAdESHeaderParameterNamesCrlId = "crlId"

// JAdESHeaderParameterNamesIssuer is the CRL issuer. Port of ISSUER.
const JAdESHeaderParameterNamesIssuer = "issuer"

// JAdESHeaderParameterNamesIssueTime is the CRL issue time. Port of ISSUE_TIME.
const JAdESHeaderParameterNamesIssueTime = "issueTime"

// JAdESHeaderParameterNamesNumber is the CRL Number. Port of NUMBER.
const JAdESHeaderParameterNamesNumber = "number"

// JAdESHeaderParameterNamesPkiOb is the PKI Object. Port of PKI_OB.
const JAdESHeaderParameterNamesPkiOb = "pkiOb"

// JAdESHeaderParameterNamesTstVD is the Timestamp Validation Data. Port of TST_VD.
const JAdESHeaderParameterNamesTstVD = "tstVD"

// JAdESHeaderParameterNamesTstVd is an alias of JAdESHeaderParameterNamesTstVD; see the file
// header note on the two spellings this ROOT chunk must provide for already-landed call sites.
const JAdESHeaderParameterNamesTstVd = "tstVD"

// JAdESHeaderParameterNamesArcTst is the Archive TimeStamp. Port of ARC_TST.
const JAdESHeaderParameterNamesArcTst = "arcTst"

// JAdESHeaderParameterNamesSigRTst is the Signature and References Timestamp. Port of
// SIG_R_TST.
const JAdESHeaderParameterNamesSigRTst = "sigRTst"

// JAdESHeaderParameterNamesRfsTst is the References Timestamp. Port of RFS_TST.
const JAdESHeaderParameterNamesRfsTst = "rfsTst"

// JAdESHeaderParameterNamesCSig is the Counter Signature. Port of C_SIG.
const JAdESHeaderParameterNamesCSig = "cSig"

// JAdESHeaderParameterNamesSigPSt is the Signature Policy Store. Port of SIG_PST.
const JAdESHeaderParameterNamesSigPSt = "sigPSt"

// JAdESHeaderParameterNamesSigPst is an alias of JAdESHeaderParameterNamesSigPSt; see the file
// header note on the two spellings this ROOT chunk must provide for already-landed call sites.
const JAdESHeaderParameterNamesSigPst = "sigPSt"

// JAdESHeaderParameterNamesSigPolDoc is the signature policy document. Port of SIG_POL_DOC.
const JAdESHeaderParameterNamesSigPolDoc = "sigPolDoc"

// JAdESHeaderParameterNamesSigPolLocalURI is the signature policy document local URI. Port of
// SIG_POL_LOCAL_URI.
const JAdESHeaderParameterNamesSigPolLocalURI = "sigPolLocalURI"

// JAdESHeaderParameterNamesAnyValData is any Validation Data. Port of ANY_VAL_DATA.
const JAdESHeaderParameterNamesAnyValData = "anyValData"

// JAdESHeaderParameterNamesExp is RFC 7519 "JSON Web Token (JWT)", 4.1.4. "exp" (Expiration
// Time) Claim.
//
// Deprecated: since DSS 6.5. Please use JWTClaimNamesExp instead. Port of EXP.
const JAdESHeaderParameterNamesExp = "exp"

// JAdESHeaderParameterNamesIat is RFC 7519 "JSON Web Token (JWT)", 4.1.6. "iat" (Issued At)
// Claim.
//
// Deprecated: since DSS 6.5. Please use JWTClaimNamesIat instead. Port of IAT.
const JAdESHeaderParameterNamesIat = "iat"
