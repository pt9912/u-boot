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
// nach einem Schreibzugriff: im Dry-Run `would <Verb>` zur Aktion,
// sonst `<Aktion>` im bisherigen Wortlaut ([driving.GenerateAction.String]).
func devcontainerWriteMessage(mode driving.PreviewMode, action driving.GenerateAction) string {
	const prefix = "generate devcontainer: "
	if mode != driving.PreviewDryRun {
		return prefix + action.String()
	}
	switch action {
	case driving.GenerateActionCreated:
		return prefix + "would create"
	case driving.GenerateActionUpdatedBlock:
		return prefix + "would update block"
	default:
		// Andere Aktionen liefert die Plan-Phase von `generate devcontainer`
		// nicht ([devcontainerAggregateAction]).
		return prefix + "would change"
	}
}
