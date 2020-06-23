package base50

import (
	"bytes"
	"testing"
)

func TestBase50AllOneByte(t *testing.T) {
	for i := 0; i < 256; i++ {
		ob := []byte{byte(i)}

		if EncodeLen(len(ob)) != 3 {
			t.Errorf("bad EncodeLen: len=%d\n",
				EncodeLen(len(ob)))
		}

		var encodedStore [3]byte
		encoded := encodedStore[:]

		encoded = Encode(encoded, ob)
		if encoded[len(encoded)-1] != '.' {
			t.Errorf("no stop character: %#x made %v (len=%d)\n",
				i, encoded, len(encoded))
		}
		encoded = encoded[:len(encoded)-1]
		if len(encoded) != 1 && len(encoded) != 2 {
			t.Errorf("bad len: %#x made %v (len=%d)\n",
				i, encoded, len(encoded))
		}

		if DecodeLen(len(encoded)) != 1 {
			t.Errorf("bad DecodeLen: len=%d\n",
				DecodeLen(len(encoded)))
		}

		var decodedStore [1]byte
		decoded := decodedStore[:]

		decoded, err := Decode(decoded, encoded)
		if err != nil {
			t.Errorf("bad err: %#x made %v\n",
				i, err)
		}
		if len(decoded) != 1 {
			t.Errorf("bad len: %#x made %v which decoded to (len=%d)\n",
				i, decoded, len(decoded))
		}

		if !bytes.Equal(ob, decoded) {
			t.Errorf("data not equal: in=%v: ut=%v\n",
				ob, decoded)
		}
	}
}

func TestBase50AllTwoByte(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping TestBase50AllTwoByte is too expensive")
	}

	for i := 0; i < 256; i++ {
		for j := 0; j < 256; j++ {
			ob := []byte{byte(i), byte(j)}

			if EncodeLen(len(ob)) != 4 {
				t.Errorf("bad EncodeLen: len=%d\n",
					EncodeLen(len(ob)))
			}

			var encodedStore [4]byte
			encoded := encodedStore[:]

			encoded = Encode(encoded, ob)
			if encoded[len(encoded)-1] != '.' {
				t.Errorf("no stop character: %#x%x made %v (len=%d)\n",
					i, j, encoded, len(encoded))
			}
			encoded = encoded[:len(encoded)-1]
			if len(encoded) != 3 {
				t.Errorf("bad len: %#x%x made %v (len=%d)\n",
					i, j, encoded, len(encoded))
			}

			if DecodeLen(len(encoded)) != 2 {
				t.Errorf("bad DecodeLen: len=%d\n",
					DecodeLen(len(encoded)))
			}

			var decodedStore [2]byte
			decoded := decodedStore[:]

			decoded, err := Decode(decoded, encoded)
			if err != nil {
				t.Errorf("bad err: %#x%x made %v\n",
					i, j, err)
			}
			if len(decoded) != 2 {
				t.Errorf("bad len: %#x%x made %v which decoded to (len=%d)\n",
					i, j, decoded, len(decoded))
			}

			if !bytes.Equal(ob, decoded) {
				t.Errorf("data not equal: in=%v: ut=%v\n",
					ob, decoded)
			}
		}
	}
}

func TestBase50AllThreeByte(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping TestBase50AllThreeByte is too expensive")
	}

	for i := 0; i < 256; i++ {
		t.Log("i=", i)
		for j := 0; j < 256; j++ {
			for k := 0; k < 256; k++ {
				ob := []byte{byte(i), byte(j), byte(k)}

				if EncodeLen(len(ob)) != 5 && EncodeLen(len(ob)) != 6 {
					t.Errorf("bad EncodeLen: len=%d\n",
						EncodeLen(len(ob)))
				}

				var encodedStore [6]byte
				encoded := encodedStore[:]

				encoded = Encode(encoded, ob)
				if encoded[len(encoded)-1] != '.' {
					t.Errorf("no stop character: %#x%x%x made %v (len=%d)\n",
						i, j, k, encoded, len(encoded))
				}
				encoded = encoded[:len(encoded)-1]
				if len(encoded) != 4 && len(encoded) != 5 {
					t.Errorf("bad len: %#x%x%x made %v (len=%d)\n",
						i, j, k, encoded, len(encoded))
				}

				if DecodeLen(len(encoded)) != 3 {
					t.Errorf("bad DecodeLen: len=%d\n",
						DecodeLen(len(encoded)))
				}

				var decodedStore [3]byte
				decoded := decodedStore[:]

				decoded, err := Decode(decoded, encoded)
				if err != nil {
					t.Errorf("bad err: %#x%x%x made %v\n",
						i, j, k, err)
				}
				if len(decoded) != 3 {
					t.Errorf("bad len: %#x%x%x made %v which decoded to (len=%d)\n",
						i, j, k, decoded, len(decoded))
				}

				if !bytes.Equal(ob, decoded) {
					t.Errorf("data not equal: in=%v: ut=%v\n",
						ob, decoded)
				}
			}
		}
	}
}

type tEncData struct {
	val       []byte
	encLen    int
	encOptLen bool
	enc       string
}

func testDataPrefix(t *testing.T, data []tEncData, prefix, encPrefix string) {
	t.Helper()

	for i := range data {
		val := append([]byte(prefix), data[i].val...)
		encLen := data[i].encLen + len(encPrefix)
		encOptLen := data[i].encOptLen
		enc := encPrefix + data[i].enc

		if encLen != EncodeLen(len(val)) {
			bad := true
			if encOptLen && (encLen-1) == EncodeLen(len(val)) {
				bad = false
			}
			if bad {
				t.Errorf("bad EncodeLen: %d: %v len=%d\n",
					i, val, EncodeLen(len(val)))
			}
		}

		encoded := make([]byte, EncodeLen(len(val)))
		encoded = Encode(encoded, val)

		if enc != string(encoded) {
			t.Errorf("data not equal: %v: %v\n tst=<%s>\n got <%s>\n",
				i, val, enc, string(encoded))
		}
		stop := encoded[len(encoded)-1] == '.'
		lenc := len(encoded)
		if stop {
			lenc--
		}
		decoded := make([]byte, DecodeLen(lenc))
		decoded, err := Decode(decoded, encoded)
		if err != nil {
			t.Errorf("bad err: %d: %v made %v\n",
				i, enc, err)
		}

		if !bytes.Equal(decoded, val) {
			t.Errorf("decoded not equal: %v: %v\n got <%s>\n",
				i, val, decoded)
		}

		if stop {
			// Now remove stop byte and try again...
			encoded = encoded[:len(encoded)-1]

			decoded, err := Decode(decoded, encoded)
			if err != nil {
				t.Errorf("bad err: %d: %v made %v\n",
					i, enc, err)
			}

			if !bytes.Equal(decoded, val) {
				t.Errorf("decoded not equal: %v: %v\n got <%s>\n",
					i, val, decoded)
			}
		}
	}
}

func testData(t *testing.T, data []tEncData) {
	t.Helper()
	testDataPrefix(t, data, "", "")
}

func testDataPrefixRev(t *testing.T, data []tEncData, prefix, encPrefix string) {
	t.Helper()

	for i := range data {
		dec := append([]byte(prefix), data[i].val...)
		decLen := data[i].encLen + len(prefix)
		val := encPrefix + data[i].enc

		if decLen > DecodeLen(len(val)) {
			t.Errorf("bad DecodeLen: %d: %v len=%d\n",
				i, val, EncodeLen(len(val)))
		}

		decoded := make([]byte, DecodeLen(len(val)))
		decoded, err := Decode(decoded, []byte(val))
		if err != nil {
			t.Errorf("bad err: %d: %v made %v\n",
				i, val, err)
		}

		if string(dec) != string(decoded) {
			t.Errorf("data not equal: %v: %v\n tst=<%s>\n got <%s>\n",
				i, val, dec, string(decoded))
		}
	}
}

func testDataRev(t *testing.T, data []tEncData) {
	t.Helper()
	testDataPrefixRev(t, data, "", "")
}

func TestBase50EncAbcdefg(t *testing.T) {
	data := []tEncData{
		{[]byte{'a'}, 3, false, "1x."},
		{[]byte("ab"), 4, false, "9yb."},
		{[]byte("abc"), 6, false, "KKuxH."},
		{[]byte("abcd"), 7, false, "KKuxP4."},
		{[]byte("abcde"), 9, false, "KKuxP0SW."},
		{[]byte("abcdef"), 10, false, "KKuxP2JF2."},
		{[]byte("abcdefg"), 10, false, "KKuxPEp1gJ"},
	}
	testData(t, data)
	testDataPrefix(t, data, "abcdefg", "KKuxPEp1gJ")
}

func TestBase50DecAbcdefg(t *testing.T) {
	data := []tEncData{
		{[]byte{'a'}, 1, false, "1x."},
		{[]byte("ab"), 2, false, "9yb."},
		{[]byte("abc"), 3, false, "KKuxH."},
		{[]byte("abcd"), 4, false, "KKuxP4."},
		{[]byte("abcde"), 5, false, "KKuxP0SW."},
		{[]byte("abcdef"), 6, false, "KKuxP2JF2."},
		{[]byte("abcdefg"), 7, false, "KKuxPEp1gJ"},
	}

	testDataRev(t, data)
	testDataPrefixRev(t, data, "a", "1x.")
}

func TestBase50Declen(t *testing.T) {
	data := []struct {
		val int
		tst int
	}{
		{0, 0},
		{1, 1},
		{2, 1},
		{3, 2},
		{4, 3},
		{5, 3},
		{6, 4},
		{7, 5},
		{8, 5},
		{9, 6},
		{10, 7},

		{10 * 2, 7 * 2},
		{(10 * 2) + 3, (7 * 2) + 2},
		{10 * 3, 7 * 3},
		{10 * 4, 7 * 4},
		{(10 * 4) + 7, (7 * 4) + 5},
		{10 * 5, 7 * 5},
		{10 * 6, 7 * 6},
		{10 * 7, 7 * 7},
	}

	for i := range data {
		val := data[i].val
		tst := data[i].tst

		if tst != DecodeLen(val) {
			t.Errorf("data not equal: %v: %v\n tst=<%v>\n got <%v>\n",
				i, val, tst, DecodeLen(val))
		}
	}
}

func TestBase50Edgecases(t *testing.T) {
	if len(Alphabet) != 50 {
		t.Errorf("Alphabet-len: %d\n",
			len(Alphabet))
	}

	data := []tEncData{
		{[]byte{0},
			3, false, "0."},
		{[]byte{0, 0},
			4, false, "000."},
		{[]byte{0, 0, 0},
			6, false, "00000."},
		{[]byte{0, 0, 0, 0},
			7, false, "000000."},
		{[]byte{0, 0, 0, 0, 0},
			9, false, "00000000."},
		{[]byte{0, 0, 0, 0, 0, 0},
			10, false, "000000000."},
		{[]byte{0, 0, 0, 0, 0, 0, 0},
			10, false, "0000000000"},

		{[]byte{0, 0, 0, 0, 0, 0, 0, 0},
			10 + 3, false, "00000000000."},
		{[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0},
			10 + 4, false, "0000000000000."},
		{[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			10 + 6, false, "000000000000000."},
		{[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			10 + 7, false, "0000000000000000."},
		{[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			10 + 9, false, "000000000000000000."},
		{[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			10 + 10, false, "0000000000000000000."},
		{[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			10 + 10, false, "00000000000000000000"},

		{[]byte{0xFF},
			3, false, "55."},
		{[]byte{0xFF, 0xFF},
			4, false, "XAh."},
		{[]byte{0xFF, 0xFF, 0xFF},
			6, false, "rxU8p."},
		{[]byte{0xFF, 0xFF, 0xFF, 0xFF},
			7, false, "rxU95J."},
		{[]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
			9, false, "rxU951du."},
		{[]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
			10, false, "rxU958NRW."},
		{[]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
			10, false, "rxU95rxU95"},

		// We want to test 50 rollover
		{[]byte{48}, 3, false, "y."},
		{[]byte{49}, 3, false, "z."},
		{[]byte{50}, 3, false, "10."},
		{[]byte{51}, 3, false, "11."},
	}
	testData(t, data)
}
