base50
======

Base50 is another way to convert binary to ASCII, it is about 70% efficient.
Every 7 bytes is turned into 2 base50 numbers, with a maximum of 10 ASCII
characters used.

It removes the "difficult" characters from it's Alphabet,
much like base58. Unlike base58 each group of 7 bytes can be converted
independently, so it can serialize a stream of data (and doesn't have
exponential performance characteristics).

If the number of bytes converted isn't evenly divisible by 7 then the final
base50 number will be shorted. Base50 does output a stop/padding character but
that is only required for input if you concatenate two streams together.

  * To install: go get github.com/james-antill/cmd/base50/base50
