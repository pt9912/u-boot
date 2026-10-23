# HAND-Variante (kein Katalogeintrag, nicht von u-boot erzeugt): glab aus dem Debian-Paket des Basis-Images (Debian 13).
RUN apt-get update \
 && apt-get install -y --no-install-recommends glab \
 && rm -rf /var/lib/apt/lists/*
