import FuzzlePieceV2 from "../contracts/FuzzlePieceV2.cdc"
import NonFungibleToken from "../contracts/NonFungibleToken.cdc"

// This scripts returns the number of FuzzlePieces currently in existence.
pub fun main(address: Address, id: UInt64): &AnyResource{FuzzlePieceV2.FuzzlePiecePublic} {   
    let account = getAccount(address) 
    let collection = account.getCapability(FuzzlePieceV2.CollectionPublicPath)!.borrow<&{FuzzlePieceV2.FuzzlePieceCollectionPublic}>()!
    let nftRef = collection.borrowFuzzlePiece(id: id)

    return nftRef!
}