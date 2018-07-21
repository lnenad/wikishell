package main

import (
	"fmt"
	"os"
	"strings"
	"wikiquestion/wikiapi"

	"net/url"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func main() {
	var pageName string

	argsWithoutProg := os.Args[1:]

	wapi := wikiapi.NewWikiAPI(&url.URL{
		Scheme: "https",
		Host:   "en.wikipedia.org",
		Path:   "/w/api.php",
	})

	app := tview.NewApplication()
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			app.Draw()
		})

	textView.SetBorder(true)

	pages := tview.NewPages()

	sectionSelectionModal := func(p tview.Primitive, width, height int) tview.Primitive {
		return tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(p, height, 1, false).
				AddItem(nil, 0, 1, false), width, 1, false).
			AddItem(nil, 0, 1, false)
	}

	sectionSelectionForm := tview.NewList().AddItem("Default", "", 'a', nil)

	pageSelectionModal := func(p tview.Primitive, width, height int) tview.Primitive {
		return tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(p, height, 1, false).
				AddItem(nil, 0, 1, false), width, 1, false).
			AddItem(nil, 0, 1, false)
	}

	pageSelectionForm := tview.NewForm().
		AddButton("Load", func() {
			textView.Clear()
			LoadArticle(pageName, wapi, textView, pages, sectionSelectionForm)
			app.SetFocus(textView)
		}).
		AddInputField("Page name", "", 20, nil, func(text string) {
			pageName = text
		})

	pages.AddPage("sectionSelection", sectionSelectionModal(sectionSelectionForm, 40, 10), true, true)
	pages.AddPage("textView", textView, true, true)

	pages.AddPage("pageSelection", pageSelectionModal(pageSelectionForm, 40, 10), true, true)

	if len(argsWithoutProg) > 0 {
		go func() { // Async loading to allow the initial render and expose the real dimensions of the window
			LoadArticle(argsWithoutProg[0], wapi, textView, pages, sectionSelectionForm)
			app.SetFocus(textView)
		}()
		pages.SwitchToPage("textView")
	} else {
		pages.SwitchToPage("pageSelection")
	}

	if err := app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlL {
			pages.SwitchToPage("pageSelection")
			app.SetFocus(pageSelectionForm)
		}
		if event.Key() == tcell.KeyCtrlS {
			fmt.Println("MRSMRSMRS")
			pages.SwitchToPage("sectionSelection")
			app.SetFocus(sectionSelectionForm)
		}

		return event
	}).SetRoot(pages, true).SetFocus(pageSelectionForm).Run(); err != nil {
		panic(err)
	}
}

func LoadArticle(pageName string, wapi *wikiapi.WikiAPI, textView *tview.TextView, pages *tview.Pages, sectionsForm *tview.List) *wikiapi.ParsedWikiResponse {
	response := wikiapi.FetchParsedText(pageName, wapi)
	page := response.Parsed.Text.Content
	_, _, w, h := pages.Box.GetInnerRect()
	contents := page.StripHTML().SplitContentByLines(w, h)
	fmt.Fprintf(textView, "%s", contents)
	pages.SwitchToPage("textView")

	sections := response.Parsed.Sections
	sectionsForm.Clear()

	for _, section := range sections {
		func(section wikiapi.Section) {
			sectionsForm.AddItem(section.Title, "", 'a', func() {
				index := FindRowIndex(&section, contents)
				textView.ScrollTo(index, 0)
				pages.SwitchToPage("textView")
			})
		}(section)
	}

	return response
}

func FindRowIndex(section *wikiapi.Section, words []string) int {
	index := -1
	target := section.Title + "[edit]"
	for _, word := range words {
		index++

		match := func(word string, index int, section *wikiapi.Section) bool {
			if strings.Trim(word, " \n") == target {
				return true
			}

			return false
		}(word, index, section)

		if match {
			return index
		}
	}

	return -1
}
