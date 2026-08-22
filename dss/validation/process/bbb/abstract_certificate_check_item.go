// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/AbstractCertificateCheckItem.java (DSS 6.5.RC1).
package bbb

import (
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// AbstractCertificateCheckItem is the abstract class to check if the given
// certificate matches one of the defined conditions. A concrete check embeds it
// instead of process.ChainItemBase and registers itself with InitChainItem the
// same way.
type AbstractCertificateCheckItem[T any] struct {
	*process.ChainItemBase[T]

	// constraint is the constraint value.
	constraint policy.CertificateApplicabilityRule
}

// NewAbstractCertificateCheckItem is the default constructor. Port of
// AbstractCertificateCheckItem(I18nProvider, T, CertificateApplicabilityRule).
func NewAbstractCertificateCheckItem[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	constraint policy.CertificateApplicabilityRule) *AbstractCertificateCheckItem[T] {
	return &AbstractCertificateCheckItem[T]{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		constraint:    constraint,
	}
}

// NewAbstractCertificateCheckItemWithId is the default constructor with Id. Port
// of AbstractCertificateCheckItem(I18nProvider, T, CertificateWrapper, CertificateApplicabilityRule).
func NewAbstractCertificateCheckItemWithId[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	certificate *diagnostic.CertificateWrapper,
	constraint policy.CertificateApplicabilityRule) *AbstractCertificateCheckItem[T] {
	return &AbstractCertificateCheckItem[T]{
		ChainItemBase: process.NewChainItemBaseWithId(i18nProvider, result, constraint, certificate.Id()),
		constraint:    constraint,
	}
}

// ProcessCertificateCheck checks the certificate, returning TRUE if the
// certificate matches the constraint. Port of
// processCertificateCheck(CertificateWrapper).
func (c *AbstractCertificateCheckItem[T]) ProcessCertificateCheck(certificate *diagnostic.CertificateWrapper) bool {
	if c.constraint == nil {
		return false
	}

	certificateExtensionsConstraint := c.constraint.CertificateExtensions()
	var expectedCertificateExtensions []string
	if certificateExtensionsConstraint != nil {
		expectedCertificateExtensions = certificateExtensionsConstraint.Values()
	}
	certificatePoliciesConstraint := c.constraint.CertificatePolicies()
	var expectedCertificatePolicies []string
	if certificatePoliciesConstraint != nil {
		expectedCertificatePolicies = certificatePoliciesConstraint.Values()
	}

	return (utils.IsCollectionNotEmpty(expectedCertificateExtensions) &&
		utils.IsCollectionNotEmpty(certificate.CertificateExtensionsOids()) &&
		process.ProcessValuesCheck(certificate.CertificateExtensionsOids(), expectedCertificateExtensions)) ||
		(utils.IsCollectionNotEmpty(expectedCertificatePolicies) &&
			utils.IsCollectionNotEmpty(certificate.CertificatePoliciesOids()) &&
			process.ProcessValuesCheck(certificate.CertificatePoliciesOids(), expectedCertificatePolicies))
}
