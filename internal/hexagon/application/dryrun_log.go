package application

import (
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// logWrite meldet ein Schreibergebnis auf Info. Im Dry-Run
// ([driving.PreviewDryRun]) ist nichts geschrieben worden: die Zeile
// trägt dann den Wortlaut `would` (Konjunktiv-Vokabular der
// Zusammenfassung), sonst `done` (Vergangenheitsform, das Ergebnis ist
// eingetreten, auch bei `--diff` allein). Pegel, Felder und Port bleiben
// gleich.
func logWrite(l driven.Logger, mode driving.PreviewMode, done, would string, args ...any) {
	if mode == driving.PreviewDryRun {
		l.Info(would, args...)
		return
	}
	l.Info(done, args...)
}

// devcontainerWriteMessage ist die Logzeile von `generate devcontainer`
// nach einem Schreibzugriff: im Dry-Run `would create` (Aktion
// [driving.GenerateActionCreated]) oder `would update block`, sonst
// `<Aktion>` ([driving.GenerateAction.String]). Die Plan-Phase liefert
// nur diese beiden Aktionen ([devcontainerAggregateAction]; `NoOp`
// endet vor der Zeile).
func devcontainerWriteMessage(mode driving.PreviewMode, action driving.GenerateAction) string {
	const prefix = "generate devcontainer: "
	if mode != driving.PreviewDryRun {
		return prefix + action.String()
	}
	if action == driving.GenerateActionCreated {
		return prefix + "would create"
	}
	return prefix + "would update block"
}
