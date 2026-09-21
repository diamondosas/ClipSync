package com.diamond.clipsync.network;

import android.content.Context;
import android.net.nsd.NsdManager;
import android.net.nsd.NsdServiceInfo;
import android.os.Build;
import android.util.Log;

import java.net.InetAddress;

public class NsdHelper {

    private static final String TAG = "ClipSyncNsd";
    private final Context context;
    private final NsdManager nsdManager;
    private final PeerManager peerManager;
    private final String serviceName;
    private final int port;

    private NsdManager.RegistrationListener registrationListener;
    private NsdManager.DiscoveryListener discoveryListener;
    private boolean isRegistered = false;
    private boolean isDiscovering = false;

    public interface OnPeerDiscoveredCallback {
        void onPeerFound(String name, String ip, int port);
    }

    private OnPeerDiscoveredCallback peerCallback;

    public NsdHelper(Context context, PeerManager peerManager, String serviceName, int port) {
        this.context = context.getApplicationContext();
        this.peerManager = peerManager;
        this.serviceName = serviceName;
        this.port = port;
        this.nsdManager = (NsdManager) this.context.getSystemService(Context.NSD_SERVICE);
    }

    public void setPeerCallback(OnPeerDiscoveredCallback callback) {
        this.peerCallback = callback;
    }

    public synchronized void start() {
        registerService();
        discoverServices();
    }

    public synchronized void stop() {
        stopDiscovery();
        unregisterService();
    }

    private void registerService() {
        if (isRegistered || nsdManager == null) return;

        NsdServiceInfo serviceInfo = new NsdServiceInfo();
        serviceInfo.setServiceName(serviceName);
        serviceInfo.setServiceType(Protocol.SERVICE_TYPE);
        serviceInfo.setPort(port);

        registrationListener = new NsdManager.RegistrationListener() {
            @Override
            public void onServiceRegistered(NsdServiceInfo info) {
                isRegistered = true;
                Log.i(TAG, "NSD Service registered: " + info.getServiceName());
            }

            @Override
            public void onRegistrationFailed(NsdServiceInfo info, int errorCode) {
                Log.e(TAG, "NSD Registration failed: " + errorCode);
            }

            @Override
            public void onServiceUnregistered(NsdServiceInfo info) {
                isRegistered = false;
                Log.i(TAG, "NSD Service unregistered");
            }

            @Override
            public void onUnregistrationFailed(NsdServiceInfo info, int errorCode) {
                Log.e(TAG, "NSD Unregistration failed: " + errorCode);
            }
        };

        try {
            nsdManager.registerService(serviceInfo, NsdManager.PROTOCOL_DNS_SD, registrationListener);
        } catch (Exception e) {
            Log.e(TAG, "Error registering NSD service: " + e.getMessage());
        }
    }

    private void discoverServices() {
        if (isDiscovering || nsdManager == null) return;

        discoveryListener = new NsdManager.DiscoveryListener() {
            @Override
            public void onDiscoveryStarted(String serviceType) {
                isDiscovering = true;
                Log.i(TAG, "NSD Discovery started for " + serviceType);
            }

            @Override
            public void onServiceFound(NsdServiceInfo serviceInfo) {
                Log.d(TAG, "NSD Service found: " + serviceInfo.getServiceName());
                if (serviceInfo.getServiceName().equals(serviceName)) {
                    // Ignore self
                    return;
                }
                resolveService(serviceInfo);
            }

            @Override
            public void onServiceLost(NsdServiceInfo serviceInfo) {
                Log.d(TAG, "NSD Service lost: " + serviceInfo.getServiceName());
            }

            @Override
            public void onDiscoveryStopped(String serviceType) {
                isDiscovering = false;
                Log.i(TAG, "NSD Discovery stopped");
            }

            @Override
            public void onStartDiscoveryFailed(String serviceType, int errorCode) {
                Log.e(TAG, "Start discovery failed: " + errorCode);
                stopDiscovery();
            }

            @Override
            public void onStopDiscoveryFailed(String serviceType, int errorCode) {
                Log.e(TAG, "Stop discovery failed: " + errorCode);
            }
        };

        try {
            nsdManager.discoverServices(Protocol.SERVICE_TYPE, NsdManager.PROTOCOL_DNS_SD, discoveryListener);
        } catch (Exception e) {
            Log.e(TAG, "Error starting NSD discovery: " + e.getMessage());
        }
    }

    private void resolveService(NsdServiceInfo serviceInfo) {
        if (nsdManager == null) return;

        try {
            nsdManager.resolveService(serviceInfo, new NsdManager.ResolveListener() {
                @Override
                public void onResolveFailed(NsdServiceInfo info, int errorCode) {
                    Log.w(TAG, "Resolve failed for " + info.getServiceName() + ": " + errorCode);
                }

                @Override
                public void onServiceResolved(NsdServiceInfo info) {
                    InetAddress host = info.getHost();
                    if (host != null) {
                        String ip = host.getHostAddress();
                        int peerPort = info.getPort();
                        String peerName = info.getServiceName();
                        Log.i(TAG, "Resolved peer: " + peerName + " @ " + ip + ":" + peerPort);
                        if (peerManager != null) {
                            peerManager.addOrUpdatePeer(peerName, ip, peerPort);
                        }
                        if (peerCallback != null) {
                            peerCallback.onPeerFound(peerName, ip, peerPort);
                        }
                    }
                }
            });
        } catch (Exception e) {
            Log.e(TAG, "Error resolving service: " + e.getMessage());
        }
    }

    private void stopDiscovery() {
        if (isDiscovering && nsdManager != null && discoveryListener != null) {
            try {
                nsdManager.stopServiceDiscovery(discoveryListener);
            } catch (Exception e) {
                Log.e(TAG, "Error stopping discovery: " + e.getMessage());
            }
            isDiscovering = false;
        }
    }

    private void unregisterService() {
        if (isRegistered && nsdManager != null && registrationListener != null) {
            try {
                nsdManager.unregisterService(registrationListener);
            } catch (Exception e) {
                Log.e(TAG, "Error unregistering service: " + e.getMessage());
            }
            isRegistered = false;
        }
    }
}
