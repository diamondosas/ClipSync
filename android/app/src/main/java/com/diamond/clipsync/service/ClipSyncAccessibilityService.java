package com.diamond.clipsync.service;

import android.accessibilityservice.AccessibilityService;
import android.content.ClipData;
import android.content.ClipboardManager;
import android.content.Context;
import android.net.wifi.WifiManager;
import android.os.Handler;
import android.os.Looper;
import android.util.Log;
import android.view.accessibility.AccessibilityEvent;

public class ClipSyncAccessibilityService extends AccessibilityService {

    private static final String TAG = "ClipSyncAccessibility";
    private ClipboardManager clipboardManager;
    private ClipboardManager.OnPrimaryClipChangedListener clipListener;
    private WifiManager.MulticastLock multicastLock;
    private final Handler mainHandler = new Handler(Looper.getMainLooper());
    private ServiceCoordinator coordinator;

    private static volatile boolean isServiceRunning = false;

    public static boolean isRunning() {
        return isServiceRunning;
    }

    @Override
    public void onServiceConnected() {
        super.onServiceConnected();
        isServiceRunning = true;
        Log.i(TAG, "ClipSync Accessibility Service connected");

        coordinator = ServiceCoordinator.getInstance(this);

        // 1. Acquire MulticastLock for Zeroconf / mDNS Wi-Fi discovery
        try {
            WifiManager wifi = (WifiManager) getApplicationContext().getSystemService(Context.WIFI_SERVICE);
            if (wifi != null) {
                multicastLock = wifi.createMulticastLock("ClipSyncMulticastLock");
                multicastLock.setReferenceCounted(true);
                multicastLock.acquire();
                Log.i(TAG, "Wi-Fi MulticastLock acquired");
            }
        } catch (Exception e) {
            Log.w(TAG, "Could not acquire MulticastLock: " + e.getMessage());
        }

        // 2. Register Clipboard Listener
        clipboardManager = (ClipboardManager) getSystemService(Context.CLIPBOARD_SERVICE);
        if (clipboardManager != null) {
            clipListener = this::checkAndForwardClipboard;
            clipboardManager.addPrimaryClipChangedListener(clipListener);
            Log.i(TAG, "Clipboard change listener registered");
        }

        // Ensure foreground network service is running
        ClipSyncForegroundService.startService(this);
    }

    @Override
    public void onAccessibilityEvent(AccessibilityEvent event) {
        if (event == null) return;

        int eventType = event.getEventType();
        if (eventType == AccessibilityEvent.TYPE_WINDOW_STATE_CHANGED
                || eventType == AccessibilityEvent.TYPE_VIEW_CLICKED
                || eventType == AccessibilityEvent.TYPE_VIEW_FOCUSED
                || eventType == AccessibilityEvent.TYPE_VIEW_TEXT_SELECTION_CHANGED) {
            // Debounced check to allow clipboard buffer to finalize
            mainHandler.postDelayed(this::checkAndForwardClipboard, 150);
        }
    }

    @Override
    public void onInterrupt() {
        Log.w(TAG, "ClipSync Accessibility Service interrupted");
    }

    private synchronized void checkAndForwardClipboard() {
        if (clipboardManager == null || !clipboardManager.hasPrimaryClip()) {
            return;
        }

        try {
            ClipData clipData = clipboardManager.getPrimaryClip();
            if (clipData != null && clipData.getItemCount() > 0) {
                ClipData.Item item = clipData.getItemAt(0);
                CharSequence text = item.getText();
                if (text != null && text.length() > 0) {
                    String str = text.toString();
                    if (coordinator != null) {
                        coordinator.onLocalClipboardCopied(str);
                    }
                }
            }
        } catch (Exception e) {
            Log.e(TAG, "Error reading clipboard: " + e.getMessage());
        }
    }

    @Override
    public void onDestroy() {
        super.onDestroy();
        isServiceRunning = false;
        if (clipboardManager != null && clipListener != null) {
            clipboardManager.removePrimaryClipChangedListener(clipListener);
        }
        if (multicastLock != null && multicastLock.isHeld()) {
            try {
                multicastLock.release();
            } catch (Exception ignored) {}
        }
        Log.i(TAG, "ClipSync Accessibility Service destroyed");
    }
}
