package safehttp

import (
	"testing"
)

func TestGet(t *testing.T) {
	resp, err := Get("https://bertha-scale.va.opower.it/v1/executions", 10, 10)

	if err != nil {
		t.Error(err)
	}

	t.Log(resp)
}


func TestGet_BadUrl(t *testing.T) {
	resp, err := Get("https://unknown.1/", 10, 10)

	if err == nil {
		t.Errorf("Should error out but got response: %s", resp)
	}
}
