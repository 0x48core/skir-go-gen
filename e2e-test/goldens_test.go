package main

import (
	"encoding/json"
	"fmt"
	"testing"

	goldens "skir-go-gen-e2e/skirout/external/gepheum/skir_golden_tests"
)

func TestUnitTestsLoaded(t *testing.T) {
	if len(goldens.UnitTests) == 0 {
		t.Fatal("UNIT_TESTS constant should not be empty")
	}
	for i, ut := range goldens.UnitTests {
		if i == 0 {
			continue
		}
		if ut.TestNumber != goldens.UnitTests[0].TestNumber+int32(i) {
			t.Fatalf(
				"Test numbers are not sequential at index %d: found %d, expected %d",
				i, ut.TestNumber, goldens.UnitTests[0].TestNumber+int32(i),
			)
		}
	}
	t.Logf("Loaded %d golden unit tests", len(goldens.UnitTests))
}

func TestPointJSONRoundTrip(t *testing.T) {
	original := goldens.Point{
		X: 10,
		Y: 20,
		Color: goldens.Color{R: 255, G: 128, B: 0},
	}
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal Point: %v", err)
	}
	var decoded goldens.Point
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal Point: %v", err)
	}
	if !original.Equal(decoded) {
		t.Errorf("Point round-trip failed: got %+v, want %+v", decoded, original)
	}
}

func TestColorJSONRoundTrip(t *testing.T) {
	original := goldens.Color{R: 100, G: 200, B: 50}
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal Color: %v", err)
	}
	var decoded goldens.Color
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal Color: %v", err)
	}
	if !original.Equal(decoded) {
		t.Errorf("Color round-trip failed: got %+v, want %+v", decoded, original)
	}
}

func TestMyEnumConstantVariantJSON(t *testing.T) {
	original := goldens.NewMyEnumOk()
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal MyEnum OK: %v", err)
	}
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		t.Fatalf("Constant enum variant should serialize to a string: %v", err)
	}
	if str != "OK" {
		t.Errorf("Expected \"OK\", got %q", str)
	}

	var decoded goldens.MyEnum
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal MyEnum: %v", err)
	}
	if decoded.Kind() != goldens.MyEnumKindOk {
		t.Errorf("Expected kind OK, got %v", decoded.Kind())
	}
}

func TestMyEnumWrapperVariantJSON(t *testing.T) {
	color := goldens.Color{R: 1, G: 2, B: 3}
	original := goldens.NewMyEnumColor(color)
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal MyEnum with Color: %v", err)
	}

	var decoded goldens.MyEnum
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal MyEnum: %v", err)
	}
	if decoded.Kind() != goldens.MyEnumKindColor {
		t.Errorf("Expected kind color, got %v", decoded.Kind())
	}
	val, ok := decoded.AsColor()
	if !ok {
		t.Fatal("Expected AsColor() to return true")
	}
	if !val.Equal(color) {
		t.Errorf("Color value mismatch: got %+v, want %+v", val, color)
	}
}

func TestMyEnumBoolVariantJSON(t *testing.T) {
	original := goldens.NewMyEnumBool(true)
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal MyEnum with bool: %v", err)
	}

	var decoded goldens.MyEnum
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal MyEnum: %v", err)
	}
	if decoded.Kind() != goldens.MyEnumKindBool {
		t.Errorf("Expected kind bool, got %v", decoded.Kind())
	}
	val, ok := decoded.AsBool()
	if !ok {
		t.Fatal("Expected AsBool() to return true")
	}
	if val != true {
		t.Errorf("Bool value mismatch: got %v, want true", val)
	}
}

func TestMyEnumUnknownJSON(t *testing.T) {
	unknown := goldens.NewMyEnumUnknown()
	data, err := json.Marshal(unknown)
	if err != nil {
		t.Fatalf("Failed to marshal MyEnum Unknown: %v", err)
	}
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		t.Fatalf("Unknown enum variant should serialize to a string: %v", err)
	}
	if str != "UNKNOWN" {
		t.Errorf("Expected \"UNKNOWN\", got %q", str)
	}

	// Deserializing an unrecognized string should produce UNKNOWN kind
	unrecognized := []byte(`"SOMETHING_NEW"`)
	var decoded goldens.MyEnum
	if err := json.Unmarshal(unrecognized, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal unrecognized enum: %v", err)
	}
	if decoded.Kind() != goldens.MyEnumKindUnknown {
		t.Errorf("Expected kind UNKNOWN for unrecognized string, got %v", decoded.Kind())
	}
}

func TestPointDefaultValues(t *testing.T) {
	p := goldens.NewPoint()
	if p.X != 0 || p.Y != 0 {
		t.Errorf("Default Point should have zero values, got X=%d, Y=%d", p.X, p.Y)
	}
}

func TestPointJSONFieldNames(t *testing.T) {
	p := goldens.Point{X: 1, Y: 2, Color: goldens.Color{R: 0, G: 0, B: 0}}
	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("Failed to unmarshal as map: %v", err)
	}
	expectedFields := []string{"x", "y", "color"}
	for _, field := range expectedFields {
		if _, ok := raw[field]; !ok {
			t.Errorf("Expected JSON field %q, but not found in %s", field, string(data))
		}
	}
}

func TestColorJSONFieldNames(t *testing.T) {
	c := goldens.Color{R: 1, G: 2, B: 3}
	data, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("Failed to unmarshal as map: %v", err)
	}
	expectedFields := []string{"r", "g", "b"}
	for _, field := range expectedFields {
		if _, ok := raw[field]; !ok {
			t.Errorf("Expected JSON field %q, but not found in %s", field, string(data))
		}
	}
}

func TestRecStructJSONRoundTrip(t *testing.T) {
	inner := goldens.RecStruct{B: true}
	original := goldens.RecStruct{A: &inner, B: false}
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal RecStruct: %v", err)
	}
	var decoded goldens.RecStruct
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal RecStruct: %v", err)
	}
	if decoded.B != original.B {
		t.Errorf("RecStruct.B mismatch: got %v, want %v", decoded.B, original.B)
	}
}

func TestRecEnumJSONRoundTrip(t *testing.T) {
	original := goldens.NewRecEnumFoo()
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal RecEnum FOO: %v", err)
	}
	var decoded goldens.RecEnum
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal RecEnum: %v", err)
	}
	if decoded.Kind() != goldens.RecEnumKindFoo {
		t.Errorf("Expected kind FOO, got %v", decoded.Kind())
	}
}

func TestReserializeValueJSONRoundTrips(t *testing.T) {
	for _, ut := range goldens.UnitTests {
		assertion := ut.Assertion
		if assertion.Kind() != goldens.AssertionKindReserializeValue {
			continue
		}
		rv, ok := assertion.AsReserializeValue()
		if !ok {
			continue
		}

		t.Run(fmt.Sprintf("test_%d", ut.TestNumber), func(t *testing.T) {
			typedValue := rv.Value
			value := extractTypedValue(t, typedValue)
			if value == nil {
				t.Skip("Unsupported TypedValue kind for Go JSON test")
				return
			}

			data, err := json.Marshal(value)
			if err != nil {
				t.Fatalf("Failed to marshal: %v", err)
			}

			if len(rv.ExpectedDenseJson) > 0 || len(rv.ExpectedReadableJson) > 0 {
				t.Logf("JSON output: %s", string(data))
			}

			roundTripped := newZeroValue(value)
			if roundTripped == nil {
				return
			}
			if err := json.Unmarshal(data, roundTripped); err != nil {
				t.Errorf("Failed to unmarshal JSON round-trip: %v (json: %s)", err, string(data))
			}
		})
	}
}

func extractTypedValue(t *testing.T, tv goldens.TypedValue) interface{} {
	t.Helper()
	switch tv.Kind() {
	case goldens.TypedValueKindPoint:
		v, _ := tv.AsPoint()
		return v
	case goldens.TypedValueKindColor:
		v, _ := tv.AsColor()
		return v
	case goldens.TypedValueKindMyEnum:
		v, _ := tv.AsMyEnum()
		return v
	case goldens.TypedValueKindRecStruct:
		v, _ := tv.AsRecStruct()
		return v
	case goldens.TypedValueKindRecEnum:
		v, _ := tv.AsRecEnum()
		return v
	default:
		return nil
	}
}

func newZeroValue(v interface{}) interface{} {
	switch v.(type) {
	case goldens.Point:
		return &goldens.Point{}
	case goldens.Color:
		return &goldens.Color{}
	case goldens.MyEnum:
		return &goldens.MyEnum{}
	case goldens.RecStruct:
		return &goldens.RecStruct{}
	case goldens.RecEnum:
		return &goldens.RecEnum{}
	default:
		return nil
	}
}
