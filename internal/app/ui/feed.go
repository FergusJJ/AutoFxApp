package ui

import "github.com/rivo/tview"

type Feed struct {
	FeedView *tview.Flex
}

func newFeed() *Feed {
	testBox := tview.NewFlex()

	return &Feed{FeedView: testBox}
}

func (feed *Feed) pushToFeed() {
	// feed.FeedView = feed.FeedView.AddItem("List item 1", "Some explanatory text", 'a', nil).
	// 	AddItem("List item 2", "Some explanatory text", 'b', nil).
	// 	AddItem("List item 3", "Some explanatory text", 'c', nil).
	// 	AddItem("List item 4", "Some explanatory text", 'd', nil)
}
