module Main exposing (main)

import Browser
import Case
import CaseTable
import Html exposing (..)
import Html.Attributes exposing (class, href, type_)
import Html.Events exposing (onClick)
import Http
import Json.Decode as JD
import NewCaseForm
import Shared exposing (classes)


main : Program () Model Msg
main =
    Browser.element
        { init = init
        , view = view
        , update = update
        , subscriptions = subscriptions
        }



-- MODEL


type alias Model =
    { newCaseForm : Maybe NewCaseForm.Model
    , caseTable : CaseTable.Model
    }


init : () -> ( Model, Cmd Msg )
init _ =
    ( Model
        Nothing
        CaseTable.init
    , Http.get
        { url = "/api/case/retrieve"
        , expect = Http.expectJson RetrieveCases (JD.dict Case.caseDecoder) -- TODO: Add decoder for validation that the key is transformable to int
        }
    )



-- UPDATE


type Msg
    = RetrieveCases RetrievedCases
    | OpenNewCaseForm
    | NewCaseFormMsg NewCaseForm.Msg
    | CaseTableMsg CaseTable.Msg


type alias RetrievedCases =
    Result Http.Error CaseTable.CasesFromServer


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        RetrieveCases result ->
            ( { model | caseTable = loadInitialCaseTable result model.caseTable }, Cmd.none )

        OpenNewCaseForm ->
            ( { model | newCaseForm = Just NewCaseForm.init }, Cmd.none )

        NewCaseFormMsg innerMsg ->
            handleNewCaseFormMsg innerMsg model

        CaseTableMsg innerMsg ->
            ( { model | caseTable = CaseTable.update innerMsg model.caseTable }, Cmd.none )


loadInitialCaseTable : RetrievedCases -> CaseTable.Model -> CaseTable.Model
loadInitialCaseTable res m =
    case res of
        Ok cs ->
            { m | cases = CaseTable.insertCases cs m.cases }

        Err _ ->
            -- TODO: Handle the error and show error message or try again
            m


handleNewCaseFormMsg : NewCaseForm.Msg -> Model -> ( Model, Cmd Msg )
handleNewCaseFormMsg msg model =
    case model.newCaseForm of
        Nothing ->
            -- There is no form so ignore the form message.
            ( model, Cmd.none )

        Just f ->
            case NewCaseForm.update msg f of
                NewCaseForm.Updated m cmd ->
                    ( { model | newCaseForm = Just m }, cmd |> Cmd.map NewCaseFormMsg )

                NewCaseForm.Saved id c ->
                    ( { model | newCaseForm = Nothing, caseTable = insertCaseToTable id c model.caseTable }, Cmd.none )

                NewCaseForm.Canceled ->
                    ( { model | newCaseForm = Nothing }, Cmd.none )


insertCaseToTable : Int -> Case.Model -> CaseTable.Model -> CaseTable.Model
insertCaseToTable id e m =
    { m | cases = CaseTable.insertCase id e m.cases }



-- SUBSCRIPTIONS


subscriptions : Model -> Sub Msg
subscriptions _ =
    Sub.none



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
