package com.diamond.clipsync.ui;

import android.Manifest;
import android.content.pm.PackageManager;
import android.os.Build;
import android.os.Bundle;
import android.view.View;

import androidx.appcompat.app.AppCompatActivity;
import androidx.core.app.ActivityCompat;
import androidx.core.content.ContextCompat;
import androidx.lifecycle.ViewModelProvider;
import androidx.viewpager2.widget.ViewPager2;

import com.diamond.clipsync.R;
import com.diamond.clipsync.service.ClipSyncForegroundService;
import com.diamond.clipsync.ui.adapter.MainPagerAdapter;
import com.diamond.clipsync.ui.dialog.ConnectDialog;
import com.diamond.clipsync.ui.dialog.HelpDialog;
import com.diamond.clipsync.ui.util.AccessibilityUtils;
import com.google.android.material.appbar.MaterialToolbar;
import com.google.android.material.bottomnavigation.BottomNavigationView;

public class MainActivity extends AppCompatActivity {

    private MainViewModel viewModel;
    private ViewPager2 viewPager;
    private BottomNavigationView bottomNav;
    private MaterialToolbar toolbar;
    private View bannerAccessibility;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);

        viewModel = new ViewModelProvider(this).get(MainViewModel.class);

        initViews();
        setupToolbar();
        setupViewPager();
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
        toolbar = findViewById(R.id.toolbar);
        viewPager = findViewById(R.id.view_pager);
        bottomNav = findViewById(R.id.bottom_navigation);
        bannerAccessibility = findViewById(R.id.banner_accessibility);

        findViewById(R.id.btn_enable_accessibility).setOnClickListener(v -> {
            AccessibilityUtils.openAccessibilitySettings(this);
        });
    }

    private void setupToolbar() {
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

    private void setupViewPager() {
        MainPagerAdapter pagerAdapter = new MainPagerAdapter(this);
        viewPager.setAdapter(pagerAdapter);
        viewPager.setOffscreenPageLimit(1);

        // Sync swipe gesture with bottom navigation and toolbar title
        viewPager.registerOnPageChangeCallback(new ViewPager2.OnPageChangeCallback() {
            @Override
            public void onPageSelected(int position) {
                super.onPageSelected(position);
                int navId = (position == 0) ? R.id.nav_devices : R.id.nav_clipboard;
                if (bottomNav.getSelectedItemId() != navId) {
                    bottomNav.setSelectedItemId(navId);
                }
                toolbar.setTitle(position == 0 ? R.string.title_devices : R.string.title_clipboard);
            }
        });

        // Sync bottom navigation clicks with ViewPager2
        bottomNav.setOnItemSelectedListener(item -> {
            int targetItem = (item.getItemId() == R.id.nav_devices) ? 0 : 1;
            if (viewPager.getCurrentItem() != targetItem) {
                viewPager.setCurrentItem(targetItem, true);
            }
            return true;
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
