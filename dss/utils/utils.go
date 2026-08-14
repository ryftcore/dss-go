// Ported from dss-utils/src/main/java/eu/europa/esig/dss/utils/Utils.java
// and dss-utils/src/main/java/eu/europa/esig/dss/utils/IUtils.java (DSS 6.5.RC1),
// with semantics cross-checked against the apache-commons adapter at
// dss-utils-apache-commons/.../ApacheCommonsUtils.java.
//
// In upstream DSS, Utils is a static facade over an IUtils implementation
// loaded via ServiceLoader (apache-commons or guava backed). This package
// collapses that indirection: every method is implemented directly against
// the Go standard library as a package-level function. There is no
// interface and no registry.
//
// Package utils is dependency-free (stdlib only).
package utils

// EmptyString is the empty string constant (Java: Utils.EMPTY_STRING).
const EmptyString = ""
