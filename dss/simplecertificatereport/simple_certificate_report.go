// Ported from dss-simple-certificate-report-jaxb/src/main/java/eu/europa/esig/dss/simplecertificatereport/SimpleCertificateReport.java
// (DSS 6.5.RC1).
//
// Several accessors below call getFirstCertificate() (Certificate is
// minOccurs="0" in the schema, so it can be nil) without a nil check before
// dereferencing it, exactly as the Java source does (XmlChainItem
// cert = getFirstCertificate(); return cert.getQwacProfile(); - no null
// guard). That is a genuine upstream NullPointerException risk on a report
// with no Certificate element; per PORTING.md's "values verbatim" this port
// preserves it rather than papering over it - Go's own nil-pointer field
// access panics in exactly the same shape.
package simplecertificatereport

import (
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/simplecertificatereport/jaxb"
)

// SimpleCertificateReport is a SimpleCertificateReport holder to fetch
// values from a JAXB SimpleCertificateReport.
type SimpleCertificateReport struct {
	simpleReport *jaxb.XmlSimpleCertificateReport
}

// NewSimpleCertificateReport builds a SimpleCertificateReport around a JAXB
// XmlSimpleCertificateReport. Port of the
// SimpleCertificateReport(XmlSimpleCertificateReport) constructor.
func NewSimpleCertificateReport(simpleReport *jaxb.XmlSimpleCertificateReport) *SimpleCertificateReport {
	return &SimpleCertificateReport{simpleReport: simpleReport}
}

// GetValidationTime returns the used validation time. Port of
// getValidationTime().
func (r *SimpleCertificateReport) GetValidationTime() *time.Time {
	return xsTime(r.simpleReport.ValidationTime)
}

// GetCertificateIds returns the list of certificate ids. Port of
// getCertificateIds().
func (r *SimpleCertificateReport) GetCertificateIds() []string {
	var ids []string
	cert := r.simpleReport.Certificate
	if cert != nil {
		ids = append(ids, cert.Id)
		for _, item := range cert.Chain {
			ids = append(ids, item.Id)
		}
	}
	return ids
}

// GetCertificateNotBefore returns the notBefore date for a given
// certificate. Port of getCertificateNotBefore(String).
func (r *SimpleCertificateReport) GetCertificateNotBefore(certificateID string) *time.Time {
	cert := r.getCertificate(certificateID)
	if cert != nil {
		return xsTime(cert.NotBefore)
	}
	return nil
}

// GetCertificateNotAfter returns the notAfter date for a given certificate.
// Port of getCertificateNotAfter(String).
func (r *SimpleCertificateReport) GetCertificateNotAfter(certificateID string) *time.Time {
	cert := r.getCertificate(certificateID)
	if cert != nil {
		return xsTime(cert.NotAfter)
	}
	return nil
}

// GetCertificateAiaUrls returns the list of AIA urls (caIssuers) for a
// given certificate. Port of getCertificateAiaUrls(String).
func (r *SimpleCertificateReport) GetCertificateAiaUrls(certificateID string) []string {
	cert := r.getCertificate(certificateID)
	if cert != nil && cert.AiaUrls != nil {
		return cert.AiaUrls.AiaUrl
	}
	return nil
}

// GetCertificateCpsUrls returns the list of CPS (Certificate Practice
// Statements) urls for a given certificate. Port of
// getCertificateCpsUrls(String).
func (r *SimpleCertificateReport) GetCertificateCpsUrls(certificateID string) []string {
	cert := r.getCertificate(certificateID)
	if cert != nil && cert.CpsUrls != nil {
		return cert.CpsUrls.CpsUrl
	}
	return nil
}

// GetCertificateCrlUrls returns the list of CRL (Certificate Revocation
// List) urls for a given certificate. Port of getCertificateCrlUrls(String).
func (r *SimpleCertificateReport) GetCertificateCrlUrls(certificateID string) []string {
	cert := r.getCertificate(certificateID)
	if cert != nil && cert.CrlUrls != nil {
		return cert.CrlUrls.CrlUrl
	}
	return nil
}

// GetCertificateOcspUrls returns the list of OCSP (Online Certificate
// Status Protocol) urls for a given certificate. Port of
// getCertificateOcspUrls(String).
func (r *SimpleCertificateReport) GetCertificateOcspUrls(certificateID string) []string {
	cert := r.getCertificate(certificateID)
	if cert != nil && cert.OcspUrls != nil {
		return cert.OcspUrls.OcspUrl
	}
	return nil
}

// GetCertificatePdsUrls returns the list of PDS (PKI Disclosure Statements)
// urls for a given certificate. Port of getCertificatePdsUrls(String).
func (r *SimpleCertificateReport) GetCertificatePdsUrls(certificateID string) []string {
	cert := r.getCertificate(certificateID)
	if cert != nil && cert.PdsUrls != nil {
		return cert.PdsUrls.PdsUrl
	}
	return nil
}

// GetCertificateCommonName returns the commonName attribute for a given
// certificate, when available. Port of getCertificateCommonName(String).
func (r *SimpleCertificateReport) GetCertificateCommonName(certificateID string) string {
	return derefStr(r.subjectField(certificateID, func(s *jaxb.XmlSubject) *string { return s.CommonName }))
}

// GetCertificateEmail returns the email attribute for a given certificate,
// when available. Port of getCertificateEmail(String).
func (r *SimpleCertificateReport) GetCertificateEmail(certificateID string) string {
	return derefStr(r.subjectField(certificateID, func(s *jaxb.XmlSubject) *string { return s.Email }))
}

// GetCertificateGivenName returns the givenName attribute for a given
// certificate, when available. Port of getCertificateGivenName(String).
func (r *SimpleCertificateReport) GetCertificateGivenName(certificateID string) string {
	return derefStr(r.subjectField(certificateID, func(s *jaxb.XmlSubject) *string { return s.GivenName }))
}

// GetCertificateLocality returns the locality attribute for a given
// certificate, when available. Port of getCertificateLocality(String).
func (r *SimpleCertificateReport) GetCertificateLocality(certificateID string) string {
	return derefStr(r.subjectField(certificateID, func(s *jaxb.XmlSubject) *string { return s.Locality }))
}

// GetCertificateState returns the state attribute for a given certificate,
// when available. Port of getCertificateState(String).
func (r *SimpleCertificateReport) GetCertificateState(certificateID string) string {
	return derefStr(r.subjectField(certificateID, func(s *jaxb.XmlSubject) *string { return s.State }))
}

// GetCertificateCountry returns the country attribute for a given
// certificate, when available. Port of getCertificateCountry(String).
func (r *SimpleCertificateReport) GetCertificateCountry(certificateID string) string {
	return derefStr(r.subjectField(certificateID, func(s *jaxb.XmlSubject) *string { return s.Country }))
}

// GetCertificateOrganizationName returns the organizationName attribute for
// a given certificate, when available. Port of
// getCertificateOrganizationName(String).
func (r *SimpleCertificateReport) GetCertificateOrganizationName(certificateID string) string {
	return derefStr(r.subjectField(certificateID, func(s *jaxb.XmlSubject) *string { return s.OrganizationName }))
}

// GetCertificateOrganizationUnit returns the organizationUnit attribute for
// a given certificate, when available. Port of
// getCertificateOrganizationUnit(String).
func (r *SimpleCertificateReport) GetCertificateOrganizationUnit(certificateID string) string {
	return derefStr(r.subjectField(certificateID, func(s *jaxb.XmlSubject) *string { return s.OrganizationUnit }))
}

// GetCertificatePseudonym returns the pseudonym attribute for a given
// certificate, when available. Port of getCertificatePseudonym(String).
func (r *SimpleCertificateReport) GetCertificatePseudonym(certificateID string) string {
	return derefStr(r.subjectField(certificateID, func(s *jaxb.XmlSubject) *string { return s.Pseudonym }))
}

// GetCertificateSurname returns the surname attribute for a given
// certificate, when available. Port of getCertificateSurname(String).
func (r *SimpleCertificateReport) GetCertificateSurname(certificateID string) string {
	return derefStr(r.subjectField(certificateID, func(s *jaxb.XmlSubject) *string { return s.Surname }))
}

// subjectField reads a field of a certificate's Subject, mirroring the
// direct cert.getSubject().getXxx() chain the Java getters use (an
// unchecked NPE risk when the certificate has no Subject; see the file
// header - this port preserves it, since getCertificate never returns a
// XmlChainItem with a nil Subject in any real report, only when
// certificateId itself does not resolve to a certificate, in which case
// there is no method call to make at all).
func (r *SimpleCertificateReport) subjectField(certificateID string, get func(*jaxb.XmlSubject) *string) *string {
	cert := r.getCertificate(certificateID)
	if cert == nil {
		return nil
	}
	return get(cert.Subject)
}

// GetCertificateIndication returns the Indication (result of validation)
// for a given certificate. Port of getCertificateIndication(String).
func (r *SimpleCertificateReport) GetCertificateIndication(certificateID string) enumerations.Indication {
	cert := r.getCertificate(certificateID)
	if cert != nil {
		return cert.Indication.Indication()
	}
	return ""
}

// GetCertificateSubIndication returns the SubIndication (result of
// validation) for a given certificate. Port of
// getCertificateSubIndication(String).
func (r *SimpleCertificateReport) GetCertificateSubIndication(certificateID string) enumerations.SubIndication {
	cert := r.getCertificate(certificateID)
	if cert != nil && cert.SubIndication != nil {
		return cert.SubIndication.SubIndication()
	}
	return ""
}

// GetCertificateRevocationDate returns the revocation date for a given
// certificate, or nil. Port of getCertificateRevocationDate(String).
func (r *SimpleCertificateReport) GetCertificateRevocationDate(certificateID string) *time.Time {
	cert := r.getCertificate(certificateID)
	if cert != nil && cert.Revocation != nil {
		return xsTime(cert.Revocation.RevocationDate)
	}
	return nil
}

// GetCertificateRevocationReason returns the revocation reason for a given
// certificate, or "". Port of getCertificateRevocationReason(String).
func (r *SimpleCertificateReport) GetCertificateRevocationReason(certificateID string) enumerations.RevocationReason {
	cert := r.getCertificate(certificateID)
	if cert != nil && cert.Revocation != nil && cert.Revocation.RevocationReason != nil {
		return cert.Revocation.RevocationReason.RevocationReason()
	}
	return ""
}

// GetX509ValidationErrors retrieves the ETSI EN 319 102-1 X.509 certificate
// validation errors for a given certificate by id. Port of
// getX509ValidationErrors(String).
func (r *SimpleCertificateReport) GetX509ValidationErrors(certificateID string) []Message {
	cert := r.getCertificate(certificateID)
	if cert != nil && cert.X509ValidationDetails != nil {
		return convertMessages(cert.X509ValidationDetails.Error)
	}
	return nil
}

// GetX509ValidationWarnings retrieves the ETSI EN 319 102-1 X.509
// certificate validation warnings for a given certificate by id. Port of
// getX509ValidationWarnings(String).
func (r *SimpleCertificateReport) GetX509ValidationWarnings(certificateID string) []Message {
	cert := r.getCertificate(certificateID)
	if cert != nil && cert.X509ValidationDetails != nil {
		return convertMessages(cert.X509ValidationDetails.Warning)
	}
	return nil
}

// GetX509ValidationInfo retrieves the ETSI EN 319 102-1 X.509 certificate
// validation information messages for a given certificate by id. Port of
// getX509ValidationInfo(String).
func (r *SimpleCertificateReport) GetX509ValidationInfo(certificateID string) []Message {
	cert := r.getCertificate(certificateID)
	if cert != nil && cert.X509ValidationDetails != nil {
		return convertMessages(cert.X509ValidationDetails.Info)
	}
	return nil
}

// GetQualificationErrorsAtIssuanceTime retrieves the qualification
// process's errors for a given certificate by id at issuance time. Port of
// getQualificationErrorsAtIssuanceTime(String).
func (r *SimpleCertificateReport) GetQualificationErrorsAtIssuanceTime(certificateID string) []Message {
	cert := r.getCertificate(certificateID)
	if cert != nil && cert.QualificationDetailsAtIssuance != nil {
		return convertMessages(cert.QualificationDetailsAtIssuance.Error)
	}
	return nil
}

// GetQualificationWarningsAtIssuanceTime retrieves the qualification
// process's warnings for a given certificate by id at issuance time. Port
// of getQualificationWarningsAtIssuanceTime(String).
func (r *SimpleCertificateReport) GetQualificationWarningsAtIssuanceTime(certificateID string) []Message {
	cert := r.getCertificate(certificateID)
	if cert != nil && cert.QualificationDetailsAtIssuance != nil {
		return convertMessages(cert.QualificationDetailsAtIssuance.Warning)
	}
	return nil
}

// GetQualificationInfoAtIssuanceTime retrieves the qualification process's
// information messages for a given certificate by id at issuance time. Port
// of getQualificationInfoAtIssuanceTime(String).
func (r *SimpleCertificateReport) GetQualificationInfoAtIssuanceTime(certificateID string) []Message {
	cert := r.getCertificate(certificateID)
	if cert != nil && cert.QualificationDetailsAtIssuance != nil {
		return convertMessages(cert.QualificationDetailsAtIssuance.Info)
	}
	return nil
}

// GetQualificationErrorsAtValidationTime retrieves the qualification
// process's errors for a given certificate by id at validation time. Port
// of getQualificationErrorsAtValidationTime(String).
func (r *SimpleCertificateReport) GetQualificationErrorsAtValidationTime(certificateID string) []Message {
	cert := r.getCertificate(certificateID)
	if cert != nil && cert.QualificationDetailsAtValidation != nil {
		return convertMessages(cert.QualificationDetailsAtValidation.Error)
	}
	return nil
}

// GetQualificationWarningsAtValidationTime retrieves the qualification
// process's warnings for a given certificate by id at validation time. Port
// of getQualificationWarningsAtValidationTime(String).
func (r *SimpleCertificateReport) GetQualificationWarningsAtValidationTime(certificateID string) []Message {
	cert := r.getCertificate(certificateID)
	if cert != nil && cert.QualificationDetailsAtValidation != nil {
		return convertMessages(cert.QualificationDetailsAtValidation.Warning)
	}
	return nil
}

// GetQualificationInfoAtValidationTime retrieves the qualification
// process's information messages for a given certificate by id at
// validation time. Port of getQualificationInfoAtValidationTime(String).
func (r *SimpleCertificateReport) GetQualificationInfoAtValidationTime(certificateID string) []Message {
	cert := r.getCertificate(certificateID)
	if cert != nil && cert.QualificationDetailsAtValidation != nil {
		return convertMessages(cert.QualificationDetailsAtValidation.Info)
	}
	return nil
}

// GetQWACValidationErrors retrieves the QWAC validation process's errors
// for a given certificate by id. Port of getQWACValidationErrors(String).
func (r *SimpleCertificateReport) GetQWACValidationErrors(certificateID string) []Message {
	cert := r.getCertificate(certificateID)
	if cert != nil && cert.QwacDetails != nil {
		return convertMessages(cert.QwacDetails.Error)
	}
	return nil
}

// GetQWACValidationWarnings retrieves the QWAC validation process's
// warnings for a given certificate by id. Port of
// getQWACValidationWarnings(String).
func (r *SimpleCertificateReport) GetQWACValidationWarnings(certificateID string) []Message {
	cert := r.getCertificate(certificateID)
	if cert != nil && cert.QwacDetails != nil {
		return convertMessages(cert.QwacDetails.Warning)
	}
	return nil
}

// GetQWACValidationInfo retrieves the QWAC validation process's information
// messages for a given certificate by id. Port of
// getQWACValidationInfo(String).
func (r *SimpleCertificateReport) GetQWACValidationInfo(certificateID string) []Message {
	cert := r.getCertificate(certificateID)
	if cert != nil && cert.QwacDetails != nil {
		return convertMessages(cert.QwacDetails.Info)
	}
	return nil
}

// GetCertificateApprovalStatusErrorsAtIssuanceTime retrieves the TS 119
// 602/605 certificate approval status process's errors for a given
// certificate by id at issuance time for the given CertificateApprovalStatus.
// Port of getCertificateApprovalStatusErrorsAtIssuanceTime(String,
// CertificateApprovalStatus).
func (r *SimpleCertificateReport) GetCertificateApprovalStatusErrorsAtIssuanceTime(certificateID string, status enumerations.CertificateApprovalStatus) []Message {
	cert := r.getCertificate(certificateID)
	if cert != nil && cert.CertificateApprovalStatusAtIssuanceTime != nil {
		x := getXmlCertificateApprovalStatus(&cert.CertificateApprovalStatusAtIssuanceTime.XmlCertificateApprovalStatusAtTime, status)
		if x != nil && x.Details != nil {
			return convertMessages(x.Details.Error)
		}
	}
	return nil
}

// GetCertificateApprovalStatusWarningsAtIssuanceTime retrieves the TS 119
// 602/605 certificate approval status process's warnings for a given
// certificate by id at issuance time for the given CertificateApprovalStatus.
// Port of getCertificateApprovalStatusWarningsAtIssuanceTime(String,
// CertificateApprovalStatus).
func (r *SimpleCertificateReport) GetCertificateApprovalStatusWarningsAtIssuanceTime(certificateID string, status enumerations.CertificateApprovalStatus) []Message {
	cert := r.getCertificate(certificateID)
	if cert != nil && cert.CertificateApprovalStatusAtIssuanceTime != nil {
		x := getXmlCertificateApprovalStatus(&cert.CertificateApprovalStatusAtIssuanceTime.XmlCertificateApprovalStatusAtTime, status)
		if x != nil && x.Details != nil {
			return convertMessages(x.Details.Warning)
		}
	}
	return nil
}

// GetCertificateApprovalStatusInfoAtIssuanceTime retrieves the TS 119
// 602/605 certificate approval status process's information messages for a
// given certificate by id at issuance time for the given
// CertificateApprovalStatus. Port of
// getCertificateApprovalStatusInfoAtIssuanceTime(String,
// CertificateApprovalStatus).
func (r *SimpleCertificateReport) GetCertificateApprovalStatusInfoAtIssuanceTime(certificateID string, status enumerations.CertificateApprovalStatus) []Message {
	cert := r.getCertificate(certificateID)
	if cert != nil && cert.CertificateApprovalStatusAtIssuanceTime != nil {
		x := getXmlCertificateApprovalStatus(&cert.CertificateApprovalStatusAtIssuanceTime.XmlCertificateApprovalStatusAtTime, status)
		if x != nil && x.Details != nil {
			return convertMessages(x.Details.Info)
		}
	}
	return nil
}

// GetCertificateApprovalStatusErrorsAtValidationTime retrieves the TS 119
// 602/605 certificate approval status process's errors for a given
// certificate by id at validation time for the given
// CertificateApprovalStatus. Port of
// getCertificateApprovalStatusErrorsAtValidationTime(String,
// CertificateApprovalStatus).
func (r *SimpleCertificateReport) GetCertificateApprovalStatusErrorsAtValidationTime(certificateID string, status enumerations.CertificateApprovalStatus) []Message {
	cert := r.getCertificate(certificateID)
	if cert != nil && cert.CertificateApprovalStatusAtValidationTime != nil {
		x := getXmlCertificateApprovalStatus(&cert.CertificateApprovalStatusAtValidationTime.XmlCertificateApprovalStatusAtTime, status)
		if x != nil && x.Details != nil {
			return convertMessages(x.Details.Error)
		}
	}
	return nil
}

// GetCertificateApprovalStatusWarningsAtValidationTime retrieves the TS 119
// 602/605 certificate approval status process's warnings for a given
// certificate by id at validation time for the given
// CertificateApprovalStatus. Port of
// getCertificateApprovalStatusWarningsAtValidationTime(String,
// CertificateApprovalStatus).
func (r *SimpleCertificateReport) GetCertificateApprovalStatusWarningsAtValidationTime(certificateID string, status enumerations.CertificateApprovalStatus) []Message {
	cert := r.getCertificate(certificateID)
	if cert != nil && cert.CertificateApprovalStatusAtValidationTime != nil {
		x := getXmlCertificateApprovalStatus(&cert.CertificateApprovalStatusAtValidationTime.XmlCertificateApprovalStatusAtTime, status)
		if x != nil && x.Details != nil {
			return convertMessages(x.Details.Warning)
		}
	}
	return nil
}

// GetCertificateApprovalStatusInfoAtValidationTime retrieves the TS 119
// 602/605 certificate approval status process's information messages for a
// given certificate by id at validation time for the given
// CertificateApprovalStatus. Port of
// getCertificateApprovalStatusInfoAtValidationTime(String,
// CertificateApprovalStatus).
func (r *SimpleCertificateReport) GetCertificateApprovalStatusInfoAtValidationTime(certificateID string, status enumerations.CertificateApprovalStatus) []Message {
	cert := r.getCertificate(certificateID)
	if cert != nil && cert.CertificateApprovalStatusAtValidationTime != nil {
		x := getXmlCertificateApprovalStatus(&cert.CertificateApprovalStatusAtValidationTime.XmlCertificateApprovalStatusAtTime, status)
		if x != nil && x.Details != nil {
			return convertMessages(x.Details.Info)
		}
	}
	return nil
}

// GetQualificationAtCertificateIssuance returns the qualification of the
// first certificate at its issuance. Port of
// getQualificationAtCertificateIssuance().
func (r *SimpleCertificateReport) GetQualificationAtCertificateIssuance() enumerations.CertificateQualification {
	cert := r.getFirstCertificate()
	if cert.QualificationAtIssuance == nil {
		return ""
	}
	return cert.QualificationAtIssuance.CertificateQualification()
}

// GetQualificationAtValidationTime returns the qualification of the first
// certificate at the validation time. Port of
// getQualificationAtValidationTime().
func (r *SimpleCertificateReport) GetQualificationAtValidationTime() enumerations.CertificateQualification {
	cert := r.getFirstCertificate()
	if cert.QualificationAtValidation == nil {
		return ""
	}
	return cert.QualificationAtValidation.CertificateQualification()
}

// GetCertificateApprovalStatusAtCertificateIssuance returns the
// qualification of the first certificate at its issuance. Port of
// getCertificateApprovalStatusAtCertificateIssuance().
func (r *SimpleCertificateReport) GetCertificateApprovalStatusAtCertificateIssuance() []enumerations.CertificateApprovalStatus {
	cert := r.getFirstCertificate()
	if cert.CertificateApprovalStatusAtIssuanceTime == nil {
		return nil
	}
	return toCertificateApprovalStatusList(&cert.CertificateApprovalStatusAtIssuanceTime.XmlCertificateApprovalStatusAtTime)
}

// GetCertificateApprovalStatusAtValidationTime returns the qualification of
// the first certificate at the validation time. Port of
// getCertificateApprovalStatusAtValidationTime().
func (r *SimpleCertificateReport) GetCertificateApprovalStatusAtValidationTime() []enumerations.CertificateApprovalStatus {
	cert := r.getFirstCertificate()
	if cert.CertificateApprovalStatusAtValidationTime == nil {
		return nil
	}
	return toCertificateApprovalStatusList(&cert.CertificateApprovalStatusAtValidationTime.XmlCertificateApprovalStatusAtTime)
}

// GetQWACProfile returns the QWAC validation result as per ETSI TS 119
// 411-5. NOTE: Applicable only when validation process is executed using
// the QWACValidator. Port of getQWACProfile().
func (r *SimpleCertificateReport) GetQWACProfile() enumerations.QWACProfile {
	cert := r.getFirstCertificate()
	if cert.QwacProfile == nil {
		return ""
	}
	return cert.QwacProfile.QWACProfile()
}

// GetTLSBindingSignature gets the TLS Certificate Binding signature. NOTE:
// Applicable only when validation process is executed using the
// QWACValidator. Port of getTLSBindingSignature().
func (r *SimpleCertificateReport) GetTLSBindingSignature() *jaxb.XmlSignature {
	cert := r.getFirstCertificate()
	return cert.TLSBindingSignature
}

// GetTLSBindingSignatureIndication gets the Indication for the TLS
// Certificate Binding signature validation. NOTE: Applicable only when
// validation process is executed using the QWACValidator. Port of
// getTLSBindingSignatureIndication().
func (r *SimpleCertificateReport) GetTLSBindingSignatureIndication() enumerations.Indication {
	sig := r.GetTLSBindingSignature()
	if sig != nil {
		return sig.Indication.Indication()
	}
	return ""
}

// GetTLSBindingSignatureSubIndication gets the SubIndication for the TLS
// Certificate Binding signature validation. NOTE: Applicable only when
// validation process is executed using the QWACValidator. Port of
// getTLSBindingSignatureSubIndication().
func (r *SimpleCertificateReport) GetTLSBindingSignatureSubIndication() enumerations.SubIndication {
	sig := r.GetTLSBindingSignature()
	if sig != nil && sig.SubIndication != nil {
		return sig.SubIndication.SubIndication()
	}
	return ""
}

// GetTLSBindingSignatureIssuerCertificate gets the issuer certificate of
// the TLS Certificate Binding signature. NOTE: Applicable only when
// validation process is executed using the QWACValidator. Port of
// getTLSBindingSignatureIssuerCertificate().
func (r *SimpleCertificateReport) GetTLSBindingSignatureIssuerCertificate() *jaxb.XmlChainItem {
	sig := r.GetTLSBindingSignature()
	if sig != nil && len(sig.Chain) > 0 {
		return sig.Chain[0]
	}
	return nil
}

// GetTLSBindingSignatureIssuerQualificationAtCertificateIssuance returns
// the qualification of the first certificate at its issuance. Port of
// getTLSBindingSignatureIssuerQualificationAtCertificateIssuance().
func (r *SimpleCertificateReport) GetTLSBindingSignatureIssuerQualificationAtCertificateIssuance() enumerations.CertificateQualification {
	item := r.GetTLSBindingSignatureIssuerCertificate()
	if item != nil && item.QualificationAtIssuance != nil {
		return item.QualificationAtIssuance.CertificateQualification()
	}
	return ""
}

// GetTLSBindingSignatureIssuerQualificationAtValidationTime returns the
// qualification of the first certificate at the validation time. Port of
// getTLSBindingSignatureIssuerQualificationAtValidationTime().
func (r *SimpleCertificateReport) GetTLSBindingSignatureIssuerQualificationAtValidationTime() enumerations.CertificateQualification {
	item := r.GetTLSBindingSignatureIssuerCertificate()
	if item != nil && item.QualificationAtValidation != nil {
		return item.QualificationAtValidation.CertificateQualification()
	}
	return ""
}

// GetTLSBindingSignatureIssuerCertificateQWACProfile gets the Indication
// for the TLS Certificate Binding signature validation. NOTE: Applicable
// only when validation process is executed using the QWACValidator. Port of
// getTLSBindingSignatureIssuerCertificateQWACProfile().
func (r *SimpleCertificateReport) GetTLSBindingSignatureIssuerCertificateQWACProfile() enumerations.QWACProfile {
	item := r.GetTLSBindingSignatureIssuerCertificate()
	if item != nil && item.QwacProfile != nil {
		return item.QwacProfile.QWACProfile()
	}
	return ""
}

// GetTrustAnchorVATNumbers returns a set of trust anchor VAT numbers. Port
// of getTrustAnchorVATNumbers(); a nil TrustServiceProviderRegistrationId
// (Java allows a HashSet<String> to hold a null element) becomes the empty
// string, the closest Go stand-in.
func (r *SimpleCertificateReport) GetTrustAnchorVATNumbers() map[string]struct{} {
	result := map[string]struct{}{}
	cert := r.getTrustAnchorCertificate()
	if cert != nil && cert.TrustAnchors != nil {
		for _, ta := range cert.TrustAnchors.TrustAnchor {
			result[derefStr(ta.TrustServiceProviderRegistrationId)] = struct{}{}
		}
	}
	return result
}

// getTrustAnchorCertificate returns the private getTrustAnchorCertificate().
func (r *SimpleCertificateReport) getTrustAnchorCertificate() *jaxb.XmlChainItem {
	cert := r.simpleReport.Certificate
	if cert == nil {
		return nil
	}
	if isTrustAnchor(cert) {
		return cert
	}
	for _, item := range cert.Chain {
		if isTrustAnchor(item) {
			return item
		}
	}
	return nil
}

func isTrustAnchor(item *jaxb.XmlChainItem) bool {
	return item != nil && item.TrustAnchors != nil && len(item.TrustAnchors.TrustAnchor) > 0
}

// getFirstCertificate is the port of the private getFirstCertificate().
func (r *SimpleCertificateReport) getFirstCertificate() *jaxb.XmlChainItem {
	return r.simpleReport.Certificate
}

// getCertificate is the port of the private getCertificate(String).
func (r *SimpleCertificateReport) getCertificate(certificateID string) *jaxb.XmlChainItem {
	if certificateID == "" {
		return nil
	}
	cert := r.simpleReport.Certificate
	if cert != nil {
		if certificateID == cert.Id {
			return cert
		}
		for _, item := range cert.Chain {
			if certificateID == item.Id {
				return item
			}
		}
	}
	sig := r.GetTLSBindingSignature()
	if sig != nil {
		for _, item := range sig.Chain {
			if certificateID == item.Id {
				return item
			}
		}
	}
	return nil
}

// GetJaxbModel returns the jaxb model of the simple certificate report.
// Port of getJaxbModel().
func (r *SimpleCertificateReport) GetJaxbModel() *jaxb.XmlSimpleCertificateReport {
	return r.simpleReport
}

// ------------------------------------------------------------------ helpers

// xsTime converts an optional jaxb.XSDateTime field to *time.Time,
// tolerating nil.
func xsTime(d *jaxb.XSDateTime) *time.Time {
	if d == nil {
		return nil
	}
	t := d.Time()
	return &t
}

// derefStr dereferences an optional *string, tolerating nil.
func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// convertMessages converts a list of XmlMessage into a list of Message.
// Port of the private convert(Collection<XmlMessage>).
func convertMessages(msgs []*jaxb.XmlMessage) []Message {
	out := make([]Message, 0, len(msgs))
	for _, m := range msgs {
		if m == nil {
			continue
		}
		out = append(out, Message{Key: derefStr(m.Key), Value: m.Value})
	}
	return out
}

// toCertificateApprovalStatusList converts a
// XmlCertificateApprovalStatusAtTime's list of XmlCertificateApprovalStatus
// into a list of enumerations.CertificateApprovalStatus. Port of the
// private toCertificateApprovalStatusList(XmlCertificateApprovalStatusAtTime).
func toCertificateApprovalStatusList(at *jaxb.XmlCertificateApprovalStatusAtTime) []enumerations.CertificateApprovalStatus {
	if at == nil {
		return nil
	}
	if len(at.CertificateApprovalStatus) == 0 {
		return []enumerations.CertificateApprovalStatus{}
	}
	out := make([]enumerations.CertificateApprovalStatus, 0, len(at.CertificateApprovalStatus))
	for _, x := range at.CertificateApprovalStatus {
		out = append(out, toCertificateApprovalStatus(x))
	}
	return out
}

// toCertificateApprovalStatus converts a single XmlCertificateApprovalStatus
// into an enumerations.CertificateApprovalStatus, falling back to a
// synthetic CERT_FOR_UNKNOWN-labelled status when no registered LoTELoader
// recognises the (listType, serviceTypeIdentifier, serviceStatus) triple.
// Port of the private
// toCertificateApprovalStatus(XmlCertificateApprovalStatus).
func toCertificateApprovalStatus(x *jaxb.XmlCertificateApprovalStatus) enumerations.CertificateApprovalStatus {
	if x == nil {
		return nil
	}
	lt := xmlListType(x.ListType)
	sti := xmlServiceTypeIdentifier(x.ServiceTypeIdentifier)
	status := xmlServiceStatus(x.ServiceStatus)
	result := enumerations.CertificateApprovalStatusFromDefinition(lt, sti, status)
	if result != nil && result.Label() != "" && enumerations.CertificateApprovalStatusEnum_CERT_FOR_UNKNOWN != result {
		return result
	}
	return enumerations.NewCertificateApprovalStatus(
		enumerations.CertificateApprovalStatusEnum_CERT_FOR_UNKNOWN.Label(), lt, sti, status)
}

// getXmlCertificateApprovalStatus finds the XmlCertificateApprovalStatus of
// at whose ListType/ServiceTypeIdentifier URIs match certificateApprovalStatus.
// Port of the private
// getXmlCertificateApprovalStatus(XmlCertificateApprovalStatusAtTime,
// CertificateApprovalStatus).
func getXmlCertificateApprovalStatus(at *jaxb.XmlCertificateApprovalStatusAtTime, certificateApprovalStatus enumerations.CertificateApprovalStatus) *jaxb.XmlCertificateApprovalStatus {
	if at == nil || certificateApprovalStatus == nil {
		return nil
	}
	for _, x := range at.CertificateApprovalStatus {
		if x == nil {
			continue
		}
		wantListType := certificateApprovalStatus.ListType()
		gotListType := xmlListType(x.ListType)
		if wantListType == nil || gotListType == nil || wantListType.URI() != gotListType.URI() {
			continue
		}
		wantSTI := certificateApprovalStatus.ServiceTypeIdentifier()
		gotSTI := xmlServiceTypeIdentifier(x.ServiceTypeIdentifier)
		if wantSTI == nil || gotSTI == nil || wantSTI.URI() != gotSTI.URI() {
			continue
		}
		return x
	}
	return nil
}

func xmlListType(v *jaxb.ListTypeValue) enumerations.ListType {
	if v == nil {
		return nil
	}
	return v.Value
}

func xmlServiceTypeIdentifier(v *jaxb.LoTEServiceTypeIdentifierValue) enumerations.LoTEServiceTypeIdentifier {
	if v == nil {
		return nil
	}
	return v.Value
}

func xmlServiceStatus(v *jaxb.LoTEServiceStatusValue) enumerations.LoTEServiceStatus {
	if v == nil {
		return nil
	}
	return v.Value
}
