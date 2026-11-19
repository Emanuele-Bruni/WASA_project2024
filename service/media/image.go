package media

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"path"
	"strings"
	"unicode"
	"unicode/utf8"
)

const MaxBytes = 1500000
const MaxDimension = 4096
const MaxPixels = 4194304

type Upload struct {
	MimeType string `json:"mimeType"`
	Data     string `json:"data"`
	Name     string `json:"name"`
}
type Image struct {
	Data     []byte
	MimeType string
	Name     string
	Width    int
	Height   int
	Digest   string
}

func Decode(upload Upload) (Image, error) {
	if upload.MimeType != "image/png" && upload.MimeType != "image/jpeg" {
		return Image{}, fmt.Errorf("usa immagini PNG o JPEG")
	}
	if len(upload.Data) > base64.StdEncoding.EncodedLen(MaxBytes) {
		return Image{}, fmt.Errorf("l’immagine supera 1500000 byte")
	}
	data, err := base64.StdEncoding.Strict().DecodeString(upload.Data)
	if err != nil || len(data) == 0 || len(data) > MaxBytes {
		return Image{}, fmt.Errorf("dati immagine non validi o troppo grandi")
	}
	config, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil || config.Width < 1 || config.Height < 1 || config.Width > MaxDimension || config.Height > MaxDimension || int64(config.Width)*int64(config.Height) > MaxPixels {
		return Image{}, fmt.Errorf("dimensioni immagine non valide: massimo 4096 per lato e 4194304 pixel")
	}
	expected := map[string]string{"png": "image/png", "jpeg": "image/jpeg"}[format]
	if expected != upload.MimeType {
		return Image{}, fmt.Errorf("il MIME dichiarato non corrisponde al contenuto")
	}
	if _, _, err = image.Decode(bytes.NewReader(data)); err != nil {
		return Image{}, fmt.Errorf("immagine danneggiata o incompleta")
	}
	name := path.Base(strings.ReplaceAll(upload.Name, "\\", "/"))
	if name == "." || name == "" {
		name = "Immagine"
	}
	if !utf8.ValidString(name) || utf8.RuneCountInString(name) > 128 {
		return Image{}, fmt.Errorf("nome immagine troppo lungo o non valido")
	}
	for _, r := range name {
		if unicode.IsControl(r) {
			return Image{}, fmt.Errorf("nome immagine non valido")
		}
	}
	digest := fmt.Sprintf("%x", sha256.Sum256(data))
	return Image{Data: data, MimeType: expected, Name: name, Width: config.Width, Height: config.Height, Digest: digest}, nil
}
