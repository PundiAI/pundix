package params

import (
	"github.com/cosmos/cosmos-sdk/server/config"
)

const CustomConfigTemplate = `
###############################################################################
###                        Custom Fx Configuration                        ###
###############################################################################
[bypass-min-fee]
# MsgTypes defines custom message types the operator may set that will bypass minimum fee checks during CheckTx.
# Example:
# ["/ibc.core.channel.v1.MsgRecvPacket", "/ibc.core.channel.v1.MsgAcknowledgement", ...]
msg-types = [{{ range .BypassMinFee.MsgTypes }}{{ printf "%q, " . }}{{end}}]

# MsgMaxGasUsage defines gas consumption threshold .Default: 300000
msg-max-gas-usage = {{ .BypassMinFee.MsgMaxGasUsage }}
`

// BypassMinFee defines custom that will bypass minimum fee checks during CheckTx.
type BypassMinFee struct {
	// MsgTypes defines custom message types the operator may set that
	// will bypass minimum fee checks during CheckTx.
	MsgTypes       []string `mapstructure:"msg-types"`
	MsgMaxGasUsage uint64   `mapstructure:"msg-max-gas-usage"`
}

// DefaultBypassMinFee returns the default BypassMinFee configuration
func DefaultBypassMinFee() BypassMinFee {
	return BypassMinFee{
		MsgTypes:       []string{},
		MsgMaxGasUsage: uint64(300_000),
	}
}

type Config struct {
	config.Config `mapstructure:",squash"`

	// BypassMinFeeMsgTypes defines custom that will bypass minimum fee checks during CheckTx.
	BypassMinFee BypassMinFee `mapstructure:"bypass-min-fee"`
}

// DefaultConfig returns server's default configuration.
func DefaultConfig() *Config {
	return &Config{
		Config:       *config.DefaultConfig(),
		BypassMinFee: DefaultBypassMinFee(),
	}
}

func DefaultConfigTemplate() string {
	return config.DefaultConfigTemplate + CustomConfigTemplate
}
