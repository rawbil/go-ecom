package users

// func TestValidateUserEmail(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		email   interface{}
// 		wantErr bool
// 	}{
// 		{
// 			name:    "empty email filter is allowed",
// 			email:   "",
// 			wantErr: false,
// 		},
// 		{
// 			name:    "valid email filter",
// 			email:   "bildadsimiyu6@gmail.com",
// 			wantErr: false,
// 		},
// 		{
// 			name:    "invalid email filter",
// 			email:   "bildadsimiyu6",
// 			wantErr: true,
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			err := ValidateUserEmail(repository.ListUsersParams{
// 				Email: test.email,
// 			})

// 			if test.wantErr {
// 				if !utils.ValidationErrorCheck("email", err) {
// 					t.Fatalf("expected email validation error, got %v", err)
// 				}
// 				return
// 			}

// 			if err != nil {
// 				t.Fatalf("expected no validation error, got %v", err)
// 			}
// 		})
// 	}
// }
