module Main exposing (main)

import Browser
import Case
import Html exposing (..)
import Html.Attributes exposing (class, href, scope, type_)
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
    , cases : List Case.Model
    }


init : Model
init =
    Model
        Nothing
        []



-- UPDATE


type Msg
    = OpenNewCaseForm
    | NewCaseMsg NewCaseForm.Msg
    | OpenCaseDetail


update : Msg -> Model -> Model
update msg model =
    case msg of
        OpenNewCaseForm ->
            { model | newCaseForm = Just NewCaseForm.init }

        NewCaseMsg innerMsg ->
            let
                innerModel : NewCaseForm.Model
                innerModel =
                    case model.newCaseForm of
                        Just v ->
                            NewCaseForm.update innerMsg v

                        Nothing ->
                            NewCaseForm.update innerMsg NewCaseForm.init
            in
            case innerModel.status of
                -- TODO: Try not to check inner status here.
                NewCaseForm.Open ->
                    { model | newCaseForm = Just innerModel }

                NewCaseForm.Saved c ->
                    { model | newCaseForm = Nothing, cases = model.cases ++ [ c ] }

                NewCaseForm.Canceled ->
                    { model | newCaseForm = Nothing }

        OpenCaseDetail ->
            -- TODO: Add this case here.
            model



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
            , caseListView model
            ]
        ]


newCaseForm : Model -> Html Msg
newCaseForm model =
    case model.newCaseForm of
        Nothing ->
            button [ type_ "button", classes "btn btn-primary btn-lg px-4 mb-5", onClick <| OpenNewCaseForm ]
                [ text "Neuer Fall" ]

        Just innerModel ->
            NewCaseForm.view innerModel |> map NewCaseMsg


caseListView : Model -> Html Msg
caseListView model =
    div []
        [ table [ classes "table table-striped" ]
            [ thead []
                [ tr []
                    [ th [ scope "col" ]
                        [ text "#" ]
                    , th [ scope "col" ]
                        [ text "Rubrum" ]
                    , th [ scope "col" ]
                        [ text "Beginn" ]
                    , th [ scope "col" ]
                        [ text "Ende" ]
                    , th [ scope "col" ]
                        [ text "Stand" ]
                    , th [ scope "col" ]
                        [ text "HV-Tage" ]
                    ]
                ]
            , tbody []
                (model.cases |> List.map caseRow)
            ]
        ]


caseRow : Case.Model -> Html Msg
caseRow c =
    tr [ onClick OpenCaseDetail ]
        [ th [ scope "row" ]
            [ text <| String.fromInt c.number ]
        , td []
            [ text c.rubrum ]
        , td []
            [ text c.beginn ]
        , td []
            [ text c.ende ]
        , td []
            [ text c.stand ]
        , td []
            [ text "" ]
        ]
