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

func LoadFromServer(baseURL string) (*Navigator, *NavigatorNode, error) {
	resp, err := http.Get(baseURL + "/zyztems/list")
	if err != nil {
		return nil, nil, fmt.Errorf("requesting zyztem list: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("unexpected status from server: %d", resp.StatusCode)
	}

	var zyztems []zyztemSummary
	if err := json.NewDecoder(resp.Body).Decode(&zyztems); err != nil {
		return nil, nil, fmt.Errorf("decoding zyztem list: %w", err)
	}

	navigator, rootNode := New()

	for _, z := range zyztems {
		childNode, err := rootNode.Add()
		if err != nil {
			return nil, nil, fmt.Errorf("adding node for zyztem %q: %w", z.ID, err)
		}
		childNode.ID = ID(z.ID)
	}

	return navigator, rootNode, nil
}
