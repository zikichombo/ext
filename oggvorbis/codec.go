// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package oggvorbis

import (
	"bufio"
	"io"

	"github.com/jfreymuth/oggvorbis"
	"zikichombo.org/codec"
	"zikichombo.org/sound"
	"zikichombo.org/sound/cil"
	"zikichombo.org/sound/freq"
	"zikichombo.org/sound/sample"
)

type oggvorbisCodec struct{ codec.NullCodec }

func (*oggvorbisCodec) Extensions() []string { return []string{".ogg"} }

func (*oggvorbisCodec) Sniff(r *bufio.Reader) bool {
	if peek, err := r.Peek(4); err == nil {
		return string(peek) == "OggS"
	}
	return false
}

func (*oggvorbisCodec) Decoder(r io.ReadCloser) (sound.Source, sample.Codec, error) {
	return NewDecoder(r)
}
func (*oggvorbisCodec) SeekingDecoder(r codec.IoReadSeekCloser) (sound.SourceSeeker, sample.Codec, error) {
	return NewSeekingDecoder(r)
}

func init() {
	codec.RegisterCodec(&oggvorbisCodec{})
}

type Decoder struct {
	dec    *oggvorbis.Reader
	closer io.Closer
	buf    []float32
	pos    int
	err    error
}

func NewDecoder(r io.ReadCloser) (sound.Source, sample.Codec, error) {
	oggr, err := oggvorbis.NewReader(r)
	d := &Decoder{}
	d.dec = oggr
	d.closer = r
	d.buf = make([]float32, 0, 8192)
	return d, codec.AnySampleCodec, err
}

func (d *Decoder) Receive(out []float64) (int, error) {
	channels := d.Channels()
	frames := len(out) / channels
	if len(out)%channels != 0 {
		return 0, sound.ErrChannelAlignment
	}
	framesRead := 0
	for d.err == nil && framesRead < frames {
		if d.pos == len(d.buf) {
			n, err := d.dec.Read(d.buf[:cap(d.buf)])
			d.buf = d.buf[:n]
			d.pos = 0
			d.err = err
		}
		n := (len(d.buf) - d.pos) / channels
		if n > frames-framesRead {
			n = frames - framesRead
		}
		in := d.buf[d.pos:]
		for ch := 0; ch < channels; ch++ {
			for i := 0; i < n; i++ {
				out[ch*frames+framesRead+i] = float64(in[i*channels+ch])
			}
		}
		d.pos += n * channels
		framesRead += n
	}
	if framesRead < frames {
		if err := cil.Compact(out, channels, framesRead); err != nil {
			return 0, err
		}
	}
	return framesRead, d.err
}

func (d *Decoder) Close() error {
	return d.closer.Close()
}

func (d *Decoder) Channels() int {
	return d.dec.Channels()
}

func (d *Decoder) SampleRate() freq.T {
	return freq.T(d.dec.SampleRate()) * freq.Hertz
}

type SeekingDecoder struct {
	Decoder
}

func NewSeekingDecoder(r codec.IoReadSeekCloser) (sound.SourceSeeker, sample.Codec, error) {
	oggr, err := oggvorbis.NewReader(r)
	d := &SeekingDecoder{}
	d.dec = oggr
	d.closer = r
	d.buf = make([]float32, 0, 8192)
	return d, codec.AnySampleCodec, err
}

func (d *SeekingDecoder) Pos() int64 {
	return d.dec.Position()
}

func (d *SeekingDecoder) Len() int64 {
	return d.dec.Length()
}

func (d *SeekingDecoder) Seek(f int64) error {
	return d.dec.SetPosition(f)
}
