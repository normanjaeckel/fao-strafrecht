module Main exposing (main)

import Browser
import Case
import CaseTable
import Html exposing (..)
import Html.Attributes exposing (class, href, type_)
import Html.Events exposing (onClick)
import NewCaseForm
import Shared exposing (classes)


main : Program () Model Msg
main =
    Browser.sandbox
        { init = init
        , view = view
        , update = update
        }



-- MODEL


type alias Model =
    { newCaseForm : Maybe NewCaseForm.Model
    , caseTable : CaseTable.Model
    }


init : Model
init =
    Model
        Nothing
        CaseTable.init



-- UPDATE


type Msg
    = OpenNewCaseForm
    | NewCaseFormMsg NewCaseForm.Msg
    | CaseTableMsg CaseTable.Msg


update : Msg -> Model -> Model
update msg model =
    case msg of
        OpenNewCaseForm ->
            { model | newCaseForm = Just NewCaseForm.init }

        NewCaseFormMsg innerMsg ->
            handleNewCaseFormMsg innerMsg model

        CaseTableMsg innerMsg ->
            handleCaseTableMsg innerMsg model


handleNewCaseFormMsg : NewCaseForm.Msg -> Model -> Model
handleNewCaseFormMsg msg model =
    case model.newCaseForm of
        Nothing ->
            -- There is no form so ignore the form message.
            model

        Just f ->
            case NewCaseForm.update msg f of
                NewCaseForm.Updated m ->
                    { model | newCaseForm = Just m }

                NewCaseForm.Saved c ->
                    { model | newCaseForm = Nothing, caseTable = insertCaseToTable c model.caseTable }

                NewCaseForm.Canceled ->
                    { model | newCaseForm = Nothing }


insertCaseToTable : Case.Model -> CaseTable.Model -> CaseTable.Model
insertCaseToTable e m =
    { m | cases = CaseTable.insertCase e m.cases }


handleCaseTableMsg : CaseTable.Msg -> Model -> Model
handleCaseTableMsg msg model =
    { model | caseTable = CaseTable.update msg model.caseTable }



-- VIEW


view : Model -> Html Msg
view model =
    div [ classes "container p-3 py-md-5" ]
        [ header [ classes "d-flex align-items-center pb-3 mb-5 border-bottom" ]
            [ a [ href "/", classes "d-flex align-items-center text-dark text-decoration-none" ]
                [ span [ class "fs-4" ]
                    [ text <| "Fachanwalt für Strafrecht" ]
                ]
            ]
        , main_ []
            [ h1 [ class "mb-5" ]
                [ text "Meine Fallliste" ]
            , newCaseForm model
            , caseTable model
            ]
        ]


newCaseForm : Model -> Html Msg
newCaseForm model =
    case model.newCaseForm of
        Nothing ->
            button [ type_ "button", classes "btn btn-primary btn-lg px-4 mb-5", onClick <| OpenNewCaseForm ]
                [ text "Neuer Fall" ]

        Just innerModel ->
            NewCaseForm.view innerModel |> map NewCaseFormMsg


caseTable : Model -> Html Msg
caseTable model =
    CaseTable.view model.caseTable |> map CaseTableMsg
