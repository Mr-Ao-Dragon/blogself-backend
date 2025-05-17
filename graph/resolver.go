//go:generate go run github.com/99designs/gqlgen generate
package graph

import (
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	OtsClient *tablestore.TableStoreClient
	TransID   *string
	errors    chan error
}

func (r *Resolver) NewTrans(table string, transPK map[string]any) *Resolver {
	transPrimaryKey := new(tablestore.PrimaryKey)
	for k, v := range transPK {
		transPrimaryKey.AddPrimaryKeyColumn(k, v)
	}
	trans := &tablestore.StartLocalTransactionRequest{
		PrimaryKey: transPrimaryKey,
		TableName:  table,
	}
	resp, err := r.OtsClient.StartLocalTransaction(trans)
	if err != nil {
		log.Fatal().Err(err).Msgf("launch Transtion failed by RequestID: %s", resp.RequestId)
		r.errors <- err
	}
	r.TransID = resp.TransactionId
	return r
}
