//go:build docker

// LH-FA-UP-005 Logs-Acceptance-Flow Pin (slice-v1-logs T4).
//
// Spec §1023-1040: `u-boot logs [service] [--follow] [--tail <n>]`
// streamt Compose-Logs an stdout. Zwei Tests pinnen die zwei
// Pflicht-Flag-Pfade gegen einen echten postgres-Compose-Stack:
//
//   - TestE2E_LHFAUP005_LogsTail: `--tail 20` produziert die
//     kanonische postgres-Boot-Phrase `database system is ready` im
//     OutputSink (Sink-Streaming durchgereicht bis zur echten
//     Compose-CLI; T1-Adapter + T2-Service + driven-Engine).
//   - TestE2E_LHFAUP005_LogsFollow: `--follow` läuft mit `ctx`-
//     Deadline; nach Timeout returnt `LogsService.Logs` `nil`
//     (SIGINT-Vertrag Schicht 2 — slice-v1-logs §SIGINT-Vertrag).
//
// Setup-Strategie: `runAcceptanceFlow` (acceptance_helpers.go:80)
// fährt init+add+up gegen postgres hoch — identisch zum LH-AK-002-
// Pfad. Die LogsService-Aufrufe operieren danach auf demselben
// `compose.yaml` im `res.dir` mit einem frisch instanziierten
// dockeradapter (stateless für Compose-Calls; das Compose-Projekt
// wird über den BaseDir/compose-Pfad adressiert, nicht über Engine-
// State). Helper-Variante A (slice-v1-logs T4-Plan-Entscheidung):
// kein Refactor von `runAcceptanceFlow` — Logs-Pfad lebt nur hier.
//
// Test-Prerequisites identisch LH-AK-002: docker CLI + compose
// plugin verfügbar; Test-Prozess + Daemon teilen Netzwerk-
// Namespace via `make test-docker` (`--network=host`).

package e2e_test

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	dockeradapter "github.com/pt9912/u-boot/internal/adapter/driven/docker"
	fsadapter "github.com/pt9912/u-boot/internal/adapter/driven/fs"
	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// postgresReadyPhrase ist die kanonische Boot-Bestätigung von
// postgres (jede Major-Version >= 9 schreibt diese Zeile in den
// stderr-Log-Stream, sobald der Server Verbindungen akzeptiert).
// Wenn diese Phrase im Sink ankommt, ist die gesamte Streaming-
// Kette adapter → use-case → caller-Sink belegt.
const postgresReadyPhrase = "database system is ready"

func TestE2E_LHFAUP005_LogsTail(t *testing.T) {
	res := runAcceptanceFlow(t, acceptanceFlow{
		projectName: "t-uboot-e2e-logs-tail",
		serviceName: "postgres",
		envKeys:     []string{"POSTGRES_USER"},
		upTimeout:   90 * time.Second,
		ctxTimeout:  3 * time.Minute,
	})

	var sink bytes.Buffer
	logsSvc := application.NewLogsService(fsadapter.New(), dockeradapter.NewEngine(), nil)

	// `--tail 20` ist großzügig genug, um die Boot-Phrase auch dann
	// einzufangen, wenn der postgres-Container in der Zwischenzeit
	// noch weitere Log-Zeilen geschrieben hat. Compose addressiert
	// den ringbuffer der Container-Stdio, nicht die Live-Quelle —
	// die `database system is ready`-Zeile bleibt am Stack-Boden
	// erhalten.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if _, err := logsSvc.Logs(ctx, driving.LogsRequest{
		BaseDir:    res.dir,
		Service:    "postgres",
		Tail:       "20",
		OutputSink: &sink,
	}); err != nil {
		t.Fatalf("Logs(--tail 20): %v", err)
	}

	if !strings.Contains(sink.String(), postgresReadyPhrase) {
		t.Errorf("logs --tail 20 sink missing %q; got:\n%s",
			postgresReadyPhrase, sink.String())
	}
}

func TestE2E_LHFAUP005_LogsFollow(t *testing.T) {
	res := runAcceptanceFlow(t, acceptanceFlow{
		projectName: "t-uboot-e2e-logs-follow",
		serviceName: "postgres",
		envKeys:     []string{"POSTGRES_USER"},
		upTimeout:   90 * time.Second,
		ctxTimeout:  3 * time.Minute,
	})

	var sink bytes.Buffer
	logsSvc := application.NewLogsService(fsadapter.New(), dockeradapter.NewEngine(), nil)

	// 8 s reicht in der Praxis: `docker compose logs --follow` flusht
	// erst den bestehenden Buffer, **dann** blockt es auf neue Zeilen.
	// Der Buffer-Flush trifft die Boot-Phrase deterministisch — der
	// Test prüft danach nur, dass der Timeout SIGINT-Vertrag-konform
	// als nil-error zurückkommt.
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	if _, err := logsSvc.Logs(ctx, driving.LogsRequest{
		BaseDir:    res.dir,
		Service:    "postgres",
		Follow:     true,
		Tail:       "all",
		OutputSink: &sink,
	}); err != nil {
		t.Errorf("Logs(--follow) returned %v after DeadlineExceeded; want nil (SIGINT-Vertrag Schicht 2)", err)
	}
	if sink.Len() == 0 {
		t.Errorf("Logs(--follow) sink is empty; expected at least the pre-block Compose-Buffer-Flush content")
	}
}
