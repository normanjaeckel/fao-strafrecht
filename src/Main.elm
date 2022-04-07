module Main exposing (main)

import Browser
import Case
import Dict
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
    , cases : Cases
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


init : Model
init =
    Model
        Nothing
        someDefaultCases
        Id
        Asc


someDefaultCases : Cases
someDefaultCases =
    let
        c1 =
            Case.Model "Schulze wg Diebstahl" "000123/2020" "" "26.04.2020" "" "" Case.Verteidiger "" "laufend"

        c2 =
            Case.Model "Müller M. wg Betrug u. a." "000245/2022" "" "10.10.2020" "" "" Case.Verteidiger "" "laufend"
    in
    Dict.singleton 1 c1 |> insertCase c2



-- UPDATE


type Msg
    = OpenNewCaseForm
    | NewCaseFormMsg NewCaseForm.Msg
    | OpenCaseDetail
    | SortCaseTable SortBy


update : Msg -> Model -> Model
update msg model =
    case msg of
        OpenNewCaseForm ->
            { model | newCaseForm = Just NewCaseForm.init }

        NewCaseFormMsg innerMsg ->
            handleNewCaseFormMsg innerMsg model

        OpenCaseDetail ->
            -- TODO: Add this case here.
            model

        SortCaseTable s ->
            changeSorting model s


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
                    { model | newCaseForm = Nothing, cases = insertCase c model.cases }

                NewCaseForm.Canceled ->
                    { model | newCaseForm = Nothing }


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
            NewCaseForm.view innerModel |> map NewCaseFormMsg


caseListView : Model -> Html Msg
caseListView model =
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
    span [ classes "float-end pe-5" ] [ text arrows ]


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
