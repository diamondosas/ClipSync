package com.diamond.clipsync.data.repository;

import android.content.Context;

import androidx.lifecycle.LiveData;

import com.diamond.clipsync.data.db.AppDatabase;
import com.diamond.clipsync.data.db.ClipDao;
import com.diamond.clipsync.data.db.ClipEntity;

import java.util.List;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

public class ClipRepository {

    private static volatile ClipRepository INSTANCE;
    private final ClipDao clipDao;
    private final ExecutorService executor;
    private static final int MAX_HISTORY_LIMIT = 100;

    private ClipRepository(Context context) {
        AppDatabase db = AppDatabase.getInstance(context);
        this.clipDao = db.clipDao();
        this.executor = Executors.newSingleThreadExecutor();
    }

    public static ClipRepository getInstance(Context context) {
        if (INSTANCE == null) {
            synchronized (ClipRepository.class) {
                if (INSTANCE == null) {
                    INSTANCE = new ClipRepository(context.getApplicationContext());
                }
            }
        }
        return INSTANCE;
    }

    public LiveData<List<ClipEntity>> getAllClipsLive() {
        return clipDao.getAllClipsLive();
    }

    public LiveData<List<ClipEntity>> searchClips(String query) {
        return clipDao.searchClips(query);
    }

    public void insertClip(final String content, final String sourceDevice, final String sourceIp) {
        if (content == null || content.trim().isEmpty()) {
            return;
        }
        executor.execute(() -> {
            ClipEntity entity = new ClipEntity(
                    content,
                    System.currentTimeMillis(),
                    false,
                    sourceDevice,
                    sourceIp
            );
            clipDao.insert(entity);
            clipDao.pruneOldUnpinned(MAX_HISTORY_LIMIT);
        });
    }

    public void togglePin(final ClipEntity clip) {
        executor.execute(() -> {
            clipDao.setPinned(clip.id, !clip.isPinned);
        });
    }

    public void deleteClip(final ClipEntity clip) {
        executor.execute(() -> {
            clipDao.delete(clip);
        });
    }

    public void clearAllUnpinned() {
        executor.execute(() -> {
            clipDao.clearUnpinned();
        });
    }
}
