package com.diamond.clipsync.data.db;

import androidx.room.ColumnInfo;
import androidx.room.Entity;
import androidx.room.PrimaryKey;

@Entity(tableName = "clips")
public class ClipEntity {
    @PrimaryKey(autoGenerate = true)
    public long id;

    @ColumnInfo(name = "content")
    public String content;

    @ColumnInfo(name = "timestamp")
    public long timestamp;

    @ColumnInfo(name = "is_pinned")
    public boolean isPinned;

    @ColumnInfo(name = "source_device")
    public String sourceDevice;

    @ColumnInfo(name = "source_ip")
    public String sourceIp;

    public ClipEntity(String content, long timestamp, boolean isPinned, String sourceDevice, String sourceIp) {
        this.content = content;
        this.timestamp = timestamp;
        this.isPinned = isPinned;
        this.sourceDevice = sourceDevice;
        this.sourceIp = sourceIp;
    }
}
