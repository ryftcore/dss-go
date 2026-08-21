// Java-oracle tests for the dto/condition ports (cert_subject_dn_attribute_condition.go,
// composite_condition.go, extended_key_usage_condition.go, key_usage_condition.go,
// policy_id_condition.go, qc_statement_condition.go).
//
// Every certificate and every expected verdict below is transcribed verbatim from the upstream
// JUnit suites (dss-tsl-validation 6.5.RC1
// eu.europa.esig.dss.tsl.dto.condition.{CertSubjectDNAttributeCondition,Composite,
// ExtendedKeyUsage,KeyUsage,PolicyId,QcStatement}ConditionTest); the Go assertions mirror them
// one for one. Test vectors are ported, not the JUnit code, per PORTING.md.
//
// The upstream tests are shared by all six conditions (they exercise the composites through the
// leaf conditions), so they live in one Go file rather than six - the per-Java-file rule governs
// the ported sources, not their oracle fixtures.
package tsl

import (
	"testing"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi"
)

// BCStyle OID constants the upstream tests reach through org.bouncycastle.asn1.x500.style.BCStyle.
const (
	conditionTestOIDCountry      = "2.5.4.6"              // BCStyle.C
	conditionTestOIDEmailAddress = "1.2.840.113549.1.9.1" // BCStyle.EmailAddress
	conditionTestOIDOrganization = "2.5.4.10"             // organizationName
	conditionTestOIDOCSPSigning  = "1.3.6.1.5.5.7.3.9"    // id-kp-OCSPSigning
)

// ETSIQCObjectIdentifiers OIDs the upstream QcStatementConditionTest reaches through
// org.bouncycastle.asn1.x509.qualified.ETSIQCObjectIdentifiers.
const (
	conditionTestOIDQcCompliance = "0.4.0.1862.1.1"
	conditionTestOIDQcSSCD       = "0.4.0.1862.1.4"
	conditionTestOIDQcPds        = "0.4.0.1862.1.5"
	conditionTestOIDQcType       = "0.4.0.1862.1.6"
)

// belgiumOCSPResponderCertificate is the certificate the CertSubjectDNAttributeCondition,
// CompositeCondition and ExtendedKeyUsageCondition upstream tests share.
const belgiumOCSPResponderCertificate = "MIIEXjCCAkagAwIBAgILBAAAAAABWLd6HkYwDQYJKoZIhvcNAQELBQAwMzELMAkGA1UEBhMCQkUxEzARBgNVBAMTCkNpdGl6ZW4gQ0ExDzANBgNVBAUTBjIwMTYzMTAeFw0xNjEyMTAxMTAwMDBaFw0xODAxMjkxMTAwMDBaMC4xHzAdBgNVBAMTFkJlbGdpdW0gT0NTUCBSZXNwb25kZXIxCzAJBgNVBAYTAkJFMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAzD0B0c4gBx/wumeE2l/Wcz5FoMSUIuRNIySH2pJ3yfKR/u/FWCOzcrJvDMdmgzR33zGb4/fZel9YlI6xcN08Yd7GkP0/WtbHUhGUPERV76Vvyrk2K/EH/IG2gtxYB+7pkA/ZZycdyjc4IxHzBOiGofP9lDkPD05GSqI7MjVf6sNkZSnHcQSKwkaCGhAshJMjHzShEsSzOgX9kXceBFPTt6Hd2prVmnMTyAwURbQ6gFHbgfxB8JLMya95U6391nGQC66ScH1GhIwd9KSn+yBY0cazJ3nIrc8wd0yGYBgPK78jN3MvAsb1ydfs7kE+Wf95z9oRMiw62Glxh/ksLS/tTQIDAQABo3gwdjAOBgNVHQ8BAf8EBAMCB4AwHQYDVR0OBBYEFBgKRBywCTroyvAErr7p657558Y9MBMGA1UdJQQMMAoGCCsGAQUFBwMJMB8GA1UdIwQYMBaAFM6Al2fQrdlOxJlqgCcikM0RNRCHMA8GCSsGAQUFBzABBQQCBQAwDQYJKoZIhvcNAQELBQADggIBAFuZrqcwt23UiiJdRst66MEBRyKbgPsQM81Uq4FVrAnV8z3l8DDUv+A29KzCPO0GnHSatqA7DNhhMzoBRC42PqCpuvrj8VEWHd43AuPOLaikE04a5tVh6DgW8b00s6Yyf/PuDHCsg2C2MqY71MUR9GcnI7ngR2SyWQGpbsf/wfjujNxEB0+SOwMDTgIAikaueHGZbYkwvlRpL6wm2ENvrE8OvKt7NlNsaWJ4KtQo0QS5Ku+Y2BDA3bX+g8eNLQkaXTycgL4X3MyE5pBOl1OW3KOjJdfyLF+Sii+JKjNf8ZQWk0xvkBEI+nhCzDXhtKAcrkTKlXE25MiUnYoRsXkXgrzYftxAMxvFOXJji/hnX5Fe/3SBAHaE+jU6yC5nk6Q9ERii8mL0nHouMlZWSiAuXtlZDFrzwtLD2ITBECe4X60BDQfb/caO2u3HcWoG1AOvGxfQB0cMmP2njCdDf8UOqryiyky4t7Jj3ghOvETjWlwMw5ObhZ8yj8p6qFAt7+EVJfpUc1gDAolS/hJoLzohbL5LnCAnUAWsFpvG3qW1ky+X0MePXi6q/boqj2tcC4IDdsYS6RHPBvzl5+yLDccrGx1s/7vQYTMNyX0dYZzuxFZxx0bttWfjqLz3hFHlAEVmLCyUkSz761CbaT9u/G4tPP4Q8ApFfSskPI57lbLWIcwP"

// skTLCertificate is the certificate the upstream CertSubjectDNAttributeConditionTest#dss1911
// case uses.
const skTLCertificate = "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"

// estonianESTEIDCertificate is the certificate the upstream KeyUsageConditionTest and
// PolicyIdConditionTest share.
const estonianESTEIDCertificate = "MIID3DCCAsSgAwIBAgIER/idhzANBgkqhkiG9w0BAQUFADBbMQswCQYDVQQGEwJFRTEiMCAGA1UEChMZQVMgU2VydGlmaXRzZWVyaW1pc2tlc2t1czEPMA0GA1UECxMGRVNURUlEMRcwFQYDVQQDEw5FU1RFSUQtU0sgMjAwNzAeFw0wODA0MDYwOTUzMDlaFw0xMjAzMDUyMjAwMDBaMIGWMQswCQYDVQQGEwJFRTEPMA0GA1UEChMGRVNURUlEMRowGAYDVQQLExFkaWdpdGFsIHNpZ25hdHVyZTEiMCAGA1UEAxMZU0lOSVZFRSxWRUlLTywzNjcwNjAyMDIxMDEQMA4GA1UEBBMHU0lOSVZFRTEOMAwGA1UEKhMFVkVJS08xFDASBgNVBAUTCzM2NzA2MDIwMjEwMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCGRN42R9e6VEHMCyvacuubjtm1+5Kk92WgIgtWA8hY8DW2iNvQJ3jOF5XlVIyIDTwl2JVKxWKhXX+8+yNFPpqAK43IINcmMfznw/KcR7jACGNuTrivA9HrvRiqDzTg5E1rktjho6OkDkdV3dgOLB2wyhVm2anNpICfrUq8c09HPwIDMMP5o4HvMIHsMA4GA1UdDwEB/wQEAwIGQDA8BgNVHR8ENTAzMDGgL6AthitodHRwOi8vd3d3LnNrLmVlL2NybHMvZXN0ZWlkL2VzdGVpZDIwMDcuY3JsMFEGA1UdIARKMEgwRgYLKwYBBAHOHwEBAQEwNzASBggrBgEFBQcCAjAGGgRub25lMCEGCCsGAQUFBwIBFhVodHRwOi8vd3d3LnNrLmVlL2Nwcy8wHwYDVR0jBBgwFoAUSAbevoyHV5WAeGP6nCMrK6A6GHUwHQYDVR0OBBYEFJAJUyDrH3rdxTStU+LDa6aHdE8dMAkGA1UdEwQCMAAwDQYJKoZIhvcNAQEFBQADggEBAA5qjfeuTdOoEtatiA9hpjDHzyqN1PROcaPrABXGqpLxcHbLVr7xmovILAjxS9fJAw28u9ZE3asRNa9xgQNTeX23mMlojJAYVbYCeIeJ6jtsRiCo34wgvO3CtVfO3+C1T8Du5XLCHa6SoT8SpCApW+Crwe+6eCZDmv2NKTjhn1wCCNO2e8HuSt+pTUNBTUB+rkvF4KO9VnuzRzT7zN7AUdW4OFF3bI+9+VmW3t9vq1zDOxNTdBkCM3zm5TRa8ZtyAPL48bW19JAcYzQLjPGORwoIRNSXdVTqX+cDiw2wbmb2IhPdxRqN9uPwU1x/ltZZ3W5GzJ1t8JeQN7PuGM0OHqE="

// The three certificates the upstream QcStatementConditionTest loads.
const (
	uaESignCertificate  = "MIIE/jCCBKSgAwIBAgIUNHSLOUor3i8EAAAAzfVDAFlVhwAwCgYIKoZIzj0EAwIwgdgxIDAeBgNVBAoMF1N0YXRlIGVudGVycHJpc2UgIkRJSUEiMTAwLgYDVQQLDCdEZXBhcnRtZW50IG9mIEVsZWN0cm9uaWMgVHJ1c3QgU2VydmljZXMxMjAwBgNVBAMMKSJESUlBIi4gUXVhbGlmaWVkIFRydXN0IFNlcnZpY2VzIFByb3ZpZGVyMRkwFwYDVQQFExBVQS00MzM5NTAzMy0xMTEwMQswCQYDVQQGEwJVQTENMAsGA1UEBwwES3lpdjEXMBUGA1UEYQwOTlRSVUEtNDMzOTUwMzMwHhcNMjIwNzA2MjEwMDAwWhcNMjQwNzA2MjA1OTU5WjCBtDEgMB4GA1UECgwXU3RhdGUgZW50ZXJwcmlzZSAiRElJQSIxGjAYBgNVBAMMEVNlcmhpaSBLb3NtaW5za3lpMRMwEQYDVQQEDApLb3NtaW5za3lpMQ8wDQYDVQQqDAZTZXJoaWkxGTAXBgNVBAUTEFRJTlVBLTMxNzIwMTEzNTMxCzAJBgNVBAYTAlVBMQ0wCwYDVQQHDARLeWl2MRcwFQYDVQRhDA5OVFJVQS00MzM5NTAzMzBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABJuLYx3sca0WzRL+rP9A3R409yHCUJ25YXHxVoLkp/+yXpVzzez+dLV0As4QImqRBrihxhVgN08AG6n9GiMw3kijggJsMIICaDAdBgNVHQ4EFgQUxQr62azBDuDYo2S/YDxDY3sZgykwHwYDVR0jBBgwFoAUtHSLOUor3i9P/LHbo7XxgjSPEA4wDgYDVR0PAQH/BAQDAgZAMEkGA1UdIARCMEAwPgYJKoYkAgEBAQICMDEwLwYIKwYBBQUHAgEWI2h0dHBzOi8vY2EuaW5mb3JtanVzdC51YS9yZWdsYW1lbnQvMAkGA1UdEwQCMAAwVAYIKwYBBQUHAQMESDBGMAgGBgQAjkYBATATBgYEAI5GAQYwCQYHBACORgEGATAOBgYEAI5GAQcwBBMCVUEwFQYIKwYBBQUHCwIwCQYHBACL7EkBATBLBgNVHR8ERDBCMECgPqA8hjpodHRwOi8vY2EuaW5mb3JtanVzdC51YS9kb3dubG9hZC9jcmxzL0NBLUI0NzQ4QjM5LUZ1bGwuY3JsMEwGA1UdLgRFMEMwQaA/oD2GO2h0dHA6Ly9jYS5pbmZvcm1qdXN0LnVhL2Rvd25sb2FkL2NybHMvQ0EtQjQ3NDhCMzktRGVsdGEuY3JsMIGFBggrBgEFBQcBAQR5MHcwMgYIKwYBBQUHMAGGJmh0dHA6Ly9jYS5pbmZvcm1qdXN0LnVhL3NlcnZpY2VzL29jc3AvMEEGCCsGAQUFBzAChjVodHRwOi8vY2EuaW5mb3JtanVzdC51YS91cGxvYWRzL2NlcnRpZmljYXRlcy9kaWlhLnA3YjBHBggrBgEFBQcBCwQ7MDkwNwYIKwYBBQUHMAOGK2h0dHA6Ly9jYS5pbmZvcm1qdXN0LnVhL3NlcnZpY2VzL3RzcC9lY2RzYS8wCgYIKoZIzj0EAwIDSAAwRQIgeFeZh5I2cd03/0dZ2FoAfwI9qqZlqEMWZY9k8J31pKoCIQClSbciBvkQc/8i5o6+TLhztWYSqd8ltzQeNgzLB0Llyg=="
	eSealCertificate    = "MIIEQDCCA+agAwIBAgIUXpKHtb9tW+oCAAAAAQAAANEAAAAwCgYIKoZIzj0EAwIwgdExNjA0BgNVBAoMLU1pbmlzdHJ5IG9mIGRpZ2l0YWwgdHJhbnNmb3JtYXRpb24gb2YgVWtyYWluZTEeMBwGA1UECwwVQWRtaW5pc3RyYXRvciBJVFMgQ0NBMSgwJgYDVQQDDB9DZW50cmFsIGNlcnRpZmljYXRpb24gYXV0aG9yaXR5MRgwFgYDVQQFDA9VQS00MzIyMDg1MS0yNTYxCzAJBgNVBAYTAlVBMQ0wCwYDVQQHDARLeWl2MRcwFQYDVQRhDA5OVFJVQS00MzIyMDg1MTAeFw0yMDAxMjEwMzA4MDBaFw0yNTAxMjEwMzA4MDBaMIHeMTYwNAYDVQQKDC1NaW5pc3RyeSBvZiBkaWdpdGFsIHRyYW5zZm9ybWF0aW9uIG9mIFVrcmFpbmUxHjAcBgNVBAsMFUFkbWluaXN0cmF0b3IgSVRTIENDQTE0MDIGA1UEAwwrT0NTUC1zZXJ2ZXIgQ2VudHJhbCBjZXJ0aWZpY2F0aW9uIGF1dGhvcml0eTEZMBcGA1UEBQwQVUEtNDMyMjA4NTEtMjAyMDELMAkGA1UEBhMCVUExDTALBgNVBAcMBEt5aXYxFzAVBgNVBGEMDk5UUlVBLTQzMjIwODUxMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE7sZ2nO//6wzDZ0Gr+fMbqdxc4aYds2hk0YTJohPkLmnainmHwMytf/sEBACJRXhdToVw8YNdONCw1/n0KsosK6OCAYswggGHMB0GA1UdDgQWBBSDWEs3pLq3kZ/T2lBjeezuljxdkzAfBgNVHSMEGDAWgBRekoe1v21b6unpHaIRk+Re/FCJkTAOBgNVHQ8BAf8EBAMCB4AwEwYDVR0lBAwwCgYIKwYBBQUHAwkwPAYDVR0gBDUwMzAxBgkqhiQCAQEBAgIwJDAiBggrBgEFBQcCARYWaHR0cHM6Ly9jem8uZ292LnVhL2NwczAJBgNVHRMEAjAAMEQGCCsGAQUFBwEDBDgwNjAIBgYEAI5GAQEwCAYGBACORgEEMBMGBgQAjkYBBjAJBgcEAI5GAQYCMAsGCSqGJAIBAQECATBHBgNVHR8EQDA+MDygOqA4hjZodHRwOi8vY3pvLmdvdi51YS9kb3dubG9hZC9jcmxzL0NBLUVDRFNBLTIwMjAtRnVsbC5jcmwwSAYDVR0uBEEwPzA9oDugOYY3aHR0cDovL2N6by5nb3YudWEvZG93bmxvYWQvY3Jscy9DQS1FQ0RTQS0yMDIwLURlbHRhLmNybDAKBggqhkjOPQQDAgNIADBFAiBR0XnBOxsiTLnVS5mbWUSlvJtNR32Zvhstc728Y5USnQIhAOGsELbe01+t0IvAZoFLw1IvCrz6Yd64kCb7puqN/HjK"
	otherOidCertificate = "MIIEyTCCBG+gAwIBAgIUXpKHtb9tW+oBAAAAAQAAAOMAAAAwCgYIKoZIzj0EAwIwgdExNjA0BgNVBAoMLU1pbmlzdHJ5IG9mIGRpZ2l0YWwgdHJhbnNmb3JtYXRpb24gb2YgVWtyYWluZTEeMBwGA1UECwwVQWRtaW5pc3RyYXRvciBJVFMgQ0NBMSgwJgYDVQQDDB9DZW50cmFsIGNlcnRpZmljYXRpb24gYXV0aG9yaXR5MRgwFgYDVQQFDA9VQS00MzIyMDg1MS0yNTYxCzAJBgNVBAYTAlVBMQ0wCwYDVQQHDARLeWl2MRcwFQYDVQRhDA5OVFJVQS00MzIyMDg1MTAeFw0yMDA2MDMwNzQzMDBaFw0yNTA2MDMwNzQzMDBaMIHYMSAwHgYDVQQKDBdTdGF0ZSBlbnRlcnByaXNlICJESUlBIjEwMC4GA1UECwwnRGVwYXJ0bWVudCBvZiBFbGVjdHJvbmljIFRydXN0IFNlcnZpY2VzMTIwMAYDVQQDDCkiRElJQSIuIFF1YWxpZmllZCBUcnVzdCBTZXJ2aWNlcyBQcm92aWRlcjEZMBcGA1UEBRMQVUEtNDMzOTUwMzMtMTExMDELMAkGA1UEBhMCVUExDTALBgNVBAcMBEt5aXYxFzAVBgNVBGEMDk5UUlVBLTQzMzk1MDMzMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAER1/iqwZrmAORrFaPbwLpTaiZ5geIS6YCgZfE0Ljr3JH092E4c4It+ZvzEKE1E/g9kf0FvbjfZ38fSlZCzJDJ0qOCAhowggIWMB0GA1UdDgQWBBS0dIs5SiveL0/8sdujtfGCNI8QDjAOBgNVHQ8BAf8EBAMCAQYwPAYDVR0gBDUwMzAxBgkqhiQCAQEBAgIwJDAiBggrBgEFBQcCARYWaHR0cHM6Ly9jem8uZ292LnVhL2NwczAtBgNVHREEJjAkghBjYS5pbmZvcm1qdXN0LnVhgRBjYUBpbmZvcm1qdXN0LnVhMBIGA1UdEwEB/wQIMAYBAf8CAQAwcgYIKwYBBQUHAQMEZjBkMAgGBgQAjkYBATAIBgYEAI5GAQQwKgYGBACORgEFMCAwHhYYaHR0cHM6Ly9jem8uZ292LnVhL2Fib3V0EwJlbjAVBggrBgEFBQcLAjAJBgcEAIvsSQECMAsGCSqGJAIBAQECATAfBgNVHSMEGDAWgBRekoe1v21b6unpHaIRk+Re/FCJkTBHBgNVHR8EQDA+MDygOqA4hjZodHRwOi8vY3pvLmdvdi51YS9kb3dubG9hZC9jcmxzL0NBLUVDRFNBLTIwMjAtRnVsbC5jcmwwSAYDVR0uBEEwPzA9oDugOYY3aHR0cDovL2N6by5nb3YudWEvZG93bmxvYWQvY3Jscy9DQS1FQ0RTQS0yMDIwLURlbHRhLmNybDA8BggrBgEFBQcBAQQwMC4wLAYIKwYBBQUHMAGGIGh0dHA6Ly9jem8uZ292LnVhL3NlcnZpY2VzL29jc3AvMAoGCCqGSM49BAMCA0gAMEUCIQDjHWMxMR4ZwopSJ+1ZpahX63DgsGbrlT2j6Hg5D4924AIgdL8OrA7dpEhHbx45FqkNNS2YvcyV325GRoy0KZiO7Fw="
)

// conditionTestCertificate loads a base64-encoded certificate, standing in for
// DSSUtils.loadCertificateFromBase64EncodedString.
func conditionTestCertificate(t *testing.T, base64Encoded string) *model.CertificateToken {
	t.Helper()
	certificate, err := spi.DSSUtilsLoadCertificateFromBase64EncodedString(base64Encoded)
	if err != nil {
		t.Fatalf("unable to load the certificate: %v", err)
	}
	return certificate
}

// TestCertSubjectDNAttributeCondition_Oracle ports CertSubjectDNAttributeConditionTest#test.
func TestCertSubjectDNAttributeCondition_Oracle(t *testing.T) {
	certificate := conditionTestCertificate(t, belgiumOCSPResponderCertificate)

	condition := NewCertSubjectDNAttributeCondition([]string{conditionTestOIDCountry})
	if !condition.Check(certificate) {
		t.Error("CertSubjectDNAttributeCondition(C) should match")
	}
	if got, want := condition.String(), "CertSubjectDNAttributeCondition: [2.5.4.6]\n"; got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}

	condition = NewCertSubjectDNAttributeCondition([]string{conditionTestOIDEmailAddress})
	if condition.Check(certificate) {
		t.Error("CertSubjectDNAttributeCondition(EmailAddress) should not match")
	}
}

// TestCertSubjectDNAttributeCondition_OracleDSS1911 ports
// CertSubjectDNAttributeConditionTest#dss1911.
func TestCertSubjectDNAttributeCondition_OracleDSS1911(t *testing.T) {
	certificate := conditionTestCertificate(t, skTLCertificate)

	allNones := NewCompositeConditionWithMatchingCriteriaIndicator(enumerations.Assert_ALL)

	none42 := NewCompositeConditionWithMatchingCriteriaIndicator(enumerations.Assert_NONE)
	none42.AddChild(NewCertSubjectDNAttributeCondition([]string{"2.5.4.42"}))

	none65 := NewCompositeConditionWithMatchingCriteriaIndicator(enumerations.Assert_NONE)
	none65.AddChild(NewCertSubjectDNAttributeCondition([]string{"2.5.4.65"}))

	none4 := NewCompositeConditionWithMatchingCriteriaIndicator(enumerations.Assert_NONE)
	// not present
	subject4 := NewCertSubjectDNAttributeCondition([]string{"2.5.4.4"})
	none4.AddChild(subject4)

	none10 := NewCompositeConditionWithMatchingCriteriaIndicator(enumerations.Assert_NONE)
	// present in cert (organizationName)
	subject10 := NewCertSubjectDNAttributeCondition([]string{conditionTestOIDOrganization})
	none10.AddChild(subject10)

	allNones.AddChild(none42)
	allNones.AddChild(none65)
	allNones.AddChild(none4)
	allNones.AddChild(none10)

	if subject4.Check(certificate) {
		t.Error("subject4 should not match")
	}
	if !none4.Check(certificate) {
		t.Error("none4 should match")
	}
	if !subject10.Check(certificate) {
		t.Error("subject10 should match")
	}
	if none10.Check(certificate) {
		t.Error("none10 should not match")
	}
	if allNones.Check(certificate) {
		t.Error("allNones should not match")
	}

	allTL := NewCompositeConditionWithMatchingCriteriaIndicator(enumerations.Assert_ALL)
	allTL.AddChild(none42)
	allTL.AddChild(none65)
	allTL.AddChild(none4)
	allTL.AddChild(subject10)
	if !allTL.Check(certificate) {
		t.Error("allTL should match")
	}
}

// TestCompositeCondition_Oracle ports CompositeConditionTest's five cases.
func TestCompositeCondition_Oracle(t *testing.T) {
	certificate := conditionTestCertificate(t, belgiumOCSPResponderCertificate)

	t.Run("default", func(t *testing.T) {
		condition := NewCompositeCondition()
		if got := condition.MatchingCriteriaIndicator(); got != enumerations.Assert_ALL {
			t.Errorf("default indicator = %q, want ALL", got)
		}
		condition.AddChild(NewCertSubjectDNAttributeCondition([]string{conditionTestOIDCountry}))
		if !condition.Check(certificate) {
			t.Error("should match")
		}
		condition.AddChild(NewCertSubjectDNAttributeCondition([]string{conditionTestOIDEmailAddress}))
		if condition.Check(certificate) {
			t.Error("should not match")
		}
	})

	t.Run("all", func(t *testing.T) {
		condition := NewCompositeConditionWithMatchingCriteriaIndicator(enumerations.Assert_ALL)
		condition.AddChild(NewCertSubjectDNAttributeCondition([]string{conditionTestOIDCountry}))
		if !condition.Check(certificate) {
			t.Error("should match")
		}
		condition.AddChild(NewCertSubjectDNAttributeCondition([]string{conditionTestOIDEmailAddress}))
		if condition.Check(certificate) {
			t.Error("should not match")
		}
	})

	t.Run("atLeastOne", func(t *testing.T) {
		condition := NewCompositeConditionWithMatchingCriteriaIndicator(enumerations.Assert_AT_LEAST_ONE)
		condition.AddChild(NewCertSubjectDNAttributeCondition([]string{conditionTestOIDCountry}))
		if !condition.Check(certificate) {
			t.Error("should match")
		}
		condition.AddChild(NewCertSubjectDNAttributeCondition([]string{conditionTestOIDEmailAddress}))
		if !condition.Check(certificate) {
			t.Error("should still match")
		}
	})

	t.Run("none", func(t *testing.T) {
		condition := NewCompositeConditionWithMatchingCriteriaIndicator(enumerations.Assert_NONE)
		condition.AddChild(NewCertSubjectDNAttributeCondition([]string{conditionTestOIDCountry}))
		if condition.Check(certificate) {
			t.Error("should not match")
		}
		condition.AddChild(NewCertSubjectDNAttributeCondition([]string{conditionTestOIDEmailAddress}))
		if condition.Check(certificate) {
			t.Error("should still not match")
		}
	})

	t.Run("multiComposites", func(t *testing.T) {
		condition := NewCompositeConditionWithMatchingCriteriaIndicator(enumerations.Assert_ALL)
		condition.AddChild(NewCertSubjectDNAttributeCondition([]string{conditionTestOIDCountry}))

		subCondition := NewCompositeConditionWithMatchingCriteriaIndicator(enumerations.Assert_ALL)
		subCondition.AddChild(NewExtendedKeyUsageCondition([]string{conditionTestOIDOCSPSigning}))

		subSubCondition := NewCompositeConditionWithMatchingCriteriaIndicator(enumerations.Assert_NONE)
		subSubCondition.AddChild(NewExtendedKeyUsageCondition([]string{"1.3.1"}))

		subCondition.AddChild(subSubCondition)
		condition.AddChild(subCondition)

		if !condition.Check(certificate) {
			t.Error("should match")
		}

		// The nesting is reflected in the rendering, one tab deeper per composite level.
		want := "CriteriaListCondition: ALL\n" +
			"\tCertSubjectDNAttributeCondition: [2.5.4.6]\n" +
			"\tCriteriaListCondition: ALL\n" +
			"\t\tExtendedKeyUsageCondition: [1.3.6.1.5.5.7.3.9]\n" +
			"\t\tCriteriaListCondition: NONE\n" +
			"\t\t\tExtendedKeyUsageCondition: [1.3.1]\n"
		if got := condition.String(); got != want {
			t.Errorf("String() =\n%q\nwant\n%q", got, want)
		}
	})
}

// TestExtendedKeyUsageCondition_Oracle ports ExtendedKeyUsageConditionTest#test.
func TestExtendedKeyUsageCondition_Oracle(t *testing.T) {
	certificateOCSP := conditionTestCertificate(t, belgiumOCSPResponderCertificate)

	condition := NewExtendedKeyUsageCondition([]string{"1.2.3"})
	if condition.Check(certificateOCSP) {
		t.Error("EKU(1.2.3) should not match")
	}

	condition = NewExtendedKeyUsageCondition([]string{conditionTestOIDOCSPSigning})
	if !condition.Check(certificateOCSP) {
		t.Error("EKU(OCSPSigning) should match")
	}

	condition = NewExtendedKeyUsageCondition([]string{conditionTestOIDOCSPSigning, "1.2.3"})
	if condition.Check(certificateOCSP) {
		t.Error("EKU(OCSPSigning, 1.2.3) should not match")
	}
}

// TestKeyUsageCondition_Oracle ports KeyUsageConditionTest#test.
func TestKeyUsageCondition_Oracle(t *testing.T) {
	certificate := conditionTestCertificate(t, estonianESTEIDCertificate)

	condition := NewKeyUsageCondition(enumerations.KeyUsageBit_DIGITAL_SIGNATURE, true)
	if condition.Check(certificate) {
		t.Error("KeyUsage(DIGITAL_SIGNATURE=true) should not match")
	}
	if got, want := condition.String(), "KeyUsageCondition: DIGITAL_SIGNATURE=true\n"; got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}

	condition2 := NewKeyUsageCondition(enumerations.KeyUsageBit_NON_REPUDIATION, true)
	if !condition2.Check(certificate) {
		t.Error("KeyUsage(NON_REPUDIATION=true) should match")
	}

	// The String overload of the Java constructor routes through KeyUsageBit.valueOf.
	fromName, err := NewKeyUsageConditionFromName("NON_REPUDIATION", true)
	if err != nil {
		t.Fatalf("NewKeyUsageConditionFromName: %v", err)
	}
	if !fromName.Equals(condition2) {
		t.Error("the (String, boolean) constructor should build the same condition")
	}
	if _, err := NewKeyUsageConditionFromName("NOT_A_KEY_USAGE", true); err == nil {
		t.Error("an unknown key usage name should be refused")
	}
}

// TestPolicyIdCondition_Oracle ports PolicyIdConditionTest#test1.
func TestPolicyIdCondition_Oracle(t *testing.T) {
	certificate := conditionTestCertificate(t, estonianESTEIDCertificate)

	condition := NewPolicyIdCondition("1.2.3")
	if condition.Check(certificate) {
		t.Error("PolicyId(1.2.3) should not match")
	}
	if got, want := condition.String(), "PolicyIdCondition: 1.2.3\n"; got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}

	condition2 := NewPolicyIdCondition("1.3.6.1.4.1.10015.1.1.1.1")
	if !condition2.Check(certificate) {
		t.Error("PolicyId(1.3.6.1.4.1.10015.1.1.1.1) should match")
	}
}

// TestQCStatementCondition_Oracle ports the seven QcStatementConditionTest cases.
func TestQCStatementCondition_Oracle(t *testing.T) {
	uaESign := conditionTestCertificate(t, uaESignCertificate)
	eSeal := conditionTestCertificate(t, eSealCertificate)
	otherOid := conditionTestCertificate(t, otherOidCertificate)

	cases := []struct {
		name                            string
		oid, qcType, legislation        string
		wantUAESign, wantESeal, wantOID bool
	}{
		{"qcCompliance", conditionTestOIDQcCompliance, "", "", true, true, true},
		{"eSignQcType", conditionTestOIDQcType, "0.4.0.1862.1.6.1", "", true, false, false},
		{"eSignQcTypeAndQcCClegislation", conditionTestOIDQcType, "0.4.0.1862.1.6.1", "UA", true, false, false},
		{"eSealQcType", conditionTestOIDQcType, "0.4.0.1862.1.6.2", "", false, true, false},
		{"eSealQcTypeAndQcCClegislation", conditionTestOIDQcType, "0.4.0.1862.1.6.2", "UA", false, false, false},
		{"qcQscd", conditionTestOIDQcSSCD, "", "", false, true, true},
		{"qcPdsOid", conditionTestOIDQcPds, "", "", false, false, true},
		{"otherOid", "1.2.804.2.1.1.1.2.1", "", "", false, true, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			condition := NewQCStatementCondition(c.oid, c.qcType, c.legislation)
			if got := condition.Check(uaESign); got != c.wantUAESign {
				t.Errorf("check(uaESignCertificate) = %v, want %v", got, c.wantUAESign)
			}
			if got := condition.Check(eSeal); got != c.wantESeal {
				t.Errorf("check(eSealCertificate) = %v, want %v", got, c.wantESeal)
			}
			if got := condition.Check(otherOid); got != c.wantOID {
				t.Errorf("check(otherOidCertificate) = %v, want %v", got, c.wantOID)
			}
		})
	}
}

// TestConditionEquality covers the equals(Object) ports, which the MRA equivalence records rely
// on to compare conditions by value rather than by reference.
func TestConditionEquality(t *testing.T) {
	if !NewPolicyIdCondition("1.2.3").Equals(NewPolicyIdCondition("1.2.3")) {
		t.Error("equal PolicyIdConditions should compare equal")
	}
	if NewPolicyIdCondition("1.2.3").Equals(NewPolicyIdCondition("1.2.4")) {
		t.Error("different PolicyIdConditions should not compare equal")
	}
	if NewPolicyIdCondition("1.2.3").Equals(NewCertSubjectDNAttributeCondition([]string{"1.2.3"})) {
		t.Error("conditions of different types should not compare equal")
	}
	if !NewQCStatementCondition("a", "b", "c").Equals(NewQCStatementCondition("a", "b", "c")) {
		t.Error("equal QCStatementConditions should compare equal")
	}
	if !NewExtendedKeyUsageCondition([]string{"1.2"}).Equals(NewExtendedKeyUsageCondition([]string{"1.2"})) {
		t.Error("equal ExtendedKeyUsageConditions should compare equal")
	}
	// Objects.equals(null, emptyList) is false in Java; the port preserves that.
	if NewExtendedKeyUsageCondition(nil).Equals(NewExtendedKeyUsageCondition([]string{})) {
		t.Error("a null OID list must not equal an empty one")
	}

	left := NewCompositeConditionWithMatchingCriteriaIndicator(enumerations.Assert_NONE)
	left.AddChild(NewPolicyIdCondition("1.2.3"))
	right := NewCompositeConditionWithMatchingCriteriaIndicator(enumerations.Assert_NONE)
	right.AddChild(NewPolicyIdCondition("1.2.3"))
	if !left.Equals(right) {
		t.Error("composites with equal children should compare equal (by value, not by reference)")
	}
	other := NewCompositeConditionWithMatchingCriteriaIndicator(enumerations.Assert_ALL)
	other.AddChild(NewPolicyIdCondition("1.2.3"))
	if left.Equals(other) {
		t.Error("composites with different indicators should not compare equal")
	}
}

// TestConditionAccessorsNeverNil covers the "possibly empty; never null" contract the final
// getters carry upstream.
func TestConditionAccessorsNeverNil(t *testing.T) {
	if got := NewCertSubjectDNAttributeCondition(nil).AttributeOids(); got == nil || len(got) != 0 {
		t.Errorf("AttributeOids() = %v, want an empty non-nil slice", got)
	}
	if got := NewExtendedKeyUsageCondition(nil).KeyPurposeIds(); got == nil || len(got) != 0 {
		t.Errorf("KeyPurposeIds() = %v, want an empty non-nil slice", got)
	}
	if got := NewCompositeCondition().Children(); got == nil || len(got) != 0 {
		t.Errorf("Children() = %v, want an empty non-nil slice", got)
	}
	// A null OID list renders as Java's literal "null".
	if got, want := NewCertSubjectDNAttributeCondition(nil).String(), "CertSubjectDNAttributeCondition: null\n"; got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}
