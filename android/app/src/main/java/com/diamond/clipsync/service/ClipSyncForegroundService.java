package com.diamond.clipsync.service;

import android.app.Notification;
import android.app.NotificationChannel;
import android.app.NotificationManager;
import android.app.PendingIntent;
import android.app.Service;
import android.content.Context;
import android.content.Intent;
import android.net.wifi.WifiManager;
import android.os.Build;
import android.os.IBinder;
import android.util.Log;

import androidx.core.app.NotificationCompat;

import com.diamond.clipsync.R;
import com.diamond.clipsync.data.model.Device;
import com.diamond.clipsync.ui.MainActivity;

import java.util.List;
import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;

public class ClipSyncForegroundService extends Service {

    private static final String TAG = "ClipSyncForegroundSvc";
    private static final String CHANNEL_ID = "clipsync_foreground_channel";
    private static final int NOTIFICATION_ID = 1001;

    private ServiceCoordinator coordinator;
    private ScheduledExecutorService scheduler;
    private NotificationManager notificationManager;
    private WifiManager.MulticastLock multicastLock;
    private WifiManager.WifiLock wifiLock;

    public static void startService(Context context) {
        Intent intent = new Intent(context, ClipSyncForegroundService.class);
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            context.startForegroundService(intent);
        } else {
            context.startService(intent);
        }
    }

    public static void stopService(Context context) {
        Intent intent = new Intent(context, ClipSyncForegroundService.class);
        context.stopService(intent);
    }

    @Override
    public void onCreate() {
        super.onCreate();
        coordinator = ServiceCoordinator.getInstance(this);
        notificationManager = (NotificationManager) getSystemService(Context.NOTIFICATION_SERVICE);

        createNotificationChannel();
        startForeground(NOTIFICATION_ID, buildNotification(0));

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

        coordinator.getPeerManager().addListener(this::updateNotificationPeerCount);

        startHeartbeatLoop();
        Log.i(TAG, "ClipSync Foreground Service created");
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

    private void createNotificationChannel() {
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            NotificationChannel channel = new NotificationChannel(
                    CHANNEL_ID,
                    getString(R.string.notification_channel_name),
                    NotificationManager.IMPORTANCE_LOW
            );
            channel.setDescription(getString(R.string.notification_channel_desc));
            channel.setShowBadge(false);
            if (notificationManager != null) {
                notificationManager.createNotificationChannel(channel);
            }
        }
    }

    private Notification buildNotification(int peerCount) {
        Intent notificationIntent = new Intent(this, MainActivity.class);
        notificationIntent.setFlags(Intent.FLAG_ACTIVITY_SINGLE_TOP);
        PendingIntent pendingIntent = PendingIntent.getActivity(
                this, 0, notificationIntent,
                PendingIntent.FLAG_IMMUTABLE | PendingIntent.FLAG_UPDATE_CURRENT
        );

        String contentText = peerCount == 0
                ? getString(R.string.notification_content)
                : getString(R.string.notification_active_peers, peerCount);

        return new NotificationCompat.Builder(this, CHANNEL_ID)
                .setContentTitle(getString(R.string.notification_title))
                .setContentText(contentText)
                .setSmallIcon(R.drawable.ic_clipboard)
                .setContentIntent(pendingIntent)
                .setOngoing(true)
                .setPriority(NotificationCompat.PRIORITY_LOW)
                .build();
    }

    private void updateNotificationPeerCount(List<Device> activePeers) {
        int count = activePeers != null ? activePeers.size() : 0;
        if (notificationManager != null) {
            notificationManager.notify(NOTIFICATION_ID, buildNotification(count));
        }
    }

    @Override
    public int onStartCommand(Intent intent, int flags, int startId) {
        return START_STICKY;
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
        Log.i(TAG, "ClipSync Foreground Service destroyed");
    }

    @Override
    public IBinder onBind(Intent intent) {
        return null;
    }
}
