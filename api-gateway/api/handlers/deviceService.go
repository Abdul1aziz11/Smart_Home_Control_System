package handlers

import (
	"api-gateway/genproto"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Create a new device
// @Description Create a new device with the given information
// @Security BearerAuth
// @Tags devices
// @Accept json
// @Produce json
// @Param device body DeviceRegister true "Device Info"
// @Param   X-User-Role header string true "User Role"
// @Success 200 {object} Device
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /devices [post]
func (h *Handlers) CreateDevice(ctx *gin.Context) {
	var device genproto.Device
	if err := ctx.BindJSON(&device); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.DeviceControl.CreateDevice(ctx, &genproto.CreateDeviceRequest{
		Device: &device,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// @Summary Update an existing device
// @Description Update a device with the given ID
// @Security BearerAuth
// @Security BearerAuth
// @Tags devices
// @Accept json
// @Produce json
// @Param id path string true "Device ID"
// @Param device body DeviceRegister true "Device Info"
// @Param   X-User-Role header string true "User Role"
// @Success 200 {object} Device
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /devices/{id} [put]
func (h *Handlers) UpdateDevice(ctx *gin.Context) {
	var device genproto.Device
	deviceID := ctx.Param("id")
	if err := ctx.BindJSON(&device); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	device.DeviceId = deviceID

	resp, err := h.DeviceControl.UpdateDevice(ctx, &genproto.UpdateDeviceRequest{
		Device: &device,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// @Summary Delete a device
// @Description Delete a device with the given ID
// @Security BearerAuth
// @Tags devices
// @Produce json
// @Param id path string true "Device ID"
// @Param   X-User-Role header string true "User Role"
// @Success 200 {object} DeleteResp
// @Failure 500 {object} map[string]string
// @Router /devices/{id} [delete]
func (h *Handlers) DeleteDevice(ctx *gin.Context) {
	deviceID := ctx.Param("id")

	resp, err := h.DeviceControl.DeleteDevice(ctx, &genproto.DeleteDeviceRequest{
		DeviceId: deviceID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// @Summary Control a device
// @Description Send command to control a device
// @Security BearerAuth
// @Tags devices
// @Accept json
// @Produce json
// @Param device body Command true "Command Info"
// @Param   X-User-Role header string true "User Role"
// @Success 200 {object} RespStatus
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /control [post]
func (h *Handlers) ControlDevice(ctx *gin.Context) {
	var command genproto.Command
	if err := ctx.BindJSON(&command); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.DeviceControl.SendCommand(ctx, &command)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// @Summary Get a device
// @Description Get details of a device with the given ID
// @Security BearerAuth
// @Tags devices
// @Produce json
// @Param id path string true "Device ID"
// @Param   X-User-Role header string true "User Role"
// @Success 200 {object} Device
// @Failure 500 {object} map[string]string
// @Router /devices/{id} [get]
func (h *Handlers) GetDevice(ctx *gin.Context) {
	deviceID := ctx.Param("id")

	resp, err := h.DeviceControl.GetDevice(ctx, &genproto.GetDeviceRequest{
		DeviceId: deviceID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
