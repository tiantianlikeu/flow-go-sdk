package crypto

import (
	"github.com/onflow/flow-go-sdk/crypto/internal/crypto"
	"github.com/onflow/flow-go-sdk/crypto/internal/crypto/hash"
)

// SignatureAlgorithm is an identifier for a signature algorithm (and parameters if applicable).
type SignatureAlgorithm int

const (
	// Supported signature algorithms
	UnknownSignatureAlgorithm SignatureAlgorithm = iota
	// BLS_BLS12381 is BLS on BLS 12-381 curve
	BLS_BLS12381
	// ECDSA_P256 is ECDSA on NIST P-256 curve
	ECDSA_P256
	// ECDSA_secp256k1 is ECDSA on secp256k1 curve
	ECDSA_secp256k1
)

// String returns the string representation of this signing algorithm.
func (f SignatureAlgorithm) String() string {
	return [...]string{"UNKNOWN", "BLS_BLS12381", "ECDSA_P256", "ECDSA_secp256k1"}[f]
}

// HashAlgorithm is an identifier for a hashing algorithm.
type HashAlgorithm int

const (
	// Supported hash algorithms
	UnknownHashAlgorithm HashAlgorithm = iota
	SHA2_256
	SHA2_384
	SHA3_256
	SHA3_384
)

// String returns the string representation of this hashing algorithm.
func (f HashAlgorithm) String() string {
	return [...]string{"UNKNOWN", "SHA2_256", "SHA2_384", "SHA3_256", "SHA3_384"}[f]
}

const (
	MinSeedLengthECDSA_P256      = crypto.KeyGenSeedMinLenECDSAP256
	MinSeedLengthECDSA_secp256k1 = crypto.KeyGenSeedMinLenECDSASecp256k1
)

// KeyType is a key format supported by Flow.
type KeyType int

const (
	UnknownKeyType KeyType = iota
	ECDSA_P256_SHA2_256
	ECDSA_P256_SHA3_256
	ECDSA_secp256k1_SHA2_256
	ECDSA_secp256k1_SHA3_256
)

// SignatureAlgorithm returns the signature algorithm for this key type.
func (k KeyType) SignatureAlgorithm() SignatureAlgorithm {
	switch k {
	case ECDSA_P256_SHA2_256, ECDSA_P256_SHA3_256:
		return ECDSA_P256
	case ECDSA_secp256k1_SHA2_256, ECDSA_secp256k1_SHA3_256:
		return ECDSA_secp256k1
	default:
		return UnknownSignatureAlgorithm
	}
}

// HashAlgorithm returns the hash algorithm for this key type.
func (k KeyType) HashAlgorithm() HashAlgorithm {
	switch k {
	case ECDSA_P256_SHA2_256, ECDSA_secp256k1_SHA2_256:
		return SHA2_256
	case ECDSA_P256_SHA3_256, ECDSA_secp256k1_SHA3_256:
		return SHA3_256
	default:
		return UnknownHashAlgorithm
	}
}

type PrivateKey struct {
	PrivateKey crypto.PrivateKey
}

func (pk PrivateKey) Sign(message []byte, hasher Hasher) ([]byte, error) {
	return pk.PrivateKey.Sign(message, hasher)
}

func (pk PrivateKey) Algorithm() SignatureAlgorithm {
	return SignatureAlgorithm(pk.PrivateKey.Algorithm())
}

func (pk PrivateKey) PublicKey() PublicKey {
	return PublicKey{PublicKey: pk.PrivateKey.PublicKey()}
}

func (pk PrivateKey) Encode() []byte {
	return pk.PrivateKey.Encode()
}

type PublicKey struct {
	PublicKey crypto.PublicKey
}

func (pk PublicKey) Algorithm() SignatureAlgorithm {
	return SignatureAlgorithm(pk.PublicKey.Algorithm())
}

func (pk PublicKey) Encode() []byte {
	return pk.PublicKey.Encode()
}

type Signer interface {
	Sign(message []byte) ([]byte, error)
}

type NaiveSigner struct {
	PrivateKey PrivateKey
	Hasher     Hasher
}

func NewNaiveSigner(privateKey PrivateKey, hashAlgo HashAlgorithm) NaiveSigner {
	return NaiveSigner{
		PrivateKey: privateKey,
		Hasher:     NewHasher(hashAlgo),
	}
}

func (s NaiveSigner) Sign(message []byte) ([]byte, error) {
	return s.PrivateKey.Sign(message, s.Hasher)
}

type MockSigner []byte

func (s MockSigner) Sign(message []byte) ([]byte, error) {
	return s, nil
}

func GeneratePrivateKey(sigAlgo SignatureAlgorithm, seed []byte) (PrivateKey, error) {
	privKey, err := crypto.GeneratePrivateKey(crypto.SigningAlgorithm(sigAlgo), seed)
	if err != nil {
		return PrivateKey{}, err
	}

	return PrivateKey{
		PrivateKey: privKey,
	}, nil
}

func DecodePrivateKey(sigAlgo SignatureAlgorithm, b []byte) (PrivateKey, error) {
	privKey, err := crypto.DecodePrivateKey(crypto.SigningAlgorithm(sigAlgo), b)
	if err != nil {
		return PrivateKey{}, err
	}

	return PrivateKey{
		PrivateKey: privKey,
	}, nil
}

func DecodePublicKey(sigAlgo SignatureAlgorithm, b []byte) (PublicKey, error) {
	pubKey, err := crypto.DecodePublicKey(crypto.SigningAlgorithm(sigAlgo), b)
	if err != nil {
		return PublicKey{}, err
	}

	return PublicKey{
		PublicKey: pubKey,
	}, nil
}

type Hasher = hash.Hasher

func NewHasher(algo HashAlgorithm) Hasher {
	switch algo {
	case SHA2_256:
		return NewSHA2_256()
	case SHA2_384:
		return NewSHA2_384()
	case SHA3_256:
		return NewSHA3_256()
	case SHA3_384:
		return NewSHA3_384()
	default:
		panic("invalid hash algorithm")
	}
}

// NewSHA2_256 returns a new instance of SHA2-256 hasher.
func NewSHA2_256() Hasher {
	return hash.NewSHA2_256()
}

// NewSHA2_384 returns a new instance of SHA2-384 hasher.
func NewSHA2_384() Hasher {
	return hash.NewSHA2_384()
}

// NewSHA3_256 returns a new instance of SHA3-256 hasher.
func NewSHA3_256() Hasher {
	return hash.NewSHA3_256()
}

// NewSHA3_384 returns a new instance of SHA3-384 hasher.
func NewSHA3_384() Hasher {
	return hash.NewSHA3_384()
}