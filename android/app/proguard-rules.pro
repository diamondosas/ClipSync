# Keep Gio UI runtime classes and native JNI bindings
-keep class org.gioui.** { *; }
-dontwarn org.gioui.**

# Keep ClipSync Accessibility Service and background components
-keep class com.clipsync.** { *; }
-dontwarn com.clipsync.**
