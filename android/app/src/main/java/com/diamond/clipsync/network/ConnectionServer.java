package com.diamond.clipsync.network;

import android.util.Log;

import java.net.DatagramPacket;
import java.net.DatagramSocket;
import java.net.InetSocketAddress;
import java.nio.charset.StandardCharsets;

public class ConnectionServer {

    private static final String TAG = "ClipSyncServer";

    public interface PacketListener {
        void onHandshakeReceived(String hostname, String ip);
        void onPingReceived(String ip);
        void onClipboardReceived(String content, String ip);
    }

    private final int port;
    private final PacketListener listener;
    private DatagramSocket socket;
    private Thread serverThread;
    private volatile boolean isRunning = false;

    public ConnectionServer(int port, PacketListener listener) {
        this.port = port;
        this.listener = listener;
    }

    public synchronized void start() {
        if (isRunning) return;
        isRunning = true;

        serverThread = new Thread(() -> {
            try {
                socket = new DatagramSocket(null);
                socket.setReuseAddress(true);
                socket.bind(new InetSocketAddress(port));
                Log.i(TAG, "ClipSync UDP Server listening on port " + port);

                byte[] buffer = new byte[65535];
                while (isRunning && !socket.isClosed()) {
                    DatagramPacket packet = new DatagramPacket(buffer, buffer.length);
                    socket.receive(packet);

                    String senderIp = packet.getAddress().getHostAddress();
                    int length = packet.getLength();
                    if (length <= 0) continue;

                    byte[] data = new byte[length];
                    System.arraycopy(packet.getData(), packet.getOffset(), data, 0, length);

                    processPacket(data, senderIp);
                }
            } catch (Exception e) {
                if (isRunning) {
                    Log.e(TAG, "Server socket error: " + e.getMessage());
                }
            } finally {
                stop();
            }
        });
        serverThread.setDaemon(true);
        serverThread.start();
    }

    private void processPacket(byte[] data, String senderIp) {
        byte[] payload = data;

        // Try decrypting if packet looks encrypted (> 16 bytes and header not 0x01/0x02/0x03)
        if (data.length > 16 && data[0] != Protocol.MSG_TYPE_HANDSHAKE
                && data[0] != Protocol.MSG_TYPE_PING
                && data[0] != Protocol.MSG_TYPE_CLIPBOARD) {
            try {
                payload = CryptoUtils.decryptCFB(data, Protocol.SECRET_KEY);
            } catch (Exception ignored) {
                // If CFB decryption fails, keep raw payload
            }
        }

        if (payload.length == 0) return;
        byte msgType = payload[0];

        switch (msgType) {
            case Protocol.MSG_TYPE_HANDSHAKE:
                String hostname = "";
                if (payload.length > 1) {
                    hostname = new String(payload, 1, payload.length - 1, StandardCharsets.UTF_8);
                }
                if (listener != null) {
                    listener.onHandshakeReceived(hostname, senderIp);
                }
                break;

            case Protocol.MSG_TYPE_PING:
                if (listener != null) {
                    listener.onPingReceived(senderIp);
                }
                break;

            case Protocol.MSG_TYPE_CLIPBOARD:
                if (payload.length > 1) {
                    String content = new String(payload, 1, payload.length - 1, StandardCharsets.UTF_8);
                    if (listener != null) {
                        listener.onClipboardReceived(content, senderIp);
                    }
                }
                break;

            default:
                Log.w(TAG, "Unknown message type 0x" + Integer.toHexString(msgType & 0xFF) + " from " + senderIp);
                break;
        }
    }

    public synchronized void stop() {
        isRunning = false;
        if (socket != null && !socket.isClosed()) {
            try {
                socket.close();
            } catch (Exception ignored) {}
            socket = null;
        }
        if (serverThread != null) {
            serverThread.interrupt();
            serverThread = null;
        }
        Log.i(TAG, "ClipSync UDP Server stopped");
    }
}
