@echo off
setlocal

:: Check if the AAR exists before trying to build
if not exist "android\app\libs\clipsync.aar" (
    echo [ERROR] clipsync.aar not found in android\app\libs
    echo Please run compile.bat first to generate the AAR.
    exit /b 1
)

echo.
echo [1/2] Building Android Debug APK...
cd android

:: Try using gradlew if it exists, otherwise fall back to gradle command
if exist "gradlew.bat" (
    call gradlew.bat assembleDebug
) else (
    echo [INFO] gradlew.bat not found. Attempting to use installed 'gradle' command...
    call gradle assembleDebug
)

if %ERRORLEVEL% neq 0 (
    echo [ERROR] Gradle build failed.
    cd ..
    exit /b %ERRORLEVEL%
)

cd ..

:: Locate and copy the generated APK to the root
echo.
echo [2/2] Locating generated APK...
set "APK_PATH=android\app\build\outputs\apk\debug\app-debug.apk"

if exist "%APK_PATH%" (
    copy "%APK_PATH%" "clipsync-debug.apk" /y
    echo.
    echo [SUCCESS] APK built and copied to: clipsync-debug.apk
) else (
    echo [ERROR] Could not find generated APK at %APK_PATH%
)

endlocal
