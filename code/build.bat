@echo off

go build -o ..\build\child.exe main.go

mkdir ..\build
pushd ..\build

::clang -g -O0 ..\code\main.c -o main.exe
::cl -FC -Zi Kernel32.lib ..\code\main.c /link
cl -FC -Zi Kernel32.lib gdi32.lib msvcrt.lib raylib.lib winmm.lib user32.lib shell32.lib ..\code\main.c ..\code\pipe\pipe_win32.c -I:e:\code\raylib\include /link /libpath:e:\code\raylib\lib /NODEFAULTLIB:libcmt /NODEFAULTLIB:msvcrtd
popd

:: gcc -o PokeTCGScraper.exe poke_tcg_scraper.c  -I include -L lib -lraylib -lgdi32 -lwinmm
