package handler

import (
	"net/http"

	"github.com/device-ms/dto"
	"github.com/device-ms/errors"
	"github.com/device-ms/util"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type getDeviceParameters struct {
	deviceID primitive.ObjectID
}

func (params *getDeviceParameters) Build(r *http.Request) error {
	var err error
	params.deviceID, err = primitive.ObjectIDFromHex(mux.Vars(r)["id"])
	if err != nil {
		return errors.InvalidParameterError("id", "invalid object id ["+mux.Vars(r)["id"]+"]")
	}

	return nil
}

func (h deviceHandler) getDevice(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := new(getDeviceParameters)
	if err := params.Build(r); err != nil {
		util.JSONErrorWithCtx(ctx, w, err, http.StatusBadRequest)
		return
	}

	device, err := h.service.DeviceController().GetDevice(ctx, params.deviceID)
	if err != nil {
		util.JSONErrorWithCtx(ctx, w, err, http.StatusInternalServerError)
		return
	}

	util.JSONReturnWithCtx(ctx, w, http.StatusOK, dto.ToDeviceDTO(device))
}
