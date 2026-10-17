package azure

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	neturl "net/url"
	"time"

	domain "github.com/YujiYabe/restaurant-order-server/backend/internal/4_domain"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/env"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/pkg"
)

const externalHTTPTimeout = 15 * time.Second

type PostCertificate struct {
	Base64Source string `json:"base64Source"`
}

// VerifyCertificate は、本人確認書類画像をAzureへ送信し、生年月日と書類番号を取得します。
func (receiver *Azure) VerifyCertificate(
	ctx context.Context,
	alcohol domain.Alcohol,
) (
	birthday string,
	number string,
	err error,
) {
	urlString := env.Env.AzureDocumentIntelligenceURL
	if urlString == "" {
		return "", "", fmt.Errorf("AZURE_DOCUMENT_INTELLIGENCE_URL is not set")
	}

	url, err := neturl.Parse(urlString)
	if err != nil {
		pkg.Logging(ctx, err)
		return
	}

	imageBytes, err := pkg.DecodeBase64DataURL(alcohol.Base64Source)
	if err != nil {
		pkg.Logging(ctx, err)
		return
	}

	postCertificate := PostCertificate{
		Base64Source: base64.StdEncoding.EncodeToString(imageBytes),
	}

	postCertificateJson, err := json.Marshal(postCertificate)
	if err != nil {
		pkg.Logging(ctx, err)
		return
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		url.String(),
		bytes.NewBuffer(postCertificateJson),
	)
	if err != nil {
		pkg.Logging(ctx, err)
		return
	}

	header := http.Header{}
	header.Set("Ocp-Apim-Subscription-Key", env.Env.AzureSubscriptionKey)
	header.Set("Content-Type", "application/json")
	req.Header = header

	client := &http.Client{Timeout: externalHTTPTimeout}
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
	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("Azure analyze returned HTTP %d", resp.StatusCode)
	}
	operationURL := resp.Header.Get("Operation-Location")
	if operationURL == "" {
		return "", "", fmt.Errorf("Azure analyze response has no Operation-Location")
	}

	if false { // デバッグ用ダンプ
		dumpResp, dumpErr := httputil.DumpResponse(resp, true)
		if dumpErr != nil {
			pkg.Logging(ctx, dumpErr)
			err = dumpErr
			return
		}
		pkg.Logging(ctx, string(dumpResp))
	}

	return receiver.pollCertificateResult(ctx, operationURL, time.Second)
}

func (receiver *Azure) pollCertificateResult(
	ctx context.Context,
	operationURL string,
	interval time.Duration,
) (
	birthday string,
	number string,
	err error,
) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for count := 0; count < 10; count++ {
		select {
		case <-ctx.Done():
			return "", "", ctx.Err()
		case <-ticker.C:
		}
		result, getErr := receiver.GetResultVerify(
			ctx,
			operationURL,
		)
		if getErr != nil {
			pkg.Logging(ctx, result)
			pkg.Logging(ctx, getErr)
			err = getErr
			return
		}
		if result.Status == "succeeded" {
			birthday, number, err = extractCertificateFields(result)
			if err != nil {
				pkg.Logging(ctx, err)
				return
			}
			return birthday, number, nil
		}
		if result.Status == "failed" {
			return "", "", fmt.Errorf("Azure document analysis failed")
		}
	}
	return "", "", fmt.Errorf("Azure document analysis timed out")
}

// extractCertificateFields は、OCR結果から必須の生年月日と書類番号を取り出します。
func extractCertificateFields(result Result) (birthday string, number string, err error) {
	if len(result.AnalyzeResult.Documents) == 0 {
		return "", "", fmt.Errorf("OCR result contains no documents")
	}

	fields := result.AnalyzeResult.Documents[0].Fields
	if fields.Birthday.ValueString == "" {
		return "", "", fmt.Errorf("OCR result contains no birthday")
	}
	if fields.Number.ValueString == "" {
		return "", "", fmt.Errorf("OCR result contains no certificate number")
	}

	return fields.Birthday.ValueString, fields.Number.ValueString, nil
}

// GetResultVerify は、Azureの解析結果URLからOCR処理結果を取得します。
func (receiver *Azure) GetResultVerify(
	ctx context.Context,
	urlString string,
) (
	result Result,
	err error,
) {
	if urlString == "" {
		return result, fmt.Errorf("Azure result URL is empty")
	}
	url, err := neturl.Parse(urlString)
	if err != nil {
		pkg.Logging(ctx, err)
		return
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		url.String(),
		nil,
	)
	if err != nil {
		pkg.Logging(ctx, err)
		return
	}

	header := http.Header{}
	header.Set("Ocp-Apim-Subscription-Key", env.Env.AzureSubscriptionKey)
	req.Header = header

	client := &http.Client{Timeout: externalHTTPTimeout}
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
		return result, fmt.Errorf("Azure result returned HTTP %d", resp.StatusCode)
	}
	byteArray, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, fmt.Errorf("read Azure result: %w", err)
	}

	err = json.Unmarshal(
		byteArray,
		&result,
	)
	if err != nil {
		pkg.Logging(ctx, err)
		return
	}

	if false {
		dumpResp, dumpErr := httputil.DumpResponse(resp, true)

		if dumpErr != nil {
			pkg.Logging(ctx, dumpErr)
			err = dumpErr
			return
		}

		pkg.Logging(ctx, string(dumpResp))
	}

	return
}
