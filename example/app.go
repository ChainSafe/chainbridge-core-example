package example

import (
	"os"
	"os/signal"
	"syscall"

	celoVoter "github.com/ChainSafe/chainbridge-celo-module/voter"
	"github.com/ChainSafe/chainbridge-core-example/example/keystore"
	"github.com/ChainSafe/chainbridge-core/chains/evm"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ChainSafe/chainbridge-core/chains/evm/listener"
	"github.com/ChainSafe/chainbridge-core/chains/evm/voter"

	"github.com/ChainSafe/chainbridge-core/chains/substrate"
	subListener "github.com/ChainSafe/chainbridge-core/chains/substrate/listener"
	subWriter "github.com/ChainSafe/chainbridge-core/chains/substrate/writer"
	"github.com/ChainSafe/chainbridge-core/crypto/sr25519"
	"github.com/ChainSafe/chainbridge-core/lvldb"
	"github.com/ChainSafe/chainbridge-core/relayer"
	subModule "github.com/ChainSafe/chainbridge-substrate-module"
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
const TestEndpointCelo = "ws://localhost:8546"

//Bridge:             0x62877dDCd49aD22f5eDfc6ac108e9a4b5D2bD88B
//Erc20 Handler:      0x3167776db165D8eA0f51790CA2bbf44Db5105ADF
func Run() error {
	errChn := make(chan error)
	stopChn := make(chan struct{})

	db, err := lvldb.NewLvlDB("./lvldbdata")
	if err != nil {
		panic(err)
	}

	ethClient, err := evmclient.NewEVMClient(TestEndpoint, AliceKp)
	if err != nil {
		panic(err)
	}

	ethClient.Configurate()

	eventHandler := listener.NewETHEventHandler(common.HexToAddress("0x62877dDCd49aD22f5eDfc6ac108e9a4b5D2bD88B"), ethClient)
	eventHandler.RegisterEventHandler("0x3167776db165D8eA0f51790CA2bbf44Db5105ADF", listener.Erc20EventHandler)
	evmListener := listener.NewEVMListener(ethClient, eventHandler, common.HexToAddress("0x62877dDCd49aD22f5eDfc6ac108e9a4b5D2bD88B"))
	mh := voter.NewEVMMessageHandler(ethClient, common.HexToAddress("0x62877dDCd49aD22f5eDfc6ac108e9a4b5D2bD88B"))
	mh.RegisterMessageHandler(common.HexToAddress("0x3167776db165D8eA0f51790CA2bbf44Db5105ADF"), voter.ERC20MessageHandler)
	evmVoter := voter.NewVoter(mh, ethClient)
	evmChain := evm.NewEVMChain(evmListener, evmVoter, db, 0)

	// Substrate setup
	kp, err := keystore.KeypairFromAddress("5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY", keystore.SubChain, "alice", true)
	if err != nil {
		panic(err)
	}
	krp := kp.(*sr25519.Keypair).AsKeyringPair()
	subC, err := subModule.NewSubstrateClient("ws://localhost:9944", krp, stopChn)
	if err != nil {
		panic(err)
	}

	subL := subListener.NewSubstrateListener(subC)
	subW := subWriter.NewSubstrateWriter(1, subC)

	// TODO: really not need this dynamic handler assignment
	subL.RegisterSubscription(relayer.FungibleTransfer, subListener.FungibleTransferHandler)
	subL.RegisterSubscription(relayer.GenericTransfer, subListener.GenericTransferHandler)
	subL.RegisterSubscription(relayer.NonFungibleTransfer, subListener.NonFungibleTransferHandler)

	subW.RegisterHandler(relayer.FungibleTransfer, subWriter.CreateFungibleProposal)
	subChain := substrate.NewSubstrateChain(subL, subW, db, 1)

	// Celo setup
	ethClientCelo, err := evmclient.NewEVMClient(TestEndpointCelo, AliceKp)
	if err != nil {
		panic(err)
	}
	ethClientCelo.Configurate()
	eventHandlerCelo := listener.NewETHEventHandler(common.HexToAddress("0x62877dDCd49aD22f5eDfc6ac108e9a4b5D2bD88B"), ethClientCelo)
	eventHandlerCelo.RegisterEventHandler("0x3167776db165D8eA0f51790CA2bbf44Db5105ADF", listener.Erc20EventHandler)
	evmListenerCelo := listener.NewEVMListener(ethClientCelo, eventHandlerCelo, common.HexToAddress("0x62877dDCd49aD22f5eDfc6ac108e9a4b5D2bD88B"))
	mhCelo := voter.NewEVMMessageHandler(ethClientCelo, common.HexToAddress("0x62877dDCd49aD22f5eDfc6ac108e9a4b5D2bD88B"))
	mhCelo.RegisterMessageHandler(common.HexToAddress("0x3167776db165D8eA0f51790CA2bbf44Db5105ADF"), celoVoter.ERC20CeloMessageHandler)
	evmVoterCelo := voter.NewVoter(mhCelo, ethClientCelo)
	celoChain := evm.NewEVMChain(evmListenerCelo, evmVoterCelo, db, 2)

	r := relayer.NewRelayer([]relayer.RelayedChain{subChain, evmChain, celoChain})

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
