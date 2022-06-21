import FuzzlePieceV2 from "../contracts/FuzzlePieceV2.cdc"

// This scripts returns the number of currently in existence.
pub fun main(): UInt64 {    
    return FuzzlePieceV2.totalSupply
}