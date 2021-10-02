package server

import (
	"bytes"
	"crypto"
	"crypto/hmac"
	_ "crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"lgc/cfg"
	"os"
)

type Token struct {
	Header *Header
	Claims *Claims
}

type Header struct {
	Alg string
	Typ string
}
type Claims struct {
	Uid  int
	Name string
}

var key []byte

// var algMap map[string]crypto.Hash
func init() {
	// algMap[crypto.SHA256.String()] = crypto.SHA256
	cfg, has := cfg.Cfg()["key"]
	if has {
		key = []byte(cfg.(string))
	} else {
		pwd, _ := os.Getwd()
		key = []byte(pwd)
	}
}

func (t *Token) Sign() ([]byte, error) {
	headerByte, err := json.Marshal(t.Header)
	if err != nil {
		return nil, err
	}
	claimsByte, err := json.Marshal(t.Claims)
	if err != nil {
		return nil, err
	}

	enc := base64.RawURLEncoding
	headerLen := enc.EncodedLen(len(headerByte))
	infoLen := headerLen + enc.EncodedLen(len(claimsByte)) + 2
	byteArr := make([]byte, infoLen)
	enc.Encode(byteArr, headerByte)
	byteArr[headerLen] = '.'
	enc.Encode(byteArr[headerLen+1:], claimsByte)
	byteArr[infoLen-1] = '.'

	// sign, err := sign(algMap[t.Header.alg], byteArr, key)
	sign, err := sign(crypto.SHA256, byteArr[:infoLen-1], key)
	if err != nil {
		return nil, err
	}
	enSign := make([]byte, enc.EncodedLen(len(sign)))
	enc.Encode(enSign, sign)

	return append(byteArr, enSign...), nil
}

func Varify(token []byte) (*Claims, error) {
	signIdx := bytes.LastIndexByte(token, '.')
	if signIdx < 0 {
		return nil, errors.New("token info error")
	}

	byteArr := token[:signIdx]
	// sign info from token
	sign, err := sign(crypto.SHA256, byteArr, key)
	if err != nil {
		return nil, err
	}
	inSign := token[signIdx+1:]
	enc := base64.RawURLEncoding
	// decode sign from token
	deSign := make([]byte, enc.DecodedLen(len(inSign)))
	enc.Decode(deSign, inSign)
	// compare 2 sign
	if !isEqual(sign, deSign) {
		return nil, errors.New("token info not equal")
	}
	// decode claims info
	claimIdx := bytes.IndexByte(token, '.')
	if claimIdx < 0 || claimIdx == signIdx {
		return nil, errors.New("not claims info")
	}
	deClaimsByte := make([]byte, enc.DecodedLen(signIdx-claimIdx))
	enc.Decode(deClaimsByte, token[claimIdx+1:signIdx])
	claims := &Claims{}
	err = json.Unmarshal(deClaimsByte, claims)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

// func (t *Token) Parse(token []byte) {}

// 根据选择的加密算法 对内容进行加密
func sign(hash crypto.Hash, in, key []byte) ([]byte, error) {
	if !hash.Available() {
		return nil, errors.New("unavailable hash alg")
	}
	hasher := hmac.New(hash.New, key)
	hasher.Write(in)
	return hasher.Sum(nil), nil
}

func parse(b []byte) {
	// obj := map[string]interface{}{}
	// fmt.Printf("%v\n", string(b))
	// json.Unmarshal(b, &obj)
	// fmt.Printf("%v\n", obj)
}

func isEqual(a, b []byte) bool {
	aLen := len(a)
	if aLen != len(b) {
		return false
	}
	for i := 0; i < aLen; i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
