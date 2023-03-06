package app

import (
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	ibckeeper "github.com/cosmos/ibc-go/v4/modules/core/keeper"
	ibcstakinginterface "github.com/cosmos/interchain-security/legacy_ibc_testing/core"
	"github.com/cosmos/interchain-security/testutil/e2e"
	ibcproviderkeeper "github.com/cosmos/interchain-security/x/ccv/provider/keeper"
)

// ProviderApp interface implementations for e2e tests

// GetProviderKeeper implements the ProviderApp interface.
func (app *OraichainApp) GetProviderKeeper() ibcproviderkeeper.Keeper { //nolint:nolintlint
	return app.ProviderKeeper
}

// GetStakingKeeper implements the TestingApp interface. Needed for ICS.
func (app *OraichainApp) GetStakingKeeper() ibcstakinginterface.StakingKeeper { //nolint:nolintlint
	return app.stakingKeeper
}

// GetIBCKeeper implements the TestingApp interface.
func (app *OraichainApp) GetIBCKeeper() *ibckeeper.Keeper { //nolint:nolintlint
	return app.ibcKeeper
}

// GetScopedIBCKeeper implements the TestingApp interface.
func (app *OraichainApp) GetScopedIBCKeeper() capabilitykeeper.ScopedKeeper { //nolint:nolintlint
	return app.scopedIBCKeeper
}

// GetE2eStakingKeeper implements the ProviderApp interface.
func (app *OraichainApp) GetE2eStakingKeeper() e2e.E2eStakingKeeper { //nolint:nolintlint
	return app.stakingKeeper
}

// GetE2eBankKeeper implements the ProviderApp interface.
func (app *OraichainApp) GetE2eBankKeeper() e2e.E2eBankKeeper { //nolint:nolintlint
	return app.bankKeeper
}

// GetE2eSlashingKeeper implements the ProviderApp interface.
func (app *OraichainApp) GetE2eSlashingKeeper() e2e.E2eSlashingKeeper { //nolint:nolintlint
	return app.slashingKeeper
}

// GetE2eDistributionKeeper implements the ProviderApp interface.
func (app *OraichainApp) GetE2eDistributionKeeper() e2e.E2eDistributionKeeper { //nolint:nolintlint
	return app.distrKeeper
}
