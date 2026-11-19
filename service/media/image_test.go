package media

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
	"strings"
	"testing"
)

func imageUpload(width, height int) Upload {
	var buffer bytes.Buffer
	_ = png.Encode(&buffer, image.NewRGBA(image.Rect(0, 0, width, height)))
	return Upload{MimeType: "image/png", Data: base64.StdEncoding.EncodeToString(buffer.Bytes()), Name: "photo.png"}
}
func TestValidateActualImageContentAndBounds(t *testing.T) {
	valid := imageUpload(2, 3)
	parsed, err := Decode(valid)
	if err != nil || parsed.Width != 2 || parsed.Height != 3 || parsed.Digest == "" {
		t.Fatalf("valid image rejected: %v", err)
	}
	cases := []Upload{valid, valid, valid, imageUpload(4097, 1), imageUpload(2049, 2049), valid, valid}
	cases[0].MimeType = "image/jpeg"
	cases[1].Data = "not base64"
	cases[2].Data = base64.StdEncoding.EncodeToString([]byte("<svg/>"))
	cases[5].Data = base64.StdEncoding.EncodeToString(make([]byte, MaxBytes+1))
	cases[6].Name = strings.Repeat("a", 129)
	for i, upload := range cases {
		if _, err = Decode(upload); err == nil {
			t.Fatalf("invalid image %d accepted", i)
		}
	}
	data, _ := base64.StdEncoding.DecodeString(valid.Data)
	valid.Data = base64.StdEncoding.EncodeToString(data[:len(data)/2])
	if _, err = Decode(valid); err == nil {
		t.Fatal("truncated PNG accepted")
	}
}
