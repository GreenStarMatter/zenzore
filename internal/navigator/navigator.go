package navigator

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Level int
type ID string

//Levels: Root, System, Device, Sensor

const (
	Root = iota
	Zyztem
	Device
	Sensor
)

func (l Level) String() string {
	switch l {
	case Root:
		return "Root"
	case Zyztem:
		return "Zyztem"
	case Device:
		return "Device"
	case Sensor:
		return "Sensor"
	default:
		return "Unknown"
	}
}

type deviceSummary struct {
	SN string `json:"DeviceSN"`
	PN string `json:"DevicePN"`
}

type sensorSummary struct {
	SN string `json:"SensorSN"`
	PN string `json:"SensorPN"`
}

type zyztemSummary struct {
	ID string `json:"ID"`
}

type NavigatorNode struct {
	ID       ID //need to determine ID
	Level    Level
	Parent   *NavigatorNode
	Children []*NavigatorNode
}

type Navigator struct {
	CurrentNode *NavigatorNode
}

func MakeID(pn, sn string) ID {
	return ID(pn + ":" + sn)
}

func SplitID(id ID) (pn, sn string, err error) {
	parts := strings.SplitN(string(id), ":", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid id %q", id)
	}
	return parts[0], parts[1], nil
}

func New() (*Navigator, *NavigatorNode) {
	rootNode := &NavigatorNode{ID: "0", Level: Root, Parent: nil, Children: make([]*NavigatorNode, 0)}
	navigator := &Navigator{CurrentNode: rootNode}
	return navigator, rootNode

}

func (parent *NavigatorNode) Add() (*NavigatorNode, error) {
	if parent.Level == Sensor {
		return nil, fmt.Errorf("canot add child node to sensor")
	}
	childLevel := parent.Level + 1
	NewNode := &NavigatorNode{ID: "0", Level: childLevel, Parent: parent, Children: make([]*NavigatorNode, 0)}
	parent.Children = append(parent.Children, NewNode)
	return NewNode, nil
}

func (parent *NavigatorNode) Remove(child *NavigatorNode) error {
	//remove child
	//remove child from parent list
	//find match index in parent.Children and pop

	for i, c := range parent.Children {
		if c == child {
			parent.Children = append(parent.Children[:i], parent.Children[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("child node %q not found under parent %q", child.ID, parent.ID)
}

func (parent *NavigatorNode) List() string {
	if len(parent.Children) == 0 {
		return "(no children)"
	}
	var sb strings.Builder
	for i, c := range parent.Children {
		fmt.Fprintf(&sb, "%d: %s\n", i, c.ID)
	}
	return sb.String()

}

func (nav *Navigator) Set(chosenNode *NavigatorNode) {
	nav.CurrentNode = chosenNode
}

func (nav *Navigator) Up() error {
	if nav.CurrentNode.Level == Root {
		return fmt.Errorf("already at root, nowhere up to go")
	}
	nav.CurrentNode = nav.CurrentNode.Parent
	return nil
}

func (nav *Navigator) Down(child *NavigatorNode) error {
	if nav.CurrentNode.Level == Sensor {
		return fmt.Errorf("at lowest level cannot go down further")
	}
	for _, c := range nav.CurrentNode.Children {
		if c == child {
			nav.CurrentNode = child
			return nil
		}
	}
	return fmt.Errorf("node %q is not a child of current node %q", child.ID, nav.CurrentNode.ID)

}

// In navigator.go

func (node *NavigatorNode) Populate(baseURL string) error {
	switch node.Level {
	case Root:
		return node.PopulateZyztems(baseURL)
	case Zyztem:
		return node.PopulateDevices(baseURL)
	case Device:
		return node.PopulateSensors(baseURL)
	}
	return nil
}

func (node *NavigatorNode) PopulateZyztems(baseURL string) error {
	if node.Level != Root {
		return fmt.Errorf("node is not root")
	}

	resp, err := http.Get(baseURL + "/zyztems")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var zyztems []zyztemSummary
	if err := json.NewDecoder(resp.Body).Decode(&zyztems); err != nil {
		return err
	}

	node.Children = nil

	for _, z := range zyztems {
		child, err := node.Add()
		if err != nil {
			return err
		}
		child.ID = ID(z.ID)
	}

	return nil
}

func (node *NavigatorNode) PopulateDevices(baseURL string) error {
	if node.Level != Zyztem {
		return fmt.Errorf("node is not a zyztem")
	}

	resp, err := http.Get(fmt.Sprintf("%s/zyztems/%s/devices", baseURL, node.ID))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var devices []deviceSummary
	if err := json.NewDecoder(resp.Body).Decode(&devices); err != nil {
		return err
	}

	// Refresh the children in case the server state changed.
	node.Children = nil

	for _, d := range devices {
		child, err := node.Add()
		if err != nil {
			return err
		}

		child.ID = MakeID(d.PN, d.SN)
	}

	return nil
}

func (node *NavigatorNode) PopulateSensors(baseURL string) error {
	if node.Level != Device {
		return fmt.Errorf("node is not a device")
	}

	devicePN, deviceSN, err := SplitID(node.ID)
	if err != nil {
		return err
	}

	zyztemID := node.Parent.ID

	url := fmt.Sprintf(
		"%s/zyztems/%s/devices/%s/%s/sensors",
		baseURL,
		zyztemID,
		devicePN,
		deviceSN,
	)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var sensors []sensorSummary
	if err := json.NewDecoder(resp.Body).Decode(&sensors); err != nil {
		return err
	}

	node.Children = nil

	for _, s := range sensors {
		child, err := node.Add()
		if err != nil {
			return err
		}

		child.ID = MakeID(s.PN, s.SN)
	}

	return nil
}

func LoadFromServer(baseURL string) (*Navigator, *NavigatorNode, error) {
	navigator, rootNode := New()
	if err := rootNode.PopulateZyztems(baseURL); err != nil {
		return nil, nil, fmt.Errorf("populating root zyztems: %w", err)
	}
	return navigator, rootNode, nil
}
