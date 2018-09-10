// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package flac

import (
	"io"

	"zikichombo.org/codec"
	"zikichombo.org/sound"
	"zikichombo.org/sound/sample"
)

type flacCodec struct {
	codec.NullCodec
}

func (f *flacCodec) Extensions() []string {
	return []string{".flc", ".flac"}
}

func (f *flacCodec) Decoder(rc io.ReadCloser) (sound.Source, sample.Codec, error) {
	dec, err := NewDecoder(rc)
	return dec, codec.AnySampleCodec, err
}

func init() {
	codec.RegisterCodec(&flacCodec{})
}
