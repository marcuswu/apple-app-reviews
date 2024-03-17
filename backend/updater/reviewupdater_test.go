package updater

import "testing"

func TestAppFileMatching(t *testing.T) {
	var tests = []struct {
		name         string
		appId        string
		expectedFile string
	}{
		{"app 1234", "1234", "App-1234.json"},
		{"app abcd", "abcd", "App-abcd.json"},
		{"app 1234abc", "1234abc", "App-1234abc.json"},
		{"app empty", "", "App-.json"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans := fileForAppId(tt.appId)
			if ans != tt.expectedFile {
				t.Errorf("got %s, want %s", ans, tt.expectedFile)
			}
		})
	}
}
