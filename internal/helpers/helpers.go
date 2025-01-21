package helpers

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/a-h/templ"
	"github.com/jackc/pgx/v5/pgtype"
)

func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func encodeUUID(src [16]byte) string {
	var buf [36]byte

	hex.Encode(buf[0:8], src[:4])
	buf[8] = '-'
	hex.Encode(buf[9:13], src[4:6])
	buf[13] = '-'
	hex.Encode(buf[14:18], src[6:8])
	buf[18] = '-'
	hex.Encode(buf[19:23], src[8:10])
	buf[23] = '-'
	hex.Encode(buf[24:], src[10:])

	return string(buf[:])
}

func GenerateSafeUrl(url string, uuid pgtype.UUID) templ.SafeURL {
	return templ.URL(GenerateUrl(url, uuid))
}

func GenerateUrl(url string, uuid pgtype.UUID) string {
	return fmt.Sprintf(url, encodeUUID(uuid.Bytes))
}

func Ternary(condition bool, valueIfTrue, valueIfFalse interface{}) interface{} {
    if condition {
        return valueIfTrue
    }
    return valueIfFalse
}


func TernaryString(condition bool, valueIfTrue string, valueIfFalse string) string {
    if condition {
        return valueIfTrue
    }
    return valueIfFalse
}
