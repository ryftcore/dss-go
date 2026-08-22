// Ported from dss-enumerations/.../X520Attributes.java (DSS 6.5.RC1).
//
// Contains attributes of a certificate's distinguished name.
package enumerations

import "fmt"

// X520Attributes represents attributes of a certificate's distinguished
// name. Implements OidDescription.
type X520Attributes string

// x520AttributesFields holds the (description, oid) pair for each constant.
type x520AttributesFields struct {
	description string
	oid         string
}

const (
	// X520AttributesObjectClass is the "objectClass" attribute (OID 2.5.4.0).
	X520AttributesObjectClass X520Attributes = "OBJECTCLASS"
	// X520AttributesAliasedEntryName is the "aliasedEntryName" attribute (OID 2.5.4.1).
	X520AttributesAliasedEntryName X520Attributes = "ALIASEDENTRYNAME"
	// X520AttributesEncryptedAliasedEntryName is the "encryptedAliasedEntryName" attribute (OID 2.5.4.1.2).
	X520AttributesEncryptedAliasedEntryName X520Attributes = "ENCRYPTEDALIASEDENTRYNAME"
	// X520AttributesKnowledgeInformation is the "knowledgeInformation" attribute (OID 2.5.4.2).
	X520AttributesKnowledgeInformation X520Attributes = "KNOWLEDGEINFORMATION"
	// X520AttributesCommonName is the "commonName" attribute (OID 2.5.4.3).
	X520AttributesCommonName X520Attributes = "COMMONNAME"
	// X520AttributesEncryptedCommonName is the "encryptedCommonName" attribute (OID 2.5.4.3.2).
	X520AttributesEncryptedCommonName X520Attributes = "ENCRYPTEDCOMMONNAME"
	// X520AttributesSurname is the "surname" attribute (OID 2.5.4.4).
	X520AttributesSurname X520Attributes = "SURNAME"
	// X520AttributesEncryptedSurname is the "encryptedSurname" attribute (OID 2.5.4.4.2).
	X520AttributesEncryptedSurname X520Attributes = "ENCRYPTEDSURNAME"
	// X520AttributesSerialNumber is the "serialNumber" attribute (OID 2.5.4.5).
	X520AttributesSerialNumber X520Attributes = "SERIALNUMBER"
	// X520AttributesEncryptedSerialNumber is the "encryptedSerialNumber" attribute (OID 2.5.4.5.2).
	X520AttributesEncryptedSerialNumber X520Attributes = "ENCRYPTEDSERIALNUMBER"
	// X520AttributesCountryName is the "countryName" attribute (OID 2.5.4.6).
	X520AttributesCountryName X520Attributes = "COUNTRYNAME"
	// X520AttributesEncryptedCountryName is the "encryptedCountryName" attribute (OID 2.5.4.6.2).
	X520AttributesEncryptedCountryName X520Attributes = "ENCRYPTEDCOUNTRYNAME"
	// X520AttributesLocalityName is the "localityName" attribute (OID 2.5.4.7).
	X520AttributesLocalityName X520Attributes = "LOCALITYNAME"
	// X520AttributesEncryptedLocalityName is the "encryptedLocalityName" attribute (OID 2.5.4.7.2).
	X520AttributesEncryptedLocalityName X520Attributes = "ENCRYPTEDLOCALITYNAME"
	// X520AttributesCollectiveLocalityName is the "collectiveLocalityName" attribute (OID 2.5.4.7.1).
	X520AttributesCollectiveLocalityName X520Attributes = "COLLECTIVELOCALITYNAME"
	// X520AttributesEncryptedCollectiveLocalityName is the "encryptedCollectiveLocalityName" attribute (OID 2.5.4.7.1.2).
	X520AttributesEncryptedCollectiveLocalityName X520Attributes = "ENCRYPTEDCOLLECTIVELOCALITYNAME"
	// X520AttributesStateOrProvinceName is the "stateOrProvinceName" attribute (OID 2.5.4.8).
	X520AttributesStateOrProvinceName X520Attributes = "STATEORPROVINCENAME"
	// X520AttributesEncryptedStateOrProvinceName is the "encryptedStateOrProvinceName" attribute (OID 2.5.4.8.2).
	X520AttributesEncryptedStateOrProvinceName X520Attributes = "ENCRYPTEDSTATEORPROVINCENAME"
	// X520AttributesCollectiveStateOrProvinceName is the "collectiveStateOrProvinceName" attribute (OID 2.5.4.8.1).
	X520AttributesCollectiveStateOrProvinceName X520Attributes = "COLLECTIVESTATEORPROVINCENAME"
	// X520AttributesEncryptedCollectiveStateOrProvinceName is the "encryptedCollectiveStateOrProvinceName" attribute (OID 2.5.4.8.1.2).
	X520AttributesEncryptedCollectiveStateOrProvinceName X520Attributes = "ENCRYPTEDCOLLECTIVESTATEORPROVINCENAME"
	// X520AttributesStreetAddress is the "streetAddress" attribute (OID 2.5.4.9).
	X520AttributesStreetAddress X520Attributes = "STREETADDRESS"
	// X520AttributesEncryptedStreetAddress is the "encryptedStreetAddress" attribute (OID 2.5.4.9.2).
	X520AttributesEncryptedStreetAddress X520Attributes = "ENCRYPTEDSTREETADDRESS"
	// X520AttributesCollectiveStreetAddress is the "collectiveStreetAddress" attribute (OID 2.5.4.9.1).
	X520AttributesCollectiveStreetAddress X520Attributes = "COLLECTIVESTREETADDRESS"
	// X520AttributesEncryptedCollectiveStreetAddress is the "encryptedCollectiveStreetAddress" attribute (OID 2.5.4.9.1.2).
	X520AttributesEncryptedCollectiveStreetAddress X520Attributes = "ENCRYPTEDCOLLECTIVESTREETADDRESS"
	// X520AttributesOrganizationName is the "organizationName" attribute (OID 2.5.4.10).
	X520AttributesOrganizationName X520Attributes = "ORGANIZATIONNAME"
	// X520AttributesEncryptedOrganizationName is the "encryptedOrganizationName" attribute (OID 2.5.4.10.2).
	X520AttributesEncryptedOrganizationName X520Attributes = "ENCRYPTEDORGANIZATIONNAME"
	// X520AttributesCollectiveOrganizationName is the "collectiveOrganizationName" attribute (OID 2.5.4.10.1).
	X520AttributesCollectiveOrganizationName X520Attributes = "COLLECTIVEORGANIZATIONNAME"
	// X520AttributesEncryptedCollectiveOrganizationName is the "encryptedCollectiveOrganizationName" attribute (OID 2.5.4.10.1.2).
	X520AttributesEncryptedCollectiveOrganizationName X520Attributes = "ENCRYPTEDCOLLECTIVEORGANIZATIONNAME"
	// X520AttributesOrganizationalUnitName is the "organizationalUnitName" attribute (OID 2.5.4.11).
	X520AttributesOrganizationalUnitName X520Attributes = "ORGANIZATIONALUNITNAME"
	// X520AttributesEncryptedOrganizationalUnitName is the "encryptedOrganizationalUnitName" attribute (OID 2.5.4.11.2).
	X520AttributesEncryptedOrganizationalUnitName X520Attributes = "ENCRYPTEDORGANIZATIONALUNITNAME"
	// X520AttributesCollectiveOrganizationalUnitName is the "collectiveOrganizationalUnitName" attribute (OID 2.5.4.11.1).
	X520AttributesCollectiveOrganizationalUnitName X520Attributes = "COLLECTIVEORGANIZATIONALUNITNAME"
	// X520AttributesEncryptedCollectiveOrganizationalUnitNam is the "encryptedCollectiveOrganizationalUnitNam" attribute (OID 2.5.4.11.1.2).
	X520AttributesEncryptedCollectiveOrganizationalUnitNam X520Attributes = "ENCRYPTEDCOLLECTIVEORGANIZATIONALUNITNAM"
	// X520AttributesTitle is the "title" attribute (OID 2.5.4.12).
	X520AttributesTitle X520Attributes = "TITLE"
	// X520AttributesEncryptedTitle is the "encryptedTitle" attribute (OID 2.5.4.12.2).
	X520AttributesEncryptedTitle X520Attributes = "ENCRYPTEDTITLE"
	// X520AttributesDescription is the "description" attribute (OID 2.5.4.13).
	X520AttributesDescription X520Attributes = "DESCRIPTION"
	// X520AttributesEncryptedDescription is the "encryptedDescription" attribute (OID 2.5.4.13.2).
	X520AttributesEncryptedDescription X520Attributes = "ENCRYPTEDDESCRIPTION"
	// X520AttributesSearchGuide is the "searchGuide" attribute (OID 2.5.4.14).
	X520AttributesSearchGuide X520Attributes = "SEARCHGUIDE"
	// X520AttributesEncryptedSearchGuide is the "encryptedSearchGuide" attribute (OID 2.5.4.14.2).
	X520AttributesEncryptedSearchGuide X520Attributes = "ENCRYPTEDSEARCHGUIDE"
	// X520AttributesBusinessCategory is the "businessCategory" attribute (OID 2.5.4.15).
	X520AttributesBusinessCategory X520Attributes = "BUSINESSCATEGORY"
	// X520AttributesEncryptedBusinessCategory is the "encryptedBusinessCategory" attribute (OID 2.5.4.15.2).
	X520AttributesEncryptedBusinessCategory X520Attributes = "ENCRYPTEDBUSINESSCATEGORY"
	// X520AttributesPostalAddress is the "postalAddress" attribute (OID 2.5.4.16).
	X520AttributesPostalAddress X520Attributes = "POSTALADDRESS"
	// X520AttributesEncryptedPostalAddress is the "encryptedPostalAddress" attribute (OID 2.5.4.16.2).
	X520AttributesEncryptedPostalAddress X520Attributes = "ENCRYPTEDPOSTALADDRESS"
	// X520AttributesCollectivePostalAddress is the "collectivePostalAddress" attribute (OID 2.5.4.16.1).
	X520AttributesCollectivePostalAddress X520Attributes = "COLLECTIVEPOSTALADDRESS"
	// X520AttributesEncryptedCollectivePostalAddress is the "encryptedCollectivePostalAddress" attribute (OID 2.5.4.16.1.2).
	X520AttributesEncryptedCollectivePostalAddress X520Attributes = "ENCRYPTEDCOLLECTIVEPOSTALADDRESS"
	// X520AttributesPostalCode is the "postalCode" attribute (OID 2.5.4.17).
	X520AttributesPostalCode X520Attributes = "POSTALCODE"
	// X520AttributesEncryptedPostalCode is the "encryptedPostalCode" attribute (OID 2.5.4.17.2).
	X520AttributesEncryptedPostalCode X520Attributes = "ENCRYPTEDPOSTALCODE"
	// X520AttributesCollectivePostalCode is the "collectivePostalCode" attribute (OID 2.5.4.17.1).
	X520AttributesCollectivePostalCode X520Attributes = "COLLECTIVEPOSTALCODE"
	// X520AttributesEncryptedCollectivePostalCode is the "encryptedCollectivePostalCode" attribute (OID 2.5.4.17.1.2).
	X520AttributesEncryptedCollectivePostalCode X520Attributes = "ENCRYPTEDCOLLECTIVEPOSTALCODE"
	// X520AttributesPostOfficeBox is the "postOfficeBox" attribute (OID 2.5.4.18).
	X520AttributesPostOfficeBox X520Attributes = "POSTOFFICEBOX"
	// X520AttributesCollectivePostOfficeBox is the "collectivePostOfficeBox" attribute (OID 2.5.4.18.1).
	X520AttributesCollectivePostOfficeBox X520Attributes = "COLLECTIVEPOSTOFFICEBOX"
	// X520AttributesEncryptedPostOfficeBox is the "encryptedPostOfficeBox" attribute (OID 2.5.4.18.2).
	X520AttributesEncryptedPostOfficeBox X520Attributes = "ENCRYPTEDPOSTOFFICEBOX"
	// X520AttributesEncryptedCollectivePostOfficeBox is the "encryptedCollectivePostOfficeBox" attribute (OID 2.5.4.18.1.2).
	X520AttributesEncryptedCollectivePostOfficeBox X520Attributes = "ENCRYPTEDCOLLECTIVEPOSTOFFICEBOX"
	// X520AttributesPhysicalDeliveryOfficeName is the "physicalDeliveryOfficeName" attribute (OID 2.5.4.19).
	X520AttributesPhysicalDeliveryOfficeName X520Attributes = "PHYSICALDELIVERYOFFICENAME"
	// X520AttributesCollectivePhysicalDeliveryOfficeName is the "collectivePhysicalDeliveryOfficeName" attribute (OID 2.5.4.19.1).
	X520AttributesCollectivePhysicalDeliveryOfficeName X520Attributes = "COLLECTIVEPHYSICALDELIVERYOFFICENAME"
	// X520AttributesEncryptedPhysicalDeliveryOfficeName is the "encryptedPhysicalDeliveryOfficeName" attribute (OID 2.5.4.19.2).
	X520AttributesEncryptedPhysicalDeliveryOfficeName X520Attributes = "ENCRYPTEDPHYSICALDELIVERYOFFICENAME"
	// X520AttributesEncryptedCollectivePhysicalDeliveryOfficeName is the "encryptedCollectivePhysicalDeliveryOfficeName" attribute (OID 2.5.4.19.1.2).
	X520AttributesEncryptedCollectivePhysicalDeliveryOfficeName X520Attributes = "ENCRYPTEDCOLLECTIVEPHYSICALDELIVERYOFFICENAME"
	// X520AttributesTelephoneNumber is the "telephoneNumber" attribute (OID 2.5.4.20).
	X520AttributesTelephoneNumber X520Attributes = "TELEPHONENUMBER"
	// X520AttributesEncryptedTelephoneNumber is the "encryptedTelephoneNumber" attribute (OID 2.5.4.20.2).
	X520AttributesEncryptedTelephoneNumber X520Attributes = "ENCRYPTEDTELEPHONENUMBER"
	// X520AttributesCollectiveTelephoneNumber is the "collectiveTelephoneNumber" attribute (OID 2.5.4.20.1).
	X520AttributesCollectiveTelephoneNumber X520Attributes = "COLLECTIVETELEPHONENUMBER"
	// X520AttributesEncryptedCollectiveTelephoneNumber is the "encryptedCollectiveTelephoneNumber" attribute (OID 2.5.4.20.1.2).
	X520AttributesEncryptedCollectiveTelephoneNumber X520Attributes = "ENCRYPTEDCOLLECTIVETELEPHONENUMBER"
	// X520AttributesTelexNumber is the "telexNumber" attribute (OID 2.5.4.21).
	X520AttributesTelexNumber X520Attributes = "TELEXNUMBER"
	// X520AttributesEncryptedTelexNumber is the "encryptedTelexNumber" attribute (OID 2.5.4.21.2).
	X520AttributesEncryptedTelexNumber X520Attributes = "ENCRYPTEDTELEXNUMBER"
	// X520AttributesCollectiveTelexNumber is the "collectiveTelexNumber" attribute (OID 2.5.4.21.1).
	X520AttributesCollectiveTelexNumber X520Attributes = "COLLECTIVETELEXNUMBER"
	// X520AttributesEncryptedCollectiveTelexNumber is the "encryptedCollectiveTelexNumber" attribute (OID 2.5.4.21.1.2).
	X520AttributesEncryptedCollectiveTelexNumber X520Attributes = "ENCRYPTEDCOLLECTIVETELEXNUMBER"
	// X520AttributesTeletexTerminalIdentifier is the "teletexTerminalIdentifier" attribute (OID 2.5.4.22).
	X520AttributesTeletexTerminalIdentifier X520Attributes = "TELETEXTERMINALIDENTIFIER"
	// X520AttributesEncryptedTeletexTerminalIdentifier is the "encryptedTeletexTerminalIdentifier" attribute (OID 2.5.4.22.2).
	X520AttributesEncryptedTeletexTerminalIdentifier X520Attributes = "ENCRYPTEDTELETEXTERMINALIDENTIFIER"
	// X520AttributesCollectiveTeletexTerminalIdentifier is the "collectiveTeletexTerminalIdentifier" attribute (OID 2.5.4.22.1).
	X520AttributesCollectiveTeletexTerminalIdentifier X520Attributes = "COLLECTIVETELETEXTERMINALIDENTIFIER"
	// X520AttributesEncryptedCollectiveTeletexTerminalIdentifier is the "encryptedCollectiveTeletexTerminalIdentifier" attribute (OID 2.5.4.22.1.2).
	X520AttributesEncryptedCollectiveTeletexTerminalIdentifier X520Attributes = "ENCRYPTEDCOLLECTIVETELETEXTERMINALIDENTIFIER"
	// X520AttributesFacsimileTelephoneNumber is the "facsimileTelephoneNumber" attribute (OID 2.5.4.23).
	X520AttributesFacsimileTelephoneNumber X520Attributes = "FACSIMILETELEPHONENUMBER"
	// X520AttributesEncryptedFacsimileTelephoneNumber is the "encryptedFacsimileTelephoneNumber" attribute (OID 2.5.4.23.2).
	X520AttributesEncryptedFacsimileTelephoneNumber X520Attributes = "ENCRYPTEDFACSIMILETELEPHONENUMBER"
	// X520AttributesCollectiveFacsimileTelephoneNumber is the "collectiveFacsimileTelephoneNumber" attribute (OID 2.5.4.23.1).
	X520AttributesCollectiveFacsimileTelephoneNumber X520Attributes = "COLLECTIVEFACSIMILETELEPHONENUMBER"
	// X520AttributesEncryptedCollectiveFacsimileTelephoneNumber is the "encryptedCollectiveFacsimileTelephoneNumber" attribute (OID 2.5.4.23.1.2).
	X520AttributesEncryptedCollectiveFacsimileTelephoneNumber X520Attributes = "ENCRYPTEDCOLLECTIVEFACSIMILETELEPHONENUMBER"
	// X520AttributesX121Address is the "x121Address" attribute (OID 2.5.4.24).
	X520AttributesX121Address X520Attributes = "X121ADDRESS"
	// X520AttributesEncryptedX121Address is the "encryptedX121Address" attribute (OID 2.5.4.24.2).
	X520AttributesEncryptedX121Address X520Attributes = "ENCRYPTEDX121ADDRESS"
	// X520AttributesInternationalISDNNumber is the "internationalISDNNumber" attribute (OID 2.5.4.25).
	X520AttributesInternationalISDNNumber X520Attributes = "INTERNATIONALISDNNUMBER"
	// X520AttributesEncryptedInternationalISDNNumber is the "encryptedInternationalISDNNumber" attribute (OID 2.5.4.25.2).
	X520AttributesEncryptedInternationalISDNNumber X520Attributes = "ENCRYPTEDINTERNATIONALISDNNUMBER"
	// X520AttributesCollectiveInternationalISDNNumber is the "collectiveInternationalISDNNumber" attribute (OID 2.5.4.25.1).
	X520AttributesCollectiveInternationalISDNNumber X520Attributes = "COLLECTIVEINTERNATIONALISDNNUMBER"
	// X520AttributesEncryptedCollectiveInternationalISDNNumber is the "encryptedCollectiveInternationalISDNNumber" attribute (OID 2.5.4.25.1.2).
	X520AttributesEncryptedCollectiveInternationalISDNNumber X520Attributes = "ENCRYPTEDCOLLECTIVEINTERNATIONALISDNNUMBER"
	// X520AttributesRegisteredAddress is the "registeredAddress" attribute (OID 2.5.4.26).
	X520AttributesRegisteredAddress X520Attributes = "REGISTEREDADDRESS"
	// X520AttributesEncryptedRegisteredAddress is the "encryptedRegisteredAddress" attribute (OID 2.5.4.26.2).
	X520AttributesEncryptedRegisteredAddress X520Attributes = "ENCRYPTEDREGISTEREDADDRESS"
	// X520AttributesDestinationIndicator is the "destinationIndicator" attribute (OID 2.5.4.27).
	X520AttributesDestinationIndicator X520Attributes = "DESTINATIONINDICATOR"
	// X520AttributesEncryptedDestinationIndicator is the "encryptedDestinationIndicator" attribute (OID 2.5.4.27.2).
	X520AttributesEncryptedDestinationIndicator X520Attributes = "ENCRYPTEDDESTINATIONINDICATOR"
	// X520AttributesPreferredDeliveryMethod is the "preferredDeliveryMethod" attribute (OID 2.5.4.28).
	X520AttributesPreferredDeliveryMethod X520Attributes = "PREFERREDDELIVERYMETHOD"
	// X520AttributesEncryptedPreferredDeliveryMethod is the "encryptedPreferredDeliveryMethod" attribute (OID 2.5.4.28.2).
	X520AttributesEncryptedPreferredDeliveryMethod X520Attributes = "ENCRYPTEDPREFERREDDELIVERYMETHOD"
	// X520AttributesPresentationAddress is the "presentationAddress" attribute (OID 2.5.4.29).
	X520AttributesPresentationAddress X520Attributes = "PRESENTATIONADDRESS"
	// X520AttributesEncryptedPresentationAddress is the "encryptedPresentationAddress" attribute (OID 2.5.4.29.2).
	X520AttributesEncryptedPresentationAddress X520Attributes = "ENCRYPTEDPRESENTATIONADDRESS"
	// X520AttributesSupportedApplicationContext is the "supportedApplicationContext" attribute (OID 2.5.4.30).
	X520AttributesSupportedApplicationContext X520Attributes = "SUPPORTEDAPPLICATIONCONTEXT"
	// X520AttributesEncryptedSupportedApplicationContext is the "encryptedSupportedApplicationContext" attribute (OID 2.5.4.30.2).
	X520AttributesEncryptedSupportedApplicationContext X520Attributes = "ENCRYPTEDSUPPORTEDAPPLICATIONCONTEXT"
	// X520AttributesMember is the "member" attribute (OID 2.5.4.31).
	X520AttributesMember X520Attributes = "MEMBER"
	// X520AttributesEncryptedMember is the "encryptedMember" attribute (OID 2.5.4.31.2).
	X520AttributesEncryptedMember X520Attributes = "ENCRYPTEDMEMBER"
	// X520AttributesOwner is the "owner" attribute (OID 2.5.4.32).
	X520AttributesOwner X520Attributes = "OWNER"
	// X520AttributesEncryptedOwner is the "encryptedOwner" attribute (OID 2.5.4.32.2).
	X520AttributesEncryptedOwner X520Attributes = "ENCRYPTEDOWNER"
	// X520AttributesRoleOccupant is the "roleOccupant" attribute (OID 2.5.4.33).
	X520AttributesRoleOccupant X520Attributes = "ROLEOCCUPANT"
	// X520AttributesEncryptedRoleOccupant is the "encryptedRoleOccupant" attribute (OID 2.5.4.33.2).
	X520AttributesEncryptedRoleOccupant X520Attributes = "ENCRYPTEDROLEOCCUPANT"
	// X520AttributesSeeAlso is the "seeAlso" attribute (OID 2.5.4.34).
	X520AttributesSeeAlso X520Attributes = "SEEALSO"
	// X520AttributesEncryptedSeeAlso is the "encryptedSeeAlso" attribute (OID 2.5.4.34.2).
	X520AttributesEncryptedSeeAlso X520Attributes = "ENCRYPTEDSEEALSO"
	// X520AttributesUserPassword is the "userPassword" attribute (OID 2.5.4.35).
	X520AttributesUserPassword X520Attributes = "USERPASSWORD"
	// X520AttributesEncryptedUserPassword is the "encryptedUserPassword" attribute (OID 2.5.4.35.2).
	X520AttributesEncryptedUserPassword X520Attributes = "ENCRYPTEDUSERPASSWORD"
	// X520AttributesUserCertificate is the "userCertificate" attribute (OID 2.5.4.36).
	X520AttributesUserCertificate X520Attributes = "USERCERTIFICATE"
	// X520AttributesEncryptedUserCertificate is the "encryptedUserCertificate" attribute (OID 2.5.4.36.2).
	X520AttributesEncryptedUserCertificate X520Attributes = "ENCRYPTEDUSERCERTIFICATE"
	// X520AttributesCACertificate is the "cACertificate" attribute (OID 2.5.4.37).
	X520AttributesCACertificate X520Attributes = "CACERTIFICATE"
	// X520AttributesEncryptedCACertificate is the "encryptedCACertificate" attribute (OID 2.5.4.37.2).
	X520AttributesEncryptedCACertificate X520Attributes = "ENCRYPTEDCACERTIFICATE"
	// X520AttributesAuthorityRevocationList is the "authorityRevocationList" attribute (OID 2.5.4.38).
	X520AttributesAuthorityRevocationList X520Attributes = "AUTHORITYREVOCATIONLIST"
	// X520AttributesEncryptedAuthorityRevocationList is the "encryptedAuthorityRevocationList" attribute (OID 2.5.4.38.2).
	X520AttributesEncryptedAuthorityRevocationList X520Attributes = "ENCRYPTEDAUTHORITYREVOCATIONLIST"
	// X520AttributesCertificateRevocationList is the "certificateRevocationList" attribute (OID 2.5.4.39).
	X520AttributesCertificateRevocationList X520Attributes = "CERTIFICATEREVOCATIONLIST"
	// X520AttributesEncryptedCertificateRevocationList is the "encryptedCertificateRevocationList" attribute (OID 2.5.4.39.2).
	X520AttributesEncryptedCertificateRevocationList X520Attributes = "ENCRYPTEDCERTIFICATEREVOCATIONLIST"
	// X520AttributesCrossCertificatePair is the "crossCertificatePair" attribute (OID 2.5.4.40).
	X520AttributesCrossCertificatePair X520Attributes = "CROSSCERTIFICATEPAIR"
	// X520AttributesEncryptedCrossCertificatePair is the "encryptedCrossCertificatePair" attribute (OID 2.5.4.40.2).
	X520AttributesEncryptedCrossCertificatePair X520Attributes = "ENCRYPTEDCROSSCERTIFICATEPAIR"
	// X520AttributesName is the "name" attribute (OID 2.5.4.41).
	X520AttributesName X520Attributes = "NAME"
	// X520AttributesGivenName is the "givenName" attribute (OID 2.5.4.42).
	X520AttributesGivenName X520Attributes = "GIVENNAME"
	// X520AttributesEncryptedGivenName is the "encryptedGivenName" attribute (OID 2.5.4.42.2).
	X520AttributesEncryptedGivenName X520Attributes = "ENCRYPTEDGIVENNAME"
	// X520AttributesInitials is the "initials" attribute (OID 2.5.4.43).
	X520AttributesInitials X520Attributes = "INITIALS"
	// X520AttributesEncryptedInitials is the "encryptedInitials" attribute (OID 2.5.4.43.2).
	X520AttributesEncryptedInitials X520Attributes = "ENCRYPTEDINITIALS"
	// X520AttributesGenerationQualifier is the "generationQualifier" attribute (OID 2.5.4.44).
	X520AttributesGenerationQualifier X520Attributes = "GENERATIONQUALIFIER"
	// X520AttributesEncryptedGenerationQualifier is the "encryptedGenerationQualifier" attribute (OID 2.5.4.44.2).
	X520AttributesEncryptedGenerationQualifier X520Attributes = "ENCRYPTEDGENERATIONQUALIFIER"
	// X520AttributesUniqueIdentifier is the "uniqueIdentifier" attribute (OID 2.5.4.45).
	X520AttributesUniqueIdentifier X520Attributes = "UNIQUEIDENTIFIER"
	// X520AttributesEncryptedUniqueIdentifier is the "encryptedUniqueIdentifier" attribute (OID 2.5.4.45.2).
	X520AttributesEncryptedUniqueIdentifier X520Attributes = "ENCRYPTEDUNIQUEIDENTIFIER"
	// X520AttributesDNQualifier is the "dnQualifier" attribute (OID 2.5.4.46).
	X520AttributesDNQualifier X520Attributes = "DNQUALIFIER"
	// X520AttributesEncryptedDNQualifier is the "encryptedDnQualifier" attribute (OID 2.5.4.46.2).
	X520AttributesEncryptedDNQualifier X520Attributes = "ENCRYPTEDDNQUALIFIER"
	// X520AttributesEnhancedSearchGuide is the "enhancedSearchGuide" attribute (OID 2.5.4.47).
	X520AttributesEnhancedSearchGuide X520Attributes = "ENHANCEDSEARCHGUIDE"
	// X520AttributesEncryptedEnhancedSearchGuide is the "encryptedEnhancedSearchGuide" attribute (OID 2.5.4.47.2).
	X520AttributesEncryptedEnhancedSearchGuide X520Attributes = "ENCRYPTEDENHANCEDSEARCHGUIDE"
	// X520AttributesProtocolInformation is the "protocolInformation" attribute (OID 2.5.4.48).
	X520AttributesProtocolInformation X520Attributes = "PROTOCOLINFORMATION"
	// X520AttributesEncryptedProtocolInformation is the "encryptedProtocolInformation" attribute (OID 2.5.4.48.2).
	X520AttributesEncryptedProtocolInformation X520Attributes = "ENCRYPTEDPROTOCOLINFORMATION"
	// X520AttributesDistinguishedName is the "distinguishedName" attribute (OID 2.5.4.49).
	X520AttributesDistinguishedName X520Attributes = "DISTINGUISHEDNAME"
	// X520AttributesEncryptedDistinguishedName is the "encryptedDistinguishedName" attribute (OID 2.5.4.49.2).
	X520AttributesEncryptedDistinguishedName X520Attributes = "ENCRYPTEDDISTINGUISHEDNAME"
	// X520AttributesUniqueMember is the "uniqueMember" attribute (OID 2.5.4.50).
	X520AttributesUniqueMember X520Attributes = "UNIQUEMEMBER"
	// X520AttributesEncryptedUniqueMember is the "encryptedUniqueMember" attribute (OID 2.5.4.50.2).
	X520AttributesEncryptedUniqueMember X520Attributes = "ENCRYPTEDUNIQUEMEMBER"
	// X520AttributesHouseIdentifier is the "houseIdentifier" attribute (OID 2.5.4.51).
	X520AttributesHouseIdentifier X520Attributes = "HOUSEIDENTIFIER"
	// X520AttributesEncryptedHouseIdentifier is the "encryptedHouseIdentifier" attribute (OID 2.5.4.51.2).
	X520AttributesEncryptedHouseIdentifier X520Attributes = "ENCRYPTEDHOUSEIDENTIFIER"
	// X520AttributesSupportedAlgorithms is the "supportedAlgorithms" attribute (OID 2.5.4.52).
	X520AttributesSupportedAlgorithms X520Attributes = "SUPPORTEDALGORITHMS"
	// X520AttributesEncryptedSupportedAlgorithms is the "encryptedSupportedAlgorithms" attribute (OID 2.5.4.52.2).
	X520AttributesEncryptedSupportedAlgorithms X520Attributes = "ENCRYPTEDSUPPORTEDALGORITHMS"
	// X520AttributesDeltaRevocationList is the "deltaRevocationList" attribute (OID 2.5.4.53).
	X520AttributesDeltaRevocationList X520Attributes = "DELTAREVOCATIONLIST"
	// X520AttributesEncryptedDeltaRevocationList is the "encryptedDeltaRevocationList" attribute (OID 2.5.4.53.2).
	X520AttributesEncryptedDeltaRevocationList X520Attributes = "ENCRYPTEDDELTAREVOCATIONLIST"
	// X520AttributesDMDName is the "dmdName" attribute (OID 2.5.4.54).
	X520AttributesDMDName X520Attributes = "DMDNAME"
	// X520AttributesEncryptedDMDName is the "encryptedDmdName" attribute (OID 2.5.4.54.2).
	X520AttributesEncryptedDMDName X520Attributes = "ENCRYPTEDDMDNAME"
	// X520AttributesClearance is the "clearance" attribute (OID 2.5.4.55).
	X520AttributesClearance X520Attributes = "CLEARANCE"
	// X520AttributesEncryptedClearance is the "encryptedClearance" attribute (OID 2.5.4.55.2).
	X520AttributesEncryptedClearance X520Attributes = "ENCRYPTEDCLEARANCE"
	// X520AttributesDefaultDirQOP is the "defaultDirQop" attribute (OID 2.5.4.56).
	X520AttributesDefaultDirQOP X520Attributes = "DEFAULTDIRQOP"
	// X520AttributesEncryptedDefaultDirQOP is the "encryptedDefaultDirQop" attribute (OID 2.5.4.56.2).
	X520AttributesEncryptedDefaultDirQOP X520Attributes = "ENCRYPTEDDEFAULTDIRQOP"
	// X520AttributesAttributeIntegrityInfo is the "attributeIntegrityInfo" attribute (OID 2.5.4.57).
	X520AttributesAttributeIntegrityInfo X520Attributes = "ATTRIBUTEINTEGRITYINFO"
	// X520AttributesEncryptedAttributeIntegrityInfo is the "encryptedAttributeIntegrityInfo" attribute (OID 2.5.4.57.2).
	X520AttributesEncryptedAttributeIntegrityInfo X520Attributes = "ENCRYPTEDATTRIBUTEINTEGRITYINFO"
	// X520AttributesAttributeCertificate is the "attributeCertificate" attribute (OID 2.5.4.58).
	X520AttributesAttributeCertificate X520Attributes = "ATTRIBUTECERTIFICATE"
	// X520AttributesEncryptedAttributeCertificate is the "encryptedAttributeCertificate" attribute (OID 2.5.4.58.2).
	X520AttributesEncryptedAttributeCertificate X520Attributes = "ENCRYPTEDATTRIBUTECERTIFICATE"
	// X520AttributesAttributeCertificateRevocationList is the "attributeCertificateRevocationList" attribute (OID 2.5.4.59).
	X520AttributesAttributeCertificateRevocationList X520Attributes = "ATTRIBUTECERTIFICATEREVOCATIONLIST"
	// X520AttributesEncryptedAttributeCertificateRevocationList is the "encryptedAttributeCertificateRevocationList" attribute (OID 2.5.4.59.2).
	X520AttributesEncryptedAttributeCertificateRevocationList X520Attributes = "ENCRYPTEDATTRIBUTECERTIFICATEREVOCATIONLIST"
	// X520AttributesConfKeyInfo is the "confKeyInfo" attribute (OID 2.5.4.60).
	X520AttributesConfKeyInfo X520Attributes = "CONFKEYINFO"
	// X520AttributesEncryptedConfKeyInfo is the "encryptedConfKeyInfo" attribute (OID 2.5.4.60.2).
	X520AttributesEncryptedConfKeyInfo X520Attributes = "ENCRYPTEDCONFKEYINFO"
	// X520AttributesAACertificate is the "aACertificate" attribute (OID 2.5.4.61).
	X520AttributesAACertificate X520Attributes = "AACERTIFICATE"
	// X520AttributesAttributeDescriptorCertificate is the "attributeDescriptorCertificate" attribute (OID 2.5.4.62).
	X520AttributesAttributeDescriptorCertificate X520Attributes = "ATTRIBUTEDESCRIPTORCERTIFICATE"
	// X520AttributesAttributeAuthorityRevocationList is the "attributeAuthorityRevocationList" attribute (OID 2.5.4.63).
	X520AttributesAttributeAuthorityRevocationList X520Attributes = "ATTRIBUTEAUTHORITYREVOCATIONLIST"
	// X520AttributesFamilyInformation is the "family-information" attribute (OID 2.5.4.64).
	X520AttributesFamilyInformation X520Attributes = "FAMILY_INFORMATION"
	// X520AttributesPseudonym is the "pseudonym" attribute (OID 2.5.4.65).
	X520AttributesPseudonym X520Attributes = "PSEUDONYM"
	// X520AttributesCommunicationsService is the "communicationsService" attribute (OID 2.5.4.66).
	X520AttributesCommunicationsService X520Attributes = "COMMUNICATIONSSERVICE"
	// X520AttributesCommunicationsNetwork is the "communicationsNetwork" attribute (OID 2.5.4.67).
	X520AttributesCommunicationsNetwork X520Attributes = "COMMUNICATIONSNETWORK"
	// X520AttributesCertificationPracticeStmt is the "certificationPracticeStmt" attribute (OID 2.5.4.68).
	X520AttributesCertificationPracticeStmt X520Attributes = "CERTIFICATIONPRACTICESTMT"
	// X520AttributesCertificatePolicy is the "certificatePolicy" attribute (OID 2.5.4.69).
	X520AttributesCertificatePolicy X520Attributes = "CERTIFICATEPOLICY"
	// X520AttributesPKIPath is the "pkiPath" attribute (OID 2.5.4.70).
	X520AttributesPKIPath X520Attributes = "PKIPATH"
	// X520AttributesPrivPolicy is the "privPolicy" attribute (OID 2.5.4.71).
	X520AttributesPrivPolicy X520Attributes = "PRIVPOLICY"
	// X520AttributesRole is the "role" attribute (OID 2.5.4.72).
	X520AttributesRole X520Attributes = "ROLE"
	// X520AttributesDelegationPath is the "delegationPath" attribute (OID 2.5.4.73).
	X520AttributesDelegationPath X520Attributes = "DELEGATIONPATH"
	// X520AttributesProtPrivPolicy is the "protPrivPolicy" attribute (OID 2.5.4.74).
	X520AttributesProtPrivPolicy X520Attributes = "PROTPRIVPOLICY"
	// X520AttributesXMLPrivilegeInfo is the "xMLPrivilegeInfo" attribute (OID 2.5.4.75).
	X520AttributesXMLPrivilegeInfo X520Attributes = "XMLPRIVILEGEINFO"
	// X520AttributesXMLPrivPolicy is the "xmlPrivPolicy" attribute (OID 2.5.4.76).
	X520AttributesXMLPrivPolicy X520Attributes = "XMLPRIVPOLICY"
	// X520AttributesUUIDPair is the "uuidpair" attribute (OID 2.5.4.77).
	X520AttributesUUIDPair X520Attributes = "UUIDPAIR"
	// X520AttributesTagOID is the "tagOid" attribute (OID 2.5.4.78).
	X520AttributesTagOID X520Attributes = "TAGOID"
	// X520AttributesUIIFormat is the "uiiFormat" attribute (OID 2.5.4.79).
	X520AttributesUIIFormat X520Attributes = "UIIFORMAT"
	// X520AttributesUIIInURN is the "uiiInUrn" attribute (OID 2.5.4.80).
	X520AttributesUIIInURN X520Attributes = "UIIINURN"
	// X520AttributesContentURL is the "contentUrl" attribute (OID 2.5.4.81).
	X520AttributesContentURL X520Attributes = "CONTENTURL"
	// X520AttributesPermission is the "permission" attribute (OID 2.5.4.82).
	X520AttributesPermission X520Attributes = "PERMISSION"
	// X520AttributesURI is the "uri" attribute (OID 2.5.4.83).
	X520AttributesURI X520Attributes = "URI"
	// X520AttributesPwdAttribute is the "pwdAttribute" attribute (OID 2.5.4.84).
	X520AttributesPwdAttribute X520Attributes = "PWDATTRIBUTE"
	// X520AttributesUserPwd is the "userPwd" attribute (OID 2.5.4.85).
	X520AttributesUserPwd X520Attributes = "USERPWD"
	// X520AttributesURN is the "urn" attribute (OID 2.5.4.86).
	X520AttributesURN X520Attributes = "URN"
	// X520AttributesURL is the "url" attribute (OID 2.5.4.87).
	X520AttributesURL X520Attributes = "URL"
	// X520AttributesUTMCoordinates is the "utmCoordinates" attribute (OID 2.5.4.88).
	X520AttributesUTMCoordinates X520Attributes = "UTMCOORDINATES"
	// X520AttributesURNC is the "urnC" attribute (OID 2.5.4.89).
	X520AttributesURNC X520Attributes = "URNC"
	// X520AttributesUII is the "uii" attribute (OID 2.5.4.90).
	X520AttributesUII X520Attributes = "UII"
	// X520AttributesEPC is the "epc" attribute (OID 2.5.4.91).
	X520AttributesEPC X520Attributes = "EPC"
	// X520AttributesTagAFI is the "tagAfi" attribute (OID 2.5.4.92).
	X520AttributesTagAFI X520Attributes = "TAGAFI"
	// X520AttributesEPCFormat is the "epcFormat" attribute (OID 2.5.4.93).
	X520AttributesEPCFormat X520Attributes = "EPCFORMAT"
	// X520AttributesEPCInURN is the "epcInUrn" attribute (OID 2.5.4.94).
	X520AttributesEPCInURN X520Attributes = "EPCINURN"
	// X520AttributesLDAPURL is the "ldapUrl" attribute (OID 2.5.4.95).
	X520AttributesLDAPURL X520Attributes = "LDAPURL"
	// X520AttributesTagLocation is the "tagLocation" attribute (OID 2.5.4.96).
	X520AttributesTagLocation X520Attributes = "TAGLOCATION"
	// X520AttributesOrganizationIdentifier is the "organizationIdentifier" attribute (OID 2.5.4.97).
	X520AttributesOrganizationIdentifier X520Attributes = "ORGANIZATIONIDENTIFIER"
	// X520AttributesCountryCode3C is the "countryCode3c" attribute (OID 2.5.4.98).
	X520AttributesCountryCode3C X520Attributes = "COUNTRYCODE3C"
	// X520AttributesCountryCode3N is the "countryCode3n" attribute (OID 2.5.4.99).
	X520AttributesCountryCode3N X520Attributes = "COUNTRYCODE3N"
	// X520AttributesDNSName is the "dnsName" attribute (OID 2.5.4.100).
	X520AttributesDNSName X520Attributes = "DNSNAME"
	// X520AttributesEEPKCertificatRevocationList is the "eepkCertificatRevocationList" attribute (OID 2.5.4.101).
	X520AttributesEEPKCertificatRevocationList X520Attributes = "EEPKCERTIFICATREVOCATIONLIST"
	// X520AttributesEEAttrCertificateRevocationList is the "eeAttrCertificateRevocationList" attribute (OID 2.5.4.102).
	X520AttributesEEAttrCertificateRevocationList X520Attributes = "EEATTRCERTIFICATEREVOCATIONLIST"
	// X520AttributesUserPwdDescription is the "userPwdDescription" attribute (OID 2.5.40.0).
	X520AttributesUserPwdDescription X520Attributes = "USERPWDDESCRIPTION"
	// X520AttributesPwdVocabularyDescription is the "pwdVocabularyDescription" attribute (OID 2.5.40.1).
	X520AttributesPwdVocabularyDescription X520Attributes = "PWDVOCABULARYDESCRIPTION"
	// X520AttributesPwdAlphabetDescription is the "pwdAlphabetDescription" attribute (OID 2.5.40.2).
	X520AttributesPwdAlphabetDescription X520Attributes = "PWDALPHABETDESCRIPTION"
	// X520AttributesPwdEncAlgDescription is the "pwdEncAlgDescription" attribute (OID 2.5.40.3).
	X520AttributesPwdEncAlgDescription X520Attributes = "PWDENCALGDESCRIPTION"
	// X520AttributesUTMCoords is the "utmCoords" attribute (OID 2.5.40.4).
	X520AttributesUTMCoords X520Attributes = "UTMCOORDS"
	// X520AttributesUIIForm is the "uiiForm" attribute (OID 2.5.40.5).
	X520AttributesUIIForm X520Attributes = "UIIFORM"
	// X520AttributesEPCForm is the "epcForm" attribute (OID 2.5.40.6).
	X520AttributesEPCForm X520Attributes = "EPCFORM"
	// X520AttributesCountryString3C is the "countryString3c" attribute (OID 2.5.40.7).
	X520AttributesCountryString3C X520Attributes = "COUNTRYSTRING3C"
	// X520AttributesCountryString3N is the "countryString3n" attribute (OID 2.5.40.8).
	X520AttributesCountryString3N X520Attributes = "COUNTRYSTRING3N"
	// X520AttributesDNSString is the "dnsString" attribute (OID 2.5.40.9).
	X520AttributesDNSString X520Attributes = "DNSSTRING"
	// X520AttributesAttributeTypeDescription is the "attributeTypeDescription" attribute (OID 1.3.6.1.4.1.1466.115.121.1.3).
	X520AttributesAttributeTypeDescription X520Attributes = "ATTRIBUTETYPEDESCRIPTION"
	// X520AttributesBitString is the "bitString" attribute (OID 1.3.6.1.4.1.1466.115.121.1.6).
	X520AttributesBitString X520Attributes = "BITSTRING"
	// X520AttributesBoolean is the "boolean" attribute (OID 1.3.6.1.4.1.1466.115.121.1.7).
	X520AttributesBoolean X520Attributes = "BOOLEAN"
	// X520AttributesX509Certificate is the "x509Certificate" attribute (OID 1.3.6.1.4.1.1466.115.121.1.8).
	X520AttributesX509Certificate X520Attributes = "X509CERTIFICATE"
	// X520AttributesX509CertificateList is the "x509CertificateList" attribute (OID 1.3.6.1.4.1.1466.115.121.1.9).
	X520AttributesX509CertificateList X520Attributes = "X509CERTIFICATELIST"
	// X520AttributesX509CertificatePair is the "x509CertificatePair" attribute (OID 1.3.6.1.4.1.1466.115.121.1.10).
	X520AttributesX509CertificatePair X520Attributes = "X509CERTIFICATEPAIR"
	// X520AttributesCountryString is the "countryString" attribute (OID 1.3.6.1.4.1.1466.115.121.1.11).
	X520AttributesCountryString X520Attributes = "COUNTRYSTRING"
	// X520AttributesDN is the "dn" attribute (OID 1.3.6.1.4.1.1466.115.121.1.12).
	X520AttributesDN X520Attributes = "DN"
	// X520AttributesDeliveryMethod is the "deliveryMethod" attribute (OID 1.3.6.1.4.1.1466.115.121.1.14).
	X520AttributesDeliveryMethod X520Attributes = "DELIVERYMETHOD"
	// X520AttributesDirectoryString is the "directoryString" attribute (OID 1.3.6.1.4.1.1466.115.121.1.15).
	X520AttributesDirectoryString X520Attributes = "DIRECTORYSTRING"
	// X520AttributesDITContentRuleDescription is the "dITContentRuleDescription" attribute (OID 1.3.6.1.4.1.1466.115.121.1.16).
	X520AttributesDITContentRuleDescription X520Attributes = "DITCONTENTRULEDESCRIPTION"
	// X520AttributesDITStructureRuleDescription is the "dITStructureRuleDescription" attribute (OID 1.3.6.1.4.1.1466.115.121.1.17).
	X520AttributesDITStructureRuleDescription X520Attributes = "DITSTRUCTURERULEDESCRIPTION"
	// X520AttributesEnhancedGuide is the "enhancedGuide" attribute (OID 1.3.6.1.4.1.1466.115.121.1.21).
	X520AttributesEnhancedGuide X520Attributes = "ENHANCEDGUIDE"
	// X520AttributesFacsimileTelephoneNr is the "facsimileTelephoneNr" attribute (OID 1.3.6.1.4.1.1466.115.121.1.22).
	X520AttributesFacsimileTelephoneNr X520Attributes = "FACSIMILETELEPHONENR"
	// X520AttributesFax is the "fax" attribute (OID 1.3.6.1.4.1.1466.115.121.1.23).
	X520AttributesFax X520Attributes = "FAX"
	// X520AttributesGeneralizedTime is the "generalizedTime" attribute (OID 1.3.6.1.4.1.1466.115.121.1.24).
	X520AttributesGeneralizedTime X520Attributes = "GENERALIZEDTIME"
	// X520AttributesGuide is the "guide" attribute (OID 1.3.6.1.4.1.1466.115.121.1.25).
	X520AttributesGuide X520Attributes = "GUIDE"
	// X520AttributesIA5String is the "ia5String" attribute (OID 1.3.6.1.4.1.1466.115.121.1.26).
	X520AttributesIA5String X520Attributes = "IA5STRING"
	// X520AttributesInteger is the "integer" attribute (OID 1.3.6.1.4.1.1466.115.121.1.27).
	X520AttributesInteger X520Attributes = "INTEGER"
	// X520AttributesJPEG is the "jpeg" attribute (OID 1.3.6.1.4.1.1466.115.121.1.28).
	X520AttributesJPEG X520Attributes = "JPEG"
	// X520AttributesMatchingRuleDescription is the "matchingRuleDescription" attribute (OID 1.3.6.1.4.1.1466.115.121.1.30).
	X520AttributesMatchingRuleDescription X520Attributes = "MATCHINGRULEDESCRIPTION"
	// X520AttributesMatchingRuleUseDescription is the "matchingRuleUseDescription" attribute (OID 1.3.6.1.4.1.1466.115.121.1.31).
	X520AttributesMatchingRuleUseDescription X520Attributes = "MATCHINGRULEUSEDESCRIPTION"
	// X520AttributesNameAndOptionalUID is the "nameAndOptionalUID" attribute (OID 1.3.6.1.4.1.1466.115.121.1.34).
	X520AttributesNameAndOptionalUID X520Attributes = "NAMEANDOPTIONALUID"
	// X520AttributesNameFormDescription is the "nameFormDescription" attribute (OID 1.3.6.1.4.1.1466.115.121.1.35).
	X520AttributesNameFormDescription X520Attributes = "NAMEFORMDESCRIPTION"
	// X520AttributesNumericString is the "numericString" attribute (OID 1.3.6.1.4.1.1466.115.121.1.36).
	X520AttributesNumericString X520Attributes = "NUMERICSTRING"
	// X520AttributesObjectClassDescription is the "objectClassDescription" attribute (OID 1.3.6.1.4.1.1466.115.121.1.37).
	X520AttributesObjectClassDescription X520Attributes = "OBJECTCLASSDESCRIPTION"
	// X520AttributesOID is the "oid" attribute (OID 1.3.6.1.4.1.1466.115.121.1.38).
	X520AttributesOID X520Attributes = "OID"
	// X520AttributesOtherMailbox is the "otherMailbox" attribute (OID 1.3.6.1.4.1.1466.115.121.1.39).
	X520AttributesOtherMailbox X520Attributes = "OTHERMAILBOX"
	// X520AttributesOctetString is the "octetString" attribute (OID 1.3.6.1.4.1.1466.115.121.1.40).
	X520AttributesOctetString X520Attributes = "OCTETSTRING"
	// X520AttributesPostalAddr is the "postalAddr" attribute (OID 1.3.6.1.4.1.1466.115.121.1.41).
	X520AttributesPostalAddr X520Attributes = "POSTALADDR"
	// X520AttributesPresentationAddr is the "presentationAddr" attribute (OID 1.3.6.1.4.1.1466.115.121.1.43).
	X520AttributesPresentationAddr X520Attributes = "PRESENTATIONADDR"
	// X520AttributesPrintableString is the "printableString" attribute (OID 1.3.6.1.4.1.1466.115.121.1.44).
	X520AttributesPrintableString X520Attributes = "PRINTABLESTRING"
	// X520AttributesSubtreeSpec is the "subtreeSpec" attribute (OID 1.3.6.1.4.1.1466.115.121.1.45).
	X520AttributesSubtreeSpec X520Attributes = "SUBTREESPEC"
	// X520AttributesX509SupportedAlgorithm is the "x509SupportedAlgorithm" attribute (OID 1.3.6.1.4.1.1466.115.121.1.49).
	X520AttributesX509SupportedAlgorithm X520Attributes = "X509SUPPORTEDALGORITHM"
	// X520AttributesTelephoneNr is the "telephoneNr" attribute (OID 1.3.6.1.4.1.1466.115.121.1.50).
	X520AttributesTelephoneNr X520Attributes = "TELEPHONENR"
	// X520AttributesTelexNr is the "telexNr" attribute (OID 1.3.6.1.4.1.1466.115.121.1.52).
	X520AttributesTelexNr X520Attributes = "TELEXNR"
	// X520AttributesUTCTime is the "utcTime" attribute (OID 1.3.6.1.4.1.1466.115.121.1.53).
	X520AttributesUTCTime X520Attributes = "UTCTIME"
	// X520AttributesLDAPSyntaxDescription is the "ldapSyntaxDescription" attribute (OID 1.3.6.1.4.1.1466.115.121.1.54).
	X520AttributesLDAPSyntaxDescription X520Attributes = "LDAPSYNTAXDESCRIPTION"
	// X520AttributesSubstringAssertion is the "substringAssertion" attribute (OID 1.3.6.1.4.1.1466.115.121.1.58).
	X520AttributesSubstringAssertion X520Attributes = "SUBSTRINGASSERTION"
	// X520AttributesEmailAddress is the "emailAddress" attribute (OID 1.2.840.113549.1.9.1).
	X520AttributesEmailAddress X520Attributes = "EMAIL_ADDRESS"
)

var x520AttributesData = map[X520Attributes]x520AttributesFields{
	X520AttributesObjectClass:                                   {"objectClass", "2.5.4.0"},
	X520AttributesAliasedEntryName:                              {"aliasedEntryName", "2.5.4.1"},
	X520AttributesEncryptedAliasedEntryName:                     {"encryptedAliasedEntryName", "2.5.4.1.2"},
	X520AttributesKnowledgeInformation:                          {"knowledgeInformation", "2.5.4.2"},
	X520AttributesCommonName:                                    {"commonName", "2.5.4.3"},
	X520AttributesEncryptedCommonName:                           {"encryptedCommonName", "2.5.4.3.2"},
	X520AttributesSurname:                                       {"surname", "2.5.4.4"},
	X520AttributesEncryptedSurname:                              {"encryptedSurname", "2.5.4.4.2"},
	X520AttributesSerialNumber:                                  {"serialNumber", "2.5.4.5"},
	X520AttributesEncryptedSerialNumber:                         {"encryptedSerialNumber", "2.5.4.5.2"},
	X520AttributesCountryName:                                   {"countryName", "2.5.4.6"},
	X520AttributesEncryptedCountryName:                          {"encryptedCountryName", "2.5.4.6.2"},
	X520AttributesLocalityName:                                  {"localityName", "2.5.4.7"},
	X520AttributesEncryptedLocalityName:                         {"encryptedLocalityName", "2.5.4.7.2"},
	X520AttributesCollectiveLocalityName:                        {"collectiveLocalityName", "2.5.4.7.1"},
	X520AttributesEncryptedCollectiveLocalityName:               {"encryptedCollectiveLocalityName", "2.5.4.7.1.2"},
	X520AttributesStateOrProvinceName:                           {"stateOrProvinceName", "2.5.4.8"},
	X520AttributesEncryptedStateOrProvinceName:                  {"encryptedStateOrProvinceName", "2.5.4.8.2"},
	X520AttributesCollectiveStateOrProvinceName:                 {"collectiveStateOrProvinceName", "2.5.4.8.1"},
	X520AttributesEncryptedCollectiveStateOrProvinceName:        {"encryptedCollectiveStateOrProvinceName", "2.5.4.8.1.2"},
	X520AttributesStreetAddress:                                 {"streetAddress", "2.5.4.9"},
	X520AttributesEncryptedStreetAddress:                        {"encryptedStreetAddress", "2.5.4.9.2"},
	X520AttributesCollectiveStreetAddress:                       {"collectiveStreetAddress", "2.5.4.9.1"},
	X520AttributesEncryptedCollectiveStreetAddress:              {"encryptedCollectiveStreetAddress", "2.5.4.9.1.2"},
	X520AttributesOrganizationName:                              {"organizationName", "2.5.4.10"},
	X520AttributesEncryptedOrganizationName:                     {"encryptedOrganizationName", "2.5.4.10.2"},
	X520AttributesCollectiveOrganizationName:                    {"collectiveOrganizationName", "2.5.4.10.1"},
	X520AttributesEncryptedCollectiveOrganizationName:           {"encryptedCollectiveOrganizationName", "2.5.4.10.1.2"},
	X520AttributesOrganizationalUnitName:                        {"organizationalUnitName", "2.5.4.11"},
	X520AttributesEncryptedOrganizationalUnitName:               {"encryptedOrganizationalUnitName", "2.5.4.11.2"},
	X520AttributesCollectiveOrganizationalUnitName:              {"collectiveOrganizationalUnitName", "2.5.4.11.1"},
	X520AttributesEncryptedCollectiveOrganizationalUnitNam:      {"encryptedCollectiveOrganizationalUnitNam", "2.5.4.11.1.2"},
	X520AttributesTitle:                                         {"title", "2.5.4.12"},
	X520AttributesEncryptedTitle:                                {"encryptedTitle", "2.5.4.12.2"},
	X520AttributesDescription:                                   {"description", "2.5.4.13"},
	X520AttributesEncryptedDescription:                          {"encryptedDescription", "2.5.4.13.2"},
	X520AttributesSearchGuide:                                   {"searchGuide", "2.5.4.14"},
	X520AttributesEncryptedSearchGuide:                          {"encryptedSearchGuide", "2.5.4.14.2"},
	X520AttributesBusinessCategory:                              {"businessCategory", "2.5.4.15"},
	X520AttributesEncryptedBusinessCategory:                     {"encryptedBusinessCategory", "2.5.4.15.2"},
	X520AttributesPostalAddress:                                 {"postalAddress", "2.5.4.16"},
	X520AttributesEncryptedPostalAddress:                        {"encryptedPostalAddress", "2.5.4.16.2"},
	X520AttributesCollectivePostalAddress:                       {"collectivePostalAddress", "2.5.4.16.1"},
	X520AttributesEncryptedCollectivePostalAddress:              {"encryptedCollectivePostalAddress", "2.5.4.16.1.2"},
	X520AttributesPostalCode:                                    {"postalCode", "2.5.4.17"},
	X520AttributesEncryptedPostalCode:                           {"encryptedPostalCode", "2.5.4.17.2"},
	X520AttributesCollectivePostalCode:                          {"collectivePostalCode", "2.5.4.17.1"},
	X520AttributesEncryptedCollectivePostalCode:                 {"encryptedCollectivePostalCode", "2.5.4.17.1.2"},
	X520AttributesPostOfficeBox:                                 {"postOfficeBox", "2.5.4.18"},
	X520AttributesCollectivePostOfficeBox:                       {"collectivePostOfficeBox", "2.5.4.18.1"},
	X520AttributesEncryptedPostOfficeBox:                        {"encryptedPostOfficeBox", "2.5.4.18.2"},
	X520AttributesEncryptedCollectivePostOfficeBox:              {"encryptedCollectivePostOfficeBox", "2.5.4.18.1.2"},
	X520AttributesPhysicalDeliveryOfficeName:                    {"physicalDeliveryOfficeName", "2.5.4.19"},
	X520AttributesCollectivePhysicalDeliveryOfficeName:          {"collectivePhysicalDeliveryOfficeName", "2.5.4.19.1"},
	X520AttributesEncryptedPhysicalDeliveryOfficeName:           {"encryptedPhysicalDeliveryOfficeName", "2.5.4.19.2"},
	X520AttributesEncryptedCollectivePhysicalDeliveryOfficeName: {"encryptedCollectivePhysicalDeliveryOfficeName", "2.5.4.19.1.2"},
	X520AttributesTelephoneNumber:                               {"telephoneNumber", "2.5.4.20"},
	X520AttributesEncryptedTelephoneNumber:                      {"encryptedTelephoneNumber", "2.5.4.20.2"},
	X520AttributesCollectiveTelephoneNumber:                     {"collectiveTelephoneNumber", "2.5.4.20.1"},
	X520AttributesEncryptedCollectiveTelephoneNumber:            {"encryptedCollectiveTelephoneNumber", "2.5.4.20.1.2"},
	X520AttributesTelexNumber:                                   {"telexNumber", "2.5.4.21"},
	X520AttributesEncryptedTelexNumber:                          {"encryptedTelexNumber", "2.5.4.21.2"},
	X520AttributesCollectiveTelexNumber:                         {"collectiveTelexNumber", "2.5.4.21.1"},
	X520AttributesEncryptedCollectiveTelexNumber:                {"encryptedCollectiveTelexNumber", "2.5.4.21.1.2"},
	X520AttributesTeletexTerminalIdentifier:                     {"teletexTerminalIdentifier", "2.5.4.22"},
	X520AttributesEncryptedTeletexTerminalIdentifier:            {"encryptedTeletexTerminalIdentifier", "2.5.4.22.2"},
	X520AttributesCollectiveTeletexTerminalIdentifier:           {"collectiveTeletexTerminalIdentifier", "2.5.4.22.1"},
	X520AttributesEncryptedCollectiveTeletexTerminalIdentifier:  {"encryptedCollectiveTeletexTerminalIdentifier", "2.5.4.22.1.2"},
	X520AttributesFacsimileTelephoneNumber:                      {"facsimileTelephoneNumber", "2.5.4.23"},
	X520AttributesEncryptedFacsimileTelephoneNumber:             {"encryptedFacsimileTelephoneNumber", "2.5.4.23.2"},
	X520AttributesCollectiveFacsimileTelephoneNumber:            {"collectiveFacsimileTelephoneNumber", "2.5.4.23.1"},
	X520AttributesEncryptedCollectiveFacsimileTelephoneNumber:   {"encryptedCollectiveFacsimileTelephoneNumber", "2.5.4.23.1.2"},
	X520AttributesX121Address:                                   {"x121Address", "2.5.4.24"},
	X520AttributesEncryptedX121Address:                          {"encryptedX121Address", "2.5.4.24.2"},
	X520AttributesInternationalISDNNumber:                       {"internationalISDNNumber", "2.5.4.25"},
	X520AttributesEncryptedInternationalISDNNumber:              {"encryptedInternationalISDNNumber", "2.5.4.25.2"},
	X520AttributesCollectiveInternationalISDNNumber:             {"collectiveInternationalISDNNumber", "2.5.4.25.1"},
	X520AttributesEncryptedCollectiveInternationalISDNNumber:    {"encryptedCollectiveInternationalISDNNumber", "2.5.4.25.1.2"},
	X520AttributesRegisteredAddress:                             {"registeredAddress", "2.5.4.26"},
	X520AttributesEncryptedRegisteredAddress:                    {"encryptedRegisteredAddress", "2.5.4.26.2"},
	X520AttributesDestinationIndicator:                          {"destinationIndicator", "2.5.4.27"},
	X520AttributesEncryptedDestinationIndicator:                 {"encryptedDestinationIndicator", "2.5.4.27.2"},
	X520AttributesPreferredDeliveryMethod:                       {"preferredDeliveryMethod", "2.5.4.28"},
	X520AttributesEncryptedPreferredDeliveryMethod:              {"encryptedPreferredDeliveryMethod", "2.5.4.28.2"},
	X520AttributesPresentationAddress:                           {"presentationAddress", "2.5.4.29"},
	X520AttributesEncryptedPresentationAddress:                  {"encryptedPresentationAddress", "2.5.4.29.2"},
	X520AttributesSupportedApplicationContext:                   {"supportedApplicationContext", "2.5.4.30"},
	X520AttributesEncryptedSupportedApplicationContext:          {"encryptedSupportedApplicationContext", "2.5.4.30.2"},
	X520AttributesMember:                                        {"member", "2.5.4.31"},
	X520AttributesEncryptedMember:                               {"encryptedMember", "2.5.4.31.2"},
	X520AttributesOwner:                                         {"owner", "2.5.4.32"},
	X520AttributesEncryptedOwner:                                {"encryptedOwner", "2.5.4.32.2"},
	X520AttributesRoleOccupant:                                  {"roleOccupant", "2.5.4.33"},
	X520AttributesEncryptedRoleOccupant:                         {"encryptedRoleOccupant", "2.5.4.33.2"},
	X520AttributesSeeAlso:                                       {"seeAlso", "2.5.4.34"},
	X520AttributesEncryptedSeeAlso:                              {"encryptedSeeAlso", "2.5.4.34.2"},
	X520AttributesUserPassword:                                  {"userPassword", "2.5.4.35"},
	X520AttributesEncryptedUserPassword:                         {"encryptedUserPassword", "2.5.4.35.2"},
	X520AttributesUserCertificate:                               {"userCertificate", "2.5.4.36"},
	X520AttributesEncryptedUserCertificate:                      {"encryptedUserCertificate", "2.5.4.36.2"},
	X520AttributesCACertificate:                                 {"cACertificate", "2.5.4.37"},
	X520AttributesEncryptedCACertificate:                        {"encryptedCACertificate", "2.5.4.37.2"},
	X520AttributesAuthorityRevocationList:                       {"authorityRevocationList", "2.5.4.38"},
	X520AttributesEncryptedAuthorityRevocationList:              {"encryptedAuthorityRevocationList", "2.5.4.38.2"},
	X520AttributesCertificateRevocationList:                     {"certificateRevocationList", "2.5.4.39"},
	X520AttributesEncryptedCertificateRevocationList:            {"encryptedCertificateRevocationList", "2.5.4.39.2"},
	X520AttributesCrossCertificatePair:                          {"crossCertificatePair", "2.5.4.40"},
	X520AttributesEncryptedCrossCertificatePair:                 {"encryptedCrossCertificatePair", "2.5.4.40.2"},
	X520AttributesName:                                          {"name", "2.5.4.41"},
	X520AttributesGivenName:                                     {"givenName", "2.5.4.42"},
	X520AttributesEncryptedGivenName:                            {"encryptedGivenName", "2.5.4.42.2"},
	X520AttributesInitials:                                      {"initials", "2.5.4.43"},
	X520AttributesEncryptedInitials:                             {"encryptedInitials", "2.5.4.43.2"},
	X520AttributesGenerationQualifier:                           {"generationQualifier", "2.5.4.44"},
	X520AttributesEncryptedGenerationQualifier:                  {"encryptedGenerationQualifier", "2.5.4.44.2"},
	X520AttributesUniqueIdentifier:                              {"uniqueIdentifier", "2.5.4.45"},
	X520AttributesEncryptedUniqueIdentifier:                     {"encryptedUniqueIdentifier", "2.5.4.45.2"},
	X520AttributesDNQualifier:                                   {"dnQualifier", "2.5.4.46"},
	X520AttributesEncryptedDNQualifier:                          {"encryptedDnQualifier", "2.5.4.46.2"},
	X520AttributesEnhancedSearchGuide:                           {"enhancedSearchGuide", "2.5.4.47"},
	X520AttributesEncryptedEnhancedSearchGuide:                  {"encryptedEnhancedSearchGuide", "2.5.4.47.2"},
	X520AttributesProtocolInformation:                           {"protocolInformation", "2.5.4.48"},
	X520AttributesEncryptedProtocolInformation:                  {"encryptedProtocolInformation", "2.5.4.48.2"},
	X520AttributesDistinguishedName:                             {"distinguishedName", "2.5.4.49"},
	X520AttributesEncryptedDistinguishedName:                    {"encryptedDistinguishedName", "2.5.4.49.2"},
	X520AttributesUniqueMember:                                  {"uniqueMember", "2.5.4.50"},
	X520AttributesEncryptedUniqueMember:                         {"encryptedUniqueMember", "2.5.4.50.2"},
	X520AttributesHouseIdentifier:                               {"houseIdentifier", "2.5.4.51"},
	X520AttributesEncryptedHouseIdentifier:                      {"encryptedHouseIdentifier", "2.5.4.51.2"},
	X520AttributesSupportedAlgorithms:                           {"supportedAlgorithms", "2.5.4.52"},
	X520AttributesEncryptedSupportedAlgorithms:                  {"encryptedSupportedAlgorithms", "2.5.4.52.2"},
	X520AttributesDeltaRevocationList:                           {"deltaRevocationList", "2.5.4.53"},
	X520AttributesEncryptedDeltaRevocationList:                  {"encryptedDeltaRevocationList", "2.5.4.53.2"},
	X520AttributesDMDName:                                       {"dmdName", "2.5.4.54"},
	X520AttributesEncryptedDMDName:                              {"encryptedDmdName", "2.5.4.54.2"},
	X520AttributesClearance:                                     {"clearance", "2.5.4.55"},
	X520AttributesEncryptedClearance:                            {"encryptedClearance", "2.5.4.55.2"},
	X520AttributesDefaultDirQOP:                                 {"defaultDirQop", "2.5.4.56"},
	X520AttributesEncryptedDefaultDirQOP:                        {"encryptedDefaultDirQop", "2.5.4.56.2"},
	X520AttributesAttributeIntegrityInfo:                        {"attributeIntegrityInfo", "2.5.4.57"},
	X520AttributesEncryptedAttributeIntegrityInfo:               {"encryptedAttributeIntegrityInfo", "2.5.4.57.2"},
	X520AttributesAttributeCertificate:                          {"attributeCertificate", "2.5.4.58"},
	X520AttributesEncryptedAttributeCertificate:                 {"encryptedAttributeCertificate", "2.5.4.58.2"},
	X520AttributesAttributeCertificateRevocationList:            {"attributeCertificateRevocationList", "2.5.4.59"},
	X520AttributesEncryptedAttributeCertificateRevocationList:   {"encryptedAttributeCertificateRevocationList", "2.5.4.59.2"},
	X520AttributesConfKeyInfo:                                   {"confKeyInfo", "2.5.4.60"},
	X520AttributesEncryptedConfKeyInfo:                          {"encryptedConfKeyInfo", "2.5.4.60.2"},
	X520AttributesAACertificate:                                 {"aACertificate", "2.5.4.61"},
	X520AttributesAttributeDescriptorCertificate:                {"attributeDescriptorCertificate", "2.5.4.62"},
	X520AttributesAttributeAuthorityRevocationList:              {"attributeAuthorityRevocationList", "2.5.4.63"},
	X520AttributesFamilyInformation:                             {"family-information", "2.5.4.64"},
	X520AttributesPseudonym:                                     {"pseudonym", "2.5.4.65"},
	X520AttributesCommunicationsService:                         {"communicationsService", "2.5.4.66"},
	X520AttributesCommunicationsNetwork:                         {"communicationsNetwork", "2.5.4.67"},
	X520AttributesCertificationPracticeStmt:                     {"certificationPracticeStmt", "2.5.4.68"},
	X520AttributesCertificatePolicy:                             {"certificatePolicy", "2.5.4.69"},
	X520AttributesPKIPath:                                       {"pkiPath", "2.5.4.70"},
	X520AttributesPrivPolicy:                                    {"privPolicy", "2.5.4.71"},
	X520AttributesRole:                                          {"role", "2.5.4.72"},
	X520AttributesDelegationPath:                                {"delegationPath", "2.5.4.73"},
	X520AttributesProtPrivPolicy:                                {"protPrivPolicy", "2.5.4.74"},
	X520AttributesXMLPrivilegeInfo:                              {"xMLPrivilegeInfo", "2.5.4.75"},
	X520AttributesXMLPrivPolicy:                                 {"xmlPrivPolicy", "2.5.4.76"},
	X520AttributesUUIDPair:                                      {"uuidpair", "2.5.4.77"},
	X520AttributesTagOID:                                        {"tagOid", "2.5.4.78"},
	X520AttributesUIIFormat:                                     {"uiiFormat", "2.5.4.79"},
	X520AttributesUIIInURN:                                      {"uiiInUrn", "2.5.4.80"},
	X520AttributesContentURL:                                    {"contentUrl", "2.5.4.81"},
	X520AttributesPermission:                                    {"permission", "2.5.4.82"},
	X520AttributesURI:                                           {"uri", "2.5.4.83"},
	X520AttributesPwdAttribute:                                  {"pwdAttribute", "2.5.4.84"},
	X520AttributesUserPwd:                                       {"userPwd", "2.5.4.85"},
	X520AttributesURN:                                           {"urn", "2.5.4.86"},
	X520AttributesURL:                                           {"url", "2.5.4.87"},
	X520AttributesUTMCoordinates:                                {"utmCoordinates", "2.5.4.88"},
	X520AttributesURNC:                                          {"urnC", "2.5.4.89"},
	X520AttributesUII:                                           {"uii", "2.5.4.90"},
	X520AttributesEPC:                                           {"epc", "2.5.4.91"},
	X520AttributesTagAFI:                                        {"tagAfi", "2.5.4.92"},
	X520AttributesEPCFormat:                                     {"epcFormat", "2.5.4.93"},
	X520AttributesEPCInURN:                                      {"epcInUrn", "2.5.4.94"},
	X520AttributesLDAPURL:                                       {"ldapUrl", "2.5.4.95"},
	X520AttributesTagLocation:                                   {"tagLocation", "2.5.4.96"},
	X520AttributesOrganizationIdentifier:                        {"organizationIdentifier", "2.5.4.97"},
	X520AttributesCountryCode3C:                                 {"countryCode3c", "2.5.4.98"},
	X520AttributesCountryCode3N:                                 {"countryCode3n", "2.5.4.99"},
	X520AttributesDNSName:                                       {"dnsName", "2.5.4.100"},
	X520AttributesEEPKCertificatRevocationList:                  {"eepkCertificatRevocationList", "2.5.4.101"},
	X520AttributesEEAttrCertificateRevocationList:               {"eeAttrCertificateRevocationList", "2.5.4.102"},
	X520AttributesUserPwdDescription:                            {"userPwdDescription", "2.5.40.0"},
	X520AttributesPwdVocabularyDescription:                      {"pwdVocabularyDescription", "2.5.40.1"},
	X520AttributesPwdAlphabetDescription:                        {"pwdAlphabetDescription", "2.5.40.2"},
	X520AttributesPwdEncAlgDescription:                          {"pwdEncAlgDescription", "2.5.40.3"},
	X520AttributesUTMCoords:                                     {"utmCoords", "2.5.40.4"},
	X520AttributesUIIForm:                                       {"uiiForm", "2.5.40.5"},
	X520AttributesEPCForm:                                       {"epcForm", "2.5.40.6"},
	X520AttributesCountryString3C:                               {"countryString3c", "2.5.40.7"},
	X520AttributesCountryString3N:                               {"countryString3n", "2.5.40.8"},
	X520AttributesDNSString:                                     {"dnsString", "2.5.40.9"},
	X520AttributesAttributeTypeDescription:                      {"attributeTypeDescription", "1.3.6.1.4.1.1466.115.121.1.3"},
	X520AttributesBitString:                                     {"bitString", "1.3.6.1.4.1.1466.115.121.1.6"},
	X520AttributesBoolean:                                       {"boolean", "1.3.6.1.4.1.1466.115.121.1.7"},
	X520AttributesX509Certificate:                               {"x509Certificate", "1.3.6.1.4.1.1466.115.121.1.8"},
	X520AttributesX509CertificateList:                           {"x509CertificateList", "1.3.6.1.4.1.1466.115.121.1.9"},
	X520AttributesX509CertificatePair:                           {"x509CertificatePair", "1.3.6.1.4.1.1466.115.121.1.10"},
	X520AttributesCountryString:                                 {"countryString", "1.3.6.1.4.1.1466.115.121.1.11"},
	X520AttributesDN:                                            {"dn", "1.3.6.1.4.1.1466.115.121.1.12"},
	X520AttributesDeliveryMethod:                                {"deliveryMethod", "1.3.6.1.4.1.1466.115.121.1.14"},
	X520AttributesDirectoryString:                               {"directoryString", "1.3.6.1.4.1.1466.115.121.1.15"},
	X520AttributesDITContentRuleDescription:                     {"dITContentRuleDescription", "1.3.6.1.4.1.1466.115.121.1.16"},
	X520AttributesDITStructureRuleDescription:                   {"dITStructureRuleDescription", "1.3.6.1.4.1.1466.115.121.1.17"},
	X520AttributesEnhancedGuide:                                 {"enhancedGuide", "1.3.6.1.4.1.1466.115.121.1.21"},
	X520AttributesFacsimileTelephoneNr:                          {"facsimileTelephoneNr", "1.3.6.1.4.1.1466.115.121.1.22"},
	X520AttributesFax:                                           {"fax", "1.3.6.1.4.1.1466.115.121.1.23"},
	X520AttributesGeneralizedTime:                               {"generalizedTime", "1.3.6.1.4.1.1466.115.121.1.24"},
	X520AttributesGuide:                                         {"guide", "1.3.6.1.4.1.1466.115.121.1.25"},
	X520AttributesIA5String:                                     {"ia5String", "1.3.6.1.4.1.1466.115.121.1.26"},
	X520AttributesInteger:                                       {"integer", "1.3.6.1.4.1.1466.115.121.1.27"},
	X520AttributesJPEG:                                          {"jpeg", "1.3.6.1.4.1.1466.115.121.1.28"},
	X520AttributesMatchingRuleDescription:                       {"matchingRuleDescription", "1.3.6.1.4.1.1466.115.121.1.30"},
	X520AttributesMatchingRuleUseDescription:                    {"matchingRuleUseDescription", "1.3.6.1.4.1.1466.115.121.1.31"},
	X520AttributesNameAndOptionalUID:                            {"nameAndOptionalUID", "1.3.6.1.4.1.1466.115.121.1.34"},
	X520AttributesNameFormDescription:                           {"nameFormDescription", "1.3.6.1.4.1.1466.115.121.1.35"},
	X520AttributesNumericString:                                 {"numericString", "1.3.6.1.4.1.1466.115.121.1.36"},
	X520AttributesObjectClassDescription:                        {"objectClassDescription", "1.3.6.1.4.1.1466.115.121.1.37"},
	X520AttributesOID:                                           {"oid", "1.3.6.1.4.1.1466.115.121.1.38"},
	X520AttributesOtherMailbox:                                  {"otherMailbox", "1.3.6.1.4.1.1466.115.121.1.39"},
	X520AttributesOctetString:                                   {"octetString", "1.3.6.1.4.1.1466.115.121.1.40"},
	X520AttributesPostalAddr:                                    {"postalAddr", "1.3.6.1.4.1.1466.115.121.1.41"},
	X520AttributesPresentationAddr:                              {"presentationAddr", "1.3.6.1.4.1.1466.115.121.1.43"},
	X520AttributesPrintableString:                               {"printableString", "1.3.6.1.4.1.1466.115.121.1.44"},
	X520AttributesSubtreeSpec:                                   {"subtreeSpec", "1.3.6.1.4.1.1466.115.121.1.45"},
	X520AttributesX509SupportedAlgorithm:                        {"x509SupportedAlgorithm", "1.3.6.1.4.1.1466.115.121.1.49"},
	X520AttributesTelephoneNr:                                   {"telephoneNr", "1.3.6.1.4.1.1466.115.121.1.50"},
	X520AttributesTelexNr:                                       {"telexNr", "1.3.6.1.4.1.1466.115.121.1.52"},
	X520AttributesUTCTime:                                       {"utcTime", "1.3.6.1.4.1.1466.115.121.1.53"},
	X520AttributesLDAPSyntaxDescription:                         {"ldapSyntaxDescription", "1.3.6.1.4.1.1466.115.121.1.54"},
	X520AttributesSubstringAssertion:                            {"substringAssertion", "1.3.6.1.4.1.1466.115.121.1.58"},
	X520AttributesEmailAddress:                                  {"emailAddress", "1.2.840.113549.1.9.1"},
}

// X520AttributesValues returns all X520Attributes constants in declaration
// order.
func X520AttributesValues() []X520Attributes {
	return []X520Attributes{
		X520AttributesObjectClass,
		X520AttributesAliasedEntryName,
		X520AttributesEncryptedAliasedEntryName,
		X520AttributesKnowledgeInformation,
		X520AttributesCommonName,
		X520AttributesEncryptedCommonName,
		X520AttributesSurname,
		X520AttributesEncryptedSurname,
		X520AttributesSerialNumber,
		X520AttributesEncryptedSerialNumber,
		X520AttributesCountryName,
		X520AttributesEncryptedCountryName,
		X520AttributesLocalityName,
		X520AttributesEncryptedLocalityName,
		X520AttributesCollectiveLocalityName,
		X520AttributesEncryptedCollectiveLocalityName,
		X520AttributesStateOrProvinceName,
		X520AttributesEncryptedStateOrProvinceName,
		X520AttributesCollectiveStateOrProvinceName,
		X520AttributesEncryptedCollectiveStateOrProvinceName,
		X520AttributesStreetAddress,
		X520AttributesEncryptedStreetAddress,
		X520AttributesCollectiveStreetAddress,
		X520AttributesEncryptedCollectiveStreetAddress,
		X520AttributesOrganizationName,
		X520AttributesEncryptedOrganizationName,
		X520AttributesCollectiveOrganizationName,
		X520AttributesEncryptedCollectiveOrganizationName,
		X520AttributesOrganizationalUnitName,
		X520AttributesEncryptedOrganizationalUnitName,
		X520AttributesCollectiveOrganizationalUnitName,
		X520AttributesEncryptedCollectiveOrganizationalUnitNam,
		X520AttributesTitle,
		X520AttributesEncryptedTitle,
		X520AttributesDescription,
		X520AttributesEncryptedDescription,
		X520AttributesSearchGuide,
		X520AttributesEncryptedSearchGuide,
		X520AttributesBusinessCategory,
		X520AttributesEncryptedBusinessCategory,
		X520AttributesPostalAddress,
		X520AttributesEncryptedPostalAddress,
		X520AttributesCollectivePostalAddress,
		X520AttributesEncryptedCollectivePostalAddress,
		X520AttributesPostalCode,
		X520AttributesEncryptedPostalCode,
		X520AttributesCollectivePostalCode,
		X520AttributesEncryptedCollectivePostalCode,
		X520AttributesPostOfficeBox,
		X520AttributesCollectivePostOfficeBox,
		X520AttributesEncryptedPostOfficeBox,
		X520AttributesEncryptedCollectivePostOfficeBox,
		X520AttributesPhysicalDeliveryOfficeName,
		X520AttributesCollectivePhysicalDeliveryOfficeName,
		X520AttributesEncryptedPhysicalDeliveryOfficeName,
		X520AttributesEncryptedCollectivePhysicalDeliveryOfficeName,
		X520AttributesTelephoneNumber,
		X520AttributesEncryptedTelephoneNumber,
		X520AttributesCollectiveTelephoneNumber,
		X520AttributesEncryptedCollectiveTelephoneNumber,
		X520AttributesTelexNumber,
		X520AttributesEncryptedTelexNumber,
		X520AttributesCollectiveTelexNumber,
		X520AttributesEncryptedCollectiveTelexNumber,
		X520AttributesTeletexTerminalIdentifier,
		X520AttributesEncryptedTeletexTerminalIdentifier,
		X520AttributesCollectiveTeletexTerminalIdentifier,
		X520AttributesEncryptedCollectiveTeletexTerminalIdentifier,
		X520AttributesFacsimileTelephoneNumber,
		X520AttributesEncryptedFacsimileTelephoneNumber,
		X520AttributesCollectiveFacsimileTelephoneNumber,
		X520AttributesEncryptedCollectiveFacsimileTelephoneNumber,
		X520AttributesX121Address,
		X520AttributesEncryptedX121Address,
		X520AttributesInternationalISDNNumber,
		X520AttributesEncryptedInternationalISDNNumber,
		X520AttributesCollectiveInternationalISDNNumber,
		X520AttributesEncryptedCollectiveInternationalISDNNumber,
		X520AttributesRegisteredAddress,
		X520AttributesEncryptedRegisteredAddress,
		X520AttributesDestinationIndicator,
		X520AttributesEncryptedDestinationIndicator,
		X520AttributesPreferredDeliveryMethod,
		X520AttributesEncryptedPreferredDeliveryMethod,
		X520AttributesPresentationAddress,
		X520AttributesEncryptedPresentationAddress,
		X520AttributesSupportedApplicationContext,
		X520AttributesEncryptedSupportedApplicationContext,
		X520AttributesMember,
		X520AttributesEncryptedMember,
		X520AttributesOwner,
		X520AttributesEncryptedOwner,
		X520AttributesRoleOccupant,
		X520AttributesEncryptedRoleOccupant,
		X520AttributesSeeAlso,
		X520AttributesEncryptedSeeAlso,
		X520AttributesUserPassword,
		X520AttributesEncryptedUserPassword,
		X520AttributesUserCertificate,
		X520AttributesEncryptedUserCertificate,
		X520AttributesCACertificate,
		X520AttributesEncryptedCACertificate,
		X520AttributesAuthorityRevocationList,
		X520AttributesEncryptedAuthorityRevocationList,
		X520AttributesCertificateRevocationList,
		X520AttributesEncryptedCertificateRevocationList,
		X520AttributesCrossCertificatePair,
		X520AttributesEncryptedCrossCertificatePair,
		X520AttributesName,
		X520AttributesGivenName,
		X520AttributesEncryptedGivenName,
		X520AttributesInitials,
		X520AttributesEncryptedInitials,
		X520AttributesGenerationQualifier,
		X520AttributesEncryptedGenerationQualifier,
		X520AttributesUniqueIdentifier,
		X520AttributesEncryptedUniqueIdentifier,
		X520AttributesDNQualifier,
		X520AttributesEncryptedDNQualifier,
		X520AttributesEnhancedSearchGuide,
		X520AttributesEncryptedEnhancedSearchGuide,
		X520AttributesProtocolInformation,
		X520AttributesEncryptedProtocolInformation,
		X520AttributesDistinguishedName,
		X520AttributesEncryptedDistinguishedName,
		X520AttributesUniqueMember,
		X520AttributesEncryptedUniqueMember,
		X520AttributesHouseIdentifier,
		X520AttributesEncryptedHouseIdentifier,
		X520AttributesSupportedAlgorithms,
		X520AttributesEncryptedSupportedAlgorithms,
		X520AttributesDeltaRevocationList,
		X520AttributesEncryptedDeltaRevocationList,
		X520AttributesDMDName,
		X520AttributesEncryptedDMDName,
		X520AttributesClearance,
		X520AttributesEncryptedClearance,
		X520AttributesDefaultDirQOP,
		X520AttributesEncryptedDefaultDirQOP,
		X520AttributesAttributeIntegrityInfo,
		X520AttributesEncryptedAttributeIntegrityInfo,
		X520AttributesAttributeCertificate,
		X520AttributesEncryptedAttributeCertificate,
		X520AttributesAttributeCertificateRevocationList,
		X520AttributesEncryptedAttributeCertificateRevocationList,
		X520AttributesConfKeyInfo,
		X520AttributesEncryptedConfKeyInfo,
		X520AttributesAACertificate,
		X520AttributesAttributeDescriptorCertificate,
		X520AttributesAttributeAuthorityRevocationList,
		X520AttributesFamilyInformation,
		X520AttributesPseudonym,
		X520AttributesCommunicationsService,
		X520AttributesCommunicationsNetwork,
		X520AttributesCertificationPracticeStmt,
		X520AttributesCertificatePolicy,
		X520AttributesPKIPath,
		X520AttributesPrivPolicy,
		X520AttributesRole,
		X520AttributesDelegationPath,
		X520AttributesProtPrivPolicy,
		X520AttributesXMLPrivilegeInfo,
		X520AttributesXMLPrivPolicy,
		X520AttributesUUIDPair,
		X520AttributesTagOID,
		X520AttributesUIIFormat,
		X520AttributesUIIInURN,
		X520AttributesContentURL,
		X520AttributesPermission,
		X520AttributesURI,
		X520AttributesPwdAttribute,
		X520AttributesUserPwd,
		X520AttributesURN,
		X520AttributesURL,
		X520AttributesUTMCoordinates,
		X520AttributesURNC,
		X520AttributesUII,
		X520AttributesEPC,
		X520AttributesTagAFI,
		X520AttributesEPCFormat,
		X520AttributesEPCInURN,
		X520AttributesLDAPURL,
		X520AttributesTagLocation,
		X520AttributesOrganizationIdentifier,
		X520AttributesCountryCode3C,
		X520AttributesCountryCode3N,
		X520AttributesDNSName,
		X520AttributesEEPKCertificatRevocationList,
		X520AttributesEEAttrCertificateRevocationList,
		X520AttributesUserPwdDescription,
		X520AttributesPwdVocabularyDescription,
		X520AttributesPwdAlphabetDescription,
		X520AttributesPwdEncAlgDescription,
		X520AttributesUTMCoords,
		X520AttributesUIIForm,
		X520AttributesEPCForm,
		X520AttributesCountryString3C,
		X520AttributesCountryString3N,
		X520AttributesDNSString,
		X520AttributesAttributeTypeDescription,
		X520AttributesBitString,
		X520AttributesBoolean,
		X520AttributesX509Certificate,
		X520AttributesX509CertificateList,
		X520AttributesX509CertificatePair,
		X520AttributesCountryString,
		X520AttributesDN,
		X520AttributesDeliveryMethod,
		X520AttributesDirectoryString,
		X520AttributesDITContentRuleDescription,
		X520AttributesDITStructureRuleDescription,
		X520AttributesEnhancedGuide,
		X520AttributesFacsimileTelephoneNr,
		X520AttributesFax,
		X520AttributesGeneralizedTime,
		X520AttributesGuide,
		X520AttributesIA5String,
		X520AttributesInteger,
		X520AttributesJPEG,
		X520AttributesMatchingRuleDescription,
		X520AttributesMatchingRuleUseDescription,
		X520AttributesNameAndOptionalUID,
		X520AttributesNameFormDescription,
		X520AttributesNumericString,
		X520AttributesObjectClassDescription,
		X520AttributesOID,
		X520AttributesOtherMailbox,
		X520AttributesOctetString,
		X520AttributesPostalAddr,
		X520AttributesPresentationAddr,
		X520AttributesPrintableString,
		X520AttributesSubtreeSpec,
		X520AttributesX509SupportedAlgorithm,
		X520AttributesTelephoneNr,
		X520AttributesTelexNr,
		X520AttributesUTCTime,
		X520AttributesLDAPSyntaxDescription,
		X520AttributesSubstringAssertion,
		X520AttributesEmailAddress,
	}
}

// X520AttributesValueOf returns the X520Attributes matching the given Java enum name.
func X520AttributesValueOf(name string) (X520Attributes, error) {
	for _, v := range X520AttributesValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", fmt.Errorf("no enum constant X520Attributes.%s", name)
}

// OID returns the OID. Implements OidDescription.
func (x X520Attributes) OID() string {
	return x520AttributesData[x].oid
}

// Description returns the description of the attribute. Implements OidDescription.
func (x X520Attributes) Description() string {
	return x520AttributesData[x].description
}

// GetUppercaseDescriptionForOids returns a map of X520 Attribute uppercase
// names (the Java enum's name()) to their corresponding OIDs.
func GetUppercaseDescriptionForOids() map[string]string {
	m := make(map[string]string, len(x520AttributesData))
	for _, v := range X520AttributesValues() {
		m[string(v)] = v.OID()
	}
	return m
}

// GetOidDescriptions returns a map of X520 Attribute OIDs to their descriptions.
func GetOidDescriptions() map[string]string {
	m := make(map[string]string, len(x520AttributesData))
	for _, v := range X520AttributesValues() {
		m[v.OID()] = v.Description()
	}
	return m
}

// compile-time interface assertion.
var _ OidDescription = X520Attributes("")
