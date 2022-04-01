module NewCaseForm exposing (Model, Msg, defaults, update, view, viewButton)

import Html exposing (..)
import Html.Attributes exposing (attribute, class, for, id, placeholder, rows, selected, type_, value)
import Html.Events exposing (onClick, onInput, onSubmit)
import Shared exposing (classes)



-- MODEL


type alias Model =
    { formOpen : Bool
    , formData : FormData
    }


type alias FormData =
    { rubrum : String
    , az : String
    , gericht : String
    , beginn : String
    , ende : String
    , gegenstand : String
    , art : Art
    , beschreibung : String
    , stand : String
    }


defaults : Model
defaults =
    Model
        False
        defaultFormData


defaultFormData : FormData
defaultFormData =
    { rubrum = ""
    , az = ""
    , gericht = ""
    , beginn = ""
    , ende = ""
    , gegenstand = ""
    , art = Verteidiger
    , beschreibung = ""
    , stand = "laufend"
    }


type Art
    = Verteidiger
    | Nebenklaeger
    | Zeugenbeistand


{-| Converts the custom model type to human readable string.
-}
artToString : Art -> String
artToString a =
    case a of
        Verteidiger ->
            "Verteidiger"

        Nebenklaeger ->
            "Nebenkläger"

        Zeugenbeistand ->
            "Zeugenbeistand"


stringToArt : String -> Art
stringToArt s =
    if s == "Verteidiger" then
        Verteidiger

    else if s == "Nebenklaeger" then
        Nebenklaeger

    else if s == "Zeugenbeistand" then
        Zeugenbeistand

    else
        defaultFormData.art



-- UPDATE


type Msg
    = Form Form
    | SaveNewCase
    | FormDataInput FormDataInput


type Form
    = Open
    | CloseAndReset


type FormDataInput
    = Rubrum String
    | Az String
    | Gericht String
    | Beginn String
    | Ende String
    | Gegenstand String
    | ArtMsg Art
    | Beschreibung String
    | Stand String


update : Msg -> Model -> Model
update msg model =
    case msg of
        Form m ->
            case m of
                Open ->
                    { model | formOpen = True }

                CloseAndReset ->
                    { model | formOpen = False, formData = defaultFormData }

        SaveNewCase ->
            -- TODO: This case is still open.
            model

        FormDataInput m ->
            { model | formData = updateFormData m model.formData }


updateFormData : FormDataInput -> FormData -> FormData
updateFormData msg formData =
    case msg of
        Rubrum v ->
            { formData | rubrum = v }

        Az v ->
            { formData | az = v }

        Gericht v ->
            { formData | gericht = v }

        Beginn v ->
            { formData | beginn = v }

        Ende v ->
            { formData | ende = v }

        Gegenstand v ->
            { formData | gegenstand = v }

        ArtMsg v ->
            { formData | art = v }

        Beschreibung v ->
            { formData | beschreibung = v }

        Stand v ->
            { formData | stand = v }



-- VIEW


view : Model -> Html Msg
view model =
    if model.formOpen then
        viewForm model.formData

    else
        viewButton


viewButton : Html Msg
viewButton =
    button [ type_ "button", classes "btn btn-primary btn-lg px-4 mb-4", onClick <| Form Open ]
        [ text "Neuer Fall" ]


viewForm : FormData -> Html Msg
viewForm formData =
    form [ onSubmit SaveNewCase ]
        [ formfields formData |> map FormDataInput
        , formButtons <| Form CloseAndReset
        ]


formfields : FormData -> Html FormDataInput
formfields formData =
    div []
        [ rubrum formData.rubrum
        , az formData.az
        , gericht formData.gericht
        , beginn formData.beginn
        , ende formData.ende
        , gegenstand formData.gegenstand
        , art formData.art
        , beschreibung formData.beschreibung
        , stand formData.stand
        ]


formButtons : msg -> Html msg
formButtons cancelMsg =
    div []
        [ button [ type_ "submit", classes "btn btn-primary" ]
            [ text "Speichern" ]
        , button
            [ type_ "button", classes "btn btn-secondary mx-2", onClick cancelMsg ]
            [ text "Abbrechen" ]
        ]


rubrum : Value -> Html FormDataInput
rubrum a =
    inputField
        "text"
        "Rubrum"
        "Beispiel: Müller, M. u. a. wegen Steuerhinterziehung. Dieses Feld wird am Ende der Kammer nicht mitgeteilt."
        Rubrum
        a


az : Value -> Html FormDataInput
az a =
    inputField "text"
        "Kanzleiaktenzeichen und Initialen"
        "Beispiel: 000234/2022 M.M."
        Az
        a


gericht : Value -> Html FormDataInput
gericht a =
    inputField "text"
        "Gericht/Behörde und Aktenzeichen"
        "Beispiel: AG Leipzig 123 Cs 456 Js 7890/2022; LG Leipzig ..."
        Gericht
        a


beginn : Value -> Html FormDataInput
beginn a =
    inputField "date"
        "Beginn"
        ""
        Beginn
        a


ende : Value -> Html FormDataInput
ende a =
    inputField "text"
        "Ende"
        "Datum der Rechtskraft/Mandatsbeendigung oder „noch anhängig“"
        Ende
        a


gegenstand : Value -> Html FormDataInput
gegenstand a =
    inputField "textarea"
        "Gegenstand"
        "Straftatvorwurf und kurzer Abriss des Lebenssachverhalts in zwei bis drei Sätzen"
        Gegenstand
        a


art : Art -> Html FormDataInput
art a =
    let
        idPrefix : String
        idPrefix =
            "NewCaseForm" ++ "Art"
    in
    div [ class "mb-3" ]
        [ label [ for (idPrefix ++ "Select"), class "form-label" ]
            [ text "Art der Tätigkeit" ]
        , select
            [ id (idPrefix ++ "Select")
            , class "form-control"
            , attribute "aria-describedby" (idPrefix ++ "Help")
            , onInput (\value -> stringToArt value |> ArtMsg)
            ]
            [ artOption Verteidiger a
            , artOption Nebenklaeger a
            , artOption Zeugenbeistand a
            ]
        , div [ id (idPrefix ++ "Help"), class "form-text" ]
            [ text "" ]
        ]


artOption : Art -> Art -> Html FormDataInput
artOption a b =
    option [ value <| artToString a, selected (a == b) ]
        [ text <| artToString a ]


beschreibung : Value -> Html FormDataInput
beschreibung a =
    inputField "textarea"
        "Beschreibung der Tätigkeit/Aufteilung der Verfahrensabschnitte/Umfang des Verfahrens"
        "Beispiele: Ermittlungsverfahren/Zwischenverfahren/Hauptverfahren; Haftprüfungsantrag, Haftbeschwerde, Beweisanträge, prozessuale Besonderheiten und prozessuale Reaktion hierauf, Verfahrensabsprache u.a.; außergewöhnlicher Aktenumfang, Haftbesuche, Gespräche mit Staatsanwaltschaft u.a"
        Beschreibung
        a


stand : Value -> Html FormDataInput
stand a =
    inputField "text"
        "Stand des Verfahrens"
        "Beispiele: laufend oder abgeschlossen, ggf. Datum der Rechtskraft von Urteilen"
        Stand
        a



-- TODO
-- - Daten der Hauptverhandlungstage (auch vor Straf- bzw. Bußgeldrichter) mit
-- Zuordnung zu dem erkennenden Gericht
-- TODO: Id might be invalid, transform it
-- HTML4: ID and NAME tokens must begin with a letter ([A-Za-z]) and may be followed by any number of letters, digits ([0-9]), hyphens ("-"), underscores ("_"), colons (":"), and periods (".").
-- HTML5: ...???


{-| Creates a Bootstrap form control div with label, input or textarea and help text.
-}
inputField : InputFieldType -> Label -> HelpText -> (String -> FormDataInput) -> Value -> Html FormDataInput
inputField t l h toMsg v =
    let
        idPrefix : String
        idPrefix =
            "NewCaseForm" ++ l

        ( tag, attrs ) =
            if t == "textarea" then
                ( textarea, [ rows 5 ] )

            else
                ( input, [ type_ t ] )
    in
    div [ class "mb-3" ]
        [ label [ for (idPrefix ++ "Input"), class "form-label" ]
            [ text l ]
        , tag
            (attrs
                ++ [ id (idPrefix ++ "Input")
                   , class "form-control"
                   , placeholder l
                   , attribute "aria-describedby" (idPrefix ++ "Help")
                   , onInput toMsg
                   , value v
                   ]
            )
            []
        , div [ id (idPrefix ++ "Help"), class "form-text" ]
            [ text h ]
        ]


type alias InputFieldType =
    String


type alias Label =
    String


type alias HelpText =
    String


type alias Value =
    String
