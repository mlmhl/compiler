package executable

import "fmt"

type ConstantPool struct {
	pool map[interface{}]int
}

func NewConstantPool() *ConstantPool {
	return &ConstantPool{
		pool: map[interface{}]int{},
	}
}

// AddIntIfAbsent add a new constant value if absent, and return
// the index of this constant value
func (pool *ConstantPool) AddIfAbsent(value interface{}) int {
	if index, ok := pool.pool[value]; ok {
		return index
	} else {
		index = len(pool.pool)
		pool.pool[value] = index
		return index
	}
}

// Encode encode the constant pool to code byte
func (pool *ConstantPool) Encode() []byte {
	buffer := []byte{'{'}
	for k, v := range(pool.pool) {
		buffer = append(buffer, fmt.Sprintf("%v:%d,", k, v)...)
	}
	if len(pool.pool) > 0 {
		buffer[len(pool.pool) - 1] = '}'
	} else {
		buffer = append(buffer, '}')
	}
	return buffer
}