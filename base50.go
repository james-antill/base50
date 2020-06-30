package base50

import (
	"fmt"
	"strings"
)

const configIndexAlphabet = false

// Alphabet is the 50 output characters used when displaying base50
//
// We need 50 characters, 26*2 + 10 = 62. So we can drop 12.
// We drop the obvious ones that look alike other things:
//  1. B ~ 8
//  2. D ~ 0
//  3. I ~ 1
//  4. O ~ 0
//  5. i ~ 1
//  6. l ~ 1
// Now things which maybe look like others:
//  1. C ~ 0, C ~ O
//  2. Q ~ 0
//  3. V ~ U
//  4. c ~ 0, c ~ o
//  5. o ~ 0
//  6. v ~ u
const Alphabet = "0123456789" +
	"A" + "EFGH" + "JKLMN" + "P" + "RSTU" + "WXYZ" +
	"ab" + "defgh" + "jk" + "mn" + "pqrstu" + "wxyz"

const configOpt = true
const panicDebug = true

// See the documentation on Encode(). Roughly 7 binary bytes fits into 10 ASCII
// bytes in base49 onwards.
func encodeInt64(dst []byte, num uint64, outb int) {
	// 0xFFFFFFFFFFFFFF = 0xFFFF_FFFF_FFFF_FF
	if panicDebug && num > 0xFFFFFFFFFFFFFF {
		panic(num)
	}
	if panicDebug && outb < 1 {
		panic(outb)
	}
	if panicDebug && outb > 10 {
		panic(outb)
	}

	//	fmt.Printf("JDBG: enc: %d %#x\n", outb, num)

	for i := outb - 1; i >= 0; i-- {
		dst[i] = Alphabet[num%50]
		num /= 50
	}

	if panicDebug && num > 0 {
		panic(num)
	}
}

// See doc. on Encode(), we encode the (upto) 7 binary bytes into a uint64
// then we'll turn that into 10 ASCII bytes.
func encodeBytes(dst, src []byte, outb int, opt uint64) int {
	num := uint64(src[0]) // 1 or 2 byte output
	if outb >= 3 {
		num <<= 8
		num += uint64(src[1])
	}
	if outb >= 4 {
		num <<= 8
		num += uint64(src[2])
	}
	if outb >= 6 {
		num <<= 8
		num += uint64(src[3])
	}
	if outb >= 7 {
		num <<= 8
		num += uint64(src[4])
	}
	if outb >= 9 {
		num <<= 8
		num += uint64(src[5])
	}
	if outb >= 10 {
		num <<= 8
		num += uint64(src[6])
	}
	if configOpt && opt > num {
		outb--
	}
	encodeInt64(dst, num, outb)

	return outb
}

// For the last group of bytes (< 7) we can output less than 10 ASCII bytes
func encodeBytesSuffix(dst, src []byte) int {
	switch len(src) {
	case 1:
		// 50**1. Optimze, Eg. 0 = 0
		return encodeBytes(dst, src, 2, 50)
	case 2:
		encodeBytes(dst, src, 3, 0)
		return 3
	case 3:
		// 50**4=6250000
		return encodeBytes(dst, src, 5, 6250000)
	case 4:
		return encodeBytes(dst, src, 6, 0)
	case 5:
		// 50**7=781250000000
		return encodeBytes(dst, src, 8, 781250000000)
	case 6:
		return encodeBytes(dst, src, 9, 0)
	default:
		break
	}

	return encodeBytes(dst, src, 10, 0)
}

// EncodeLen for every 3.5 bytes of input we have 5 bytes output and
// we might need an extra byte of "padding". Eg. 0x00 = "0." | 0xFF = "55."
func EncodeLen(x int) int {
	rem := x % 7
	whole := (x / 7) * 10
	switch rem {
	case 0:
		break
	case 1:
		return whole + 2 + 1 // Could be -1
	case 2:
		return whole + 3 + 1
	case 3:
		return whole + 5 + 1 // Could be -1
	case 4:
		return whole + 6 + 1
	case 5:
		return whole + 8 + 1 // Could be -1
	case 6:
		return whole + 9 + 1

	default:
		if panicDebug {
			panic(rem)
		}
	}

	return whole // No padding
}

// For every 7 bytes of binary data we produce 10 bytes of ASCII data.
// Basic math is that 0xFFFF_FFFF_FFFF_FF = (16**14)-1 = 72057594037927935
//                            zzzzz_zzzzz = (50**10)-1 = 97656249999999999
// For the complete list of bytes to base50 we have:
//  1   = 0xFF                =               255
//          zz                =              2499
//  2   = 0xFFFF              =             65535
//          zzz               =            124999
//  3   = 0xFFFF_FF           =          16777215
//          zzzz_z            =         312499999
//  4   = 0xFFFF_FFFF         =        4294967295
//          zzzz_zz           =       15624999999
//  5   = 0xFFFF_FFFF_FF      =     1099511627775
//          zzzz_zzzz         =    39062499999999
//  6   = 0xFFFF_FFFF_FFFF    =   281474976710655
//          zzzz_zzzz_z       =  1953124999999999
//  7   = 0xFFFF_FFFF_FFFF_FF = 72057594037927935
//          zzzz_zzzz_zz      = 97656249999999999

// Encode encodes src into EncodedLen(len(src)) bytes of dst. As a convenience,
// it returns the number of bytes written to dst, but this value is always
// EncodedLen(len(src)). Encode implements base50 encoding
func Encode(dst, src []byte) []byte {
	idx := 0

	// Get 7 bytes at once, just to make life easier...
	for len(src) >= 7 {
		_ = encodeBytes(dst[idx:], src, 10, 0)
		src = src[7:]
		idx += 10
	}

	if len(src) > 0 {
		i := encodeBytesSuffix(dst[idx:], src)
		idx += i
		dst[idx] = '.'
		idx++
	}

	return dst[:idx]
}

// EncodeToBytes returns the base50 encoding of src
func EncodeToBytes(src []byte) []byte {
	dst := make([]byte, EncodeLen(len(src)))
	return Encode(dst, src)
}

// EncodeToString returns the base50 encoding of src as a string
func EncodeToString(src []byte) string {
	return string(EncodeToBytes(src))
}

// from50Char converts a base50 character into its value and a success flag.
// in theory we could index the alphabet, but this should be faster...
func from50Char(c byte) (uint64, bool) {
	if configIndexAlphabet {
		idx := strings.IndexByte(Alphabet, c)
		if idx == -1 {
			return uint64(c), false
		}
		return uint64(idx), true
	}

	switch {
	case '0' <= c && c <= '9':
		return uint64(c - '0'), true
	case 'A' == c:
		return uint64(c-'A') + 10, true
	case 'E' <= c && c <= 'H':
		return uint64(c-'E') + 10 + 1, true
	case 'J' <= c && c <= 'N':
		return uint64(c-'J') + 10 + 1 + 4, true
	case 'P' == c:
		return uint64(c-'P') + 10 + 1 + 4 + 5, true
	case 'R' <= c && c <= 'U':
		return uint64(c-'R') + 10 + 1 + 4 + 5 + 1, true
	case 'W' <= c && c <= 'Z':
		return uint64(c-'W') + 10 + 1 + 4 + 5 + 1 + 4, true
	case 'a' <= c && c <= 'b':
		return uint64(c-'a') + 10 + 1 + 4 + 5 + 1 + 4 + 4, true
	case 'd' <= c && c <= 'h':
		return uint64(c-'d') + 10 + 1 + 4 + 5 + 1 + 4 + 4 + 2, true
	case 'j' <= c && c <= 'k':
		return uint64(c-'j') + 10 + 1 + 4 + 5 + 1 + 4 + 4 + 2 + 5, true
	case 'm' <= c && c <= 'n':
		return uint64(c-'m') + 10 + 1 + 4 + 5 + 1 + 4 + 4 + 2 + 5 + 2, true
	case 'p' <= c && c <= 'u':
		return uint64(c-'p') + 10 + 1 + 4 + 5 + 1 + 4 + 4 + 2 + 5 + 2 + 2, true
	case 'w' <= c && c <= 'z':
		return uint64(c-'w') + 10 + 1 + 4 + 5 + 1 + 4 + 4 + 2 + 5 + 2 + 2 + 6, true
	}

	return 0, false
}

// skipChar returns true for characters we should skip when decoding,
// mostly whitespace
func skipChar(c byte) bool {
	switch c {
	case '\t':
	case '\n':
	case '\r':
	case ' ':
	case '_': // Allows 0xFFFF_FFFF type stuff.
		return true
	}

	return false
}

// InvalidByteError values describe errors resulting from an invalid byte in a base50 string.
type InvalidByteError byte

func (e InvalidByteError) Error() string {
	return fmt.Sprintf("base50: invalid byte: %#U", rune(e))
}

// InvalidTotalError values describe errors resulting from an invalid series of
// bytes in a base50 string (the value is greater than 16**14).
type InvalidTotalError uint64

func (e InvalidTotalError) Error() string {
	num := uint64(e)
	if num > 0xFFFF_FFFF_FFFF_FF {
		return fmt.Sprintf("base50: invalid num > 0xFFFF_FFFF_FFFF_FF: %x", num)
	}
	return fmt.Sprintf("base50: invalid encoding (Eg. 56 should be 056): %x",
		num)
}

// DecodeLen for every 10 bytes of input we have 7 bytes output, apart from
// the last group.
func DecodeLen(x int) int {
	// return ((x+1) * 10) / 7

	rem := x % 10
	whole := (x / 10) * 7
	switch rem {
	case 0:
		break

	case 1:
		fallthrough // Again, or -1
	case 2:
		return whole + 1
	case 3:
		return whole + 2
	case 4:
		fallthrough // Again. or -1
	case 5:
		return whole + 3
	case 6:
		return whole + 4
	case 7:
		fallthrough // Again, or -1
	case 8:
		return whole + 5
	case 9:
		return whole + 6

	default:
		if panicDebug {
			panic(rem)
		}
	}

	return whole // No extra
}

// Decode decodes src into DecodedLen(len(src)) bytes, returning the actual
// number of bytes written to dst.
//
// Decode expects that src contains only base50 characters, or whitespace/underbar
// or the stop character if you've concatenated multiple encodings together.
// Decode also expects that src has a correct encoding (Eg. 56 is not valid).
// If the input is malformed, Decode returns the number of bytes decoded before
// the error.
func Decode(dst, src []byte) ([]byte, error) {
	count := 0
	odst := dst

	for len(src) > 0 {
		var v [10]uint64
		var ok [10]bool

		var Tbuf [10]byte
		nsrc := Tbuf[:0]
		for i := 0; len(src) > 0 && len(nsrc) < 10; {
			c := src[0]
			src = src[1:]

			if skipChar(c) {
				continue
			}
			if c == '.' {
				break
			}

			nsrc = append(nsrc, c)

			if v[i], ok[i] = from50Char(c); !ok[i] {
				return odst[:count], InvalidByteError(c)
			}

			i++
		}

		var num uint64

		for i := 0; i < len(nsrc) && i < 10; i++ {
			num *= 50
			num += v[i]
		}
		if num > 0xFFFFFFFFFFFFFF {
			return odst[:count], InvalidTotalError(num)
		}
		enum := num // Save the original num, for errors.

		//		fmt.Printf("JDBG: dec: %d %#x\n", len(nsrc), num)

		// See the table on Encode()
		if len(nsrc) >= 10 {
			dst[6] = byte(num & 0xFF)
			num >>= 8
			count++
		}
		if len(nsrc) >= 9 {
			dst[5] = byte(num & 0xFF)
			num >>= 8
			count++
		}
		if len(nsrc) >= 7 {
			dst[4] = byte(num & 0xFF)
			num >>= 8
			count++
		}
		if len(nsrc) >= 6 {
			dst[3] = byte(num & 0xFF)
			num >>= 8
			count++
		}
		if len(nsrc) >= 4 {
			dst[2] = byte(num & 0xFF)
			num >>= 8
			count++
		}
		if len(nsrc) >= 3 {
			dst[1] = byte(num & 0xFF)
			num >>= 8
			count++
		}
		dst[0] = byte(num & 0xFF)
		num >>= 8
		count++

		if panicDebug && len(nsrc) >= 10 && num > 0 {
			panic(num)
		}
		if num > 0 {
			return odst[:count], InvalidTotalError(enum)
		}

		dst = odst[count:]
	}

	return odst[:count], nil
}

// DecodeString returns the bytes represented by the base50 string s.
//
// DecodeString expects that src contains only base50 characters, or whitespace/underbar
// or the stop character if you've concatenated multiple encodings together.
// Decode also expects that src has a correct encoding (Eg. 56 is not valid).
// If the input is malformed, DecodeString returns the number of bytes decoded before
// the error.
func DecodeString(s string) ([]byte, error) {
	src := []byte(s)
	// We can use the source slice itself as the destination
	// because we read the "number" first and then write. And src always >.
	return Decode(src, src)
}
