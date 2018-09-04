// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package flac

import "zikichombo.org/sound"

// Assert that flac.Decoder implements the sound.Source interface.
var _ sound.Source = (*Decoder)(nil)
