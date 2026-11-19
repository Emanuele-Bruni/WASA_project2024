package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

const maxBodyBytes = 2 * 1024 * 1024
const maxTextBytes = 2000000

// readJSON accepts one UTF-8 JSON object with exact, unique field names.
func readJSON(w http.ResponseWriter, r *http.Request, destination interface{}, optional bool) bool {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			writeError(w, 413, "Il corpo della richiesta supera 2 MiB")
		} else {
			writeError(w, 400, "Impossibile leggere la richiesta")
		}
		return false
	}
	if len(bytes.TrimSpace(data)) == 0 && optional {
		return true
	}
	media, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || media != "application/json" || (params["charset"] != "" && !strings.EqualFold(params["charset"], "utf-8")) {
		writeError(w, 415, "Usa Content-Type: application/json con codifica UTF-8")
		return false
	}
	if !utf8.Valid(data) {
		writeError(w, 400, "Il JSON deve essere UTF-8 valido")
		return false
	}
	if err := validateObjectFields(data, reflect.TypeOf(destination).Elem()); err != nil {
		writeError(w, 400, "Oggetto JSON non valido: campi sconosciuti, ripetuti, null o documenti aggiuntivi")
		return false
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err = decoder.Decode(destination); err != nil {
		writeError(w, 400, "Tipi dei campi JSON non validi")
		return false
	}

	return true
}

func noRequestBody(w http.ResponseWriter, r *http.Request) bool {
	return readJSON(w, r, &struct{}{}, true)
}

func validateName(w http.ResponseWriter, value *string, field string, min, max int) bool {
	for _, r := range *value {
		if unicode.IsControl(r) || r == '\u2028' || r == '\u2029' {
			writeError(w, 400, "Il nome non può contenere caratteri di controllo", field)
			return false
		}
	}
	*value = strings.TrimSpace(*value)
	if length := utf8.RuneCountInString(*value); length < min || length > max {
		writeError(w, 400, "Lunghezza del nome non valida", field)
		return false
	}
	return true
}
func validateMessage(w http.ResponseWriter, value string) bool {
	if strings.TrimSpace(value) == "" || len(value) > maxTextBytes {
		writeError(w, 400, "Il messaggio deve contenere testo e non superare 2000000 byte", "text")
		return false
	}
	for _, r := range value {
		if unicode.IsControl(r) && r != '\n' && r != '\r' && r != '\t' {
			writeError(w, 400, "Il messaggio contiene caratteri di controllo non validi", "text")
			return false
		}
	}
	return true
}

var emojiPattern = regexp.MustCompile(`^(?:[\x{1F1E6}-\x{1F1FF}]{2}|[0-9#*]\x{FE0F}?\x{20E3}|\p{So}(?:\x{FE0F}|[\x{1F3FB}-\x{1F3FF}]){0,2}(?:\x{200D}\p{So}(?:\x{FE0F}|[\x{1F3FB}-\x{1F3FF}]){0,2})*)$`)

func validateEmoji(w http.ResponseWriter, value string) bool {
	if utf8.RuneCountInString(value) > 16 || !emojiPattern.MatchString(value) {
		writeError(w, 400, "Reazione non valida: usa un simbolo o una sequenza emoji supportata", "emoji")
		return false
	}
	return true
}

func validateObjectFields(data []byte, kind reflect.Type) error {
	for kind.Kind() == reflect.Ptr {
		kind = kind.Elem()
	}
	allowed := make(map[string]reflect.Type)
	for i := 0; i < kind.NumField(); i++ {
		allowed[strings.Split(kind.Field(i).Tag.Get("json"), ",")[0]] = kind.Field(i).Type
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	token, err := decoder.Token()
	if err != nil || token != json.Delim('{') {
		return errors.New("object required")
	}
	seen := make(map[string]bool)
	for decoder.More() {
		key, err := decoder.Token()
		if err != nil {
			return err
		}
		name, ok := key.(string)
		field, exists := allowed[name]
		if !ok || !exists || seen[name] {
			return errors.New("unknown or repeated field")
		}
		seen[name] = true
		var value json.RawMessage
		if err = decoder.Decode(&value); err != nil || bytes.Equal(value, []byte("null")) {
			return errors.New("invalid value")
		}
		for field.Kind() == reflect.Ptr {
			field = field.Elem()
		}
		if field.Kind() == reflect.Struct {
			if err = validateObjectFields(value, field); err != nil {
				return err
			}
		}
	}
	if _, err = decoder.Token(); err != nil {
		return err
	}
	var extra interface{}
	if err = decoder.Decode(&extra); err != io.EOF {
		return errors.New("additional JSON")
	}
	return nil
}
