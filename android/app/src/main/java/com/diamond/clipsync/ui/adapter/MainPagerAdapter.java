package com.diamond.clipsync.ui.adapter;

import androidx.annotation.NonNull;
import androidx.fragment.app.Fragment;
import androidx.fragment.app.FragmentActivity;
import androidx.viewpager2.adapter.FragmentStateAdapter;

import com.diamond.clipsync.ui.fragment.ClipboardFragment;
import com.diamond.clipsync.ui.fragment.DevicesFragment;

public class MainPagerAdapter extends FragmentStateAdapter {

    public MainPagerAdapter(@NonNull FragmentActivity fragmentActivity) {
        super(fragmentActivity);
    }

    @NonNull
    @Override
    public Fragment createFragment(int position) {
        if (position == 0) {
            return new DevicesFragment();
        } else {
            return new ClipboardFragment();
        }
    }

    @Override
    public int getItemCount() {
        return 2;
    }
}
