// Ported from dss-token/src/main/java/eu/europa/esig/dss/token/predicate/DSSKeyEntryPredicate.java (DSS 6.5.RC1).
package token

// DSSKeyEntryPredicate filters DSSPrivateKeyEntry values considered by
// AbstractKeyStoreTokenConnection.Keys.
//
// DEVIATION: Java declares a marker interface extending java.util.function.Predicate<T>. Go
// represents that single-method functional interface directly as a function type, its Test(T)
// method becoming a plain call.
type DSSKeyEntryPredicate func(entry DSSPrivateKeyEntry) bool
