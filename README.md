[![Go Report Card](https://goreportcard.com/badge/github.com/james-antill/base50?style=flat-square)](https://goreportcard.com/report/github.com/james-antill/base50)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/james-antill/base50)
[![Release](https://img.shields.io/github/release/james-antill/base50.svg?style=flat-square)](https://github.com/james-antill/base50/releases/latest)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

base50
======

Base50 is another way to convert binary to ASCII, it is about 70% efficient.
Every 7 bytes is turned into a base50 number, with a maximum of 10 ASCII
characters used (all but the last one will be 10 characters).

It removes the "difficult" characters from it's Alphabet,
much like base58. Unlike base58 each group of 7 bytes can be converted
independently, so it can serialize a stream of data (and doesn't have
exponential performance characteristics).

If the number of bytes converted isn't evenly divisible by 7 then the final
base50 number will be shortened. Base50 does output a stop/padding character but
that is only required for input if you concatenate two streams together.

  * To install: go get github.com/james-antill/base50/cmd/base50

Example output
==============

Given the current sha256sum of the base50.go file:

  * base16
  * * ffe62174a2d8149e1ffa2701135610fde6407e425b20ff1c021480da60e3954e
  * base32
  * * 77TCC5FC3AKJ4H72E4ARGVQQ7XTEA7SCLMQP6HACCSANUYHDSVHA====
  * base36
  * * 6DLW81G3WXNT53STCHEDCNDN8QTES2XVUZE7TXAPNCAFCY5Q32
  * base50
  * * rwdnuFSFPFXXSNjMp4Ny2rXsm1K9pT5XnkA5a9FmKFm9XH.
  * base58
  * * JDvVWzaZ3UdQCpsdUb3nSC2As4S8C2zxAitBEtD552xM
  * base64
  * * /+YhdKLYFJ4f+icBE1YQ/eZAfkJbIP8cAhSA2mDjlU4=
  * base64url
  * * _-YhdKLYFJ4f-icBE1YQ_eZAfkJbIP8cAhSA2mDjlU4=

