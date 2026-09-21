package com.diamond.clipsync.ui;

import android.app.Application;

import androidx.annotation.NonNull;
import androidx.lifecycle.AndroidViewModel;
import androidx.lifecycle.LiveData;
import androidx.lifecycle.MutableLiveData;
import androidx.lifecycle.Transformations;

import com.diamond.clipsync.data.db.ClipEntity;
import com.diamond.clipsync.data.model.Device;
import com.diamond.clipsync.data.repository.ClipRepository;
import com.diamond.clipsync.network.PeerManager;
import com.diamond.clipsync.service.ServiceCoordinator;

import java.util.List;

public class MainViewModel extends AndroidViewModel {

    private final ClipRepository repository;
    private final ServiceCoordinator coordinator;
    private final MutableLiveData<String> searchQuery = new MutableLiveData<>("");
    private final MutableLiveData<List<Device>> devices = new MutableLiveData<>();
    private final LiveData<List<ClipEntity>> clips;

    public MainViewModel(@NonNull Application application) {
        super(application);
        this.coordinator = ServiceCoordinator.getInstance(application);
        this.repository = coordinator.getRepository();

        this.clips = Transformations.switchMap(searchQuery, query -> {
            if (query == null || query.trim().isEmpty()) {
                return repository.getAllClipsLive();
            } else {
                return repository.searchClips(query.trim());
            }
        });

        // Listen for peer changes
        this.coordinator.getPeerManager().addListener(devices::postValue);
    }

    public LiveData<List<ClipEntity>> getClips() {
        return clips;
    }

    public LiveData<List<Device>> getDevices() {
        return devices;
    }

    public void setSearchQuery(String query) {
        searchQuery.setValue(query);
    }

    public void togglePin(ClipEntity clip) {
        repository.togglePin(clip);
    }

    public void deleteClip(ClipEntity clip) {
        repository.deleteClip(clip);
    }

    public void clearAllUnpinned() {
        repository.clearAllUnpinned();
    }

    public void syncCurrentClip(String content) {
        coordinator.onLocalClipboardCopied(content);
    }

    public void connectManual(String ip) {
        coordinator.connectManual(ip);
    }
}
