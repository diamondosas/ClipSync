package com.diamond.clipsync;

import android.app.Application;
import android.util.Log;
import org.gioui.Gio;

/**
 * ClipSyncApp is the Android Application entry point.
 * It initializes the Gio Go runtime and loads the native libgio.so library
 * before any Activity or Background Service starts.
 */
public class ClipSyncApp extends Application {
    private static final String TAG = "ClipSyncApp";

    @Override
    public void onCreate() {
        super.onCreate();
        Log.i(TAG, "ClipSyncApp onCreate - initializing Gio runtime");
        try {
            Gio.init(this);
            Log.i(TAG, "Gio runtime successfully initialized");
        } catch (Throwable t) {
            Log.e(TAG, "Failed to initialize Gio runtime: " + t.getMessage(), t);
        }
    }
}
