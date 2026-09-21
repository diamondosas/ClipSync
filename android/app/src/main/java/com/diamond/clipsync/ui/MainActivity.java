package com.diamond.clipsync.ui;

import android.Manifest;
import android.content.ClipData;
import android.content.ClipboardManager;
import android.content.Context;
import android.content.pm.PackageManager;
import android.os.Build;
import android.os.Bundle;
import android.text.Editable;
import android.text.TextWatcher;
import android.view.MenuItem;
import android.view.View;
import android.widget.EditText;
import android.widget.ImageView;
import android.widget.TextView;
import android.widget.Toast;

import androidx.appcompat.app.AppCompatActivity;
import androidx.core.app.ActivityCompat;
import androidx.core.content.ContextCompat;
import androidx.lifecycle.ViewModelProvider;
import androidx.recyclerview.widget.LinearLayoutManager;
import androidx.recyclerview.widget.RecyclerView;

import com.diamond.clipsync.R;
import com.diamond.clipsync.data.db.ClipEntity;
import com.diamond.clipsync.service.ClipSyncForegroundService;
import com.diamond.clipsync.ui.adapter.ClipAdapter;
import com.diamond.clipsync.ui.adapter.DeviceAdapter;
import com.diamond.clipsync.ui.dialog.ConnectDialog;
import com.diamond.clipsync.ui.dialog.HelpDialog;
import com.diamond.clipsync.ui.util.AccessibilityUtils;
import com.google.android.material.appbar.MaterialToolbar;
import com.google.android.material.bottomnavigation.BottomNavigationView;
import com.google.android.material.button.MaterialButton;

public class MainActivity extends AppCompatActivity {

    private MainViewModel viewModel;
    private DeviceAdapter deviceAdapter;
    private ClipAdapter clipAdapter;

    private View layoutDevicesTab;
    private View layoutClipboardTab;
    private View layoutDevicesEmpty;
    private View layoutClipsEmpty;
    private View bannerAccessibility;

    private TextView tvDevicesSearching;
    private TextView tvDevicesConnectedCount;
    private TextView tvClipsCount;
    private TextView tvClipsEmptyText;
    private EditText editSearch;
    private ImageView btnClearSearch;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);

        viewModel = new ViewModelProvider(this).get(MainViewModel.class);

        initViews();
        setupToolbar();
        setupBottomNav();
        setupDevicesTab();
        setupClipboardTab();
        observeViewModel();
        checkPermissions();

        // Start background foreground engine
        ClipSyncForegroundService.startService(this);
    }

    @Override
    protected void onResume() {
        super.onResume();
        updateAccessibilityBanner();
    }

    private void initViews() {
        layoutDevicesTab = findViewById(R.id.layout_devices_tab);
        layoutClipboardTab = findViewById(R.id.layout_clipboard_tab);
        layoutDevicesEmpty = findViewById(R.id.layout_devices_empty);
        layoutClipsEmpty = findViewById(R.id.layout_clips_empty);
        bannerAccessibility = findViewById(R.id.banner_accessibility);

        tvDevicesSearching = findViewById(R.id.tv_devices_searching);
        tvDevicesConnectedCount = findViewById(R.id.tv_devices_connected_count);
        tvClipsCount = findViewById(R.id.tv_clips_count);
        tvClipsEmptyText = findViewById(R.id.tv_clips_empty_text);
        editSearch = findViewById(R.id.edit_search);
        btnClearSearch = findViewById(R.id.btn_clear_search);

        findViewById(R.id.btn_enable_accessibility).setOnClickListener(v -> {
            AccessibilityUtils.openAccessibilitySettings(this);
        });

        findViewById(R.id.btn_empty_connect_manual).setOnClickListener(v -> showConnectDialog());
    }

    private void setupToolbar() {
        MaterialToolbar toolbar = findViewById(R.id.toolbar);
        toolbar.setOnMenuItemClickListener(item -> {
            if (item.getItemId() == R.id.action_connect) {
                showConnectDialog();
                return true;
            } else if (item.getItemId() == R.id.action_help) {
                new HelpDialog(this).show();
                return true;
            }
            return false;
        });
    }

    private void setupBottomNav() {
        BottomNavigationView bottomNav = findViewById(R.id.bottom_navigation);
        bottomNav.setOnItemSelectedListener(item -> {
            MaterialToolbar toolbar = findViewById(R.id.toolbar);
            if (item.getItemId() == R.id.nav_devices) {
                layoutDevicesTab.setVisibility(View.VISIBLE);
                layoutClipboardTab.setVisibility(View.GONE);
                toolbar.setTitle(R.string.title_devices);
                return true;
            } else if (item.getItemId() == R.id.nav_clipboard) {
                layoutDevicesTab.setVisibility(View.GONE);
                layoutClipboardTab.setVisibility(View.VISIBLE);
                toolbar.setTitle(R.string.title_clipboard);
                return true;
            }
            return false;
        });
    }

    private void setupDevicesTab() {
        RecyclerView recyclerDevices = findViewById(R.id.recycler_devices);
        recyclerDevices.setLayoutManager(new LinearLayoutManager(this));
        deviceAdapter = new DeviceAdapter();
        recyclerDevices.setAdapter(deviceAdapter);
    }

    private void setupClipboardTab() {
        RecyclerView recyclerClips = findViewById(R.id.recycler_clips);
        recyclerClips.setLayoutManager(new LinearLayoutManager(this));

        clipAdapter = new ClipAdapter(new ClipAdapter.ClipActionListener() {
            @Override
            public void onTogglePin(ClipEntity clip) {
                viewModel.togglePin(clip);
            }

            @Override
            public void onDeleteClip(ClipEntity clip) {
                viewModel.deleteClip(clip);
            }

            @Override
            public void onClipCopied(ClipEntity clip) {
                // Clipboard updated
            }
        });
        recyclerClips.setAdapter(clipAdapter);

        // Instant search text listener
        editSearch.addTextChangedListener(new TextWatcher() {
            @Override
            public void beforeTextChanged(CharSequence s, int start, int count, int after) {}

            @Override
            public void onTextChanged(CharSequence s, int start, int before, int count) {
                String query = s != null ? s.toString() : "";
                btnClearSearch.setVisibility(query.isEmpty() ? View.GONE : View.VISIBLE);
                viewModel.setSearchQuery(query);
            }

            @Override
            public void afterTextChanged(Editable s) {}
        });

        btnClearSearch.setOnClickListener(v -> editSearch.setText(""));

        // "Sync Phone Clip" button action
        MaterialButton btnSync = findViewById(R.id.btn_sync_clipboard);
        btnSync.setOnClickListener(v -> {
            try {
                ClipboardManager cm = (ClipboardManager) getSystemService(Context.CLIPBOARD_SERVICE);
                if (cm != null && cm.hasPrimaryClip() && cm.getPrimaryClip().getItemCount() > 0) {
                    CharSequence text = cm.getPrimaryClip().getItemAt(0).getText();
                    if (text != null && text.length() > 0) {
                        viewModel.syncCurrentClip(text.toString());
                        Toast.makeText(this, R.string.toast_synced, Toast.LENGTH_SHORT).show();
                        return;
                    }
                }
                Toast.makeText(this, "Clipboard is empty", Toast.LENGTH_SHORT).show();
            } catch (Exception e) {
                Toast.makeText(this, "Could not access clipboard: " + e.getMessage(), Toast.LENGTH_SHORT).show();
            }
        });

        // "Clear All" unpinned button action
        MaterialButton btnClearAll = findViewById(R.id.btn_clear_all_clips);
        btnClearAll.setOnClickListener(v -> viewModel.clearAllUnpinned());
    }

    private void observeViewModel() {
        // Observe Devices
        viewModel.getDevices().observe(this, devices -> {
            deviceAdapter.updateDevices(devices);
            int count = devices != null ? devices.size() : 0;
            tvDevicesConnectedCount.setText(getString(R.string.devices_connected_format, count));

            if (count == 0) {
                layoutDevicesEmpty.setVisibility(View.VISIBLE);
                tvDevicesSearching.setText(R.string.searching_devices);
            } else {
                layoutDevicesEmpty.setVisibility(View.GONE);
                tvDevicesSearching.setText("Local devices on Wi-Fi:");
            }
        });

        // Observe Clips
        viewModel.getClips().observe(this, clips -> {
            clipAdapter.updateClips(clips);
            int count = clips != null ? clips.size() : 0;
            String query = editSearch.getText() != null ? editSearch.getText().toString().trim() : "";

            if (query.isEmpty()) {
                tvClipsCount.setText(getString(R.string.clips_count_format, count));
            } else {
                tvClipsCount.setText(getString(R.string.clips_filter_format, count, count));
            }

            if (count == 0) {
                layoutClipsEmpty.setVisibility(View.VISIBLE);
                if (!query.isEmpty()) {
                    tvClipsEmptyText.setText(getString(R.string.no_clips_matched, query));
                } else {
                    tvClipsEmptyText.setText(R.string.no_clips_empty);
                }
            } else {
                layoutClipsEmpty.setVisibility(View.GONE);
            }
        });
    }

    private void updateAccessibilityBanner() {
        boolean enabled = AccessibilityUtils.isAccessibilityServiceEnabled(this);
        bannerAccessibility.setVisibility(enabled ? View.GONE : View.VISIBLE);
    }

    private void showConnectDialog() {
        new ConnectDialog(this, ip -> {
            viewModel.connectManual(ip);
        }).show();
    }

    private void checkPermissions() {
        // Request POST_NOTIFICATIONS on Android 13+
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.TIRAMISU) {
            if (ContextCompat.checkSelfPermission(this, Manifest.permission.POST_NOTIFICATIONS)
                    != PackageManager.PERMISSION_GRANTED) {
                ActivityCompat.requestPermissions(this,
                        new String[]{Manifest.permission.POST_NOTIFICATIONS}, 101);
            }
        }
    }
}
