package example

import (
	"os"
	"os/signal"
	"syscall"

	evmClient "github.com/ChainSafe/chainbridge-eth-module"
	subClient "github.com/ChainSafe/chainbridge-substrate-module"

	"github.com/ChainSafe/chainbridge-core/chains"
	"github.com/ChainSafe/chainbridge-core/config"

	"github.com/ChainSafe/chainbridge-core/keystore"

	"github.com/ChainSafe/chainbridge-core/lvldb"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
)

var AliceKp = keystore.TestKeyRing.EthereumKeys[keystore.AliceKey]
var BobKp = keystore.TestKeyRing.EthereumKeys[keystore.BobKey]
var EveKp = keystore.TestKeyRing.EthereumKeys[keystore.EveKey]

var AliceKpSub = keystore.TestKeyRing.SubstrateKeys[keystore.AliceKey]
var BobKpSub = keystore.TestKeyRing.SubstrateKeys[keystore.BobKey]
var EveKpSub = keystore.TestKeyRing.SubstrateKeys[keystore.EveKey]

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

	db, err := lvldb.NewLvlDB("./lvldbdata")
	if err != nil {
		panic(err)
	}

	cfg, err := chains.GetConfig(".", "fullConfig")
	if err != nil {
		panic(err)
	}
	log.Info().Msgf("%v", cfg.Chains)

	ethClient := evmClient.NewEVMClient()
	substrateClient := subClient.NewSubstrateClient()
	r, err := config.InitializeRelayer(cfg, ethClient, substrateClient, AliceKp, db, ethClient.ReturnErc20HandlerFabric, stopChn)
	if err != nil {
		panic(err)
	}

	// relayedChains := make([]relayer.RelayedChain, len(cfg.Chains))
	// for index, chainConfig := range cfg.Chains {

	// 	if chainConfig.Type == "ethereum" {
	// 		evmChain, err := evmClient.InitializeEthChain(&chainConfig, db, AliceKp)
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		relayedChains[index] = evmChain
	// 	} else if chainConfig.Type == "substrate" {
	// 		subChain, err := subClient.InitializeSubChain(&chainConfig, db, stopChn)
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		relayedChains[index] = subChain
	// 	} else {
	// 		return errors.New("unrecognized Chain Type")
	// 	}

	// }

	// r := relayer.NewRelayer(relayedChains)

	// config.InitializeRelayer(
	// 	cfg,
	// 	ethClient,
	// 	substrateClient,
	// 	AliceKp,
	// 	kvdb,
	// 	stopChn)

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
