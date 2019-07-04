package s3update_test

import (
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws/awserr"

	"github.com/heetch/s3update"
)

// TestUpdateDisabled tests that calling AutoUpdate() with en empty Updater doesn't lead to errors
// when a specific environment variable is set.
func TestUpdateDisabled(t *testing.T) {
	os.Setenv("S3UPDATE_DISABLED", "true")
	defer os.Setenv("S3UPDATE_DISABLED", "")

	err := s3update.AutoUpdate(s3update.Updater{})
	if err != nil {
		t.Fatal(err)
	}
}

// TestEmptyUpdateErrors tests that calling AutoUpdate() with an Updater where not all fields have values leads to errors.
func TestEmptyUpdateErrors(t *testing.T) {
	u := s3update.Updater{}
	err := s3update.AutoUpdate(u)
	expectedErr := "no version set"
	if err == nil {
		t.Fatal("Expected an error")
	} else if err.Error() != expectedErr {
		t.Fatalf("Expected error message \"%v\", but was \"%v\"", expectedErr, err)
	}

	u.CurrentVersion = "test"
	err = s3update.AutoUpdate(u)
	expectedErr = "no bucket set"
	if err == nil {
		t.Fatal("Expected an error")
	} else if err.Error() != expectedErr {
		t.Fatalf("Expected error message \"%v\", but was \"%v\"", expectedErr, err)
	}

	u.S3Bucket = "test"
	err = s3update.AutoUpdate(u)
	expectedErr = "no s3 region"
	if err == nil {
		t.Fatal("Expected an error")
	} else if err.Error() != expectedErr {
		t.Fatalf("Expected error message \"%v\", but was \"%v\"", expectedErr, err)
	}

	u.S3Region = "test"
	err = s3update.AutoUpdate(u)
	expectedErr = "no s3ReleaseKey set"
	if err == nil {
		t.Fatal("Expected an error")
	} else if err.Error() != expectedErr {
		t.Fatalf("Expected error message \"%v\", but was \"%v\"", expectedErr, err)
	}

	u.S3ReleaseKey = "test"
	err = s3update.AutoUpdate(u)
	expectedErr = "no s3VersionKey set"
	if err == nil {
		t.Fatal("Expected an error")
	} else if err.Error() != expectedErr {
		t.Fatalf("Expected error message \"%v\", but was \"%v\"", expectedErr, err)
	}
}

// TestUpdateErrors tests specific errors when calling AutoUpdate() with an Updater with invalid values.
func TestUpdateErrors(t *testing.T) {
	u := s3update.Updater{
		CurrentVersion: "test",
		S3Bucket:       "test",
		S3Region:       "test",
		S3ReleaseKey:   "test",
		S3VersionKey:   "test",
	}
	err := s3update.AutoUpdate(u)
	expectedErr := "invalid local version"
	if err == nil {
		t.Fatal("Expected an error")
	} else if err.Error() != expectedErr {
		t.Fatalf("Expected error message \"%v\", but was \"%v\"", expectedErr, err)
	}

	u.CurrentVersion = "0"
	err = s3update.AutoUpdate(u)
	expectedErr = "invalid local version"
	if err == nil {
		t.Fatal("Expected an error")
	} else if err.Error() != expectedErr {
		t.Fatalf("Expected error message \"%v\", but was \"%v\"", expectedErr, err)
	}

	u.CurrentVersion = "1"
	err = s3update.AutoUpdate(u)
	if err == nil {
		t.Fatal("Expected an error")
	} else if awsErr, ok := err.(awserr.Error); ok {
		// TODO: Check for type awserr.RequestFailure or find out if the aws package contains an error constant with this code
		if awsErr.Code() != "RequestError" {
			t.Fatalf("Expected awserr.Error.Code \"RequestError\", but got %v", awsErr.Code())
		}
	} else {
		t.Fatalf("An awserr.Error was expected, but got %T", err)
	}
}
