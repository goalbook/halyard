package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/nu7hatch/gouuid"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var config *Config

type Config struct {
	ServerPort         string
	EtcdURL            string
	SecurityKey        string
	GithubAPIToken     string
	DockerRegistryAuth string
}

func readConfig(c *Config) error {
	c.ServerPort = os.Getenv("SERVER_PORT")
	if c.ServerPort == "" {
		return errors.New("Server Port not provided ($SERVER_PORT)")
	}
	c.EtcdURL = os.Getenv("ETCD_URL")
	if c.EtcdURL == "" {
		return errors.New("etcd URL not provided ($ETCD_URL)")
	}
	c.SecurityKey = os.Getenv("SECURITY_KEY")
	if c.SecurityKey == "" {
		return errors.New("Security Key not provided ($SECURITY_KEY)")
	}
	c.GithubAPIToken = os.Getenv("GITHUB_API_TOKEN")
	if c.GithubAPIToken == "" {
		return errors.New("Github API Token not provied ($GITHUB_API_TOKEN)")
	}
	c.DockerRegistryAuth = os.Getenv("DOCKER_REGISTRY_AUTH")
	if c.DockerRegistryAuth == "" {
		return errors.New("Docker Registry Auth not provided ($DOCKER_REGISTRY_AUTH)")
	}
	return nil
}

// {
//   'repository': 'mynamespace/repository',
//   'namespace': 'mynamespace',
//   'name': 'repository',
//   'docker_url': 'quay.io/mynamespace/repository',
//   'homepage': 'https://quay.io/repository/mynamespace/repository/build?current=some-fake-build',
//   'visibility': 'public',

//   'build_id': build_uuid,
//   'build_name': 'some-fake-build',
//   'docker_tags': ['latest', 'foo', 'bar'],
//   'trigger_kind': 'github'

// }

type QuayBuildSuccessHook struct {
	Repository  string   `json:"repository"`
	Namespace   string   `json:"namespace"`
	Name        string   `json:"name"`
	DockerURL   string   `json:"docker_url"`
	Homepage    string   `json:"homepage"`
	Visibility  string   `json:"visibility"`
	BuildId     string   `json:"build_id"`
	BuildName   string   `json:"build_name"`
	DockerTags  []string `json:"docker_tags"`
	TriggerKind string   `json:"trigger_kind"`
}

func main() {
	config := &Config{}
	err := readConfig(config)
	if err != nil {
		log.Panicf("Error occurred in config: %v", err)
	}

	r := mux.NewRouter()

	r.Path("/healthcheck").Methods("GET").HandlerFunc(HealthCheckGetHandler)

	s := r.PathPrefix(fmt.Sprintf("/%s", config.SecurityKey)).Subrouter()

	s.Path("/success").Methods("POST").HandlerFunc(QuayBuildSuccessPostHandler)

	http.Handle("/", r)
	http.ListenAndServe(":"+config.ServerPort, handlers.LoggingHandler(os.Stdout, http.DefaultServeMux))
}

func HealthCheckGetHandler(w http.ResponseWriter, r *http.Request) {
	res := map[string]string{
		"status": "healthy",
	}
	WriteResponseJSON(w, 200, res)
}

// We create a photo from robohash and save it s3
func QuayBuildSuccessPostHandler(w http.ResponseWriter, r *http.Request) {
	body := QuayBuildSuccessHook{}
	_ = DecodeJSONBody(r, &body)

	res := map[string]string{}
	WriteResponseJSON(w, 200, res)
}

func UUID() string {
	uuid, _ := uuid.NewV4()
	return uuid.String()
}

func DecodeJSONBody(r *http.Request, v interface{}) error {
	if r.Method == "POST" || r.Method == "PUT" {
		if strings.ToLower(r.Header.Get("Content-Type")) == "application/json;charset=utf-8" {
			bytes, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return err
			}
			return DecodeJSON(bytes, v)
		}
	}
	return nil
}

func DecodeJSON(body []byte, v interface{}) error {
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(v); err != nil {
		return err
	}
	return nil
}

func EncodeJSONBody(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func WriteResponse(w http.ResponseWriter, status int, body []byte) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	w.Write(body)
	w.Write([]byte("\n"))
	return
}

func WriteResponseJSON(w http.ResponseWriter, status int, v interface{}) {
	body, _ := json.Marshal(v)
	WriteResponse(w, status, body)
	return
}
