**Vorgang:** [`slice-v1-devcontainer-env-datei`](../../../../done/slice-v1-devcontainer-env-datei.md)
**Fund:** `--env-file=…` kam als erster Eintrag in `runArgs` aller erzeugten Konfigurationen. Der Implementer fuhr `make gates` (grün) und gezielt drei Docker-Tests; `TestE2E_SandboxDevcontainer_NestedDocker` (`sandbox_devcontainer_docker_test.go`) pinnte `runArgs = [--privileged]` und blieb rot, bis der Volllauf von `make test-docker` (3738 `PASS`, 1 `FAIL`) ihn zeigte. Behoben im selben Slice; der gezielte Lauf ist grün (64,8 s).

**Beleg-Qualität:** belegt durch das Log des Volllaufs (Scratchpad der Sitzung) und den Nachzug im Test. Keine Schärfung entschieden.
