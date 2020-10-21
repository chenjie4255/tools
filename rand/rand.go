package rand

import (
	"crypto/rand"
	"fmt"
	"math"
	mathrand "math/rand"
	"time"

	"encoding/base64"
	"encoding/hex"

	"github.com/chenjie4255/tools/errcode"
	"github.com/chenjie4255/errors"
)

func init() {
	mathrand.Seed(time.Now().Unix())
}

func Str(length int) (string, error) {
	if length <= 2 {
		return "", errors.NewWithTag("invalid random string length", errcode.ParamError)
	}

	buf := make([]byte, length/2)

	if _, err := rand.Read(buf); err != nil {
		return "", err
	}

	result := hex.EncodeToString(buf)
	return result, nil
}

func Digital(length int) (string, error) {
	if length < 2 {
		return "", errors.NewWithTag("invalid random string length", errcode.ParamError)
	}

	max := math.Pow10(length)
	val := mathrand.Int63n(int64(max))

	format := fmt.Sprintf("%%0%dd", length)
	return fmt.Sprintf(format, val), nil
}

func B64Str(byteLength int) (string, error) {
	if byteLength <= 0 {
		return "", errors.NewWithTag("invalid random string length", errcode.ParamError)
	}

	buf := make([]byte, byteLength)

	if _, err := rand.Read(buf); err != nil {
		return "", err
	}

	result := base64.RawURLEncoding.EncodeToString(buf)
	return result, nil
}

func Bytes(length int) ([]byte, error) {
	buf := make([]byte, length)
	if _, err := rand.Read(buf); err != nil {
		return nil, err
	}

	return buf, nil
}

func RedeemCode(length int) (string, error) {
	data, err := Bytes(length)
	if err != nil {
		return "", err
	}

	return Encode2NormalText(data), nil
}
