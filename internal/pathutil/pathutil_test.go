package pathutil

import (
	"testing"
)

func TestGetProjectPath_Valid(t *testing.T) {
	got, _ := GetProjectPath("/")
	want := "/"

	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestGetProject_PathEmpty(t *testing.T) {
	got, err := GetProjectPath("")
	want := ""
	wantErrMsg := "missing project path argument"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != wantErrMsg {
		t.Errorf("got %s, want %s", err, wantErrMsg)
	}
}

func TestGetProjectPath_Blank(t *testing.T) {
	got, err := GetProjectPath("")
	want := ""
	wantErrMsg := "missing project path argument"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != wantErrMsg {
		t.Errorf("got %s, want %s", err, wantErrMsg)
	}
}

func TestValidateDir_Valid(t *testing.T) {
	got := ValidateDir("/")
	want := true

	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestValidateDir_Blank(t *testing.T) {
	got := ValidateDir("")
	want := false

	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestValidateFile_Valid(t *testing.T) {
	got := ValidateFile("../../test_data/test.md")
	want := true
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}

}

func TestValidateFile_Invalid(t *testing.T) {
	got := ValidateFile("/")
	want := false
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}

}

func TestValidateFile_BadFiile(t *testing.T) {
	got := ValidateFile("lemons")
	want := false
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}

}
