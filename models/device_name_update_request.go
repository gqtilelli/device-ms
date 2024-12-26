// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// DeviceNameUpdateRequest DeviceNameUpdateRequest
//
// swagger:model DeviceNameUpdateRequest
type DeviceNameUpdateRequest struct {

	// The name of the device
	Name string `json:"name,omitempty"`
}

// Validate validates this device name update request
func (m *DeviceNameUpdateRequest) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this device name update request based on context it is used
func (m *DeviceNameUpdateRequest) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *DeviceNameUpdateRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DeviceNameUpdateRequest) UnmarshalBinary(b []byte) error {
	var res DeviceNameUpdateRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
