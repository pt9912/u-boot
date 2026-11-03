# im-test-docker-ergaenzung.sh — M4-Ergänzung zu slice-v1-bau-nur-devcontainer: die
# Fehlschläge aus roh-test-docker.txt, deren Meldung im 14-Zeilen-Ausschnitt fehlt, je
# mit Assert-Kontext laufen (die Messung liest; keine Behebung). Aufruf siehe README.
# TEST_FILTER (Go-regex, Standard: die sechs ungenannten) und TEST_MODE
# (assert = Assert-Zeilen je Fehlschlag; roh = alle Meldungszeilen ungefiltert nach Quelle).
set -u
cd /workspaces/u-boot
filter="${TEST_FILTER:-TestLauncher_Docker_S5_GuardFalsifiedIntegrityBreaksTheBuild|TestLauncher_Docker_S6_ValueIsNeverInImageLabelFileOrCommandLine|TestE2E_ForgeCLI_PinMismatchBreaksBuild|TestE2E_SandboxDevcontainer_Ollama_PinMismatchBreaksBuild|TestE2E_Toolchain_PinMismatchBreaksBuild|TestE2E_SandboxDevcontainer_StatusFiles_VFS}"
mode="${TEST_MODE:-assert}"
sudo sh /usr/local/bin/u-boot-docker-init >/tmp/docker-init.log 2>&1
for i in $(seq 1 60); do docker info >/dev/null 2>&1 && break; sleep 1; done
git config --global --add safe.directory /workspaces/u-boot
echo "== gezielter Lauf: $filter"
t0=$(date +%s); go test -v -tags docker -timeout 20m -run "$filter" ./... >/tmp/ergaenzung.log 2>&1; rc=$?; t1=$(date +%s)
echo "ergaenzung-exit=$rc wall=$(awk -v a="$t1" -v b="$t0" 'BEGIN{printf "%.1f", a-b}')s"
if [ "$mode" = "roh" ]; then
  echo "== alle Meldungszeilen des Laufes (ohne RUN-Texte)"
  grep -E 'lacks|another reason|the value was on|want [0-9]+|leftover|proof|restricts unprivileged' /tmp/ergaenzung.log | cut -c1-240 | head -40
fi
echo "== Zählung"
printf 'PASS=%s FAIL=%s SKIP=%s\n' \
  "$(grep -c -E '^--- PASS' /tmp/ergaenzung.log)" \
  "$(grep -c -E '^--- FAIL' /tmp/ergaenzung.log)" \
  "$(grep -c -E '^--- SKIP' /tmp/ergaenzung.log)"