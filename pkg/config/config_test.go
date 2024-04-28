package config

import (
	"os"
	"testing"
)

func TestReadConfig(t *testing.T) {
	l := NewIConfig()
	dir, _ := os.Getwd()
	if err := l.ReadConfig(dir, "test_config"); err != nil {
		t.Fatalf(`Failed to ReadConfig(. , test_config). Got : '%s', want : nil`, err.Error())
	}
}

func TestRead(t *testing.T) {
	l := NewIConfig()

	dir, _ := os.Getwd()
	if err := l.ReadConfig(dir, "test_config"); err != nil {
		t.Fatalf(`Failed to ReadConfig(. , test_config). Got : '%s', want : nil`, err.Error())
	}

	if value, err := l.ReadKey("foo"); err != nil {
		t.Fatalf(`Failed to ReadKey(foo). Got : '%s', want : nil`, err.Error())
	} else {
		if value.(string) != "bar" {
			t.Fatalf(`value error: got %s, want bar`, value)
		}
	}
}
func TestWrite(t *testing.T) {
	l := NewIConfig()

	dir, _ := os.Getwd()
	if err := l.ReadConfig(dir, "test_config"); err != nil {
		t.Fatalf(`Failed to ReadConfig(. , test_config). Got : '%s', want : nil`, err.Error())
	}

	if err := l.WriteKeyValue("winnie", "pooh"); err != nil {
		t.Fatalf(`Failed to WriteKeyValue(winnie, pooh). Got : '%s', want : nil`, err.Error())
	} else {
		if value, err := l.ReadKey("winnie"); err != nil {
			t.Fatalf(`Failed to ReadKey(winnie). Got : '%s', want : nil`, err.Error())
		} else {
			if value.(string) != "pooh" {
				t.Fatalf(`value error: got %s, want pooh`, value)
			}
		}
	}
}
