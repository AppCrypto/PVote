# Context-bound ballot regression tests

Run:

```sh
go test ./crypto/RBPVSS -run 'TestReplay|TestContext|TestTranscriptBoundDual' -count=1 -v
```

Voting callers use `ShareWithContext`, `DVerifyWithContext`, and
`AggregateVerifiedWithContext`. The expected `BallotContext` comes from the
application's election state and authenticated voter index. Its parameter
fingerprint covers all tallier keys, the range-issuer key, score domain, and
public setup. The election identifier must be unique across all elections sharing the
parameters; applications spanning contracts or chains must namespace it.

The tests cover cross-voter copying, cross-election replay with reused keys,
range-proof transplantation, coordinate changes, parameter substitution,
duplicate acceptance, unchanged aggregate state after rejection, valid equal
scores, and completion of an independently generated earlier-closing election.
The dual-code test constructs offsets that pass a fixed old challenge, then
checks rejection after deriving the challenge from the modified transcript.

`replayElection` in the test file is a small model of the authenticated
acceptance boundary. It receives the authenticated index rather than trusting
an index carried inside a ballot. Cryptographic re-verification in the same
context remains valid; the acceptance map prevents counting it twice.

These are Go verifier and acceptance-model tests, not EVM execution or wallet
authentication tests. Existing Solidity bindings and the standalone
`Share`/`DVerify` entry points retain their legacy format; they do not implement
the new context binding. Existing benchmark results therefore are not new
measurements of the context-bound path. These tests provide regression coverage,
not a ballot-privacy proof.
