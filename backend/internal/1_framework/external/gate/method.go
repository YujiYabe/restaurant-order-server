package gate

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	neturl "net/url"
	"time"

	"github.com/YujiYabe/restaurant-order-server/backend/internal/env"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/pkg"
)

// SendPostTransactionInfo は、取引情報を共通サービスへ送信します。
func (receiver *Gate) SendPostTransactionInfo(
	ctx context.Context,
	transactionHeadId string,
) (
	err error,
) {
	if shouldSkipTransactionInfo(env.Env.Environment) {
		return nil
	}

	url, err := transactionInfoURL(
		env.Env.GateTransactionInfoBaseURL,
		env.Env.GateTransactionInfoPath,
	)
	if err != nil {
		pkg.Logging(ctx, err)
		return
	}

	data := struct {
		StoreId           string `json:"storeId"`
		TransactionHeadId string `json:"transactionHeadId"`
	}{
		StoreId:           env.Env.EnvUnique.StoreID,
		TransactionHeadId: transactionHeadId,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		pkg.Logging(ctx, err)
		return
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		url.String(),
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		pkg.Logging(ctx, err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		pkg.Logging(ctx, err)
		return
	}

	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			pkg.Logging(ctx, closeErr)
			if err == nil {
				err = closeErr
			}
		}
	}()
	if !isSuccessfulStatus(resp.StatusCode) {
		return fmt.Errorf("transaction info endpoint returned HTTP %d", resp.StatusCode)
	}

	return
}

func transactionInfoURL(baseURL string, path string) (*neturl.URL, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("GATE_TRANSACTION_INFO_BASE_URL is not set")
	}
	if path == "" {
		return nil, fmt.Errorf("GATE_TRANSACTION_INFO_PATH is not set")
	}

	url, err := neturl.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("parse transaction info base URL: %w", err)
	}
	if url.Scheme == "" || url.Host == "" {
		return nil, fmt.Errorf("GATE_TRANSACTION_INFO_BASE_URL must be an absolute URL")
	}

	reference, err := neturl.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("parse transaction info path: %w", err)
	}
	if reference.IsAbs() || reference.Host != "" {
		return nil, fmt.Errorf("GATE_TRANSACTION_INFO_PATH must be a relative path")
	}

	return url.ResolveReference(reference), nil
}

func shouldSkipTransactionInfo(environment string) bool {
	return environment == env.LOCAL
}

func isSuccessfulStatus(statusCode int) bool {
	return statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices
}
