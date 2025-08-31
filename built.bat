@echo off
setlocal

:: =============================================================================
:: Go Application Build Script
:: =============================================================================

:: --- Configuration ---
:: Use %~dp0 to make paths relative to the script's location, not the CWD.
SET "MAIN_FILE=%~dp0cmd\wamp\wamp.go"
SET "OUTPUT_NAME=wamp"

:: --- Script Logic ---

:: Check if an argument was provided.
IF "%~1"=="" (
    ECHO ERROR: No build environment specified.
    ECHO.
    ECHO Usage: %~n0 [dev^|prod]
    GOTO :EOF
)

:: Convert the first argument to lowercase for case-insensitive comparison.
SET "BUILD_ENV=%~1"

ECHO Building for '%BUILD_ENV%' environment...

:: Build for the 'dev' environment.
IF /I "%BUILD_ENV%"=="dev" (
    ECHO Compiling with debug symbols...
    go build -o "%~dp0build\%OUTPUT_NAME%-dev.exe" "%MAIN_FILE%"
    IF %ERRORLEVEL% EQU 0 (
        ECHO.
        ECHO Development build successful! Output: build\%OUTPUT_NAME%-dev.exe
    ) ELSE (
        ECHO.
        ECHO Development build FAILED.
    )
    GOTO :EOF
)

:: Build for the 'prod' environment.
IF /I "%BUILD_ENV%"=="prod" (
    ECHO Compiling and optimizing for production...
    :: -ldflags="-s -w" strips debug information to reduce the binary size.
    go build -ldflags="-s -w" -o "%~dp0build\%OUTPUT_NAME%.exe" "%MAIN_FILE%"
    IF %ERRORLEVEL% EQU 0 (
        ECHO.
        ECHO Production build successful! Output: build\%OUTPUT_NAME%.exe
    ) ELSE (
        ECHO.
        ECHO Production build FAILED.
    )
    GOTO :EOF
)

:: Handle invalid arguments.
ECHO ERROR: Invalid argument '%BUILD_ENV%'.
ECHO.
ECHO Please use 'dev' or 'prod'.
ECHO Usage: %~n0 [dev^|prod]

:EOF
endlocal