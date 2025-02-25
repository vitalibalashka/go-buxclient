package transports

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/BuxOrg/bux"
	"github.com/mrz1836/go-datastore"
)

// NewXpub will register an xPub
func (h *TransportHTTP) NewXpub(ctx context.Context, rawXPub string, metadata *bux.Metadata) error {

	// Adding a xpub needs to be signed by an admin key
	if h.adminXPriv == nil {
		return ErrAdminKey
	}

	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldMetadata: processMetadata(metadata),
		FieldXpubKey:  rawXPub,
	})
	if err != nil {
		return err
	}

	var xPubData bux.Xpub

	return h.doHTTPRequest(
		ctx, http.MethodPost, "/xpub", jsonStr, h.adminXPriv, true, &xPubData,
	)
}

// RegisterXpub alias for NewXpub
func (h *TransportHTTP) RegisterXpub(ctx context.Context, rawXPub string, metadata *bux.Metadata) error {
	return h.NewXpub(ctx, rawXPub, metadata)
}

// AdminGetStatus get whether admin key is valid
func (h *TransportHTTP) AdminGetStatus(ctx context.Context) (bool, error) {

	var status bool
	if err := h.doHTTPRequest(
		ctx, http.MethodGet, "/admin/status", nil, h.xPriv, true, &status,
	); err != nil {
		return false, err
	}
	if h.debug {
		log.Printf("admin status: %v\n", status)
	}

	return status, nil
}

// AdminGetStats get admin stats
func (h *TransportHTTP) AdminGetStats(ctx context.Context) (*bux.AdminStats, error) {

	var stats *bux.AdminStats
	if err := h.doHTTPRequest(
		ctx, http.MethodGet, "/admin/stats", nil, h.xPriv, true, &stats,
	); err != nil {
		return nil, err
	}
	if h.debug {
		log.Printf("admin stats: %v\n", stats)
	}

	return stats, nil
}

// AdminGetAccessKeys get all access keys filtered by conditions
func (h *TransportHTTP) AdminGetAccessKeys(ctx context.Context, conditions map[string]interface{},
	metadata *bux.Metadata, queryParams *datastore.QueryParams) ([]*bux.AccessKey, error) {

	var models []*bux.AccessKey
	if err := h.adminGetModels(ctx, conditions, metadata, queryParams, "/admin/access-keys/search", &models); err != nil {
		return nil, err
	}

	return models, nil
}

// AdminGetAccessKeysCount get a count of all the access keys filtered by conditions
func (h *TransportHTTP) AdminGetAccessKeysCount(ctx context.Context, conditions map[string]interface{},
	metadata *bux.Metadata) (int64, error) {

	return h.adminCount(ctx, conditions, metadata, "/admin/access-keys/count")
}

// AdminGetBlockHeaders get all block headers filtered by conditions
func (h *TransportHTTP) AdminGetBlockHeaders(ctx context.Context, conditions map[string]interface{},
	metadata *bux.Metadata, queryParams *datastore.QueryParams) ([]*bux.BlockHeader, error) {

	var models []*bux.BlockHeader
	if err := h.adminGetModels(ctx, conditions, metadata, queryParams, "/admin/block-headers/search", &models); err != nil {
		return nil, err
	}

	return models, nil
}

// AdminGetBlockHeadersCount get a count of all the block headers filtered by conditions
func (h *TransportHTTP) AdminGetBlockHeadersCount(ctx context.Context, conditions map[string]interface{},
	metadata *bux.Metadata) (int64, error) {

	return h.adminCount(ctx, conditions, metadata, "/admin/block-headers/count")
}

// AdminGetDestinations get all block destinations filtered by conditions
func (h *TransportHTTP) AdminGetDestinations(ctx context.Context, conditions map[string]interface{},
	metadata *bux.Metadata, queryParams *datastore.QueryParams) ([]*bux.Destination, error) {

	var models []*bux.Destination
	if err := h.adminGetModels(ctx, conditions, metadata, queryParams, "/admin/destinations/search", &models); err != nil {
		return nil, err
	}

	return models, nil
}

// AdminGetDestinationsCount get a count of all the destinations filtered by conditions
func (h *TransportHTTP) AdminGetDestinationsCount(ctx context.Context, conditions map[string]interface{},
	metadata *bux.Metadata) (int64, error) {

	return h.adminCount(ctx, conditions, metadata, "/admin/destinations/count")
}

// AdminGetPaymail get a paymail by address
func (h *TransportHTTP) AdminGetPaymail(ctx context.Context, address string) (*bux.PaymailAddress, error) {

	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldAddress: address,
	})
	if err != nil {
		return nil, err
	}

	var model *bux.PaymailAddress
	if err = h.doHTTPRequest(
		ctx, http.MethodGet, "/admin/paymail/get", jsonStr, h.xPriv, true, &model,
	); err != nil {
		return nil, err
	}
	if h.debug {
		log.Printf("admin get paymail: %v\n", model)
	}

	return model, nil
}

// AdminGetPaymails get all block paymails filtered by conditions
func (h *TransportHTTP) AdminGetPaymails(ctx context.Context, conditions map[string]interface{},
	metadata *bux.Metadata, queryParams *datastore.QueryParams) ([]*bux.PaymailAddress, error) {

	var models []*bux.PaymailAddress
	if err := h.adminGetModels(ctx, conditions, metadata, queryParams, "/admin/paymails/search", &models); err != nil {
		return nil, err
	}

	return models, nil
}

// AdminGetPaymailsCount get a count of all the paymails filtered by conditions
func (h *TransportHTTP) AdminGetPaymailsCount(ctx context.Context, conditions map[string]interface{},
	metadata *bux.Metadata) (int64, error) {

	return h.adminCount(ctx, conditions, metadata, "/admin/paymails/count")
}

// AdminCreatePaymail create a new paymail for a xpub
func (h *TransportHTTP) AdminCreatePaymail(ctx context.Context, xPubID string, address string, publicName string, avatar string) (*bux.PaymailAddress, error) {

	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldXpubID:     xPubID,
		FieldAddress:    address,
		FieldPublicName: publicName,
		FieldAvatar:     avatar,
	})
	if err != nil {
		return nil, err
	}

	var model *bux.PaymailAddress
	if err = h.doHTTPRequest(
		ctx, http.MethodPost, "/admin/paymail/create", jsonStr, h.xPriv, true, &model,
	); err != nil {
		return nil, err
	}
	if h.debug {
		log.Printf("admin create paymail: %v\n", model)
	}

	return model, nil
}

// AdminDeletePaymail delete a paymail address from the database
func (h *TransportHTTP) AdminDeletePaymail(ctx context.Context, address string) (*bux.PaymailAddress, error) {

	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldAddress: address,
	})
	if err != nil {
		return nil, err
	}

	var model *bux.PaymailAddress
	if err = h.doHTTPRequest(
		ctx, http.MethodPost, "/admin/paymail/delete", jsonStr, h.xPriv, true, &model,
	); err != nil {
		return nil, err
	}
	if h.debug {
		log.Printf("admin delete paymail: %v\n", model)
	}

	return model, nil
}

// AdminGetTransactions get all block transactions filtered by conditions
func (h *TransportHTTP) AdminGetTransactions(ctx context.Context, conditions map[string]interface{},
	metadata *bux.Metadata, queryParams *datastore.QueryParams) ([]*bux.Transaction, error) {

	var models []*bux.Transaction
	if err := h.adminGetModels(ctx, conditions, metadata, queryParams, "/admin/transactions/search", &models); err != nil {
		return nil, err
	}

	return models, nil
}

// AdminGetTransactionsCount get a count of all the transactions filtered by conditions
func (h *TransportHTTP) AdminGetTransactionsCount(ctx context.Context, conditions map[string]interface{},
	metadata *bux.Metadata) (int64, error) {

	return h.adminCount(ctx, conditions, metadata, "/admin/transactions/count")
}

// AdminGetUtxos get all block utxos filtered by conditions
func (h *TransportHTTP) AdminGetUtxos(ctx context.Context, conditions map[string]interface{},
	metadata *bux.Metadata, queryParams *datastore.QueryParams) ([]*bux.Utxo, error) {

	var models []*bux.Utxo
	if err := h.adminGetModels(ctx, conditions, metadata, queryParams, "/admin/utxos/search", &models); err != nil {
		return nil, err
	}

	return models, nil
}

// AdminGetUtxosCount get a count of all the utxos filtered by conditions
func (h *TransportHTTP) AdminGetUtxosCount(ctx context.Context, conditions map[string]interface{},
	metadata *bux.Metadata) (int64, error) {

	return h.adminCount(ctx, conditions, metadata, "/admin/utxos/count")
}

// AdminGetXPubs get all block xpubs filtered by conditions
func (h *TransportHTTP) AdminGetXPubs(ctx context.Context, conditions map[string]interface{},
	metadata *bux.Metadata, queryParams *datastore.QueryParams) ([]*bux.Xpub, error) {

	var models []*bux.Xpub
	if err := h.adminGetModels(ctx, conditions, metadata, queryParams, "/admin/xpubs/search", &models); err != nil {
		return nil, err
	}

	return models, nil
}

// AdminGetXPubsCount get a count of all the xpubs filtered by conditions
func (h *TransportHTTP) AdminGetXPubsCount(ctx context.Context, conditions map[string]interface{},
	metadata *bux.Metadata) (int64, error) {

	return h.adminCount(ctx, conditions, metadata, "/admin/xpubs/count")
}

func (h *TransportHTTP) adminGetModels(ctx context.Context, conditions map[string]interface{},
	metadata *bux.Metadata, queryParams *datastore.QueryParams, path string, models interface{}) error {

	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldConditions:  conditions,
		FieldMetadata:    processMetadata(metadata),
		FieldQueryParams: queryParams,
	})
	if err != nil {
		return err
	}

	if err = h.doHTTPRequest(
		ctx, http.MethodGet, path, jsonStr, h.xPriv, true, &models,
	); err != nil {
		return err
	}
	if h.debug {
		log.Printf(path+": %v\n", models)
	}

	return nil
}

func (h *TransportHTTP) adminCount(ctx context.Context, conditions map[string]interface{}, metadata *bux.Metadata, path string) (int64, error) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldConditions: conditions,
		FieldMetadata:   processMetadata(metadata),
	})
	if err != nil {
		return 0, err
	}

	var count int64
	if err = h.doHTTPRequest(
		ctx, http.MethodGet, path, jsonStr, h.xPriv, true, &count,
	); err != nil {
		return 0, err
	}
	if h.debug {
		log.Printf(path+": %v\n", count)
	}

	return count, nil
}

// AdminRecordTransaction will record a transaction as an admin
func (h *TransportHTTP) AdminRecordTransaction(ctx context.Context, hex string) (*bux.Transaction, error) {

	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldHex: hex,
	})
	if err != nil {
		return nil, err
	}

	var transaction bux.Transaction
	if err = h.doHTTPRequest(
		ctx, http.MethodPost, "/admin/transactions/record", jsonStr, h.xPriv, h.signRequest, &transaction,
	); err != nil {
		return nil, err
	}
	if h.debug {
		log.Printf("transaction: %s\n", transaction.ID)
	}

	return &transaction, nil
}
