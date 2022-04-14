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
    = Cases (Dict.Dict Int Case.Model) Expanded


type Expanded
    = Expanded (Maybe Int)


getExp : Cases -> Maybe Int
getExp (Cases _ (Expanded exp)) =
    exp


isExpanded : Cases -> Bool
isExpanded c =
    case getExp c of
        Just _ ->
            True

        Nothing ->
            False


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
    Cases (Dict.singleton 1 c1) (Expanded Nothing)
        |> insertCase c2
        |> insertCase c3



-- UPDATE


type Msg
    = OpenCaseDetail (Maybe Int)
    | SortCases SortBy


{-| Processes the messages of this module.
-}
update : Msg -> Model -> Model
update msg model =
    case msg of
        OpenCaseDetail innerMsg ->
            { model | cases = changedDetail model.cases innerMsg }

        SortCases innerMsg ->
            if isExpanded model.cases then
                model

            else
                { model | sorting = changeSorting model.sorting innerMsg }


insertCase : Case.Model -> Cases -> Cases
insertCase e (Cases c exp) =
    let
        newId : Int
        newId =
            case Dict.keys c |> List.maximum of
                Nothing ->
                    1

                Just max ->
                    max + 1
    in
    Cases
        (Dict.insert newId e c)
        exp


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


changedDetail : Cases -> Maybe Int -> Cases
changedDetail (Cases c (Expanded exp)) id =
    -- We have five cases. Test them.
    -- N N -> N
    -- N 42 -> 42
    -- 10 N -> 10
    -- 10 10 -> N
    -- 10 42 -> 42
    let
        value =
            case exp of
                Nothing ->
                    id

                Just currentId ->
                    case id of
                        Nothing ->
                            id

                        Just i ->
                            if currentId == i then
                                Nothing

                            else
                                id
    in
    value |> Expanded |> Cases c



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
                (sortCases model.cases model.sorting |> caseRows)
            ]
        ]



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


type SortedCases
    = SortedCases (List ( Int, Case.Model )) Expanded


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

        result : List ( Int, Case.Model )
        result =
            case s.sortDir of
                Asc ->
                    sort cases

                Desc ->
                    List.reverse (sort cases)
    in
    SortedCases result <| Expanded <| getExp cases


sortById : Cases -> List ( Int, Case.Model )
sortById (Cases c _) =
    c |> Dict.toList |> List.sortBy (\n -> Tuple.first n)


sortByStringField : (Case.Model -> String) -> Cases -> List ( Int, Case.Model )
sortByStringField fn (Cases cases _) =
    let
        sortFn : ( Int, Case.Model ) -> String
        sortFn =
            \elem ->
                Tuple.second elem |> fn
    in
    Dict.toList cases |> List.sortBy sortFn


caseRows : SortedCases -> List (Html Msg)
caseRows (SortedCases s exp) =
    List.map (caseRow exp) <| s


caseRow : Expanded -> ( Int, Case.Model ) -> Html Msg
caseRow (Expanded exp) t =
    let
        id : Int
        id =
            Tuple.first t

        c : Case.Model
        c =
            Tuple.second t
    in
    case exp of
        Nothing ->
            caseLine id c

        Just i ->
            if i /= id then
                caseLine id c

            else
                caseForm id c


caseLine : Int -> Case.Model -> Html Msg
caseLine id c =
    tr [ onClick <| OpenCaseDetail <| Just id ]
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


caseForm : Int -> Case.Model -> Html Msg
caseForm _ _ =
    div [ onClick <| OpenCaseDetail <| Nothing ]
        [ text "Byebye" ]
