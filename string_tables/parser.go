package string_tables

import (
	"math"

	"github.com/davecgh/go-spew/spew"
	"github.com/elobuff/d2rp/core/utils"
)

func Parse(bytes []byte, numEntries, maxEntries int32, isFixedSize bool, numBits int32) map[int]*StringTableItem {
	spew.Dump(bytes)
	result := map[int]*StringTableItem{}
	lastEntry := -1
	history := make([]string, 0, 32)
	br := utils.NewBitReader(bytes)
	br.SeekBits(1, utils.Begin)
	for i := int32(0); i < numEntries; i++ {
		item := &StringTableItem{}
		entryIndex := lastEntry + 1
		if !br.ReadBoolean() {
			entryIndex = int(br.ReadUBits(int(math.Log(float64(maxEntries)) / math.Log(2))))
		}
		lastEntry = entryIndex
		spew.Dump(lastEntry)
		if br.ReadBoolean() {
			value := ""
			substringcheck := br.ReadBoolean()
			if substringcheck {
				index := int(br.ReadUBits(5))
				bytestocopy := int(br.ReadUBits(5))
				spew.Dump(index)
				spew.Dump(bytestocopy)
				spew.Dump(history)
				spew.Dump(history[index])
				value = history[index][0:bytestocopy] + br.ReadString()
			} else {
				value = br.ReadString()
			}
			item.Str = value
			history = append(history, value)
		}
		if br.ReadBoolean() {
			if isFixedSize {
				item.Data = []byte{byte(br.ReadBits(int(numBits)))}
			} else {
				length := int(br.ReadUBits(14))
				item.Data = br.ReadBytes(length)
			}
		}
		if len(history) > 32 {
			newHistory := make([]string, 0, 32)
			copy(newHistory, history[1:])
			history = newHistory
		}
		result[entryIndex] = item
	}
	return result
}
