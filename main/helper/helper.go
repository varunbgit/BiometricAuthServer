package helper

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/fxamacker/cbor/v2"
)

// Struct to represent the decoded clientDataJSON
type ClientDataJSON struct {
	Challenge string `json:"challenge"`
	Origin    string `json:"origin"`
	Type      string `json:"type"`
}

// Define a struct for CBOR parsing based on the structure of the attestationObject
type AttestationObject struct {
	AuthData []byte                 `cbor:"authData"` // Authentication data
	Fmt      string                 `cbor:"fmt"`      // Attestation format
	AttStmt  map[string]interface{} `cbor:"attStmt"`  // Attestation statement
}

func ParseAttestationObject(decoded []byte) (*AttestationObject, error) {
	var attestation AttestationObject
	err := cbor.Unmarshal(decoded, &attestation)
	if err != nil {
		fmt.Println("error parseAttestationObject:  ", err)
		return nil, fmt.Errorf("failed to parse CBOR attestationObject: %w", err)
	}
	return &attestation, nil
}

func ProcessAuthData(authData []byte) error {
	if len(authData) < 37 {
		return fmt.Errorf("authData is too short")
	}

	// The first 32 bytes are the RP ID hash (SHA-256 hash of the RP ID)
	rpIDHash := authData[:32]
	fmt.Printf("RP ID Hash: %s\n", hex.EncodeToString(rpIDHash))

	// The next byte contains flags (e.g., whether the user is present and whether the attestation data is included)
	flags := authData[32]
	fmt.Printf("Flags: %08b\n", flags)

	// The next 4 bytes are the signature counter
	signatureCounter := authData[33:37]
	fmt.Printf("Signature Counter: %d\n", bytesToUint32(signatureCounter))

	// If the flags indicate attestation data is present, extract it
	if flags&0x40 != 0 {
		// Process the attestation data (e.g., public key, AAGUID, credential ID)
		attestationData := authData[37:]
		fmt.Printf("Attestation Data: %x\n", attestationData)
		// Here you can parse the attestation data further, such as extracting the public key.
	}

	return nil
}

func bytesToUint32(b []byte) uint32 {
	return uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3])
}

// Function to extract the credential ID from the authData
func ExtractCredentialID(authData []byte) (string, error) {
	// Ensure authData is large enough to contain attestation data
	if len(authData) < 37+16+2 {
		return "", fmt.Errorf("authData is too short")
	}

	// Step 1: Skip the first 37 bytes (RP ID hash, flags, and signature counter)
	offset := 37

	// Step 2: Skip the AAGUID (16 bytes)
	offset += 16

	// Step 3: Read the next 2 bytes to get the Credential ID length
	credIDLen := binary.BigEndian.Uint16(authData[offset : offset+2])
	offset += 2

	// Step 4: Extract the Credential ID of length `credIDLen`
	if int(credIDLen) > len(authData)-offset {
		return "", fmt.Errorf("authData too short for credential ID of length %d", credIDLen)
	}
	credentialID := authData[offset : offset+int(credIDLen)]

	// Return the Credential ID as a hex string (or Base64, or any desired format)
	return fmt.Sprintf("%x", credentialID), nil
}
