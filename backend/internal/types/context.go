package types

type ContextKey string

// Auth
const JwtTokenKey ContextKey = "jwt_token"
const UserEmail ContextKey = "user_email"
const UserID ContextKey = "user_id"

// Collection
const CollectionID ContextKey = "collection_id"
const CollectionData ContextKey = "collection_data"

// Card
const CardID ContextKey = "card_id"
const CardData ContextKey = "card_data"
