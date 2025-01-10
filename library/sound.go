package library

import (
	"bytes"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

// SoundByPath returns a sound by path based on embedded assets
func SoundByPath(path string) ([]byte, error) {
	data, err := embeddedAssets.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	stream, err := mp3.DecodeWithSampleRate(44100, bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(stream)
	if err != nil {
		return nil, fmt.Errorf("read from stream: %w", err)
	}

	return buf.Bytes(), nil
}
