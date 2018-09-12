// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

// +build listen

package flac

import (
	"os"
	"testing"

	"zikichombo.org/codec"
	"zikichombo.org/sio"
	"zikichombo.org/sound"
)

// Assert that flac.Decoder implements the sound.Source interface.
var _ sound.Source = (*Decoder)(nil)

func TestRegister(t *testing.T) {
	co, err := codec.CodecFor(".flac", nil)
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.Open("0.flac")
	if err != nil {
		t.Fatal(err)
	}
	dec, _, err := co.Decoder(f)
	if err != nil {
		t.Fatal(err)
	}
	defer dec.Close()
	sio.Play(dec)
}
