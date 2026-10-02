# Huang et al. tally baseline

This package implements only the tally-relevant path described by Huang et al.:
per-coordinate exponential-ElGamal encryption, public aggregation of the active
ballots, election-key decryption, and bounded discrete-log decoding. It is used
to measure the cost of rescanning all accepted ciphertexts at the deadline.

It does **not** claim to implement the paper's complete ballot proof,
registration, or threshold abstention-recovery protocol. Consequently, the
paper labels this curve as a tally-path experiment rather than a complete
protocol benchmark.
