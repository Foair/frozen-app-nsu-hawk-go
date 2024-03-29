package codec

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/xml"
	"math/big"
)

type xmlRsaKey struct {
	Modulus, Exponent, P, Q, DP, DQ, InverseQ, D string
}

type xmlRsaPubKey struct {
	XMLName           xml.Name `xml:"RSAKeyValue"`
	Modulus, Exponent string
}

// 解码 Base64 字符串并返回字节
func b64decode(str string) []byte {
	encoded, _ := base64.StdEncoding.DecodeString(str)
	return encoded
}

// 从 Base64 字符串中解析大整数
func b64ToBigint(str string) *big.Int {
	bigInt := &big.Int{}
	bigInt.SetBytes(b64decode(str))
	return bigInt
}

// XMLToPrivateKey 从 XML 中导出私钥
func XMLToPrivateKey(xmlPriKey string) *rsa.PrivateKey {
	xmlKey := xmlRsaKey{}
	_ = xml.Unmarshal([]byte(xmlPriKey), &xmlKey)
	key := &rsa.PrivateKey{
		PublicKey: rsa.PublicKey{
			N: b64ToBigint(xmlKey.Modulus),
			E: int(b64ToBigint(xmlKey.Exponent).Int64()),
		},
		D:      b64ToBigint(xmlKey.D),
		Primes: []*big.Int{b64ToBigint(xmlKey.P), b64ToBigint(xmlKey.Q)},
		Precomputed: rsa.PrecomputedValues{
			Dp:        b64ToBigint(xmlKey.DP),
			Dq:        b64ToBigint(xmlKey.DQ),
			Qinv:      b64ToBigint(xmlKey.InverseQ),
			CRTValues: ([]rsa.CRTValue)(nil), // 一个 []rsa.CRTValue 类型的空指针
		},
	}
	return key
}

// XMLPublicKey 从 XML 中导出公钥（XML 格式）
func XMLPublicKey(xmlPriKey string) string {
	xmlKey := xmlRsaPubKey{}
	_ = xml.Unmarshal([]byte(xmlPriKey), &xmlKey)
	bytes, _ := xml.Marshal(xmlKey)
	return string(bytes)
}
