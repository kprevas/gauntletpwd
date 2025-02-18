function charValues(baseAddr, t)
    table.insert(t, memory.read_u32_be(baseAddr + 0x134) .. ',')
    table.insert(t, tostring(memory.read_u8(baseAddr + 0x96) > 0) .. ',')
    table.insert(t, tostring(memory.read_u8(baseAddr + 0x98) > 0) .. ',')
    table.insert(t, tostring(memory.read_u8(baseAddr + 0x9a) > 0) .. ',')
    table.insert(t, tostring(memory.read_u8(baseAddr + 0x9c) > 0) .. ',')
    table.insert(t, tostring(memory.read_u8(baseAddr + 0x9e) > 0) .. ',')
    table.insert(t, tostring(memory.read_u8(baseAddr + 0xa0) > 0) .. ',')
    local towers = memory.read_u8(baseAddr + 0x16d)
    table.insert(t, tostring(bit.check(towers, 0)) .. ',')
    table.insert(t, tostring(bit.check(towers, 1)) .. ',')
    table.insert(t, tostring(bit.check(towers, 2)) .. ',')
    table.insert(t, tostring(bit.check(towers, 3)) .. ',')
    table.insert(t, memory.read_u32_be(baseAddr + 0x124) .. ',')
    table.insert(t, memory.read_u32_be(baseAddr + 0x128) .. ',')
    table.insert(t, memory.read_u16_be(baseAddr + 0x2a0) .. ',')
    table.insert(t, memory.read_u16_be(baseAddr + 0x2a2) .. ',')
    table.insert(t, memory.read_u16_be(baseAddr + 0x2a4) .. ',')
    table.insert(t, memory.read_u16_be(baseAddr + 0x2a6) .. ',')
    table.insert(t, memory.read_u16_be(baseAddr + 0x2a8) .. ',')
    table.insert(t, memory.read_u16_be(baseAddr + 0x2aa) .. ',')
    table.insert(t, memory.read_u16_be(baseAddr + 0x298) .. ',')
    table.insert(t, memory.read_u16_be(baseAddr + 0x29a) .. ',')
    table.insert(t, memory.read_u16_be(baseAddr + 0x29c) .. ',')
    table.insert(t, tostring(memory.read_u16_be(baseAddr + 0x29e) > 0) .. ',')
    local inventory = memory.read_u8(baseAddr + 0x166)
    table.insert(t, tostring(bit.check(inventory, 0)) .. ',')
    table.insert(t, tostring(bit.check(inventory, 1)) .. ',')
    table.insert(t, tostring(bit.check(inventory, 2)) .. ',')
    table.insert(t, tostring(bit.check(inventory, 3)) .. ',')
    table.insert(t, tostring(bit.check(inventory, 4)) .. ',')
    table.insert(t, tostring(bit.check(inventory, 5)) .. ',')
    table.insert(t, memory.read_u8(baseAddr + 0x1a2) .. ',')
    table.insert(t, memory.read_u32_be(baseAddr + 0x20c) // 0x3c .. ',')
    table.insert(t, memory.read_u32_be(baseAddr + 0x1fc) .. ',')
    table.insert(t, memory.read_u32_be(baseAddr + 0x200) .. ',')
    table.insert(t, memory.read_u16_be(baseAddr + 0xf0) .. ',')
    table.insert(t, memory.read_u16_be(baseAddr + 0xb4) .. ',')
    table.insert(t, memory.read_u16_be(baseAddr + 0x14a) .. ',')
    table.insert(t, memory.read_u16_be(baseAddr + 0xcc) .. ',')
    table.insert(t, memory.read_u16_be(baseAddr + 0xce) .. ',')
end

local t = {}
table.insert(t, '{')

table.insert(t, '\'warrior\': [')
charValues(0xFF0040, t)
table.insert(t, '],')

table.insert(t, '\'valkyrie\': [')
charValues(0xFF0304, t)
table.insert(t, '],')

table.insert(t, '\'wizard\': [')
charValues(0xFF05C8, t)
table.insert(t, '],')

table.insert(t, '\'elf\': [')
charValues(0xFF088C, t)
table.insert(t, '],')

table.insert(t, '}')

console.log(table.concat(t, ' '))