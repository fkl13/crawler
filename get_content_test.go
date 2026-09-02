package main

import (
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func TestGetHeadingFromHTMLBasic(t *testing.T) {
	inputBody := "<html><body><h1>Test Title</h1></body></html>"
	actual := getHeadingFromHTML(inputBody)
	expected := "Test Title"

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestGetFirstParagraphFromHTMLMainPriority(t *testing.T) {
	inputBody := `<html><body>
		<p>Outside paragraph.</p>
		<main>
			<p>Main paragraph.</p>
		</main>
	</body></html>`
	actual := getFirstParagraphFromHTML(inputBody)
	expected := "Main paragraph."

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestGetFirstParagraphFromHTMLFallback(t *testing.T) {
	inputBody := `<html><body>
		<p>First paragraph outside main.</p>
		<p>Second paragraph outside main.</p>
	</body></html>`
	actual := getFirstParagraphFromHTML(inputBody)
	expected := "First paragraph outside main."

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestGetHeadingFromHTMLNoHeading(t *testing.T) {
	inputBody := `<html><body><p></p></body></html>`
	actual := getHeadingFromHTML(inputBody)
	expected := ""

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestGetHeadingFromHTMLH2Fallback(t *testing.T) {
	inputBody := "<html><body><h2>Fallback Title</h2></body></html>"
	actual := getHeadingFromHTML(inputBody)
	expected := "Fallback Title"

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestGetURLsFromHTML(t *testing.T) {
	cases := []struct {
		name          string
		inputURL      string
		inputBody     string
		expected      []string
		errorContains string
	}{
		{
			name:     "absolute URL",
			inputURL: "https://crawler-test.com",
			inputBody: `
<html>
	<body>
		<a href="https://crawler-test.com">
			<span>Boot.dev</span>
		</a>
	</body>
</html>
`,
			expected: []string{"https://crawler-test.com"},
		},
		{
			name:     "relative URL",
			inputURL: "https://crawler-test.com",
			inputBody: `
<html>
	<body>
		<a href="/path/one">
			<span>Boot.dev</span>
		</a>
	</body>
</html>
`,
			expected: []string{"https://crawler-test.com/path/one"},
		},
		{
			name:     "absolute and relative URLs",
			inputURL: "https://crawler-test.com",
			inputBody: `
<html>
	<body>
		<a href="/path/one">
			<span>Boot.dev</span>
		</a>
		<a href="https://other.com/path/one">
			<span>Boot.dev</span>
		</a>
	</body>
</html>
`,
			expected: []string{"https://crawler-test.com/path/one", "https://other.com/path/one"},
		},
		{
			name:     "no href",
			inputURL: "https://crawler-test.com",
			inputBody: `
<html>
	<body>
		<a>
			<span>Boot.dev</span>
		</a>
	</body>
</html>
`,
			expected: nil,
		},
		{
			name:     "bad HTML",
			inputURL: "https://crawler-test.com",
			inputBody: `
<html body>
	<a href="path/one">
		<span>Boot.dev</span>
	</a>
</html body>
`,
			expected: []string{"https://crawler-test.com/path/one"},
		},
		{
			name:     "invalid href URL",
			inputURL: "https://crawler-test.com",
			inputBody: `
<html>
	<body>
		<a href=":\\invalidURL">
			<span>Boot.dev</span>
		</a>
	</body>
</html>
`,
			expected: nil,
		},
	}

	for i, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			baseURL, err := url.Parse(tc.inputURL)
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: couldn't parse input URL: %v", i, tc.name, err)
				return
			}

			actual, err := getURLsFromHTML(tc.inputBody, baseURL)

			if err != nil && !strings.Contains(err.Error(), tc.errorContains) {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			} else if err != nil && tc.errorContains == "" {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			} else if err == nil && tc.errorContains != "" {
				t.Errorf("Test %v - '%s' FAIL: expected error containing '%v', got none.", i, tc.name, tc.errorContains)
				return
			}

			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Test %v - '%s' FAIL: expected URLs %v, got URLs %v", i, tc.name, tc.expected, actual)
				return
			}
		})
	}
}

func TestGetImagesFromHTMLAbsolute(t *testing.T) {
	inputURL := "https://crawler-test.com"
	inputBody := `<html><body><img src="https://crawler-test.com/logo.png" alt="Logo"></body></html>`

	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	actual, err := getImagesFromHTML(inputBody, parsedURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"https://crawler-test.com/logo.png"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestGetImagesFromHTMLRelative(t *testing.T) {
	inputURL := "https://crawler-test.com"
	inputBody := `<html><body><img src="/logo.png" alt="Logo"></body></html>`

	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	actual, err := getImagesFromHTML(inputBody, parsedURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"https://crawler-test.com/logo.png"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestGetImagesFromHTMLMultiple(t *testing.T) {
	inputURL := "https://crawler-test.com"
	inputBody := `<html><body>
		<img src="/logo.png" alt="Logo">
		<img src="https://cdn.boot.dev/banner.jpg">
	</body></html>`

	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	actual, err := getImagesFromHTML(inputBody, parsedURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{
		"https://crawler-test.com/logo.png",
		"https://cdn.boot.dev/banner.jpg",
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestExtractPageData(t *testing.T) {
	cases := []struct {
		name    string
		pageURL string
		html    string
		want    PageData
	}{
		{
			name:    "basic: h1, main paragraph, relative link and img",
			pageURL: "https://crawler-test.com",
			html: `
<html>
  <body>
    <h1>Hello World</h1>
    <main><p>First paragraph inside main.</p></main>
    <a href="/about">About</a>
    <img src="/logo.png" alt="Logo">
  </body>
</html>`,
			want: PageData{
				URL:            "https://crawler-test.com",
				Heading:        "Hello World",
				FirstParagraph: "First paragraph inside main.",
				OutgoingLinks:  []string{"https://crawler-test.com/about"},
				ImageURLs:      []string{"https://crawler-test.com/logo.png"},
			},
		},
		{
			name:    "fallback paragraph when no <main>",
			pageURL: "https://crawler-test.com",
			html: `
<html>
  <body>
    <h1>Title</h1>
    <p>Outside paragraph wins.</p>
    <a href="/x">x</a>
    <img src="/img.png">
  </body>
</html>`,
			want: PageData{
				URL:            "https://crawler-test.com",
				Heading:        "Title",
				FirstParagraph: "Outside paragraph wins.",
				OutgoingLinks:  []string{"https://crawler-test.com/x"},
				ImageURLs:      []string{"https://crawler-test.com/img.png"},
			},
		},
		{
			name:    "malformed HTML still parsed; absolute link and image",
			pageURL: "https://crawler-test.com",
			html: `
<html body>
  <h1>Messy</h1>
  <a href="https://other.com/path">Other</a>
  <img src="https://cdn.boot.dev/banner.jpg">
</html body>`,
			want: PageData{
				URL:            "https://crawler-test.com",
				Heading:        "Messy",
				FirstParagraph: "", // no <p> present
				OutgoingLinks:  []string{"https://other.com/path"},
				ImageURLs:      []string{"https://cdn.boot.dev/banner.jpg"},
			},
		},
		{
			name:    "no h1 and no paragraph",
			pageURL: "https://crawler-test.com",
			html: `
<html>
  <body>
    <a href="/only-link">Only link</a>
    <img src="/only.png">
  </body>
</html>`,
			want: PageData{
				URL:            "https://crawler-test.com",
				Heading:        "",
				FirstParagraph: "",
				OutgoingLinks:  []string{"https://crawler-test.com/only-link"},
				ImageURLs:      []string{"https://crawler-test.com/only.png"},
			},
		},
		{
			name:    "multiple links and images preserve order",
			pageURL: "https://crawler-test.com",
			html: `
<html><body>
  <h1>t</h1>
  <main><p>p</p></main>
  <a href="/a1">a1</a>
  <a href="https://x.dev/a2">a2</a>
  <img src="/i1.png">
  <img src="https://x.dev/i2.png">
</body></html>`,
			want: PageData{
				URL:            "https://crawler-test.com",
				Heading:        "t",
				FirstParagraph: "p",
				OutgoingLinks: []string{
					"https://crawler-test.com/a1",
					"https://x.dev/a2",
				},
				ImageURLs: []string{
					"https://crawler-test.com/i1.png",
					"https://x.dev/i2.png",
				},
			},
		},
		{
			name:    "invalid base URL → empty link/image slices",
			pageURL: `:\\invalidBaseURL`,
			html: `
<html>
  <body>
    <h1>Title</h1>
    <p>Paragraph</p>
    <a href="/path">path</a>
    <img src="/logo.png">
  </body>
</html>`,
			want: PageData{
				URL:            `:\\invalidBaseURL`,
				Heading:        "Title",
				FirstParagraph: "Paragraph",
				OutgoingLinks:  nil,
				ImageURLs:      nil,
			},
		},
	}

	for _, tc := range cases {
		tc := tc // shadow the loop variable.
		t.Run(tc.name, func(t *testing.T) {
			got := extractPageData(tc.html, tc.pageURL)

			if got.URL != tc.want.URL {
				t.Errorf("URL: want %q, got %q", tc.want.URL, got.URL)
			}
			if got.Heading != tc.want.Heading {
				t.Errorf("Heading: want %q, got %q", tc.want.Heading, got.Heading)
			}
			if got.FirstParagraph != tc.want.FirstParagraph {
				t.Errorf("FirstParagraph: want %q, got %q", tc.want.FirstParagraph, got.FirstParagraph)
			}
			if !reflect.DeepEqual(got.OutgoingLinks, tc.want.OutgoingLinks) {
				t.Errorf("OutgoingLinks: want %v, got %v", tc.want.OutgoingLinks, got.OutgoingLinks)
			}
			if !reflect.DeepEqual(got.ImageURLs, tc.want.ImageURLs) {
				t.Errorf("ImageURLs: want %v, got %v", tc.want.ImageURLs, got.ImageURLs)
			}
		})
	}
}
