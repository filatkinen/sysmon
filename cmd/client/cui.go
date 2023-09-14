package main

import (
	"errors"
	"fmt"
	"sync"

	"github.com/GoryMoon/gocui"
	"github.com/fatih/color"
	"github.com/filatkinen/sysmon/internal/client"
	"github.com/filatkinen/sysmon/internal/model"
)

type ClientView struct {
	g            *gocui.Gui
	params       []string
	data         []model.DataToClientStamp
	lock         sync.Mutex
	idx          int
	isLoadedData bool
}

func NewClientView() (*ClientView, error) {
	g, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		return nil, err
	}

	g.Highlight = true
	g.Cursor = true
	g.SelFgColor = gocui.ColorGreen

	return &ClientView{g: g, lock: sync.Mutex{}, params: append([]string(nil), "No data yet...")}, nil
}

func (c *ClientView) Start() error {
	c.g.SetManagerFunc(c.layout)

	if err := c.keybindings(c.g); err != nil {
		return err
	}
	if err := c.g.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		return err
	}
	return nil
}

func (c *ClientView) Stop() {
	c.Close()
}

func (c *ClientView) Close() {
	c.g.Close()
}

func (c *ClientView) GetData(m []model.DataToClientStamp) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if !c.isLoadedData {
		c.params = nil
		if len(m) == 0 {
			c.params = append(c.params, "No data...")
			c.isLoadedData = false
		} else {
			for i := range m {
				c.params = append(c.params, m[i].Name)
			}
			c.isLoadedData = true
		}
		_ = c.g.DeleteView("side")
		_ = c.printParams()
	}
	c.data = m
	_ = c.printData()
	c.g.Update(func(gui *gocui.Gui) error {
		return nil
	})
}

func (c *ClientView) printData() error {
	out, err := c.g.View("main")
	if err != nil {
		return err
	}
	out.Clear()
	if len(c.data) > 0 && c.idx < len(c.data) {
		d := c.data[c.idx].Data
		headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
		columnFmt := color.New(color.FgYellow).SprintfFunc()
		tbl := client.NewTable(d[0]...)
		tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
		tbl.WithWriter(out)
		for i := range d {
			if i == 0 {
				continue
			}
			tbl.AddRow(d[i]...)
		}
		tbl.Print()
	}
	return nil
}

func (c *ClientView) printParams() error {
	if v, err := c.g.SetView("side", 0, 5, 30, len(c.params)+6, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = "Params"
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		for _, val := range c.params {
			fmt.Fprintln(v, val)
		}
		if _, err := c.g.SetCurrentView("side"); err != nil {
			return err
		}
	}
	return nil
}

func (c *ClientView) nextView(g *gocui.Gui, v *gocui.View) error {
	if v == nil || v.Name() == "main" {
		_, err := g.SetCurrentView("side")
		return err
	}
	_, err := g.SetCurrentView("main")
	return err
}

func (c *ClientView) cursorDown(_ *gocui.Gui, v *gocui.View) error {
	if v != nil { //nolint:nestif
		cx, cy := v.Cursor()
		if cy < len(c.params)-1 {
			err := v.SetCursor(cx, cy+1)
			if err != nil {
				ox, oy := v.Origin()
				if err := v.SetOrigin(ox, oy+1); err != nil {
					return err
				}
			}
			c.idx = cy + 1
			_ = c.printData()
		}
	}
	return nil
}

func (c *ClientView) cursorUp(_ *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		err := v.SetCursor(cx, cy-1)
		if err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
		if cy > 0 {
			c.idx = cy - 1
		}
		_ = c.printData()
	}
	return nil
}

func (c *ClientView) quit(_ *gocui.Gui, _ *gocui.View) error {
	return gocui.ErrQuit
}

func (c *ClientView) keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("side", gocui.KeyTab, gocui.ModNone, c.nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyTab, gocui.ModNone, c.nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("side", gocui.KeyArrowDown, gocui.ModNone, c.cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("side", gocui.KeyArrowUp, gocui.ModNone, c.cursorUp); err != nil {
		return err
	}
	err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, c.quit)
	return err
}

func (c *ClientView) layout(g *gocui.Gui) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	maxX, maxY := g.Size()
	g.Cursor = false
	if v, err := g.SetView("keybinding", 0, 0, 30, 4, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = "KeyBindins"
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		fmt.Fprintln(v, "TAB -switch")
		fmt.Fprintln(v, "Arrow Up/Down:choose param")
		fmt.Fprintln(v, "Contol-C :exit")
	}

	err := c.printParams()
	if err != nil {
		return err
	}
	if v, err := g.SetView("main", 30, 0, maxX-1, maxY-1, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = "Counters"
		v.Editable = true
	}
	return nil
}
