// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/source/LOTLSource.java (DSS 6.5.RC1).
//
// The two OtherTSLPointerType predicate fields carry the named OtherTSLPointerPredicate interface
// for the reason tl_source.go's header spells out: TLPredicateFactory, whose two factories seed
// them, produces exactly that type.
package tsl

// LOTLSource represents a List of Trusted Lists source.
type LOTLSource struct {
	TLSource

	// pivotSupport enables/disables pivot LOTL support.
	pivotSupport bool

	// mraSupport enables/disables MRA (Mutual Recognition Agreement) LOTL support.
	mraSupport bool

	// lotlPredicate filters the LOTL.
	//
	// Default: filter the XML European list of trusted list (LOTL).
	lotlPredicate OtherTSLPointerPredicate

	// tlPredicate filters the TLs.
	//
	// Default: filter all XML trusted lists (TL) for European countries.
	tlPredicate OtherTSLPointerPredicate

	// signingCertificatesAnnouncementPredicate optionally filters the URL where the provided
	// signing certificates are defined.
	signingCertificatesAnnouncementPredicate *LOTLSigningCertificatesAnnouncementSchemeInformationURI
}

// NewLOTLSource is the default constructor, instantiating the object with the minimal EU
// configuration. Port of LOTLSource(), whose two predicate field initialisers run before the
// (empty) constructor body.
func NewLOTLSource() *LOTLSource {
	source := &LOTLSource{TLSource: *NewTLSource()}
	source.pivotSupport = false
	source.mraSupport = false
	source.lotlPredicate = TLPredicateFactoryCreateEULOTLPredicate()
	source.tlPredicate = TLPredicateFactoryCreateEUTLPredicate()
	return source
}

// IsPivotSupport gets whether the LOTL configuration supports pivots. Port of isPivotSupport().
func (s *LOTLSource) IsPivotSupport() bool {
	return s.pivotSupport
}

// SetPivotSupport sets whether the LOTLSource shall support pivots. Port of
// setPivotSupport(boolean).
func (s *LOTLSource) SetPivotSupport(pivotSupport bool) {
	s.pivotSupport = pivotSupport
}

// IsMraSupport gets whether the LOTL configuration supports MRA (Mutual Recognition Agreement).
// Port of isMraSupport().
func (s *LOTLSource) IsMraSupport() bool {
	return s.mraSupport
}

// SetMraSupport sets whether the LOTL shall support the MRA (Mutual Recognition Agreement) scheme
// defining trust service equivalence mapping between the LOTL and a third-country Trusted List.
//
// Setting this condition to true allows processing a LOTL containing pointers to third-country
// trusted lists, including a special scheme for transitioning the qualification scope rules.
//
// Default: false (LOTL MRA is not supported). Port of setMraSupport(boolean).
func (s *LOTLSource) SetMraSupport(mraSupport bool) {
	s.mraSupport = mraSupport
}

// LotlPredicate gets the LOTL filtering predicate. Port of getLotlPredicate().
func (s *LOTLSource) LotlPredicate() OtherTSLPointerPredicate {
	return s.lotlPredicate
}

// SetLotlPredicate sets the LOTL filtering predicate. Port of setLotlPredicate(Predicate).
func (s *LOTLSource) SetLotlPredicate(lotlPredicate OtherTSLPointerPredicate) {
	s.lotlPredicate = lotlPredicate
}

// TlPredicate gets the TL filtering predicate. Port of getTlPredicate().
func (s *LOTLSource) TlPredicate() OtherTSLPointerPredicate {
	return s.tlPredicate
}

// SetTlPredicate sets the TL filtering predicate. Port of setTlPredicate(Predicate).
func (s *LOTLSource) SetTlPredicate(tlPredicate OtherTSLPointerPredicate) {
	s.tlPredicate = tlPredicate
}

// SigningCertificatesAnnouncementPredicate gets the
// LOTLSigningCertificatesAnnouncementSchemeInformationURI. Port of
// getSigningCertificatesAnnouncementPredicate().
func (s *LOTLSource) SigningCertificatesAnnouncementPredicate() *LOTLSigningCertificatesAnnouncementSchemeInformationURI {
	return s.signingCertificatesAnnouncementPredicate
}

// SetSigningCertificatesAnnouncementPredicate sets the
// LOTLSigningCertificatesAnnouncementSchemeInformationURI. Port of
// setSigningCertificatesAnnouncementPredicate(LOTLSigningCertificatesAnnouncementSchemeInformationURI).
func (s *LOTLSource) SetSigningCertificatesAnnouncementPredicate(
	signingCertificatesAnnouncementPredicate *LOTLSigningCertificatesAnnouncementSchemeInformationURI) {
	s.signingCertificatesAnnouncementPredicate = signingCertificatesAnnouncementPredicate
}
