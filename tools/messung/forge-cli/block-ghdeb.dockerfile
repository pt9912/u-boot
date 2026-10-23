# HAND-Variante (kein Katalogeintrag, nicht von u-boot erzeugt): gh aus dem Debian-Paket des Basis-Images (Debian 13).
RUN apt-get update \
 && apt-get install -y --no-install-recommends gh \
 && rm -rf /var/lib/apt/lists/*
