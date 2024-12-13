package bitmap

type Bitmap struct {
	bits []byte
	size int
}

func NewBitmap(size int) *Bitmap {

	if size == 0 {
		size = 250 // 默认值
	}
	return &Bitmap{
		bits: make([]byte, size),
		size: size * 8, // 是bit
	}
}

func (b *Bitmap) Set(id string) {
	// 计算id的哈希值
	idx := hash(id) % b.size
	// 计算这个位置在哪个字节
	byteIdx := idx / 8
	// 在这个字节中的bit位置
	bitIdx := idx % 8

	// 先左移，再设置
	b.bits[byteIdx] |= 1 << bitIdx
}

func (b *Bitmap) IsSet(id string) bool {
	// 计算id的哈希值
	idx := hash(id) % b.size
	// 计算这个位置在哪个字节
	byteIdx := idx / 8
	// 在这个字节中的bit位置
	bitIdx := idx % 8

	// 判断1
	return (b.bits[byteIdx] & (1 << bitIdx)) != 0
}

// Export 导出和加载都是为了存储在数据库中做准备
func (b *Bitmap) Export() []byte {
	return b.bits
}

func (b *Bitmap) Load(bits []byte) *Bitmap {
	if len(bits) == 0 {
		return NewBitmap(0)
	}

	return &Bitmap{
		bits: bits,
		size: len(bits) * 8,
	}
}

func hash(id string) int {
	// 使用 BKDR 哈希算法
	seed := 131313
	hash := 0
	for _, c := range id {
		hash = hash*seed + int(c)
	}
	return hash & 0x7FFFFFFF
}
