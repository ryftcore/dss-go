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
	// X520Attributes_OBJECTCLASS is the "objectClass" attribute (OID 2.5.4.0).
	X520Attributes_OBJECTCLASS X520Attributes = "OBJECTCLASS"
	// X520Attributes_ALIASEDENTRYNAME is the "aliasedEntryName" attribute (OID 2.5.4.1).
	X520Attributes_ALIASEDENTRYNAME X520Attributes = "ALIASEDENTRYNAME"
	// X520Attributes_ENCRYPTEDALIASEDENTRYNAME is the "encryptedAliasedEntryName" attribute (OID 2.5.4.1.2).
	X520Attributes_ENCRYPTEDALIASEDENTRYNAME X520Attributes = "ENCRYPTEDALIASEDENTRYNAME"
	// X520Attributes_KNOWLEDGEINFORMATION is the "knowledgeInformation" attribute (OID 2.5.4.2).
	X520Attributes_KNOWLEDGEINFORMATION X520Attributes = "KNOWLEDGEINFORMATION"
	// X520Attributes_COMMONNAME is the "commonName" attribute (OID 2.5.4.3).
	X520Attributes_COMMONNAME X520Attributes = "COMMONNAME"
	// X520Attributes_ENCRYPTEDCOMMONNAME is the "encryptedCommonName" attribute (OID 2.5.4.3.2).
	X520Attributes_ENCRYPTEDCOMMONNAME X520Attributes = "ENCRYPTEDCOMMONNAME"
	// X520Attributes_SURNAME is the "surname" attribute (OID 2.5.4.4).
	X520Attributes_SURNAME X520Attributes = "SURNAME"
	// X520Attributes_ENCRYPTEDSURNAME is the "encryptedSurname" attribute (OID 2.5.4.4.2).
	X520Attributes_ENCRYPTEDSURNAME X520Attributes = "ENCRYPTEDSURNAME"
	// X520Attributes_SERIALNUMBER is the "serialNumber" attribute (OID 2.5.4.5).
	X520Attributes_SERIALNUMBER X520Attributes = "SERIALNUMBER"
	// X520Attributes_ENCRYPTEDSERIALNUMBER is the "encryptedSerialNumber" attribute (OID 2.5.4.5.2).
	X520Attributes_ENCRYPTEDSERIALNUMBER X520Attributes = "ENCRYPTEDSERIALNUMBER"
	// X520Attributes_COUNTRYNAME is the "countryName" attribute (OID 2.5.4.6).
	X520Attributes_COUNTRYNAME X520Attributes = "COUNTRYNAME"
	// X520Attributes_ENCRYPTEDCOUNTRYNAME is the "encryptedCountryName" attribute (OID 2.5.4.6.2).
	X520Attributes_ENCRYPTEDCOUNTRYNAME X520Attributes = "ENCRYPTEDCOUNTRYNAME"
	// X520Attributes_LOCALITYNAME is the "localityName" attribute (OID 2.5.4.7).
	X520Attributes_LOCALITYNAME X520Attributes = "LOCALITYNAME"
	// X520Attributes_ENCRYPTEDLOCALITYNAME is the "encryptedLocalityName" attribute (OID 2.5.4.7.2).
	X520Attributes_ENCRYPTEDLOCALITYNAME X520Attributes = "ENCRYPTEDLOCALITYNAME"
	// X520Attributes_COLLECTIVELOCALITYNAME is the "collectiveLocalityName" attribute (OID 2.5.4.7.1).
	X520Attributes_COLLECTIVELOCALITYNAME X520Attributes = "COLLECTIVELOCALITYNAME"
	// X520Attributes_ENCRYPTEDCOLLECTIVELOCALITYNAME is the "encryptedCollectiveLocalityName" attribute (OID 2.5.4.7.1.2).
	X520Attributes_ENCRYPTEDCOLLECTIVELOCALITYNAME X520Attributes = "ENCRYPTEDCOLLECTIVELOCALITYNAME"
	// X520Attributes_STATEORPROVINCENAME is the "stateOrProvinceName" attribute (OID 2.5.4.8).
	X520Attributes_STATEORPROVINCENAME X520Attributes = "STATEORPROVINCENAME"
	// X520Attributes_ENCRYPTEDSTATEORPROVINCENAME is the "encryptedStateOrProvinceName" attribute (OID 2.5.4.8.2).
	X520Attributes_ENCRYPTEDSTATEORPROVINCENAME X520Attributes = "ENCRYPTEDSTATEORPROVINCENAME"
	// X520Attributes_COLLECTIVESTATEORPROVINCENAME is the "collectiveStateOrProvinceName" attribute (OID 2.5.4.8.1).
	X520Attributes_COLLECTIVESTATEORPROVINCENAME X520Attributes = "COLLECTIVESTATEORPROVINCENAME"
	// X520Attributes_ENCRYPTEDCOLLECTIVESTATEORPROVINCENAME is the "encryptedCollectiveStateOrProvinceName" attribute (OID 2.5.4.8.1.2).
	X520Attributes_ENCRYPTEDCOLLECTIVESTATEORPROVINCENAME X520Attributes = "ENCRYPTEDCOLLECTIVESTATEORPROVINCENAME"
	// X520Attributes_STREETADDRESS is the "streetAddress" attribute (OID 2.5.4.9).
	X520Attributes_STREETADDRESS X520Attributes = "STREETADDRESS"
	// X520Attributes_ENCRYPTEDSTREETADDRESS is the "encryptedStreetAddress" attribute (OID 2.5.4.9.2).
	X520Attributes_ENCRYPTEDSTREETADDRESS X520Attributes = "ENCRYPTEDSTREETADDRESS"
	// X520Attributes_COLLECTIVESTREETADDRESS is the "collectiveStreetAddress" attribute (OID 2.5.4.9.1).
	X520Attributes_COLLECTIVESTREETADDRESS X520Attributes = "COLLECTIVESTREETADDRESS"
	// X520Attributes_ENCRYPTEDCOLLECTIVESTREETADDRESS is the "encryptedCollectiveStreetAddress" attribute (OID 2.5.4.9.1.2).
	X520Attributes_ENCRYPTEDCOLLECTIVESTREETADDRESS X520Attributes = "ENCRYPTEDCOLLECTIVESTREETADDRESS"
	// X520Attributes_ORGANIZATIONNAME is the "organizationName" attribute (OID 2.5.4.10).
	X520Attributes_ORGANIZATIONNAME X520Attributes = "ORGANIZATIONNAME"
	// X520Attributes_ENCRYPTEDORGANIZATIONNAME is the "encryptedOrganizationName" attribute (OID 2.5.4.10.2).
	X520Attributes_ENCRYPTEDORGANIZATIONNAME X520Attributes = "ENCRYPTEDORGANIZATIONNAME"
	// X520Attributes_COLLECTIVEORGANIZATIONNAME is the "collectiveOrganizationName" attribute (OID 2.5.4.10.1).
	X520Attributes_COLLECTIVEORGANIZATIONNAME X520Attributes = "COLLECTIVEORGANIZATIONNAME"
	// X520Attributes_ENCRYPTEDCOLLECTIVEORGANIZATIONNAME is the "encryptedCollectiveOrganizationName" attribute (OID 2.5.4.10.1.2).
	X520Attributes_ENCRYPTEDCOLLECTIVEORGANIZATIONNAME X520Attributes = "ENCRYPTEDCOLLECTIVEORGANIZATIONNAME"
	// X520Attributes_ORGANIZATIONALUNITNAME is the "organizationalUnitName" attribute (OID 2.5.4.11).
	X520Attributes_ORGANIZATIONALUNITNAME X520Attributes = "ORGANIZATIONALUNITNAME"
	// X520Attributes_ENCRYPTEDORGANIZATIONALUNITNAME is the "encryptedOrganizationalUnitName" attribute (OID 2.5.4.11.2).
	X520Attributes_ENCRYPTEDORGANIZATIONALUNITNAME X520Attributes = "ENCRYPTEDORGANIZATIONALUNITNAME"
	// X520Attributes_COLLECTIVEORGANIZATIONALUNITNAME is the "collectiveOrganizationalUnitName" attribute (OID 2.5.4.11.1).
	X520Attributes_COLLECTIVEORGANIZATIONALUNITNAME X520Attributes = "COLLECTIVEORGANIZATIONALUNITNAME"
	// X520Attributes_ENCRYPTEDCOLLECTIVEORGANIZATIONALUNITNAM is the "encryptedCollectiveOrganizationalUnitNam" attribute (OID 2.5.4.11.1.2).
	X520Attributes_ENCRYPTEDCOLLECTIVEORGANIZATIONALUNITNAM X520Attributes = "ENCRYPTEDCOLLECTIVEORGANIZATIONALUNITNAM"
	// X520Attributes_TITLE is the "title" attribute (OID 2.5.4.12).
	X520Attributes_TITLE X520Attributes = "TITLE"
	// X520Attributes_ENCRYPTEDTITLE is the "encryptedTitle" attribute (OID 2.5.4.12.2).
	X520Attributes_ENCRYPTEDTITLE X520Attributes = "ENCRYPTEDTITLE"
	// X520Attributes_DESCRIPTION is the "description" attribute (OID 2.5.4.13).
	X520Attributes_DESCRIPTION X520Attributes = "DESCRIPTION"
	// X520Attributes_ENCRYPTEDDESCRIPTION is the "encryptedDescription" attribute (OID 2.5.4.13.2).
	X520Attributes_ENCRYPTEDDESCRIPTION X520Attributes = "ENCRYPTEDDESCRIPTION"
	// X520Attributes_SEARCHGUIDE is the "searchGuide" attribute (OID 2.5.4.14).
	X520Attributes_SEARCHGUIDE X520Attributes = "SEARCHGUIDE"
	// X520Attributes_ENCRYPTEDSEARCHGUIDE is the "encryptedSearchGuide" attribute (OID 2.5.4.14.2).
	X520Attributes_ENCRYPTEDSEARCHGUIDE X520Attributes = "ENCRYPTEDSEARCHGUIDE"
	// X520Attributes_BUSINESSCATEGORY is the "businessCategory" attribute (OID 2.5.4.15).
	X520Attributes_BUSINESSCATEGORY X520Attributes = "BUSINESSCATEGORY"
	// X520Attributes_ENCRYPTEDBUSINESSCATEGORY is the "encryptedBusinessCategory" attribute (OID 2.5.4.15.2).
	X520Attributes_ENCRYPTEDBUSINESSCATEGORY X520Attributes = "ENCRYPTEDBUSINESSCATEGORY"
	// X520Attributes_POSTALADDRESS is the "postalAddress" attribute (OID 2.5.4.16).
	X520Attributes_POSTALADDRESS X520Attributes = "POSTALADDRESS"
	// X520Attributes_ENCRYPTEDPOSTALADDRESS is the "encryptedPostalAddress" attribute (OID 2.5.4.16.2).
	X520Attributes_ENCRYPTEDPOSTALADDRESS X520Attributes = "ENCRYPTEDPOSTALADDRESS"
	// X520Attributes_COLLECTIVEPOSTALADDRESS is the "collectivePostalAddress" attribute (OID 2.5.4.16.1).
	X520Attributes_COLLECTIVEPOSTALADDRESS X520Attributes = "COLLECTIVEPOSTALADDRESS"
	// X520Attributes_ENCRYPTEDCOLLECTIVEPOSTALADDRESS is the "encryptedCollectivePostalAddress" attribute (OID 2.5.4.16.1.2).
	X520Attributes_ENCRYPTEDCOLLECTIVEPOSTALADDRESS X520Attributes = "ENCRYPTEDCOLLECTIVEPOSTALADDRESS"
	// X520Attributes_POSTALCODE is the "postalCode" attribute (OID 2.5.4.17).
	X520Attributes_POSTALCODE X520Attributes = "POSTALCODE"
	// X520Attributes_ENCRYPTEDPOSTALCODE is the "encryptedPostalCode" attribute (OID 2.5.4.17.2).
	X520Attributes_ENCRYPTEDPOSTALCODE X520Attributes = "ENCRYPTEDPOSTALCODE"
	// X520Attributes_COLLECTIVEPOSTALCODE is the "collectivePostalCode" attribute (OID 2.5.4.17.1).
	X520Attributes_COLLECTIVEPOSTALCODE X520Attributes = "COLLECTIVEPOSTALCODE"
	// X520Attributes_ENCRYPTEDCOLLECTIVEPOSTALCODE is the "encryptedCollectivePostalCode" attribute (OID 2.5.4.17.1.2).
	X520Attributes_ENCRYPTEDCOLLECTIVEPOSTALCODE X520Attributes = "ENCRYPTEDCOLLECTIVEPOSTALCODE"
	// X520Attributes_POSTOFFICEBOX is the "postOfficeBox" attribute (OID 2.5.4.18).
	X520Attributes_POSTOFFICEBOX X520Attributes = "POSTOFFICEBOX"
	// X520Attributes_COLLECTIVEPOSTOFFICEBOX is the "collectivePostOfficeBox" attribute (OID 2.5.4.18.1).
	X520Attributes_COLLECTIVEPOSTOFFICEBOX X520Attributes = "COLLECTIVEPOSTOFFICEBOX"
	// X520Attributes_ENCRYPTEDPOSTOFFICEBOX is the "encryptedPostOfficeBox" attribute (OID 2.5.4.18.2).
	X520Attributes_ENCRYPTEDPOSTOFFICEBOX X520Attributes = "ENCRYPTEDPOSTOFFICEBOX"
	// X520Attributes_ENCRYPTEDCOLLECTIVEPOSTOFFICEBOX is the "encryptedCollectivePostOfficeBox" attribute (OID 2.5.4.18.1.2).
	X520Attributes_ENCRYPTEDCOLLECTIVEPOSTOFFICEBOX X520Attributes = "ENCRYPTEDCOLLECTIVEPOSTOFFICEBOX"
	// X520Attributes_PHYSICALDELIVERYOFFICENAME is the "physicalDeliveryOfficeName" attribute (OID 2.5.4.19).
	X520Attributes_PHYSICALDELIVERYOFFICENAME X520Attributes = "PHYSICALDELIVERYOFFICENAME"
	// X520Attributes_COLLECTIVEPHYSICALDELIVERYOFFICENAME is the "collectivePhysicalDeliveryOfficeName" attribute (OID 2.5.4.19.1).
	X520Attributes_COLLECTIVEPHYSICALDELIVERYOFFICENAME X520Attributes = "COLLECTIVEPHYSICALDELIVERYOFFICENAME"
	// X520Attributes_ENCRYPTEDPHYSICALDELIVERYOFFICENAME is the "encryptedPhysicalDeliveryOfficeName" attribute (OID 2.5.4.19.2).
	X520Attributes_ENCRYPTEDPHYSICALDELIVERYOFFICENAME X520Attributes = "ENCRYPTEDPHYSICALDELIVERYOFFICENAME"
	// X520Attributes_ENCRYPTEDCOLLECTIVEPHYSICALDELIVERYOFFICENAME is the "encryptedCollectivePhysicalDeliveryOfficeName" attribute (OID 2.5.4.19.1.2).
	X520Attributes_ENCRYPTEDCOLLECTIVEPHYSICALDELIVERYOFFICENAME X520Attributes = "ENCRYPTEDCOLLECTIVEPHYSICALDELIVERYOFFICENAME"
	// X520Attributes_TELEPHONENUMBER is the "telephoneNumber" attribute (OID 2.5.4.20).
	X520Attributes_TELEPHONENUMBER X520Attributes = "TELEPHONENUMBER"
	// X520Attributes_ENCRYPTEDTELEPHONENUMBER is the "encryptedTelephoneNumber" attribute (OID 2.5.4.20.2).
	X520Attributes_ENCRYPTEDTELEPHONENUMBER X520Attributes = "ENCRYPTEDTELEPHONENUMBER"
	// X520Attributes_COLLECTIVETELEPHONENUMBER is the "collectiveTelephoneNumber" attribute (OID 2.5.4.20.1).
	X520Attributes_COLLECTIVETELEPHONENUMBER X520Attributes = "COLLECTIVETELEPHONENUMBER"
	// X520Attributes_ENCRYPTEDCOLLECTIVETELEPHONENUMBER is the "encryptedCollectiveTelephoneNumber" attribute (OID 2.5.4.20.1.2).
	X520Attributes_ENCRYPTEDCOLLECTIVETELEPHONENUMBER X520Attributes = "ENCRYPTEDCOLLECTIVETELEPHONENUMBER"
	// X520Attributes_TELEXNUMBER is the "telexNumber" attribute (OID 2.5.4.21).
	X520Attributes_TELEXNUMBER X520Attributes = "TELEXNUMBER"
	// X520Attributes_ENCRYPTEDTELEXNUMBER is the "encryptedTelexNumber" attribute (OID 2.5.4.21.2).
	X520Attributes_ENCRYPTEDTELEXNUMBER X520Attributes = "ENCRYPTEDTELEXNUMBER"
	// X520Attributes_COLLECTIVETELEXNUMBER is the "collectiveTelexNumber" attribute (OID 2.5.4.21.1).
	X520Attributes_COLLECTIVETELEXNUMBER X520Attributes = "COLLECTIVETELEXNUMBER"
	// X520Attributes_ENCRYPTEDCOLLECTIVETELEXNUMBER is the "encryptedCollectiveTelexNumber" attribute (OID 2.5.4.21.1.2).
	X520Attributes_ENCRYPTEDCOLLECTIVETELEXNUMBER X520Attributes = "ENCRYPTEDCOLLECTIVETELEXNUMBER"
	// X520Attributes_TELETEXTERMINALIDENTIFIER is the "teletexTerminalIdentifier" attribute (OID 2.5.4.22).
	X520Attributes_TELETEXTERMINALIDENTIFIER X520Attributes = "TELETEXTERMINALIDENTIFIER"
	// X520Attributes_ENCRYPTEDTELETEXTERMINALIDENTIFIER is the "encryptedTeletexTerminalIdentifier" attribute (OID 2.5.4.22.2).
	X520Attributes_ENCRYPTEDTELETEXTERMINALIDENTIFIER X520Attributes = "ENCRYPTEDTELETEXTERMINALIDENTIFIER"
	// X520Attributes_COLLECTIVETELETEXTERMINALIDENTIFIER is the "collectiveTeletexTerminalIdentifier" attribute (OID 2.5.4.22.1).
	X520Attributes_COLLECTIVETELETEXTERMINALIDENTIFIER X520Attributes = "COLLECTIVETELETEXTERMINALIDENTIFIER"
	// X520Attributes_ENCRYPTEDCOLLECTIVETELETEXTERMINALIDENTIFIER is the "encryptedCollectiveTeletexTerminalIdentifier" attribute (OID 2.5.4.22.1.2).
	X520Attributes_ENCRYPTEDCOLLECTIVETELETEXTERMINALIDENTIFIER X520Attributes = "ENCRYPTEDCOLLECTIVETELETEXTERMINALIDENTIFIER"
	// X520Attributes_FACSIMILETELEPHONENUMBER is the "facsimileTelephoneNumber" attribute (OID 2.5.4.23).
	X520Attributes_FACSIMILETELEPHONENUMBER X520Attributes = "FACSIMILETELEPHONENUMBER"
	// X520Attributes_ENCRYPTEDFACSIMILETELEPHONENUMBER is the "encryptedFacsimileTelephoneNumber" attribute (OID 2.5.4.23.2).
	X520Attributes_ENCRYPTEDFACSIMILETELEPHONENUMBER X520Attributes = "ENCRYPTEDFACSIMILETELEPHONENUMBER"
	// X520Attributes_COLLECTIVEFACSIMILETELEPHONENUMBER is the "collectiveFacsimileTelephoneNumber" attribute (OID 2.5.4.23.1).
	X520Attributes_COLLECTIVEFACSIMILETELEPHONENUMBER X520Attributes = "COLLECTIVEFACSIMILETELEPHONENUMBER"
	// X520Attributes_ENCRYPTEDCOLLECTIVEFACSIMILETELEPHONENUMBER is the "encryptedCollectiveFacsimileTelephoneNumber" attribute (OID 2.5.4.23.1.2).
	X520Attributes_ENCRYPTEDCOLLECTIVEFACSIMILETELEPHONENUMBER X520Attributes = "ENCRYPTEDCOLLECTIVEFACSIMILETELEPHONENUMBER"
	// X520Attributes_X121ADDRESS is the "x121Address" attribute (OID 2.5.4.24).
	X520Attributes_X121ADDRESS X520Attributes = "X121ADDRESS"
	// X520Attributes_ENCRYPTEDX121ADDRESS is the "encryptedX121Address" attribute (OID 2.5.4.24.2).
	X520Attributes_ENCRYPTEDX121ADDRESS X520Attributes = "ENCRYPTEDX121ADDRESS"
	// X520Attributes_INTERNATIONALISDNNUMBER is the "internationalISDNNumber" attribute (OID 2.5.4.25).
	X520Attributes_INTERNATIONALISDNNUMBER X520Attributes = "INTERNATIONALISDNNUMBER"
	// X520Attributes_ENCRYPTEDINTERNATIONALISDNNUMBER is the "encryptedInternationalISDNNumber" attribute (OID 2.5.4.25.2).
	X520Attributes_ENCRYPTEDINTERNATIONALISDNNUMBER X520Attributes = "ENCRYPTEDINTERNATIONALISDNNUMBER"
	// X520Attributes_COLLECTIVEINTERNATIONALISDNNUMBER is the "collectiveInternationalISDNNumber" attribute (OID 2.5.4.25.1).
	X520Attributes_COLLECTIVEINTERNATIONALISDNNUMBER X520Attributes = "COLLECTIVEINTERNATIONALISDNNUMBER"
	// X520Attributes_ENCRYPTEDCOLLECTIVEINTERNATIONALISDNNUMBER is the "encryptedCollectiveInternationalISDNNumber" attribute (OID 2.5.4.25.1.2).
	X520Attributes_ENCRYPTEDCOLLECTIVEINTERNATIONALISDNNUMBER X520Attributes = "ENCRYPTEDCOLLECTIVEINTERNATIONALISDNNUMBER"
	// X520Attributes_REGISTEREDADDRESS is the "registeredAddress" attribute (OID 2.5.4.26).
	X520Attributes_REGISTEREDADDRESS X520Attributes = "REGISTEREDADDRESS"
	// X520Attributes_ENCRYPTEDREGISTEREDADDRESS is the "encryptedRegisteredAddress" attribute (OID 2.5.4.26.2).
	X520Attributes_ENCRYPTEDREGISTEREDADDRESS X520Attributes = "ENCRYPTEDREGISTEREDADDRESS"
	// X520Attributes_DESTINATIONINDICATOR is the "destinationIndicator" attribute (OID 2.5.4.27).
	X520Attributes_DESTINATIONINDICATOR X520Attributes = "DESTINATIONINDICATOR"
	// X520Attributes_ENCRYPTEDDESTINATIONINDICATOR is the "encryptedDestinationIndicator" attribute (OID 2.5.4.27.2).
	X520Attributes_ENCRYPTEDDESTINATIONINDICATOR X520Attributes = "ENCRYPTEDDESTINATIONINDICATOR"
	// X520Attributes_PREFERREDDELIVERYMETHOD is the "preferredDeliveryMethod" attribute (OID 2.5.4.28).
	X520Attributes_PREFERREDDELIVERYMETHOD X520Attributes = "PREFERREDDELIVERYMETHOD"
	// X520Attributes_ENCRYPTEDPREFERREDDELIVERYMETHOD is the "encryptedPreferredDeliveryMethod" attribute (OID 2.5.4.28.2).
	X520Attributes_ENCRYPTEDPREFERREDDELIVERYMETHOD X520Attributes = "ENCRYPTEDPREFERREDDELIVERYMETHOD"
	// X520Attributes_PRESENTATIONADDRESS is the "presentationAddress" attribute (OID 2.5.4.29).
	X520Attributes_PRESENTATIONADDRESS X520Attributes = "PRESENTATIONADDRESS"
	// X520Attributes_ENCRYPTEDPRESENTATIONADDRESS is the "encryptedPresentationAddress" attribute (OID 2.5.4.29.2).
	X520Attributes_ENCRYPTEDPRESENTATIONADDRESS X520Attributes = "ENCRYPTEDPRESENTATIONADDRESS"
	// X520Attributes_SUPPORTEDAPPLICATIONCONTEXT is the "supportedApplicationContext" attribute (OID 2.5.4.30).
	X520Attributes_SUPPORTEDAPPLICATIONCONTEXT X520Attributes = "SUPPORTEDAPPLICATIONCONTEXT"
	// X520Attributes_ENCRYPTEDSUPPORTEDAPPLICATIONCONTEXT is the "encryptedSupportedApplicationContext" attribute (OID 2.5.4.30.2).
	X520Attributes_ENCRYPTEDSUPPORTEDAPPLICATIONCONTEXT X520Attributes = "ENCRYPTEDSUPPORTEDAPPLICATIONCONTEXT"
	// X520Attributes_MEMBER is the "member" attribute (OID 2.5.4.31).
	X520Attributes_MEMBER X520Attributes = "MEMBER"
	// X520Attributes_ENCRYPTEDMEMBER is the "encryptedMember" attribute (OID 2.5.4.31.2).
	X520Attributes_ENCRYPTEDMEMBER X520Attributes = "ENCRYPTEDMEMBER"
	// X520Attributes_OWNER is the "owner" attribute (OID 2.5.4.32).
	X520Attributes_OWNER X520Attributes = "OWNER"
	// X520Attributes_ENCRYPTEDOWNER is the "encryptedOwner" attribute (OID 2.5.4.32.2).
	X520Attributes_ENCRYPTEDOWNER X520Attributes = "ENCRYPTEDOWNER"
	// X520Attributes_ROLEOCCUPANT is the "roleOccupant" attribute (OID 2.5.4.33).
	X520Attributes_ROLEOCCUPANT X520Attributes = "ROLEOCCUPANT"
	// X520Attributes_ENCRYPTEDROLEOCCUPANT is the "encryptedRoleOccupant" attribute (OID 2.5.4.33.2).
	X520Attributes_ENCRYPTEDROLEOCCUPANT X520Attributes = "ENCRYPTEDROLEOCCUPANT"
	// X520Attributes_SEEALSO is the "seeAlso" attribute (OID 2.5.4.34).
	X520Attributes_SEEALSO X520Attributes = "SEEALSO"
	// X520Attributes_ENCRYPTEDSEEALSO is the "encryptedSeeAlso" attribute (OID 2.5.4.34.2).
	X520Attributes_ENCRYPTEDSEEALSO X520Attributes = "ENCRYPTEDSEEALSO"
	// X520Attributes_USERPASSWORD is the "userPassword" attribute (OID 2.5.4.35).
	X520Attributes_USERPASSWORD X520Attributes = "USERPASSWORD"
	// X520Attributes_ENCRYPTEDUSERPASSWORD is the "encryptedUserPassword" attribute (OID 2.5.4.35.2).
	X520Attributes_ENCRYPTEDUSERPASSWORD X520Attributes = "ENCRYPTEDUSERPASSWORD"
	// X520Attributes_USERCERTIFICATE is the "userCertificate" attribute (OID 2.5.4.36).
	X520Attributes_USERCERTIFICATE X520Attributes = "USERCERTIFICATE"
	// X520Attributes_ENCRYPTEDUSERCERTIFICATE is the "encryptedUserCertificate" attribute (OID 2.5.4.36.2).
	X520Attributes_ENCRYPTEDUSERCERTIFICATE X520Attributes = "ENCRYPTEDUSERCERTIFICATE"
	// X520Attributes_CACERTIFICATE is the "cACertificate" attribute (OID 2.5.4.37).
	X520Attributes_CACERTIFICATE X520Attributes = "CACERTIFICATE"
	// X520Attributes_ENCRYPTEDCACERTIFICATE is the "encryptedCACertificate" attribute (OID 2.5.4.37.2).
	X520Attributes_ENCRYPTEDCACERTIFICATE X520Attributes = "ENCRYPTEDCACERTIFICATE"
	// X520Attributes_AUTHORITYREVOCATIONLIST is the "authorityRevocationList" attribute (OID 2.5.4.38).
	X520Attributes_AUTHORITYREVOCATIONLIST X520Attributes = "AUTHORITYREVOCATIONLIST"
	// X520Attributes_ENCRYPTEDAUTHORITYREVOCATIONLIST is the "encryptedAuthorityRevocationList" attribute (OID 2.5.4.38.2).
	X520Attributes_ENCRYPTEDAUTHORITYREVOCATIONLIST X520Attributes = "ENCRYPTEDAUTHORITYREVOCATIONLIST"
	// X520Attributes_CERTIFICATEREVOCATIONLIST is the "certificateRevocationList" attribute (OID 2.5.4.39).
	X520Attributes_CERTIFICATEREVOCATIONLIST X520Attributes = "CERTIFICATEREVOCATIONLIST"
	// X520Attributes_ENCRYPTEDCERTIFICATEREVOCATIONLIST is the "encryptedCertificateRevocationList" attribute (OID 2.5.4.39.2).
	X520Attributes_ENCRYPTEDCERTIFICATEREVOCATIONLIST X520Attributes = "ENCRYPTEDCERTIFICATEREVOCATIONLIST"
	// X520Attributes_CROSSCERTIFICATEPAIR is the "crossCertificatePair" attribute (OID 2.5.4.40).
	X520Attributes_CROSSCERTIFICATEPAIR X520Attributes = "CROSSCERTIFICATEPAIR"
	// X520Attributes_ENCRYPTEDCROSSCERTIFICATEPAIR is the "encryptedCrossCertificatePair" attribute (OID 2.5.4.40.2).
	X520Attributes_ENCRYPTEDCROSSCERTIFICATEPAIR X520Attributes = "ENCRYPTEDCROSSCERTIFICATEPAIR"
	// X520Attributes_NAME is the "name" attribute (OID 2.5.4.41).
	X520Attributes_NAME X520Attributes = "NAME"
	// X520Attributes_GIVENNAME is the "givenName" attribute (OID 2.5.4.42).
	X520Attributes_GIVENNAME X520Attributes = "GIVENNAME"
	// X520Attributes_ENCRYPTEDGIVENNAME is the "encryptedGivenName" attribute (OID 2.5.4.42.2).
	X520Attributes_ENCRYPTEDGIVENNAME X520Attributes = "ENCRYPTEDGIVENNAME"
	// X520Attributes_INITIALS is the "initials" attribute (OID 2.5.4.43).
	X520Attributes_INITIALS X520Attributes = "INITIALS"
	// X520Attributes_ENCRYPTEDINITIALS is the "encryptedInitials" attribute (OID 2.5.4.43.2).
	X520Attributes_ENCRYPTEDINITIALS X520Attributes = "ENCRYPTEDINITIALS"
	// X520Attributes_GENERATIONQUALIFIER is the "generationQualifier" attribute (OID 2.5.4.44).
	X520Attributes_GENERATIONQUALIFIER X520Attributes = "GENERATIONQUALIFIER"
	// X520Attributes_ENCRYPTEDGENERATIONQUALIFIER is the "encryptedGenerationQualifier" attribute (OID 2.5.4.44.2).
	X520Attributes_ENCRYPTEDGENERATIONQUALIFIER X520Attributes = "ENCRYPTEDGENERATIONQUALIFIER"
	// X520Attributes_UNIQUEIDENTIFIER is the "uniqueIdentifier" attribute (OID 2.5.4.45).
	X520Attributes_UNIQUEIDENTIFIER X520Attributes = "UNIQUEIDENTIFIER"
	// X520Attributes_ENCRYPTEDUNIQUEIDENTIFIER is the "encryptedUniqueIdentifier" attribute (OID 2.5.4.45.2).
	X520Attributes_ENCRYPTEDUNIQUEIDENTIFIER X520Attributes = "ENCRYPTEDUNIQUEIDENTIFIER"
	// X520Attributes_DNQUALIFIER is the "dnQualifier" attribute (OID 2.5.4.46).
	X520Attributes_DNQUALIFIER X520Attributes = "DNQUALIFIER"
	// X520Attributes_ENCRYPTEDDNQUALIFIER is the "encryptedDnQualifier" attribute (OID 2.5.4.46.2).
	X520Attributes_ENCRYPTEDDNQUALIFIER X520Attributes = "ENCRYPTEDDNQUALIFIER"
	// X520Attributes_ENHANCEDSEARCHGUIDE is the "enhancedSearchGuide" attribute (OID 2.5.4.47).
	X520Attributes_ENHANCEDSEARCHGUIDE X520Attributes = "ENHANCEDSEARCHGUIDE"
	// X520Attributes_ENCRYPTEDENHANCEDSEARCHGUIDE is the "encryptedEnhancedSearchGuide" attribute (OID 2.5.4.47.2).
	X520Attributes_ENCRYPTEDENHANCEDSEARCHGUIDE X520Attributes = "ENCRYPTEDENHANCEDSEARCHGUIDE"
	// X520Attributes_PROTOCOLINFORMATION is the "protocolInformation" attribute (OID 2.5.4.48).
	X520Attributes_PROTOCOLINFORMATION X520Attributes = "PROTOCOLINFORMATION"
	// X520Attributes_ENCRYPTEDPROTOCOLINFORMATION is the "encryptedProtocolInformation" attribute (OID 2.5.4.48.2).
	X520Attributes_ENCRYPTEDPROTOCOLINFORMATION X520Attributes = "ENCRYPTEDPROTOCOLINFORMATION"
	// X520Attributes_DISTINGUISHEDNAME is the "distinguishedName" attribute (OID 2.5.4.49).
	X520Attributes_DISTINGUISHEDNAME X520Attributes = "DISTINGUISHEDNAME"
	// X520Attributes_ENCRYPTEDDISTINGUISHEDNAME is the "encryptedDistinguishedName" attribute (OID 2.5.4.49.2).
	X520Attributes_ENCRYPTEDDISTINGUISHEDNAME X520Attributes = "ENCRYPTEDDISTINGUISHEDNAME"
	// X520Attributes_UNIQUEMEMBER is the "uniqueMember" attribute (OID 2.5.4.50).
	X520Attributes_UNIQUEMEMBER X520Attributes = "UNIQUEMEMBER"
	// X520Attributes_ENCRYPTEDUNIQUEMEMBER is the "encryptedUniqueMember" attribute (OID 2.5.4.50.2).
	X520Attributes_ENCRYPTEDUNIQUEMEMBER X520Attributes = "ENCRYPTEDUNIQUEMEMBER"
	// X520Attributes_HOUSEIDENTIFIER is the "houseIdentifier" attribute (OID 2.5.4.51).
	X520Attributes_HOUSEIDENTIFIER X520Attributes = "HOUSEIDENTIFIER"
	// X520Attributes_ENCRYPTEDHOUSEIDENTIFIER is the "encryptedHouseIdentifier" attribute (OID 2.5.4.51.2).
	X520Attributes_ENCRYPTEDHOUSEIDENTIFIER X520Attributes = "ENCRYPTEDHOUSEIDENTIFIER"
	// X520Attributes_SUPPORTEDALGORITHMS is the "supportedAlgorithms" attribute (OID 2.5.4.52).
	X520Attributes_SUPPORTEDALGORITHMS X520Attributes = "SUPPORTEDALGORITHMS"
	// X520Attributes_ENCRYPTEDSUPPORTEDALGORITHMS is the "encryptedSupportedAlgorithms" attribute (OID 2.5.4.52.2).
	X520Attributes_ENCRYPTEDSUPPORTEDALGORITHMS X520Attributes = "ENCRYPTEDSUPPORTEDALGORITHMS"
	// X520Attributes_DELTAREVOCATIONLIST is the "deltaRevocationList" attribute (OID 2.5.4.53).
	X520Attributes_DELTAREVOCATIONLIST X520Attributes = "DELTAREVOCATIONLIST"
	// X520Attributes_ENCRYPTEDDELTAREVOCATIONLIST is the "encryptedDeltaRevocationList" attribute (OID 2.5.4.53.2).
	X520Attributes_ENCRYPTEDDELTAREVOCATIONLIST X520Attributes = "ENCRYPTEDDELTAREVOCATIONLIST"
	// X520Attributes_DMDNAME is the "dmdName" attribute (OID 2.5.4.54).
	X520Attributes_DMDNAME X520Attributes = "DMDNAME"
	// X520Attributes_ENCRYPTEDDMDNAME is the "encryptedDmdName" attribute (OID 2.5.4.54.2).
	X520Attributes_ENCRYPTEDDMDNAME X520Attributes = "ENCRYPTEDDMDNAME"
	// X520Attributes_CLEARANCE is the "clearance" attribute (OID 2.5.4.55).
	X520Attributes_CLEARANCE X520Attributes = "CLEARANCE"
	// X520Attributes_ENCRYPTEDCLEARANCE is the "encryptedClearance" attribute (OID 2.5.4.55.2).
	X520Attributes_ENCRYPTEDCLEARANCE X520Attributes = "ENCRYPTEDCLEARANCE"
	// X520Attributes_DEFAULTDIRQOP is the "defaultDirQop" attribute (OID 2.5.4.56).
	X520Attributes_DEFAULTDIRQOP X520Attributes = "DEFAULTDIRQOP"
	// X520Attributes_ENCRYPTEDDEFAULTDIRQOP is the "encryptedDefaultDirQop" attribute (OID 2.5.4.56.2).
	X520Attributes_ENCRYPTEDDEFAULTDIRQOP X520Attributes = "ENCRYPTEDDEFAULTDIRQOP"
	// X520Attributes_ATTRIBUTEINTEGRITYINFO is the "attributeIntegrityInfo" attribute (OID 2.5.4.57).
	X520Attributes_ATTRIBUTEINTEGRITYINFO X520Attributes = "ATTRIBUTEINTEGRITYINFO"
	// X520Attributes_ENCRYPTEDATTRIBUTEINTEGRITYINFO is the "encryptedAttributeIntegrityInfo" attribute (OID 2.5.4.57.2).
	X520Attributes_ENCRYPTEDATTRIBUTEINTEGRITYINFO X520Attributes = "ENCRYPTEDATTRIBUTEINTEGRITYINFO"
	// X520Attributes_ATTRIBUTECERTIFICATE is the "attributeCertificate" attribute (OID 2.5.4.58).
	X520Attributes_ATTRIBUTECERTIFICATE X520Attributes = "ATTRIBUTECERTIFICATE"
	// X520Attributes_ENCRYPTEDATTRIBUTECERTIFICATE is the "encryptedAttributeCertificate" attribute (OID 2.5.4.58.2).
	X520Attributes_ENCRYPTEDATTRIBUTECERTIFICATE X520Attributes = "ENCRYPTEDATTRIBUTECERTIFICATE"
	// X520Attributes_ATTRIBUTECERTIFICATEREVOCATIONLIST is the "attributeCertificateRevocationList" attribute (OID 2.5.4.59).
	X520Attributes_ATTRIBUTECERTIFICATEREVOCATIONLIST X520Attributes = "ATTRIBUTECERTIFICATEREVOCATIONLIST"
	// X520Attributes_ENCRYPTEDATTRIBUTECERTIFICATEREVOCATIONLIST is the "encryptedAttributeCertificateRevocationList" attribute (OID 2.5.4.59.2).
	X520Attributes_ENCRYPTEDATTRIBUTECERTIFICATEREVOCATIONLIST X520Attributes = "ENCRYPTEDATTRIBUTECERTIFICATEREVOCATIONLIST"
	// X520Attributes_CONFKEYINFO is the "confKeyInfo" attribute (OID 2.5.4.60).
	X520Attributes_CONFKEYINFO X520Attributes = "CONFKEYINFO"
	// X520Attributes_ENCRYPTEDCONFKEYINFO is the "encryptedConfKeyInfo" attribute (OID 2.5.4.60.2).
	X520Attributes_ENCRYPTEDCONFKEYINFO X520Attributes = "ENCRYPTEDCONFKEYINFO"
	// X520Attributes_AACERTIFICATE is the "aACertificate" attribute (OID 2.5.4.61).
	X520Attributes_AACERTIFICATE X520Attributes = "AACERTIFICATE"
	// X520Attributes_ATTRIBUTEDESCRIPTORCERTIFICATE is the "attributeDescriptorCertificate" attribute (OID 2.5.4.62).
	X520Attributes_ATTRIBUTEDESCRIPTORCERTIFICATE X520Attributes = "ATTRIBUTEDESCRIPTORCERTIFICATE"
	// X520Attributes_ATTRIBUTEAUTHORITYREVOCATIONLIST is the "attributeAuthorityRevocationList" attribute (OID 2.5.4.63).
	X520Attributes_ATTRIBUTEAUTHORITYREVOCATIONLIST X520Attributes = "ATTRIBUTEAUTHORITYREVOCATIONLIST"
	// X520Attributes_FAMILY_INFORMATION is the "family-information" attribute (OID 2.5.4.64).
	X520Attributes_FAMILY_INFORMATION X520Attributes = "FAMILY_INFORMATION"
	// X520Attributes_PSEUDONYM is the "pseudonym" attribute (OID 2.5.4.65).
	X520Attributes_PSEUDONYM X520Attributes = "PSEUDONYM"
	// X520Attributes_COMMUNICATIONSSERVICE is the "communicationsService" attribute (OID 2.5.4.66).
	X520Attributes_COMMUNICATIONSSERVICE X520Attributes = "COMMUNICATIONSSERVICE"
	// X520Attributes_COMMUNICATIONSNETWORK is the "communicationsNetwork" attribute (OID 2.5.4.67).
	X520Attributes_COMMUNICATIONSNETWORK X520Attributes = "COMMUNICATIONSNETWORK"
	// X520Attributes_CERTIFICATIONPRACTICESTMT is the "certificationPracticeStmt" attribute (OID 2.5.4.68).
	X520Attributes_CERTIFICATIONPRACTICESTMT X520Attributes = "CERTIFICATIONPRACTICESTMT"
	// X520Attributes_CERTIFICATEPOLICY is the "certificatePolicy" attribute (OID 2.5.4.69).
	X520Attributes_CERTIFICATEPOLICY X520Attributes = "CERTIFICATEPOLICY"
	// X520Attributes_PKIPATH is the "pkiPath" attribute (OID 2.5.4.70).
	X520Attributes_PKIPATH X520Attributes = "PKIPATH"
	// X520Attributes_PRIVPOLICY is the "privPolicy" attribute (OID 2.5.4.71).
	X520Attributes_PRIVPOLICY X520Attributes = "PRIVPOLICY"
	// X520Attributes_ROLE is the "role" attribute (OID 2.5.4.72).
	X520Attributes_ROLE X520Attributes = "ROLE"
	// X520Attributes_DELEGATIONPATH is the "delegationPath" attribute (OID 2.5.4.73).
	X520Attributes_DELEGATIONPATH X520Attributes = "DELEGATIONPATH"
	// X520Attributes_PROTPRIVPOLICY is the "protPrivPolicy" attribute (OID 2.5.4.74).
	X520Attributes_PROTPRIVPOLICY X520Attributes = "PROTPRIVPOLICY"
	// X520Attributes_XMLPRIVILEGEINFO is the "xMLPrivilegeInfo" attribute (OID 2.5.4.75).
	X520Attributes_XMLPRIVILEGEINFO X520Attributes = "XMLPRIVILEGEINFO"
	// X520Attributes_XMLPRIVPOLICY is the "xmlPrivPolicy" attribute (OID 2.5.4.76).
	X520Attributes_XMLPRIVPOLICY X520Attributes = "XMLPRIVPOLICY"
	// X520Attributes_UUIDPAIR is the "uuidpair" attribute (OID 2.5.4.77).
	X520Attributes_UUIDPAIR X520Attributes = "UUIDPAIR"
	// X520Attributes_TAGOID is the "tagOid" attribute (OID 2.5.4.78).
	X520Attributes_TAGOID X520Attributes = "TAGOID"
	// X520Attributes_UIIFORMAT is the "uiiFormat" attribute (OID 2.5.4.79).
	X520Attributes_UIIFORMAT X520Attributes = "UIIFORMAT"
	// X520Attributes_UIIINURN is the "uiiInUrn" attribute (OID 2.5.4.80).
	X520Attributes_UIIINURN X520Attributes = "UIIINURN"
	// X520Attributes_CONTENTURL is the "contentUrl" attribute (OID 2.5.4.81).
	X520Attributes_CONTENTURL X520Attributes = "CONTENTURL"
	// X520Attributes_PERMISSION is the "permission" attribute (OID 2.5.4.82).
	X520Attributes_PERMISSION X520Attributes = "PERMISSION"
	// X520Attributes_URI is the "uri" attribute (OID 2.5.4.83).
	X520Attributes_URI X520Attributes = "URI"
	// X520Attributes_PWDATTRIBUTE is the "pwdAttribute" attribute (OID 2.5.4.84).
	X520Attributes_PWDATTRIBUTE X520Attributes = "PWDATTRIBUTE"
	// X520Attributes_USERPWD is the "userPwd" attribute (OID 2.5.4.85).
	X520Attributes_USERPWD X520Attributes = "USERPWD"
	// X520Attributes_URN is the "urn" attribute (OID 2.5.4.86).
	X520Attributes_URN X520Attributes = "URN"
	// X520Attributes_URL is the "url" attribute (OID 2.5.4.87).
	X520Attributes_URL X520Attributes = "URL"
	// X520Attributes_UTMCOORDINATES is the "utmCoordinates" attribute (OID 2.5.4.88).
	X520Attributes_UTMCOORDINATES X520Attributes = "UTMCOORDINATES"
	// X520Attributes_URNC is the "urnC" attribute (OID 2.5.4.89).
	X520Attributes_URNC X520Attributes = "URNC"
	// X520Attributes_UII is the "uii" attribute (OID 2.5.4.90).
	X520Attributes_UII X520Attributes = "UII"
	// X520Attributes_EPC is the "epc" attribute (OID 2.5.4.91).
	X520Attributes_EPC X520Attributes = "EPC"
	// X520Attributes_TAGAFI is the "tagAfi" attribute (OID 2.5.4.92).
	X520Attributes_TAGAFI X520Attributes = "TAGAFI"
	// X520Attributes_EPCFORMAT is the "epcFormat" attribute (OID 2.5.4.93).
	X520Attributes_EPCFORMAT X520Attributes = "EPCFORMAT"
	// X520Attributes_EPCINURN is the "epcInUrn" attribute (OID 2.5.4.94).
	X520Attributes_EPCINURN X520Attributes = "EPCINURN"
	// X520Attributes_LDAPURL is the "ldapUrl" attribute (OID 2.5.4.95).
	X520Attributes_LDAPURL X520Attributes = "LDAPURL"
	// X520Attributes_TAGLOCATION is the "tagLocation" attribute (OID 2.5.4.96).
	X520Attributes_TAGLOCATION X520Attributes = "TAGLOCATION"
	// X520Attributes_ORGANIZATIONIDENTIFIER is the "organizationIdentifier" attribute (OID 2.5.4.97).
	X520Attributes_ORGANIZATIONIDENTIFIER X520Attributes = "ORGANIZATIONIDENTIFIER"
	// X520Attributes_COUNTRYCODE3C is the "countryCode3c" attribute (OID 2.5.4.98).
	X520Attributes_COUNTRYCODE3C X520Attributes = "COUNTRYCODE3C"
	// X520Attributes_COUNTRYCODE3N is the "countryCode3n" attribute (OID 2.5.4.99).
	X520Attributes_COUNTRYCODE3N X520Attributes = "COUNTRYCODE3N"
	// X520Attributes_DNSNAME is the "dnsName" attribute (OID 2.5.4.100).
	X520Attributes_DNSNAME X520Attributes = "DNSNAME"
	// X520Attributes_EEPKCERTIFICATREVOCATIONLIST is the "eepkCertificatRevocationList" attribute (OID 2.5.4.101).
	X520Attributes_EEPKCERTIFICATREVOCATIONLIST X520Attributes = "EEPKCERTIFICATREVOCATIONLIST"
	// X520Attributes_EEATTRCERTIFICATEREVOCATIONLIST is the "eeAttrCertificateRevocationList" attribute (OID 2.5.4.102).
	X520Attributes_EEATTRCERTIFICATEREVOCATIONLIST X520Attributes = "EEATTRCERTIFICATEREVOCATIONLIST"
	// X520Attributes_USERPWDDESCRIPTION is the "userPwdDescription" attribute (OID 2.5.40.0).
	X520Attributes_USERPWDDESCRIPTION X520Attributes = "USERPWDDESCRIPTION"
	// X520Attributes_PWDVOCABULARYDESCRIPTION is the "pwdVocabularyDescription" attribute (OID 2.5.40.1).
	X520Attributes_PWDVOCABULARYDESCRIPTION X520Attributes = "PWDVOCABULARYDESCRIPTION"
	// X520Attributes_PWDALPHABETDESCRIPTION is the "pwdAlphabetDescription" attribute (OID 2.5.40.2).
	X520Attributes_PWDALPHABETDESCRIPTION X520Attributes = "PWDALPHABETDESCRIPTION"
	// X520Attributes_PWDENCALGDESCRIPTION is the "pwdEncAlgDescription" attribute (OID 2.5.40.3).
	X520Attributes_PWDENCALGDESCRIPTION X520Attributes = "PWDENCALGDESCRIPTION"
	// X520Attributes_UTMCOORDS is the "utmCoords" attribute (OID 2.5.40.4).
	X520Attributes_UTMCOORDS X520Attributes = "UTMCOORDS"
	// X520Attributes_UIIFORM is the "uiiForm" attribute (OID 2.5.40.5).
	X520Attributes_UIIFORM X520Attributes = "UIIFORM"
	// X520Attributes_EPCFORM is the "epcForm" attribute (OID 2.5.40.6).
	X520Attributes_EPCFORM X520Attributes = "EPCFORM"
	// X520Attributes_COUNTRYSTRING3C is the "countryString3c" attribute (OID 2.5.40.7).
	X520Attributes_COUNTRYSTRING3C X520Attributes = "COUNTRYSTRING3C"
	// X520Attributes_COUNTRYSTRING3N is the "countryString3n" attribute (OID 2.5.40.8).
	X520Attributes_COUNTRYSTRING3N X520Attributes = "COUNTRYSTRING3N"
	// X520Attributes_DNSSTRING is the "dnsString" attribute (OID 2.5.40.9).
	X520Attributes_DNSSTRING X520Attributes = "DNSSTRING"
	// X520Attributes_ATTRIBUTETYPEDESCRIPTION is the "attributeTypeDescription" attribute (OID 1.3.6.1.4.1.1466.115.121.1.3).
	X520Attributes_ATTRIBUTETYPEDESCRIPTION X520Attributes = "ATTRIBUTETYPEDESCRIPTION"
	// X520Attributes_BITSTRING is the "bitString" attribute (OID 1.3.6.1.4.1.1466.115.121.1.6).
	X520Attributes_BITSTRING X520Attributes = "BITSTRING"
	// X520Attributes_BOOLEAN is the "boolean" attribute (OID 1.3.6.1.4.1.1466.115.121.1.7).
	X520Attributes_BOOLEAN X520Attributes = "BOOLEAN"
	// X520Attributes_X509CERTIFICATE is the "x509Certificate" attribute (OID 1.3.6.1.4.1.1466.115.121.1.8).
	X520Attributes_X509CERTIFICATE X520Attributes = "X509CERTIFICATE"
	// X520Attributes_X509CERTIFICATELIST is the "x509CertificateList" attribute (OID 1.3.6.1.4.1.1466.115.121.1.9).
	X520Attributes_X509CERTIFICATELIST X520Attributes = "X509CERTIFICATELIST"
	// X520Attributes_X509CERTIFICATEPAIR is the "x509CertificatePair" attribute (OID 1.3.6.1.4.1.1466.115.121.1.10).
	X520Attributes_X509CERTIFICATEPAIR X520Attributes = "X509CERTIFICATEPAIR"
	// X520Attributes_COUNTRYSTRING is the "countryString" attribute (OID 1.3.6.1.4.1.1466.115.121.1.11).
	X520Attributes_COUNTRYSTRING X520Attributes = "COUNTRYSTRING"
	// X520Attributes_DN is the "dn" attribute (OID 1.3.6.1.4.1.1466.115.121.1.12).
	X520Attributes_DN X520Attributes = "DN"
	// X520Attributes_DELIVERYMETHOD is the "deliveryMethod" attribute (OID 1.3.6.1.4.1.1466.115.121.1.14).
	X520Attributes_DELIVERYMETHOD X520Attributes = "DELIVERYMETHOD"
	// X520Attributes_DIRECTORYSTRING is the "directoryString" attribute (OID 1.3.6.1.4.1.1466.115.121.1.15).
	X520Attributes_DIRECTORYSTRING X520Attributes = "DIRECTORYSTRING"
	// X520Attributes_DITCONTENTRULEDESCRIPTION is the "dITContentRuleDescription" attribute (OID 1.3.6.1.4.1.1466.115.121.1.16).
	X520Attributes_DITCONTENTRULEDESCRIPTION X520Attributes = "DITCONTENTRULEDESCRIPTION"
	// X520Attributes_DITSTRUCTURERULEDESCRIPTION is the "dITStructureRuleDescription" attribute (OID 1.3.6.1.4.1.1466.115.121.1.17).
	X520Attributes_DITSTRUCTURERULEDESCRIPTION X520Attributes = "DITSTRUCTURERULEDESCRIPTION"
	// X520Attributes_ENHANCEDGUIDE is the "enhancedGuide" attribute (OID 1.3.6.1.4.1.1466.115.121.1.21).
	X520Attributes_ENHANCEDGUIDE X520Attributes = "ENHANCEDGUIDE"
	// X520Attributes_FACSIMILETELEPHONENR is the "facsimileTelephoneNr" attribute (OID 1.3.6.1.4.1.1466.115.121.1.22).
	X520Attributes_FACSIMILETELEPHONENR X520Attributes = "FACSIMILETELEPHONENR"
	// X520Attributes_FAX is the "fax" attribute (OID 1.3.6.1.4.1.1466.115.121.1.23).
	X520Attributes_FAX X520Attributes = "FAX"
	// X520Attributes_GENERALIZEDTIME is the "generalizedTime" attribute (OID 1.3.6.1.4.1.1466.115.121.1.24).
	X520Attributes_GENERALIZEDTIME X520Attributes = "GENERALIZEDTIME"
	// X520Attributes_GUIDE is the "guide" attribute (OID 1.3.6.1.4.1.1466.115.121.1.25).
	X520Attributes_GUIDE X520Attributes = "GUIDE"
	// X520Attributes_IA5STRING is the "ia5String" attribute (OID 1.3.6.1.4.1.1466.115.121.1.26).
	X520Attributes_IA5STRING X520Attributes = "IA5STRING"
	// X520Attributes_INTEGER is the "integer" attribute (OID 1.3.6.1.4.1.1466.115.121.1.27).
	X520Attributes_INTEGER X520Attributes = "INTEGER"
	// X520Attributes_JPEG is the "jpeg" attribute (OID 1.3.6.1.4.1.1466.115.121.1.28).
	X520Attributes_JPEG X520Attributes = "JPEG"
	// X520Attributes_MATCHINGRULEDESCRIPTION is the "matchingRuleDescription" attribute (OID 1.3.6.1.4.1.1466.115.121.1.30).
	X520Attributes_MATCHINGRULEDESCRIPTION X520Attributes = "MATCHINGRULEDESCRIPTION"
	// X520Attributes_MATCHINGRULEUSEDESCRIPTION is the "matchingRuleUseDescription" attribute (OID 1.3.6.1.4.1.1466.115.121.1.31).
	X520Attributes_MATCHINGRULEUSEDESCRIPTION X520Attributes = "MATCHINGRULEUSEDESCRIPTION"
	// X520Attributes_NAMEANDOPTIONALUID is the "nameAndOptionalUID" attribute (OID 1.3.6.1.4.1.1466.115.121.1.34).
	X520Attributes_NAMEANDOPTIONALUID X520Attributes = "NAMEANDOPTIONALUID"
	// X520Attributes_NAMEFORMDESCRIPTION is the "nameFormDescription" attribute (OID 1.3.6.1.4.1.1466.115.121.1.35).
	X520Attributes_NAMEFORMDESCRIPTION X520Attributes = "NAMEFORMDESCRIPTION"
	// X520Attributes_NUMERICSTRING is the "numericString" attribute (OID 1.3.6.1.4.1.1466.115.121.1.36).
	X520Attributes_NUMERICSTRING X520Attributes = "NUMERICSTRING"
	// X520Attributes_OBJECTCLASSDESCRIPTION is the "objectClassDescription" attribute (OID 1.3.6.1.4.1.1466.115.121.1.37).
	X520Attributes_OBJECTCLASSDESCRIPTION X520Attributes = "OBJECTCLASSDESCRIPTION"
	// X520Attributes_OID is the "oid" attribute (OID 1.3.6.1.4.1.1466.115.121.1.38).
	X520Attributes_OID X520Attributes = "OID"
	// X520Attributes_OTHERMAILBOX is the "otherMailbox" attribute (OID 1.3.6.1.4.1.1466.115.121.1.39).
	X520Attributes_OTHERMAILBOX X520Attributes = "OTHERMAILBOX"
	// X520Attributes_OCTETSTRING is the "octetString" attribute (OID 1.3.6.1.4.1.1466.115.121.1.40).
	X520Attributes_OCTETSTRING X520Attributes = "OCTETSTRING"
	// X520Attributes_POSTALADDR is the "postalAddr" attribute (OID 1.3.6.1.4.1.1466.115.121.1.41).
	X520Attributes_POSTALADDR X520Attributes = "POSTALADDR"
	// X520Attributes_PRESENTATIONADDR is the "presentationAddr" attribute (OID 1.3.6.1.4.1.1466.115.121.1.43).
	X520Attributes_PRESENTATIONADDR X520Attributes = "PRESENTATIONADDR"
	// X520Attributes_PRINTABLESTRING is the "printableString" attribute (OID 1.3.6.1.4.1.1466.115.121.1.44).
	X520Attributes_PRINTABLESTRING X520Attributes = "PRINTABLESTRING"
	// X520Attributes_SUBTREESPEC is the "subtreeSpec" attribute (OID 1.3.6.1.4.1.1466.115.121.1.45).
	X520Attributes_SUBTREESPEC X520Attributes = "SUBTREESPEC"
	// X520Attributes_X509SUPPORTEDALGORITHM is the "x509SupportedAlgorithm" attribute (OID 1.3.6.1.4.1.1466.115.121.1.49).
	X520Attributes_X509SUPPORTEDALGORITHM X520Attributes = "X509SUPPORTEDALGORITHM"
	// X520Attributes_TELEPHONENR is the "telephoneNr" attribute (OID 1.3.6.1.4.1.1466.115.121.1.50).
	X520Attributes_TELEPHONENR X520Attributes = "TELEPHONENR"
	// X520Attributes_TELEXNR is the "telexNr" attribute (OID 1.3.6.1.4.1.1466.115.121.1.52).
	X520Attributes_TELEXNR X520Attributes = "TELEXNR"
	// X520Attributes_UTCTIME is the "utcTime" attribute (OID 1.3.6.1.4.1.1466.115.121.1.53).
	X520Attributes_UTCTIME X520Attributes = "UTCTIME"
	// X520Attributes_LDAPSYNTAXDESCRIPTION is the "ldapSyntaxDescription" attribute (OID 1.3.6.1.4.1.1466.115.121.1.54).
	X520Attributes_LDAPSYNTAXDESCRIPTION X520Attributes = "LDAPSYNTAXDESCRIPTION"
	// X520Attributes_SUBSTRINGASSERTION is the "substringAssertion" attribute (OID 1.3.6.1.4.1.1466.115.121.1.58).
	X520Attributes_SUBSTRINGASSERTION X520Attributes = "SUBSTRINGASSERTION"
	// X520Attributes_EMAIL_ADDRESS is the "emailAddress" attribute (OID 1.2.840.113549.1.9.1).
	X520Attributes_EMAIL_ADDRESS X520Attributes = "EMAIL_ADDRESS"
)

var x520AttributesData = map[X520Attributes]x520AttributesFields{
	X520Attributes_OBJECTCLASS:                                   {"objectClass", "2.5.4.0"},
	X520Attributes_ALIASEDENTRYNAME:                              {"aliasedEntryName", "2.5.4.1"},
	X520Attributes_ENCRYPTEDALIASEDENTRYNAME:                     {"encryptedAliasedEntryName", "2.5.4.1.2"},
	X520Attributes_KNOWLEDGEINFORMATION:                          {"knowledgeInformation", "2.5.4.2"},
	X520Attributes_COMMONNAME:                                    {"commonName", "2.5.4.3"},
	X520Attributes_ENCRYPTEDCOMMONNAME:                           {"encryptedCommonName", "2.5.4.3.2"},
	X520Attributes_SURNAME:                                       {"surname", "2.5.4.4"},
	X520Attributes_ENCRYPTEDSURNAME:                              {"encryptedSurname", "2.5.4.4.2"},
	X520Attributes_SERIALNUMBER:                                  {"serialNumber", "2.5.4.5"},
	X520Attributes_ENCRYPTEDSERIALNUMBER:                         {"encryptedSerialNumber", "2.5.4.5.2"},
	X520Attributes_COUNTRYNAME:                                   {"countryName", "2.5.4.6"},
	X520Attributes_ENCRYPTEDCOUNTRYNAME:                          {"encryptedCountryName", "2.5.4.6.2"},
	X520Attributes_LOCALITYNAME:                                  {"localityName", "2.5.4.7"},
	X520Attributes_ENCRYPTEDLOCALITYNAME:                         {"encryptedLocalityName", "2.5.4.7.2"},
	X520Attributes_COLLECTIVELOCALITYNAME:                        {"collectiveLocalityName", "2.5.4.7.1"},
	X520Attributes_ENCRYPTEDCOLLECTIVELOCALITYNAME:               {"encryptedCollectiveLocalityName", "2.5.4.7.1.2"},
	X520Attributes_STATEORPROVINCENAME:                           {"stateOrProvinceName", "2.5.4.8"},
	X520Attributes_ENCRYPTEDSTATEORPROVINCENAME:                  {"encryptedStateOrProvinceName", "2.5.4.8.2"},
	X520Attributes_COLLECTIVESTATEORPROVINCENAME:                 {"collectiveStateOrProvinceName", "2.5.4.8.1"},
	X520Attributes_ENCRYPTEDCOLLECTIVESTATEORPROVINCENAME:        {"encryptedCollectiveStateOrProvinceName", "2.5.4.8.1.2"},
	X520Attributes_STREETADDRESS:                                 {"streetAddress", "2.5.4.9"},
	X520Attributes_ENCRYPTEDSTREETADDRESS:                        {"encryptedStreetAddress", "2.5.4.9.2"},
	X520Attributes_COLLECTIVESTREETADDRESS:                       {"collectiveStreetAddress", "2.5.4.9.1"},
	X520Attributes_ENCRYPTEDCOLLECTIVESTREETADDRESS:              {"encryptedCollectiveStreetAddress", "2.5.4.9.1.2"},
	X520Attributes_ORGANIZATIONNAME:                              {"organizationName", "2.5.4.10"},
	X520Attributes_ENCRYPTEDORGANIZATIONNAME:                     {"encryptedOrganizationName", "2.5.4.10.2"},
	X520Attributes_COLLECTIVEORGANIZATIONNAME:                    {"collectiveOrganizationName", "2.5.4.10.1"},
	X520Attributes_ENCRYPTEDCOLLECTIVEORGANIZATIONNAME:           {"encryptedCollectiveOrganizationName", "2.5.4.10.1.2"},
	X520Attributes_ORGANIZATIONALUNITNAME:                        {"organizationalUnitName", "2.5.4.11"},
	X520Attributes_ENCRYPTEDORGANIZATIONALUNITNAME:               {"encryptedOrganizationalUnitName", "2.5.4.11.2"},
	X520Attributes_COLLECTIVEORGANIZATIONALUNITNAME:              {"collectiveOrganizationalUnitName", "2.5.4.11.1"},
	X520Attributes_ENCRYPTEDCOLLECTIVEORGANIZATIONALUNITNAM:      {"encryptedCollectiveOrganizationalUnitNam", "2.5.4.11.1.2"},
	X520Attributes_TITLE:                                         {"title", "2.5.4.12"},
	X520Attributes_ENCRYPTEDTITLE:                                {"encryptedTitle", "2.5.4.12.2"},
	X520Attributes_DESCRIPTION:                                   {"description", "2.5.4.13"},
	X520Attributes_ENCRYPTEDDESCRIPTION:                          {"encryptedDescription", "2.5.4.13.2"},
	X520Attributes_SEARCHGUIDE:                                   {"searchGuide", "2.5.4.14"},
	X520Attributes_ENCRYPTEDSEARCHGUIDE:                          {"encryptedSearchGuide", "2.5.4.14.2"},
	X520Attributes_BUSINESSCATEGORY:                              {"businessCategory", "2.5.4.15"},
	X520Attributes_ENCRYPTEDBUSINESSCATEGORY:                     {"encryptedBusinessCategory", "2.5.4.15.2"},
	X520Attributes_POSTALADDRESS:                                 {"postalAddress", "2.5.4.16"},
	X520Attributes_ENCRYPTEDPOSTALADDRESS:                        {"encryptedPostalAddress", "2.5.4.16.2"},
	X520Attributes_COLLECTIVEPOSTALADDRESS:                       {"collectivePostalAddress", "2.5.4.16.1"},
	X520Attributes_ENCRYPTEDCOLLECTIVEPOSTALADDRESS:              {"encryptedCollectivePostalAddress", "2.5.4.16.1.2"},
	X520Attributes_POSTALCODE:                                    {"postalCode", "2.5.4.17"},
	X520Attributes_ENCRYPTEDPOSTALCODE:                           {"encryptedPostalCode", "2.5.4.17.2"},
	X520Attributes_COLLECTIVEPOSTALCODE:                          {"collectivePostalCode", "2.5.4.17.1"},
	X520Attributes_ENCRYPTEDCOLLECTIVEPOSTALCODE:                 {"encryptedCollectivePostalCode", "2.5.4.17.1.2"},
	X520Attributes_POSTOFFICEBOX:                                 {"postOfficeBox", "2.5.4.18"},
	X520Attributes_COLLECTIVEPOSTOFFICEBOX:                       {"collectivePostOfficeBox", "2.5.4.18.1"},
	X520Attributes_ENCRYPTEDPOSTOFFICEBOX:                        {"encryptedPostOfficeBox", "2.5.4.18.2"},
	X520Attributes_ENCRYPTEDCOLLECTIVEPOSTOFFICEBOX:              {"encryptedCollectivePostOfficeBox", "2.5.4.18.1.2"},
	X520Attributes_PHYSICALDELIVERYOFFICENAME:                    {"physicalDeliveryOfficeName", "2.5.4.19"},
	X520Attributes_COLLECTIVEPHYSICALDELIVERYOFFICENAME:          {"collectivePhysicalDeliveryOfficeName", "2.5.4.19.1"},
	X520Attributes_ENCRYPTEDPHYSICALDELIVERYOFFICENAME:           {"encryptedPhysicalDeliveryOfficeName", "2.5.4.19.2"},
	X520Attributes_ENCRYPTEDCOLLECTIVEPHYSICALDELIVERYOFFICENAME: {"encryptedCollectivePhysicalDeliveryOfficeName", "2.5.4.19.1.2"},
	X520Attributes_TELEPHONENUMBER:                               {"telephoneNumber", "2.5.4.20"},
	X520Attributes_ENCRYPTEDTELEPHONENUMBER:                      {"encryptedTelephoneNumber", "2.5.4.20.2"},
	X520Attributes_COLLECTIVETELEPHONENUMBER:                     {"collectiveTelephoneNumber", "2.5.4.20.1"},
	X520Attributes_ENCRYPTEDCOLLECTIVETELEPHONENUMBER:            {"encryptedCollectiveTelephoneNumber", "2.5.4.20.1.2"},
	X520Attributes_TELEXNUMBER:                                   {"telexNumber", "2.5.4.21"},
	X520Attributes_ENCRYPTEDTELEXNUMBER:                          {"encryptedTelexNumber", "2.5.4.21.2"},
	X520Attributes_COLLECTIVETELEXNUMBER:                         {"collectiveTelexNumber", "2.5.4.21.1"},
	X520Attributes_ENCRYPTEDCOLLECTIVETELEXNUMBER:                {"encryptedCollectiveTelexNumber", "2.5.4.21.1.2"},
	X520Attributes_TELETEXTERMINALIDENTIFIER:                     {"teletexTerminalIdentifier", "2.5.4.22"},
	X520Attributes_ENCRYPTEDTELETEXTERMINALIDENTIFIER:            {"encryptedTeletexTerminalIdentifier", "2.5.4.22.2"},
	X520Attributes_COLLECTIVETELETEXTERMINALIDENTIFIER:           {"collectiveTeletexTerminalIdentifier", "2.5.4.22.1"},
	X520Attributes_ENCRYPTEDCOLLECTIVETELETEXTERMINALIDENTIFIER:  {"encryptedCollectiveTeletexTerminalIdentifier", "2.5.4.22.1.2"},
	X520Attributes_FACSIMILETELEPHONENUMBER:                      {"facsimileTelephoneNumber", "2.5.4.23"},
	X520Attributes_ENCRYPTEDFACSIMILETELEPHONENUMBER:             {"encryptedFacsimileTelephoneNumber", "2.5.4.23.2"},
	X520Attributes_COLLECTIVEFACSIMILETELEPHONENUMBER:            {"collectiveFacsimileTelephoneNumber", "2.5.4.23.1"},
	X520Attributes_ENCRYPTEDCOLLECTIVEFACSIMILETELEPHONENUMBER:   {"encryptedCollectiveFacsimileTelephoneNumber", "2.5.4.23.1.2"},
	X520Attributes_X121ADDRESS:                                   {"x121Address", "2.5.4.24"},
	X520Attributes_ENCRYPTEDX121ADDRESS:                          {"encryptedX121Address", "2.5.4.24.2"},
	X520Attributes_INTERNATIONALISDNNUMBER:                       {"internationalISDNNumber", "2.5.4.25"},
	X520Attributes_ENCRYPTEDINTERNATIONALISDNNUMBER:              {"encryptedInternationalISDNNumber", "2.5.4.25.2"},
	X520Attributes_COLLECTIVEINTERNATIONALISDNNUMBER:             {"collectiveInternationalISDNNumber", "2.5.4.25.1"},
	X520Attributes_ENCRYPTEDCOLLECTIVEINTERNATIONALISDNNUMBER:    {"encryptedCollectiveInternationalISDNNumber", "2.5.4.25.1.2"},
	X520Attributes_REGISTEREDADDRESS:                             {"registeredAddress", "2.5.4.26"},
	X520Attributes_ENCRYPTEDREGISTEREDADDRESS:                    {"encryptedRegisteredAddress", "2.5.4.26.2"},
	X520Attributes_DESTINATIONINDICATOR:                          {"destinationIndicator", "2.5.4.27"},
	X520Attributes_ENCRYPTEDDESTINATIONINDICATOR:                 {"encryptedDestinationIndicator", "2.5.4.27.2"},
	X520Attributes_PREFERREDDELIVERYMETHOD:                       {"preferredDeliveryMethod", "2.5.4.28"},
	X520Attributes_ENCRYPTEDPREFERREDDELIVERYMETHOD:              {"encryptedPreferredDeliveryMethod", "2.5.4.28.2"},
	X520Attributes_PRESENTATIONADDRESS:                           {"presentationAddress", "2.5.4.29"},
	X520Attributes_ENCRYPTEDPRESENTATIONADDRESS:                  {"encryptedPresentationAddress", "2.5.4.29.2"},
	X520Attributes_SUPPORTEDAPPLICATIONCONTEXT:                   {"supportedApplicationContext", "2.5.4.30"},
	X520Attributes_ENCRYPTEDSUPPORTEDAPPLICATIONCONTEXT:          {"encryptedSupportedApplicationContext", "2.5.4.30.2"},
	X520Attributes_MEMBER:                                        {"member", "2.5.4.31"},
	X520Attributes_ENCRYPTEDMEMBER:                               {"encryptedMember", "2.5.4.31.2"},
	X520Attributes_OWNER:                                         {"owner", "2.5.4.32"},
	X520Attributes_ENCRYPTEDOWNER:                                {"encryptedOwner", "2.5.4.32.2"},
	X520Attributes_ROLEOCCUPANT:                                  {"roleOccupant", "2.5.4.33"},
	X520Attributes_ENCRYPTEDROLEOCCUPANT:                         {"encryptedRoleOccupant", "2.5.4.33.2"},
	X520Attributes_SEEALSO:                                       {"seeAlso", "2.5.4.34"},
	X520Attributes_ENCRYPTEDSEEALSO:                              {"encryptedSeeAlso", "2.5.4.34.2"},
	X520Attributes_USERPASSWORD:                                  {"userPassword", "2.5.4.35"},
	X520Attributes_ENCRYPTEDUSERPASSWORD:                         {"encryptedUserPassword", "2.5.4.35.2"},
	X520Attributes_USERCERTIFICATE:                               {"userCertificate", "2.5.4.36"},
	X520Attributes_ENCRYPTEDUSERCERTIFICATE:                      {"encryptedUserCertificate", "2.5.4.36.2"},
	X520Attributes_CACERTIFICATE:                                 {"cACertificate", "2.5.4.37"},
	X520Attributes_ENCRYPTEDCACERTIFICATE:                        {"encryptedCACertificate", "2.5.4.37.2"},
	X520Attributes_AUTHORITYREVOCATIONLIST:                       {"authorityRevocationList", "2.5.4.38"},
	X520Attributes_ENCRYPTEDAUTHORITYREVOCATIONLIST:              {"encryptedAuthorityRevocationList", "2.5.4.38.2"},
	X520Attributes_CERTIFICATEREVOCATIONLIST:                     {"certificateRevocationList", "2.5.4.39"},
	X520Attributes_ENCRYPTEDCERTIFICATEREVOCATIONLIST:            {"encryptedCertificateRevocationList", "2.5.4.39.2"},
	X520Attributes_CROSSCERTIFICATEPAIR:                          {"crossCertificatePair", "2.5.4.40"},
	X520Attributes_ENCRYPTEDCROSSCERTIFICATEPAIR:                 {"encryptedCrossCertificatePair", "2.5.4.40.2"},
	X520Attributes_NAME:                                          {"name", "2.5.4.41"},
	X520Attributes_GIVENNAME:                                     {"givenName", "2.5.4.42"},
	X520Attributes_ENCRYPTEDGIVENNAME:                            {"encryptedGivenName", "2.5.4.42.2"},
	X520Attributes_INITIALS:                                      {"initials", "2.5.4.43"},
	X520Attributes_ENCRYPTEDINITIALS:                             {"encryptedInitials", "2.5.4.43.2"},
	X520Attributes_GENERATIONQUALIFIER:                           {"generationQualifier", "2.5.4.44"},
	X520Attributes_ENCRYPTEDGENERATIONQUALIFIER:                  {"encryptedGenerationQualifier", "2.5.4.44.2"},
	X520Attributes_UNIQUEIDENTIFIER:                              {"uniqueIdentifier", "2.5.4.45"},
	X520Attributes_ENCRYPTEDUNIQUEIDENTIFIER:                     {"encryptedUniqueIdentifier", "2.5.4.45.2"},
	X520Attributes_DNQUALIFIER:                                   {"dnQualifier", "2.5.4.46"},
	X520Attributes_ENCRYPTEDDNQUALIFIER:                          {"encryptedDnQualifier", "2.5.4.46.2"},
	X520Attributes_ENHANCEDSEARCHGUIDE:                           {"enhancedSearchGuide", "2.5.4.47"},
	X520Attributes_ENCRYPTEDENHANCEDSEARCHGUIDE:                  {"encryptedEnhancedSearchGuide", "2.5.4.47.2"},
	X520Attributes_PROTOCOLINFORMATION:                           {"protocolInformation", "2.5.4.48"},
	X520Attributes_ENCRYPTEDPROTOCOLINFORMATION:                  {"encryptedProtocolInformation", "2.5.4.48.2"},
	X520Attributes_DISTINGUISHEDNAME:                             {"distinguishedName", "2.5.4.49"},
	X520Attributes_ENCRYPTEDDISTINGUISHEDNAME:                    {"encryptedDistinguishedName", "2.5.4.49.2"},
	X520Attributes_UNIQUEMEMBER:                                  {"uniqueMember", "2.5.4.50"},
	X520Attributes_ENCRYPTEDUNIQUEMEMBER:                         {"encryptedUniqueMember", "2.5.4.50.2"},
	X520Attributes_HOUSEIDENTIFIER:                               {"houseIdentifier", "2.5.4.51"},
	X520Attributes_ENCRYPTEDHOUSEIDENTIFIER:                      {"encryptedHouseIdentifier", "2.5.4.51.2"},
	X520Attributes_SUPPORTEDALGORITHMS:                           {"supportedAlgorithms", "2.5.4.52"},
	X520Attributes_ENCRYPTEDSUPPORTEDALGORITHMS:                  {"encryptedSupportedAlgorithms", "2.5.4.52.2"},
	X520Attributes_DELTAREVOCATIONLIST:                           {"deltaRevocationList", "2.5.4.53"},
	X520Attributes_ENCRYPTEDDELTAREVOCATIONLIST:                  {"encryptedDeltaRevocationList", "2.5.4.53.2"},
	X520Attributes_DMDNAME:                                       {"dmdName", "2.5.4.54"},
	X520Attributes_ENCRYPTEDDMDNAME:                              {"encryptedDmdName", "2.5.4.54.2"},
	X520Attributes_CLEARANCE:                                     {"clearance", "2.5.4.55"},
	X520Attributes_ENCRYPTEDCLEARANCE:                            {"encryptedClearance", "2.5.4.55.2"},
	X520Attributes_DEFAULTDIRQOP:                                 {"defaultDirQop", "2.5.4.56"},
	X520Attributes_ENCRYPTEDDEFAULTDIRQOP:                        {"encryptedDefaultDirQop", "2.5.4.56.2"},
	X520Attributes_ATTRIBUTEINTEGRITYINFO:                        {"attributeIntegrityInfo", "2.5.4.57"},
	X520Attributes_ENCRYPTEDATTRIBUTEINTEGRITYINFO:               {"encryptedAttributeIntegrityInfo", "2.5.4.57.2"},
	X520Attributes_ATTRIBUTECERTIFICATE:                          {"attributeCertificate", "2.5.4.58"},
	X520Attributes_ENCRYPTEDATTRIBUTECERTIFICATE:                 {"encryptedAttributeCertificate", "2.5.4.58.2"},
	X520Attributes_ATTRIBUTECERTIFICATEREVOCATIONLIST:            {"attributeCertificateRevocationList", "2.5.4.59"},
	X520Attributes_ENCRYPTEDATTRIBUTECERTIFICATEREVOCATIONLIST:   {"encryptedAttributeCertificateRevocationList", "2.5.4.59.2"},
	X520Attributes_CONFKEYINFO:                                   {"confKeyInfo", "2.5.4.60"},
	X520Attributes_ENCRYPTEDCONFKEYINFO:                          {"encryptedConfKeyInfo", "2.5.4.60.2"},
	X520Attributes_AACERTIFICATE:                                 {"aACertificate", "2.5.4.61"},
	X520Attributes_ATTRIBUTEDESCRIPTORCERTIFICATE:                {"attributeDescriptorCertificate", "2.5.4.62"},
	X520Attributes_ATTRIBUTEAUTHORITYREVOCATIONLIST:              {"attributeAuthorityRevocationList", "2.5.4.63"},
	X520Attributes_FAMILY_INFORMATION:                            {"family-information", "2.5.4.64"},
	X520Attributes_PSEUDONYM:                                     {"pseudonym", "2.5.4.65"},
	X520Attributes_COMMUNICATIONSSERVICE:                         {"communicationsService", "2.5.4.66"},
	X520Attributes_COMMUNICATIONSNETWORK:                         {"communicationsNetwork", "2.5.4.67"},
	X520Attributes_CERTIFICATIONPRACTICESTMT:                     {"certificationPracticeStmt", "2.5.4.68"},
	X520Attributes_CERTIFICATEPOLICY:                             {"certificatePolicy", "2.5.4.69"},
	X520Attributes_PKIPATH:                                       {"pkiPath", "2.5.4.70"},
	X520Attributes_PRIVPOLICY:                                    {"privPolicy", "2.5.4.71"},
	X520Attributes_ROLE:                                          {"role", "2.5.4.72"},
	X520Attributes_DELEGATIONPATH:                                {"delegationPath", "2.5.4.73"},
	X520Attributes_PROTPRIVPOLICY:                                {"protPrivPolicy", "2.5.4.74"},
	X520Attributes_XMLPRIVILEGEINFO:                              {"xMLPrivilegeInfo", "2.5.4.75"},
	X520Attributes_XMLPRIVPOLICY:                                 {"xmlPrivPolicy", "2.5.4.76"},
	X520Attributes_UUIDPAIR:                                      {"uuidpair", "2.5.4.77"},
	X520Attributes_TAGOID:                                        {"tagOid", "2.5.4.78"},
	X520Attributes_UIIFORMAT:                                     {"uiiFormat", "2.5.4.79"},
	X520Attributes_UIIINURN:                                      {"uiiInUrn", "2.5.4.80"},
	X520Attributes_CONTENTURL:                                    {"contentUrl", "2.5.4.81"},
	X520Attributes_PERMISSION:                                    {"permission", "2.5.4.82"},
	X520Attributes_URI:                                           {"uri", "2.5.4.83"},
	X520Attributes_PWDATTRIBUTE:                                  {"pwdAttribute", "2.5.4.84"},
	X520Attributes_USERPWD:                                       {"userPwd", "2.5.4.85"},
	X520Attributes_URN:                                           {"urn", "2.5.4.86"},
	X520Attributes_URL:                                           {"url", "2.5.4.87"},
	X520Attributes_UTMCOORDINATES:                                {"utmCoordinates", "2.5.4.88"},
	X520Attributes_URNC:                                          {"urnC", "2.5.4.89"},
	X520Attributes_UII:                                           {"uii", "2.5.4.90"},
	X520Attributes_EPC:                                           {"epc", "2.5.4.91"},
	X520Attributes_TAGAFI:                                        {"tagAfi", "2.5.4.92"},
	X520Attributes_EPCFORMAT:                                     {"epcFormat", "2.5.4.93"},
	X520Attributes_EPCINURN:                                      {"epcInUrn", "2.5.4.94"},
	X520Attributes_LDAPURL:                                       {"ldapUrl", "2.5.4.95"},
	X520Attributes_TAGLOCATION:                                   {"tagLocation", "2.5.4.96"},
	X520Attributes_ORGANIZATIONIDENTIFIER:                        {"organizationIdentifier", "2.5.4.97"},
	X520Attributes_COUNTRYCODE3C:                                 {"countryCode3c", "2.5.4.98"},
	X520Attributes_COUNTRYCODE3N:                                 {"countryCode3n", "2.5.4.99"},
	X520Attributes_DNSNAME:                                       {"dnsName", "2.5.4.100"},
	X520Attributes_EEPKCERTIFICATREVOCATIONLIST:                  {"eepkCertificatRevocationList", "2.5.4.101"},
	X520Attributes_EEATTRCERTIFICATEREVOCATIONLIST:               {"eeAttrCertificateRevocationList", "2.5.4.102"},
	X520Attributes_USERPWDDESCRIPTION:                            {"userPwdDescription", "2.5.40.0"},
	X520Attributes_PWDVOCABULARYDESCRIPTION:                      {"pwdVocabularyDescription", "2.5.40.1"},
	X520Attributes_PWDALPHABETDESCRIPTION:                        {"pwdAlphabetDescription", "2.5.40.2"},
	X520Attributes_PWDENCALGDESCRIPTION:                          {"pwdEncAlgDescription", "2.5.40.3"},
	X520Attributes_UTMCOORDS:                                     {"utmCoords", "2.5.40.4"},
	X520Attributes_UIIFORM:                                       {"uiiForm", "2.5.40.5"},
	X520Attributes_EPCFORM:                                       {"epcForm", "2.5.40.6"},
	X520Attributes_COUNTRYSTRING3C:                               {"countryString3c", "2.5.40.7"},
	X520Attributes_COUNTRYSTRING3N:                               {"countryString3n", "2.5.40.8"},
	X520Attributes_DNSSTRING:                                     {"dnsString", "2.5.40.9"},
	X520Attributes_ATTRIBUTETYPEDESCRIPTION:                      {"attributeTypeDescription", "1.3.6.1.4.1.1466.115.121.1.3"},
	X520Attributes_BITSTRING:                                     {"bitString", "1.3.6.1.4.1.1466.115.121.1.6"},
	X520Attributes_BOOLEAN:                                       {"boolean", "1.3.6.1.4.1.1466.115.121.1.7"},
	X520Attributes_X509CERTIFICATE:                               {"x509Certificate", "1.3.6.1.4.1.1466.115.121.1.8"},
	X520Attributes_X509CERTIFICATELIST:                           {"x509CertificateList", "1.3.6.1.4.1.1466.115.121.1.9"},
	X520Attributes_X509CERTIFICATEPAIR:                           {"x509CertificatePair", "1.3.6.1.4.1.1466.115.121.1.10"},
	X520Attributes_COUNTRYSTRING:                                 {"countryString", "1.3.6.1.4.1.1466.115.121.1.11"},
	X520Attributes_DN:                                            {"dn", "1.3.6.1.4.1.1466.115.121.1.12"},
	X520Attributes_DELIVERYMETHOD:                                {"deliveryMethod", "1.3.6.1.4.1.1466.115.121.1.14"},
	X520Attributes_DIRECTORYSTRING:                               {"directoryString", "1.3.6.1.4.1.1466.115.121.1.15"},
	X520Attributes_DITCONTENTRULEDESCRIPTION:                     {"dITContentRuleDescription", "1.3.6.1.4.1.1466.115.121.1.16"},
	X520Attributes_DITSTRUCTURERULEDESCRIPTION:                   {"dITStructureRuleDescription", "1.3.6.1.4.1.1466.115.121.1.17"},
	X520Attributes_ENHANCEDGUIDE:                                 {"enhancedGuide", "1.3.6.1.4.1.1466.115.121.1.21"},
	X520Attributes_FACSIMILETELEPHONENR:                          {"facsimileTelephoneNr", "1.3.6.1.4.1.1466.115.121.1.22"},
	X520Attributes_FAX:                                           {"fax", "1.3.6.1.4.1.1466.115.121.1.23"},
	X520Attributes_GENERALIZEDTIME:                               {"generalizedTime", "1.3.6.1.4.1.1466.115.121.1.24"},
	X520Attributes_GUIDE:                                         {"guide", "1.3.6.1.4.1.1466.115.121.1.25"},
	X520Attributes_IA5STRING:                                     {"ia5String", "1.3.6.1.4.1.1466.115.121.1.26"},
	X520Attributes_INTEGER:                                       {"integer", "1.3.6.1.4.1.1466.115.121.1.27"},
	X520Attributes_JPEG:                                          {"jpeg", "1.3.6.1.4.1.1466.115.121.1.28"},
	X520Attributes_MATCHINGRULEDESCRIPTION:                       {"matchingRuleDescription", "1.3.6.1.4.1.1466.115.121.1.30"},
	X520Attributes_MATCHINGRULEUSEDESCRIPTION:                    {"matchingRuleUseDescription", "1.3.6.1.4.1.1466.115.121.1.31"},
	X520Attributes_NAMEANDOPTIONALUID:                            {"nameAndOptionalUID", "1.3.6.1.4.1.1466.115.121.1.34"},
	X520Attributes_NAMEFORMDESCRIPTION:                           {"nameFormDescription", "1.3.6.1.4.1.1466.115.121.1.35"},
	X520Attributes_NUMERICSTRING:                                 {"numericString", "1.3.6.1.4.1.1466.115.121.1.36"},
	X520Attributes_OBJECTCLASSDESCRIPTION:                        {"objectClassDescription", "1.3.6.1.4.1.1466.115.121.1.37"},
	X520Attributes_OID:                                           {"oid", "1.3.6.1.4.1.1466.115.121.1.38"},
	X520Attributes_OTHERMAILBOX:                                  {"otherMailbox", "1.3.6.1.4.1.1466.115.121.1.39"},
	X520Attributes_OCTETSTRING:                                   {"octetString", "1.3.6.1.4.1.1466.115.121.1.40"},
	X520Attributes_POSTALADDR:                                    {"postalAddr", "1.3.6.1.4.1.1466.115.121.1.41"},
	X520Attributes_PRESENTATIONADDR:                              {"presentationAddr", "1.3.6.1.4.1.1466.115.121.1.43"},
	X520Attributes_PRINTABLESTRING:                               {"printableString", "1.3.6.1.4.1.1466.115.121.1.44"},
	X520Attributes_SUBTREESPEC:                                   {"subtreeSpec", "1.3.6.1.4.1.1466.115.121.1.45"},
	X520Attributes_X509SUPPORTEDALGORITHM:                        {"x509SupportedAlgorithm", "1.3.6.1.4.1.1466.115.121.1.49"},
	X520Attributes_TELEPHONENR:                                   {"telephoneNr", "1.3.6.1.4.1.1466.115.121.1.50"},
	X520Attributes_TELEXNR:                                       {"telexNr", "1.3.6.1.4.1.1466.115.121.1.52"},
	X520Attributes_UTCTIME:                                       {"utcTime", "1.3.6.1.4.1.1466.115.121.1.53"},
	X520Attributes_LDAPSYNTAXDESCRIPTION:                         {"ldapSyntaxDescription", "1.3.6.1.4.1.1466.115.121.1.54"},
	X520Attributes_SUBSTRINGASSERTION:                            {"substringAssertion", "1.3.6.1.4.1.1466.115.121.1.58"},
	X520Attributes_EMAIL_ADDRESS:                                 {"emailAddress", "1.2.840.113549.1.9.1"},
}

func X520AttributesValues() []X520Attributes {
	return []X520Attributes{
		X520Attributes_OBJECTCLASS,
		X520Attributes_ALIASEDENTRYNAME,
		X520Attributes_ENCRYPTEDALIASEDENTRYNAME,
		X520Attributes_KNOWLEDGEINFORMATION,
		X520Attributes_COMMONNAME,
		X520Attributes_ENCRYPTEDCOMMONNAME,
		X520Attributes_SURNAME,
		X520Attributes_ENCRYPTEDSURNAME,
		X520Attributes_SERIALNUMBER,
		X520Attributes_ENCRYPTEDSERIALNUMBER,
		X520Attributes_COUNTRYNAME,
		X520Attributes_ENCRYPTEDCOUNTRYNAME,
		X520Attributes_LOCALITYNAME,
		X520Attributes_ENCRYPTEDLOCALITYNAME,
		X520Attributes_COLLECTIVELOCALITYNAME,
		X520Attributes_ENCRYPTEDCOLLECTIVELOCALITYNAME,
		X520Attributes_STATEORPROVINCENAME,
		X520Attributes_ENCRYPTEDSTATEORPROVINCENAME,
		X520Attributes_COLLECTIVESTATEORPROVINCENAME,
		X520Attributes_ENCRYPTEDCOLLECTIVESTATEORPROVINCENAME,
		X520Attributes_STREETADDRESS,
		X520Attributes_ENCRYPTEDSTREETADDRESS,
		X520Attributes_COLLECTIVESTREETADDRESS,
		X520Attributes_ENCRYPTEDCOLLECTIVESTREETADDRESS,
		X520Attributes_ORGANIZATIONNAME,
		X520Attributes_ENCRYPTEDORGANIZATIONNAME,
		X520Attributes_COLLECTIVEORGANIZATIONNAME,
		X520Attributes_ENCRYPTEDCOLLECTIVEORGANIZATIONNAME,
		X520Attributes_ORGANIZATIONALUNITNAME,
		X520Attributes_ENCRYPTEDORGANIZATIONALUNITNAME,
		X520Attributes_COLLECTIVEORGANIZATIONALUNITNAME,
		X520Attributes_ENCRYPTEDCOLLECTIVEORGANIZATIONALUNITNAM,
		X520Attributes_TITLE,
		X520Attributes_ENCRYPTEDTITLE,
		X520Attributes_DESCRIPTION,
		X520Attributes_ENCRYPTEDDESCRIPTION,
		X520Attributes_SEARCHGUIDE,
		X520Attributes_ENCRYPTEDSEARCHGUIDE,
		X520Attributes_BUSINESSCATEGORY,
		X520Attributes_ENCRYPTEDBUSINESSCATEGORY,
		X520Attributes_POSTALADDRESS,
		X520Attributes_ENCRYPTEDPOSTALADDRESS,
		X520Attributes_COLLECTIVEPOSTALADDRESS,
		X520Attributes_ENCRYPTEDCOLLECTIVEPOSTALADDRESS,
		X520Attributes_POSTALCODE,
		X520Attributes_ENCRYPTEDPOSTALCODE,
		X520Attributes_COLLECTIVEPOSTALCODE,
		X520Attributes_ENCRYPTEDCOLLECTIVEPOSTALCODE,
		X520Attributes_POSTOFFICEBOX,
		X520Attributes_COLLECTIVEPOSTOFFICEBOX,
		X520Attributes_ENCRYPTEDPOSTOFFICEBOX,
		X520Attributes_ENCRYPTEDCOLLECTIVEPOSTOFFICEBOX,
		X520Attributes_PHYSICALDELIVERYOFFICENAME,
		X520Attributes_COLLECTIVEPHYSICALDELIVERYOFFICENAME,
		X520Attributes_ENCRYPTEDPHYSICALDELIVERYOFFICENAME,
		X520Attributes_ENCRYPTEDCOLLECTIVEPHYSICALDELIVERYOFFICENAME,
		X520Attributes_TELEPHONENUMBER,
		X520Attributes_ENCRYPTEDTELEPHONENUMBER,
		X520Attributes_COLLECTIVETELEPHONENUMBER,
		X520Attributes_ENCRYPTEDCOLLECTIVETELEPHONENUMBER,
		X520Attributes_TELEXNUMBER,
		X520Attributes_ENCRYPTEDTELEXNUMBER,
		X520Attributes_COLLECTIVETELEXNUMBER,
		X520Attributes_ENCRYPTEDCOLLECTIVETELEXNUMBER,
		X520Attributes_TELETEXTERMINALIDENTIFIER,
		X520Attributes_ENCRYPTEDTELETEXTERMINALIDENTIFIER,
		X520Attributes_COLLECTIVETELETEXTERMINALIDENTIFIER,
		X520Attributes_ENCRYPTEDCOLLECTIVETELETEXTERMINALIDENTIFIER,
		X520Attributes_FACSIMILETELEPHONENUMBER,
		X520Attributes_ENCRYPTEDFACSIMILETELEPHONENUMBER,
		X520Attributes_COLLECTIVEFACSIMILETELEPHONENUMBER,
		X520Attributes_ENCRYPTEDCOLLECTIVEFACSIMILETELEPHONENUMBER,
		X520Attributes_X121ADDRESS,
		X520Attributes_ENCRYPTEDX121ADDRESS,
		X520Attributes_INTERNATIONALISDNNUMBER,
		X520Attributes_ENCRYPTEDINTERNATIONALISDNNUMBER,
		X520Attributes_COLLECTIVEINTERNATIONALISDNNUMBER,
		X520Attributes_ENCRYPTEDCOLLECTIVEINTERNATIONALISDNNUMBER,
		X520Attributes_REGISTEREDADDRESS,
		X520Attributes_ENCRYPTEDREGISTEREDADDRESS,
		X520Attributes_DESTINATIONINDICATOR,
		X520Attributes_ENCRYPTEDDESTINATIONINDICATOR,
		X520Attributes_PREFERREDDELIVERYMETHOD,
		X520Attributes_ENCRYPTEDPREFERREDDELIVERYMETHOD,
		X520Attributes_PRESENTATIONADDRESS,
		X520Attributes_ENCRYPTEDPRESENTATIONADDRESS,
		X520Attributes_SUPPORTEDAPPLICATIONCONTEXT,
		X520Attributes_ENCRYPTEDSUPPORTEDAPPLICATIONCONTEXT,
		X520Attributes_MEMBER,
		X520Attributes_ENCRYPTEDMEMBER,
		X520Attributes_OWNER,
		X520Attributes_ENCRYPTEDOWNER,
		X520Attributes_ROLEOCCUPANT,
		X520Attributes_ENCRYPTEDROLEOCCUPANT,
		X520Attributes_SEEALSO,
		X520Attributes_ENCRYPTEDSEEALSO,
		X520Attributes_USERPASSWORD,
		X520Attributes_ENCRYPTEDUSERPASSWORD,
		X520Attributes_USERCERTIFICATE,
		X520Attributes_ENCRYPTEDUSERCERTIFICATE,
		X520Attributes_CACERTIFICATE,
		X520Attributes_ENCRYPTEDCACERTIFICATE,
		X520Attributes_AUTHORITYREVOCATIONLIST,
		X520Attributes_ENCRYPTEDAUTHORITYREVOCATIONLIST,
		X520Attributes_CERTIFICATEREVOCATIONLIST,
		X520Attributes_ENCRYPTEDCERTIFICATEREVOCATIONLIST,
		X520Attributes_CROSSCERTIFICATEPAIR,
		X520Attributes_ENCRYPTEDCROSSCERTIFICATEPAIR,
		X520Attributes_NAME,
		X520Attributes_GIVENNAME,
		X520Attributes_ENCRYPTEDGIVENNAME,
		X520Attributes_INITIALS,
		X520Attributes_ENCRYPTEDINITIALS,
		X520Attributes_GENERATIONQUALIFIER,
		X520Attributes_ENCRYPTEDGENERATIONQUALIFIER,
		X520Attributes_UNIQUEIDENTIFIER,
		X520Attributes_ENCRYPTEDUNIQUEIDENTIFIER,
		X520Attributes_DNQUALIFIER,
		X520Attributes_ENCRYPTEDDNQUALIFIER,
		X520Attributes_ENHANCEDSEARCHGUIDE,
		X520Attributes_ENCRYPTEDENHANCEDSEARCHGUIDE,
		X520Attributes_PROTOCOLINFORMATION,
		X520Attributes_ENCRYPTEDPROTOCOLINFORMATION,
		X520Attributes_DISTINGUISHEDNAME,
		X520Attributes_ENCRYPTEDDISTINGUISHEDNAME,
		X520Attributes_UNIQUEMEMBER,
		X520Attributes_ENCRYPTEDUNIQUEMEMBER,
		X520Attributes_HOUSEIDENTIFIER,
		X520Attributes_ENCRYPTEDHOUSEIDENTIFIER,
		X520Attributes_SUPPORTEDALGORITHMS,
		X520Attributes_ENCRYPTEDSUPPORTEDALGORITHMS,
		X520Attributes_DELTAREVOCATIONLIST,
		X520Attributes_ENCRYPTEDDELTAREVOCATIONLIST,
		X520Attributes_DMDNAME,
		X520Attributes_ENCRYPTEDDMDNAME,
		X520Attributes_CLEARANCE,
		X520Attributes_ENCRYPTEDCLEARANCE,
		X520Attributes_DEFAULTDIRQOP,
		X520Attributes_ENCRYPTEDDEFAULTDIRQOP,
		X520Attributes_ATTRIBUTEINTEGRITYINFO,
		X520Attributes_ENCRYPTEDATTRIBUTEINTEGRITYINFO,
		X520Attributes_ATTRIBUTECERTIFICATE,
		X520Attributes_ENCRYPTEDATTRIBUTECERTIFICATE,
		X520Attributes_ATTRIBUTECERTIFICATEREVOCATIONLIST,
		X520Attributes_ENCRYPTEDATTRIBUTECERTIFICATEREVOCATIONLIST,
		X520Attributes_CONFKEYINFO,
		X520Attributes_ENCRYPTEDCONFKEYINFO,
		X520Attributes_AACERTIFICATE,
		X520Attributes_ATTRIBUTEDESCRIPTORCERTIFICATE,
		X520Attributes_ATTRIBUTEAUTHORITYREVOCATIONLIST,
		X520Attributes_FAMILY_INFORMATION,
		X520Attributes_PSEUDONYM,
		X520Attributes_COMMUNICATIONSSERVICE,
		X520Attributes_COMMUNICATIONSNETWORK,
		X520Attributes_CERTIFICATIONPRACTICESTMT,
		X520Attributes_CERTIFICATEPOLICY,
		X520Attributes_PKIPATH,
		X520Attributes_PRIVPOLICY,
		X520Attributes_ROLE,
		X520Attributes_DELEGATIONPATH,
		X520Attributes_PROTPRIVPOLICY,
		X520Attributes_XMLPRIVILEGEINFO,
		X520Attributes_XMLPRIVPOLICY,
		X520Attributes_UUIDPAIR,
		X520Attributes_TAGOID,
		X520Attributes_UIIFORMAT,
		X520Attributes_UIIINURN,
		X520Attributes_CONTENTURL,
		X520Attributes_PERMISSION,
		X520Attributes_URI,
		X520Attributes_PWDATTRIBUTE,
		X520Attributes_USERPWD,
		X520Attributes_URN,
		X520Attributes_URL,
		X520Attributes_UTMCOORDINATES,
		X520Attributes_URNC,
		X520Attributes_UII,
		X520Attributes_EPC,
		X520Attributes_TAGAFI,
		X520Attributes_EPCFORMAT,
		X520Attributes_EPCINURN,
		X520Attributes_LDAPURL,
		X520Attributes_TAGLOCATION,
		X520Attributes_ORGANIZATIONIDENTIFIER,
		X520Attributes_COUNTRYCODE3C,
		X520Attributes_COUNTRYCODE3N,
		X520Attributes_DNSNAME,
		X520Attributes_EEPKCERTIFICATREVOCATIONLIST,
		X520Attributes_EEATTRCERTIFICATEREVOCATIONLIST,
		X520Attributes_USERPWDDESCRIPTION,
		X520Attributes_PWDVOCABULARYDESCRIPTION,
		X520Attributes_PWDALPHABETDESCRIPTION,
		X520Attributes_PWDENCALGDESCRIPTION,
		X520Attributes_UTMCOORDS,
		X520Attributes_UIIFORM,
		X520Attributes_EPCFORM,
		X520Attributes_COUNTRYSTRING3C,
		X520Attributes_COUNTRYSTRING3N,
		X520Attributes_DNSSTRING,
		X520Attributes_ATTRIBUTETYPEDESCRIPTION,
		X520Attributes_BITSTRING,
		X520Attributes_BOOLEAN,
		X520Attributes_X509CERTIFICATE,
		X520Attributes_X509CERTIFICATELIST,
		X520Attributes_X509CERTIFICATEPAIR,
		X520Attributes_COUNTRYSTRING,
		X520Attributes_DN,
		X520Attributes_DELIVERYMETHOD,
		X520Attributes_DIRECTORYSTRING,
		X520Attributes_DITCONTENTRULEDESCRIPTION,
		X520Attributes_DITSTRUCTURERULEDESCRIPTION,
		X520Attributes_ENHANCEDGUIDE,
		X520Attributes_FACSIMILETELEPHONENR,
		X520Attributes_FAX,
		X520Attributes_GENERALIZEDTIME,
		X520Attributes_GUIDE,
		X520Attributes_IA5STRING,
		X520Attributes_INTEGER,
		X520Attributes_JPEG,
		X520Attributes_MATCHINGRULEDESCRIPTION,
		X520Attributes_MATCHINGRULEUSEDESCRIPTION,
		X520Attributes_NAMEANDOPTIONALUID,
		X520Attributes_NAMEFORMDESCRIPTION,
		X520Attributes_NUMERICSTRING,
		X520Attributes_OBJECTCLASSDESCRIPTION,
		X520Attributes_OID,
		X520Attributes_OTHERMAILBOX,
		X520Attributes_OCTETSTRING,
		X520Attributes_POSTALADDR,
		X520Attributes_PRESENTATIONADDR,
		X520Attributes_PRINTABLESTRING,
		X520Attributes_SUBTREESPEC,
		X520Attributes_X509SUPPORTEDALGORITHM,
		X520Attributes_TELEPHONENR,
		X520Attributes_TELEXNR,
		X520Attributes_UTCTIME,
		X520Attributes_LDAPSYNTAXDESCRIPTION,
		X520Attributes_SUBSTRINGASSERTION,
		X520Attributes_EMAIL_ADDRESS,
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
