# Etap 1: Budowanie

# Warstwa 1: Bazowy obraz buildowy — zawiera kompilator Go oraz Alpine Linux
FROM golang:1.23.8-alpine AS builder

# Warstwa 2: Ustawienia środowiska
# CGO_ENABLED=0 - wyłącza zależności C
# GOOS=linux, GOARCH=amd64 - targetujemy Linuksa 64-bit
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Warstwa 3: Informacja o autorze obrazu (standard OCI)
LABEL org.opencontainers.image.authors="Wiktoria Matacz"

# Warstwa 4: Katalog roboczy w kontenerze
WORKDIR /app

# Warstwa 5: Kopiujemy pliki zależności, by cache był lepiej wykorzystywany
COPY go.mod ./
RUN go mod download

# Warstwa 6: Kopiujemy cały projekt do kontenera
COPY . .

# Warstwa 7: Kompilacja aplikacji do pliku binarnego
RUN go build -o app_pogoda .


#Etap 2: Obraz

# Warstwa 8: Scratch — minimalny obraz.
FROM scratch

# Warstwa 9: Informacja o autorze (dla końcowego obrazu)
LABEL org.opencontainers.image.authors="Wiktoria Matacz"

# Warstwa 10: Kopiujemy skompilowaną binarkę z etapu `builder`
COPY --from=builder /app/app_pogoda /app_pogoda

# Warstwa 11: Dodajemy certyfikaty SSL (wymagane przez http.Get dla HTTPS)
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Warstwa 12: Otwieramy port aplikacji
EXPOSE 8080

# Warstwa 13: HEALTHCHECK sprawdza, czy aplikacja działa (ping na localhost:8080 co 30 sekund)
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
  CMD wget --spider -q http://localhost:8080/ || exit 1

# Warstwa 14: ENTRYPOINT ustawia aplikację jako główny proces kontenera
ENTRYPOINT ["/app_pogoda"]