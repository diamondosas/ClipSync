package com.diamond.clipsync.network;

import android.os.Handler;
import android.os.Looper;

import com.diamond.clipsync.data.model.Device;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

public class PeerManager {

    public interface PeerChangeListener {
        void onPeersChanged(List<Device> activePeers);
    }

    private final Map<String, Device> peers = new ConcurrentHashMap<>();
    private final List<PeerChangeListener> listeners = Collections.synchronizedList(new ArrayList<>());
    private final Handler mainHandler = new Handler(Looper.getMainLooper());

    public void addListener(PeerChangeListener listener) {
        if (listener != null && !listeners.contains(listener)) {
            listeners.add(listener);
            listener.onPeersChanged(getActivePeers());
        }
    }

    public void removeListener(PeerChangeListener listener) {
        listeners.remove(listener);
    }

    public void addOrUpdatePeer(String name, String ip, int port) {
        Device dev = peers.get(ip);
        if (dev == null) {
            dev = new Device(name, ip, port);
            peers.put(ip, dev);
        } else {
            if (name != null && !name.isEmpty()) {
                dev.setName(name);
            }
            dev.setPort(port);
            dev.touch();
        }
        notifyListeners();
    }

    public void touchPeer(String ip) {
        Device dev = peers.get(ip);
        if (dev != null) {
            dev.touch();
            notifyListeners();
        }
    }

    public void pruneDeadPeers(long timeoutMs) {
        long now = System.currentTimeMillis();
        boolean changed = false;
        for (Map.Entry<String, Device> entry : peers.entrySet()) {
            Device dev = entry.getValue();
            if (now - dev.getLastSeen() > timeoutMs) {
                if (dev.isAlive()) {
                    dev.setAlive(false);
                    changed = true;
                }
            }
        }
        if (changed) {
            notifyListeners();
        }
    }

    public List<Device> getActivePeers() {
        List<Device> active = new ArrayList<>();
        for (Device dev : peers.values()) {
            if (dev.isAlive()) {
                active.add(dev);
            }
        }
        return active;
    }

    public List<Device> getAllPeers() {
        return new ArrayList<>(peers.values());
    }

    private void notifyListeners() {
        final List<Device> copy = getActivePeers();
        mainHandler.post(() -> {
            synchronized (listeners) {
                for (PeerChangeListener l : listeners) {
                    l.onPeersChanged(copy);
                }
            }
        });
    }
}
