// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package backend

import (
	"code.google.com/p/log4go"
	"fmt"
	. "github.com/quarnster/util/text"
	"io/ioutil"
)

type (
	WatchedFile interface {
		Name() string
		Reload()
	}

	WatchedUserFile struct {
		view *View
	}

	WatchedSettingFile struct {
		path string
	}
)

func NewWatchedUserFile(view *View) *WatchedUserFile {
	return &WatchedUserFile{view}
}

func (o WatchedUserFile) String() string {
	return fmt.Sprintf("%s (%d)", o.Name(), o.view.Id())
}

func (o *WatchedUserFile) Name() string {
	return o.view.Buffer().FileName()
}

func (o *WatchedUserFile) Reload() {
	log4go.Finest("\"%v\".Reload()", o)

	view := o.view
	filename := o.Name()

	if saving, ok := view.Settings().Get("lime.saving", false).(bool); ok && saving {
		// This reload was triggered by ourselves saving to this file, so don't reload it
		return
	}
	if !GetEditor().Frontend().OkCancelDialog("File was changed by another program, reload?", "reload") {
		return
	}

	if d, err := ioutil.ReadFile(filename); err != nil {
		log4go.Error("Could not read file: %s\n. Error was: %v", filename, err)
	} else {
		edit := view.BeginEdit()
		end := view.Buffer().Size()
		view.Replace(edit, Region{0, end}, string(d))
		view.EndEdit(edit)
	}
}

func NewWatchedSettingFile(path string) *WatchedSettingFile {
	return &WatchedSettingFile{path}
}

func (o *WatchedSettingFile) Name() string {
	return o.path
}

func (o *WatchedSettingFile) Reload() {
	editor := GetEditor()
	editor.loadSetting(o.path)
}
