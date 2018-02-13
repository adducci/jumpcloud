package password

import (
	"bytes"
	"testing"
)

//Test stringToBytes
func TestPassword1(t *testing.T) {
	cases := []struct {
		in   string
		want []byte
	}{
		{"angryMonkey", []byte("angryMonkey")},
	}
	for _, c := range cases {
		got := stringToBytes(c.in)
		if !bytes.Equal(got, c.want) {
			t.Errorf("stringToBytes(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}

//Test bytesToBase64String
func TestPassword2(t *testing.T) {
	cases := []struct {
		in   []byte
		want string
	}{
		{[]byte("angryMonkey"), "YW5ncnlNb25rZXk="},
	}
	for _, c := range cases {
		got := bytesToBase64String(c.in)
		if got != c.want {
			t.Errorf("bytesToBase64String(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}

//Test Encrypt
func TestPassword3(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"angryMonkey", "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="},
	}
	for _, c := range cases {
		got := Encrypt(c.in)
		if got != c.want {
			t.Errorf("Encrypt(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}
