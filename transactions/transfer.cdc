import FuzzlePieceV2 from "../contracts/FuzzlePieceV2.cdc"
import NonFungibleToken from "../contracts/NonFungibleToken.cdc"

transaction(recipientAddr: Address, id: UInt64) {

     // local variable for storing the minter reference
    let sender: &FuzzlePieceV2.Collection
    let targetAccount: &AnyResource{NonFungibleToken.CollectionPublic}
    prepare(signer: AuthAccount) {

        // get the recipients public account object
        let recipient = getAccount(recipientAddr)

        // borrow a reference to the NFTMinter resource in storage
        self.sender = signer.borrow<&FuzzlePieceV2.Collection>(from: FuzzlePieceV2.CollectionStoragePath)
            ?? panic("Could not borrow a reference to the NFT minter")

        self.targetAccount = recipient.getCapability(FuzzlePieceV2.CollectionPublicPath)!.borrow<&{NonFungibleToken.CollectionPublic}>()!
    }

    execute { 
        let nft <- self.sender.withdraw(withdrawID: id)
        self.targetAccount.deposit(token: <-nft)
    }

}