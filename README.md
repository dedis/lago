# LAGO: Lattice Cryptography Library in Golang

This package provides a toolbox of lattice-based cryptographic primitives for Go. The library is still at an experimental stage and should be used for research purposes only.

The LAGO subpackages from the lowest to the highest abstraction level and their provided functionalities are as follows:

- `bigint`: Modular arithmetic operations for big integers.
- `polynomial`: Modular arithmetic operations for polynomials, Number Theoretic Transformation (NTT).
- `ring`: Modular arithmetic operations for polynomials over rings, Gaussian sampling.
- `crypto`: Fan-Vercauteren (FV) homomorphic encryption/decryption.
- `encoding`: Encode/decode messages to/from plaintexts.

## Examples

[main.go](https://github.com/dedis/student_18_lattices/blob/master/main.go) gives an example on how to use this library.
In each subpackage you can find additional test files documenting further usage approaches.

## License

The LAGO Source code is released under MIT license, see the file [LICENSE](https://github.com/dedis/lago/blob/master/LICENSE) for the full text.
