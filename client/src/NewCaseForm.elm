module NewCaseForm exposing (Action(..), Model, Msg, init, update, view)

import Case
import Html exposing (..)
import Html.Attributes exposing (attribute, class, classList, disabled, for, id, placeholder, rows, selected, type_, value)
import Html.Events exposing (onClick, onInput, onSubmit)
import Http
import Json.Decode as JD
import Shared exposing (classes)



-- MODEL


{-| Model controls the form to create a new case.

    formData :      Contains data from all form fields.
    invalidFields : Holds booleans for all fields that have been filled with
                    invalid data (i. e. fields that must not be empty).
    caseOnSaving : Holds the case that was sent to server while we are waiting
                   for the response. The form is "locked" meanwhile.

-}
type alias Model =
    { formData : FormData
    , invalidFields : InvalidFields
    , caseOnSaving : Maybe Case.Model
    }


type alias FormData =
    { rubrum : String
    , az : String
    , gericht : String
    , beginn : String
    , ende : String
    , gegenstand : String
    , art : Case.Art
    , beschreibung : String
    , stand : String
    }


{-| Only these fields can be invalid.
-}
type alias InvalidFields =
    { rubrum : Bool
    , beginn : Bool
    , stand : Bool
    }


{-| Initializes empty form.
-}
init : Model
init =
    Model
        defaultFormData
        defaultInvalidFields
        Nothing


defaultFormData : FormData
defaultFormData =
    { rubrum = ""
    , az = ""
    , gericht = ""
    , beginn = ""
    , ende = ""
    , gegenstand = ""
    , art = Case.defaultArt
    , beschreibung = ""
    , stand = "laufend"
    }


defaultInvalidFields : InvalidFields
defaultInvalidFields =
    InvalidFields False False False



-- UPDATE


type Msg
    = FormDataMsg FormDataInput
    | Save
    | Cancel
    | FromServer (Result Http.Error Int)


type FormDataInput
    = Rubrum String
    | Az String
    | Gericht String
    | Beginn String
    | Ende String
    | Gegenstand String
    | ArtMsg Case.Art
    | Beschreibung String
    | Stand String


{-| The return value of the update function.
-}
type Action
    = Canceled
    | Saved Int Case.Model
    | Updated Model (Cmd Msg)


{-| Processes the messages of this module and provides also an OutMsg for the
parent.
-}
update : Msg -> Model -> Action
update msg model =
    case msg of
        FormDataMsg m ->
            Updated { model | formData = updateFormData m model.formData } Cmd.none

        Save ->
            save model

        Cancel ->
            Canceled

        FromServer res ->
            case res of
                Ok id ->
                    case model.caseOnSaving of
                        Just c ->
                            Saved id c

                        Nothing ->
                            -- When the "form lock" is gone, just close the form.
                            Canceled

                Err _ ->
                    -- Received error so do not save something and remove the "form lock".
                    -- TODO: Show error message.
                    Updated { model | caseOnSaving = Nothing } Cmd.none


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


{-| If the form is invalid, we just fill the invalidFields property. If the form
is valid, we create the new Case and send it to the server.
-}
save : Model -> Action
save model =
    let
        f : FormData
        f =
            model.formData

        v : InvalidFields
        v =
            formValidate f
    in
    if formIsInvalid v then
        Updated { model | invalidFields = v } Cmd.none

    else
        let
            c : Case.Model
            c =
                Case.Model
                    f.rubrum
                    f.az
                    f.gericht
                    f.beginn
                    f.ende
                    f.gegenstand
                    f.art
                    f.beschreibung
                    f.stand

            cmd : Cmd Msg
            cmd =
                Http.post
                    { url = "/api/case/new"
                    , body = Http.jsonBody (Case.caseEncoder c)
                    , expect = Http.expectJson FromServer (JD.field "id" JD.int)
                    }
        in
        Updated { model | caseOnSaving = Just c } cmd


formValidate : FormData -> InvalidFields
formValidate f =
    InvalidFields
        (f.rubrum == "")
        (f.beginn == "")
        (f.stand == "")


formIsInvalid : InvalidFields -> Bool
formIsInvalid i =
    i.rubrum || i.beginn || i.stand



-- VIEW


{-| Show the form with save and cancel button.
-}
view : Model -> Html Msg
view model =
    div []
        [ form
            [ onSubmit Save, class "mb-5" ]
            [ formfields model.formData model.invalidFields |> map FormDataMsg
            , formButtons model.caseOnSaving
            ]
        , hr [ classes "col-4 mb-5" ] []
        ]


formfields : FormData -> InvalidFields -> Html FormDataInput
formfields formData invalidFields =
    div []
        [ rubrum formData.rubrum invalidFields.rubrum
        , az formData.az
        , gericht formData.gericht
        , beginn formData.beginn invalidFields.beginn
        , ende formData.ende
        , gegenstand formData.gegenstand
        , art formData.art
        , beschreibung formData.beschreibung
        , stand formData.stand invalidFields.stand
        ]


formButtons : Maybe Case.Model -> Html Msg
formButtons caseOnSaving =
    let
        d : Bool
        d =
            case caseOnSaving of
                Just _ ->
                    True

                Nothing ->
                    False
    in
    div []
        [ button [ type_ "submit", classes "btn btn-primary", disabled d ]
            [ text "Speichern" ]
        , button
            [ type_ "button", classes "btn btn-secondary mx-2", disabled d, onClick Cancel ]
            [ text "Abbrechen" ]
        ]



-- All form fields with helper methods follow:


rubrum : Value -> IsInvalid -> Html FormDataInput
rubrum a i =
    inputField
        "text"
        "Rubrum"
        "Erforderliche Angabe. Beispiel: M??ller, M. u. a. wegen Steuerhinterziehung. Dieses Feld wird am Ende der Kammer nicht mitgeteilt."
        i
        Rubrum
        a


az : Value -> Html FormDataInput
az a =
    inputField "text"
        "Kanzleiaktenzeichen und Initialen"
        "Beispiel: 000234/2022 M.M."
        False
        Az
        a


gericht : Value -> Html FormDataInput
gericht a =
    inputField "text"
        "Gericht/Beh??rde und Aktenzeichen"
        "Beispiel: AG Leipzig 123 Cs 456 Js 7890/2022; LG Leipzig ..."
        False
        Gericht
        a


beginn : Value -> IsInvalid -> Html FormDataInput
beginn a i =
    inputField "date"
        "Beginn"
        "Erforderliche Angabe."
        i
        Beginn
        a


ende : Value -> Html FormDataInput
ende a =
    inputField "text"
        "Ende"
        "Datum der Rechtskraft/Mandatsbeendigung oder ???noch anh??ngig???"
        False
        Ende
        a


gegenstand : Value -> Html FormDataInput
gegenstand a =
    inputField "textarea"
        "Gegenstand"
        "Straftatvorwurf und kurzer Abriss des Lebenssachverhalts in zwei bis drei S??tzen"
        False
        Gegenstand
        a


art : Case.Art -> Html FormDataInput
art a =
    let
        idPrefix : String
        idPrefix =
            "NewCaseForm" ++ "Art"
    in
    div [ class "mb-3" ]
        [ label [ for (idPrefix ++ "Select"), class "form-label" ]
            [ text "Art der T??tigkeit" ]
        , select
            [ id (idPrefix ++ "Select")
            , class "form-control"
            , attribute "aria-describedby" (idPrefix ++ "Help")
            , onInput (\value -> Case.stringToArt value |> ArtMsg)
            ]
            [ artOption Case.Verteidiger a
            , artOption Case.Nebenklaeger a
            , artOption Case.Zeugenbeistand a
            , artOption Case.Adhaesionsklaeger a
            ]
        , div [ id (idPrefix ++ "Help"), class "form-text" ]
            [ text "" ]
        ]


artOption : Case.Art -> Case.Art -> Html FormDataInput
artOption a b =
    option [ value <| Case.artToString a, selected (a == b) ]
        [ text <| Case.artToString a ]


beschreibung : Value -> Html FormDataInput
beschreibung a =
    inputField "textarea"
        "Beschreibung der T??tigkeit/Aufteilung der Verfahrensabschnitte/Umfang des Verfahrens"
        "Beispiele: Ermittlungsverfahren/Zwischenverfahren/Hauptverfahren; Haftpr??fungsantrag, Haftbeschwerde, Beweisantr??ge, prozessuale Besonderheiten und prozessuale Reaktion hierauf, Verfahrensabsprache u.a.; au??ergew??hnlicher Aktenumfang, Haftbesuche, Gespr??che mit Staatsanwaltschaft u.a"
        False
        Beschreibung
        a


stand : Value -> IsInvalid -> Html FormDataInput
stand a i =
    inputField "text"
        "Stand des Verfahrens"
        "Erforderliche Angabe. Beispiele: laufend oder abgeschlossen, ggf. Datum der Rechtskraft von Urteilen"
        i
        Stand
        a



-- TODO
-- - Daten der Hauptverhandlungstage (auch vor Straf- bzw. Bu??geldrichter) mit Zuordnung zu dem erkennenden Gericht
-- TODO: Id might be invalid, transform it
-- HTML4: ID and NAME tokens must begin with a letter ([A-Za-z]) and may be followed by any number of letters, digits ([0-9]), hyphens ("-"), underscores ("_"), colons (":"), and periods (".").
-- HTML5: ...???


{-| Creates a Bootstrap form control div with label, input or textarea and help
text. This is a helper for almost all form fields.
-}
inputField :
    InputFieldType
    -> Label
    -> HelpText
    -> IsInvalid
    -> FormDataInputVariant
    -> Value
    -> Html FormDataInput
inputField t l h i toMsg v =
    let
        idPrefix : String
        idPrefix =
            "NewCaseForm" ++ l

        ( tag, attrs ) =
            if t == "textarea" then
                ( textarea, [ rows 5 ] )

            else
                ( input, [ type_ t, classList [ ( "is-invalid", i ) ] ] )
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


type alias IsInvalid =
    Bool


type alias FormDataInputVariant =
    String -> FormDataInput


type alias Value =
    String
