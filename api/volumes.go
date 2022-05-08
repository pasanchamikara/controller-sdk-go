package api

// Mount is the definition of PATCH /v2/apps/<app_id>/volumes/<name>/path/.
type Mount struct {
	Values map[string]string `json:"values"`
}

// Unmount is the definition of PATCH /v2/apps/<app_id>/volumes/<name>/path/.
type Unmount struct {
	Values map[string]interface{} `json:"values"`
}

// Volume is the structure of an app's volume.
type Volume struct {
	// Owner is the app owner.
	Owner string `json:"owner,omitempty"`
	// App is the app the tls settings apply to and cannot be updated.
	App string `json:"app,omitempty"`
	// Created is the time that the volume was created and cannot be updated.
	Created string `json:"created,omitempty"`
	// Updated is the last time the TLS settings was changed and cannot be updated.
	Updated string `json:"updated,omitempty"`
	// UUID is a unique string reflecting the volume in its current state.
	// It changes every time the volume is changed and cannot be updated.
	UUID string `json:"uuid,omitempty"`
	// Volume's name
	Name string `json:"name,omitempty"`
	//Volume's size
	Size string `json:"size,omitempty"`
	// mount application's path
	Path map[string]interface{} `json:"path,omitempty"`
}

type Volumes []Volume
