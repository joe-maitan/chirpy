package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func testMakeAndValidateJWT(t *testing.T) {
	newUUID := uuid.UUID{}
	newJWT, err := MakeJWT(newUUID, "secret!", time.Hour)
	if err != nil {
		t.Errorf("MakeJWT failed: %s", err)
	}

	fetchedUUID, err := ValidateJWT(newJWT, "secret!")
	if err != nil {
		t.Errorf("ValidateJWT failed: %s", err)
	}

	if fetchedUUID != newUUID {
		t.Errorf("testMakeAndValidateJWT - The UUID's are different.")
	}
} // End testCreateJWT() func