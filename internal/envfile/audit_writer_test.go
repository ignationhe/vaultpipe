package envfile

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func singleEntryLog() AuditLog {
	return AuditLog{
		Entries: []AuditEntry{
			{
				Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
				Key:       "DB_HOST",
				Action:    AuditAdded,
				OldValue:  "",
				NewValue:  "***",
				Source:    "vault://secret/app",
			},
		},
	}
}

func TestWriteAuditLog_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteAuditLog(singleEntryLog(), &buf, AuditFormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected DB_HOST in output, got:\n%s", out)
	}
	if !strings.Contains(out, "added") {
		t.Errorf("expected 'added' in output, got:\n%s", out)
	}
	if !strings.Contains(out, "TIMESTAMP") {
		t.Errorf("expected header in output, got:\n%s", out)
	}
}

func TestWriteAuditLog_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteAuditLog(singleEntryLog(), &buf, AuditFormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var entries []AuditEntry
	if err := json.Unmarshal(buf.Bytes(), &entries); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(entries) != 1 || entries[0].Key != "DB_HOST" {
		t.Errorf("unexpected entries: %+v", entries)
	}
}

func TestWriteAuditLogToFile_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")
	if err := WriteAuditLogToFile(singleEntryLog(), path, AuditFormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("file not created: %v", err)
	}
	if !strings.Contains(string(data), "DB_HOST") {
		t.Errorf("expected DB_HOST in file, got:\n%s", data)
	}
}

func TestWriteAuditLogToFile_Appends(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")
	for i := 0; i < 2; i++ {
		if err := WriteAuditLogToFile(singleEntryLog(), path, AuditFormatText); err != nil {
			t.Fatalf("write %d failed: %v", i, err)
		}
	}
	data, _ := os.ReadFile(path)
	count := strings.Count(string(data), "DB_HOST")
	if count != 2 {
		t.Errorf("expected 2 occurrences of DB_HOST (append), got %d", count)
	}
}
