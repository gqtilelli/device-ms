package handler

import (
	"encoding/json"
	"net/http"

	"github.com/device-ms/dto"
	"github.com/device-ms/errors"
	"github.com/device-ms/util"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type updateDeviceRequest struct {
	dto.UpdateDeviceRequestDTO
}

// Build builds the update request dto
func (req *updateDeviceRequest) Build(r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		return errors.DecodeError(err)
	}

	req.DeviceID, err = primitive.ObjectIDFromHex(mux.Vars(r)["id"])
	if err != nil {
		return errors.InvalidParameterError("id", "invalid object id ["+mux.Vars(r)["id"]+"]")
	}

	return req.Validate()
}

// Validate validates the update request dto
func (req updateDeviceRequest) Validate() error {
	if req.Brand == "" {
		return errors.RequiredParameterError("brand", "body")
	}
	if !req.Brand.IsValid() {
		return errors.InvalidParameterError("brand", "invalid value ["+string(req.Brand)+"]")
	}
	return nil
}

func (h deviceHandler) updateDevice(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := new(updateDeviceRequest)
	if err := req.Build(r); err != nil {
		util.JSONErrorWithCtx(ctx, w, err, http.StatusBadRequest)
		return
	}

	device, err := req.ToModel()
	if err != nil {
		util.JSONErrorWithCtx(ctx, w, err, http.StatusBadRequest)
		return
	}

	err = h.service.DeviceController().Update(ctx, device)
	if err != nil {
		util.JSONErrorWithCtx(ctx, w, err, http.StatusInternalServerError)
		return
	}

	util.JSONReturnWithCtx(ctx, w, http.StatusNoContent, nil)
}
