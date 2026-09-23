package com.diamond.clipsync.ui.fragment;

import android.content.ClipData;
import android.content.ClipboardManager;
import android.content.Context;
import android.os.Bundle;
import android.text.Editable;
import android.text.TextWatcher;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.EditText;
import android.widget.ImageView;
import android.widget.TextView;
import android.widget.Toast;

import androidx.annotation.NonNull;
import androidx.annotation.Nullable;
import androidx.fragment.app.Fragment;
import androidx.lifecycle.ViewModelProvider;
import androidx.recyclerview.widget.LinearLayoutManager;
import androidx.recyclerview.widget.RecyclerView;

import com.diamond.clipsync.R;
import com.diamond.clipsync.data.db.ClipEntity;
import com.diamond.clipsync.ui.MainViewModel;
import com.diamond.clipsync.ui.adapter.ClipAdapter;
import com.google.android.material.button.MaterialButton;

public class ClipboardFragment extends Fragment {

    private MainViewModel viewModel;
    private ClipAdapter clipAdapter;

    private TextView tvClipsCount;
    private TextView tvClipsEmptyText;
    private View layoutClipsEmpty;
    private EditText editSearch;
    private ImageView btnClearSearch;

    @Nullable
    @Override
    public View onCreateView(@NonNull LayoutInflater inflater, @Nullable ViewGroup container, @Nullable Bundle savedInstanceState) {
        return inflater.inflate(R.layout.fragment_clipboard, container, false);
    }

    @Override
    public void onViewCreated(@NonNull View view, @Nullable Bundle savedInstanceState) {
        super.onViewCreated(view, savedInstanceState);

        viewModel = new ViewModelProvider(requireActivity()).get(MainViewModel.class);

        tvClipsCount = view.findViewById(R.id.tv_clips_count);
        tvClipsEmptyText = view.findViewById(R.id.tv_clips_empty_text);
        layoutClipsEmpty = view.findViewById(R.id.layout_clips_empty);
        editSearch = view.findViewById(R.id.edit_search);
        btnClearSearch = view.findViewById(R.id.btn_clear_search);

        RecyclerView recyclerClips = view.findViewById(R.id.recycler_clips);
        recyclerClips.setLayoutManager(new LinearLayoutManager(requireContext()));

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
        MaterialButton btnSync = view.findViewById(R.id.btn_sync_clipboard);
        btnSync.setOnClickListener(v -> {
            try {
                ClipboardManager cm = (ClipboardManager) requireContext().getSystemService(Context.CLIPBOARD_SERVICE);
                if (cm != null && cm.hasPrimaryClip() && cm.getPrimaryClip().getItemCount() > 0) {
                    CharSequence text = cm.getPrimaryClip().getItemAt(0).getText();
                    if (text != null && text.length() > 0) {
                        viewModel.syncCurrentClip(text.toString());
                        Toast.makeText(requireContext(), R.string.toast_synced, Toast.LENGTH_SHORT).show();
                        return;
                    }
                }
                Toast.makeText(requireContext(), "Clipboard is empty", Toast.LENGTH_SHORT).show();
            } catch (Exception e) {
                Toast.makeText(requireContext(), "Could not access clipboard: " + e.getMessage(), Toast.LENGTH_SHORT).show();
            }
        });

        // "Clear All" unpinned button action
        MaterialButton btnClearAll = view.findViewById(R.id.btn_clear_all_clips);
        btnClearAll.setOnClickListener(v -> viewModel.clearAllUnpinned());

        observeViewModel();
    }

    private void observeViewModel() {
        viewModel.getClips().observe(getViewLifecycleOwner(), clips -> {
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
}
