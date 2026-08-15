// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/aov/TimestampAlgorithmObsolescenceValidation.java (DSS 6.5.RC1).
//
// Performs cryptographic validation of a timestamp token.
package aov

import (
	"time"

	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
)

// TimestampAlgorithmObsolescenceValidation performs cryptographic validation
// of a timestamp token.
type TimestampAlgorithmObsolescenceValidation struct {
	SignatureAlgorithmObsolescenceValidation[*diagnostic.TimestampWrapper]
}

// NewTimestampAlgorithmObsolescenceValidation is the default constructor.
func NewTimestampAlgorithmObsolescenceValidation(i18nProvider *i18n.I18nProvider, token *diagnostic.TimestampWrapper,
	validationDate time.Time, validationPolicy policy.ValidationPolicy) *TimestampAlgorithmObsolescenceValidation {
	c := &TimestampAlgorithmObsolescenceValidation{}
	c.InitAlgorithmObsolescenceValidation(i18nProvider, token, enumerations.Context_TIMESTAMP, validationDate, validationPolicy, c)
	c.InitChainBase(c)
	return c
}
