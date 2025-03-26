package mend

import (
	"fmt"
	"io"
	"path/filepath"
	"slices"
	"strings"

	"github.com/bbfh-dev/mend.html/mend/attrs"
	"github.com/bbfh-dev/mend.html/mend/settings"
	"github.com/bbfh-dev/mend.html/mend/tags"
	"golang.org/x/net/html"
)

const PREFIX = "mend:"
const PREFIX_LEN = len(PREFIX)

type Template struct {
	Name   string
	Params string

	Root tags.NodeWithChildren

	currentLine  int
	currentToken html.Token
	currentText  string
	currentAttrs attrs.Attributes
	// A list of current parents from greatest to closest
	breadcrumbs []tags.NodeWithChildren
}

func NewTemplate(name string, params string) *Template {
	root := tags.NewRootNode()
	return &Template{
		Name:         name,
		Params:       params,
		Root:         root,
		currentLine:  1,
		currentToken: html.Token{},
		currentAttrs: attrs.Attributes{},
		breadcrumbs:  []tags.NodeWithChildren{root},
	}
}

func (template *Template) lastBreadcrumb() tags.NodeWithChildren {
	return template.breadcrumbs[len(template.breadcrumbs)-1]
}

func (template *Template) append(nodes ...tags.Node) {
	template.lastBreadcrumb().Add(nodes...)
}

func (template *Template) appendLevel(node tags.NodeWithChildren) {
	template.append(node)
	template.breadcrumbs = append(template.breadcrumbs, node)
}

func (template *Template) Parse(reader io.Reader) error {
	tokenizer := html.NewTokenizer(reader)

loop:
	for {
		tokenType := tokenizer.Next()
		switch tokenType {

		case html.ErrorToken:
			if tokenizer.Err() == io.EOF {
				break loop
			}
			return fmt.Errorf(
				"(%s) %w",
				filepath.Base(template.Name),
				tokenizer.Err(),
			)

		case html.TextToken:
			template.currentLine += strings.Count(template.currentToken.Data, "\n")

		}

		template.currentToken = tokenizer.Token()
		template.currentAttrs = attrs.New(template.currentToken.Attr)
		template.currentText = strings.TrimSpace(template.currentToken.Data)

		err := template.Process(tokenType)
		if err != nil {
			return fmt.Errorf(
				"(%s:%d) %w",
				filepath.Base(template.Name),
				template.currentLine,
				err,
			)
		}
	}

	return nil
}

func (template *Template) Process(tokenType html.TokenType) error {
	switch tokenType {

	case html.DoctypeToken:
		template.append(tags.NewDoctypeNode(template.currentText))

	case html.CommentToken:
		if settings.KeepComments {
			template.append(tags.NewCommentNode(template.currentText))
		}

	case html.TextToken:
		if len(template.currentText) == 0 {
			break
		}

		var builder strings.Builder
		for _, line := range strings.Split(template.currentText, "\n") {
			builder.WriteString(strings.TrimSpace(line))
			builder.WriteString(" ")
		}
		template.append(tags.NewTextNode(builder.String()))

	case html.SelfClosingTagToken:
		if !strings.HasPrefix(template.currentText, PREFIX) {
			node := tags.NewVoidNode(template.currentText, template.currentAttrs)
			template.append(node)
			break
		}

	case html.StartTagToken:
		if !strings.HasPrefix(template.currentText, PREFIX) {
			// Is it actually a self-closing tag with wrong syntax?
			if slices.Contains(attrs.SelfClosingTags, template.currentText) {
				return template.Process(html.SelfClosingTagToken)
			}

			node := tags.NewTagNode(template.currentText, template.currentAttrs)
			template.appendLevel(node)
			break
		}

	case html.EndTagToken:
		if len(template.breadcrumbs) == 1 {
			break
		}
		template.breadcrumbs = template.breadcrumbs[:len(template.breadcrumbs)-1]
	}

	return nil
}
