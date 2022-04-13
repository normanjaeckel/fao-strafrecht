module CaseTable exposing (Model, Msg, init, insertCase, update, view)

import Case
import Dict
import Html exposing (..)
import Html.Attributes exposing (scope)
import Html.Events exposing (onClick)
import List exposing (sortBy)
import Shared exposing (classes)



-- MODEL


{-| Model controls the table with all cases.

    cases :   All cases (dictionary from id to case)
    sortBy :  Marks the sorting column
    sortDir : Sorting direction

-}
type alias Model =
    { cases : Cases
    , sorting : Sorting
    }


type Cases
    = Cases (Dict.Dict Int Case.Model)


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
        someDefaultCases
        (Sorting Id Asc)


someDefaultCases : Cases
someDefaultCases =
    -- TODO: Remove this after the server can provide defaults
    let
        c1 =
            Case.Model "Schulze wg. Diebstahl" "000123/2020" "" "26.04.2020" "" "" Case.Verteidiger "" "laufend"

        c2 =
            Case.Model "Maller M. wg Betrug u. a." "000245/2022" "" "10.10.2020" "" "" Case.Nebenklaeger "" "laufend"

        c3 =
            Case.Model "Meier wg. Steuerhinterziehung" "000333/2022" "" "11.10.2020" "" "" Case.Verteidiger "" "laufend"
    in
    Cases (Dict.singleton 1 c1)
        |> insertCase c2
        |> insertCase c3



-- UPDATE


type Msg
    = OpenCaseDetail
    | SortCaseTable SortBy


{-| Processes the messages of this module.
-}
update : Msg -> Model -> Model
update msg model =
    case msg of
        OpenCaseDetail ->
            -- TODO: Add this case here.
            model

        SortCaseTable innerMsg ->
            { model | sorting = changeSorting model.sorting innerMsg }


insertCase : Case.Model -> Cases -> Cases
insertCase e (Cases c) =
    let
        newId : Int
        newId =
            case Dict.keys c |> List.maximum of
                Nothing ->
                    1

                Just max ->
                    max + 1
    in
    Dict.insert newId e c |> Cases


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
                (sortCases model.cases model.sorting |> caseRow)
            ]
        ]



-- Helpers for header and body of the table follow:


caseListHeader : String -> Model -> SortBy -> Html Msg
caseListHeader txt model sortBy =
    th [ scope "col", onClick <| SortCaseTable sortBy ]
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


type SortedCases
    = SortedCases (List ( Int, Case.Model ))


sortCases : Cases -> Sorting -> SortedCases
sortCases (Cases cases) s =
    let
        l1 : List ( Int, Case.Model )
        l1 =
            Dict.toList cases

        l2 : List ( Int, Case.Model )
        l2 =
            case s.sortBy of
                Id ->
                    List.sortBy (\n -> Tuple.first n) l1

                Rubrum ->
                    List.sortBy
                        (\n ->
                            let
                                c =
                                    Tuple.second n
                            in
                            c.rubrum
                        )
                        l1

                Beginn ->
                    List.sortBy
                        (\n ->
                            let
                                c =
                                    Tuple.second n
                            in
                            c.beginn
                        )
                        l1

                Ende ->
                    List.sortBy
                        (\n ->
                            let
                                c =
                                    Tuple.second n
                            in
                            c.ende
                        )
                        l1

                Stand ->
                    List.sortBy
                        (\n ->
                            let
                                c =
                                    Tuple.second n
                            in
                            c.stand
                        )
                        l1

        l3 : List ( Int, Case.Model )
        l3 =
            case s.sortDir of
                Asc ->
                    l2

                Desc ->
                    List.reverse l2
    in
    SortedCases l3


caseRow : SortedCases -> List (Html Msg)
caseRow (SortedCases s) =
    let
        fn : ( Int, Case.Model ) -> Html Msg
        fn =
            \elem ->
                let
                    id : Int
                    id =
                        Tuple.first elem

                    c : Case.Model
                    c =
                        Tuple.second elem
                in
                tr [ onClick OpenCaseDetail ]
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
    in
    List.map fn s
