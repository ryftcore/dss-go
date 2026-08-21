// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/claim/DeviceKeyClaimWrapper.java (DSS 6.5.RC1).
package diagnostic

import "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"

// DeviceKeyClaimWrapper provides user-friendly access to the information present within a claim
// representing a wallet holder's key.
type DeviceKeyClaimWrapper struct {
	ClaimWrapper

	// wrapped is the concrete XmlDeviceKeyClaim, shadowing the promoted (synthetic) field of
	// the embedded ClaimWrapper; see the covariant-getWrapped note in claim_wrapper.go.
	wrapped *jaxb.XmlDeviceKeyClaim
}

// NewDeviceKeyClaimWrapper is the default constructor. Port of DeviceKeyClaimWrapper(XmlDeviceKeyClaim).
func NewDeviceKeyClaimWrapper(wrapped *jaxb.XmlDeviceKeyClaim) *DeviceKeyClaimWrapper {
	return NewDeviceKeyClaimWrapperWithParent(wrapped, nil)
}

// NewDeviceKeyClaimWrapperWithParent is the constructor with a parent provided. Port of
// DeviceKeyClaimWrapper(XmlDeviceKeyClaim, ClaimWrapper).
func NewDeviceKeyClaimWrapperWithParent(wrapped *jaxb.XmlDeviceKeyClaim, parent *ClaimWrapper) *DeviceKeyClaimWrapper {
	return &DeviceKeyClaimWrapper{
		ClaimWrapper: *NewClaimWrapperWithParent(claimBase(wrapped.XmlClaimContent, wrapped.XmlClaimAttrs), parent),
		wrapped:      wrapped,
	}
}

// PublicKey gets the public key provided within the claim. Port of getPublicKey().
func (w *DeviceKeyClaimWrapper) PublicKey() []byte {
	if w.wrapped.PublicKey != nil {
		return []byte(*w.wrapped.PublicKey)
	}
	return nil
}

// Certificates gets a list of certificate tokens. Port of getCertificates().
func (w *DeviceKeyClaimWrapper) Certificates() []*CertificateWrapper {
	x509Certificates := w.wrapped.X509Certificate
	if len(x509Certificates) == 0 {
		return nil
	}
	result := make([]*CertificateWrapper, 0, len(x509Certificates))
	for _, x := range x509Certificates {
		result = append(result, NewCertificateWrapper(x.Certificate))
	}
	return result
}

// CertificateDigests gets a list of certificate digests. Port of getCertificateDigests().
func (w *DeviceKeyClaimWrapper) CertificateDigests() []*jaxb.XmlDigestAlgoAndValue {
	return w.wrapped.DigestAlgoAndValue
}

// KIDs gets a list of certificate key identifiers. Port of getKIDs().
func (w *DeviceKeyClaimWrapper) KIDs() []string {
	return w.wrapped.KID
}

// X509URLs gets a list of certificate access URLs. Port of getX509URLs().
func (w *DeviceKeyClaimWrapper) X509URLs() []string {
	return w.wrapped.X509Url
}

// AuthorizedNamespaces gets a list namespaces the key is authorized to sign or MAC. Port of
// getAuthorizedNamespaces().
func (w *DeviceKeyClaimWrapper) AuthorizedNamespaces() []string {
	if w.wrapped.KeyAuthorizations != nil && w.wrapped.KeyAuthorizations.AuthorizedNamespace != nil {
		return w.wrapped.KeyAuthorizations.AuthorizedNamespace
	}
	return nil
}

// AuthorizedDataElements gets a map of namespaces and corresponding data elements the key is
// authorized to sign or MAC. Port of getAuthorizedDataElements().
func (w *DeviceKeyClaimWrapper) AuthorizedDataElements() map[string][]string {
	if w.wrapped.KeyAuthorizations == nil || w.wrapped.KeyAuthorizations.AuthorizedDataElements == nil {
		return map[string][]string{}
	}
	result := make(map[string][]string, len(w.wrapped.KeyAuthorizations.AuthorizedDataElements))
	for _, xmlAuthorizedDataElements := range w.wrapped.KeyAuthorizations.AuthorizedDataElements {
		var namespace string
		if xmlAuthorizedDataElements.Namespace != nil {
			namespace = *xmlAuthorizedDataElements.Namespace
		}
		result[namespace] = xmlAuthorizedDataElements.DataElement
	}
	return result
}

// Wrapped is the covariant override of ClaimWrapper.Wrapped(). Port of the covariant getWrapped().
func (w *DeviceKeyClaimWrapper) Wrapped() *jaxb.XmlDeviceKeyClaim { return w.wrapped }

// AsClaim views the wrapper as its ClaimWrapper base type; isList/getList/isMap/getMap are not
// overridden here. See the package note in claim_wrapper.go.
func (w *DeviceKeyClaimWrapper) AsClaim() *ClaimWrapper { return &w.ClaimWrapper }
