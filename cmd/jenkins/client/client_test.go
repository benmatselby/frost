package client_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/benmatselby/frost/cmd/jenkins/client"
	jenkins "github.com/benmatselby/frost/cmd/jenkins/client"
)

const (
	baseURLPath     = "/testing"
	getJobsURL      = "/api/json"
	getJobsResponse = `{
		"views": [
			{
				"name": "View One",
				"jobs": [
					{
						"_class": "hudson.model.FreeStyleProject",
						"name": "frost-one",
						"fullDisplayName": "frost",
						"color": "blue",
						"lastBuild": {
							"_class": "hudson.model.FreeStyleBuild",
							"number": 377,
							"result": "SUCCESS",
							"timestamp": 1536547140829
						}
					}
				]
			},
			{
				"name": "View Two",
				"jobs": [
					{
						"_class": "hudson.model.FreeStyleProject",
						"name": "frost-two",
						"fullDisplayName": "frost",
						"color": "blue",
						"lastBuild": {
							"_class": "hudson.model.FreeStyleBuild",
							"number": 377,
							"result": "SUCCESS",
							"timestamp": 1536547140829
						}
					},
					{
						"_class": "org.jenkinsci.plugins.workflow.job.WorkflowJob",
						"name": "frost-three",
						"fullDisplayName": "frost > master",
						"color": "blue",
						"lastBuild": {
							"_class": "org.jenkinsci.plugins.workflow.job.WorkflowRun",
							"number": 9,
							"result": "SUCCESS",
							"timestamp": 1536124741734
						}
					}
				]
			}
		]
	}`
)

// Pulled from https://github.com/google/go-github/blob/master/github/github_test.go
func setup() (client *client.Client, mux *http.ServeMux, serverURL string, teardown func()) {
	// mux is the HTTP request multiplexer used with the test server.
	mux = http.NewServeMux()

	// We want to ensure that tests catch mistakes where the endpoint URL is
	// specified as absolute rather than relative. It only makes a difference
	// when there's a non-empty base URL path. So, use that. See issue #752.
	apiHandler := http.NewServeMux()
	apiHandler.Handle(baseURLPath+"/", http.StripPrefix(baseURLPath, mux))
	apiHandler.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(os.Stderr, "FAIL: Client.Baseurl path prefix is not preserved in the request URL:")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "\t"+req.URL.String())
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "\tDid you accidentally use an absolute endpoint URL rather than relative?")
		http.Error(w, "Client.BaseURL path prefix is not preserved in the request URL.", http.StatusInternalServerError)
	})

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(apiHandler)

	client = jenkins.New(server.URL+"/"+baseURLPath, "USERNAME", "TOKEN")

	return client, mux, server.URL, server.Close
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

func TestJenkins_New(t *testing.T) {
	URI := "http://localhost"
	username := "jenkins"
	password := "secret"

	client := jenkins.New(URI, username, password)

	if client.Baseurl != URI {
		t.Fatalf("expected base URL %s; got %v", URI, client.Baseurl)
	}
}

func TestJenkins_GetJobs(t *testing.T) {
	tt := []struct {
		testName string
		URL      string
		viewName string
		response string
		err      string
	}{
		{testName: "can handle a bad response", URL: getJobsURL, viewName: "RANDON_VIEW", response: "{badresponse}", err: "decoding json response for jenkins views from api/json?tree=views%5Bname%2Cjobs%5Bname%2CfullDisplayName%2ClastBuild%5Bnumber%2Ctimestamp%2Cresult%5D%5D%5D&depth=1 failed: invalid character 'b' looking for beginning of object key string"},
		{testName: "cannot find the view", URL: getJobsURL, viewName: "RANDON_VIEW", response: getJobsResponse, err: "unable to get named jenkins view: RANDON_VIEW"},
		{testName: "we find the view we want and return it", URL: getJobsURL, viewName: "View One", response: getJobsResponse, err: ""},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			c, mux, _, teardown := setup()
			defer teardown()

			mux.HandleFunc(tc.URL, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				json := tc.response
				fmt.Fprint(w, json)
			})

			jobs, err := c.GetJobs(tc.viewName)

			if tc.err == "" {
				if len(jobs) > 1 {
					t.Fatalf("expected 1 job, got '%v'", len(jobs))
				}

				if jobs[0].Name != "frost-one" {
					t.Fatalf("expected the Name to be 'frost'; got '%v'", jobs[0].Name)
				}
			} else {
				if err.Error() != tc.err {
					t.Fatalf("expected error '%s'; got '%v'", tc.err, err.Error())
				}
			}
		})
	}
}
