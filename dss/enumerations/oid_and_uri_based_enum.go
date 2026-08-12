// Ported from dss-enumerations/.../OidAndUriBasedEnum.java (DSS 6.5.RC1).
package enumerations

// OidAndUriBasedEnum joins the attributes of OID and URI based enums.
type OidAndUriBasedEnum interface {
	OidBasedEnum
	UriBasedEnum
}
