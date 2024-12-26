package handler

import (
	"net/http"

	"github.com/device-ms/errors"
	"github.com/device-ms/util"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type deleteDeviceParameters struct {
	deviceID primitive.ObjectID
}

func (params *deleteDeviceParameters) Build(r *http.Request) error {
	var err error
	params.deviceID, err = primitive.ObjectIDFromHex(mux.Vars(r)["id"])
	if err != nil {
		return errors.InvalidParameterError("id", "invalid object id ["+mux.Vars(r)["id"]+"]")
	}

	return nil
}

func (h deviceHandler) deleteDevice(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := new(getDeviceParameters)
	if err := params.Build(r); err != nil {
		util.JSONErrorWithCtx(ctx, w, err, http.StatusBadRequest)
		return
	}

	err := h.service.DeviceController().Delete(ctx, params.deviceID)
	if err != nil {
		util.JSONErrorWithCtx(ctx, w, err, http.StatusInternalServerError)
		return
	}

	util.JSONReturnWithCtx(ctx, w, http.StatusNoContent, nil)
}
