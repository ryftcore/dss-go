// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/job/DocumentListInfo.java (DSS 6.5.RC1).
package job

// DocumentListInfo contains a validation result for a document list. P is the parent
// DocumentInfo type and C the child DocumentInfo type, mirroring Java's
// "DocumentListInfo<P extends DocumentInfo<P>, C extends DocumentInfo<P>>".
type DocumentListInfo[P any, C any] interface {
	DocumentInfo[P]

	// ChildrenInfos returns a list of DocumentInfo summaries for documents referenced from
	// the current document. Port of getChildrenInfos().
	ChildrenInfos() []C
}
