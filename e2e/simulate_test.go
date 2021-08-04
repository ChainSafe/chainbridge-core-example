package e2e

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
)

type SimTestSuite struct {
	suite.Suite
	client ChainClient
}

const TestEndpoint = "ws://localhost:8545"

func TestRunSimTests(t *testing.T) {
	suite.Run(t, new(SimTestSuite))
}
func (s *SimTestSuite) SetupSuite()    {}
func (s *SimTestSuite) TearDownSuite() {}

//bridge: 0x3f709398808af36ADBA86ACC617FeB7F5B7B193E
//erc20: 0x3167776db165D8eA0f51790CA2bbf44Db5105ADF
//erc20handler: 0x2B6Ab4b880A45a07d83Cf4d664Df4Ab85705Bc07
//from: 0xff93B45308FD417dF303D6515aB04D9e89a750Ca
func (s *SimTestSuite) SetupTest() {
	ethClient, err := evmclient.NewEVMClientFromParams(TestEndpoint, AliceKp.PrivateKey())
	if err != nil {
		panic(err)
	}
	s.client = ethClient
}

func (s *SimTestSuite) TestSimulate() {
	b, err := Simulate(s.client, big.NewInt(1), common.HexToHash("0x123"), AliceKp.CommonAddress())
	if err != nil {
		panic(err)
	}
	log.Debug().Msgf("LEN SIM %v", len(b))
}
