module Case exposing (Art(..), Model, artToString, caseDecoder, caseEncoder, defaultArt, stringToArt)

import Json.Decode as JD
import Json.Decode.Pipeline as JP
import Json.Encode as JE


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


caseDecoder : JD.Decoder Model
caseDecoder =
    JD.succeed Model
        |> JP.required "Rubrum" JD.string
        |> JP.required "Az" JD.string
        |> JP.required "Gericht" JD.string
        |> JP.required "Beginn" JD.string
        |> JP.required "Ende" JD.string
        |> JP.required "Gegenstand" JD.string
        |> JP.required "Art" artDecoder
        |> JP.required "Beschreibung" JD.string
        |> JP.required "Stand" JD.string


artDecoder : JD.Decoder Art
artDecoder =
    JD.string
        |> JD.andThen (\s -> stringToArt s |> JD.succeed)


caseEncoder : Model -> JE.Value
caseEncoder m =
    JE.object
        [ ( "Rubrum", JE.string m.rubrum )
        , ( "Az", JE.string m.az )
        , ( "Gericht", JE.string m.gericht )
        , ( "Beginn", JE.string m.beginn )
        , ( "Ende", JE.string m.ende )
        , ( "Gegenstand", JE.string m.gegenstand )
        , ( "Art", JE.string <| artToString <| m.art )
        , ( "Beschreibung", JE.string m.beschreibung )
        , ( "Stand", JE.string m.stand )
        ]
