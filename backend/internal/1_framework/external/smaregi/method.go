package smaregi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	netUrl "net/url"
	"strconv"
	"strings"
	"time"

	domain "github.com/YujiYabe/restaurant-order-server/backend/internal/4_domain"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/env"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/pkg"
)

var (
	contractID   = env.Env.ContractId
	clientID     = env.Env.ClientId
	clientSecret = env.Env.ClientSecret
	authURL      = env.Env.AuthURL + contractID + "/token" //
	apiBaseURL   = env.Env.ApiBaseURL + contractID
	scope        = "pos.transactions:read pos.transactions:write"
	grantType    = "client_credentials"
)

// SendPostTransactionTemporary は、注文内容をスマレジへ仮販売登録し、取引情報を返します。
func (receiver *Smaregi) SendPostTransactionTemporary(
	ctx context.Context,
	reqPostTransactionTemporary domain.ReqPostTransactionTemporary,
) (
	smaregiPostResponseEnquiry domain.SmaregiPostResponseEnquiry,
	err error,
) {
	smaregiPostResponseEnquiry = domain.SmaregiPostResponseEnquiry{}
	if env.Env.Environment == env.LOCAL {
		return localTransactionResponse(env.Env.StoreID, reqPostTransactionTemporary), nil
	}

	url, err := netUrl.Parse(
		apiBaseURL + "/pos/transactions/temporaries",
	)
	if err != nil {
		pkg.Logging(ctx, err)
		return
	}

	details := []Detail{}
	for _, domainDetail := range reqPostTransactionTemporary.Details {

		details = append(
			details,
			Detail{
				TransactionDetailId:       strconv.Itoa(domainDetail.TransactionDetailId),
				SalesPrice:                strconv.Itoa(domainDetail.SalesPrice),
				Quantity:                  strconv.Itoa(domainDetail.Quantity),
				ProductId:                 strconv.Itoa(domainDetail.ProductId),
				TransactionDetailDivision: "1",
			},
		)
	}

	postTransactionTemporaryReq := &PostTransactionTemporaryReq{
		// 固定値
		TransactionHeadDivision: transactionHeadDivision,
		Status:                  status,
		StoreId:                 env.Env.StoreID,
		TerminalId:              terminalId,

		Subtotal: strconv.Itoa(reqPostTransactionTemporary.Subtotal), // 小計 明細合計と異なる場合エラー
		Total:    strconv.Itoa(reqPostTransactionTemporary.Total),    // 合計
		// 小計     - 小計値引き            - ポイント値引き + 端数値引額    + 外税額     - 免税額 を設定
		// 計算式：Subtotal - SubtotalDiscountPrice - PointDiscount + RoundingPrice + TaxExclude - TaxFree

		// 変動値
		TerminalTranId:       strconv.Itoa(reqPostTransactionTemporary.TerminalTranId), // 端末取引ID (連携元システムで任意の連番を設定してください）
		TerminalTranDateTime: reqPostTransactionTemporary.TerminalTranDateTime,         // 端末取引日時
		CustomerId:           strconv.Itoa(reqPostTransactionTemporary.CustomerId),     // 会員ID
		CustomerCode:         strconv.Itoa(reqPostTransactionTemporary.CustomerId),     // 会員コード
		Details:              details,
	}

	postTransactionTemporaryReqJson, err := json.Marshal(postTransactionTemporaryReq)
	if err != nil {
		pkg.Logging(ctx, err)
		return
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		url.String(),
		bytes.NewBuffer(postTransactionTemporaryReqJson),
	)
	if err != nil {
		pkg.Logging(ctx, err)
		return
	}

	err = receiver.refreshAccessToken(ctx)
	if err != nil {
		pkg.Logging(ctx, err)
		return
	}

	req.Header.Add("Authorization", "Bearer "+receiver.accessToken)

	client := &http.Client{Timeout: 15 * time.Second}
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

	if false {
		dumpResp, dumpErr := httputil.DumpResponse(resp, true)
		if dumpErr != nil {
			pkg.Logging(ctx, dumpErr)
			err = dumpErr
			return
		}

		pkg.Logging(ctx, string(dumpResp))
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("unable to get this url : http status %d", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		pkg.Logging(ctx, err)
		return
	}

	postTransactionTemporaryRes := PostTransactionTemporaryRes{}
	err = json.Unmarshal(
		body,
		&postTransactionTemporaryRes,
	)
	if err != nil {
		pkg.Logging(ctx, err)
		return
	}
	if postTransactionTemporaryRes.CustomerId == nil || *postTransactionTemporaryRes.CustomerId == "" {
		return smaregiPostResponseEnquiry, fmt.Errorf("Smaregi response has no customerId")
	}
	if postTransactionTemporaryRes.TransactionHeadId == "" {
		return smaregiPostResponseEnquiry, fmt.Errorf("Smaregi response has no transactionHeadId")
	}

	smaregiPostResponseEnquiry.CustomerId, err = strconv.Atoi(*postTransactionTemporaryRes.CustomerId)
	if err != nil {
		pkg.Logging(ctx, err)
		return
	}

	smaregiPostResponseEnquiry.TransactionHeadId = postTransactionTemporaryRes.TransactionHeadId

	return
}

func localTransactionResponse(
	storeID string,
	request domain.ReqPostTransactionTemporary,
) domain.SmaregiPostResponseEnquiry {
	return domain.SmaregiPostResponseEnquiry{
		CustomerId: request.CustomerId,
		TransactionHeadId: fmt.Sprintf(
			"local-%s-%d-%d",
			storeID,
			request.CustomerId,
			request.TerminalTranId,
		),
	}
}

// refreshAccessToken は、スマレジのアクセストークンが未取得または期限間近の場合に更新します。
func (receiver *Smaregi) refreshAccessToken(
	ctx context.Context,
) (
	err error,
) {
	if receiver.accessToken == "" {
		err = receiver.updateAccessToken(ctx)
		return
	}

	// 現在時刻が前回認証時刻より50分過ぎていれば更新
	now := time.Now()
	if now.After(receiver.accessTokenExpired.Add(50 * time.Minute)) {
		err = receiver.updateAccessToken(ctx)
		return
	}

	return
}

// updateAccessToken は、スマレジ認証APIから新しいアクセストークンを取得します。
func (receiver *Smaregi) updateAccessToken(
	ctx context.Context,
) (
	err error,
) {
	values := netUrl.Values{}
	values.Add("grant_type", grantType)
	values.Add("scope", scope)

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		authURL,
		strings.NewReader(values.Encode()),
	)
	if err != nil {
		pkg.Logging(ctx, err)
		return
	}

	// Content-Type 設定
	req.Header.Set(
		"Content-Type",
		"application/x-www-form-urlencoded",
	)
	req.SetBasicAuth(clientID, clientSecret) //

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
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("Smaregi token endpoint returned HTTP %d", resp.StatusCode)
	}

	byteArray, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read Smaregi token response: %w", err)
	}
	var cre credentialInfo

	err = json.Unmarshal(
		byteArray,
		&cre,
	)
	if err != nil {
		pkg.Logging(ctx, err)
		return
	}

	if cre.AccessToken == "" {
		return fmt.Errorf("Smaregi token response has no access token")
	}
	receiver.accessToken = cre.AccessToken
	receiver.accessTokenExpired = time.Now()

	return
}
