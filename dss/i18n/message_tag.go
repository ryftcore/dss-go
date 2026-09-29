// Ported from dss-i18n/.../i18n/MessageTag.java (DSS 6.5.RC1).

package i18n

import "fmt"

// MessageTag contains a list of possible message tags.
// NOTE: all message tags shall be listed in the dss-messages.properties file.
type MessageTag string

// MessageTag constants, one per message key in dss-messages.properties.
// Each corresponds to a single ETSI EN 319 102-1 validation building-block
// check (e.g. "BBB_FC_IEFF" is a format-checking sub-check); the "_ANS"
// suffix (and numbered "_ANS1", "_ANS2", ...) marks the negative-answer
// variant of the message immediately preceding it. Individual constants are
// intentionally undocumented beyond their name — see dss-messages.properties
// (embedded via Provider) for the exact human-readable text each
// resolves to.
const (
	MessageTagBBBFCIEFF                                 MessageTag = "BBB_FC_IEFF"
	MessageTagBBBFCIEFFANS                              MessageTag = "BBB_FC_IEFF_ANS"
	MessageTagBBBFCICFD                                 MessageTag = "BBB_FC_ICFD"
	MessageTagBBBFCICFDANS                              MessageTag = "BBB_FC_ICFD_ANS"
	MessageTagBBBFCISD                                  MessageTag = "BBB_FC_ISD"
	MessageTagBBBFCISDANS                               MessageTag = "BBB_FC_ISD_ANS"
	MessageTagBBBFCISRIA                                MessageTag = "BBB_FC_ISRIA"
	MessageTagBBBFCISRIAANS                             MessageTag = "BBB_FC_ISRIA_ANS"
	MessageTagBBBFCIOSIP                                MessageTag = "BBB_FC_IOSIP"
	MessageTagBBBFCIOSIPANS                             MessageTag = "BBB_FC_IOSIP_ANS"
	MessageTagBBBFCDASTHVBR                             MessageTag = "BBB_FC_DASTHVBR"
	MessageTagBBBFCDASTHVBRANS                          MessageTag = "BBB_FC_DASTHVBR_ANS"
	MessageTagBBBFCDBTOOST                              MessageTag = "BBB_FC_DBTOOST"
	MessageTagBBBFCDBTOOSTANS                           MessageTag = "BBB_FC_DBTOOST_ANS"
	MessageTagBBBFCIBRV                                 MessageTag = "BBB_FC_IBRV"
	MessageTagBBBFCIBRVANS                              MessageTag = "BBB_FC_IBRV_ANS"
	MessageTagBBBFCISDC                                 MessageTag = "BBB_FC_ISDC"
	MessageTagBBBFCISDCANS                              MessageTag = "BBB_FC_ISDC_ANS"
	MessageTagBBBFCDSFREAP                              MessageTag = "BBB_FC_DSFREAP"
	MessageTagBBBFCDSFREAPANS                           MessageTag = "BBB_FC_DSFREAP_ANS"
	MessageTagBBBFCIAOD                                 MessageTag = "BBB_FC_IAOD"
	MessageTagBBBFCIAODANS                              MessageTag = "BBB_FC_IAOD_ANS"
	MessageTagBBBFCIVDBSFR                              MessageTag = "BBB_FC_IVDBSFR"
	MessageTagBBBFCIVDBSFRANS                           MessageTag = "BBB_FC_IVDBSFR_ANS"
	MessageTagBBBFCISVADMDPD                            MessageTag = "BBB_FC_ISVADMDPD"
	MessageTagBBBFCISVADMDPDANS                         MessageTag = "BBB_FC_ISVADMDPD_ANS"
	MessageTagBBBFCISVAFMDPD                            MessageTag = "BBB_FC_ISVAFMDPD"
	MessageTagBBBFCISVAFMDPDANS                         MessageTag = "BBB_FC_ISVAFMDPD_ANS"
	MessageTagBBBFCISVASFLD                             MessageTag = "BBB_FC_ISVASFLD"
	MessageTagBBBFCISVASFLDANS                          MessageTag = "BBB_FC_ISVASFLD_ANS"
	MessageTagBBBFCDSCNFFSM                             MessageTag = "BBB_FC_DSCNFFSM"
	MessageTagBBBFCDSCNFFSMANS                          MessageTag = "BBB_FC_DSCNFFSM_ANS"
	MessageTagBBBFCDSCNACMDM                            MessageTag = "BBB_FC_DSCNACMDM"
	MessageTagBBBFCDSCNACMDMANS                         MessageTag = "BBB_FC_DSCNACMDM_ANS"
	MessageTagBBBFCDSCNUOM                              MessageTag = "BBB_FC_DSCNUOM"
	MessageTagBBBFCDSCNUOMANS                           MessageTag = "BBB_FC_DSCNUOM_ANS"
	MessageTagBBBFCIECKSCDA                             MessageTag = "BBB_FC_IECKSCDA"
	MessageTagBBBFCIECKSCDAANS1                         MessageTag = "BBB_FC_IECKSCDA_ANS1"
	MessageTagBBBFCIECKSCDAANS2                         MessageTag = "BBB_FC_IECKSCDA_ANS2"
	MessageTagBBBFCIECKSCDAANS3                         MessageTag = "BBB_FC_IECKSCDA_ANS3"
	MessageTagBBBFCIECKSCDAANS4                         MessageTag = "BBB_FC_IECKSCDA_ANS4"
	MessageTagBBBFCIECKSCDAANS5                         MessageTag = "BBB_FC_IECKSCDA_ANS5"
	MessageTagBBBFCDDAPDFAF                             MessageTag = "BBB_FC_DDAPDFAF"
	MessageTagBBBFCDDAPDFAFANS                          MessageTag = "BBB_FC_DDAPDFAF_ANS"
	MessageTagBBBFCIDPDFAC                              MessageTag = "BBB_FC_IDPDFAC"
	MessageTagBBBFCIDPDFACANS                           MessageTag = "BBB_FC_IDPDFAC_ANS"
	MessageTagBBBFCIECTF                                MessageTag = "BBB_FC_IECTF"
	MessageTagBBBFCIECTFANS                             MessageTag = "BBB_FC_IECTF_ANS"
	MessageTagBBBFCISFCS                                MessageTag = "BBB_FC_ISFCS"
	MessageTagBBBFCISFCSANS                             MessageTag = "BBB_FC_ISFCS_ANS"
	MessageTagBBBFCITFCS                                MessageTag = "BBB_FC_ITFCS"
	MessageTagBBBFCITFCSANS                             MessageTag = "BBB_FC_ITFCS_ANS"
	MessageTagBBBFCIMFCS                                MessageTag = "BBB_FC_IMFCS"
	MessageTagBBBFCIMFCSANS                             MessageTag = "BBB_FC_IMFCS_ANS"
	MessageTagBBBFCITZCP                                MessageTag = "BBB_FC_ITZCP"
	MessageTagBBBFCITZCPANS                             MessageTag = "BBB_FC_ITZCP_ANS"
	MessageTagBBBFCITEZCF                               MessageTag = "BBB_FC_ITEZCF"
	MessageTagBBBFCITEZCFANS                            MessageTag = "BBB_FC_ITEZCF_ANS"
	MessageTagBBBFCITMFP                                MessageTag = "BBB_FC_ITMFP"
	MessageTagBBBFCITMFPANS                             MessageTag = "BBB_FC_ITMFP_ANS"
	MessageTagBBBFCIEMCF                                MessageTag = "BBB_FC_IEMCF"
	MessageTagBBBFCIEMCFANS                             MessageTag = "BBB_FC_IEMCF_ANS"
	MessageTagBBBFCIMFPASiCE                            MessageTag = "BBB_FC_IMFP_ASICE"
	MessageTagBBBFCIMFPASiCEANS                         MessageTag = "BBB_FC_IMFP_ASICE_ANS"
	MessageTagBBBFCISFPASiCE                            MessageTag = "BBB_FC_ISFP_ASICE"
	MessageTagBBBFCISFPASiCEANS                         MessageTag = "BBB_FC_ISFP_ASICE_ANS"
	MessageTagBBBFCISFPASiCS                            MessageTag = "BBB_FC_ISFP_ASICS"
	MessageTagBBBFCISFPASiCSANS                         MessageTag = "BBB_FC_ISFP_ASICS_ANS"
	MessageTagBBBFCISFPASTFORAMC                        MessageTag = "BBB_FC_ISFP_ASTFORAMC"
	MessageTagBBBFCISFPASTFORAMCANS                     MessageTag = "BBB_FC_ISFP_ASTFORAMC_ANS"
	MessageTagBBBFCIAHIV                                MessageTag = "BBB_FC_IAHIV"
	MessageTagBBBFCIAHIVANS                             MessageTag = "BBB_FC_IAHIV_ANS"
	MessageTagBBBCVIRDOF                                MessageTag = "BBB_CV_IRDOF"
	MessageTagBBBCVIRDOFANS                             MessageTag = "BBB_CV_IRDOF_ANS"
	MessageTagBBBCVTSPIRDOF                             MessageTag = "BBB_CV_TSP_IRDOF"
	MessageTagBBBCVTSPIRDOFANS                          MessageTag = "BBB_CV_TSP_IRDOF_ANS"
	MessageTagBBBCVCSCSSVF                              MessageTag = "BBB_CV_CS_CSSVF"
	MessageTagBBBCVCSCSSVFANS                           MessageTag = "BBB_CV_CS_CSSVF_ANS"
	MessageTagBBBCVIMEOF                                MessageTag = "BBB_CV_IMEOF"
	MessageTagBBBCVIMEOFANS                             MessageTag = "BBB_CV_IMEOF_ANS"
	MessageTagBBBCVEAASDCBF                             MessageTag = "BBB_CV_EAA_SDCBF"
	MessageTagBBBCVEAASDCBFANS                          MessageTag = "BBB_CV_EAA_SDCBF_ANS"
	MessageTagBBBCVEAANSDCBF                            MessageTag = "BBB_CV_EAA_NSDCBF"
	MessageTagBBBCVEAANSDCBFANS                         MessageTag = "BBB_CV_EAA_NSDCBF_ANS"
	MessageTagBBBCVERIODOF                              MessageTag = "BBB_CV_ER_IODOF"
	MessageTagBBBCVERIODOFANS                           MessageTag = "BBB_CV_ER_IODOF_ANS"
	MessageTagBBBCVERHasSDoc                            MessageTag = "BBB_CV_ER_HASSDOC"
	MessageTagBBBCVERHasSDocANS                         MessageTag = "BBB_CV_ER_HASSDOC_ANS"
	MessageTagBBBCVERDFHVLCDOG                          MessageTag = "BBB_CV_ER_DFHVLCDOG"
	MessageTagBBBCVERDFHVLCDOGANS                       MessageTag = "BBB_CV_ER_DFHVLCDOG_ANS"
	MessageTagBBBCVERATSRF                              MessageTag = "BBB_CV_ER_ATSRF"
	MessageTagBBBCVERATSRFANS                           MessageTag = "BBB_CV_ER_ATSRF_ANS"
	MessageTagBBBCVERATSSRF                             MessageTag = "BBB_CV_ER_ATSSRF"
	MessageTagBBBCVERATSSRFANS                          MessageTag = "BBB_CV_ER_ATSSRF_ANS"
	MessageTagBBBCVIRDOI                                MessageTag = "BBB_CV_IRDOI"
	MessageTagBBBCVIRDOIANS                             MessageTag = "BBB_CV_IRDOI_ANS"
	MessageTagBBBCVTSPIRDOI                             MessageTag = "BBB_CV_TSP_IRDOI"
	MessageTagBBBCVTSPIRDOIANS                          MessageTag = "BBB_CV_TSP_IRDOI_ANS"
	MessageTagBBBCVCSCSPS                               MessageTag = "BBB_CV_CS_CSPS"
	MessageTagBBBCVCSCSPSANS                            MessageTag = "BBB_CV_CS_CSPS_ANS"
	MessageTagBBBCVIMEDOI                               MessageTag = "BBB_CV_IMEDOI"
	MessageTagBBBCVIMEDOIANS                            MessageTag = "BBB_CV_IMEDOI_ANS"
	MessageTagBBBCVEAASDCBI                             MessageTag = "BBB_CV_EAA_SDCBI"
	MessageTagBBBCVEAASDCBIANS                          MessageTag = "BBB_CV_EAA_SDCBI_ANS"
	MessageTagBBBCVEAANSDCBI                            MessageTag = "BBB_CV_EAA_NSDCBI"
	MessageTagBBBCVEAANSDCBIANS                         MessageTag = "BBB_CV_EAA_NSDCBI_ANS"
	MessageTagBBBCVERATSRI                              MessageTag = "BBB_CV_ER_ATSRI"
	MessageTagBBBCVERATSRIANS                           MessageTag = "BBB_CV_ER_ATSRI_ANS"
	MessageTagBBBCVERATSSRI                             MessageTag = "BBB_CV_ER_ATSSRI"
	MessageTagBBBCVERATSSRIANS                          MessageTag = "BBB_CV_ER_ATSSRI_ANS"
	MessageTagBBBCVISMEC                                MessageTag = "BBB_CV_ISMEC"
	MessageTagBBBCVISMECANS                             MessageTag = "BBB_CV_ISMEC_ANS"
	MessageTagBBBCVISMECANS2                            MessageTag = "BBB_CV_ISMEC_ANS_2"
	MessageTagBBBCVAAMEF                                MessageTag = "BBB_CV_AAMEF"
	MessageTagBBBCVAAMEFANS                             MessageTag = "BBB_CV_AAMEF_ANS"
	MessageTagBBBCVERTSTRN                              MessageTag = "BBB_CV_ER_TST_RN"
	MessageTagBBBCVERTSTRNANS1                          MessageTag = "BBB_CV_ER_TST_RN_ANS_1"
	MessageTagBBBCVERTSTRNANS2                          MessageTag = "BBB_CV_ER_TST_RN_ANS_2"
	MessageTagBBBCVDRNMND                               MessageTag = "BBB_CV_DRNMND"
	MessageTagBBBCVDRNMNDANS                            MessageTag = "BBB_CV_DRNMND_ANS"
	MessageTagBBBCVDMENMND                              MessageTag = "BBB_CV_DMENMND"
	MessageTagBBBCVDMENMNDANS                           MessageTag = "BBB_CV_DMENMND_ANS"
	MessageTagBBBCVISI                                  MessageTag = "BBB_CV_ISI"
	MessageTagBBBCVISIC                                 MessageTag = "BBB_CV_ISIC"
	MessageTagBBBCVISIR                                 MessageTag = "BBB_CV_ISIR"
	MessageTagBBBCVISIT                                 MessageTag = "BBB_CV_ISIT"
	MessageTagBBBCVISIANS                               MessageTag = "BBB_CV_ISI_ANS"
	MessageTagBBBCVIAFS                                 MessageTag = "BBB_CV_IAFS"
	MessageTagBBBCVIAFSANS                              MessageTag = "BBB_CV_IAFS_ANS"
	MessageTagBBBICSISCI                                MessageTag = "BBB_ICS_ISCI"
	MessageTagBBBICSISCIANS                             MessageTag = "BBB_ICS_ISCI_ANS"
	MessageTagBBBICSISASCP                              MessageTag = "BBB_ICS_ISASCP"
	MessageTagBBBICSISASCPANS                           MessageTag = "BBB_ICS_ISASCP_ANS"
	MessageTagBBBICSISASCPU                             MessageTag = "BBB_ICS_ISASCPU"
	MessageTagBBBICSISASCPUANS                          MessageTag = "BBB_ICS_ISASCPU_ANS"
	MessageTagBBBICSISACDP                              MessageTag = "BBB_ICS_ISACDP"
	MessageTagBBBICSISACDPANS                           MessageTag = "BBB_ICS_ISACDP_ANS"
	MessageTagBBBICSICDVV                               MessageTag = "BBB_ICS_ICDVV"
	MessageTagBBBICSICDVVANS                            MessageTag = "BBB_ICS_ICDVV_ANS"
	MessageTagBBBICSICDVVS                              MessageTag = "BBB_ICS_ICDVVS"
	MessageTagBBBICSICDVVSANS                           MessageTag = "BBB_ICS_ICDVVS_ANS"
	MessageTagBBBICSAIDNASNE                            MessageTag = "BBB_ICS_AIDNASNE"
	MessageTagBBBICSAIDNASNEANS                         MessageTag = "BBB_ICS_AIDNASNE_ANS"
	MessageTagBBBICSISAKIDP                             MessageTag = "BBB_ICS_ISAKIDP"
	MessageTagBBBICSISAKIDPANS                          MessageTag = "BBB_ICS_ISAKIDP_ANS"
	MessageTagBBBICSDKIDVM                              MessageTag = "BBB_ICS_DKIDVM"
	MessageTagBBBICSDKIDVMANS                           MessageTag = "BBB_ICS_DKIDVM_ANS"
	MessageTagBBBICSISAX509UP                           MessageTag = "BBB_ICS_ISAX509UP"
	MessageTagBBBICSISAX509UPANS                        MessageTag = "BBB_ICS_ISAX509UP_ANS"
	MessageTagBBBICSISAX509UA                           MessageTag = "BBB_ICS_ISAX509UA"
	MessageTagBBBICSISAX509UAANS                        MessageTag = "BBB_ICS_ISAX509UA_ANS"
	MessageTagBBBRFCNUP                                 MessageTag = "BBB_RFC_NUP"
	MessageTagBBBRFCNUPANS                              MessageTag = "BBB_RFC_NUP_ANS"
	MessageTagBBBRFCIRIF                                MessageTag = "BBB_RFC_IRIF"
	MessageTagBBBRFCIRIFTUNU                            MessageTag = "BBB_RFC_IRIF_TUNU"
	MessageTagBBBRFCIRIFANS                             MessageTag = "BBB_RFC_IRIF_ANS"
	MessageTagADESTROBVPIIC                             MessageTag = "ADEST_ROBVPIIC"
	MessageTagADESTROBVPIICANS                          MessageTag = "ADEST_ROBVPIIC_ANS"
	MessageTagADESTROTVPIIC                             MessageTag = "ADEST_ROTVPIIC"
	MessageTagADESTROTVPIICANS                          MessageTag = "ADEST_ROTVPIIC_ANS"
	MessageTagADESTRORPIIC                              MessageTag = "ADEST_RORPIIC"
	MessageTagADESTRORPIICANS                           MessageTag = "ADEST_RORPIIC_ANS"
	MessageTagBSVIFCRC                                  MessageTag = "BSV_IFCRC"
	MessageTagBSVIFCRCANS                               MessageTag = "BSV_IFCRC_ANS"
	MessageTagBSVIISCRC                                 MessageTag = "BSV_IISCRC"
	MessageTagBSVIISCRCANS                              MessageTag = "BSV_IISCRC_ANS"
	MessageTagBSVIVCIRC                                 MessageTag = "BSV_IVCIRC"
	MessageTagBSVIVCIRCANS                              MessageTag = "BSV_IVCIRC_ANS"
	MessageTagBSVIXCVRC                                 MessageTag = "BSV_IXCVRC"
	MessageTagBSVIXCVRCANS                              MessageTag = "BSV_IXCVRC_ANS"
	MessageTagBSVISCRAVTC                               MessageTag = "BSV_ISCRAVTC"
	MessageTagBSVISCRAVTCANS                            MessageTag = "BSV_ISCRAVTC_ANS"
	MessageTagBSVIVTAVRSC                               MessageTag = "BSV_IVTAVRSC"
	MessageTagBSVIVTAVRSCANS                            MessageTag = "BSV_IVTAVRSC_ANS"
	MessageTagBSVISCCTC                                 MessageTag = "BSV_ISCCTC"
	MessageTagBSVISCCTCANS                              MessageTag = "BSV_ISCCTC_ANS"
	MessageTagBSVICTGTNASCRT                            MessageTag = "BSV_ICTGTNASCRT"
	MessageTagBSVICTGTNASCRTANS                         MessageTag = "BSV_ICTGTNASCRT_ANS"
	MessageTagBSVICTGTNASCET                            MessageTag = "BSV_ICTGTNASCET"
	MessageTagBSVICTGTNASCETANS                         MessageTag = "BSV_ICTGTNASCET_ANS"
	MessageTagBSVICVRC                                  MessageTag = "BSV_ICVRC"
	MessageTagBSVICVRCANS                               MessageTag = "BSV_ICVRC_ANS"
	MessageTagBSVISAVRC                                 MessageTag = "BSV_ISAVRC"
	MessageTagBSVISAVRCANS                              MessageTag = "BSV_ISAVRC_ANS"
	MessageTagBSVICTGTNACCET                            MessageTag = "BSV_ICTGTNACCET"
	MessageTagBSVICTGTNACCETANS                         MessageTag = "BSV_ICTGTNACCET_ANS"
	MessageTagBSVIEAAAVRC                               MessageTag = "BSV_IEAAAVRC"
	MessageTagBSVIEAAAVRCANS                            MessageTag = "BSV_IEAAAVRC_ANS"
	MessageTagLTVABSV                                   MessageTag = "LTV_ABSV"
	MessageTagLTVABSVANS                                MessageTag = "LTV_ABSV_ANS"
	MessageTagLTVISCKNR                                 MessageTag = "LTV_ISCKNR"
	MessageTagLTVISCKNRANS0                             MessageTag = "LTV_ISCKNR_ANS0"
	MessageTagLTVISCKNRANS1                             MessageTag = "LTV_ISCKNR_ANS1"
	MessageTagArchLTVV                                  MessageTag = "ARCH_LTVV"
	MessageTagArchLTVVANS                               MessageTag = "ARCH_LTVV_ANS"
	MessageTagArchLTAIVMP                               MessageTag = "ARCH_LTAIVMP"
	MessageTagArchLTAIVMPANS                            MessageTag = "ARCH_LTAIVMP_ANS"
	MessageTagArchIRTVBBA                               MessageTag = "ARCH_IRTVBBA"
	MessageTagArchIRTVBBAANS                            MessageTag = "ARCH_IRTVBBA_ANS"
	MessageTagArchICHFCRLPOET                           MessageTag = "ARCH_ICHFCRLPOET"
	MessageTagArchICHFCRLPOETANS                        MessageTag = "ARCH_ICHFCRLPOET_ANS"
	MessageTagACCM                                      MessageTag = "ACCM"
	MessageTagACCMANS                                   MessageTag = "ACCM_ANS"
	MessageTagASCCMCAA                                  MessageTag = "ASCCM_CAA"
	MessageTagASCCMCAAANS                               MessageTag = "ASCCM_CAA_ANS"
	MessageTagASCCMDAA                                  MessageTag = "ASCCM_DAA"
	MessageTagASCCMDAAANS                               MessageTag = "ASCCM_DAA_ANS"
	MessageTagASCCMDAAANS2                              MessageTag = "ASCCM_DAA_ANS_2"
	MessageTagASCCMAPKSA                                MessageTag = "ASCCM_APKSA"
	MessageTagASCCMAPKSAANS                             MessageTag = "ASCCM_APKSA_ANS"
	MessageTagASCCMAPKSAANS2                            MessageTag = "ASCCM_APKSA_ANS_2"
	MessageTagASCCMAR                                   MessageTag = "ASCCM_AR"
	MessageTagASCCMARANSANR                             MessageTag = "ASCCM_AR_ANS_ANR"
	MessageTagASCCMARANSANR2                            MessageTag = "ASCCM_AR_ANS_ANR_2"
	MessageTagASCCMARANSAKSNR                           MessageTag = "ASCCM_AR_ANS_AKSNR"
	MessageTagASCCMARANSAKSNR2                          MessageTag = "ASCCM_AR_ANS_AKSNR_2"
	MessageTagASCCMPKSK                                 MessageTag = "ASCCM_PKSK"
	MessageTagASCCMPKSKANS                              MessageTag = "ASCCM_PKSK_ANS"
	MessageTagACCMDescWithID                            MessageTag = "ACCM_DESC_WITH_ID"
	MessageTagACCMDescWithIDResult                      MessageTag = "ACCM_DESC_WITH_ID_RESULT"
	MessageTagACCMDescWithName                          MessageTag = "ACCM_DESC_WITH_NAME"
	MessageTagACCMPosSigSig                             MessageTag = "ACCM_POS_SIG_SIG"
	MessageTagACCMPosTSTSig                             MessageTag = "ACCM_POS_TST_SIG"
	MessageTagACCMPosRevocSig                           MessageTag = "ACCM_POS_REVOC_SIG"
	MessageTagACCMPosEVRecord                           MessageTag = "ACCM_POS_EV_RECORD"
	MessageTagACCMPosEAA                                MessageTag = "ACCM_POS_EAA"
	MessageTagACCMPosEAARev                             MessageTag = "ACCM_POS_EAA_REV"
	MessageTagACCMPosCNTRSig                            MessageTag = "ACCM_POS_CNTR_SIG"
	MessageTagACCMPosCNTRSigPL                          MessageTag = "ACCM_POS_CNTR_SIG_PL"
	MessageTagACCMPosConDig                             MessageTag = "ACCM_POS_CON_DIG"
	MessageTagACCMPosEAAKB                              MessageTag = "ACCM_POS_EAA_KB"
	MessageTagACCMPosEAAND                              MessageTag = "ACCM_POS_EAA_ND"
	MessageTagACCMPosEAANSD                             MessageTag = "ACCM_POS_EAA_NSD"
	MessageTagACCMPosEAANSDPL                           MessageTag = "ACCM_POS_EAA_NSD_PL"
	MessageTagACCMPosEAAOSDC                            MessageTag = "ACCM_POS_EAA_OSDC"
	MessageTagACCMPosEAAOSDCPL                          MessageTag = "ACCM_POS_EAA_OSDC_PL"
	MessageTagACCMPosEAAPD                              MessageTag = "ACCM_POS_EAA_PD"
	MessageTagACCMPosEAASD                              MessageTag = "ACCM_POS_EAA_SD"
	MessageTagACCMPosEAASDPL                            MessageTag = "ACCM_POS_EAA_SD_PL"
	MessageTagACCMPosERADO                              MessageTag = "ACCM_POS_ER_ADO"
	MessageTagACCMPosERADOPL                            MessageTag = "ACCM_POS_ER_ADO_PL"
	MessageTagACCMPosEROr                               MessageTag = "ACCM_POS_ER_OR"
	MessageTagACCMPosEROrPL                             MessageTag = "ACCM_POS_ER_OR_PL"
	MessageTagACCMPosERTST                              MessageTag = "ACCM_POS_ER_TST"
	MessageTagACCMPosERTSTSeq                           MessageTag = "ACCM_POS_ER_TST_SEQ"
	MessageTagACCMPosERMSTSig                           MessageTag = "ACCM_POS_ER_MST_SIG"
	MessageTagACCMPosJWS                                MessageTag = "ACCM_POS_JWS"
	MessageTagACCMPosCose                               MessageTag = "ACCM_POS_COSE"
	MessageTagACCMPosKey                                MessageTag = "ACCM_POS_KEY"
	MessageTagACCMPosKeyPL                              MessageTag = "ACCM_POS_KEY_PL"
	MessageTagACCMPosMan                                MessageTag = "ACCM_POS_MAN"
	MessageTagACCMPosManPL                              MessageTag = "ACCM_POS_MAN_PL"
	MessageTagACCMPosManENT                             MessageTag = "ACCM_POS_MAN_ENT"
	MessageTagACCMPosManENTPL                           MessageTag = "ACCM_POS_MAN_ENT_PL"
	MessageTagACCMPosMesDig                             MessageTag = "ACCM_POS_MES_DIG"
	MessageTagACCMPosMessImp                            MessageTag = "ACCM_POS_MESS_IMP"
	MessageTagACCMPosRef                                MessageTag = "ACCM_POS_REF"
	MessageTagACCMPosRefPL                              MessageTag = "ACCM_POS_REF_PL"
	MessageTagACCMPosSigDENT                            MessageTag = "ACCM_POS_SIG_D_ENT"
	MessageTagACCMPosSigDENTPL                          MessageTag = "ACCM_POS_SIG_D_ENT_PL"
	MessageTagACCMPosSigValAndPrt                       MessageTag = "ACCM_POS_SIG_VAL_AND_PRT"
	MessageTagACCMPosSIGNDObj                           MessageTag = "ACCM_POS_SIGND_OBJ"
	MessageTagACCMPosSIGNDPrt                           MessageTag = "ACCM_POS_SIGND_PRT"
	MessageTagACCMPosSigntrPrt                          MessageTag = "ACCM_POS_SIGNTR_PRT"
	MessageTagACCMPosCertChainSig                       MessageTag = "ACCM_POS_CERT_CHAIN_SIG"
	MessageTagACCMPosCertChainTST                       MessageTag = "ACCM_POS_CERT_CHAIN_TST"
	MessageTagACCMPosCertChainRevoc                     MessageTag = "ACCM_POS_CERT_CHAIN_REVOC"
	MessageTagACCMPosCertChainEAARev                    MessageTag = "ACCM_POS_CERT_CHAIN_EAA_REV"
	MessageTagACCMPosCertChain                          MessageTag = "ACCM_POS_CERT_CHAIN"
	MessageTagACCMPosSigCertRef                         MessageTag = "ACCM_POS_SIG_CERT_REF"
	MessageTagBBBSAVDSCACRCC                            MessageTag = "BBB_SAV_DSCACRCC"
	MessageTagBBBSAVDSCACRCCANS                         MessageTag = "BBB_SAV_DSCACRCC_ANS"
	MessageTagBBBSAVACPCCRSCA                           MessageTag = "BBB_SAV_ACPCCRSCA"
	MessageTagBBBSAVACPCCRSCAANS                        MessageTag = "BBB_SAV_ACPCCRSCA_ANS"
	MessageTagBBBSAVISVA                                MessageTag = "BBB_SAV_ISVA"
	MessageTagBBBSAVISVAANS                             MessageTag = "BBB_SAV_ISVA_ANS"
	MessageTagBBBSAVISSV                                MessageTag = "BBB_SAV_ISSV"
	MessageTagBBBSAVISSVANS                             MessageTag = "BBB_SAV_ISSV_ANS"
	MessageTagBBBSAVICERRM                              MessageTag = "BBB_SAV_ICERRM"
	MessageTagBBBSAVICERRMANS                           MessageTag = "BBB_SAV_ICERRM_ANS"
	MessageTagBBBSAVICRM                                MessageTag = "BBB_SAV_ICRM"
	MessageTagBBBSAVICRMANS                             MessageTag = "BBB_SAV_ICRM_ANS"
	MessageTagBBBSAVISQPCTP                             MessageTag = "BBB_SAV_ISQPCTP"
	MessageTagBBBSAVISQPCTPANS                          MessageTag = "BBB_SAV_ISQPCTP_ANS"
	MessageTagBBBSAVISQPCHP                             MessageTag = "BBB_SAV_ISQPCHP"
	MessageTagBBBSAVISQPCHPANS                          MessageTag = "BBB_SAV_ISQPCHP_ANS"
	MessageTagBBBSAVISQPCIP                             MessageTag = "BBB_SAV_ISQPCIP"
	MessageTagBBBSAVISQPCIPANS                          MessageTag = "BBB_SAV_ISQPCIP_ANS"
	MessageTagBBBSAVISQPCTSIP                           MessageTag = "BBB_SAV_ISQPCTSIP"
	MessageTagBBBSAVISQPCTSIPANS                        MessageTag = "BBB_SAV_ISQPCTSIP_ANS"
	MessageTagBBBSAVISQPSTYPP                           MessageTag = "BBB_SAV_ISQPSTYPP"
	MessageTagBBBSAVISQPSTYPPANS                        MessageTag = "BBB_SAV_ISQPSTYPP_ANS"
	MessageTagBBBSAVISQPSLP                             MessageTag = "BBB_SAV_ISQPSLP"
	MessageTagBBBSAVISQPSLPANS                          MessageTag = "BBB_SAV_ISQPSLP_ANS"
	MessageTagBBBSAVISQPSTP                             MessageTag = "BBB_SAV_ISQPSTP"
	MessageTagBBBSAVISQPSTPANS                          MessageTag = "BBB_SAV_ISQPSTP_ANS"
	MessageTagBBBSAVISQPSTWSCVR                         MessageTag = "BBB_SAV_ISQPSTWSCVR"
	MessageTagBBBSAVISQPSTWSCVRANS                      MessageTag = "BBB_SAV_ISQPSTWSCVR_ANS"
	MessageTagBBBSAVISQPXTIP                            MessageTag = "BBB_SAV_ISQPXTIP"
	MessageTagBBBSAVISQPXTIPANS                         MessageTag = "BBB_SAV_ISQPXTIP_ANS"
	MessageTagBBBSAVIUQPCSP                             MessageTag = "BBB_SAV_IUQPCSP"
	MessageTagBBBSAVIUQPCSPANS                          MessageTag = "BBB_SAV_IUQPCSP_ANS"
	MessageTagBBBSAVIUQPSTSP                            MessageTag = "BBB_SAV_IUQPSTSP"
	MessageTagBBBSAVIUQPSTSPANS                         MessageTag = "BBB_SAV_IUQPSTSP_ANS"
	MessageTagBBBSAVIUQPVDTSP                           MessageTag = "BBB_SAV_IUQPVDTSP"
	MessageTagBBBSAVIUQPVDTSPANS                        MessageTag = "BBB_SAV_IUQPVDTSP_ANS"
	MessageTagBBBSAVIUQPVDROTSP                         MessageTag = "BBB_SAV_IUQPVDROTSP"
	MessageTagBBBSAVIUQPVDROTSPANS                      MessageTag = "BBB_SAV_IUQPVDROTSP_ANS"
	MessageTagBBBSAVIUQPATSP                            MessageTag = "BBB_SAV_IUQPATSP"
	MessageTagBBBSAVIUQPATSPANS                         MessageTag = "BBB_SAV_IUQPATSP_ANS"
	MessageTagBBBSAVICTVS                               MessageTag = "BBB_SAV_ICTVS"
	MessageTagBBBSAVICTVSANS                            MessageTag = "BBB_SAV_ICTVS_ANS"
	MessageTagBBBSAVIDTSP                               MessageTag = "BBB_SAV_IDTSP"
	MessageTagBBBSAVIDTSPANS                            MessageTag = "BBB_SAV_IDTSP_ANS"
	MessageTagBBBSAVITVS                                MessageTag = "BBB_SAV_ITVS"
	MessageTagBBBSAVITVSANS                             MessageTag = "BBB_SAV_ITVS_ANS"
	MessageTagBBBSAVIVTTSTP                             MessageTag = "BBB_SAV_IVTTSTP"
	MessageTagBBBSAVIVTTSTPANS                          MessageTag = "BBB_SAV_IVTTSTP_ANS"
	MessageTagBBBSAVIVLTATSTP                           MessageTag = "BBB_SAV_IVLTATSTP"
	MessageTagBBBSAVIVLTATSTPANS                        MessageTag = "BBB_SAV_IVLTATSTP_ANS"
	MessageTagBBBSAVISQPMDOSPP                          MessageTag = "BBB_SAV_ISQPMDOSPP"
	MessageTagBBBSAVISQPMDOSPPANS                       MessageTag = "BBB_SAV_ISQPMDOSPP_ANS"
	MessageTagBBBSAVDMICTSTMCMI                         MessageTag = "BBB_SAV_DMICTSTMCMI"
	MessageTagBBBSAVDMICTSTMCMIANS                      MessageTag = "BBB_SAV_DMICTSTMCMI_ANS"
	MessageTagBBBTavITSAP                               MessageTag = "BBB_TAV_ITSAP"
	MessageTagBBBTavITSAPANS                            MessageTag = "BBB_TAV_ITSAP_ANS"
	MessageTagBBBTavDTSAVM                              MessageTag = "BBB_TAV_DTSAVM"
	MessageTagBBBTavDTSAVMANS                           MessageTag = "BBB_TAV_DTSAVM_ANS"
	MessageTagBBBTavDTSAOM                              MessageTag = "BBB_TAV_DTSAOM"
	MessageTagBBBTavDTSAOMANS                           MessageTag = "BBB_TAV_DTSAOM_ANS"
	MessageTagBBBVCIISPK                                MessageTag = "BBB_VCI_ISPK"
	MessageTagBBBVCIISPKANS                             MessageTag = "BBB_VCI_ISPK_ANS"
	MessageTagBBBVCIISPA                                MessageTag = "BBB_VCI_ISPA"
	MessageTagBBBVCIISPAANS                             MessageTag = "BBB_VCI_ISPA_ANS"
	MessageTagBBBVCIISPSUPP                             MessageTag = "BBB_VCI_ISPSUPP"
	MessageTagBBBVCIISPSUPPANS                          MessageTag = "BBB_VCI_ISPSUPP_ANS"
	MessageTagBBBVCIISPM                                MessageTag = "BBB_VCI_ISPM"
	MessageTagBBBVCIISPMANS                             MessageTag = "BBB_VCI_ISPM_ANS"
	MessageTagBBBVCIIZHSP                               MessageTag = "BBB_VCI_IZHSP"
	MessageTagBBBVCIIZHSPANS                            MessageTag = "BBB_VCI_IZHSP_ANS"
	MessageTagBBBXCVSub                                 MessageTag = "BBB_XCV_SUB"
	MessageTagBBBXCVSubANS                              MessageTag = "BBB_XCV_SUB_ANS"
	MessageTagBBBXCVSubANS2                             MessageTag = "BBB_XCV_SUB_ANS_2"
	MessageTagBBBXCVRFC                                 MessageTag = "BBB_XCV_RFC"
	MessageTagBBBXCVRFCANS                              MessageTag = "BBB_XCV_RFC_ANS"
	MessageTagBBBXCVRAC                                 MessageTag = "BBB_XCV_RAC"
	MessageTagBBBXCVRACANS                              MessageTag = "BBB_XCV_RAC_ANS"
	MessageTagBBBXCVCCCBB                               MessageTag = "BBB_XCV_CCCBB"
	MessageTagBBBXCVCCCBBANS                            MessageTag = "BBB_XCV_CCCBB_ANS"
	MessageTagBBBXCVCCCBBSigANS                         MessageTag = "BBB_XCV_CCCBB_SIG_ANS"
	MessageTagBBBXCVCCCBBTSPANS                         MessageTag = "BBB_XCV_CCCBB_TSP_ANS"
	MessageTagBBBXCVCCCBBRevANS                         MessageTag = "BBB_XCV_CCCBB_REV_ANS"
	MessageTagBBBXCVCMDCIPI                             MessageTag = "BBB_XCV_CMDCIPI"
	MessageTagBBBXCVCMDCIPIANS                          MessageTag = "BBB_XCV_CMDCIPI_ANS"
	MessageTagBBBXCVCMDCIQC                             MessageTag = "BBB_XCV_CMDCIQC"
	MessageTagBBBXCVCMDCIQCANS                          MessageTag = "BBB_XCV_CMDCIQC_ANS"
	MessageTagBBBXCVCMDCIQSCD                           MessageTag = "BBB_XCV_CMDCIQSCD"
	MessageTagBBBXCVCMDCIQSCDANS                        MessageTag = "BBB_XCV_CMDCIQSCD_ANS"
	MessageTagBBBXCVCMDCIITLP                           MessageTag = "BBB_XCV_CMDCIITLP"
	MessageTagBBBXCVCMDCIITLPANS                        MessageTag = "BBB_XCV_CMDCIITLP_ANS"
	MessageTagBBBXCVCMDCIITNP                           MessageTag = "BBB_XCV_CMDCIITNP"
	MessageTagBBBXCVCMDCIITNPANS                        MessageTag = "BBB_XCV_CMDCIITNP_ANS"
	MessageTagBBBXCVCMDCICQCC                           MessageTag = "BBB_XCV_CMDCICQCC"
	MessageTagBBBXCVCMDCICQCCANS                        MessageTag = "BBB_XCV_CMDCICQCC_ANS"
	MessageTagBBBXCVCMDCICQCLVA                         MessageTag = "BBB_XCV_CMDCICQCLVA"
	MessageTagBBBXCVCMDCICQCLVAANS                      MessageTag = "BBB_XCV_CMDCICQCLVA_ANS"
	MessageTagBBBXCVCMDCICQCLVHAC                       MessageTag = "BBB_XCV_CMDCICQCLVHAC"
	MessageTagBBBXCVCMDCICQCLVHACANS                    MessageTag = "BBB_XCV_CMDCICQCLVHAC_ANS"
	MessageTagBBBXCVCMDCICQCERPA                        MessageTag = "BBB_XCV_CMDCICQCERPA"
	MessageTagBBBXCVCMDCICQCERPAANS                     MessageTag = "BBB_XCV_CMDCICQCERPA_ANS"
	MessageTagBBBXCVCMDCICSQCSSCD                       MessageTag = "BBB_XCV_CMDCICSQCSSCD"
	MessageTagBBBXCVCMDCICSQCSSCDANS                    MessageTag = "BBB_XCV_CMDCICSQCSSCD_ANS"
	MessageTagBBBXCVCMDCICQCPDSLA                       MessageTag = "BBB_XCV_CMDCICQCPDSLA"
	MessageTagBBBXCVCMDCICQCPDSLAANS                    MessageTag = "BBB_XCV_CMDCICQCPDSLA_ANS"
	MessageTagBBBXCVCMDCICQCTA                          MessageTag = "BBB_XCV_CMDCICQCTA"
	MessageTagBBBXCVCMDCICQCTAANS                       MessageTag = "BBB_XCV_CMDCICQCTA_ANS"
	MessageTagBBBXCVCMDCDCQCCLCEC                       MessageTag = "BBB_XCV_CMDCDCQCCLCEC"
	MessageTagBBBXCVCMDCDCQCCLCECANS                    MessageTag = "BBB_XCV_CMDCDCQCCLCEC_ANS"
	MessageTagBBBXCVCMDCDCQCCLCECANSEU                  MessageTag = "BBB_XCV_CMDCDCQCCLCEC_ANS_EU"
	MessageTagBBBXCVCMDCSCSIA                           MessageTag = "BBB_XCV_CMDCSCSIA"
	MessageTagBBBXCVCMDCSCSIAANS                        MessageTag = "BBB_XCV_CMDCSCSIA_ANS"
	MessageTagBBBXCVCMDCICQCRA                          MessageTag = "BBB_XCV_CMDCICQCRA"
	MessageTagBBBXCVCMDCICQCRAANS                       MessageTag = "BBB_XCV_CMDCICQCRA_ANS"
	MessageTagBBBXCVCMDCICQCNA                          MessageTag = "BBB_XCV_CMDCICQCNA"
	MessageTagBBBXCVCMDCICQCNAANS                       MessageTag = "BBB_XCV_CMDCICQCNA_ANS"
	MessageTagBBBXCVCMDCICQCIA                          MessageTag = "BBB_XCV_CMDCICQCIA"
	MessageTagBBBXCVCMDCICQCIAANS                       MessageTag = "BBB_XCV_CMDCICQCIA_ANS"
	MessageTagBBBXCVCMDCDCQCQSCDLSA                     MessageTag = "BBB_XCV_CMDCDCQCQSCDLSA"
	MessageTagBBBXCVCMDCDCQCQSCDLSAANS                  MessageTag = "BBB_XCV_CMDCDCQCQSCDLSA_ANS"
	MessageTagBBBXCVCMDCDCQCIMSA                        MessageTag = "BBB_XCV_CMDCDCQCIMSA"
	MessageTagBBBXCVCMDCDCQCIMSAANS                     MessageTag = "BBB_XCV_CMDCDCQCIMSA_ANS"
	MessageTagBBBXCVCMDCPSBCLA                          MessageTag = "BBB_XCV_CMDCPSBCLA"
	MessageTagBBBXCVCMDCPSBCLAANS                       MessageTag = "BBB_XCV_CMDCPSBCLA_ANS"
	MessageTagBBBXCVCMDCPSBASIA                         MessageTag = "BBB_XCV_CMDCPSBASIA"
	MessageTagBBBXCVCMDCPSBASIAANS                      MessageTag = "BBB_XCV_CMDCPSBASIA_ANS"
	MessageTagBBBXCVCMDCPSBLIA                          MessageTag = "BBB_XCV_CMDCPSBLIA"
	MessageTagBBBXCVCMDCPSBLIAANS                       MessageTag = "BBB_XCV_CMDCPSBLIA_ANS"
	MessageTagBBBXCVDCCUCE                              MessageTag = "BBB_XCV_DCCUCE"
	MessageTagBBBXCVDCCUCEANS                           MessageTag = "BBB_XCV_DCCUCE_ANS"
	MessageTagBBBXCVDCCFCE                              MessageTag = "BBB_XCV_DCCFCE"
	MessageTagBBBXCVDCCFCEANS                           MessageTag = "BBB_XCV_DCCFCE_ANS"
	MessageTagBBBXCVDCSBSINC                            MessageTag = "BBB_XCV_DCSBSINC"
	MessageTagBBBXCVDCSBSINCANS                         MessageTag = "BBB_XCV_DCSBSINC_ANS"
	MessageTagBBBXCVICAC                                MessageTag = "BBB_XCV_ICAC"
	MessageTagBBBXCVICACANS                             MessageTag = "BBB_XCV_ICAC_ANS"
	MessageTagBBBXCVICPDV                               MessageTag = "BBB_XCV_ICPDV"
	MessageTagBBBXCVICPDVANS                            MessageTag = "BBB_XCV_ICPDV_ANS"
	MessageTagBBBXCVICPTV                               MessageTag = "BBB_XCV_ICPTV"
	MessageTagBBBXCVICPTVANS                            MessageTag = "BBB_XCV_ICPTV_ANS"
	MessageTagBBBXCVIAKIP                               MessageTag = "BBB_XCV_IAKIP"
	MessageTagBBBXCVIAKIPANS                            MessageTag = "BBB_XCV_IAKIP_ANS"
	MessageTagBBBXCVISKIP                               MessageTag = "BBB_XCV_ISKIP"
	MessageTagBBBXCVISKIPANS                            MessageTag = "BBB_XCV_ISKIP_ANS"
	MessageTagBBBXCVICNRAEV                             MessageTag = "BBB_XCV_ICNRAEV"
	MessageTagBBBXCVICNRAEVANS                          MessageTag = "BBB_XCV_ICNRAEV_ANS"
	MessageTagBBBXCVIVTBCTSD                            MessageTag = "BBB_XCV_IVTBCTSD"
	MessageTagBBBXCVIVTBCTSDANS                         MessageTag = "BBB_XCV_IVTBCTSD_ANS"
	MessageTagBBBXCVICTIVRSC                            MessageTag = "BBB_XCV_ICTIVRSC"
	MessageTagBBBXCVICTIVRSCANS                         MessageTag = "BBB_XCV_ICTIVRSC_ANS"
	MessageTagBBBXCVICTIVRCIRI                          MessageTag = "BBB_XCV_ICTIVRCIRI"
	MessageTagBBBXCVICTIVRCIRIANS                       MessageTag = "BBB_XCV_ICTIVRCIRI_ANS"
	MessageTagBBBXCVIRDCSFC                             MessageTag = "BBB_XCV_IRDCSFC"
	MessageTagBBBXCVIRDCSFCANS                          MessageTag = "BBB_XCV_IRDCSFC_ANS"
	MessageTagBBBXCVIRDPFC                              MessageTag = "BBB_XCV_IRDPFC"
	MessageTagBBBXCVIRDPFCANS                           MessageTag = "BBB_XCV_IRDPFC_ANS"
	MessageTagBBBXCVIRDPFRC                             MessageTag = "BBB_XCV_IRDPFRC"
	MessageTagBBBXCVIRDPFRCANS                          MessageTag = "BBB_XCV_IRDPFRC_ANS"
	MessageTagBBBXCVIARDPFC                             MessageTag = "BBB_XCV_IARDPFC"
	MessageTagBBBXCVIARDPFCANS                          MessageTag = "BBB_XCV_IARDPFC_ANS"
	MessageTagBBBVTSIRDPFC                              MessageTag = "BBB_VTS_IRDPFC"
	MessageTagBBBVTSIRDPFCANS                           MessageTag = "BBB_VTS_IRDPFC_ANS"
	MessageTagBBBXCVISCOH                               MessageTag = "BBB_XCV_ISCOH"
	MessageTagBBBXCVISCOHANS                            MessageTag = "BBB_XCV_ISCOH_ANS"
	MessageTagBBBXCVISCUKN                              MessageTag = "BBB_XCV_ISCUKN"
	MessageTagBBBXCVISCUKNANS                           MessageTag = "BBB_XCV_ISCUKN_ANS"
	MessageTagBBBXCVISCR                                MessageTag = "BBB_XCV_ISCR"
	MessageTagBBBXCVISCRANS                             MessageTag = "BBB_XCV_ISCR_ANS"
	MessageTagBBBXCVISCGKU                              MessageTag = "BBB_XCV_ISCGKU"
	MessageTagBBBXCVISCGKUANS                           MessageTag = "BBB_XCV_ISCGKU_ANS"
	MessageTagBBBXCVISCGKUANSCert                       MessageTag = "BBB_XCV_ISCGKU_ANS_CERT"
	MessageTagBBBXCVISCGEKU                             MessageTag = "BBB_XCV_ISCGEKU"
	MessageTagBBBXCVISCGEKUANS                          MessageTag = "BBB_XCV_ISCGEKU_ANS"
	MessageTagBBBXCVISCGEKUANSCert                      MessageTag = "BBB_XCV_ISCGEKU_ANS_CERT"
	MessageTagBBBXCVICSI                                MessageTag = "BBB_XCV_ICSI"
	MessageTagBBBXCVICSIANS                             MessageTag = "BBB_XCV_ICSI_ANS"
	MessageTagBBBXCVIOTAA                               MessageTag = "BBB_XCV_IOTAA"
	MessageTagBBBXCVIOTAAANS                            MessageTag = "BBB_XCV_IOTAA_ANS"
	MessageTagBBBXCVHPCCVVT                             MessageTag = "BBB_XCV_HPCCVVT"
	MessageTagBBBXCVHPCCVVTANS                          MessageTag = "BBB_XCV_HPCCVVT_ANS"
	MessageTagBBBXCVPseudoUse                           MessageTag = "BBB_XCV_PSEUDO_USE"
	MessageTagBBBXCVPseudoUseANS                        MessageTag = "BBB_XCV_PSEUDO_USE_ANS"
	MessageTagBBBXCVAIAPres                             MessageTag = "BBB_XCV_AIA_PRES"
	MessageTagBBBXCVAIAPresANS                          MessageTag = "BBB_XCV_AIA_PRES_ANS"
	MessageTagBBBXCVRevocPres                           MessageTag = "BBB_XCV_REVOC_PRES"
	MessageTagBBBXCVRevocPresANS                        MessageTag = "BBB_XCV_REVOC_PRES_ANS"
	MessageTagBBBXCVRevocThisUpdatePresent              MessageTag = "BBB_XCV_REVOC_THIS_UPDATE_PRESENT"
	MessageTagBBBXCVRevocThisUpdatePresentANS           MessageTag = "BBB_XCV_REVOC_THIS_UPDATE_PRESENT_ANS"
	MessageTagBBBXCVRevocIssuerKnown                    MessageTag = "BBB_XCV_REVOC_ISSUER_KNOWN"
	MessageTagBBBXCVRevocIssuerKnownANS                 MessageTag = "BBB_XCV_REVOC_ISSUER_KNOWN_ANS"
	MessageTagBBBXCVRevocIssuerValidAtProd              MessageTag = "BBB_XCV_REVOC_ISSUER_VALID_AT_PROD"
	MessageTagBBBXCVRevocIssuerValidAtProdANS           MessageTag = "BBB_XCV_REVOC_ISSUER_VALID_AT_PROD_ANS"
	MessageTagBBBXCVRevocAfterCertNotBefore             MessageTag = "BBB_XCV_REVOC_AFTER_CERT_NOT_BEFORE"
	MessageTagBBBXCVRevocAfterCertNotBeforeANS          MessageTag = "BBB_XCV_REVOC_AFTER_CERT_NOT_BEFORE_ANS"
	MessageTagBBBXCVRevocHasCertInfo                    MessageTag = "BBB_XCV_REVOC_HAS_CERT_INFO"
	MessageTagBBBXCVRevocHasCertInfoANS                 MessageTag = "BBB_XCV_REVOC_HAS_CERT_INFO_ANS"
	MessageTagBBBXCVRevocRespIDMatch                    MessageTag = "BBB_XCV_REVOC_RESPID_MATCH"
	MessageTagBBBXCVRevocRespIDMatchANS                 MessageTag = "BBB_XCV_REVOC_RESPID_MATCH_ANS"
	MessageTagBBBXCVRevocCertHashPresent                MessageTag = "BBB_XCV_REVOC_CERT_HASH_PRESENT"
	MessageTagBBBXCVRevocCertHashPresentANS             MessageTag = "BBB_XCV_REVOC_CERT_HASH_PRESENT_ANS"
	MessageTagBBBXCVRevocCertHashMatch                  MessageTag = "BBB_XCV_REVOC_CERT_HASH_MATCH"
	MessageTagBBBXCVRevocCertHashMatchANS               MessageTag = "BBB_XCV_REVOC_CERT_HASH_MATCH_ANS"
	MessageTagBBBXCVRevocSelfIssuedOCSP                 MessageTag = "BBB_XCV_REVOC_SELF_ISSUED_OCSP"
	MessageTagBBBXCVRevocSelfIssuedOCSPANS              MessageTag = "BBB_XCV_REVOC_SELF_ISSUED_OCSP_ANS"
	MessageTagBBBXCVDCIDNMSDNIC                         MessageTag = "BBB_XCV_DCIDNMSDNIC"
	MessageTagBBBXCVDCIDNMSDNICANS                      MessageTag = "BBB_XCV_DCIDNMSDNIC_ANS"
	MessageTagBBBXCVISCGCOUN                            MessageTag = "BBB_XCV_ISCGCOUN"
	MessageTagBBBXCVISCGCOUNANS                         MessageTag = "BBB_XCV_ISCGCOUN_ANS"
	MessageTagBBBXCVISCGLOC                             MessageTag = "BBB_XCV_ISCGLOC"
	MessageTagBBBXCVISCGLOCANS                          MessageTag = "BBB_XCV_ISCGLOC_ANS"
	MessageTagBBBXCVISCGST                              MessageTag = "BBB_XCV_ISCGST"
	MessageTagBBBXCVISCGSTANS                           MessageTag = "BBB_XCV_ISCGST_ANS"
	MessageTagBBBXCVISCGORGAN                           MessageTag = "BBB_XCV_ISCGORGAN"
	MessageTagBBBXCVISCGORGANANS                        MessageTag = "BBB_XCV_ISCGORGAN_ANS"
	MessageTagBBBXCVISCGORGAU                           MessageTag = "BBB_XCV_ISCGORGAU"
	MessageTagBBBXCVISCGORGAUANS                        MessageTag = "BBB_XCV_ISCGORGAU_ANS"
	MessageTagBBBXCVISCGORGAI                           MessageTag = "BBB_XCV_ISCGORGAI"
	MessageTagBBBXCVISCGORGAIANS                        MessageTag = "BBB_XCV_ISCGORGAI_ANS"
	MessageTagBBBXCVISCGSURN                            MessageTag = "BBB_XCV_ISCGSURN"
	MessageTagBBBXCVISCGSURNANS                         MessageTag = "BBB_XCV_ISCGSURN_ANS"
	MessageTagBBBXCVISCGGIVEN                           MessageTag = "BBB_XCV_ISCGGIVEN"
	MessageTagBBBXCVISCGGIVENANS                        MessageTag = "BBB_XCV_ISCGGIVEN_ANS"
	MessageTagBBBXCVISCGPSEUDO                          MessageTag = "BBB_XCV_ISCGPSEUDO"
	MessageTagBBBXCVISCGPSEUDOANS                       MessageTag = "BBB_XCV_ISCGPSEUDO_ANS"
	MessageTagBBBXCVISCGCOMMONN                         MessageTag = "BBB_XCV_ISCGCOMMONN"
	MessageTagBBBXCVISCGCOMMONNANS                      MessageTag = "BBB_XCV_ISCGCOMMONN_ANS"
	MessageTagBBBXCVISCGTITLE                           MessageTag = "BBB_XCV_ISCGTITLE"
	MessageTagBBBXCVISCGTITLEANS                        MessageTag = "BBB_XCV_ISCGTITLE_ANS"
	MessageTagBBBXCVISCGEMAIL                           MessageTag = "BBB_XCV_ISCGEMAIL"
	MessageTagBBBXCVISCGEMAILANS                        MessageTag = "BBB_XCV_ISCGEMAIL_ANS"
	MessageTagBBBXCVISSSC                               MessageTag = "BBB_XCV_ISSSC"
	MessageTagBBBXCVISSSCANS                            MessageTag = "BBB_XCV_ISSSC_ANS"
	MessageTagBBBXCVISNSSC                              MessageTag = "BBB_XCV_ISNSSC"
	MessageTagBBBXCVISNSSCANS                           MessageTag = "BBB_XCV_ISNSSC_ANS"
	MessageTagBBBXCVIRDC                                MessageTag = "BBB_XCV_IRDC"
	MessageTagBBBXCVIRDCANS                             MessageTag = "BBB_XCV_IRDC_ANS"
	MessageTagXCVTSLESP                                 MessageTag = "XCV_TSL_ESP"
	MessageTagXCVTSLESPANS                              MessageTag = "XCV_TSL_ESP_ANS"
	MessageTagXCVTSLESPSigANS                           MessageTag = "XCV_TSL_ESP_SIG_ANS"
	MessageTagXCVTSLESPTSPANS                           MessageTag = "XCV_TSL_ESP_TSP_ANS"
	MessageTagXCVTSLESPRevANS                           MessageTag = "XCV_TSL_ESP_REV_ANS"
	MessageTagXCVTSLETIP                                MessageTag = "XCV_TSL_ETIP"
	MessageTagXCVTSLETIPANS                             MessageTag = "XCV_TSL_ETIP_ANS"
	MessageTagXCVTSLETIPSigANS                          MessageTag = "XCV_TSL_ETIP_SIG_ANS"
	MessageTagXCVTSLETIPTSPANS                          MessageTag = "XCV_TSL_ETIP_TSP_ANS"
	MessageTagXCVTSLETIPRevANS                          MessageTag = "XCV_TSL_ETIP_REV_ANS"
	MessageTagPCVIVTSC                                  MessageTag = "PCV_IVTSC"
	MessageTagPCVIVTSCANS                               MessageTag = "PCV_IVTSC_ANS"
	MessageTagPCVICCSVTSF                               MessageTag = "PCV_ICCSVTSF"
	MessageTagPCVICCSVTSFANS                            MessageTag = "PCV_ICCSVTSF_ANS"
	MessageTagPSVIPCVA                                  MessageTag = "PSV_IPCVA"
	MessageTagPSVIPCVAANS                               MessageTag = "PSV_IPCVA_ANS"
	MessageTagPSVIPCVC                                  MessageTag = "PSV_IPCVC"
	MessageTagPSVIPCVCANS                               MessageTag = "PSV_IPCVC_ANS"
	MessageTagPSVIPSVC                                  MessageTag = "PSV_IPSVC"
	MessageTagPSVIPSVCANS                               MessageTag = "PSV_IPSVC_ANS"
	MessageTagPSVIPTVC                                  MessageTag = "PSV_IPTVC"
	MessageTagPSVIPTVCANS                               MessageTag = "PSV_IPTVC_ANS"
	MessageTagPSVITPOCOBCT                              MessageTag = "PSV_ITPOCOBCT"
	MessageTagPSVITPOSVAOBCT                            MessageTag = "PSV_ITPOSVAOBCT"
	MessageTagPSVITPOSVAOBCTANS                         MessageTag = "PSV_ITPOSVAOBCT_ANS"
	MessageTagPSVITPORDAOBCT                            MessageTag = "PSV_ITPORDAOBCT"
	MessageTagPSVITPOOBCTANS                            MessageTag = "PSV_ITPOOBCT_ANS"
	MessageTagPSVITPRISCNARTCAC                         MessageTag = "PSV_ITPRISCNARTCAC"
	MessageTagPSVITPRISCNARTCACANS                      MessageTag = "PSV_ITPRISCNARTCAC_ANS"
	MessageTagPSVICRDIT                                 MessageTag = "PSV_ICRDIT"
	MessageTagPSVICRDITANS                              MessageTag = "PSV_ICRDIT_ANS"
	MessageTagPSVIPCRIAIDBEDC                           MessageTag = "PSV_IPCRIAIDBEDC"
	MessageTagPSVIPCRIAIDBEDCANS                        MessageTag = "PSV_IPCRIAIDBEDC_ANS"
	MessageTagPSVICTD                                   MessageTag = "PSV_ICTD"
	MessageTagPSVICTDANS                                MessageTag = "PSV_ICTD_ANS"
	MessageTagPSVISDDTA                                 MessageTag = "PSV_ISDDTA"
	MessageTagPSVISDDTAANS                              MessageTag = "PSV_ISDDTA_ANS"
	MessageTagPSVHRDBIBCT                               MessageTag = "PSV_HRDBIBCT"
	MessageTagPSVHRDBIBCTANS                            MessageTag = "PSV_HRDBIBCT_ANS"
	MessageTagPSVDIURDSCHPVR                            MessageTag = "PSV_DIURDSCHPVR"
	MessageTagPSVDIURDSCHPVRANS                         MessageTag = "PSV_DIURDSCHPVR_ANS"
	MessageTagTSVASTPTCT                                MessageTag = "TSV_ASTPTCT"
	MessageTagTSVASTPTCTANS                             MessageTag = "TSV_ASTPTCT_ANS"
	MessageTagTSVIBSTAIDOSC                             MessageTag = "TSV_IBSTAIDOSC"
	MessageTagTSVIBSTAIDOSCANS                          MessageTag = "TSV_IBSTAIDOSC_ANS"
	MessageTagTSVIBSTBCEC                               MessageTag = "TSV_IBSTBCEC"
	MessageTagTSVIBSTBCECANS                            MessageTag = "TSV_IBSTBCEC_ANS"
	MessageTagTSVISCNVABST                              MessageTag = "TSV_ISCNVABST"
	MessageTagTSVISCNVABSTANS                           MessageTag = "TSV_ISCNVABST_ANS"
	MessageTagADESTIRTPTBST                             MessageTag = "ADEST_IRTPTBST"
	MessageTagADESTIRTPTBSTANS                          MessageTag = "ADEST_IRTPTBST_ANS"
	MessageTagADESTISTPTBST                             MessageTag = "ADEST_ISTPTBST"
	MessageTagADESTISTPTBSTANS                          MessageTag = "ADEST_ISTPTBST_ANS"
	MessageTagADESTVFDTAOCSTANS                         MessageTag = "ADEST_VFDTAOCST_ANS"
	MessageTagADESTISTPTDABST                           MessageTag = "ADEST_ISTPTDABST"
	MessageTagADESTISTPTDABSTANS                        MessageTag = "ADEST_ISTPTDABST_ANS"
	MessageTagADESTIBSVPSC                              MessageTag = "ADEST_IBSVPSC"
	MessageTagADESTIBSVPSCANS                           MessageTag = "ADEST_IBSVPSC_ANS"
	MessageTagADESTIBSVPTC                              MessageTag = "ADEST_IBSVPTC"
	MessageTagADESTIBSVPTCANS                           MessageTag = "ADEST_IBSVPTC_ANS"
	MessageTagADESTIBSVPTADC                            MessageTag = "ADEST_IBSVPTADC"
	MessageTagADESTIBSVPTADCANS                         MessageTag = "ADEST_IBSVPTADC_ANS"
	MessageTagADESTIRERVPC                              MessageTag = "ADEST_IRERVPC"
	MessageTagADESTIRERVPCANS                           MessageTag = "ADEST_IRERVPC_ANS"
	MessageTagEAACertLoTEReached                        MessageTag = "EAA_CERT_LOTE_REACHED"
	MessageTagEAACertLoTEReachedANS                     MessageTag = "EAA_CERT_LOTE_REACHED_ANS"
	MessageTagEAACertTrustAnchorListReached             MessageTag = "EAA_CERT_TRUST_ANCHOR_LIST_REACHED"
	MessageTagEAACertTrustAnchorListReachedANS          MessageTag = "EAA_CERT_TRUST_ANCHOR_LIST_REACHED_ANS"
	MessageTagEAADPEAAP                                 MessageTag = "EAA_DPEAAP"
	MessageTagEAADPEAAPANS                              MessageTag = "EAA_DPEAAP_ANS"
	MessageTagEAADLEEAAP                                MessageTag = "EAA_DLEEAAP"
	MessageTagEAADLEEAAPANS                             MessageTag = "EAA_DLEEAAP_ANS"
	MessageTagEAAKBRC                                   MessageTag = "EAA_KBRC"
	MessageTagEAAKBRCANS                                MessageTag = "EAA_KBRC_ANS"
	MessageTagEAAKBSP                                   MessageTag = "EAA_KBSP"
	MessageTagEAAKBSPANS                                MessageTag = "EAA_KBSP_ANS"
	MessageTagEAAClaims                                 MessageTag = "EAA_CLAIMS"
	MessageTagEAAClaimsANS                              MessageTag = "EAA_CLAIMS_ANS"
	MessageTagEAAClaimsInfo                             MessageTag = "EAA_CLAIMS_INFO"
	MessageTagEAASupportedClaims                        MessageTag = "EAA_SUPPORTED_CLAIMS"
	MessageTagEAASupportedClaimsANS                     MessageTag = "EAA_SUPPORTED_CLAIMS_ANS"
	MessageTagEAAUnsupportedClaims                      MessageTag = "EAA_UNSUPPORTED_CLAIMS"
	MessageTagEAAAcceptableType                         MessageTag = "EAA_ACCEPTABLE_TYPE"
	MessageTagEAAAcceptableTypeANS                      MessageTag = "EAA_ACCEPTABLE_TYPE_ANS"
	MessageTagEAAIdentifierPresent                      MessageTag = "EAA_IDENTIFIER_PRESENT"
	MessageTagEAAIdentifierPresentANS                   MessageTag = "EAA_IDENTIFIER_PRESENT_ANS"
	MessageTagEAAIssuanceDatePresent                    MessageTag = "EAA_ISSUANCE_DATE_PRESENT"
	MessageTagEAAIssuanceDatePresentANS                 MessageTag = "EAA_ISSUANCE_DATE_PRESENT_ANS"
	MessageTagEAANBFPresent                             MessageTag = "EAA_NBF_PRESENT"
	MessageTagEAANBFPresentANS                          MessageTag = "EAA_NBF_PRESENT_ANS"
	MessageTagEAAExpPresent                             MessageTag = "EAA_EXP_PRESENT"
	MessageTagEAAExpPresentANS                          MessageTag = "EAA_EXP_PRESENT_ANS"
	MessageTagEAAAIDPresent                             MessageTag = "EAA_AID_PRESENT"
	MessageTagEAAAIDPresentANS                          MessageTag = "EAA_AID_PRESENT_ANS"
	MessageTagEAAAEDPresent                             MessageTag = "EAA_AED_PRESENT"
	MessageTagEAAAEDPresentANS                          MessageTag = "EAA_AED_PRESENT_ANS"
	MessageTagEAASigPresent                             MessageTag = "EAA_SIG_PRESENT"
	MessageTagEAASigPresentANS                          MessageTag = "EAA_SIG_PRESENT_ANS"
	MessageTagEAASigQual                                MessageTag = "EAA_SIG_QUAL"
	MessageTagEAASigQualANS                             MessageTag = "EAA_SIG_QUAL_ANS"
	MessageTagEAACATEAA                                 MessageTag = "EAA_CAT_EAA"
	MessageTagEAACATEAAANS1                             MessageTag = "EAA_CAT_EAA_ANS_1"
	MessageTagEAACATEAAANS2                             MessageTag = "EAA_CAT_EAA_ANS_2"
	MessageTagEAACATPubEAA                              MessageTag = "EAA_CAT_PUBEAA"
	MessageTagEAACATPubEAAANS                           MessageTag = "EAA_CAT_PUBEAA_ANS"
	MessageTagEAACATQEAA                                MessageTag = "EAA_CAT_QEAA"
	MessageTagEAACATQEAAANS                             MessageTag = "EAA_CAT_QEAA_ANS"
	MessageTagEAAQCPSB                                  MessageTag = "EAA_QC_PSB"
	MessageTagEAAQCPSBANS                               MessageTag = "EAA_QC_PSB_ANS"
	MessageTagEAAQualConclusive                         MessageTag = "EAA_QUAL_CONCLUSIVE"
	MessageTagEAAQualConclusiveANS                      MessageTag = "EAA_QUAL_CONCLUSIVE_ANS"
	MessageTagEAAETSI194721                             MessageTag = "EAA_ETSI194721"
	MessageTagEAAETSI194721ANS                          MessageTag = "EAA_ETSI194721_ANS"
	MessageTagEAAVTITVR                                 MessageTag = "EAA_VT_ITVR"
	MessageTagEAAVTITVRANS                              MessageTag = "EAA_VT_ITVR_ANS"
	MessageTagEAAVTITVRValidity                         MessageTag = "EAA_VT_ITVR_VALIDITY"
	MessageTagEAANowBeforeNBF                           MessageTag = "EAA_NOW_BEFORE_NBF"
	MessageTagEAANowAfterExp                            MessageTag = "EAA_NOW_AFTER_EXP"
	MessageTagEAAVTIAVR                                 MessageTag = "EAA_VT_IAVR"
	MessageTagEAAVTIAVRANS                              MessageTag = "EAA_VT_IAVR_ANS"
	MessageTagEAAVTIAVRValidity                         MessageTag = "EAA_VT_IAVR_VALIDITY"
	MessageTagEAANowBeforeADI                           MessageTag = "EAA_NOW_BEFORE_ADI"
	MessageTagEAANowAfterADE                            MessageTag = "EAA_NOW_AFTER_ADE"
	MessageTagEAAADSDJWTConformance                     MessageTag = "EAA_AD_SDJWT_CONFORMANCE"
	MessageTagEAAShortLivedStatusPresent                MessageTag = "EAA_SHORT_LIVED_STATUS_PRESENT"
	MessageTagEAAMandatoryStatusAbsent                  MessageTag = "EAA_MANDATORY_STATUS_ABSENT"
	MessageTagEAARevSDJWTConformance                    MessageTag = "EAA_REV_SDJWT_CONFORMANCE"
	MessageTagEAAMDocIssuingAuthority                   MessageTag = "EAA_MDOC_ISSUING_AUTHORITY"
	MessageTagEAASDJWTIssuingAuthority                  MessageTag = "EAA_SDJWT_ISSUING_AUTHORITY"
	MessageTagEAAMDocDocumentNumberAbsent               MessageTag = "EAA_MDOC_DOCUMENT_NUMBER_ABSENT"
	MessageTagEAASub                                    MessageTag = "EAA_SUB"
	MessageTagEAASubANS                                 MessageTag = "EAA_SUB_ANS"
	MessageTagEAASubPSE                                 MessageTag = "EAA_SUB_PSE"
	MessageTagEAASubPSEANS                              MessageTag = "EAA_SUB_PSE_ANS"
	MessageTagEAACAT                                    MessageTag = "EAA_CAT"
	MessageTagEAACATANS                                 MessageTag = "EAA_CAT_ANS"
	MessageTagEAAISSCOUN                                MessageTag = "EAA_ISS_COUN"
	MessageTagEAAISSCOUNANS                             MessageTag = "EAA_ISS_COUN_ANS"
	MessageTagEAAISSAuth                                MessageTag = "EAA_ISS_AUTH"
	MessageTagEAAISSAuthANS                             MessageTag = "EAA_ISS_AUTH_ANS"
	MessageTagEAAISSRegID                               MessageTag = "EAA_ISS_REG_ID"
	MessageTagEAAISSRegIDANS                            MessageTag = "EAA_ISS_REG_ID_ANS"
	MessageTagEAARevPR                                  MessageTag = "EAA_REV_PR"
	MessageTagEAARevPRANS                               MessageTag = "EAA_REV_PR_ANS"
	MessageTagEAARevAV                                  MessageTag = "EAA_REV_AV"
	MessageTagEAARevAVANS                               MessageTag = "EAA_REV_AV_ANS"
	MessageTagEAARevACC                                 MessageTag = "EAA_REV_ACC"
	MessageTagEAARevACCANS                              MessageTag = "EAA_REV_ACC_ANS"
	MessageTagEAARevACCFND                              MessageTag = "EAA_REV_ACC_FND"
	MessageTagEAARevACCFNDANS                           MessageTag = "EAA_REV_ACC_FND_ANS"
	MessageTagEAARevNotRev                              MessageTag = "EAA_REV_NOT_REV"
	MessageTagEAARevNotRevANS                           MessageTag = "EAA_REV_NOT_REV_ANS"
	MessageTagEAARevNotOnHold                           MessageTag = "EAA_REV_NOT_ON_HOLD"
	MessageTagEAARevNotOnHoldANS                        MessageTag = "EAA_REV_NOT_ON_HOLD_ANS"
	MessageTagEAASHLVD                                  MessageTag = "EAA_SH_LVD"
	MessageTagEAASHLVDANS                               MessageTag = "EAA_SH_LVD_ANS"
	MessageTagEAAOTU                                    MessageTag = "EAA_OTU"
	MessageTagEAAOTUANS                                 MessageTag = "EAA_OTU_ANS"
	MessageTagEAAPseudoUsed                             MessageTag = "EAA_PSEUDO_USED"
	MessageTagEAAPseudoUsedANS                          MessageTag = "EAA_PSEUDO_USED_ANS"
	MessageTagSDJWTEAAVCTPresent                        MessageTag = "SDJWT_EAA_VCT_PRESENT"
	MessageTagSDJWTEAAVCTPresentANS                     MessageTag = "SDJWT_EAA_VCT_PRESENT_ANS"
	MessageTagSDJWTEAAVCTIntPresent                     MessageTag = "SDJWT_EAA_VCT_INT_PRESENT"
	MessageTagSDJWTEAAVCTIntPresentANS                  MessageTag = "SDJWT_EAA_VCT_INT_PRESENT_ANS"
	MessageTagEAARevType                                MessageTag = "EAA_REV_TYPE"
	MessageTagEAARevTypeANS                             MessageTag = "EAA_REV_TYPE_ANS"
	MessageTagEAARevKnown                               MessageTag = "EAA_REV_KNOWN"
	MessageTagEAARevKnownANS                            MessageTag = "EAA_REV_KNOWN_ANS"
	MessageTagEAARevISS                                 MessageTag = "EAA_REV_ISS"
	MessageTagEAARevISSANS                              MessageTag = "EAA_REV_ISS_ANS"
	MessageTagEAARevExp                                 MessageTag = "EAA_REV_EXP"
	MessageTagEAARevExpANS                              MessageTag = "EAA_REV_EXP_ANS"
	MessageTagEAARevNotExp                              MessageTag = "EAA_REV_NOT_EXP"
	MessageTagEAARevNotExpANS                           MessageTag = "EAA_REV_NOT_EXP_ANS"
	MessageTagEAARevSub                                 MessageTag = "EAA_REV_SUB"
	MessageTagEAARevSubANS                              MessageTag = "EAA_REV_SUB_ANS"
	MessageTagEAARevSubMatch                            MessageTag = "EAA_REV_SUB_MATCH"
	MessageTagEAARevSubMatchANS                         MessageTag = "EAA_REV_SUB_MATCH_ANS"
	MessageTagEAARevISSValid                            MessageTag = "EAA_REV_ISS_VALID"
	MessageTagEAARevISSValidANS                         MessageTag = "EAA_REV_ISS_VALID_ANS"
	MessageTagEAARevTime                                MessageTag = "EAA_REV_TIME"
	MessageTagEAARevISSCert                             MessageTag = "EAA_REV_ISS_CERT"
	MessageTagQualTLExp                                 MessageTag = "QUAL_TL_EXP"
	MessageTagQualTLExpANS                              MessageTag = "QUAL_TL_EXP_ANS"
	MessageTagQualTLFresh                               MessageTag = "QUAL_TL_FRESH"
	MessageTagQualTLFreshANS                            MessageTag = "QUAL_TL_FRESH_ANS"
	MessageTagQualTLVersion                             MessageTag = "QUAL_TL_VERSION"
	MessageTagQualTLVersionANS                          MessageTag = "QUAL_TL_VERSION_ANS"
	MessageTagQualTLWS                                  MessageTag = "QUAL_TL_WS"
	MessageTagQualTLWSANS                               MessageTag = "QUAL_TL_WS_ANS"
	MessageTagQualTLSV                                  MessageTag = "QUAL_TL_SV"
	MessageTagQualTLSVANS                               MessageTag = "QUAL_TL_SV_ANS"
	MessageTagQualTLIMRA                                MessageTag = "QUAL_TL_IMRA"
	MessageTagQualTLIMRAANS                             MessageTag = "QUAL_TL_IMRA_ANS"
	MessageTagQualTLIMRAANSV1                           MessageTag = "QUAL_TL_IMRA_ANS_V1"
	MessageTagQualTLIMRAANSV2                           MessageTag = "QUAL_TL_IMRA_ANS_V2"
	MessageTagQualTLServCONS                            MessageTag = "QUAL_TL_SERV_CONS"
	MessageTagQualTLServCONSANS0                        MessageTag = "QUAL_TL_SERV_CONS_ANS0"
	MessageTagQualTLServCONSANS1                        MessageTag = "QUAL_TL_SERV_CONS_ANS1"
	MessageTagQualTLServCONSANS2                        MessageTag = "QUAL_TL_SERV_CONS_ANS2"
	MessageTagQualTLServCONSANS3                        MessageTag = "QUAL_TL_SERV_CONS_ANS3"
	MessageTagQualTLServCONSANS3A                       MessageTag = "QUAL_TL_SERV_CONS_ANS3A"
	MessageTagQualTLServCONSANS3B                       MessageTag = "QUAL_TL_SERV_CONS_ANS3B"
	MessageTagQualTLServCONSANS3C                       MessageTag = "QUAL_TL_SERV_CONS_ANS3C"
	MessageTagQualTLServCONSANS4                        MessageTag = "QUAL_TL_SERV_CONS_ANS4"
	MessageTagQualTLServCONSANS5                        MessageTag = "QUAL_TL_SERV_CONS_ANS5"
	MessageTagQualTLServCONSANS6                        MessageTag = "QUAL_TL_SERV_CONS_ANS6"
	MessageTagQualTLServCONSANS7                        MessageTag = "QUAL_TL_SERV_CONS_ANS7"
	MessageTagQualCertTrustedListReached                MessageTag = "QUAL_CERT_TRUSTED_LIST_REACHED"
	MessageTagQualCertTrustedListReachedANS             MessageTag = "QUAL_CERT_TRUSTED_LIST_REACHED_ANS"
	MessageTagQualTrustedListAccept                     MessageTag = "QUAL_TRUSTED_LIST_ACCEPT"
	MessageTagQualTrustedListAcceptANS                  MessageTag = "QUAL_TRUSTED_LIST_ACCEPT_ANS"
	MessageTagQualListOfTrustedListsAccept              MessageTag = "QUAL_LIST_OF_TRUSTED_LISTS_ACCEPT"
	MessageTagQualListOfTrustedListsAcceptANS           MessageTag = "QUAL_LIST_OF_TRUSTED_LISTS_ACCEPT_ANS"
	MessageTagQualValidTrustedListPresent               MessageTag = "QUAL_VALID_TRUSTED_LIST_PRESENT"
	MessageTagQualValidTrustedListPresentANS            MessageTag = "QUAL_VALID_TRUSTED_LIST_PRESENT_ANS"
	MessageTagQualCertTypeAtST                          MessageTag = "QUAL_CERT_TYPE_AT_ST"
	MessageTagQualCertTypeAtSTANS                       MessageTag = "QUAL_CERT_TYPE_AT_ST_ANS"
	MessageTagQualCertTypeAtCC                          MessageTag = "QUAL_CERT_TYPE_AT_CC"
	MessageTagQualCertTypeAtCCANS                       MessageTag = "QUAL_CERT_TYPE_AT_CC_ANS"
	MessageTagQualCertTypeAtVT                          MessageTag = "QUAL_CERT_TYPE_AT_VT"
	MessageTagQualCertTypeAtVTANS                       MessageTag = "QUAL_CERT_TYPE_AT_VT_ANS"
	MessageTagQualQCAtST                                MessageTag = "QUAL_QC_AT_ST"
	MessageTagQualQCAtSTANS                             MessageTag = "QUAL_QC_AT_ST_ANS"
	MessageTagQualQCAtCC                                MessageTag = "QUAL_QC_AT_CC"
	MessageTagQualQCAtCCANS                             MessageTag = "QUAL_QC_AT_CC_ANS"
	MessageTagQualQCAtVT                                MessageTag = "QUAL_QC_AT_VT"
	MessageTagQualQCAtVTANS                             MessageTag = "QUAL_QC_AT_VT_ANS"
	MessageTagQualQSCDAtST                              MessageTag = "QUAL_QSCD_AT_ST"
	MessageTagQualQSCDAtSTANS                           MessageTag = "QUAL_QSCD_AT_ST_ANS"
	MessageTagQualQSCDAtCC                              MessageTag = "QUAL_QSCD_AT_CC"
	MessageTagQualQSCDAtCCANS                           MessageTag = "QUAL_QSCD_AT_CC_ANS"
	MessageTagQualQSCDAtVT                              MessageTag = "QUAL_QSCD_AT_VT"
	MessageTagQualQSCDAtVTANS                           MessageTag = "QUAL_QSCD_AT_VT_ANS"
	MessageTagQualUniqueCert                            MessageTag = "QUAL_UNIQUE_CERT"
	MessageTagQualUniqueCertANS                         MessageTag = "QUAL_UNIQUE_CERT_ANS"
	MessageTagQualIsAdES                                MessageTag = "QUAL_IS_ADES"
	MessageTagQualIsAdESInd                             MessageTag = "QUAL_IS_ADES_IND"
	MessageTagQualIsAdESINV                             MessageTag = "QUAL_IS_ADES_INV"
	MessageTagQualHasMETS                               MessageTag = "QUAL_HAS_METS"
	MessageTagQualHasMETSANS                            MessageTag = "QUAL_HAS_METS_ANS"
	MessageTagQualHasMETSAtTime                         MessageTag = "QUAL_HAS_METS_ATTIME"
	MessageTagQualHasMETSAtTimeANS                      MessageTag = "QUAL_HAS_METS_ATTIME_ANS"
	MessageTagQualHasMETSHCCECBA                        MessageTag = "QUAL_HAS_METS_HCCECBA"
	MessageTagQualHasMETSHCCECBAANS                     MessageTag = "QUAL_HAS_METS_HCCECBA_ANS"
	MessageTagQualHasMETSHCCECBAANS2                    MessageTag = "QUAL_HAS_METS_HCCECBA_ANS_2"
	MessageTagQualHasMETSHCCECBAANS3                    MessageTag = "QUAL_HAS_METS_HCCECBA_ANS_3"
	MessageTagQualHasCAQC                               MessageTag = "QUAL_HAS_CAQC"
	MessageTagQualHasCAQCANS                            MessageTag = "QUAL_HAS_CAQC_ANS"
	MessageTagQualHasCAQCANS2                           MessageTag = "QUAL_HAS_CAQC_ANS_2"
	MessageTagQualHasAtTime                             MessageTag = "QUAL_HAS_ATTIME"
	MessageTagQualHasAtTimeANS                          MessageTag = "QUAL_HAS_ATTIME_ANS"
	MessageTagQualHasTSCertType                         MessageTag = "QUAL_HAS_TS_CERT_TYPE"
	MessageTagQualHasTSCertTypeANS                      MessageTag = "QUAL_HAS_TS_CERT_TYPE_ANS"
	MessageTagQualHasConf                               MessageTag = "QUAL_HAS_CONF"
	MessageTagQualHasConfANS                            MessageTag = "QUAL_HAS_CONF_ANS"
	MessageTagQualHasQEAA                               MessageTag = "QUAL_HAS_QEAA"
	MessageTagQualHasQEAAANS                            MessageTag = "QUAL_HAS_QEAA_ANS"
	MessageTagQualHasQTST                               MessageTag = "QUAL_HAS_QTST"
	MessageTagQualHasQTSTANS                            MessageTag = "QUAL_HAS_QTST_ANS"
	MessageTagQualIsTrustCertMatchService               MessageTag = "QUAL_IS_TRUST_CERT_MATCH_SERVICE"
	MessageTagQualIsTrustCertMatchServiceANS0           MessageTag = "QUAL_IS_TRUST_CERT_MATCH_SERVICE_ANS0"
	MessageTagQualIsTrustCertMatchServiceANS1           MessageTag = "QUAL_IS_TRUST_CERT_MATCH_SERVICE_ANS1"
	MessageTagQualIsTrustCertMatchServiceANS2           MessageTag = "QUAL_IS_TRUST_CERT_MATCH_SERVICE_ANS2"
	MessageTagQualHasGranted                            MessageTag = "QUAL_HAS_GRANTED"
	MessageTagQualHasGrantedANS                         MessageTag = "QUAL_HAS_GRANTED_ANS"
	MessageTagQualHasGrantedANS2                        MessageTag = "QUAL_HAS_GRANTED_ANS_2"
	MessageTagQualHasGrantedAt                          MessageTag = "QUAL_HAS_GRANTED_AT"
	MessageTagQualHasGrantedAtANS                       MessageTag = "QUAL_HAS_GRANTED_AT_ANS"
	MessageTagQualHasConsistentByQC                     MessageTag = "QUAL_HAS_CONSISTENT_BY_QC"
	MessageTagQualHasConsistentByQCANS                  MessageTag = "QUAL_HAS_CONSISTENT_BY_QC_ANS"
	MessageTagQualHasConsistentByQSCD                   MessageTag = "QUAL_HAS_CONSISTENT_BY_QSCD"
	MessageTagQualHasConsistentByQSCDANS                MessageTag = "QUAL_HAS_CONSISTENT_BY_QSCD_ANS"
	MessageTagQualHasCertTypeCoverage                   MessageTag = "QUAL_HAS_CERT_TYPE_COVERAGE"
	MessageTagQualHasCertTypeCoverageANS                MessageTag = "QUAL_HAS_CERT_TYPE_COVERAGE_ANS"
	MessageTagQualHasValidCAQC                          MessageTag = "QUAL_HAS_VALID_CAQC"
	MessageTagQualHasValidCAQCANS                       MessageTag = "QUAL_HAS_VALID_CAQC_ANS"
	MessageTagQualHasOnlyOne                            MessageTag = "QUAL_HAS_ONLY_ONE"
	MessageTagQualHasOnlyOneANS                         MessageTag = "QUAL_HAS_ONLY_ONE_ANS"
	MessageTagQWACValid                                 MessageTag = "QWAC_VALID"
	MessageTagQWACValidANS                              MessageTag = "QWAC_VALID_ANS"
	MessageTagQWACValidANS2                             MessageTag = "QWAC_VALID_ANS_2"
	MessageTagQWACCertQualConclusive                    MessageTag = "QWAC_CERT_QUAL_CONCLUSIVE"
	MessageTagQWACCertQualConclusiveANS                 MessageTag = "QWAC_CERT_QUAL_CONCLUSIVE_ANS"
	MessageTagQWACIsWSAAtTime                           MessageTag = "QWAC_IS_WSA_AT_TIME"
	MessageTagQWACIsWSAAtTimeANS                        MessageTag = "QWAC_IS_WSA_AT_TIME_ANS"
	MessageTagQWACCertPolicy                            MessageTag = "QWAC_CERT_POLICY"
	MessageTagQWACCertPolicyANS                         MessageTag = "QWAC_CERT_POLICY_ANS"
	MessageTagQWACValPeriod                             MessageTag = "QWAC_VAL_PERIOD"
	MessageTagQWACValPeriodANS                          MessageTag = "QWAC_VAL_PERIOD_ANS"
	MessageTagQWACDomainName                            MessageTag = "QWAC_DOMAIN_NAME"
	MessageTagQWACDomainNameANS                         MessageTag = "QWAC_DOMAIN_NAME_ANS"
	MessageTagQWAC2ExtKeyUsage                          MessageTag = "QWAC2_EXT_KEY_USAGE"
	MessageTagQWAC2ExtKeyUsageANS                       MessageTag = "QWAC2_EXT_KEY_USAGE_ANS"
	MessageTagTLSCertBindingURL                         MessageTag = "TLS_CERT_BINDING_URL"
	MessageTagTLSCertBindingURLANS                      MessageTag = "TLS_CERT_BINDING_URL_ANS"
	MessageTagTLSCertBindingSig                         MessageTag = "TLS_CERT_BINDING_SIG"
	MessageTagTLSCertBindingSigANS                      MessageTag = "TLS_CERT_BINDING_SIG_ANS"
	MessageTagTLSCertBindingSigForm                     MessageTag = "TLS_CERT_BINDING_SIG_FORM"
	MessageTagTLSCertBindingSigFormANS                  MessageTag = "TLS_CERT_BINDING_SIG_FORM_ANS"
	MessageTagTLSCertBindingSigSer                      MessageTag = "TLS_CERT_BINDING_SIG_SER"
	MessageTagTLSCertBindingSigSerANS                   MessageTag = "TLS_CERT_BINDING_SIG_SER_ANS"
	MessageTagTLSCertBindingSigExp                      MessageTag = "TLS_CERT_BINDING_SIG_EXP"
	MessageTagTLSCertBindingSigExpANS                   MessageTag = "TLS_CERT_BINDING_SIG_EXP_ANS"
	MessageTagTLSCertBindingSigExpiryDate               MessageTag = "TLS_CERT_BINDING_SIG_EXPIRY_DATE"
	MessageTagTLSCertBindingSigExpiryDateANS            MessageTag = "TLS_CERT_BINDING_SIG_EXPIRY_DATE_ANS"
	MessageTagTLSCertBindingQWAC2                       MessageTag = "TLS_CERT_BINDING_QWAC2"
	MessageTagTLSCertBindingQWAC2ANS                    MessageTag = "TLS_CERT_BINDING_QWAC2_ANS"
	MessageTagTLSCertBindingSigValid                    MessageTag = "TLS_CERT_BINDING_SIG_VALID"
	MessageTagTLSCertBindingSigValidANS                 MessageTag = "TLS_CERT_BINDING_SIG_VALID_ANS"
	MessageTagTLSCertBindingCertIdentified              MessageTag = "TLS_CERT_BINDING_CERT_IDENTIFIED"
	MessageTagTLSCertBindingCertIdentifiedANS           MessageTag = "TLS_CERT_BINDING_CERT_IDENTIFIED_ANS"
	MessageTagCertUsageLoTEAccept                       MessageTag = "CERT_USAGE_LOTE_ACCEPT"
	MessageTagCertUsageLoTEAcceptANS                    MessageTag = "CERT_USAGE_LOTE_ACCEPT_ANS"
	MessageTagCertUsageLoLoTEAccept                     MessageTag = "CERT_USAGE_LOLOTE_ACCEPT"
	MessageTagCertUsageLoLoTEAcceptANS                  MessageTag = "CERT_USAGE_LOLOTE_ACCEPT_ANS"
	MessageTagCertUsageValidLoTEPresent                 MessageTag = "CERT_USAGE_VALID_LOTE_PRESENT"
	MessageTagCertUsageValidLoTEPresentANS              MessageTag = "CERT_USAGE_VALID_LOTE_PRESENT_ANS"
	MessageTagCertUsageHasAtTime                        MessageTag = "CERT_USAGE_HAS_ATTIME"
	MessageTagCertUsageHasAtTimeANS                     MessageTag = "CERT_USAGE_HAS_ATTIME_ANS"
	MessageTagCertUsageListTypeKnown                    MessageTag = "CERT_USAGE_LIST_TYPE_KNOWN"
	MessageTagCertUsageListTypeKnownANS                 MessageTag = "CERT_USAGE_LIST_TYPE_KNOWN_ANS"
	MessageTagCertUsageStatus                           MessageTag = "CERT_USAGE_STATUS"
	MessageTagCertUsageStatusANS                        MessageTag = "CERT_USAGE_STATUS_ANS"
	MessageTagCertUsageStatusCONS                       MessageTag = "CERT_USAGE_STATUS_CONS"
	MessageTagCertUsageStatusCONSANS                    MessageTag = "CERT_USAGE_STATUS_CONS_ANS"
	MessageTagCertUsageStatusKnown                      MessageTag = "CERT_USAGE_STATUS_KNOWN"
	MessageTagCertUsageStatusKnownANS                   MessageTag = "CERT_USAGE_STATUS_KNOWN_ANS"
	MessageTagCertUsageSti                              MessageTag = "CERT_USAGE_STI"
	MessageTagCertUsageStiANS                           MessageTag = "CERT_USAGE_STI_ANS"
	MessageTagCertUsageStiKnown                         MessageTag = "CERT_USAGE_STI_KNOWN"
	MessageTagCertUsageStiKnownANS                      MessageTag = "CERT_USAGE_STI_KNOWN_ANS"
	MessageTagPIDDocumentType                           MessageTag = "PID_DOCUMENT_TYPE"
	MessageTagPIDDocumentTypeANS                        MessageTag = "PID_DOCUMENT_TYPE_ANS"
	MessageTagPIDLoTETypePIDProviders                   MessageTag = "PID_LOTE_TYPE_PID_PROVIDERS"
	MessageTagPIDLoTETypePIDProvidersANS                MessageTag = "PID_LOTE_TYPE_PID_PROVIDERS_ANS"
	MessageTagPIDStiPIDIssuance                         MessageTag = "PID_STI_PID_ISSUANCE"
	MessageTagPIDStiPIDIssuanceANS                      MessageTag = "PID_STI_PID_ISSUANCE_ANS"
	MessageTagPIDProviderAtIssuanceTime                 MessageTag = "PID_PROVIDER_AT_ISSUANCE_TIME"
	MessageTagPIDProviderAtIssuanceTimeANS              MessageTag = "PID_PROVIDER_AT_ISSUANCE_TIME_ANS"
	MessageTagPIDProviderAtValidationTime               MessageTag = "PID_PROVIDER_AT_VALIDATION_TIME"
	MessageTagPIDProviderAtValidationTimeANS            MessageTag = "PID_PROVIDER_AT_VALIDATION_TIME_ANS"
	MessageTagBBBAccept                                 MessageTag = "BBB_ACCEPT"
	MessageTagBBBAcceptANS                              MessageTag = "BBB_ACCEPT_ANS"
	MessageTagTSTTypeContentTST                         MessageTag = "TST_TYPE_CONTENT_TST"
	MessageTagTSTTypeSignatureTST                       MessageTag = "TST_TYPE_SIGNATURE_TST"
	MessageTagTSTTypeVDTST                              MessageTag = "TST_TYPE_VD_TST"
	MessageTagTSTTypeDocTST                             MessageTag = "TST_TYPE_DOC_TST"
	MessageTagTSTTypeContainerTST                       MessageTag = "TST_TYPE_CONTAINER_TST"
	MessageTagTSTTypeArchiveTST                         MessageTag = "TST_TYPE_ARCHIVE_TST"
	MessageTagTSTTypeERTST                              MessageTag = "TST_TYPE_ER_TST"
	MessageTagTSTTypeRefERATST                          MessageTag = "TST_TYPE_REF_ER_ATST"
	MessageTagTSTTypeRefERATSTSeq                       MessageTag = "TST_TYPE_REF_ER_ATST_SEQ"
	MessageTagEmpty                                     MessageTag = "EMPTY"
	MessageTagCertificate                               MessageTag = "CERTIFICATE"
	MessageTagCACertificate                             MessageTag = "CA_CERTIFICATE"
	MessageTagSigningCertificate                        MessageTag = "SIGNING_CERTIFICATE"
	MessageTagRevocation                                MessageTag = "REVOCATION"
	MessageTagRevocationSigCert                         MessageTag = "REVOCATION_SIG_CERT"
	MessageTagRevocationCACert                          MessageTag = "REVOCATION_CA_CERT"
	MessageTagSignature                                 MessageTag = "SIGNATURE"
	MessageTagTimestamp                                 MessageTag = "TIMESTAMP"
	MessageTagTimestampSigCert                          MessageTag = "TIMESTAMP_SIG_CERT"
	MessageTagTimestampCACert                           MessageTag = "TIMESTAMP_CA_CERT"
	MessageTagEAARev                                    MessageTag = "EAA_REV"
	MessageTagEAARevSigCert                             MessageTag = "EAA_REV_SIG_CERT"
	MessageTagEAARevCACert                              MessageTag = "EAA_REV_CA_CERT"
	MessageTagAcceptableRevocation                      MessageTag = "ACCEPTABLE_REVOCATION"
	MessageTagBasicSignatureValidationResult            MessageTag = "BASIC_SIGNATURE_VALIDATION_RESULT"
	MessageTagBESTSignatureTimeCertNotAfter             MessageTag = "BEST_SIGNATURE_TIME_CERT_NOT_AFTER"
	MessageTagBESTSignatureTimeCertNotBefore            MessageTag = "BEST_SIGNATURE_TIME_CERT_NOT_BEFORE"
	MessageTagBESTSignatureTimeCertRevocation           MessageTag = "BEST_SIGNATURE_TIME_CERT_REVOCATION"
	MessageTagBESTSignatureTimeCertSuspension           MessageTag = "BEST_SIGNATURE_TIME_CERT_SUSPENSION"
	MessageTagCertificateID                             MessageTag = "CERTIFICATE_ID"
	MessageTagCertificateRevocationFound                MessageTag = "CERTIFICATE_REVOCATION_FOUND"
	MessageTagCertificateRevocationNotFound             MessageTag = "CERTIFICATE_REVOCATION_NOT_FOUND"
	MessageTagCertificateSunsetDate                     MessageTag = "CERTIFICATE_SUNSET_DATE"
	MessageTagCertificateSunsetDateTrustAnchor          MessageTag = "CERTIFICATE_SUNSET_DATE_TRUST_ANCHOR"
	MessageTagCertificateSunsetDateValid                MessageTag = "CERTIFICATE_SUNSET_DATE_VALID"
	MessageTagCertificateType                           MessageTag = "CERTIFICATE_TYPE"
	MessageTagCertificateValidity                       MessageTag = "CERTIFICATE_VALIDITY"
	MessageTagCertificateUsage                          MessageTag = "CERTIFICATE_USAGE"
	MessageTagCertificateUsageListType                  MessageTag = "CERTIFICATE_USAGE_LIST_TYPE"
	MessageTagCertificateUsageStatus                    MessageTag = "CERTIFICATE_USAGE_STATUS"
	MessageTagCertificateUsageStatuses                  MessageTag = "CERTIFICATE_USAGE_STATUSES"
	MessageTagCertificateUsageSti                       MessageTag = "CERTIFICATE_USAGE_STI"
	MessageTagCertificateUsageSTIS                      MessageTag = "CERTIFICATE_USAGE_STIS"
	MessageTagControlTime                               MessageTag = "CONTROL_TIME"
	MessageTagControlTimeAlone                          MessageTag = "CONTROL_TIME_ALONE"
	MessageTagControlTimeWithPOE                        MessageTag = "CONTROL_TIME_WITH_POE"
	MessageTagControlTimeWithTrustAnchor                MessageTag = "CONTROL_TIME_WITH_TRUST_ANCHOR"
	MessageTagCryptographicCheckFailure                 MessageTag = "CRYPTOGRAPHIC_CHECK_FAILURE"
	MessageTagCryptographicCheckFailureWithID           MessageTag = "CRYPTOGRAPHIC_CHECK_FAILURE_WITH_ID"
	MessageTagCryptographicCheckFailureWithRef          MessageTag = "CRYPTOGRAPHIC_CHECK_FAILURE_WITH_REF"
	MessageTagCryptographicCheckFailureWithRefWithName  MessageTag = "CRYPTOGRAPHIC_CHECK_FAILURE_WITH_REF_WITH_NAME"
	MessageTagCryptographicCheckFailureWithRefWithNames MessageTag = "CRYPTOGRAPHIC_CHECK_FAILURE_WITH_REF_WITH_NAMES"
	MessageTagCryptographicCheckSuccess                 MessageTag = "CRYPTOGRAPHIC_CHECK_SUCCESS"
	MessageTagCryptographicCheckSuccessKeySize          MessageTag = "CRYPTOGRAPHIC_CHECK_SUCCESS_KEY_SIZE"
	MessageTagCryptographicCheckSuccessDM               MessageTag = "CRYPTOGRAPHIC_CHECK_SUCCESS_DM"
	MessageTagCryptographicCheckSuccessDMWithID         MessageTag = "CRYPTOGRAPHIC_CHECK_SUCCESS_DM_WITH_ID"
	MessageTagCryptographicCheckSuccessDMWithName       MessageTag = "CRYPTOGRAPHIC_CHECK_SUCCESS_DM_WITH_NAME"
	MessageTagCryptographicCheckSuccessDMWithNames      MessageTag = "CRYPTOGRAPHIC_CHECK_SUCCESS_DM_WITH_NAMES"
	MessageTagEvidenceRecordValidation                  MessageTag = "EVIDENCE_RECORD_VALIDATION"
	MessageTagExtendedKeyUsage                          MessageTag = "EXTENDED_KEY_USAGE"
	MessageTagKeyUsage                                  MessageTag = "KEY_USAGE"
	MessageTagLastAcceptableRevocation                  MessageTag = "LAST_ACCEPTABLE_REVOCATION"
	MessageTagListOfTrustedEntities                     MessageTag = "LIST_OF_TRUSTED_ENTITIES"
	MessageTagPseudo                                    MessageTag = "PSEUDO"
	MessageTagQWACExpiryExp                             MessageTag = "QWAC_EXPIRY_EXP"
	MessageTagQWACExpiryTLSCert                         MessageTag = "QWAC_EXPIRY_TLS_CERT"
	MessageTagQWACExpirySignCert                        MessageTag = "QWAC_EXPIRY_SIGN_CERT"
	MessageTagReference                                 MessageTag = "REFERENCE"
	MessageTagReferenceNameCheck                        MessageTag = "REFERENCE_NAME_CHECK"
	MessageTagReferencesWithNames                       MessageTag = "REFERENCES_WITH_NAMES"
	MessageTagRevocationAcceptanceCheck                 MessageTag = "REVOCATION_ACCEPTANCE_CHECK"
	MessageTagRevocationCertHashOK                      MessageTag = "REVOCATION_CERT_HASH_OK"
	MessageTagRevocationCertHashOKID                    MessageTag = "REVOCATION_CERT_HASH_OK_ID"
	MessageTagRevocationCertValidity                    MessageTag = "REVOCATION_CERT_VALIDITY"
	MessageTagRevocationCheck                           MessageTag = "REVOCATION_CHECK"
	MessageTagRevocationConsistent                      MessageTag = "REVOCATION_CONSISTENT"
	MessageTagRevocationConsistentCRL                   MessageTag = "REVOCATION_CONSISTENT_CRL"
	MessageTagRevocationConsistentOCSP                  MessageTag = "REVOCATION_CONSISTENT_OCSP"
	MessageTagRevocationConsistentTL                    MessageTag = "REVOCATION_CONSISTENT_TL"
	MessageTagRevocationInfo                            MessageTag = "REVOCATION_INFO"
	MessageTagRevocationNotAfterAfter                   MessageTag = "REVOCATION_NOT_AFTER_AFTER"
	MessageTagRevocationNotAfterAfterID                 MessageTag = "REVOCATION_NOT_AFTER_AFTER_ID"
	MessageTagRevocationProducedAtCertValidity          MessageTag = "REVOCATION_PRODUCED_AT_CERT_VALIDITY"
	MessageTagRevocationProducedAtOutOfBounds           MessageTag = "REVOCATION_PRODUCED_AT_OUT_OF_BOUNDS"
	MessageTagRevocationProducedAtOutOfBoundsID         MessageTag = "REVOCATION_PRODUCED_AT_OUT_OF_BOUNDS_ID"
	MessageTagRevocationReason                          MessageTag = "REVOCATION_REASON"
	MessageTagRevocationThisUpdateControlTime           MessageTag = "REVOCATION_THIS_UPDATE_CONTROL_TIME"
	MessageTagSignatureAlgorithmWithKeySize             MessageTag = "SIGNATURE_ALGORITHM_WITH_KEY_SIZE"
	MessageTagSignatureID                               MessageTag = "SIGNATURE_ID"
	MessageTagStructuralValidationFailure               MessageTag = "STRUCTURAL_VALIDATION_FAILURE"
	MessageTagTimestampAndCertificateNotAfter           MessageTag = "TIMESTAMP_AND_CERTIFICATE_NOT_AFTER"
	MessageTagTimestampAndCryptoConstraintsExpiration   MessageTag = "TIMESTAMP_AND_CRYPTO_CONSTRAINTS_EXPIRATION"
	MessageTagTimestampAndRevocationTime                MessageTag = "TIMESTAMP_AND_REVOCATION_TIME"
	MessageTagTimestampValidation                       MessageTag = "TIMESTAMP_VALIDATION"
	MessageTagTokenID                                   MessageTag = "TOKEN_ID"
	MessageTagTrustServiceName                          MessageTag = "TRUST_SERVICE_NAME"
	MessageTagTrustedServiceStatus                      MessageTag = "TRUSTED_SERVICE_STATUS"
	MessageTagTrustedServiceType                        MessageTag = "TRUSTED_SERVICE_TYPE"
	MessageTagTrustedList                               MessageTag = "TRUSTED_LIST"
	MessageTagValidationTime                            MessageTag = "VALIDATION_TIME"
	MessageTagCryptographicVerification                 MessageTag = "CRYPTOGRAPHIC_VERIFICATION"
	MessageTagFormatChecking                            MessageTag = "FORMAT_CHECKING"
	MessageTagIdentificationOfTheSigningCertificate     MessageTag = "IDENTIFICATION_OF_THE_SIGNING_CERTIFICATE"
	MessageTagPastSignatureValidation                   MessageTag = "PAST_SIGNATURE_VALIDATION"
	MessageTagPastCertificateValidation                 MessageTag = "PAST_CERTIFICATE_VALIDATION"
	MessageTagRevocationFreshnessChecker                MessageTag = "REVOCATION_FRESHNESS_CHECKER"
	MessageTagSignatureAcceptanceValidation             MessageTag = "SIGNATURE_ACCEPTANCE_VALIDATION"
	MessageTagValidationContextInitialization           MessageTag = "VALIDATION_CONTEXT_INITIALIZATION"
	MessageTagValidationTimeSliding                     MessageTag = "VALIDATION_TIME_SLIDING"
	MessageTagX509CertificateValidation                 MessageTag = "X509_CERTIFICATE_VALIDATION"
	MessageTagResults                                   MessageTag = "RESULTS"
	MessageTagEEAType                                   MessageTag = "EEA_TYPE"
	MessageTagEAAAcceptanceValidation                   MessageTag = "EAA_ACCEPTANCE_VALIDATION"
	MessageTagAOV                                       MessageTag = "AOV"
	MessageTagCertQualification                         MessageTag = "CERT_QUALIFICATION"
	MessageTagCertQualificationAtTime                   MessageTag = "CERT_QUALIFICATION_AT_TIME"
	MessageTagCertUsageAtTime                           MessageTag = "CERT_USAGE_AT_TIME"
	MessageTagCertUsages                                MessageTag = "CERT_USAGES"
	MessageTagCC                                        MessageTag = "CC"
	MessageTagCRS                                       MessageTag = "CRS"
	MessageTagDAAV                                      MessageTag = "DAAV"
	MessageTagEAAQualification                          MessageTag = "EAA_QUALIFICATION"
	MessageTagEAAQualificationProcess                   MessageTag = "EAA_QUALIFICATION_PROCESS"
	MessageTagLoTE                                      MessageTag = "LOTE"
	MessageTagLoLoTE                                    MessageTag = "LOLOTE"
	MessageTagLOTL                                      MessageTag = "LOTL"
	MessageTagPIDQualificationProcess                   MessageTag = "PID_QUALIFICATION_PROCESS"
	MessageTagPSVCRS                                    MessageTag = "PSV_CRS"
	MessageTagQWACValidation                            MessageTag = "QWAC_VALIDATION"
	MessageTagQWACValidationProfile                     MessageTag = "QWAC_VALIDATION_PROFILE"
	MessageTagRAC                                       MessageTag = "RAC"
	MessageTagSigQualification                          MessageTag = "SIG_QUALIFICATION"
	MessageTagSubXCV                                    MessageTag = "SUB_XCV"
	MessageTagTL                                        MessageTag = "TL"
	MessageTagTSTQualification                          MessageTag = "TST_QUALIFICATION"
	MessageTagTSTQualificationAtTime                    MessageTag = "TST_QUALIFICATION_AT_TIME"
	MessageTagVPBS                                      MessageTag = "VPBS"
	MessageTagVPEAA                                     MessageTag = "VPEAA"
	MessageTagVPER                                      MessageTag = "VPER"
	MessageTagVPFLTVD                                   MessageTag = "VPFLTVD"
	MessageTagVPFRVC                                    MessageTag = "VPFRVC"
	MessageTagVpfswatsp                                 MessageTag = "VPFSWATSP"
	MessageTagVpftsp                                    MessageTag = "VPFTSP"
	MessageTagVpftspwatsp                               MessageTag = "VPFTSPWATSP"
	MessageTagVTSCRS                                    MessageTag = "VTS_CRS"
	MessageTagVTBESTSignatureTime                       MessageTag = "VT_BEST_SIGNATURE_TIME"
	MessageTagVTCertificateIssuanceTime                 MessageTag = "VT_CERTIFICATE_ISSUANCE_TIME"
	MessageTagVTValidationTime                          MessageTag = "VT_VALIDATION_TIME"
	MessageTagVTTSTGenerationTime                       MessageTag = "VT_TST_GENERATION_TIME"
	MessageTagVTTSTPOETime                              MessageTag = "VT_TST_POE_TIME"
	MessageTagQWAC1Profile                              MessageTag = "QWAC1_PROFILE"
	MessageTagQWAC2Profile                              MessageTag = "QWAC2_PROFILE"
	MessageTagTLSByQWAC2Profile                         MessageTag = "TLS_BY_QWAC2_PROFILE"
	MessageTagSemanticsTotalPassed                      MessageTag = "SEMANTICS_TOTAL_PASSED"
	MessageTagSemanticsPassed                           MessageTag = "SEMANTICS_PASSED"
	MessageTagSemanticsTotalFailed                      MessageTag = "SEMANTICS_TOTAL_FAILED"
	MessageTagSemanticsFailed                           MessageTag = "SEMANTICS_FAILED"
	MessageTagSemanticsIndeterminate                    MessageTag = "SEMANTICS_INDETERMINATE"
	MessageTagSemanticsNoSignatureFound                 MessageTag = "SEMANTICS_NO_SIGNATURE_FOUND"
	MessageTagSemanticsFormatFailure                    MessageTag = "SEMANTICS_FORMAT_FAILURE"
	MessageTagSemanticsHashFailure                      MessageTag = "SEMANTICS_HASH_FAILURE"
	MessageTagSemanticsSigCryptoFailure                 MessageTag = "SEMANTICS_SIG_CRYPTO_FAILURE"
	MessageTagSemanticsRevoked                          MessageTag = "SEMANTICS_REVOKED"
	MessageTagSemanticsExpired                          MessageTag = "SEMANTICS_EXPIRED"
	MessageTagSemanticsNotYetValid                      MessageTag = "SEMANTICS_NOT_YET_VALID"
	MessageTagSemanticsSigConstraintsFailure            MessageTag = "SEMANTICS_SIG_CONSTRAINTS_FAILURE"
	MessageTagSemanticsChainConstraintsFailure          MessageTag = "SEMANTICS_CHAIN_CONSTRAINTS_FAILURE"
	MessageTagSemanticsCertificateChainGeneralFailure   MessageTag = "SEMANTICS_CERTIFICATE_CHAIN_GENERAL_FAILURE"
	MessageTagSemanticsCryptoConstraintsFailure         MessageTag = "SEMANTICS_CRYPTO_CONSTRAINTS_FAILURE"
	MessageTagSemanticsPolicyProcessingError            MessageTag = "SEMANTICS_POLICY_PROCESSING_ERROR"
	MessageTagSemanticsSignaturePolicyNotAvailable      MessageTag = "SEMANTICS_SIGNATURE_POLICY_NOT_AVAILABLE"
	MessageTagSemanticsTimestampOrderFailure            MessageTag = "SEMANTICS_TIMESTAMP_ORDER_FAILURE"
	MessageTagSemanticsNoSigningCertificateFound        MessageTag = "SEMANTICS_NO_SIGNING_CERTIFICATE_FOUND"
	MessageTagSemanticsNoCertificateChainFound          MessageTag = "SEMANTICS_NO_CERTIFICATE_CHAIN_FOUND"
	MessageTagSemanticsNoCertificateChainFoundNoPOE     MessageTag = "SEMANTICS_NO_CERTIFICATE_CHAIN_FOUND_NO_POE"
	MessageTagSemanticsRevokedNoPOE                     MessageTag = "SEMANTICS_REVOKED_NO_POE"
	MessageTagSemanticsRevokedCANoPOE                   MessageTag = "SEMANTICS_REVOKED_CA_NO_POE"
	MessageTagSemanticsOutOfBoundsNotRevoked            MessageTag = "SEMANTICS_OUT_OF_BOUNDS_NOT_REVOKED"
	MessageTagSemanticsOutOfBoundsNoPOE                 MessageTag = "SEMANTICS_OUT_OF_BOUNDS_NO_POE"
	MessageTagSemanticsRevocationOutOfBoundsNoPOE       MessageTag = "SEMANTICS_REVOCATION_OUT_OF_BOUNDS_NO_POE"
	MessageTagSemanticsCryptoConstraintsFailureNoPOE    MessageTag = "SEMANTICS_CRYPTO_CONSTRAINTS_FAILURE_NO_POE"
	MessageTagSemanticsNoPOE                            MessageTag = "SEMANTICS_NO_POE"
	MessageTagSemanticsTryLater                         MessageTag = "SEMANTICS_TRY_LATER"
	MessageTagSemanticsSignedDataNotFound               MessageTag = "SEMANTICS_SIGNED_DATA_NOT_FOUND"
	MessageTagSemanticsEAAConstraintsFailure            MessageTag = "SEMANTICS_EAA_CONSTRAINTS_FAILURE"
)

// MessageTagValues returns all MessageTag constants in declaration order.
func MessageTagValues() []MessageTag {
	return []MessageTag{
		MessageTagBBBFCIEFF,
		MessageTagBBBFCIEFFANS,
		MessageTagBBBFCICFD,
		MessageTagBBBFCICFDANS,
		MessageTagBBBFCISD,
		MessageTagBBBFCISDANS,
		MessageTagBBBFCISRIA,
		MessageTagBBBFCISRIAANS,
		MessageTagBBBFCIOSIP,
		MessageTagBBBFCIOSIPANS,
		MessageTagBBBFCDASTHVBR,
		MessageTagBBBFCDASTHVBRANS,
		MessageTagBBBFCDBTOOST,
		MessageTagBBBFCDBTOOSTANS,
		MessageTagBBBFCIBRV,
		MessageTagBBBFCIBRVANS,
		MessageTagBBBFCISDC,
		MessageTagBBBFCISDCANS,
		MessageTagBBBFCDSFREAP,
		MessageTagBBBFCDSFREAPANS,
		MessageTagBBBFCIAOD,
		MessageTagBBBFCIAODANS,
		MessageTagBBBFCIVDBSFR,
		MessageTagBBBFCIVDBSFRANS,
		MessageTagBBBFCISVADMDPD,
		MessageTagBBBFCISVADMDPDANS,
		MessageTagBBBFCISVAFMDPD,
		MessageTagBBBFCISVAFMDPDANS,
		MessageTagBBBFCISVASFLD,
		MessageTagBBBFCISVASFLDANS,
		MessageTagBBBFCDSCNFFSM,
		MessageTagBBBFCDSCNFFSMANS,
		MessageTagBBBFCDSCNACMDM,
		MessageTagBBBFCDSCNACMDMANS,
		MessageTagBBBFCDSCNUOM,
		MessageTagBBBFCDSCNUOMANS,
		MessageTagBBBFCIECKSCDA,
		MessageTagBBBFCIECKSCDAANS1,
		MessageTagBBBFCIECKSCDAANS2,
		MessageTagBBBFCIECKSCDAANS3,
		MessageTagBBBFCIECKSCDAANS4,
		MessageTagBBBFCIECKSCDAANS5,
		MessageTagBBBFCDDAPDFAF,
		MessageTagBBBFCDDAPDFAFANS,
		MessageTagBBBFCIDPDFAC,
		MessageTagBBBFCIDPDFACANS,
		MessageTagBBBFCIECTF,
		MessageTagBBBFCIECTFANS,
		MessageTagBBBFCISFCS,
		MessageTagBBBFCISFCSANS,
		MessageTagBBBFCITFCS,
		MessageTagBBBFCITFCSANS,
		MessageTagBBBFCIMFCS,
		MessageTagBBBFCIMFCSANS,
		MessageTagBBBFCITZCP,
		MessageTagBBBFCITZCPANS,
		MessageTagBBBFCITEZCF,
		MessageTagBBBFCITEZCFANS,
		MessageTagBBBFCITMFP,
		MessageTagBBBFCITMFPANS,
		MessageTagBBBFCIEMCF,
		MessageTagBBBFCIEMCFANS,
		MessageTagBBBFCIMFPASiCE,
		MessageTagBBBFCIMFPASiCEANS,
		MessageTagBBBFCISFPASiCE,
		MessageTagBBBFCISFPASiCEANS,
		MessageTagBBBFCISFPASiCS,
		MessageTagBBBFCISFPASiCSANS,
		MessageTagBBBFCISFPASTFORAMC,
		MessageTagBBBFCISFPASTFORAMCANS,
		MessageTagBBBFCIAHIV,
		MessageTagBBBFCIAHIVANS,
		MessageTagBBBCVIRDOF,
		MessageTagBBBCVIRDOFANS,
		MessageTagBBBCVTSPIRDOF,
		MessageTagBBBCVTSPIRDOFANS,
		MessageTagBBBCVCSCSSVF,
		MessageTagBBBCVCSCSSVFANS,
		MessageTagBBBCVIMEOF,
		MessageTagBBBCVIMEOFANS,
		MessageTagBBBCVEAASDCBF,
		MessageTagBBBCVEAASDCBFANS,
		MessageTagBBBCVEAANSDCBF,
		MessageTagBBBCVEAANSDCBFANS,
		MessageTagBBBCVERIODOF,
		MessageTagBBBCVERIODOFANS,
		MessageTagBBBCVERHasSDoc,
		MessageTagBBBCVERHasSDocANS,
		MessageTagBBBCVERDFHVLCDOG,
		MessageTagBBBCVERDFHVLCDOGANS,
		MessageTagBBBCVERATSRF,
		MessageTagBBBCVERATSRFANS,
		MessageTagBBBCVERATSSRF,
		MessageTagBBBCVERATSSRFANS,
		MessageTagBBBCVIRDOI,
		MessageTagBBBCVIRDOIANS,
		MessageTagBBBCVTSPIRDOI,
		MessageTagBBBCVTSPIRDOIANS,
		MessageTagBBBCVCSCSPS,
		MessageTagBBBCVCSCSPSANS,
		MessageTagBBBCVIMEDOI,
		MessageTagBBBCVIMEDOIANS,
		MessageTagBBBCVEAASDCBI,
		MessageTagBBBCVEAASDCBIANS,
		MessageTagBBBCVEAANSDCBI,
		MessageTagBBBCVEAANSDCBIANS,
		MessageTagBBBCVERATSRI,
		MessageTagBBBCVERATSRIANS,
		MessageTagBBBCVERATSSRI,
		MessageTagBBBCVERATSSRIANS,
		MessageTagBBBCVISMEC,
		MessageTagBBBCVISMECANS,
		MessageTagBBBCVISMECANS2,
		MessageTagBBBCVAAMEF,
		MessageTagBBBCVAAMEFANS,
		MessageTagBBBCVERTSTRN,
		MessageTagBBBCVERTSTRNANS1,
		MessageTagBBBCVERTSTRNANS2,
		MessageTagBBBCVDRNMND,
		MessageTagBBBCVDRNMNDANS,
		MessageTagBBBCVDMENMND,
		MessageTagBBBCVDMENMNDANS,
		MessageTagBBBCVISI,
		MessageTagBBBCVISIC,
		MessageTagBBBCVISIR,
		MessageTagBBBCVISIT,
		MessageTagBBBCVISIANS,
		MessageTagBBBCVIAFS,
		MessageTagBBBCVIAFSANS,
		MessageTagBBBICSISCI,
		MessageTagBBBICSISCIANS,
		MessageTagBBBICSISASCP,
		MessageTagBBBICSISASCPANS,
		MessageTagBBBICSISASCPU,
		MessageTagBBBICSISASCPUANS,
		MessageTagBBBICSISACDP,
		MessageTagBBBICSISACDPANS,
		MessageTagBBBICSICDVV,
		MessageTagBBBICSICDVVANS,
		MessageTagBBBICSICDVVS,
		MessageTagBBBICSICDVVSANS,
		MessageTagBBBICSAIDNASNE,
		MessageTagBBBICSAIDNASNEANS,
		MessageTagBBBICSISAKIDP,
		MessageTagBBBICSISAKIDPANS,
		MessageTagBBBICSDKIDVM,
		MessageTagBBBICSDKIDVMANS,
		MessageTagBBBICSISAX509UP,
		MessageTagBBBICSISAX509UPANS,
		MessageTagBBBICSISAX509UA,
		MessageTagBBBICSISAX509UAANS,
		MessageTagBBBRFCNUP,
		MessageTagBBBRFCNUPANS,
		MessageTagBBBRFCIRIF,
		MessageTagBBBRFCIRIFTUNU,
		MessageTagBBBRFCIRIFANS,
		MessageTagADESTROBVPIIC,
		MessageTagADESTROBVPIICANS,
		MessageTagADESTROTVPIIC,
		MessageTagADESTROTVPIICANS,
		MessageTagADESTRORPIIC,
		MessageTagADESTRORPIICANS,
		MessageTagBSVIFCRC,
		MessageTagBSVIFCRCANS,
		MessageTagBSVIISCRC,
		MessageTagBSVIISCRCANS,
		MessageTagBSVIVCIRC,
		MessageTagBSVIVCIRCANS,
		MessageTagBSVIXCVRC,
		MessageTagBSVIXCVRCANS,
		MessageTagBSVISCRAVTC,
		MessageTagBSVISCRAVTCANS,
		MessageTagBSVIVTAVRSC,
		MessageTagBSVIVTAVRSCANS,
		MessageTagBSVISCCTC,
		MessageTagBSVISCCTCANS,
		MessageTagBSVICTGTNASCRT,
		MessageTagBSVICTGTNASCRTANS,
		MessageTagBSVICTGTNASCET,
		MessageTagBSVICTGTNASCETANS,
		MessageTagBSVICVRC,
		MessageTagBSVICVRCANS,
		MessageTagBSVISAVRC,
		MessageTagBSVISAVRCANS,
		MessageTagBSVICTGTNACCET,
		MessageTagBSVICTGTNACCETANS,
		MessageTagBSVIEAAAVRC,
		MessageTagBSVIEAAAVRCANS,
		MessageTagLTVABSV,
		MessageTagLTVABSVANS,
		MessageTagLTVISCKNR,
		MessageTagLTVISCKNRANS0,
		MessageTagLTVISCKNRANS1,
		MessageTagArchLTVV,
		MessageTagArchLTVVANS,
		MessageTagArchLTAIVMP,
		MessageTagArchLTAIVMPANS,
		MessageTagArchIRTVBBA,
		MessageTagArchIRTVBBAANS,
		MessageTagArchICHFCRLPOET,
		MessageTagArchICHFCRLPOETANS,
		MessageTagACCM,
		MessageTagACCMANS,
		MessageTagASCCMCAA,
		MessageTagASCCMCAAANS,
		MessageTagASCCMDAA,
		MessageTagASCCMDAAANS,
		MessageTagASCCMDAAANS2,
		MessageTagASCCMAPKSA,
		MessageTagASCCMAPKSAANS,
		MessageTagASCCMAPKSAANS2,
		MessageTagASCCMAR,
		MessageTagASCCMARANSANR,
		MessageTagASCCMARANSANR2,
		MessageTagASCCMARANSAKSNR,
		MessageTagASCCMARANSAKSNR2,
		MessageTagASCCMPKSK,
		MessageTagASCCMPKSKANS,
		MessageTagACCMDescWithID,
		MessageTagACCMDescWithIDResult,
		MessageTagACCMDescWithName,
		MessageTagACCMPosSigSig,
		MessageTagACCMPosTSTSig,
		MessageTagACCMPosRevocSig,
		MessageTagACCMPosEVRecord,
		MessageTagACCMPosEAA,
		MessageTagACCMPosEAARev,
		MessageTagACCMPosCNTRSig,
		MessageTagACCMPosCNTRSigPL,
		MessageTagACCMPosConDig,
		MessageTagACCMPosEAAKB,
		MessageTagACCMPosEAAND,
		MessageTagACCMPosEAANSD,
		MessageTagACCMPosEAANSDPL,
		MessageTagACCMPosEAAOSDC,
		MessageTagACCMPosEAAOSDCPL,
		MessageTagACCMPosEAAPD,
		MessageTagACCMPosEAASD,
		MessageTagACCMPosEAASDPL,
		MessageTagACCMPosERADO,
		MessageTagACCMPosERADOPL,
		MessageTagACCMPosEROr,
		MessageTagACCMPosEROrPL,
		MessageTagACCMPosERTST,
		MessageTagACCMPosERTSTSeq,
		MessageTagACCMPosERMSTSig,
		MessageTagACCMPosJWS,
		MessageTagACCMPosCose,
		MessageTagACCMPosKey,
		MessageTagACCMPosKeyPL,
		MessageTagACCMPosMan,
		MessageTagACCMPosManPL,
		MessageTagACCMPosManENT,
		MessageTagACCMPosManENTPL,
		MessageTagACCMPosMesDig,
		MessageTagACCMPosMessImp,
		MessageTagACCMPosRef,
		MessageTagACCMPosRefPL,
		MessageTagACCMPosSigDENT,
		MessageTagACCMPosSigDENTPL,
		MessageTagACCMPosSigValAndPrt,
		MessageTagACCMPosSIGNDObj,
		MessageTagACCMPosSIGNDPrt,
		MessageTagACCMPosSigntrPrt,
		MessageTagACCMPosCertChainSig,
		MessageTagACCMPosCertChainTST,
		MessageTagACCMPosCertChainRevoc,
		MessageTagACCMPosCertChainEAARev,
		MessageTagACCMPosCertChain,
		MessageTagACCMPosSigCertRef,
		MessageTagBBBSAVDSCACRCC,
		MessageTagBBBSAVDSCACRCCANS,
		MessageTagBBBSAVACPCCRSCA,
		MessageTagBBBSAVACPCCRSCAANS,
		MessageTagBBBSAVISVA,
		MessageTagBBBSAVISVAANS,
		MessageTagBBBSAVISSV,
		MessageTagBBBSAVISSVANS,
		MessageTagBBBSAVICERRM,
		MessageTagBBBSAVICERRMANS,
		MessageTagBBBSAVICRM,
		MessageTagBBBSAVICRMANS,
		MessageTagBBBSAVISQPCTP,
		MessageTagBBBSAVISQPCTPANS,
		MessageTagBBBSAVISQPCHP,
		MessageTagBBBSAVISQPCHPANS,
		MessageTagBBBSAVISQPCIP,
		MessageTagBBBSAVISQPCIPANS,
		MessageTagBBBSAVISQPCTSIP,
		MessageTagBBBSAVISQPCTSIPANS,
		MessageTagBBBSAVISQPSTYPP,
		MessageTagBBBSAVISQPSTYPPANS,
		MessageTagBBBSAVISQPSLP,
		MessageTagBBBSAVISQPSLPANS,
		MessageTagBBBSAVISQPSTP,
		MessageTagBBBSAVISQPSTPANS,
		MessageTagBBBSAVISQPSTWSCVR,
		MessageTagBBBSAVISQPSTWSCVRANS,
		MessageTagBBBSAVISQPXTIP,
		MessageTagBBBSAVISQPXTIPANS,
		MessageTagBBBSAVIUQPCSP,
		MessageTagBBBSAVIUQPCSPANS,
		MessageTagBBBSAVIUQPSTSP,
		MessageTagBBBSAVIUQPSTSPANS,
		MessageTagBBBSAVIUQPVDTSP,
		MessageTagBBBSAVIUQPVDTSPANS,
		MessageTagBBBSAVIUQPVDROTSP,
		MessageTagBBBSAVIUQPVDROTSPANS,
		MessageTagBBBSAVIUQPATSP,
		MessageTagBBBSAVIUQPATSPANS,
		MessageTagBBBSAVICTVS,
		MessageTagBBBSAVICTVSANS,
		MessageTagBBBSAVIDTSP,
		MessageTagBBBSAVIDTSPANS,
		MessageTagBBBSAVITVS,
		MessageTagBBBSAVITVSANS,
		MessageTagBBBSAVIVTTSTP,
		MessageTagBBBSAVIVTTSTPANS,
		MessageTagBBBSAVIVLTATSTP,
		MessageTagBBBSAVIVLTATSTPANS,
		MessageTagBBBSAVISQPMDOSPP,
		MessageTagBBBSAVISQPMDOSPPANS,
		MessageTagBBBSAVDMICTSTMCMI,
		MessageTagBBBSAVDMICTSTMCMIANS,
		MessageTagBBBTavITSAP,
		MessageTagBBBTavITSAPANS,
		MessageTagBBBTavDTSAVM,
		MessageTagBBBTavDTSAVMANS,
		MessageTagBBBTavDTSAOM,
		MessageTagBBBTavDTSAOMANS,
		MessageTagBBBVCIISPK,
		MessageTagBBBVCIISPKANS,
		MessageTagBBBVCIISPA,
		MessageTagBBBVCIISPAANS,
		MessageTagBBBVCIISPSUPP,
		MessageTagBBBVCIISPSUPPANS,
		MessageTagBBBVCIISPM,
		MessageTagBBBVCIISPMANS,
		MessageTagBBBVCIIZHSP,
		MessageTagBBBVCIIZHSPANS,
		MessageTagBBBXCVSub,
		MessageTagBBBXCVSubANS,
		MessageTagBBBXCVSubANS2,
		MessageTagBBBXCVRFC,
		MessageTagBBBXCVRFCANS,
		MessageTagBBBXCVRAC,
		MessageTagBBBXCVRACANS,
		MessageTagBBBXCVCCCBB,
		MessageTagBBBXCVCCCBBANS,
		MessageTagBBBXCVCCCBBSigANS,
		MessageTagBBBXCVCCCBBTSPANS,
		MessageTagBBBXCVCCCBBRevANS,
		MessageTagBBBXCVCMDCIPI,
		MessageTagBBBXCVCMDCIPIANS,
		MessageTagBBBXCVCMDCIQC,
		MessageTagBBBXCVCMDCIQCANS,
		MessageTagBBBXCVCMDCIQSCD,
		MessageTagBBBXCVCMDCIQSCDANS,
		MessageTagBBBXCVCMDCIITLP,
		MessageTagBBBXCVCMDCIITLPANS,
		MessageTagBBBXCVCMDCIITNP,
		MessageTagBBBXCVCMDCIITNPANS,
		MessageTagBBBXCVCMDCICQCC,
		MessageTagBBBXCVCMDCICQCCANS,
		MessageTagBBBXCVCMDCICQCLVA,
		MessageTagBBBXCVCMDCICQCLVAANS,
		MessageTagBBBXCVCMDCICQCLVHAC,
		MessageTagBBBXCVCMDCICQCLVHACANS,
		MessageTagBBBXCVCMDCICQCERPA,
		MessageTagBBBXCVCMDCICQCERPAANS,
		MessageTagBBBXCVCMDCICSQCSSCD,
		MessageTagBBBXCVCMDCICSQCSSCDANS,
		MessageTagBBBXCVCMDCICQCPDSLA,
		MessageTagBBBXCVCMDCICQCPDSLAANS,
		MessageTagBBBXCVCMDCICQCTA,
		MessageTagBBBXCVCMDCICQCTAANS,
		MessageTagBBBXCVCMDCDCQCCLCEC,
		MessageTagBBBXCVCMDCDCQCCLCECANS,
		MessageTagBBBXCVCMDCDCQCCLCECANSEU,
		MessageTagBBBXCVCMDCSCSIA,
		MessageTagBBBXCVCMDCSCSIAANS,
		MessageTagBBBXCVCMDCICQCRA,
		MessageTagBBBXCVCMDCICQCRAANS,
		MessageTagBBBXCVCMDCICQCNA,
		MessageTagBBBXCVCMDCICQCNAANS,
		MessageTagBBBXCVCMDCICQCIA,
		MessageTagBBBXCVCMDCICQCIAANS,
		MessageTagBBBXCVCMDCDCQCQSCDLSA,
		MessageTagBBBXCVCMDCDCQCQSCDLSAANS,
		MessageTagBBBXCVCMDCDCQCIMSA,
		MessageTagBBBXCVCMDCDCQCIMSAANS,
		MessageTagBBBXCVCMDCPSBCLA,
		MessageTagBBBXCVCMDCPSBCLAANS,
		MessageTagBBBXCVCMDCPSBASIA,
		MessageTagBBBXCVCMDCPSBASIAANS,
		MessageTagBBBXCVCMDCPSBLIA,
		MessageTagBBBXCVCMDCPSBLIAANS,
		MessageTagBBBXCVDCCUCE,
		MessageTagBBBXCVDCCUCEANS,
		MessageTagBBBXCVDCCFCE,
		MessageTagBBBXCVDCCFCEANS,
		MessageTagBBBXCVDCSBSINC,
		MessageTagBBBXCVDCSBSINCANS,
		MessageTagBBBXCVICAC,
		MessageTagBBBXCVICACANS,
		MessageTagBBBXCVICPDV,
		MessageTagBBBXCVICPDVANS,
		MessageTagBBBXCVICPTV,
		MessageTagBBBXCVICPTVANS,
		MessageTagBBBXCVIAKIP,
		MessageTagBBBXCVIAKIPANS,
		MessageTagBBBXCVISKIP,
		MessageTagBBBXCVISKIPANS,
		MessageTagBBBXCVICNRAEV,
		MessageTagBBBXCVICNRAEVANS,
		MessageTagBBBXCVIVTBCTSD,
		MessageTagBBBXCVIVTBCTSDANS,
		MessageTagBBBXCVICTIVRSC,
		MessageTagBBBXCVICTIVRSCANS,
		MessageTagBBBXCVICTIVRCIRI,
		MessageTagBBBXCVICTIVRCIRIANS,
		MessageTagBBBXCVIRDCSFC,
		MessageTagBBBXCVIRDCSFCANS,
		MessageTagBBBXCVIRDPFC,
		MessageTagBBBXCVIRDPFCANS,
		MessageTagBBBXCVIRDPFRC,
		MessageTagBBBXCVIRDPFRCANS,
		MessageTagBBBXCVIARDPFC,
		MessageTagBBBXCVIARDPFCANS,
		MessageTagBBBVTSIRDPFC,
		MessageTagBBBVTSIRDPFCANS,
		MessageTagBBBXCVISCOH,
		MessageTagBBBXCVISCOHANS,
		MessageTagBBBXCVISCUKN,
		MessageTagBBBXCVISCUKNANS,
		MessageTagBBBXCVISCR,
		MessageTagBBBXCVISCRANS,
		MessageTagBBBXCVISCGKU,
		MessageTagBBBXCVISCGKUANS,
		MessageTagBBBXCVISCGKUANSCert,
		MessageTagBBBXCVISCGEKU,
		MessageTagBBBXCVISCGEKUANS,
		MessageTagBBBXCVISCGEKUANSCert,
		MessageTagBBBXCVICSI,
		MessageTagBBBXCVICSIANS,
		MessageTagBBBXCVIOTAA,
		MessageTagBBBXCVIOTAAANS,
		MessageTagBBBXCVHPCCVVT,
		MessageTagBBBXCVHPCCVVTANS,
		MessageTagBBBXCVPseudoUse,
		MessageTagBBBXCVPseudoUseANS,
		MessageTagBBBXCVAIAPres,
		MessageTagBBBXCVAIAPresANS,
		MessageTagBBBXCVRevocPres,
		MessageTagBBBXCVRevocPresANS,
		MessageTagBBBXCVRevocThisUpdatePresent,
		MessageTagBBBXCVRevocThisUpdatePresentANS,
		MessageTagBBBXCVRevocIssuerKnown,
		MessageTagBBBXCVRevocIssuerKnownANS,
		MessageTagBBBXCVRevocIssuerValidAtProd,
		MessageTagBBBXCVRevocIssuerValidAtProdANS,
		MessageTagBBBXCVRevocAfterCertNotBefore,
		MessageTagBBBXCVRevocAfterCertNotBeforeANS,
		MessageTagBBBXCVRevocHasCertInfo,
		MessageTagBBBXCVRevocHasCertInfoANS,
		MessageTagBBBXCVRevocRespIDMatch,
		MessageTagBBBXCVRevocRespIDMatchANS,
		MessageTagBBBXCVRevocCertHashPresent,
		MessageTagBBBXCVRevocCertHashPresentANS,
		MessageTagBBBXCVRevocCertHashMatch,
		MessageTagBBBXCVRevocCertHashMatchANS,
		MessageTagBBBXCVRevocSelfIssuedOCSP,
		MessageTagBBBXCVRevocSelfIssuedOCSPANS,
		MessageTagBBBXCVDCIDNMSDNIC,
		MessageTagBBBXCVDCIDNMSDNICANS,
		MessageTagBBBXCVISCGCOUN,
		MessageTagBBBXCVISCGCOUNANS,
		MessageTagBBBXCVISCGLOC,
		MessageTagBBBXCVISCGLOCANS,
		MessageTagBBBXCVISCGST,
		MessageTagBBBXCVISCGSTANS,
		MessageTagBBBXCVISCGORGAN,
		MessageTagBBBXCVISCGORGANANS,
		MessageTagBBBXCVISCGORGAU,
		MessageTagBBBXCVISCGORGAUANS,
		MessageTagBBBXCVISCGORGAI,
		MessageTagBBBXCVISCGORGAIANS,
		MessageTagBBBXCVISCGSURN,
		MessageTagBBBXCVISCGSURNANS,
		MessageTagBBBXCVISCGGIVEN,
		MessageTagBBBXCVISCGGIVENANS,
		MessageTagBBBXCVISCGPSEUDO,
		MessageTagBBBXCVISCGPSEUDOANS,
		MessageTagBBBXCVISCGCOMMONN,
		MessageTagBBBXCVISCGCOMMONNANS,
		MessageTagBBBXCVISCGTITLE,
		MessageTagBBBXCVISCGTITLEANS,
		MessageTagBBBXCVISCGEMAIL,
		MessageTagBBBXCVISCGEMAILANS,
		MessageTagBBBXCVISSSC,
		MessageTagBBBXCVISSSCANS,
		MessageTagBBBXCVISNSSC,
		MessageTagBBBXCVISNSSCANS,
		MessageTagBBBXCVIRDC,
		MessageTagBBBXCVIRDCANS,
		MessageTagXCVTSLESP,
		MessageTagXCVTSLESPANS,
		MessageTagXCVTSLESPSigANS,
		MessageTagXCVTSLESPTSPANS,
		MessageTagXCVTSLESPRevANS,
		MessageTagXCVTSLETIP,
		MessageTagXCVTSLETIPANS,
		MessageTagXCVTSLETIPSigANS,
		MessageTagXCVTSLETIPTSPANS,
		MessageTagXCVTSLETIPRevANS,
		MessageTagPCVIVTSC,
		MessageTagPCVIVTSCANS,
		MessageTagPCVICCSVTSF,
		MessageTagPCVICCSVTSFANS,
		MessageTagPSVIPCVA,
		MessageTagPSVIPCVAANS,
		MessageTagPSVIPCVC,
		MessageTagPSVIPCVCANS,
		MessageTagPSVIPSVC,
		MessageTagPSVIPSVCANS,
		MessageTagPSVIPTVC,
		MessageTagPSVIPTVCANS,
		MessageTagPSVITPOCOBCT,
		MessageTagPSVITPOSVAOBCT,
		MessageTagPSVITPOSVAOBCTANS,
		MessageTagPSVITPORDAOBCT,
		MessageTagPSVITPOOBCTANS,
		MessageTagPSVITPRISCNARTCAC,
		MessageTagPSVITPRISCNARTCACANS,
		MessageTagPSVICRDIT,
		MessageTagPSVICRDITANS,
		MessageTagPSVIPCRIAIDBEDC,
		MessageTagPSVIPCRIAIDBEDCANS,
		MessageTagPSVICTD,
		MessageTagPSVICTDANS,
		MessageTagPSVISDDTA,
		MessageTagPSVISDDTAANS,
		MessageTagPSVHRDBIBCT,
		MessageTagPSVHRDBIBCTANS,
		MessageTagPSVDIURDSCHPVR,
		MessageTagPSVDIURDSCHPVRANS,
		MessageTagTSVASTPTCT,
		MessageTagTSVASTPTCTANS,
		MessageTagTSVIBSTAIDOSC,
		MessageTagTSVIBSTAIDOSCANS,
		MessageTagTSVIBSTBCEC,
		MessageTagTSVIBSTBCECANS,
		MessageTagTSVISCNVABST,
		MessageTagTSVISCNVABSTANS,
		MessageTagADESTIRTPTBST,
		MessageTagADESTIRTPTBSTANS,
		MessageTagADESTISTPTBST,
		MessageTagADESTISTPTBSTANS,
		MessageTagADESTVFDTAOCSTANS,
		MessageTagADESTISTPTDABST,
		MessageTagADESTISTPTDABSTANS,
		MessageTagADESTIBSVPSC,
		MessageTagADESTIBSVPSCANS,
		MessageTagADESTIBSVPTC,
		MessageTagADESTIBSVPTCANS,
		MessageTagADESTIBSVPTADC,
		MessageTagADESTIBSVPTADCANS,
		MessageTagADESTIRERVPC,
		MessageTagADESTIRERVPCANS,
		MessageTagEAACertLoTEReached,
		MessageTagEAACertLoTEReachedANS,
		MessageTagEAACertTrustAnchorListReached,
		MessageTagEAACertTrustAnchorListReachedANS,
		MessageTagEAADPEAAP,
		MessageTagEAADPEAAPANS,
		MessageTagEAADLEEAAP,
		MessageTagEAADLEEAAPANS,
		MessageTagEAAKBRC,
		MessageTagEAAKBRCANS,
		MessageTagEAAKBSP,
		MessageTagEAAKBSPANS,
		MessageTagEAAClaims,
		MessageTagEAAClaimsANS,
		MessageTagEAAClaimsInfo,
		MessageTagEAASupportedClaims,
		MessageTagEAASupportedClaimsANS,
		MessageTagEAAUnsupportedClaims,
		MessageTagEAAAcceptableType,
		MessageTagEAAAcceptableTypeANS,
		MessageTagEAAIdentifierPresent,
		MessageTagEAAIdentifierPresentANS,
		MessageTagEAAIssuanceDatePresent,
		MessageTagEAAIssuanceDatePresentANS,
		MessageTagEAANBFPresent,
		MessageTagEAANBFPresentANS,
		MessageTagEAAExpPresent,
		MessageTagEAAExpPresentANS,
		MessageTagEAAAIDPresent,
		MessageTagEAAAIDPresentANS,
		MessageTagEAAAEDPresent,
		MessageTagEAAAEDPresentANS,
		MessageTagEAASigPresent,
		MessageTagEAASigPresentANS,
		MessageTagEAASigQual,
		MessageTagEAASigQualANS,
		MessageTagEAACATEAA,
		MessageTagEAACATEAAANS1,
		MessageTagEAACATEAAANS2,
		MessageTagEAACATPubEAA,
		MessageTagEAACATPubEAAANS,
		MessageTagEAACATQEAA,
		MessageTagEAACATQEAAANS,
		MessageTagEAAQCPSB,
		MessageTagEAAQCPSBANS,
		MessageTagEAAQualConclusive,
		MessageTagEAAQualConclusiveANS,
		MessageTagEAAETSI194721,
		MessageTagEAAETSI194721ANS,
		MessageTagEAAVTITVR,
		MessageTagEAAVTITVRANS,
		MessageTagEAAVTITVRValidity,
		MessageTagEAANowBeforeNBF,
		MessageTagEAANowAfterExp,
		MessageTagEAAVTIAVR,
		MessageTagEAAVTIAVRANS,
		MessageTagEAAVTIAVRValidity,
		MessageTagEAANowBeforeADI,
		MessageTagEAANowAfterADE,
		MessageTagEAAADSDJWTConformance,
		MessageTagEAAShortLivedStatusPresent,
		MessageTagEAAMandatoryStatusAbsent,
		MessageTagEAARevSDJWTConformance,
		MessageTagEAAMDocIssuingAuthority,
		MessageTagEAASDJWTIssuingAuthority,
		MessageTagEAAMDocDocumentNumberAbsent,
		MessageTagEAASub,
		MessageTagEAASubANS,
		MessageTagEAASubPSE,
		MessageTagEAASubPSEANS,
		MessageTagEAACAT,
		MessageTagEAACATANS,
		MessageTagEAAISSCOUN,
		MessageTagEAAISSCOUNANS,
		MessageTagEAAISSAuth,
		MessageTagEAAISSAuthANS,
		MessageTagEAAISSRegID,
		MessageTagEAAISSRegIDANS,
		MessageTagEAARevPR,
		MessageTagEAARevPRANS,
		MessageTagEAARevAV,
		MessageTagEAARevAVANS,
		MessageTagEAARevACC,
		MessageTagEAARevACCANS,
		MessageTagEAARevACCFND,
		MessageTagEAARevACCFNDANS,
		MessageTagEAARevNotRev,
		MessageTagEAARevNotRevANS,
		MessageTagEAARevNotOnHold,
		MessageTagEAARevNotOnHoldANS,
		MessageTagEAASHLVD,
		MessageTagEAASHLVDANS,
		MessageTagEAAOTU,
		MessageTagEAAOTUANS,
		MessageTagEAAPseudoUsed,
		MessageTagEAAPseudoUsedANS,
		MessageTagSDJWTEAAVCTPresent,
		MessageTagSDJWTEAAVCTPresentANS,
		MessageTagSDJWTEAAVCTIntPresent,
		MessageTagSDJWTEAAVCTIntPresentANS,
		MessageTagEAARevType,
		MessageTagEAARevTypeANS,
		MessageTagEAARevKnown,
		MessageTagEAARevKnownANS,
		MessageTagEAARevISS,
		MessageTagEAARevISSANS,
		MessageTagEAARevExp,
		MessageTagEAARevExpANS,
		MessageTagEAARevNotExp,
		MessageTagEAARevNotExpANS,
		MessageTagEAARevSub,
		MessageTagEAARevSubANS,
		MessageTagEAARevSubMatch,
		MessageTagEAARevSubMatchANS,
		MessageTagEAARevISSValid,
		MessageTagEAARevISSValidANS,
		MessageTagEAARevTime,
		MessageTagEAARevISSCert,
		MessageTagQualTLExp,
		MessageTagQualTLExpANS,
		MessageTagQualTLFresh,
		MessageTagQualTLFreshANS,
		MessageTagQualTLVersion,
		MessageTagQualTLVersionANS,
		MessageTagQualTLWS,
		MessageTagQualTLWSANS,
		MessageTagQualTLSV,
		MessageTagQualTLSVANS,
		MessageTagQualTLIMRA,
		MessageTagQualTLIMRAANS,
		MessageTagQualTLIMRAANSV1,
		MessageTagQualTLIMRAANSV2,
		MessageTagQualTLServCONS,
		MessageTagQualTLServCONSANS0,
		MessageTagQualTLServCONSANS1,
		MessageTagQualTLServCONSANS2,
		MessageTagQualTLServCONSANS3,
		MessageTagQualTLServCONSANS3A,
		MessageTagQualTLServCONSANS3B,
		MessageTagQualTLServCONSANS3C,
		MessageTagQualTLServCONSANS4,
		MessageTagQualTLServCONSANS5,
		MessageTagQualTLServCONSANS6,
		MessageTagQualTLServCONSANS7,
		MessageTagQualCertTrustedListReached,
		MessageTagQualCertTrustedListReachedANS,
		MessageTagQualTrustedListAccept,
		MessageTagQualTrustedListAcceptANS,
		MessageTagQualListOfTrustedListsAccept,
		MessageTagQualListOfTrustedListsAcceptANS,
		MessageTagQualValidTrustedListPresent,
		MessageTagQualValidTrustedListPresentANS,
		MessageTagQualCertTypeAtST,
		MessageTagQualCertTypeAtSTANS,
		MessageTagQualCertTypeAtCC,
		MessageTagQualCertTypeAtCCANS,
		MessageTagQualCertTypeAtVT,
		MessageTagQualCertTypeAtVTANS,
		MessageTagQualQCAtST,
		MessageTagQualQCAtSTANS,
		MessageTagQualQCAtCC,
		MessageTagQualQCAtCCANS,
		MessageTagQualQCAtVT,
		MessageTagQualQCAtVTANS,
		MessageTagQualQSCDAtST,
		MessageTagQualQSCDAtSTANS,
		MessageTagQualQSCDAtCC,
		MessageTagQualQSCDAtCCANS,
		MessageTagQualQSCDAtVT,
		MessageTagQualQSCDAtVTANS,
		MessageTagQualUniqueCert,
		MessageTagQualUniqueCertANS,
		MessageTagQualIsAdES,
		MessageTagQualIsAdESInd,
		MessageTagQualIsAdESINV,
		MessageTagQualHasMETS,
		MessageTagQualHasMETSANS,
		MessageTagQualHasMETSAtTime,
		MessageTagQualHasMETSAtTimeANS,
		MessageTagQualHasMETSHCCECBA,
		MessageTagQualHasMETSHCCECBAANS,
		MessageTagQualHasMETSHCCECBAANS2,
		MessageTagQualHasMETSHCCECBAANS3,
		MessageTagQualHasCAQC,
		MessageTagQualHasCAQCANS,
		MessageTagQualHasCAQCANS2,
		MessageTagQualHasAtTime,
		MessageTagQualHasAtTimeANS,
		MessageTagQualHasTSCertType,
		MessageTagQualHasTSCertTypeANS,
		MessageTagQualHasConf,
		MessageTagQualHasConfANS,
		MessageTagQualHasQEAA,
		MessageTagQualHasQEAAANS,
		MessageTagQualHasQTST,
		MessageTagQualHasQTSTANS,
		MessageTagQualIsTrustCertMatchService,
		MessageTagQualIsTrustCertMatchServiceANS0,
		MessageTagQualIsTrustCertMatchServiceANS1,
		MessageTagQualIsTrustCertMatchServiceANS2,
		MessageTagQualHasGranted,
		MessageTagQualHasGrantedANS,
		MessageTagQualHasGrantedANS2,
		MessageTagQualHasGrantedAt,
		MessageTagQualHasGrantedAtANS,
		MessageTagQualHasConsistentByQC,
		MessageTagQualHasConsistentByQCANS,
		MessageTagQualHasConsistentByQSCD,
		MessageTagQualHasConsistentByQSCDANS,
		MessageTagQualHasCertTypeCoverage,
		MessageTagQualHasCertTypeCoverageANS,
		MessageTagQualHasValidCAQC,
		MessageTagQualHasValidCAQCANS,
		MessageTagQualHasOnlyOne,
		MessageTagQualHasOnlyOneANS,
		MessageTagQWACValid,
		MessageTagQWACValidANS,
		MessageTagQWACValidANS2,
		MessageTagQWACCertQualConclusive,
		MessageTagQWACCertQualConclusiveANS,
		MessageTagQWACIsWSAAtTime,
		MessageTagQWACIsWSAAtTimeANS,
		MessageTagQWACCertPolicy,
		MessageTagQWACCertPolicyANS,
		MessageTagQWACValPeriod,
		MessageTagQWACValPeriodANS,
		MessageTagQWACDomainName,
		MessageTagQWACDomainNameANS,
		MessageTagQWAC2ExtKeyUsage,
		MessageTagQWAC2ExtKeyUsageANS,
		MessageTagTLSCertBindingURL,
		MessageTagTLSCertBindingURLANS,
		MessageTagTLSCertBindingSig,
		MessageTagTLSCertBindingSigANS,
		MessageTagTLSCertBindingSigForm,
		MessageTagTLSCertBindingSigFormANS,
		MessageTagTLSCertBindingSigSer,
		MessageTagTLSCertBindingSigSerANS,
		MessageTagTLSCertBindingSigExp,
		MessageTagTLSCertBindingSigExpANS,
		MessageTagTLSCertBindingSigExpiryDate,
		MessageTagTLSCertBindingSigExpiryDateANS,
		MessageTagTLSCertBindingQWAC2,
		MessageTagTLSCertBindingQWAC2ANS,
		MessageTagTLSCertBindingSigValid,
		MessageTagTLSCertBindingSigValidANS,
		MessageTagTLSCertBindingCertIdentified,
		MessageTagTLSCertBindingCertIdentifiedANS,
		MessageTagCertUsageLoTEAccept,
		MessageTagCertUsageLoTEAcceptANS,
		MessageTagCertUsageLoLoTEAccept,
		MessageTagCertUsageLoLoTEAcceptANS,
		MessageTagCertUsageValidLoTEPresent,
		MessageTagCertUsageValidLoTEPresentANS,
		MessageTagCertUsageHasAtTime,
		MessageTagCertUsageHasAtTimeANS,
		MessageTagCertUsageListTypeKnown,
		MessageTagCertUsageListTypeKnownANS,
		MessageTagCertUsageStatus,
		MessageTagCertUsageStatusANS,
		MessageTagCertUsageStatusCONS,
		MessageTagCertUsageStatusCONSANS,
		MessageTagCertUsageStatusKnown,
		MessageTagCertUsageStatusKnownANS,
		MessageTagCertUsageSti,
		MessageTagCertUsageStiANS,
		MessageTagCertUsageStiKnown,
		MessageTagCertUsageStiKnownANS,
		MessageTagPIDDocumentType,
		MessageTagPIDDocumentTypeANS,
		MessageTagPIDLoTETypePIDProviders,
		MessageTagPIDLoTETypePIDProvidersANS,
		MessageTagPIDStiPIDIssuance,
		MessageTagPIDStiPIDIssuanceANS,
		MessageTagPIDProviderAtIssuanceTime,
		MessageTagPIDProviderAtIssuanceTimeANS,
		MessageTagPIDProviderAtValidationTime,
		MessageTagPIDProviderAtValidationTimeANS,
		MessageTagBBBAccept,
		MessageTagBBBAcceptANS,
		MessageTagTSTTypeContentTST,
		MessageTagTSTTypeSignatureTST,
		MessageTagTSTTypeVDTST,
		MessageTagTSTTypeDocTST,
		MessageTagTSTTypeContainerTST,
		MessageTagTSTTypeArchiveTST,
		MessageTagTSTTypeERTST,
		MessageTagTSTTypeRefERATST,
		MessageTagTSTTypeRefERATSTSeq,
		MessageTagEmpty,
		MessageTagCertificate,
		MessageTagCACertificate,
		MessageTagSigningCertificate,
		MessageTagRevocation,
		MessageTagRevocationSigCert,
		MessageTagRevocationCACert,
		MessageTagSignature,
		MessageTagTimestamp,
		MessageTagTimestampSigCert,
		MessageTagTimestampCACert,
		MessageTagEAARev,
		MessageTagEAARevSigCert,
		MessageTagEAARevCACert,
		MessageTagAcceptableRevocation,
		MessageTagBasicSignatureValidationResult,
		MessageTagBESTSignatureTimeCertNotAfter,
		MessageTagBESTSignatureTimeCertNotBefore,
		MessageTagBESTSignatureTimeCertRevocation,
		MessageTagBESTSignatureTimeCertSuspension,
		MessageTagCertificateID,
		MessageTagCertificateRevocationFound,
		MessageTagCertificateRevocationNotFound,
		MessageTagCertificateSunsetDate,
		MessageTagCertificateSunsetDateTrustAnchor,
		MessageTagCertificateSunsetDateValid,
		MessageTagCertificateType,
		MessageTagCertificateValidity,
		MessageTagCertificateUsage,
		MessageTagCertificateUsageListType,
		MessageTagCertificateUsageStatus,
		MessageTagCertificateUsageStatuses,
		MessageTagCertificateUsageSti,
		MessageTagCertificateUsageSTIS,
		MessageTagControlTime,
		MessageTagControlTimeAlone,
		MessageTagControlTimeWithPOE,
		MessageTagControlTimeWithTrustAnchor,
		MessageTagCryptographicCheckFailure,
		MessageTagCryptographicCheckFailureWithID,
		MessageTagCryptographicCheckFailureWithRef,
		MessageTagCryptographicCheckFailureWithRefWithName,
		MessageTagCryptographicCheckFailureWithRefWithNames,
		MessageTagCryptographicCheckSuccess,
		MessageTagCryptographicCheckSuccessKeySize,
		MessageTagCryptographicCheckSuccessDM,
		MessageTagCryptographicCheckSuccessDMWithID,
		MessageTagCryptographicCheckSuccessDMWithName,
		MessageTagCryptographicCheckSuccessDMWithNames,
		MessageTagEvidenceRecordValidation,
		MessageTagExtendedKeyUsage,
		MessageTagKeyUsage,
		MessageTagLastAcceptableRevocation,
		MessageTagListOfTrustedEntities,
		MessageTagPseudo,
		MessageTagQWACExpiryExp,
		MessageTagQWACExpiryTLSCert,
		MessageTagQWACExpirySignCert,
		MessageTagReference,
		MessageTagReferenceNameCheck,
		MessageTagReferencesWithNames,
		MessageTagRevocationAcceptanceCheck,
		MessageTagRevocationCertHashOK,
		MessageTagRevocationCertHashOKID,
		MessageTagRevocationCertValidity,
		MessageTagRevocationCheck,
		MessageTagRevocationConsistent,
		MessageTagRevocationConsistentCRL,
		MessageTagRevocationConsistentOCSP,
		MessageTagRevocationConsistentTL,
		MessageTagRevocationInfo,
		MessageTagRevocationNotAfterAfter,
		MessageTagRevocationNotAfterAfterID,
		MessageTagRevocationProducedAtCertValidity,
		MessageTagRevocationProducedAtOutOfBounds,
		MessageTagRevocationProducedAtOutOfBoundsID,
		MessageTagRevocationReason,
		MessageTagRevocationThisUpdateControlTime,
		MessageTagSignatureAlgorithmWithKeySize,
		MessageTagSignatureID,
		MessageTagStructuralValidationFailure,
		MessageTagTimestampAndCertificateNotAfter,
		MessageTagTimestampAndCryptoConstraintsExpiration,
		MessageTagTimestampAndRevocationTime,
		MessageTagTimestampValidation,
		MessageTagTokenID,
		MessageTagTrustServiceName,
		MessageTagTrustedServiceStatus,
		MessageTagTrustedServiceType,
		MessageTagTrustedList,
		MessageTagValidationTime,
		MessageTagCryptographicVerification,
		MessageTagFormatChecking,
		MessageTagIdentificationOfTheSigningCertificate,
		MessageTagPastSignatureValidation,
		MessageTagPastCertificateValidation,
		MessageTagRevocationFreshnessChecker,
		MessageTagSignatureAcceptanceValidation,
		MessageTagValidationContextInitialization,
		MessageTagValidationTimeSliding,
		MessageTagX509CertificateValidation,
		MessageTagResults,
		MessageTagEEAType,
		MessageTagEAAAcceptanceValidation,
		MessageTagAOV,
		MessageTagCertQualification,
		MessageTagCertQualificationAtTime,
		MessageTagCertUsageAtTime,
		MessageTagCertUsages,
		MessageTagCC,
		MessageTagCRS,
		MessageTagDAAV,
		MessageTagEAAQualification,
		MessageTagEAAQualificationProcess,
		MessageTagLoTE,
		MessageTagLoLoTE,
		MessageTagLOTL,
		MessageTagPIDQualificationProcess,
		MessageTagPSVCRS,
		MessageTagQWACValidation,
		MessageTagQWACValidationProfile,
		MessageTagRAC,
		MessageTagSigQualification,
		MessageTagSubXCV,
		MessageTagTL,
		MessageTagTSTQualification,
		MessageTagTSTQualificationAtTime,
		MessageTagVPBS,
		MessageTagVPEAA,
		MessageTagVPER,
		MessageTagVPFLTVD,
		MessageTagVPFRVC,
		MessageTagVpfswatsp,
		MessageTagVpftsp,
		MessageTagVpftspwatsp,
		MessageTagVTSCRS,
		MessageTagVTBESTSignatureTime,
		MessageTagVTCertificateIssuanceTime,
		MessageTagVTValidationTime,
		MessageTagVTTSTGenerationTime,
		MessageTagVTTSTPOETime,
		MessageTagQWAC1Profile,
		MessageTagQWAC2Profile,
		MessageTagTLSByQWAC2Profile,
		MessageTagSemanticsTotalPassed,
		MessageTagSemanticsPassed,
		MessageTagSemanticsTotalFailed,
		MessageTagSemanticsFailed,
		MessageTagSemanticsIndeterminate,
		MessageTagSemanticsNoSignatureFound,
		MessageTagSemanticsFormatFailure,
		MessageTagSemanticsHashFailure,
		MessageTagSemanticsSigCryptoFailure,
		MessageTagSemanticsRevoked,
		MessageTagSemanticsExpired,
		MessageTagSemanticsNotYetValid,
		MessageTagSemanticsSigConstraintsFailure,
		MessageTagSemanticsChainConstraintsFailure,
		MessageTagSemanticsCertificateChainGeneralFailure,
		MessageTagSemanticsCryptoConstraintsFailure,
		MessageTagSemanticsPolicyProcessingError,
		MessageTagSemanticsSignaturePolicyNotAvailable,
		MessageTagSemanticsTimestampOrderFailure,
		MessageTagSemanticsNoSigningCertificateFound,
		MessageTagSemanticsNoCertificateChainFound,
		MessageTagSemanticsNoCertificateChainFoundNoPOE,
		MessageTagSemanticsRevokedNoPOE,
		MessageTagSemanticsRevokedCANoPOE,
		MessageTagSemanticsOutOfBoundsNotRevoked,
		MessageTagSemanticsOutOfBoundsNoPOE,
		MessageTagSemanticsRevocationOutOfBoundsNoPOE,
		MessageTagSemanticsCryptoConstraintsFailureNoPOE,
		MessageTagSemanticsNoPOE,
		MessageTagSemanticsTryLater,
		MessageTagSemanticsSignedDataNotFound,
		MessageTagSemanticsEAAConstraintsFailure,
	}
}

// messageTagsByName indexes every MessageTag by its Java enum name, so MessageTagValueOf and
// MessageTagGetSemantic (called by the report builders once per indication/sub-indication) are a
// map lookup instead of building the ~1000-entry MessageTagValues slice and scanning it on every
// call. Built once from MessageTagValues in declaration order; a duplicate name (none exists)
// would keep its first entry, as the linear scan it replaces did.
var messageTagsByName = func() map[string]MessageTag {
	values := MessageTagValues()
	byName := make(map[string]MessageTag, len(values))
	for _, v := range values {
		if _, exists := byName[string(v)]; !exists {
			byName[string(v)] = v
		}
	}
	return byName
}()

// MessageTagValueOf returns the MessageTag matching the given Java enum name.
func MessageTagValueOf(name string) (MessageTag, error) {
	if v, ok := messageTagsByName[name]; ok {
		return v, nil
	}
	return "", fmt.Errorf("no enum constant MessageTag.%s", name)
}

// Id returns the id code of the referred message. Ports MessageTag#getId.
func (m MessageTag) Id() string {
	return string(m)
}

// MessageTagGetSemantic returns a semantics MessageTag to be used by ETSI
// Indication/SubIndication. Ports the static MessageTag#getSemantic method;
// Java's null return ("no matching tag") becomes the ok=false result.
func MessageTagGetSemantic(etsiCode string) (MessageTag, bool) {
	expectedEnumValue := "SEMANTICS_" + etsiCode
	messageTag, ok := messageTagsByName[expectedEnumValue]
	return messageTag, ok
}
