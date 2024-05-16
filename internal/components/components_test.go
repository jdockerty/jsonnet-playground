package components

import (
	"context"
	"fmt"
	"io"
	"strings"
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
		_ = RootPage("").Render(context.Background(), w)
		_ = w.Close()
	}()
	doc, err := goquery.NewDocumentFromReader(r)
	assert.Nil(err)

	textarea := doc.Find("#jsonnet-input").Get(0)
	foundPlaceholder := false
	foundAutofocus := false
	foundOnKeyDown := false
	foundHtmxGet := false
	for _, attr := range textarea.Attr {

		if attr.Key == "autofocus" {
			foundAutofocus = true
		}

		if attr.Key == "onkeydown" && strings.Contains(attr.Val, "allowTabs") {
			foundOnKeyDown = true
		}

		if attr.Val == "Type your Jsonnet here..." {
			foundPlaceholder = true
		}
	}

	assert.True(foundPlaceholder, "Placeholder text not found")
	assert.True(foundAutofocus, "autofocus should be on the textarea")
	assert.True(foundOnKeyDown, "onkeydown should contain the allowTabs script")
	assert.False(foundHtmxGet, "No hx-get should be present")

	assert.Equal("Jsonnet Playground", doc.Find("h1").Text())
}

func TestRootPageWithShare(t *testing.T) {
	assert := assert.New(t)
	fakePath := "fake"

	r, w := io.Pipe()
	go func() {
		_ = RootPage(fakePath).Render(context.Background(), w)
		_ = w.Close()
	}()
	doc, err := goquery.NewDocumentFromReader(r)
	assert.Nil(err)

	textarea := doc.Find("#jsonnet-input").Get(0)
	foundHtmxGet := false
	foundHtmxTrigger := false
	for _, attr := range textarea.Attr {

		if attr.Key == "hx-get" && attr.Val == fmt.Sprintf("/api/share/%s", fakePath) {
			foundHtmxGet = true
		}

		if attr.Key == "hx-trigger" && attr.Val == "load" {
			foundHtmxTrigger = true
		}
	}

	assert.True(foundHtmxGet, "hx-get attribute should exist with the /api/share/%s path", fakePath)
	assert.True(foundHtmxTrigger, "hx-trigger attribute should exist")

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
