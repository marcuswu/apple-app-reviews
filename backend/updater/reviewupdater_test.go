package updater

import (
	"os"
	"testing"
	"time"
)

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

type appWithAge struct {
	id           string
	ageInSeconds int
}

func setupNextAppTest(apps []appWithAge) {
	for _, app := range apps {
		file := fileForAppId(app.id)
		os.Remove(file)
		os.Create(file)
		os.Chtimes(file, time.Now(), time.Now().Add(time.Duration(-app.ageInSeconds)*time.Second))
	}
}

func teardownNextAppTest(apps []appWithAge) {
	for _, app := range apps {
		os.Remove(fileForAppId(app.id))
	}
}

func TestNextApp(t *testing.T) {
	tests := []struct {
		name     string
		input    []appWithAge
		expected string
		error    bool
	}{
		{"one app no refresh", []appWithAge{{"123456789", 300}}, "123456789", true},
		{"one app with refresh", []appWithAge{{"123456789", 601}}, "123456789", false},
		{"three app no refresh", []appWithAge{
			{"123456789", 301},
			{"987654321", 302},
			{"1234567890", 101},
		}, "987654321", true},
		{"three app last refresh", []appWithAge{
			{"123456789", 601},
			{"987654321", 302},
			{"1234567890", 1001},
		}, "1234567890", false},
		{"three app first refresh", []appWithAge{
			{"123456789", 801},
			{"987654321", 302},
			{"1234567890", 701},
		}, "123456789", false},
	}

	for _, test := range tests {
		setupNextAppTest(test.input)
		defer teardownNextAppTest(test.input)

		files := make([]string, 0, len(test.input))
		for _, app := range test.input {
			files = append(files, fileForAppId(app.id))
		}
		next, err := nextApp(files)

		if !test.error && err != nil {
			t.Errorf("test \"%s\" should not encounter an error, but did: %s", test.name, err)
		}
		if test.error && err == nil {
			t.Errorf("test \"%s\" should encounter an error, but did not", test.name)
		}
		if test.expected != next {
			t.Errorf("test \"%s\" expected %s, but found %s", test.name, test.expected, next)
		}
	}
}
