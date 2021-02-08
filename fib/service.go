package fib

import (
	"math/big"
)

type Service interface {
	Fib(uint64) *big.Int
}

type service struct {
	cache [2]*big.Int
}

func (svc service) Fib(n uint64) *big.Int {
	if n == 0 {
		return big.NewInt(0)
	} else if n == 1 {
		return big.NewInt(1)
	}

	svc.cache = [2]*big.Int{big.NewInt(1), big.NewInt(1)}
	for i := uint64(2); i < n; i++ {
		svc.cache[i%2] = svc.cache[0].Add(svc.cache[0], svc.cache[1])
	}

	return svc.cache[n%2]
}

func NewService() Service {
	return &service{}
}
