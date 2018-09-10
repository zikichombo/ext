package oggvorbis

import (
	"os"
	"testing"

	"zikichombo.org/codec"
	"zikichombo.org/sio"
)

func TestRegister(t *testing.T) {
	file, err := os.Open("path/to/some/vorbis/file.ogg")
	if err != nil {
		t.Fatal(err)
	}
	source, _, err := codec.Decoder(file, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = sio.Play(source)
	if err != nil {
		t.Fatal(err)
	}
}
