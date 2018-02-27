// Package wkbhex 实现非扩展二进制对于 string 的编码和解码
package wkbhex

import (
	"encoding/binary"
	"encoding/hex"

	"github.com/chengxiaoer/geomGo"
	"github.com/chengxiaoer/geomGo/encoding/wkb"
)

// Encode函数 将任意几何图形编码为 string
func Encode(g geom.T, byteOrder binary.ByteOrder) (string, error) {
	wkb, err := wkb.Marshal(g, byteOrder)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(wkb), nil
}

// Decode函数 将几何图形从 string 中解析出来
func Decode(s string) (geom.T, error) {
	data, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return wkb.Unmarshal(data)
}
