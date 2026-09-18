# Keep Gio UI runtime classes and native JNI bindings
-keep class org.gioui.** { *; }
-dontwarn org.gioui.**

# Keep ClipSync Application and packages
-keep class com.diamond.clipsync.** { *; }
-dontwarn com.diamond.clipsync.**

# Keep ClipSync Accessibility Service and background components
-keep class com.clipsync.** { *; }
-dontwarn com.clipsync.**

# Keep all native methods
-keepclasseswithmembernames class * {
    native <methods>;
}
