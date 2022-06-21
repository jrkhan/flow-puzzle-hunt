import NonFungibleToken from "./NonFungibleToken.cdc"
import MetadataViews from "./MetadataViews.cdc"

pub contract FuzzlePieceV2: NonFungibleToken {

    // Standard NFT events
    
    // Event that emitted when the NFT contract is initialized
    //
    pub event ContractInitialized()

    // Event that is emitted when a token is withdrawn,
    // indicating the owner of the collection that it was withdrawn from.
    //
    // If the collection is not in an account's storage, `from` will be `nil`.
    //
    pub event Withdraw(id: UInt64, from: Address?)

    // Event that emitted when a token is deposited to a collection.
    //
    // It indicates the owner of the collection that it was deposited to.
    //
    pub event Deposit(id: UInt64, to: Address?)
    
    pub event Minted(id: UInt64)


    // Named Paths
    pub let MinterStoragePath: StoragePath
    pub let CollectionStoragePath: StoragePath
    pub let CollectionPublicPath: PublicPath


    pub var totalSupply: UInt64

    access(self) var images: {String: String} 

    pub resource interface FuzzlePiecePublic {
        pub fun getDisplayName(): String 
        pub fun getId(): UInt64
        pub fun getPuzzleId(): UInt64
        pub fun getPieceId(): UInt64
    }

    pub resource NFT:  NonFungibleToken.INFT, FuzzlePiecePublic, MetadataViews.Resolver {
        pub let id: UInt64
        pub let displayName: String
        pub let puzzleId: UInt64?
        pub let pieceId: UInt64?

        init(id: UInt64, displayName: String, puzzleId: UInt64, pieceId: UInt64) {
            self.id = id
            self.displayName = displayName
            self.puzzleId = puzzleId
            self.pieceId = pieceId
        }

        pub fun getId(): UInt64 {
            return self.id
        }

        pub fun getPuzzleId(): UInt64 {
            return self.puzzleId ?? 0
        }

        pub fun getPieceId(): UInt64 {
            return self.pieceId ?? 0
        }

        pub fun name(): String {
            return self.displayName
        }

        pub fun description(): String {    
            return self.displayName
        }

        pub fun getDisplayName(): String {
            return self.displayName
        }

        pub fun getViews(): [Type] {
            return [
                Type<MetadataViews.Display>()
            ]
        }

        pub fun imageCID(): String {
            return FuzzlePieceV2.images[self.displayName]!
        }

        pub fun resolveView(_ view: Type): AnyStruct? {
            switch view {
                case Type<MetadataViews.Display>():
                    return MetadataViews.Display(
                        name: self.name(),
                        description: self.description(),
                        thumbnail: MetadataViews.IPFSFile(
                            cid: self.imageCID(), 
                            path: "sm.png"
                        )
                    )
            }
            return nil
        }
    }

    pub resource interface FuzzlePieceCollectionPublic {
        pub fun deposit(token: @NonFungibleToken.NFT)
        pub fun getIDs(): [UInt64]
        pub fun borrowNFT(id: UInt64): &NonFungibleToken.NFT
        pub fun borrowFuzzlePiece(id: UInt64): &AnyResource{FuzzlePieceV2.FuzzlePiecePublic}? {
            // If the result isn't nil, the id of the returned reference
            // should be the same as the argument to the function
            post {
                (result == nil) || (result?.getId() == id):
                    "Cannot borrow FuzzlePiece reference: The ID of the returned reference is incorrect"
            }
        }
    }

    pub resource Collection:  NonFungibleToken.Provider, FuzzlePieceCollectionPublic, NonFungibleToken.Receiver, NonFungibleToken.CollectionPublic, MetadataViews.ResolverCollection {
        pub var ownedNFTs: @{UInt64: NonFungibleToken.NFT}
    
        init() {
            self.ownedNFTs <- {}
        }

        // destructor
        destroy() {
            destroy self.ownedNFTs
        }

        pub fun borrowViewResolver(id: UInt64): &AnyResource{MetadataViews.Resolver} {
            // use an auth reference so we can downcast on the next line
            let nft = &self.ownedNFTs[id] as auth &NonFungibleToken.NFT?
            return nft! as! &FuzzlePieceV2.NFT
        }

        pub fun getIDs(): [UInt64] {
            return self.ownedNFTs.keys
        }

        pub fun deposit(token: @NonFungibleToken.NFT) {
            let token <- token as! @FuzzlePieceV2.NFT
            let id: UInt64 = token.id
            let oldToken <- self.ownedNFTs[id] <- token
            emit Deposit(id: id, to: self.owner?.address)

            destroy oldToken
        }

        pub fun borrowNFT(id: UInt64): &NonFungibleToken.NFT {
            let optNFT = &self.ownedNFTs[id] as &NonFungibleToken.NFT?
            return optNFT!
        }

        pub fun withdraw(withdrawID: UInt64): @NonFungibleToken.NFT {
            let token <- self.ownedNFTs.remove(key: withdrawID) ?? panic("missing NFT")
            emit Withdraw(id: token.id, from: self.owner?.address)

            return <-token
        }    
    
        pub fun borrowFuzzlePiece(id: UInt64): &FuzzlePieceV2.NFT? {
            panic("TODO")
        }
}
    

    pub fun createEmptyCollection(): @NonFungibleToken.Collection {
        return <- create Collection()
    }

    // NFTMinter
    // Resource that an admin or something similar would own to be
    // able to mint new NFTs
    //
    pub resource NFTMinter {

        // mintNFT
        // Mints a new NFT with a new ID
        // and deposit it in the recipients collection using their collection reference
        //
        pub fun mintNFT(
            recipient: &{NonFungibleToken.CollectionPublic}, 
            displayName: String,
            puzzleId: UInt64,
            pieceId: UInt64
        ) {
            // deposit it in the recipient's account using their reference
            recipient.deposit(token: <-create FuzzlePieceV2.NFT(
                id: FuzzlePieceV2.totalSupply, 
                displayName: displayName,
                puzzleId: puzzleId,
                pieceId: pieceId
            ))

            emit Minted(
                id: FuzzlePieceV2.totalSupply
            )

            FuzzlePieceV2.totalSupply = FuzzlePieceV2.totalSupply + (1 as UInt64)
        }

        pub fun setImage(key: String, val: String) {
            FuzzlePieceV2.images.insert(key: key, val)
        }
    }


    init() {
        self.totalSupply = 0
        self.images = {}

        self.CollectionStoragePath = /storage/fuzzlePieceCollectionV2
        self.CollectionPublicPath = /public/fuzzlePieceCollectionV2
        self.MinterStoragePath = /storage/fuzzlePieceMinterV2

        self.account.save(<- create NFTMinter(), to: self.MinterStoragePath)
    }
}

