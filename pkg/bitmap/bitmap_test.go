package bitmap

import (
	"testing"
)

func TestBitmap_Set(t *testing.T) {
	b := NewBitmap(20)

	b.Set("pppp")
	b.Set("222")
	b.Set("pppp")
	b.Set("ccc")
	b.Set("eee")
	b.Set("xxx")
	b.Set("fff")

	for _, bit := range b.bits {
		t.Logf("%b, %v", bit, bit)
	}
}
