package main

import (
	"bytes"
	"crypto/rand"
	"encoding/asn1"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/cmscore"
	"github.com/ryftcore/dss-go/dss/model"
)

// httpTSPSource is a minimal RFC 3161 time-stamping client over HTTP, for
// the -tsa flag of "sign" and "extend". It exists only in the CLI: the dss
// module itself ships no HTTP TSA client (see the dss package doc, "Network
// access") because Java DSS's own OnlineTSPSource is a thin wrapper around
// exactly this request/response pair, and reproducing that wrapper here -
// not inside the library - keeps the facade's "no network unless asked"
// promise intact while still giving the CLI a real one.
//
// It implements [dss.TSPSource] (== spi/validation.TSPSource), the same
// interface [github.com/ryftcore/dss-go/dss/spi/validation.KeyEntityTSPSource]
// implements for the offline, self-hosted TSA the library's own tests use.
type httpTSPSource struct {
	url    string
	client *http.Client
}

// displayURL is the TSA URL as error messages show it: with any password in
// its userinfo (https://user:secret@tsa.example/, which net/http turns into
// HTTP Basic credentials) replaced by "xxxxx", so that a failure printed to
// stderr never discloses the TSA credentials.
func (t *httpTSPSource) displayURL() string {
	u, err := url.Parse(t.url)
	if err != nil {
		return "the TSA"
	}
	return u.Redacted()
}

// newHTTPTSPSource returns a TSPSource that requests a time-stamp token from
// url over HTTP POST, following RFC 3161 section 3.4 (the same protocol
// every public TSA speaks, including the CA/Browser Forum ones).
func newHTTPTSPSource(url string) *httpTSPSource {
	return &httpTSPSource{url: url, client: &http.Client{Timeout: 30 * time.Second}}
}

// The RFC 3161 request/response ASN.1 shapes this client needs. Only the
// fields the client sends or reads are declared; everything else is either
// absent (an OPTIONAL field simply not present in the Go struct) or captured
// generically as asn1.RawValue without being interpreted, which is enough to
// extract the one thing the facade needs: the encoded time-stamp token.

// tsaAlgorithmIdentifier is AlgorithmIdentifier with its OPTIONAL parameters
// omitted, which is valid DER for an absent OPTIONAL field and is what every
// TSA this client has been tested against accepts for a plain hash OID.
type tsaAlgorithmIdentifier struct {
	Algorithm asn1.ObjectIdentifier
}

// tsaMessageImprint is MessageImprint (RFC 3161 section 2.4.1).
type tsaMessageImprint struct {
	HashAlgorithm tsaAlgorithmIdentifier
	HashedMessage []byte
}

// tsaRequest is TimeStampReq (RFC 3161 section 2.4.1). reqPolicy and the
// [0] extensions are never sent, so they are omitted from the struct
// entirely rather than encoded empty.
type tsaRequest struct {
	Version        int
	MessageImprint tsaMessageImprint
	Nonce          *big.Int `asn1:"optional"`
	CertReq        bool
}

// tsaPKIStatusInfo is PKIStatusInfo (RFC 3161 section 2.4.2), with
// statusString and failInfo captured generically: this client only needs
// the numeric status to decide whether a token follows.
type tsaPKIStatusInfo struct {
	Status       int
	StatusString asn1.RawValue `asn1:"optional"`
	FailInfo     asn1.RawValue `asn1:"optional"`
}

// tsaResponse is TimeStampResp (RFC 3161 section 2.4.2). Token, when
// present, is the encoded TimeStampToken (a CMS ContentInfo) - exactly the
// bytes [dss.TSPSource.TimeStampResponse] returns, byte for byte, since a
// time-stamp token is opaque to the client that requests it.
type tsaResponse struct {
	Status tsaPKIStatusInfo
	Token  asn1.RawValue `asn1:"optional"`
}

// PKIStatus values this client distinguishes (RFC 3161 section 2.4.2):
// granted and grantedWithMods both carry a usable token, everything else
// does not.
const (
	pkiStatusGranted         = 0
	pkiStatusGrantedWithMods = 1
	tsaMaxResponseBytes      = 1 << 20 // 1 MiB: generous for a time-stamp token, small enough to bound a hostile response.
	tsaRequestContentType    = "application/timestamp-query"
)

// TimeStampResponse implements [dss.TSPSource]. It builds a TimeStampReq
// over digest, POSTs it to the configured URL, and extracts the token from
// the TimeStampResp.
func (t *httpTSPSource) TimeStampResponse(digestAlgorithm enumerations.DigestAlgorithm, digest []byte) (*model.TimestampBinary, error) {
	oid, err := parseOID(digestAlgorithm.OID())
	if err != nil {
		return nil, fmt.Errorf("tsa: digest algorithm %s: %w", digestAlgorithm, err)
	}
	nonce, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 64))
	if err != nil {
		return nil, fmt.Errorf("tsa: generating a nonce: %w", err)
	}

	reqBytes, err := asn1.Marshal(tsaRequest{
		Version: 1,
		MessageImprint: tsaMessageImprint{
			HashAlgorithm: tsaAlgorithmIdentifier{Algorithm: oid},
			HashedMessage: digest,
		},
		Nonce:   nonce,
		CertReq: true,
	})
	if err != nil {
		return nil, fmt.Errorf("tsa: encoding the time-stamp request: %w", err)
	}

	httpReq, err := http.NewRequest(http.MethodPost, t.url, bytes.NewReader(reqBytes))
	if err != nil {
		return nil, fmt.Errorf("tsa: building the HTTP request for %s: invalid TSA URL", t.displayURL())
	}
	httpReq.Header.Set("Content-Type", tsaRequestContentType)

	// net/http redacts the URL password in the errors Do returns.
	httpResp, err := t.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("tsa: requesting a time-stamp: %w", err)
	}
	defer httpResp.Body.Close()

	respBytes, err := io.ReadAll(io.LimitReader(httpResp.Body, tsaMaxResponseBytes))
	if err != nil {
		return nil, fmt.Errorf("tsa: reading the response from %s: %w", t.displayURL(), err)
	}
	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tsa: %s returned HTTP %d", t.displayURL(), httpResp.StatusCode)
	}

	var resp tsaResponse
	if _, err := asn1.Unmarshal(respBytes, &resp); err != nil {
		return nil, fmt.Errorf("tsa: parsing the response from %s: %w", t.displayURL(), err)
	}
	if resp.Status.Status != pkiStatusGranted && resp.Status.Status != pkiStatusGrantedWithMods {
		return nil, fmt.Errorf("tsa: %s refused the request (PKIStatus %d)", t.displayURL(), resp.Status.Status)
	}
	if len(resp.Token.FullBytes) == 0 {
		return nil, fmt.Errorf("tsa: %s reported success but returned no time-stamp token", t.displayURL())
	}
	if err := checkTimeStampToken(resp.Token.FullBytes, oid, digest, nonce); err != nil {
		return nil, fmt.Errorf("tsa: %s: %w", t.displayURL(), err)
	}
	return model.NewTimestampBinary(resp.Token.FullBytes), nil
}

// checkTimeStampToken matches the token a TSA returned against the request it
// answers, the checks BouncyCastle's TimeStampResponse.validate(request) runs
// for Java DSS's OnlineTSPSource: the token must be a TimeStampToken whose
// TSTInfo covers exactly the requested message imprint (hash algorithm and
// digest) and echoes the request's nonce (RFC 3161 section 2.4.2: the nonce
// "MUST be present if the nonce field was present in the TimeStampReq").
// Without it a token for other data - or a replayed one - would be embedded
// in the signature as if it time-stamped this one.
func checkTimeStampToken(token []byte, oid asn1.ObjectIdentifier, digest []byte, nonce *big.Int) error {
	parsed, err := cmscore.ParseTimeStampToken(token)
	if err != nil {
		return fmt.Errorf("the returned time-stamp token cannot be parsed: %w", err)
	}
	info := parsed.TSTInfo()
	imprint := info.MessageImprint
	if imprint == nil || imprint.HashAlgorithm == nil || !imprint.HashAlgorithm.Algorithm.Equal(oid) {
		return fmt.Errorf("the time-stamp token is for a different message imprint algorithm")
	}
	if !bytes.Equal(imprint.HashedMessage, digest) {
		return fmt.Errorf("the time-stamp token is for a different message imprint digest")
	}
	if info.Nonce == nil || info.Nonce.Cmp(nonce) != 0 {
		return fmt.Errorf("the time-stamp token carries the wrong nonce value")
	}
	return nil
}

// parseOID parses a dotted-decimal OID string, e.g. "2.16.840.1.101.3.4.2.1",
// into an asn1.ObjectIdentifier. [enumerations.DigestAlgorithm.OID] always
// returns a value of this shape for a known algorithm.
func parseOID(s string) (asn1.ObjectIdentifier, error) {
	parts := strings.Split(s, ".")
	oid := make(asn1.ObjectIdentifier, len(parts))
	for i, part := range parts {
		n, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("invalid OID %q", s)
		}
		oid[i] = n
	}
	return oid, nil
}
