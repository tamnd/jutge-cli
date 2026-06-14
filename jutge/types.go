package jutge

// Problem is a Jutge.org programming problem.
type Problem struct {
	Rank   int    `json:"rank"   csv:"rank"   tsv:"rank"`
	Code   string `json:"code"   csv:"code"   tsv:"code"`
	Title  string `json:"title"  csv:"title"  tsv:"title"`
	URL    string `json:"url"    csv:"url"    tsv:"url"`
}
