package com.example.javaapp;

import android.app.Activity;
import android.content.Intent;
import android.os.Bundle;
import android.widget.Button;

public class MainActivity extends Activity {
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        
        Button btn = new Button(this);
        btn.setText("Open Gio screen");
        btn.setOnClickListener(v ->
            startActivity(new Intent(this, org.gioui.GioActivity.class))
        );
        
        setContentView(btn);
    }
}
