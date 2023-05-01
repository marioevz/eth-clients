package validator

type ValidatorKeys struct {
	// ValidatorKeystoreJSON encodes an EIP-2335 keystore, serialized in JSON
	// https://github.com/ethereum/EIPs/blob/master/EIPS/eip-2335.md
	ValidatorKeystoreJSON []byte
	// ValidatorKeystorePass holds the secret used for ValidatorKeystoreJSON
	ValidatorKeystorePass string
	// ValidatorSecretKey is the serialized secret key for validator duties
	ValidatorSecretKey [32]byte
	// ValidatorSecretKey is the serialized pubkey derived from ValidatorSecretKey
	ValidatorPubkey [48]byte
	// WithdrawalSecretKey is the serialized secret key for withdrawing stake
	WithdrawalSecretKey [32]byte
	// WithdrawalPubkey is the serialized pubkey derived from WithdrawalSecretKey
	WithdrawalPubkey [48]byte
}
