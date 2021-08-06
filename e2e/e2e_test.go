package e2e_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ChainSafe/chainbridge-core-example/e2e"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ChainSafe/chainbridge-core/keystore"
	"github.com/ChainSafe/chainbridge-core/relayer"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
)

const TestEndpoint = "ws://localhost:8546"
const TestEndpoint2 = "ws://localhost:8547"

type IntegrationTestSuite struct {
	suite.Suite
	client            e2e.ChainClient
	client2           e2e.ChainClient
	bridgeAddr        common.Address
	erc20HandlerAddr  common.Address
	erc20ContractAddr common.Address
}

func TestRunE2ETests(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) SetupSuite()    {}
func (s *IntegrationTestSuite) TearDownSuite() {}

func (s *IntegrationTestSuite) SetupTest() {
	ethClient, err := evmclient.NewEVMClientFromParams(TestEndpoint, e2e.AliceKp.PrivateKey(), big.NewInt(e2e.DefaultMaxGasPrice))
	if err != nil {
		panic(err)
	}
	ethClient2, err := evmclient.NewEVMClientFromParams(TestEndpoint2, e2e.AliceKp.PrivateKey(), big.NewInt(e2e.DefaultMaxGasPrice))
	if err != nil {
		panic(err)
	}
	b, err := ethClient.LatestBlock()
	if err != nil {
		panic(err)
	}
	log.Debug().Msgf("Latest block %s", b.String())
	bridgeAddr, erc20Addr, erc20HandlerAddr, err := e2e.DeployAndPrepareEnv(ethClient, 1, big.NewInt(1))
	if err != nil {
		panic(err)
	}
	s.client = ethClient
	s.client2 = ethClient2
	s.bridgeAddr = bridgeAddr
	s.erc20ContractAddr = erc20Addr
	s.erc20HandlerAddr = erc20HandlerAddr
	//Contract addresses should be the same
	_, _, _, err = e2e.DeployAndPrepareEnv(ethClient2, 2, big.NewInt(1))
	if err != nil {
		panic(err)
	}
}

func (s *IntegrationTestSuite) TestDeposit() {
	dstAddr := keystore.TestKeyRing.EthereumKeys[keystore.BobKey].CommonAddress()
	senderBalBefore, err := e2e.GetBalance(s.client, s.erc20ContractAddr, e2e.AliceKp.CommonAddress())
	s.Nil(err)
	destBalanceBefore, err := e2e.GetBalance(s.client2, s.erc20ContractAddr, dstAddr)
	s.Nil(err)

	b, err := s.client2.LatestBlock()
	if err != nil {
		panic(err)
	}
	amountToDeposit := big.NewInt(1000000)
	resourceID := e2e.SliceTo32Bytes(append(common.LeftPadBytes(s.erc20ContractAddr.Bytes(), 31), uint8(0)))
	err = e2e.Deposit(s.client, s.bridgeAddr, dstAddr, amountToDeposit, resourceID, 2)
	s.Nil(err)

	//Wait 30 seconds for relayer vote
	time.Sleep(60 * time.Second)
	senderBalAfter, err := e2e.GetBalance(s.client, s.erc20ContractAddr, e2e.AliceKp.CommonAddress())
	s.Nil(err)
	s.Equal(-1, senderBalAfter.Cmp(senderBalBefore))

	ba, err := s.client2.LatestBlock()
	if err != nil {
		panic(err)
	}
	//wait for vote log event
	proposalEvent := "ProposalEvent(uint8,uint64,uint8,bytes32,bytes32)"
	evts, err := s.client2.FetchEventLogs(context.Background(), s.bridgeAddr, proposalEvent, b, ba)
	var passedEventFound bool
	for _, evt := range evts {
		status := evt.Topics[3].Big().Uint64()
		if uint8(relayer.ProposalStatusPassed) == uint8(status) {
			passedEventFound = true
		}
	}
	s.True(passedEventFound)
	s.Equal(senderBalBefore.Cmp(big.NewInt(0).Add(senderBalAfter, amountToDeposit)), 0)

	//Wait 30 seconds for relayer to execute
	time.Sleep(30 * time.Second)

	ba, err = s.client2.LatestBlock()
	s.Nil(err)
	queryExecute, err := s.client2.FetchEventLogs(context.Background(), s.bridgeAddr, proposalEvent, b, ba)
	s.Nil(err)
	var executedEventFound bool
	for _, evt := range queryExecute {
		status := evt.Topics[3].Big().Uint64()
		if uint8(relayer.ProposalStatusExecuted) == uint8(status) {
			executedEventFound = true
		}
	}
	s.True(executedEventFound)
	//
	destBalanceAfter, err := e2e.GetBalance(s.client2, s.erc20ContractAddr, dstAddr)
	s.Nil(err)
	//Balance has increased
	s.Equal(1, destBalanceAfter.Cmp(destBalanceBefore))
}
