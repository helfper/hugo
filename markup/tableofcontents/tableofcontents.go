// Copyright 2019 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tableofcontents

import (
	"strings"
)

// Header holds the data about a header and its children.
type Header struct {
	ID   string
	Text string

	Headers []Header
}

// IsZero is true when no ID or Text is set.
func (h Header) IsZero() bool {
	return h.ID == "" && h.Text == ""
}

// Root holds the top level (h1) headers.
type Root struct {
	Headers []Header
}

// AddAt adds the header into the given level, starting at 1.
func (toc *Root) AddAt(h Header, level int) {
	headers := &toc.Headers
	for i := 1; i < level; i++ {
		if len(*headers) == 0 {
			*headers = append(*headers, Header{})
		}
		headers = &(*headers)[len(*headers)-1].Headers
	}
	*headers = append(*headers, h)
}

// ToHTML renders the ToC as HTML.
func (toc Root) ToHTML(startLevel, stopLevel int, ordered bool) string {
	b := &tocBuilder{
		s:          strings.Builder{},
		root:       toc,
		startLevel: startLevel,
		stopLevel:  stopLevel,
		tag:        orderedTag[ordered],
	}
	b.Build()
	return b.s.String()
}

type tocBuilder struct {
	s    strings.Builder
	root Root

	startLevel int
	stopLevel  int
	tag        string
}

func (b *tocBuilder) Build() {
	b.writeNav(b.root.Headers)
}

func (b *tocBuilder) writeNav(hs []Header) {
	b.s.WriteString("<nav id=\"TableOfContents\">")
	b.writeHeaders(1, 0, hs)
	b.s.WriteString("</nav>")
}

func (b *tocBuilder) writeHeaders(level, indent int, hs []Header) {
	if level < b.startLevel {
		for _, h := range hs {
			b.writeHeaders(level+1, indent, h.Headers)
		}
		return
	}

	if b.stopLevel != -1 && level > b.stopLevel {
		return
	}

	hasChildren := len(hs) > 0

	if hasChildren {
		b.s.WriteString("\n")
		b.indent(indent + 1)
		b.s.WriteString("<" + b.tag + ">\n")
	}

	for _, h := range hs {
		b.writeHeader(level+1, indent+2, h)
	}

	if hasChildren {
		b.indent(indent + 1)
		b.s.WriteString("</" + b.tag + ">")
		b.s.WriteString("\n")
		b.indent(indent)
	}

}
func (b *tocBuilder) writeHeader(level, indent int, h Header) {
	b.indent(indent)
	b.s.WriteString("<li>")
	if !h.IsZero() {
		b.s.WriteString("<a href=\"#" + h.ID + "\">" + h.Text + "</a>")
	}
	b.writeHeaders(level, indent, h.Headers)
	b.s.WriteString("</li>\n")
}

func (b *tocBuilder) indent(n int) {
	for i := 0; i < n; i++ {
		b.s.WriteString("  ")
	}
}

var orderedTag = map[bool]string{true: "ol", false: "ul"}

// DefaultConfig is the default ToC configuration.
var DefaultConfig = Config{
	StartLevel: 2,
	EndLevel:   3,
	Ordered:    false,
}

// Config holds ToC configuration.
type Config struct {
	// Heading start level to include in the table of contents, starting
	// at h1 (inclusive).
	StartLevel int

	// Heading end level, inclusive, to include in the table of contents.
	// Default is 3, a value of -1 will include everything.
	EndLevel int

	// Whether to produce a ordered list or not.
	Ordered bool
}
