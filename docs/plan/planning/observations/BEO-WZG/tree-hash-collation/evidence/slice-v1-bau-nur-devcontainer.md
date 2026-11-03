**Vorgang:** [`slice-v1-bau-nur-devcontainer`](../../../../done/slice-v1-bau-nur-devcontainer.md)
**Fund:** Verifier V-1 (MEDIUM) — der Gate-Nachweis-Hash (`tools/harness/working-tree-hash.sh`) pinnte die Collation nicht: Host-uutils-sort (0.10.0) vs. Container-GNU-sort (9.7) sortieren die Datei-Zeilen unterschiedlich (786 abweichende Zeilen; Host-Hash `48c16b81…` vs. Container-Stempel `c5149757…`); seit dem Vorgang stempelt `record-gates` im Container und der Stop-Hook liest am Host — nicht-streng verschleiert, streng falsch-rot je Turn.

**Beleg-Qualität:** beide Kontexte nach dem `LC_ALL=C`-Pin je Probe gelaufen — Host- und Container-Hash **byte-gleich** (`22d96ce4…` je Wert). Der Pin ist die Berichtigung; die Hash-Funktion selbst ist unverändert.

**Zähler:** 1 Vorgang (dieses); die 3×-Schwelle steht aus.