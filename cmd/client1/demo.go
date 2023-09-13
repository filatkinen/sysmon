// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"github.com/GoryMoon/gocui"
	"github.com/filatkinen/sysmon/internal/model"
	"log"
)

type ClientView struct {
	g      *gocui.Gui
	params []string
	data   []model.DataToClientStamp
}

func NewClientView() (*ClientView, error) {
	g, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		return nil, err
	}

	g.Highlight = true
	g.Cursor = true
	g.SelFgColor = gocui.ColorGreen

	return &ClientView{g: g}, nil
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

func (c *ClientView) Stop() error {
	c.Close()
	return nil
}

func (c *ClientView) Close() {
	c.g.Close()
}

func (c *ClientView) nextView(g *gocui.Gui, v *gocui.View) error {
	if v == nil || v.Name() == "main" {
		_, err := g.SetCurrentView("side")
		return err
	}
	_, err := g.SetCurrentView("main")
	return err
}

func (c *ClientView) GetData(m []model.DataToClientStamp) error {
	c.params = nil
	for i := range m {
		c.params = append(c.params, m[i].Name)
	}

}

func (c *ClientView) cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *ClientView) cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *ClientView) quit(g *gocui.Gui, v *gocui.View) error {
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
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, c.quit); err != nil {
		return err
	}
	return nil
}

func (c *ClientView) layout(g *gocui.Gui) error {
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

	if v, err := g.SetView("side", 0, 5, 30, 4+4, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = "Params"
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		fmt.Fprintln(v, "Item 1")
		fmt.Fprintln(v, "Item 2")
		fmt.Fprintln(v, "Item 3")
		if _, err := g.SetCurrentView("side"); err != nil {
			return err
		}

	}
	if v, err := g.SetView("main", 30, 0, maxX-1, maxY-1, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		for i := 1; i < maxY+20; i++ {
			fmt.Fprintf(v, "%d\n", i)
		}
		v.Title = "counters"
		v.Editable = true
		//v.Wrap = true
		//v.Autoscroll =

	}
	return nil
}

func main() {

	c, err := NewClientView()
	if err != nil {
		log.Println(err)
		return
	}
	defer c.Close()
	err = c.Start()
	if err != nil {
		log.Println(err)
	}
}
