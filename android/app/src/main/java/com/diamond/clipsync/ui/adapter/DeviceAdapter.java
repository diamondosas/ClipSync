package com.diamond.clipsync.ui.adapter;

import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.TextView;

import androidx.annotation.NonNull;
import androidx.recyclerview.widget.RecyclerView;

import com.diamond.clipsync.R;
import com.diamond.clipsync.data.model.Device;

import java.util.ArrayList;
import java.util.List;

public class DeviceAdapter extends RecyclerView.Adapter<DeviceAdapter.DeviceViewHolder> {

    private final List<Device> devices = new ArrayList<>();

    public void updateDevices(List<Device> newDevices) {
        devices.clear();
        if (newDevices != null) {
            devices.addAll(newDevices);
        }
        notifyDataSetChanged();
    }

    @NonNull
    @Override
    public DeviceViewHolder onCreateViewHolder(@NonNull ViewGroup parent, int viewType) {
        View view = LayoutInflater.from(parent.getContext())
                .inflate(R.layout.item_device_card, parent, false);
        return new DeviceViewHolder(view);
    }

    @Override
    public void onBindViewHolder(@NonNull DeviceViewHolder holder, int position) {
        Device device = devices.get(position);
        holder.tvName.setText(device.getName());
        holder.tvIp.setText("IP: " + device.getIp());
        holder.tvStatus.setText("● Online");
    }

    @Override
    public int getItemCount() {
        return devices.size();
    }

    static class DeviceViewHolder extends RecyclerView.ViewHolder {
        final TextView tvName;
        final TextView tvIp;
        final TextView tvStatus;

        public DeviceViewHolder(@NonNull View itemView) {
            super(itemView);
            tvName = itemView.findViewById(R.id.tv_device_name);
            tvIp = itemView.findViewById(R.id.tv_device_ip);
            tvStatus = itemView.findViewById(R.id.tv_device_status);
        }
    }
}
