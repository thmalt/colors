@echo off

setlocal

if "%~1"=="all" (
    set GEN_LUT=1
)

go generate ./gen