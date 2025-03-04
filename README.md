# slitherlink_solver

Przy testowaniu programu korzystano z Ubuntu 20.04.6 LTS, dlatego poniższa instrukcja
przeznaczona jest na ten właśnie system operacyjny.

W celu uruchomienia programu należy zainstalować język Go w wersji 1.23.2 (można w
tym celu skorzystać ze strony https://go.dev/doc/install). Po upewnieniu się, że język został
zainstalowany poprawnie, należy w głównym folderze projektu uruchomić komendę
go build .

W ten sposób powinien zostać wygenerowany program slitherlink_solver. Aby go
uruchomić, należy wywołać komendę

```
./slitherlink_solver arg1 arg2 arg3 arg4 arg5
```

gdzie jako wartości argumentów należy podać:

* `arg1` – służy do wyboru wariantu uruchamianych łamigłówek. Program działa dla
poniższych wartości:
  - `s` – plansza zbudowana z kwadratów
  - `h` – plansza zbudowana z sześciokątów foremnych
  - `t` – plansza zbudowana z trójkątów foremnych

* `arg2` – służy do uruchomienia programu z heurystyką lub bez. Może przyjmować wartości
on (heurystyka włączona) lub off (heurystyka wyłączona).

* `arg3` – służy do wyboru numeru heurystyki z zakresu 0-4. Rodzaje heurystyk opisano
w rozdziale 6. W przypadku wyboru heurystyki o numerze 0, algorytm nie będzie
wyszukiwał wzorców.

* `arg4` – argument opcjonalny. Służy do uruchamiania pojedynczej instancji zagadki. Jako
wartość przyjmuje numer zagadki do uruchomienia.

* `arg5` – argument opcjonalny. Jeśli przyjmie wartość d, program zostanie uruchomiony
w trybie debugowania. W takim wypadku drukowane będą wszystkie logi, a każdy
odwiedzony stan planszy zostanie wyświetlony w konsoli. Ponadto po każdej iteracji
program zatrzyma się na 1000ms w celu ułatwienia zapoznawania się z danymi.

Jako wynik działania programu, w konsoli zostanie wyświetlona plansza ze znalezionym
rozwiązaniem, a także dodatkowe informacje o działaniu algorytmu: czas działania,
liczba odwiedzonych stanów oraz średnia i maksymalna głębokość przeszukiwanego drzewa
rozwiązań.
