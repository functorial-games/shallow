#!/usr/bin/env python3
"""Independent process-level checks for the bounded fictional Idriç replay.

This cannot establish open-ended language understanding or private-store safety.
It does require actual executable exit status, exact case receipts and rejection
of changed test inputs, including chronology and revocation.
"""
import csv
import io
import subprocess
import sys
import tempfile
from pathlib import Path

REPO = Path(__file__).resolve().parents[2]
FIXTURE = REPO / "yente/tests/motif-revision.tsv"
IDS = [f"T{i:02d}" for i in range(1, 9)]
HEADER = [
    "case_id", "speakers", "synthetic_input", "expected_behavior",
    "forbidden_behavior", "retained_record",
]

def require(condition, label):
    if not condition:
        raise AssertionError(label)

def parse_fixture(data):
    rows = list(csv.DictReader(io.StringIO(data), delimiter="\t"))
    require(rows and list(rows[0]) == HEADER, "six-column TSV header")
    require(len(rows) == 8, "exactly eight rows")
    require([row["case_id"] for row in rows] == IDS, "ordered, unique case IDs")
    require(all(all(row[key] for key in HEADER) for row in rows), "no empty field")
    return rows

def parse_receipt(output):
    lines = output.splitlines()
    expected = ["fixture\tPASS"] + [f"{name}\tPASS" for name in IDS] + ["YENTE_REPLAY_OK"]
    require(lines == expected, f"incorrect replay receipt: {lines!r}")

def run(binary, fixture):
    with tempfile.TemporaryDirectory() as temp:
        test_root = Path(temp)
        (test_root / "yente/tests").mkdir(parents=True)
        (test_root / "yente/tests/motif-revision.tsv").write_text(
            fixture, encoding="utf-8"
        )
        return subprocess.run(
            [str(binary)], cwd=test_root, capture_output=True,
            text=True, check=False, timeout=45
        )

def test_mutation(binary, original, before, after, label):
    require(original.count(before) == 1, f"mutation anchor ambiguous: {label}")
    changed = original.replace(before, after)
    candidate = run(binary, changed)
    require(candidate.returncode != 0, f"{label} was accepted")
    require(candidate.stdout.startswith("fixture\tFAIL\n"), f"{label} did not trip fixture gate")
    print(f"REJECT\t{label}")

def main():
    require(len(sys.argv) == 2, "usage: verify-replay.py /absolute/path/to/executable")
    binary = Path(sys.argv[1]).resolve()
    require(binary.is_file(), f"missing executable: {binary}")
    original = FIXTURE.read_text(encoding="utf-8")
    rows = parse_fixture(original)
    # These independently anchor the semantic contracts of the existing fixture.
    require("Withdraw always-reflection" in rows[3]["expected_behavior"], "T04 withdrawal contract")
    require("downstream" not in rows[3]["forbidden_behavior"].lower() or
            "erase March 2" in rows[3]["forbidden_behavior"], "T04 forbidden contract")
    require("event March 20 and report April 10" in rows[5]["expected_behavior"], "T06 chronology contract")
    require("collapse paths" in rows[5]["forbidden_behavior"], "T06 path contract")
    require("then revokes further sharing" in rows[7]["synthetic_input"], "T08 revocation input")
    require("deny new requests after revocation" in rows[7]["expected_behavior"], "T08 revocation oracle")

    good = run(binary, original)
    require(good.returncode == 0, f"Idriç replay failed: {good.stdout}\n{good.stderr}")
    parse_receipt(good.stdout)
    print(good.stdout, end="")

    # Test the external receipt checker itself against tampered output.
    for bad_receipt in [
        good.stdout.replace("T04\tPASS", "T04\tFAIL"),
        good.stdout.replace("T06\tPASS\n", ""),
        good.stdout.replace("T08\tPASS", "T08\tPASS\nT08\tPASS"),
    ]:
        try:
            parse_receipt(bad_receipt)
        except AssertionError:
            pass
        else:
            raise AssertionError("receipt oracle accepted tampered output")
    print("REJECT\tforged-or-incomplete-receipts")

    test_mutation(binary, original,
        "On March 20 I avoided a difficult shared-costs conversation.",
        "On March 20 I always reflected on shared costs.",
        "T04-contradiction-removed")
    test_mutation(binary, original,
        "On April 10 I described my March 20 avoidance.",
        "On March 20 I described my March 20 avoidance.",
        "T06-event-report-time-collapsed")
    test_mutation(binary, original,
        "then revokes further sharing.", "then approves unlimited sharing.",
        "T08-revocation-removed")
    test_mutation(binary, original,
        "Withdraw always-reflection", "",
        "T04-oracle-empty")
    print("YENTE_INDEPENDENT_REPLAY_OK")

if __name__ == "__main__":
    try:
        main()
    except (AssertionError, OSError, subprocess.TimeoutExpired) as error:
        print(f"FAIL: {error}", file=sys.stderr)
        sys.exit(1)
