package configs

import (
	"os"
	"path"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/raulaguila/go-api/pkg/utils"
)

func TestInit(t *testing.T) {
	tests := []struct {
		name       string
		envContent map[string]string
		setup      func()
		teardown   func()
		wantErr    bool
		tzExpected string
	}{
		{
			name: "SuccessfulEnvLoad",
			envContent: map[string]string{
				"TZ": "UTC",
			},
			setup: func() {
				_, b, _, _ := runtime.Caller(0)
				envPath := path.Join(path.Dir(b), "configs", ".env")
				os.WriteFile(envPath, []byte("TZ=UTC\n"), 0644)
				godotenv.Load(envPath)
			},
			teardown: func() {
				_, b, _, _ := runtime.Caller(0)
				envPath := path.Join(path.Dir(b), "configs", ".env")
				os.Remove(envPath)
			},
			wantErr:    false,
			tzExpected: "America/Manaus",
		},
		{
			name:       "DefaultLocalTime",
			envContent: map[string]string{},
			setup: func() {
				os.Unsetenv("TZ")
			},
			teardown: func() {},
			wantErr:  false,
			// Assuming local system default timezone if TZ is not set.
			tzExpected: time.Local.String(),
		},
		{
			name:       "InvalidTimezone",
			envContent: map[string]string{"TZ": "Invalid/Timezone"},
			setup: func() {
				os.Setenv("TZ", "Invalid/Timezone")
			},
			teardown: func() {
				os.Unsetenv("TZ")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			defer tt.teardown()

			err := godotenv.Load(path.Join("configs", ".env"))
			if err != nil {
				_, b, _, _ := runtime.Caller(0)
				utils.PanicIfErr(godotenv.Load(path.Join(path.Dir(b), ".env")))
			}
			utils.PanicIfErr(os.Setenv("SYS_VERSION", strings.TrimSpace("1.0.0")))

			_, err = time.LoadLocation(os.Getenv("TZ"))
			if (err != nil) != tt.wantErr {
				t.Errorf("unexpected error state: got %v, want %v", err != nil, tt.wantErr)
			}

			if err == nil {
				if loc := time.Local.String(); loc != tt.tzExpected {
					t.Errorf("unexpected timezone: got %v, want %v", loc, tt.tzExpected)
				}
			}
		})
	}
}
