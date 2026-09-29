package validation

import (
	"errors"
	"testing"

	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/x509/revocation"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/exception"
)

// panickingCRLSource is a CRL source whose RevocationToken panics with a fixed value.
type panickingCRLSource struct{ value any }

func (s *panickingCRLSource) RevocationToken(certificateToken, issuerCertificateToken *model.CertificateToken) spi.RevocationToken[revocation.CRL] {
	panic(s.value)
}

// panickingOCSPSource is the OCSP counterpart of panickingCRLSource.
type panickingOCSPSource struct{ value any }

func (s *panickingOCSPSource) RevocationToken(certificateToken, issuerCertificateToken *model.CertificateToken) spi.RevocationToken[revocation.OCSP] {
	panic(s.value)
}

// TestLoadingStrategyCatchesTheWholeDSSExceptionFamily pins the recover() of checkCRL/checkOCSP to
// Java's catch (DSSException): DSSException and every subclass (here DSSExternalResourceException
// and DSSDataLoaderMultipleException, the ones a revocation source built on a data loader raises)
// make the source count as having no token, whereas anything else keeps propagating (T22D4-SEC-001).
func TestLoadingStrategyCatchesTheWholeDSSExceptionFamily(t *testing.T) {
	caught := map[string]error{
		"DSSException":                       model.NewDSSError("plain"),
		"wrapped DSSException":               errors.Join(errors.New("context"), model.NewDSSError("wrapped")),
		"DSSExternalResourceException":       exception.NewDSSExternalResourceException("external"),
		"DSSExternalResourceException+cause": exception.NewDSSExternalResourceExceptionMessageCause("external", errors.New("io")),
		"DSSDataLoaderMultipleException": exception.NewDSSDataLoaderMultipleException(
			map[string]error{"http://a": errors.New("a"), "http://b": errors.New("b")}),
	}
	for name, value := range caught {
		t.Run(name, func(t *testing.T) {
			strategy := &RevocationDataLoadingStrategy{}
			strategy.setCrlSource(&panickingCRLSource{value: value})
			strategy.setOcspSource(&panickingOCSPSource{value: value})
			if token := strategy.checkCRL(nil, nil); token != nil {
				t.Errorf("checkCRL returned %v, want no token", token)
			}
			if token := strategy.checkOCSP(nil, nil); token != nil {
				t.Errorf("checkOCSP returned %v, want no token", token)
			}
		})
	}

	// Not a DSSException upstream (a NullPointerException, a ClassCastException, ...): propagates.
	for _, value := range []any{"a plain string", errors.New("a plain error")} {
		strategy := &RevocationDataLoadingStrategy{}
		strategy.setCrlSource(&panickingCRLSource{value: value})
		func() {
			defer func() {
				if recovered := recover(); recovered != value {
					t.Errorf("checkCRL recovered %v, want %v to propagate", recovered, value)
				}
			}()
			strategy.checkCRL(nil, nil)
			t.Errorf("checkCRL swallowed %v", value)
		}()
	}
}
