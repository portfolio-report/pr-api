package argon2

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	hash, err := HashPasswordDefault("secret")
	require.Nil(t, err)

	regexHash := regexp.MustCompile(`^\$argon2id\$v=19\$m=65536,t=1,p=2\$[A-Za-z0-9+/]{22}\$[A-Za-z0-9+/]{43}$`)
	assert.Regexp(t, regexHash, hash)

	hash2, err := HashPasswordDefault("secret")
	require.Nil(t, err)

	assert.NotEqual(t, hash, hash2)
}

func TestVerifyPassword(t *testing.T) {
	testCases := []struct {
		name     string
		hash     string
		password string
		valid    bool
	}{
		{
			name:     "valid argon2id",
			hash:     "$argon2id$v=19$m=4096,t=2,p=1$ckg3QktpRlR2RFUxeVBkOA$QAZA+hzrcOKtz5RggyNdZQ",
			password: "secret",
			valid:    true,
		},
		{
			name:     "invalid argon2id",
			hash:     "$argon2id$v=19$m=4096,t=2,p=1$ckg3QktpRlR2RFUxeVBkOA$QAZA+hzrcOKtz5RggyNdZQ",
			password: "wrong",
			valid:    false,
		},
		{
			name:     "valid argon2i",
			hash:     "$argon2i$v=19$m=4096,t=2,p=1$eVAzSUY4QWZoVkloNnBKQQ$VHPIMSUFkMmOkB/hMh/rFQ",
			password: "secret",
			valid:    true,
		},
		{
			name:     "invalid argon2i",
			hash:     "$argon2i$v=19$m=4096,t=2,p=1$eVAzSUY4QWZoVkloNnBKQQ$VHPIMSUFkMmOkB/hMh/rFQ",
			password: "wrong",
			valid:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			valid, err := VerifyPassword(tc.password, tc.hash)
			require.Nil(t, err)
			assert.Equal(t, tc.valid, valid)
		})
	}

	invalidHashes := []string{
		"",
		"secret",
		"$argon2d$v=19$m=16,t=2,p=1$NkJBV2RBaUFDaVU2aHRvNQ$eOEVxJRX6Uj7vjT2e4eawA",
		"$argon2id$v$m=4096,t=2,p=1$ckg3QktpRlR2RFUxeVBkOA$QAZA+hzrcOKtz5RggyNdZQ",
		"$argon2id$v=1$m=4096,t=2,p=1$ckg3QktpRlR2RFUxeVBkOA$QAZA+hzrcOKtz5RggyNdZQ",
		"$argon2id$v=19$m,t,p$ckg3QktpRlR2RFUxeVBkOA$QAZA+hzrcOKtz5RggyNdZQ",
		"$argon2id$v=19$m=4096,t=2,p=1$_$QAZA+hzrcOKtz5RggyNdZQ",
		"$argon2id$v=19$m=4096,t=2,p=1$ckg3QktpRlR2RFUxeVBkOA$_",
	}

	for _, invalidHash := range invalidHashes {
		t.Run("invalid hash "+invalidHash, func(t *testing.T) {
			valid, err := VerifyPassword("secret", invalidHash)
			assert.Error(t, err)
			assert.False(t, valid)
		})
	}
}

func BenchmarkHashPasswordDefault(b *testing.B) {
	for n := 0; n < b.N; n++ {
		HashPasswordDefault("secret")
	}
}

func BenchmarkVerifyPassword(b *testing.B) {
	for n := 0; n < b.N; n++ {
		VerifyPassword("secret", "$argon2id$v=19$m=4096,t=2,p=1$ckg3QktpRlR2RFUxeVBkOA$QAZA+hzrcOKtz5RggyNdZQ")
	}
}
