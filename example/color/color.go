package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/gosuri/uiprogress"
	"github.com/gosuri/uiprogress/util/strutil"
)

var steps = []string{
	color.CyanString("downloading source"),
	color.CyanString("installing deps"),
	color.CyanString("compiling"),
	color.CyanString("packaging"),
	color.CyanString("seeding database"),
	color.CyanString("deploying"),
	color.CyanString("staring servers"),
}

// FIXME: Breaks on Windows - keeps appending more lines
func main() {
	fmt.Fprintf(color.Output, "apps: deployment started: %s\n", color.MagentaString("app1, app2"))
	p := uiprogress.New()
	p.SetOut(color.Output)
	p.Start()

	var wg sync.WaitGroup
	wg.Add(1)
	go deploy(p, color.MagentaString("app1"), &wg)
	wg.Add(1)
	go deploy(p, color.MagentaString("app2"), &wg)
	wg.Wait()

	fmt.Fprintf(color.Output, "apps: successfully deployed: %s\n", color.MagentaString("app1, app2"))
}

func deploy(p *uiprogress.Progress, app string, wg *sync.WaitGroup) {
	defer wg.Done()
	bar := p.AddBar(len(steps)).AppendCompleted().PrependElapsed()
	bar.Width = 50

	// prepend the deploy step to the bar
	bar.PrependFunc(func(b *uiprogress.Bar) string {
		return strutil.Resize(app+": "+steps[b.Current()-1], 22)
	})

	rand.Seed(500)
	for bar.Incr() {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(2000)))
	}
}
