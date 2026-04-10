@echo off
echo Compiling Main App for Android (AAR)...
gogio -target android -buildmode archive -appid com.diamond.clipsync -o android/app/libs/clipsync.aar ./
echo.
echo Done.