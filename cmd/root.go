package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/cosmos/cosmos-sdk/version"
	tmcfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/node"
	tmstore "github.com/tendermint/tendermint/store"

	"github.com/pundix/pundix/server/grpc/base/gasprice"

	appparams "github.com/pundix/pundix/app/params"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	sdkCfg "github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/snapshots"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingcli "github.com/cosmos/cosmos-sdk/x/auth/vesting/client/cli"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/pundix/pundix/app"
	"github.com/pundix/pundix/app/cli"
	pxtypes "github.com/pundix/pundix/types"
)

// NewRootCmd creates a new root command for simd. It is called once in the
// main function.
func NewRootCmd() *cobra.Command {
	pxtypes.SetConfig()

	encodingConfig := app.MakeEncodingConfig()
	initClientCtx := client.Context{}.
		WithCodec(encodingConfig.Codec).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithOutput(os.Stdout).
		WithAccountRetriever(types.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastBlock).
		WithHomeDir(app.DefaultNodeHome).
		WithViper("")

	rootCmd := &cobra.Command{
		Use:   pxtypes.Name + "d",
		Short: "PundiX Chain App",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			// set the default command outputs
			cmd.SetOut(cmd.OutOrStdout())
			cmd.SetErr(cmd.ErrOrStderr())

			initClientCtx, err := client.ReadPersistentCommandFlags(initClientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			initClientCtx, err = sdkCfg.ReadFromClientConfig(initClientCtx)
			if err != nil {
				return err
			}

			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}

			if cmd.Name() != "start" && len(initClientCtx.ChainID) > 0 {
				pxtypes.SetChainId(initClientCtx.ChainID)
			}

			if f := cmd.Flags().Lookup(flags.FlagGasPrices); f != nil {
				gasPricesStr, _ := cmd.Flags().GetString(flags.FlagGasPrices)
				gasPrices, err := sdk.ParseCoinsNormalized(gasPricesStr)
				if err != nil {
					return err
				}
				if err := f.Value.Set(gasPrices.String()); err != nil {
					panic(err)
				}
			}

			customTemplate, customConfig := initAppConfig(fmt.Sprintf("2000000000000%s", pxtypes.StakingBondDenom()))

			return server.InterceptConfigsPreRunHandler(cmd, customTemplate, customConfig)
		},
	}
	initRootCmd(rootCmd, encodingConfig)
	overwriteFlagDefaults(rootCmd, map[string]string{
		flags.FlagChainID:        pxtypes.ChainId(),
		flags.FlagKeyringBackend: keyring.BackendOS,
		flags.FlagGas:            "auto",
		flags.FlagGasAdjustment:  "1.5",
		flags.FlagGasPrices:      "0.000002PUNDIX",
	})
	return rootCmd
}

func initAppConfig(minGasPrice string) (string, interface{}) {
	srvCfg := config.DefaultConfig()
	srvCfg.MinGasPrices = minGasPrice
	customAppConfig := appparams.Config{
		Config:       *srvCfg,
		BypassMinFee: appparams.DefaultBypassMinFee(),
	}
	return appparams.DefaultConfigTemplate(), customAppConfig
}

func initRootCmd(rootCmd *cobra.Command, encodingConfig appparams.EncodingConfig) {
	sdkCfgCmd := sdkCfg.Cmd()
	sdkCfgCmd.AddCommand(cli.UpdateCfgCmd(), cli.AppTomlCmd(), cli.ConfigTomlCmd())

	rootCmd.AddCommand(
		InitCmd(app.DefaultNodeHome, app.NewDefAppGenesisByDenom, app.CustomConsensusParams()),
		cli.CollectGenTxsCmd(banktypes.GenesisBalancesIterator{}, app.DefaultNodeHome),
		cli.GenTxCmd(app.ModuleBasics, encodingConfig.TxConfig, banktypes.GenesisBalancesIterator{}, app.DefaultNodeHome),
		genutilcli.ValidateGenesisCmd(app.ModuleBasics),
		cli.AddGenesisAccountCmd(app.DefaultNodeHome),
		tmcli.NewCompletionCmd(rootCmd, true),
		cli.Debug(),
		sdkCfgCmd,
		cli.DataCmd(),
		cli.PreUpgradeCmd(),
	)

	appCreator := appCreator{encodingConfig}
	server.AddCommands(rootCmd, app.DefaultNodeHome, appCreator.newApp, appCreator.appExport, addModuleInitFlags)

	// add keybase, auxiliary RPC, query, and tx child commands
	rootCmd.AddCommand(
		keys.Commands(app.DefaultNodeHome),
		cli.StatusCommand(),
		queryCommand(),
		txCommand(),
	)
}

func addModuleInitFlags(startCmd *cobra.Command) {
	crisis.AddModuleInitFlags(startCmd)
	startCmd.Flags().StringSlice(cli.FlagLogFilter, nil, `The logging filter can discard custom log type (ABCIQuery)`)

	startCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		serverCtx := server.GetServerContextFromCmd(cmd)

		if zeroLog, ok := serverCtx.Logger.(server.ZeroLogWrapper); ok {
			filterLogTypes, _ := cmd.Flags().GetStringSlice(cli.FlagLogFilter)
			if len(filterLogTypes) > 0 {
				serverCtx.Logger = cli.NewFxZeroLogWrapper(zeroLog, filterLogTypes)
			}
		}

		// Bind flags to the Context's Viper so the app construction can set
		// options accordingly.
		if err := serverCtx.Viper.BindPFlags(cmd.Flags()); err != nil {
			return err
		}

		if _, err := server.GetPruningOptionsFromFlags(serverCtx.Viper); err != nil {
			return err
		}

		genesisDoc, err := tmtypes.GenesisDocFromFile(serverCtx.Config.GenesisFile())
		if err != nil {
			return err
		}
		if err = checkMainnetAndBlock(genesisDoc, serverCtx.Config); err != nil {
			return err
		}
		pxtypes.SetChainId(genesisDoc.ChainID)
		return nil
	}
}

func queryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "query",
		Aliases:                    []string{"q"},
		Short:                      "Querying subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		authcmd.GetAccountCmd(),
		cli.ValidatorCommand(),
		cli.BlockCommand(),
		cli.QueryTxsByEventsCmd(),
		cli.QueryTxCmd(),
		cli.QueryStoreCmd(),
		cli.QueryValidatorByConsAddr(),
		cli.QueryBlockResultsCmd(),
		gasprice.QueryCmd(),
	)

	app.ModuleBasics.AddQueryCommands(cmd)
	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

func txCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "tx",
		Short:                      "Transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		authcmd.GetSignCommand(),
		authcmd.GetSignBatchCommand(),
		authcmd.GetMultiSignCommand(),
		authcmd.GetValidateSignaturesCommand(),
		authcmd.GetBroadcastCommand(),
		authcmd.GetEncodeCommand(),
		authcmd.GetDecodeCommand(),
		vestingcli.GetTxCmd(),
	)

	app.ModuleBasics.AddTxCommands(cmd)

	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

type appCreator struct {
	encCfg appparams.EncodingConfig
}

// newApp is an AppCreator
func (a appCreator) newApp(logger log.Logger, db dbm.DB, traceStore io.Writer, appOpts servertypes.AppOptions) servertypes.Application {
	var cache sdk.MultiStorePersistentCache

	if cast.ToBool(appOpts.Get(server.FlagInterBlockCache)) {
		cache = store.NewCommitKVStoreCacheManager()
	}

	skipUpgradeHeights := make(map[int64]bool)
	for _, h := range cast.ToIntSlice(appOpts.Get(server.FlagUnsafeSkipUpgrades)) {
		skipUpgradeHeights[int64(h)] = true
	}

	pruningOpts, err := server.GetPruningOptionsFromFlags(appOpts)
	if err != nil {
		panic(err)
	}

	snapshotDir := filepath.Join(cast.ToString(appOpts.Get(flags.FlagHome)), "data", "snapshots")
	snapshotDB, err := sdk.NewLevelDB("metadata", snapshotDir)
	if err != nil {
		panic(err)
	}
	snapshotStore, err := snapshots.NewStore(snapshotDB, snapshotDir)
	if err != nil {
		panic(err)
	}

	gasPricesStr := cast.ToString(appOpts.Get(server.FlagMinGasPrices))
	gasPrices, err := sdk.ParseCoinsNormalized(gasPricesStr)
	if err != nil {
		panic(err)
	}

	return app.NewPundixApp(
		logger, db, traceStore, true, skipUpgradeHeights,
		cast.ToString(appOpts.Get(flags.FlagHome)),
		cast.ToUint(appOpts.Get(server.FlagInvCheckPeriod)),
		a.encCfg,
		// this line is used by starport scaffolding # stargate/root/appArgument
		appOpts,
		baseapp.SetPruning(pruningOpts),
		baseapp.SetMinGasPrices(gasPrices.String()),
		baseapp.SetMinRetainBlocks(cast.ToUint64(appOpts.Get(server.FlagMinRetainBlocks))),
		baseapp.SetHaltHeight(cast.ToUint64(appOpts.Get(server.FlagHaltHeight))),
		baseapp.SetHaltTime(cast.ToUint64(appOpts.Get(server.FlagHaltTime))),
		baseapp.SetInterBlockCache(cache),
		baseapp.SetTrace(cast.ToBool(appOpts.Get(server.FlagTrace))),
		baseapp.SetIndexEvents(cast.ToStringSlice(appOpts.Get(server.FlagIndexEvents))),
		baseapp.SetSnapshotStore(snapshotStore),
		baseapp.SetSnapshotInterval(cast.ToUint64(appOpts.Get(server.FlagStateSyncSnapshotInterval))),
		baseapp.SetSnapshotKeepRecent(cast.ToUint32(appOpts.Get(server.FlagStateSyncSnapshotKeepRecent))),
	)
}

// appExport creates a new simapp (optionally at a given height)
func (a appCreator) appExport(
	logger log.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool, jailAllowedAddrs []string,
	appOpts servertypes.AppOptions,
) (servertypes.ExportedApp, error) {
	var anApp *app.PundixApp

	homePath, ok := appOpts.Get(flags.FlagHome).(string)
	if !ok || homePath == "" {
		return servertypes.ExportedApp{}, errors.New("application home not set")
	}

	var loadLatest bool
	if height == -1 {
		loadLatest = true
	}

	anApp = app.NewPundixApp(
		logger,
		db,
		traceStore,
		loadLatest,
		map[int64]bool{},
		homePath,
		cast.ToUint(appOpts.Get(server.FlagInvCheckPeriod)),
		a.encCfg,
		appOpts,
	)

	if height == -1 {
		if err := anApp.LoadHeight(height); err != nil {
			return servertypes.ExportedApp{}, err
		}
	}

	return anApp.ExportAppStateAndValidators(forZeroHeight, jailAllowedAddrs)
}

func overwriteFlagDefaults(c *cobra.Command, defaults map[string]string) {
	set := func(s *pflag.FlagSet, key, val string) {
		if f := s.Lookup(key); f != nil {
			f.DefValue = val
			if err := f.Value.Set(val); err != nil {
				panic(err)
			}
			if key == flags.FlagGasPrices {
				f.Usage = "Gas prices in decimal format to determine the transaction fee"
			}
			if key == flags.FlagGas {
				f.Usage = "gas limit to set per-transaction; set to 'auto' to calculate sufficient gas automatically"
			}
		}
	}
	for key, val := range defaults {
		set(c.Flags(), key, val)
		set(c.PersistentFlags(), key, val)
	}
	for _, c := range c.Commands() {
		overwriteFlagDefaults(c, defaults)
	}
}

func checkMainnetAndBlock(genesisDoc *tmtypes.GenesisDoc, config *tmcfg.Config) error {
	if genesisDoc.InitialHeight > 1 || genesisDoc.ChainID != pxtypes.MainnetChainId || config.StateSync.Enable {
		return nil
	}
	genesisTime, err := time.Parse("2006-01-02T15:04:05Z", "2021-10-13T10:00:00Z")
	if err != nil {
		return err
	}
	blockStoreDB, err := node.DefaultDBProvider(&node.DBContext{ID: "blockstore", Config: config})
	if err != nil {
		return err
	}
	defer blockStoreDB.Close()
	blockStore := tmstore.NewBlockStore(blockStoreDB)
	if genesisDoc.GenesisTime.Equal(genesisTime) && blockStore.Height() <= 0 {
		return fmt.Errorf("invalid version: Sync block from scratch please use use v0.1.x current %s", version.Version)
	}
	return nil
}
