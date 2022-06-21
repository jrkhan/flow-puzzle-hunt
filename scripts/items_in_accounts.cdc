import FuzzlePieceV2 from "../contracts/FuzzlePieceV2.cdc"
import NonFungibleToken from "../contracts/NonFungibleToken.cdc"

// This scripts returns the number of FuzzlePieces currently in existence.
pub fun main(address: Address): [UInt64] {   

    let account = getAccount(address) 
    let collection = account.getCapability(FuzzlePieceV2.CollectionPublicPath)!.borrow<&{NonFungibleToken.CollectionPublic}>()!
    return collection.getIDs()
}