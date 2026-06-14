package yesno

import (
	"testing"
)

// These tests are offline: they exercise the URI driver's pure string functions.
// HTTP behaviour is covered in yesno_test.go.

func TestDomainInfo(t *testing.T) {
	info := Domain{}.Info()
	if info.Scheme != "yesno" {
		t.Errorf("Scheme = %q, want yesno", info.Scheme)
	}
	if len(info.Hosts) == 0 || info.Hosts[0] != Host {
		t.Errorf("Hosts = %v, want [%s]", info.Hosts, Host)
	}
	if info.Identity.Binary != "yesno" {
		t.Errorf("Identity.Binary = %q, want yesno", info.Identity.Binary)
	}
}

func TestClassify(t *testing.T) {
	cases := []struct{ in, typ, id string }{
		{"yes", "answer", "yes"},
		{"no", "answer", "no"},
	}
	for _, tc := range cases {
		typ, id, err := Domain{}.Classify(tc.in)
		if err != nil || typ != tc.typ || id != tc.id {
			t.Errorf("Classify(%q) = (%q, %q, %v), want (%q, %q, nil)",
				tc.in, typ, id, err, tc.typ, tc.id)
		}
	}
}

func TestClassifyEmptyErrors(t *testing.T) {
	_, _, err := Domain{}.Classify("")
	if err == nil {
		t.Error("expected error for empty input, got nil")
	}
}

func TestLocate(t *testing.T) {
	got, err := Domain{}.Locate("answer", "anything")
	want := BaseURL + "/api"
	if err != nil || got != want {
		t.Errorf("Locate = (%q, %v), want (%q, nil)", got, err, want)
	}
}

func TestLocateUnknownType(t *testing.T) {
	_, err := Domain{}.Locate("unknown", "foo")
	if err == nil {
		t.Error("expected error for unknown type, got nil")
	}
}
