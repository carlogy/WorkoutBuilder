package auth

import (
	"math/rand"
	"testing"
	"time"

	"github.com/carlogy/WorkoutBuilder/internal/auth"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()<>,.;/?=+-_"

func generateRandomString(length int) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	bSlice := make([]byte, length)
	for i := range bSlice {
		bSlice[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(bSlice)
}

func TestAuthHashPasswordInvalidLen(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "Password Len 73",
			input: generateRandomString(73),
			want:  "password is longer than the accepted limit",
		},
		{
			name:  "Password Len 90",
			input: generateRandomString(90),
			want:  "password is longer than the accepted limit",
		},
	}

	for _, tc := range testCases {
		got, err := auth.HashPassword(tc.input)

		if got != "" || err.Error() != tc.want {

			t.Errorf("TestName:\t%s\nErr:\t%s\tWanted:\t%s", tc.name, err.Error(), tc.want)

		}
	}
}

func TestHashingPW(t *testing.T) {

	password := generateRandomString(10)
	hashPassword, err := auth.HashPassword(password)

	if err != nil {
		t.Errorf("Expecting nil and got %s", err.Error())
	}

	if hashPassword == "" || len(hashPassword) == 30 {
		t.Error("Hashed password is invalid")
	}

	password = generateRandomString(20)
	hashPassword, err = auth.HashPassword(password)

	if err != nil {
		t.Errorf("Expecting nil and got %s", err.Error())
	}

	if hashPassword == "" || len(hashPassword) == 30 {
		t.Error("Hashed password is invalid")
	}
}

func TestCompareMatchingPWHash(t *testing.T) {

	pw := generateRandomString(10)
	hashPW, err := auth.HashPassword(pw)
	if err != nil {
		t.Error("Error hashing pw")
	}

	if hashPW == "" || len(hashPW) == 30 {
		t.Error("Hashed password invalid")
	}

	err = auth.CheckPasswordHash(pw, hashPW)

	if err != nil {
		t.Errorf("PW don't match, PW: %s\tHashPW: %s\tError:%v\n", pw, hashPW, err.Error())
	}

	pw = generateRandomString(20)
	hashPW, err = auth.HashPassword(pw)

	if err != nil {
		t.Error("Error hashing pw")
	}

	if hashPW == "" || len(hashPW) == 30 {
		t.Error("Hashed password invalid")
	}

	err = auth.CheckPasswordHash(pw, hashPW)

	if err != nil {
		t.Errorf("PW don't match, PW: %s\tHashPW: %s\tError:%v\n", pw, hashPW, err.Error())
	}
}

func TestCompareNonMatcthingPWHashPW(t *testing.T) {

	pw := generateRandomString(10)
	hashPW, err := auth.HashPassword(pw)
	if err != nil {
		t.Error("Error hashing pw")
	}

	if hashPW == "" || len(hashPW) == 30 {
		t.Error("Hashed password invalid")
	}

	err = auth.CheckPasswordHash("Shouldn't Match", hashPW)

	if err == nil {
		t.Errorf("PW don't match, PW: %s\tHashPW: %s\n", pw, hashPW)
	}

	pw = generateRandomString(20)
	hashPW, err = auth.HashPassword(pw)

	if err != nil {
		t.Error("Error hashing pw")
	}

	if hashPW == "" || len(hashPW) == 30 {
		t.Error("Hashed password invalid")
	}

	err = auth.CheckPasswordHash("Shoulnd't match", hashPW)

	if err == nil {
		t.Errorf("PW don't match, PW: %s\tHashPW: %s\n", pw, hashPW)
	}
}
