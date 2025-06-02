# ProgApChZAD2

Treść zadania

Opracować łańcuch (pipeline) w usłudzie GitHub Actions, który zbuduje obraz kontenera na podstawie Dockerfile-a oraz kodów źródłowych aplikacji opracowanej jako rozwiązanie zadania nr 1 a następnie prześle go do publicznego repozytorium autora na Github (ghcr.io). Proces budowania 
obrazu opisany w łańcuchu GHAction powinien dodatkowo spełniać następujące warunki: 
a. Obraz wspierać ma dwie architektury: linux/arm64 oraz linux/amd64. 
b. Wykorzystywane mają być (wysyłanie i pobieranie) dane cache (eksporter: registry oraz backend-u registry w trybie max). Te dane cache powinny być przechowywane w dedykowanym, publicznym repozytorium autora na DockerHub.       c. Ma być wykonany test CVE obrazu, który zapewni, że obraz zostanie przesłany do publicznego repozytorium obrazów na GitHub tylko wtedy gdy nie będzie zawierał zagrożeń sklasyfikowanych jako krytyczne lub wysokie. 
W opisie rozwiązania należy krótko przedstawić przyjęty sposób tagowania obrazów i danych cache. 
Uzasadnienie (z ewentualnym powołaniem się na źródła) tego wyboru będzie „nagrodzone” dodatkowymi punktami. 

Etapy pracy:

1. Uruchomienie workflow
Pipeline uruchamia się automatycznie w dwóch przypadkach:
- Push do gałęzi main
- Ręczne uruchomienie z poziomu GitHub UI (workflow_dispatch)
![image](https://github.com/user-attachments/assets/39f7562c-c988-4a97-bf93-5a415ae27359)
2. (linux/amd64, linux/arm64)
  ![image](https://github.com/user-attachments/assets/e9b1b218-e216-40fb-96df-720ca8dcdefa)
3. Strategia tagowania i konfiguracja cache
  
- Tag latest – wskazuje zawsze najnowszy obraz z gałęzi main, typowe dla środowisk deweloperskich.

- Oddzielenie obrazu i cache – umożliwia lepsze zarządzanie oraz przyspiesza budowanie dzięki wykorzystywaniu warstw cache.

Dwa rejestry:

- GHCR (GitHub Container Registry) – zintegrowany z repozytorium, bezpieczne i automatyczne uwierzytelnianie.

- DockerHub – wykorzystywany tylko do przechowywania cache, co zmniejsza wykorzystanie limitów w GHCR.

Uzasadnienie wyboru cache:

- Szybsze budowanie obrazów – warstwy Dockera, które nie zmieniły się między kolejnymi buildami, nie są budowane od nowa.

- Oszczędność zasobów CI – krótszy czas działania workflow i mniejsze zużycie minut GitHub Actions.

- Wydajniejsze wykorzystanie rejestru – warstwy obrazu są przechowywane w zewnętrznym rejestrze (DockerHub), co pozwala na ich ponowne użycie w kolejnych uruchomieniach pipeline’u.

- Stabilność buildów – powtarzalność i deterministyczne wyniki budowania, niezależne od zmian w środowisku.

4.Skanowanie CVE (Trivy)
![image](https://github.com/user-attachments/assets/581a6499-d8ca-499f-9ffd-22f3d6276769)
- Lekkość i szybkość działania – szybki czas skanowania, niskie wymagania zasobowe.

- Skanowanie bezpośrednio z obrazu Dockera – nie trzeba uruchamiać kontenera, wystarczy wskazać tag.

- Wysoka dokładność i aktualna baza CVE – Trivy korzysta z aktualizowanej bazy podatności z NVD, Red Hat, GitHub Advisory Database itd.

- Integracja z GitHub Actions – dostępna oficjalna akcja, łatwa konfiguracja.

5. Sekrety repozytorium
- DOCKER_USERNAME - nazwa użytkownika DockerHub
- DOCKERHUB_TOKEN - token dostępu DockerHub
- GITHUB_TOKEN - automatycznie dostępny w GitHub Actions

6. Weryfikacja działania pipeline’u
- Pipeline uruchamia się automatycznie przy pushu na main lub ręcznie przez GitHub Actions.

- Obraz Dockera budowany jest dla architektur linux/amd64 i linux/arm64 z użyciem QEMU i Buildx.

- Cache jest pobierany i zapisywany w repozytorium na DockerHub, co przyspiesza kolejne buildy.

- Obraz skanowany jest narzędziem Trivy pod kątem krytycznych i wysokich luk bezpieczeństwa.

- Obraz jest publikowany w GHCR tylko, gdy skan zakończy się sukcesem (brak krytycznych podatności).
