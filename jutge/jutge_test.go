package jutge_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tamnd/jutge-cli/jutge"
)

const fakeProblemsHTML = `
<html><body>
<table>
<tr><td>&nbsp;<a href='/problems/P10051_en'>P10051</a>&nbsp; Worst path <span>pdf</span></td></tr>
<tr><td>&nbsp;<a href='/problems/P10052_en'>P10052</a>&nbsp; Best path <span>pdf</span></td></tr>
<tr><td>&nbsp;<a href='/problems/X10001_en'>X10001</a>&nbsp; Sum of digits <span>pdf</span></td></tr>
</table>
</body></html>
`

func newTestClient(ts *httptest.Server) *jutge.Client {
	cfg := jutge.DefaultConfig()
	cfg.BaseURL = ts.URL
	cfg.Rate = 0
	return jutge.NewClient(cfg)
}

func TestListProblems(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fakeProblemsHTML))
	}))
	defer ts.Close()

	c := newTestClient(ts)
	problems, err := c.ListProblems(context.Background(), 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(problems) != 3 {
		t.Fatalf("want 3 problems, got %d", len(problems))
	}
	if problems[0].Code != "P10051" {
		t.Errorf("first code = %q, want P10051", problems[0].Code)
	}
	if problems[0].Title != "Worst path" {
		t.Errorf("first title = %q, want 'Worst path'", problems[0].Title)
	}
}

func TestListProblemsLimit(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fakeProblemsHTML))
	}))
	defer ts.Close()

	c := newTestClient(ts)
	problems, err := c.ListProblems(context.Background(), 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(problems) != 2 {
		t.Fatalf("want 2 problems with limit, got %d", len(problems))
	}
}

func TestSearch(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fakeProblemsHTML))
	}))
	defer ts.Close()

	c := newTestClient(ts)
	results, err := c.Search(context.Background(), "path", 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 {
		t.Fatalf("want 2 results for 'path', got %d", len(results))
	}
}
