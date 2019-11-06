package test

import (
	"net/http"
	"testing"
)

// TestDownload validates the http Get function can download content.
// 表组测试，推荐使用
func TestDownload2(t *testing.T) {
	var tests = []struct {
		url        string
		statusCode int
	}{
		{
			"http://www.goinggo.net/feeds/posts/default?alt=rss",
			http.StatusOK,
		},
		{
			"http://rss.cnn.com/rss/cnn_topstbadurl.rss",
			http.StatusNotFound,
		},
	}

	t.Log("Given the need to test downloading content.")
	{
		for _, test := range tests {

			t.Logf("\tWhen checking \"%s\" for status code \"%d\"",
				test.url, test.statusCode)
			{
				resp, err := http.Get(test.url)
				if err != nil {
					t.Fatal("\t\tShould be able to make the Get call.",
						ballotX, err)
				}
				t.Log("\t\tShould be able to make the Get call.",
					checkMark)

				defer resp.Body.Close()

				if resp.StatusCode == test.statusCode {
					t.Logf("\t\tShould receive a \"%d\" status. %v",
						test.statusCode, checkMark)
				} else {
					t.Errorf("\t\tShould receive a \"%d\" status. %v %v",
						test.statusCode, ballotX, resp.StatusCode)
				}
			}
		}

	}
}
