# 3PVSS reproduction

This package implements the individual-commitment construction
`Lambda_RO^indv` and the binary multi-candidate voting protocol in Baghery,
Makri, and Suryakari, *Packed Pre-Constructed PVSS for Randomness Generation
and E-Voting*.

`Share`, `Verify`, `OptimisticReconstruct`, and `PessimisticReconstruct`
implement the complete native `Lambda_RO^indv` interface. The pessimistic path
decrypts a qualified set of `r=f+l` shares, executes each DLEQ proof, and
interpolates every packed secret encoding. `ShareBinary` and `VerifyBinary`
execute the packed PDL proof and one
Fiat--Shamir binary OR proof per coordinate. `ShareScore` and `VerifyScore` are
the explicitly labeled black-box score adaptation: a score in `[0,2^kappa)`
is decomposed into `kappa` independently shared binary planes. `Tally`
homomorphically accumulates each plane, verifies reconstruction shares, and
performs bounded discrete-log decoding before binary weighting.

The implementation uses the same BN254 G1 backend as the other baselines.
`PublicBallotBytes` follows the communication accounting stated in the 3PVSS
paper, so figures compare the cited wire format rather than Go object sizes.
