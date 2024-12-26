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

type updateDeviceBrandRequest struct {
	dto.UpdateDeviceBrandRequestDTO
}

// Build builds the update device brand request dto
func (req *updateDeviceBrandRequest) Build(r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		return errors.DecodeError(err)
	}

	req.DeviceID, err = primitive.ObjectIDFromHex(mux.Vars(r)["id"])
	if err != nil {
		return errors.InvalidParameterError("id", "invalid object id ["+mux.Vars(r)["id"]+"]")
	}

	return nil
}

func (h deviceHandler) updateDeviceBrand(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := new(updateDeviceBrandRequest)
	if err := req.Build(r); err != nil {
		util.JSONErrorWithCtx(ctx, w, err, http.StatusBadRequest)
		return
	}

	device, err := req.ToModel()
	if err != nil {
		util.JSONErrorWithCtx(ctx, w, err, http.StatusBadRequest)
		return
	}

	err = h.service.DeviceController().UpdateBrand(ctx, device.ID, device.Brand)
	if err != nil {
		util.JSONErrorWithCtx(ctx, w, err, http.StatusInternalServerError)
		return
	}

	util.JSONReturnWithCtx(ctx, w, http.StatusNoContent, nil)
}
