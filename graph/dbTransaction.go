package graph

import (
	"errors"
	"os"
	"sync"
	"time"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/matoous/go-nanoid/v2"
	"github.com/rs/zerolog"
)

var log = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).With().Timestamp().Logger()

func (r *Resolver) TranWithNewUser(ghTokens []string) *Resolver {
	var wg sync.WaitGroup
	nanoid := func() string {
		nanoid, err := gonanoid.New()
		if err != nil {
			log.Error().Err(err).Msg("Generate userid failed")
			r.errors <- err
			return ""
		}
		return nanoid
	}
	reqCtx := &tablestore.BatchWriteRowRequest{}
	reqCtx.SetTraceID(*r.TransID)
	reqCtx.IsAtomic = true
	change := new(tablestore.PutRowChange)
	genUserID := func() (result *tablestore.PrimaryKey) {
		result = new(tablestore.PrimaryKey)
		result.AddPrimaryKeyColumn("userid", nanoid())
		return result
	}
	insertData := func(userToken *string) {
		defer wg.Done()
		change = &tablestore.PutRowChange{
			TableName:     "user",
			PrimaryKey:    genUserID(),
			Columns:       nil,
			Condition:     nil,
			ReturnType:    0,
			TransactionId: r.TransID,
		}
		change.AddColumnWithTimestamp("githubtoken", *userToken, time.Now().Unix())
		change.SetCondition(tablestore.RowExistenceExpectation_IGNORE)
		reqCtx.AddRowChange(change)
	}
	for _, v := range ghTokens {

		wg.Add(1)
		go insertData(&v)
	}
	wg.Wait()
	resp, err := r.OtsClient.BatchWriteRow(reqCtx)
	if err != nil {
		log.Error().Err(err).Msgf("submit data failed by RequestID: %s", resp.RequestId)
		r.errors <- err
		return r
	}
	return r
}
func (r *Resolver) Commit() error {
	// process error
	if len(r.errors) >= 0 {
		for err := range r.errors {
			log.Error().Err(err).Stack().Timestamp().Msg("")
		}
		if AbortTrans, err := r.OtsClient.AbortTransaction(&tablestore.AbortTransactionRequest{TransactionId: r.TransID}); err != nil {
			log.Fatal().Errs("commit transaction", []error{errors.New("abort transaction failed"), err}).Timestamp().Stack().Msg("the program will be exit.")
			// return err
		} else {
			log.Error().Errs("commit transaction", []error{errors.New("commit transaction failed"), err}).Timestamp().Stack().Msgf("now abort transaction, request ID is: %s", AbortTrans.RequestId)
			return err
		}
	}
	// commit trans
	if commitTrans, err := r.OtsClient.CommitTransaction(&tablestore.CommitTransactionRequest{TransactionId: r.TransID}); err != nil {
		// commit failed, abort and drop trans
		r.errors <- err
		return r.Commit()
	} else {
		log.Info().Timestamp().Msgf("commit transaction success, request ID is: %s", commitTrans.RequestId)
	}
	return nil
}
