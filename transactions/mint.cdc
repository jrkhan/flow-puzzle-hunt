import NonFungibleToken from "../contracts/NonFungibleToken.cdc"
import FuzzlePieceV2 from "../contracts/FuzzlePieceV2.cdc"

// This transction uses the NFTMinter resource to mint a new NFT.
//
// It must be run with the account that has the minter resource
// stored at path /storage/NFTMinter.

transaction(recipient: Address, displayName: String, puzzleId: UInt64, pieceId: UInt64) {

    // local variable for storing the minter reference
    let minter: &FuzzlePieceV2.NFTMinter

    prepare(signer: AuthAccount) {

        // borrow a reference to the NFTMinter resource in storage
        self.minter = signer.borrow<&FuzzlePieceV2.NFTMinter>(from: FuzzlePieceV2.MinterStoragePath)
            ?? panic("Could not borrow a reference to the NFT minter")
    }

    execute {
        // get the public account object for the recipient
        let recipient = getAccount(recipient)

        // borrow the recipient's public NFT collection reference
        let receiver = recipient
            .getCapability(FuzzlePieceV2.CollectionPublicPath)!
            .borrow<&{NonFungibleToken.CollectionPublic}>()
            ?? panic("Could not get receiver reference to the NFT Collection")

        // mint the NFT and deposit it to the recipient's collection
        self.minter.mintNFT(
            recipient: receiver,
            displayName: displayName,
            puzzleId: puzzleId,
            pieceId: pieceId
        )
    }
}