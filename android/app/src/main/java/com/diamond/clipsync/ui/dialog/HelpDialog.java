package com.diamond.clipsync.ui.dialog;

import android.app.Dialog;
import android.content.Context;
import android.os.Bundle;
import android.view.Window;

import androidx.annotation.NonNull;

import com.diamond.clipsync.R;

public class HelpDialog extends Dialog {

    public HelpDialog(@NonNull Context context) {
        super(context);
    }

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        requestWindowFeature(Window.FEATURE_NO_TITLE);
        setContentView(R.layout.dialog_help);

        findViewById(R.id.btn_dialog_close).setOnClickListener(v -> dismiss());
    }
}
