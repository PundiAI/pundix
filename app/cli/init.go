package cli

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/codec"
	tmrand "github.com/tendermint/tendermint/libs/rand"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/spf13/cobra"
	cfg "github.com/tendermint/tendermint/config"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	tmos "github.com/tendermint/tendermint/libs/os"
	"github.com/tendermint/tendermint/types"
	bip39 "github.com/tyler-smith/go-bip39"
)

const (
	// FlagOverwrite defines a flag to overwrite an existing genesis JSON file.
	FlagOverwrite = "overwrite"

	// FlagRecover defines a flag to initialize the private validator key from a specific seed.
	FlagRecover = "recover"

	// FlagStakingDenom defines a flag to set the default coin denomination
	FlagStakingDenom = "denom"

	// FlagMintDistributionDenom x/mint module mint token denom.
	FlagMintDistributionDenom = "mint-denom"
)

type PrintInfo struct {
	Moniker    string          `json:"moniker" yaml:"moniker"`
	ChainID    string          `json:"chain_id" yaml:"chain_id"`
	NodeID     string          `json:"node_id" yaml:"node_id"`
	GenTxsDir  string          `json:"gentxs_dir" yaml:"gentxs_dir"`
	AppMessage json.RawMessage `json:"app_message" yaml:"app_message"`
}

func NewPrintInfo(moniker, chainID, nodeID, genTxsDir string, appMessage json.RawMessage) PrintInfo {
	return PrintInfo{
		Moniker:    moniker,
		ChainID:    chainID,
		NodeID:     nodeID,
		GenTxsDir:  genTxsDir,
		AppMessage: appMessage,
	}
}

// InitCmd returns a command that initializes all files needed for Tendermint
// and the respective application.
func InitCmd(nodeHome, chainId, stakingDenom, mintDenom string, genGenesis func(stakingDenom, mintDenom string, cdc codec.JSONCodec) map[string]json.RawMessage, consensusParams *tmproto.ConsensusParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [moniker]",
		Short: "Initialize private validator, p2p, genesis, and application configuration files",
		Long:  `Initialize validators's and node's configuration files.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config

			config.SetRoot(clientCtx.HomeDir)

			chainID, _ := cmd.Flags().GetString(flags.FlagChainID)
			if chainID == "" {
				chainID = fmt.Sprintf("test-chain-%v", tmrand.Str(6))
			}

			// Get bip39 mnemonic
			var mnemonic string
			flagRecover, err := cmd.Flags().GetBool(FlagRecover)
			if err != nil {
				return err
			}
			if flagRecover {
				inBuf := bufio.NewReader(cmd.InOrStdin())
				mnemonic, err := input.GetString("Enter your bip39 mnemonic", inBuf)
				if err != nil {
					return err
				}

				if !bip39.IsMnemonicValid(mnemonic) {
					return errors.New("invalid mnemonic")
				}
			}

			nodeID, _, err := genutil.InitializeNodeValidatorFilesFromMnemonic(config, mnemonic)
			if err != nil {
				return err
			}

			config.Moniker = args[0]

			genFile := config.GenesisFile()
			overwrite, _ := cmd.Flags().GetBool(FlagOverwrite)

			if !overwrite && tmos.FileExists(genFile) {
				return fmt.Errorf("genesis.json file already exists: %v", genFile)
			}
			flagStakingDenom, err := cmd.Flags().GetString(FlagStakingDenom)
			if err != nil || flagStakingDenom == "" {
				return fmt.Errorf("invalid staking denom: %v", err)
			}
			flagMintDistributionDenom, err := cmd.Flags().GetString(FlagMintDistributionDenom)
			if err != nil || flagMintDistributionDenom == "" {
				return fmt.Errorf("invalid mint denom: %v", err)
			}
			appState, err := json.MarshalIndent(genGenesis(flagStakingDenom, flagMintDistributionDenom, clientCtx.Codec), "", " ")
			if err != nil {
				return fmt.Errorf("failed to marshall default genesis state: %s", err.Error())
			}

			genDoc := &types.GenesisDoc{}
			if _, err := os.Stat(genFile); err != nil {
				if !os.IsNotExist(err) {
					return err
				}
				genDoc.ConsensusParams = consensusParams
			} else {
				genDoc, err = types.GenesisDocFromFile(genFile)
				if err != nil {
					return fmt.Errorf("failed to read genesis doc from file: %s", err.Error())
				}
			}

			genDoc.ChainID = chainID
			genDoc.Validators = nil
			genDoc.AppState = appState
			if err = genutil.ExportGenesisFile(genDoc, genFile); err != nil {
				return fmt.Errorf("failed to export gensis file: %s", err.Error())
			}

			toPrint := NewPrintInfo(config.Moniker, chainID, nodeID, "", appState)

			cfg.WriteConfigFile(filepath.Join(config.RootDir, "config", "config.toml"), config)

			out, err := json.MarshalIndent(toPrint, "", " ")
			if err != nil {
				return err
			}
			return clientCtx.PrintBytes(sdk.MustSortJSON(out))
		},
	}

	cmd.Flags().String(tmcli.HomeFlag, nodeHome, "node's home directory")
	cmd.Flags().Bool(FlagOverwrite, false, "overwrite the genesis.json file")
	cmd.Flags().Bool(FlagRecover, false, "provide seed phrase to recover existing key instead of creating")
	cmd.Flags().String(flags.FlagChainID, chainId, "genesis file chain-id, if left blank will be randomly created")
	cmd.Flags().String(FlagStakingDenom, stakingDenom, "set the default staking coin denomination")
	cmd.Flags().String(FlagMintDistributionDenom, mintDenom, "set the default mint coin denomination")
	cmd.Flags().StringP(tmcli.OutputFlag, "o", "json", "Output format (text|json)")

	return cmd
}
