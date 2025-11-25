package backend

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestInitCreatesPaths(t *testing.T) {
	root := t.TempDir()
	b := New(root)
	b.Init()

	for _, p := range b.Paths.All() {
		if stat, err := os.Stat(p); err != nil {
			t.Fatalf("expected path %s to exist: %v", p, err)
		} else if !stat.IsDir() {
			t.Fatalf("expected path %s to be directory", p)
		}
	}
}

func TestCreateDailyNoteAddsMetadata(t *testing.T) {
	root := t.TempDir()
	b := New(root)
	b.Init()

	date := time.Date(2024, time.January, 2, 0, 0, 0, 0, time.UTC)
	filename, err := b.CreateNew_DailyNote(date)
	if err != nil {
		t.Fatalf("CreateNew_DailyNote returned error: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(b.DailyNotes, filename))
	if err != nil {
		t.Fatalf("reading created note failed: %v", err)
	}

	text := string(content)
	for _, snippet := range []string{"icon = \"󰃭\"", "tags = [\"daily\"]", "created = 2024-01-02T00:00:00Z", "# Jan 2 2024"} {
		if !strings.Contains(text, snippet) {
			t.Fatalf("expected metadata %q in daily note", snippet)
		}
	}
}

func TestCreateRandomNoteGeneratesUniqueNames(t *testing.T) {
	root := t.TempDir()
	b := New(root)
	b.Init()

	first, err := b.CreateNew_RandomNote()
	if err != nil {
		t.Fatalf("CreateNew_RandomNote returned error: %v", err)
	}

	second, err := b.CreateNew_RandomNote()
	if err != nil {
		t.Fatalf("second CreateNew_RandomNote returned error: %v", err)
	}

	if first == second {
		t.Fatalf("expected unique filenames, got %s for both", first)
	}

	if !strings.HasSuffix(first, "-01.md") || !strings.HasSuffix(second, "-02.md") {
		t.Fatalf("expected sequential suffixes, got %s and %s", first, second)
	}
}

func TestTrashNoteMovesToTypedTrash(t *testing.T) {
	root := t.TempDir()
	b := New(root)
	b.Init()

	fname := "example.md"
	notePath := filepath.Join(b.RandomNotes, fname)
	if err := os.WriteFile(notePath, []byte("# note"), 0o600); err != nil {
		t.Fatalf("failed to write seed note: %v", err)
	}

	if err := b.TrashNote(fname, b.RandomNotes); err != nil {
		t.Fatalf("TrashNote returned error: %v", err)
	}

	if _, err := os.Stat(notePath); !os.IsNotExist(err) {
		t.Fatalf("expected original note to be moved, got err=%v", err)
	}

	trashedPath := filepath.Join(b.TrashRandom, fname)
	if _, err := os.Stat(trashedPath); err != nil {
		t.Fatalf("expected trashed note at %s: %v", trashedPath, err)
	}
}

func TestParseMarkdownNoteReadsDefaults(t *testing.T) {
	root := t.TempDir()
	noteName := "sample.md"
	notePath := filepath.Join(root, noteName)

	body := "+++\ncreated = 2024-04-03T10:00:00Z\nicon = \"🧠\"\n+++\n\n# Custom Title\nContent"
	if err := os.WriteFile(notePath, []byte(body), 0o600); err != nil {
		t.Fatalf("failed to write note: %v", err)
	}

	title, icon, tags, created, err := ParseMarkdownNote(root, noteName)
	if err != nil {
		t.Fatalf("ParseMarkdownNote returned error: %v", err)
	}

	if title != "Custom Title" {
		t.Fatalf("expected title 'Custom Title', got %s", title)
	}
	if icon != "🧠" {
		t.Fatalf("expected icon '🧠', got %s", icon)
	}
	if len(tags) != 0 {
		t.Fatalf("expected empty tags, got %v", tags)
	}
	if created.IsZero() || created.Year() != 2024 {
		t.Fatalf("expected created time parsed, got %v", created)
	}
}
