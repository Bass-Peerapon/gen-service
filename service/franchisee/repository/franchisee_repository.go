package repository

import (
	"git.innovasive.co.th/backend/psql"
	"github.com/Bass-Peerapon/gen-service/service/franchisee"
)

type franchiseeRepository struct {
	client *psql.Client
}

func NewfranchiseeRepository(client *psql.Client) franchisee.FranchiseeRepository {
	return &franchiseeRepository{
		client: client,
	}
}
