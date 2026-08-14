// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pdf/PAdESConstants.java (DSS 6.5.RC1).
//
// eu.europa.esig.dss.pdf is the one Java package of dss-pades that landed in no s5b manifest
// (see pdf_object.go's header). PAdESConstantsSignaturePKCS7SubFilter and
// PAdESConstantsSignaturePKCS7SHA1SubFilter already exist in the landed pades_signature.go
// (SCOPE-TS/SIGN chunk), matching the PAdESConstants<Name> naming this file completes; they are
// not redeclared here. Everything else is a plain string constant, so the Java "utility class of
// public static final String fields" collapses to a single const block; Java field names map to
// Go identifiers by dropping underscores and title-casing each word, exactly as the two
// already-landed constants (SIGNATURE_PKCS7_SUBFILTER -> SignaturePKCS7SubFilter) do.
package pades

const (
	// PAdESConstantsSignatureType is the 'Sig' signature type. Port of SIGNATURE_TYPE.
	PAdESConstantsSignatureType = "Sig"
	// PAdESConstantsSignatureDefaultFilter is the 'Adobe.PPKLite' filter. Port of SIGNATURE_DEFAULT_FILTER.
	PAdESConstantsSignatureDefaultFilter = "Adobe.PPKLite"
	// PAdESConstantsSignatureDefaultSubFilter is the 'ETSI.CAdES.detached' SubFilter.
	// Port of SIGNATURE_DEFAULT_SUBFILTER.
	PAdESConstantsSignatureDefaultSubFilter = "ETSI.CAdES.detached"

	// PAdESConstantsTimestampType is the 'DocTimeStamp' signature type. Port of TIMESTAMP_TYPE.
	PAdESConstantsTimestampType = "DocTimeStamp"
	// PAdESConstantsTimestampDefaultFilter is the 'Adobe.PPKLite' filter. Port of TIMESTAMP_DEFAULT_FILTER.
	PAdESConstantsTimestampDefaultFilter = "Adobe.PPKLite"
	// PAdESConstantsTimestampDefaultSubFilter is the 'ETSI.RFC3161' SubFilter.
	// Port of TIMESTAMP_DEFAULT_SUBFILTER.
	PAdESConstantsTimestampDefaultSubFilter = "ETSI.RFC3161"

	// PAdESConstantsDssDictionaryName is 'DSS'. Port of DSS_DICTIONARY_NAME.
	PAdESConstantsDssDictionaryName = "DSS"
	// PAdESConstantsCertArrayNameDss is 'Certs'. Port of CERT_ARRAY_NAME_DSS.
	PAdESConstantsCertArrayNameDss = "Certs"
	// PAdESConstantsOcspArrayNameDss is 'OCSPs'. Port of OCSP_ARRAY_NAME_DSS.
	PAdESConstantsOcspArrayNameDss = "OCSPs"
	// PAdESConstantsCrlArrayNameDss is 'CRLs'. Port of CRL_ARRAY_NAME_DSS.
	PAdESConstantsCrlArrayNameDss = "CRLs"

	// PAdESConstantsVriDictionaryName is 'VRI'. Port of VRI_DICTIONARY_NAME.
	PAdESConstantsVriDictionaryName = "VRI"
	// PAdESConstantsCertArrayNameVri is 'Cert'. Port of CERT_ARRAY_NAME_VRI.
	PAdESConstantsCertArrayNameVri = "Cert"
	// PAdESConstantsOcspArrayNameVri is 'OCSP'. Port of OCSP_ARRAY_NAME_VRI.
	PAdESConstantsOcspArrayNameVri = "OCSP"
	// PAdESConstantsCrlArrayNameVri is 'CRL'. Port of CRL_ARRAY_NAME_VRI.
	PAdESConstantsCrlArrayNameVri = "CRL"

	// PAdESConstantsTuDictionaryNameVri is 'TU'. Port of TU_DICTIONARY_NAME_VRI.
	PAdESConstantsTuDictionaryNameVri = "TU"
	// PAdESConstantsTsDictionaryNameVri is 'TS'. Port of TS_DICTIONARY_NAME_VRI.
	PAdESConstantsTsDictionaryNameVri = "TS"

	// Field names.

	// PAdESConstantsAcroFormName is 'AcroForm'. Port of ACRO_FORM_NAME.
	PAdESConstantsAcroFormName = "AcroForm"
	// PAdESConstantsActionName is 'Action'. Port of ACTION_NAME.
	PAdESConstantsActionName = "Action"
	// PAdESConstantsAnnotFlag is 'F' (Annotation flag). Port of ANNOT_FLAG.
	PAdESConstantsAnnotFlag = "F"
	// PAdESConstantsAnnotsName is 'Annots'. Port of ANNOTS_NAME.
	PAdESConstantsAnnotsName = "Annots"
	// PAdESConstantsAppearanceDictionaryName is 'AP' (Appearance dictionary).
	// Port of APPEARANCE_DICTIONARY_NAME.
	PAdESConstantsAppearanceDictionaryName = "AP"
	// PAdESConstantsAsName is 'AS'. Port of AS_NAME.
	PAdESConstantsAsName = "AS"
	// PAdESConstantsByteRangeName is 'ByteRange'. Port of BYTE_RANGE_NAME.
	PAdESConstantsByteRangeName = "ByteRange"
	// PAdESConstantsCatalogName is 'Catalog'. Port of CATALOG_NAME.
	PAdESConstantsCatalogName = "Catalog"
	// PAdESConstantsContactInfoName is 'ContactInfo'. Port of CONTACT_INFO_NAME.
	PAdESConstantsContactInfoName = "ContactInfo"
	// PAdESConstantsContentsName is 'Contents'. Port of CONTENTS_NAME.
	PAdESConstantsContentsName = "Contents"
	// PAdESConstantsDataName is 'Data'. Port of DATA_NAME.
	PAdESConstantsDataName = "Data"
	// PAdESConstantsDocMdpName is 'DocMDP'. Port of DOC_MDP_NAME.
	PAdESConstantsDocMdpName = "DocMDP"
	// PAdESConstantsDocumentAppearanceName is 'DA' (Document-wide appearance).
	// Port of DOCUMENT_APPEARANCE_NAME.
	PAdESConstantsDocumentAppearanceName = "DA"
	// PAdESConstantsDocumentResourcesName is 'DR' (Document-wide resources).
	// Port of DOCUMENT_RESOURCES_NAME.
	PAdESConstantsDocumentResourcesName = "DR"
	// PAdESConstantsExtensionsName is 'Extensions'. Port of EXTENSIONS_NAME.
	PAdESConstantsExtensionsName = "Extensions"
	// PAdESConstantsBaseVersionName is 'BaseVersion'. Port of BASE_VERSION_NAME.
	PAdESConstantsBaseVersionName = "BaseVersion"
	// PAdESConstantsExtensionLevelName is 'ExtensionLevel'. Port of EXTENSION_LEVEL_NAME.
	PAdESConstantsExtensionLevelName = "ExtensionLevel"
	// PAdESConstantsExtensionRevisionName is 'ExtensionRevision'. Port of EXTENSION_REVISION_NAME.
	PAdESConstantsExtensionRevisionName = "ExtensionRevision"
	// PAdESConstantsURLName is 'URL'. Port of URL_NAME.
	PAdESConstantsURLName = "URL"
	// PAdESConstantsFieldMdpName is 'FieldMDP'. Port of FIELD_MDP_NAME.
	PAdESConstantsFieldMdpName = "FieldMDP"
	// PAdESConstantsFieldsName is 'Fields'. Port of FIELDS_NAME.
	PAdESConstantsFieldsName = "Fields"
	// PAdESConstantsFieldNameName is 'T' (Field name). Port of FIELD_NAME_NAME.
	PAdESConstantsFieldNameName = "T"
	// PAdESConstantsFilterName is 'Filter'. Port of FILTER_NAME.
	PAdESConstantsFilterName = "Filter"
	// PAdESConstantsFontName is 'Font'. Port of FONT_NAME.
	PAdESConstantsFontName = "Font"
	// PAdESConstantsItextName is 'ITXT' (iText identifier). Port of ITEXT_NAME.
	PAdESConstantsItextName = "ITXT"
	// PAdESConstantsKidsName is 'Kids'. Port of KIDS_NAME.
	PAdESConstantsKidsName = "Kids"
	// PAdESConstantsLengthName is 'Length'. Port of LENGTH_NAME.
	PAdESConstantsLengthName = "Length"
	// PAdESConstantsLocationName is 'Location'. Port of LOCATION_NAME.
	PAdESConstantsLocationName = "Location"
	// PAdESConstantsLockName is 'Lock'. Port of LOCK_NAME.
	PAdESConstantsLockName = "Lock"
	// PAdESConstantsMetadataName is 'Metadata'. Port of METADATA_NAME.
	PAdESConstantsMetadataName = "Metadata"
	// PAdESConstantsNormalAppearanceName is 'N'. Port of NORMAL_APPEARANCE_NAME.
	PAdESConstantsNormalAppearanceName = "N"
	// PAdESConstantsNameName is 'Name'. Port of NAME_NAME.
	PAdESConstantsNameName = "Name"
	// PAdESConstantsNamesName is 'Names'. Port of NAMES_NAME.
	PAdESConstantsNamesName = "Names"
	// PAdESConstantsOutputIntentsName is 'OutputIntents'. Port of OUTPUT_INTENTS_NAME.
	PAdESConstantsOutputIntentsName = "OutputIntents"
	// PAdESConstantsParentName is 'Parent'. Port of PARENT_NAME.
	PAdESConstantsParentName = "Parent"
	// PAdESConstantsPermissionsName is 'P' (Permissions). Port of PERMISSIONS_NAME.
	PAdESConstantsPermissionsName = "P"
	// PAdESConstantsPermsName is 'Perms'. Port of PERMS_NAME.
	PAdESConstantsPermsName = "Perms"
	// PAdESConstantsPieceInfoName is 'PieceInfo'. Port of PIECE_INFO_NAME.
	PAdESConstantsPieceInfoName = "PieceInfo"
	// PAdESConstantsReasonName is 'Reason'. Port of REASON_NAME.
	PAdESConstantsReasonName = "Reason"
	// PAdESConstantsReferenceName is 'Reference'. Port of REFERENCE_NAME.
	PAdESConstantsReferenceName = "Reference"
	// PAdESConstantsRootName is 'Root'. Port of ROOT_NAME.
	PAdESConstantsRootName = "Root"
	// PAdESConstantsSigningDateName is 'M' (Signing date). Port of SIGNING_DATE_NAME.
	PAdESConstantsSigningDateName = "M"
	// PAdESConstantsStreamName is 'stream'. Port of STREAM_NAME.
	PAdESConstantsStreamName = "stream"
	// PAdESConstantsSigFieldLockName is 'SigFieldLock'. Port of SIG_FIELD_LOCK_NAME.
	PAdESConstantsSigFieldLockName = "SigFieldLock"
	// PAdESConstantsSigFlagsName is 'SigFlags'. Port of SIG_FLAGS_NAME.
	PAdESConstantsSigFlagsName = "SigFlags"
	// PAdESConstantsSigRefName is 'SigRef'. Port of SIG_REF_NAME.
	PAdESConstantsSigRefName = "SigRef"
	// PAdESConstantsStructTreeRootName is 'StructTreeRoot'. Port of STRUCT_TREE_ROOT_NAME.
	PAdESConstantsStructTreeRootName = "StructTreeRoot"
	// PAdESConstantsStructTreeRootIDTreeName is 'IDTree'. Port of STRUCT_TREE_ROOT_ID_TREE_NAME.
	PAdESConstantsStructTreeRootIDTreeName = "IDTree"
	// PAdESConstantsStructTreeRootKName is 'K'. Port of STRUCT_TREE_ROOT_K_NAME.
	PAdESConstantsStructTreeRootKName = "K"
	// PAdESConstantsStructTreeRootParentTreeName is 'ParentTree'.
	// Port of STRUCT_TREE_ROOT_PARENT_TREE_NAME.
	PAdESConstantsStructTreeRootParentTreeName = "ParentTree"
	// PAdESConstantsStructTreeRootParentTreeNextKeyName is 'ParentTreeNextKey'.
	// Port of STRUCT_TREE_ROOT_PARENT_TREE_NEXT_KEY_NAME.
	PAdESConstantsStructTreeRootParentTreeNextKeyName = "ParentTreeNextKey"
	// PAdESConstantsSubFilterName is 'SubFilter'. Port of SUB_FILTER_NAME.
	PAdESConstantsSubFilterName = "SubFilter"
	// PAdESConstantsTypeName is 'Type'. Port of TYPE_NAME.
	PAdESConstantsTypeName = "Type"
	// PAdESConstantsTransformMethodName is 'TransformMethod'. Port of TRANSFORM_METHOD_NAME.
	PAdESConstantsTransformMethodName = "TransformMethod"
	// PAdESConstantsTransformParamsName is 'TransformParams'. Port of TRANSFORM_PARAMS_NAME.
	PAdESConstantsTransformParamsName = "TransformParams"
	// PAdESConstantsURName is 'UR' (User rights). Port of UR_NAME.
	PAdESConstantsURName = "UR"
	// PAdESConstantsUR3Name is 'UR3' (User rights). Port of UR3_NAME.
	PAdESConstantsUR3Name = "UR3"
	// PAdESConstantsValueName is 'V' (Value). Port of VALUE_NAME.
	PAdESConstantsValueName = "V"
	// PAdESConstantsVersionName is 'Version'. Port of VERSION_NAME.
	PAdESConstantsVersionName = "Version"

	// Build properties dictionary.

	// PAdESConstantsApp is 'App' (Application software). Port of APP.
	PAdESConstantsApp = "App"
	// PAdESConstantsPropBuild is 'Prop_Build' (Build properties). Port of PROP_BUILD.
	PAdESConstantsPropBuild = "Prop_Build"
	// PAdESConstantsVersionDefault is 'V=1.2'. Port of VERSION_DEFAULT.
	PAdESConstantsVersionDefault = "1.2"
)
