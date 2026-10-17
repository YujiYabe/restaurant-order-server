//go:build e2e

package e2e

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

var (
	testBaseURL = resolveBaseURL()
	testClient  = newHTTPClient(5 * time.Second)
)

func resolveBaseURL() string {
	baseURL := strings.TrimSpace(os.Getenv("E2E_BASE_URL"))
	if baseURL == "" {
		baseURL = "https://127.0.0.1:6789"
	}
	return mustLoopbackURL(strings.TrimRight(baseURL, "/"))
}

func mustLoopbackURL(baseURL string) string {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		panic(fmt.Sprintf("E2E_BASE_URL is invalid: %s", baseURL))
	}
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		panic(fmt.Sprintf("E2E_BASE_URL must use http or https: %s", baseURL))
	}
	hostname := strings.ToLower(parsedURL.Hostname())
	if hostname == "localhost" {
		return baseURL
	}
	ipAddress := net.ParseIP(hostname)
	if ipAddress == nil || !ipAddress.IsLoopback() {
		panic(fmt.Sprintf("E2E_BASE_URL must point to localhost or a loopback IP: %s", baseURL))
	}
	return baseURL
}

func newHTTPClient(timeout time.Duration) *http.Client {
	transport, ok := http.DefaultTransport.(*http.Transport)
	if !ok {
		panic("http.DefaultTransport is not *http.Transport")
	}
	transport = transport.Clone()
	if strings.HasPrefix(testBaseURL, "https://") {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //nolint:gosec // ローカルの自己署名証明書専用
	}
	return &http.Client{Transport: transport, Timeout: timeout}
}

func requireDeviceKey(t *testing.T) string {
	t.Helper()
	deviceKey := strings.TrimSpace(os.Getenv("E2E_DEVICE_KEY"))
	if deviceKey == "" {
		t.Fatal("E2E_DEVICE_KEY must match the FROM_DEVICE_KEY used by the local backend")
	}
	return deviceKey
}

func requireBillingTestConfiguration(t *testing.T) (janCode int, gateKey string) {
	t.Helper()
	if os.Getenv("E2E_LOCAL_BILLING_STUBS") != "true" {
		t.Fatal("E2E_LOCAL_BILLING_STUBS=true is required to confirm local billing stubs and a disposable SQLite database")
	}

	janCodeText := strings.TrimSpace(os.Getenv("E2E_JAN_CODE"))
	janCode, err := strconv.Atoi(janCodeText)
	if err != nil || janCode <= 0 {
		t.Fatalf("E2E_JAN_CODE must be a positive JAN code: %q", janCodeText)
	}

	gateKey = strings.TrimSpace(os.Getenv("E2E_GATE_KEY"))
	if gateKey == "" {
		t.Fatal("E2E_GATE_KEY must match the FROM_GATE_KEY used by the local backend")
	}
	return janCode, gateKey
}

func request(
	t *testing.T,
	method string,
	path string,
	body any,
) *http.Response {
	t.Helper()
	requestURL, err := url.Parse(testBaseURL + path)
	if err != nil {
		t.Fatalf("parse request URL: %v", err)
	}
	query := requestURL.Query()
	query.Set("device_key", requireDeviceKey(t))
	requestURL.RawQuery = query.Encode()

	var reader io.Reader
	if body != nil {
		data, marshalErr := json.Marshal(body)
		if marshalErr != nil {
			t.Fatalf("marshal request: %v", marshalErr)
		}
		reader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(t.Context(), method, requestURL.String(), reader)
	if err != nil {
		t.Fatalf("build request: %v", err)
	}
	req.SetBasicAuth("e2e", "e2e")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	response, err := testClient.Do(req)
	if err != nil {
		t.Fatalf("request %s %s: %v", method, path, err)
	}
	return response
}

func requestAndExpectStatus(t *testing.T, method string, path string, body any, want int) {
	t.Helper()
	response := request(t, method, path, body)
	defer closeResponseBody(t, response)
	if response.StatusCode == want {
		return
	}
	responseBody, _ := io.ReadAll(io.LimitReader(response.Body, 4096))
	t.Fatalf("status = %d, want %d, body=%s", response.StatusCode, want, strings.TrimSpace(string(responseBody)))
}

func closeResponseBody(t *testing.T, response *http.Response) {
	t.Helper()
	if err := response.Body.Close(); err != nil {
		t.Errorf("close response body: %v", err)
	}
}

func waitForBackend(t *testing.T, timeout time.Duration) {
	t.Helper()
	deviceKey := requireDeviceKey(t)
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		requestURL := testBaseURL + "/v1/place?device_key=" + url.QueryEscape(deviceKey)
		req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, requestURL, nil)
		if err != nil {
			t.Fatalf("build readiness request: %v", err)
		}
		req.SetBasicAuth("e2e", "e2e")
		response, err := testClient.Do(req)
		if err == nil && response.StatusCode == http.StatusOK {
			_ = response.Body.Close()
			return
		}
		if response != nil {
			_ = response.Body.Close()
		}
		time.Sleep(200 * time.Millisecond)
	}
	t.Fatalf("backend did not become ready at %s", testBaseURL)
}

func getPlaces(t *testing.T) map[string]*placeResponse {
	t.Helper()
	response := request(t, http.MethodGet, "/v1/place", nil)
	defer closeResponseBody(t, response)
	if response.StatusCode != http.StatusOK {
		t.Fatalf("GET /v1/place status = %d", response.StatusCode)
	}
	places := map[string]*placeResponse{}
	if err := json.NewDecoder(response.Body).Decode(&places); err != nil {
		t.Fatalf("decode places: %v", err)
	}
	return places
}

func resetPlaces(t *testing.T) {
	t.Helper()
	requestAndExpectStatus(t, http.MethodGet, "/v1/admin/place/reset", nil, http.StatusOK)
}

func patchPlaceStatus(t *testing.T, placeName string, status int) {
	t.Helper()
	path := "/v1/place/" + url.PathEscape(placeName)
	requestAndExpectStatus(t, http.MethodPatch, path, map[string]int{"place_status": status}, http.StatusOK)
}

func patchOrderStatus(t *testing.T, placeName string, orderNumber int, status int) {
	t.Helper()
	path := fmt.Sprintf("/v1/place/%s/order/%d", url.PathEscape(placeName), orderNumber)
	requestAndExpectStatus(t, http.MethodPatch, path, map[string]int{"order_status": status}, http.StatusOK)
}

func completeBilling(t *testing.T, transactionHeadID string, gateKey string) {
	t.Helper()
	requestBody, err := json.Marshal(map[string]string{"transactionHeadId": transactionHeadID})
	if err != nil {
		t.Fatalf("marshal Smaregi webhook: %v", err)
	}
	req, err := http.NewRequestWithContext(
		t.Context(),
		http.MethodPost,
		testBaseURL+"/v1/smaregi",
		bytes.NewReader(requestBody),
	)
	if err != nil {
		t.Fatalf("build Smaregi webhook request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", gateKey)
	response, err := testClient.Do(req)
	if err != nil {
		t.Fatalf("send Smaregi webhook: %v", err)
	}
	defer closeResponseBody(t, response)
	if response.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(io.LimitReader(response.Body, 4096))
		t.Fatalf("Smaregi webhook status = %d, want %d, body=%s", response.StatusCode, http.StatusOK, strings.TrimSpace(string(responseBody)))
	}
}

type sseStream struct {
	response *http.Response
	scanner  *bufio.Scanner
	cancel   context.CancelFunc
}

func openPlaceSSE(t *testing.T, placeName string) *sseStream {
	t.Helper()
	ctx, cancel := context.WithCancel(t.Context())
	endpoint := testBaseURL + "/v1/sse/app?kind=back_yard&place_name=" + url.QueryEscape(placeName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		cancel()
		t.Fatalf("build SSE request: %v", err)
	}
	client := newHTTPClient(0)
	response, err := client.Do(req)
	if err != nil {
		cancel()
		t.Fatalf("open SSE: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(response.Body, 4096))
		statusCode := response.StatusCode
		closeResponseBody(t, response)
		cancel()
		t.Fatalf("SSE status = %d, want %d, body=%s", statusCode, http.StatusOK, strings.TrimSpace(string(body)))
	}
	return &sseStream{response: response, scanner: bufio.NewScanner(response.Body), cancel: cancel}
}

func (stream *sseStream) close() {
	stream.cancel()
	_ = stream.response.Body.Close()
}

func (stream *sseStream) nextPlace(t *testing.T) placeResponse {
	t.Helper()
	result := make(chan placeResponse, 1)
	errors := make(chan error, 1)
	go func() {
		for stream.scanner.Scan() {
			line := stream.scanner.Text()
			if !strings.HasPrefix(line, "data: ") {
				continue
			}
			var place placeResponse
			if err := json.Unmarshal([]byte(strings.TrimPrefix(line, "data: ")), &place); err != nil {
				errors <- err
				return
			}
			result <- place
			return
		}
		errors <- fmt.Errorf("SSE stream ended: %w", stream.scanner.Err())
	}()

	select {
	case place := <-result:
		return place
	case err := <-errors:
		t.Fatal(err)
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for SSE event")
	}
	return placeResponse{}
}
