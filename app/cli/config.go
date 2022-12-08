package cli

import (
	"fmt"
	"path/filepath"
	"strings"

	appparams "github.com/pundix/pundix/app/params"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/config"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tmcfg "github.com/tendermint/tendermint/config"
	tmcli "github.com/tendermint/tendermint/libs/cli"
)

const (
	configFileName = "config.toml"
	appFileName    = "app.toml"
)

func UpdateCfgCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update app.toml and config.toml files to the latest version, default only missing parts are added",
		Args:  cobra.RangeArgs(0, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			serverCtx := server.GetServerContextFromCmd(cmd)
			rootDir := serverCtx.Config.RootDir
			fileName := filepath.Join(rootDir, "config", configFileName)
			tmcfg.WriteConfigFile(fileName, serverCtx.Config)
			serverCtx.Logger.Info("Update config.toml is successful", "fileName", fileName)

			config.SetConfigTemplate(appparams.DefaultConfigTemplate())
			appConfig := appparams.DefaultConfig()
			if err := serverCtx.Viper.Unmarshal(appConfig); err != nil {
				return err
			}
			fileName = filepath.Join(rootDir, "config", appFileName)
			config.WriteConfigFile(fileName, appConfig)
			serverCtx.Logger.Info("Update app.toml is successful", "fileName", fileName)
			return nil
		},
	}
	return cmd
}

func AppTomlCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "app.toml [key] [value]",
		Short: "Create or query an `config/app.toml` file",
		Args:  cobra.RangeArgs(0, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			config.SetConfigTemplate(appparams.DefaultConfigTemplate())
			return runConfigCmd(cmd, append([]string{appFileName}, args...))
		},
	}
	cmd.Flags().StringP(tmcli.OutputFlag, "o", "text", "Output format (text|json)")
	return cmd
}

func ConfigTomlCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config.toml [key] [value]",
		Short: "Create or query an `config/config.toml` file",
		Args:  cobra.RangeArgs(0, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigCmd(cmd, append([]string{configFileName}, args...))
		},
	}
	cmd.Flags().StringP(tmcli.OutputFlag, "o", "text", "Output format (text|json)")
	return cmd
}

func runConfigCmd(cmd *cobra.Command, args []string) error {
	serverCtx := server.GetServerContextFromCmd(cmd)
	clientCtx := client.GetClientContextFromCmd(cmd)

	configName := filepath.Join(serverCtx.Config.RootDir, "config", args[0])
	cfg, err := newConfig(serverCtx.Viper, configName)
	if err != nil {
		return err
	}

	// is len(args) == 1, get config file content
	if len(args) == 1 {
		return cfg.output(clientCtx)
	}

	// 2. is len(args) == 2, get config key value
	if len(args) == 2 {
		return output(clientCtx, serverCtx.Viper.Get(args[1]))
	}

	serverCtx.Viper.Set(args[1], args[2])
	return cfg.save()
}

type cmdConfig interface {
	save() error
	output(ctx client.Context) error
}

var (
	_ cmdConfig = &appTomlConfig{}
	_ cmdConfig = &configTomlConfig{}
)

type appTomlConfig struct {
	v          *viper.Viper
	config     *appparams.Config
	configName string
}

func (a *appTomlConfig) output(ctx client.Context) error {
	return output(ctx, a.config)
}

func (a *appTomlConfig) save() error {
	if err := a.v.Unmarshal(a.config, func(decoderConfig *mapstructure.DecoderConfig) {
		decoderConfig.ZeroFields = true
	}); err != nil {
		return err
	}
	config.WriteConfigFile(a.configName, a.config)
	return nil
}

type configTomlConfig struct {
	v          *viper.Viper
	config     *tmcfg.Config
	configName string
}

func (c *configTomlConfig) output(ctx client.Context) error {
	type outputConfig struct {
		tmcfg.BaseConfig `mapstructure:",squash"`
		RPC              tmcfg.RPCConfig             `mapstructure:"rpc"`
		P2P              tmcfg.P2PConfig             `mapstructure:"p2p"`
		Mempool          tmcfg.MempoolConfig         `mapstructure:"mempool"`
		StateSync        tmcfg.StateSyncConfig       `mapstructure:"statesync"`
		FastSync         tmcfg.FastSyncConfig        `mapstructure:"fastsync"`
		Consensus        tmcfg.ConsensusConfig       `mapstructure:"consensus"`
		Storage          tmcfg.StorageConfig         `mapstructure:"storage"`
		TxIndex          tmcfg.TxIndexConfig         `mapstructure:"tx_index"`
		Instrumentation  tmcfg.InstrumentationConfig `mapstructure:"instrumentation"`
	}
	return output(ctx, outputConfig{
		BaseConfig:      c.config.BaseConfig,
		RPC:             *c.config.RPC,
		P2P:             *c.config.P2P,
		Mempool:         *c.config.Mempool,
		StateSync:       *c.config.StateSync,
		FastSync:        *c.config.FastSync,
		Consensus:       *c.config.Consensus,
		Storage:         *c.config.Storage,
		TxIndex:         *c.config.TxIndex,
		Instrumentation: *c.config.Instrumentation,
	})
}

func (c *configTomlConfig) save() error {
	if err := c.v.Unmarshal(c.config, func(decoderConfig *mapstructure.DecoderConfig) {
		decoderConfig.ZeroFields = true
	}); err != nil {
		return err
	}
	tmcfg.WriteConfigFile(c.configName, c.config)
	return nil
}

func newConfig(v *viper.Viper, configName string) (cmdConfig, error) {
	if strings.HasSuffix(configName, appFileName) {
		var configData = appparams.Config{}
		if err := v.Unmarshal(&configData); err != nil {
			return nil, err
		}
		return &appTomlConfig{config: &configData, v: v, configName: configName}, nil
	} else if strings.HasSuffix(configName, configFileName) {
		var configData = tmcfg.Config{}
		if err := v.Unmarshal(&configData); err != nil {
			return nil, err
		}
		return &configTomlConfig{config: &configData, v: v, configName: configName}, nil
	} else {
		return nil, fmt.Errorf("invalid config file: %s, (support: %v)", configName, strings.Join([]string{appFileName, configFileName}, "/"))
	}
}

func output(ctx client.Context, content interface{}) error {
	var mapData map[string]interface{}
	if err := mapstructure.Decode(content, &mapData); err != nil {
		var data interface{}
		if err := mapstructure.Decode(content, &data); err != nil {
			return err
		}
		return PrintOutput(ctx, data)
	}
	return PrintOutput(ctx, mapData)
}
