package com.diamond.clipsync.ui.adapter;

import android.content.ClipData;
import android.content.ClipboardManager;
import android.content.Context;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.TextView;
import android.widget.Toast;

import androidx.annotation.NonNull;
import androidx.cardview.widget.CardView;
import androidx.core.content.ContextCompat;
import androidx.recyclerview.widget.RecyclerView;

import com.diamond.clipsync.R;
import com.diamond.clipsync.data.db.ClipEntity;
import com.diamond.clipsync.ui.util.TextFormatter;
import com.google.android.material.button.MaterialButton;

import java.util.ArrayList;
import java.util.HashSet;
import java.util.List;
import java.util.Set;

public class ClipAdapter extends RecyclerView.Adapter<ClipAdapter.ClipViewHolder> {

    public interface ClipActionListener {
        void onTogglePin(ClipEntity clip);
        void onDeleteClip(ClipEntity clip);
        void onClipCopied(ClipEntity clip);
    }

    private final List<ClipEntity> clips = new ArrayList<>();
    private final Set<Long> expandedIds = new HashSet<>();
    private final ClipActionListener listener;

    public ClipAdapter(ClipActionListener listener) {
        this.listener = listener;
    }

    public void updateClips(List<ClipEntity> newClips) {
        clips.clear();
        if (newClips != null) {
            clips.addAll(newClips);
        }
        notifyDataSetChanged();
    }

    @NonNull
    @Override
    public ClipViewHolder onCreateViewHolder(@NonNull ViewGroup parent, int viewType) {
        View view = LayoutInflater.from(parent.getContext())
                .inflate(R.layout.item_clip_card, parent, false);
        return new ClipViewHolder(view);
    }

    @Override
    public void onBindViewHolder(@NonNull ClipViewHolder holder, int position) {
        ClipEntity clip = clips.get(position);
        Context context = holder.itemView.getContext();

        String cleanedText = TextFormatter.cleanDisplayText(clip.content);
        holder.tvContent.setText(cleanedText);

        // Pinned state styling
        if (clip.isPinned) {
            holder.cardClip.setCardBackgroundColor(ContextCompat.getColor(context, R.color.surface_pinned));
            holder.tvPinnedBadge.setVisibility(View.VISIBLE);
            holder.btnPin.setText(R.string.pinned_action);
            holder.btnPin.setTextColor(ContextCompat.getColor(context, R.color.text_primary));
            holder.btnPin.setBackgroundColor(ContextCompat.getColor(context, R.color.accent_brown));
        } else {
            holder.cardClip.setCardBackgroundColor(ContextCompat.getColor(context, R.color.surface_card));
            holder.tvPinnedBadge.setVisibility(View.GONE);
            holder.btnPin.setText(R.string.pin_action);
            holder.btnPin.setTextColor(ContextCompat.getColor(context, R.color.text_muted));
            holder.btnPin.setBackgroundColor(ContextCompat.getColor(context, R.color.bg_dark));
        }

        // Source device caption
        if (clip.sourceDevice != null && !clip.sourceDevice.isEmpty()) {
            holder.tvSource.setText("From: " + clip.sourceDevice);
            holder.tvSource.setVisibility(View.VISIBLE);
        } else {
            holder.tvSource.setVisibility(View.GONE);
        }

        // Expand / Collapse for long text (> 80 chars or >= 3 lines)
        boolean isLong = cleanedText.length() > 80 || cleanedText.contains("\n");
        boolean isExpanded = expandedIds.contains(clip.id);

        if (isLong) {
            holder.btnShowMore.setVisibility(View.VISIBLE);
            if (isExpanded) {
                holder.tvContent.setMaxLines(Integer.MAX_VALUE);
                holder.btnShowMore.setText(R.string.show_less);
            } else {
                holder.tvContent.setMaxLines(3);
                holder.btnShowMore.setText(R.string.show_more);
            }
            holder.btnShowMore.setOnClickListener(v -> {
                if (expandedIds.contains(clip.id)) {
                    expandedIds.remove(clip.id);
                } else {
                    expandedIds.add(clip.id);
                }
                notifyItemChanged(holder.getBindingAdapterPosition());
            });
        } else {
            holder.btnShowMore.setVisibility(View.GONE);
            holder.tvContent.setMaxLines(3);
        }

        // Tap-to-copy card action
        holder.itemView.setOnClickListener(v -> {
            try {
                ClipboardManager cm = (ClipboardManager) context.getSystemService(Context.CLIPBOARD_SERVICE);
                if (cm != null) {
                    cm.setPrimaryClip(ClipData.newPlainText("ClipSync", clip.content));
                    Toast.makeText(context, R.string.toast_copied, Toast.LENGTH_SHORT).show();
                }
                if (listener != null) {
                    listener.onClipCopied(clip);
                }
            } catch (Exception ignored) {}
        });

        // Pin button action
        holder.btnPin.setOnClickListener(v -> {
            if (listener != null) {
                listener.onTogglePin(clip);
            }
        });

        // Delete button action
        holder.btnDelete.setOnClickListener(v -> {
            if (listener != null) {
                listener.onDeleteClip(clip);
            }
        });
    }

    @Override
    public int getItemCount() {
        return clips.size();
    }

    static class ClipViewHolder extends RecyclerView.ViewHolder {
        final CardView cardClip;
        final TextView tvContent;
        final TextView btnShowMore;
        final TextView tvPinnedBadge;
        final TextView tvSource;
        final MaterialButton btnPin;
        final MaterialButton btnDelete;

        public ClipViewHolder(@NonNull View itemView) {
            super(itemView);
            cardClip = itemView.findViewById(R.id.card_clip);
            tvContent = itemView.findViewById(R.id.tv_clip_content);
            btnShowMore = itemView.findViewById(R.id.btn_show_more);
            tvPinnedBadge = itemView.findViewById(R.id.tv_pinned_badge);
            tvSource = itemView.findViewById(R.id.tv_clip_source);
            btnPin = itemView.findViewById(R.id.btn_pin);
            btnDelete = itemView.findViewById(R.id.btn_delete);
        }
    }
}
