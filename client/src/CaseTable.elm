module CaseTable exposing (Model, Msg, init, insertCase, update, view)

import Case
import Dict
import Html exposing (..)
import Html.Attributes exposing (scope)
import Html.Events exposing (onClick)
import Shared exposing (classes)



-- MODEL


{-| Model controls the table with all cases.

    cases :   All cases (dictionary from id to case)
    sortBy :  Marks the sorting column
    sortDir : Sorting direction

-}
type alias Model =
    { cases : Cases
    , sortBy : SortBy
    , sortDir : SortDir
    }


type alias Cases =
    Dict.Dict Int Case.Model


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
        Id
        Asc


someDefaultCases : Cases
someDefaultCases =
    -- TODO: Remove this after the server can provide defaults
    let
        c1 =
            Case.Model "Schulze wg. Diebstahl" "000123/2020" "" "26.04.2020" "" "" Case.Verteidiger "" "laufend"

        c2 =
            Case.Model "Müller M. wg Betrug u. a." "000245/2022" "" "10.10.2020" "" "" Case.Verteidiger "" "laufend"
    in
    Dict.singleton 1 c1 |> insertCase c2



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

        SortCaseTable s ->
            changeSorting model s


insertCase : Case.Model -> Cases -> Cases
insertCase e c =
    let
        newId : Int
        newId =
            case Dict.keys c |> List.maximum of
                Nothing ->
                    1

                Just max ->
                    max + 1
    in
    Dict.insert newId e c


changeSorting : Model -> SortBy -> Model
changeSorting model s =
    if model.sortBy == s then
        case model.sortDir of
            Asc ->
                { model | sortDir = Desc }

            Desc ->
                { model | sortDir = Asc }

    else
        { model | sortBy = s, sortDir = Asc }



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
                (model.cases |> Dict.toList |> sortCases model.sortBy model.sortDir |> List.map caseRow)
            ]
        ]



-- Helpers for header and body of the table follow:


caseListHeader : String -> Model -> SortBy -> Html Msg
caseListHeader txt model sortBy =
    th [ scope "col", onClick <| SortCaseTable sortBy ]
        [ text txt, sortArrows model.sortBy model.sortDir sortBy ]


sortArrows : SortBy -> SortDir -> SortBy -> Html msg
sortArrows current dir field =
    let
        arrows : String
        arrows =
            if current == field then
                case dir of
                    Asc ->
                        "▴ ▿"

                    Desc ->
                        "▵ ▾"

            else
                "▵ ▿"
    in
    span [ classes "float-end pe-5 default-cursor" ] [ text arrows ]


sortCases : SortBy -> SortDir -> List ( Int, Case.Model ) -> List ( Int, Case.Model )
sortCases s d l1 =
    let
        l2 : List ( Int, Case.Model )
        l2 =
            case s of
                Id ->
                    List.sortBy (\n -> Tuple.first n) l1

                Rubrum ->
                    List.sortBy
                        (\n ->
                            let
                                r =
                                    Tuple.second n
                            in
                            r.rubrum
                        )
                        l1

                Beginn ->
                    List.sortBy
                        (\n ->
                            let
                                r =
                                    Tuple.second n
                            in
                            r.beginn
                        )
                        l1

                Ende ->
                    List.sortBy
                        (\n ->
                            let
                                r =
                                    Tuple.second n
                            in
                            r.ende
                        )
                        l1

                Stand ->
                    List.sortBy
                        (\n ->
                            let
                                r =
                                    Tuple.second n
                            in
                            r.stand
                        )
                        l1
    in
    case d of
        Asc ->
            l2

        Desc ->
            List.reverse l2


caseRow : ( Int, Case.Model ) -> Html Msg
caseRow ( id, c ) =
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