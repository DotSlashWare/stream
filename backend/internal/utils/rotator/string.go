package rotator

import (
	"log"
	"time"
)

var AlphanumericCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// StringRotator cycles through a predefined list of strings at specified intervals.
type StringRotator struct {
	strings       []string
	frequency     time.Duration
	currentString string
	nextRotation  time.Time
}

// NewStringRotator initializes a new StringRotator with the provided list of strings and rotation frequency.
func NewStringRotator(stringList []string, rotationFrequency time.Duration) *StringRotator {
	if len(stringList) == 0 {
		log.Println("Warning: StringRotator initialized with an empty string list. Using fallback value.")
		stringList = []string{"fallback"}
	}
	
	return &StringRotator{
		strings:   stringList,
		frequency: rotationFrequency,
	}
}

// Rotate returns a new string from the rotator's list based on the rotation frequency.
func (rotator *StringRotator) GetString() string {
	if rotator.currentString != "" && time.Now().Before(rotator.nextRotation) {
		return rotator.currentString
	}

	if len(rotator.strings) == 0 {
		log.Println("Warning: StringRotator has an empty string list.")
		return ""
	}

	rotator.currentString = rotator.strings[time.Now().UnixNano()%int64(len(rotator.strings))]
	rotator.nextRotation = time.Now().Add(rotator.frequency)
	return rotator.currentString
}

// RandomStringRotator generates random strings at specified intervals.
type RandomStringRotator struct {
	length        int
	charset       string
	frequency     time.Duration
	currentString string
	nextRotation  time.Time
}

// NewRandomStringRotator initializes a new RandomStringRotator with the specified length, character set, and rotation frequency.
func NewRandomStringRotator(length int, charset string, rotationFrequency time.Duration) *RandomStringRotator {
	if length <= 0 {
		length = 16
		log.Println("Warning: Invalid length for RandomStringRotator. Defaulting to 16.")
	}

	if charset == "" {
		charset = AlphanumericCharset
		log.Println("Warning: Empty charset for RandomStringRotator. Defaulting to alphanumeric characters.")
	}

	return &RandomStringRotator{
		length:    length,
		charset:   charset,
		frequency: rotationFrequency,
	}
}

// Rotate returns a new random string based on the rotator's configuration.
func (rotator *RandomStringRotator) GetString() string {
	if rotator.currentString != "" && time.Now().Before(rotator.nextRotation) {
		return rotator.currentString
	}

	newString := make([]byte, rotator.length)
	for i := range newString {
		newString[i] = rotator.charset[time.Now().UnixNano()%int64(len(rotator.charset))]
	}
	rotator.currentString = string(newString)
	rotator.nextRotation = time.Now().Add(rotator.frequency)
	return rotator.currentString
}
