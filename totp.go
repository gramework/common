package common

import (
	"database/sql/driver"
	"encoding/base32"
	"encoding/json"
	"fmt"
	"image"
	"sync"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"

	otp "github.com/hgfischer/go-otp"
)

type totpLength uint8

const (
	// TOTPLengthShort is the default TOTP password length.
	TOTPLengthShort = 6
	// TOTPLengthLong is the strongest length, but it will *NOT* work with Google Authenticator
	TOTPLengthLong = 8
)

// TOTPState handles database/sql interface implementation,
// recovery password generation and validation
type TOTPState struct {
	valid     bool
	mu        *sync.RWMutex
	generator *otp.TOTP
	state     *totpInternalState
}

type totpInternalState struct {
	Secret            string   `json:"s"`
	RecoveryPasswords []string `json:"rp"`
	Label             string   `json:"l"`
	Issuer            string   `json:"i"`
	User              string   `json:"u"`
}

func NewTOTP(label, issuer, user string, length ...totpLength) *TOTPState {
	totp := &TOTPState{}
	totp.mu = &sync.RWMutex{}
	totp.state = &totpInternalState{
		Label:  label,
		Issuer: issuer,
		User:   user,
	}
	for {
		secretBytes, err := GenerateKey()
		if err == nil {
			totp.state.Secret = base32.StdEncoding.EncodeToString(secretBytes)
			break
		}
	}

	for i := 0; i < 12; i++ {
		for {
			rawRP, err := GenerateKey()
			strRP := fmt.Sprintf("%x", rawRP)
			if err == nil && len(strRP) > 24 {
				totp.state.RecoveryPasswords = append(totp.state.RecoveryPasswords, strRP[:24])
				break
			}
		}
	}

	totp.generator = &otp.TOTP{
		Secret:         totp.state.Secret,
		IsBase32Secret: true,
	}
	totp.valid = true
	return totp
}

func (totp *TOTPState) url() string {
	// see https://github.com/google/google-authenticator/wiki/Key-Uri-Format
	return fmt.Sprintf(
		"otpauth://totp/%s:%s?secret=%s&issuer=%s",
		totp.state.Label,
		totp.state.User,
		totp.state.Secret,
		totp.state.Issuer,
	)
}

func (totp *TOTPState) unscaledQR() (barcode.Barcode, error) {
	b, err := qr.Encode(totp.url(), qr.H, qr.Auto)
	return b, err
}

func (totp *TOTPState) QR() (image.Image, error) {
	b, err := totp.unscaledQR()
	if err != nil {
		return nil, err
	}

	b, err = barcode.Scale(b, 300, 300)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (totp *TOTPState) QRScaled(width, height uint16) (image.Image, error) {
	b, err := totp.unscaledQR()
	if err != nil {
		return nil, err
	}

	b, err = barcode.Scale(b, int(width), int(height))
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (totp *TOTPState) IsValidState() bool {
	return totp.valid
}

func (totp *TOTPState) Token() string {
	return totp.generator.Now().Get()
}

func (totp *TOTPState) Verify(token string) bool {
	return totp.generator.Verify(token)
}

func (totp *TOTPState) GetRecoveryPasswords() []string {
	totp.mu.RLock()
	rp := make([]string, len(totp.state.RecoveryPasswords))
	copy(rp, totp.state.RecoveryPasswords)
	totp.mu.RUnlock()
	return rp
}

func (totp *TOTPState) InvalidateRecoveryPassword(rp string) (wasValid bool) {
	totp.mu.Lock()
	for i, pw := range totp.state.RecoveryPasswords {
		if rp == pw {
			totp.state.RecoveryPasswords = append(totp.state.RecoveryPasswords[:i], totp.state.RecoveryPasswords[i+1:]...)
			wasValid = true
			break
		}
	}
	totp.mu.Unlock()
	return
}

// Value implements the driver Valuer interface.
func (totp *TOTPState) Value() (driver.Value, error) {
	totp.mu.RLock()
	b, err := json.Marshal(totp.state)
	totp.mu.RUnlock()
	return b, err
}

// Scan implements the Scanner interface.
func (totp *TOTPState) Scan(value interface{}) error {
	totp.mu.Lock()
	if value == nil {
		totp.valid = false
		totp.mu.Unlock()
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		totp.mu.Unlock()
		return ErrCouldNotScan
	}
	err := json.Unmarshal(b, totp.state)
	if err != nil || len(totp.state.Secret) == 0 {
		totp.valid = false
		log.WithError(err).Error("unmarshaling error")
		totp.mu.Unlock()
		return ErrCouldNotScan
	}
	if len(totp.state.Secret) > 0 {
		totp.valid = true
		totp.generator = &otp.TOTP{
			Secret:         totp.state.Secret,
			IsBase32Secret: true,
		}
	}
	totp.mu.Unlock()
	return nil
}
