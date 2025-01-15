package encrypts

import (
	"accgo/utils/logs"
	"testing"
)

// go test -v -count=1 -run TestEncrypt ./utils/encrypts/
func TestEncrypt(t *testing.T) {
	token := "gotoehell"

	srcs := []string{
		"18833338888",
		"123456",
		"我了个去",
		"我记得那时候身边的朋友",
	}

	for _, src := range srcs {
		enc, err := GoAhead(src, token)
		if err != nil {
			logs.Warn(err).Str("token", token).Str("src", src).Msg("encrypt fail")
			continue
		}
		dec, err := ComeBack(enc, token)
		if err != nil {
			logs.Warn(err).Str("token", token).Str("src", src).Msg("decrypt fail")
			continue
		}

		if dec == src {
			logs.Info().Str("token", token).Str("src", src).Str("enc", enc).Str("dec", dec).Msg("encrypt/decrypt succ")
		} else {
			logs.Info().Str("token", token).Str("src", src).Str("enc", enc).Str("dec", dec).Msg("encrypt/decrypt fail")
			t.Fail()
		}
	}
}
