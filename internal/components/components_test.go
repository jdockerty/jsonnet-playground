package components

import (
	"context"
	"io"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
)

func TestHeader(t *testing.T) {
	assert := assert.New(t)

	r, w := io.Pipe()
	go func() {
		_ = heading().Render(context.Background(), w)
		_ = w.Close()
	}()
	doc, err := goquery.NewDocumentFromReader(r)
	assert.Nil(err)

	assert.Equal("Jsonnet Playground", doc.Find("title").Text())
}

func TestRootPage(t *testing.T) {
	assert := assert.New(t)

	r, w := io.Pipe()
	go func() {
		_ = RootPage().Render(context.Background(), w)
		_ = w.Close()
	}()
	doc, err := goquery.NewDocumentFromReader(r)
	assert.Nil(err)

	assert.Equal("Jsonnet Playground", doc.Find("h1").Text())
}

func TestJsonnetDisplay(t *testing.T) {
	assert := assert.New(t)

	r, w := io.Pipe()
	go func() {
		_ = jsonnetDisplay("").Render(context.Background(), w)
		_ = w.Close()
	}()
	doc, err := goquery.NewDocumentFromReader(r)
	assert.Nil(err)

	assert.Equal(2, len(doc.Find("textarea").Nodes))
}
