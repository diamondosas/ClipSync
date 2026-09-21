package com.diamond.clipsync.service;

import android.content.ClipData;
import android.content.ClipboardManager;
import android.content.Context;
import android.os.Build;
import android.os.Handler;
import android.os.Looper;
import android.util.Log;
import android.widget.Toast;

import com.diamond.clipsync.data.model.Device;
import com.diamond.clipsync.data.repository.ClipRepository;
import com.diamond.clipsync.network.ConnectionClient;
import com.diamond.clipsync.network.ConnectionServer;
import com.diamond.clipsync.network.CryptoUtils;
import com.diamond.clipsync.network.NsdHelper;
import com.diamond.clipsync.network.PeerManager;
import com.diamond.clipsync.network.Protocol;

import java.util.List;

public class ServiceCoordinator implements ConnectionServer.PacketListener {

    private static final String TAG = "ClipSyncCoordinator";
    private static volatile ServiceCoordinator INSTANCE;

    private final Context context;
    private final ClipRepository repository;
    private final PeerManager peerManager;
    private final ConnectionClient connectionClient;
    private final ConnectionServer connectionServer;
    private final NsdHelper nsdHelper;
    private final ClipboardManager clipboardManager;
    private final Handler mainHandler = new Handler(Looper.getMainLooper());

    private final String localDeviceName;

    // Two-Way Echo / Ping-Pong Cancellation State
    private volatile String lastProcessedLocalHash = "";
    private volatile String lastReceivedRemoteHash = "";
    private volatile long lastRemoteReceiveTime = 0L;
    private static final long ECHO_GRACE_PERIOD_MS = 2000L;

    public interface StatusListener {
        void onToast(String message);
    }

    private StatusListener statusListener;

    private ServiceCoordinator(Context context) {
        this.context = context.getApplicationContext();
        this.repository = ClipRepository.getInstance(this.context);
        this.peerManager = new PeerManager();
        this.clipboardManager = (ClipboardManager) this.context.getSystemService(Context.CLIPBOARD_SERVICE);

        String model = Build.MODEL != null ? Build.MODEL : "Android";
        this.localDeviceName = model + "-ClipSync";

        this.connectionClient = new ConnectionClient(peerManager, Protocol.DEFAULT_PORT, localDeviceName);
        this.connectionServer = new ConnectionServer(Protocol.DEFAULT_PORT, this);
        this.nsdHelper = new NsdHelper(this.context, peerManager, localDeviceName, Protocol.DEFAULT_PORT);

        this.nsdHelper.setPeerCallback((name, ip, port) -> {
            // Automatically send handshake when discovering peer
            connectionClient.sendHandshake(ip, port);
        });
    }

    public static ServiceCoordinator getInstance(Context context) {
        if (INSTANCE == null) {
            synchronized (ServiceCoordinator.class) {
                if (INSTANCE == null) {
                    INSTANCE = new ServiceCoordinator(context.getApplicationContext());
                }
            }
        }
        return INSTANCE;
    }

    public void setStatusListener(StatusListener listener) {
        this.statusListener = listener;
    }

    public void startNetworkEngine() {
        connectionServer.start();
        nsdHelper.start();
        Log.i(TAG, "Network engine started");
    }

    public void stopNetworkEngine() {
        nsdHelper.stop();
        connectionServer.stop();
        Log.i(TAG, "Network engine stopped");
    }

    public PeerManager getPeerManager() {
        return peerManager;
    }

    public ClipRepository getRepository() {
        return repository;
    }

    public void connectManual(String ip) {
        if (ip == null || ip.trim().isEmpty()) return;
        String cleanIp = ip.trim();
        peerManager.addOrUpdatePeer(cleanIp, cleanIp, Protocol.DEFAULT_PORT);
        connectionClient.sendHandshake(cleanIp, Protocol.DEFAULT_PORT);
        showToast("Connecting to " + cleanIp + "…");
    }

    /**
     * Invoked when AccessibilityService or local app copies new text.
     */
    public synchronized void onLocalClipboardCopied(String content) {
        if (content == null || content.isEmpty()) return;

        String hash = CryptoUtils.sha256(content);
        long now = System.currentTimeMillis();

        // 1. Echo filter: Ignore if this exact clip was received from a remote peer recently
        if (now - lastRemoteReceiveTime < ECHO_GRACE_PERIOD_MS && hash.equals(lastReceivedRemoteHash)) {
            Log.d(TAG, "Echo suppressed: clip was received from remote peer recently");
            return;
        }

        // 2. Deduplication filter: Ignore if identical to last processed local clip
        if (hash.equals(lastProcessedLocalHash)) {
            return;
        }
        lastProcessedLocalHash = hash;

        Log.i(TAG, "Local clip copied (" + content.length() + " chars): broadcasting to peers");

        // 3. Store in Room database
        repository.insertClip(content, "Phone", "127.0.0.1");

        // 4. Broadcast to connected peers
        connectionClient.broadcastClipboard(content);

        showToast("Phone clip synced to devices");
    }

    // --- PacketListener Callbacks (from ConnectionServer) ---

    @Override
    public void onHandshakeReceived(String hostname, String ip) {
        Log.i(TAG, "Handshake received from: " + hostname + " (" + ip + ")");
        peerManager.addOrUpdatePeer(hostname, ip, Protocol.DEFAULT_PORT);
        // Reply with handshake if needed
        connectionClient.sendPing(ip, Protocol.DEFAULT_PORT);
    }

    @Override
    public void onPingReceived(String ip) {
        peerManager.touchPeer(ip);
    }

    @Override
    public void onClipboardReceived(String content, String ip) {
        if (content == null || content.isEmpty()) return;

        String hash = CryptoUtils.sha256(content);
        lastReceivedRemoteHash = hash;
        lastRemoteReceiveTime = System.currentTimeMillis();
        lastProcessedLocalHash = hash;

        Log.i(TAG, "Received clipboard (" + content.length() + " chars) from " + ip);

        // 1. Save to Room database
        String peerName = ip;
        for (Device d : peerManager.getActivePeers()) {
            if (d.getIp().equals(ip)) {
                peerName = d.getName();
                break;
            }
        }
        repository.insertClip(content, peerName, ip);

        // 2. Write to system clipboard on UI thread
        mainHandler.post(() -> {
            try {
                if (clipboardManager != null) {
                    ClipData clip = ClipData.newPlainText("ClipSync", content);
                    clipboardManager.setPrimaryClip(clip);
                }
            } catch (Exception e) {
                Log.e(TAG, "Error writing to clipboard: " + e.getMessage());
            }
        });

        showToast("Synced new clip from " + peerName);
    }

    public void pingAllPeers() {
        connectionClient.pingAllPeers();
    }

    public void pruneDeadPeers() {
        peerManager.pruneDeadPeers(Protocol.PEER_TIMEOUT_MS);
    }

    private void showToast(final String message) {
        mainHandler.post(() -> {
            if (statusListener != null) {
                statusListener.onToast(message);
            } else {
                try {
                    Toast.makeText(context, message, Toast.LENGTH_SHORT).show();
                } catch (Exception ignored) {}
            }
        });
    }
}
