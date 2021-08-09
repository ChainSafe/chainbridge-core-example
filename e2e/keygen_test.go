// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package e2e

import (
	"testing"

	"github.com/ChainSafe/chainbridge-core/crypto/secp256k1"
	"github.com/rs/zerolog/log"

	"github.com/stretchr/testify/suite"
)

type KGTestSuite struct {
	suite.Suite
}

func TestRunKGTests(t *testing.T) {
	suite.Run(t, new(KGTestSuite))
}

func (s *KGTestSuite) SetupSuite()    {}
func (s *KGTestSuite) TearDownSuite() {}

func (s *KGTestSuite) SetupTest() {}

func (s *KGTestSuite) TestGenKey() {
	kp, err := secp256k1.GenerateKeypair()
	s.Nil(err)
	log.Debug().Msgf("Addr: %s,  PrivKey %x", kp.CommonAddress().String(), kp.Encode())
}
