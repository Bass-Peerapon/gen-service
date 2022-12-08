package http

import "github.com/Bass-Peerapon/gen-service/service/franchisee"

type franchiseeHandler struct {
	franchiseeUs franchisee.FranchiseeUsecase
}

func NewfranchiseeHandler(franchiseeUs franchisee.FranchiseeUsecase) franchisee.FranchiseeRepository {
	return &franchiseeHandler{
		franchiseeUs: franchiseeUs,
	}
}
