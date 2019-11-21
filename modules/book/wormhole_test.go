package book

import "testing"

func TestNewWormhole(t *testing.T) {
	wor := NewWormhole()

	wor.CaptureTags()
}
