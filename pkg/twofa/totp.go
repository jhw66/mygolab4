package twofa

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image/png"
	"strings"
	"time"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

type Enrollment struct {
	Secret        string
	Issuer        string
	AccountName   string
	OtpAuthUrl    string
	QrCodeDataUrl string
}

func GenerateEnrollment(issuer, accountName, qrCodeLevel string, qrCodeSize int) (*Enrollment, error) {
	issuer = strings.TrimSpace(issuer)
	accountName = strings.TrimSpace(accountName)
	if issuer == "" || accountName == "" {
		return nil, errors.New("issuer和accountName不能为空")
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: accountName,
	})
	if err != nil {
		return nil, err
	}

	dataURL, err := qrCodeDataUrl(key.URL(), qrCodeLevel, qrCodeSize)
	if err != nil {
		return nil, err
	}

	return &Enrollment{
		Secret:        key.Secret(),
		Issuer:        issuer,
		AccountName:   accountName,
		OtpAuthUrl:    key.URL(),
		QrCodeDataUrl: dataURL,
	}, nil

}
func qrCodeDataUrl(otpAuthUrl, qrCodeLevel string, qrCodeSize int) (string, error) {
	if qrCodeSize <= 0 {
		qrCodeSize = 256
	}

	level, err := parseQRCodeLevel(qrCodeLevel)
	if err != nil {
		return "", err
	}

	code, err := qr.Encode(otpAuthUrl, level, qr.Auto)
	if err != nil {
		return "", err
	}

	code, err = barcode.Scale(code, qrCodeSize, qrCodeSize)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, code); err != nil {
		return "", err
	}

	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func parseQRCodeLevel(level string) (qr.ErrorCorrectionLevel, error) {
	switch strings.ToUpper(strings.TrimSpace(level)) {
	case "", "M":
		return qr.M, nil
	case "L":
		return qr.L, nil
	case "Q":
		return qr.Q, nil
	case "H":
		return qr.H, nil
	default:
		return qr.M, errors.New("invalid qr code level")
	}
}

func ValidateCode(secret, code string) bool {
	code = strings.TrimSpace(code)
	secret = strings.TrimSpace(secret)
	if secret == "" || code == "" {
		return false
	}

	valid, err := totp.ValidateCustom(code, secret, time.Now(), totp.ValidateOpts{
		Period:    30,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})
	return err == nil && valid
}
