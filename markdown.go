package shorthand

import (
	// 3rd Party packages
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

var (
	// Go Markdown setup.
	extensions = parser.NoIntraEmphasis |
		parser.Tables |
		parser.FencedCode |
		parser.Autolink |
		parser.Strikethrough |
		parser.SpaceHeadings
	htmlFlags = html.CommonFlags | html.HrefTargetBlank
	opts      = html.RendererOptions{Flags: htmlFlags}
	renderer  = html.NewRenderer(opts)
)

// MarkdownToHTML wraps the Gomarkdown markdown library
// to allow changing out the markdown processing in the future
// by swaping out the contents of this function.
func MarkdownToHTML(src []byte) []byte {
	parser := parser.NewWithExtensions(extensions)
	return markdown.ToHTML(src, parser, renderer)
}
