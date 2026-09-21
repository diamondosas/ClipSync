package com.diamond.clipsync.ui.util;

import android.accessibilityservice.AccessibilityServiceInfo;
import android.content.Context;
import android.content.Intent;
import android.provider.Settings;
import android.view.accessibility.AccessibilityManager;

import com.diamond.clipsync.service.ClipSyncAccessibilityService;

import java.util.List;

public class AccessibilityUtils {

    public static boolean isAccessibilityServiceEnabled(Context context) {
        if (ClipSyncAccessibilityService.isRunning()) {
            return true;
        }

        AccessibilityManager am = (AccessibilityManager) context.getSystemService(Context.ACCESSIBILITY_SERVICE);
        if (am == null) return false;

        List<AccessibilityServiceInfo> enabledServices =
                am.getEnabledAccessibilityServiceList(AccessibilityServiceInfo.FEEDBACK_ALL_MASK);

        String targetPackage = context.getPackageName();
        for (AccessibilityServiceInfo info : enabledServices) {
            if (info.getId() != null && info.getId().contains(targetPackage)) {
                return true;
            }
        }
        return false;
    }

    public static void openAccessibilitySettings(Context context) {
        Intent intent = new Intent(Settings.ACTION_ACCESSIBILITY_SETTINGS);
        intent.addFlags(Intent.FLAG_ACTIVITY_NEW_TASK);
        context.startActivity(intent);
    }
}
