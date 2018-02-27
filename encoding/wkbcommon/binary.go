// Package wkbcommon 包含了 WKB and EWKB 编码相关的公共代码
package wkbcommon

import (
	"encoding/binary"
	"io"
	"math"
)

func readFloat(buf []byte, byteOrder binary.ByteOrder) float64 {
	u := byteOrder.Uint64(buf)
	return math.Float64frombits(u)
}

// ReadUInt32函数 从 r 中读取一个 uint32.
func ReadUInt32(r io.Reader, byteOrder binary.ByteOrder) (uint32, error) {
	var buf [4]byte
	if _, err := io.ReadFull(r, buf[:]); err != nil {
		return 0, err
	}
	return byteOrder.Uint32(buf[:]), nil
}

// ReadFloatArray函数 从 r 中读取一个 []float64.
func ReadFloatArray(r io.Reader, byteOrder binary.ByteOrder, array []float64) error {
	buf := make([]byte, 8*len(array))
	if _, err := io.ReadFull(r, buf); err != nil {
		return err
	}
	// 转换为浮动数组
	for i := range array {
		array[i] = readFloat(buf[8*i:], byteOrder)
	}
	return nil
}

// ReadByte函数 从 r 中读取一个 byte.
func ReadByte(r io.Reader) (byte, error) {
	var buf [1]byte
	if _, err := r.Read(buf[:]); err != nil {
		return 0, err
	}
	return buf[0], nil
}

func writeFloat(buf []byte, byteOrder binary.ByteOrder, value float64) {
	u := math.Float64bits(value)
	byteOrder.PutUint64(buf, u)
}

// WriteFloatArray函数  向 w 写入 一个[]float64.
func WriteFloatArray(w io.Writer, byteOrder binary.ByteOrder, array []float64) error {
	buf := make([]byte, 8*len(array))
	for i, f := range array {
		writeFloat(buf[8*i:], byteOrder, f)
	}
	_, err := w.Write(buf)
	return err
}

// WriteUInt32函数 向 w 写入一个 uint32.
func WriteUInt32(w io.Writer, byteOrder binary.ByteOrder, value uint32) error {
	var buf [4]byte
	byteOrder.PutUint32(buf[:], value)
	_, err := w.Write(buf[:])
	return err
}

// WriteByte函数 向 w 写入一个 byte.
func WriteByte(w io.Writer, value byte) error {
	var buf [1]byte
	buf[0] = value
	_, err := w.Write(buf[:])
	return err
}
