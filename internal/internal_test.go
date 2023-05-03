package internal_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ikorchynskyi/terraform-version-inspect/internal"
)

var ReleaseListResponse = []string{
	`[
		{
			"timestamp_created": "2023-04-26T18:26:07.947Z",
			"version": "1.4.6"
		},
		{
			"timestamp_created": "2023-04-12T18:20:00.211Z",
			"version": "1.4.5"
		},
		{
			"timestamp_created": "2023-04-05T17:00:04.540Z",
			"version": "1.5.0-alpha20230405"
		},
		{
			"timestamp_created": "2023-03-30T19:52:25.048Z",
			"version": "1.4.4"
		},
		{
			"timestamp_created": "2023-03-30T16:12:15.000Z",
			"version": "1.4.3"
		},
		{
			"timestamp_created": "2023-03-16T13:42:34.377Z",
			"version": "1.4.2"
		},
		{
			"timestamp_created": "2023-03-15T14:51:34.470Z",
			"version": "1.4.1"
		},
		{
			"timestamp_created": "2023-03-08T18:13:21.723Z",
			"version": "1.4.0"
		},
		{
			"timestamp_created": "2023-02-24T19:19:59.594Z",
			"version": "1.4.0-rc1"
		},
		{
			"timestamp_created": "2023-02-15T17:30:44.822Z",
			"version": "1.3.9"
		},
		{
			"timestamp_created": "2023-02-15T15:34:21.761Z",
			"version": "1.4.0-beta2"
		},
		{
			"timestamp_created": "2023-02-09T19:50:21.840Z",
			"version": "1.3.8"
		},
		{
			"timestamp_created": "2023-02-09T16:59:39.118Z",
			"version": "1.4.0-beta1"
		},
		{
			"timestamp_created": "2023-01-04T15:34:05.718Z",
			"version": "1.3.7"
		},
		{
			"timestamp_created": "2022-12-07T15:08:26.965Z",
			"version": "1.4.0-alpha20221207"
		},
		{
			"timestamp_created": "2022-11-30T20:56:54.970Z",
			"version": "1.3.6"
		},
		{
			"timestamp_created": "2022-11-17T20:01:35.078Z",
			"version": "1.3.5"
		},
		{
			"timestamp_created": "2022-11-09T18:27:45.160Z",
			"version": "1.4.0-alpha20221109"
		},
		{
			"timestamp_created": "2022-11-02T16:31:11.692Z",
			"version": "1.3.4"
		},
		{
			"timestamp_created": "2022-10-19T17:54:15.486Z",
			"version": "1.3.3"
		}
	]`,
	`[
		{
			"timestamp_created": "2022-10-06T16:57:24.231Z",
			"version": "1.3.2"
		},
		{
			"timestamp_created": "2022-09-28T14:00:09.590Z",
			"version": "1.3.1"
		},
		{
			"timestamp_created": "2022-09-21T13:58:58.183Z",
			"version": "1.3.0"
		},
		{
			"timestamp_created": "2022-09-14T18:00:23.849Z",
			"version": "1.3.0-rc1"
		},
		{
			"timestamp_created": "2022-09-07T21:03:28.128Z",
			"version": "1.2.9"
		},
		{
			"timestamp_created": "2022-08-31T13:18:11.524Z",
			"version": "1.3.0-beta1"
		},
		{
			"timestamp_created": "2022-08-24T14:34:02.743Z",
			"version": "1.2.8"
		},
		{
			"timestamp_created": "2022-08-17T15:37:24.720Z",
			"version": "1.3.0-alpha20220817"
		},
		{
			"timestamp_created": "2022-08-10T17:57:06.837Z",
			"version": "1.2.7"
		},
		{
			"timestamp_created": "2022-08-03T17:25:52.543Z",
			"version": "1.3.0-alpha20220803"
		},
		{
			"timestamp_created": "2022-07-27T15:34:10.807Z",
			"version": "1.2.6"
		},
		{
			"timestamp_created": "2022-07-13T10:34:39.031Z",
			"version": "1.2.5"
		},
		{
			"timestamp_created": "2022-07-06T18:26:28.440Z",
			"version": "1.3.0-alpha20220706"
		},
		{
			"timestamp_created": "2022-06-29T17:59:34.438Z",
			"version": "1.2.4"
		},
		{
			"timestamp_created": "2022-06-22T19:34:37.343Z",
			"version": "1.3.0-alpha20220622"
		},
		{
			"timestamp_created": "2022-06-15T17:56:17.787Z",
			"version": "1.2.3"
		},
		{
			"timestamp_created": "2022-06-08T17:30:57.209Z",
			"version": "1.3.0-alpha20220608"
		},
		{
			"timestamp_created": "2022-06-01T16:58:55.636Z",
			"version": "1.2.2"
		},
		{
			"timestamp_created": "2022-05-23T22:49:08.229Z",
			"version": "1.2.1"
		},
		{
			"timestamp_created": "2022-05-18T21:47:46.272Z",
			"version": "1.2.0"
		}
	]`,
	`[
		{
			"timestamp_created": "2022-05-11T18:57:37.675Z",
			"version": "1.2.0-rc2"
		},
		{
			"timestamp_created": "2022-05-04T17:00:45.000Z",
			"version": "1.2.0-rc1"
		},
		{
			"timestamp_created": "2022-04-27T19:29:14.000Z",
			"version": "1.2.0-beta1"
		},
		{
			"timestamp_created": "2022-04-20T13:46:15.000Z",
			"version": "1.1.9"
		},
		{
			"timestamp_created": "2022-04-13T18:30:29.000Z",
			"version": "1.2.0-alpha20220413"
		},
		{
			"timestamp_created": "2022-04-07T17:05:13.000Z",
			"version": "1.1.8"
		},
		{
			"timestamp_created": "2022-03-28T10:27:08.000Z",
			"version": "1.2.0-alpha-20220328"
		},
		{
			"timestamp_created": "2022-03-02T19:32:54.000Z",
			"version": "1.1.7"
		},
		{
			"timestamp_created": "2022-02-16T18:44:06.000Z",
			"version": "1.1.6"
		},
		{
			"timestamp_created": "2022-02-02T20:46:14.000Z",
			"version": "1.1.5"
		}
	]`,
}

var ReleaseListAfter = []string{
	"",
	"2022-10-19T17:54:15.486Z",
	"2022-05-18T21:47:46.272Z",
	"2022-02-02T20:46:14.000Z",
}

func TestGetReleases(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if after := r.URL.Query().Get("after"); after != "" {
			if ReleaseListAfter[0] == "" {
				t.Errorf("Expected empty after time, got '%s'", after)
			}
			afterTime, err := time.Parse(time.RFC3339Nano, after)
			if err != nil {
				t.Errorf("Expected correct after time, got '%s'", after)
			}
			expectedTime, _ := time.Parse(time.RFC3339Nano, ReleaseListAfter[0])
			fmt.Printf("%v, %v, %v", afterTime, expectedTime, err)
			if !afterTime.Equal(expectedTime) {
				t.Errorf("Expected '%s' after time, got '%s'", afterTime.Format(time.RFC3339Nano), expectedTime.Format(time.RFC3339Nano))
			}
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)

		if len(ReleaseListResponse) > 0 {
			w.Write([]byte(ReleaseListResponse[0]))
			ReleaseListResponse = ReleaseListResponse[1:]
			ReleaseListAfter = ReleaseListAfter[1:]
		} else {
			w.Write([]byte(`[]`))
		}
	}))
	defer server.Close()

	releases, _ := internal.GetReleases(server.URL)
	expected := []string{
		"1.4.6",
		"1.4.5",
		"1.5.0-alpha20230405",
		"1.4.4",
		"1.4.3",
		"1.4.2",
		"1.4.1",
		"1.4.0",
		"1.4.0-rc1",
		"1.3.9",
		"1.4.0-beta2",
		"1.3.8",
		"1.4.0-beta1",
		"1.3.7",
		"1.4.0-alpha20221207",
		"1.3.6",
		"1.3.5",
		"1.4.0-alpha20221109",
		"1.3.4",
		"1.3.3",
		"1.3.2",
		"1.3.1",
		"1.3.0",
		"1.3.0-rc1",
		"1.2.9",
		"1.3.0-beta1",
		"1.2.8",
		"1.3.0-alpha20220817",
		"1.2.7",
		"1.3.0-alpha20220803",
		"1.2.6",
		"1.2.5",
		"1.3.0-alpha20220706",
		"1.2.4",
		"1.3.0-alpha20220622",
		"1.2.3",
		"1.3.0-alpha20220608",
		"1.2.2",
		"1.2.1",
		"1.2.0",
		"1.2.0-rc2",
		"1.2.0-rc1",
		"1.2.0-beta1",
		"1.1.9",
		"1.2.0-alpha20220413",
		"1.1.8",
		"1.2.0-alpha-20220328",
		"1.1.7",
		"1.1.6",
		"1.1.5",
	}

	if len(expected) != len(releases) {
		t.Fatalf("Expected %d releases, got %d", len(expected), len(releases))
	}

	for i, v := range releases {
		if expected[i] != v.Version {
			t.Errorf("Expected release version '%s', got '%s'", expected[i], v.Version)
		}
	}
}
