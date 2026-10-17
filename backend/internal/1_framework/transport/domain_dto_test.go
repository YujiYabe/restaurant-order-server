package transport

import (
	"encoding/json"
	"strings"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestOrderJSONUsesExistingWireNames(t *testing.T) {
	data, err := json.Marshal(Order{JANCode: 123, Count: 2})
	if err != nil {
		t.Fatal(err)
	}
	if got := string(data); !strings.Contains(got, `"jan_code":123`) || !strings.Contains(got, `"count":2`) {
		t.Fatalf("unexpected JSON: %s", got)
	}
}

func TestPlaceYAMLConvertsToDomain(t *testing.T) {
	var places PlaceList
	if err := yaml.Unmarshal([]byte("A1:\n  place_name: A1\n  place_status: 1\n  capacity: 4\n"), &places); err != nil {
		t.Fatal(err)
	}
	place := places.Domain()["A1"]
	if place == nil || place.PlaceName != "A1" || place.PlaceStatus != 1 || place.Capacity != 4 {
		t.Fatalf("unexpected place: %#v", place)
	}
}
