package lib

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExportImportRoundtrip(t *testing.T) {
	BuildSchedule()

	home := t.TempDir()
	t.Setenv("HOME", home)

	if err := EnsureResultsFile(); err != nil {
		t.Fatal(err)
	}
	if err := EnsureParticipantsDir(); err != nil {
		t.Fatal(err)
	}

	results, err := resultsPath()
	if err != nil {
		t.Fatal(err)
	}
	resultsContent := "1 2 1\n.\n.\n.\n.\n.\n.\n.\n"
	if err := os.WriteFile(results, []byte(resultsContent), 0o644); err != nil {
		t.Fatal(err)
	}

	participants, err := participantsDir()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(participants, "alice"), []byte(emptyResultsContent()), 0o644); err != nil {
		t.Fatal(err)
	}

	archivePath := filepath.Join(home, "backup.tar.gz")
	if err := ExportData(archivePath); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(results, []byte(emptyResultsContent()), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(participants, "alice")); err != nil {
		t.Fatal(err)
	}

	if err := ImportData(archivePath); err != nil {
		t.Fatal(err)
	}

	got, err := os.ReadFile(results)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != resultsContent {
		t.Fatalf("unexpected results after import: %q", string(got))
	}
	if _, err := os.Stat(filepath.Join(participants, "alice")); err != nil {
		t.Fatalf("expected participant file after import: %v", err)
	}
}

func TestValidateArchiveEntryRejectsTraversal(t *testing.T) {
	if _, err := validateArchiveEntry("../results"); err == nil {
		t.Fatal("expected error for path traversal")
	}
	if _, err := validateArchiveEntry("other/file"); err == nil {
		t.Fatal("expected error for unexpected path")
	}
}
