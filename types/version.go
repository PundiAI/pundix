package types

import (
	"os"
	"sync"
)

// testnet constant
const (
	TestnetChainId          = "payalebar"
	testnetMintDenom        = "bsc0x0BEdB58eC8D603E71556ef8aA4014c68DBd57AD7"
	testnetStakingBondDenom = "ibc/169A52CA4862329131348484982CE75B3D6CC78AFB94C3107026C70CB66E7B2E"
)

// mainnet constant
const (
	MainnetChainId          = "PUNDIX"
	mainnetMintDenom        = "bsc0x29a63F4B209C29B4DC47f06FFA896F32667DAD2C"
	mainnetStakingBondDenom = "ibc/55367B7B6572631B78A93C66EF9FDFCE87CDE372CC4ED7848DA78C1EB1DCDD78"
)

var (
	chainId = MainnetChainId
	once    sync.Once
)

func SetChainId(id string) {
	if id != MainnetChainId && id != TestnetChainId {
		panic("invalid chainId: " + id)
	}
	once.Do(func() {
		chainId = id
	})
}

func ChainId() string {
	return chainId
}

func MintDenom() string {
	if denom := os.Getenv("LOCAL_MINT_DENOM"); len(denom) > 0 {
		return denom
	}
	if TestnetChainId == chainId {
		return testnetMintDenom
	}
	return mainnetMintDenom
}

func StakingBondDenom() string {
	if denom := os.Getenv("LOCAL_STAKING_BOND_DENOM"); len(denom) > 0 {
		return denom
	}
	if TestnetChainId == chainId {
		return testnetStakingBondDenom
	}
	return mainnetStakingBondDenom
}
