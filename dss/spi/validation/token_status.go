// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/status/TokenStatus.java (DSS 6.5.RC1).
//
// SCC flattening: Java spi.validation.status.TokenStatus lands in this same Go package.
package validation

import (
	"sort"

	"github.com/ryftcore/dss-go/dss/alert"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/utils"
)

// TokenStatus contains tokens concerned by an occurred event and corresponding information
// about them.
//
// Judgment call: upstream keys relatedTokenMap in a java.util.HashMap, whose iteration order
// is unspecified. RelatedTokens sorts by DSSIDAsString for determinism, mirroring the same
// documented deviation taken by alert.ObjectStatus.objectMapToString for its own map.
type TokenStatus struct {
	*alert.ObjectStatus

	// relatedTokenMap is a map between tokens concerned by the check and corresponding error
	// explanation messages.
	relatedTokenMap map[model.Token]string
}

// NewTokenStatus instantiates a TokenStatus with an empty map of related tokens. Ports the
// default constructor.
func NewTokenStatus() *TokenStatus {
	return &TokenStatus{
		ObjectStatus:    alert.NewObjectStatus(),
		relatedTokenMap: make(map[model.Token]string),
	}
}

// AddRelatedTokenAndErrorMessage adds a concerned token and information about the occurred
// event. Port of addRelatedTokenAndErrorMessage(Token, String).
func (s *TokenStatus) AddRelatedTokenAndErrorMessage(token model.Token, errorMessage string) {
	s.ObjectStatus.AddRelatedObjectIdentifierAndErrorMessage(token.DSSIDAsString(), errorMessage)
	s.relatedTokenMap[token] = errorMessage
}

// RelatedTokens returns a collection of tokens concerned by failure of the processed check.
// Port of getRelatedTokens().
func (s *TokenStatus) RelatedTokens() []model.Token {
	tokens := make([]model.Token, 0, len(s.relatedTokenMap))
	for token := range s.relatedTokenMap {
		tokens = append(tokens, token)
	}
	sort.Slice(tokens, func(i, j int) bool {
		return tokens[i].DSSIDAsString() < tokens[j].DSSIDAsString()
	})
	return tokens
}

// MessageForToken returns the error message for the given token. Port of getMessageForToken(Token).
func (s *TokenStatus) MessageForToken(token model.Token) string {
	return s.relatedTokenMap[token]
}

// IsEmpty reports whether the status has no message and no related tokens. Port of the
// isEmpty() override.
func (s *TokenStatus) IsEmpty() bool {
	return s.ObjectStatus.IsEmpty() && utils.IsMapEmpty(s.relatedTokenMap)
}

// compile-time assertion: a TokenStatus is an alert.Status.
var _ alert.Status = (*TokenStatus)(nil)
