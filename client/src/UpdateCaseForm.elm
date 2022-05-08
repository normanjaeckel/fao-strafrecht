module UpdateCaseForm exposing (Model, Msg, init, update, view)

import Case
import Html exposing (..)
import Html.Attributes exposing (disabled, type_)
import Html.Events exposing (onClick, onSubmit)
import Shared exposing (classes)



-- MODEL


type alias Model =
    { id : Int
    }


init : Int -> Model
init id =
    { id = id }



-- UPDATE


type Msg
    = FormDataMsg FormDataInput
    | Save
    | Cancel


update : Msg -> Model -> Model
update _ m =
    m



--| FromServer (Result Http.Error Int)


type FormDataInput
    = Rubrum String
    | Az String
    | Gericht String
    | Beginn String
    | Ende String
    | Gegenstand String
    | ArtMsg Case.Art
    | Beschreibung String
    | Stand String



-- VIEW


view : Model -> Html Msg
view _ =
    let
        d : Bool
        d =
            False
    in
    form [ onSubmit Save ]
        [ button [ type_ "submit", classes "btn btn-primary", disabled d ]
            [ text "Speichern" ]
        , button
            [ type_ "button", classes "btn btn-secondary mx-2", disabled d, onClick Cancel ]
            [ text "Abbrechen" ]
        ]
