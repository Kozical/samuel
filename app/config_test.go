package app

import "testing"

var unmarshalConfigTests = []struct {
	id       int
	yaml     []byte
	expected *Config
}{
	{1, []byte("%"), nil},
	{2, []byte("---\n"), &Config{}},
	{3, []byte("---\ncrt_path: C:\\test\\path"), &Config{CrtPath: "C:\\test\\path"}},
	{4, []byte("---\nkey_path: C:\\test\\path"), &Config{KeyPath: "C:\\test\\path"}},
	{5, []byte("---\ncrt_path: C:\\test\\path\nkey_path: C:\\test\\path"), &Config{CrtPath: "C:\\test\\path", KeyPath: "C:\\test\\path"}},
	{5, []byte("---\ncrt_path: C:\\test\\path\nkey_path: C:\\test\\path\nendpoint: 127.0.0.1:8104"), &Config{CrtPath: "C:\\test\\path", KeyPath: "C:\\test\\path", Endpoint: "127.0.0.1:8104"}},
}

func TestUnmarshalConfig(t *testing.T) {
	for _, tt := range unmarshalConfigTests {
		actual, err := unmarshalConfig(tt.yaml)
		if err != nil && tt.expected != nil {
			t.Errorf("[id: %d] unmarshalConfig: failed with %s", tt.id, err)
			continue
		}
		if err != nil && tt.expected == nil {
			continue
		}
		if actual.KeyPath != tt.expected.KeyPath ||
			actual.CrtPath != tt.expected.CrtPath {
			t.Errorf("[id: %d] unmarshalConfig: expected %q, actual %q", tt.id, tt.expected, actual)
		}
	}
}
