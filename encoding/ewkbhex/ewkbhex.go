// Package ewkbhex 实现了扩展的著名二进制 字符串的编码和解码
package ewkbhex

import (
	"encoding/binary"
	"encoding/hex"

	"github.com/chengxiaoer/geomGo"
	"github.com/chengxiaoer/geomGo/encoding/ewkb"
)

// Encode 将任意几何图形编码成二进制字符串
func Encode(g geom.T, byteOrder binary.ByteOrder) (string, error) {
	ewkb, err := ewkb.Marshal(g, byteOrder)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(ewkb), nil
}

// Decode 从二进制字符串中解码出几何图形.
func Decode(s string) (geom.T, error) {
	data, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return ewkb.Unmarshal(data)
}
