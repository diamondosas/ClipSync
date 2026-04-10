@echo off
setlocal enabledelayedexpansion

:menu
cls
echo =======================================================
echo     ClipSync Android Build System
echo =======================================================
echo  1. RUN ALL (Full Build: AAR + APK)
echo  2. [Step 1] Compile Go to AAR (AAR only)
echo  3. [Step 2] Build Android APK (Gradle only)
echo  4. Exit
echo =======================================================
set /p choice="Select an option (1-4): "

if "%choice%"=="1" goto run_full_build
if "%choice%"=="2" (
    call :compile_aar
    pause
    goto menu
)
if "%choice%"=="3" (
    call :build_apk
    pause
    goto menu
)
if "%choice%"=="4" exit /b 0

echo Invalid choice. Try again.
pause
goto menu

:run_full_build
echo.
echo === Starting Full Build ===
call :compile_aar
if !errorlevel! neq 0 (
    echo.
    echo [ERROR] Step 1 failed. Aborting full build.
    pause
    goto menu
)
call :build_apk
if !errorlevel! neq 0 (
    echo.
    echo [ERROR] Step 2 failed. Aborting full build.
    pause
    goto menu
)
echo.
echo === Full Build SUCCESS ===
pause
goto menu

:: --- SUBROUTINES ---

:compile_aar
echo.
echo =======================================================
echo [1/3] Compiling Go to Android (AAR)...
echo =======================================================
:: Build the main app AAR
call gogio -target android -buildmode archive -appid com.diamond.clipsync -o android/app/libs/clipsync.aar ./
set "AAR_ERR=!errorlevel!"
if !AAR_ERR! neq 0 (
    echo.
    echo [ERROR] gogio compilation failed. 
    echo Please ensure gogio is installed and NDK is configured.
    exit /b !AAR_ERR!
)

:: Build the learn app AAR (secondary task)
echo Compiling Learn AAR...
call gogio -target android -buildmode archive -appid com.diamond.clipsync -o learn/learn.aar ./learn/learn.go

echo.
echo [Step 1] SUCCESS: AAR generated in android/app/libs/clipsync.aar
exit /b 0

:build_apk
if not exist "android\app\libs\clipsync.aar" (
    echo.
    echo [ERROR] clipsync.aar not found in android\app\libs\
    echo Please run Step 1 (AAR compilation) first.
    exit /b 1
)

echo.
echo =======================================================
echo [2/3] Building Android APK (Gradle)...
echo =======================================================
cd android

:: Use gradlew if available, otherwise fallback to system gradle
if exist "gradlew.bat" (
    call gradlew.bat assembleDebug
) else (
    echo [INFO] gradlew.bat not found. Attempting to use installed 'gradle' command...
    call gradle assembleDebug
)
set "GRADLE_ERR=!errorlevel!"
cd ..

if !GRADLE_ERR! neq 0 (
    echo.
    echo [ERROR] Gradle build failed. Check logs above.
    exit /b !GRADLE_ERR!
)

echo.
echo =======================================================
echo [3/3] Finalizing APK...
echo =======================================================
set "APK_PATH=android\app\build\outputs\apk\debug\app-debug.apk"

if exist "%APK_PATH%" (
    copy "%APK_PATH%" "clipsync-debug.apk" /y
    echo.
    echo [SUCCESS] APK copied to root as: clipsync-debug.apk
) else (
    echo [ERROR] Could not find generated APK at %APK_PATH%
    exit /b 1
)
exit /b 0
