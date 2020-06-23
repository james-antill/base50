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

const panicDebug = true

// See the documentation on Encode(). Roughly 3.5 bin bytes fits into 5 ASCII
// bytes in base49 onwards.
func encodeInt(dst []byte, num uint32, outb int) {
	if panicDebug && num > 0xFFFFFFF {
		panic(num)
	}
	if panicDebug && outb < 1 {
		panic(outb)
	}
	if panicDebug && outb > 5 {
		panic(outb)
	}

	for i := outb - 1; i >= 0; i-- {
		dst[i] = Alphabet[num%50]
		num /= 50
	}

	if panicDebug && num > 0 {
		panic(num)
	}
}

func encodeBytes1(dst, src []byte, outb int) {
	num := uint32(src[0])
	if outb > 2 {
		num <<= 8
		num += uint32(src[1])
	}
	if outb > 3 {
		num <<= 8
		num += uint32(src[2])
	}
	if outb > 4 {
		num <<= 4
		num += ((uint32(src[3]&0xF0) >> 4) & 0x0F)
	}
	encodeInt(dst, num, outb)
}

func encodeBytes1Suffix(dst, src []byte) (int, []byte) {
	switch len(src) {
	case 1:
		if src[0] < 50 { // Optimze, Eg. 0 = 0
			encodeBytes1(dst, src, 1)
			return 1, src[1:]
		}
		encodeBytes1(dst, src, 2)
		return 2, src[1:]
	case 2:
		encodeBytes1(dst, src, 3)
		return 3, src[2:]
	case 3:
		if false { // If the number is small enough we could optimize here.
			encodeBytes1(dst, src, 4)
			return 4, src[3:]
		}
		var t [4]byte
		nsrc := t[:]
		copy(nsrc, src)
		encodeBytes1(dst, nsrc, 5)
		return 5, src[3:]
	default:
		encodeBytes1(dst, src, 5)
		return 5, src[3:]
	}
}

func encodeBytes2(dst, src []byte, outb int) {
	var num uint32
	num = uint32(src[0]) & 0x0F
	if outb > 1 {
		num <<= 8
		num += uint32(src[1])
	}
	if outb > 3 {
		num <<= 8
		num += uint32(src[2])
	}
	if outb > 4 {
		num <<= 8
		num += uint32(src[3])
	}
	encodeInt(dst, num, outb)
}

func encodeBytes2Suffix(dst, src []byte) int {
	switch len(src) {
	case 0:
		return 0
	case 1:
		encodeBytes2(dst, src, 1)
		return 1
	case 2:
		if false { // If the number is small enough we could optimize here.
			encodeBytes2(dst, src, 2)
			return 2
		}
		encodeBytes2(dst, src, 3)
		return 3
	case 3:
		encodeBytes2(dst, src, 4)
		return 4
	default:
		encodeBytes2(dst, src, 5)
		return 5
	}
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
		return whole + 5 + 1
	case 4:
		return whole + 6 + 1
	case 5:
		return whole + 8 + 1
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
// 3.5* = 0xFFFF_FFF          =         268435455
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

	// Get 2x 3.5 bytes at once, just to make life easier...
	for len(src) >= 7 {
		encodeBytes1(dst[idx:], src, 5)
		src = src[3:]
		idx += 5
		encodeBytes2(dst[idx:], src, 5)
		src = src[4:]
		idx += 5
	}

	if len(src) > 0 {
		i, ns := encodeBytes1Suffix(dst[idx:], src)
		idx += i
		src = ns
		idx += encodeBytes2Suffix(dst[idx:], src)
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
func from50Char(c byte) (uint32, bool) {
	if configIndexAlphabet {
		idx := strings.IndexByte(Alphabet, c)
		if idx == -1 {
			return uint32(c), false
		}
		return uint32(idx), true
	}

	switch {
	case '0' <= c && c <= '9':
		return uint32(c - '0'), true
	case 'A' == c:
		return uint32(c-'A') + 10, true
	case 'E' <= c && c <= 'H':
		return uint32(c-'E') + 10 + 1, true
	case 'J' <= c && c <= 'N':
		return uint32(c-'J') + 10 + 1 + 4, true
	case 'P' == c:
		return uint32(c-'P') + 10 + 1 + 4 + 5, true
	case 'R' <= c && c <= 'U':
		return uint32(c-'R') + 10 + 1 + 4 + 5 + 1, true
	case 'W' <= c && c <= 'Z':
		return uint32(c-'W') + 10 + 1 + 4 + 5 + 1 + 4, true
	case 'a' <= c && c <= 'b':
		return uint32(c-'a') + 10 + 1 + 4 + 5 + 1 + 4 + 4, true
	case 'd' <= c && c <= 'h':
		return uint32(c-'d') + 10 + 1 + 4 + 5 + 1 + 4 + 4 + 2, true
	case 'j' <= c && c <= 'k':
		return uint32(c-'j') + 10 + 1 + 4 + 5 + 1 + 4 + 4 + 2 + 5, true
	case 'm' <= c && c <= 'n':
		return uint32(c-'m') + 10 + 1 + 4 + 5 + 1 + 4 + 4 + 2 + 5 + 2, true
	case 'p' <= c && c <= 'u':
		return uint32(c-'p') + 10 + 1 + 4 + 5 + 1 + 4 + 4 + 2 + 5 + 2 + 2, true
	case 'w' <= c && c <= 'z':
		return uint32(c-'w') + 10 + 1 + 4 + 5 + 1 + 4 + 4 + 2 + 5 + 2 + 2 + 6, true
	}

	return 0, false
}

// InvalidByteError values describe errors resulting from an invalid byte in a base50 string.
type InvalidByteError byte

func (e InvalidByteError) Error() string {
	return fmt.Sprintf("base50: invalid byte: %#U", rune(e))
}

// InvalidTotalError values describe errors resulting from an invalid series of bytes in a base50 string.
type InvalidTotalError uint32

func (e InvalidTotalError) Error() string {
	return fmt.Sprintf("base50: invalid total > 0xFFFF_FFF: %x", int32(e))
}

// DecodeLen for every 5 bytes of input we have 3.5 bytes output
func DecodeLen(x int) int {
	// return ((x+1) * 10) / 7

	rem := x % 10
	whole := (x / 10) * 7
	switch rem {
	case 0:
		break
	case 1:
		fallthrough
	case 2:
		return whole + 1
	case 3:
		return whole + 2
	case 4:
		fallthrough
	case 5:
		return whole + 3
	case 6:
		return whole + 4
	case 7:
		fallthrough
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
// Decode expects that src contains only base50 characters and that src has even length. If the input is malformed, Decode returns the number of bytes decoded before the error.
func Decode(dst, src []byte) ([]byte, error) {
	count := 0

	odst := dst

	for len(src) >= 10 {
		var v [10]uint32
		var ok [10]bool

		if len(src) == 10 && src[9] == '.' { // End marker, so we can concat
			break
		}

		for i := 0; i < 10; i++ {
			if v[i], ok[i] = from50Char(src[i]); !ok[i] {
				return odst[:count], InvalidByteError(v[i])
			}
		}

		var num uint32
		for i := 0; i < 5; i++ {
			num *= 50
			num += v[i]
		}
		if num > 0xFFFFFFF {
			return odst[:count], InvalidTotalError(num)
		}

		dst[3] = byte(num & 0xF)
		dst[3] <<= 4 // It's the high 4 bits, stored in the low 4 bits of num.
		num >>= 4
		dst[2] = byte(num & 0xFF)
		num >>= 8
		dst[1] = byte(num & 0xFF)
		num >>= 8
		dst[0] = byte(num)

		num = 0
		for i := 5; i < 10; i++ {
			num *= 50
			num += v[i]
		}
		if num > 0xFFFFFFF {
			return odst[:count], InvalidTotalError(num)
		}

		dst[6] = byte(num & 0xFF)
		num >>= 8
		dst[5] = byte(num & 0xFF)
		num >>= 8
		dst[4] = byte(num & 0xFF)
		num >>= 8
		// The low four bits from above, which were the high 4 bits in num.
		dst[3] += byte(num)

		dst = dst[7:]
		src = src[10:]
		count += 7
	}

	if len(src) > 0 {
		var v [10]uint32
		var ok [10]bool

		if src[len(src)-1] == '.' { // End marker, so we can concat
			src = src[:len(src)-1]
		}

		for i := 0; i < len(src); i++ {
			if v[i], ok[i] = from50Char(src[i]); !ok[i] {
				return odst[:count], InvalidByteError(v[i])
			}
		}

		var num uint32
		i := 0
		for ; i < len(src) && i < 5; i++ {
			num *= 50
			num += v[i]
		}
		if num > 0xFFFFFFF {
			return odst[:count], InvalidTotalError(num)
		}

		// See the table on Encode()
		if len(src) > 5 {
			dst[3] = byte(num & 0xF)
			dst[3] <<= 4 // It's the high 4 bits, stored in the low 4 bits of num.
			count++
		}
		if len(src) > 4 {
			num >>= 4
		}
		if len(src) > 3 {
			dst[2] = byte(num & 0xFF)
			num >>= 8
			count++
		}
		if len(src) > 2 {
			dst[1] = byte(num & 0xFF)
			num >>= 8
			count++
		}
		dst[0] = byte(num)
		count++

		num = 0
		i = 5
		for ; i < len(src) && i < 10; i++ {
			num *= 50
			num += v[i]
		}
		if num > 0xFFFFFFF {
			return odst[:count], InvalidTotalError(num)
		}

		if len(src) > 9 { // ??
			dst[6] = byte(num & 0xFF)
			num >>= 8
			count++
		}
		if len(src) > 8 {
			dst[5] = byte(num & 0xFF)
			num >>= 8
			count++
		}
		if len(src) > 6 {
			dst[4] = byte(num & 0xFF)
			num >>= 8
			count++
		}
		if len(src) > 5 {
			// The low four bits from above, which were the high 4 bits in num.
			dst[3] += byte(num)
		}
	}

	return odst[:count], nil
}
