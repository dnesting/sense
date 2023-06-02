package sense

import (
	"context"

	"github.com/dnesting/sense/internal/client"
)

type Device struct {
	ID       string
	Name     string
	Type     string
	Make     string
	Model    string
	Location string
}

// I'm not entirely sure what the relationship between these fields is, so
// just pick one that seems reasonable.
func getType(d client.Device) string {
	tags := deref(d.Tags)

	for _, tag := range []string{
		"UserDeviceType",
		"Type",
		"DefaultUserDeviceType",
	} {
		if s := stringOrEmpty(tags[tag]); s != "" {
			return s
		}
	}
	return ""
}

func stringOrEmpty(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

// GetDevices returns a list of devices known to the given monitor.
func (s *Client) GetDevices(ctx context.Context, monitorID int, includeMerged bool) (devs []Device, err error) {
	res, err1 := s.client.GetDevicesWithResponse(
		ctx,
		monitorID,
		&client.GetDevicesParams{
			IncludeMerged: &includeMerged,
		})
	if err := client.Ensure(err1, "GetDevices", res, 200); err != nil {
		return nil, err
	}
	for _, d := range deref(res.JSON200.Devices) {
		devs = append(devs, Device{
			ID:       deref(d.Id),
			Name:     deref(d.Name),
			Type:     getType(d),
			Make:     deref(d.Make),
			Model:    deref(d.Model),
			Location: deref(d.Location),
		})
	}
	return devs, nil
}
