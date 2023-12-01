package main

import (
	"os"
	"testing"
)

func Test_ParseEml(t *testing.T) {
	expected := Email{
		From:    "alice@local.net",
		To:      "bob@local.net",
		Date:    "Thu, 16 Jan 2020 18:40:05 +0000",
		Subject: "test 1",
	}
	emlBs, err := os.ReadFile("fixtures/test1.eml")
	if err != nil {
		t.Errorf("Failed to read test email file")
	}

	actual, err := ParseEml(string(emlBs))
	if err != nil {
		t.Error("Failed to parse email string")
	}

	if actual != expected {
		t.Errorf("Expected %v, but got %v", expected, actual)
	}
}
