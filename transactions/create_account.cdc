import NonFungibleToken from "../contracts/NonFungibleToken.cdc"
import FuzzlePieceV2 from "../contracts/FuzzlePieceV2.cdc"
import MetadataViews from "../contracts/MetadataViews.cdc"

// This transaction configures an account to hold Kitty Items.

transaction {
    prepare(signer: AuthAccount) {
        // if the account doesn't already have a collection
        if signer.borrow<&FuzzlePieceV2.Collection>(from: FuzzlePieceV2.CollectionStoragePath) == nil {

            // create a new empty collection
            let collection <- FuzzlePieceV2.createEmptyCollection()
            
            // save it to the account
            signer.save(<-collection, to: FuzzlePieceV2.CollectionStoragePath)

            // create a public capability for the collection
            signer.link<&FuzzlePieceV2.Collection{NonFungibleToken.CollectionPublic, FuzzlePieceV2.FuzzlePieceCollectionPublic, MetadataViews.ResolverCollection}>(FuzzlePieceV2.CollectionPublicPath, target: FuzzlePieceV2.CollectionStoragePath)
        }
    }
}