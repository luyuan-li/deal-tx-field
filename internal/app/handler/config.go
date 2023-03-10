package handlers

import (
	"github.com/kaifei-bianjie/common-parser/codec"
	"github.com/luyuan-li/deal-tx-field/internal/app/config"
)

const (

	// Bech32ChainPrefix defines the prefix of this chain
	Bech32ChainPrefix = "i"

	// PrefixAcc is the prefix for account
	PrefixAcc = "a"

	// PrefixValidator is the prefix for validator keys
	PrefixValidator = "v"

	// PrefixConsensus is the prefix for consensus keys
	PrefixConsensus = "c"

	// PrefixPublic is the prefix for public
	PrefixPublic = "p"

	// PrefixAddress is the prefix for address
	PrefixAddress = "a"
)

func initBech32Prefix(conf *config.Config) {
	var (
		// Bech32PrefixAccAddr defines the Bech32 prefix of an account's address
		Bech32PrefixAccAddr = conf.Server.Bech32AccPrefix
		// Bech32PrefixAccPub defines the Bech32 prefix of an account's public key
		Bech32PrefixAccPub = conf.Server.Bech32AccPrefix + PrefixAcc + PrefixPublic
		// Bech32PrefixValAddr defines the Bech32 prefix of a validator's operator address
		Bech32PrefixValAddr = conf.Server.Bech32AccPrefix + PrefixValidator + PrefixAddress
		// Bech32PrefixValPub defines the Bech32 prefix of a validator's operator public key
		Bech32PrefixValPub = conf.Server.Bech32AccPrefix + PrefixValidator + PrefixPublic
		// Bech32PrefixConsAddr defines the Bech32 prefix of a consensus node address
		Bech32PrefixConsAddr = conf.Server.Bech32AccPrefix + PrefixConsensus + PrefixAddress
		// Bech32PrefixConsPub defines the Bech32 prefix of a consensus node public key
		Bech32PrefixConsPub = conf.Server.Bech32AccPrefix + PrefixConsensus + PrefixPublic
	)
	codec.SetBech32Prefix(Bech32PrefixAccAddr, Bech32PrefixAccPub, Bech32PrefixValAddr,
		Bech32PrefixValPub, Bech32PrefixConsAddr, Bech32PrefixConsPub)
}
