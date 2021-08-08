package utils

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

const HexPrefix = "0x"

// hex string <-> int
func HexString2Int(s string) (int, error) {
	s = strings.ToLower(s)
	i, err := strconv.ParseInt(strings.TrimPrefix(s, HexPrefix), 16, 64)
	if err != nil {
		return 0, err
	}
	return int(i), nil
}

// func HexString2BigInt(s string) (f *big.Float, b int, err error) {
// 	return big.ParseFloat(s, 16, 0, big.ToNearestEven)
// }

func Int2HexString(i int) string {
	return fmt.Sprintf(HexPrefix+"%x", i)
}

// hex string <-> big int
func HexString2BigInt(str string) (big.Int, error) {
	i := big.Int{}
	_, err := fmt.Sscan(HexPrefix+strings.TrimPrefix(strings.ToLower(str), HexPrefix), &i)
	return i, err
}

func BigIntToHexString(bigInt big.Int) string {
	if bigInt.BitLen() == 0 {
		return HexPrefix + "0"
	}

	return HexPrefix + strings.TrimPrefix(fmt.Sprintf("%x", bigInt.Bytes()), "0")
}

// hex string -> decimal
func HexString2Decimal(str string, exp int32) decimal.Decimal {
	i, _ := HexString2BigInt(str)
	return decimal.NewFromBigInt(&i, exp)
}

func String2Decimal(str string) decimal.Decimal {
	d, _ := decimal.NewFromString(str)
	return d
}

// bytes <-> hex string
func Bytes2HexString(bytes []byte) string {
	return HexPrefix + hex.EncodeToString(bytes)
}

func HexString2Bytes(str string) []byte {
	str = strings.TrimPrefix(strings.ToLower(str), HexPrefix)
	if len(str)%2 == 1 {
		str = "0" + str
	}
	b, _ := hex.DecodeString(str)
	return b
}

// decimal <-> big int
func DecimalToBigInt(d decimal.Decimal) *big.Int {
	n := new(big.Int)
	n, _ = n.SetString(d.Floor().String(), 0)
	return n
}
