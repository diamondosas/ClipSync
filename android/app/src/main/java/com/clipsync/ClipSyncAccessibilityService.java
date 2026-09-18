package com.clipsync;

import android.accessibilityservice.AccessibilityService;
import android.content.ClipData;
import android.content.ClipboardManager;
import android.content.Context;
import android.net.wifi.WifiManager;
import android.os.Handler;
import android.os.Looper;
import android.util.Log;
import android.view.accessibility.AccessibilityEvent;

import java.net.DatagramPacket;
import java.net.DatagramSocket;
import java.net.InetAddress;
import java.nio.charset.StandardCharsets;

/**
 * ClipSyncAccessibilityService runs continuously in the background on Android.
 * It monitors system clipboard changes across all apps and forwards copied text
 * to the Go ClipSync background engine via localhost UDP (127.0.0.1:9998).
 */
public class ClipSyncAccessibilityService extends AccessibilityService {

    private static final String TAG = "ClipSyncAccessService";
    private static final int LOCAL_BRIDGE_PORT = 9998;

    private ClipboardManager clipboardManager;
    private ClipboardManager.OnPrimaryClipChangedListener clipListener;
    private WifiManager.MulticastLock multicastLock;
    private String lastCopiedText = "";
    private final Handler mainHandler = new Handler(Looper.getMainLooper());

    @Override
    public void onServiceConnected() {
        super.onServiceConnected();
        Log.i(TAG, "ClipSync Accessibility Service connected");

        // Ensure Gio Go runtime is initialized if service started before UI
        try {
            org.gioui.Gio.init(getApplicationContext());
        } catch (Throwable t) {
            Log.w(TAG, "Gio.init in AccessibilityService: " + t.getMessage());
        }

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
            clipListener = new ClipboardManager.OnPrimaryClipChangedListener() {
                @Override
                public void onPrimaryClipChanged() {
                    checkAndForwardClipboard();
                }
            };
            clipboardManager.addPrimaryClipChangedListener(clipListener);
            Log.i(TAG, "Clipboard change listener registered");
        }
    }

    @Override
    public void onAccessibilityEvent(AccessibilityEvent event) {
        // Fallback: check clipboard on window/focus changes
        if (event != null) {
            int eventType = event.getEventType();
            if (eventType == AccessibilityEvent.TYPE_WINDOW_STATE_CHANGED ||
                eventType == AccessibilityEvent.TYPE_VIEW_CLICKED) {
                mainHandler.postDelayed(new Runnable() {
                    @Override
                    public void run() {
                        checkAndForwardClipboard();
                    }
                }, 150);
            }
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
                CharSequence text = clipData.getItemAt(0).getText();
                if (text != null && text.length() > 0) {
                    String str = text.toString();
                    if (!str.equals(lastCopiedText)) {
                        lastCopiedText = str;
                        sendToLocalGoEngine(str);
                    }
                }
            }
        } catch (Exception e) {
            Log.e(TAG, "Error reading clipboard: " + e.getMessage());
        }
    }

    private void sendToLocalGoEngine(final String text) {
        new Thread(new Runnable() {
            @Override
            public void run() {
                try (DatagramSocket socket = new DatagramSocket()) {
                    byte[] data = text.getBytes(StandardCharsets.UTF_8);
                    InetAddress localhost = InetAddress.getByName("127.0.0.1");
                    DatagramPacket packet = new DatagramPacket(data, data.length, localhost, LOCAL_BRIDGE_PORT);
                    socket.send(packet);
                    Log.i(TAG, "Forwarded " + data.length + " bytes to Go engine");
                } catch (Exception e) {
                    Log.e(TAG, "Failed to send packet to Go engine: " + e.getMessage());
                }
            }
        }).start();
    }

    @Override
    public void onDestroy() {
        super.onDestroy();
        if (clipboardManager != null && clipListener != null) {
            clipboardManager.removePrimaryClipChangedListener(clipListener);
        }
        if (multicastLock != null && multicastLock.isHeld()) {
            multicastLock.release();
        }
        Log.i(TAG, "ClipSync Accessibility Service destroyed");
    }
}
