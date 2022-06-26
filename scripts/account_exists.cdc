import FuzzlePieceV2 from "../contracts/FuzzlePieceV2.cdc"
import NonFungibleToken from "../contracts/NonFungibleToken.cdc"


// This scripts returns the number of FuzzlePieces currently in existence.
pub fun main(address: Address): Bool {   
    let account = getAccount(address) 
    let capability = account.getCapability<&{NonFungibleToken.CollectionPublic}>(FuzzlePieceV2.CollectionPublicPath) 
    return capability.borrow() != nil
}