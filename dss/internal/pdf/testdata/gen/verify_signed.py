#!/usr/bin/env python3
"""verify_signed.py — the Java side of the writer check (DESIGN.md §6.3 KAT-C,
Go -> Java direction).

TestWriterOverCorpus appends a signature increment to twelve representative
upstream documents and asserts the Go-side invariants itself. This script adds
pdfbox's opinion of the same bytes: that it parses the result at all, that it
sees exactly one more signature than the input had, that the appended
signature's /ByteRange is the one the file's bytes actually justify, and that
none of the prior signatures or revisions moved.

Run BY HAND, like PdfOracle.java. Full procedure:

    SP=/tmp/pdfaudit && mkdir -p $SP
    M2=$HOME/.m2/repository
    CP=$M2/org/apache/pdfbox/pdfbox/3.0.7/pdfbox-3.0.7.jar\\
    :$M2/org/apache/pdfbox/pdfbox-io/3.0.7/pdfbox-io-3.0.7.jar\\
    :$M2/org/apache/pdfbox/fontbox/3.0.7/fontbox-3.0.7.jar\\
    :$M2/commons-logging/commons-logging/1.3.5/commons-logging-1.3.5.jar
    export LANG=C.UTF-8 LC_ALL=C.UTF-8   # a corpus filename is non-ASCII

    # 1. oracle over the untouched upstream corpus
    javac -nowarn -cp "$CP" -d $SP/classes PdfOracle.java
    find /home/user/dss-upstream -name '*.pdf' -not -path '*/target/*' \\
        | sed 's|^/home/user/dss-upstream/||' | sort > $SP/list.txt
    java -Dfile.encoding=UTF-8 -Dsun.jnu.encoding=UTF-8 \\
         -Dorg.apache.commons.logging.Log=org.apache.commons.logging.impl.NoOpLog \\
         -cp "$CP:$SP/classes" PdfOracle \\
         /home/user/dss-upstream $SP/before.tsv $SP/before.man $SP/list.txt

    # 2. Go writer produces the signed outputs
    cd /home/user/esig/dss
    PDF_CORPUS_DIR=/home/user/dss-upstream PDF_WRITE_OUT=$SP/signed \\
        go test ./internal/pdf/ -run TestWriterOverCorpus -count=1 -v

    # 3. oracle over the signed outputs
    (cd $SP/signed && ls > $SP/signed.txt)
    java -Dfile.encoding=UTF-8 -Dsun.jnu.encoding=UTF-8 \\
         -Dorg.apache.commons.logging.Log=org.apache.commons.logging.impl.NoOpLog \\
         -cp "$CP:$SP/classes" PdfOracle \\
         $SP/signed $SP/after.tsv $SP/after.man $SP/signed.txt

    # 4. this script
    python3 verify_signed.py $SP/before.tsv $SP/after.tsv $SP/signed

Exit status is 0 only when every signed document passes every check.
"""

import os
import re
import sys

# A signature record in the oracle's `sigs` field is "<name>:EMPTY" or
# "<name>:valueKey=...,covers=<bool>". Records are space-separated, but a
# record's own /ByteRange contains spaces, so records cannot be split naively.
SIG_RECORD = re.compile(
    r"(?:^| )([^ :]*):(EMPTY|valueKey=.*?covers=(?:true|false))(?= [^ :]*:|$)"
)
BYTE_RANGE = re.compile(r"BR=\[([0-9 ]+)\]")


def load(path):
    """Read an oracle TSV into {path: {field: value}}."""
    out = {}
    with open(path, encoding="utf-8") as fh:
        for line in fh:
            rec = dict(
                f.split("=", 1) for f in line.rstrip("\n").split("\t") if "=" in f
            )
            out[rec["path"]] = rec
    return out


def signatures(field):
    return SIG_RECORD.findall(field or "")


def check(name, before, after, data):
    """Return the list of problems with one signed document."""
    problems = []
    sigs_before = signatures(before.get("sigs", ""))
    sigs_after = signatures(after.get("sigs", ""))
    names_before = [n for n, _ in sigs_before]

    if len(sigs_after) != len(sigs_before) + 1:
        problems.append(
            "signature count %d -> %d, want +1" % (len(sigs_before), len(sigs_after))
        )

    # --- the appended signature ------------------------------------------
    new = [(n, v) for n, v in sigs_after if n not in names_before]
    if len(new) != 1:
        problems.append("%d signatures are new, want exactly 1" % len(new))
    else:
        value = new[0][1]
        match = BYTE_RANGE.search(value)
        if not match:
            problems.append("the appended signature has no /ByteRange")
        else:
            br = [int(x) for x in match.group(1).split()]
            if len(br) != 4:
                problems.append("/ByteRange has %d values, want 4" % len(br))
            else:
                gap_start, gap_end = br[0] + br[1], br[2]
                if br[0] != 0:
                    problems.append("/ByteRange[0] is %d, an append must start at 0" % br[0])
                if br[2] + br[3] != len(data):
                    problems.append(
                        "/ByteRange ends at %d, the file is %d bytes"
                        % (br[2] + br[3], len(data))
                    )
                elif data[gap_start : gap_start + 1] != b"<" or data[gap_end - 1 : gap_end] != b">":
                    problems.append("the /ByteRange gap does not bracket a hex string")
        if "covers=true" not in value:
            problems.append("pdfbox does not report the signature as covering the document")
        if "SubFilter=ETSI.CAdES.detached" not in value:
            problems.append("/SubFilter is not ETSI.CAdES.detached")
        if "Filter=Adobe.PPKLite" not in value:
            problems.append("/Filter is not Adobe.PPKLite")

    # --- everything that was already there must not have moved ------------
    by_name_after = dict(sigs_after)
    for sig_name, value in sigs_before:
        if sig_name not in by_name_after:
            problems.append("prior signature field %r vanished" % sig_name)
            continue
        # `covers` is expected to flip: the file grew, so a signature that used
        # to cover the whole document no longer does. Everything else is frozen.
        strip = lambda s: re.sub(r",covers=(?:true|false)$", "", s)
        if strip(value) != strip(by_name_after[sig_name]):
            problems.append("prior signature field %r changed" % sig_name)

    revs_before, revs_after = int(before["eofRevisions"]), int(after["eofRevisions"])
    if revs_after != revs_before + 1:
        problems.append("revisions %d -> %d, want +1" % (revs_before, revs_after))
    ends_before = before.get("revisionEnds", "")
    if ends_before and not after.get("revisionEnds", "").startswith(ends_before + " "):
        problems.append("prior revision boundaries moved")

    lost = set(before.get("objectKeys", "").split()) - set(after.get("objectKeys", "").split())
    if lost:
        problems.append("%d prior object keys are no longer resolvable" % len(lost))

    err = after.get("error", "")
    if err.startswith("!ERROR"):
        problems.append("pdfbox failed to parse our output: %s" % err)

    return problems


def main(argv):
    if len(argv) != 4:
        sys.exit(__doc__)
    before, after, signed_dir = load(argv[1]), load(argv[2]), argv[3]

    failures = 0
    for name in sorted(after):
        # PDF_WRITE_OUT flattens the upstream path by replacing '/' with '_'.
        origin = [p for p in before if p.replace("/", "_") == name]
        if len(origin) != 1:
            print("FAIL %-62s cannot match to an input document" % name[:62])
            failures += 1
            continue
        with open(os.path.join(signed_dir, name), "rb") as fh:
            data = fh.read()
        problems = check(name, before[origin[0]], after[name], data)
        print(
            "%s %-62s %s"
            % ("FAIL" if problems else "ok  ", name[:62], "; ".join(problems))
        )
        failures += bool(problems)

    print("\n%d/%d signed documents verified by pdfbox" % (len(after) - failures, len(after)))
    return 1 if failures else 0


if __name__ == "__main__":
    sys.exit(main(sys.argv))
