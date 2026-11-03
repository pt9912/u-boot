**Vorgang:** [`slice-v1-release-cut-v0.20.0`](../../../../done/slice-v1-release-cut-v0.20.0.md)
**Fund:** die Zustandsmarker-Klasse in der Release-Cut-Familie — **dritter Cut in Folge, immer dieselben zwei Träger**:

- **F-1 (Review, MEDIUM):** der §7-Marker des Plans trug „Noch nicht geschlossen; der Slice liegt in `open/` (Plan)." über den Tag hinaus — nach dem Tag ist der Inhalt-Commit vor dem `git mv` fällig (`AGENTS.md` §3.3 Fall 2), der Marker wandert sonst mit.
- **F-3 (Review, LOW):** die Bestand-Zeile `open/README.md` trug „… nur Plan" — sie existierte seit dem Plan-Commit (die F-5-Lehre des [`slice-v1-release-cut-v0.19.0`](../../../../done/slice-v1-release-cut-v0.19.0.md) war getragen), aber sie beschrieb den Zustand weiter statt den Tag-Stand.
- **F-2/V-1:** die Nächste-Wellen-Zeile trug die Vorschlags-Version `v0.20.0` über den Tag hinaus — die Version ist mit dem Release verbraucht.

Dazu aus dem Verify-Report (je LOW): **V-2** der Review-Report trug die `repo.mk`-Fundstellen 66/67 vertauscht (am Code ist 67 die `VERSION ?=`-Zeile); **V-5** die zwei Bericht-Dateien fehlten als §3-Zeilen (die Bericht-Dateien-Familie der Klasse, wiederum).

**Beleg-Qualität:** belegt durch die Reports ([Review](../../../../../../reviews/2026-11-03-slice-v1-release-cut-v0.20.0.md), [Verification](../../../../../../reviews/2026-11-03-slice-v1-release-cut-v0.20.0-verify.md)). Alle Stellen berichtigt: F-2/V-1 und F-4 in der Fix-Runde `f652ed4f`, F-4 (Lead) in `f652ed4f`, V-2/V-5/Bestand- und Marker-Texte in der Berichtigungs- und Closure-Runde (§7 durch den Closure-Block abgelöst — F-1); die Sensor-Frage „Bestand-Zeile und §7-Marker nach dem Tag gegen die Verzeichnis-Lage" liegt beim Projektinhaber.