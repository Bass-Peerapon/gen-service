package usecase

import "github.com/Bass-Peerapon/gen-service/service/franchisee"

type franchiseeUsecase struct {
	franchiseeRepo franchisee.FranchiseeRepository
}

func NewfranchiseeUsecase(franchiseeRepo franchisee.FranchiseeRepository) franchisee.FranchiseeUsecase {
	return &franchiseeUsecase{
		franchiseeRepo: franchiseeRepo,
	}
}
