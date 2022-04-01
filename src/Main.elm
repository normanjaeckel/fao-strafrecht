module Main exposing (main)

import Browser
import Html exposing (..)
import Html.Attributes exposing (class, href)
import NewCaseForm
import Shared exposing (classes)


main : Program () Model Msg
main =
    Browser.sandbox
        { init = init
        , view = view
        , update = update
        }


type alias Model =
    { newCase : NewCaseForm.Model
    , foo : String
    }


init : Model
init =
    Model NewCaseForm.defaults ""


type Msg
    = NewCaseMsg NewCaseForm.Msg


update : Msg -> Model -> Model
update msg model =
    case msg of
        NewCaseMsg innerMsg ->
            { model | newCase = NewCaseForm.update innerMsg model.newCase }


view : Model -> Html Msg
view model =
    div [ classes "container p-3 py-md-5" ]
        [ header [ classes "d-flex align-items-center pb-3 mb-5 border-bottom" ]
            [ a [ href "/", classes "d-flex align-items-center text-dark text-decoration-none" ]
                [ span [ class "fs-4" ]
                    [ text <| "Fachanwalt für Strafrecht" ++ model.newCase.formData.rubrum ]
                ]
            ]
        , main_ []
            [ h1 [ class "mb-4" ]
                [ text "Meine Fallliste" ]
            , NewCaseForm.view model.newCase |> map NewCaseMsg
            ]
        ]
