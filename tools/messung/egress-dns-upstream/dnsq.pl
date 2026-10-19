#!/usr/bin/perl
# Sonde fuer D5: schickt nacheinander DNS-Anfragen (A, UDP, Port 53) mit zufaelligen Labels vor einem Suffix und wartet je Anfrage auf die Antwort.
# Aufruf: dnsq.pl <server> <suffix> <sekunden> <labels> [label-laenge]
#   <labels> Labels je Anfrage (je <label-laenge> Hex-Zeichen, Vorgabe 63); der Name bleibt unter 253 Zeichen.
# Ausgabe: queries= antworten= rcode0= sekunden= zeichen= (Label-Zeichen, die in Anfragen gingen); Nutzlast-Zahlen rechnet das README.
# Nur Kernmodule (IO::Socket::INET, Time::HiRes, IO::Select).
use strict; use warnings;
use IO::Socket::INET; use IO::Select; use Time::HiRes qw(time);
my ($srv, $suffix, $secs, $nlab, $llen) = @ARGV;
die "usage: dnsq.pl server suffix sekunden labels [laenge]\n" unless defined $nlab;
$llen //= 63;
my $s = IO::Socket::INET->new(PeerAddr => $srv, PeerPort => 53, Proto => 'udp') or die "socket: $!\n";
my $sel = IO::Select->new($s);
my ($q, $a, $ok, $chars) = (0, 0, 0, 0);
my $end = time() + $secs;
my $t0 = time();
while (time() < $end) {
  my @labels = map { join '', map { sprintf '%x', int(rand 16) } 1 .. $llen } 1 .. $nlab;
  my $name = join('.', @labels) . ".$suffix";
  die "Name zu lang\n" if length($name) > 253;
  my $id = int(rand 65536);
  my $pkt = pack('nnnnnn', $id, 0x0100, 1, 0, 0, 0)
          . join('', map { pack('C', length $_) . $_ } split /\./, $name) . "\0" . pack('nn', 1, 1);
  $s->send($pkt) or next;
  $q++; $chars += length(join '', @labels);
  if ($sel->can_read(1.0)) {
    my $buf; $s->recv($buf, 512) or next;
    my ($rid, $flags) = unpack('nn', $buf);
    next unless $rid == $id;
    $a++; $ok++ if ($flags & 0xf) == 0;
  }
}
printf "queries=%d antworten=%d rcode0=%d sekunden=%.2f zeichen=%d\n", $q, $a, $ok, time() - $t0, $chars;
