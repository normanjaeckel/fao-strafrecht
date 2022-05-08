module CaseTable exposing (CasesFromServer, Model, Msg, init, insertCase, insertCases, update, view)

import Case
import Dict
import Html exposing (..)
import Html.Attributes exposing (colspan, scope)
import Html.Events exposing (onClick)
import List exposing (sortBy)
import Shared exposing (classes)
import UpdateCaseForm



-- MODEL


{-| Model controls the table with all cases.

    cases :   All cases (dictionary from id to case)
    sorting :  Marks the sorting column and sorting direction

-}
type alias Model =
    { cases : Cases
    , sorting : Sorting
    , updateCaseForm : Maybe UpdateCaseForm.Model
    }


type alias Cases =
    Dict.Dict Int Case.Model


type alias Sorting =
    { sortBy : SortBy
    , sortDir : SortDir
    }


type SortBy
    = Id
    | Rubrum
    | Beginn
    | Ende
    | Stand


type SortDir
    = Asc
    | Desc


{-| Initializes empty table with default sorting values.
-}
init : Model
init =
    Model
        Dict.empty
        (Sorting Id Asc)
        Nothing



-- UPDATE


type Msg
    = OpenCaseDetail Int
    | SortCases SortBy
    | UpdateCaseFormMsg UpdateCaseForm.Msg


{-| Processes the messages of this module.
-}
update : Msg -> Model -> Model
update msg model =
    case msg of
        OpenCaseDetail id ->
            { model | updateCaseForm = Just <| UpdateCaseForm.init id }

        SortCases innerMsg ->
            case model.updateCaseForm of
                Just _ ->
                    model

                Nothing ->
                    { model | sorting = changeSorting model.sorting innerMsg }

        UpdateCaseFormMsg innerMsg ->
            case model.updateCaseForm of
                Nothing ->
                    model

                Just ucfM ->
                    { model | updateCaseForm = Just <| UpdateCaseForm.update innerMsg ucfM }


insertCase : Int -> Case.Model -> Cases -> Cases
insertCase id e c =
    Dict.insert id e c


type alias CasesFromServer =
    Dict.Dict String Case.Model


insertCases : CasesFromServer -> Cases -> Cases
insertCases elements c =
    let
        foldFn : String -> Case.Model -> Dict.Dict Int Case.Model -> Dict.Dict Int Case.Model
        foldFn =
            \id case_ updatedCases ->
                updatedCases
                    |> Dict.insert
                        (String.toInt id |> Maybe.withDefault 0)
                        case_
    in
    Dict.foldl foldFn c elements


changeSorting : Sorting -> SortBy -> Sorting
changeSorting sorting innerMsg =
    if sorting.sortBy == innerMsg then
        case sorting.sortDir of
            Asc ->
                { sorting | sortDir = Desc }

            Desc ->
                { sorting | sortDir = Asc }

    else
        { sorting | sortBy = innerMsg, sortDir = Asc }



-- VIEW


{-| Shows the table with all messages
-}
view : Model -> Html Msg
view model =
    if not (Dict.isEmpty model.cases) then
        div []
            [ table [ classes "table table-striped" ]
                [ thead []
                    [ tr []
                        [ caseListHeader "#" model Id
                        , caseListHeader "Rubrum" model Rubrum
                        , caseListHeader "Beginn" model Beginn
                        , caseListHeader "Ende" model Ende
                        , caseListHeader "Stand" model Stand
                        , th [ scope "col" ]
                            [ text "HV-Tage" ]
                        ]
                    ]
                , tbody []
                    (sortCases model.cases model.sorting |> caseRows model.updateCaseForm)
                ]
            ]

    else
        div [] []



-- Helpers for header and body of the table follow:


caseListHeader : String -> Model -> SortBy -> Html Msg
caseListHeader txt model sortBy =
    th [ scope "col", onClick <| SortCases sortBy ]
        [ text txt, sortArrows model.sorting sortBy ]


sortArrows : Sorting -> SortBy -> Html msg
sortArrows s field =
    let
        arrows : String
        arrows =
            if s.sortBy == field then
                case s.sortDir of
                    Asc ->
                        "▴ ▿"

                    Desc ->
                        "▵ ▾"

            else
                "▵ ▿"
    in
    span [ classes "float-end pe-5 default-cursor" ] [ text arrows ]


type alias SortedCases =
    List ( Int, Case.Model )


sortCases : Cases -> Sorting -> SortedCases
sortCases cases s =
    let
        sort : Cases -> List ( Int, Case.Model )
        sort =
            case s.sortBy of
                Id ->
                    sortById

                Rubrum ->
                    sortByStringField .rubrum

                Beginn ->
                    sortByStringField .beginn

                Ende ->
                    sortByStringField .ende

                Stand ->
                    sortByStringField .stand
    in
    case s.sortDir of
        Asc ->
            sort cases

        Desc ->
            List.reverse (sort cases)


sortById : Cases -> List ( Int, Case.Model )
sortById c =
    c |> Dict.toList |> List.sortBy (\n -> Tuple.first n)


sortByStringField : (Case.Model -> String) -> Cases -> List ( Int, Case.Model )
sortByStringField fn c =
    let
        sortFn : ( Int, Case.Model ) -> String
        sortFn =
            \elem ->
                Tuple.second elem |> fn
    in
    Dict.toList c |> List.sortBy sortFn


caseRows : Maybe UpdateCaseForm.Model -> SortedCases -> List (Html Msg)
caseRows m s =
    List.map (caseRow m) <| s


caseRow : Maybe UpdateCaseForm.Model -> ( Int, Case.Model ) -> Html Msg
caseRow m t =
    let
        id : Int
        id =
            Tuple.first t

        c : Case.Model
        c =
            Tuple.second t
    in
    case m of
        Nothing ->
            caseLine id c

        Just ucfM ->
            if ucfM.id /= id then
                caseLine id c

            else
                caseForm ucfM


caseLine : Int -> Case.Model -> Html Msg
caseLine id c =
    tr [ onClick <| OpenCaseDetail <| id ]
        [ th [ scope "row" ]
            [ text <| String.fromInt id ]
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


caseForm : UpdateCaseForm.Model -> Html Msg
caseForm model =
    tr []
        [ td [ colspan 6 ]
            [ UpdateCaseForm.view model |> map UpdateCaseFormMsg ]
        ]
