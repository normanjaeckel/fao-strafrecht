module Case exposing (Art(..), Model, artToString, defaultArt, stringToArt)


type alias Model =
    { rubrum : String
    , az : String
    , gericht : String
    , beginn : String
    , ende : String
    , gegenstand : String
    , art : Art
    , beschreibung : String
    , stand : String
    }


{-| One of the possible roles of the lawyer in a criminal law case.
-}
type Art
    = Verteidiger
    | Nebenklaeger
    | Zeugenbeistand
    | Adhaesionsklaeger


defaultArt : Art
defaultArt =
    Verteidiger


{-| Converts the custom model type to human readable string.
-}
artToString : Art -> String
artToString a =
    case a of
        Verteidiger ->
            "Verteidiger"

        Nebenklaeger ->
            "Nebenkläger"

        Zeugenbeistand ->
            "Zeugenbeistand"

        Adhaesionsklaeger ->
            "Adhäsionskläger"


{-| Converts a string to the corresponding Art value. In case of missmatch it
returns the default value.
-}
stringToArt : String -> Art
stringToArt s =
    if s == "Verteidiger" then
        Verteidiger

    else if s == "Nebenkläger" then
        Nebenklaeger

    else if s == "Zeugenbeistand" then
        Zeugenbeistand

    else if s == "Adhäsionskläger" then
        Adhaesionsklaeger

    else
        defaultArt
