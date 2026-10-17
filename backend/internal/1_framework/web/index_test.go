package web

import (
	"strings"
	"testing"

	domain "github.com/YujiYabe/restaurant-order-server/backend/internal/4_domain"
)

func TestBuildIndexPageCreatesPlaceLinksGroupedByPlaceGroup(t *testing.T) {
	t.Parallel()
	placeList := domain.PlaceList{
		"b2":    {PlaceGroup: "B", PlaceName: "B 2", Permission: "permission/b"},
		"a2":    {PlaceGroup: "A", PlaceName: "A2", Permission: "permission a2"},
		"a1":    {PlaceGroup: "A", PlaceName: "A1", Permission: "permission a"},
		"empty": nil,
	}

	page := buildIndexPage(placeList, "device key")
	if len(page.PlaceGroups) != 2 {
		t.Fatalf("place group count = %d, want 2", len(page.PlaceGroups))
	}
	if page.PlaceGroups[0].Name != "A" || page.PlaceGroups[1].Name != "B" {
		t.Fatalf("place groups are not sorted: %#v", page.PlaceGroups)
	}
	if len(page.PlaceGroups[0].Places) != 2 || page.PlaceGroups[0].Places[0].Name != "A1" || page.PlaceGroups[0].Places[1].Name != "A2" {
		t.Fatalf("places in group A are not sorted: %#v", page.PlaceGroups[0].Places)
	}
	if !strings.Contains(page.PlaceGroups[0].Places[0].WaiterURL, "device_key=device+key") {
		t.Fatalf("waiter URL does not preserve device key: %s", page.PlaceGroups[0].Places[0].WaiterURL)
	}
	if !strings.Contains(page.PlaceGroups[1].Places[0].MobileOrderURL, "permission=permission%2Fb") {
		t.Fatalf("mobile order URL does not include escaped permission: %s", page.PlaceGroups[1].Places[0].MobileOrderURL)
	}
}

func TestResolveIndexDeviceKey(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		environment   string
		requestedKey  string
		configuredKey string
		want          string
	}{
		{name: "local default", environment: "local", configuredKey: "configured", want: "configured"},
		{name: "local requested", environment: "local", requestedKey: "requested", configuredKey: "configured", want: "requested"},
		{name: "production requires request", environment: "prod", configuredKey: "configured", want: ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := resolveIndexDeviceKey(test.environment, test.requestedKey, test.configuredKey)
			if got != test.want {
				t.Fatalf("resolveIndexDeviceKey() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestBuildIndexURL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		environment   string
		configuredURL string
		port          string
		want          string
	}{
		{name: "local", environment: "local", configuredURL: "192.0.2.1:6789", port: "6789", want: "https://localhost:6789/"},
		{name: "production host", environment: "prod", configuredURL: "example.com", port: "443", want: "https://example.com/"},
		{name: "configured scheme", environment: "dev", configuredURL: "https://dev.example.com/", port: "443", want: "https://dev.example.com/"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := buildIndexURL(test.environment, test.configuredURL, test.port)
			if got != test.want {
				t.Fatalf("buildIndexURL() = %q, want %q", got, test.want)
			}
		})
	}
}
