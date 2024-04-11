package tests

import (
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"
	"route256.ozon.ru/project/cart/tests/integration/service"
	"testing"
)

func TestSmokeSuite(t *testing.T) {
	defer goleak.VerifyNone(t)
	t.Parallel()
	suite.Run(t, new(service.Suit))
}
