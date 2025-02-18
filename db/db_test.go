package db

import (
	"os"
	"testing"
)

func TestInitDB(t *testing.T) {
	// Define test cases as an array of struct
	tests := []struct {
		name          string
		dbPassword    string
		expectedError bool
	}{
		{
			name:          "Valid DB password",
			dbPassword:    "qawzsx1", // Correct password
			expectedError: false,     // Expect no error
		},
		{
			name:          "Empty DB password",
			dbPassword:    "",   // Empty password (this will fail)
			expectedError: true, // Expect an error
		},
		{
			name:          "Invalid DB password",
			dbPassword:    "wrongpassword", // Incorrect password
			expectedError: true,            // Expect an error
		},
	}

	// Iterate through test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("DB_PASSWORD", tt.dbPassword)
			dbPassword := os.Getenv("DB_PASSWORD")
			if dbPassword == "" && tt.dbPassword != "" {
				t.Fatal("DB_PASSWORD is not set correctly")
			}
			t.Log("DB_PASSWORD:", dbPassword)

			DB = nil
			err := InitDB()

			if tt.expectedError && err == nil {
				t.Fatalf("Expected error but got none")
			}
			if !tt.expectedError && err != nil {
				t.Fatalf("Expected no error but got: %v", err)
			}

			// If no error expected, check if the DB is reachable
			if !tt.expectedError && DB != nil {
				if err := DB.Ping(); err != nil {
					t.Fatalf("Failed to ping the database: %v", err)
				}
			}
		})
	}
}
