package main

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
	"log"
	"time"
)

func main() {
	{
		// initialize a controllable Chrome instance
		ctx, cancel := chromedp.NewContext(
			context.Background(),
		)

		// to release the browser resources when
		// it is no longer needed
		defer cancel()

		var html string
		err := chromedp.Run(ctx,
			// visit the target page
			chromedp.Navigate("https://www.naver.com/"),
			// wait for the page to load
			chromedp.Sleep(2000*time.Millisecond),
			// extract the raw HTML from the page
			chromedp.ActionFunc(func(ctx context.Context) error {
				// select the root node on the page
				rootNode, err := dom.GetDocument().Do(ctx)
				if err != nil {
					return err
				}
				html, err = dom.GetOuterHTML().WithNodeID(rootNode.NodeID).Do(ctx)
				return err
			}),
		)
		if err != nil {
			log.Fatal("Error while performing the automation logic:", err)
		}

		fmt.Println(html)
	}

	{
		opts := append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", false),
		)
		ctx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
		ctx, cancel := chromedp.NewContext(
			ctx,
			//chromedp.WithDebugf(log.Printf),
		)
		defer cancel()

		var res string
		url := "https://toscrape.com"
		err := chromedp.Run(ctx, chromedp.Navigate(url), chromedp.WaitVisible("body"), chromedp.Sleep(4000*time.Millisecond), chromedp.TextContent(`body`, &res))

		if err != nil {
			fmt.Println("Error navigating to the website:", err)
			return
		}

		fmt.Println("Website loaded successfully:", res)
	}
}
