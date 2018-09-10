// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

// TODO: add seek support; only enabled in the dev branch of mewkiz/flac.

// Note, ZikiChombo uses the terminology frame to refer to a single sample per
// channel. Thus, given an audio source with 100 samples and 2 channels, there
// exist 50 frames. This is not to be confused with a FLAC frame, which is a
// container for audio samples; where each FLAC frame contains one subframe per
// channel, and each subframe contains the audio samples of a given channel.

package flac

import (
	"io"

	"github.com/mewkiz/flac"
	"github.com/mewkiz/flac/frame"
	"zikichombo.org/sound"
	"zikichombo.org/sound/cil"
	"zikichombo.org/sound/freq"
	"zikichombo.org/sound/sample"
)

// Decoder encapsulates state for decoding and seeking a FLAC audio stream.
type Decoder struct {
	stream *flac.Stream
	frame  *frame.Frame // current frame.
	i      int          // index of current sample in subframe(s).
}

// NewDecoder creates a decoder from a FLAC audio stream (seekable, readable).
func NewDecoder(r io.Reader) (*Decoder, error) {
	stream, err := flac.New(r)
	if err != nil {
		return nil, err
	}
	d := &Decoder{
		stream: stream,
	}
	return d, nil
}

// Receive places samples in dst.
//
// Receive returns the number of frames placed in dst together with
// any error.  Receive may use all of dst as scatch space.
//
// Receive returns a non-nil error if and only if it returns 0 frames
// received.  Receive may return 0 < n < len(dst)/Channels() frames only
// if the subsequent call will return (0, io.EOF).  As a result, the
// caller need not be concerned with whether or not the data is "ready".
//
// Receive returns multi-channel data in de-interleaved format.
// If len(dst) is not a multiple of Channels(), then Receive returns
// ErrChannelAlignment.  If Receive returns fewer than len(dst)/Channels()
// frames, then the deinterleaved data of n frames is arranged in
// the prefix dst[:n*Channels()].
func (d *Decoder) Receive(dst []float64) (int, error) {
	nC := d.Channels()
	if len(dst)%nC != 0 {
		return 0, sound.ErrChannelAlignment
	}
	// number of frames (samples per channel) read.
	n := 0
	// number of frames (samples per channel) to read.
	nF := len(dst) / nC
	bps := int(d.stream.Info.BitsPerSample)
	for n < nF {
		if d.frame == nil {
			frame, err := d.stream.ParseNext()
			if err != nil {
				// Compact channel-interleaved samples if n < nF.
				if err := cil.Compact(dst, nC, n); err != nil {
					return n, err
				}
				return n, err
			}
			d.frame = frame
			d.i = 0
		}
		samplesLeft := len(d.frame.Subframes[0].Samples[d.i:])
		j := nF
		if j > samplesLeft {
			j = samplesLeft
		}
		if n+j > nF {
			j = nF - n
		}
		for c := 0; c < nC; c++ {
			sample.Int32sToFloats(dst[c*nF+n:c*nF+n+j], d.frame.Subframes[c].Samples[d.i:d.i+j], bps)
		}
		d.i += j
		n += j
		if len(d.frame.Subframes[0].Samples[d.i:]) == 0 {
			d.frame = nil
			d.i = 0
		}
	}
	return n, nil
}

// Channels returns the number of channels of the FLAC stream.
func (d *Decoder) Channels() int {
	return int(d.stream.Info.NChannels)
}

// SampleRate returns the sample rate of the FLAC stream.
func (d *Decoder) SampleRate() freq.T {
	return freq.T(d.stream.Info.SampleRate) * freq.Hertz
}

// Close closes the underlying FLAC stream of the decoder.
func (d *Decoder) Close() error {
	return d.stream.Close()
}
