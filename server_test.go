package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestQuayBuildSuccessPostHandler(t *testing.T) {
	req, _ := http.NewRequest("POST",
		"http://localhost:8080/success",
		bytes.NewBuffer(MockQuayBuildSuccessHookJSON()))
	req.Header.Add("Content-Type", "application/json;charset=utf-8")

	w := httptest.NewRecorder()

	QuayBuildSuccessPostHandler(w, req)

	v := map[string]string{}

	_ = DecodeJSON(w.Body.Bytes(), &v)

	if v["status"] != "ok" {
		t.Error(fmt.Sprintf("Did not return correct response %v", v))
	}
}

func MockQuayBuildSuccessHook() QuayBuildSuccessHook {
	body := QuayBuildSuccessHook{}
	body.Repository = "mynamespace/repository"
	body.Namespace = "mynamespace"
	body.Name = "repository"
	body.DockerURL = "quay.io/mynamespace/repository"
	body.Homepage = "https://quay.io/repository/mynamespace/repository/build?current=some-fake-build"
	body.Visibility = "public"
	body.BuildId = "dba8aa95-4de5-4a09-8d87-1527eaa4856a"
	body.BuildName = "some-fake-build"
	body.DockerTags = []string{"latest", "foo", "bar"}
	body.TriggerKind = "github"

	return body
}

func MockQuayBuildSuccessHookJSON() []byte {
	b, _ := json.Marshal(MockQuayBuildSuccessHook())
	return b
}
