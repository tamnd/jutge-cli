// Package jutge scrapes the Jutge.org problem archive.
package jutge

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// Config controls the HTTP client behaviour.
type Config struct {
	BaseURL   string
	Rate      time.Duration
	Timeout   time.Duration
	Retries   int
	UserAgent string
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		BaseURL:   "https://jutge.org",
		Rate:      500 * time.Millisecond,
		Timeout:   60 * time.Second,
		Retries:   3,
		UserAgent: "jutge-cli/0.1 (github.com/tamnd/jutge-cli)",
	}
}

// Client fetches Jutge data.
type Client struct {
	cfg  Config
	http *http.Client
	last time.Time
}

// NewClient creates a Client from cfg.
func NewClient(cfg Config) *Client {
	return &Client{
		cfg:  cfg,
		http: &http.Client{Timeout: cfg.Timeout},
	}
}

var problemRe = regexp.MustCompile(`href='/problems/([A-Z]\d+)_en'>[A-Z]\d+</a>\s*&nbsp;\s*([^\n<]+)`)

// ListProblems fetches all English problems from the Jutge problem list page.
func (c *Client) ListProblems(ctx context.Context, limit int) ([]Problem, error) {
	body, err := c.fetch(ctx, "/problems")
	if err != nil {
		return nil, err
	}
	matches := problemRe.FindAllStringSubmatch(string(body), -1)
	var problems []Problem
	for i, m := range matches {
		if limit > 0 && i >= limit {
			break
		}
		code := m[1]
		title := strings.TrimSpace(m[2])
		problems = append(problems, Problem{
			Rank:  i + 1,
			Code:  code,
			Title: title,
			URL:   fmt.Sprintf("%s/problems/%s_en", c.cfg.BaseURL, code),
		})
	}
	return problems, nil
}

// Search returns problems whose title or code contains query (case-insensitive).
func (c *Client) Search(ctx context.Context, query string, limit int) ([]Problem, error) {
	all, err := c.ListProblems(ctx, 0)
	if err != nil {
		return nil, err
	}
	q := strings.ToLower(query)
	var out []Problem
	for _, p := range all {
		if strings.Contains(strings.ToLower(p.Title), q) || strings.Contains(strings.ToLower(p.Code), q) {
			out = append(out, p)
			if limit > 0 && len(out) >= limit {
				break
			}
		}
	}
	for i := range out {
		out[i].Rank = i + 1
	}
	return out, nil
}

func (c *Client) fetch(ctx context.Context, path string) ([]byte, error) {
	url := c.cfg.BaseURL + path
	var last error
	for attempt := 0; attempt <= c.cfg.Retries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(time.Duration(attempt) * time.Second):
			}
		}
		c.pace()
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("User-Agent", c.cfg.UserAgent)
		resp, err := c.http.Do(req)
		if err != nil {
			last = err
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 {
			last = fmt.Errorf("HTTP %d", resp.StatusCode)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("HTTP %d for %s", resp.StatusCode, url)
		}
		return io.ReadAll(resp.Body)
	}
	return nil, fmt.Errorf("all retries failed for %s: %w", url, last)
}

func (c *Client) pace() {
	if c.cfg.Rate <= 0 {
		return
	}
	if wait := c.cfg.Rate - time.Since(c.last); wait > 0 {
		time.Sleep(wait)
	}
	c.last = time.Now()
}
