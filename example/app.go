package example

import (
	"os"
	"os/signal"
	"syscall"

	evmClient "github.com/ChainSafe/chainbridge-eth-module"
	subClient "github.com/ChainSafe/chainbridge-substrate-module"
	"github.com/spf13/viper"

	subClientConfig "github.com/ChainSafe/chainbridge-substrate-module/config"

	"github.com/ChainSafe/chainbridge-core-example/example/keystore"
	"github.com/ChainSafe/chainbridge-core/chains/evm"
	"github.com/ChainSafe/chainbridge-core/chains/evm/listener"
	"github.com/ChainSafe/chainbridge-core/chains/evm/writer"
	"github.com/ChainSafe/chainbridge-core/config"
	"github.com/ChainSafe/chainbridge-core/crypto/sr25519"

	"github.com/ChainSafe/chainbridge-core/chains/substrate"
	subListener "github.com/ChainSafe/chainbridge-core/chains/substrate/listener"
	subWriter "github.com/ChainSafe/chainbridge-core/chains/substrate/writer"
	"github.com/ChainSafe/chainbridge-core/lvldb"
	"github.com/ChainSafe/chainbridge-core/relayer"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
)

var AliceKp = keystore.TestKeyRing.EthereumKeys[keystore.AliceKey]
var BobKp = keystore.TestKeyRing.EthereumKeys[keystore.BobKey]
var EveKp = keystore.TestKeyRing.EthereumKeys[keystore.EveKey]

var (
	DefaultRelayerAddresses = []common.Address{
		common.HexToAddress(keystore.TestKeyRing.EthereumKeys[keystore.AliceKey].Address()),
		common.HexToAddress(keystore.TestKeyRing.EthereumKeys[keystore.BobKey].Address()),
		common.HexToAddress(keystore.TestKeyRing.EthereumKeys[keystore.CharlieKey].Address()),
		common.HexToAddress(keystore.TestKeyRing.EthereumKeys[keystore.DaveKey].Address()),
		common.HexToAddress(keystore.TestKeyRing.EthereumKeys[keystore.EveKey].Address()),
	}
)

const DefaultGasLimit = 6721975
const DefaultGasPrice = 20000000000

const TestEndpoint = "ws://localhost:8545"
const TestEndpoint2 = "ws://localhost:8546"

//Bridge:             0x62877dDCd49aD22f5eDfc6ac108e9a4b5D2bD88B
//Erc20 Handler:      0x3167776db165D8eA0f51790CA2bbf44Db5105ADF
func Run() error {
	errChn := make(chan error)
	stopChn := make(chan struct{})

	db, err := lvldb.NewLvlDB(viper.GetString(config.BlockstoreFlagName))
	if err != nil {
		panic(err)
	}

	ethClient := evmClient.NewEVMClient(AliceKp)
	err = ethClient.Configurate(".", "config")
	if err != nil {
		panic(err)
	}
	ethConfig := ethClient.GetConfig()

	evmListener := listener.NewEVMListener(ethClient)
	evmListener.RegisterHandlerFabric(ethClient.GetConfig().SharedEVMConfig.Erc20Handler, ethClient.ReturnErc20HandlerFabric)

	evmWriter := writer.NewWriter(ethClient)
	evmWriter.RegisterProposalHandler(ethClient.GetConfig().SharedEVMConfig.Erc20Handler, writer.ERC20ProposalHandler)

	// TODO: get rid of bridge and id params
	evmChain := evm.NewEVMChain(
		evmListener,
		evmWriter,
		db,
		ethConfig.SharedEVMConfig.Bridge,
		*ethConfig.SharedEVMConfig.GeneralChainConfig.Id,
		&ethConfig.SharedEVMConfig)
	if err != nil {
		panic(err)
	}

	// TODO: this can be removed once grabbing the keypair is done inside the module
	// This most likely can only be done by moving the keystore package into core to be used by the module
	subCfg, err := subClientConfig.GetConfig(".", "subConfig")
	if err != nil {
		panic(err)
	}

	kp, err := keystore.KeypairFromAddress(subCfg.GeneralChainConfig.From, keystore.SubChain, "alice", true)
	if err != nil {
		panic(err)
	}
	krp := kp.(*sr25519.Keypair).AsKeyringPair()

	subC := subClient.NewSubstrateClient(krp, stopChn)
	err = subC.Configurate(".", "subConfig")
	if err != nil {
		panic(err)
	}
	subConfig := subC.GetConfig()

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
		*subConfig.SharedSubstrateConfig.GeneralChainConfig.Id,
		&subConfig.SharedSubstrateConfig)

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
