// Java-oracle tests for tl_validator_task.go and tl_validation_result.go.
//
// Every fixture, certificate and expected Indication/SubIndication below is transcribed verbatim
// from the upstream JUnit suite (dss-tsl-validation 6.5.RC1
// eu.europa.esig.dss.tsl.validation.TLValidatorTaskTest). testdata/eu-lotl.xml,
// testdata/eu-lotl-broken-sig.xml and testdata/eu-lotl-no-sig.xml are the upstream
// src/test/resources files of the same names, copied unchanged.
//
// This is harness contract (B) - TL/LOTL signature-validation verdicts - exercised end to end
// through the frozen XAdES validator.
//
// Test vectors are ported, not the JUnit code, per PORTING.md.
package tsl

import (
	"strings"
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/corpustest"
	"github.com/ryftcore/dss-go/dss/model"
	dsspolicy "github.com/ryftcore/dss-go/dss/policy"
	"github.com/ryftcore/dss-go/dss/spi"
	validationpolicy "github.com/ryftcore/dss-go/dss/validation/policy"
)

// INTEGRATION FLAG (see this batch's porter notes): Java discovers the ETSI
// ValidationPolicyFactory through ServiceLoader against dss-policy-jaxb, which reaches
// dss-tsl-validation's classpath transitively via dss-validation. The Go port of
// ValidationPolicyLoader replaces ServiceLoader with an explicit registry that NOTHING in the
// production tree populates yet (see validation/policy/validation_policy_loader.go's header),
// so TLValidatorTask#trustedListValidationPolicy fails with upstream's
// "no suitable ValidationPolicyFactory has been found" message until a composing application
// registers one. This init does what such an application would do; it mirrors the same init in
// validation/policy/validation_policy_loader_test.go.
func init() {
	validationpolicy.RegisterValidationPolicyFactory(dsspolicy.NewEtsiValidationPolicyFactory())
}

const (
	tlValidatorTaskSigningCertificate = "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"
	tlValidatorTaskWrongCertificate   = "MIIFvjCCA6agAwIBAgIQALwvYx2O1YN6UxQOi3Bx3jANBgkqhkiG9w0BAQUFADBbMQswCQYDVQQGEwJFUzEoMCYGA1UECgwfRElSRUNDSU9OIEdFTkVSQUwgREUgTEEgUE9MSUNJQTEMMAoGA1UECwwDQ05QMRQwEgYDVQQDDAtBQyBSQUlaIERHUDAeFw0wNzAxMjUxMjA1MDhaFw0zNzAxMjUxMjA1MDhaMFsxCzAJBgNVBAYTAkVTMSgwJgYDVQQKDB9ESVJFQ0NJT04gR0VORVJBTCBERSBMQSBQT0xJQ0lBMQwwCgYDVQQLDANDTlAxFDASBgNVBAMMC0FDIFJBSVogREdQMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAgBD1t16zMJxvoxuIDlyt6pfgzPmmfJMFvPyoj0AOxjyxu6f77K/thV/pMatQqjGae3Yj83upv7YFygq/jU02EeEIeQQEf+QJ+B+LX+oGLPbU5g8/W1eFcnXC4Jg2ipP7L2qcEfA180AsT1UqmHTc7kRI3N6yJZZiHkM4hpjf3vgsCxUQtXw+XAZYtaRbjFO69tTSdbpbXN4fvOQwHNlenF1GMxsih7tgGUwRlY2EVfh7EGYvXt2mtpHiEIeSp1s2WBxzgiWU1IufiDo18olZj859oHkNBD0sx6LVPPun/sINuM1M6aBRwc725cMgZmIyNDOHZkqExL8DNUiTzXYzqr7R/X+kn59RYLwIEmfRQLkKxyYlZeFbuOI5n7Uz3vKANcTbUuCymA0+ZA9ESlrz8kA6fHV0+fMePUBYnociJO5fFX/jxtScOqrQt+K+gGm4TubalBoL7ECGzs3CmKtnuyOH+KFO/8q71Fxhn3WqlKgO7dBUhp0I/7dr4R2bF4ry1NnqZWObCuBfKqyL80Dx+6zaGsTo7UBLNdcA4sXArJoAMUqHb/77rqu45dWJIhQA5V3qolwowwuTdZwC1ec2AWwA6gMf2uchNJsPWWmQrkXvkhu2rI756cKwgR7y22517q/B9MNx7InsZbMbOWUwQuei3UcoIgCFs2TWCbhxHNkCAwEAAaN+MHwwDwYDVR0TAQH/BAUwAwEB/zAOBgNVHQ8BAf8EBAMCAQYwHQYDVR0OBBYEFA6cduGiLokzQfLjPmxbFkW9vYaOMDoGA1UdIAQzMDEwLwYEVR0gADAnMCUGCCsGAQUFBwIBFhlodHRwOi8vd3d3LnBvbGljaWEuZXMvZHBjMA0GCSqGSIb3DQEBBQUAA4ICAQBslvw3pwCj21vCctyL7YOrmfINjJFp4TNFfNnDwSsuonqOjwppXCEFJ6MkOeCUOy9vXziNoYtoDd/tXAn++9975d7PB9vXnu7ErHRx+e74obKpqfBoVv9fwPp0bObO3YbTq9EGPLM8mbcUEivPlL2mQ7tk78z2p8gpytcCZRc08Jd5m+AeYPrHUDeF6ZIlnH7SIrtP3Bp8zwnNIFbNtkyrCyWtN8Ajo3RXqecM/bs+YgGzjVbDToQUBkBCuoG3XU+QYSQ79yZsvjTCsFKBYnXXijiGZSokx33iauY0PIyaNu/ulMloSNUwWZ5WBPqJXWlkZ+deApxZLXJLFMSTjFeFdpZUgOC1wrRkxXidWQwr4566fYWhYH0w+hwK9gD6NEsMA3D7NOPCTCOx9Qst5848RsJVJ4F+ZFmT4iyTYLyglkNkeB+tSXVyC9Lg+Tvay85VyeZMSZ3PpGmpNzaQxVZl9XCfs8R6Ew4pG91eOA0BjsI1ZHY7H9e5Pomup/jTA6JwlCYooEiBM31Gdwe/3oUFNzB+NvOWdwb+ZG6va70j98EdipGWoLvjv/oJlFN2q1Nrt/u7whKp+VsVOjuZMrSpw9C+Ec4yiLha5RRiXnHX1cqwT694KIDQZIgqQChQDeDqrvCphtdHdxFQ5NBzt2HKhaSh8ggDdOdpH451rB45Jg=="
	tlValidatorTaskCACertificate      = "MIIGjzCCBHegAwIBAgIQXh8tIPyKnFy/9o9+/M2RrTANBgkqhkiG9w0BAQsFADAoMQswCQYDVQQGEwJCRTEZMBcGA1UEAxMQQmVsZ2l1bSBSb290IENBNDAeFw0xNzExMTcxMDAwMDBaFw0yOTA3MTcxMDAwMDBaMGQxCzAJBgNVBAYTAkJFMREwDwYDVQQHEwhCcnVzc2VsczEcMBoGA1UEChMTQ2VydGlwb3N0IE4uVi4vUy5BLjETMBEGA1UEAxMKQ2l0aXplbiBDQTEPMA0GA1UEBRMGMjAxODAzMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEA3qHulx/7LvzTQKfLxLlHEcKco2n1Xfy47Z8Kqjt/ECvVFQxmFNa9X78O/EBttg2dLJgxKGhHOo0UtJAszF9872Z/SDQxV+DnfpLE0rew7Y65MDzCT6F0d2Yirett8TVIUVzuHDGmPeVixTvqrdK6xjsmJ98Hfe0kxfFEw8zh3gky7E7d4Sn9IF0sTbOlgFDmnGU/d5XxBdtlMYVTT2I22kYz5Q97Rm2Hx+GTf7fvNsg/eUpAlvNBIz9s8q+lK/MjN+gt8dA6Ec0S9bY/tocUwSEIf84vnPoS0u7HGt1gERvYOr5zoCC2TEcP++nHR+I8asriJA7/j/3Jza6qz/ioV7nCEHQ3jt9+jrwG1AnxQni6RHpRCL9KYx2Riwk4aAsjstuI2b501LONbEC8C1wBqi53sgbHS3xl5L8EA1iy78rEAW5bimF3HWYSa5U5deiUCcMZVKY3ts6lXGqiPBUqQG/Ug5gazIkDsijJ4On+oUoG5KFQ7UhDHj7AtP77iqlBtRB97VI7amRyOLoYbENFPqZOm/EHSRpo/C0dSBQDGg83iSge67Kn14Al5Zw7BefrSN6rUC6BsioxukDi+7WOztqiYme3UVpQl+WOeE7Z7g/ytcyZ+4Sl/TRqOa9YCS1S925kxq7+QXW2TozgqQOac06V+4H9Oi2c8okWqekJDWMCAwEAAaOCAXcwggFzMA4GA1UdDwEB/wQEAwIBBjASBgNVHRMBAf8ECDAGAQH/AgEAMHAGCCsGAQUFBwEBBGQwYjA2BggrBgEFBQcwAoYqaHR0cDovL2NlcnRzLmVpZC5iZWxnaXVtLmJlL2JlbGdpdW1yczQuY3J0MCgGCCsGAQUFBzABhhxodHRwOi8vb2NzcC5laWQuYmVsZ2l1bS5iZS8yMB8GA1UdIwQYMBaAFGfo8U5Ps7XzB28InAyD2XrZW+dJMEMGA1UdIAQ8MDowOAYGYDgMAQECMC4wLAYIKwYBBQUHAgEWIGh0dHA6Ly9yZXBvc2l0b3J5LmVpZC5iZWxnaXVtLmJlMDcGA1UdHwQwMC4wLKAqoCiGJmh0dHA6Ly9jcmwuZWlkLmJlbGdpdW0uYmUvYmVsZ2l1bTQuY3JsMB0GA1UdDgQWBBTZNCE+OkIlbeG7v75HTPgCreMvVDAdBgNVHSUEFjAUBggrBgEFBQcDAgYIKwYBBQUHAwQwDQYJKoZIhvcNAQELBQADggIBADiZCmd+Gol1t9GxK+vf9eii1oiXDYqIOJAgCiXR9S+d5CN0KMwXGczlVySKEjCVyjkWPXAqj/pb1gbI9GDNNlD3lNadnwNwaIiJYVcNX9eip2V70JEzvJ6OQoQUjPtbsKgeYo/S8IjJ6YsSFYf+Bt6oWBa1getGAJNvaCHDM9mnmKyQNW6cu/qWPnRBAY1wwX7b3+OtP9wkOrqflUlCKmJ3SzxaPM3B/nft3QsoT6pvKekZxo+Bq1Ae2sAqC1IswqEunXVKKHllHz3HIKNLHVPPjKaMMVAGJj73A5BKmXRHwskccN8FxioY6Xav+putUVg/IialA7UBiQQZF3udm+vTBiIXQhwjNnbmI5k3naJH54R0tlxWMoV8gX6s7zZR4GKX+elTTBdvI+FRCsAOJ0ApVckDRcq0g+uAPwYQ9YvaSeBHr9DoYT/4MwkiOwgtdXe83EPQioNy4JntDxLh/9VbVJqDXJghmxDoOOtf5CwmkYMJCiET6+G61tBIRclo4CUy7GW56PYgS98V8ZEro33llB74ew9+EH6nw6EfNcf216PfVC6dDN905MWxWIPv4lYitLK2168nkIj2CXcGACt58bszr+kM6GEJv1qPHQ2zoR+nY7fJjNSB3JxBkKJ3HTlV9SzJ0UULUoLfW6iMN2r/pHAO/Nq32ajAFQyySjpI"
)

// tlValidatorTaskCertificateSource ports the test's private getCertificateSource(List).
func tlValidatorTaskCertificateSource(t *testing.T, base64Certificates ...string) spi.CertificateSource {
	t.Helper()
	source := spi.NewCommonCertificateSource()
	for _, base64Certificate := range base64Certificates {
		certificate, err := spi.DSSUtilsLoadCertificateFromBase64EncodedString(base64Certificate)
		if err != nil {
			t.Fatalf("unable to load the certificate: %v", err)
		}
		source.AddCertificate(certificate)
	}
	return &source
}

func tlValidatorTaskDocument(t *testing.T, name string) model.DSSDocument {
	t.Helper()
	document, err := model.NewFileDocument(corpustest.Path(t, name))
	if err != nil {
		t.Fatalf("unable to load %s: %v", name, err)
	}
	return document
}

// TestTLValidatorTask_Oracle ports the eight TLValidatorTaskTest cases.
func TestTLValidatorTask_Oracle(t *testing.T) {
	t.Run("correctCert", func(t *testing.T) {
		trustedList := tlValidatorTaskDocument(t, "eu-lotl.xml")
		source := tlValidatorTaskCertificateSource(t, tlValidatorTaskSigningCertificate)
		result, err := NewTLValidatorTask(trustedList, source).Get()
		if err != nil {
			t.Fatalf("Get: %v", err)
		}
		if result.Indication() == "" {
			t.Error("an Indication is expected")
		}
		if result.SigningTime().IsZero() {
			t.Error("a signing time is expected")
		}
		if result.SigningCertificate() == nil {
			t.Fatal("a signing certificate is expected")
		}
		if !result.SigningCertificate().Equals(source.Certificates()[0]) {
			t.Error("the signing certificate should be the announced one")
		}
		// The announced signer is the actual one, so the chain is found.
		if got, want := result.Indication(), enumerations.IndicationTotalPassed; got != want {
			t.Errorf("Indication = %q, want %q", got, want)
		}
	})

	t.Run("wrongCert", func(t *testing.T) {
		trustedList := tlValidatorTaskDocument(t, "eu-lotl.xml")
		source := tlValidatorTaskCertificateSource(t, tlValidatorTaskWrongCertificate)
		result, err := NewTLValidatorTask(trustedList, source).Get()
		if err != nil {
			t.Fatalf("Get: %v", err)
		}
		if got, want := result.Indication(), enumerations.IndicationIndeterminate; got != want {
			t.Errorf("Indication = %q, want %q", got, want)
		}
		if got, want := result.SubIndication(), enumerations.SubIndicationNoCertificateChainFound; got != want {
			t.Errorf("SubIndication = %q, want %q", got, want)
		}
		if result.SigningTime().IsZero() {
			t.Error("a signing time is expected")
		}
		if result.SigningCertificate() == nil {
			t.Fatal("a signing certificate is expected")
		}
		if result.SigningCertificate().Equals(source.Certificates()[0]) {
			t.Error("the recovered signing certificate must not be the wrong announced one")
		}
	})

	t.Run("noCert", func(t *testing.T) {
		trustedList := tlValidatorTaskDocument(t, "eu-lotl.xml")
		source := tlValidatorTaskCertificateSource(t)
		result, err := NewTLValidatorTask(trustedList, source).Get()
		if err != nil {
			t.Fatalf("Get: %v", err)
		}
		if got, want := result.Indication(), enumerations.IndicationIndeterminate; got != want {
			t.Errorf("Indication = %q, want %q", got, want)
		}
		if got, want := result.SubIndication(), enumerations.SubIndicationNoCertificateChainFound; got != want {
			t.Errorf("SubIndication = %q, want %q", got, want)
		}
		if result.SigningTime().IsZero() {
			t.Error("a signing time is expected")
		}
		if result.SigningCertificate() == nil {
			t.Error("a signing certificate is expected")
		}
		if got := result.PotentialSigners(); len(got) != 0 {
			t.Errorf("PotentialSigners() = %v, want none", got)
		}
	})

	t.Run("caCert", func(t *testing.T) {
		trustedList := tlValidatorTaskDocument(t, "eu-lotl.xml")
		source := tlValidatorTaskCertificateSource(t, tlValidatorTaskCACertificate)
		result, err := NewTLValidatorTask(trustedList, source).Get()
		if err != nil {
			t.Fatalf("Get: %v", err)
		}
		if got, want := result.Indication(), enumerations.IndicationIndeterminate; got != want {
			t.Errorf("Indication = %q, want %q", got, want)
		}
		if got, want := result.SubIndication(), enumerations.SubIndicationNoCertificateChainFound; got != want {
			t.Errorf("SubIndication = %q, want %q", got, want)
		}
		if result.SigningCertificate() == nil {
			t.Fatal("a signing certificate is expected")
		}
		if result.SigningCertificate().Equals(source.Certificates()[0]) {
			t.Error("the recovered signing certificate must not be the announced CA")
		}
	})

	t.Run("brokenTL", func(t *testing.T) {
		trustedList := tlValidatorTaskDocument(t, "eu-lotl-broken-sig.xml")
		source := tlValidatorTaskCertificateSource(t, tlValidatorTaskSigningCertificate)
		result, err := NewTLValidatorTask(trustedList, source).Get()
		if err != nil {
			t.Fatalf("Get: %v", err)
		}
		if got, want := result.Indication(), enumerations.IndicationTotalFailed; got != want {
			t.Errorf("Indication = %q, want %q", got, want)
		}
		if got, want := result.SubIndication(), enumerations.SubIndicationHashFailure; got != want {
			t.Errorf("SubIndication = %q, want %q", got, want)
		}
		if result.SigningTime().IsZero() {
			t.Error("a signing time is expected")
		}
		if result.SigningCertificate() == nil {
			t.Fatal("a signing certificate is expected")
		}
		if !result.SigningCertificate().Equals(source.Certificates()[0]) {
			t.Error("the signing certificate should be the announced one")
		}
	})

	t.Run("noSig", func(t *testing.T) {
		trustedList := tlValidatorTaskDocument(t, "eu-lotl-no-sig.xml")
		source := tlValidatorTaskCertificateSource(t, tlValidatorTaskSigningCertificate)
		_, err := NewTLValidatorTask(trustedList, source).Get()
		if err == nil {
			t.Fatal("an unsigned Trusted List should be refused")
		}
		if got, want := err.Error(), "Number of signatures must be equal to 1 (currently : 0)"; got != want {
			t.Errorf("error = %q, want %q", got, want)
		}
	})

	t.Run("notXML", func(t *testing.T) {
		trustedList := model.NewInMemoryDocument([]byte{})
		source := tlValidatorTaskCertificateSource(t, tlValidatorTaskSigningCertificate)
		task := NewTLValidatorTask(trustedList, source)
		failed := false
		func() {
			defer func() {
				if recover() != nil {
					failed = true
				}
			}()
			if _, err := task.Get(); err != nil {
				failed = true
			}
		}()
		if !failed {
			t.Error("a non-XML document should be refused")
		}
	})

	t.Run("nullArguments", func(t *testing.T) {
		assertPanicMessage := func(name, want string, call func()) {
			t.Helper()
			defer func() {
				recovered := recover()
				if recovered == nil {
					t.Errorf("%s should have panicked", name)
					return
				}
				if got, ok := recovered.(string); !ok || !strings.Contains(got, want) {
					t.Errorf("%s panic = %v, want %q", name, recovered, want)
				}
			}()
			call()
		}
		trustedList := tlValidatorTaskDocument(t, "eu-lotl.xml")
		assertPanicMessage("nil certificate source", "The certificate source is null", func() {
			NewTLValidatorTask(trustedList, nil)
		})
		assertPanicMessage("nil document", "The document is null", func() {
			source := spi.NewCommonCertificateSource()
			NewTLValidatorTask(nil, &source)
		})
	})
}
