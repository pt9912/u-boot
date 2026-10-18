# Pin und Lieferkette der Podman-Quelle altern: Version, Prüfsummen und Signaturschlüssel sind ein Stand, kein Zustand

**Sub-Area:** `internal/hexagon/application` (Kürzel APP, Greenfield laut [`harness/conventions.md`](../../../../../../harness/conventions.md#mr-006)) — das Template mit den Pins liegt dort (`templates/devcontainer/Dockerfile.tmpl`, Block `NestedPodman`).

Das erzeugte Dockerfile bei `devcontainer.sandbox.nestedRuntime: podman` lädt Podman 6.1.4 aus `mgoltzsche/podman-static` und prüft es gegen einen Pin (Version, SHA-256 je Architektur, Primär-Fingerprint; [`ADR-0023`](../../../../adr/0023-podman-aus-podman-static-im-sandbox-devcontainer.md) Entscheidung 2). Drei Dinge altern, ohne dass ein Test rot wird:

- **Pin-Alterung:** „Neueste Version“ gilt nur bis zum nächsten Podman-Release; eine Hebung ist Handarbeit (Slice mit Wiederholung der Abnahmetests, [`ADR-0023`](../../../../adr/0023-podman-aus-podman-static-im-sandbox-devcontainer.md) Entscheidung 9).
- **Schlüsselablauf:** Der Signaturschlüssel läuft am **2027-03-08** ab. Nach dem Ablauf besteht `gpg --verify` weiter (`EXPKEYSIG`, Exit-Code 0, `VALIDSIG` unverändert); die Vorlage erzwingt weder Ablauf noch Widerruf, die Integrität trägt die SHA-256. Ein Ablauf bricht den Bau also nicht, er macht die Signaturprüfung unbemerkt schwächer.
- **Lieferkette Einzelperson:** Die Quelle ist ein Community-Build einer Person und ihrer CI ([`ADR-0023`](../../../../adr/0023-podman-aus-podman-static-im-sandbox-devcontainer.md) Entscheidung 7); eine Rotation des Schlüssels bricht den Bau (gewollt), ein kompromittiertes Konto nicht vorhersehbar.

**Trigger:** Pin-Hebung samt Schlüsselprüfung **vor 2027-03-08**; ferner die Re-Evaluierungs-Trigger der ADR (neue Hauptversion bei `podman-static`, Sicherheitslücke ohne Release binnen 14 Tagen, Rotation oder Ablauf des Schlüssels ohne bestätigten Nachfolger, Messbefund auf `amd64`, Docker Desktop oder cgroup v1). **Grenze:** Das Datum hat keinen Wächter; Leser ist die Slice-Planung (§8, Sichtungs-Schritt) jedes Slice, der `internal/hexagon/application` berührt, und der Release-Cut. Kein Carveout: Es ist nichts gelockert, die Prüfung gilt voll.
