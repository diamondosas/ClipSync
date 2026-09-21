package com.diamond.clipsync.ui.dialog;

import android.app.Dialog;
import android.content.Context;
import android.os.Bundle;
import android.view.Window;
import android.widget.EditText;
import android.widget.Toast;

import androidx.annotation.NonNull;

import com.diamond.clipsync.R;

public class ConnectDialog extends Dialog {

    public interface OnConnectListener {
        void onConnect(String ip);
    }

    private final OnConnectListener listener;

    public ConnectDialog(@NonNull Context context, OnConnectListener listener) {
        super(context);
        this.listener = listener;
    }

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        requestWindowFeature(Window.FEATURE_NO_TITLE);
        setContentView(R.layout.dialog_connect);

        EditText editIp = findViewById(R.id.edit_connect_ip);

        findViewById(R.id.btn_dialog_cancel).setOnClickListener(v -> dismiss());

        findViewById(R.id.btn_dialog_connect).setOnClickListener(v -> {
            String ip = editIp.getText() != null ? editIp.getText().toString().trim() : "";
            if (ip.isEmpty()) {
                Toast.makeText(getContext(), "Please enter a valid IP address", Toast.LENGTH_SHORT).show();
                return;
            }
            if (listener != null) {
                listener.onConnect(ip);
            }
            dismiss();
        });
    }
}
