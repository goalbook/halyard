package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/coreos/go-etcd/etcd"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestQuayBuildSuccessPostHandler(t *testing.T) {
	config := &Config{}
	err := readConfig(config)
	if err != nil {
		log.Panicf("Error occurred in config: %v", err)
	}

	etcdClient = etcd.NewClient([]string{config.EtcdURL})

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

	mock := MockQuayBuildSuccessHook()
	result, err := etcdClient.Get(fmt.Sprintf("/containers/%s/latest/build", mock.DockerURL), false, false)

	if err != nil {
		t.Error(fmt.Sprintf("Error getting from etcd: %v", err))
	}
	if result.Node.Value != mock.BuildId {
		t.Error(fmt.Sprintf("Did not write etcd: %v", result.Node.Value))
	}
}

func MockQuayBuildSuccessHook() QuayBuildSuccessHook {
	body := QuayBuildSuccessHook{}
	body.Repository = "mynamespace/repository"
	body.Namespace = "mynamespace"
	body.Name = "repository"
	body.DockerURL = "quay.io/mynamespace/repository"
	body.DockerTags = []string{"master", "latest"}
	body.Homepage = "https://quay.io/repository/mynamespace/repository/build?current=some-fake-build"
	body.Visibility = "public"
	body.ImageId = "c2cbdf995d089eaa5c33c9ebf37dd1e61311503f30f530edbf9c7f2f6c2be441"
	body.BuildId = "dba8aa95-4de5-4a09-8d87-1527eaa4856a"
	body.BuildName = "some-fake-build"
	body.TriggerId = "8e42ea6b-8883-42a1-b199-75cdb68ac3ec"
	body.TriggerKind = "github"

	return body
}

func MockQuayBuildSuccessHookJSON() []byte {
	b, _ := json.Marshal(MockQuayBuildSuccessHook())
	return b
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

//{"build_id": "fake-build-id", "trigger_kind": "GitHub", "name": "halyard", "repository": "goalbook/halyard", "namespace": "goalbook", "docker_url": "quay.io/goalbook/halyard", "visibility": "public", "docker_tags": ["latest", "foo", "bar"], "build_name": "some-fake-build", "image_id": "1245657346", "trigger_metadata": {"default_branch": "master", "ref": "refs/heads/somebranch", "commit_sha": "42d4a62c53350993ea41069e9f2cfdefb0df097d"}, "homepage": "https://quay.io/repository/goalbook/halyard/build?current=fake-build-id"}

// {
//     "build_id": "6da6f88a-4f13-49f8-abe3-0ab71ea26092",
//     "trigger_kind": "github",
//     "name": "halyard",
//     "repository": "goalbook/halyard",
//     "namespace": "goalbook",
//     "docker_url": "quay.io/goalbook/halyard",
//     "visibility": "public",
//     "docker_tags": ["master", "latest"],
//     "build_name": "9978c03",
//     "image_id": "c2cbdf995d089eaa5c33c9ebf37dd1e61311503f30f530edbf9c7f2f6c2be441",
//     "trigger_metadata": {
//         "default_branch": "master",
//         "ref": "refs/heads/master",
//         "commit_sha": "9978c03c1351cb24ca86a00e15894a6f3cd9af1d",
//         "commit_info": {
//             "url": "https://github.com/goalbook/halyard/commit/9978c03c1351cb24ca86a00e15894a6f3cd9af1d",
//             "date": "Sat, 28 Mar 2015 20:13:05 GMT",
//             "message": "only check for app/json not utf for quay",
//             "committer": {
//                 "username":
//                 "danieljyoo",
//                 "url": "https://github.com/danieljyoo",
//                 "avatar_url": "https://avatars.githubusercontent.com/u/17211?v=3"
//             }, "author": {
//                 "username":
//                 "danieljyoo",
//                 "url": "https://github.com/danieljyoo",
//                 "avatar_url": "https://avatars.githubusercontent.com/u/17211?v=3"
//             }
//         }
//     },
//     "trigger_id": "8e42ea6b-8883-42a1-b199-75cdb68ac3ec",
//     "homepage": "https://quay.io/repository/goalbook/halyard/build?current=6da6f88a-4f13-49f8-abe3-0ab71ea26092"
// }
