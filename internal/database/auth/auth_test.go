package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestAdd(t *testing.T) {
	given_id, _ := uuid.NewUUID()
	tokenString, _ := MakeJWT(given_id, "testSecret", time.Duration(time.Minute))
	valid_id, err := ValidateJWT(tokenString, "testSecret")
	if given_id != valid_id {
		t.Errorf("IDs do not match: got %v, want %v", valid_id, given_id)
	}

	given_id, _ = uuid.NewUUID()
	tokenString, _ = MakeJWT(given_id, "testSecret", time.Duration(time.Minute))
	valid_id, err = ValidateJWT(tokenString, "testSecretDifferent")
	if err == nil {
		t.Error("Error should not be nil for different tokenSecret strings")
	}

	given_id, _ = uuid.NewUUID()
	tokenString, _ = MakeJWT(given_id, "testSecret", time.Duration(time.Nanosecond))

	time.Sleep(time.Second)
	valid_id, err = ValidateJWT(tokenString, "testSecret")
	if err == nil {
		t.Error("Error should not be nil for expiered tokens")
	}
}
