package example

import (
	"os"
	"os/signal"
	"syscall"

	evmClient "github.com/ChainSafe/chainbridge-eth-module"
	subClient "github.com/ChainSafe/chainbridge-substrate-module"
	"github.com/spf13/viper"

	"github.com/ChainSafe/chainbridge-core/chains/evm"
	"github.com/ChainSafe/chainbridge-core/chains/evm/listener"
	"github.com/ChainSafe/chainbridge-core/chains/evm/writer"
	"github.com/ChainSafe/chainbridge-core/config"

	"github.com/ChainSafe/chainbridge-core/chains/substrate"
	subListener "github.com/ChainSafe/chainbridge-core/chains/substrate/listener"
	subWriter "github.com/ChainSafe/chainbridge-core/chains/substrate/writer"
	"github.com/ChainSafe/chainbridge-core/lvldb"
	"github.com/ChainSafe/chainbridge-core/relayer"
	"github.com/rs/zerolog/log"
)

//Bridge:             0x62877dDCd49aD22f5eDfc6ac108e9a4b5D2bD88B
//Erc20 Handler:      0x3167776db165D8eA0f51790CA2bbf44Db5105ADF
func Run() error {
	errChn := make(chan error)
	stopChn := make(chan struct{})

	db, err := lvldb.NewLvlDB(viper.GetString(config.BlockstoreFlagName))
	if err != nil {
		panic(err)
	}

	ethClient := evmClient.NewEVMClient()
	err = ethClient.Configurate(".", "config")
	if err != nil {
		panic(err)
	}
	ethSharedConfig := ethClient.GetConfig().SharedEVMConfig

	evmListener := listener.NewEVMListener(ethClient)
	evmListener.RegisterHandlerFabric(ethSharedConfig.Erc20Handler, ethClient.ReturnErc20HandlerFabric)

	evmWriter := writer.NewWriter(ethClient)
	evmWriter.RegisterProposalHandler(ethSharedConfig.Erc20Handler, writer.ERC20ProposalHandler)

	// TODO: get rid of bridge and id params
	evmChain := evm.NewEVMChain(
		evmListener,
		evmWriter,
		db,
		ethSharedConfig.Bridge,
		*ethSharedConfig.GeneralChainConfig.Id,
		&ethSharedConfig)
	if err != nil {
		panic(err)
	}

	subC := subClient.NewSubstrateClient(stopChn)
	err = subC.Configurate(".", "subConfig")
	if err != nil {
		panic(err)
	}
	subSharedConfig := subC.GetConfig().SharedSubstrateConfig

	subL := subListener.NewSubstrateListener(subC)
	subW := subWriter.NewSubstrateWriter(1, subC)

	// TODO: really not need this dynamic handler assignment
	subL.RegisterSubscription(relayer.FungibleTransfer, subListener.FungibleTransferHandler)
	subL.RegisterSubscription(relayer.GenericTransfer, subListener.GenericTransferHandler)
	subL.RegisterSubscription(relayer.NonFungibleTransfer, subListener.NonFungibleTransferHandler)

	subW.RegisterHandler(relayer.FungibleTransfer, subWriter.CreateFungibleProposal)

	// TODO: get rid of id param
	subChain := substrate.NewSubstrateChain(
		subL,
		subW,
		db,
		*subSharedConfig.GeneralChainConfig.Id,
		&subSharedConfig)

	r := relayer.NewRelayer([]relayer.RelayedChain{evmChain, subChain})

	go r.Start(stopChn, errChn)

	sysErr := make(chan os.Signal, 1)
	signal.Notify(sysErr,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGHUP,
		syscall.SIGQUIT)

	select {
	case err := <-errChn:
		log.Error().Err(err).Msg("failed to listen and serve")
		close(stopChn)
		return err
	case sig := <-sysErr:
		log.Info().Msgf("terminating got [%v] signal", sig)
		return nil
	}
}
