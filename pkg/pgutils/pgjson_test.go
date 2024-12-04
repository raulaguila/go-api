package pgutils

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"testing"
)

func TestJSONB_Value(t *testing.T) {
	tests := []struct {
		name string
		jb   JSONB
		want driver.Value
		err  error
	}{
		{"empty map", JSONB{}, []byte("{}"), nil},
		{"single element", JSONB{"key": "value"}, []byte(`{"key":"value"}`), nil},
		{"nested map", JSONB{"outer": map[string]any{"inner": "value"}}, []byte(`{"outer":{"inner":"value"}}`), nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.jb.Value()
			if (err != nil) != (tt.err != nil) {
				t.Errorf("JSONB.Value() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !jsonEqual(got.([]byte), tt.want.([]byte)) {
				t.Errorf("JSONB.Value() = %v, want %v", string(got.([]byte)), string(tt.want.([]byte)))
			}
		})
	}
}

func TestJSONB_Scan(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		want    JSONB
		wantErr error
	}{
		{"valid json", []byte(`{"key":"value"}`), JSONB{"key": "value"}, nil},
		{"invalid json", []byte(`invalid`), JSONB{}, errors.New("type assertion to []byte failed")},
		{"non-byte input", "string", JSONB{}, errors.New("type assertion to []byte failed")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var jb JSONB
			err := jb.Scan(tt.input)
			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("JSONB.Scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && !jsonEqualMaps(jb, tt.want) {
				t.Errorf("JSONB.Scan() = %v, want %v", jb, tt.want)
			}
		})
	}
}

func jsonEqual(a, b []byte) bool {
	var o1, o2 map[string]any
	if err1, err2 := json.Unmarshal(a, &o1), json.Unmarshal(b, &o2); err1 != nil || err2 != nil {
		return false
	}
	return jsonEqualMaps(o1, o2)
}

func jsonEqualMaps(m1, m2 map[string]any) bool {
	return len(m1) == len(m2) && jsonUnorderedEqual(m1, m2) && jsonUnorderedEqual(m2, m1)
}

func jsonUnorderedEqual(m1, m2 map[string]any) bool {
	for k, v1 := range m1 {
		if v2, exists := m2[k]; !exists || !jsonDeepEqual(v1, v2) {
			return false
		}
	}
	return true
}

func jsonDeepEqual(v1, v2 any) bool {
	switch v1 := v1.(type) {
	case map[string]any:
		v2, ok := v2.(map[string]any)
		return ok && jsonEqualMaps(v1, v2)
	case []any:
		v2, ok := v2.([]any)
		if !ok || len(v1) != len(v2) {
			return false
		}
		for i := range v1 {
			if !jsonDeepEqual(v1[i], v2[i]) {
				return false
			}
		}
		return true
	default:
		return v1 == v2
	}
}
