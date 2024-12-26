package handler

import (
	"net/http"

	"github.com/device-ms/dto"
	"github.com/device-ms/errors"
	"github.com/device-ms/model"
	"github.com/device-ms/util"
)

type searchDevicesRequest struct {
	dto.SearchDevicesRequestDTO
}

// Build builds the creation dto
func (req *searchDevicesRequest) Build(r *http.Request) error {
	req.Brand = model.Brand(r.URL.Query().Get("brand"))

	return req.Validate()
}

// Validate validates the creation dto
func (req searchDevicesRequest) Validate() error {
	if req.Brand != "" {
		if !req.Brand.IsValid() {
			return errors.InvalidParameterError("brand", "invalid value ["+string(req.Brand)+"]")
		}
	}
	return nil
}

func (dh deviceHandler) getDevices(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := new(searchDevicesRequest)
	if err := req.Build(r); err != nil {
		util.JSONErrorWithCtx(ctx, w, err, http.StatusBadRequest)
		return
	}

	var res []dto.DeviceDTO
	var err error
	if req.Brand != "" {
		res, err = dh.service.DeviceController().GetDevicesByBrand(ctx, req.Brand)
		if err != nil {
			util.JSONErrorWithCtx(ctx, w, err, http.StatusInternalServerError)
			return
		}
	} else {
		res, err = dh.service.DeviceController().GetDevices(ctx)
		if err != nil {
			util.JSONErrorWithCtx(ctx, w, err, http.StatusInternalServerError)
			return
		}
	}

	util.JSONReturnWithCtx(ctx, w, http.StatusOK, res)
}
