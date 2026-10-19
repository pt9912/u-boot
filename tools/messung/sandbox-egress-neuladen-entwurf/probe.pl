#!/usr/bin/perl
# Messsonde im Container: protokolliert nur Zustandswechsel, mit Epochenzeit in Sekunden (Millisekunden-Aufloesung).
# Aufruf: probe.pl <sekunden> <dns-name> [egress|dns|both]   (zwei Prozesse getrennt, damit die Wartezeit der
#         DNS-Sonde die Aufloesung der Egress-Sonde nicht begrenzt)
#   egress : UDP-send an 1.1.1.1:9. Ein `drop` im OUTPUT-Haken liefert dem Aufrufer einen Fehler (EPERM), ein offener Pfad
#            liefert Erfolg; "offen" heisst also: kein Filter-Regelwerk haelt den Verkehr zurueck.
#   dns    : Anfrage fuer <dns-name> (Typ A) an 127.0.0.1:53, den lokalen dnsmasq; Antwort mit Rcode, keine Antwort
#            (Timeout 20 ms) oder abgelehnter Port (kein Listener).
use strict; use warnings;
use Socket;
use IO::Select;
use Time::HiRes qw(time);
my ($dur, $name, $mode) = @ARGV;
$mode ||= 'both';
my $t0 = time;
printf "t0 epoch=%.3f dauer=%ss dns-name=%s modus=%s\n", $t0, $dur, $name, $mode;
sub query {
  my $id = int(rand(65535));
  my $p = pack('nnnnnn', $id, 0x0100, 1, 0, 0, 0)
        . join('', map { pack('C/a*', $_) } split(/\./, $name)) . pack('C', 0) . pack('nn', 1, 1);
  return $p;
}
my ($le, $ld) = ('', '');
sub note { my ($kind, $state, $lastref) = @_; if ($state ne $$lastref) { $$lastref = $state; printf "%.3f %s %s\n", time, $kind, $state } }
my $iter = 0;
while (time - $t0 < $dur) {
  $iter++;
  if ($mode ne 'dns') {
    socket(my $e, PF_INET, SOCK_DGRAM, 0) or die "socket: $!";
    connect($e, sockaddr_in(9, inet_aton('1.1.1.1'))) or die "connect: $!";
    my $ok = send($e, 'x', 0);
    note('egress', defined $ok ? 'offen' : "gesperrt ($!)", \$le);
    close $e;
  }
  next if $mode eq 'egress';

  socket(my $d, PF_INET, SOCK_DGRAM, 0) or die "socket: $!";
  connect($d, sockaddr_in(53, inet_aton('127.0.0.1'))) or die "connect: $!";
  my $s = defined send($d, query(), 0) ? 'ok' : "sendfehler ($!)";
  my $state;
  if ($s ne 'ok') { $state = $s }
  else {
    my $sel = IO::Select->new($d);
    if ($sel->can_read(0.02)) {
      my $buf;
      if (defined recv($d, $buf, 512, 0) && length($buf) >= 4) { $state = 'antwort rcode=' . (unpack('n', substr($buf, 2, 2)) & 0xf) }
      else { $state = "leseFehler ($!)" }
    } else { $state = 'timeout' }
  }
  note('dns', $state, \$ld);
  close $d;
}
printf "%.3f ende iterationen=%d\n", time, $iter;
