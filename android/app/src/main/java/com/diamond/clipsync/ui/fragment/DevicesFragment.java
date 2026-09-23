package com.diamond.clipsync.ui.fragment;

import android.os.Bundle;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.TextView;

import androidx.annotation.NonNull;
import androidx.annotation.Nullable;
import androidx.fragment.app.Fragment;
import androidx.lifecycle.ViewModelProvider;
import androidx.recyclerview.widget.LinearLayoutManager;
import androidx.recyclerview.widget.RecyclerView;

import com.diamond.clipsync.R;
import com.diamond.clipsync.ui.MainViewModel;
import com.diamond.clipsync.ui.adapter.DeviceAdapter;
import com.diamond.clipsync.ui.dialog.ConnectDialog;

public class DevicesFragment extends Fragment {

    private MainViewModel viewModel;
    private DeviceAdapter deviceAdapter;

    private TextView tvDevicesSearching;
    private TextView tvDevicesConnectedCount;
    private View layoutDevicesEmpty;

    @Nullable
    @Override
    public View onCreateView(@NonNull LayoutInflater inflater, @Nullable ViewGroup container, @Nullable Bundle savedInstanceState) {
        return inflater.inflate(R.layout.fragment_devices, container, false);
    }

    @Override
    public void onViewCreated(@NonNull View view, @Nullable Bundle savedInstanceState) {
        super.onViewCreated(view, savedInstanceState);

        viewModel = new ViewModelProvider(requireActivity()).get(MainViewModel.class);

        tvDevicesSearching = view.findViewById(R.id.tv_devices_searching);
        tvDevicesConnectedCount = view.findViewById(R.id.tv_devices_connected_count);
        layoutDevicesEmpty = view.findViewById(R.id.layout_devices_empty);

        RecyclerView recyclerDevices = view.findViewById(R.id.recycler_devices);
        recyclerDevices.setLayoutManager(new LinearLayoutManager(requireContext()));
        deviceAdapter = new DeviceAdapter();
        recyclerDevices.setAdapter(deviceAdapter);

        view.findViewById(R.id.btn_empty_connect_manual).setOnClickListener(v -> {
            new ConnectDialog(requireContext(), ip -> viewModel.connectManual(ip)).show();
        });

        observeViewModel();
    }

    private void observeViewModel() {
        viewModel.getDevices().observe(getViewLifecycleOwner(), devices -> {
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
    }
}
