package com.diamond.clipsync.network;

import java.nio.charset.StandardCharsets;

public class Protocol {
    public static final int DEFAULT_PORT = 9999;
    public static final String SERVICE_TYPE = "_clipsync._tcp";
    public static final String SERVICE_DOMAIN = "local.";
    public static final long PING_INTERVAL_MS = 2000L;
    public static final long PEER_TIMEOUT_MS = 5000L;

    // Default shared secret key (16 bytes) matching desktop and android
    public static final byte[] SECRET_KEY = "clipboardsyncapp".getBytes(StandardCharsets.UTF_8);

    // Wire message type headers
    public static final byte MSG_TYPE_HANDSHAKE = 0x01;
    public static final byte MSG_TYPE_PING      = 0x02;
    public static final byte MSG_TYPE_CLIPBOARD = 0x03;

    public static byte[] encodeHandshake(String hostname) {
        byte[] hostBytes = hostname.getBytes(StandardCharsets.UTF_8);
        byte[] packet = new byte[1 + hostBytes.length];
        packet[0] = MSG_TYPE_HANDSHAKE;
        System.arraycopy(hostBytes, 0, packet, 1, hostBytes.length);
        return packet;
    }

    public static byte[] encodePing() {
        return new byte[]{MSG_TYPE_PING};
    }

    public static byte[] encodeClipboard(String text) {
        byte[] textBytes = text.getBytes(StandardCharsets.UTF_8);
        byte[] packet = new byte[1 + textBytes.length];
        packet[0] = MSG_TYPE_CLIPBOARD;
        System.arraycopy(textBytes, 0, packet, 1, textBytes.length);
        return packet;
    }
}
