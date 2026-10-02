# PriScore reproduction

This package implements the protocol interfaces and Algorithms 1--4 in
Yang et al., *PriScore: Blockchain-Based Self-Tallying Election System
Supporting Score Voting* (IEEE TIFS 2021):

`Setup`, `KeyGen`, `KeyDerive`, `Commit`, `VerifyCommit`, `Vote`,
`VerifyVote`, `SelfTally`, `RecoverForAbsent`, `VerifyRecovery`, and
`TallyWithAbort`.

The paper's multiplicative group is instantiated with BN254 G1. Its
one-out-of-`(B+1)` statements are Fiat--Shamir OR proofs, and witnesses are
shared between the commitment and encrypted-ballot equations. The sum proofs,
signed pairwise masks, abort correction values, and correction proofs are all
executed rather than represented by operation counters.

`priscore_test.go` runs the normal and abort tally paths and rejects a modified
ballot. The cross-protocol runner is `cmd/full-comparison-bench`.
