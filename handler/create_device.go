package handler

import (
	"encoding/json"
	"net/http"

	"github.com/device-ms/dto"
	"github.com/device-ms/errors"
	"github.com/device-ms/util"
)

type createDeviceRequest struct {
	dto.CreateDeviceRequestDTO
}

// Build builds the creation dto
func (req *createDeviceRequest) Build(r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return errors.DecodeError(err)
	}

	return req.Validate()
}

// Validate validates the creation dto
func (req createDeviceRequest) Validate() error {
	if req.Brand == "" {
		return errors.RequiredParameterError("brand", "body")
	}
	if !req.Brand.IsValid() {
		return errors.InvalidParameterError("brand", "invalid value ["+string(req.Brand)+"]")
	}
	return nil
}

func (h deviceHandler) createDevice(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := new(createDeviceRequest)
	if err := req.Build(r); err != nil {
		util.JSONErrorWithCtx(ctx, w, err, http.StatusBadRequest)
		return
	}

	device := req.ToModel()

	err := h.service.DeviceController().Create(ctx, device)
	if err != nil {
		util.JSONErrorWithCtx(ctx, w, err, http.StatusInternalServerError)
		return
	}

	util.JSONReturnWithCtx(ctx, w, http.StatusCreated, dto.ToCreatedDeviceResponseDTO(device))
}
