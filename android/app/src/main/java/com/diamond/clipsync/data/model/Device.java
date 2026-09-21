package com.diamond.clipsync.data.model;

import java.util.Objects;

public class Device {
    private String name;
    private String ip;
    private int port;
    private boolean isAlive;
    private long lastSeen;

    public Device(String name, String ip, int port) {
        this.name = (name != null && !name.isEmpty()) ? name : ip;
        this.ip = ip;
        this.port = port;
        this.isAlive = true;
        this.lastSeen = System.currentTimeMillis();
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getIp() {
        return ip;
    }

    public void setIp(String ip) {
        this.ip = ip;
    }

    public int getPort() {
        return port;
    }

    public void setPort(int port) {
        this.port = port;
    }

    public boolean isAlive() {
        return isAlive;
    }

    public void setAlive(boolean alive) {
        isAlive = alive;
    }

    public long getLastSeen() {
        return lastSeen;
    }

    public void setLastSeen(long lastSeen) {
        this.lastSeen = lastSeen;
    }

    public void touch() {
        this.lastSeen = System.currentTimeMillis();
        this.isAlive = true;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        Device device = (Device) o;
        return Objects.equals(ip, device.ip);
    }

    @Override
    public int hashCode() {
        return Objects.hash(ip);
    }
}
