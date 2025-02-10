package bloom_filter

import (
	"hash"
	"hash/fnv"
	"math"
)

type BloomFilter struct {
	bitset      []bool        // 位数组
	size        uint64        // 位数组大小
	hashFuncs   []hash.Hash64 // 哈希函数集合
	numHashFunc int           // 哈希函数数量
}

// 创建布隆过滤器
func NewBloomFilter(expectedElements int, falsePositiveRate float64) *BloomFilter {
	// 计算位数组大小和哈希函数数量
	m, k := calculateParameters(expectedElements, falsePositiveRate)
	return &BloomFilter{
		bitset:      make([]bool, m),
		size:        m,
		hashFuncs:   initHashFunctions(k),
		numHashFunc: k,
	}
}

// 计算最优的 m（位数组大小） 和 k（哈希函数数量）
func calculateParameters(n int, p float64) (uint64, int) {
	m := uint64(-float64(n) * math.Log(p) / (math.Ln2 * math.Ln2))
	k := int(math.Ceil(math.Ln2 * float64(m) / float64(n)))
	return m, k
}

// 初始化哈希函数（使用不同种子的FNV算法）
func initHashFunctions(k int) []hash.Hash64 {
	var funcs []hash.Hash64
	for i := 0; i < k; i++ {
		fnvHash := fnv.New64a()
		fnvHash.Write([]byte{byte(i)}) // 用不同种子生成不同哈希
		funcs = append(funcs, fnvHash)
	}
	return funcs
}

// 添加元素到布隆过滤器
func (bf *BloomFilter) Add(element string) {
	for _, h := range bf.hashFuncs {
		h.Reset()
		h.Write([]byte(element))
		index := h.Sum64() % bf.size
		bf.bitset[index] = true
	}
}

// 检查元素是否存在（可能误判）
func (bf *BloomFilter) Contains(element string) bool {
	for _, h := range bf.hashFuncs {
		h.Reset()
		h.Write([]byte(element))
		index := h.Sum64() % bf.size
		if !bf.bitset[index] {
			return false
		}
	}
	return true
}
