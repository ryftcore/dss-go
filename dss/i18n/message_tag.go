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
// (embedded via I18nProvider) for the exact human-readable text each
// resolves to.
const (
	MessageTag_BBB_FC_IEFF                                     MessageTag = "BBB_FC_IEFF"
	MessageTag_BBB_FC_IEFF_ANS                                 MessageTag = "BBB_FC_IEFF_ANS"
	MessageTag_BBB_FC_ICFD                                     MessageTag = "BBB_FC_ICFD"
	MessageTag_BBB_FC_ICFD_ANS                                 MessageTag = "BBB_FC_ICFD_ANS"
	MessageTag_BBB_FC_ISD                                      MessageTag = "BBB_FC_ISD"
	MessageTag_BBB_FC_ISD_ANS                                  MessageTag = "BBB_FC_ISD_ANS"
	MessageTag_BBB_FC_ISRIA                                    MessageTag = "BBB_FC_ISRIA"
	MessageTag_BBB_FC_ISRIA_ANS                                MessageTag = "BBB_FC_ISRIA_ANS"
	MessageTag_BBB_FC_IOSIP                                    MessageTag = "BBB_FC_IOSIP"
	MessageTag_BBB_FC_IOSIP_ANS                                MessageTag = "BBB_FC_IOSIP_ANS"
	MessageTag_BBB_FC_DASTHVBR                                 MessageTag = "BBB_FC_DASTHVBR"
	MessageTag_BBB_FC_DASTHVBR_ANS                             MessageTag = "BBB_FC_DASTHVBR_ANS"
	MessageTag_BBB_FC_DBTOOST                                  MessageTag = "BBB_FC_DBTOOST"
	MessageTag_BBB_FC_DBTOOST_ANS                              MessageTag = "BBB_FC_DBTOOST_ANS"
	MessageTag_BBB_FC_IBRV                                     MessageTag = "BBB_FC_IBRV"
	MessageTag_BBB_FC_IBRV_ANS                                 MessageTag = "BBB_FC_IBRV_ANS"
	MessageTag_BBB_FC_ISDC                                     MessageTag = "BBB_FC_ISDC"
	MessageTag_BBB_FC_ISDC_ANS                                 MessageTag = "BBB_FC_ISDC_ANS"
	MessageTag_BBB_FC_DSFREAP                                  MessageTag = "BBB_FC_DSFREAP"
	MessageTag_BBB_FC_DSFREAP_ANS                              MessageTag = "BBB_FC_DSFREAP_ANS"
	MessageTag_BBB_FC_IAOD                                     MessageTag = "BBB_FC_IAOD"
	MessageTag_BBB_FC_IAOD_ANS                                 MessageTag = "BBB_FC_IAOD_ANS"
	MessageTag_BBB_FC_IVDBSFR                                  MessageTag = "BBB_FC_IVDBSFR"
	MessageTag_BBB_FC_IVDBSFR_ANS                              MessageTag = "BBB_FC_IVDBSFR_ANS"
	MessageTag_BBB_FC_ISVADMDPD                                MessageTag = "BBB_FC_ISVADMDPD"
	MessageTag_BBB_FC_ISVADMDPD_ANS                            MessageTag = "BBB_FC_ISVADMDPD_ANS"
	MessageTag_BBB_FC_ISVAFMDPD                                MessageTag = "BBB_FC_ISVAFMDPD"
	MessageTag_BBB_FC_ISVAFMDPD_ANS                            MessageTag = "BBB_FC_ISVAFMDPD_ANS"
	MessageTag_BBB_FC_ISVASFLD                                 MessageTag = "BBB_FC_ISVASFLD"
	MessageTag_BBB_FC_ISVASFLD_ANS                             MessageTag = "BBB_FC_ISVASFLD_ANS"
	MessageTag_BBB_FC_DSCNFFSM                                 MessageTag = "BBB_FC_DSCNFFSM"
	MessageTag_BBB_FC_DSCNFFSM_ANS                             MessageTag = "BBB_FC_DSCNFFSM_ANS"
	MessageTag_BBB_FC_DSCNACMDM                                MessageTag = "BBB_FC_DSCNACMDM"
	MessageTag_BBB_FC_DSCNACMDM_ANS                            MessageTag = "BBB_FC_DSCNACMDM_ANS"
	MessageTag_BBB_FC_DSCNUOM                                  MessageTag = "BBB_FC_DSCNUOM"
	MessageTag_BBB_FC_DSCNUOM_ANS                              MessageTag = "BBB_FC_DSCNUOM_ANS"
	MessageTag_BBB_FC_IECKSCDA                                 MessageTag = "BBB_FC_IECKSCDA"
	MessageTag_BBB_FC_IECKSCDA_ANS1                            MessageTag = "BBB_FC_IECKSCDA_ANS1"
	MessageTag_BBB_FC_IECKSCDA_ANS2                            MessageTag = "BBB_FC_IECKSCDA_ANS2"
	MessageTag_BBB_FC_IECKSCDA_ANS3                            MessageTag = "BBB_FC_IECKSCDA_ANS3"
	MessageTag_BBB_FC_IECKSCDA_ANS4                            MessageTag = "BBB_FC_IECKSCDA_ANS4"
	MessageTag_BBB_FC_IECKSCDA_ANS5                            MessageTag = "BBB_FC_IECKSCDA_ANS5"
	MessageTag_BBB_FC_DDAPDFAF                                 MessageTag = "BBB_FC_DDAPDFAF"
	MessageTag_BBB_FC_DDAPDFAF_ANS                             MessageTag = "BBB_FC_DDAPDFAF_ANS"
	MessageTag_BBB_FC_IDPDFAC                                  MessageTag = "BBB_FC_IDPDFAC"
	MessageTag_BBB_FC_IDPDFAC_ANS                              MessageTag = "BBB_FC_IDPDFAC_ANS"
	MessageTag_BBB_FC_IECTF                                    MessageTag = "BBB_FC_IECTF"
	MessageTag_BBB_FC_IECTF_ANS                                MessageTag = "BBB_FC_IECTF_ANS"
	MessageTag_BBB_FC_ISFCS                                    MessageTag = "BBB_FC_ISFCS"
	MessageTag_BBB_FC_ISFCS_ANS                                MessageTag = "BBB_FC_ISFCS_ANS"
	MessageTag_BBB_FC_ITFCS                                    MessageTag = "BBB_FC_ITFCS"
	MessageTag_BBB_FC_ITFCS_ANS                                MessageTag = "BBB_FC_ITFCS_ANS"
	MessageTag_BBB_FC_IMFCS                                    MessageTag = "BBB_FC_IMFCS"
	MessageTag_BBB_FC_IMFCS_ANS                                MessageTag = "BBB_FC_IMFCS_ANS"
	MessageTag_BBB_FC_ITZCP                                    MessageTag = "BBB_FC_ITZCP"
	MessageTag_BBB_FC_ITZCP_ANS                                MessageTag = "BBB_FC_ITZCP_ANS"
	MessageTag_BBB_FC_ITEZCF                                   MessageTag = "BBB_FC_ITEZCF"
	MessageTag_BBB_FC_ITEZCF_ANS                               MessageTag = "BBB_FC_ITEZCF_ANS"
	MessageTag_BBB_FC_ITMFP                                    MessageTag = "BBB_FC_ITMFP"
	MessageTag_BBB_FC_ITMFP_ANS                                MessageTag = "BBB_FC_ITMFP_ANS"
	MessageTag_BBB_FC_IEMCF                                    MessageTag = "BBB_FC_IEMCF"
	MessageTag_BBB_FC_IEMCF_ANS                                MessageTag = "BBB_FC_IEMCF_ANS"
	MessageTag_BBB_FC_IMFP_ASICE                               MessageTag = "BBB_FC_IMFP_ASICE"
	MessageTag_BBB_FC_IMFP_ASICE_ANS                           MessageTag = "BBB_FC_IMFP_ASICE_ANS"
	MessageTag_BBB_FC_ISFP_ASICE                               MessageTag = "BBB_FC_ISFP_ASICE"
	MessageTag_BBB_FC_ISFP_ASICE_ANS                           MessageTag = "BBB_FC_ISFP_ASICE_ANS"
	MessageTag_BBB_FC_ISFP_ASICS                               MessageTag = "BBB_FC_ISFP_ASICS"
	MessageTag_BBB_FC_ISFP_ASICS_ANS                           MessageTag = "BBB_FC_ISFP_ASICS_ANS"
	MessageTag_BBB_FC_ISFP_ASTFORAMC                           MessageTag = "BBB_FC_ISFP_ASTFORAMC"
	MessageTag_BBB_FC_ISFP_ASTFORAMC_ANS                       MessageTag = "BBB_FC_ISFP_ASTFORAMC_ANS"
	MessageTag_BBB_FC_IAHIV                                    MessageTag = "BBB_FC_IAHIV"
	MessageTag_BBB_FC_IAHIV_ANS                                MessageTag = "BBB_FC_IAHIV_ANS"
	MessageTag_BBB_CV_IRDOF                                    MessageTag = "BBB_CV_IRDOF"
	MessageTag_BBB_CV_IRDOF_ANS                                MessageTag = "BBB_CV_IRDOF_ANS"
	MessageTag_BBB_CV_TSP_IRDOF                                MessageTag = "BBB_CV_TSP_IRDOF"
	MessageTag_BBB_CV_TSP_IRDOF_ANS                            MessageTag = "BBB_CV_TSP_IRDOF_ANS"
	MessageTag_BBB_CV_CS_CSSVF                                 MessageTag = "BBB_CV_CS_CSSVF"
	MessageTag_BBB_CV_CS_CSSVF_ANS                             MessageTag = "BBB_CV_CS_CSSVF_ANS"
	MessageTag_BBB_CV_IMEOF                                    MessageTag = "BBB_CV_IMEOF"
	MessageTag_BBB_CV_IMEOF_ANS                                MessageTag = "BBB_CV_IMEOF_ANS"
	MessageTag_BBB_CV_EAA_SDCBF                                MessageTag = "BBB_CV_EAA_SDCBF"
	MessageTag_BBB_CV_EAA_SDCBF_ANS                            MessageTag = "BBB_CV_EAA_SDCBF_ANS"
	MessageTag_BBB_CV_EAA_NSDCBF                               MessageTag = "BBB_CV_EAA_NSDCBF"
	MessageTag_BBB_CV_EAA_NSDCBF_ANS                           MessageTag = "BBB_CV_EAA_NSDCBF_ANS"
	MessageTag_BBB_CV_ER_IODOF                                 MessageTag = "BBB_CV_ER_IODOF"
	MessageTag_BBB_CV_ER_IODOF_ANS                             MessageTag = "BBB_CV_ER_IODOF_ANS"
	MessageTag_BBB_CV_ER_HASSDOC                               MessageTag = "BBB_CV_ER_HASSDOC"
	MessageTag_BBB_CV_ER_HASSDOC_ANS                           MessageTag = "BBB_CV_ER_HASSDOC_ANS"
	MessageTag_BBB_CV_ER_DFHVLCDOG                             MessageTag = "BBB_CV_ER_DFHVLCDOG"
	MessageTag_BBB_CV_ER_DFHVLCDOG_ANS                         MessageTag = "BBB_CV_ER_DFHVLCDOG_ANS"
	MessageTag_BBB_CV_ER_ATSRF                                 MessageTag = "BBB_CV_ER_ATSRF"
	MessageTag_BBB_CV_ER_ATSRF_ANS                             MessageTag = "BBB_CV_ER_ATSRF_ANS"
	MessageTag_BBB_CV_ER_ATSSRF                                MessageTag = "BBB_CV_ER_ATSSRF"
	MessageTag_BBB_CV_ER_ATSSRF_ANS                            MessageTag = "BBB_CV_ER_ATSSRF_ANS"
	MessageTag_BBB_CV_IRDOI                                    MessageTag = "BBB_CV_IRDOI"
	MessageTag_BBB_CV_IRDOI_ANS                                MessageTag = "BBB_CV_IRDOI_ANS"
	MessageTag_BBB_CV_TSP_IRDOI                                MessageTag = "BBB_CV_TSP_IRDOI"
	MessageTag_BBB_CV_TSP_IRDOI_ANS                            MessageTag = "BBB_CV_TSP_IRDOI_ANS"
	MessageTag_BBB_CV_CS_CSPS                                  MessageTag = "BBB_CV_CS_CSPS"
	MessageTag_BBB_CV_CS_CSPS_ANS                              MessageTag = "BBB_CV_CS_CSPS_ANS"
	MessageTag_BBB_CV_IMEDOI                                   MessageTag = "BBB_CV_IMEDOI"
	MessageTag_BBB_CV_IMEDOI_ANS                               MessageTag = "BBB_CV_IMEDOI_ANS"
	MessageTag_BBB_CV_EAA_SDCBI                                MessageTag = "BBB_CV_EAA_SDCBI"
	MessageTag_BBB_CV_EAA_SDCBI_ANS                            MessageTag = "BBB_CV_EAA_SDCBI_ANS"
	MessageTag_BBB_CV_EAA_NSDCBI                               MessageTag = "BBB_CV_EAA_NSDCBI"
	MessageTag_BBB_CV_EAA_NSDCBI_ANS                           MessageTag = "BBB_CV_EAA_NSDCBI_ANS"
	MessageTag_BBB_CV_ER_ATSRI                                 MessageTag = "BBB_CV_ER_ATSRI"
	MessageTag_BBB_CV_ER_ATSRI_ANS                             MessageTag = "BBB_CV_ER_ATSRI_ANS"
	MessageTag_BBB_CV_ER_ATSSRI                                MessageTag = "BBB_CV_ER_ATSSRI"
	MessageTag_BBB_CV_ER_ATSSRI_ANS                            MessageTag = "BBB_CV_ER_ATSSRI_ANS"
	MessageTag_BBB_CV_ISMEC                                    MessageTag = "BBB_CV_ISMEC"
	MessageTag_BBB_CV_ISMEC_ANS                                MessageTag = "BBB_CV_ISMEC_ANS"
	MessageTag_BBB_CV_ISMEC_ANS_2                              MessageTag = "BBB_CV_ISMEC_ANS_2"
	MessageTag_BBB_CV_AAMEF                                    MessageTag = "BBB_CV_AAMEF"
	MessageTag_BBB_CV_AAMEF_ANS                                MessageTag = "BBB_CV_AAMEF_ANS"
	MessageTag_BBB_CV_ER_TST_RN                                MessageTag = "BBB_CV_ER_TST_RN"
	MessageTag_BBB_CV_ER_TST_RN_ANS_1                          MessageTag = "BBB_CV_ER_TST_RN_ANS_1"
	MessageTag_BBB_CV_ER_TST_RN_ANS_2                          MessageTag = "BBB_CV_ER_TST_RN_ANS_2"
	MessageTag_BBB_CV_DRNMND                                   MessageTag = "BBB_CV_DRNMND"
	MessageTag_BBB_CV_DRNMND_ANS                               MessageTag = "BBB_CV_DRNMND_ANS"
	MessageTag_BBB_CV_DMENMND                                  MessageTag = "BBB_CV_DMENMND"
	MessageTag_BBB_CV_DMENMND_ANS                              MessageTag = "BBB_CV_DMENMND_ANS"
	MessageTag_BBB_CV_ISI                                      MessageTag = "BBB_CV_ISI"
	MessageTag_BBB_CV_ISIC                                     MessageTag = "BBB_CV_ISIC"
	MessageTag_BBB_CV_ISIR                                     MessageTag = "BBB_CV_ISIR"
	MessageTag_BBB_CV_ISIT                                     MessageTag = "BBB_CV_ISIT"
	MessageTag_BBB_CV_ISI_ANS                                  MessageTag = "BBB_CV_ISI_ANS"
	MessageTag_BBB_CV_IAFS                                     MessageTag = "BBB_CV_IAFS"
	MessageTag_BBB_CV_IAFS_ANS                                 MessageTag = "BBB_CV_IAFS_ANS"
	MessageTag_BBB_ICS_ISCI                                    MessageTag = "BBB_ICS_ISCI"
	MessageTag_BBB_ICS_ISCI_ANS                                MessageTag = "BBB_ICS_ISCI_ANS"
	MessageTag_BBB_ICS_ISASCP                                  MessageTag = "BBB_ICS_ISASCP"
	MessageTag_BBB_ICS_ISASCP_ANS                              MessageTag = "BBB_ICS_ISASCP_ANS"
	MessageTag_BBB_ICS_ISASCPU                                 MessageTag = "BBB_ICS_ISASCPU"
	MessageTag_BBB_ICS_ISASCPU_ANS                             MessageTag = "BBB_ICS_ISASCPU_ANS"
	MessageTag_BBB_ICS_ISACDP                                  MessageTag = "BBB_ICS_ISACDP"
	MessageTag_BBB_ICS_ISACDP_ANS                              MessageTag = "BBB_ICS_ISACDP_ANS"
	MessageTag_BBB_ICS_ICDVV                                   MessageTag = "BBB_ICS_ICDVV"
	MessageTag_BBB_ICS_ICDVV_ANS                               MessageTag = "BBB_ICS_ICDVV_ANS"
	MessageTag_BBB_ICS_ICDVVS                                  MessageTag = "BBB_ICS_ICDVVS"
	MessageTag_BBB_ICS_ICDVVS_ANS                              MessageTag = "BBB_ICS_ICDVVS_ANS"
	MessageTag_BBB_ICS_AIDNASNE                                MessageTag = "BBB_ICS_AIDNASNE"
	MessageTag_BBB_ICS_AIDNASNE_ANS                            MessageTag = "BBB_ICS_AIDNASNE_ANS"
	MessageTag_BBB_ICS_ISAKIDP                                 MessageTag = "BBB_ICS_ISAKIDP"
	MessageTag_BBB_ICS_ISAKIDP_ANS                             MessageTag = "BBB_ICS_ISAKIDP_ANS"
	MessageTag_BBB_ICS_DKIDVM                                  MessageTag = "BBB_ICS_DKIDVM"
	MessageTag_BBB_ICS_DKIDVM_ANS                              MessageTag = "BBB_ICS_DKIDVM_ANS"
	MessageTag_BBB_ICS_ISAX509UP                               MessageTag = "BBB_ICS_ISAX509UP"
	MessageTag_BBB_ICS_ISAX509UP_ANS                           MessageTag = "BBB_ICS_ISAX509UP_ANS"
	MessageTag_BBB_ICS_ISAX509UA                               MessageTag = "BBB_ICS_ISAX509UA"
	MessageTag_BBB_ICS_ISAX509UA_ANS                           MessageTag = "BBB_ICS_ISAX509UA_ANS"
	MessageTag_BBB_RFC_NUP                                     MessageTag = "BBB_RFC_NUP"
	MessageTag_BBB_RFC_NUP_ANS                                 MessageTag = "BBB_RFC_NUP_ANS"
	MessageTag_BBB_RFC_IRIF                                    MessageTag = "BBB_RFC_IRIF"
	MessageTag_BBB_RFC_IRIF_TUNU                               MessageTag = "BBB_RFC_IRIF_TUNU"
	MessageTag_BBB_RFC_IRIF_ANS                                MessageTag = "BBB_RFC_IRIF_ANS"
	MessageTag_ADEST_ROBVPIIC                                  MessageTag = "ADEST_ROBVPIIC"
	MessageTag_ADEST_ROBVPIIC_ANS                              MessageTag = "ADEST_ROBVPIIC_ANS"
	MessageTag_ADEST_ROTVPIIC                                  MessageTag = "ADEST_ROTVPIIC"
	MessageTag_ADEST_ROTVPIIC_ANS                              MessageTag = "ADEST_ROTVPIIC_ANS"
	MessageTag_ADEST_RORPIIC                                   MessageTag = "ADEST_RORPIIC"
	MessageTag_ADEST_RORPIIC_ANS                               MessageTag = "ADEST_RORPIIC_ANS"
	MessageTag_BSV_IFCRC                                       MessageTag = "BSV_IFCRC"
	MessageTag_BSV_IFCRC_ANS                                   MessageTag = "BSV_IFCRC_ANS"
	MessageTag_BSV_IISCRC                                      MessageTag = "BSV_IISCRC"
	MessageTag_BSV_IISCRC_ANS                                  MessageTag = "BSV_IISCRC_ANS"
	MessageTag_BSV_IVCIRC                                      MessageTag = "BSV_IVCIRC"
	MessageTag_BSV_IVCIRC_ANS                                  MessageTag = "BSV_IVCIRC_ANS"
	MessageTag_BSV_IXCVRC                                      MessageTag = "BSV_IXCVRC"
	MessageTag_BSV_IXCVRC_ANS                                  MessageTag = "BSV_IXCVRC_ANS"
	MessageTag_BSV_ISCRAVTC                                    MessageTag = "BSV_ISCRAVTC"
	MessageTag_BSV_ISCRAVTC_ANS                                MessageTag = "BSV_ISCRAVTC_ANS"
	MessageTag_BSV_IVTAVRSC                                    MessageTag = "BSV_IVTAVRSC"
	MessageTag_BSV_IVTAVRSC_ANS                                MessageTag = "BSV_IVTAVRSC_ANS"
	MessageTag_BSV_ISCCTC                                      MessageTag = "BSV_ISCCTC"
	MessageTag_BSV_ISCCTC_ANS                                  MessageTag = "BSV_ISCCTC_ANS"
	MessageTag_BSV_ICTGTNASCRT                                 MessageTag = "BSV_ICTGTNASCRT"
	MessageTag_BSV_ICTGTNASCRT_ANS                             MessageTag = "BSV_ICTGTNASCRT_ANS"
	MessageTag_BSV_ICTGTNASCET                                 MessageTag = "BSV_ICTGTNASCET"
	MessageTag_BSV_ICTGTNASCET_ANS                             MessageTag = "BSV_ICTGTNASCET_ANS"
	MessageTag_BSV_ICVRC                                       MessageTag = "BSV_ICVRC"
	MessageTag_BSV_ICVRC_ANS                                   MessageTag = "BSV_ICVRC_ANS"
	MessageTag_BSV_ISAVRC                                      MessageTag = "BSV_ISAVRC"
	MessageTag_BSV_ISAVRC_ANS                                  MessageTag = "BSV_ISAVRC_ANS"
	MessageTag_BSV_ICTGTNACCET                                 MessageTag = "BSV_ICTGTNACCET"
	MessageTag_BSV_ICTGTNACCET_ANS                             MessageTag = "BSV_ICTGTNACCET_ANS"
	MessageTag_BSV_IEAAAVRC                                    MessageTag = "BSV_IEAAAVRC"
	MessageTag_BSV_IEAAAVRC_ANS                                MessageTag = "BSV_IEAAAVRC_ANS"
	MessageTag_LTV_ABSV                                        MessageTag = "LTV_ABSV"
	MessageTag_LTV_ABSV_ANS                                    MessageTag = "LTV_ABSV_ANS"
	MessageTag_LTV_ISCKNR                                      MessageTag = "LTV_ISCKNR"
	MessageTag_LTV_ISCKNR_ANS0                                 MessageTag = "LTV_ISCKNR_ANS0"
	MessageTag_LTV_ISCKNR_ANS1                                 MessageTag = "LTV_ISCKNR_ANS1"
	MessageTag_ARCH_LTVV                                       MessageTag = "ARCH_LTVV"
	MessageTag_ARCH_LTVV_ANS                                   MessageTag = "ARCH_LTVV_ANS"
	MessageTag_ARCH_LTAIVMP                                    MessageTag = "ARCH_LTAIVMP"
	MessageTag_ARCH_LTAIVMP_ANS                                MessageTag = "ARCH_LTAIVMP_ANS"
	MessageTag_ARCH_IRTVBBA                                    MessageTag = "ARCH_IRTVBBA"
	MessageTag_ARCH_IRTVBBA_ANS                                MessageTag = "ARCH_IRTVBBA_ANS"
	MessageTag_ARCH_ICHFCRLPOET                                MessageTag = "ARCH_ICHFCRLPOET"
	MessageTag_ARCH_ICHFCRLPOET_ANS                            MessageTag = "ARCH_ICHFCRLPOET_ANS"
	MessageTag_ACCM                                            MessageTag = "ACCM"
	MessageTag_ACCM_ANS                                        MessageTag = "ACCM_ANS"
	MessageTag_ASCCM_CAA                                       MessageTag = "ASCCM_CAA"
	MessageTag_ASCCM_CAA_ANS                                   MessageTag = "ASCCM_CAA_ANS"
	MessageTag_ASCCM_DAA                                       MessageTag = "ASCCM_DAA"
	MessageTag_ASCCM_DAA_ANS                                   MessageTag = "ASCCM_DAA_ANS"
	MessageTag_ASCCM_DAA_ANS_2                                 MessageTag = "ASCCM_DAA_ANS_2"
	MessageTag_ASCCM_APKSA                                     MessageTag = "ASCCM_APKSA"
	MessageTag_ASCCM_APKSA_ANS                                 MessageTag = "ASCCM_APKSA_ANS"
	MessageTag_ASCCM_APKSA_ANS_2                               MessageTag = "ASCCM_APKSA_ANS_2"
	MessageTag_ASCCM_AR                                        MessageTag = "ASCCM_AR"
	MessageTag_ASCCM_AR_ANS_ANR                                MessageTag = "ASCCM_AR_ANS_ANR"
	MessageTag_ASCCM_AR_ANS_ANR_2                              MessageTag = "ASCCM_AR_ANS_ANR_2"
	MessageTag_ASCCM_AR_ANS_AKSNR                              MessageTag = "ASCCM_AR_ANS_AKSNR"
	MessageTag_ASCCM_AR_ANS_AKSNR_2                            MessageTag = "ASCCM_AR_ANS_AKSNR_2"
	MessageTag_ASCCM_PKSK                                      MessageTag = "ASCCM_PKSK"
	MessageTag_ASCCM_PKSK_ANS                                  MessageTag = "ASCCM_PKSK_ANS"
	MessageTag_ACCM_DESC_WITH_ID                               MessageTag = "ACCM_DESC_WITH_ID"
	MessageTag_ACCM_DESC_WITH_ID_RESULT                        MessageTag = "ACCM_DESC_WITH_ID_RESULT"
	MessageTag_ACCM_DESC_WITH_NAME                             MessageTag = "ACCM_DESC_WITH_NAME"
	MessageTag_ACCM_POS_SIG_SIG                                MessageTag = "ACCM_POS_SIG_SIG"
	MessageTag_ACCM_POS_TST_SIG                                MessageTag = "ACCM_POS_TST_SIG"
	MessageTag_ACCM_POS_REVOC_SIG                              MessageTag = "ACCM_POS_REVOC_SIG"
	MessageTag_ACCM_POS_EV_RECORD                              MessageTag = "ACCM_POS_EV_RECORD"
	MessageTag_ACCM_POS_EAA                                    MessageTag = "ACCM_POS_EAA"
	MessageTag_ACCM_POS_EAA_REV                                MessageTag = "ACCM_POS_EAA_REV"
	MessageTag_ACCM_POS_CNTR_SIG                               MessageTag = "ACCM_POS_CNTR_SIG"
	MessageTag_ACCM_POS_CNTR_SIG_PL                            MessageTag = "ACCM_POS_CNTR_SIG_PL"
	MessageTag_ACCM_POS_CON_DIG                                MessageTag = "ACCM_POS_CON_DIG"
	MessageTag_ACCM_POS_EAA_KB                                 MessageTag = "ACCM_POS_EAA_KB"
	MessageTag_ACCM_POS_EAA_ND                                 MessageTag = "ACCM_POS_EAA_ND"
	MessageTag_ACCM_POS_EAA_NSD                                MessageTag = "ACCM_POS_EAA_NSD"
	MessageTag_ACCM_POS_EAA_NSD_PL                             MessageTag = "ACCM_POS_EAA_NSD_PL"
	MessageTag_ACCM_POS_EAA_OSDC                               MessageTag = "ACCM_POS_EAA_OSDC"
	MessageTag_ACCM_POS_EAA_OSDC_PL                            MessageTag = "ACCM_POS_EAA_OSDC_PL"
	MessageTag_ACCM_POS_EAA_PD                                 MessageTag = "ACCM_POS_EAA_PD"
	MessageTag_ACCM_POS_EAA_SD                                 MessageTag = "ACCM_POS_EAA_SD"
	MessageTag_ACCM_POS_EAA_SD_PL                              MessageTag = "ACCM_POS_EAA_SD_PL"
	MessageTag_ACCM_POS_ER_ADO                                 MessageTag = "ACCM_POS_ER_ADO"
	MessageTag_ACCM_POS_ER_ADO_PL                              MessageTag = "ACCM_POS_ER_ADO_PL"
	MessageTag_ACCM_POS_ER_OR                                  MessageTag = "ACCM_POS_ER_OR"
	MessageTag_ACCM_POS_ER_OR_PL                               MessageTag = "ACCM_POS_ER_OR_PL"
	MessageTag_ACCM_POS_ER_TST                                 MessageTag = "ACCM_POS_ER_TST"
	MessageTag_ACCM_POS_ER_TST_SEQ                             MessageTag = "ACCM_POS_ER_TST_SEQ"
	MessageTag_ACCM_POS_ER_MST_SIG                             MessageTag = "ACCM_POS_ER_MST_SIG"
	MessageTag_ACCM_POS_JWS                                    MessageTag = "ACCM_POS_JWS"
	MessageTag_ACCM_POS_COSE                                   MessageTag = "ACCM_POS_COSE"
	MessageTag_ACCM_POS_KEY                                    MessageTag = "ACCM_POS_KEY"
	MessageTag_ACCM_POS_KEY_PL                                 MessageTag = "ACCM_POS_KEY_PL"
	MessageTag_ACCM_POS_MAN                                    MessageTag = "ACCM_POS_MAN"
	MessageTag_ACCM_POS_MAN_PL                                 MessageTag = "ACCM_POS_MAN_PL"
	MessageTag_ACCM_POS_MAN_ENT                                MessageTag = "ACCM_POS_MAN_ENT"
	MessageTag_ACCM_POS_MAN_ENT_PL                             MessageTag = "ACCM_POS_MAN_ENT_PL"
	MessageTag_ACCM_POS_MES_DIG                                MessageTag = "ACCM_POS_MES_DIG"
	MessageTag_ACCM_POS_MESS_IMP                               MessageTag = "ACCM_POS_MESS_IMP"
	MessageTag_ACCM_POS_REF                                    MessageTag = "ACCM_POS_REF"
	MessageTag_ACCM_POS_REF_PL                                 MessageTag = "ACCM_POS_REF_PL"
	MessageTag_ACCM_POS_SIG_D_ENT                              MessageTag = "ACCM_POS_SIG_D_ENT"
	MessageTag_ACCM_POS_SIG_D_ENT_PL                           MessageTag = "ACCM_POS_SIG_D_ENT_PL"
	MessageTag_ACCM_POS_SIG_VAL_AND_PRT                        MessageTag = "ACCM_POS_SIG_VAL_AND_PRT"
	MessageTag_ACCM_POS_SIGND_OBJ                              MessageTag = "ACCM_POS_SIGND_OBJ"
	MessageTag_ACCM_POS_SIGND_PRT                              MessageTag = "ACCM_POS_SIGND_PRT"
	MessageTag_ACCM_POS_SIGNTR_PRT                             MessageTag = "ACCM_POS_SIGNTR_PRT"
	MessageTag_ACCM_POS_CERT_CHAIN_SIG                         MessageTag = "ACCM_POS_CERT_CHAIN_SIG"
	MessageTag_ACCM_POS_CERT_CHAIN_TST                         MessageTag = "ACCM_POS_CERT_CHAIN_TST"
	MessageTag_ACCM_POS_CERT_CHAIN_REVOC                       MessageTag = "ACCM_POS_CERT_CHAIN_REVOC"
	MessageTag_ACCM_POS_CERT_CHAIN_EAA_REV                     MessageTag = "ACCM_POS_CERT_CHAIN_EAA_REV"
	MessageTag_ACCM_POS_CERT_CHAIN                             MessageTag = "ACCM_POS_CERT_CHAIN"
	MessageTag_ACCM_POS_SIG_CERT_REF                           MessageTag = "ACCM_POS_SIG_CERT_REF"
	MessageTag_BBB_SAV_DSCACRCC                                MessageTag = "BBB_SAV_DSCACRCC"
	MessageTag_BBB_SAV_DSCACRCC_ANS                            MessageTag = "BBB_SAV_DSCACRCC_ANS"
	MessageTag_BBB_SAV_ACPCCRSCA                               MessageTag = "BBB_SAV_ACPCCRSCA"
	MessageTag_BBB_SAV_ACPCCRSCA_ANS                           MessageTag = "BBB_SAV_ACPCCRSCA_ANS"
	MessageTag_BBB_SAV_ISVA                                    MessageTag = "BBB_SAV_ISVA"
	MessageTag_BBB_SAV_ISVA_ANS                                MessageTag = "BBB_SAV_ISVA_ANS"
	MessageTag_BBB_SAV_ISSV                                    MessageTag = "BBB_SAV_ISSV"
	MessageTag_BBB_SAV_ISSV_ANS                                MessageTag = "BBB_SAV_ISSV_ANS"
	MessageTag_BBB_SAV_ICERRM                                  MessageTag = "BBB_SAV_ICERRM"
	MessageTag_BBB_SAV_ICERRM_ANS                              MessageTag = "BBB_SAV_ICERRM_ANS"
	MessageTag_BBB_SAV_ICRM                                    MessageTag = "BBB_SAV_ICRM"
	MessageTag_BBB_SAV_ICRM_ANS                                MessageTag = "BBB_SAV_ICRM_ANS"
	MessageTag_BBB_SAV_ISQPCTP                                 MessageTag = "BBB_SAV_ISQPCTP"
	MessageTag_BBB_SAV_ISQPCTP_ANS                             MessageTag = "BBB_SAV_ISQPCTP_ANS"
	MessageTag_BBB_SAV_ISQPCHP                                 MessageTag = "BBB_SAV_ISQPCHP"
	MessageTag_BBB_SAV_ISQPCHP_ANS                             MessageTag = "BBB_SAV_ISQPCHP_ANS"
	MessageTag_BBB_SAV_ISQPCIP                                 MessageTag = "BBB_SAV_ISQPCIP"
	MessageTag_BBB_SAV_ISQPCIP_ANS                             MessageTag = "BBB_SAV_ISQPCIP_ANS"
	MessageTag_BBB_SAV_ISQPCTSIP                               MessageTag = "BBB_SAV_ISQPCTSIP"
	MessageTag_BBB_SAV_ISQPCTSIP_ANS                           MessageTag = "BBB_SAV_ISQPCTSIP_ANS"
	MessageTag_BBB_SAV_ISQPSTYPP                               MessageTag = "BBB_SAV_ISQPSTYPP"
	MessageTag_BBB_SAV_ISQPSTYPP_ANS                           MessageTag = "BBB_SAV_ISQPSTYPP_ANS"
	MessageTag_BBB_SAV_ISQPSLP                                 MessageTag = "BBB_SAV_ISQPSLP"
	MessageTag_BBB_SAV_ISQPSLP_ANS                             MessageTag = "BBB_SAV_ISQPSLP_ANS"
	MessageTag_BBB_SAV_ISQPSTP                                 MessageTag = "BBB_SAV_ISQPSTP"
	MessageTag_BBB_SAV_ISQPSTP_ANS                             MessageTag = "BBB_SAV_ISQPSTP_ANS"
	MessageTag_BBB_SAV_ISQPSTWSCVR                             MessageTag = "BBB_SAV_ISQPSTWSCVR"
	MessageTag_BBB_SAV_ISQPSTWSCVR_ANS                         MessageTag = "BBB_SAV_ISQPSTWSCVR_ANS"
	MessageTag_BBB_SAV_ISQPXTIP                                MessageTag = "BBB_SAV_ISQPXTIP"
	MessageTag_BBB_SAV_ISQPXTIP_ANS                            MessageTag = "BBB_SAV_ISQPXTIP_ANS"
	MessageTag_BBB_SAV_IUQPCSP                                 MessageTag = "BBB_SAV_IUQPCSP"
	MessageTag_BBB_SAV_IUQPCSP_ANS                             MessageTag = "BBB_SAV_IUQPCSP_ANS"
	MessageTag_BBB_SAV_IUQPSTSP                                MessageTag = "BBB_SAV_IUQPSTSP"
	MessageTag_BBB_SAV_IUQPSTSP_ANS                            MessageTag = "BBB_SAV_IUQPSTSP_ANS"
	MessageTag_BBB_SAV_IUQPVDTSP                               MessageTag = "BBB_SAV_IUQPVDTSP"
	MessageTag_BBB_SAV_IUQPVDTSP_ANS                           MessageTag = "BBB_SAV_IUQPVDTSP_ANS"
	MessageTag_BBB_SAV_IUQPVDROTSP                             MessageTag = "BBB_SAV_IUQPVDROTSP"
	MessageTag_BBB_SAV_IUQPVDROTSP_ANS                         MessageTag = "BBB_SAV_IUQPVDROTSP_ANS"
	MessageTag_BBB_SAV_IUQPATSP                                MessageTag = "BBB_SAV_IUQPATSP"
	MessageTag_BBB_SAV_IUQPATSP_ANS                            MessageTag = "BBB_SAV_IUQPATSP_ANS"
	MessageTag_BBB_SAV_ICTVS                                   MessageTag = "BBB_SAV_ICTVS"
	MessageTag_BBB_SAV_ICTVS_ANS                               MessageTag = "BBB_SAV_ICTVS_ANS"
	MessageTag_BBB_SAV_IDTSP                                   MessageTag = "BBB_SAV_IDTSP"
	MessageTag_BBB_SAV_IDTSP_ANS                               MessageTag = "BBB_SAV_IDTSP_ANS"
	MessageTag_BBB_SAV_ITVS                                    MessageTag = "BBB_SAV_ITVS"
	MessageTag_BBB_SAV_ITVS_ANS                                MessageTag = "BBB_SAV_ITVS_ANS"
	MessageTag_BBB_SAV_IVTTSTP                                 MessageTag = "BBB_SAV_IVTTSTP"
	MessageTag_BBB_SAV_IVTTSTP_ANS                             MessageTag = "BBB_SAV_IVTTSTP_ANS"
	MessageTag_BBB_SAV_IVLTATSTP                               MessageTag = "BBB_SAV_IVLTATSTP"
	MessageTag_BBB_SAV_IVLTATSTP_ANS                           MessageTag = "BBB_SAV_IVLTATSTP_ANS"
	MessageTag_BBB_SAV_ISQPMDOSPP                              MessageTag = "BBB_SAV_ISQPMDOSPP"
	MessageTag_BBB_SAV_ISQPMDOSPP_ANS                          MessageTag = "BBB_SAV_ISQPMDOSPP_ANS"
	MessageTag_BBB_SAV_DMICTSTMCMI                             MessageTag = "BBB_SAV_DMICTSTMCMI"
	MessageTag_BBB_SAV_DMICTSTMCMI_ANS                         MessageTag = "BBB_SAV_DMICTSTMCMI_ANS"
	MessageTag_BBB_TAV_ITSAP                                   MessageTag = "BBB_TAV_ITSAP"
	MessageTag_BBB_TAV_ITSAP_ANS                               MessageTag = "BBB_TAV_ITSAP_ANS"
	MessageTag_BBB_TAV_DTSAVM                                  MessageTag = "BBB_TAV_DTSAVM"
	MessageTag_BBB_TAV_DTSAVM_ANS                              MessageTag = "BBB_TAV_DTSAVM_ANS"
	MessageTag_BBB_TAV_DTSAOM                                  MessageTag = "BBB_TAV_DTSAOM"
	MessageTag_BBB_TAV_DTSAOM_ANS                              MessageTag = "BBB_TAV_DTSAOM_ANS"
	MessageTag_BBB_VCI_ISPK                                    MessageTag = "BBB_VCI_ISPK"
	MessageTag_BBB_VCI_ISPK_ANS                                MessageTag = "BBB_VCI_ISPK_ANS"
	MessageTag_BBB_VCI_ISPA                                    MessageTag = "BBB_VCI_ISPA"
	MessageTag_BBB_VCI_ISPA_ANS                                MessageTag = "BBB_VCI_ISPA_ANS"
	MessageTag_BBB_VCI_ISPSUPP                                 MessageTag = "BBB_VCI_ISPSUPP"
	MessageTag_BBB_VCI_ISPSUPP_ANS                             MessageTag = "BBB_VCI_ISPSUPP_ANS"
	MessageTag_BBB_VCI_ISPM                                    MessageTag = "BBB_VCI_ISPM"
	MessageTag_BBB_VCI_ISPM_ANS                                MessageTag = "BBB_VCI_ISPM_ANS"
	MessageTag_BBB_VCI_IZHSP                                   MessageTag = "BBB_VCI_IZHSP"
	MessageTag_BBB_VCI_IZHSP_ANS                               MessageTag = "BBB_VCI_IZHSP_ANS"
	MessageTag_BBB_XCV_SUB                                     MessageTag = "BBB_XCV_SUB"
	MessageTag_BBB_XCV_SUB_ANS                                 MessageTag = "BBB_XCV_SUB_ANS"
	MessageTag_BBB_XCV_SUB_ANS_2                               MessageTag = "BBB_XCV_SUB_ANS_2"
	MessageTag_BBB_XCV_RFC                                     MessageTag = "BBB_XCV_RFC"
	MessageTag_BBB_XCV_RFC_ANS                                 MessageTag = "BBB_XCV_RFC_ANS"
	MessageTag_BBB_XCV_RAC                                     MessageTag = "BBB_XCV_RAC"
	MessageTag_BBB_XCV_RAC_ANS                                 MessageTag = "BBB_XCV_RAC_ANS"
	MessageTag_BBB_XCV_CCCBB                                   MessageTag = "BBB_XCV_CCCBB"
	MessageTag_BBB_XCV_CCCBB_ANS                               MessageTag = "BBB_XCV_CCCBB_ANS"
	MessageTag_BBB_XCV_CCCBB_SIG_ANS                           MessageTag = "BBB_XCV_CCCBB_SIG_ANS"
	MessageTag_BBB_XCV_CCCBB_TSP_ANS                           MessageTag = "BBB_XCV_CCCBB_TSP_ANS"
	MessageTag_BBB_XCV_CCCBB_REV_ANS                           MessageTag = "BBB_XCV_CCCBB_REV_ANS"
	MessageTag_BBB_XCV_CMDCIPI                                 MessageTag = "BBB_XCV_CMDCIPI"
	MessageTag_BBB_XCV_CMDCIPI_ANS                             MessageTag = "BBB_XCV_CMDCIPI_ANS"
	MessageTag_BBB_XCV_CMDCIQC                                 MessageTag = "BBB_XCV_CMDCIQC"
	MessageTag_BBB_XCV_CMDCIQC_ANS                             MessageTag = "BBB_XCV_CMDCIQC_ANS"
	MessageTag_BBB_XCV_CMDCIQSCD                               MessageTag = "BBB_XCV_CMDCIQSCD"
	MessageTag_BBB_XCV_CMDCIQSCD_ANS                           MessageTag = "BBB_XCV_CMDCIQSCD_ANS"
	MessageTag_BBB_XCV_CMDCIITLP                               MessageTag = "BBB_XCV_CMDCIITLP"
	MessageTag_BBB_XCV_CMDCIITLP_ANS                           MessageTag = "BBB_XCV_CMDCIITLP_ANS"
	MessageTag_BBB_XCV_CMDCIITNP                               MessageTag = "BBB_XCV_CMDCIITNP"
	MessageTag_BBB_XCV_CMDCIITNP_ANS                           MessageTag = "BBB_XCV_CMDCIITNP_ANS"
	MessageTag_BBB_XCV_CMDCICQCC                               MessageTag = "BBB_XCV_CMDCICQCC"
	MessageTag_BBB_XCV_CMDCICQCC_ANS                           MessageTag = "BBB_XCV_CMDCICQCC_ANS"
	MessageTag_BBB_XCV_CMDCICQCLVA                             MessageTag = "BBB_XCV_CMDCICQCLVA"
	MessageTag_BBB_XCV_CMDCICQCLVA_ANS                         MessageTag = "BBB_XCV_CMDCICQCLVA_ANS"
	MessageTag_BBB_XCV_CMDCICQCLVHAC                           MessageTag = "BBB_XCV_CMDCICQCLVHAC"
	MessageTag_BBB_XCV_CMDCICQCLVHAC_ANS                       MessageTag = "BBB_XCV_CMDCICQCLVHAC_ANS"
	MessageTag_BBB_XCV_CMDCICQCERPA                            MessageTag = "BBB_XCV_CMDCICQCERPA"
	MessageTag_BBB_XCV_CMDCICQCERPA_ANS                        MessageTag = "BBB_XCV_CMDCICQCERPA_ANS"
	MessageTag_BBB_XCV_CMDCICSQCSSCD                           MessageTag = "BBB_XCV_CMDCICSQCSSCD"
	MessageTag_BBB_XCV_CMDCICSQCSSCD_ANS                       MessageTag = "BBB_XCV_CMDCICSQCSSCD_ANS"
	MessageTag_BBB_XCV_CMDCICQCPDSLA                           MessageTag = "BBB_XCV_CMDCICQCPDSLA"
	MessageTag_BBB_XCV_CMDCICQCPDSLA_ANS                       MessageTag = "BBB_XCV_CMDCICQCPDSLA_ANS"
	MessageTag_BBB_XCV_CMDCICQCTA                              MessageTag = "BBB_XCV_CMDCICQCTA"
	MessageTag_BBB_XCV_CMDCICQCTA_ANS                          MessageTag = "BBB_XCV_CMDCICQCTA_ANS"
	MessageTag_BBB_XCV_CMDCDCQCCLCEC                           MessageTag = "BBB_XCV_CMDCDCQCCLCEC"
	MessageTag_BBB_XCV_CMDCDCQCCLCEC_ANS                       MessageTag = "BBB_XCV_CMDCDCQCCLCEC_ANS"
	MessageTag_BBB_XCV_CMDCDCQCCLCEC_ANS_EU                    MessageTag = "BBB_XCV_CMDCDCQCCLCEC_ANS_EU"
	MessageTag_BBB_XCV_CMDCSCSIA                               MessageTag = "BBB_XCV_CMDCSCSIA"
	MessageTag_BBB_XCV_CMDCSCSIA_ANS                           MessageTag = "BBB_XCV_CMDCSCSIA_ANS"
	MessageTag_BBB_XCV_CMDCICQCRA                              MessageTag = "BBB_XCV_CMDCICQCRA"
	MessageTag_BBB_XCV_CMDCICQCRA_ANS                          MessageTag = "BBB_XCV_CMDCICQCRA_ANS"
	MessageTag_BBB_XCV_CMDCICQCNA                              MessageTag = "BBB_XCV_CMDCICQCNA"
	MessageTag_BBB_XCV_CMDCICQCNA_ANS                          MessageTag = "BBB_XCV_CMDCICQCNA_ANS"
	MessageTag_BBB_XCV_CMDCICQCIA                              MessageTag = "BBB_XCV_CMDCICQCIA"
	MessageTag_BBB_XCV_CMDCICQCIA_ANS                          MessageTag = "BBB_XCV_CMDCICQCIA_ANS"
	MessageTag_BBB_XCV_CMDCDCQCQSCDLSA                         MessageTag = "BBB_XCV_CMDCDCQCQSCDLSA"
	MessageTag_BBB_XCV_CMDCDCQCQSCDLSA_ANS                     MessageTag = "BBB_XCV_CMDCDCQCQSCDLSA_ANS"
	MessageTag_BBB_XCV_CMDCDCQCIMSA                            MessageTag = "BBB_XCV_CMDCDCQCIMSA"
	MessageTag_BBB_XCV_CMDCDCQCIMSA_ANS                        MessageTag = "BBB_XCV_CMDCDCQCIMSA_ANS"
	MessageTag_BBB_XCV_CMDCPSBCLA                              MessageTag = "BBB_XCV_CMDCPSBCLA"
	MessageTag_BBB_XCV_CMDCPSBCLA_ANS                          MessageTag = "BBB_XCV_CMDCPSBCLA_ANS"
	MessageTag_BBB_XCV_CMDCPSBASIA                             MessageTag = "BBB_XCV_CMDCPSBASIA"
	MessageTag_BBB_XCV_CMDCPSBASIA_ANS                         MessageTag = "BBB_XCV_CMDCPSBASIA_ANS"
	MessageTag_BBB_XCV_CMDCPSBLIA                              MessageTag = "BBB_XCV_CMDCPSBLIA"
	MessageTag_BBB_XCV_CMDCPSBLIA_ANS                          MessageTag = "BBB_XCV_CMDCPSBLIA_ANS"
	MessageTag_BBB_XCV_DCCUCE                                  MessageTag = "BBB_XCV_DCCUCE"
	MessageTag_BBB_XCV_DCCUCE_ANS                              MessageTag = "BBB_XCV_DCCUCE_ANS"
	MessageTag_BBB_XCV_DCCFCE                                  MessageTag = "BBB_XCV_DCCFCE"
	MessageTag_BBB_XCV_DCCFCE_ANS                              MessageTag = "BBB_XCV_DCCFCE_ANS"
	MessageTag_BBB_XCV_DCSBSINC                                MessageTag = "BBB_XCV_DCSBSINC"
	MessageTag_BBB_XCV_DCSBSINC_ANS                            MessageTag = "BBB_XCV_DCSBSINC_ANS"
	MessageTag_BBB_XCV_ICAC                                    MessageTag = "BBB_XCV_ICAC"
	MessageTag_BBB_XCV_ICAC_ANS                                MessageTag = "BBB_XCV_ICAC_ANS"
	MessageTag_BBB_XCV_ICPDV                                   MessageTag = "BBB_XCV_ICPDV"
	MessageTag_BBB_XCV_ICPDV_ANS                               MessageTag = "BBB_XCV_ICPDV_ANS"
	MessageTag_BBB_XCV_ICPTV                                   MessageTag = "BBB_XCV_ICPTV"
	MessageTag_BBB_XCV_ICPTV_ANS                               MessageTag = "BBB_XCV_ICPTV_ANS"
	MessageTag_BBB_XCV_IAKIP                                   MessageTag = "BBB_XCV_IAKIP"
	MessageTag_BBB_XCV_IAKIP_ANS                               MessageTag = "BBB_XCV_IAKIP_ANS"
	MessageTag_BBB_XCV_ISKIP                                   MessageTag = "BBB_XCV_ISKIP"
	MessageTag_BBB_XCV_ISKIP_ANS                               MessageTag = "BBB_XCV_ISKIP_ANS"
	MessageTag_BBB_XCV_ICNRAEV                                 MessageTag = "BBB_XCV_ICNRAEV"
	MessageTag_BBB_XCV_ICNRAEV_ANS                             MessageTag = "BBB_XCV_ICNRAEV_ANS"
	MessageTag_BBB_XCV_IVTBCTSD                                MessageTag = "BBB_XCV_IVTBCTSD"
	MessageTag_BBB_XCV_IVTBCTSD_ANS                            MessageTag = "BBB_XCV_IVTBCTSD_ANS"
	MessageTag_BBB_XCV_ICTIVRSC                                MessageTag = "BBB_XCV_ICTIVRSC"
	MessageTag_BBB_XCV_ICTIVRSC_ANS                            MessageTag = "BBB_XCV_ICTIVRSC_ANS"
	MessageTag_BBB_XCV_ICTIVRCIRI                              MessageTag = "BBB_XCV_ICTIVRCIRI"
	MessageTag_BBB_XCV_ICTIVRCIRI_ANS                          MessageTag = "BBB_XCV_ICTIVRCIRI_ANS"
	MessageTag_BBB_XCV_IRDCSFC                                 MessageTag = "BBB_XCV_IRDCSFC"
	MessageTag_BBB_XCV_IRDCSFC_ANS                             MessageTag = "BBB_XCV_IRDCSFC_ANS"
	MessageTag_BBB_XCV_IRDPFC                                  MessageTag = "BBB_XCV_IRDPFC"
	MessageTag_BBB_XCV_IRDPFC_ANS                              MessageTag = "BBB_XCV_IRDPFC_ANS"
	MessageTag_BBB_XCV_IRDPFRC                                 MessageTag = "BBB_XCV_IRDPFRC"
	MessageTag_BBB_XCV_IRDPFRC_ANS                             MessageTag = "BBB_XCV_IRDPFRC_ANS"
	MessageTag_BBB_XCV_IARDPFC                                 MessageTag = "BBB_XCV_IARDPFC"
	MessageTag_BBB_XCV_IARDPFC_ANS                             MessageTag = "BBB_XCV_IARDPFC_ANS"
	MessageTag_BBB_VTS_IRDPFC                                  MessageTag = "BBB_VTS_IRDPFC"
	MessageTag_BBB_VTS_IRDPFC_ANS                              MessageTag = "BBB_VTS_IRDPFC_ANS"
	MessageTag_BBB_XCV_ISCOH                                   MessageTag = "BBB_XCV_ISCOH"
	MessageTag_BBB_XCV_ISCOH_ANS                               MessageTag = "BBB_XCV_ISCOH_ANS"
	MessageTag_BBB_XCV_ISCUKN                                  MessageTag = "BBB_XCV_ISCUKN"
	MessageTag_BBB_XCV_ISCUKN_ANS                              MessageTag = "BBB_XCV_ISCUKN_ANS"
	MessageTag_BBB_XCV_ISCR                                    MessageTag = "BBB_XCV_ISCR"
	MessageTag_BBB_XCV_ISCR_ANS                                MessageTag = "BBB_XCV_ISCR_ANS"
	MessageTag_BBB_XCV_ISCGKU                                  MessageTag = "BBB_XCV_ISCGKU"
	MessageTag_BBB_XCV_ISCGKU_ANS                              MessageTag = "BBB_XCV_ISCGKU_ANS"
	MessageTag_BBB_XCV_ISCGKU_ANS_CERT                         MessageTag = "BBB_XCV_ISCGKU_ANS_CERT"
	MessageTag_BBB_XCV_ISCGEKU                                 MessageTag = "BBB_XCV_ISCGEKU"
	MessageTag_BBB_XCV_ISCGEKU_ANS                             MessageTag = "BBB_XCV_ISCGEKU_ANS"
	MessageTag_BBB_XCV_ISCGEKU_ANS_CERT                        MessageTag = "BBB_XCV_ISCGEKU_ANS_CERT"
	MessageTag_BBB_XCV_ICSI                                    MessageTag = "BBB_XCV_ICSI"
	MessageTag_BBB_XCV_ICSI_ANS                                MessageTag = "BBB_XCV_ICSI_ANS"
	MessageTag_BBB_XCV_IOTAA                                   MessageTag = "BBB_XCV_IOTAA"
	MessageTag_BBB_XCV_IOTAA_ANS                               MessageTag = "BBB_XCV_IOTAA_ANS"
	MessageTag_BBB_XCV_HPCCVVT                                 MessageTag = "BBB_XCV_HPCCVVT"
	MessageTag_BBB_XCV_HPCCVVT_ANS                             MessageTag = "BBB_XCV_HPCCVVT_ANS"
	MessageTag_BBB_XCV_PSEUDO_USE                              MessageTag = "BBB_XCV_PSEUDO_USE"
	MessageTag_BBB_XCV_PSEUDO_USE_ANS                          MessageTag = "BBB_XCV_PSEUDO_USE_ANS"
	MessageTag_BBB_XCV_AIA_PRES                                MessageTag = "BBB_XCV_AIA_PRES"
	MessageTag_BBB_XCV_AIA_PRES_ANS                            MessageTag = "BBB_XCV_AIA_PRES_ANS"
	MessageTag_BBB_XCV_REVOC_PRES                              MessageTag = "BBB_XCV_REVOC_PRES"
	MessageTag_BBB_XCV_REVOC_PRES_ANS                          MessageTag = "BBB_XCV_REVOC_PRES_ANS"
	MessageTag_BBB_XCV_REVOC_THIS_UPDATE_PRESENT               MessageTag = "BBB_XCV_REVOC_THIS_UPDATE_PRESENT"
	MessageTag_BBB_XCV_REVOC_THIS_UPDATE_PRESENT_ANS           MessageTag = "BBB_XCV_REVOC_THIS_UPDATE_PRESENT_ANS"
	MessageTag_BBB_XCV_REVOC_ISSUER_KNOWN                      MessageTag = "BBB_XCV_REVOC_ISSUER_KNOWN"
	MessageTag_BBB_XCV_REVOC_ISSUER_KNOWN_ANS                  MessageTag = "BBB_XCV_REVOC_ISSUER_KNOWN_ANS"
	MessageTag_BBB_XCV_REVOC_ISSUER_VALID_AT_PROD              MessageTag = "BBB_XCV_REVOC_ISSUER_VALID_AT_PROD"
	MessageTag_BBB_XCV_REVOC_ISSUER_VALID_AT_PROD_ANS          MessageTag = "BBB_XCV_REVOC_ISSUER_VALID_AT_PROD_ANS"
	MessageTag_BBB_XCV_REVOC_AFTER_CERT_NOT_BEFORE             MessageTag = "BBB_XCV_REVOC_AFTER_CERT_NOT_BEFORE"
	MessageTag_BBB_XCV_REVOC_AFTER_CERT_NOT_BEFORE_ANS         MessageTag = "BBB_XCV_REVOC_AFTER_CERT_NOT_BEFORE_ANS"
	MessageTag_BBB_XCV_REVOC_HAS_CERT_INFO                     MessageTag = "BBB_XCV_REVOC_HAS_CERT_INFO"
	MessageTag_BBB_XCV_REVOC_HAS_CERT_INFO_ANS                 MessageTag = "BBB_XCV_REVOC_HAS_CERT_INFO_ANS"
	MessageTag_BBB_XCV_REVOC_RESPID_MATCH                      MessageTag = "BBB_XCV_REVOC_RESPID_MATCH"
	MessageTag_BBB_XCV_REVOC_RESPID_MATCH_ANS                  MessageTag = "BBB_XCV_REVOC_RESPID_MATCH_ANS"
	MessageTag_BBB_XCV_REVOC_CERT_HASH_PRESENT                 MessageTag = "BBB_XCV_REVOC_CERT_HASH_PRESENT"
	MessageTag_BBB_XCV_REVOC_CERT_HASH_PRESENT_ANS             MessageTag = "BBB_XCV_REVOC_CERT_HASH_PRESENT_ANS"
	MessageTag_BBB_XCV_REVOC_CERT_HASH_MATCH                   MessageTag = "BBB_XCV_REVOC_CERT_HASH_MATCH"
	MessageTag_BBB_XCV_REVOC_CERT_HASH_MATCH_ANS               MessageTag = "BBB_XCV_REVOC_CERT_HASH_MATCH_ANS"
	MessageTag_BBB_XCV_REVOC_SELF_ISSUED_OCSP                  MessageTag = "BBB_XCV_REVOC_SELF_ISSUED_OCSP"
	MessageTag_BBB_XCV_REVOC_SELF_ISSUED_OCSP_ANS              MessageTag = "BBB_XCV_REVOC_SELF_ISSUED_OCSP_ANS"
	MessageTag_BBB_XCV_DCIDNMSDNIC                             MessageTag = "BBB_XCV_DCIDNMSDNIC"
	MessageTag_BBB_XCV_DCIDNMSDNIC_ANS                         MessageTag = "BBB_XCV_DCIDNMSDNIC_ANS"
	MessageTag_BBB_XCV_ISCGCOUN                                MessageTag = "BBB_XCV_ISCGCOUN"
	MessageTag_BBB_XCV_ISCGCOUN_ANS                            MessageTag = "BBB_XCV_ISCGCOUN_ANS"
	MessageTag_BBB_XCV_ISCGLOC                                 MessageTag = "BBB_XCV_ISCGLOC"
	MessageTag_BBB_XCV_ISCGLOC_ANS                             MessageTag = "BBB_XCV_ISCGLOC_ANS"
	MessageTag_BBB_XCV_ISCGST                                  MessageTag = "BBB_XCV_ISCGST"
	MessageTag_BBB_XCV_ISCGST_ANS                              MessageTag = "BBB_XCV_ISCGST_ANS"
	MessageTag_BBB_XCV_ISCGORGAN                               MessageTag = "BBB_XCV_ISCGORGAN"
	MessageTag_BBB_XCV_ISCGORGAN_ANS                           MessageTag = "BBB_XCV_ISCGORGAN_ANS"
	MessageTag_BBB_XCV_ISCGORGAU                               MessageTag = "BBB_XCV_ISCGORGAU"
	MessageTag_BBB_XCV_ISCGORGAU_ANS                           MessageTag = "BBB_XCV_ISCGORGAU_ANS"
	MessageTag_BBB_XCV_ISCGORGAI                               MessageTag = "BBB_XCV_ISCGORGAI"
	MessageTag_BBB_XCV_ISCGORGAI_ANS                           MessageTag = "BBB_XCV_ISCGORGAI_ANS"
	MessageTag_BBB_XCV_ISCGSURN                                MessageTag = "BBB_XCV_ISCGSURN"
	MessageTag_BBB_XCV_ISCGSURN_ANS                            MessageTag = "BBB_XCV_ISCGSURN_ANS"
	MessageTag_BBB_XCV_ISCGGIVEN                               MessageTag = "BBB_XCV_ISCGGIVEN"
	MessageTag_BBB_XCV_ISCGGIVEN_ANS                           MessageTag = "BBB_XCV_ISCGGIVEN_ANS"
	MessageTag_BBB_XCV_ISCGPSEUDO                              MessageTag = "BBB_XCV_ISCGPSEUDO"
	MessageTag_BBB_XCV_ISCGPSEUDO_ANS                          MessageTag = "BBB_XCV_ISCGPSEUDO_ANS"
	MessageTag_BBB_XCV_ISCGCOMMONN                             MessageTag = "BBB_XCV_ISCGCOMMONN"
	MessageTag_BBB_XCV_ISCGCOMMONN_ANS                         MessageTag = "BBB_XCV_ISCGCOMMONN_ANS"
	MessageTag_BBB_XCV_ISCGTITLE                               MessageTag = "BBB_XCV_ISCGTITLE"
	MessageTag_BBB_XCV_ISCGTITLE_ANS                           MessageTag = "BBB_XCV_ISCGTITLE_ANS"
	MessageTag_BBB_XCV_ISCGEMAIL                               MessageTag = "BBB_XCV_ISCGEMAIL"
	MessageTag_BBB_XCV_ISCGEMAIL_ANS                           MessageTag = "BBB_XCV_ISCGEMAIL_ANS"
	MessageTag_BBB_XCV_ISSSC                                   MessageTag = "BBB_XCV_ISSSC"
	MessageTag_BBB_XCV_ISSSC_ANS                               MessageTag = "BBB_XCV_ISSSC_ANS"
	MessageTag_BBB_XCV_ISNSSC                                  MessageTag = "BBB_XCV_ISNSSC"
	MessageTag_BBB_XCV_ISNSSC_ANS                              MessageTag = "BBB_XCV_ISNSSC_ANS"
	MessageTag_BBB_XCV_IRDC                                    MessageTag = "BBB_XCV_IRDC"
	MessageTag_BBB_XCV_IRDC_ANS                                MessageTag = "BBB_XCV_IRDC_ANS"
	MessageTag_XCV_TSL_ESP                                     MessageTag = "XCV_TSL_ESP"
	MessageTag_XCV_TSL_ESP_ANS                                 MessageTag = "XCV_TSL_ESP_ANS"
	MessageTag_XCV_TSL_ESP_SIG_ANS                             MessageTag = "XCV_TSL_ESP_SIG_ANS"
	MessageTag_XCV_TSL_ESP_TSP_ANS                             MessageTag = "XCV_TSL_ESP_TSP_ANS"
	MessageTag_XCV_TSL_ESP_REV_ANS                             MessageTag = "XCV_TSL_ESP_REV_ANS"
	MessageTag_XCV_TSL_ETIP                                    MessageTag = "XCV_TSL_ETIP"
	MessageTag_XCV_TSL_ETIP_ANS                                MessageTag = "XCV_TSL_ETIP_ANS"
	MessageTag_XCV_TSL_ETIP_SIG_ANS                            MessageTag = "XCV_TSL_ETIP_SIG_ANS"
	MessageTag_XCV_TSL_ETIP_TSP_ANS                            MessageTag = "XCV_TSL_ETIP_TSP_ANS"
	MessageTag_XCV_TSL_ETIP_REV_ANS                            MessageTag = "XCV_TSL_ETIP_REV_ANS"
	MessageTag_PCV_IVTSC                                       MessageTag = "PCV_IVTSC"
	MessageTag_PCV_IVTSC_ANS                                   MessageTag = "PCV_IVTSC_ANS"
	MessageTag_PCV_ICCSVTSF                                    MessageTag = "PCV_ICCSVTSF"
	MessageTag_PCV_ICCSVTSF_ANS                                MessageTag = "PCV_ICCSVTSF_ANS"
	MessageTag_PSV_IPCVA                                       MessageTag = "PSV_IPCVA"
	MessageTag_PSV_IPCVA_ANS                                   MessageTag = "PSV_IPCVA_ANS"
	MessageTag_PSV_IPCVC                                       MessageTag = "PSV_IPCVC"
	MessageTag_PSV_IPCVC_ANS                                   MessageTag = "PSV_IPCVC_ANS"
	MessageTag_PSV_IPSVC                                       MessageTag = "PSV_IPSVC"
	MessageTag_PSV_IPSVC_ANS                                   MessageTag = "PSV_IPSVC_ANS"
	MessageTag_PSV_IPTVC                                       MessageTag = "PSV_IPTVC"
	MessageTag_PSV_IPTVC_ANS                                   MessageTag = "PSV_IPTVC_ANS"
	MessageTag_PSV_ITPOCOBCT                                   MessageTag = "PSV_ITPOCOBCT"
	MessageTag_PSV_ITPOSVAOBCT                                 MessageTag = "PSV_ITPOSVAOBCT"
	MessageTag_PSV_ITPOSVAOBCT_ANS                             MessageTag = "PSV_ITPOSVAOBCT_ANS"
	MessageTag_PSV_ITPORDAOBCT                                 MessageTag = "PSV_ITPORDAOBCT"
	MessageTag_PSV_ITPOOBCT_ANS                                MessageTag = "PSV_ITPOOBCT_ANS"
	MessageTag_PSV_ITPRISCNARTCAC                              MessageTag = "PSV_ITPRISCNARTCAC"
	MessageTag_PSV_ITPRISCNARTCAC_ANS                          MessageTag = "PSV_ITPRISCNARTCAC_ANS"
	MessageTag_PSV_ICRDIT                                      MessageTag = "PSV_ICRDIT"
	MessageTag_PSV_ICRDIT_ANS                                  MessageTag = "PSV_ICRDIT_ANS"
	MessageTag_PSV_IPCRIAIDBEDC                                MessageTag = "PSV_IPCRIAIDBEDC"
	MessageTag_PSV_IPCRIAIDBEDC_ANS                            MessageTag = "PSV_IPCRIAIDBEDC_ANS"
	MessageTag_PSV_ICTD                                        MessageTag = "PSV_ICTD"
	MessageTag_PSV_ICTD_ANS                                    MessageTag = "PSV_ICTD_ANS"
	MessageTag_PSV_ISDDTA                                      MessageTag = "PSV_ISDDTA"
	MessageTag_PSV_ISDDTA_ANS                                  MessageTag = "PSV_ISDDTA_ANS"
	MessageTag_PSV_HRDBIBCT                                    MessageTag = "PSV_HRDBIBCT"
	MessageTag_PSV_HRDBIBCT_ANS                                MessageTag = "PSV_HRDBIBCT_ANS"
	MessageTag_PSV_DIURDSCHPVR                                 MessageTag = "PSV_DIURDSCHPVR"
	MessageTag_PSV_DIURDSCHPVR_ANS                             MessageTag = "PSV_DIURDSCHPVR_ANS"
	MessageTag_TSV_ASTPTCT                                     MessageTag = "TSV_ASTPTCT"
	MessageTag_TSV_ASTPTCT_ANS                                 MessageTag = "TSV_ASTPTCT_ANS"
	MessageTag_TSV_IBSTAIDOSC                                  MessageTag = "TSV_IBSTAIDOSC"
	MessageTag_TSV_IBSTAIDOSC_ANS                              MessageTag = "TSV_IBSTAIDOSC_ANS"
	MessageTag_TSV_IBSTBCEC                                    MessageTag = "TSV_IBSTBCEC"
	MessageTag_TSV_IBSTBCEC_ANS                                MessageTag = "TSV_IBSTBCEC_ANS"
	MessageTag_TSV_ISCNVABST                                   MessageTag = "TSV_ISCNVABST"
	MessageTag_TSV_ISCNVABST_ANS                               MessageTag = "TSV_ISCNVABST_ANS"
	MessageTag_ADEST_IRTPTBST                                  MessageTag = "ADEST_IRTPTBST"
	MessageTag_ADEST_IRTPTBST_ANS                              MessageTag = "ADEST_IRTPTBST_ANS"
	MessageTag_ADEST_ISTPTBST                                  MessageTag = "ADEST_ISTPTBST"
	MessageTag_ADEST_ISTPTBST_ANS                              MessageTag = "ADEST_ISTPTBST_ANS"
	MessageTag_ADEST_VFDTAOCST_ANS                             MessageTag = "ADEST_VFDTAOCST_ANS"
	MessageTag_ADEST_ISTPTDABST                                MessageTag = "ADEST_ISTPTDABST"
	MessageTag_ADEST_ISTPTDABST_ANS                            MessageTag = "ADEST_ISTPTDABST_ANS"
	MessageTag_ADEST_IBSVPSC                                   MessageTag = "ADEST_IBSVPSC"
	MessageTag_ADEST_IBSVPSC_ANS                               MessageTag = "ADEST_IBSVPSC_ANS"
	MessageTag_ADEST_IBSVPTC                                   MessageTag = "ADEST_IBSVPTC"
	MessageTag_ADEST_IBSVPTC_ANS                               MessageTag = "ADEST_IBSVPTC_ANS"
	MessageTag_ADEST_IBSVPTADC                                 MessageTag = "ADEST_IBSVPTADC"
	MessageTag_ADEST_IBSVPTADC_ANS                             MessageTag = "ADEST_IBSVPTADC_ANS"
	MessageTag_ADEST_IRERVPC                                   MessageTag = "ADEST_IRERVPC"
	MessageTag_ADEST_IRERVPC_ANS                               MessageTag = "ADEST_IRERVPC_ANS"
	MessageTag_EAA_CERT_LOTE_REACHED                           MessageTag = "EAA_CERT_LOTE_REACHED"
	MessageTag_EAA_CERT_LOTE_REACHED_ANS                       MessageTag = "EAA_CERT_LOTE_REACHED_ANS"
	MessageTag_EAA_CERT_TRUST_ANCHOR_LIST_REACHED              MessageTag = "EAA_CERT_TRUST_ANCHOR_LIST_REACHED"
	MessageTag_EAA_CERT_TRUST_ANCHOR_LIST_REACHED_ANS          MessageTag = "EAA_CERT_TRUST_ANCHOR_LIST_REACHED_ANS"
	MessageTag_EAA_DPEAAP                                      MessageTag = "EAA_DPEAAP"
	MessageTag_EAA_DPEAAP_ANS                                  MessageTag = "EAA_DPEAAP_ANS"
	MessageTag_EAA_DLEEAAP                                     MessageTag = "EAA_DLEEAAP"
	MessageTag_EAA_DLEEAAP_ANS                                 MessageTag = "EAA_DLEEAAP_ANS"
	MessageTag_EAA_KBRC                                        MessageTag = "EAA_KBRC"
	MessageTag_EAA_KBRC_ANS                                    MessageTag = "EAA_KBRC_ANS"
	MessageTag_EAA_KBSP                                        MessageTag = "EAA_KBSP"
	MessageTag_EAA_KBSP_ANS                                    MessageTag = "EAA_KBSP_ANS"
	MessageTag_EAA_CLAIMS                                      MessageTag = "EAA_CLAIMS"
	MessageTag_EAA_CLAIMS_ANS                                  MessageTag = "EAA_CLAIMS_ANS"
	MessageTag_EAA_CLAIMS_INFO                                 MessageTag = "EAA_CLAIMS_INFO"
	MessageTag_EAA_SUPPORTED_CLAIMS                            MessageTag = "EAA_SUPPORTED_CLAIMS"
	MessageTag_EAA_SUPPORTED_CLAIMS_ANS                        MessageTag = "EAA_SUPPORTED_CLAIMS_ANS"
	MessageTag_EAA_UNSUPPORTED_CLAIMS                          MessageTag = "EAA_UNSUPPORTED_CLAIMS"
	MessageTag_EAA_ACCEPTABLE_TYPE                             MessageTag = "EAA_ACCEPTABLE_TYPE"
	MessageTag_EAA_ACCEPTABLE_TYPE_ANS                         MessageTag = "EAA_ACCEPTABLE_TYPE_ANS"
	MessageTag_EAA_IDENTIFIER_PRESENT                          MessageTag = "EAA_IDENTIFIER_PRESENT"
	MessageTag_EAA_IDENTIFIER_PRESENT_ANS                      MessageTag = "EAA_IDENTIFIER_PRESENT_ANS"
	MessageTag_EAA_ISSUANCE_DATE_PRESENT                       MessageTag = "EAA_ISSUANCE_DATE_PRESENT"
	MessageTag_EAA_ISSUANCE_DATE_PRESENT_ANS                   MessageTag = "EAA_ISSUANCE_DATE_PRESENT_ANS"
	MessageTag_EAA_NBF_PRESENT                                 MessageTag = "EAA_NBF_PRESENT"
	MessageTag_EAA_NBF_PRESENT_ANS                             MessageTag = "EAA_NBF_PRESENT_ANS"
	MessageTag_EAA_EXP_PRESENT                                 MessageTag = "EAA_EXP_PRESENT"
	MessageTag_EAA_EXP_PRESENT_ANS                             MessageTag = "EAA_EXP_PRESENT_ANS"
	MessageTag_EAA_AID_PRESENT                                 MessageTag = "EAA_AID_PRESENT"
	MessageTag_EAA_AID_PRESENT_ANS                             MessageTag = "EAA_AID_PRESENT_ANS"
	MessageTag_EAA_AED_PRESENT                                 MessageTag = "EAA_AED_PRESENT"
	MessageTag_EAA_AED_PRESENT_ANS                             MessageTag = "EAA_AED_PRESENT_ANS"
	MessageTag_EAA_SIG_PRESENT                                 MessageTag = "EAA_SIG_PRESENT"
	MessageTag_EAA_SIG_PRESENT_ANS                             MessageTag = "EAA_SIG_PRESENT_ANS"
	MessageTag_EAA_SIG_QUAL                                    MessageTag = "EAA_SIG_QUAL"
	MessageTag_EAA_SIG_QUAL_ANS                                MessageTag = "EAA_SIG_QUAL_ANS"
	MessageTag_EAA_CAT_EAA                                     MessageTag = "EAA_CAT_EAA"
	MessageTag_EAA_CAT_EAA_ANS_1                               MessageTag = "EAA_CAT_EAA_ANS_1"
	MessageTag_EAA_CAT_EAA_ANS_2                               MessageTag = "EAA_CAT_EAA_ANS_2"
	MessageTag_EAA_CAT_PUBEAA                                  MessageTag = "EAA_CAT_PUBEAA"
	MessageTag_EAA_CAT_PUBEAA_ANS                              MessageTag = "EAA_CAT_PUBEAA_ANS"
	MessageTag_EAA_CAT_QEAA                                    MessageTag = "EAA_CAT_QEAA"
	MessageTag_EAA_CAT_QEAA_ANS                                MessageTag = "EAA_CAT_QEAA_ANS"
	MessageTag_EAA_QC_PSB                                      MessageTag = "EAA_QC_PSB"
	MessageTag_EAA_QC_PSB_ANS                                  MessageTag = "EAA_QC_PSB_ANS"
	MessageTag_EAA_QUAL_CONCLUSIVE                             MessageTag = "EAA_QUAL_CONCLUSIVE"
	MessageTag_EAA_QUAL_CONCLUSIVE_ANS                         MessageTag = "EAA_QUAL_CONCLUSIVE_ANS"
	MessageTag_EAA_ETSI194721                                  MessageTag = "EAA_ETSI194721"
	MessageTag_EAA_ETSI194721_ANS                              MessageTag = "EAA_ETSI194721_ANS"
	MessageTag_EAA_VT_ITVR                                     MessageTag = "EAA_VT_ITVR"
	MessageTag_EAA_VT_ITVR_ANS                                 MessageTag = "EAA_VT_ITVR_ANS"
	MessageTag_EAA_VT_ITVR_VALIDITY                            MessageTag = "EAA_VT_ITVR_VALIDITY"
	MessageTag_EAA_NOW_BEFORE_NBF                              MessageTag = "EAA_NOW_BEFORE_NBF"
	MessageTag_EAA_NOW_AFTER_EXP                               MessageTag = "EAA_NOW_AFTER_EXP"
	MessageTag_EAA_VT_IAVR                                     MessageTag = "EAA_VT_IAVR"
	MessageTag_EAA_VT_IAVR_ANS                                 MessageTag = "EAA_VT_IAVR_ANS"
	MessageTag_EAA_VT_IAVR_VALIDITY                            MessageTag = "EAA_VT_IAVR_VALIDITY"
	MessageTag_EAA_NOW_BEFORE_ADI                              MessageTag = "EAA_NOW_BEFORE_ADI"
	MessageTag_EAA_NOW_AFTER_ADE                               MessageTag = "EAA_NOW_AFTER_ADE"
	MessageTag_EAA_AD_SDJWT_CONFORMANCE                        MessageTag = "EAA_AD_SDJWT_CONFORMANCE"
	MessageTag_EAA_SHORT_LIVED_STATUS_PRESENT                  MessageTag = "EAA_SHORT_LIVED_STATUS_PRESENT"
	MessageTag_EAA_MANDATORY_STATUS_ABSENT                     MessageTag = "EAA_MANDATORY_STATUS_ABSENT"
	MessageTag_EAA_REV_SDJWT_CONFORMANCE                       MessageTag = "EAA_REV_SDJWT_CONFORMANCE"
	MessageTag_EAA_MDOC_ISSUING_AUTHORITY                      MessageTag = "EAA_MDOC_ISSUING_AUTHORITY"
	MessageTag_EAA_SDJWT_ISSUING_AUTHORITY                     MessageTag = "EAA_SDJWT_ISSUING_AUTHORITY"
	MessageTag_EAA_MDOC_DOCUMENT_NUMBER_ABSENT                 MessageTag = "EAA_MDOC_DOCUMENT_NUMBER_ABSENT"
	MessageTag_EAA_SUB                                         MessageTag = "EAA_SUB"
	MessageTag_EAA_SUB_ANS                                     MessageTag = "EAA_SUB_ANS"
	MessageTag_EAA_SUB_PSE                                     MessageTag = "EAA_SUB_PSE"
	MessageTag_EAA_SUB_PSE_ANS                                 MessageTag = "EAA_SUB_PSE_ANS"
	MessageTag_EAA_CAT                                         MessageTag = "EAA_CAT"
	MessageTag_EAA_CAT_ANS                                     MessageTag = "EAA_CAT_ANS"
	MessageTag_EAA_ISS_COUN                                    MessageTag = "EAA_ISS_COUN"
	MessageTag_EAA_ISS_COUN_ANS                                MessageTag = "EAA_ISS_COUN_ANS"
	MessageTag_EAA_ISS_AUTH                                    MessageTag = "EAA_ISS_AUTH"
	MessageTag_EAA_ISS_AUTH_ANS                                MessageTag = "EAA_ISS_AUTH_ANS"
	MessageTag_EAA_ISS_REG_ID                                  MessageTag = "EAA_ISS_REG_ID"
	MessageTag_EAA_ISS_REG_ID_ANS                              MessageTag = "EAA_ISS_REG_ID_ANS"
	MessageTag_EAA_REV_PR                                      MessageTag = "EAA_REV_PR"
	MessageTag_EAA_REV_PR_ANS                                  MessageTag = "EAA_REV_PR_ANS"
	MessageTag_EAA_REV_AV                                      MessageTag = "EAA_REV_AV"
	MessageTag_EAA_REV_AV_ANS                                  MessageTag = "EAA_REV_AV_ANS"
	MessageTag_EAA_REV_ACC                                     MessageTag = "EAA_REV_ACC"
	MessageTag_EAA_REV_ACC_ANS                                 MessageTag = "EAA_REV_ACC_ANS"
	MessageTag_EAA_REV_ACC_FND                                 MessageTag = "EAA_REV_ACC_FND"
	MessageTag_EAA_REV_ACC_FND_ANS                             MessageTag = "EAA_REV_ACC_FND_ANS"
	MessageTag_EAA_REV_NOT_REV                                 MessageTag = "EAA_REV_NOT_REV"
	MessageTag_EAA_REV_NOT_REV_ANS                             MessageTag = "EAA_REV_NOT_REV_ANS"
	MessageTag_EAA_REV_NOT_ON_HOLD                             MessageTag = "EAA_REV_NOT_ON_HOLD"
	MessageTag_EAA_REV_NOT_ON_HOLD_ANS                         MessageTag = "EAA_REV_NOT_ON_HOLD_ANS"
	MessageTag_EAA_SH_LVD                                      MessageTag = "EAA_SH_LVD"
	MessageTag_EAA_SH_LVD_ANS                                  MessageTag = "EAA_SH_LVD_ANS"
	MessageTag_EAA_OTU                                         MessageTag = "EAA_OTU"
	MessageTag_EAA_OTU_ANS                                     MessageTag = "EAA_OTU_ANS"
	MessageTag_EAA_PSEUDO_USED                                 MessageTag = "EAA_PSEUDO_USED"
	MessageTag_EAA_PSEUDO_USED_ANS                             MessageTag = "EAA_PSEUDO_USED_ANS"
	MessageTag_SDJWT_EAA_VCT_PRESENT                           MessageTag = "SDJWT_EAA_VCT_PRESENT"
	MessageTag_SDJWT_EAA_VCT_PRESENT_ANS                       MessageTag = "SDJWT_EAA_VCT_PRESENT_ANS"
	MessageTag_SDJWT_EAA_VCT_INT_PRESENT                       MessageTag = "SDJWT_EAA_VCT_INT_PRESENT"
	MessageTag_SDJWT_EAA_VCT_INT_PRESENT_ANS                   MessageTag = "SDJWT_EAA_VCT_INT_PRESENT_ANS"
	MessageTag_EAA_REV_TYPE                                    MessageTag = "EAA_REV_TYPE"
	MessageTag_EAA_REV_TYPE_ANS                                MessageTag = "EAA_REV_TYPE_ANS"
	MessageTag_EAA_REV_KNOWN                                   MessageTag = "EAA_REV_KNOWN"
	MessageTag_EAA_REV_KNOWN_ANS                               MessageTag = "EAA_REV_KNOWN_ANS"
	MessageTag_EAA_REV_ISS                                     MessageTag = "EAA_REV_ISS"
	MessageTag_EAA_REV_ISS_ANS                                 MessageTag = "EAA_REV_ISS_ANS"
	MessageTag_EAA_REV_EXP                                     MessageTag = "EAA_REV_EXP"
	MessageTag_EAA_REV_EXP_ANS                                 MessageTag = "EAA_REV_EXP_ANS"
	MessageTag_EAA_REV_NOT_EXP                                 MessageTag = "EAA_REV_NOT_EXP"
	MessageTag_EAA_REV_NOT_EXP_ANS                             MessageTag = "EAA_REV_NOT_EXP_ANS"
	MessageTag_EAA_REV_SUB                                     MessageTag = "EAA_REV_SUB"
	MessageTag_EAA_REV_SUB_ANS                                 MessageTag = "EAA_REV_SUB_ANS"
	MessageTag_EAA_REV_SUB_MATCH                               MessageTag = "EAA_REV_SUB_MATCH"
	MessageTag_EAA_REV_SUB_MATCH_ANS                           MessageTag = "EAA_REV_SUB_MATCH_ANS"
	MessageTag_EAA_REV_ISS_VALID                               MessageTag = "EAA_REV_ISS_VALID"
	MessageTag_EAA_REV_ISS_VALID_ANS                           MessageTag = "EAA_REV_ISS_VALID_ANS"
	MessageTag_EAA_REV_TIME                                    MessageTag = "EAA_REV_TIME"
	MessageTag_EAA_REV_ISS_CERT                                MessageTag = "EAA_REV_ISS_CERT"
	MessageTag_QUAL_TL_EXP                                     MessageTag = "QUAL_TL_EXP"
	MessageTag_QUAL_TL_EXP_ANS                                 MessageTag = "QUAL_TL_EXP_ANS"
	MessageTag_QUAL_TL_FRESH                                   MessageTag = "QUAL_TL_FRESH"
	MessageTag_QUAL_TL_FRESH_ANS                               MessageTag = "QUAL_TL_FRESH_ANS"
	MessageTag_QUAL_TL_VERSION                                 MessageTag = "QUAL_TL_VERSION"
	MessageTag_QUAL_TL_VERSION_ANS                             MessageTag = "QUAL_TL_VERSION_ANS"
	MessageTag_QUAL_TL_WS                                      MessageTag = "QUAL_TL_WS"
	MessageTag_QUAL_TL_WS_ANS                                  MessageTag = "QUAL_TL_WS_ANS"
	MessageTag_QUAL_TL_SV                                      MessageTag = "QUAL_TL_SV"
	MessageTag_QUAL_TL_SV_ANS                                  MessageTag = "QUAL_TL_SV_ANS"
	MessageTag_QUAL_TL_IMRA                                    MessageTag = "QUAL_TL_IMRA"
	MessageTag_QUAL_TL_IMRA_ANS                                MessageTag = "QUAL_TL_IMRA_ANS"
	MessageTag_QUAL_TL_IMRA_ANS_V1                             MessageTag = "QUAL_TL_IMRA_ANS_V1"
	MessageTag_QUAL_TL_IMRA_ANS_V2                             MessageTag = "QUAL_TL_IMRA_ANS_V2"
	MessageTag_QUAL_TL_SERV_CONS                               MessageTag = "QUAL_TL_SERV_CONS"
	MessageTag_QUAL_TL_SERV_CONS_ANS0                          MessageTag = "QUAL_TL_SERV_CONS_ANS0"
	MessageTag_QUAL_TL_SERV_CONS_ANS1                          MessageTag = "QUAL_TL_SERV_CONS_ANS1"
	MessageTag_QUAL_TL_SERV_CONS_ANS2                          MessageTag = "QUAL_TL_SERV_CONS_ANS2"
	MessageTag_QUAL_TL_SERV_CONS_ANS3                          MessageTag = "QUAL_TL_SERV_CONS_ANS3"
	MessageTag_QUAL_TL_SERV_CONS_ANS3A                         MessageTag = "QUAL_TL_SERV_CONS_ANS3A"
	MessageTag_QUAL_TL_SERV_CONS_ANS3B                         MessageTag = "QUAL_TL_SERV_CONS_ANS3B"
	MessageTag_QUAL_TL_SERV_CONS_ANS3C                         MessageTag = "QUAL_TL_SERV_CONS_ANS3C"
	MessageTag_QUAL_TL_SERV_CONS_ANS4                          MessageTag = "QUAL_TL_SERV_CONS_ANS4"
	MessageTag_QUAL_TL_SERV_CONS_ANS5                          MessageTag = "QUAL_TL_SERV_CONS_ANS5"
	MessageTag_QUAL_TL_SERV_CONS_ANS6                          MessageTag = "QUAL_TL_SERV_CONS_ANS6"
	MessageTag_QUAL_TL_SERV_CONS_ANS7                          MessageTag = "QUAL_TL_SERV_CONS_ANS7"
	MessageTag_QUAL_CERT_TRUSTED_LIST_REACHED                  MessageTag = "QUAL_CERT_TRUSTED_LIST_REACHED"
	MessageTag_QUAL_CERT_TRUSTED_LIST_REACHED_ANS              MessageTag = "QUAL_CERT_TRUSTED_LIST_REACHED_ANS"
	MessageTag_QUAL_TRUSTED_LIST_ACCEPT                        MessageTag = "QUAL_TRUSTED_LIST_ACCEPT"
	MessageTag_QUAL_TRUSTED_LIST_ACCEPT_ANS                    MessageTag = "QUAL_TRUSTED_LIST_ACCEPT_ANS"
	MessageTag_QUAL_LIST_OF_TRUSTED_LISTS_ACCEPT               MessageTag = "QUAL_LIST_OF_TRUSTED_LISTS_ACCEPT"
	MessageTag_QUAL_LIST_OF_TRUSTED_LISTS_ACCEPT_ANS           MessageTag = "QUAL_LIST_OF_TRUSTED_LISTS_ACCEPT_ANS"
	MessageTag_QUAL_VALID_TRUSTED_LIST_PRESENT                 MessageTag = "QUAL_VALID_TRUSTED_LIST_PRESENT"
	MessageTag_QUAL_VALID_TRUSTED_LIST_PRESENT_ANS             MessageTag = "QUAL_VALID_TRUSTED_LIST_PRESENT_ANS"
	MessageTag_QUAL_CERT_TYPE_AT_ST                            MessageTag = "QUAL_CERT_TYPE_AT_ST"
	MessageTag_QUAL_CERT_TYPE_AT_ST_ANS                        MessageTag = "QUAL_CERT_TYPE_AT_ST_ANS"
	MessageTag_QUAL_CERT_TYPE_AT_CC                            MessageTag = "QUAL_CERT_TYPE_AT_CC"
	MessageTag_QUAL_CERT_TYPE_AT_CC_ANS                        MessageTag = "QUAL_CERT_TYPE_AT_CC_ANS"
	MessageTag_QUAL_CERT_TYPE_AT_VT                            MessageTag = "QUAL_CERT_TYPE_AT_VT"
	MessageTag_QUAL_CERT_TYPE_AT_VT_ANS                        MessageTag = "QUAL_CERT_TYPE_AT_VT_ANS"
	MessageTag_QUAL_QC_AT_ST                                   MessageTag = "QUAL_QC_AT_ST"
	MessageTag_QUAL_QC_AT_ST_ANS                               MessageTag = "QUAL_QC_AT_ST_ANS"
	MessageTag_QUAL_QC_AT_CC                                   MessageTag = "QUAL_QC_AT_CC"
	MessageTag_QUAL_QC_AT_CC_ANS                               MessageTag = "QUAL_QC_AT_CC_ANS"
	MessageTag_QUAL_QC_AT_VT                                   MessageTag = "QUAL_QC_AT_VT"
	MessageTag_QUAL_QC_AT_VT_ANS                               MessageTag = "QUAL_QC_AT_VT_ANS"
	MessageTag_QUAL_QSCD_AT_ST                                 MessageTag = "QUAL_QSCD_AT_ST"
	MessageTag_QUAL_QSCD_AT_ST_ANS                             MessageTag = "QUAL_QSCD_AT_ST_ANS"
	MessageTag_QUAL_QSCD_AT_CC                                 MessageTag = "QUAL_QSCD_AT_CC"
	MessageTag_QUAL_QSCD_AT_CC_ANS                             MessageTag = "QUAL_QSCD_AT_CC_ANS"
	MessageTag_QUAL_QSCD_AT_VT                                 MessageTag = "QUAL_QSCD_AT_VT"
	MessageTag_QUAL_QSCD_AT_VT_ANS                             MessageTag = "QUAL_QSCD_AT_VT_ANS"
	MessageTag_QUAL_UNIQUE_CERT                                MessageTag = "QUAL_UNIQUE_CERT"
	MessageTag_QUAL_UNIQUE_CERT_ANS                            MessageTag = "QUAL_UNIQUE_CERT_ANS"
	MessageTag_QUAL_IS_ADES                                    MessageTag = "QUAL_IS_ADES"
	MessageTag_QUAL_IS_ADES_IND                                MessageTag = "QUAL_IS_ADES_IND"
	MessageTag_QUAL_IS_ADES_INV                                MessageTag = "QUAL_IS_ADES_INV"
	MessageTag_QUAL_HAS_METS                                   MessageTag = "QUAL_HAS_METS"
	MessageTag_QUAL_HAS_METS_ANS                               MessageTag = "QUAL_HAS_METS_ANS"
	MessageTag_QUAL_HAS_METS_ATTIME                            MessageTag = "QUAL_HAS_METS_ATTIME"
	MessageTag_QUAL_HAS_METS_ATTIME_ANS                        MessageTag = "QUAL_HAS_METS_ATTIME_ANS"
	MessageTag_QUAL_HAS_METS_HCCECBA                           MessageTag = "QUAL_HAS_METS_HCCECBA"
	MessageTag_QUAL_HAS_METS_HCCECBA_ANS                       MessageTag = "QUAL_HAS_METS_HCCECBA_ANS"
	MessageTag_QUAL_HAS_METS_HCCECBA_ANS_2                     MessageTag = "QUAL_HAS_METS_HCCECBA_ANS_2"
	MessageTag_QUAL_HAS_METS_HCCECBA_ANS_3                     MessageTag = "QUAL_HAS_METS_HCCECBA_ANS_3"
	MessageTag_QUAL_HAS_CAQC                                   MessageTag = "QUAL_HAS_CAQC"
	MessageTag_QUAL_HAS_CAQC_ANS                               MessageTag = "QUAL_HAS_CAQC_ANS"
	MessageTag_QUAL_HAS_CAQC_ANS_2                             MessageTag = "QUAL_HAS_CAQC_ANS_2"
	MessageTag_QUAL_HAS_ATTIME                                 MessageTag = "QUAL_HAS_ATTIME"
	MessageTag_QUAL_HAS_ATTIME_ANS                             MessageTag = "QUAL_HAS_ATTIME_ANS"
	MessageTag_QUAL_HAS_TS_CERT_TYPE                           MessageTag = "QUAL_HAS_TS_CERT_TYPE"
	MessageTag_QUAL_HAS_TS_CERT_TYPE_ANS                       MessageTag = "QUAL_HAS_TS_CERT_TYPE_ANS"
	MessageTag_QUAL_HAS_CONF                                   MessageTag = "QUAL_HAS_CONF"
	MessageTag_QUAL_HAS_CONF_ANS                               MessageTag = "QUAL_HAS_CONF_ANS"
	MessageTag_QUAL_HAS_QEAA                                   MessageTag = "QUAL_HAS_QEAA"
	MessageTag_QUAL_HAS_QEAA_ANS                               MessageTag = "QUAL_HAS_QEAA_ANS"
	MessageTag_QUAL_HAS_QTST                                   MessageTag = "QUAL_HAS_QTST"
	MessageTag_QUAL_HAS_QTST_ANS                               MessageTag = "QUAL_HAS_QTST_ANS"
	MessageTag_QUAL_IS_TRUST_CERT_MATCH_SERVICE                MessageTag = "QUAL_IS_TRUST_CERT_MATCH_SERVICE"
	MessageTag_QUAL_IS_TRUST_CERT_MATCH_SERVICE_ANS0           MessageTag = "QUAL_IS_TRUST_CERT_MATCH_SERVICE_ANS0"
	MessageTag_QUAL_IS_TRUST_CERT_MATCH_SERVICE_ANS1           MessageTag = "QUAL_IS_TRUST_CERT_MATCH_SERVICE_ANS1"
	MessageTag_QUAL_IS_TRUST_CERT_MATCH_SERVICE_ANS2           MessageTag = "QUAL_IS_TRUST_CERT_MATCH_SERVICE_ANS2"
	MessageTag_QUAL_HAS_GRANTED                                MessageTag = "QUAL_HAS_GRANTED"
	MessageTag_QUAL_HAS_GRANTED_ANS                            MessageTag = "QUAL_HAS_GRANTED_ANS"
	MessageTag_QUAL_HAS_GRANTED_ANS_2                          MessageTag = "QUAL_HAS_GRANTED_ANS_2"
	MessageTag_QUAL_HAS_GRANTED_AT                             MessageTag = "QUAL_HAS_GRANTED_AT"
	MessageTag_QUAL_HAS_GRANTED_AT_ANS                         MessageTag = "QUAL_HAS_GRANTED_AT_ANS"
	MessageTag_QUAL_HAS_CONSISTENT_BY_QC                       MessageTag = "QUAL_HAS_CONSISTENT_BY_QC"
	MessageTag_QUAL_HAS_CONSISTENT_BY_QC_ANS                   MessageTag = "QUAL_HAS_CONSISTENT_BY_QC_ANS"
	MessageTag_QUAL_HAS_CONSISTENT_BY_QSCD                     MessageTag = "QUAL_HAS_CONSISTENT_BY_QSCD"
	MessageTag_QUAL_HAS_CONSISTENT_BY_QSCD_ANS                 MessageTag = "QUAL_HAS_CONSISTENT_BY_QSCD_ANS"
	MessageTag_QUAL_HAS_CERT_TYPE_COVERAGE                     MessageTag = "QUAL_HAS_CERT_TYPE_COVERAGE"
	MessageTag_QUAL_HAS_CERT_TYPE_COVERAGE_ANS                 MessageTag = "QUAL_HAS_CERT_TYPE_COVERAGE_ANS"
	MessageTag_QUAL_HAS_VALID_CAQC                             MessageTag = "QUAL_HAS_VALID_CAQC"
	MessageTag_QUAL_HAS_VALID_CAQC_ANS                         MessageTag = "QUAL_HAS_VALID_CAQC_ANS"
	MessageTag_QUAL_HAS_ONLY_ONE                               MessageTag = "QUAL_HAS_ONLY_ONE"
	MessageTag_QUAL_HAS_ONLY_ONE_ANS                           MessageTag = "QUAL_HAS_ONLY_ONE_ANS"
	MessageTag_QWAC_VALID                                      MessageTag = "QWAC_VALID"
	MessageTag_QWAC_VALID_ANS                                  MessageTag = "QWAC_VALID_ANS"
	MessageTag_QWAC_VALID_ANS_2                                MessageTag = "QWAC_VALID_ANS_2"
	MessageTag_QWAC_CERT_QUAL_CONCLUSIVE                       MessageTag = "QWAC_CERT_QUAL_CONCLUSIVE"
	MessageTag_QWAC_CERT_QUAL_CONCLUSIVE_ANS                   MessageTag = "QWAC_CERT_QUAL_CONCLUSIVE_ANS"
	MessageTag_QWAC_IS_WSA_AT_TIME                             MessageTag = "QWAC_IS_WSA_AT_TIME"
	MessageTag_QWAC_IS_WSA_AT_TIME_ANS                         MessageTag = "QWAC_IS_WSA_AT_TIME_ANS"
	MessageTag_QWAC_CERT_POLICY                                MessageTag = "QWAC_CERT_POLICY"
	MessageTag_QWAC_CERT_POLICY_ANS                            MessageTag = "QWAC_CERT_POLICY_ANS"
	MessageTag_QWAC_VAL_PERIOD                                 MessageTag = "QWAC_VAL_PERIOD"
	MessageTag_QWAC_VAL_PERIOD_ANS                             MessageTag = "QWAC_VAL_PERIOD_ANS"
	MessageTag_QWAC_DOMAIN_NAME                                MessageTag = "QWAC_DOMAIN_NAME"
	MessageTag_QWAC_DOMAIN_NAME_ANS                            MessageTag = "QWAC_DOMAIN_NAME_ANS"
	MessageTag_QWAC2_EXT_KEY_USAGE                             MessageTag = "QWAC2_EXT_KEY_USAGE"
	MessageTag_QWAC2_EXT_KEY_USAGE_ANS                         MessageTag = "QWAC2_EXT_KEY_USAGE_ANS"
	MessageTag_TLS_CERT_BINDING_URL                            MessageTag = "TLS_CERT_BINDING_URL"
	MessageTag_TLS_CERT_BINDING_URL_ANS                        MessageTag = "TLS_CERT_BINDING_URL_ANS"
	MessageTag_TLS_CERT_BINDING_SIG                            MessageTag = "TLS_CERT_BINDING_SIG"
	MessageTag_TLS_CERT_BINDING_SIG_ANS                        MessageTag = "TLS_CERT_BINDING_SIG_ANS"
	MessageTag_TLS_CERT_BINDING_SIG_FORM                       MessageTag = "TLS_CERT_BINDING_SIG_FORM"
	MessageTag_TLS_CERT_BINDING_SIG_FORM_ANS                   MessageTag = "TLS_CERT_BINDING_SIG_FORM_ANS"
	MessageTag_TLS_CERT_BINDING_SIG_SER                        MessageTag = "TLS_CERT_BINDING_SIG_SER"
	MessageTag_TLS_CERT_BINDING_SIG_SER_ANS                    MessageTag = "TLS_CERT_BINDING_SIG_SER_ANS"
	MessageTag_TLS_CERT_BINDING_SIG_EXP                        MessageTag = "TLS_CERT_BINDING_SIG_EXP"
	MessageTag_TLS_CERT_BINDING_SIG_EXP_ANS                    MessageTag = "TLS_CERT_BINDING_SIG_EXP_ANS"
	MessageTag_TLS_CERT_BINDING_SIG_EXPIRY_DATE                MessageTag = "TLS_CERT_BINDING_SIG_EXPIRY_DATE"
	MessageTag_TLS_CERT_BINDING_SIG_EXPIRY_DATE_ANS            MessageTag = "TLS_CERT_BINDING_SIG_EXPIRY_DATE_ANS"
	MessageTag_TLS_CERT_BINDING_QWAC2                          MessageTag = "TLS_CERT_BINDING_QWAC2"
	MessageTag_TLS_CERT_BINDING_QWAC2_ANS                      MessageTag = "TLS_CERT_BINDING_QWAC2_ANS"
	MessageTag_TLS_CERT_BINDING_SIG_VALID                      MessageTag = "TLS_CERT_BINDING_SIG_VALID"
	MessageTag_TLS_CERT_BINDING_SIG_VALID_ANS                  MessageTag = "TLS_CERT_BINDING_SIG_VALID_ANS"
	MessageTag_TLS_CERT_BINDING_CERT_IDENTIFIED                MessageTag = "TLS_CERT_BINDING_CERT_IDENTIFIED"
	MessageTag_TLS_CERT_BINDING_CERT_IDENTIFIED_ANS            MessageTag = "TLS_CERT_BINDING_CERT_IDENTIFIED_ANS"
	MessageTag_CERT_USAGE_LOTE_ACCEPT                          MessageTag = "CERT_USAGE_LOTE_ACCEPT"
	MessageTag_CERT_USAGE_LOTE_ACCEPT_ANS                      MessageTag = "CERT_USAGE_LOTE_ACCEPT_ANS"
	MessageTag_CERT_USAGE_LOLOTE_ACCEPT                        MessageTag = "CERT_USAGE_LOLOTE_ACCEPT"
	MessageTag_CERT_USAGE_LOLOTE_ACCEPT_ANS                    MessageTag = "CERT_USAGE_LOLOTE_ACCEPT_ANS"
	MessageTag_CERT_USAGE_VALID_LOTE_PRESENT                   MessageTag = "CERT_USAGE_VALID_LOTE_PRESENT"
	MessageTag_CERT_USAGE_VALID_LOTE_PRESENT_ANS               MessageTag = "CERT_USAGE_VALID_LOTE_PRESENT_ANS"
	MessageTag_CERT_USAGE_HAS_ATTIME                           MessageTag = "CERT_USAGE_HAS_ATTIME"
	MessageTag_CERT_USAGE_HAS_ATTIME_ANS                       MessageTag = "CERT_USAGE_HAS_ATTIME_ANS"
	MessageTag_CERT_USAGE_LIST_TYPE_KNOWN                      MessageTag = "CERT_USAGE_LIST_TYPE_KNOWN"
	MessageTag_CERT_USAGE_LIST_TYPE_KNOWN_ANS                  MessageTag = "CERT_USAGE_LIST_TYPE_KNOWN_ANS"
	MessageTag_CERT_USAGE_STATUS                               MessageTag = "CERT_USAGE_STATUS"
	MessageTag_CERT_USAGE_STATUS_ANS                           MessageTag = "CERT_USAGE_STATUS_ANS"
	MessageTag_CERT_USAGE_STATUS_CONS                          MessageTag = "CERT_USAGE_STATUS_CONS"
	MessageTag_CERT_USAGE_STATUS_CONS_ANS                      MessageTag = "CERT_USAGE_STATUS_CONS_ANS"
	MessageTag_CERT_USAGE_STATUS_KNOWN                         MessageTag = "CERT_USAGE_STATUS_KNOWN"
	MessageTag_CERT_USAGE_STATUS_KNOWN_ANS                     MessageTag = "CERT_USAGE_STATUS_KNOWN_ANS"
	MessageTag_CERT_USAGE_STI                                  MessageTag = "CERT_USAGE_STI"
	MessageTag_CERT_USAGE_STI_ANS                              MessageTag = "CERT_USAGE_STI_ANS"
	MessageTag_CERT_USAGE_STI_KNOWN                            MessageTag = "CERT_USAGE_STI_KNOWN"
	MessageTag_CERT_USAGE_STI_KNOWN_ANS                        MessageTag = "CERT_USAGE_STI_KNOWN_ANS"
	MessageTag_PID_DOCUMENT_TYPE                               MessageTag = "PID_DOCUMENT_TYPE"
	MessageTag_PID_DOCUMENT_TYPE_ANS                           MessageTag = "PID_DOCUMENT_TYPE_ANS"
	MessageTag_PID_LOTE_TYPE_PID_PROVIDERS                     MessageTag = "PID_LOTE_TYPE_PID_PROVIDERS"
	MessageTag_PID_LOTE_TYPE_PID_PROVIDERS_ANS                 MessageTag = "PID_LOTE_TYPE_PID_PROVIDERS_ANS"
	MessageTag_PID_STI_PID_ISSUANCE                            MessageTag = "PID_STI_PID_ISSUANCE"
	MessageTag_PID_STI_PID_ISSUANCE_ANS                        MessageTag = "PID_STI_PID_ISSUANCE_ANS"
	MessageTag_PID_PROVIDER_AT_ISSUANCE_TIME                   MessageTag = "PID_PROVIDER_AT_ISSUANCE_TIME"
	MessageTag_PID_PROVIDER_AT_ISSUANCE_TIME_ANS               MessageTag = "PID_PROVIDER_AT_ISSUANCE_TIME_ANS"
	MessageTag_PID_PROVIDER_AT_VALIDATION_TIME                 MessageTag = "PID_PROVIDER_AT_VALIDATION_TIME"
	MessageTag_PID_PROVIDER_AT_VALIDATION_TIME_ANS             MessageTag = "PID_PROVIDER_AT_VALIDATION_TIME_ANS"
	MessageTag_BBB_ACCEPT                                      MessageTag = "BBB_ACCEPT"
	MessageTag_BBB_ACCEPT_ANS                                  MessageTag = "BBB_ACCEPT_ANS"
	MessageTag_TST_TYPE_CONTENT_TST                            MessageTag = "TST_TYPE_CONTENT_TST"
	MessageTag_TST_TYPE_SIGNATURE_TST                          MessageTag = "TST_TYPE_SIGNATURE_TST"
	MessageTag_TST_TYPE_VD_TST                                 MessageTag = "TST_TYPE_VD_TST"
	MessageTag_TST_TYPE_DOC_TST                                MessageTag = "TST_TYPE_DOC_TST"
	MessageTag_TST_TYPE_CONTAINER_TST                          MessageTag = "TST_TYPE_CONTAINER_TST"
	MessageTag_TST_TYPE_ARCHIVE_TST                            MessageTag = "TST_TYPE_ARCHIVE_TST"
	MessageTag_TST_TYPE_ER_TST                                 MessageTag = "TST_TYPE_ER_TST"
	MessageTag_TST_TYPE_REF_ER_ATST                            MessageTag = "TST_TYPE_REF_ER_ATST"
	MessageTag_TST_TYPE_REF_ER_ATST_SEQ                        MessageTag = "TST_TYPE_REF_ER_ATST_SEQ"
	MessageTag_EMPTY                                           MessageTag = "EMPTY"
	MessageTag_CERTIFICATE                                     MessageTag = "CERTIFICATE"
	MessageTag_CA_CERTIFICATE                                  MessageTag = "CA_CERTIFICATE"
	MessageTag_SIGNING_CERTIFICATE                             MessageTag = "SIGNING_CERTIFICATE"
	MessageTag_REVOCATION                                      MessageTag = "REVOCATION"
	MessageTag_REVOCATION_SIG_CERT                             MessageTag = "REVOCATION_SIG_CERT"
	MessageTag_REVOCATION_CA_CERT                              MessageTag = "REVOCATION_CA_CERT"
	MessageTag_SIGNATURE                                       MessageTag = "SIGNATURE"
	MessageTag_TIMESTAMP                                       MessageTag = "TIMESTAMP"
	MessageTag_TIMESTAMP_SIG_CERT                              MessageTag = "TIMESTAMP_SIG_CERT"
	MessageTag_TIMESTAMP_CA_CERT                               MessageTag = "TIMESTAMP_CA_CERT"
	MessageTag_EAA_REV                                         MessageTag = "EAA_REV"
	MessageTag_EAA_REV_SIG_CERT                                MessageTag = "EAA_REV_SIG_CERT"
	MessageTag_EAA_REV_CA_CERT                                 MessageTag = "EAA_REV_CA_CERT"
	MessageTag_ACCEPTABLE_REVOCATION                           MessageTag = "ACCEPTABLE_REVOCATION"
	MessageTag_BASIC_SIGNATURE_VALIDATION_RESULT               MessageTag = "BASIC_SIGNATURE_VALIDATION_RESULT"
	MessageTag_BEST_SIGNATURE_TIME_CERT_NOT_AFTER              MessageTag = "BEST_SIGNATURE_TIME_CERT_NOT_AFTER"
	MessageTag_BEST_SIGNATURE_TIME_CERT_NOT_BEFORE             MessageTag = "BEST_SIGNATURE_TIME_CERT_NOT_BEFORE"
	MessageTag_BEST_SIGNATURE_TIME_CERT_REVOCATION             MessageTag = "BEST_SIGNATURE_TIME_CERT_REVOCATION"
	MessageTag_BEST_SIGNATURE_TIME_CERT_SUSPENSION             MessageTag = "BEST_SIGNATURE_TIME_CERT_SUSPENSION"
	MessageTag_CERTIFICATE_ID                                  MessageTag = "CERTIFICATE_ID"
	MessageTag_CERTIFICATE_REVOCATION_FOUND                    MessageTag = "CERTIFICATE_REVOCATION_FOUND"
	MessageTag_CERTIFICATE_REVOCATION_NOT_FOUND                MessageTag = "CERTIFICATE_REVOCATION_NOT_FOUND"
	MessageTag_CERTIFICATE_SUNSET_DATE                         MessageTag = "CERTIFICATE_SUNSET_DATE"
	MessageTag_CERTIFICATE_SUNSET_DATE_TRUST_ANCHOR            MessageTag = "CERTIFICATE_SUNSET_DATE_TRUST_ANCHOR"
	MessageTag_CERTIFICATE_SUNSET_DATE_VALID                   MessageTag = "CERTIFICATE_SUNSET_DATE_VALID"
	MessageTag_CERTIFICATE_TYPE                                MessageTag = "CERTIFICATE_TYPE"
	MessageTag_CERTIFICATE_VALIDITY                            MessageTag = "CERTIFICATE_VALIDITY"
	MessageTag_CERTIFICATE_USAGE                               MessageTag = "CERTIFICATE_USAGE"
	MessageTag_CERTIFICATE_USAGE_LIST_TYPE                     MessageTag = "CERTIFICATE_USAGE_LIST_TYPE"
	MessageTag_CERTIFICATE_USAGE_STATUS                        MessageTag = "CERTIFICATE_USAGE_STATUS"
	MessageTag_CERTIFICATE_USAGE_STATUSES                      MessageTag = "CERTIFICATE_USAGE_STATUSES"
	MessageTag_CERTIFICATE_USAGE_STI                           MessageTag = "CERTIFICATE_USAGE_STI"
	MessageTag_CERTIFICATE_USAGE_STIS                          MessageTag = "CERTIFICATE_USAGE_STIS"
	MessageTag_CONTROL_TIME                                    MessageTag = "CONTROL_TIME"
	MessageTag_CONTROL_TIME_ALONE                              MessageTag = "CONTROL_TIME_ALONE"
	MessageTag_CONTROL_TIME_WITH_POE                           MessageTag = "CONTROL_TIME_WITH_POE"
	MessageTag_CONTROL_TIME_WITH_TRUST_ANCHOR                  MessageTag = "CONTROL_TIME_WITH_TRUST_ANCHOR"
	MessageTag_CRYPTOGRAPHIC_CHECK_FAILURE                     MessageTag = "CRYPTOGRAPHIC_CHECK_FAILURE"
	MessageTag_CRYPTOGRAPHIC_CHECK_FAILURE_WITH_ID             MessageTag = "CRYPTOGRAPHIC_CHECK_FAILURE_WITH_ID"
	MessageTag_CRYPTOGRAPHIC_CHECK_FAILURE_WITH_REF            MessageTag = "CRYPTOGRAPHIC_CHECK_FAILURE_WITH_REF"
	MessageTag_CRYPTOGRAPHIC_CHECK_FAILURE_WITH_REF_WITH_NAME  MessageTag = "CRYPTOGRAPHIC_CHECK_FAILURE_WITH_REF_WITH_NAME"
	MessageTag_CRYPTOGRAPHIC_CHECK_FAILURE_WITH_REF_WITH_NAMES MessageTag = "CRYPTOGRAPHIC_CHECK_FAILURE_WITH_REF_WITH_NAMES"
	MessageTag_CRYPTOGRAPHIC_CHECK_SUCCESS                     MessageTag = "CRYPTOGRAPHIC_CHECK_SUCCESS"
	MessageTag_CRYPTOGRAPHIC_CHECK_SUCCESS_KEY_SIZE            MessageTag = "CRYPTOGRAPHIC_CHECK_SUCCESS_KEY_SIZE"
	MessageTag_CRYPTOGRAPHIC_CHECK_SUCCESS_DM                  MessageTag = "CRYPTOGRAPHIC_CHECK_SUCCESS_DM"
	MessageTag_CRYPTOGRAPHIC_CHECK_SUCCESS_DM_WITH_ID          MessageTag = "CRYPTOGRAPHIC_CHECK_SUCCESS_DM_WITH_ID"
	MessageTag_CRYPTOGRAPHIC_CHECK_SUCCESS_DM_WITH_NAME        MessageTag = "CRYPTOGRAPHIC_CHECK_SUCCESS_DM_WITH_NAME"
	MessageTag_CRYPTOGRAPHIC_CHECK_SUCCESS_DM_WITH_NAMES       MessageTag = "CRYPTOGRAPHIC_CHECK_SUCCESS_DM_WITH_NAMES"
	MessageTag_EVIDENCE_RECORD_VALIDATION                      MessageTag = "EVIDENCE_RECORD_VALIDATION"
	MessageTag_EXTENDED_KEY_USAGE                              MessageTag = "EXTENDED_KEY_USAGE"
	MessageTag_KEY_USAGE                                       MessageTag = "KEY_USAGE"
	MessageTag_LAST_ACCEPTABLE_REVOCATION                      MessageTag = "LAST_ACCEPTABLE_REVOCATION"
	MessageTag_LIST_OF_TRUSTED_ENTITIES                        MessageTag = "LIST_OF_TRUSTED_ENTITIES"
	MessageTag_PSEUDO                                          MessageTag = "PSEUDO"
	MessageTag_QWAC_EXPIRY_EXP                                 MessageTag = "QWAC_EXPIRY_EXP"
	MessageTag_QWAC_EXPIRY_TLS_CERT                            MessageTag = "QWAC_EXPIRY_TLS_CERT"
	MessageTag_QWAC_EXPIRY_SIGN_CERT                           MessageTag = "QWAC_EXPIRY_SIGN_CERT"
	MessageTag_REFERENCE                                       MessageTag = "REFERENCE"
	MessageTag_REFERENCE_NAME_CHECK                            MessageTag = "REFERENCE_NAME_CHECK"
	MessageTag_REFERENCES_WITH_NAMES                           MessageTag = "REFERENCES_WITH_NAMES"
	MessageTag_REVOCATION_ACCEPTANCE_CHECK                     MessageTag = "REVOCATION_ACCEPTANCE_CHECK"
	MessageTag_REVOCATION_CERT_HASH_OK                         MessageTag = "REVOCATION_CERT_HASH_OK"
	MessageTag_REVOCATION_CERT_HASH_OK_ID                      MessageTag = "REVOCATION_CERT_HASH_OK_ID"
	MessageTag_REVOCATION_CERT_VALIDITY                        MessageTag = "REVOCATION_CERT_VALIDITY"
	MessageTag_REVOCATION_CHECK                                MessageTag = "REVOCATION_CHECK"
	MessageTag_REVOCATION_CONSISTENT                           MessageTag = "REVOCATION_CONSISTENT"
	MessageTag_REVOCATION_CONSISTENT_CRL                       MessageTag = "REVOCATION_CONSISTENT_CRL"
	MessageTag_REVOCATION_CONSISTENT_OCSP                      MessageTag = "REVOCATION_CONSISTENT_OCSP"
	MessageTag_REVOCATION_CONSISTENT_TL                        MessageTag = "REVOCATION_CONSISTENT_TL"
	MessageTag_REVOCATION_INFO                                 MessageTag = "REVOCATION_INFO"
	MessageTag_REVOCATION_NOT_AFTER_AFTER                      MessageTag = "REVOCATION_NOT_AFTER_AFTER"
	MessageTag_REVOCATION_NOT_AFTER_AFTER_ID                   MessageTag = "REVOCATION_NOT_AFTER_AFTER_ID"
	MessageTag_REVOCATION_PRODUCED_AT_CERT_VALIDITY            MessageTag = "REVOCATION_PRODUCED_AT_CERT_VALIDITY"
	MessageTag_REVOCATION_PRODUCED_AT_OUT_OF_BOUNDS            MessageTag = "REVOCATION_PRODUCED_AT_OUT_OF_BOUNDS"
	MessageTag_REVOCATION_PRODUCED_AT_OUT_OF_BOUNDS_ID         MessageTag = "REVOCATION_PRODUCED_AT_OUT_OF_BOUNDS_ID"
	MessageTag_REVOCATION_REASON                               MessageTag = "REVOCATION_REASON"
	MessageTag_REVOCATION_THIS_UPDATE_CONTROL_TIME             MessageTag = "REVOCATION_THIS_UPDATE_CONTROL_TIME"
	MessageTag_SIGNATURE_ALGORITHM_WITH_KEY_SIZE               MessageTag = "SIGNATURE_ALGORITHM_WITH_KEY_SIZE"
	MessageTag_SIGNATURE_ID                                    MessageTag = "SIGNATURE_ID"
	MessageTag_STRUCTURAL_VALIDATION_FAILURE                   MessageTag = "STRUCTURAL_VALIDATION_FAILURE"
	MessageTag_TIMESTAMP_AND_CERTIFICATE_NOT_AFTER             MessageTag = "TIMESTAMP_AND_CERTIFICATE_NOT_AFTER"
	MessageTag_TIMESTAMP_AND_CRYPTO_CONSTRAINTS_EXPIRATION     MessageTag = "TIMESTAMP_AND_CRYPTO_CONSTRAINTS_EXPIRATION"
	MessageTag_TIMESTAMP_AND_REVOCATION_TIME                   MessageTag = "TIMESTAMP_AND_REVOCATION_TIME"
	MessageTag_TIMESTAMP_VALIDATION                            MessageTag = "TIMESTAMP_VALIDATION"
	MessageTag_TOKEN_ID                                        MessageTag = "TOKEN_ID"
	MessageTag_TRUST_SERVICE_NAME                              MessageTag = "TRUST_SERVICE_NAME"
	MessageTag_TRUSTED_SERVICE_STATUS                          MessageTag = "TRUSTED_SERVICE_STATUS"
	MessageTag_TRUSTED_SERVICE_TYPE                            MessageTag = "TRUSTED_SERVICE_TYPE"
	MessageTag_TRUSTED_LIST                                    MessageTag = "TRUSTED_LIST"
	MessageTag_VALIDATION_TIME                                 MessageTag = "VALIDATION_TIME"
	MessageTag_CRYPTOGRAPHIC_VERIFICATION                      MessageTag = "CRYPTOGRAPHIC_VERIFICATION"
	MessageTag_FORMAT_CHECKING                                 MessageTag = "FORMAT_CHECKING"
	MessageTag_IDENTIFICATION_OF_THE_SIGNING_CERTIFICATE       MessageTag = "IDENTIFICATION_OF_THE_SIGNING_CERTIFICATE"
	MessageTag_PAST_SIGNATURE_VALIDATION                       MessageTag = "PAST_SIGNATURE_VALIDATION"
	MessageTag_PAST_CERTIFICATE_VALIDATION                     MessageTag = "PAST_CERTIFICATE_VALIDATION"
	MessageTag_REVOCATION_FRESHNESS_CHECKER                    MessageTag = "REVOCATION_FRESHNESS_CHECKER"
	MessageTag_SIGNATURE_ACCEPTANCE_VALIDATION                 MessageTag = "SIGNATURE_ACCEPTANCE_VALIDATION"
	MessageTag_VALIDATION_CONTEXT_INITIALIZATION               MessageTag = "VALIDATION_CONTEXT_INITIALIZATION"
	MessageTag_VALIDATION_TIME_SLIDING                         MessageTag = "VALIDATION_TIME_SLIDING"
	MessageTag_X509_CERTIFICATE_VALIDATION                     MessageTag = "X509_CERTIFICATE_VALIDATION"
	MessageTag_RESULTS                                         MessageTag = "RESULTS"
	MessageTag_EEA_TYPE                                        MessageTag = "EEA_TYPE"
	MessageTag_EAA_ACCEPTANCE_VALIDATION                       MessageTag = "EAA_ACCEPTANCE_VALIDATION"
	MessageTag_AOV                                             MessageTag = "AOV"
	MessageTag_CERT_QUALIFICATION                              MessageTag = "CERT_QUALIFICATION"
	MessageTag_CERT_QUALIFICATION_AT_TIME                      MessageTag = "CERT_QUALIFICATION_AT_TIME"
	MessageTag_CERT_USAGE_AT_TIME                              MessageTag = "CERT_USAGE_AT_TIME"
	MessageTag_CERT_USAGES                                     MessageTag = "CERT_USAGES"
	MessageTag_CC                                              MessageTag = "CC"
	MessageTag_CRS                                             MessageTag = "CRS"
	MessageTag_DAAV                                            MessageTag = "DAAV"
	MessageTag_EAA_QUALIFICATION                               MessageTag = "EAA_QUALIFICATION"
	MessageTag_EAA_QUALIFICATION_PROCESS                       MessageTag = "EAA_QUALIFICATION_PROCESS"
	MessageTag_LOTE                                            MessageTag = "LOTE"
	MessageTag_LOLOTE                                          MessageTag = "LOLOTE"
	MessageTag_LOTL                                            MessageTag = "LOTL"
	MessageTag_PID_QUALIFICATION_PROCESS                       MessageTag = "PID_QUALIFICATION_PROCESS"
	MessageTag_PSV_CRS                                         MessageTag = "PSV_CRS"
	MessageTag_QWAC_VALIDATION                                 MessageTag = "QWAC_VALIDATION"
	MessageTag_QWAC_VALIDATION_PROFILE                         MessageTag = "QWAC_VALIDATION_PROFILE"
	MessageTag_RAC                                             MessageTag = "RAC"
	MessageTag_SIG_QUALIFICATION                               MessageTag = "SIG_QUALIFICATION"
	MessageTag_SUB_XCV                                         MessageTag = "SUB_XCV"
	MessageTag_TL                                              MessageTag = "TL"
	MessageTag_TST_QUALIFICATION                               MessageTag = "TST_QUALIFICATION"
	MessageTag_TST_QUALIFICATION_AT_TIME                       MessageTag = "TST_QUALIFICATION_AT_TIME"
	MessageTag_VPBS                                            MessageTag = "VPBS"
	MessageTag_VPEAA                                           MessageTag = "VPEAA"
	MessageTag_VPER                                            MessageTag = "VPER"
	MessageTag_VPFLTVD                                         MessageTag = "VPFLTVD"
	MessageTag_VPFRVC                                          MessageTag = "VPFRVC"
	MessageTag_VPFSWATSP                                       MessageTag = "VPFSWATSP"
	MessageTag_VPFTSP                                          MessageTag = "VPFTSP"
	MessageTag_VPFTSPWATSP                                     MessageTag = "VPFTSPWATSP"
	MessageTag_VTS_CRS                                         MessageTag = "VTS_CRS"
	MessageTag_VT_BEST_SIGNATURE_TIME                          MessageTag = "VT_BEST_SIGNATURE_TIME"
	MessageTag_VT_CERTIFICATE_ISSUANCE_TIME                    MessageTag = "VT_CERTIFICATE_ISSUANCE_TIME"
	MessageTag_VT_VALIDATION_TIME                              MessageTag = "VT_VALIDATION_TIME"
	MessageTag_VT_TST_GENERATION_TIME                          MessageTag = "VT_TST_GENERATION_TIME"
	MessageTag_VT_TST_POE_TIME                                 MessageTag = "VT_TST_POE_TIME"
	MessageTag_QWAC1_PROFILE                                   MessageTag = "QWAC1_PROFILE"
	MessageTag_QWAC2_PROFILE                                   MessageTag = "QWAC2_PROFILE"
	MessageTag_TLS_BY_QWAC2_PROFILE                            MessageTag = "TLS_BY_QWAC2_PROFILE"
	MessageTag_SEMANTICS_TOTAL_PASSED                          MessageTag = "SEMANTICS_TOTAL_PASSED"
	MessageTag_SEMANTICS_PASSED                                MessageTag = "SEMANTICS_PASSED"
	MessageTag_SEMANTICS_TOTAL_FAILED                          MessageTag = "SEMANTICS_TOTAL_FAILED"
	MessageTag_SEMANTICS_FAILED                                MessageTag = "SEMANTICS_FAILED"
	MessageTag_SEMANTICS_INDETERMINATE                         MessageTag = "SEMANTICS_INDETERMINATE"
	MessageTag_SEMANTICS_NO_SIGNATURE_FOUND                    MessageTag = "SEMANTICS_NO_SIGNATURE_FOUND"
	MessageTag_SEMANTICS_FORMAT_FAILURE                        MessageTag = "SEMANTICS_FORMAT_FAILURE"
	MessageTag_SEMANTICS_HASH_FAILURE                          MessageTag = "SEMANTICS_HASH_FAILURE"
	MessageTag_SEMANTICS_SIG_CRYPTO_FAILURE                    MessageTag = "SEMANTICS_SIG_CRYPTO_FAILURE"
	MessageTag_SEMANTICS_REVOKED                               MessageTag = "SEMANTICS_REVOKED"
	MessageTag_SEMANTICS_EXPIRED                               MessageTag = "SEMANTICS_EXPIRED"
	MessageTag_SEMANTICS_NOT_YET_VALID                         MessageTag = "SEMANTICS_NOT_YET_VALID"
	MessageTag_SEMANTICS_SIG_CONSTRAINTS_FAILURE               MessageTag = "SEMANTICS_SIG_CONSTRAINTS_FAILURE"
	MessageTag_SEMANTICS_CHAIN_CONSTRAINTS_FAILURE             MessageTag = "SEMANTICS_CHAIN_CONSTRAINTS_FAILURE"
	MessageTag_SEMANTICS_CERTIFICATE_CHAIN_GENERAL_FAILURE     MessageTag = "SEMANTICS_CERTIFICATE_CHAIN_GENERAL_FAILURE"
	MessageTag_SEMANTICS_CRYPTO_CONSTRAINTS_FAILURE            MessageTag = "SEMANTICS_CRYPTO_CONSTRAINTS_FAILURE"
	MessageTag_SEMANTICS_POLICY_PROCESSING_ERROR               MessageTag = "SEMANTICS_POLICY_PROCESSING_ERROR"
	MessageTag_SEMANTICS_SIGNATURE_POLICY_NOT_AVAILABLE        MessageTag = "SEMANTICS_SIGNATURE_POLICY_NOT_AVAILABLE"
	MessageTag_SEMANTICS_TIMESTAMP_ORDER_FAILURE               MessageTag = "SEMANTICS_TIMESTAMP_ORDER_FAILURE"
	MessageTag_SEMANTICS_NO_SIGNING_CERTIFICATE_FOUND          MessageTag = "SEMANTICS_NO_SIGNING_CERTIFICATE_FOUND"
	MessageTag_SEMANTICS_NO_CERTIFICATE_CHAIN_FOUND            MessageTag = "SEMANTICS_NO_CERTIFICATE_CHAIN_FOUND"
	MessageTag_SEMANTICS_NO_CERTIFICATE_CHAIN_FOUND_NO_POE     MessageTag = "SEMANTICS_NO_CERTIFICATE_CHAIN_FOUND_NO_POE"
	MessageTag_SEMANTICS_REVOKED_NO_POE                        MessageTag = "SEMANTICS_REVOKED_NO_POE"
	MessageTag_SEMANTICS_REVOKED_CA_NO_POE                     MessageTag = "SEMANTICS_REVOKED_CA_NO_POE"
	MessageTag_SEMANTICS_OUT_OF_BOUNDS_NOT_REVOKED             MessageTag = "SEMANTICS_OUT_OF_BOUNDS_NOT_REVOKED"
	MessageTag_SEMANTICS_OUT_OF_BOUNDS_NO_POE                  MessageTag = "SEMANTICS_OUT_OF_BOUNDS_NO_POE"
	MessageTag_SEMANTICS_REVOCATION_OUT_OF_BOUNDS_NO_POE       MessageTag = "SEMANTICS_REVOCATION_OUT_OF_BOUNDS_NO_POE"
	MessageTag_SEMANTICS_CRYPTO_CONSTRAINTS_FAILURE_NO_POE     MessageTag = "SEMANTICS_CRYPTO_CONSTRAINTS_FAILURE_NO_POE"
	MessageTag_SEMANTICS_NO_POE                                MessageTag = "SEMANTICS_NO_POE"
	MessageTag_SEMANTICS_TRY_LATER                             MessageTag = "SEMANTICS_TRY_LATER"
	MessageTag_SEMANTICS_SIGNED_DATA_NOT_FOUND                 MessageTag = "SEMANTICS_SIGNED_DATA_NOT_FOUND"
	MessageTag_SEMANTICS_EAA_CONSTRAINTS_FAILURE               MessageTag = "SEMANTICS_EAA_CONSTRAINTS_FAILURE"
)

// MessageTagValues returns all MessageTag constants in declaration order.
func MessageTagValues() []MessageTag {
	return []MessageTag{
		MessageTag_BBB_FC_IEFF,
		MessageTag_BBB_FC_IEFF_ANS,
		MessageTag_BBB_FC_ICFD,
		MessageTag_BBB_FC_ICFD_ANS,
		MessageTag_BBB_FC_ISD,
		MessageTag_BBB_FC_ISD_ANS,
		MessageTag_BBB_FC_ISRIA,
		MessageTag_BBB_FC_ISRIA_ANS,
		MessageTag_BBB_FC_IOSIP,
		MessageTag_BBB_FC_IOSIP_ANS,
		MessageTag_BBB_FC_DASTHVBR,
		MessageTag_BBB_FC_DASTHVBR_ANS,
		MessageTag_BBB_FC_DBTOOST,
		MessageTag_BBB_FC_DBTOOST_ANS,
		MessageTag_BBB_FC_IBRV,
		MessageTag_BBB_FC_IBRV_ANS,
		MessageTag_BBB_FC_ISDC,
		MessageTag_BBB_FC_ISDC_ANS,
		MessageTag_BBB_FC_DSFREAP,
		MessageTag_BBB_FC_DSFREAP_ANS,
		MessageTag_BBB_FC_IAOD,
		MessageTag_BBB_FC_IAOD_ANS,
		MessageTag_BBB_FC_IVDBSFR,
		MessageTag_BBB_FC_IVDBSFR_ANS,
		MessageTag_BBB_FC_ISVADMDPD,
		MessageTag_BBB_FC_ISVADMDPD_ANS,
		MessageTag_BBB_FC_ISVAFMDPD,
		MessageTag_BBB_FC_ISVAFMDPD_ANS,
		MessageTag_BBB_FC_ISVASFLD,
		MessageTag_BBB_FC_ISVASFLD_ANS,
		MessageTag_BBB_FC_DSCNFFSM,
		MessageTag_BBB_FC_DSCNFFSM_ANS,
		MessageTag_BBB_FC_DSCNACMDM,
		MessageTag_BBB_FC_DSCNACMDM_ANS,
		MessageTag_BBB_FC_DSCNUOM,
		MessageTag_BBB_FC_DSCNUOM_ANS,
		MessageTag_BBB_FC_IECKSCDA,
		MessageTag_BBB_FC_IECKSCDA_ANS1,
		MessageTag_BBB_FC_IECKSCDA_ANS2,
		MessageTag_BBB_FC_IECKSCDA_ANS3,
		MessageTag_BBB_FC_IECKSCDA_ANS4,
		MessageTag_BBB_FC_IECKSCDA_ANS5,
		MessageTag_BBB_FC_DDAPDFAF,
		MessageTag_BBB_FC_DDAPDFAF_ANS,
		MessageTag_BBB_FC_IDPDFAC,
		MessageTag_BBB_FC_IDPDFAC_ANS,
		MessageTag_BBB_FC_IECTF,
		MessageTag_BBB_FC_IECTF_ANS,
		MessageTag_BBB_FC_ISFCS,
		MessageTag_BBB_FC_ISFCS_ANS,
		MessageTag_BBB_FC_ITFCS,
		MessageTag_BBB_FC_ITFCS_ANS,
		MessageTag_BBB_FC_IMFCS,
		MessageTag_BBB_FC_IMFCS_ANS,
		MessageTag_BBB_FC_ITZCP,
		MessageTag_BBB_FC_ITZCP_ANS,
		MessageTag_BBB_FC_ITEZCF,
		MessageTag_BBB_FC_ITEZCF_ANS,
		MessageTag_BBB_FC_ITMFP,
		MessageTag_BBB_FC_ITMFP_ANS,
		MessageTag_BBB_FC_IEMCF,
		MessageTag_BBB_FC_IEMCF_ANS,
		MessageTag_BBB_FC_IMFP_ASICE,
		MessageTag_BBB_FC_IMFP_ASICE_ANS,
		MessageTag_BBB_FC_ISFP_ASICE,
		MessageTag_BBB_FC_ISFP_ASICE_ANS,
		MessageTag_BBB_FC_ISFP_ASICS,
		MessageTag_BBB_FC_ISFP_ASICS_ANS,
		MessageTag_BBB_FC_ISFP_ASTFORAMC,
		MessageTag_BBB_FC_ISFP_ASTFORAMC_ANS,
		MessageTag_BBB_FC_IAHIV,
		MessageTag_BBB_FC_IAHIV_ANS,
		MessageTag_BBB_CV_IRDOF,
		MessageTag_BBB_CV_IRDOF_ANS,
		MessageTag_BBB_CV_TSP_IRDOF,
		MessageTag_BBB_CV_TSP_IRDOF_ANS,
		MessageTag_BBB_CV_CS_CSSVF,
		MessageTag_BBB_CV_CS_CSSVF_ANS,
		MessageTag_BBB_CV_IMEOF,
		MessageTag_BBB_CV_IMEOF_ANS,
		MessageTag_BBB_CV_EAA_SDCBF,
		MessageTag_BBB_CV_EAA_SDCBF_ANS,
		MessageTag_BBB_CV_EAA_NSDCBF,
		MessageTag_BBB_CV_EAA_NSDCBF_ANS,
		MessageTag_BBB_CV_ER_IODOF,
		MessageTag_BBB_CV_ER_IODOF_ANS,
		MessageTag_BBB_CV_ER_HASSDOC,
		MessageTag_BBB_CV_ER_HASSDOC_ANS,
		MessageTag_BBB_CV_ER_DFHVLCDOG,
		MessageTag_BBB_CV_ER_DFHVLCDOG_ANS,
		MessageTag_BBB_CV_ER_ATSRF,
		MessageTag_BBB_CV_ER_ATSRF_ANS,
		MessageTag_BBB_CV_ER_ATSSRF,
		MessageTag_BBB_CV_ER_ATSSRF_ANS,
		MessageTag_BBB_CV_IRDOI,
		MessageTag_BBB_CV_IRDOI_ANS,
		MessageTag_BBB_CV_TSP_IRDOI,
		MessageTag_BBB_CV_TSP_IRDOI_ANS,
		MessageTag_BBB_CV_CS_CSPS,
		MessageTag_BBB_CV_CS_CSPS_ANS,
		MessageTag_BBB_CV_IMEDOI,
		MessageTag_BBB_CV_IMEDOI_ANS,
		MessageTag_BBB_CV_EAA_SDCBI,
		MessageTag_BBB_CV_EAA_SDCBI_ANS,
		MessageTag_BBB_CV_EAA_NSDCBI,
		MessageTag_BBB_CV_EAA_NSDCBI_ANS,
		MessageTag_BBB_CV_ER_ATSRI,
		MessageTag_BBB_CV_ER_ATSRI_ANS,
		MessageTag_BBB_CV_ER_ATSSRI,
		MessageTag_BBB_CV_ER_ATSSRI_ANS,
		MessageTag_BBB_CV_ISMEC,
		MessageTag_BBB_CV_ISMEC_ANS,
		MessageTag_BBB_CV_ISMEC_ANS_2,
		MessageTag_BBB_CV_AAMEF,
		MessageTag_BBB_CV_AAMEF_ANS,
		MessageTag_BBB_CV_ER_TST_RN,
		MessageTag_BBB_CV_ER_TST_RN_ANS_1,
		MessageTag_BBB_CV_ER_TST_RN_ANS_2,
		MessageTag_BBB_CV_DRNMND,
		MessageTag_BBB_CV_DRNMND_ANS,
		MessageTag_BBB_CV_DMENMND,
		MessageTag_BBB_CV_DMENMND_ANS,
		MessageTag_BBB_CV_ISI,
		MessageTag_BBB_CV_ISIC,
		MessageTag_BBB_CV_ISIR,
		MessageTag_BBB_CV_ISIT,
		MessageTag_BBB_CV_ISI_ANS,
		MessageTag_BBB_CV_IAFS,
		MessageTag_BBB_CV_IAFS_ANS,
		MessageTag_BBB_ICS_ISCI,
		MessageTag_BBB_ICS_ISCI_ANS,
		MessageTag_BBB_ICS_ISASCP,
		MessageTag_BBB_ICS_ISASCP_ANS,
		MessageTag_BBB_ICS_ISASCPU,
		MessageTag_BBB_ICS_ISASCPU_ANS,
		MessageTag_BBB_ICS_ISACDP,
		MessageTag_BBB_ICS_ISACDP_ANS,
		MessageTag_BBB_ICS_ICDVV,
		MessageTag_BBB_ICS_ICDVV_ANS,
		MessageTag_BBB_ICS_ICDVVS,
		MessageTag_BBB_ICS_ICDVVS_ANS,
		MessageTag_BBB_ICS_AIDNASNE,
		MessageTag_BBB_ICS_AIDNASNE_ANS,
		MessageTag_BBB_ICS_ISAKIDP,
		MessageTag_BBB_ICS_ISAKIDP_ANS,
		MessageTag_BBB_ICS_DKIDVM,
		MessageTag_BBB_ICS_DKIDVM_ANS,
		MessageTag_BBB_ICS_ISAX509UP,
		MessageTag_BBB_ICS_ISAX509UP_ANS,
		MessageTag_BBB_ICS_ISAX509UA,
		MessageTag_BBB_ICS_ISAX509UA_ANS,
		MessageTag_BBB_RFC_NUP,
		MessageTag_BBB_RFC_NUP_ANS,
		MessageTag_BBB_RFC_IRIF,
		MessageTag_BBB_RFC_IRIF_TUNU,
		MessageTag_BBB_RFC_IRIF_ANS,
		MessageTag_ADEST_ROBVPIIC,
		MessageTag_ADEST_ROBVPIIC_ANS,
		MessageTag_ADEST_ROTVPIIC,
		MessageTag_ADEST_ROTVPIIC_ANS,
		MessageTag_ADEST_RORPIIC,
		MessageTag_ADEST_RORPIIC_ANS,
		MessageTag_BSV_IFCRC,
		MessageTag_BSV_IFCRC_ANS,
		MessageTag_BSV_IISCRC,
		MessageTag_BSV_IISCRC_ANS,
		MessageTag_BSV_IVCIRC,
		MessageTag_BSV_IVCIRC_ANS,
		MessageTag_BSV_IXCVRC,
		MessageTag_BSV_IXCVRC_ANS,
		MessageTag_BSV_ISCRAVTC,
		MessageTag_BSV_ISCRAVTC_ANS,
		MessageTag_BSV_IVTAVRSC,
		MessageTag_BSV_IVTAVRSC_ANS,
		MessageTag_BSV_ISCCTC,
		MessageTag_BSV_ISCCTC_ANS,
		MessageTag_BSV_ICTGTNASCRT,
		MessageTag_BSV_ICTGTNASCRT_ANS,
		MessageTag_BSV_ICTGTNASCET,
		MessageTag_BSV_ICTGTNASCET_ANS,
		MessageTag_BSV_ICVRC,
		MessageTag_BSV_ICVRC_ANS,
		MessageTag_BSV_ISAVRC,
		MessageTag_BSV_ISAVRC_ANS,
		MessageTag_BSV_ICTGTNACCET,
		MessageTag_BSV_ICTGTNACCET_ANS,
		MessageTag_BSV_IEAAAVRC,
		MessageTag_BSV_IEAAAVRC_ANS,
		MessageTag_LTV_ABSV,
		MessageTag_LTV_ABSV_ANS,
		MessageTag_LTV_ISCKNR,
		MessageTag_LTV_ISCKNR_ANS0,
		MessageTag_LTV_ISCKNR_ANS1,
		MessageTag_ARCH_LTVV,
		MessageTag_ARCH_LTVV_ANS,
		MessageTag_ARCH_LTAIVMP,
		MessageTag_ARCH_LTAIVMP_ANS,
		MessageTag_ARCH_IRTVBBA,
		MessageTag_ARCH_IRTVBBA_ANS,
		MessageTag_ARCH_ICHFCRLPOET,
		MessageTag_ARCH_ICHFCRLPOET_ANS,
		MessageTag_ACCM,
		MessageTag_ACCM_ANS,
		MessageTag_ASCCM_CAA,
		MessageTag_ASCCM_CAA_ANS,
		MessageTag_ASCCM_DAA,
		MessageTag_ASCCM_DAA_ANS,
		MessageTag_ASCCM_DAA_ANS_2,
		MessageTag_ASCCM_APKSA,
		MessageTag_ASCCM_APKSA_ANS,
		MessageTag_ASCCM_APKSA_ANS_2,
		MessageTag_ASCCM_AR,
		MessageTag_ASCCM_AR_ANS_ANR,
		MessageTag_ASCCM_AR_ANS_ANR_2,
		MessageTag_ASCCM_AR_ANS_AKSNR,
		MessageTag_ASCCM_AR_ANS_AKSNR_2,
		MessageTag_ASCCM_PKSK,
		MessageTag_ASCCM_PKSK_ANS,
		MessageTag_ACCM_DESC_WITH_ID,
		MessageTag_ACCM_DESC_WITH_ID_RESULT,
		MessageTag_ACCM_DESC_WITH_NAME,
		MessageTag_ACCM_POS_SIG_SIG,
		MessageTag_ACCM_POS_TST_SIG,
		MessageTag_ACCM_POS_REVOC_SIG,
		MessageTag_ACCM_POS_EV_RECORD,
		MessageTag_ACCM_POS_EAA,
		MessageTag_ACCM_POS_EAA_REV,
		MessageTag_ACCM_POS_CNTR_SIG,
		MessageTag_ACCM_POS_CNTR_SIG_PL,
		MessageTag_ACCM_POS_CON_DIG,
		MessageTag_ACCM_POS_EAA_KB,
		MessageTag_ACCM_POS_EAA_ND,
		MessageTag_ACCM_POS_EAA_NSD,
		MessageTag_ACCM_POS_EAA_NSD_PL,
		MessageTag_ACCM_POS_EAA_OSDC,
		MessageTag_ACCM_POS_EAA_OSDC_PL,
		MessageTag_ACCM_POS_EAA_PD,
		MessageTag_ACCM_POS_EAA_SD,
		MessageTag_ACCM_POS_EAA_SD_PL,
		MessageTag_ACCM_POS_ER_ADO,
		MessageTag_ACCM_POS_ER_ADO_PL,
		MessageTag_ACCM_POS_ER_OR,
		MessageTag_ACCM_POS_ER_OR_PL,
		MessageTag_ACCM_POS_ER_TST,
		MessageTag_ACCM_POS_ER_TST_SEQ,
		MessageTag_ACCM_POS_ER_MST_SIG,
		MessageTag_ACCM_POS_JWS,
		MessageTag_ACCM_POS_COSE,
		MessageTag_ACCM_POS_KEY,
		MessageTag_ACCM_POS_KEY_PL,
		MessageTag_ACCM_POS_MAN,
		MessageTag_ACCM_POS_MAN_PL,
		MessageTag_ACCM_POS_MAN_ENT,
		MessageTag_ACCM_POS_MAN_ENT_PL,
		MessageTag_ACCM_POS_MES_DIG,
		MessageTag_ACCM_POS_MESS_IMP,
		MessageTag_ACCM_POS_REF,
		MessageTag_ACCM_POS_REF_PL,
		MessageTag_ACCM_POS_SIG_D_ENT,
		MessageTag_ACCM_POS_SIG_D_ENT_PL,
		MessageTag_ACCM_POS_SIG_VAL_AND_PRT,
		MessageTag_ACCM_POS_SIGND_OBJ,
		MessageTag_ACCM_POS_SIGND_PRT,
		MessageTag_ACCM_POS_SIGNTR_PRT,
		MessageTag_ACCM_POS_CERT_CHAIN_SIG,
		MessageTag_ACCM_POS_CERT_CHAIN_TST,
		MessageTag_ACCM_POS_CERT_CHAIN_REVOC,
		MessageTag_ACCM_POS_CERT_CHAIN_EAA_REV,
		MessageTag_ACCM_POS_CERT_CHAIN,
		MessageTag_ACCM_POS_SIG_CERT_REF,
		MessageTag_BBB_SAV_DSCACRCC,
		MessageTag_BBB_SAV_DSCACRCC_ANS,
		MessageTag_BBB_SAV_ACPCCRSCA,
		MessageTag_BBB_SAV_ACPCCRSCA_ANS,
		MessageTag_BBB_SAV_ISVA,
		MessageTag_BBB_SAV_ISVA_ANS,
		MessageTag_BBB_SAV_ISSV,
		MessageTag_BBB_SAV_ISSV_ANS,
		MessageTag_BBB_SAV_ICERRM,
		MessageTag_BBB_SAV_ICERRM_ANS,
		MessageTag_BBB_SAV_ICRM,
		MessageTag_BBB_SAV_ICRM_ANS,
		MessageTag_BBB_SAV_ISQPCTP,
		MessageTag_BBB_SAV_ISQPCTP_ANS,
		MessageTag_BBB_SAV_ISQPCHP,
		MessageTag_BBB_SAV_ISQPCHP_ANS,
		MessageTag_BBB_SAV_ISQPCIP,
		MessageTag_BBB_SAV_ISQPCIP_ANS,
		MessageTag_BBB_SAV_ISQPCTSIP,
		MessageTag_BBB_SAV_ISQPCTSIP_ANS,
		MessageTag_BBB_SAV_ISQPSTYPP,
		MessageTag_BBB_SAV_ISQPSTYPP_ANS,
		MessageTag_BBB_SAV_ISQPSLP,
		MessageTag_BBB_SAV_ISQPSLP_ANS,
		MessageTag_BBB_SAV_ISQPSTP,
		MessageTag_BBB_SAV_ISQPSTP_ANS,
		MessageTag_BBB_SAV_ISQPSTWSCVR,
		MessageTag_BBB_SAV_ISQPSTWSCVR_ANS,
		MessageTag_BBB_SAV_ISQPXTIP,
		MessageTag_BBB_SAV_ISQPXTIP_ANS,
		MessageTag_BBB_SAV_IUQPCSP,
		MessageTag_BBB_SAV_IUQPCSP_ANS,
		MessageTag_BBB_SAV_IUQPSTSP,
		MessageTag_BBB_SAV_IUQPSTSP_ANS,
		MessageTag_BBB_SAV_IUQPVDTSP,
		MessageTag_BBB_SAV_IUQPVDTSP_ANS,
		MessageTag_BBB_SAV_IUQPVDROTSP,
		MessageTag_BBB_SAV_IUQPVDROTSP_ANS,
		MessageTag_BBB_SAV_IUQPATSP,
		MessageTag_BBB_SAV_IUQPATSP_ANS,
		MessageTag_BBB_SAV_ICTVS,
		MessageTag_BBB_SAV_ICTVS_ANS,
		MessageTag_BBB_SAV_IDTSP,
		MessageTag_BBB_SAV_IDTSP_ANS,
		MessageTag_BBB_SAV_ITVS,
		MessageTag_BBB_SAV_ITVS_ANS,
		MessageTag_BBB_SAV_IVTTSTP,
		MessageTag_BBB_SAV_IVTTSTP_ANS,
		MessageTag_BBB_SAV_IVLTATSTP,
		MessageTag_BBB_SAV_IVLTATSTP_ANS,
		MessageTag_BBB_SAV_ISQPMDOSPP,
		MessageTag_BBB_SAV_ISQPMDOSPP_ANS,
		MessageTag_BBB_SAV_DMICTSTMCMI,
		MessageTag_BBB_SAV_DMICTSTMCMI_ANS,
		MessageTag_BBB_TAV_ITSAP,
		MessageTag_BBB_TAV_ITSAP_ANS,
		MessageTag_BBB_TAV_DTSAVM,
		MessageTag_BBB_TAV_DTSAVM_ANS,
		MessageTag_BBB_TAV_DTSAOM,
		MessageTag_BBB_TAV_DTSAOM_ANS,
		MessageTag_BBB_VCI_ISPK,
		MessageTag_BBB_VCI_ISPK_ANS,
		MessageTag_BBB_VCI_ISPA,
		MessageTag_BBB_VCI_ISPA_ANS,
		MessageTag_BBB_VCI_ISPSUPP,
		MessageTag_BBB_VCI_ISPSUPP_ANS,
		MessageTag_BBB_VCI_ISPM,
		MessageTag_BBB_VCI_ISPM_ANS,
		MessageTag_BBB_VCI_IZHSP,
		MessageTag_BBB_VCI_IZHSP_ANS,
		MessageTag_BBB_XCV_SUB,
		MessageTag_BBB_XCV_SUB_ANS,
		MessageTag_BBB_XCV_SUB_ANS_2,
		MessageTag_BBB_XCV_RFC,
		MessageTag_BBB_XCV_RFC_ANS,
		MessageTag_BBB_XCV_RAC,
		MessageTag_BBB_XCV_RAC_ANS,
		MessageTag_BBB_XCV_CCCBB,
		MessageTag_BBB_XCV_CCCBB_ANS,
		MessageTag_BBB_XCV_CCCBB_SIG_ANS,
		MessageTag_BBB_XCV_CCCBB_TSP_ANS,
		MessageTag_BBB_XCV_CCCBB_REV_ANS,
		MessageTag_BBB_XCV_CMDCIPI,
		MessageTag_BBB_XCV_CMDCIPI_ANS,
		MessageTag_BBB_XCV_CMDCIQC,
		MessageTag_BBB_XCV_CMDCIQC_ANS,
		MessageTag_BBB_XCV_CMDCIQSCD,
		MessageTag_BBB_XCV_CMDCIQSCD_ANS,
		MessageTag_BBB_XCV_CMDCIITLP,
		MessageTag_BBB_XCV_CMDCIITLP_ANS,
		MessageTag_BBB_XCV_CMDCIITNP,
		MessageTag_BBB_XCV_CMDCIITNP_ANS,
		MessageTag_BBB_XCV_CMDCICQCC,
		MessageTag_BBB_XCV_CMDCICQCC_ANS,
		MessageTag_BBB_XCV_CMDCICQCLVA,
		MessageTag_BBB_XCV_CMDCICQCLVA_ANS,
		MessageTag_BBB_XCV_CMDCICQCLVHAC,
		MessageTag_BBB_XCV_CMDCICQCLVHAC_ANS,
		MessageTag_BBB_XCV_CMDCICQCERPA,
		MessageTag_BBB_XCV_CMDCICQCERPA_ANS,
		MessageTag_BBB_XCV_CMDCICSQCSSCD,
		MessageTag_BBB_XCV_CMDCICSQCSSCD_ANS,
		MessageTag_BBB_XCV_CMDCICQCPDSLA,
		MessageTag_BBB_XCV_CMDCICQCPDSLA_ANS,
		MessageTag_BBB_XCV_CMDCICQCTA,
		MessageTag_BBB_XCV_CMDCICQCTA_ANS,
		MessageTag_BBB_XCV_CMDCDCQCCLCEC,
		MessageTag_BBB_XCV_CMDCDCQCCLCEC_ANS,
		MessageTag_BBB_XCV_CMDCDCQCCLCEC_ANS_EU,
		MessageTag_BBB_XCV_CMDCSCSIA,
		MessageTag_BBB_XCV_CMDCSCSIA_ANS,
		MessageTag_BBB_XCV_CMDCICQCRA,
		MessageTag_BBB_XCV_CMDCICQCRA_ANS,
		MessageTag_BBB_XCV_CMDCICQCNA,
		MessageTag_BBB_XCV_CMDCICQCNA_ANS,
		MessageTag_BBB_XCV_CMDCICQCIA,
		MessageTag_BBB_XCV_CMDCICQCIA_ANS,
		MessageTag_BBB_XCV_CMDCDCQCQSCDLSA,
		MessageTag_BBB_XCV_CMDCDCQCQSCDLSA_ANS,
		MessageTag_BBB_XCV_CMDCDCQCIMSA,
		MessageTag_BBB_XCV_CMDCDCQCIMSA_ANS,
		MessageTag_BBB_XCV_CMDCPSBCLA,
		MessageTag_BBB_XCV_CMDCPSBCLA_ANS,
		MessageTag_BBB_XCV_CMDCPSBASIA,
		MessageTag_BBB_XCV_CMDCPSBASIA_ANS,
		MessageTag_BBB_XCV_CMDCPSBLIA,
		MessageTag_BBB_XCV_CMDCPSBLIA_ANS,
		MessageTag_BBB_XCV_DCCUCE,
		MessageTag_BBB_XCV_DCCUCE_ANS,
		MessageTag_BBB_XCV_DCCFCE,
		MessageTag_BBB_XCV_DCCFCE_ANS,
		MessageTag_BBB_XCV_DCSBSINC,
		MessageTag_BBB_XCV_DCSBSINC_ANS,
		MessageTag_BBB_XCV_ICAC,
		MessageTag_BBB_XCV_ICAC_ANS,
		MessageTag_BBB_XCV_ICPDV,
		MessageTag_BBB_XCV_ICPDV_ANS,
		MessageTag_BBB_XCV_ICPTV,
		MessageTag_BBB_XCV_ICPTV_ANS,
		MessageTag_BBB_XCV_IAKIP,
		MessageTag_BBB_XCV_IAKIP_ANS,
		MessageTag_BBB_XCV_ISKIP,
		MessageTag_BBB_XCV_ISKIP_ANS,
		MessageTag_BBB_XCV_ICNRAEV,
		MessageTag_BBB_XCV_ICNRAEV_ANS,
		MessageTag_BBB_XCV_IVTBCTSD,
		MessageTag_BBB_XCV_IVTBCTSD_ANS,
		MessageTag_BBB_XCV_ICTIVRSC,
		MessageTag_BBB_XCV_ICTIVRSC_ANS,
		MessageTag_BBB_XCV_ICTIVRCIRI,
		MessageTag_BBB_XCV_ICTIVRCIRI_ANS,
		MessageTag_BBB_XCV_IRDCSFC,
		MessageTag_BBB_XCV_IRDCSFC_ANS,
		MessageTag_BBB_XCV_IRDPFC,
		MessageTag_BBB_XCV_IRDPFC_ANS,
		MessageTag_BBB_XCV_IRDPFRC,
		MessageTag_BBB_XCV_IRDPFRC_ANS,
		MessageTag_BBB_XCV_IARDPFC,
		MessageTag_BBB_XCV_IARDPFC_ANS,
		MessageTag_BBB_VTS_IRDPFC,
		MessageTag_BBB_VTS_IRDPFC_ANS,
		MessageTag_BBB_XCV_ISCOH,
		MessageTag_BBB_XCV_ISCOH_ANS,
		MessageTag_BBB_XCV_ISCUKN,
		MessageTag_BBB_XCV_ISCUKN_ANS,
		MessageTag_BBB_XCV_ISCR,
		MessageTag_BBB_XCV_ISCR_ANS,
		MessageTag_BBB_XCV_ISCGKU,
		MessageTag_BBB_XCV_ISCGKU_ANS,
		MessageTag_BBB_XCV_ISCGKU_ANS_CERT,
		MessageTag_BBB_XCV_ISCGEKU,
		MessageTag_BBB_XCV_ISCGEKU_ANS,
		MessageTag_BBB_XCV_ISCGEKU_ANS_CERT,
		MessageTag_BBB_XCV_ICSI,
		MessageTag_BBB_XCV_ICSI_ANS,
		MessageTag_BBB_XCV_IOTAA,
		MessageTag_BBB_XCV_IOTAA_ANS,
		MessageTag_BBB_XCV_HPCCVVT,
		MessageTag_BBB_XCV_HPCCVVT_ANS,
		MessageTag_BBB_XCV_PSEUDO_USE,
		MessageTag_BBB_XCV_PSEUDO_USE_ANS,
		MessageTag_BBB_XCV_AIA_PRES,
		MessageTag_BBB_XCV_AIA_PRES_ANS,
		MessageTag_BBB_XCV_REVOC_PRES,
		MessageTag_BBB_XCV_REVOC_PRES_ANS,
		MessageTag_BBB_XCV_REVOC_THIS_UPDATE_PRESENT,
		MessageTag_BBB_XCV_REVOC_THIS_UPDATE_PRESENT_ANS,
		MessageTag_BBB_XCV_REVOC_ISSUER_KNOWN,
		MessageTag_BBB_XCV_REVOC_ISSUER_KNOWN_ANS,
		MessageTag_BBB_XCV_REVOC_ISSUER_VALID_AT_PROD,
		MessageTag_BBB_XCV_REVOC_ISSUER_VALID_AT_PROD_ANS,
		MessageTag_BBB_XCV_REVOC_AFTER_CERT_NOT_BEFORE,
		MessageTag_BBB_XCV_REVOC_AFTER_CERT_NOT_BEFORE_ANS,
		MessageTag_BBB_XCV_REVOC_HAS_CERT_INFO,
		MessageTag_BBB_XCV_REVOC_HAS_CERT_INFO_ANS,
		MessageTag_BBB_XCV_REVOC_RESPID_MATCH,
		MessageTag_BBB_XCV_REVOC_RESPID_MATCH_ANS,
		MessageTag_BBB_XCV_REVOC_CERT_HASH_PRESENT,
		MessageTag_BBB_XCV_REVOC_CERT_HASH_PRESENT_ANS,
		MessageTag_BBB_XCV_REVOC_CERT_HASH_MATCH,
		MessageTag_BBB_XCV_REVOC_CERT_HASH_MATCH_ANS,
		MessageTag_BBB_XCV_REVOC_SELF_ISSUED_OCSP,
		MessageTag_BBB_XCV_REVOC_SELF_ISSUED_OCSP_ANS,
		MessageTag_BBB_XCV_DCIDNMSDNIC,
		MessageTag_BBB_XCV_DCIDNMSDNIC_ANS,
		MessageTag_BBB_XCV_ISCGCOUN,
		MessageTag_BBB_XCV_ISCGCOUN_ANS,
		MessageTag_BBB_XCV_ISCGLOC,
		MessageTag_BBB_XCV_ISCGLOC_ANS,
		MessageTag_BBB_XCV_ISCGST,
		MessageTag_BBB_XCV_ISCGST_ANS,
		MessageTag_BBB_XCV_ISCGORGAN,
		MessageTag_BBB_XCV_ISCGORGAN_ANS,
		MessageTag_BBB_XCV_ISCGORGAU,
		MessageTag_BBB_XCV_ISCGORGAU_ANS,
		MessageTag_BBB_XCV_ISCGORGAI,
		MessageTag_BBB_XCV_ISCGORGAI_ANS,
		MessageTag_BBB_XCV_ISCGSURN,
		MessageTag_BBB_XCV_ISCGSURN_ANS,
		MessageTag_BBB_XCV_ISCGGIVEN,
		MessageTag_BBB_XCV_ISCGGIVEN_ANS,
		MessageTag_BBB_XCV_ISCGPSEUDO,
		MessageTag_BBB_XCV_ISCGPSEUDO_ANS,
		MessageTag_BBB_XCV_ISCGCOMMONN,
		MessageTag_BBB_XCV_ISCGCOMMONN_ANS,
		MessageTag_BBB_XCV_ISCGTITLE,
		MessageTag_BBB_XCV_ISCGTITLE_ANS,
		MessageTag_BBB_XCV_ISCGEMAIL,
		MessageTag_BBB_XCV_ISCGEMAIL_ANS,
		MessageTag_BBB_XCV_ISSSC,
		MessageTag_BBB_XCV_ISSSC_ANS,
		MessageTag_BBB_XCV_ISNSSC,
		MessageTag_BBB_XCV_ISNSSC_ANS,
		MessageTag_BBB_XCV_IRDC,
		MessageTag_BBB_XCV_IRDC_ANS,
		MessageTag_XCV_TSL_ESP,
		MessageTag_XCV_TSL_ESP_ANS,
		MessageTag_XCV_TSL_ESP_SIG_ANS,
		MessageTag_XCV_TSL_ESP_TSP_ANS,
		MessageTag_XCV_TSL_ESP_REV_ANS,
		MessageTag_XCV_TSL_ETIP,
		MessageTag_XCV_TSL_ETIP_ANS,
		MessageTag_XCV_TSL_ETIP_SIG_ANS,
		MessageTag_XCV_TSL_ETIP_TSP_ANS,
		MessageTag_XCV_TSL_ETIP_REV_ANS,
		MessageTag_PCV_IVTSC,
		MessageTag_PCV_IVTSC_ANS,
		MessageTag_PCV_ICCSVTSF,
		MessageTag_PCV_ICCSVTSF_ANS,
		MessageTag_PSV_IPCVA,
		MessageTag_PSV_IPCVA_ANS,
		MessageTag_PSV_IPCVC,
		MessageTag_PSV_IPCVC_ANS,
		MessageTag_PSV_IPSVC,
		MessageTag_PSV_IPSVC_ANS,
		MessageTag_PSV_IPTVC,
		MessageTag_PSV_IPTVC_ANS,
		MessageTag_PSV_ITPOCOBCT,
		MessageTag_PSV_ITPOSVAOBCT,
		MessageTag_PSV_ITPOSVAOBCT_ANS,
		MessageTag_PSV_ITPORDAOBCT,
		MessageTag_PSV_ITPOOBCT_ANS,
		MessageTag_PSV_ITPRISCNARTCAC,
		MessageTag_PSV_ITPRISCNARTCAC_ANS,
		MessageTag_PSV_ICRDIT,
		MessageTag_PSV_ICRDIT_ANS,
		MessageTag_PSV_IPCRIAIDBEDC,
		MessageTag_PSV_IPCRIAIDBEDC_ANS,
		MessageTag_PSV_ICTD,
		MessageTag_PSV_ICTD_ANS,
		MessageTag_PSV_ISDDTA,
		MessageTag_PSV_ISDDTA_ANS,
		MessageTag_PSV_HRDBIBCT,
		MessageTag_PSV_HRDBIBCT_ANS,
		MessageTag_PSV_DIURDSCHPVR,
		MessageTag_PSV_DIURDSCHPVR_ANS,
		MessageTag_TSV_ASTPTCT,
		MessageTag_TSV_ASTPTCT_ANS,
		MessageTag_TSV_IBSTAIDOSC,
		MessageTag_TSV_IBSTAIDOSC_ANS,
		MessageTag_TSV_IBSTBCEC,
		MessageTag_TSV_IBSTBCEC_ANS,
		MessageTag_TSV_ISCNVABST,
		MessageTag_TSV_ISCNVABST_ANS,
		MessageTag_ADEST_IRTPTBST,
		MessageTag_ADEST_IRTPTBST_ANS,
		MessageTag_ADEST_ISTPTBST,
		MessageTag_ADEST_ISTPTBST_ANS,
		MessageTag_ADEST_VFDTAOCST_ANS,
		MessageTag_ADEST_ISTPTDABST,
		MessageTag_ADEST_ISTPTDABST_ANS,
		MessageTag_ADEST_IBSVPSC,
		MessageTag_ADEST_IBSVPSC_ANS,
		MessageTag_ADEST_IBSVPTC,
		MessageTag_ADEST_IBSVPTC_ANS,
		MessageTag_ADEST_IBSVPTADC,
		MessageTag_ADEST_IBSVPTADC_ANS,
		MessageTag_ADEST_IRERVPC,
		MessageTag_ADEST_IRERVPC_ANS,
		MessageTag_EAA_CERT_LOTE_REACHED,
		MessageTag_EAA_CERT_LOTE_REACHED_ANS,
		MessageTag_EAA_CERT_TRUST_ANCHOR_LIST_REACHED,
		MessageTag_EAA_CERT_TRUST_ANCHOR_LIST_REACHED_ANS,
		MessageTag_EAA_DPEAAP,
		MessageTag_EAA_DPEAAP_ANS,
		MessageTag_EAA_DLEEAAP,
		MessageTag_EAA_DLEEAAP_ANS,
		MessageTag_EAA_KBRC,
		MessageTag_EAA_KBRC_ANS,
		MessageTag_EAA_KBSP,
		MessageTag_EAA_KBSP_ANS,
		MessageTag_EAA_CLAIMS,
		MessageTag_EAA_CLAIMS_ANS,
		MessageTag_EAA_CLAIMS_INFO,
		MessageTag_EAA_SUPPORTED_CLAIMS,
		MessageTag_EAA_SUPPORTED_CLAIMS_ANS,
		MessageTag_EAA_UNSUPPORTED_CLAIMS,
		MessageTag_EAA_ACCEPTABLE_TYPE,
		MessageTag_EAA_ACCEPTABLE_TYPE_ANS,
		MessageTag_EAA_IDENTIFIER_PRESENT,
		MessageTag_EAA_IDENTIFIER_PRESENT_ANS,
		MessageTag_EAA_ISSUANCE_DATE_PRESENT,
		MessageTag_EAA_ISSUANCE_DATE_PRESENT_ANS,
		MessageTag_EAA_NBF_PRESENT,
		MessageTag_EAA_NBF_PRESENT_ANS,
		MessageTag_EAA_EXP_PRESENT,
		MessageTag_EAA_EXP_PRESENT_ANS,
		MessageTag_EAA_AID_PRESENT,
		MessageTag_EAA_AID_PRESENT_ANS,
		MessageTag_EAA_AED_PRESENT,
		MessageTag_EAA_AED_PRESENT_ANS,
		MessageTag_EAA_SIG_PRESENT,
		MessageTag_EAA_SIG_PRESENT_ANS,
		MessageTag_EAA_SIG_QUAL,
		MessageTag_EAA_SIG_QUAL_ANS,
		MessageTag_EAA_CAT_EAA,
		MessageTag_EAA_CAT_EAA_ANS_1,
		MessageTag_EAA_CAT_EAA_ANS_2,
		MessageTag_EAA_CAT_PUBEAA,
		MessageTag_EAA_CAT_PUBEAA_ANS,
		MessageTag_EAA_CAT_QEAA,
		MessageTag_EAA_CAT_QEAA_ANS,
		MessageTag_EAA_QC_PSB,
		MessageTag_EAA_QC_PSB_ANS,
		MessageTag_EAA_QUAL_CONCLUSIVE,
		MessageTag_EAA_QUAL_CONCLUSIVE_ANS,
		MessageTag_EAA_ETSI194721,
		MessageTag_EAA_ETSI194721_ANS,
		MessageTag_EAA_VT_ITVR,
		MessageTag_EAA_VT_ITVR_ANS,
		MessageTag_EAA_VT_ITVR_VALIDITY,
		MessageTag_EAA_NOW_BEFORE_NBF,
		MessageTag_EAA_NOW_AFTER_EXP,
		MessageTag_EAA_VT_IAVR,
		MessageTag_EAA_VT_IAVR_ANS,
		MessageTag_EAA_VT_IAVR_VALIDITY,
		MessageTag_EAA_NOW_BEFORE_ADI,
		MessageTag_EAA_NOW_AFTER_ADE,
		MessageTag_EAA_AD_SDJWT_CONFORMANCE,
		MessageTag_EAA_SHORT_LIVED_STATUS_PRESENT,
		MessageTag_EAA_MANDATORY_STATUS_ABSENT,
		MessageTag_EAA_REV_SDJWT_CONFORMANCE,
		MessageTag_EAA_MDOC_ISSUING_AUTHORITY,
		MessageTag_EAA_SDJWT_ISSUING_AUTHORITY,
		MessageTag_EAA_MDOC_DOCUMENT_NUMBER_ABSENT,
		MessageTag_EAA_SUB,
		MessageTag_EAA_SUB_ANS,
		MessageTag_EAA_SUB_PSE,
		MessageTag_EAA_SUB_PSE_ANS,
		MessageTag_EAA_CAT,
		MessageTag_EAA_CAT_ANS,
		MessageTag_EAA_ISS_COUN,
		MessageTag_EAA_ISS_COUN_ANS,
		MessageTag_EAA_ISS_AUTH,
		MessageTag_EAA_ISS_AUTH_ANS,
		MessageTag_EAA_ISS_REG_ID,
		MessageTag_EAA_ISS_REG_ID_ANS,
		MessageTag_EAA_REV_PR,
		MessageTag_EAA_REV_PR_ANS,
		MessageTag_EAA_REV_AV,
		MessageTag_EAA_REV_AV_ANS,
		MessageTag_EAA_REV_ACC,
		MessageTag_EAA_REV_ACC_ANS,
		MessageTag_EAA_REV_ACC_FND,
		MessageTag_EAA_REV_ACC_FND_ANS,
		MessageTag_EAA_REV_NOT_REV,
		MessageTag_EAA_REV_NOT_REV_ANS,
		MessageTag_EAA_REV_NOT_ON_HOLD,
		MessageTag_EAA_REV_NOT_ON_HOLD_ANS,
		MessageTag_EAA_SH_LVD,
		MessageTag_EAA_SH_LVD_ANS,
		MessageTag_EAA_OTU,
		MessageTag_EAA_OTU_ANS,
		MessageTag_EAA_PSEUDO_USED,
		MessageTag_EAA_PSEUDO_USED_ANS,
		MessageTag_SDJWT_EAA_VCT_PRESENT,
		MessageTag_SDJWT_EAA_VCT_PRESENT_ANS,
		MessageTag_SDJWT_EAA_VCT_INT_PRESENT,
		MessageTag_SDJWT_EAA_VCT_INT_PRESENT_ANS,
		MessageTag_EAA_REV_TYPE,
		MessageTag_EAA_REV_TYPE_ANS,
		MessageTag_EAA_REV_KNOWN,
		MessageTag_EAA_REV_KNOWN_ANS,
		MessageTag_EAA_REV_ISS,
		MessageTag_EAA_REV_ISS_ANS,
		MessageTag_EAA_REV_EXP,
		MessageTag_EAA_REV_EXP_ANS,
		MessageTag_EAA_REV_NOT_EXP,
		MessageTag_EAA_REV_NOT_EXP_ANS,
		MessageTag_EAA_REV_SUB,
		MessageTag_EAA_REV_SUB_ANS,
		MessageTag_EAA_REV_SUB_MATCH,
		MessageTag_EAA_REV_SUB_MATCH_ANS,
		MessageTag_EAA_REV_ISS_VALID,
		MessageTag_EAA_REV_ISS_VALID_ANS,
		MessageTag_EAA_REV_TIME,
		MessageTag_EAA_REV_ISS_CERT,
		MessageTag_QUAL_TL_EXP,
		MessageTag_QUAL_TL_EXP_ANS,
		MessageTag_QUAL_TL_FRESH,
		MessageTag_QUAL_TL_FRESH_ANS,
		MessageTag_QUAL_TL_VERSION,
		MessageTag_QUAL_TL_VERSION_ANS,
		MessageTag_QUAL_TL_WS,
		MessageTag_QUAL_TL_WS_ANS,
		MessageTag_QUAL_TL_SV,
		MessageTag_QUAL_TL_SV_ANS,
		MessageTag_QUAL_TL_IMRA,
		MessageTag_QUAL_TL_IMRA_ANS,
		MessageTag_QUAL_TL_IMRA_ANS_V1,
		MessageTag_QUAL_TL_IMRA_ANS_V2,
		MessageTag_QUAL_TL_SERV_CONS,
		MessageTag_QUAL_TL_SERV_CONS_ANS0,
		MessageTag_QUAL_TL_SERV_CONS_ANS1,
		MessageTag_QUAL_TL_SERV_CONS_ANS2,
		MessageTag_QUAL_TL_SERV_CONS_ANS3,
		MessageTag_QUAL_TL_SERV_CONS_ANS3A,
		MessageTag_QUAL_TL_SERV_CONS_ANS3B,
		MessageTag_QUAL_TL_SERV_CONS_ANS3C,
		MessageTag_QUAL_TL_SERV_CONS_ANS4,
		MessageTag_QUAL_TL_SERV_CONS_ANS5,
		MessageTag_QUAL_TL_SERV_CONS_ANS6,
		MessageTag_QUAL_TL_SERV_CONS_ANS7,
		MessageTag_QUAL_CERT_TRUSTED_LIST_REACHED,
		MessageTag_QUAL_CERT_TRUSTED_LIST_REACHED_ANS,
		MessageTag_QUAL_TRUSTED_LIST_ACCEPT,
		MessageTag_QUAL_TRUSTED_LIST_ACCEPT_ANS,
		MessageTag_QUAL_LIST_OF_TRUSTED_LISTS_ACCEPT,
		MessageTag_QUAL_LIST_OF_TRUSTED_LISTS_ACCEPT_ANS,
		MessageTag_QUAL_VALID_TRUSTED_LIST_PRESENT,
		MessageTag_QUAL_VALID_TRUSTED_LIST_PRESENT_ANS,
		MessageTag_QUAL_CERT_TYPE_AT_ST,
		MessageTag_QUAL_CERT_TYPE_AT_ST_ANS,
		MessageTag_QUAL_CERT_TYPE_AT_CC,
		MessageTag_QUAL_CERT_TYPE_AT_CC_ANS,
		MessageTag_QUAL_CERT_TYPE_AT_VT,
		MessageTag_QUAL_CERT_TYPE_AT_VT_ANS,
		MessageTag_QUAL_QC_AT_ST,
		MessageTag_QUAL_QC_AT_ST_ANS,
		MessageTag_QUAL_QC_AT_CC,
		MessageTag_QUAL_QC_AT_CC_ANS,
		MessageTag_QUAL_QC_AT_VT,
		MessageTag_QUAL_QC_AT_VT_ANS,
		MessageTag_QUAL_QSCD_AT_ST,
		MessageTag_QUAL_QSCD_AT_ST_ANS,
		MessageTag_QUAL_QSCD_AT_CC,
		MessageTag_QUAL_QSCD_AT_CC_ANS,
		MessageTag_QUAL_QSCD_AT_VT,
		MessageTag_QUAL_QSCD_AT_VT_ANS,
		MessageTag_QUAL_UNIQUE_CERT,
		MessageTag_QUAL_UNIQUE_CERT_ANS,
		MessageTag_QUAL_IS_ADES,
		MessageTag_QUAL_IS_ADES_IND,
		MessageTag_QUAL_IS_ADES_INV,
		MessageTag_QUAL_HAS_METS,
		MessageTag_QUAL_HAS_METS_ANS,
		MessageTag_QUAL_HAS_METS_ATTIME,
		MessageTag_QUAL_HAS_METS_ATTIME_ANS,
		MessageTag_QUAL_HAS_METS_HCCECBA,
		MessageTag_QUAL_HAS_METS_HCCECBA_ANS,
		MessageTag_QUAL_HAS_METS_HCCECBA_ANS_2,
		MessageTag_QUAL_HAS_METS_HCCECBA_ANS_3,
		MessageTag_QUAL_HAS_CAQC,
		MessageTag_QUAL_HAS_CAQC_ANS,
		MessageTag_QUAL_HAS_CAQC_ANS_2,
		MessageTag_QUAL_HAS_ATTIME,
		MessageTag_QUAL_HAS_ATTIME_ANS,
		MessageTag_QUAL_HAS_TS_CERT_TYPE,
		MessageTag_QUAL_HAS_TS_CERT_TYPE_ANS,
		MessageTag_QUAL_HAS_CONF,
		MessageTag_QUAL_HAS_CONF_ANS,
		MessageTag_QUAL_HAS_QEAA,
		MessageTag_QUAL_HAS_QEAA_ANS,
		MessageTag_QUAL_HAS_QTST,
		MessageTag_QUAL_HAS_QTST_ANS,
		MessageTag_QUAL_IS_TRUST_CERT_MATCH_SERVICE,
		MessageTag_QUAL_IS_TRUST_CERT_MATCH_SERVICE_ANS0,
		MessageTag_QUAL_IS_TRUST_CERT_MATCH_SERVICE_ANS1,
		MessageTag_QUAL_IS_TRUST_CERT_MATCH_SERVICE_ANS2,
		MessageTag_QUAL_HAS_GRANTED,
		MessageTag_QUAL_HAS_GRANTED_ANS,
		MessageTag_QUAL_HAS_GRANTED_ANS_2,
		MessageTag_QUAL_HAS_GRANTED_AT,
		MessageTag_QUAL_HAS_GRANTED_AT_ANS,
		MessageTag_QUAL_HAS_CONSISTENT_BY_QC,
		MessageTag_QUAL_HAS_CONSISTENT_BY_QC_ANS,
		MessageTag_QUAL_HAS_CONSISTENT_BY_QSCD,
		MessageTag_QUAL_HAS_CONSISTENT_BY_QSCD_ANS,
		MessageTag_QUAL_HAS_CERT_TYPE_COVERAGE,
		MessageTag_QUAL_HAS_CERT_TYPE_COVERAGE_ANS,
		MessageTag_QUAL_HAS_VALID_CAQC,
		MessageTag_QUAL_HAS_VALID_CAQC_ANS,
		MessageTag_QUAL_HAS_ONLY_ONE,
		MessageTag_QUAL_HAS_ONLY_ONE_ANS,
		MessageTag_QWAC_VALID,
		MessageTag_QWAC_VALID_ANS,
		MessageTag_QWAC_VALID_ANS_2,
		MessageTag_QWAC_CERT_QUAL_CONCLUSIVE,
		MessageTag_QWAC_CERT_QUAL_CONCLUSIVE_ANS,
		MessageTag_QWAC_IS_WSA_AT_TIME,
		MessageTag_QWAC_IS_WSA_AT_TIME_ANS,
		MessageTag_QWAC_CERT_POLICY,
		MessageTag_QWAC_CERT_POLICY_ANS,
		MessageTag_QWAC_VAL_PERIOD,
		MessageTag_QWAC_VAL_PERIOD_ANS,
		MessageTag_QWAC_DOMAIN_NAME,
		MessageTag_QWAC_DOMAIN_NAME_ANS,
		MessageTag_QWAC2_EXT_KEY_USAGE,
		MessageTag_QWAC2_EXT_KEY_USAGE_ANS,
		MessageTag_TLS_CERT_BINDING_URL,
		MessageTag_TLS_CERT_BINDING_URL_ANS,
		MessageTag_TLS_CERT_BINDING_SIG,
		MessageTag_TLS_CERT_BINDING_SIG_ANS,
		MessageTag_TLS_CERT_BINDING_SIG_FORM,
		MessageTag_TLS_CERT_BINDING_SIG_FORM_ANS,
		MessageTag_TLS_CERT_BINDING_SIG_SER,
		MessageTag_TLS_CERT_BINDING_SIG_SER_ANS,
		MessageTag_TLS_CERT_BINDING_SIG_EXP,
		MessageTag_TLS_CERT_BINDING_SIG_EXP_ANS,
		MessageTag_TLS_CERT_BINDING_SIG_EXPIRY_DATE,
		MessageTag_TLS_CERT_BINDING_SIG_EXPIRY_DATE_ANS,
		MessageTag_TLS_CERT_BINDING_QWAC2,
		MessageTag_TLS_CERT_BINDING_QWAC2_ANS,
		MessageTag_TLS_CERT_BINDING_SIG_VALID,
		MessageTag_TLS_CERT_BINDING_SIG_VALID_ANS,
		MessageTag_TLS_CERT_BINDING_CERT_IDENTIFIED,
		MessageTag_TLS_CERT_BINDING_CERT_IDENTIFIED_ANS,
		MessageTag_CERT_USAGE_LOTE_ACCEPT,
		MessageTag_CERT_USAGE_LOTE_ACCEPT_ANS,
		MessageTag_CERT_USAGE_LOLOTE_ACCEPT,
		MessageTag_CERT_USAGE_LOLOTE_ACCEPT_ANS,
		MessageTag_CERT_USAGE_VALID_LOTE_PRESENT,
		MessageTag_CERT_USAGE_VALID_LOTE_PRESENT_ANS,
		MessageTag_CERT_USAGE_HAS_ATTIME,
		MessageTag_CERT_USAGE_HAS_ATTIME_ANS,
		MessageTag_CERT_USAGE_LIST_TYPE_KNOWN,
		MessageTag_CERT_USAGE_LIST_TYPE_KNOWN_ANS,
		MessageTag_CERT_USAGE_STATUS,
		MessageTag_CERT_USAGE_STATUS_ANS,
		MessageTag_CERT_USAGE_STATUS_CONS,
		MessageTag_CERT_USAGE_STATUS_CONS_ANS,
		MessageTag_CERT_USAGE_STATUS_KNOWN,
		MessageTag_CERT_USAGE_STATUS_KNOWN_ANS,
		MessageTag_CERT_USAGE_STI,
		MessageTag_CERT_USAGE_STI_ANS,
		MessageTag_CERT_USAGE_STI_KNOWN,
		MessageTag_CERT_USAGE_STI_KNOWN_ANS,
		MessageTag_PID_DOCUMENT_TYPE,
		MessageTag_PID_DOCUMENT_TYPE_ANS,
		MessageTag_PID_LOTE_TYPE_PID_PROVIDERS,
		MessageTag_PID_LOTE_TYPE_PID_PROVIDERS_ANS,
		MessageTag_PID_STI_PID_ISSUANCE,
		MessageTag_PID_STI_PID_ISSUANCE_ANS,
		MessageTag_PID_PROVIDER_AT_ISSUANCE_TIME,
		MessageTag_PID_PROVIDER_AT_ISSUANCE_TIME_ANS,
		MessageTag_PID_PROVIDER_AT_VALIDATION_TIME,
		MessageTag_PID_PROVIDER_AT_VALIDATION_TIME_ANS,
		MessageTag_BBB_ACCEPT,
		MessageTag_BBB_ACCEPT_ANS,
		MessageTag_TST_TYPE_CONTENT_TST,
		MessageTag_TST_TYPE_SIGNATURE_TST,
		MessageTag_TST_TYPE_VD_TST,
		MessageTag_TST_TYPE_DOC_TST,
		MessageTag_TST_TYPE_CONTAINER_TST,
		MessageTag_TST_TYPE_ARCHIVE_TST,
		MessageTag_TST_TYPE_ER_TST,
		MessageTag_TST_TYPE_REF_ER_ATST,
		MessageTag_TST_TYPE_REF_ER_ATST_SEQ,
		MessageTag_EMPTY,
		MessageTag_CERTIFICATE,
		MessageTag_CA_CERTIFICATE,
		MessageTag_SIGNING_CERTIFICATE,
		MessageTag_REVOCATION,
		MessageTag_REVOCATION_SIG_CERT,
		MessageTag_REVOCATION_CA_CERT,
		MessageTag_SIGNATURE,
		MessageTag_TIMESTAMP,
		MessageTag_TIMESTAMP_SIG_CERT,
		MessageTag_TIMESTAMP_CA_CERT,
		MessageTag_EAA_REV,
		MessageTag_EAA_REV_SIG_CERT,
		MessageTag_EAA_REV_CA_CERT,
		MessageTag_ACCEPTABLE_REVOCATION,
		MessageTag_BASIC_SIGNATURE_VALIDATION_RESULT,
		MessageTag_BEST_SIGNATURE_TIME_CERT_NOT_AFTER,
		MessageTag_BEST_SIGNATURE_TIME_CERT_NOT_BEFORE,
		MessageTag_BEST_SIGNATURE_TIME_CERT_REVOCATION,
		MessageTag_BEST_SIGNATURE_TIME_CERT_SUSPENSION,
		MessageTag_CERTIFICATE_ID,
		MessageTag_CERTIFICATE_REVOCATION_FOUND,
		MessageTag_CERTIFICATE_REVOCATION_NOT_FOUND,
		MessageTag_CERTIFICATE_SUNSET_DATE,
		MessageTag_CERTIFICATE_SUNSET_DATE_TRUST_ANCHOR,
		MessageTag_CERTIFICATE_SUNSET_DATE_VALID,
		MessageTag_CERTIFICATE_TYPE,
		MessageTag_CERTIFICATE_VALIDITY,
		MessageTag_CERTIFICATE_USAGE,
		MessageTag_CERTIFICATE_USAGE_LIST_TYPE,
		MessageTag_CERTIFICATE_USAGE_STATUS,
		MessageTag_CERTIFICATE_USAGE_STATUSES,
		MessageTag_CERTIFICATE_USAGE_STI,
		MessageTag_CERTIFICATE_USAGE_STIS,
		MessageTag_CONTROL_TIME,
		MessageTag_CONTROL_TIME_ALONE,
		MessageTag_CONTROL_TIME_WITH_POE,
		MessageTag_CONTROL_TIME_WITH_TRUST_ANCHOR,
		MessageTag_CRYPTOGRAPHIC_CHECK_FAILURE,
		MessageTag_CRYPTOGRAPHIC_CHECK_FAILURE_WITH_ID,
		MessageTag_CRYPTOGRAPHIC_CHECK_FAILURE_WITH_REF,
		MessageTag_CRYPTOGRAPHIC_CHECK_FAILURE_WITH_REF_WITH_NAME,
		MessageTag_CRYPTOGRAPHIC_CHECK_FAILURE_WITH_REF_WITH_NAMES,
		MessageTag_CRYPTOGRAPHIC_CHECK_SUCCESS,
		MessageTag_CRYPTOGRAPHIC_CHECK_SUCCESS_KEY_SIZE,
		MessageTag_CRYPTOGRAPHIC_CHECK_SUCCESS_DM,
		MessageTag_CRYPTOGRAPHIC_CHECK_SUCCESS_DM_WITH_ID,
		MessageTag_CRYPTOGRAPHIC_CHECK_SUCCESS_DM_WITH_NAME,
		MessageTag_CRYPTOGRAPHIC_CHECK_SUCCESS_DM_WITH_NAMES,
		MessageTag_EVIDENCE_RECORD_VALIDATION,
		MessageTag_EXTENDED_KEY_USAGE,
		MessageTag_KEY_USAGE,
		MessageTag_LAST_ACCEPTABLE_REVOCATION,
		MessageTag_LIST_OF_TRUSTED_ENTITIES,
		MessageTag_PSEUDO,
		MessageTag_QWAC_EXPIRY_EXP,
		MessageTag_QWAC_EXPIRY_TLS_CERT,
		MessageTag_QWAC_EXPIRY_SIGN_CERT,
		MessageTag_REFERENCE,
		MessageTag_REFERENCE_NAME_CHECK,
		MessageTag_REFERENCES_WITH_NAMES,
		MessageTag_REVOCATION_ACCEPTANCE_CHECK,
		MessageTag_REVOCATION_CERT_HASH_OK,
		MessageTag_REVOCATION_CERT_HASH_OK_ID,
		MessageTag_REVOCATION_CERT_VALIDITY,
		MessageTag_REVOCATION_CHECK,
		MessageTag_REVOCATION_CONSISTENT,
		MessageTag_REVOCATION_CONSISTENT_CRL,
		MessageTag_REVOCATION_CONSISTENT_OCSP,
		MessageTag_REVOCATION_CONSISTENT_TL,
		MessageTag_REVOCATION_INFO,
		MessageTag_REVOCATION_NOT_AFTER_AFTER,
		MessageTag_REVOCATION_NOT_AFTER_AFTER_ID,
		MessageTag_REVOCATION_PRODUCED_AT_CERT_VALIDITY,
		MessageTag_REVOCATION_PRODUCED_AT_OUT_OF_BOUNDS,
		MessageTag_REVOCATION_PRODUCED_AT_OUT_OF_BOUNDS_ID,
		MessageTag_REVOCATION_REASON,
		MessageTag_REVOCATION_THIS_UPDATE_CONTROL_TIME,
		MessageTag_SIGNATURE_ALGORITHM_WITH_KEY_SIZE,
		MessageTag_SIGNATURE_ID,
		MessageTag_STRUCTURAL_VALIDATION_FAILURE,
		MessageTag_TIMESTAMP_AND_CERTIFICATE_NOT_AFTER,
		MessageTag_TIMESTAMP_AND_CRYPTO_CONSTRAINTS_EXPIRATION,
		MessageTag_TIMESTAMP_AND_REVOCATION_TIME,
		MessageTag_TIMESTAMP_VALIDATION,
		MessageTag_TOKEN_ID,
		MessageTag_TRUST_SERVICE_NAME,
		MessageTag_TRUSTED_SERVICE_STATUS,
		MessageTag_TRUSTED_SERVICE_TYPE,
		MessageTag_TRUSTED_LIST,
		MessageTag_VALIDATION_TIME,
		MessageTag_CRYPTOGRAPHIC_VERIFICATION,
		MessageTag_FORMAT_CHECKING,
		MessageTag_IDENTIFICATION_OF_THE_SIGNING_CERTIFICATE,
		MessageTag_PAST_SIGNATURE_VALIDATION,
		MessageTag_PAST_CERTIFICATE_VALIDATION,
		MessageTag_REVOCATION_FRESHNESS_CHECKER,
		MessageTag_SIGNATURE_ACCEPTANCE_VALIDATION,
		MessageTag_VALIDATION_CONTEXT_INITIALIZATION,
		MessageTag_VALIDATION_TIME_SLIDING,
		MessageTag_X509_CERTIFICATE_VALIDATION,
		MessageTag_RESULTS,
		MessageTag_EEA_TYPE,
		MessageTag_EAA_ACCEPTANCE_VALIDATION,
		MessageTag_AOV,
		MessageTag_CERT_QUALIFICATION,
		MessageTag_CERT_QUALIFICATION_AT_TIME,
		MessageTag_CERT_USAGE_AT_TIME,
		MessageTag_CERT_USAGES,
		MessageTag_CC,
		MessageTag_CRS,
		MessageTag_DAAV,
		MessageTag_EAA_QUALIFICATION,
		MessageTag_EAA_QUALIFICATION_PROCESS,
		MessageTag_LOTE,
		MessageTag_LOLOTE,
		MessageTag_LOTL,
		MessageTag_PID_QUALIFICATION_PROCESS,
		MessageTag_PSV_CRS,
		MessageTag_QWAC_VALIDATION,
		MessageTag_QWAC_VALIDATION_PROFILE,
		MessageTag_RAC,
		MessageTag_SIG_QUALIFICATION,
		MessageTag_SUB_XCV,
		MessageTag_TL,
		MessageTag_TST_QUALIFICATION,
		MessageTag_TST_QUALIFICATION_AT_TIME,
		MessageTag_VPBS,
		MessageTag_VPEAA,
		MessageTag_VPER,
		MessageTag_VPFLTVD,
		MessageTag_VPFRVC,
		MessageTag_VPFSWATSP,
		MessageTag_VPFTSP,
		MessageTag_VPFTSPWATSP,
		MessageTag_VTS_CRS,
		MessageTag_VT_BEST_SIGNATURE_TIME,
		MessageTag_VT_CERTIFICATE_ISSUANCE_TIME,
		MessageTag_VT_VALIDATION_TIME,
		MessageTag_VT_TST_GENERATION_TIME,
		MessageTag_VT_TST_POE_TIME,
		MessageTag_QWAC1_PROFILE,
		MessageTag_QWAC2_PROFILE,
		MessageTag_TLS_BY_QWAC2_PROFILE,
		MessageTag_SEMANTICS_TOTAL_PASSED,
		MessageTag_SEMANTICS_PASSED,
		MessageTag_SEMANTICS_TOTAL_FAILED,
		MessageTag_SEMANTICS_FAILED,
		MessageTag_SEMANTICS_INDETERMINATE,
		MessageTag_SEMANTICS_NO_SIGNATURE_FOUND,
		MessageTag_SEMANTICS_FORMAT_FAILURE,
		MessageTag_SEMANTICS_HASH_FAILURE,
		MessageTag_SEMANTICS_SIG_CRYPTO_FAILURE,
		MessageTag_SEMANTICS_REVOKED,
		MessageTag_SEMANTICS_EXPIRED,
		MessageTag_SEMANTICS_NOT_YET_VALID,
		MessageTag_SEMANTICS_SIG_CONSTRAINTS_FAILURE,
		MessageTag_SEMANTICS_CHAIN_CONSTRAINTS_FAILURE,
		MessageTag_SEMANTICS_CERTIFICATE_CHAIN_GENERAL_FAILURE,
		MessageTag_SEMANTICS_CRYPTO_CONSTRAINTS_FAILURE,
		MessageTag_SEMANTICS_POLICY_PROCESSING_ERROR,
		MessageTag_SEMANTICS_SIGNATURE_POLICY_NOT_AVAILABLE,
		MessageTag_SEMANTICS_TIMESTAMP_ORDER_FAILURE,
		MessageTag_SEMANTICS_NO_SIGNING_CERTIFICATE_FOUND,
		MessageTag_SEMANTICS_NO_CERTIFICATE_CHAIN_FOUND,
		MessageTag_SEMANTICS_NO_CERTIFICATE_CHAIN_FOUND_NO_POE,
		MessageTag_SEMANTICS_REVOKED_NO_POE,
		MessageTag_SEMANTICS_REVOKED_CA_NO_POE,
		MessageTag_SEMANTICS_OUT_OF_BOUNDS_NOT_REVOKED,
		MessageTag_SEMANTICS_OUT_OF_BOUNDS_NO_POE,
		MessageTag_SEMANTICS_REVOCATION_OUT_OF_BOUNDS_NO_POE,
		MessageTag_SEMANTICS_CRYPTO_CONSTRAINTS_FAILURE_NO_POE,
		MessageTag_SEMANTICS_NO_POE,
		MessageTag_SEMANTICS_TRY_LATER,
		MessageTag_SEMANTICS_SIGNED_DATA_NOT_FOUND,
		MessageTag_SEMANTICS_EAA_CONSTRAINTS_FAILURE,
	}
}

// MessageTagValueOf returns the MessageTag matching the given Java enum name.
func MessageTagValueOf(name string) (MessageTag, error) {
	for _, v := range MessageTagValues() {
		if string(v) == name {
			return v, nil
		}
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
	for _, messageTag := range MessageTagValues() {
		if string(messageTag) == expectedEnumValue {
			return messageTag, true
		}
	}
	return "", false
}
