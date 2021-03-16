#!/usr/bin/env bash

# shellcheck disable=SC2046

set -e

export LOCAL_MINT_DENOM="PURSE"
export LOCAL_STAKING_BOND_DENOM="PUNDIX"

if [[ "$1" == "init" ]]; then
  if [ -d ~/.pundix ]; then
    read -p "Are you sure you want to delete all the data and start over? [y/N] " input
    if [[ "$input" != "y" && "$input" != "Y" ]]; then
      exit 1
    fi
    rm -rf ~/.pundix
  fi
  # Initialize private validator, p2p, genesis, and application configuration files
  pundixd init --chain-id="PUNDIX" --denom="PUNDIX" --mint-denom="PURSE" local

  # update pundix client config
  pundixd config config.toml instrumentation.prometheus true
  pundixd config config.toml rpc.laddr tcp://0.0.0.0:26657
  # consensus
  pundixd config config.toml consensus.timeout_commit 1s

  # update pundix client config
  pundixd config app.toml grpc-web.enable false
  pundixd config app.toml telemetry.enabled true
  pundixd config app.toml telemetry.prometheus-retention-time 60
  pundixd config app.toml api.enable true
  pundixd config app.toml minimum-gas-prices "2000000000000PUNDIX"

  # update pundix client config
  pundixd config chain-id PUNDIX
  pundixd config keyring-backend test
  pundixd config output json
  pundixd config broadcast-mode "block"

  pundixd keys add fx1
  pundixd add-genesis-account fx1 1000000000000000000000PUNDIX
  pundixd gentx fx1 100000000000000000000PUNDIX --chain-id PUNDIX \
    --moniker="pundix-val-1" \
    --commission-max-change-rate=0.01 \
    --commission-max-rate=0.2 \
    --commission-rate=0.03 \
    --details="Details A PUNDIX self-hosted validator." \
    --security-contact="contact@pundix.com" \
    --website="https://pundix.com"
  pundixd collect-gentxs
fi

pundixd start --log_filter='ABCIQuery'
