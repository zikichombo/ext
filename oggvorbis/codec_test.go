package oggvorbis

import (
	"testing"

	"zikichombo.org/codec"
)

func TestRegister(t *testing.T) {
	err := codec.RegisterCodec(Codec)
	if err != nil {
		t.Error(err)
	}

	/*
		file, err := os.Open("path/to/some/vorbis/file.ogg")
		if err != nil {
			panic(err)
		}
		source, _, err := codec.Decoder(file, nil)
		if err != nil {
			panic(err)
		}
		err = sio.Play(source)
		if err != nil {
			panic(err)
		}
	*/
}
