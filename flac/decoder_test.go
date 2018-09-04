package flac

import "zikichombo.org/sound"

// Assert that flac.Decoder implements the sound.Source interface.
var _ sound.Source = (*Decoder)(nil)
