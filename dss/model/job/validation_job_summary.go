// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/job/ValidationJobSummary.java (DSS 6.5.RC1).
package job

// ValidationJobSummary contains a summary of the validation job, with validation information
// for every document or/and a document list. D is the DocumentInfo type, L the
// DocumentListInfo type, mirroring Java's
// "ValidationJobSummary<D extends DocumentInfo<L>, L extends DocumentListInfo<L, D>>".
type ValidationJobSummary[D any, L any] interface {
	// DocumentListInfos returns a list of document list infos. Port of
	// getDocumentListInfos().
	DocumentListInfos() []L
	// OtherDocumentInfos returns a list of other document infos. Port of
	// getOtherDocumentInfos().
	OtherDocumentInfos() []D
}
