package com.diamond.clipsync.network;

import android.util.Log;

import com.diamond.clipsync.data.model.Device;

import java.net.DatagramPacket;
import java.net.DatagramSocket;
import java.net.InetAddress;
import java.util.List;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

public class ConnectionClient {

    private static final String TAG = "ClipSyncClient";
    private final PeerManager peerManager;
    private final int defaultPort;
    private final String localDeviceName;
    private final ExecutorService sendExecutor;

    public ConnectionClient(PeerManager peerManager, int defaultPort, String localDeviceName) {
        this.peerManager = peerManager;
        this.defaultPort = defaultPort;
        this.localDeviceName = localDeviceName;
        this.sendExecutor = Executors.newCachedThreadPool();
    }

    public void sendHandshake(final String ip, final int port) {
        sendExecutor.execute(() -> {
            byte[] packet = Protocol.encodeHandshake(localDeviceName);
            sendPacket(ip, port > 0 ? port : defaultPort, packet);
        });
    }

    public void sendPing(final String ip, final int port) {
        sendExecutor.execute(() -> {
            byte[] packet = Protocol.encodePing();
            sendPacket(ip, port > 0 ? port : defaultPort, packet);
        });
    }

    public void pingAllPeers() {
        if (peerManager == null) return;
        List<Device> activePeers = peerManager.getActivePeers();
        for (Device peer : activePeers) {
            sendPing(peer.getIp(), peer.getPort());
        }
    }

    public void broadcastClipboard(final String content) {
        if (peerManager == null || content == null || content.isEmpty()) return;
        final List<Device> activePeers = peerManager.getActivePeers();
        if (activePeers.isEmpty()) {
            Log.d(TAG, "No active peers to broadcast clipboard to");
            return;
        }

        sendExecutor.execute(() -> {
            byte[] packet = Protocol.encodeClipboard(content);
            for (Device peer : activePeers) {
                sendPacket(peer.getIp(), peer.getPort() > 0 ? peer.getPort() : defaultPort, packet);
            }
        });
    }

    private void sendPacket(String ip, int port, byte[] data) {
        try (DatagramSocket socket = new DatagramSocket()) {
            socket.setSoTimeout(3000);
            InetAddress targetAddr = InetAddress.getByName(ip);

            // Send standard raw packet
            DatagramPacket outPacket = new DatagramPacket(data, data.length, targetAddr, port);
            socket.send(outPacket);

            // Also send AES CFB encrypted variant if payload is clipboard or handshake for compatibility
            try {
                byte[] enc = CryptoUtils.encryptCFB(data, Protocol.SECRET_KEY);
                DatagramPacket encPacket = new DatagramPacket(enc, enc.length, targetAddr, port);
                socket.send(encPacket);
            } catch (Exception ignored) {}

            Log.d(TAG, "Sent packet (" + data.length + " bytes) to " + ip + ":" + port);
        } catch (Exception e) {
            Log.e(TAG, "Failed sending packet to " + ip + ":" + port + " -> " + e.getMessage());
        }
    }

    public void shutdown() {
        sendExecutor.shutdown();
    }
}
