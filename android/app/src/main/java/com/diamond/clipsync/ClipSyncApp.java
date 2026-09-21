package com.diamond.clipsync;

import android.app.Application;
import android.util.Log;

import com.diamond.clipsync.data.db.AppDatabase;
import com.diamond.clipsync.service.ServiceCoordinator;

public class ClipSyncApp extends Application {

    private static final String TAG = "ClipSyncApp";

    @Override
    public void onCreate() {
        super.onCreate();
        Log.i(TAG, "Initializing ClipSync Application");

        // Pre-initialize Room Database
        AppDatabase.getInstance(this);

        // Pre-initialize ServiceCoordinator
        ServiceCoordinator.getInstance(this);
    }
}
