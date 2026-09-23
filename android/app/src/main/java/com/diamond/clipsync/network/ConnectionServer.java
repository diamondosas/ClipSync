package com.diamond.clipsync.network;

import android.util.Log;

import java.net.DatagramPacket;
import java.net.DatagramSocket;
import java.net.InetSocketAddress;
import java.nio.ByteBuffer;
import java.nio.CharBuffer;
import java.nio.charset.CharsetDecoder;
import java.nio.charset.CodingErrorAction;
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
                byte[] dec = CryptoUtils.decryptCFB(data, Protocol.SECRET_KEY);
                if (dec != null && dec.length > 0 && (dec[0] == Protocol.MSG_TYPE_HANDSHAKE
                        || dec[0] == Protocol.MSG_TYPE_PING
                        || dec[0] == Protocol.MSG_TYPE_CLIPBOARD)) {
                    payload = dec;
                }
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
                    hostname = decodeUtf8Safe(payload, 1, payload.length - 1);
                    if (hostname == null) {
                        hostname = new String(payload, 1, payload.length - 1, StandardCharsets.UTF_8);
                    }
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
                    String content = decodeUtf8Safe(payload, 1, payload.length - 1);
                    if (content == null || !isCleanText(content)) {
                        Log.w(TAG, "Discarded binary/scrambled clipboard payload from " + senderIp + " (" + (payload.length - 1) + " bytes)");
                        break;
                    }
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

    private static String decodeUtf8Safe(byte[] data, int offset, int length) {
        try {
            CharsetDecoder decoder = StandardCharsets.UTF_8.newDecoder();
            decoder.onMalformedInput(CodingErrorAction.REPORT);
            decoder.onUnmappableCharacter(CodingErrorAction.REPORT);
            ByteBuffer buf = ByteBuffer.wrap(data, offset, length);
            CharBuffer charBuf = decoder.decode(buf);
            return charBuf.toString();
        } catch (Exception e) {
            return null;
        }
    }

    private static boolean isCleanText(String s) {
        if (s == null || s.isEmpty()) return false;
        int controlCount = 0;
        int length = s.length();
        for (int i = 0; i < length; i++) {
            char c = s.charAt(i);
            if (c == '\n' || c == '\r' || c == '\t') continue;
            if (c == '\uFFFD' || c < 0x20 || (c >= 0x7f && c <= 0x9f)) {
                controlCount++;
            }
        }
        return controlCount == 0 || (length > 10 && ((double) controlCount / length) < 0.05);
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
