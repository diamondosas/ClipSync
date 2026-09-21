package com.diamond.clipsync.data.db;

import androidx.lifecycle.LiveData;
import androidx.room.Dao;
import androidx.room.Delete;
import androidx.room.Insert;
import androidx.room.OnConflictStrategy;
import androidx.room.Query;

import java.util.List;

@Dao
public interface ClipDao {

    @Query("SELECT * FROM clips ORDER BY is_pinned DESC, timestamp DESC")
    LiveData<List<ClipEntity>> getAllClipsLive();

    @Query("SELECT * FROM clips ORDER BY is_pinned DESC, timestamp DESC")
    List<ClipEntity> getAllClipsSync();

    @Query("SELECT * FROM clips WHERE content LIKE '%' || :query || '%' ORDER BY is_pinned DESC, timestamp DESC")
    LiveData<List<ClipEntity>> searchClips(String query);

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    long insert(ClipEntity clip);

    @Delete
    void delete(ClipEntity clip);

    @Query("DELETE FROM clips WHERE is_pinned = 0")
    void clearUnpinned();

    @Query("UPDATE clips SET is_pinned = :isPinned WHERE id = :id")
    void setPinned(long id, boolean isPinned);

    @Query("SELECT COUNT(*) FROM clips")
    int getCount();

    @Query("DELETE FROM clips WHERE is_pinned = 0 AND id NOT IN (SELECT id FROM clips ORDER BY is_pinned DESC, timestamp DESC LIMIT :limit)")
    void pruneOldUnpinned(int limit);
}
