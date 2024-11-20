package main

import (
	"encoding/hex"
	"strconv"
)
import "syscall/js"

type PlayerData struct {
	name string

	health int32

	extraMagicPower bool
	extraShotPower  bool
	extraShotSpeed  bool
	extraFightPower bool
	extraDefense    bool
	extraSpeed      bool

	earthTower bool
	waterTower bool
	fireTower  bool
	windTower  bool

	gold       int32
	experience int32

	levelsMagic      int16
	levelsShotPower  int16
	levelsShotSpeed  int16
	levelsFightPower int16
	levelsDefense    int16
	levelsSpeed      int16

	equippedWeapon    int16
	equippedArmor     int16
	equippedAccessory int16
	defendOrb         bool

	healDrink  bool
	warpWing   bool
	floatRing  bool
	fightRing  bool
	healRing   bool
	mirrorRing bool

	shopkeeperConversations int8

	timePlayed     int32
	monstersKilled int32

	unknown0x240 int32
	unknown0x130 int16

	playerType int16

	unknown0x18A int16

	keys    int16
	potions int16
}

var checksumValues = []int{
	0x089, 0x142, 0x20c, 0x0d0, 0x128, 0x241, 0x08a, 0x144, 0x218, 0x0e0,
	0x109, 0x242, 0x08c, 0x150, 0x228, 0x0c1, 0x10a, 0x244, 0x098, 0x160,
	0x209, 0x0c2, 0x10c, 0x250, 0x0a8, 0x141, 0x20a, 0x0c4, 0x118, 0x260,
}

type Bitpack struct {
	// running checksum (register D6)
	Checksum int
	// bit position (register D7)
	Position int
	// target buffer
	Buffer []byte
}

func bitIdxToByteIdxAndMask(bitIdx int) (int, byte) {
	byteIdx := bitIdx >> 3
	localBitIdx := 7 - (bitIdx & 7)
	mask := byte(1) << localBitIdx
	return byteIdx, mask
}

func (b *Bitpack) ReadBit(pos int) int {
	if pos < 0 || pos >= len(b.Buffer)*8 {
		panic("out of range")
	}
	byteIdx, bitMask := bitIdxToByteIdxAndMask(pos)
	if b.Buffer[byteIdx]&bitMask == bitMask {
		b.Checksum ^= checksumValues[pos%30]
		return 1
	}
	return 0
}

func (b *Bitpack) WriteBit(pos int, bit int) {
	if pos < 0 || pos >= len(b.Buffer)*8 {
		panic("out of range")
	}
	byteIdx, mask := bitIdxToByteIdxAndMask(pos)
	if bit > 0 {
		b.Checksum ^= checksumValues[pos%30]
		b.Buffer[byteIdx] |= mask
	} else {
		b.Buffer[byteIdx] &= ^mask
	}
}

func (b *Bitpack) Pack(value int, numBits int) {
	js.Global().Get("console").Call("log", "pack", value, numBits)
	for i := numBits - 1; i >= 0; i-- {
		mask := int(1) << i
		bit := 0
		if (value & mask) > 0 {
			bit = 1
		}
		b.WriteBit(b.Position, bit)
		b.Position++
	}
}

func (b *Bitpack) Unpack(numBits int) int {
	result := 0
	for i := numBits - 1; i >= 0; i-- {
		result |= b.ReadBit(b.Position) << i
		b.Position++
	}
	return result
}

var table = []int16{
	0x00, 0x0d, 0x19, 0x24, 0x2e, 0x37, 0x3f, 0x46, 0x4c, 0x51, 0x55, 0x58, 0x5A}

func packKeysAndPotions(playerData *PlayerData) int16 {
	result := table[playerData.keys]
	result += playerData.potions
	return result
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func encodePassword(playerData *PlayerData) (result string) {
	defer func() {
		if r := recover(); r != nil {
			result = "ERROR"
		}
	}()

	buffer := make([]byte, 19)
	b := Bitpack{Checksum: 0, Position: 0, Buffer: buffer}

	b.Pack(int(playerData.health/0x64), 0xa)

	b.Pack(boolToInt(playerData.extraMagicPower), 0x1)
	b.Pack(boolToInt(playerData.extraShotPower), 0x1)
	b.Pack(boolToInt(playerData.extraShotSpeed), 0x1)
	b.Pack(boolToInt(playerData.extraFightPower), 0x1)
	b.Pack(boolToInt(playerData.extraDefense), 0x1)
	b.Pack(boolToInt(playerData.extraSpeed), 0x1)

	towers := 0
	if playerData.earthTower {
		towers += 1
	}
	if playerData.waterTower {
		towers += 2
	}
	if playerData.fireTower {
		towers += 4
	}
	if playerData.windTower {
		towers += 8
	}
	b.Pack(towers, 0x4)

	b.Pack(int(playerData.gold), 0x10)
	b.Pack(int(playerData.experience), 0x10)

	b.Pack(int(playerData.levelsMagic), 0x3)
	b.Pack(int(playerData.levelsShotPower), 0x3)
	b.Pack(int(playerData.levelsShotSpeed), 0x3)
	b.Pack(int(playerData.levelsFightPower), 0x3)
	b.Pack(int(playerData.levelsDefense), 0x3)
	b.Pack(int(playerData.levelsSpeed), 0x3)

	b.Pack(int(playerData.equippedWeapon), 0x3)
	b.Pack(int(playerData.equippedArmor), 0x2)
	b.Pack(int(playerData.equippedAccessory), 0x2)
	b.Pack(boolToInt(playerData.defendOrb), 0x1)

	inventory := 0
	if playerData.healDrink {
		inventory += 1
	}
	if playerData.warpWing {
		inventory += 2
	}
	if playerData.floatRing {
		inventory += 4
	}
	if playerData.fightRing {
		inventory += 8
	}
	if playerData.healRing {
		inventory += 16
	}
	if playerData.mirrorRing {
		inventory += 32
	}
	b.Pack(inventory, 0x6)

	b.Pack(int(packKeysAndPotions(playerData)), 0x7)

	shopkeeperFlags := ([]int8{0, 1, 3, 7, 15, 31})[playerData.shopkeeperConversations]
	b.Pack(int(shopkeeperFlags), 0x3)
	b.Pack(int(min(playerData.timePlayed/0x3c, 0x176f)), 0xd)
	b.Pack(int(min(playerData.monstersKilled, 0x186a0)), 0x11)

	b.Pack(int(min(playerData.unknown0x240, 0x7f)), 0x7)
	b.Pack(int(playerData.unknown0x130), 0x2)
	b.Pack(int(playerData.playerType), 0x2)
	b.Pack(int(min(playerData.unknown0x18A, 0x1f)), 0x5)

	b.Pack(b.Checksum, 0xa)

	logBuffer(buffer)

	b = Bitpack{Checksum: 0, Position: 0, Buffer: buffer}

	password := ""
	for n := 0x1d; n >= 0; n-- {
		d0 := b.Unpack(0x5) + n
		nameChar := playerData.name[(0x1d-n)%len(playerData.name)]
		d0 += int(nameChar)
		d0 &= 0x1f
		password += string("J2H=K7+U0W9GTR3F4:6LC-1Y8EXMD5PA"[d0])
	}

	return password
}

func logBuffer(buffer []byte) {
	s := hex.EncodeToString(buffer)
	for i := 4; i < len(s); i += 5 {
		s = s[:i] + " " + s[i:]
	}
	js.Global().Get("console").Call("warn", s)
}

func toInt(str string) int {
	i, _ := strconv.Atoi(str)
	return i
}

func encodeWrapper() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		js.Global().Get("console").Call("warn", args[5].String())
		playerData := PlayerData{
			name:                    args[0].String(),
			health:                  int32(toInt(args[1].String())),
			extraMagicPower:         args[2].Bool(),
			extraShotPower:          args[3].Bool(),
			extraShotSpeed:          args[4].Bool(),
			extraFightPower:         args[5].Bool(),
			extraDefense:            args[6].Bool(),
			extraSpeed:              args[7].Bool(),
			earthTower:              args[8].Bool(),
			waterTower:              args[9].Bool(),
			fireTower:               args[10].Bool(),
			windTower:               args[11].Bool(),
			gold:                    int32(toInt(args[12].String())),
			experience:              int32(toInt(args[13].String())),
			levelsMagic:             int16(toInt(args[14].String())),
			levelsShotPower:         int16(toInt(args[15].String())),
			levelsShotSpeed:         int16(toInt(args[16].String())),
			levelsFightPower:        int16(toInt(args[17].String())),
			levelsDefense:           int16(toInt(args[18].String())),
			levelsSpeed:             int16(toInt(args[19].String())),
			equippedWeapon:          int16(toInt(args[20].String())),
			equippedArmor:           int16(toInt(args[21].String())),
			equippedAccessory:       int16(toInt(args[22].String())),
			defendOrb:               args[23].Bool(),
			healDrink:               args[24].Bool(),
			warpWing:                args[25].Bool(),
			floatRing:               args[26].Bool(),
			fightRing:               args[27].Bool(),
			healRing:                args[28].Bool(),
			mirrorRing:              args[29].Bool(),
			shopkeeperConversations: int8(toInt(args[30].String())),
			timePlayed:              int32(toInt(args[31].String())),
			monstersKilled:          int32(toInt(args[32].String())),
			unknown0x240:            int32(toInt(args[33].String())),
			unknown0x130:            int16(toInt(args[34].String())),
			playerType:              int16(toInt(args[35].String())),
			unknown0x18A:            int16(toInt(args[36].String())),
			keys:                    int16(toInt(args[37].String())),
			potions:                 int16(toInt(args[38].String())),
		}
		return encodePassword(&playerData)
	})
}

func main() {
	js.Global().Set("encode", encodeWrapper())
	<-make(chan struct{})
}
