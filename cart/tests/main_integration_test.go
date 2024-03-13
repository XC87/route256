package tests

import (
	"github.com/stretchr/testify/suite"
	"route256.ozon.ru/project/cart/tests/integration/service"
	"testing"
)

func TestSmokeSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(service.Suit))
}
