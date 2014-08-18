# Yasha

Dota 2 (Source 1) replay parser written in Go.

We're working on both maintaining the parser for replays of both Source 1 and Source 2 Engines.

## About Source 1

This parser is able to parse replays recorded since about 2012.

## About Source 2

The upcoming Source 2 Engine requires several incompatible changes, so we work on it in the `source2` branch.

See also [Sange, the Dota 2 (Source 2) parser](http://github.com/dotabuff/sange)

## Installation

Simple as:

    $ go get github.com/dotabuff/yasha

And in your code:

    import "github.com/dotabuff/yasha"

Please be aware that you _can't import_ Sange and Yasha in the same binary!
The Protocol Buffer definitions conflict, and will panic.

## Usage

See the `examples` directory.

## License

MIT, see the LICENSE file.

## Help

If you have any questions, just ask manveru in the #dota2replay channel on QuakeNet.
