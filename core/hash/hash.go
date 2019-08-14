package hash

func Hash(value string) uint32 {
	var seed uint32

	for i := 0; i < len(value); i++ {
		seed = hashCombine(seed, uint32(value[i]))
	}

	//fmt.Printf("RPC ID: %s, RPC Hash: %d\n", value, seed)

	return seed
}

func hashCombine(seed uint32, value uint32) uint32 {
	seed ^= value + 0x9e3779b9 + (seed << 6) + (seed >> 2)
	return seed
}
