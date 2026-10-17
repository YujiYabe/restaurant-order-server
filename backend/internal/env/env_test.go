package env

import "testing"

// TestInitEnvOperateLoadsSecretsFromEnvironment は、対象仕様の回帰動作を検証します。
func TestInitEnvOperateLoadsSecretsFromEnvironment(t *testing.T) {
	values := map[string]string{
		"BASIC_AUTH_USERNAME":             "test-user",
		"BASIC_AUTH_PASSWORD":             "test-password",
		"FROM_GATE_KEY":                   "test-gate-key",
		"FROM_DEVICE_KEY":                 "test-device-key",
		"CONTRACT_ID":                     "test-contract-id",
		"CLIENT_ID":                       "test-client-id",
		"CLIENT_SECRET":                   "test-client-secret",
		"GATE_TRANSACTION_INFO_BASE_URL":  "https://example.invalid",
		"GATE_TRANSACTION_INFO_PATH":      "/transactions",
		"AZURE_DOCUMENT_INTELLIGENCE_URL": "https://example.invalid/analyze",
		"AZURE_SUBSCRIPTION_KEY":          "test-azure-key",
	}
	for key, value := range values {
		t.Setenv(key, value)
	}

	initEnvOperate()

	if Env.BasicAuthUserName != values["BASIC_AUTH_USERNAME"] ||
		Env.BasicAuthPassword != values["BASIC_AUTH_PASSWORD"] ||
		Env.FromGateKey != values["FROM_GATE_KEY"] ||
		Env.FromDeviceKey != values["FROM_DEVICE_KEY"] ||
		Env.ContractId != values["CONTRACT_ID"] ||
		Env.ClientId != values["CLIENT_ID"] ||
		Env.ClientSecret != values["CLIENT_SECRET"] ||
		Env.GateTransactionInfoBaseURL != values["GATE_TRANSACTION_INFO_BASE_URL"] ||
		Env.GateTransactionInfoPath != values["GATE_TRANSACTION_INFO_PATH"] ||
		Env.AzureDocumentIntelligenceURL != values["AZURE_DOCUMENT_INTELLIGENCE_URL"] ||
		Env.AzureSubscriptionKey != values["AZURE_SUBSCRIPTION_KEY"] {
		t.Fatal("runtime secrets were not loaded from environment variables")
	}
}

func TestInitEnvOperateLoadsServerPortFromEnvironment(t *testing.T) {
	t.Setenv("SERVER_PORT", "7443")

	initEnvOperate()

	if Env.ServerPort != "7443" {
		t.Fatalf("ServerPort = %q, want %q", Env.ServerPort, "7443")
	}
}
