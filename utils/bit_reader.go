package utils

import (
	"bytes"
	"encoding/binary"
	"math"
	"strconv"

	"github.com/elobuff/d2rp/core"
	"github.com/elobuff/d2rp/core/send_tables"
	dota "github.com/elobuff/d2rp/dota"
)

func flag(prop dota.CSVCMsg_SendTableSendpropT) send_tables.Flag {
	return send_tables.Flag(prop.GetFlags())
}

const (
	CoordIntegerBits            = 14
	CoordFractionalBits         = 5
	CoordDenominator            = (1 << CoordFractionalBits)
	CoordResolution     float64 = (1.0 / CoordDenominator)

	NormalFractionalBits         = 11
	NormalDenominator            = ((1 << NormalFractionalBits) - 1)
	NormalResolution     float64 = (1.0 / NormalDenominator)
)

type BitReader struct {
	buffer     []byte
	currentBit int
}

func NewBitReader(buffer []byte) *BitReader {
	if len(buffer) == 0 {
		panic("empty buffer?")
	}
	return &BitReader{buffer: buffer}
}

func (br *BitReader) Length() int      { return len(br.buffer) }
func (br *BitReader) CurrentBit() int  { return br.currentBit }
func (br *BitReader) CurrentByte() int { return br.currentBit / 8 }
func (br *BitReader) BitsLeft() int    { return (len(br.buffer) * 8) - br.currentBit }
func (br *BitReader) BytesLeft() int   { return len(br.buffer) - br.CurrentByte() }

type SeekOrigin int

const (
	Current SeekOrigin = iota
	Begin
	End
)

func (br *BitReader) SeekBits(offset int, origin SeekOrigin) {
	if origin == Current {
		br.currentBit += offset
	} else if origin == Begin {
		br.currentBit = offset
	} else if origin == End {
		br.currentBit = (len(br.buffer) * 8) - offset
	}
	if br.currentBit < 0 || br.currentBit > (len(br.buffer)*8) {
		panic("out of range")
	}
}

func (br *BitReader) ReadUBitsByteAligned(nBits int) uint {
	if nBits%8 != 0 {
		panic("Must be multple of 8")
	}
	if br.currentBit%8 != 0 {
		panic("Current bit is not byte-aligned")
	}
	var result uint
	for i := 0; i < nBits/8; i++ {
		result += uint(br.buffer[br.CurrentByte()] << (uint(i) * 8))
		br.currentBit += 8
	}
	return result
}

func (br *BitReader) ReadUBitsNotByteAligned(nBits int) uint {
	bitOffset := br.currentBit % 8
	nBitsToRead := bitOffset + nBits
	nBytesToRead := nBitsToRead / 8
	if nBitsToRead%8 != 0 {
		nBytesToRead += 1
	}

	var currentValue uint64
	for i := 0; i < nBytesToRead; i++ {
		b := br.buffer[br.CurrentByte()+1]
		currentValue += (uint64(b) << (uint64(i) * 8))
	}
	currentValue >>= uint(bitOffset)
	currentValue &= ((1 << uint(nBits)) - 1)
	br.currentBit += nBits
	return uint(currentValue)
}
func (br *BitReader) ReadVarInt() uint {
	var b uint = 0x80
	var count int
	var result uint
	for (b & 0x80) == 0x80 {
		if count == 5 {
			return result
		} else if br.CurrentByte() >= len(br.buffer) {
			return result
		}
		b = br.ReadUBits(8)
		result |= (b & 0x7f) << uint(7*count)
		count++
	}
	return result
}
func (br *BitReader) ReadUBits(nBits int) uint {
	if nBits <= 0 || nBits > 32 {
		panic("Value must be a positive integer between 1 and 32 inclusive.")
	}
	if br.CurrentBit()+nBits > br.Length()*8 {
		panic("Out of range")
	}
	if br.currentBit%8 == 0 && nBits%8 == 0 {
		return br.ReadUBitsByteAligned(nBits)
	}
	return br.ReadUBitsNotByteAligned(nBits)
}
func (br *BitReader) ReadBits(nBits int) int {
	result := br.ReadUBits(nBits - 1)
	if br.ReadBoolean() {
		result = -((1 << uint(nBits-1)) - result)
	}
	return int(result)
}
func (br *BitReader) ReadBoolean() bool {
	if br.CurrentBit()+1 > br.Length()*8 {
		panic("Out of range")
	}
	currentByte := br.CurrentBit() / 8
	bitOffset := br.CurrentBit() % 8
	result := br.buffer[currentByte]&(1<<uint(bitOffset)) != 0
	br.currentBit++
	return result
}
func (br *BitReader) ReadByte() byte {
	return byte(br.ReadUBits(8))
}
func (br *BitReader) ReadSByte() int8 {
	return int8(br.ReadBits(8))
}
func (br *BitReader) ReadBytes(nBytes int) []byte {
	if nBytes <= 0 {
		panic("Must be positive integer: nBytes")
	}
	result := make([]byte, nBytes)
	for i := 0; i < nBytes; i++ {
		result[i] = br.ReadByte()
	}
	return result
}
func (br *BitReader) ReadBitFloat() float32 {
	b := bytes.NewBuffer(br.ReadBytes(4))
	var f float32
	binary.Read(b, binary.LittleEndian, &f)
	return f
}
func (br *BitReader) ReadBitNormal() float64 {
	signbit := br.ReadBoolean()
	fractval := float64(br.ReadUBits(NormalFractionalBits))
	value := fractval * NormalResolution
	if signbit {
		value = -value
	}
	return value
}
func (br *BitReader) ReadBitCellCoord(bits int, integral, lowPrecision bool) float64 {
	intval := 0
	fractval := 0
	value := 0.0
	if integral {
		value = float64(br.ReadBits(bits))
	} else {
		intval = br.ReadBits(bits)
		if lowPrecision {
			fractval = br.ReadBits(3)
			value = float64(intval) + (float64(fractval) * (1.0 / (1 << 3)))
		} else {
			fractval = br.ReadBits(5)
			value = float64(intval) + (float64(fractval) * (1.0 / (1 << 5)))
		}
	}
	return value
}
func (br *BitReader) ReadBitCoord() float64 {
	intFlag := br.ReadBoolean()
	fractFlag := br.ReadBoolean()
	value := 0.0
	if intFlag || fractFlag {
		negative := br.ReadBoolean()
		if intFlag {
			value += float64(br.ReadUBits(CoordIntegerBits)) + 1
		}
		if fractFlag {
			value += float64(br.ReadUBits(CoordFractionalBits)) * CoordResolution
		}
		if negative {
			value = -value
		}
	}
	return value
}
func (br *BitReader) ReadString() string {
	bs := []byte{}
	for {
		b := br.ReadByte()
		if b == 0 {
			break
		}
		bs = append(bs, b)
	}
	return string(bs)
}

func (br *BitReader) ReadFloat(prop dota.CSVCMsg_SendTableSendpropT) float64 {
	if result, ok := br.ReadSpecialFloat(prop); ok {
		return result
	}
	dwInterp := float64(br.ReadUBits(int(prop.GetNumBits())))
	bits := 1 << uint(prop.GetNumBits())
	result := dwInterp / float64(bits-1)
	low, high := float64(prop.GetLowValue()), float64(prop.GetHighValue())
	return low + (high-low)*result
}

func (br *BitReader) ReadLengthPrefixedString() string {
	stringLength := int(br.ReadUBits(9))
	if stringLength > 0 {
		return string(br.ReadBytes(stringLength))
	}
	return ""
}

func (br *BitReader) ReadVector(prop dota.CSVCMsg_SendTableSendpropT) *core.Vector {
	result := &core.Vector{
		X: br.ReadFloat(prop),
		Y: br.ReadFloat(prop),
	}
	if (flag(prop) & send_tables.SPROP_NORMAL) == 0 {
		result.Z = br.ReadFloat(prop)
	} else {
		signbit := br.ReadBoolean()
		v0v0v1v1 := result.X*result.X + result.Y*result.Y
		if v0v0v1v1 < 1.0 {
			result.Z = math.Sqrt(1.0 - v0v0v1v1)
		} else {
			result.Z = 0.0
		}
		if signbit {
			result.Z *= -1.0
		}
	}

	return result
}

func (br *BitReader) ReadVectorXY(prop dota.CSVCMsg_SendTableSendpropT) *core.Vector {
	return &core.Vector{
		X: br.ReadFloat(prop),
		Y: br.ReadFloat(prop),
	}
}

func (br *BitReader) ReadInt(prop dota.CSVCMsg_SendTableSendpropT) int {
	if (flag(prop) & send_tables.SPROP_UNSIGNED) != 0 {
		return int(br.ReadUBits(int(prop.GetNumBits())))
	}
	return br.ReadBits(int(prop.GetNumBits()))
}

func (br *BitReader) ReadSpecialFloat(prop dota.CSVCMsg_SendTableSendpropT) (float64, bool) {
	flags := flag(prop)
	if (flags & send_tables.SPROP_COORD) != 0 {
		return br.ReadBitCoord(), true
	} else if (flags & send_tables.SPROP_COORD_MP) != 0 {
		panic("wtf")
	} else if (flags & send_tables.SPROP_COORD_MP_INTEGRAL) != 0 {
		panic("wtf")
	} else if (flags & send_tables.SPROP_COORD_MP_LOWPRECISION) != 0 {
		panic("wtf")
	} else if (flags & send_tables.SPROP_CELL_COORD) != 0 {
		return br.ReadBitCellCoord(int(prop.GetNumBits()), false, false), true
	} else if (flags & send_tables.SPROP_CELL_COORD_INTEGRAL) != 0 {
		return br.ReadBitCellCoord(int(prop.GetNumBits()), true, false), true
	} else if (flags & send_tables.SPROP_CELL_COORD_LOWPRECISION) != 0 {
		return br.ReadBitCellCoord(int(prop.GetNumBits()), false, true), true
	} else if (flags & send_tables.SPROP_NOSCALE) != 0 {
		return float64(br.ReadBitFloat()), true
	} else if (flags & send_tables.SPROP_NORMAL) != 0 {
		return br.ReadBitNormal(), true
	}
	return 0, false
}

func (br *BitReader) ReadInt64(prop dota.CSVCMsg_SendTableSendpropT) uint64 {
	var low, high uint
	if (flag(prop) & send_tables.SPROP_UNSIGNED) != 0 {
		low = br.ReadUBits(32)
		high = br.ReadUBits(32)
	} else {
		br.SeekBits(1, Current)
		low = br.ReadUBits(32)
		high = br.ReadUBits(31)
	}
	res := uint64(high)
	res = (res << 32)
	res = res | uint64(low)
	return res
}

func (br *BitReader) ReadNextEntityIndex(oldEntity int) int {
	ret := br.ReadUBits(4)
	more1 := br.ReadBoolean()
	more2 := br.ReadBoolean()
	if more1 {
		ret += (br.ReadUBits(4) << 4)
	}
	if more2 {
		ret += (br.ReadUBits(8) << 4)
	}
	return oldEntity + 1 + int(ret)
}

func (br *BitReader) ReadPropertiesIndex() []int {
	props := []int{}
	prop := -1
	for {
		if br.ReadBoolean() {
			prop += 1
			props = append(props, prop)
		} else {
			value := br.ReadVarInt()
			if value == 16383 {
				break
			}
			prop += (int(value) + 1)
			props = append(props, prop)
		}
	}
	return props
}

func (br *BitReader) ReadPropertiesValues(mapping []dota.CSVCMsg_SendTableSendpropT, multiples map[string]int, indices []int) map[string]interface{} {
	values := map[string]interface{}{}
	for _, index := range indices {
		prop := mapping[index]
		multiple := multiples[prop.GetDtName()+"."+prop.GetVarName()] > 1
		elements := 1
		if (flag(prop) & send_tables.SPROP_INSIDEARRAY) != 0 {
			elements = int(br.ReadUBits(6))
		}
		for k := 0; k < elements; k++ {
			key := prop.GetDtName() + "." + prop.GetVarName()
			if multiple {
				key += ("-" + strconv.Itoa(index))
			}
			if elements > 1 {
				key += ("-" + strconv.Itoa(k))
			}
			switch send_tables.DPTType(prop.GetType()) {
			case send_tables.DPT_Int:
				if (flag(prop) & send_tables.SPROP_ENCODED_AGAINST_TICKCOUNT) != 0 {
					values[key] = br.ReadVarInt()
				} else {
					values[key] = br.ReadInt(prop)
				}
			case send_tables.DPT_Float:
				if (flag(prop) & send_tables.SPROP_ENCODED_AGAINST_TICKCOUNT) != 0 {
					panic("SPROP_ENCODED_AGAINST_TICKCOUNT")
				} else {
					values[key] = br.ReadFloat(prop)
				}
			case send_tables.DPT_Vector:
				if (flag(prop) & send_tables.SPROP_ENCODED_AGAINST_TICKCOUNT) != 0 {
					panic("SPROP_ENCODED_AGAINST_TICKCOUNT")
				} else {
					values[key] = br.ReadVector(prop)
				}
			case send_tables.DPT_VectorXY:
				if (flag(prop) & send_tables.SPROP_ENCODED_AGAINST_TICKCOUNT) != 0 {
					panic("SPROP_ENCODED_AGAINST_TICKCOUNT")
				} else {
					values[key] = br.ReadVectorXY(prop)
				}
			case send_tables.DPT_String:
				if (flag(prop) & send_tables.SPROP_ENCODED_AGAINST_TICKCOUNT) != 0 {
					panic("SPROP_ENCODED_AGAINST_TICKCOUNT")
				} else {
					values[key] = br.ReadLengthPrefixedString()
				}
			case send_tables.DPT_Int64:
				if (flag(prop) & send_tables.SPROP_ENCODED_AGAINST_TICKCOUNT) != 0 {
					panic("SPROP_ENCODED_AGAINST_TICKCOUNT")
				} else {
					values[key] = br.ReadInt64(prop)
				}
			}
		}
	}
	return values
}
