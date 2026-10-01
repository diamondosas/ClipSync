package com.diamond.clipsync.service;

import android.content.BroadcastReceiver;
import android.content.Context;
import android.content.Intent;

/**
 * @deprecated Auto-start on boot is disabled per configuration.
 */
@Deprecated
public class BootReceiver extends BroadcastReceiver {
    @Override
    public void onReceive(Context context, Intent intent) {
        // Disabled: do not start on boot
    }
}
