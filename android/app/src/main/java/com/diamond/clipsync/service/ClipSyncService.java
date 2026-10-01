package com.diamond.clipsync.service;

import android.app.Service;
import android.content.Context;
import android.content.Intent;
import android.net.wifi.WifiManager;
import android.os.IBinder;
import android.util.Log;

import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;

public class ClipSyncService extends Service {

    private static final String TAG = "ClipSyncService";

    private ServiceCoordinator coordinator;
    private ScheduledExecutorService scheduler;
    private WifiManager.MulticastLock multicastLock;
    private WifiManager.WifiLock wifiLock;

    public static void startService(Context context) {
        Intent intent = new Intent(context, ClipSyncService.class);
        context.startService(intent);
    }

    public static void stopService(Context context) {
        Intent intent = new Intent(context, ClipSyncService.class);
        context.stopService(intent);
    }

    @Override
    public void onCreate() {
        super.onCreate();
        coordinator = ServiceCoordinator.getInstance(this);

        try {
            WifiManager wm = (WifiManager) getApplicationContext().getSystemService(Context.WIFI_SERVICE);
            if (wm != null) {
                multicastLock = wm.createMulticastLock("ClipSyncSvcMulticastLock");
                multicastLock.setReferenceCounted(true);
                multicastLock.acquire();

                wifiLock = wm.createWifiLock(WifiManager.WIFI_MODE_FULL_HIGH_PERF, "ClipSyncSvcWifiLock");
                wifiLock.acquire();
            }
        } catch (Exception e) {
            Log.w(TAG, "Failed acquiring Wi-Fi locks: " + e.getMessage());
        }

        coordinator.startNetworkEngine();
        startHeartbeatLoop();
        Log.i(TAG, "ClipSync Service created (no persistent notification)");
    }

    private void startHeartbeatLoop() {
        scheduler = Executors.newSingleThreadScheduledExecutor();

        // Heartbeat Ping Loop every 2s
        scheduler.scheduleAtFixedRate(() -> {
            try {
                coordinator.pingAllPeers();
            } catch (Exception e) {
                Log.e(TAG, "Ping error: " + e.getMessage());
            }
        }, 1, 2, TimeUnit.SECONDS);

        // Dead Peer Prune Loop every 1s
        scheduler.scheduleAtFixedRate(() -> {
            try {
                coordinator.pruneDeadPeers();
            } catch (Exception e) {
                Log.e(TAG, "Prune error: " + e.getMessage());
            }
        }, 2, 1, TimeUnit.SECONDS);
    }

    @Override
    public int onStartCommand(Intent intent, int flags, int startId) {
        // Do not recreate service if killed by the system
        return START_NOT_STICKY;
    }

    @Override
    public void onTaskRemoved(Intent rootIntent) {
        super.onTaskRemoved(rootIntent);
        Log.i(TAG, "App swiped away from recent apps: stopping ClipSync");
        stopSelf();
    }

    @Override
    public void onDestroy() {
        super.onDestroy();
        if (scheduler != null) {
            scheduler.shutdownNow();
            scheduler = null;
        }
        if (coordinator != null) {
            coordinator.stopNetworkEngine();
        }
        if (multicastLock != null && multicastLock.isHeld()) {
            try {
                multicastLock.release();
            } catch (Exception ignored) {}
            multicastLock = null;
        }
        if (wifiLock != null && wifiLock.isHeld()) {
            try {
                wifiLock.release();
            } catch (Exception ignored) {}
            wifiLock = null;
        }
        Log.i(TAG, "ClipSync Service destroyed");
    }

    @Override
    public IBinder onBind(Intent intent) {
        return null;
    }
}
