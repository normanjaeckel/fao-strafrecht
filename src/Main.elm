module Main exposing (main)

import Browser
import Case
import Html exposing (..)
import Html.Attributes exposing (class, href, scope)
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


type alias Model =
    { newCase : NewCaseForm.Model
    , cases : List Case.Model
    }


init : Model
init =
    Model
        NewCaseForm.defaults
        []


type Msg
    = NewCaseMsg NewCaseForm.Msg
    | OpenCaseDetail


update : Msg -> Model -> Model
update msg model =
    case msg of
        NewCaseMsg innerMsg ->
            let
                innerModel : NewCaseForm.Model
                innerModel =
                    NewCaseForm.update innerMsg model.newCase
            in
            case innerModel.save of
                Nothing ->
                    { model | newCase = innerModel }

                Just c ->
                    let
                        updatedInnerModel =
                            { innerModel | save = Nothing }
                    in
                    { model | newCase = updatedInnerModel, cases = model.cases ++ [ c ] }

        OpenCaseDetail ->
            -- TODO: Add this case here.
            model


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
            [ h1 [ class "mb-5" ]
                [ text "Meine Fallliste" ]
            , NewCaseForm.view model.newCase |> map NewCaseMsg
            , caseListView model
            ]
        ]


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
