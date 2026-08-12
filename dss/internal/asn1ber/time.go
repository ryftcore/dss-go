package asn1ber

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ParseUTCTime ports ASN1UTCTime#getDate, including getAdjustedTime's two-digit year rule: a
// year of 50 or more is 19xx, below 50 it is 20xx.
func ParseUTCTime(value string) (time.Time, error) {
	if len(value) < 10 {
		return time.Time{}, fmt.Errorf("invalid UTCTime: %q", value)
	}
	year, err := strconv.Atoi(value[0:2])
	if err != nil {
		return time.Time{}, err
	}
	century := "20"
	if year >= 50 {
		century = "19"
	}
	return parseTimeDigits(century+value[0:2], value[2:])
}

// ParseGeneralizedTime ports ASN1GeneralizedTime#getDate for the YYYYMMDDHHMM[SS[.f+]] forms,
// with an optional Z or +/-HH[MM] offset.
func ParseGeneralizedTime(value string) (time.Time, error) {
	if len(value) < 10 {
		return time.Time{}, fmt.Errorf("invalid GeneralizedTime: %q", value)
	}
	return parseTimeDigits(value[0:4], value[4:])
}

// parseTimeDigits decodes MMDDHHMM[SS[.f+]][zone] against an already resolved four-digit year.
//
// DEVIATION: a time without a zone is read as UTC. BouncyCastle interprets it in the JVM's
// default time zone, which makes the parsed instant depend on the host configuration;
// certificates and revocation data always carry the "Z" suffix, so the two agree in practice.
func parseTimeDigits(year string, rest string) (time.Time, error) {
	offset := time.Duration(0)
	location := time.UTC
	switch {
	case strings.HasSuffix(rest, "Z"):
		rest = rest[:len(rest)-1]
	default:
		if index := strings.LastIndexAny(rest, "+-"); index > 0 {
			zone := rest[index:]
			rest = rest[:index]
			sign := time.Duration(1)
			if zone[0] == '-' {
				sign = -1
			}
			zone = zone[1:]
			if len(zone) != 2 && len(zone) != 4 && !(len(zone) == 5 && zone[2] == ':') {
				return time.Time{}, fmt.Errorf("invalid time zone: %q", zone)
			}
			hours, err := strconv.Atoi(zone[0:2])
			if err != nil {
				return time.Time{}, err
			}
			minutes := 0
			if len(zone) >= 4 {
				minutes, err = strconv.Atoi(zone[len(zone)-2:])
				if err != nil {
					return time.Time{}, err
				}
			}
			offset = sign * (time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute)
		}
	}

	fraction := time.Duration(0)
	if index := strings.IndexAny(rest, ".,"); index >= 0 {
		digits := rest[index+1:]
		rest = rest[:index]
		// Java keeps millisecond precision; a longer fraction is truncated.
		for len(digits) < 3 {
			digits += "0"
		}
		milliseconds, err := strconv.Atoi(digits[0:3])
		if err != nil {
			return time.Time{}, err
		}
		fraction = time.Duration(milliseconds) * time.Millisecond
	}

	if len(rest) != 6 && len(rest) != 8 && len(rest) != 10 {
		return time.Time{}, fmt.Errorf("invalid time digits: %q", year+rest)
	}
	layoutSource := year + rest
	layout := "2006010215"
	switch len(rest) {
	case 8:
		layout = "200601021504"
	case 10:
		layout = "20060102150405"
	}
	parsed, err := time.ParseInLocation(layout, layoutSource, location)
	if err != nil {
		return time.Time{}, err
	}
	return parsed.Add(fraction).Add(-offset), nil
}
