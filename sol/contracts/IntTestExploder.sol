// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract IntTestExploder {
    event Debug(string message);

    // ✅ Should work
    uint256 public u256 = 256;
    uint128 public u128 = 128;
    uint64 public u64 = 64;
    uint32 public u32 = 32;
    uint16 public u16 = 16;
//    // ❓ Mystery zone
    uint8 public u8 = 8;
//
//    // ✅ Should work
    int256 public i256 = -256;
    int128 public i128 = -128;
    int64 public i64 = -64;
    int32 public i32 = -32;
    int16 public i16 = -16;
    // ❓ Mystery zone
    int8 public i8 = -8;

    constructor() {
        emit Debug("Constructor started");

        // Access each value once, to make sure they are initialized and not optimized out
//        uint256 scratch = u256 + u128 + u64 + u32 + u16 + u8;
//        scratch += uint256(int256(i256 + i128 + i64 + i32 + i16 + i8));
//
//        emit Debug("Constructor done");
    }
}
