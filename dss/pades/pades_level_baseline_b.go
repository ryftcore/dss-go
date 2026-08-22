// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/signature/PAdESLevelBaselineB.java (DSS 6.5.RC1).
//
// Java extends cades.CAdESLevelBaselineB and overrides four of its protected add*() hooks with
// empty bodies plus addSignedAttributes (super + message-imprint). Go has no method overriding
// across embedding, and cades.LevelBaselineB.AddSignedAttributes calls its own methods on
// the base receiver, so the two entry points the CMS builder uses - SignedAttributes and
// AddSignedAttributes - are re-implemented here with the four PAdES no-ops folded in. The
// remaining eight attribute builders are the embedded base's, untouched, so the produced DER
// stays identical to CAdES for them.
//
// The empty overrides are still declared (as identity functions) so that a caller holding a
// *PAdESLevelBaselineB observes upstream's behaviour on a direct call too.
package pades

import (
	"github.com/ryftcore/dss-go/dss/cades"
	"github.com/ryftcore/dss-go/dss/internal/asn1ber"
	"github.com/ryftcore/dss-go/dss/internal/cmscore"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/utils"
)

// LevelBaselineB holds the PAdES Baseline B signature profile.
type LevelBaselineB struct {
	cades.LevelBaselineB

	// messageDigest is the message digest computed on the PAdES revision.
	messageDigest model.DSSMessageDigest
}

// NewLevelBaselineB is the default constructor.
// Port of PAdESLevelBaselineB(DSSMessageDigest).
func NewLevelBaselineB(messageDigest model.DSSMessageDigest) *LevelBaselineB {
	return &LevelBaselineB{
		LevelBaselineB: *cades.NewLevelBaselineB(),
		messageDigest:  messageDigest,
	}
}

// SignedAttributes generates and returns the signed attributes table.
// Port of the inherited getSignedAttributes(CAdESSignatureParameters); re-declared because the
// embedded implementation would call the base's AddSignedAttributes rather than this type's.
func (b *LevelBaselineB) SignedAttributes(parameters *cades.SignatureParameters) (cmscore.Attributes, error) {
	if utils.IsArrayNotEmpty(parameters.SignedData()) {
		// Upstream logs "Using explicit SignedAttributes from parameter".
		return cades.UtilsAttributesFromByteArray(parameters.SignedData())
	}
	return b.AddSignedAttributes(parameters, cmscore.Attributes{})
}

// AddSignedAttributes adds the signed attributes of a PAdES Baseline B signature: the CAdES set
// minus signing-time, signer-location, content-identifier and mime-type, plus the
// message-imprint. Port of the protected #addSignedAttributes.
func (b *LevelBaselineB) AddSignedAttributes(parameters *cades.SignatureParameters,
	signedAttributes cmscore.Attributes) (cmscore.Attributes, error) {
	var err error
	// super.addSignedAttributes(parameters, signedAttributes), with this type's overrides in
	// the same order the base calls them.
	if signedAttributes, err = b.AddSigningCertificateAttribute(parameters, signedAttributes); err != nil {
		return nil, err
	}
	if signedAttributes, err = b.AddSigningTimeAttribute(parameters, signedAttributes); err != nil {
		return nil, err
	}
	if signedAttributes, err = b.AddSignerAttribute(parameters, signedAttributes); err != nil {
		return nil, err
	}
	if signedAttributes, err = b.AddSignaturePolicyId(parameters, signedAttributes); err != nil {
		return nil, err
	}
	if signedAttributes, err = b.AddContentHints(parameters, signedAttributes); err != nil {
		return nil, err
	}
	if signedAttributes, err = b.AddMimeType(parameters, signedAttributes); err != nil {
		return nil, err
	}
	if signedAttributes, err = b.AddContentIdentifier(parameters, signedAttributes); err != nil {
		return nil, err
	}
	if signedAttributes, err = b.AddCommitmentType(parameters, signedAttributes); err != nil {
		return nil, err
	}
	if signedAttributes, err = b.AddSignerLocation(parameters, signedAttributes); err != nil {
		return nil, err
	}
	if signedAttributes, err = b.AddContentTimestamps(parameters, signedAttributes); err != nil {
		return nil, err
	}
	return b.AddMessageImprint(parameters, signedAttributes)
}

// AddSigningTimeAttribute is a no-op: in PAdES the signing time is not included, per
// ETSI TS 102 778-3 V1.2.1 (2010-07), 4.5.3 signing-time Attribute.
// Port of the empty override of #addSigningTimeAttribute.
func (b *LevelBaselineB) AddSigningTimeAttribute(_ *cades.SignatureParameters,
	signedAttributes cmscore.Attributes) (cmscore.Attributes, error) {
	return signedAttributes, nil
}

// AddSignerLocation is a no-op: in PAdES the role is in the signature dictionary.
// Port of the empty override of #addSignerLocation.
func (b *LevelBaselineB) AddSignerLocation(_ *cades.SignatureParameters,
	signedAttributes cmscore.Attributes) (cmscore.Attributes, error) {
	return signedAttributes, nil
}

// AddContentIdentifier is a no-op: the attribute is prohibited in PAdES B.
// Port of the empty override of #addContentIdentifier.
func (b *LevelBaselineB) AddContentIdentifier(_ *cades.SignatureParameters,
	signedAttributes cmscore.Attributes) (cmscore.Attributes, error) {
	return signedAttributes, nil
}

// AddMimeType is a no-op: the attribute is skipped for PAdES.
// Port of the empty override of #addMimeType.
func (b *LevelBaselineB) AddMimeType(_ *cades.SignatureParameters,
	signedAttributes cmscore.Attributes) (cmscore.Attributes, error) {
	return signedAttributes, nil
}

// AddMessageImprint adds the message-digest attribute carrying the digest computed on the PDF
// revision's ByteRange, i.e. new Attribute(CMSAttributes.messageDigest,
// new DERSet(new DEROctetString(messageDigest.getValue()))).
// Port of the protected #addMessageImprint.
func (b *LevelBaselineB) AddMessageImprint(_ *cades.SignatureParameters,
	signedAttributes cmscore.Attributes) (cmscore.Attributes, error) {
	octetString := asn1ber.WriteTLV(asn1ber.TagOctetString, b.messageDigest.Value())
	return append(signedAttributes, cmscore.NewAttribute(cmscore.OIDMessageDigest, octetString)), nil
}
